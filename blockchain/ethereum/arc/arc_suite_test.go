package arc_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEthereum(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Arc Suite")
}
