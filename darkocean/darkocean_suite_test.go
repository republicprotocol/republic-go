package darkocean_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/darknode"
)

func TestDarkocean(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Darkocean Suite")
}

const (
	GanacheRPC        = "http://localhost:8545"
	NumberOfDarkNodes = 24
)

var testnetEnv darknode.TestnetEnv

var _ = BeforeSuite(func() {
	var err error
	testnetEnv, err = darknode.NewTestnet(NumberOfDarkNodes, NumberOfDarkNodes)
	go testnetEnv.Run()
	time.Sleep(2 * time.Second)
	Expect(err).ShouldNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	testnetEnv.Teardown()
})
