package main

import (
	"html/template"
	"htmx/star"
	"htmx/templates"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", mainPage)
	mux.HandleFunc("/click", star.AddTemporaryStar)
	http.ListenAndServe("localhost:1234", mux)
}

func mainPage(rp http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFS(templates.TemplatesFolder, "index.html")
	if err != nil {
		log.Println(err)
		return
	}
	err = tmpl.ExecuteTemplate(rp, "index", nil)
	if err != nil {
		log.Println(err)
	}
}


