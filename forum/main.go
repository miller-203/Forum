package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const port = ":7006"

type Text struct {
	ErrorNum int
	ErrorMes string
}

type User struct {
	ID             int
	Email          string
	Username       string
	Password       string
	Password_check string
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		t, err := template.ParseFiles("templates/errorPage.html")
		if err != nil {
			errorHandler(w, r, http.StatusInternalServerError)
			return
		}
		em := "HTTP status 404: Page Not Found"
		p := Text{ErrorNum: status, ErrorMes: em}
		t.Execute(w, p)
	}
	if status == http.StatusInternalServerError {
		t, err := template.ParseFiles("templates/errorPage.html")
		if err != nil {
			fmt.Fprint(w, "HTTP status 500: Internal Server Error -missing errorPage.html file")
		}
		em := "HTTP status 500: Internal Server Error"
		p := Text{ErrorNum: status, ErrorMes: em}
		t.Execute(w, p)
	}
	if status == http.StatusBadRequest {
		t, err := template.ParseFiles("templates/errorPage.html")
		if err != nil {
			fmt.Fprint(w, "HTTP status 500: Internal Server Error -missing errorPage.html file")
		}
		em := "HTTP status 400: Bad Request! Please select artist from the Home Page"
		p := Text{ErrorNum: status, ErrorMes: em}
		t.Execute(w, p)
	}
	if status == http.StatusForbidden {
		t, err := template.ParseFiles("templates/errorPage.html")
		if err != nil {
			fmt.Fprint(w, "HTTP status 500: Internal Server Error -missing errorPage.html file")
		}
		em := "HTTP status 403: Forbidden! Please login first"
		p := Text{ErrorNum: status, ErrorMes: em}
		t.Execute(w, p)
	}
}

func registerPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	tmpl, err := template.ParseFiles("templates/register.html")
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
	User_name := r.FormValue("username")
	Password := r.FormValue("password")
	fmt.Println(User_name, Password)
}

func displayPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	tmpl, err := template.ParseFiles("templates/post.html")
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", loginPage)
	http.HandleFunc("/register", registerPage)
	http.HandleFunc("/post", displayPost)
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("css"))))

	fmt.Println("Server listening at http://localhost" + port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
