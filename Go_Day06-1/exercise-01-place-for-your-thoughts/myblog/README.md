Давайте разберем каждую функцию в вашем коде и объясним, что она делает.

### 1. `add(x, y int) int` и `sub(x, y int) int`

```go
func add(x, y int) int {
    return x + y
}

func sub(x, y int) int {
    return x - y
}
```

Эти функции предназначены для использования в шаблонах HTML (через `template.FuncMap`). Они возвращают сумму и разность двух целых чисел соответственно. Например, их можно использовать в шаблонах для вычисления значений прямо в HTML-коде.

### 2. `adminHandler(w http.ResponseWriter, r *http.Request)`

```go
func adminHandler(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("templates/admin.html"))
    tmpl.Execute(w, nil)
}
```

Эта функция обрабатывает запросы к админской странице (`/admin`). Она загружает и отображает HTML-шаблон `admin.html` из папки `templates`.

- `template.Must` проверяет, не произошла ли ошибка при загрузке шаблона, и если ошибка есть, программа завершится с паникой.
- `tmpl.Execute(w, nil)` применяет шаблон и отправляет результат клиенту.

### 3. `postArticleHandler(w http.ResponseWriter, r *http.Request)`

```go
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
```

Эта функция обрабатывает POST-запросы на добавление статьи через админскую панель. Вот как она работает:

- Извлекает данные из формы (`title` и `content`).
- Вставляет новую статью в базу данных с использованием SQL-запроса.
- Если вставка прошла успешно, перенаправляет пользователя на главную страницу (`/`).
- Если произошла ошибка, возвращает клиенту сообщение об ошибке и код состояния 500.

### 4. `articleHandler(w http.ResponseWriter, r *http.Request)`

```go
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
```

Эта функция обрабатывает запросы на просмотр статьи по ID (`/article/{id}`). Вот ее действия:

- Извлекает ID статьи из URL.
- Выполняет SQL-запрос для получения данных о статье (название и содержимое).
- Если статья найдена, она рендерится в HTML с использованием встроенного шаблона.
- Использует `blackfriday` для преобразования содержимого статьи из Markdown в HTML.
- Возвращает рендеренную страницу клиенту.

### 5. `homeHandler(w http.ResponseWriter, r *http.Request)`

```go
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
```

Эта функция обрабатывает запросы на главную страницу (`/`). Она отвечает за отображение списка статей с поддержкой пагинации:

- Извлекает номер страницы из URL (по умолчанию страница 1).
- Выполняет SQL-запрос для получения списка статей с ограничением по количеству (`pageSize = 3`) и с учетом смещения (`offset`).
- Собирает данные о статьях и вычисляет общее количество страниц.
- Рендерит шаблон `home.html`, передавая в него список статей, текущую страницу и общее количество страниц.
- Шаблон использует функции `add` и `sub` для вычисления значений в HTML.

### 6. `countTotalArticles() int`

```go
func countTotalArticles() int {
    var count int
    err := db.QueryRow("SELECT COUNT(*) FROM articles").Scan(&count)
    if err != nil {
        log.Println("Error counting articles:", err)
        return 0
    }
    return count
}
```

Эта функция возвращает общее количество статей в базе данных. Она используется для расчета количества страниц в `homeHandler`.

### 7. `main()`

```go
func main() {
    var err error
    db, err = sql.Open("postgres", "user=postgres password=postgres dbname=postgres sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }

    r := mux.NewRouter()
    r.HandleFunc("/", homeHandler).Methods("GET")
    r.HandleFunc("/article/{id:[0-9]+}", articleHandler).Methods("GET")
    r.HandleFunc("/admin", adminHandler).Methods("GET")
    r.HandleFunc("/admin/post", postArticleHandler).Methods("POST")

    r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))))
    r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("./images/"))))
    r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("./js/"))))

    log.Println("Server starting on port 8888...")
    http.ListenAndServe(":8888", nil)
}
```

Эта функция инициализирует приложение:

- Подключается к базе данных PostgreSQL.
- Настраивает маршруты (URL-адреса) с помощью `mux.Router`.
- Определяет обработчики для главной страницы, страницы статьи, админской панели и обработки новых статей.
- Настраивает обработку статических файлов (CSS, изображения, JS).
- Запускает HTTP-сервер на порту 8888.

### Общая работа приложения

Приложение представляет собой простой веб-сайт, на котором можно:

- Просматривать статьи на главной странице.
- Переходить к полному содержимому статьи по ее ID.
- Создавать новые статьи через админскую панель.
- Пагинация реализована на главной странице, чтобы показывать ограниченное количество статей на одной странице.