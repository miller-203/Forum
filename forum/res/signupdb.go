package forum

import (
	"database/sql"
	forum "forum/forum/utils"
	"time"

	"html/template"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db *sql.DB
)

type Posts struct {
	Title   string
	Content string
}
type User struct {
	Username string
	Email    string
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		if len(username) < 3 || len(email) < 3 || len(password) < 3 {
			http.Error(w, "username or email too short.", http.StatusBadRequest)
			return
		}
		if !forum.IsValidUsername(username) {
			http.Error(w, "Invalid username. Only alphanumeric characters and underscores are allowed.", http.StatusBadRequest)
			return
		}

		hashedPassword, err := forum.HashPassword(password)
		if err != nil {
			log.Println("Error hashing password:", err)
			http.Error(w, "Error during sign-up", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec(`
            INSERT INTO users (username, email, password) 
            VALUES (?, ?, ?)`,
			username, email, hashedPassword)

		if err != nil {
			log.Println("Error during sign-up:", err)
			http.Error(w, "Error during sign-up", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
	tmpl := template.Must(template.ParseFiles("login/signup.html"))
	tmpl.Execute(w, nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		email := strings.TrimSpace(r.FormValue("email"))
		password := r.FormValue("password")

		var dbPassword, username string
		err := db.QueryRow("SELECT password, username FROM users WHERE email = ?", email).Scan(&dbPassword, &username)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Email not found", http.StatusUnauthorized)
			} else {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		if !forum.VerifyPassword(dbPassword, password) {
			http.Error(w, "Incorrect password", http.StatusUnauthorized)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "username",
			Value: username,
			Path:  "/",
		})

		http.Redirect(w, r, "/forum", http.StatusSeeOther)
	}
	tmpl := template.Must(template.ParseFiles("login/login.html"))
	tmpl.Execute(w, nil)
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the username from the cookie
	cookie, err := r.Cookie("username")
	if err != nil {
		log.Printf("Cookie error: %v", err)
		http.Redirect(w, r, "/forum", http.StatusSeeOther)
		return
	}
	username := cookie.Value

	// Connect to SQLite database
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Printf("Database connection error: %v", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Query the specific user by username
	row := db.QueryRow("SELECT username, email FROM users WHERE username = ?", username)

	var user User
	err = row.Scan(&user.Username, &user.Email)
	if err == sql.ErrNoRows {
		log.Printf("No user found for username: %v", username)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, "Error fetching user data", http.StatusInternalServerError)
		return
	}

	// Parse and execute template
	tmpl, err := template.ParseFiles("forum/profile.html")
	if err != nil {
		log.Printf("Template parsing error: %v", err)
		http.Error(w, "Template parsing error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, user)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		return
	}
}
func AddPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		cookie, err := r.Cookie("username")
		if err != nil || cookie.Value == "" {
			http.Error(w, "User not authenticated", http.StatusUnauthorized)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		// Parse form data

		// Retrieve title and content from the form
		title := r.FormValue("title")
		postContent := r.FormValue("content")
		Categories := r.FormValue("category")

		// Check if the post content or title is empty
		if postContent == "" || title == "" {
			http.Error(w, "Post content and title cannot be empty", http.StatusBadRequest)
			return
		}

		// Retrieve userID from the cookie

		// Insert the new post into the database
		_, err = db.Exec(`
            INSERT INTO posts (username, title, content,Categories, created_at)
            VALUES (?, ?, ?, ?, ?)`,
			cookie.Value, title, postContent, Categories, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			http.Error(w, "Failed to insert post", http.StatusInternalServerError)
			return
		}

		// Redirect to the forum page after successful insertion
		http.Redirect(w, r, "/forum", http.StatusSeeOther)
		return
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
