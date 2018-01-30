package rpc_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-rpc"
)

var _ = Describe("Data serializing and deserialization", func() {
	var keyPair identity.KeyPair
	var multiAddressString string
	var err error

	BeforeEach(func() {
		keyPair, err = identity.NewKeyPair()
		Ω(err).ShouldNot(HaveOccurred())
		multiAddressString = "/ip4/192.168.0.1/tcp/80/republic/8MHzQ7ZQDvvT8Nqo3HLQQDZvfcHJYB"
	})

	Context("serialize and deserialize address ", func() {
		It("should be able to serialize an identity.address to an rpc.Address", func() {
			address := keyPair.Address()
			serializedAddress := rpc.SerializeAddress(address)
			Ω(*serializedAddress).Should(Equal(rpc.Address{Address:address.String()}))

			newAddress, err:= rpc.DeserializeAddress(serializedAddress)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(newAddress).Should(Equal(address))
		})

		It("should be deserialize the rpc.Address to an identity.address", func() {
			address := &rpc.Address{Address:keyPair.Address().String()}
			deserializedAddress, err := rpc.DeserializeAddress(address)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(deserializedAddress).Should(Equal(keyPair.Address()))

			newAddress := rpc.SerializeAddress(deserializedAddress)
			Ω(newAddress).Should(Equal(address))
		})
	})

	Context("serialize and deserialize multiaddress ", func() {
		It("should be able to serialize an identity.multiaddress to an rpc.Multiaddress", func() {
			multiAddress, err := identity.NewMultiAddressFromString(multiAddressString)
			Ω(err).ShouldNot(HaveOccurred())
			serializedMulti := rpc.SerializeMultiAddress(multiAddress)
			Ω(*serializedMulti).Should(Equal(rpc.MultiAddress{Multi:multiAddress.String()}))

			newMultiAddress, err:= rpc.DeserializeMultiAddress(serializedMulti)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(newMultiAddress).Should(Equal(multiAddress))
		})

		It("should be deserialize the rpc.Address to an identity.address", func() {
			rpcMultiAddress := &rpc.MultiAddress{Multi:multiAddressString}
			deserializedMulti, err:= rpc.DeserializeMultiAddress(rpcMultiAddress)
			Ω(err).ShouldNot(HaveOccurred())

			newMultiAddress:= rpc.SerializeMultiAddress(deserializedMulti)
			Ω(newMultiAddress).Should(Equal(rpcMultiAddress))
		})
	})

	Context("serialize and deserialize multiaddresses ", func() {
		It("should be able to serialize an identity.multiaddresses to an rpc.Multiaddresses", func() {

		})

		It("should be deserialize the rpc.Addresses to an identity.addresses", func() {

		})
	})

	Context("serialize and deserialize order fragment ", func() {
		It("should be able to serialize an compute.orderFragment to an rpc.orderFragment", func() {

		})

		It("should be deserialize the rpc.orderFragment to an compute.orderFragment", func() {

		})
	})

	Context("serialize and deserialize result fragment ", func() {
		It("should be able to serialize an compute.resultFragment to an rpc.resultFragment", func() {

		})

		It("should be deserialize the rpc.resultFragment to an compute.resultFragment", func() {

		})
	})

})