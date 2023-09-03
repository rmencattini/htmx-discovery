package main

import (
	"database/sql"
	"html/template"
	"htmx/star"
	"htmx/templates"
	"log"
	"math"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var startTime int64

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
	mux.HandleFunc("/click", addStar)
	startTime = time.Now().UnixMilli()
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
			star.InsertStarInDb()
		}
	}
}

func addStar(rp http.ResponseWriter, _ *http.Request) {
    temporyStar, err := star.InsertStarInDb()
    if err != nil {
        log.Println(err)
        return
    }
    star.InsertStarsInTemplate([]star.Star{temporyStar}, rp)
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
	rotate := math.Mod(float64(time.Now().UnixMilli()-startTime)/1000.0, 360)

	err = tmpl.ExecuteTemplate(rp, "index", struct {
		Star        []star.Star
		RotateStart float64
		RotateEnd   float64
	}{
		Star:        stars,
		RotateStart: rotate,
		RotateEnd:   rotate + 360,
	})
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
