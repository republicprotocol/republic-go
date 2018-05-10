package smpcer_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/rpc/smpcer"

	"github.com/jbenet/go-base58"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/delta"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/shamir"
	"github.com/republicprotocol/republic-go/stackint"
)

var _ = Describe("Marshal", func() {

	Context("MarshalOrderFragment method", func() {

		It(" should return a network representation of the order fragment", func() {
			crypter := crypto.NewWeakCrypter()

			orderFragment, _, err := createNewOrderFragment()
			Expect(err).ShouldNot(HaveOccurred())

			marshalledFragment, err := MarshalOrderFragment("", &crypter, &orderFragment)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(marshalledFragment).ToNot(BeNil())
			Expect(base58.Encode(marshalledFragment.OrderFragmentId)).To(Equal(orderFragment.ID.String()))
		})
	})

	Context("UnmarshalOrderFragment method", func() {

		It(" should convert a network representation of an OrderFragment into an order.Fragment", func() {
			crypter := crypto.NewWeakCrypter()

			orderFragment, _, err := createNewOrderFragment()
			Expect(err).ShouldNot(HaveOccurred())

			marshalledFragment, err := MarshalOrderFragment("", &crypter, &orderFragment)
			Expect(err).ShouldNot(HaveOccurred())

			unmarshalledFragment, err := UnmarshalOrderFragment(&crypter, marshalledFragment)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(unmarshalledFragment).ToNot(BeNil())

			Expect(unmarshalledFragment.ID).To(Equal(orderFragment.ID))
			Expect(unmarshalledFragment.OrderID).To(Equal(orderFragment.OrderID))
			Expect(unmarshalledFragment.OrderType).To(Equal(orderFragment.OrderType))
			Expect(unmarshalledFragment.OrderParity).To(Equal(orderFragment.OrderParity))
			Expect(unmarshalledFragment.OrderExpiry.Unix()).To(Equal(orderFragment.OrderExpiry.Unix()))
			Expect(unmarshalledFragment.FstCodeShare).To(Equal(orderFragment.FstCodeShare))
			Expect(unmarshalledFragment.SndCodeShare).To(Equal(orderFragment.SndCodeShare))
			Expect(unmarshalledFragment.PriceShare).To(Equal(orderFragment.PriceShare))
			Expect(unmarshalledFragment.MaxVolumeShare).To(Equal(orderFragment.MaxVolumeShare))
			Expect(unmarshalledFragment.MinVolumeShare).To(Equal(orderFragment.MinVolumeShare))
		})
	})

	Context("MarshalDeltaFragment method", func() {

		It(" should convert delta.Fragment into a RPC protobuf object", func() {
			deltaFragment, err := createNewDeltaFragment()
			Expect(err).ShouldNot(HaveOccurred())

			marshalledFragment := MarshalDeltaFragment(&deltaFragment)
			Expect(marshalledFragment).ToNot(BeNil())
			Expect(base58.Encode(marshalledFragment.DeltaFragmentId)).To(Equal(deltaFragment.ID.String()))
		})
	})

	Context("UnmarshalDeltaFragment method", func() {

		It(" should convert a RPC protobuf object into an delta.Fragment", func() {
			deltaFragment, err := createNewDeltaFragment()
			Expect(err).ShouldNot(HaveOccurred())

			marshalledFragment := MarshalDeltaFragment(&deltaFragment)

			unmarshalledFragment, err := UnmarshalDeltaFragment(marshalledFragment)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(unmarshalledFragment).ToNot(BeNil())
			Expect(unmarshalledFragment).To(Equal(deltaFragment))
		})
	})
})

// createNewOrderFragment creates a test order fragment
func createNewOrderFragment() (order.Fragment, stackint.Int1024, error) {

	price := 1000000000000
	testOrder := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour),
		order.CurrencyCodeETH, order.CurrencyCodeBTC, stackint.FromUint(uint(price)), stackint.FromUint(uint(price)),
		stackint.FromUint(uint(price)), stackint.FromUint(1))

	prime, err := stackint.FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111")
	if err != nil {
		return order.Fragment{}, stackint.Int1024{}, fmt.Errorf("cannot generate prime: %v", err)
	}

	share := shamir.Share{
		Key:   int64(1),
		Value: prime,
	}

	fragment := order.NewFragment(
		testOrder.ID,
		testOrder.Type,
		testOrder.Parity,
		share,
		share,
		share,
		share,
		share,
	)

	return *fragment, prime, nil
}

// CreateNewDeltaFragment creates a test delta fragment
func createNewDeltaFragment() (delta.Fragment, error) {

	fragment, prime, err := createNewOrderFragment()
	if err != nil {
		return delta.Fragment{}, err
	}

	deltaFragment := delta.NewDeltaFragment(&fragment, &fragment, &prime)

	return deltaFragment, nil
}
