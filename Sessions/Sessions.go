package sessions

import (
	Sessions "WebTemplate/Database/Sessions"
	"context"
	"fmt"
	"net/http"
)

func SessionMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		if len(request.URL.Path) >= 7 && request.URL.Path[:7] == "/static" {
			next.ServeHTTP(responseWriter, request)
			return
		}

		if len(request.URL.Path) >= 6 && request.URL.Path[:6] == "/robot" {
			next.ServeHTTP(responseWriter, request)
			return
		}

		if len(request.URL.Path) >= 12 && request.URL.Path[:12] == "/favicon.ico" {
			next.ServeHTTP(responseWriter, request)
			return
		}

		fmt.Println("URL of the request is ", request.URL.String())
		fmt.Println("I am called.")
		cookie, err := request.Cookie("session_id")
		var sessionId string

		if err == nil {
			sessionId = cookie.Value
		}

		if !Sessions.DoesSessionExistsInDatabase(sessionId) {
			sessionId = ""
		}

		ctx := request.Context()
		ctx = context.WithValue(ctx, "sessionId", sessionId)

		next.ServeHTTP(responseWriter, request.WithContext(ctx))
	})
}
