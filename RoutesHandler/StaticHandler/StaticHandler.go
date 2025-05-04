// This file is to handle all the static routes, except for the home page.
package StaticHandler

import (
	"net/http"
)

func HandleRoutes() {
	http.HandleFunc("GET /robots.txt", func(responseWriter http.ResponseWriter, request *http.Request) {
		http.ServeFile(responseWriter, request, "./static/robots.txt")
	})

	http.HandleFunc("GET /static/", func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Add("Service-Worker-Allowed", "/")
		http.ServeFile(responseWriter, request, "."+request.URL.Path)
	})
}
