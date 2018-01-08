package identity

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/jbenet/go-base58"
	"github.com/multiformats/go-multihash"
)

// Errors returned by the package.
var (
	ErrFailToDecode       = fmt.Errorf("fail to decode the string")
	ErrWrongAddressLength = fmt.Errorf("wrong address length")
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

// IDLength is the number of bytes in an ID.
const IDLength = 20

// An ID is a slice of 20 bytes that can be converted into an Address.
// It must always be example 20 bytes.
type ID []byte

// NewID generates a new ID by generating a random KeyPair. It returns the ID,
// and the KeyPair, or an error. It is most commonly used for testing.
func NewID() (ID, KeyPair, error) {
	keyPair, err := NewKeyPair()
	if err != nil {
		return nil, keyPair, err
	}
	return keyPair.ID(), keyPair, nil
}

// AddressLength is the number of bytes in an Address.
const AddressLength = 30

// An Address string is generated from an ID.
type Address string

// Addresses is an alias.
type Addresses []Address

// NewAddress generates a new Address by generating a random KeyPair. It
// returns the Address, and the KeyPair, or an error.  It is most commonly used
// for testing.
func NewAddress() (Address, KeyPair, error) {
	keyPair, err := NewKeyPair()
	if err != nil {
		return "", keyPair, err
	}
	return keyPair.Address(), keyPair, nil
}

// Distance uses a bitwise XOR to calculate distance between two Addresses.
func (address Address) Distance(other Address) ([]byte, error) {
	// Check the length of both Addresses.
	if len(address) != AddressLength || len(other) != AddressLength {
		return nil, ErrWrongAddressLength
	}
	// Decode both Addresses into bytes.
	idByte := base58.DecodeAlphabet(string(address), base58.BTCAlphabet)
	if len(idByte) == 0 {
		return nil, ErrFailToDecode
	}
	otherByte := base58.DecodeAlphabet(string(other), base58.BTCAlphabet)
	if len(otherByte) == 0 {
		return nil, ErrFailToDecode
	}
	// Get distance by using the XOR operation.
	xor := make([]byte, IDLength)
	for i := 0; i < IDLength; i++ {
		xor[i] = idByte[i+2] ^ otherByte[i+2]
	}
	return xor, nil
}

// SamePrefixLength returns the number of prefix bits that match with those of
// another Address, excluding the first 2 bytes.
func (address Address) SamePrefixLength(other Address) (int, error) {
	// Calculator the distance between two Addresses.
	diff, err := address.Distance(other)
	if err != nil {
		return -1, err
	}
	// Calculate the same prefix bits.
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

// MultiAddress returns the Republic multi-address of the Address. It can be
// appended with other MultiAddresses.
func (address Address) MultiAddress() (MultiAddress, error) {
	return NewMultiAddress(fmt.Sprintf("/republic/%s", string(address)))
}

// Closer returns true if the left Address is closer to the target than the
// right Address, otherwise it returns false. In the case that both Addresses
// are equal distances from the target, it returns true.
func Closer(left, right, target Address) (bool, error) {
	leftDist, err := left.Distance(target)
	if err != nil {
		return false, err
	}
	rightDist, err := right.Distance(target)
	if err != nil {
		return false, err
	}
	for i := 0; i < IDLength; i++ {
		if leftDist[i] < rightDist[i] {
			return true, nil
		} else if leftDist[i] > rightDist[i] {
			return false, nil
		}
	}
	// If the Addresses are the same, return true.
	return true, nil
}
