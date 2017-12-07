package crypto

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
)

// SECP256K1 is an implementation of the Identity interface. It is the elliptic
// key pair scheme used by Ethereum.
type SECP256K1 struct {
	ECDSA    *ecdsa.PrivateKey
	ECDSAPub *ecdsa.PublicKey
}

// NewSECP256K1 returns a new SECP256K1 containing a randomly generated public
// and private key pair, or an error,
func NewSECP256K1() (SECP256K1, error) {
	ecdsa, err := crypto.GenerateKey()
	if err != nil {
		return SECP256K1{
			ECDSA:    nil,
			ECDSAPub: nil,
		}, err
	}
	return SECP256K1{
		ECDSA:    ecdsa,
		ECDSAPub: &ecdsa.PublicKey,
	}, nil
}

// NewSECP256K1FromPrivateKey returns a new SECP256K1 from a given private key,
// or an error.
func NewSECP256K1FromPrivateKey(privateKey []byte) (SECP256K1, error) {
	ecdsa, err := crypto.ToECDSA(privateKey)
	if err != nil {
		return SECP256K1{}, err
	}
	return SECP256K1{
		ECDSA:    ecdsa,
		ECDSAPub: &ecdsa.PublicKey,
	}, nil
}

// NewSECP256K1FromJSONFile reads a JSON file and marshals it into a SECP256K1
// identity. Returns the unmarshaled SECP256K1, or an error.
func NewSECP256K1FromJSONFile(filename string) (SECP256K1, error) {
	// Open the file.
	file, err := os.Open(filename)
	if err != nil {
		return SECP256K1{}, err
	}
	defer file.Close()

	// Unmarshal the file into a SECP256K1 key pair.
	id := SECP256K1{}
	if err := json.NewDecoder(file).Decode(&id); err != nil {
		return id, err
	}
	return id, nil
}

// PublicAddress returns the public address of a SECP256K1 public key. The
// public key is hashed using Keccak256, then the hash is Base64 encoded and
// the first 42 runes are returned.
func (id SECP256K1) PublicAddress() string {
	// Hash the public key.
	hash := crypto.Keccak256(id.PublicKey())
	// Convert to runes to preserve UTF-8 character counting.
	runes := []rune(base64.StdEncoding.EncodeToString(hash))
	// Return the first 42 runes.
	return string(runes[0:42])
}

// PublicKey returns the public key as a slice of bytes.
func (id SECP256K1) PublicKey() []byte {
	if id.ECDSAPub == nil && id.ECDSA == nil {
		return nil
	}
	if id.ECDSAPub == nil {
		return crypto.FromECDSAPub(&id.ECDSA.PublicKey)
	}
	return crypto.FromECDSAPub(id.ECDSAPub)
}

// PrivateKey returns the private key as a slice of bytes.
func (id SECP256K1) PrivateKey() []byte {
	if id.ECDSA == nil {
		return nil
	}
	return crypto.FromECDSA(id.ECDSA)
}

// Sign a payload. The payload is hashed using Keccak256 and then signed using
// the SECP256K1 private key. Returns the hash of the payload, the signed hash
// of the payload, or an error.
func (id SECP256K1) Sign(payload []byte) ([]byte, []byte, error) {
	// Hash the payload and sign the hash.
	hash := crypto.Keccak256(payload)
	signedHash, err := crypto.Sign(hash, id.ECDSA)
	if err != nil {
		return nil, nil, err
	}
	// Return the unsigned hash and the signed hash.
	return hash, signedHash, nil
}

// Verify that a hash was signed by the SECP256K1 private key. The unsigned
// hash, and the signed hash that will be verified, are used as input. Returns
// true if the signed hash was signed by the SECP256K1 private key, false
// otherwise, or an error.
func (id SECP256K1) Verify(hash []byte, signedHash []byte) (bool, error) {
	// Recover the public key using the hash and signed hash.
	publicKey, err := crypto.Ecrecover(hash, signedHash)
	if err != nil {
		return false, err
	}
	// Verify that the public keys match.
	return bytes.Equal(id.PublicKey(), publicKey), nil
}

// MarshalJSON implements the json.Marshaler interface.
func (id SECP256K1) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{
		"ecdsa_pub": base64.StdEncoding.EncodeToString(id.PublicKey()),
		"ecdsa":     base64.StdEncoding.EncodeToString(id.PrivateKey()),
	}
	return json.Marshal(data)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (id SECP256K1) UnmarshalJSON(data []byte) error {
	// Unmarshal data into a struct.
	v := struct {
		Public  string `json:"ecdsa_pub"`
		Private string `json:"ecdsa"`
	}{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	// Decode the private key.
	ecdsaBlob, err := base64.StdEncoding.DecodeString(v.Private)
	if err != nil {
		return err
	}
	ecdsa, err := crypto.ToECDSA(ecdsaBlob)
	if err != nil {
		return err
	}
	// Set the private and public key.
	id.ECDSA = ecdsa
	id.ECDSAPub = &ecdsa.PublicKey
	return nil
}
