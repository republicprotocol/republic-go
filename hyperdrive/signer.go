package hyperdrive

import "errors"

// A Hash is a Keccak256 Hash of some data.
type Hash [32]byte

// Signatures that are being collected, valid once a threshold of unique
// Signatures has been reached.
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

// A Signature produced by a Replica. The current implementation uses SECP256K1
// ECDSA private keys for producing signatures.
type Signature [65]byte

// A Signer can produce Signatures of some Hash of data.
type Signer interface {
	Sign(Hash) (Signature, error)
}

// WeakSigner produces Signatures by returning its ID.
type WeakSigner struct {
	ID [32]byte
}

// NewWeakSigner returns a new WeakSigner.
func NewWeakSigner(id [32]byte) WeakSigner {
	return WeakSigner{
		ID: id,
	}
}

// Sign implements the Signer interface.
func (signer *WeakSigner) Sign(hash Hash) (Signature, error) {
	signature := [65]byte{}
	copy(signature[:], signer.ID[:])
	return Signature(signature), nil
}

// ErrorSigner returns errors instead of producing Signatures.
type ErrorSigner struct {
}

// NewErrorSigner returns a new ErrorSigner.
func NewErrorSigner() ErrorSigner {
	return ErrorSigner{}
}

// Sign implements the Signer interface.
func (signer *ErrorSigner) Sign(hash Hash) (Signature, error) {
	return [65]byte{}, errors.New("cannot use error signer to sign")
}
