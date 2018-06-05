package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	ethSecp256k1 "github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/jbenet/go-base58"
	"github.com/multiformats/go-multihash"
	"github.com/republicprotocol/republic-go/identity"
)

// EcdsaKey for signing and verifying hashes.
type EcdsaKey struct {
	*ecdsa.PrivateKey
	address string
}

// RandomEcdsaKey using a secp256k1 s256 curve.
func RandomEcdsaKey() (EcdsaKey, error) {
	privateKey, err := ecdsa.GenerateKey(ethSecp256k1.S256(), rand.Reader)
	if err != nil {
		return EcdsaKey{}, err
	}
	return EcdsaKey{
		PrivateKey: privateKey,
		address:    addressFromPublickey(&privateKey.PublicKey),
	}, nil
}

// NewEcdsaKey returns an EcdsaKey from an existing private key. It does not
// verify that the private key was generated correctly.
func NewEcdsaKey(privateKey *ecdsa.PrivateKey) EcdsaKey {
	return EcdsaKey{
		PrivateKey: privateKey,
		address:    addressFromPublickey(&privateKey.PublicKey),
	}
}

// Sign implements the Signer interface. It uses the ecdsa.PrivateKey to sign
// the data without performing any kind of preprocessing of the data. If the
// data is not exactly 32 bytes, an error is returned.
func (key *EcdsaKey) Sign(data []byte) ([]byte, error) {
	return ethCrypto.Sign(data, key.PrivateKey)
}

// Verify implements the Verifier interface. It uses its own address as the
// expected signatory.
func (key *EcdsaKey) Verify(data []byte, signature []byte) error {
	return NewEcdsaVerifier(key.address).Verify(data, signature)
}

// Address of the EcdsaKey. An Address is generated in the same way as an
// Ethereum address, but instead of a hex encoding, it uses a base58 encoding
// of a Keccak256 multihash.
func (key *EcdsaKey) Address() string {
	return key.address
}

// Equal returns true if two EcdsaKeys are exactly equal. The name of the
// elliptic.Curve is not checked.
func (key *EcdsaKey) Equal(rhs *EcdsaKey) bool {
	return key.address == rhs.address &&
		key.D.Cmp(rhs.D) == 0 &&
		key.X.Cmp(rhs.X) == 0 &&
		key.Y.Cmp(rhs.Y) == 0 &&
		key.Curve.Params().P.Cmp(rhs.Curve.Params().P) == 0 &&
		key.Curve.Params().N.Cmp(rhs.Curve.Params().N) == 0 &&
		key.Curve.Params().B.Cmp(rhs.Curve.Params().B) == 0 &&
		key.Curve.Params().Gx.Cmp(rhs.Curve.Params().Gx) == 0 &&
		key.Curve.Params().Gy.Cmp(rhs.Curve.Params().Gy) == 0 &&
		key.Curve.Params().BitSize == rhs.Curve.Params().BitSize
}

// MarshalJSON implements the json.Marshaler interface. The EcdsaKey is
// formatted according to the Republic Protocol Keystore specification.
func (key EcdsaKey) MarshalJSON() ([]byte, error) {
	jsonKey := map[string]interface{}{}
	// Private key
	jsonKey["d"] = key.D.Bytes()

	// Public key
	ethAddress, err := republicAddressToEthAddress(key.address)
	if err != nil {
		return []byte{}, err
	}

	jsonKey["address"] = ethAddress.Hex()
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

	key.address = addressFromPublickey(&key.PublicKey)
	return nil
}

// EcdsaVerifier is used to verify signatures produced by an EcdsaKey.
type EcdsaVerifier struct {
	addr string
}

// NewEcdsaVerifier returns an EcdsaVerifier that expects the signatory of all
// signatures that it checks to equal the given address.
func NewEcdsaVerifier(addr string) EcdsaVerifier {
	return EcdsaVerifier{
		addr: addr,
	}
}

// Verify implements the Verifier interface.
func (verifier EcdsaVerifier) Verify(data []byte, signature []byte) error {
	if data == nil || len(data) == 0 {
		return ErrNilData
	}
	if signature == nil || len(signature) == 0 {
		return ErrNilSignature
	}
	addrRecovered, err := RecoverAddress(data, signature)
	if err != nil {
		return err
	}
	if addrRecovered != verifier.addr {
		return ErrInvalidSignature
	}
	return nil
}

// RecoverAddress used to produce a signature.
func RecoverAddress(data []byte, signature []byte) (string, error) {

	// Returns 65-byte uncompress pubkey (0x04 | X | Y)
	publicKey, err := ethCrypto.Ecrecover(data, signature)
	if err != nil {
		return "", err
	}

	// Address from an EcdsaKey
	privateKey := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: ethSecp256k1.S256(),
			X:     big.NewInt(0).SetBytes(publicKey[1:33]),
			Y:     big.NewInt(0).SetBytes(publicKey[33:65]),
		},
	}
	key := EcdsaKey{
		PrivateKey: privateKey,
		address:    addressFromPublickey(&privateKey.PublicKey),
	}
	return key.Address(), nil
}

func addressFromPublickey(publicKey *ecdsa.PublicKey) string {
	bytes := elliptic.Marshal(ethSecp256k1.S256(), publicKey.X, publicKey.Y)
	hash := ethCrypto.Keccak256(bytes[1:]) // Keccak256 hash
	hash = hash[(len(hash) - 20):]         // Take the last 20 bytes
	addr := make([]byte, 2, 22)            // Create the multihash address
	addr[0] = multihash.KECCAK_256         // Set the keccak256 byte
	addr[1] = 20                           // Set the length byte
	addr = append(addr, hash...)           // Append the data
	return base58.EncodeAlphabet(addr, base58.BTCAlphabet)
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

// Convert republic address to ethereum address
func republicAddressToEthAddress(repAddress string) (common.Address, error) {
	addByte := base58.DecodeAlphabet(repAddress, base58.BTCAlphabet)[2:]
	if len(addByte) == 0 {
		return common.Address{}, errors.New("fail to decode the address")
	}
	address := common.BytesToAddress(addByte)
	return address, nil
}

// Convert republic address to ethereum address
func ethAddressToRepublicAddress(ethAddress string) identity.Address {
	address := common.HexToAddress(ethAddress)
	addr := make([]byte, 2, 22)
	addr[0] = multihash.KECCAK_256
	addr[1] = 20
	addr = append(addr, address.Bytes()...)
	return identity.Address(base58.EncodeAlphabet(addr, base58.BTCAlphabet))
}
