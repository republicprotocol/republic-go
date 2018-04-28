package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"

	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	ethSecp256k1 "github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/jbenet/go-base58"
	"github.com/multiformats/go-multihash"
)

// EcdsaKey for signing and verifying hashes.
type EcdsaKey struct {
	*ecdsa.PrivateKey
}

// RandomEcdsaKey using a secp256k1 s256 curve.
func RandomEcdsaKey() (EcdsaKey, error) {
	privateKey, err := ecdsa.GenerateKey(ethSecp256k1.S256(), rand.Reader)
	if err != nil {
		return EcdsaKey{}, err
	}
	return EcdsaKey{
		PrivateKey: privateKey,
	}, nil
}

// NewEcdsaKey returns an EcdsaKey from an existing private key. It does not
// verify that the private key was generated correctly.
func NewEcdsaKey(privateKey *ecdsa.PrivateKey) EcdsaKey {
	return EcdsaKey{
		PrivateKey: privateKey,
	}
}

// Address of the EcdsaKey. An Address is generated in the same way as an
// Ethereum address, but instead of a hex encoding, it uses a base58 encoding
// of a Keccak256 multihash.
func (key *EcdsaKey) Address() string {
	bytes := elliptic.Marshal(s256, key.PublicKey.X, key.PublicKey.Y)
	hash := ethCrypto.Keccak256(bytes[1:]) // Keccak256 hash
	hash = hash[(len(hash) - 20):]         // Take the last 20 bytes
	addr := make([]byte, 2, 22)            // Create the multihash address
	addr[0] = multihash.KECCAK_256         // Set the keccak256 byte
	addr[1] = 20                           // Set the length byte
	addr = append(addr, hash...)           // Append the data
	return base58.EncodeAlphabet(hash, base58.BTCAlphabet)
}

// VerifySignature of a hash matches an address.
func VerifySignature(hasher Hasher, signature []byte, addr string) error {
	if signature == nil {
		return ErrInvalidSignature
	}
	addrRecovered, err := RecoverAddress(hasher, signature)
	if err != nil {
		return err
	}
	if addr != addrRecovered {
		return ErrInvalidSignature
	}
	return nil
}

// RecoverAddress used to produce a signature.
func RecoverAddress(hasher Hasher, signature []byte) (string, error) {

	// Returns 65-byte uncompress pubkey (0x04 | X | Y)
	publicKey, err := ethCrypto.Ecrecover(hasher.Hash(), signature)
	if err != nil {
		return "", err
	}

	// Address from an EcdsaKey
	key := EcdsaKey{
		PrivateKey: &ecdsa.PrivateKey{
			PublicKey: ecdsa.PublicKey{
				Curve: ethSecp256k1.S256(),
				X:     big.NewInt(0).SetBytes(publicKey[1:33]),
				Y:     big.NewInt(0).SetBytes(publicKey[33:65]),
			},
		},
	}
	return key.Address(), nil
}

// s256 is unused
var s256 *elliptic.CurveParams

func init() {
	// See SEC 2 section 2.7.1
	// curve parameters taken from:
	// http://www.secg.org/collateral/sec2_final.pdf
	s256.P, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F", 16)
	s256.N, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)
	s256.B, _ = new(big.Int).SetString("0000000000000000000000000000000000000000000000000000000000000007", 16)
	s256.Gx, _ = new(big.Int).SetString("79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798", 16)
	s256.Gy, _ = new(big.Int).SetString("483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8", 16)
	s256.BitSize = 256
}
