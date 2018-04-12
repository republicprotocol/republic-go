package node_test

import node "github.com/republicprotocol/republic-go/dark-node"

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Configurations", func() {
	It("should be able to load config from a json file", func() {
		// config, err := node.LoadConfig("./test_configs/test_config.json")
		// Ω(err).Should(BeNil())
		// Ω(config.RepublicKeyPair.PrivateKey.D.Int64()).Should(Equal(int64(8252748218128302958)))
		// Ω(config.EthereumPrivateKey).Should(Equal("3a1076bf45ab87712ad64ccb3b10217737f7faacbf2872e88fdd9a537d8fe266"))
		// Ω(config.MultiAddress.String()).Should(Equal("/ip4/127.0.0.1/tcp/4000/republic/8MKZ8JwCU9m9affPWHZ9rxp2azXNnE"))
		// Ω(config.BootstrapMultiAddresses[0].String()).Should(Equal("/ip4/127.0.0.1/tcp/3000/republic/8MHarRJdvWd7SsTJE8vRVfj2jb5cWS"))
		// Ω(config.BootstrapMultiAddresses[1].String()).Should(Equal("/ip4/127.0.0.1/tcp/3001/republic/8MKDGUTgKtkymyKTH28xeMxiCnJ9xy"))
		// Ω(config.BootstrapMultiAddresses[2].String()).Should(Equal("/ip4/127.0.0.1/tcp/3002/republic/8MGg76n7RfC6tuw23PYf85VFyM8Zto"))
		// Ω(config.BootstrapMultiAddresses[3].String()).Should(Equal("/ip4/127.0.0.1/tcp/3003/republic/8MJ38m8Nzknh3gVj7QiMjuejmHBMSf"))
	})

	Context("negative tests", func() {
		It("should return an error when trying to open an non-existent file", func() {
			_, err := node.LoadConfig("non-existent.json")
			Ω(err).Should(HaveOccurred())
		})

		It("should return an error when loading a wrong-formatted json file", func() {
			_, err := node.LoadConfig("./test_configs/wrong_test_config.json")
			Ω(err).Should(HaveOccurred())
		})
	})
})
