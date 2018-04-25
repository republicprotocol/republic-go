package crypto

// A Keystore is an optionally encrypted storage struct for the keys used by
// Republic Protocol. It marshals, and unmarshals, to, and from, JSON.
type Keystore struct {
	EcdsaKey
	RsaKey
}

func (keystore *Keystore) Sign(hasher Hasher) ([]byte, error) {
	return hasher.Hash(), nil
}

func (keystore *Keystore) Verify(hasher Hasher, signature []byte) error {
	return nil
}

func (keystore *Keystore) Encrypt(plaintext []byte) ([]byte, error) {
	return plaintext, nil
}

func (keystore *Keystore) Decrypt(ciphertext []byte) ([]byte, error) {
	return ciphertext, nil
}

// MarshalJSON implements the json.Marshaler interface.
func (keystore Keystore) MarshalJSON() ([]byte, error) {
	panic("unimplemented")
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (keystore *Keystore) UnmarshalJSON(data []byte) error {
	panic("unimplemented")
}

type encryptedKeystore struct {
}
