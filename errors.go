package sss

import (
	"fmt"
	"math/big"
)

// An NKError is used when the number of shares needed to reconstruct a secret
// is larger than the number of shares generated when encoding the secret.
type NKError string

// NewNKError returns a new NKError. The number of shares generated when
// encoding the secret is given by N, and the number of shares required for
// reconstruction is K.
func NewNKError(n, k int64) NKError {
	return NKError(fmt.Sprintf("expected n = %v to be greater than or equal to k = %v", n, k))
}

// Error implements the Error interface for NKError.
func (err NKError) Error() string {
	return string(err)
}

// A FiniteFieldError is used when a secret is outside the finite field being
// used by the secret encoder.
type FiniteFieldError string

// NewFiniteFieldError returns a new FiniteFieldError.
func NewFiniteFieldError(secret *big.Int) FiniteFieldError {
	return FiniteFieldError(fmt.Sprintf("expected secret = %v to be within the finite field", secret))
}

// Error implements the Error interface for FiniteFieldError.
func (err FiniteFieldError) Error() string {
	return string(err)
}
