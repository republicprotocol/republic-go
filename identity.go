package identity

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/jbenet/go-base58"
	"github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multihash"
)

const IDLength = 20

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
func (keyPair KeyPair) PublicID() []byte {
	bytes := elliptic.Marshal(secp256k1.S256(), keyPair.PublicKey.X, keyPair.PublicKey.Y)
	hash := crypto.Keccak256(bytes)
	return hash[:IDLength]
}

// PublicAddress returns the Republic Address of this KeyPair. The Address is
// the Base58 encoding of the MultiHash of the Republic ID.
func (keyPair KeyPair) PublicAddress() string {
	hash := make([]byte, 0, 2+IDLength)
	hash = append(hash, multihash.KECCAK_256, IDLength)
	hash = append(hash, keyPair.PublicID()...)
	return base58.EncodeAlphabet(hash, base58.BTCAlphabet)
}

// MultiAddress returns the Republic multi address of the KeyPair.
// It can be encapsulated by other multiaddress
func (keyPair KeyPair) MultiAddress() (multiaddr.Multiaddr, error) {
	return NewMultiaddr("/republic/" + keyPair.PublicAddress())
}
