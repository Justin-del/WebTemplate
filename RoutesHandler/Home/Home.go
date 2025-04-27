package home

import (
	"html/template"
	"net/http"
)

func HandleRoutes() {
	http.HandleFunc("GET /", func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/" {
			http.NotFound(responseWriter, request)
			return
		}
		t, _ := template.ParseFiles("./templates/base.html", "./templates/index.html")
		t.ExecuteTemplate(responseWriter, "index.html", nil)
	})
}
