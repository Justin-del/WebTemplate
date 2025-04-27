package main

import (
	"WebTemplate/Database"
	RoutesHandler "WebTemplate/RoutesHandler"
	"net/http"
	"os"
)

func main() {

	os.Setenv("CGO_ENABLED", "1")
	Database.CreateTablesIfNotExists()

	RoutesHandler.HandleRoutes()

	http.ListenAndServe(":8080", nil)
}
