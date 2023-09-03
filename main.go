package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"htmx/star"
	"htmx/templates"
	"log"
	"net/http"
)

func main() {
	db, err := sql.Open("sqlite3", "test.db")

	if err != nil {
		log.Fatal(err)
	}

	initDatabase(db)
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/", mainPage)
	mux.HandleFunc("/stars", star.AddExistingStars)
	mux.HandleFunc("/click", star.AddTemporaryStar)
	http.ListenAndServe("localhost:1234", mux)
}

func mainPage(rp http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFS(templates.TemplatesFolder, "index.html")
	if err != nil {
		log.Println(err)
		return
	}

	tmpl, err = tmpl.New("star").ParseFS(templates.TemplatesFolder, "star.html")
	if err != nil {
		log.Println(err)
		return
	}
	stars := star.GetAllStars()

	err = tmpl.ExecuteTemplate(rp, "index", struct{ Star []star.Star }{Star: stars})
	if err != nil {
		log.Println(err)
	}
}

func initDatabase(db *sql.DB) {
	createRowsIfNotExist := `
        CREATE TABLE IF NOT EXISTS star(
            id TEXT PRIMARY KEY,
            created_at INTEGER,
            time INTEGER,
            top INTEGER,
            left INTEGER,
            star_type TEXT,
            rotate INTEGER,
            scale REAL
        );
    `
	_, err := db.Exec(createRowsIfNotExist)
	if err != nil {
		log.Println(err)
	}
}
