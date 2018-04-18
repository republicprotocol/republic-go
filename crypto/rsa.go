package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
)

type RsaKeyPair struct {
	*rsa.PublicKey
	*rsa.PrivateKey
}

func NewRsaKeyPair() (RsaKeyPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return RsaKeyPair{}, err
	}
	publicKey := privateKey.PublicKey
	return RsaKeyPair{
		PublicKey:  &publicKey,
		PrivateKey: privateKey,
	}, nil
}

func Encrypt(pubKey *rsa.PublicKey, msg []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, pubKey, msg)
}

func Decrypt(privKey *rsa.PrivateKey, cipherText []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, privKey, cipherText)
}

func PublicKeyToBytes(publicKey *rsa.PublicKey) ([]byte, error) {
	return json.Marshal(publicKey)
}

func BytesToPublicKey(publicKey []byte) (*rsa.PublicKey, error) {
	var pubKey rsa.PublicKey
	err := json.Unmarshal(publicKey, &pubKey)
	if err != nil {
		return &rsa.PublicKey{}, err
	}
	return &pubKey, nil
}
