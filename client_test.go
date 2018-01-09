package swarm_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-swarm"
	"github.com/republicprotocol/go-identity"
	"fmt"
)

var _ = Describe("Client", func() {
	var keyPair identity.KeyPair
	var err error

	BeforeEach(func() {
		keyPair, err = identity.NewKeyPair()
		Ω(err).ShouldNot(HaveOccurred())
		multi, err := identity.NewMultiAddress(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d/republic/%s", 3000, keyPair.Address()))
		Ω(err).ShouldNot(HaveOccurred())
		_, _, err = swarm.NewNodeClient(multi)
		Ω(err).ShouldNot(HaveOccurred())
	})

	Context("Ping", func() {
		It("should be able to ping a target node and get its multiaddress", func() {
			target,_, err :=  identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			multi,err  := target.MultiAddress()
			Ω(err).ShouldNot(HaveOccurred())

		})
	})

	Describe("Malformed configuration files", func() {
		It("should be finished in the future", func() {
			Ω(false).Should(Equal(true))
		})
	})

})
