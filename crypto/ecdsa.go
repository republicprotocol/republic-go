package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"

	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	ethSecp256k1 "github.com/ethereum/go-ethereum/crypto/secp256k1"
	base58 "github.com/jbenet/go-base58"
	multihash "github.com/multiformats/go-multihash"
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

// Address of the EcdsaKey. This is the last 20 bytes of the ecdsa.PublicKey
// after a keccak256 hash and then encoded as a base58 string.
func (key *EcdsaKey) Address() string {
	bytes := elliptic.Marshal(s256, key.PublicKey.X, key.PublicKey.Y)
	hash := ethCrypto.Keccak256(bytes)
	bytes = hash[(len(hash) - 20):]

	hash = make([]byte, 2, 22)
	hash[0], hash[1] = multihash.KECCAK_256, 20
	hash = append(hash, bytes...)
	return base58.EncodeAlphabet(hash, base58.BTCAlphabet)
}

// // RecoverSigner calculates the signing public key given signable data and its signature
// func RecoverSigner(hasher Hasher, signature []byte) (ID, error) {
// 	hash := hasher.Hash()

// 	publicKey, err := c
// 	if err != nil {
// 		return publicKey, err
// 	}

// 	// Returns 65-byte uncompress pubkey (0x04 | X | Y)
// 	pubkey, err := ethcrypto.Ecrecover(hash, signature)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Convert to KeyPair before calculating ID
// 	id := KeyPair{
// 		nil, &ecdsa.PublicKey{
// 			Curve: secp256k1.S256(),
// 			X:     big.NewInt(0).SetBytes(pubkey[1:33]),
// 			Y:     big.NewInt(0).SetBytes(pubkey[33:65]),
// 		},
// 	}.ID()

// 	return id, nil
// }

// // VerifySignature verifies that the data's signature has been signed by the provided
// // ID's private key
// func VerifySignature(hasher crypto.Hasher, signature []byte, id ID) error {
// 	if signature == nil {
// 		return crypto.ErrInvalidSignature
// 	}
// 	signer, err := RecoverSigner(hasher, signature)
// 	if err != nil {
// 		return err
// 	}
// 	if !bytes.Equal(signer, id) {
// 		return crypto.ErrInvalidSignature
// 	}
// 	return nil
// }

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
