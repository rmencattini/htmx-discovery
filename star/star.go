package star

import (
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
    "htmx/templates"
)

type Star struct {
        Time string
        Top int
        Left int
        StarType string
        Rotate int
        Scale float32
}

var POSSIBLE_ROTATE [6]int = [6]int{ 0, 30, 15, 0, 60, 0 }
// TODO: add new type of star and range of color
// TODO: store the star in case of reload
//  * load all stars on refresh page
// TODO: make the star pulsing
func AddTemporaryStar(rp http.ResponseWriter, _ *http.Request) {
    tmpl, err := template.ParseFS(templates.TemplatesFolder, "star.html")
    if err != nil {
        log.Println(err)
        return
    }
    timeout := strconv.Itoa(10 + rand.Intn(10)) + "s"
    top := rand.Intn(90)
    left := rand.Intn(90)
    scale := float32(rand.Intn(4)) + rand.Float32()
    err = tmpl.ExecuteTemplate(rp, "star", Star{
        Time: timeout,
        Top: top,
        Left: left,
        StarType: "north-star",
        Rotate: POSSIBLE_ROTATE[rand.Intn(6)],
        Scale: scale,
    })
    if err != nil {
        log.Println(err)
    }
}
