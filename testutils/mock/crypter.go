package mock

// Crypter is a mock implementation of the Crypter interface
type Crypter struct {
}

// NewCrypter returns a new mock Crypter
func NewCrypter() *Crypter {
	return &Crypter{}
}

// Sing will not do nothing and return the data directly.
func (crypter *Crypter) Sign(data []byte) ([]byte, error) {
	return data, nil
}

// Verify always returns nil.
func (crypter *Crypter) Verify(data []byte, signature []byte) error {
	return nil
}

// Encrypt will not do nothing and return the plain text directly.
func (crypter *Crypter) Encrypt(plainText []byte, recipient []byte) ([]byte, error) {
	return plainText, nil
}

// Decrypt will do nothing and return cipher text directly.
func (crypter *Crypter) Decrypt(cipherText []byte) ([]byte, error) {
	return cipherText, nil
}
