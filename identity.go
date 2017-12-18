package identity

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"

	"github.com/multiformats/go-multihash"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/jbenet/go-base58"
	"github.com/multiformats/go-multiaddr"
)

// IDLength is the number of bytes in an ID.
const IDLength = 20

// An ID is a slice of 20 bytes that can be converted into an Address. It must
// always be example 20 bytes.
type ID []byte

// An Address is generated from an ID.
type Address string

// KeyPair contains an ECDSA key pair using a SECP256K1 S256 elliptic curve.
type KeyPair struct {
	*ecdsa.PrivateKey
	*ecdsa.PublicKey
}

// NewKeyPair generates a new ECDSA key pair using a SECP256K1 S256 elliptic
// curve. It returns a RepublicID that uses this key pair.
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

// PublicID returns the Republic ID of this KeyPair. The ID is the first 20
// bytes of Keccak256 hash of the public key.
func (keyPair KeyPair) PublicID() ID {
	bytes := elliptic.Marshal(secp256k1.S256(), keyPair.PublicKey.X, keyPair.PublicKey.Y)
	return crypto.Keccak256(bytes)[:20]
}

// PublicAddress returns the Republic Address of this KeyPair. The Address is
// the Base58 encoding of the MultiHash of the Republic ID.
func (keyPair KeyPair) PublicAddress() Address {
	id := keyPair.PublicID()
	hash := make([]byte, 2, 20)
	hash[0] = multihash.KECCAK_256
	hash[1] = IDLength
	hash = append(hash, id...)
	return Address(base58.EncodeAlphabet(hash, base58.BTCAlphabet))
}

// MultiAddress returns the Republic multi address of the KeyPair.
// It can be encapsulated by other multiaddress
func (keyPair KeyPair) MultiAddress() (multiaddr.Multiaddr, error) {
	return multiaddr.NewMultiaddr(fmt.Sprintf("/republic/%s", string(keyPair.PublicAddress())))
}
