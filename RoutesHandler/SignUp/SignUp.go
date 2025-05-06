package SignUp

import (
	AuthenticationChallenges "WebTemplate/Database/AuthenticationChallenges"
	Users "WebTemplate/Database/Users"
	TemplateParser "WebTemplate/TemplateParser"
	WebAuthn "WebTemplate/Utils/WebAuthn"
	Globals "WebTemplate/Globals"
	"encoding/json"
	"net/http"
)

func HandleRoutes() {
	http.HandleFunc("GET /signUp", func(responseWriter http.ResponseWriter, request *http.Request) {
		TemplateParser.ExecuteTemplate("signUp", "Sign Up", responseWriter, request)
	})

	http.HandleFunc("GET /signUp/RegistrationData", func(responseWriter http.ResponseWriter, request *http.Request) {
		challenge := AuthenticationChallenges.CreateNewChallenge()

		data := WebAuthn.RegistrationData{
			Challenge:               challenge,
			RP:                      WebAuthn.RP,
			SupportedCoseAlgorithms: WebAuthn.ListOfSupportedCoseAlgorithms,
			TimeoutInMinutes:        WebAuthn.TimeoutInMinutes,
		}

		responseWriter.Header().Set("Content-Type", "application/json")
		json.NewEncoder(responseWriter).Encode(data)
	})

	http.HandleFunc("POST /signUp/{challengeId}/{userId}", func(responseWriter http.ResponseWriter, request *http.Request) {

		var publicKeyCredential map[string]any
		err := json.NewDecoder(request.Body).Decode(&publicKeyCredential)

		if err != nil {
			http.Error(responseWriter, "Error", 400)
		}

		isOperationSuccesful := WebAuthn.SignUp(Globals.OriginOfServer, request.PathValue("userId"), request.PathValue("challengeId"), AuthenticationChallenges.DeleteChallengeByID, publicKeyCredential, Users.AddUserIntoDatabaseWithCredentials)

		if isOperationSuccesful {
			responseWriter.WriteHeader(200)
		} else {
			http.Error(responseWriter, "Error", 400)
		}
	})
}
