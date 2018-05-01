package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
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

// Sign implements the Signer interface. It uses the ecdsa.PrivateKey to sign
// the hash produced by a Hasher.
func (key *EcdsaKey) Sign(hasher Hasher) ([]byte, error) {
	return ethCrypto.Sign(hasher.Hash(), key.PrivateKey)
}

// Address of the EcdsaKey. An Address is generated in the same way as an
// Ethereum address, but instead of a hex encoding, it uses a base58 encoding
// of a Keccak256 multihash.
func (key *EcdsaKey) Address() string {
	bytes := elliptic.Marshal(ethSecp256k1.S256(), key.PublicKey.X, key.PublicKey.Y)
	hash := ethCrypto.Keccak256(bytes[1:]) // Keccak256 hash
	hash = hash[(len(hash) - 20):]         // Take the last 20 bytes
	addr := make([]byte, 2, 22)            // Create the multihash address
	addr[0] = multihash.KECCAK_256         // Set the keccak256 byte
	addr[1] = 20                           // Set the length byte
	addr = append(addr, hash...)           // Append the data
	return base58.EncodeAlphabet(addr, base58.BTCAlphabet)
}

// MarshalJSON implements the json.Marshaler interface. The EcdsaKey is
// formatted according to the Republic Protocol Keystore specification.
func (key *EcdsaKey) MarshalJSON() ([]byte, error) {
	jsonKey := map[string]interface{}{}
	// Private key
	jsonKey["d"] = key.D.Bytes()

	// Public key
	jsonKey["x"] = key.X.Bytes()
	jsonKey["y"] = key.Y.Bytes()

	// Curve
	jsonKey["curveParams"] = map[string]interface{}{
		"p":    ethSecp256k1.S256().P.Bytes(),  // the order of the underlying field
		"n":    ethSecp256k1.S256().N.Bytes(),  // the order of the base point
		"b":    ethSecp256k1.S256().B.Bytes(),  // the constant of the curve equation
		"x":    ethSecp256k1.S256().Gx.Bytes(), // (x,y) of the base point
		"y":    ethSecp256k1.S256().Gy.Bytes(),
		"bits": ethSecp256k1.S256().BitSize, // the size of the underlying field
		"name": "s256",                      // the canonical name of the curve
	}
	return json.Marshal(jsonKey)
}

// UnmarshalJSON implements the json.Unmarshaler interface. An EcdsaKey is
// created from data that is assumed to be compliant with the Republic Protocol
// Keystore specification. The use of secp256k1 s256 curve is not checked.
func (key *EcdsaKey) UnmarshalJSON(data []byte) error {
	jsonKey := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &jsonKey); err != nil {
		return err
	}

	var err error

	// Private key
	key.PrivateKey = new(ecdsa.PrivateKey)
	key.PrivateKey.D, err = unmarshalBigIntFromMap(jsonKey, "d")
	if err != nil {
		return err
	}

	// Public key
	key.PrivateKey.PublicKey = ecdsa.PublicKey{}
	key.PrivateKey.PublicKey.X, err = unmarshalBigIntFromMap(jsonKey, "x")
	if err != nil {
		return err
	}
	key.PrivateKey.PublicKey.Y, err = unmarshalBigIntFromMap(jsonKey, "y")
	if err != nil {
		return err
	}

	// Curve
	if jsonVal, ok := jsonKey["curveParams"]; ok {
		curveParams := elliptic.CurveParams{}
		jsonCurveParams := map[string]json.RawMessage{}
		if err := json.Unmarshal(jsonVal, &jsonCurveParams); err != nil {
			return err
		}
		curveParams.P, err = unmarshalBigIntFromMap(jsonCurveParams, "p")
		if err != nil {
			return err
		}
		curveParams.N, err = unmarshalBigIntFromMap(jsonCurveParams, "n")
		if err != nil {
			return err
		}
		curveParams.B, err = unmarshalBigIntFromMap(jsonCurveParams, "b")
		if err != nil {
			return err
		}
		curveParams.Gx, err = unmarshalBigIntFromMap(jsonCurveParams, "x")
		if err != nil {
			return err
		}
		curveParams.Gy, err = unmarshalBigIntFromMap(jsonCurveParams, "y")
		if err != nil {
			return err
		}
		curveParams.BitSize, err = unmarshalIntFromMap(jsonCurveParams, "bits")
		if err != nil {
			return err
		}
		curveParams.Name, err = unmarshalStringFromMap(jsonCurveParams, "name")
		if err != nil {
			return err
		}
		key.PrivateKey.Curve = &curveParams
	} else {
		return fmt.Errorf("curveParams is nil")
	}
	return nil
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
var s256 elliptic.CurveParams

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
