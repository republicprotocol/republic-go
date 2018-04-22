package interop_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoAtom(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Atom Protocol Suite")
}
