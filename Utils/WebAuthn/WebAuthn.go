package webauthn

import (
	globals "WebTemplate/Globals"
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"

	"github.com/fxamacker/cbor/v2"
)

type RelyingParty struct {
	Name string
	Id   string
}

type Challenge struct {
	Id        int
	Challenge []byte
}

/*
Represents the data that is needed by navigator.credentials.create in the frontend
*/
type RegistrationData struct {
	Challenge               Challenge
	RP                      RelyingParty
	SupportedCoseAlgorithms []int64
	TimeoutInMinutes        int
}

/*
Represents the data that is needed by navigator.credentials.get in the frontend
*/
type AuthenticationData struct {
	Challenge        Challenge
	RelyingPartyId   string
	TimeoutInMinutes int
}

var ListOfSupportedCoseAlgorithms = []int64{-8, -7, -257}
var TimeoutInMinutes = 5

var RP = RelyingParty{
	Name: globals.ApplicationName,
	Id:   globals.DomainName,
}

type ClientData struct {
	Challenge []byte
	Type      string
	Origin    string
}

type AttestationObject struct {
	AttStmt  cbor.RawMessage
	AuthData []byte
	Fmt      string
}

func ParseClientDataJson(clientDataJSON string) ClientData {
	clientDataByteArray, err := base64.RawURLEncoding.DecodeString(clientDataJSON)
	if err != nil {
		return ClientData{}
	}

	var clientData map[string]any
	json.Unmarshal(clientDataByteArray, &clientData)

	challenge, err := base64.RawURLEncoding.DecodeString(clientData["challenge"].(string))
	if err != nil {
		return ClientData{}
	}

	origin := clientData["origin"].(string)

	type_string := clientData["type"].(string)

	return ClientData{
		Challenge: challenge,
		Type:      type_string,
		Origin:    origin,
	}

}

func ParseAttestationObject(attestationObjectString string) AttestationObject {
	decodedAttestationObjectString, err := base64.RawURLEncoding.DecodeString(attestationObjectString)
	if err != nil {
		return AttestationObject{}
	}

	var attestationObject AttestationObject
	cbor.NewDecoder(bytes.NewReader(decodedAttestationObjectString)).Decode(&attestationObject)
	return attestationObject
}

func IsClientDataJSONCorrect(correct_origin string, correct_challenge []byte, correct_type string, clientDataJSON string) bool {
	clientData := ParseClientDataJson(clientDataJSON)
	return clientData.Origin == correct_origin && bytes.Equal(clientData.Challenge, correct_challenge) && clientData.Type == correct_type
}

func Sha256Hash(input []byte) []byte {
	algorithm := sha256.New()
	algorithm.Write(input)
	return algorithm.Sum(nil)
}

func IsRpIdHashCorrect(rpIdHash []byte) bool {
	correctHash := Sha256Hash([]byte(RP.Id))
	return bytes.Equal(rpIdHash, correctHash)
}

func IsPublicKeyAlgorithmSupported(publicKeyAlgorithm any) bool {

	isPublicKeyAlgorithmSupported := false

	for i := 0; i < len(ListOfSupportedCoseAlgorithms); i += 1 {
		if publicKeyAlgorithm == ListOfSupportedCoseAlgorithms[i] {
			isPublicKeyAlgorithmSupported = true
		}
	}

	return isPublicKeyAlgorithmSupported
}

func AreFlagsValid(flags byte) bool {
	userPresentBit := flags & 1
	userVerificationBit := (flags >> 2) & 1
	backupEligibilityBit := (flags >> 3) & 1
	backupStateBit := (flags >> 4) & 1

	isValidBackup := (backupEligibilityBit == 0 && backupStateBit == 0) || (backupEligibilityBit == 1 && backupStateBit == 0) || (backupEligibilityBit == 1 && backupStateBit == 1)

	return userPresentBit == 1 && userVerificationBit == 1 && isValidBackup
}

func IsSignatureCounterValid(signatureCounterFromServer uint32, signatureCounterFromAuthenticator uint32) bool {
	if signatureCounterFromServer != 0 || signatureCounterFromAuthenticator != 0 {
		return signatureCounterFromAuthenticator > signatureCounterFromServer
	}
	return true
}

/*
Assuming that the following are true:

	Let options be a new CredentialCreationOptions structure configured to the Relying Party’s needs for the ceremony.

	options.mediation is not set to conditional.
	attestation is not cared.

	Returns true, if the operation is succesful and false the operation is unsuccesful.

	Also, please note that functionToSaveCredentialsIntoDatabase should return true if the operation is succesful and false if the operation is not succesful.
*/
func SaveCredentialsIntoDatabaseIfAuthDataIsValid(userId string, userName string, authData []byte, functionToSaveCredentialsIntoDatabase func(userId string, userName string, credentialId []byte, credentialPublicKey []byte, signatureCounter uint32) bool) bool {
	hash := authData[0:32]

	credentialIdLength := binary.BigEndian.Uint16(authData[53:55])

	var publicKeyMap map[int64]any
	credentialPublicKey := authData[55+credentialIdLength:]
	cbor.NewDecoder(bytes.NewReader(credentialPublicKey)).Decode(&publicKeyMap)

	isAuthDataValid := IsRpIdHashCorrect(hash) && AreFlagsValid(authData[32]) && credentialIdLength <= 1023 && IsPublicKeyAlgorithmSupported(publicKeyMap[3])
	if !isAuthDataValid {
		return false
	}

	credentialId := authData[55:(55 + credentialIdLength)]

	signatureCounter := binary.BigEndian.Uint32(authData[33:37])

	isOperationSuccesful := functionToSaveCredentialsIntoDatabase(userId, userName, credentialId, credentialPublicKey, signatureCounter)
	return isOperationSuccesful
}

func IsEmptyAttestationObject(attestationObject AttestationObject) bool {
	return len(attestationObject.AttStmt) == 0 && len(attestationObject.AuthData) == 0 && attestationObject.Fmt == ""
}

/*
functionToSaveCredentialsIntoDatabase should return true if the operation is succesful and false if the operation is not succesful.
The SignUp function will return true if the operation is a success, and false if the operation is a failure.
Also, it is of the caller's responsibility to ensure that the challenge gets deleted after it is used.
*/
func SignUp(originOfServer string, userId string, challengeId string, functionToGetCorrectChallenge func(id any) []byte, json map[string]any, functionToSaveCredentialsIntoDatabase func(userId string, userName string, credentialId []byte, credentialPublicKey []byte, signatureCounter uint32) bool) bool {

	publicKeyCredential := json["credential"].(map[string]any)
	userName := json["username"].(string)

	response := publicKeyCredential["response"].(map[string]any)
	clientDataJSON := response["clientDataJSON"].(string)

	correctChallenge := functionToGetCorrectChallenge(challengeId)

	if !IsClientDataJSONCorrect(originOfServer, correctChallenge, "webauthn.create", clientDataJSON) {
		return false
	}

	attestationObjectString := response["attestationObject"].(string)
	attestationObject := ParseAttestationObject(attestationObjectString)

	if IsEmptyAttestationObject(attestationObject) {
		return false
	}

	isOperationSuccesful := SaveCredentialsIntoDatabaseIfAuthDataIsValid(userId, userName, attestationObject.AuthData, functionToSaveCredentialsIntoDatabase)

	return isOperationSuccesful
}

/*
Returns a non empty string, which is the user id if the authentication process succeeds and an empty string if the authentication process fails.
It is of the caller's responsibility to ensure that the challenge gets deleted after the authentication process succeeds.
*/
func Authenticate(publicKeyCredential map[string]any, functionToGetPublicKeyAndSignatureCounter func(credential_id []byte, user_id []byte) ([]byte, uint32), functionToGetCorrectChallenge func(challenge_id any) []byte, challenge_id string, functionToUpdateSignatureCounter func(credential_id []byte, signature_counter uint32)) string {

	response := publicKeyCredential["response"].(map[string]any)
	rawId := publicKeyCredential["rawId"].(string)
	signature := response["signature"].(string)
	userHandle := response["userHandle"].(string)

	decodedRawId, err1 := base64.RawURLEncoding.DecodeString(rawId)
	decodedSignature, err2 := base64.RawURLEncoding.DecodeString(signature)
	decodedUserHandle, err3 := base64.RawURLEncoding.DecodeString(userHandle)

	if err1 != nil {
		return ""
	}

	if err2 != nil {
		return ""
	}

	if err3 != nil {
		return ""
	}

	publicKey, signatureCounterFromServer := functionToGetPublicKeyAndSignatureCounter(decodedRawId, decodedUserHandle)

	correct_challenge := functionToGetCorrectChallenge(challenge_id)

	var clientDataJSON string = response["clientDataJSON"].(string)

	isClientDataJSONCorrect := IsClientDataJSONCorrect(globals.OriginOfServer, correct_challenge, "webauthn.get", clientDataJSON)

	authData := response["authenticatorData"].(string)
	decodedAuthData, _ := base64.RawURLEncoding.DecodeString(authData)

	rpIdHash := decodedAuthData[0:32]
	decodedClientData, err4 := base64.RawURLEncoding.DecodeString(clientDataJSON)

	if err4 != nil {
		return ""
	}

	signatureCounterFromAuthenticator := binary.BigEndian.Uint32(decodedAuthData[33:37])

	isAuthenticated := IsRpIdHashCorrect(rpIdHash) && isClientDataJSONCorrect && AreFlagsValid(decodedAuthData[32]) && IsSignatureVerified(append(decodedAuthData, Sha256Hash(decodedClientData)...), decodedSignature, publicKey) && IsSignatureCounterValid(signatureCounterFromServer, signatureCounterFromAuthenticator)

	functionToUpdateSignatureCounter(decodedRawId, signatureCounterFromAuthenticator)

	if isAuthenticated {
		return string(decodedUserHandle)
	} else {
		return ""
	}
}
