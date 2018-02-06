package identity

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/jbenet/go-base58"
	"github.com/multiformats/go-multihash"
)

// KeyPair contains an ECDSA key pair created using a SECP256K1 S256 elliptic
// curve.
type KeyPair struct {
	*ecdsa.PrivateKey
	*ecdsa.PublicKey
}

// NewKeyPair generates a new ECDSA key pair using a SECP256K1 S256 elliptic
// curve. It returns a randomly generated KeyPair, or an error.
func NewKeyPair() (KeyPair, error) {
	key, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	if err != nil {
		return KeyPair{}, err
	}
	return KeyPair{
		PrivateKey: key,
		PublicKey:  &key.PublicKey,
	}, nil
}

// NewKeyPairFromPrivateKey a new ECDSA key pair using a given private key. It
// does not validate that this private key was generated correctly.
func NewKeyPairFromPrivateKey(key *ecdsa.PrivateKey) (KeyPair, error) {
	return KeyPair{
		PrivateKey: key,
		PublicKey:  &key.PublicKey,
	}, nil
}

// ID returns the Republic ID of the KeyPair.
func (keyPair KeyPair) ID() ID {
	bytes := elliptic.Marshal(secp256k1.S256(), keyPair.PublicKey.X, keyPair.PublicKey.Y)
	return crypto.Keccak256(bytes)[:IDLength]
}

// Address returns the Republic Address of the KeyPair.
func (keyPair KeyPair) Address() Address {
	id := keyPair.ID()
	hash := make([]byte, 2, IDLength+2)
	hash[0], hash[1] = multihash.KECCAK_256, IDLength
	hash = append(hash, id...)
	return Address(base58.EncodeAlphabet(hash, base58.BTCAlphabet))
}

// MarshalJSON implements the json.Marshaler interface.
func (keyPair KeyPair) MarshalJSON() ([]byte, error) {
	keyPairAsBytes := crypto.FromECDSA(keyPair.PrivateKey)
	keyPairAsString := base58.EncodeAlphabet(keyPairAsBytes, base58.BTCAlphabet)
	return json.Marshal(keyPairAsString)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (keyPair *KeyPair) UnmarshalJSON(data []byte) error {
	keyPairAsString := ""
	if err := json.Unmarshal(data, &keyPairAsString); err != nil {
		return err
	}
	keyPairAsBytes := base58.DecodeAlphabet(keyPairAsString, base58.BTCAlphabet)
	privateKey, err := crypto.ToECDSA(keyPairAsBytes)
	if err != nil {
		return err
	}
	keyPair.PrivateKey = privateKey
	keyPair.PublicKey = &privateKey.PublicKey
	return nil
}
