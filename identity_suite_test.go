package identity_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAddress(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Republic Identity Suite")
}
