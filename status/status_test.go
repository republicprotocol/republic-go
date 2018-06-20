package status_test

import (
	"fmt"
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/dht"
	"github.com/republicprotocol/republic-go/identity"
	. "github.com/republicprotocol/republic-go/status"
	"github.com/republicprotocol/republic-go/testutils"
)

var _ = Describe("Status", func() {
	testStr := "someRandomString"

	Context("when writing and reading to the provider", func() {
		var prov Provider
		var confAddr identity.Address
		var _dht dht.DHT

		BeforeEach(func() {
			var err error
			confAddr, err = testutils.RandomAddress()
			Expect(err).ShouldNot(HaveOccurred())
			_dht = dht.NewDHT(confAddr, 64)
			Expect(err).ShouldNot(HaveOccurred())
			prov = NewProvider(&_dht)
		})

		It("should store network information correctly", func() {
			err := prov.WriteNetwork(testStr)
			Expect(err).ShouldNot(HaveOccurred())
			network, err := prov.Network()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(network).Should(Equal(testStr))
		})

		It("should store ethereum address correctly", func() {
			err := prov.WriteEthereumAddress(testStr)
			Expect(err).ShouldNot(HaveOccurred())
			ethAddr, err := prov.EthereumAddress()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(ethAddr).Should(Equal(testStr))
		})

		It("should store multi-address correctly", func() {
			ipAddr := "123.123.123.123"
			maString := fmt.Sprintf("/ip4/%s/tcp/18514/republic/%s", ipAddr, confAddr)
			multiAddr, err := identity.NewMultiAddressFromString(maString)
			Expect(err).ShouldNot(HaveOccurred())
			err = prov.WriteMultiAddress(multiAddr)
			Expect(err).ShouldNot(HaveOccurred())
			readMultiAddr, err := prov.MultiAddress()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(multiAddr).Should(Equal(readMultiAddr))
		})

		It("should store the public key correctly", func() {
			rsa, err := crypto.RandomRsaKey()
			Expect(err).ShouldNot(HaveOccurred())
			pk, err := crypto.BytesFromRsaPublicKey(&rsa.PublicKey)
			if err != nil {
				log.Fatalf("could not determine public key: %v", err)
			}
			err = prov.WritePublicKey(pk)
			Expect(err).ShouldNot(HaveOccurred())
			readPublicKey, err := prov.PublicKey()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(pk).Should(Equal(readPublicKey))
		})

		It("should return the correct number of peers", func() {
			// shoud have zero by default
			peers, err := prov.Peers()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(peers).Should(Equal(0))

			// should return 1 after adding a peer
			multiAddr, err := testutils.RandomMultiAddress()
			_dht.UpdateMultiAddress(multiAddr)
			Expect(err).ShouldNot(HaveOccurred())
			peers, err = prov.Peers()
			Expect(peers).Should(Equal(1))

			// should return 2 after adding another peer
			multiAddr, err = testutils.RandomMultiAddress()
			_dht.UpdateMultiAddress(multiAddr)
			Expect(err).ShouldNot(HaveOccurred())
			peers, err = prov.Peers()
			Expect(peers).Should(Equal(2))

			// should return 1 after removing a peer
			_dht.RemoveMultiAddress(multiAddr)
			Expect(err).ShouldNot(HaveOccurred())
			peers, err = prov.Peers()
			Expect(peers).Should(Equal(1))
		})
	})

})
