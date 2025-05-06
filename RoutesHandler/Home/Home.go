package home

import (
	TemplateParser "TodoApp/RoutesHandler/TemplateParser"
	"net/http"
)

func HandleRoutes() {
	http.HandleFunc("GET /", func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/" {
			http.NotFound(responseWriter, request)
			return
		}
		TemplateParser.ParseTemplate("index", "Home", responseWriter, request)
	})
}
