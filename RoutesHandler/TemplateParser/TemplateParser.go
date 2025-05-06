package templateparser

import (
	Sessions "TodoApp/Database/Sessions"
	"html/template"
	"net/http"
	"strings"
)

var templatesFolder string = "templates"
var baseTemplate string = "base"

func ParseTemplate(templateName string, pageName string, responseWriter http.ResponseWriter, request *http.Request) {
	ParseTemplateWithAdditionalData(templateName, pageName, responseWriter, request, nil)
}

func ParseTemplateWithAdditionalData(templateName string, pageName string, responseWriter http.ResponseWriter, request *http.Request, additionalData any) {
	cookie, err := request.Cookie("session_id")

	var sessionId string

	if err == nil {
		sessionId = cookie.Value
	}

	type Data struct {
		IsLoggedIn     bool
		PageName       string
		AdditionalData any
	}

	data := Data{
		IsLoggedIn:     Sessions.DoesSessionExistsInDatabase(sessionId),
		PageName:       pageName,
		AdditionalData: additionalData,
	}

	if data.IsLoggedIn {
		//Disable caching of sensitive data.
		responseWriter.Header().Add("Cache-control", "no-store")
	}

	responseWriter.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
	responseWriter.Header().Add("Cross-Origin-Opener-Policy", "same-origin")
	responseWriter.Header().Add("X-Frame-Options", "deny")

	t, err := template.ParseFiles("./templates/"+baseTemplate+".html", "./templates/"+templateName+".html")

	lastSlashIndex := strings.LastIndex(templateName, "/")
	if lastSlashIndex != -1 {
		templateName = templateName[lastSlashIndex+1:]
	}
	err = t.ExecuteTemplate(responseWriter, templateName+".html", data)

}
