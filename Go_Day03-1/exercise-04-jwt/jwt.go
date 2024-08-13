package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"errors"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type Page struct {
	Title     string
	Total     string
	Items     []string
	Name      []string
	Address   []string
	Phone     []string
	PageNum   int
	TotalPage int
	PrevPage  int
	NextPage  int
	LastPage  int
}

type PageDataApi struct {
	Title     string
	Total     string
	Items     []Source
	PageNum   int
	TotalPage int
	PrevPage  int
	NextPage  int
	LastPage  int
}

type PageRecommendation struct {
	Title string
	Items []Source
}

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Source struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Address  string   `json:"address"`
	Phone    string   `json:"phone"`
	Location Location `json:"location"`
}

type Hits struct {
	Total struct {
		Value    int    `json:"value"`
		Relation string `json:"relation"`
	} `json:"total"`
	MaxScore interface{} `json:"max_score"`
	Hits     []struct {
		Index  string      `json:"_index"`
		Type   string      `json:"_type"`
		ID     string      `json:"_id"`
		Score  interface{} `json:"_score"`
		Source Source      `json:"_source"`
		Sort   interface{} `json:"sort"`
	} `json:"hits"`
}

type Response struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits Hits `json:"hits"`
}

type UnsignedResponse struct {
	Message interface{} `json:"message"`
}

type SignedResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func readFromFile(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Ошибка чтения файла: %v", err)
	}
	return data, nil
}

func sendPutRequestMaxResult(client *http.Client) error {
	url := "http://localhost:9200/places/_settings"
	requestBody := []byte(`{"index":{"max_result_window":20000}}`)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Неуспешный запрос, код состояния: %d", res.StatusCode)
	}

	return nil
}

func extractDataFromSource(places []Source, totalRecord int, offset int, chooseField string) []string {
	var newSlice []string
	for n := offset; n < len(places) && n < offset+10; n++ {
		if chooseField == "Address" {
			newSlice = append(newSlice, places[n].Address)
		} else if chooseField == "Phone" {
			newSlice = append(newSlice, places[n].Phone)
		} else if chooseField == "Name" {
			newSlice = append(newSlice, places[n].Name)
		}
	}
	return newSlice
}

func handlerCommon(c *gin.Context, indexHtml string, places []Source, totalRecord int) {
	tmpl := template.Must(template.New("indexHtml").Parse(indexHtml))
	pageSize := 10
	// pageNum := c.Param("pageNum")
	pageNum := c.DefaultQuery("page", "")

	if pageNum == "" {
		c.String(http.StatusBadRequest, "Недопустимый номер страницы")
		return
	}

	pageNumInt := 1
	fmt.Sscanf(pageNum, "%d", &pageNumInt)

	if pageNumInt < 1 || pageNumInt > (totalRecord+pageSize-1)/pageSize {
		c.String(http.StatusBadRequest, "Invalid 'page' value: 'foo'")
		return
	}

	offset := (pageNumInt - 1) * pageSize

	pageName := extractDataFromSource(places, totalRecord, offset, "Name")
	pageAddress := extractDataFromSource(places, totalRecord, offset, "Address")
	pagePhone := extractDataFromSource(places, totalRecord, offset, "Phone")

	pageData := Page{
		Title:     "Places",
		Total:     "Total: " + strconv.Itoa(totalRecord),
		Items:     pageName,
		Name:      pageName,
		Address:   pageAddress,
		Phone:     pagePhone,
		PageNum:   pageNumInt,
		TotalPage: (totalRecord + pageSize - 1) / pageSize,
		PrevPage:  max(pageNumInt-1, 1),
		NextPage:  min(pageNumInt+1, (totalRecord+pageSize-1)/pageSize),
		LastPage:  (totalRecord + pageSize - 1) / pageSize,
	}

	tmpl.Execute(c.Writer, pageData)
}

func handlerLogin(c *gin.Context) {
	indexHtmlBytes, err2 := ioutil.ReadFile("authorization.html")
	if err2 != nil {
		fmt.Println("Ошибка при чтении файла:", err2)
		return
	}

	indexHtml := string(indexHtmlBytes)

	tmpl := template.Must(template.New("indexHtml").Parse(indexHtml))
	tmpl.Execute(c.Writer, "")
	login := c.PostForm("login")
	password := c.PostForm("password")

	// Делайте что-то с полученными значениями логина и пароля

	// Пример вывода значений для проверки
	fmt.Println("!Логин:", login)
	fmt.Println("!Пароль:", password)

	// if login == "valid_username" && password == "valid_password" {
	// 	// Redirect to another page upon successful login
	// 	c.Redirect(http.StatusFound, "/success") // Replace "/success" with the actual URL you want to redirect to.
	// 	return
	// }

	// // Отправка ответа клиенту
	// if login != "" && password != "" {
	// 	// c.String(http.StatusOK, "Авторизация успешна")
	// 	c.Redirect(http.StatusFound, "/api/recommend")
	// }
	// // c.Redirect(http.StatusFound, "/api/recommend")
}

func handlerLoginPOST(c *gin.Context) {
	login := c.PostForm("login")
	password := c.PostForm("password")
	var statusauthorization string
	var redirPath string
	if login == "user" && password == "123" {
		statusauthorization = "Авторизация успешна"
		redirPath = "/api/get_token"
	} else {
		statusauthorization = "Авторизация неуспешна"
		redirPath = "/login"
	}

	htmlResponse := fmt.Sprintf(`
	<html>
	<body>
		<p>%s</p>
		<script>
			setTimeout(function () {
				window.location.href = "%s";
			}, 2000); 
		</script>
	</body>
	</html>
	`, statusauthorization, redirPath)
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Writer.WriteString(htmlResponse)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func handlerApi(c *gin.Context, places []Source, totalRecord int) {
	pageSize := 10
	pageNum := c.DefaultQuery("page", "")

	if pageNum == "" {
		c.String(http.StatusBadRequest, "Недопустимый номер страницы")
		return
	}

	pageNumInt := 1
	fmt.Sscanf(pageNum, "%d", &pageNumInt)

	if pageNumInt < 1 || pageNumInt > (totalRecord+pageSize-1)/pageSize {
		c.String(http.StatusBadRequest, "Invalid 'page' value: 'foo'")
		return
	}

	offset := (pageNumInt - 1) * pageSize
	limit := pageSize

	pageItems := places[offset:min(offset+limit, totalRecord)]

	pageData := PageDataApi{
		Title:     "Places",
		Total:     strconv.Itoa(totalRecord),
		Items:     pageItems,
		PageNum:   pageNumInt,
		TotalPage: (totalRecord + pageSize - 1) / pageSize,
		PrevPage:  max(pageNumInt-1, 1),
		NextPage:  min(pageNumInt+1, (totalRecord+pageSize-1)/pageSize),
		LastPage:  (totalRecord + pageSize - 1) / pageSize,
	}

	jsonData, err := json.MarshalIndent(pageData, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.String(http.StatusOK, string(jsonData))
}

func handlerRecommend(c *gin.Context, h *http.Client) {
	lat := c.DefaultQuery("lat", "0.0")
	lon := c.DefaultQuery("lon", "0.0")

	if lat == "" || lon == "" {
		c.String(http.StatusBadRequest, "Invalid 'lat' or 'lon' value: 'foo'")
		return
	}

	url := "http://localhost:9200/places/_search"
	payload := []byte(fmt.Sprintf(`{
		"query": {
			"match_all": {}
		},
		"sort": [
			{
				"_geo_distance": {
					"location": {
						"lat": %s,
						"lon": %s
					},
					"order": "asc",
					"unit": "km",
					"mode": "min",
					"distance_type": "arc",
					"ignore_unmapped": true
				}
			}
		],
		"size": 3
	}`, lat, lon))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := h.Do(req)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return
	}

	var response Response
	err2 := json.Unmarshal([]byte(responseBody), &response)
	handleError(err2)

	places := extractDataFromResponse(response, 3) //  source

	pageData := PageRecommendation{
		Title: "Places",
		Items: places,
	}

	jsonData, err := json.MarshalIndent(pageData, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.String(http.StatusOK, string(jsonData))
}

func extractDataFromResponse(response Response, totalRecord int) []Source {
	places := []Source{}
	for _, hit := range response.Hits.Hits {
		place := Source{
			ID:       hit.Source.ID,
			Name:     hit.Source.Name,
			Address:  hit.Source.Address,
			Phone:    hit.Source.Phone,
			Location: Location{Lat: hit.Source.Location.Lat, Lon: hit.Source.Location.Lon},
		}
		places = append(places, place)
	}
	return places
}

func login(c *gin.Context) {
	// type login struct {
	// 	Username string `json:"username,omitempty"`
	// }

	// loginParams := login{}
	// c.ShouldBindJSON(&loginParams)

	// loginStr := c.DefaultQuery("login", "")

	// if loginParams.Username == "mike" || loginParams.Username == "rama" {
	// fmt.Println(loginStr)
	// if loginStr == "mike" {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": "correctlogin",
		"nbf":  time.Date(2018, 01, 01, 12, 0, 0, 0, time.UTC).Unix(),
	})

	tokenStr, err := token.SignedString([]byte("supersaucysecret"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, UnsignedResponse{
			Message: err.Error(),
		})
		return
	}

	// Сохранение токена в куки
	c.SetCookie("jwt_token", tokenStr, 3600, "/", "127.0.0.1", false, true)
	c.Redirect(http.StatusFound, "/api/recommend")
	// c.JSON(http.StatusOK, SignedResponse{
	// 	Token:   tokenStr,
	// 	Message: "logged in",
	// })

	return

	// }

	// c.JSON(http.StatusBadRequest, UnsignedResponse{
	// 	Message: "bad username",
	// })
}

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("bad header value given")
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("incorrectly formatted authorizationorization header")
	}

	return jwtToken[1], nil
}

func parseToken(jwtToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, OK := token.Method.(*jwt.SigningMethodHMAC); !OK {
			return nil, errors.New("bad signed method received")
		}
		return []byte("supersaucysecret"), nil
	})

	if err != nil {
		return nil, errors.New("bad jwt token")
	}

	return token, nil
}

func jwtTokenCheck(c *gin.Context) {
	// Извлечение JWT-токена из куки
	cookie, err := c.Cookie("jwt_token")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
			Message: "JWT token not found in cookie",
		})
		return
	}

	// Парсинг JWT-токена
	token, err := parseToken(cookie)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
			Message: "bad JWT token in cookie",
		})
		return
	}

	_, OK := token.Claims.(jwt.MapClaims)
	if !OK {
		c.AbortWithStatusJSON(http.StatusInternalServerError, UnsignedResponse{
			Message: "unable to parse claims from JWT token",
		})
		return
	}

	c.Next()
}

func main() {
	httpClient := &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       5 * time.Second,
	}

	if err := sendPutRequestMaxResult(httpClient); err != nil {
		log.Fatal(err)
	}

	res, err := httpClient.Get("http://localhost:9200/places/_search?size=20000&sort=id:asc")
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var response Response
	err2 := json.Unmarshal([]byte(body), &response)
	handleError(err2)

	totalRecord := response.Hits.Total.Value

	places := extractDataFromResponse(response, totalRecord)

	indexHtmlBytes, err := ioutil.ReadFile("index.html")
	if err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
		return
	}

	indexHtml := string(indexHtmlBytes)

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		handlerCommon(c, indexHtml, places, totalRecord)
	})

	router.GET("/api/places", func(c *gin.Context) {
		handlerApi(c, places, totalRecord)
	})

	// router.GET("/api/recommend", func(c *gin.Context) {
	// 	handlerRecommend(c, httpClient)
	// })

	router.GET("/login", func(c *gin.Context) {
		handlerLogin(c)
	})

	router.POST("/login", func(c *gin.Context) {
		handlerLoginPOST(c)
	})

	// router.GET("/login", login)
	router.GET("/api/get_token", login)
	router.Use(jwtTokenCheck).GET("/api/recommend", func(c *gin.Context) {
		handlerRecommend(c, httpClient)
	})

	router.Run(":8888")
}

// http://127.0.0.1:8888/?page=3

// http://127.0.0.1:8888/api/places?page=3

// http://127.0.0.1:8888/api/recommend?lat=55.674&lon=37.666

// http://127.0.0.1:8888/login

// http://127.0.0.1:8888/api/get_token

// http://127.0.0.1:8888/login?login=mike

// http://127.0.0.1:8888/api/get_token?login=mike

// http://127.0.0.1:8888/api/recommend?lat=55.674&lon=37.666
