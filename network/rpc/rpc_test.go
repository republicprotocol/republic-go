package rpc_test

import (
	"time"

	"github.com/republicprotocol/republic-go/order"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/network/rpc"
	"github.com/republicprotocol/republic-go/shamir"
	"github.com/republicprotocol/republic-go/stackint"
)

const defaultTimeout = time.Second

type mockServer struct {
	identity.MultiAddress
}

type mockClient struct {
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
			serializedMulti := rpc.SerializeMultiAddress(multiAddress, nil)
			Ω(*serializedMulti).Should(Equal(rpc.MultiAddress{MultiAddress: multiAddress.String()}))

			newMultiAddress, _, err := rpc.DeserializeMultiAddress(serializedMulti)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(newMultiAddress).Should(Equal(multiAddress))
		})

		It("should be able to deserialize identity.MultiAddress", func() {
			rpcMultiAddress := &rpc.MultiAddress{MultiAddress: multiAddressString}
			deserializedMulti, _, err := rpc.DeserializeMultiAddress(rpcMultiAddress)
			Ω(err).ShouldNot(HaveOccurred())

			newMultiAddress := rpc.SerializeMultiAddress(deserializedMulti, nil)
			Ω(newMultiAddress).Should(Equal(rpcMultiAddress))
		})
	})

	Context("identity.MultiAddresses", func() {
		It("should be able to serialize and deserialize identity.MultiAddresses", func() {
			// multiAddress1, err := identity.NewMultiAddressFromString(multiAddressString)
			// Ω(err).ShouldNot(HaveOccurred())
			// multiAddress2, err := identity.NewMultiAddressFromString(multiAddressString)
			// Ω(err).ShouldNot(HaveOccurred())
			// multiAddresses := identity.MultiAddresses{multiAddress1, multiAddress2}

			// rpcMultiAddresses := rpc.SerializeMultiAddresses(multiAddresses, nil)
			// newMultiAddresses, _, err := rpc.DeserializeMultiAddresses(rpcMultiAddresses)
			// Ω(err).ShouldNot(HaveOccurred())
			// Ω(multiAddresses).Should(Equal(newMultiAddresses))
		})

		It("should return an error when deserializing a malformed identity.MultiAddresses", func() {
			// wrongMultiAddress := "/ip4/192.168.0.1/"
			// wrongMultiAddresses := rpc.MultiAddresses{Multis: []*rpc.MultiAddress{{wrongMultiAddress}}}
			// _, _, err := rpc.DeserializeMultiAddresses(&wrongMultiAddresses)
			// Ω(err).Should(HaveOccurred())
		})
	})

	Context("compute.OrderFragment", func() {
		var wrongOrderFragment rpc.OrderFragment
		shamirShare := shamir.Share{Key: 1, Value: stackint.Zero()}

		BeforeEach(func() {
			wrongOrderFragment = rpc.OrderFragment{
				// To:             &rpc.Address{},
				// From:           &rpc.MultiAddress{},
				Id:             []byte("id"),
				OrderId:        []byte("orderID"),
				OrderType:      1,
				OrderParity:    1,
				FstCodeShare:   shamir.ToBytes(shamirShare),
				SndCodeShare:   shamir.ToBytes(shamirShare),
				PriceShare:     shamir.ToBytes(shamirShare),
				MaxVolumeShare: shamir.ToBytes(shamirShare),
				MinVolumeShare: shamir.ToBytes(shamirShare),
			}
		})

		It("should be able to serialize and deserialize compute.OrderFragment", func() {
			price := stackint.FromUint(10)
			maxVolume := stackint.FromUint(1000)
			minVolume := stackint.FromUint(100)
			nonce := stackint.Zero()

			prime, _ := stackint.FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111")

			fragments, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, &price, &maxVolume, &minVolume, &nonce).Split(2, 1, &prime)
			Ω(err).ShouldNot(HaveOccurred())

			for _, orderFragment := range fragments {
				rpcOrderFragment := rpc.SerializeOrderFragment(orderFragment)
				newOrderFragment, err := rpc.DeserializeOrderFragment(rpcOrderFragment)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(*orderFragment).Should(Equal(*newOrderFragment))
			}
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

	//Context("compute.Result", func() {
	//	It("should be able to serialize and deserialize a compute.Result", func() {
	//		result := &compute.Final{
	//			ID:          []byte("resultID"),
	//			BuyOrderID:  []byte("BuyOrderID"),
	//			SellOrderID: []byte("SellOrderID"),
	//			FstCode:     big.NewInt(0),
	//			SndCode:     big.NewInt(0),
	//			Price:       big.NewInt(0),
	//			MaxVolume:   big.NewInt(0),
	//			MinVolume:   big.NewInt(0),
	//		}
	//		serializedResult := rpc.SerializeFinal(result)
	//		newResult := rpc.DeserializeFinal(serializedResult)
	//		Ω(*result).Should(Equal(*newResult))
	//	})
	//})
	Context("delta fragment", func() {
		It("should be able to serialize and deserialize compute.DeltaFragment", func() {

		})
	})

	Context("atom.Atom", func() {
		It("should be able to serialize and deserialize atom.Atom", func() {
			// a := atom.Atom{
			// 	Ledger:     atom.Ledger(0),
			// 	LedgerData: []byte("data"),
			// 	Signature:  []byte("signature"),
			// }
			// rpcAtom := rpc.SerializeAtom(a)
			// newAtom := rpc.DeserializeAtom(rpcAtom)
			// Ω(newAtom).Should(Equal(a))
		})
	})

})
