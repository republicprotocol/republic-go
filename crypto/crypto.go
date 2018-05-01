package crypto

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
)

var ErrUnimplemented = errors.New("unimplemented")
var ErrInvalidSignature = errors.New("invalid signature")

type Hasher interface {
	Hash() []byte
}

type Hash32 [32]byte

func NewHash32(data []byte) Hash32 {
	hash32 := Hash32{}
	for i := 0; i < 32 && i < len(data); i++ {
		hash32[i] = data[i]
	}
	return hash32
}

func (hash32 Hash32) Hash() []byte {
	return hash32[:]
}

type Hash65 [65]byte

func NewHash65(data []byte) Hash65 {
	hash65 := Hash65{}
	for i := 0; i < 65 && i < len(data); i++ {
		hash65[i] = data[i]
	}
	return hash65
}

func (hash65 Hash65) Hash() []byte {
	return hash65[:]
}

type Signer interface {
	Sign(Hasher) ([]byte, error)
}

type Verifier interface {
	Verify(Hasher, []byte) error
}

type Encrypter interface {
	Encrypt(string, []byte) ([]byte, error)
	Decrypt(string, []byte) ([]byte, error)
}

type Crypter interface {
	Signer
	Verifier
	Encrypter
}

// WeakCrypter implements the Crypter interface. It signs hashes by
// immediately returning the hash, and verifies all signatures as correct. It
// returns the plain text immediately when encrypting, and returns the cipher
// text immediately when decrypting. A WeakCrypter must only be used during
// testing.
type WeakCrypter struct{}

// NewWeakCrypter returns a new WeakCrypter. A WeakCrypter must only be used
// during testing.
func NewWeakCrypter() WeakCrypter {
	return WeakCrypter{}
}

// Sign implements the Crypter interface. It returns the hash of the Hasher and
// no error.
func (crypter *WeakCrypter) Sign(hasher Hasher) ([]byte, error) {
	return hasher.Hash(), nil
}

// Verify implements the Crypter interface. It returns no error.
func (crypter *WeakCrypter) Verify(hasher Hasher, signature []byte) error {
	return nil
}

// Encrypt implements the Crypter interface. It returns the plain text and no
// error.
func (crypter *WeakCrypter) Encrypt(addr string, plainText []byte) ([]byte, error) {
	return plainText, nil
}

// Decrypt implements the Crypter interface. It returns the cipher text and no
// error.
func (crypter *WeakCrypter) Decrypt(addr string, cipherText []byte) ([]byte, error) {
	return cipherText, nil
}

// ErrCrypter implements the Crypter interface and immediately returns an error
// from all method. An ErrCrypter must only be used during testing.
type ErrCrypter struct{}

// NewErrCrypter returns a new ErrCrypter. An ErrCrypter must only be used
// during testing.
func NewErrCrypter() ErrCrypter {
	return ErrCrypter{}
}

// Signer implements the Crypter interface. It returns an empty byte slice and
// ErrUnimplemented.
func (crypter *ErrCrypter) Signer(hasher Hasher) ([]byte, error) {
	return []byte{}, ErrUnimplemented
}

// Verify implements the Crypter interface. It returns an empty byte slice and
// ErrUnimplemented.
func (crypter *ErrCrypter) Verify(hasher Hasher, signature []byte) error {
	return ErrUnimplemented
}

// Encrypt implements the Crypter interface. It returns an empty byte slice and
// ErrUnimplemented.
func (crypter *ErrCrypter) Encrypt(addr string, plainText []byte) ([]byte, error) {
	return []byte{}, ErrUnimplemented
}

// Decrypt implements the Crypter interface. It returns an empty byte slice and
// ErrUnimplemented.
func (crypter *ErrCrypter) Decrypt(addr string, cipherText []byte) ([]byte, error) {
	return []byte{}, ErrUnimplemented
}

func unmarshalStringFromMap(m map[string]json.RawMessage, k string) (string, error) {
	if val, ok := m[k]; ok {
		str := ""
		if err := json.Unmarshal(val, &str); err != nil {
			return "", err
		}
		return str, nil
	}
	return "", fmt.Errorf("%s is nil", k)
}

func unmarshalIntFromMap(m map[string]json.RawMessage, k string) (int, error) {
	if val, ok := m[k]; ok {
		i := 0
		if err := json.Unmarshal(val, &i); err != nil {
			return 0, err
		}
		return i, nil
	}
	return 0, fmt.Errorf("%s is nil", k)
}

func unmarshalBigIntFromMap(m map[string]json.RawMessage, k string) (*big.Int, error) {
	if val, ok := m[k]; ok {
		bytes := []byte{}
		if err := json.Unmarshal(val, &bytes); err != nil {
			return nil, err
		}
		return big.NewInt(0).SetBytes(bytes), nil
	}
	return nil, fmt.Errorf("%s is nil", k)
}

func unmarshalBigIntsFromMap(m map[string]json.RawMessage, k string) ([]*big.Int, error) {
	bigInts := []*big.Int{}
	if val, ok := m[k]; ok {
		vals := []json.RawMessage{}
		if err := json.Unmarshal(val, &vals); err != nil {
			return bigInts, err
		}
		for _, val := range vals {
			bytes := []byte{}
			if err := json.Unmarshal(val, &bytes); err != nil {
				return bigInts, err
			}
			bigInts = append(bigInts, big.NewInt(0).SetBytes(bytes))
		}
		return bigInts, nil
	}
	return bigInts, fmt.Errorf("%s is nil", k)
}
