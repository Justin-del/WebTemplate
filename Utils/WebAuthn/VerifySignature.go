package webauthn

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rsa"
	"math/big"

	"github.com/fxamacker/cbor/v2"
)

func isES256SignatureVerified(data []byte, signature []byte, public_key_map map[int64]any) bool {

	publicKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     big.NewInt(0).SetBytes(public_key_map[-2].([]byte)),
		Y:     big.NewInt(0).SetBytes(public_key_map[-3].([]byte)),
	}

	return ecdsa.VerifyASN1(publicKey, Sha256Hash(data), signature)
}

func isRS256SignatureVerified(data []byte, signature []byte, public_key_map map[int64]any) bool {

	publicKey := &rsa.PublicKey{
		N: big.NewInt(0).SetBytes(public_key_map[-1].([]byte)),
		E: int(big.NewInt(0).SetBytes(public_key_map[-2].([]byte)).Int64()),
	}

	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, Sha256Hash(data), signature) == nil
}

func IsSignatureVerified(data []byte, signature []byte, public_key []byte) bool {
	var public_key_map map[int64]any
	err := cbor.NewDecoder(bytes.NewReader(public_key)).Decode(&public_key_map)

	if err != nil {
		return false
	}

	return (public_key_map[3].(int64) == -7 && public_key_map[-1].(uint64) == 1 && isES256SignatureVerified(data, signature, public_key_map)) || (public_key_map[3].(int64) == -8 && public_key_map[-1].(uint64) == 6 && ed25519.Verify(public_key, data, signature)) || isRS256SignatureVerified(data, signature, public_key_map)
}
