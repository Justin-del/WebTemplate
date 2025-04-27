package SignUp

import (
	AuthenticationChallenges "WebTemplate/Database/AuthenticationChallenges"
	users "WebTemplate/Database/Users"
	webauthn "WebTemplate/Utils/WebAuthn"
	"WebTemplate/globals"
	"encoding/json"
	"html/template"
	"net/http"
)

func HandleRoutes() {
	http.HandleFunc("GET /SignUp", func(responseWriter http.ResponseWriter, request *http.Request) {
		AuthenticationChallenges.DeleteAnyExpiredChallenges()
		t, _ := template.ParseFiles("./templates/base.html", "./templates/SignUp.html")
		t.ExecuteTemplate(responseWriter, "SignUp.html", nil)
	})

	http.HandleFunc("GET /SignUp/RegistrationData", func(responseWriter http.ResponseWriter, request *http.Request) {
		AuthenticationChallenges.DeleteAnyExpiredChallenges()
		challenge := AuthenticationChallenges.CreateNewChallenge()

		data := webauthn.RegistrationData{
			Challenge:               challenge,
			RP:                      webauthn.RP,
			SupportedCoseAlgorithms: webauthn.ListOfSupportedCoseAlgorithms,
			TimeoutInMinutes:        webauthn.TimeoutInMinutes,
		}

		responseWriter.Header().Set("Content-Type", "application/json")
		json.NewEncoder(responseWriter).Encode(data)
	})

	http.HandleFunc("POST /SignUp/{challengeId}/{userId}", func(responseWriter http.ResponseWriter, request *http.Request) {
		AuthenticationChallenges.DeleteAnyExpiredChallenges()

		var publicKeyCredential map[string]any
		err := json.NewDecoder(request.Body).Decode(&publicKeyCredential)

		if err != nil {
			http.Error(responseWriter, "Error", 400)
		}

		isOperationSuccesful := webauthn.SignUp(globals.OriginOfServer, request.PathValue("userId"), request.PathValue("challengeId"), AuthenticationChallenges.DeleteChallengeByID, publicKeyCredential, users.AddUserIntoDatabaseWithCredentials)

		if isOperationSuccesful {
			responseWriter.WriteHeader(200)
		} else {
			http.Error(responseWriter, "Error", 400)
		}
	})
}
