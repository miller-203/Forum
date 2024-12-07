package forum

import (
	"net/http"
	"text/template"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Post struct {
	ID         int
	Username   string
	Title      string
	Content    string
	Categories string
	CreatedAt  string
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("username")
	if err != nil || cookie.Value == "" {
		// If no user is logged in, redirect to login page
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Handle POST request for adding a new post
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")
		categories := r.FormValue("category")

		if title == "" || content == "" {
			http.Error(w, "Title and content cannot be empty", http.StatusBadRequest)
			return
		}

		// Insert the new post into the database
		_, err := db.Exec(`
            INSERT INTO posts (username, title ,content,categories, created_at)
            VALUES (?, ?, ?, ?, ?)`,
			cookie.Value, title, content, categories, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			http.Error(w, "Failed to insert post", http.StatusInternalServerError)
			return
		}

		// Redirect to the forum page to see the new post
		http.Redirect(w, r, "/forum", http.StatusSeeOther)
		return
	}

	// Handle GET request to display posts
	rows, err := db.Query("SELECT id, username, title, content, Categories, created_at FROM posts ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Username, &post.Title, &post.Content, &post.Categories, &post.CreatedAt); err != nil {
			http.Error(w, "Failed to scan post", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error reading posts", http.StatusInternalServerError)
		return
	}

	// Pass the posts data to the template
	data := struct {
		Username string
		Posts    []Post
	}{
		Username: cookie.Value,
		Posts:    posts,
	}

	tmpl, err := template.ParseFiles("forum/forum.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}
