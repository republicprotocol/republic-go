package identity

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

// A Signable struct is able to be signed by a KeyPair
type Signable interface {
	SerializeForSigning() []byte
}

// The Signature type represents the signature of the hash of Signable data
type Signature = []byte

// ErrInvalidSignature indicates that a signature could not be verified against the provided Identity
var ErrInvalidSignature = fmt.Errorf("failed to verify signature")

// Sign hashes and signs the Signable data
func (keyPair *KeyPair) Sign(data Signable) (Signature, error) {
	hash := crypto.Keccak256(data.SerializeForSigning())

	return crypto.Sign(hash, keyPair.PrivateKey)
}

// RecoverSigner calculates the signing public key given signable data and its signature
func RecoverSigner(data Signable, signature Signature) (ID, error) {
	hash := crypto.Keccak256(data.SerializeForSigning())

	// Returns 65-byte uncompress pubkey (0x04 | X | Y)
	pubkey, err := crypto.Ecrecover(hash, signature)
	if err != nil {
		return nil, err
	}

	// Convert to KeyPair before calculating ID
	id := KeyPair{
		nil, &ecdsa.PublicKey{
			Curve: secp256k1.S256(),
			X:     big.NewInt(0).SetBytes(pubkey[1:33]),
			Y:     big.NewInt(0).SetBytes(pubkey[33:65]),
		},
	}.ID()

	return id, nil
}

// VerifySignature verifies that the data's signature corresponds to the provided ID
func VerifySignature(data Signable, signature Signature, id ID) error {
	signer, err := RecoverSigner(data, signature)
	if err != nil {
		return err
	}
	// TODO: Don't convert to string to compare
	if signer.String() != id.String() {
		return ErrInvalidSignature
	}
	return nil
}
