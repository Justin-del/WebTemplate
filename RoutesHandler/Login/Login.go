package login

import (
	AuthenticationChallenges "WebTemplate/Database/AuthenticationChallenges"
	Credentials "WebTemplate/Database/Credentials"
	Sessions "WebTemplate/Database/Sessions"
	TemplateParser "WebTemplate/RoutesHandler/TemplateParser"
	webauthn "WebTemplate/Utils/WebAuthn"
	"encoding/json"
	"net/http"
)

func HandleRoutes() {
	http.HandleFunc("GET /login", func(responseWriter http.ResponseWriter, request *http.Request) {
		TemplateParser.ExecuteTemplate("login", "login", responseWriter, request)
	})

	http.HandleFunc("GET /login/AuthenticationData", func(responseWriter http.ResponseWriter, request *http.Request) {
		challenge := AuthenticationChallenges.CreateNewChallenge()

		data := webauthn.AuthenticationData{
			Challenge:        challenge,
			RelyingPartyId:   webauthn.RP.Id,
			TimeoutInMinutes: webauthn.TimeoutInMinutes,
		}

		responseWriter.Header().Set("Content-Type", "application/json")
		json.NewEncoder(responseWriter).Encode(data)
	})

	http.HandleFunc("POST /login/{challengeId}", func(responseWriter http.ResponseWriter, request *http.Request) {
		var publicKeyCredential map[string]any
		json.NewDecoder(request.Body).Decode(&publicKeyCredential)

		userId := webauthn.Authenticate(publicKeyCredential, Credentials.GetPublicKeyAndSignatureCounter, AuthenticationChallenges.DeleteChallengeByID, request.PathValue("challengeId"), Credentials.UpdateSignatureCounter)

		if userId == "" {
			responseWriter.WriteHeader(401)
		} else {
			sessionId := Sessions.CreateASession(userId)

			http.SetCookie(responseWriter, &http.Cookie{
				Name:     "session_id",
				Value:    sessionId,
				Path:     "/",
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			})

			responseWriter.WriteHeader(200)
		}

	})

}
