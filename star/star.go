package star

import (
	"database/sql"
	"html/template"
	"htmx/templates"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Star struct {
	Id       string
	Time     float32
	Top      int
	Left     int
	StarType string
	Rotate   int
	Scale    float32
}

var POSSIBLE_ROTATE [6]int = [6]int{0, 30, 15, 0, 60, 0}

func AddTemporaryStar(rp http.ResponseWriter, _ *http.Request) {
	timeout := float32(5 + rand.Intn(10))
	top := rand.Intn(90)
	left := rand.Intn(90)
	scale := float32(rand.Intn(4)) + rand.Float32()
	star := Star{
		Id:       uuid.NewString(),
		Time:     timeout,
		Top:      top,
		Left:     left,
		StarType: "north-star",
		Rotate:   POSSIBLE_ROTATE[rand.Intn(6)],
		Scale:    scale,
	}
	err := InsertStarInDB(star)
	if err != nil {
		return
	}
	insertStarsInTemplate([]Star{star}, rp)
}

func AddExistingStars(rp http.ResponseWriter, _ *http.Request) {
	stars := GetAllStars()
	insertStarsInTemplate(stars, rp)
}

func GetAllStars() []Star {
	results := []Star{}
	db, err := sql.Open("sqlite3", "test.db")

	if err != nil {
		log.Println(err)
		return results
	}
	currentTime := time.Now().UnixMilli()
	_, err = db.Exec("delete from star where ((time * 1000) + created_at) < ?", currentTime)
	rows, err := db.Query("SELECT id, time, top, left, star_type, rotate, scale, created_at from star")
	if err != nil {
		log.Println(err)
		return results
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		var time float32
		var top int
		var left int
		var starType string
		var rotate int
		var scale float32
		var created_at int64

		err = rows.Scan(&id, &time, &top, &left, &starType, &rotate, &scale, &created_at)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, Star{
			Id:       id,
			Top:      top,
			Time:     time - float32(currentTime-created_at)/1000,
			Left:     left,
			StarType: starType,
			Rotate:   rotate,
			Scale:    scale,
		})
	}
	return results
}

func insertStarsInTemplate(stars []Star, rp http.ResponseWriter) {
	tmpl, err := template.ParseFS(templates.TemplatesFolder, "star.html")
	if err != nil {
		log.Println(err)
		return
	}
	for _, star := range stars {
		err = tmpl.ExecuteTemplate(rp, "star", star)
		if err != nil {
			log.Println(err)
		}
	}
}

func InsertStarInDB(star Star) error {
	db, err := sql.Open("sqlite3", "test.db")

	if err != nil {
		log.Println(err)
		return err
	}
	_, err = db.Exec("INSERT INTO star(id, time, top, left, star_type, rotate, scale, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		star.Id, star.Time, star.Top, star.Left, star.StarType, star.Rotate, star.Scale, time.Now().UnixMilli())

	if err != nil {
		log.Println(err)
		return err
	}
	defer db.Close()
	return err
}
