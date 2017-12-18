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

// Error for not successfully decoding the string in base58
var ErrFailToDecode = fmt.Errorf("fail to decode the string")

// IDLength is the number of bytes in an ID.
const IDLength = 20

// An ID is a slice of 20 bytes that can be converted into an Address. It must
// always be example 20 bytes.
type ID []byte

// NewID generates a new ID from a key value pair
func NewID() (ID, error)  {
	keyPair, err := NewKeyPair()
	if err != nil {
		return nil, err
	}
	return keyPair.PublicID(),nil
}

// AddressLength is the number of bytes in an Address.
const AddressLength = 30

// An Address is generated from an ID.
type Address string

// NewAddress generates a new Address from a key value pair
func NewAddress() (Address,error)  {
	keyPair, err := NewKeyPair()
	if err != nil {
		return "", err
	}
	return keyPair.PublicAddress(),nil
}

// Use xor to calculate distance between two Addresses
func (address Address) Distance(other Address) ([]byte ,error) {
	// Decode both addresses into bytes
	idByte := base58.Decode(string(address))
	if len(idByte) == 0 {
		return nil, ErrFailToDecode
	}
	otherByte := base58.Decode(string(other))
	if len(otherByte) == 0 {
		return nil, ErrFailToDecode
	}

	xor := make([]byte, IDLength)
	for i := 2; i < 2+IDLength; i++ {
		xor[i] = idByte[i] ^ otherByte[i]
	}
	return xor,nil
}

// Number of same prefix bits with another address excluding the first 2 bytes
func (address Address) SamePrefixLen(other Address) (int,error) {
	diff, err := address.Distance(other)
	if err != nil {
		return -1, err
	}

	ret := 0
	for i := 0; i < IDLength; i++ {
		if diff[i] == uint8(0) {
			ret += 8
		} else {
			bits := fmt.Sprintf("%08b", diff[i])
			for j := 0; j < len(bits); j++ {
				if bits[j] == '1' {
					return ret,nil
				}
				ret++
			}
		}
	}
	return ret,nil
}

// MultiAddress returns the Republic multi address of the KeyPair.
// It can be encapsulated by other multiaddress
func (address Address) MultiAddress() (multiaddr.Multiaddr, error) {
	return multiaddr.NewMultiaddr(fmt.Sprintf("/republic/%s", string(address)))
}

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
	return Address(base58.Encode(hash))
}

