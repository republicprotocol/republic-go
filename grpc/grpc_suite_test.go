package grpc_test

import (
	"crypto/rand"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRpc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "gRPC Suite")
}

// mockVerifier is a simple implementation of the crypto.Verifier interface
// that accepts all signatures.
type mockVerifier struct {
}

func (verifier mockVerifier) Verify(data []byte, signature []byte) error {
	return nil
}

// mockSigner is a simple implementation of the crypto.Signer interface that
// produces random signatures.
type mockSigner struct {
}

func (signer mockSigner) Sign(data []byte) ([]byte, error) {
	signature := [65]byte{}
	_, err := rand.Read(signature[:])
	return signature[:], err
}
