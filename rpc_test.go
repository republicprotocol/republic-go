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

	Context("when using an address", func() {
		It("should be able to serialize", func() {
			address := keyPair.Address()
			serializedAddress := rpc.SerializeAddress(address)
			Ω(*serializedAddress).Should(Equal(rpc.Address{Address: address.String()}))

			newAddress := rpc.DeserializeAddress(serializedAddress)
			Ω(newAddress).Should(Equal(address))
		})

		It("should be able to deserialize", func() {
			address := &rpc.Address{Address: keyPair.Address().String()}
			deserializedAddress := rpc.DeserializeAddress(address)
			Ω(deserializedAddress).Should(Equal(keyPair.Address()))

			newAddress := rpc.SerializeAddress(deserializedAddress)
			Ω(newAddress).Should(Equal(address))
		})
	})

	Context("when using a multi-address", func() {
		It("should be able to serialize", func() {
			multiAddress, err := identity.NewMultiAddressFromString(multiAddressString)
			Ω(err).ShouldNot(HaveOccurred())
			serializedMulti := rpc.SerializeMultiAddress(multiAddress)
			Ω(*serializedMulti).Should(Equal(rpc.MultiAddress{Multi: multiAddress.String()}))

			newMultiAddress, err := rpc.DeserializeMultiAddress(serializedMulti)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(newMultiAddress).Should(Equal(multiAddress))
		})

		It("should be able to deserialize", func() {
			rpcMultiAddress := &rpc.MultiAddress{Multi: multiAddressString}
			deserializedMulti, err := rpc.DeserializeMultiAddress(rpcMultiAddress)
			Ω(err).ShouldNot(HaveOccurred())

			newMultiAddress := rpc.SerializeMultiAddress(deserializedMulti)
			Ω(newMultiAddress).Should(Equal(rpcMultiAddress))
		})
	})

	Context("when using multi-addresses", func() {
		It("should be able to serialize and deserialize", func() {
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

		It("should return an error when deserializing malformed multi-addresses", func() {
			wrongMultiaddress := "/ip4/192.168.0.1/"
			wrongMultiaddresses := rpc.MultiAddresses{Multis: []*rpc.MultiAddress{{wrongMultiaddress}}}
			_, err := rpc.DeserializeMultiAddresses(&wrongMultiaddresses)
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("when using an order fragment", func() {

		It("should be able to serialize and deserialize", func() {
			sssShare := sss.Share{Key: 1, Value: &big.Int{}}
			orderFragment := compute.NewOrderFragment([]byte("orderID"), compute.OrderTypeIBBO, compute.OrderParityBuy,
				sssShare, sssShare, sssShare, sssShare, sssShare)
			rpcOrderFragment := rpc.SerializeOrderFragment(orderFragment)
			newOrderFragment, err := rpc.DeserializeOrderFragment(rpcOrderFragment)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(*orderFragment).Should(Equal(*newOrderFragment))
		})

		It("should return an error when deserializing an order fragment with a malformed first code share", func() {
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

		It("should return an error when deserializing an order fragment with a malformed second code share", func() {
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

		It("should return an error when deserializing an order fragment with a malformed price share", func() {
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

		It("should return an error when deserializing an order fragment with a malformed maximum volume share", func() {
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

		It("should return an error when deserializing an order fragment with a malformed minimum volume share", func() {
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

	Context("when using a result fragment", func() {
		It("should be able to serialize and deserialize", func() {
			sssShare := sss.Share{Key: 1, Value: &big.Int{}}
			resultFragment := compute.NewResultFragment([]byte("butOrderID"), []byte("sellOrderID"),
				[]byte("butOrderFragmentID"), []byte("sellOrderFragmentID"),
				sssShare, sssShare, sssShare, sssShare, sssShare)
			rpcResultFragment := rpc.SerializeResultFragment(resultFragment)
			newResultFragment, err := rpc.DeserializeResultFragment(rpcResultFragment)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(*resultFragment).Should(Equal(*newResultFragment))
		})

		It("should return an error when deserializing a result fragment with a malformed first code share", func() {
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

		It("should return an error when deserializing a result fragment with a malformed second code share", func() {
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

		It("should return an error when deserializing a result fragment with a malformed price share", func() {
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

		It("should return an error when deserializing a result fragment with a malformed maximum volume share", func() {
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

		It("should return an error when deserializing a result fragment with a malformed minimum volume share", func() {
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

	Context("when using a trading atom", func() {
		It("should be able to serialize and deserialize", func() {
			a := atom.Atom{
				ID:         []byte{},
				Lock:       []byte{},
				Fst:        atom.LedgerBitcoin,
				FstAddress: "0x",
				Snd:        atom.LedgerEthereum,
				SndAddress: "0x",
			}
			rpcAtom := rpc.SerializeAtom(a)
			newAtom := rpc.DeserializeAtom(rpcAtom)
			Ω(newAtom).Should(Equal(a))
		})
	})
})
