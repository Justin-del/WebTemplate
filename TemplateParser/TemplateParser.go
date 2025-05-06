package templateparser

import (
	Sessions "WebTemplate/Database/Sessions"
	"html/template"
	"net/http"
)

var signUpTemplateFiles []string = []string{"./templates/base.html", "./templates/SignUp.html"}
var loginTemplateFiles []string = []string{"./templates/base.html", "./templates/Login.html"}
var deleteAccountTemplateFiles []string = []string{"./templates/base.html", "./templates/DeleteAccount.html"}
var homeTemplateFiles []string = []string{"./templates/base.html", "./templates/index.html"}
var authorizedTemplateFiles []string = []string{"./templates/base.html", "./templates/Authorized.html"}

/*
The key will be the name of the template file to execute (without HTML extension) .
*/
var templatesMap map[string]*template.Template = make(map[string]*template.Template)

var shouldReparseTemplateOnEveryRequest = false

// You can set this variable to true before calling ExecuteTemplate or ExecuteTemplateWithAdditionalData if you know that the user is already logged in. This will prevent database lookup for checking if the user is already logged in when ExecuteTemplate or ExecuteTemplateWithAdditionalData is called.
var IsLoggedIn *bool = nil

func InitTemplatesMap() {
	templatesMap["SignUp"], _ = template.ParseFiles(signUpTemplateFiles...)
	templatesMap["Login"], _ = template.ParseFiles(loginTemplateFiles...)
	templatesMap["DeleteAccount"], _ = template.ParseFiles(deleteAccountTemplateFiles...)
	templatesMap["index"], _ = template.ParseFiles(homeTemplateFiles...)
	templatesMap["Authorized"], _ = template.ParseFiles(authorizedTemplateFiles...)
}

/*
The templateName is the name of the template file that you would like to execute (without the HTML extension).
*/
func ExecuteTemplate(templateName string, pageName string, responseWriter http.ResponseWriter, request *http.Request) {
	ExecuteTemplateWithAdditionalData(templateName, pageName, responseWriter, request, nil)
}

/*
The templateName is the name of the template file that you would like to execute (without the HTML extension).
*/
func ExecuteTemplateWithAdditionalData(templateName string, pageName string, responseWriter http.ResponseWriter, request *http.Request, additionalData any) {
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

	var data Data

	if IsLoggedIn == nil {
		data = Data{
			IsLoggedIn:     Sessions.DoesSessionExistsInDatabase(sessionId),
			PageName:       pageName,
			AdditionalData: additionalData,
		}
	} else {
		data = Data{
			IsLoggedIn:     *IsLoggedIn,
			PageName:       pageName,
			AdditionalData: additionalData,
		}
	}

	//This variable must be set to nil after it is used because whether the user is logged in or not when this function is called again is unknown.
	IsLoggedIn = nil

	if data.IsLoggedIn {
		//Disable caching of sensitive data.
		responseWriter.Header().Add("Cache-control", "no-store")
	}

	responseWriter.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
	responseWriter.Header().Add("Cross-Origin-Opener-Policy", "same-origin")
	responseWriter.Header().Add("X-Frame-Options", "deny")

	if shouldReparseTemplateOnEveryRequest {
		InitTemplatesMap()
	}

	err = templatesMap[templateName].ExecuteTemplate(responseWriter, templateName+".html", additionalData)
}
