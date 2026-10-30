package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
)

//go:embed templates/*
var templatesFS embed.FS

func main() {
	if err := implementfetch(); err != nil {
		log.Fatal(err)
	}

	var err error
	tmpl, err = template.ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("GET /styling.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, templatesFS, "templates/styling.css")
	})
	http.HandleFunc("GET /{$}", indexhandler)
	http.HandleFunc("GET /artist", artistHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("server is running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
