package status_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/cmd/darknode/config"
	"github.com/republicprotocol/republic-go/dht"
	. "github.com/republicprotocol/republic-go/status"
)

var _ = Describe("Status", func() {
	testStr := "someRandomString"

	Context("when writing to the provider", func() {
		var prov Provider

		BeforeEach(func() {
			config, err := config.NewConfigFromJSONFile("./test_config.json")
			Expect(err).ShouldNot(HaveOccurred())
			dht := dht.NewDHT(config.Address, 64)
			prov = NewProvider(&dht)
		})

		Context("when writing ethereum address", func() {

			It("should not trigger an error", func() {
				err := prov.WriteEthereumAddress(testStr)
				Expect(err).ShouldNot(HaveOccurred())
				ethAddr, err := prov.EthereumAddress()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(ethAddr).Should(Equal(testStr))
			})

		})
	})

})
