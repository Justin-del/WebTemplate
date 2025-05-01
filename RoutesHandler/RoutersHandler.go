package routeshandler

import (
	Authorized "WebTemplate/RoutesHandler/Authorized"
	DeleteAccount "WebTemplate/RoutesHandler/DeleteAccount"
	Home "WebTemplate/RoutesHandler/Home"
	Login "WebTemplate/RoutesHandler/Login"
	Logout "WebTemplate/RoutesHandler/Logout"
	SignUp "WebTemplate/RoutesHandler/SignUp"
	StaticHandler "WebTemplate/RoutesHandler/StaticHandler"
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
