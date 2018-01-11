package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
    b := make([]byte, n)
    _, err := rand.Read(b)
    // Note that err == nil only if we read len(b) bytes.
    if err != nil {
        return nil, err
    }

    return b, nil
}

// Generates a secret key pair, to be used in atomic swaps
// returns (Secret, SecretLock, error) 
// It will return an error if something goes wrong during 
// random string generation.
func GenerateSecretPair() (string, string, error) {
	token, err := GenerateRandomBytes(32)
	if err != nil {
		return "", "", err
	}
	h := sha256.New()
	h.Write(token)
	return hex.EncodeToString(token), hex.EncodeToString(h.Sum(nil)[:32]), nil
}