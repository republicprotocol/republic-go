package identity

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"golang.org/x/crypto/sha3"
)

// RepublicID contains an ECDSA key pair using a SECP256K1 S256 elliptic curve.
type RepublicID struct {
	*ecdsa.PrivateKey
	*ecdsa.PublicKey
}

// NewRepublicID generates a new ECDSA key pair using a SECP256K1 S256 elliptic
// curve. It returns a RepublicID that uses this key pair.
func NewRepublicID() (RepublicID, error) {
	key, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	if err != nil {
		return RepublicID{}, err
	}
	return RepublicID{
		PrivateKey: key,
		PublicKey:  &key.PublicKey,
	}, nil
}

// PublicAddress implements the Identifier interface.
func (id RepublicID) PublicAddress() string {
	return ""
}

// PublicAddressInBytes implements the Identifier interface.
func (id RepublicID) PublicAddressInBytes() []byte {
	bytes := elliptic.Marshal(secp256k1.S256(), id.PublicKey.X, id.PublicKey.Y)
	hash := sha3.Sum256(bytes)
	return hash[:20]
}
