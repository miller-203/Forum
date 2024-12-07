package main

import (
	forum "forum/forum/res"
	"log"
	"net/http"
)

func main() {

	fs := http.FileServer(http.Dir("./login/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	err := forum.InitDB()
	if err != nil {
		log.Fatal("Error initializing database:", err)
	}
	http.HandleFunc("/logout", forum.LogoutHandler)
	// http.HandleFunc("/forum", forum.ForumHandler)
	http.HandleFunc("/forum", forum.GetPosts)
	http.HandleFunc("/", forum.LoginHandler)

	http.HandleFunc("/profile", forum.ProfileHandler)

	// http.HandleFunc("/posts", forum.AddPost)
	http.HandleFunc("/signup", forum.SignupHandler)
	forum.Print()
	http.ListenAndServe(":8080", nil)

}
