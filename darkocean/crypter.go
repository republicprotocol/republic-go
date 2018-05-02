package darkocean

import "github.com/republicprotocol/republic-go/crypto"

// Crypter is an implementation of the crypto.Crypter interface. In addition to
// standard signature verification, the Crypter uses a DarkOcean to verify that
// the signatory is correctly registered to the network. It also uses the
// DarkOcean acquire the necessary rsa.PublicKeys for encryption.
type Crypter struct {
	keystore crypto.Keystore
	ocean    DarkOcean
}

// NewCrypter returns a new Crypter that uses a crypto.Keystore to identify
// itself when signing and decrypting messages. It uses a DarkOcean to identify
// others when verifying and encrypting messages.
func NewCrypter(keystore crypto.Keystore, ocean DarkOcean) Crypter {
	return Crypter{
		keystore: keystore,
		ocean:    ocean,
	}
}

func (crypter *Crypter) Sign(hasher crypto.Hasher) ([]byte, error) {
	panic("unimplemented")
}

func (crypter *Crypter) Verify(hasher crypto.Hasher, signature []byte) error {
	panic("unimplemented")
}

func (crypter *Crypter) Encrypt(addr string, plainText []byte) ([]byte, error) {
	panic("unimplemented")
}

func (crypter *Crypter) Decrypt(cipherText []byte) ([]byte, error) {
	panic("unimplemented")
}
