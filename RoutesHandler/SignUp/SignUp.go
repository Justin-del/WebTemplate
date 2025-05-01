package SignUp

import (
	AuthenticationChallenges "WebTemplate/Database/AuthenticationChallenges"
	users "WebTemplate/Database/Users"
	TemplateParser "WebTemplate/RoutesHandler/TemplateParser"
	webauthn "WebTemplate/Utils/WebAuthn"
	"WebTemplate/globals"
	"encoding/json"
	"net/http"
)

func HandleRoutes() {
	http.HandleFunc("GET /SignUp", func(responseWriter http.ResponseWriter, request *http.Request) {
		TemplateParser.ParseTemplate("SignUp", responseWriter, request)
	})

	http.HandleFunc("GET /SignUp/RegistrationData", func(responseWriter http.ResponseWriter, request *http.Request) {
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
