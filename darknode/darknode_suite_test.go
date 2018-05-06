package darknode_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/darknode"
)

func TestDarknode(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Darknode Suite")
}

const (
	GanacheRPC                 = "http://localhost:8545"
	NumberOfDarkNodes          = 6
	NumberOfBootstrapDarkNodes = 6
)

var env TestnetEnv

var _ = BeforeSuite(func() {
	var err error
	env, err = NewTestnet(NumberOfDarkNodes, NumberOfBootstrapDarkNodes)
	go env.Run()
	time.Sleep(10 * time.Second)
	Expect(err).ShouldNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	env.Teardown()
})
