package identity

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

// A Hash is a Keccak256 Hash of some data.
type Hash [32]byte

// A Signable struct is able to be signed by a KeyPair
// Hash must return a 32-byte []byte array
type Signable interface {
	Hash() Hash
}

// A Signer can produce Signatures of some Hash of data.
type Signer interface {
	Sign(Hash) (Signature, error)
}

// A Verifier is used to verify Messages.
type Verifier interface {
	VerifyProposer(Signature) error
	VerifySignatures(Signatures) error
}

// The Signature type represents the signature of the hash of Signable data
type Signature = [65]byte

// ErrInvalidSignature indicates that a signature could not be verified against the provided Identity
var ErrInvalidSignature = fmt.Errorf("failed to verify signature")

// ErrSignData indicates that we fail to sign the data.
var ErrSignData = fmt.Errorf("failed to sign the data")

// Sign hashes and signs the Signable data
// If the Hash() function defined does not correctly hash the struct,
// it may allow for chosen plaintext attacks on the keypair's private key
func (keyPair *KeyPair) Sign(data Signable) (Signature, error) {
	hash := data.Hash()

	signed, err := crypto.Sign(hash[:], keyPair.PrivateKey)
	if err != nil {
		return Signature{}, err
	}

	var signature [65]byte
	copied := copy(signature[:], signed[:])
	if copied != 65 {
		return Signature{}, ErrSignData
	}

	return signature, nil
}

// RecoverSigner calculates the signing public key given signable data and its signature
func RecoverSigner(data Signable, signature Signature) (ID, error) {
	hash := data.Hash()

	// Returns 65-byte uncompressed publicKey (0x04 | X | Y)
	publicKey, err := crypto.Ecrecover(hash[:], signature[:])
	if err != nil {
		return nil, err
	}

	// Convert to KeyPair before calculating ID
	id := KeyPair{
		nil, &ecdsa.PublicKey{
			Curve: secp256k1.S256(),
			X:     big.NewInt(0).SetBytes(publicKey[1:33]),
			Y:     big.NewInt(0).SetBytes(publicKey[33:65]),
		},
	}.ID()

	return id, nil
}

// VerifySignature verifies that the data's signature has been signed by the provided
// ID's private key
func VerifySignature(data Signable, signature Signature, id ID) error {
	signer, err := RecoverSigner(data, signature)
	if err != nil {
		return err
	}
	if !bytes.Equal(signer, id) {
		return ErrInvalidSignature
	}
	return nil
}

// Signatures is a slice of Signature.
type Signatures []Signature

// Merge Signatures together and avoid duplication. Returns the merged
// Signatures without modifying the inputs.
func (signatures Signatures) Merge(others Signatures) Signatures {
	merger := map[Signature]struct{}{}
	for i := range signatures {
		merger[signatures[i]] = struct{}{}
	}
	for i := range others {
		merger[others[i]] = struct{}{}
	}

	i := 0
	mergedSignatures := make(Signatures, len(merger))
	for key := range merger {
		mergedSignatures[i] = key
		i++
	}
	return mergedSignatures
}
