package topology

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-x"
	"github.com/republicprotocol/go-x/rpc"
)

var _ = Describe("Pair topologies", func() {

	Context("when pinging", func() {
		It("should update their DHTs", func() {
			μ.Lock()
			defer μ.Unlock()

			// Create the left Node.
			keyPair, err := identity.NewKeyPair()
			Ω(err).ShouldNot(HaveOccurred())
			multiAddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/republic/%s", 3000, keyPair.Address()))
			Ω(err).ShouldNot(HaveOccurred())
			left, err := x.NewNode(&x.Config{
				KeyPair:      keyPair,
				MultiAddress: multiAddress,
				MultiAddresses:        make(identity.MultiAddresses, 0, numberOfNodes-1),
			})
			Ω(err).ShouldNot(HaveOccurred())

			// Create the right Node.
			keyPair, err = identity.NewKeyPair()
			Ω(err).ShouldNot(HaveOccurred())
			multiAddress, err = identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/republic/%s", 3001, keyPair.Address()))
			Ω(err).ShouldNot(HaveOccurred())
			right, err := x.NewNode(&x.Config{
				KeyPair:      keyPair,
				MultiAddress: multiAddress,
				MultiAddresses:        make(identity.MultiAddresses, 0, numberOfNodes-1),
			})
			Ω(err).ShouldNot(HaveOccurred())

			// Start the left and right Node.
			go func() {
				defer GinkgoRecover()
				Ω(left.Serve()).ShouldNot(HaveOccurred())
			}()
			go func() {
				defer GinkgoRecover()
				Ω(right.Serve()).ShouldNot(HaveOccurred())
			}()
			defer left.Stop()
			defer right.Stop()

			time.Sleep(startTimeDelay)

			// Ping the left Node from the right Node.
			err = x.Ping(left.MultiAddress, &rpc.MultiAddress{Multi: right.MultiAddress.String()})
			Ω(err).ShouldNot(HaveOccurred())

			// Ping the right Node from the left Node.
			err = x.Ping(right.MultiAddress, &rpc.MultiAddress{Multi: left.MultiAddress.String()})
			Ω(err).ShouldNot(HaveOccurred())

			Ω(len(left.DHT.MultiAddresses())).Should(Equal(1))
			Ω(len(right.DHT.MultiAddresses())).Should(Equal(1))
		})
	})
})
