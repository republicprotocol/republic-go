package rpc_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-order-compute"
	"github.com/republicprotocol/go-rpc"
	"github.com/republicprotocol/go-sss"
	"math/big"
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
			Ω(*serializedAddress).Should(Equal(rpc.Address{Address: address.String()}))

			newAddress, err := rpc.DeserializeAddress(serializedAddress)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(newAddress).Should(Equal(address))
		})

		It("should be deserialize the rpc.Address to an identity.address", func() {
			address := &rpc.Address{Address: keyPair.Address().String()}
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
			Ω(*serializedMulti).Should(Equal(rpc.MultiAddress{Multi: multiAddress.String()}))

			newMultiAddress, err := rpc.DeserializeMultiAddress(serializedMulti)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(newMultiAddress).Should(Equal(multiAddress))
		})

		It("should be deserialize the rpc.Address to an identity.address", func() {
			rpcMultiAddress := &rpc.MultiAddress{Multi: multiAddressString}
			deserializedMulti, err := rpc.DeserializeMultiAddress(rpcMultiAddress)
			Ω(err).ShouldNot(HaveOccurred())

			newMultiAddress := rpc.SerializeMultiAddress(deserializedMulti)
			Ω(newMultiAddress).Should(Equal(rpcMultiAddress))
		})
	})

	Context("serialize and deserialize multiaddresses ", func() {
		It("should be able to serialize an identity.multiaddresses to an rpc.Multiaddresses", func() {
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

		It("should return an error when deserialize a wrong multiaddresses", func() {
			wrongMultiaddress := "/ip4/192.168.0.1/"
			wrongMultiaddresses := rpc.MultiAddresses{Multis: []*rpc.MultiAddress{{wrongMultiaddress}}}
			_, err := rpc.DeserializeMultiAddresses(&wrongMultiaddresses)
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("serialize and deserialize order fragment ", func() {
		It("should be able to serialize and deserialize between order fragment and rpc.orderFragment", func() {
			sssShare := sss.Share{Key: 1, Value: &big.Int{}}
			orderFragment := compute.NewOrderFragment([]byte("orderID"), compute.OrderTypeIBBO, compute.OrderParityBuy,
				sssShare, sssShare, sssShare, sssShare, sssShare)
			rpcOrderFragment := rpc.SerializeOrderFragment(orderFragment)
			newOrderFragment, err := rpc.DeserializeOrderFragment(rpcOrderFragment)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(*orderFragment).Should(Equal(*newOrderFragment))
		})

		It("should return an error when deserializing a order fragment with wrong first code share", func() {
			sssShare := sss.Share{Key: 1, Value: &big.Int{}}
			wrongOrderFragment := rpc.OrderFragment{
				To:             &rpc.Address{},
				From:           &rpc.MultiAddress{},
				Id:             []byte("id"),
				OrderId:        []byte("orderID"),
				OrderType:      1,
				OrderParity:    1,
				FstCodeShare:   []byte(""),
				SndCodeShare:   sss.ToBytes(sssShare),
				PriceShare:     sss.ToBytes(sssShare),
				MaxVolumeShare: sss.ToBytes(sssShare),
				MinVolumeShare: sss.ToBytes(sssShare),
			}
			_, err := rpc.DeserializeOrderFragment(&wrongOrderFragment)
			Ω(err).Should(HaveOccurred())
		})

		It("should return an error when deserializing a order fragment with wrong second code share", func() {
			sssShare := sss.Share{Key: 1, Value: &big.Int{}}
			wrongOrderFragment := rpc.OrderFragment{
				To:             &rpc.Address{},
				From:           &rpc.MultiAddress{},
				Id:             []byte("id"),
				OrderId:        []byte("orderID"),
				OrderType:      1,
				OrderParity:    1,
				FstCodeShare:   sss.ToBytes(sssShare),
				SndCodeShare:   []byte(""),
				PriceShare:     sss.ToBytes(sssShare),
				MaxVolumeShare: sss.ToBytes(sssShare),
				MinVolumeShare: sss.ToBytes(sssShare),
			}
			_, err := rpc.DeserializeOrderFragment(&wrongOrderFragment)
			Ω(err).Should(HaveOccurred())
		})

		It("should return an error when deserializing a order fragment with wrong price share", func() {
			sssShare := sss.Share{Key: 1, Value: &big.Int{}}
			wrongOrderFragment := rpc.OrderFragment{
				To:             &rpc.Address{},
				From:           &rpc.MultiAddress{},
				Id:             []byte("id"),
				OrderId:        []byte("orderID"),
				OrderType:      1,
				OrderParity:    1,
				FstCodeShare:   sss.ToBytes(sssShare),
				SndCodeShare:   sss.ToBytes(sssShare),
				PriceShare:     []byte(""),
				MaxVolumeShare: sss.ToBytes(sssShare),
				MinVolumeShare: sss.ToBytes(sssShare),
			}
			_, err := rpc.DeserializeOrderFragment(&wrongOrderFragment)
			Ω(err).Should(HaveOccurred())
		})

		It("should return an error when deserializing a order fragment with wrong max volume share", func() {
			sssShare := sss.Share{Key: 1, Value: &big.Int{}}
			wrongOrderFragment := rpc.OrderFragment{
				To:             &rpc.Address{},
				From:           &rpc.MultiAddress{},
				Id:             []byte("id"),
				OrderId:        []byte("orderID"),
				OrderType:      1,
				OrderParity:    1,
				FstCodeShare:   sss.ToBytes(sssShare),
				SndCodeShare:   sss.ToBytes(sssShare),
				PriceShare:     sss.ToBytes(sssShare),
				MaxVolumeShare: []byte(""),
				MinVolumeShare: sss.ToBytes(sssShare),
			}
			_, err := rpc.DeserializeOrderFragment(&wrongOrderFragment)
			Ω(err).Should(HaveOccurred())
		})

		It("should return an error when deserializing a order fragment with wrong min volume share", func() {
			sssShare := sss.Share{Key: 1, Value: &big.Int{}}
			wrongOrderFragment := rpc.OrderFragment{
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
				MinVolumeShare: []byte(""),
			}
			_, err := rpc.DeserializeOrderFragment(&wrongOrderFragment)
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("serialize and deserialize result fragment ", func() {
		It("should be able to serialize and deserialize between result fragment and rpc.orderFragment", func() {
			sssShare := sss.Share{Key: 1, Value: &big.Int{}}
			resultFragment := compute.NewResultFragment([]byte("butOrderID"), []byte("sellOrderID"),
				[]byte("butOrderFragmentID"), []byte("sellOrderFragmentID"),
				sssShare, sssShare, sssShare, sssShare, sssShare)
			rpcResultFragment := rpc.SerializeResultFragment(resultFragment)
			newResultFragment, err := rpc.DeserializeResultFragment(rpcResultFragment)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(*resultFragment).Should(Equal(*newResultFragment))
		})

		It("should return an error when deserializing a result fragment with wrong first code share", func() {
			sssShare := sss.Share{Key: 1, Value: &big.Int{}}
			wrongResultFragment := rpc.ResultFragment{
				To:                  &rpc.Address{},
				From:                &rpc.MultiAddress{},
				Id:                  []byte("id"),
				BuyOrderId:          []byte("buyOrderID"),
				SellOrderId:         []byte("sellOrderID"),
				BuyOrderFragmentId:  []byte("buyOrderFragmentID"),
				SellOrderFragmentId: []byte("sellOrderFragmentID"),
				FstCodeShare:        []byte(""),
				SndCodeShare:        sss.ToBytes(sssShare),
				PriceShare:          sss.ToBytes(sssShare),
				MaxVolumeShare:      sss.ToBytes(sssShare),
				MinVolumeShare:      sss.ToBytes(sssShare),
			}
			_, err := rpc.DeserializeResultFragment(&wrongResultFragment)
			Ω(err).Should(HaveOccurred())
		})

		It("should return an error when deserializing a result fragment with wrong second code share", func() {
			sssShare := sss.Share{Key: 1, Value: &big.Int{}}
			wrongResultFragment := rpc.ResultFragment{
				To:                  &rpc.Address{},
				From:                &rpc.MultiAddress{},
				Id:                  []byte("id"),
				BuyOrderId:          []byte("buyOrderID"),
				SellOrderId:         []byte("sellOrderID"),
				BuyOrderFragmentId:  []byte("buyOrderFragmentID"),
				SellOrderFragmentId: []byte("sellOrderFragmentID"),
				FstCodeShare:        sss.ToBytes(sssShare),
				SndCodeShare:        []byte(""),
				PriceShare:          sss.ToBytes(sssShare),
				MaxVolumeShare:      sss.ToBytes(sssShare),
				MinVolumeShare:      sss.ToBytes(sssShare),
			}
			_, err := rpc.DeserializeResultFragment(&wrongResultFragment)
			Ω(err).Should(HaveOccurred())
		})

		It("should return an error when deserializing a result fragment with wrong price share", func() {
			sssShare := sss.Share{Key: 1, Value: &big.Int{}}
			wrongResultFragment := rpc.ResultFragment{
				To:                  &rpc.Address{},
				From:                &rpc.MultiAddress{},
				Id:                  []byte("id"),
				BuyOrderId:          []byte("buyOrderID"),
				SellOrderId:         []byte("sellOrderID"),
				BuyOrderFragmentId:  []byte("buyOrderFragmentID"),
				SellOrderFragmentId: []byte("sellOrderFragmentID"),
				FstCodeShare:        sss.ToBytes(sssShare),
				SndCodeShare:        sss.ToBytes(sssShare),
				PriceShare:          []byte(""),
				MaxVolumeShare:      sss.ToBytes(sssShare),
				MinVolumeShare:      sss.ToBytes(sssShare),
			}
			_, err := rpc.DeserializeResultFragment(&wrongResultFragment)
			Ω(err).Should(HaveOccurred())
		})
		It("should return an error when deserializing a result fragment with wrong max volume share", func() {
			sssShare := sss.Share{Key: 1, Value: &big.Int{}}
			wrongResultFragment := rpc.ResultFragment{
				To:                  &rpc.Address{},
				From:                &rpc.MultiAddress{},
				Id:                  []byte("id"),
				BuyOrderId:          []byte("buyOrderID"),
				SellOrderId:         []byte("sellOrderID"),
				BuyOrderFragmentId:  []byte("buyOrderFragmentID"),
				SellOrderFragmentId: []byte("sellOrderFragmentID"),
				FstCodeShare:        sss.ToBytes(sssShare),
				SndCodeShare:        sss.ToBytes(sssShare),
				PriceShare:          sss.ToBytes(sssShare),
				MaxVolumeShare:      []byte(""),
				MinVolumeShare:      sss.ToBytes(sssShare),
			}
			_, err := rpc.DeserializeResultFragment(&wrongResultFragment)
			Ω(err).Should(HaveOccurred())
		})
		It("should return an error when deserializing a result fragment with wrong first code share", func() {
			sssShare := sss.Share{Key: 1, Value: &big.Int{}}
			wrongResultFragment := rpc.ResultFragment{
				To:                  &rpc.Address{},
				From:                &rpc.MultiAddress{},
				Id:                  []byte("id"),
				BuyOrderId:          []byte("buyOrderID"),
				SellOrderId:         []byte("sellOrderID"),
				BuyOrderFragmentId:  []byte("buyOrderFragmentID"),
				SellOrderFragmentId: []byte("sellOrderFragmentID"),
				FstCodeShare:        sss.ToBytes(sssShare),
				SndCodeShare:        sss.ToBytes(sssShare),
				PriceShare:          sss.ToBytes(sssShare),
				MaxVolumeShare:      sss.ToBytes(sssShare),
				MinVolumeShare:      []byte(""),
			}
			_, err := rpc.DeserializeResultFragment(&wrongResultFragment)
			Ω(err).Should(HaveOccurred())
		})

	})

})
