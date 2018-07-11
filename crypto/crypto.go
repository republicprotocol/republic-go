package crypto

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	ethCrypto "github.com/ethereum/go-ethereum/crypto"
)

// ErrNilData is returned when a Verifier encounters a nil, or empty, data.
// Signing nil data is not in itself erroneous but it is rarely a reasonable
// action to take.
var ErrNilData = errors.New("nil signature")

// ErrNilSignature is returned when a Verifier encounters a nil, or empty,
// signature.
var ErrNilSignature = errors.New("nil signature")

// ErrInvalidSignature is returned when a signature is invalid, or the
// recovered signatory does not match the required signatory.
var ErrInvalidSignature = errors.New("invalid signature")

// A Signer can consume bytes and produce a unique signature for those bytes.
// This is usually done using the private key of an asymmetrical cryptography
// scheme.
type Signer interface {
	Sign(data []byte) ([]byte, error)
}

// A Verifier can consume bytes and a signature for those bytes, and extract
// signatory.
type Verifier interface {
	Verify(data []byte, signature []byte) error
}

// An Encrypter can consume a plain text, and produce a cipher text for that
// can only be decrypted by a specific recipient.
type Encrypter interface {
	Encrypt(recipient string, plainText []byte) ([]byte, error)
}

// A Decrypter can consume a cipher text and produce a plain text.
type Decrypter interface {
	Decrypt(cipherText []byte) ([]byte, error)
}

// Keccak256 of all bytes concatenated together.
func Keccak256(data ...[]byte) []byte {
	return ethCrypto.Keccak256(data...)
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
