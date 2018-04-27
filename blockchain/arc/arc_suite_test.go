package arc_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/blockchain/test"
	"github.com/republicprotocol/republic-go/blockchain/test/ganache"
)

func TestArc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Arc Suite")
}

var _ = test.SkipCIBeforeSuite(func() {
	_, err := ganache.StartAndConnect()
	Expect(err).ShouldNot(HaveOccurred())
})

var _ = test.SkipCIAfterSuite(func() {
	ganache.Stop()
})
