package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
)

// RsaKeyPair contains a RSA key pair
type RsaKeyPair struct {
	*rsa.PublicKey
	*rsa.PrivateKey
}

// NewRsaKeyPair generates a new RSA key pair.
// It returns a randomly generated KeyPair, or an error.
// It precomputes some useful values to save time on decryption.
func NewRsaKeyPair() (RsaKeyPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	privateKey.Precompute()
	if err != nil {
		return RsaKeyPair{}, err
	}
	publicKey := privateKey.PublicKey
	return RsaKeyPair{
		PublicKey:  &publicKey,
		PrivateKey: privateKey,
	}, nil
}

// NewRsaKeyPairFromPrivateKey a new RSA key pair using a given private key. It
// does not validate that this private key was generated correctly.
// It precomputes some useful values to save time on decryption.
func NewRsaKeyPairFromPrivateKey(privKey *rsa.PrivateKey) RsaKeyPair {
	privKey.Precompute()
	return RsaKeyPair{
		PublicKey:  &privKey.PublicKey,
		PrivateKey: privKey,
	}
}

// Encrypt encrypts a message with the provided RSA public key.
// It returns the cipher text, or an error.
func Encrypt(pubKey *rsa.PublicKey, msg []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, pubKey, msg)
}

// Decrypt decrypts the given cipher text with the provided RSA private key.
// It returns the decrypted message, or an error.
func Decrypt(privKey *rsa.PrivateKey, cipherText []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, privKey, cipherText)
}

// PublicKeyToBytes converts Public Key to a byte array.
// It returns the byte array, or an error.
func PublicKeyToBytes(publicKey *rsa.PublicKey) ([]byte, error) {
	return json.Marshal(publicKey)
}

// BytesToPublicKey converts a byte array to a RSA Public Key
// It returns pointer to a RSA public key
func BytesToPublicKey(publicKey []byte) (*rsa.PublicKey, error) {
	var pubKey rsa.PublicKey
	err := json.Unmarshal(publicKey, &pubKey)
	if err != nil {
		return &rsa.PublicKey{}, err
	}
	return &pubKey, nil
}
