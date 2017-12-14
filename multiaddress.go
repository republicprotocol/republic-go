package identity

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"golang.org/x/crypto/sha3"
)

// PrivateIdentifier is a Republic ID that contains both the private and public
// ECDSA keys.
type PrivateIdentifier struct {
	*ecdsa.PrivateKey
	*ecdsa.PublicKey
}

// NewPrivateIdentifier generates a new ECDSA key pair using a SECP256K1 S256
// elliptic curve. It returns a PrivateIdentifier that uses this key pair.
func NewPrivateIdentifier() (PrivateIdentifier, error) {
	key, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	if err != nil {
		return PrivateIdentifier{}, err
	}
	return PrivateIdentifier{
		PrivateKey: key,
		PublicKey:  &key.PublicKey,
	}, nil
}

// PublicAddress implements the Identifier interface.
func (id PrivateIdentifier) PublicAddress() string {
	return ""
}

// PublicAddressInBytes implements the Identifier interface.
func (id PrivateIdentifier) PublicAddressInBytes() []byte {
	bytes := elliptic.Marshal(secp256k1.S256(), id.PublicKey.X, id.PublicKey.Y)
	hash := sha3.Sum256(bytes)
	return hash[:20]
}
