package crypto

import "crypto/rsa"

type RsaKeyPair struct {
	*rsa.PublicKey
	*rsa.PrivateKey
}

func NewRsaKeyPair() RsaKeyPair {
	privateKey, err = rsa.GenerateKey()

}

func KeyPairFromPublicKey(pubKey *rsa.PublicKey) RsaKeyPair {
	return RsaKeyPair{
		PublicKey: pubKey,
	}
}

func (kp *RsaKeyPair) Encrypt(msg []byte) ([]byte, error) {

}

func (kp *RsaKeyPair) Decrypt(cipherText []byte) ([]byte, error) {

}
