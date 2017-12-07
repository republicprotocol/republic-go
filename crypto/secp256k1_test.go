package crypto

import (
	"encoding/base64"
	"testing"
)

func TestSECP256K1SignAndVerify(t *testing.T) {
	// Generate a key pair and payload.
	payload := []byte("this-is-super-secret")
	id, err := NewSECP256K1()
	if err != nil {
		t.Fatalf("failed to generate SECP256K1: %v", err)
	}

	// Sign the payload.
	hash, signedHash, err := id.Sign(payload)
	if err != nil {
		t.Fatalf("failed to signed payload using SECP256K1: %v", err)
	}

	// Verify the signature.
	verified, err := id.Verify(hash, signedHash)
	if err != nil {
		t.Fatalf("failed to verify signed hash using SECP256K1: %v", err)
	}
	if !verified {
		t.Fatalf("failed to verify signed hash using SECP256K1")
	}
}

func TestSECP256K1GenerateAndNew(t *testing.T) {
	// Generate a key pair.
	id, err := NewSECP256K1()
	if err != nil {
		t.Fatalf("failed to generate SECP256K1: %v", err)
	}
	privateKey := id.PrivateKey()
	publicKey := id.PublicKey()

	// Create a new key pair from the generated key pair.
	id, err = NewSECP256K1FromPrivateKey(privateKey)
	if err != nil {
		t.Fatalf("failed to regenerate SECP256K1: %v", err)
	}
	if string(id.PrivateKey()) != string(privateKey) {
		t.Fatalf("expected matching private keys")
	}
	if string(id.PublicKey()) != string(publicKey) {
		t.Fatalf("expected matching public keys")
	}
}
func TestSECP256K1Base64(t *testing.T) {
	// Generate a key pair.
	id, err := NewSECP256K1()
	if err != nil {
		t.Fatalf("failed to generate SECP256K1: %v", err)
	}
	privateKey := id.PrivateKey()
	publicKey := id.PublicKey()

	// Encode and decode the key pair.
	encodedPrivateKey := base64.StdEncoding.EncodeToString(privateKey)
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(encodedPrivateKey)
	if err != nil {
		t.Fatalf("failed to decode SECP256K1: %v", err)
	}

	// Create a new key pair from the decoded data.
	id, err = NewSECP256K1FromPrivateKey(decodedPrivateKey)
	if err != nil {
		t.Fatalf("failed to regenerate SECP256K1: %v", err)
	}
	if string(id.PrivateKey()) != string(privateKey) {
		t.Fatalf("expected matching private keys")
	}
	if string(id.PublicKey()) != string(publicKey) {
		t.Fatalf("expected matching public keys")
	}
}

func TestSECP256K1JSON(t *testing.T) {
	// Generate a key pair.
	id, err := NewSECP256K1()
	if err != nil {
		t.Fatalf("failed to generate SECP256K1: %v", err)
	}
	privateKey := id.PrivateKey()
	publicKey := id.PublicKey()

	// Marshal and unmarshal the key pair.
	blob, err := id.MarshalJSON()
	if err != nil {
		t.Fatalf("failed to marshal SECP256K1: %v", err)
	}
	if err := id.UnmarshalJSON(blob); err != nil {
		t.Fatalf("failed to unmarshal SECP256K1: %v", err)
	}
	if string(id.PrivateKey()) != string(privateKey) {
		t.Fatalf("expected matching private keys")
	}
	if string(id.PublicKey()) != string(publicKey) {
		t.Fatalf("expected matching public keys")
	}
}
