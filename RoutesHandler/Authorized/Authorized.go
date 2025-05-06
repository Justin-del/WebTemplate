package authorized

import (
	Sessions "WebTemplate/Database/Sessions"
	TemplateParser "WebTemplate/TemplateParser"
	"net/http"
)

func HandleRoutes() {

	http.HandleFunc("GET /authorized", func(responseWriter http.ResponseWriter, request *http.Request) {
		//An example of how to control access to authorized pages.

		cookie, err := request.Cookie("session_id")
		if err != nil {
			http.Redirect(responseWriter, request, "/login", http.StatusSeeOther)
			return
		}

		sessionId := cookie.Value

		if Sessions.DoesSessionExistsInDatabase(sessionId) {
			TemplateParser.ExecuteTemplate("Authorized", "Authorized", responseWriter, request)
		} else {
			http.Redirect(responseWriter, request, "/login", http.StatusSeeOther)
		}

	})
}
