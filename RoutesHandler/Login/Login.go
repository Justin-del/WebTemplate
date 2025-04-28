package login

import (
	AuthenticationChallenges "WebTemplate/Database/AuthenticationChallenges"
	Credentials "WebTemplate/Database/Credentials"
	Sessions "WebTemplate/Database/Sessions"
	TemplateParser "WebTemplate/RoutesHandler/TemplateParser"
	webauthn "WebTemplate/Utils/WebAuthn"
	"WebTemplate/globals"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
)

func HandleRoutes() {
	http.HandleFunc("GET /login", func(responseWriter http.ResponseWriter, request *http.Request) {
		TemplateParser.ParseTemplate("login", responseWriter, request)
	})

	http.HandleFunc("GET /login/AuthenticationData", func(responseWriter http.ResponseWriter, request *http.Request) {
		AuthenticationChallenges.DeleteAnyExpiredChallenges()
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
		AuthenticationChallenges.DeleteAnyExpiredChallenges()
		var publicKeyCredential map[string]any
		json.NewDecoder(request.Body).Decode(&publicKeyCredential)

		response := publicKeyCredential["response"].(map[string]any)
		rawId := publicKeyCredential["rawId"].(string)
		signature := response["signature"].(string)
		userHandle := response["userHandle"].(string)

		decodedRawId, err1 := base64.RawURLEncoding.DecodeString(rawId)
		decodedSignature, err2 := base64.RawURLEncoding.DecodeString(signature)
		decodedUserHandle, err3 := base64.RawURLEncoding.DecodeString(userHandle)

		if err1 != nil {
			http.Error(responseWriter, "Bad request", 400)
		}

		if err2 != nil {
			http.Error(responseWriter, "Bad request", 400)
		}

		if err3 != nil {
			http.Error(responseWriter, "Bad request", 400)
		}

		publicKey, signatureCounterFromServer := Credentials.GetPublicKeyAndSignatureCounter(decodedRawId, decodedUserHandle)

		correct_challenge := AuthenticationChallenges.DeleteChallengeByID(request.PathValue("challengeId"))

		var clientDataJSON string = response["clientDataJSON"].(string)

		isClientDataJSONCorrect := webauthn.IsClientDataJSONCorrect(globals.OriginOfServer, correct_challenge, "webauthn.get", clientDataJSON)

		authData := response["authenticatorData"].(string)
		decodedAuthData, _ := base64.RawURLEncoding.DecodeString(authData)

		rpIdHash := decodedAuthData[0:32]
		decodedClientData, err4 := base64.RawURLEncoding.DecodeString(clientDataJSON)

		if err4 != nil {
			http.Error(responseWriter, "Bad request", 400)
		}

		signatureCounterFromAuthenticator := binary.BigEndian.Uint32(decodedAuthData[33:37])

		isAuthenticated := webauthn.IsRpIdHashCorrect(rpIdHash) && isClientDataJSONCorrect && webauthn.AreFlagsValid(decodedAuthData[32]) && webauthn.IsSignatureVerified(append(decodedAuthData, webauthn.Sha256Hash(decodedClientData)...), decodedSignature, publicKey) && webauthn.IsSignatureCounterValid(signatureCounterFromServer, signatureCounterFromAuthenticator)

		Credentials.UpdateSignatureCounter(decodedRawId, signatureCounterFromAuthenticator)
		fmt.Println(isAuthenticated)

		if isAuthenticated {
			sessionId := Sessions.CreateASession(string(decodedUserHandle))
			fmt.Println("The session id is ", sessionId)

			http.SetCookie(responseWriter, &http.Cookie{
				Name:     "session_id",
				Value:    sessionId,
				Path:     "/",
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			})

			responseWriter.WriteHeader(200)
		} else {
			responseWriter.WriteHeader(401)
		}
	})

}
