package testutils

type MockCrypter struct {
}

func NewMockCrypter() *MockCrypter {
	return &MockCrypter{}
}

func (crypter *MockCrypter) Sign(data []byte) ([]byte, error) {
	return data, nil
}

func (crypter *MockCrypter) Verify(data []byte, signature []byte) error {
	return nil
}

func (crypter *MockCrypter) Encrypt(addr string, plainText []byte) ([]byte, error) {
	return plainText, nil
}

func (crypter *MockCrypter) Decrypt(cipherText []byte) ([]byte, error) {
	return cipherText, nil
}
