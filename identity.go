package identity

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/jbenet/go-base58"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/multiformats/go-multihash"
)

const HashLength = 20

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
	hash := []byte{
		multihash.KECCAK_256,
		HashLength,
	}
	hash = append(hash, id.PublicAddressInBytes()...)
	return base58.EncodeAlphabet(hash, base58.BTCAlphabet)
}

// PublicAddressInBytes implements the Identifier interface.
func (id RepublicID) PublicAddressInBytes() []byte {
	bytes := elliptic.Marshal(secp256k1.S256(), id.PublicKey.X, id.PublicKey.Y)
	hash := crypto.Keccak256(bytes)
	return hash[:HashLength]
}


