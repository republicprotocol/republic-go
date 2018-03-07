package dnr_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoDarkNodeRegistrar(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GoDarkNodeRegistrar Suite")
}
