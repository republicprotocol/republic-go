package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/binary"
	"encoding/json"
	"math/big"
)

// RsaKey for encrypting and decrypting sensitive data that must be transported
// between actors in the network.
type RsaKey struct {
	*rsa.PrivateKey
}

// RandomRsaKey using 2048 bits, with precomputed values for improved
// performance.
func RandomRsaKey() (RsaKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	privateKey.Precompute()
	if err != nil {
		return RsaKey{}, err
	}
	return RsaKey{
		PrivateKey: privateKey,
	}, nil
}

// NewRsaKey returns an RsaKey from an existing private key. It does not verify
// that the private key was generated correctly. It precomputes values for
// improved performance
func NewRsaKey(privateKey *rsa.PrivateKey) RsaKey {
	privateKey.Precompute()
	return RsaKey{
		PrivateKey: privateKey,
	}
}

// Encrypt a plain text message and return the cipher text.
func (key *RsaKey) Encrypt(plainText []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, &key.PublicKey, plainText)
}

// Decrypt a cipher text and return the plain text message.
func (key *RsaKey) Decrypt(cipherText []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, key.PrivateKey, cipherText)
}

// MarshalJSON implements the json.Marshaler interface. The RsaKey is formatted
// according to the Republic Protocol Keystore specification.
func (key RsaKey) MarshalJSON() ([]byte, error) {
	jsonKey := map[string]interface{}{}
	// Private key
	jsonKey["d"] = key.D.Bytes()
	jsonKey["primes"] = [][]byte{}
	for _, p := range key.Primes {
		jsonKey["primes"] = append(jsonKey["primes"].([][]byte), p.Bytes())
	}
	// Public key
	jsonKey["n"] = key.N.Bytes()
	jsonKey["e"] = key.E
	return json.Marshal(jsonKey)
}

// UnmarshalJSON implements the json.Unmarshaler interface. An RsaKey is
// created from data that is assumed to be compliant with the Republic Protocol
// Keystore specification. The RsaKey will be precomputed.
func (key *RsaKey) UnmarshalJSON(data []byte) error {
	jsonKey := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &jsonKey); err != nil {
		return err
	}

	var err error

	// Private key
	key.PrivateKey = new(rsa.PrivateKey)
	key.PrivateKey.D, err = unmarshalBigIntFromMap(jsonKey, "d")
	if err != nil {
		return err
	}
	key.PrivateKey.Primes, err = unmarshalBigIntsFromMap(jsonKey, "primes")
	if err != nil {
		return err
	}

	// Public key
	key.PrivateKey.PublicKey = rsa.PublicKey{}
	key.PrivateKey.PublicKey.N, err = unmarshalBigIntFromMap(jsonKey, "n")
	if err != nil {
		return err
	}
	key.PrivateKey.PublicKey.E, err = unmarshalIntFromMap(jsonKey, "e")
	if err != nil {
		return err
	}

	key.Precompute()
	return nil
}

// BytesFromRsaPublicKey by using the Republic Protocol Keystore specification
// for binary marshaling.
func BytesFromRsaPublicKey(publicKey *rsa.PublicKey) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, int64(publicKey.E))
	binary.Write(buf, binary.BigEndian, publicKey.N.Bytes())
	return buf.Bytes()
}

// RsaPublicKeyFromBytes decodes a slice of bytes into an rsa.PublicKey. It
// assumes that the bytes slice is compliant with the Republic Protocol
// Keystore specification.
func RsaPublicKeyFromBytes(data []byte) rsa.PublicKey {
	reader := bytes.NewReader(data)
	e := int64(0)
	binary.Read(reader, binary.BigEndian, &e)
	n := make([]byte, reader.Len())
	binary.Read(reader, binary.BigEndian, n)
	return rsa.PublicKey{
		E: int(e),
		N: big.NewInt(0).SetBytes(n),
	}
}
