package authorized

import (
	TemplateParser "WebTemplate/TemplateParser"
	"net/http"
)

func HandleRoutes() {

	http.HandleFunc("GET /authorized", func(responseWriter http.ResponseWriter, request *http.Request) {
		//An example of how to control access to authorized pages.
		sessionId := request.Context().Value("sessionId").(string)
		if sessionId == "" {
			http.Redirect(responseWriter, request, "/login", http.StatusSeeOther)
			return
		}

		TemplateParser.ExecuteTemplate("Authorized", "Authorized", responseWriter, request)
	})
}
