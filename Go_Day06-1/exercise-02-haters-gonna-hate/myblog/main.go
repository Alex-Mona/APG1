package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"

	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
	"github.com/russross/blackfriday/v2"
	"golang.org/x/time/rate"
)

var (
	db *sql.DB

	// Создаем глобальный rate limiter на все запросы
	limiter = rate.NewLimiter(1, 3) // 1 запрос в секунду, с максимум 3 в бэклоге

	// Мьютекс для безопасного использования rate limiter в многопоточной среде
	mu sync.Mutex
)

// Определяем функции для шаблонов
func add(x, y int) int {
	return x + y
}

func sub(x, y int) int {
	return x - y
}

// Middleware для ограничения частоты запросов
func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		allowed := limiter.Allow()
		mu.Unlock()

		if !allowed {
			http.Error(w, "429 Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/admin.html"))
	tmpl.Execute(w, nil)
}

func postArticleHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	content := r.FormValue("content")

	_, err := db.Exec("INSERT INTO articles (title, content) VALUES ($1, $2)", title, content)
	if err != nil {
		http.Error(w, "Unable to post article", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func articleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var title, content string
	err = db.QueryRow("SELECT title, content FROM articles WHERE id = $1", id).Scan(&title, &content)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	htmlContent := blackfriday.Run([]byte(content))
	tmpl := template.Must(template.New("article").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
		    <title>{{.Title}}</title>
		    <link rel="stylesheet" type="text/css" href="/css/main.css">
		</head>
		<body>
		    <h1>{{.Title}}</h1>
		    <div>{{.Content}}</div>
		    <a href="/">Back</a>
		</body>
		</html>
	`))

	data := struct {
		Title   string
		Content template.HTML
	}{
		Title:   title,
		Content: template.HTML(htmlContent),
	}

	tmpl.Execute(w, data)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}
	pageNumber, err := strconv.Atoi(page)
	if err != nil || pageNumber < 1 {
		pageNumber = 1
	}

	const pageSize = 3
	offset := (pageNumber - 1) * pageSize

	rows, err := db.Query("SELECT id, title FROM articles ORDER BY created_at DESC LIMIT $1 OFFSET $2", pageSize, offset)
	if err != nil {
		log.Printf("Error fetching articles: %v", err)
		http.Error(w, "Unable to fetch articles", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var articles []struct {
		ID    int
		Title string
	}
	for rows.Next() {
		var a struct {
			ID    int
			Title string
		}
		if err := rows.Scan(&a.ID, &a.Title); err != nil {
			log.Printf("Error scanning article: %v", err)
			http.Error(w, "Unable to process articles", http.StatusInternalServerError)
			return
		}
		articles = append(articles, a)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		http.Error(w, "Unable to process articles", http.StatusInternalServerError)
		return
	}

	totalRows := countTotalArticles()
	totalPages := (totalRows + pageSize - 1) / pageSize

	// Передаем функции в шаблон
	funcMap := template.FuncMap{
		"add": add,
		"sub": sub,
	}
	tmpl := template.Must(template.New("home.html").Funcs(funcMap).ParseFiles("templates/home.html"))
	data := struct {
		Articles  []struct{ ID int; Title string }
		Page      int
		TotalPages int
	}{
		Articles:  articles,
		Page:      pageNumber,
		TotalPages: totalPages,
	}
	if err := tmpl.ExecuteTemplate(w, "home.html", data); err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Unable to render page", http.StatusInternalServerError)
	}
}

func countTotalArticles() int {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM articles").Scan(&count)
	if err != nil {
		log.Println("Error counting articles:", err)
		return 0
	}
	return count
}

func main() {
	var err error
	// Подключение к базе данных
	db, err = sql.Open("postgres", "user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	// Настройка маршрутизации
	r := mux.NewRouter()

	// Добавление rate limiter middleware ко всем маршрутам
	r.Use(rateLimitMiddleware)

	// Обработчики маршрутов
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/article/{id:[0-9]+}", articleHandler).Methods("GET")
	r.HandleFunc("/admin", adminHandler).Methods("GET")
	r.HandleFunc("/admin/post", postArticleHandler).Methods("POST")

	// Обработка статических файлов
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("./images/"))))
	r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("./js/"))))

	// Запуск сервера
	log.Println("Server starting on port 8888...")
	if err := http.ListenAndServe(":8888", r); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
