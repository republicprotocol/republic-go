package identity

import (
	"crypto/ecdsa"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

// TODO: Instead of Hash(), change to Bytes() / Serialize() or similar

// A Signable struct must implement the Hash function
type Signable interface {
	Hash() []byte
}

// The Signature type represents the signature of the hash of Signable data
type Signature = []byte

// Sign hashes and signs the Signable data
func (keyPair *KeyPair) Sign(data Signable) (Signature, error) {
	hash := data.Hash()

	return crypto.Sign(hash, keyPair.PrivateKey)
}

// RecoverSigner calculates the signing public key given signable data and its signature
func (keyPair *KeyPair) RecoverSigner(data Signable, signature Signature) (KeyPair, error) {
	hash := data.Hash()

	// Returns 65-byte uncompress pubkey (0x04 | X | Y)
	pubkey, err := crypto.Ecrecover(hash, signature)
	if err != nil {
		return KeyPair{}, err
	}

	return KeyPair{
		nil, &ecdsa.PublicKey{
			Curve: secp256k1.S256(),
			X:     big.NewInt(0).SetBytes(pubkey[1:33]),
			Y:     big.NewInt(0).SetBytes(pubkey[33:65]),
		},
	}, nil
}

// VerifySignature verifies that the data's signature corresponds to the provided ID
func (keyPair *KeyPair) VerifySignature(data Signable, signature Signature, id ID) error {
	signer, err := keyPair.RecoverSigner(data, signature)
	if err != nil {
		return err
	}
	// TODO: Don't convert to string to compare
	if signer.ID().String() != id.String() {
		return errors.New("invalid signature")
	}
	return nil
}
