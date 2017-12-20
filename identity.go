package identity

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"

	"github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multihash"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/jbenet/go-base58"
)


var (
	// Error for not successfully decoding the string in base58
	ErrFailToDecode = fmt.Errorf("fail to decode the string")
	// Error for wrong address length
	ErrWrongAddressLength = fmt.Errorf("wrong address length")
	)

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
	return crypto.Keccak256(bytes)[:IDLength]
}

// PublicAddress returns the Republic Address of this KeyPair. The Address is
// the Base58 encoding of the MultiHash of the Republic ID.
func (keyPair KeyPair) PublicAddress() Address {
	id := keyPair.PublicID()
	hash := make([]byte, 2, IDLength+2)
	// first two byte represent the hash function and length of the ID
	hash[0], hash[1] = multihash.KECCAK_256, IDLength
	hash = append(hash, id...)
	return Address(base58.Encode(hash))
}

// IDLength is the number of bytes in an ID.
const IDLength = 20

// An ID is a slice of 20 bytes that can be converted into an Address.
// It must always be example 20 bytes.
type ID []byte

// NewID generates a new ID from a key value pair
func NewID() (ID, error) {
	keyPair, err := NewKeyPair()
	if err != nil {
		return nil, err
	}
	return keyPair.PublicID(), nil
}

// AddressLength is the number of bytes in an Address.
const AddressLength = 30

// An Address is generated from an ID.
type Address string

// NewAddress generates a new Address from a key value pair
func NewAddress() (Address, error) {
	keyPair, err := NewKeyPair()
	if err != nil {
		return "", err
	}
	return keyPair.PublicAddress(), nil
}

// Use xor to calculate distance between two Addresses
func (address Address) Distance(other Address) ([]byte, error) {

	// Check the length of both addresses
	if len(address) != AddressLength || len(other) != AddressLength {
		return nil, ErrWrongAddressLength
	}

	// Decode both addresses into bytes
	idByte := base58.Decode(string(address))
	if len(idByte) == 0 {
		return nil, ErrFailToDecode
	}
	otherByte := base58.Decode(string(other))
	if len(otherByte) == 0 {
		return nil, ErrFailToDecode
	}

	// Get distance by xor operation
	xor := make([]byte, IDLength)
	for i := 0; i < IDLength; i++ {
		xor[i] = idByte[i+2] ^ otherByte[i+2]
	}
	return xor, nil
}

// SamePrefixLen returns the number of same prefix bits with another
// address excluding the first 2 bytes
func (address Address) SamePrefixLen(other Address) (int, error) {
	// Calculator the distance between two addresses
	diff, err := address.Distance(other)
	if err != nil {
		return -1, err
	}

	// Calculate the same prefix bits
	ret := 0
	for i := 0; i < IDLength; i++ {
		if diff[i] == uint8(0) {
			ret += 8
		} else {
			bits := fmt.Sprintf("%08b", diff[i])
			for j := 0; j < len(bits); j++ {
				if bits[j] == '1' {
					return ret, nil
				}
				ret++
			}
		}
	}
	return ret, nil
}

// MultiAddress returns the Republic multi address of the address.
// So that we can join it with other multiaddress
func (address Address) MultiAddress() (multiaddr.Multiaddr, error) {
	return multiaddr.NewMultiaddr(fmt.Sprintf("/republic/%s", string(address)))
}

// Closer returns true if address1 is closer to the target than
func Closer(address1, address2, target Address) (bool, error) {
	distance1, err := address1.Distance(target)
	if err != nil {
		return false, err
	}
	distance2, err := address2.Distance(target)
	if err != nil {
		return false, err
	}

	for i := 0; i < IDLength; i++ {
		if distance1[i] < distance2[i] {
			return true, nil
		} else if distance1[i] > distance2[i] {
			return false, nil
		}
	}
	// If the addresses are the same, return true
	return true, nil
}
