package routeshandler

import (
	Authorized "TodoApp/RoutesHandler/Todo"
	DeleteAccount "TodoApp/RoutesHandler/DeleteAccount"
	Home "TodoApp/RoutesHandler/Home"
	Login "TodoApp/RoutesHandler/Login"
	Logout "TodoApp/RoutesHandler/Logout"
	SignUp "TodoApp/RoutesHandler/SignUp"
	StaticHandler "TodoApp/RoutesHandler/StaticHandler"
)

func HandleRoutes() {
	SignUp.HandleRoutes()
	Home.HandleRoutes()
	StaticHandler.HandleRoutes()
	Login.HandleRoutes()
	Logout.HandleRoutes()
	DeleteAccount.HandleRoutes()
	Authorized.HandleRoutes()
}
