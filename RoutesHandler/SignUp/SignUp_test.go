package SignUp

import (
	webauthn "WebTemplate/Utils/WebAuthn"
	"WebTemplate/globals"
	"crypto/rand"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/google/uuid"
)

var supportedPublicKey = webauthn.CreateMockPublicKey(-7)
var unsupportedPublicKey = webauthn.CreateMockPublicKey(123456789)

func TestCannotSignUpWithAReplayAttack(t *testing.T) {
	registrationData := webauthn.GetRegistrationData()

	clientData := webauthn.CreateClientData(registrationData.Challenge.Challenge, "webauthn.create", globals.OriginOfServer)

	userId := uuid.New().String()
	credentialId := uuid.New().String()
	postUrl := globals.OriginOfServer + "/SignUp/" + strconv.Itoa(registrationData.Challenge.Id) + "/" + userId

	// First request
	resp, err := http.Post(postUrl, "application/json", strings.NewReader(webauthn.CreateMockPublicKeyCredential(clientData, supportedPublicKey, registrationData.RP.Id, true, true, false, false, credentialId)))
	if err != nil {
		t.Fatalf("Failed to make first POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200 for the first request, got %d", resp.StatusCode)
	}

	// Replay the same request
	respReplay, err := http.Post(postUrl, "application/json", strings.NewReader(webauthn.CreateMockPublicKeyCredential(clientData, supportedPublicKey, webauthn.RP.Id, true, true, false, false, credentialId)))
	if err != nil {
		t.Fatalf("Failed to replay POST request: %v", err)
	}
	defer respReplay.Body.Close()

	if respReplay.StatusCode != 400 {
		t.Fatalf("Expected status code 400, got %d", respReplay.StatusCode)
	}
}

func TestCannotSignUpIfClientDataIsIncorrect(t *testing.T) {
	for i := range 3 {
		registrationData := webauthn.GetRegistrationData()
		var clientData map[string]any

		if i == 0 {
			//Incorrect challenge
			clientData = webauthn.CreateClientData(append(registrationData.Challenge.Challenge,[]byte{0}...), "webauthn.create", globals.OriginOfServer)
		} else if i == 1 {
			//Incorrect type
			clientData = webauthn.CreateClientData(registrationData.Challenge.Challenge, "webauthn.bola", globals.OriginOfServer)
		} else if i == 2 {
			//Incorrect origin
			clientData = webauthn.CreateClientData(registrationData.Challenge.Challenge, "webauthn.create", globals.OriginOfServer+"hehe")
		}

		userId := uuid.New().String()
		credentialId := uuid.New().String()
		postUrl := globals.OriginOfServer + "/SignUp/" + strconv.Itoa(registrationData.Challenge.Id) + "/" + userId

		resp, err := http.Post(postUrl, "application/json", strings.NewReader(webauthn.CreateMockPublicKeyCredential(clientData, supportedPublicKey, webauthn.RP.Id, true, true, false, false, credentialId)))

		if err != nil {
			t.Fatalf("Failed to make POST request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 400 {
			t.Fatalf("Expected status code 400, got %d", resp.StatusCode)
		}
	}
}

func TestCannotSignUpIfPublicKeyCredentialIsIncorrect(t *testing.T) {
	var clientDatas [](map[string]any) = [](map[string]any){}
	var registrationDatas []webauthn.RegistrationData = []webauthn.RegistrationData{}

	//The number should corresponds to the number of incorrect public key credentials.
	//The for loop is necessary because a challenge can only be used once by each request so there should be a different challenge for each request.
	for range 5 {
		registrationData := webauthn.GetRegistrationData()
		registrationDatas = append(registrationDatas, registrationData)

		clientData := webauthn.CreateClientData(registrationData.Challenge.Challenge, "webauthn.create", globals.OriginOfServer)
		clientDatas = append(clientDatas, clientData)
	}

	var incorrectPublicKeyCredentials []string = []string{
		webauthn.CreateMockPublicKeyCredential(clientDatas[0], supportedPublicKey, registrationDatas[0].RP.Id+"?", true, true, false, false, uuid.New().String()), /*Test for incorrect relying party hash*/
		webauthn.CreateMockPublicKeyCredential(clientDatas[1], unsupportedPublicKey, registrationDatas[1].RP.Id, true, true, false, false, uuid.New().String()),   /*test for unsupported public key algorithm*/
		webauthn.CreateMockPublicKeyCredential(clientDatas[2], supportedPublicKey, registrationDatas[2].RP.Id, false, true, false, false, uuid.New().String()),    /* test for user present bit is false */
		webauthn.CreateMockPublicKeyCredential(clientDatas[3], supportedPublicKey, registrationDatas[3].RP.Id, true, false, false, false, uuid.New().String()),    /* test for user verification bit is false */
		webauthn.CreateMockPublicKeyCredential(clientDatas[4], supportedPublicKey, registrationDatas[4].RP.Id, true, true, false, true, uuid.New().String()),      /* test for backup state is set but backup eligibility is not set.*/
	}

	for i, credential := range incorrectPublicKeyCredentials {
		userId := uuid.New().String()
		postUrl := globals.OriginOfServer + "/SignUp/" + strconv.Itoa(registrationDatas[i].Challenge.Id) + "/" + userId
		resp, err := http.Post(postUrl, "application/json", strings.NewReader(credential))
		if err != nil {
			t.Fatalf("Failed to make POST request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 400 {
			t.Fatalf("Expected status code 400, got %d", resp.StatusCode)
		}
	}

}

func TestForValidBackupStateAndBackupElibility(t *testing.T) {
	var clientDatas [](map[string]any) = [](map[string]any){}
	var registrationDatas []webauthn.RegistrationData = []webauthn.RegistrationData{}

	//The number should corresponds to the number of correct public key credentials.
	//The for loop is necessary because a challenge can only be used once by each request so there should be a different challenge for each request.
	for range 3 {
		registrationData := webauthn.GetRegistrationData()
		registrationDatas = append(registrationDatas, registrationData)

		clientData := webauthn.CreateClientData(registrationData.Challenge.Challenge, "webauthn.create", globals.OriginOfServer)
		clientDatas = append(clientDatas, clientData)
	}

	var publicKeyCredentials []string = []string{
		webauthn.CreateMockPublicKeyCredential(clientDatas[0], supportedPublicKey, registrationDatas[0].RP.Id, true, true, false, false, uuid.New().String()), //backupEligibility = false & backupState = false
		webauthn.CreateMockPublicKeyCredential(clientDatas[1], supportedPublicKey, registrationDatas[1].RP.Id, true, true, true, false, uuid.New().String()),  //backupEligibity = true & backupState = false
		webauthn.CreateMockPublicKeyCredential(clientDatas[2], supportedPublicKey, registrationDatas[2].RP.Id, true, true, true, true, uuid.New().String()),   //backupEligibility = true & backupState = true
	}

	for i, credential := range publicKeyCredentials {
		userId := uuid.New().String()
		postUrl := globals.OriginOfServer + "/SignUp/" + strconv.Itoa(registrationDatas[i].Challenge.Id) + "/" + userId
		resp, err := http.Post(postUrl, "application/json", strings.NewReader(credential))
		if err != nil {
			t.Fatalf("Failed to make POST request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
		}
	}
}

func TestCannotSignUpIfCredentialIdIsGreaterThan1023Bytes(t *testing.T) {
	registrationData := webauthn.GetRegistrationData()
	clientData := webauthn.CreateClientData(registrationData.Challenge.Challenge, "webauthn.create", globals.OriginOfServer)

	// Generate a random credential ID of 1024 bytes
	credentialIdBytes := make([]byte, 1024)
	_, err := rand.Read(credentialIdBytes)
	if err != nil {
		t.Fatalf("Failed to generate random credential ID: %v", err)
	}
	credentialId := string(credentialIdBytes)

	userId := uuid.New().String()
	postUrl := globals.OriginOfServer + "/SignUp/" + strconv.Itoa(registrationData.Challenge.Id) + "/" + userId

	resp, err := http.Post(postUrl, "application/json", strings.NewReader(webauthn.CreateMockPublicKeyCredential(clientData, supportedPublicKey, registrationData.RP.Id, true, true, false, false, credentialId)))
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 400 {
		t.Fatalf("Expected status code 400, got %d", resp.StatusCode)
	}
}

func TestCannotSignUpIfCredentialIdIsAlreadyRegisteredForAUser(t *testing.T) {
	registrationData := webauthn.GetRegistrationData()

	clientData := webauthn.CreateClientData(registrationData.Challenge.Challenge, "webauthn.create", globals.OriginOfServer)

	userId := uuid.New().String()
	postUrl := globals.OriginOfServer + "/SignUp/" + strconv.Itoa(registrationData.Challenge.Id) + "/" + userId

	credentialId := uuid.New().String()

	// First request
	resp, err := http.Post(postUrl, "application/json", strings.NewReader(webauthn.CreateMockPublicKeyCredential(clientData, supportedPublicKey, registrationData.RP.Id, true, true, false, false, credentialId)))
	if err != nil {
		t.Fatalf("Failed to make first POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200 for the first request, got %d", resp.StatusCode)
	}

	//Sign up again with the same credentialId.
	registrationData = webauthn.GetRegistrationData()

	clientData = webauthn.CreateClientData(registrationData.Challenge.Challenge, "webauthn.create", globals.OriginOfServer)

	userId = uuid.New().String()
	postUrl = globals.OriginOfServer + "/SignUp/" + strconv.Itoa(registrationData.Challenge.Id) + "/" + userId

	// Sign up again with the same credentialId.
	resp, err = http.Post(postUrl, "application/json", strings.NewReader(webauthn.CreateMockPublicKeyCredential(clientData, supportedPublicKey, registrationData.RP.Id, true, true, false, false, credentialId)))
	if err != nil {
		t.Fatalf("Failed to make first POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 400 {
		t.Fatalf("Expected status code 400 , got %d", resp.StatusCode)
	}
}
