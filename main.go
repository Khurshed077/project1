package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

type Article struct {
	Id, Category_id         uint16
	Title, Anons, Full_text string
	Image                   sql.NullString
}

type Category struct {
	ID   int
	Name string
}

var (
	db    *sql.DB
	store = sessions.NewCookieStore([]byte("super-secret-key"))
	tpl   *template.Template
)

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

// Создание таблицы

func initDB() {
	var err error
	db, err = sql.Open("sqlite", "/data/database.db")
	if err != nil {
		log.Fatal("Ошибка открытия SQLite: ", err)
	}

	// Создаём таблицы, если их нет
	schema := `
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS articles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    anons TEXT NOT NULL,
    full_text TEXT NOT NULL,
    image TEXT,
    category_id INTEGER,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
);
`
	_, err = db.Exec(schema)
	if err != nil {
		log.Fatal("Ошибка создания таблиц: ", err)
	}

	// если нет котегорий добавим их
	db.Exec(`INSERT INTO categories (name) SELECT 'Новости' WHERE NOT EXISTS (SELECT 1 FROM categories WHERE name='Новости');`)
	db.Exec(`INSERT INTO categories (name) SELECT 'Статьи' WHERE NOT EXISTS (SELECT 1 FROM categories WHERE name='Статьи');`)
	db.Exec(`INSERT INTO categories (name) SELECT 'Игры' WHERE NOT EXISTS (SELECT 1 FROM categories WHERE name='Игры');`)

	log.Println("SQLite (modernc) инициализирован")
}

//Пользователи

func getUser(r *http.Request) string {
	sess, _ := store.Get(r, "session")
	if username, ok := sess.Values["username"].(string); ok {
		return username
	}
	return ""
}

// Главная страница

func index(w http.ResponseWriter, r *http.Request) {
	sortCol := r.URL.Query().Get("sort")

	var rows *sql.Rows
	var err error

	if sortCol == "" {
		rows, err = db.Query(`SELECT id, title, anons, full_text, image, category_id FROM articles`)
	} else {
		rows, err = db.Query(
			`SELECT id, title, anons, full_text, image, category_id FROM articles WHERE category_id = ?`,
			sortCol,
		)
	}
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var posts []Article
	for rows.Next() {
		var p Article
		err = rows.Scan(&p.Id, &p.Title, &p.Anons, &p.Full_text, &p.Image, &p.Category_id)
		if err != nil {
			panic(err)
		}
		if !p.Image.Valid {
			p.Image.String = ""
		}
		posts = append(posts, p)
	}

	// категории
	crows, err := db.Query("SELECT id, name FROM categories")
	if err != nil {
		panic(err)
	}
	defer crows.Close()

	var categories []Category
	for crows.Next() {
		var c Category
		crows.Scan(&c.ID, &c.Name)
		categories = append(categories, c)
	}

	tpl.ExecuteTemplate(w, "index", map[string]interface{}{
		"User":       getUser(r),
		"Posts":      posts,
		"Categories": categories,
	})
}

func show_post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var post Article
	err := db.QueryRow(
		`SELECT id, title, anons, full_text, image, category_id FROM articles WHERE id = ?`,
		vars["id"],
	).Scan(&post.Id, &post.Title, &post.Anons, &post.Full_text, &post.Image, &post.Category_id)

	if err != nil {
		http.Error(w, "Статья не найдена", http.StatusNotFound)
		return
	}

	if !post.Image.Valid {
		post.Image.String = ""
	}

	tpl.ExecuteTemplate(w, "show", map[string]interface{}{
		"User": getUser(r),
		"Post": post,
	})
}

func create(w http.ResponseWriter, r *http.Request) {
	if getUser(r) == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	rows, err := db.Query("SELECT id, name FROM categories")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var c Category
		rows.Scan(&c.ID, &c.Name)
		categories = append(categories, c)
	}

	tpl.ExecuteTemplate(w, "create", map[string]interface{}{
		"User":       getUser(r),
		"Categories": categories,
	})
}

func save_article(w http.ResponseWriter, r *http.Request) {
	if getUser(r) == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	r.ParseMultipartForm(10 << 20)

	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")
	category_id := r.FormValue("category_id")

	imagePath := ""
	file, handler, err := r.FormFile("image")
	if err == nil {
		defer file.Close()

		uploadDir := "/data/uploads"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.Mkdir(uploadDir, os.ModePerm)
		}

		dst, err := os.Create(uploadDir + "/" + handler.Filename)
		if err != nil {
			panic(err)
		}
		defer dst.Close()
		io.Copy(dst, file)
		imagePath = "uploads/" + handler.Filename
	}

	_, err = db.Exec(
		`INSERT INTO articles (title, anons, full_text, image, category_id)
		 VALUES (?, ?, ?, ?, ?)`,
		title, anons, full_text, imagePath, category_id,
	)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func register(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "register", nil)
}

func registerPost(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if username == "" || email == "" || password == "" {
		fmt.Fprintf(w, "Все поля обязательны!")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Ошибка хеширования", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(
		"INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)",
		username, email, string(hash),
	)

	if err != nil {
		fmt.Fprintf(w, "Ошибка: %v", err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func login(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "login", nil)
}

func loginPost(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	var id int
	var username, passwordHash string

	err := db.QueryRow(
		`SELECT id, username, password_hash FROM users WHERE email = ?`,
		email,
	).Scan(&id, &username, &passwordHash)

	if err != nil {
		fmt.Fprintf(w, "Пользователь не найден")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		fmt.Fprintf(w, "Неверный пароль")
		return
	}

	sess, _ := store.Get(r, "session")
	sess.Values["user_id"] = id
	sess.Values["username"] = username
	sess.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func logout(w http.ResponseWriter, r *http.Request) {
	sess, _ := store.Get(r, "session")
	sess.Options.MaxAge = -1
	sess.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Пути

func handleFunc() {
	rtr := mux.NewRouter()

	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/create", create).Methods("GET")
	rtr.HandleFunc("/save_article", save_article).Methods("POST")
	rtr.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET")
	rtr.HandleFunc("/register", register).Methods("GET")
	rtr.HandleFunc("/register", registerPost).Methods("POST")
	rtr.HandleFunc("/login", login).Methods("GET")
	rtr.HandleFunc("/login", loginPost).Methods("POST")
	rtr.HandleFunc("/logout", logout).Methods("GET")

	rtr.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	rtr.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("/data/uploads"))))

	http.Handle("/", rtr)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	initDB()
	handleFunc()
}
