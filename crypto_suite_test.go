package crypto_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoCrypto(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GoCrypto Suite")
}
