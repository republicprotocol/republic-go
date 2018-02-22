package rpc_test

import (
	"math/big"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-atom"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-order-compute"
	"github.com/republicprotocol/go-rpc"
	"github.com/republicprotocol/go-sss"
)

const defaultTimeout = time.Second

type mockServer struct {
	identity.MultiAddress
}

var _ = Describe("Data serialization and deserialization", func() {
	var keyPair identity.KeyPair
	var multiAddressString string
	var err error

	BeforeEach(func() {
		keyPair, err = identity.NewKeyPair()
		Ω(err).ShouldNot(HaveOccurred())
		multiAddressString = "/ip4/192.168.0.1/tcp/80/republic/8MHzQ7ZQDvvT8Nqo3HLQQDZvfcHJYB"
	})

	Context("identity.address", func() {
		It("should be able to serialize identity.address", func() {
			address := keyPair.Address()
			serializedAddress := rpc.SerializeAddress(address)
			Ω(*serializedAddress).Should(Equal(rpc.Address{Address: address.String()}))

			newAddress := rpc.DeserializeAddress(serializedAddress)
			Ω(newAddress).Should(Equal(address))
		})

		It("should be able to deserialize identity.address", func() {
			address := &rpc.Address{Address: keyPair.Address().String()}
			deserializedAddress := rpc.DeserializeAddress(address)
			Ω(deserializedAddress).Should(Equal(keyPair.Address()))

			newAddress := rpc.SerializeAddress(deserializedAddress)
			Ω(newAddress).Should(Equal(address))
		})
	})

	Context("identity.MultiAddress", func() {
		It("should be able to serialize identity.MultiAddress", func() {
			multiAddress, err := identity.NewMultiAddressFromString(multiAddressString)
			Ω(err).ShouldNot(HaveOccurred())
			serializedMulti := rpc.SerializeMultiAddress(multiAddress)
			Ω(*serializedMulti).Should(Equal(rpc.MultiAddress{Multi: multiAddress.String()}))

			newMultiAddress, err := rpc.DeserializeMultiAddress(serializedMulti)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(newMultiAddress).Should(Equal(multiAddress))
		})

		It("should be able to deserialize identity.MultiAddress", func() {
			rpcMultiAddress := &rpc.MultiAddress{Multi: multiAddressString}
			deserializedMulti, err := rpc.DeserializeMultiAddress(rpcMultiAddress)
			Ω(err).ShouldNot(HaveOccurred())

			newMultiAddress := rpc.SerializeMultiAddress(deserializedMulti)
			Ω(newMultiAddress).Should(Equal(rpcMultiAddress))
		})
	})

	Context("identity.MultiAddresses", func() {
		It("should be able to serialize and deserialize identity.MultiAddresses", func() {
			multiAddress1, err := identity.NewMultiAddressFromString(multiAddressString)
			Ω(err).ShouldNot(HaveOccurred())
			multiAddress2, err := identity.NewMultiAddressFromString(multiAddressString)
			Ω(err).ShouldNot(HaveOccurred())
			multiAddresses := identity.MultiAddresses{multiAddress1, multiAddress2}

			rpcMultiAddresses := rpc.SerializeMultiAddresses(multiAddresses)
			newMultiAddresses, err := rpc.DeserializeMultiAddresses(rpcMultiAddresses)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(multiAddresses).Should(Equal(newMultiAddresses))
		})

		It("should return an error when deserializing a malformed identity.MultiAddresses", func() {
			wrongMultiaddress := "/ip4/192.168.0.1/"
			wrongMultiaddresses := rpc.MultiAddresses{Multis: []*rpc.MultiAddress{{wrongMultiaddress}}}
			_, err := rpc.DeserializeMultiAddresses(&wrongMultiaddresses)
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("compute.OrderFragment", func() {
		var wrongOrderFragment rpc.OrderFragment
		sssShare := sss.Share{Key: 1, Value: &big.Int{}}

		BeforeEach(func() {
			wrongOrderFragment = rpc.OrderFragment{
				To:             &rpc.Address{},
				From:           &rpc.MultiAddress{},
				Id:             []byte("id"),
				OrderId:        []byte("orderID"),
				OrderType:      1,
				OrderParity:    1,
				FstCodeShare:   sss.ToBytes(sssShare),
				SndCodeShare:   sss.ToBytes(sssShare),
				PriceShare:     sss.ToBytes(sssShare),
				MaxVolumeShare: sss.ToBytes(sssShare),
				MinVolumeShare: sss.ToBytes(sssShare),
			}
		})

		It("should be able to serialize and deserialize compute.OrderFragment", func() {
			orderFragment := compute.NewOrderFragment([]byte("orderID"), compute.OrderTypeIBBO, compute.OrderParityBuy,
				sssShare, sssShare, sssShare, sssShare, sssShare)
			rpcOrderFragment := rpc.SerializeOrderFragment(orderFragment)
			newOrderFragment, err := rpc.DeserializeOrderFragment(rpcOrderFragment)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(*orderFragment).Should(Equal(*newOrderFragment))
		})

		It("should return an error when deserializing an compute.OrderFragment with a malformed FstCodeShare", func() {
			wrongOrderFragment.FstCodeShare = []byte("")
			_, err := rpc.DeserializeOrderFragment(&wrongOrderFragment)
			Ω(err).Should(HaveOccurred())
		})

		It("should return an error when deserializing an compute.OrderFragment with a malformed SndCodeShare", func() {
			wrongOrderFragment.SndCodeShare = []byte("")
			_, err := rpc.DeserializeOrderFragment(&wrongOrderFragment)
			Ω(err).Should(HaveOccurred())
		})

		It("should return an error when deserializing an compute.OrderFragment with a malformed PriceShare", func() {
			wrongOrderFragment.PriceShare = []byte("")
			_, err := rpc.DeserializeOrderFragment(&wrongOrderFragment)
			Ω(err).Should(HaveOccurred())
		})

		It("should return an error when deserializing an compute.OrderFragment with a malformed MaxVolumeShare", func() {
			wrongOrderFragment.MaxVolumeShare = []byte("")
			_, err := rpc.DeserializeOrderFragment(&wrongOrderFragment)
			Ω(err).Should(HaveOccurred())
		})

		It("should return an error when deserializing an compute.OrderFragment with a malformed MinVolumeShare", func() {
			wrongOrderFragment.MinVolumeShare = []byte("")
			_, err := rpc.DeserializeOrderFragment(&wrongOrderFragment)
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("compute.Result", func() {
		It("should be able to serialize and deserialize a compute.Result", func() {
			result := &compute.Final{
				ID:          []byte("resultID"),
				BuyOrderID:  []byte("BuyOrderID"),
				SellOrderID: []byte("SellOrderID"),
				FstCode:     big.NewInt(0),
				SndCode:     big.NewInt(0),
				Price:       big.NewInt(0),
				MaxVolume:   big.NewInt(0),
				MinVolume:   big.NewInt(0),
			}
			serializedResult := rpc.SerializeFinal(result)
			newResult := rpc.DeserializeFinal(serializedResult)
			Ω(*result).Should(Equal(*newResult))
		})
	})

	Context("atom.Atom", func() {
		It("should be able to serialize and deserialize atom.Atom", func() {
			a := atom.Atom{
				Ledger:     atom.Ledger(0),
				LedgerData: []byte("data"),
				Signature:  []byte("signature"),
			}
			rpcAtom := rpc.SerializeAtom(a)
			newAtom := rpc.DeserializeAtom(rpcAtom)
			Ω(newAtom).Should(Equal(a))
		})
	})
})
