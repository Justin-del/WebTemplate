package webauthn

import (
	"WebTemplate/globals"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"net/http"

	"github.com/fxamacker/cbor/v2"
)

func GetRegistrationData() RegistrationData {
	response, err := http.Get(globals.OriginOfServer + "/SignUp/RegistrationData")
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	var registrationData RegistrationData

	err = json.NewDecoder(response.Body).Decode(&registrationData)

	if err != nil {
		panic(err)
	}

	return registrationData

}

/*
t is short for type,
challenge is not base64 url encoded.
*/
func CreateClientData(challenge []byte, t string, origin string) map[string]any {
	var clientData map[string]any = make(map[string]any)

	clientData["challenge"] = base64.RawURLEncoding.EncodeToString(challenge)
	clientData["origin"] = origin
	clientData["type"] = t

	return clientData
}

func CreateMockPublicKey(publicKeyAlgorithm int64) []byte {
	public_key := make(map[int64]any)

	public_key[1] = 2
	public_key[3] = publicKeyAlgorithm
	public_key[-1] = 1
	public_key[-2] = make([]byte, 32)
	public_key[-3] = make([]byte, 32)

	encodedPublicKey, _ := cbor.Marshal(public_key)
	return encodedPublicKey
}

func CreateMockAuthData(publicKey []byte, relyingPartyId string, userPresent bool, userVerified bool, backupEligibility bool, backupState bool, credentialId string) []byte {
	mockAuthData := make([]byte, 10000)
	copy(mockAuthData[0:32], Sha256Hash([]byte(relyingPartyId)))

	var flags byte = 0

	if userPresent {
		//Set the bit for User Present
		flags = SetNthBitTo1(flags, 0)
	}

	if userVerified {
		//Set the bit for User Verified
		flags = SetNthBitTo1(flags, 2)
	}

	if backupEligibility {
		//Set the bit for backup eligibility
		flags = SetNthBitTo1(flags, 3)
	}

	if backupState {
		//Set the bit for backup state
		flags = SetNthBitTo1(flags, 4)
	}

	mockAuthData[32] = flags

	lengthOfIdToBeUsedForTesting := uint16(len(credentialId))
	binary.BigEndian.PutUint16(mockAuthData[53:55], lengthOfIdToBeUsedForTesting)
	copy(mockAuthData[55:55+lengthOfIdToBeUsedForTesting], []byte(credentialId))
	copy(mockAuthData[55+lengthOfIdToBeUsedForTesting:], publicKey)
	return mockAuthData
}

func CreateMockAttestationObject(publicKey []byte, relyingPartyId string, userPresent bool, userVerified bool, backupEligibility bool, backupState bool, credentialId string) string {

	attestationObject := AttestationObject{
		AttStmt:  []byte{160},
		AuthData: CreateMockAuthData(publicKey, relyingPartyId, userPresent, userVerified, backupEligibility, backupState, credentialId),
		Fmt:      "none",
	}

	encodedAttestationObject, _ := cbor.Marshal(attestationObject)
	return base64.RawURLEncoding.EncodeToString(encodedAttestationObject)
}

/*
The 0th bit is the least signficant bit.
*/
func SetNthBitTo1(b byte, n int) byte {
	if n < 0 || n > 7 {
		panic("bit index out of range")
	}
	b |= (1 << n)
	return b
}

func CreateMockPublicKeyCredential(clientData map[string]any, publicKey []byte, relyingPartyId string, userPresent bool, userVerified bool, backupEligibility bool, backupState bool, credentialId string) string {
	publicKeyCredential := make(map[string]any)

	publicKeyCredential["authenticatorAttachment"] = "platform"
	publicKeyCredential["clientExtensionResults"] = make(map[string]any)
	publicKeyCredential["id"] = credentialId
	publicKeyCredential["rawId"] = credentialId

	response := make(map[string]any)

	clientDataJSON, _ := json.Marshal(clientData)

	response["clientDataJSON"] = base64.RawURLEncoding.EncodeToString(clientDataJSON)

	response["attestationObject"] = CreateMockAttestationObject(publicKey, relyingPartyId, userPresent, userVerified, backupEligibility, backupState, credentialId)
	
	publicKeyCredential["response"] = response

	jsonString, _ := json.Marshal(publicKeyCredential)

	return string(jsonString)
}
