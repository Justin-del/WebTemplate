package SignUp

import (
	AuthenticationChallenges "WebTemplate/Database/AuthenticationChallenges"
	users "WebTemplate/Database/Users"
	globals "WebTemplate/Globals"
	TemplateParser "WebTemplate/TemplateParser"
	webauthn "WebTemplate/Utils/WebAuthn"
	"encoding/json"
	"net/http"
)

func HandleRoutes() {
	http.HandleFunc("GET /signUp", func(responseWriter http.ResponseWriter, request *http.Request) {
		TemplateParser.ExecuteTemplate("signUp", "Sign Up", responseWriter, request)
	})

	http.HandleFunc("GET /signUp/RegistrationData", func(responseWriter http.ResponseWriter, request *http.Request) {
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

	http.HandleFunc("POST /signUp/{challengeId}/{userId}", func(responseWriter http.ResponseWriter, request *http.Request) {

		var decodedJSON map[string]any
		err := json.NewDecoder(request.Body).Decode(&decodedJSON)

		if err != nil {
			http.Error(responseWriter, "Error", 400)
		}

		isOperationSuccesful := webauthn.SignUp(globals.OriginOfServer, request.PathValue("userId"), request.PathValue("challengeId"), AuthenticationChallenges.DeleteChallengeByID, decodedJSON, users.AddUserIntoDatabaseWithCredentials)

		if isOperationSuccesful {
			responseWriter.WriteHeader(200)
		} else {
			http.Error(responseWriter, "Error", 400)
		}
	})

	http.HandleFunc("POST /signUp/isUsernameTaken", func(responseWriter http.ResponseWriter, request *http.Request) {
		var decodedJSON map[string]any
		err := json.NewDecoder(request.Body).Decode(&decodedJSON)

		if err != nil {
			http.Error(responseWriter, "Error", 400)
		}

		json.NewEncoder(responseWriter).Encode(map[string]bool{"isUsernameTaken": users.DoesUserNameExistsInDatabase(decodedJSON["username"].(string))})

	})
}
