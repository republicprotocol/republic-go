package darknode_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/blockchain/test/ganache"
)

func TestGoMiner(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Darknode Suite")
}

var _ = BeforeSuite(func() {
	_, err := ganache.StartAndConnect()
	Expect(err).ShouldNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	ganache.Stop()
})
