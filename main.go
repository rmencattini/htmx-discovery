package main

import (
	"database/sql"
	"html/template"
	"htmx/star"
	"htmx/templates"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "test.db")

	if err != nil {
		log.Fatal(err)
	}

	initDatabase(db)
	defer db.Close()

	mux := http.NewServeMux()

    go createStar()

	mux.HandleFunc("/", mainPage)
	mux.HandleFunc("/stars", star.AddExistingStars)
	mux.HandleFunc("/click", star.AddTemporaryStar)
	http.ListenAndServe("localhost:1234", mux)
}

func createStar() {
    for {
        stars := star.GetAllStars()
        maxStarsNumber := 30

        diff := maxStarsNumber - len(stars)
        if diff <= 0 {
            time.Sleep(100 * time.Millisecond)
            continue
        } else {
            timeout := float32(5 + rand.Intn(10))
            top := rand.Intn(90)
            left := rand.Intn(90)
            scale := float32(rand.Intn(4)) + rand.Float32()
            temporaryStar := star.Star{
                Id:       uuid.NewString(),
                Time:     timeout,
                Top:      top,
                Left:     left,
                StarType: "north-star",
                Rotate:   star.POSSIBLE_ROTATE[rand.Intn(6)],
                Scale:    scale,
            }
            err := star.InsertStarInDB(temporaryStar)
            if err != nil {
                log.Println(err)
            }

        }
    }
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
