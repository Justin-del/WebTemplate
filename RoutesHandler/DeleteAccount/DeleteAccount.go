package deleteaccount

import (
	AuthenticationChallenges "WebTemplate/Database/AuthenticationChallenges"
	Credentials "WebTemplate/Database/Credentials"
	Users "WebTemplate/Database/Users"
	TemplateParser "WebTemplate/TemplateParser"
	webauthn "WebTemplate/Utils/WebAuthn"
	"encoding/json"
	"net/http"
)

func HandleRoutes() {
	http.HandleFunc("GET /deleteAccount", func(responseWriter http.ResponseWriter, request *http.Request) {
		TemplateParser.ExecuteTemplate("DeleteAccount", "Delete Account", responseWriter, request)
	})

	http.HandleFunc("POST /deleteAccount/{challengeId}", func(responseWriter http.ResponseWriter, request *http.Request) {
		var publicKeyCredential map[string]any
		json.NewDecoder(request.Body).Decode(&publicKeyCredential)

		userId := webauthn.Authenticate(publicKeyCredential, Credentials.GetPublicKeyAndSignatureCounter, AuthenticationChallenges.DeleteChallengeByID, request.PathValue("challengeId"), Credentials.UpdateSignatureCounter)

		if userId == "" {
			responseWriter.WriteHeader(401)
		} else {
			Users.DeleteUser(userId)
			responseWriter.WriteHeader(200)
		}

	})
}
