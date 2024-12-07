package forum

import (
	"net/http"
	"text/template"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "err", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/" {

		http.Error(w, "404", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("login/index.html")
	if err != nil {

		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}

	tmpl.Execute(w, nil)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear the username cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "username",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	// Redirect to the login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
