package templateparser

import (
	Sessions "WebTemplate/Database/Sessions"
	"html/template"
	"net/http"
)

var templatesFolder string = "templates"
var baseTemplate string = "base"

func ParseTemplate(templateName string, responseWriter http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("session_id")

	var sessionId string

	if err == nil {
		sessionId = cookie.Value
	}

	type Data struct {
		IsLoggedIn bool
	}

	data := Data{
		IsLoggedIn: Sessions.DoesSessionExistsInDatabase(sessionId),
	}

	if data.IsLoggedIn {
		//Disable caching of sensitive data.
		responseWriter.Header().Add("Cache-control", "no-store")
	}

	responseWriter.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
	responseWriter.Header().Add("Cross-Origin-Opener-Policy", "same-origin")
	responseWriter.Header().Add("X-Frame-Options", "deny")

	t, err := template.ParseFiles("./templates/"+baseTemplate+".html", "./templates/"+templateName+".html")

	err = t.ExecuteTemplate(responseWriter, templateName+".html", data)
}
