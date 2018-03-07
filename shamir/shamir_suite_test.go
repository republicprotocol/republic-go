package shamir_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestShamir(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shamir's Secret Sharing Suite")
}
