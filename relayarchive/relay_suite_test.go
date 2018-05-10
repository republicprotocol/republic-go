package relay_test

import (
	"testing"
	"time"

	"github.com/republicprotocol/republic-go/relay"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/darknode"
)

func TestRelay(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Relay Suite")
}

const (
	GanacheRPC             = "http://localhost:8545"
	NumberOfDarkNodes      = 24
	NumberOfBootstrapNodes = 5
	NumberOfRelays         = 1
)

var darknodeTestnetEnv darknode.TestnetEnv
var relayTestNetEnv relay.TestnetEnv

var _ = BeforeSuite(func() {
	var err error
	darknodeTestnetEnv, err = darknode.NewTestnet(NumberOfDarkNodes, NumberOfBootstrapNodes)
	go darknodeTestnetEnv.Run()
	Expect(err).ShouldNot(HaveOccurred())
	time.Sleep(10 * time.Second)
	relayTestNetEnv, err = relay.NewTestnet(NumberOfRelays, darknodeTestnetEnv.DarknodeRegistry, 3000+NumberOfDarkNodes, darknodeTestnetEnv.BootstrapMultiAddrs, darknodeTestnetEnv.Darknodes[0].Config.Keystore)
	go relayTestNetEnv.Run(3000 + NumberOfDarkNodes)
	time.Sleep(5 * time.Second)
	Expect(err).ShouldNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	darknodeTestnetEnv.Teardown()
})
