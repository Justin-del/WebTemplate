package main

import (
	"WebTemplate/Database"
	RoutesHandler "WebTemplate/RoutesHandler"
	TemplateParser "WebTemplate/TemplateParser"
	"net/http"
	"os"
)

func main() {

	os.Setenv("CGO_ENABLED", "1")
	Database.CreateTablesIfNotExists()

	RoutesHandler.HandleRoutes()
	TemplateParser.InitTemplatesMap()

	http.ListenAndServe(":8080", nil)
}
