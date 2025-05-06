package main

import (
	"TodoApp/Database"
	RoutesHandler "TodoApp/RoutesHandler"
	"net/http"
	"os"
)

func main() {

	os.Setenv("CGO_ENABLED", "1")
	Database.CreateTablesIfNotExists()

	RoutesHandler.HandleRoutes()

	http.ListenAndServe(":8080", nil)
}
