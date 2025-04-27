package routeshandler

import (
	Home "WebTemplate/RoutesHandler/Home"
	SignUp "WebTemplate/RoutesHandler/SignUp"
	StaticHandler "WebTemplate/RoutesHandler/StaticHandler"
	Login "WebTemplate/RoutesHandler/Login"
	Logout "WebTemplate/RoutesHandler/Logout"
	Authorized "WebTemplate/RoutesHandler/Authorized"
)

func HandleRoutes() {
	SignUp.HandleRoutes()
	Home.HandleRoutes()
	StaticHandler.HandleRoutes()
	Login.HandleRoutes()
	Logout.HandleRoutes()
	Authorized.HandleRoutes()
}
