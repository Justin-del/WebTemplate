package logout

import (
	Sessions "TodoApp/Database/Sessions"
	"net/http"
)

func HandleRoutes() {
	http.HandleFunc("POST /logout", func(responseWriter http.ResponseWriter, request *http.Request) {
		cookie, err := request.Cookie("session_id")
		if err != nil {
			return
		}
		sessionId := cookie.Value
		Sessions.DeleteSession(sessionId)
		responseWriter.Header().Add("Hx-Redirect", "/login")
	})
}
