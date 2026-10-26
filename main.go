package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	var err error
	err = implementfetch()
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err = template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal(err)
	}

	artistTmpl, err = template.ParseFiles("templates/artist.html")
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/styling.css", http.FileServer(http.Dir("templates")))
	http.HandleFunc("/", indexhandler)
	http.HandleFunc("/artist", artistHandler)

	log.Println("server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
