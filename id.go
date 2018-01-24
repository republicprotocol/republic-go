package identity

// IDLength is the number of bytes in an ID.
const IDLength = 20

// An ID is a slice of 20 bytes that can be converted into an Address.
// It must always be example 20 bytes.
type ID []byte

// NewID generates a new ID by generating a random KeyPair. It returns the ID,
// and the KeyPair, or an error. It is most commonly used for testing.
func NewID() (ID, KeyPair, error) {
	keyPair, err := NewKeyPair()
	if err != nil {
		return nil, keyPair, err
	}
	return keyPair.ID(), keyPair, nil
}

// String returns the ID as a string.
func (id ID) String() string {
	return string(id)
}