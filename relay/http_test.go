package relay_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/contracts/dnr"
	"github.com/republicprotocol/republic-go/dark"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/relay"
	"github.com/republicprotocol/republic-go/stackint"
)

// GetPools return dark pools from a mock dnr
func GetPools() dark.Pools {
	log, err := logger.NewLogger(logger.Options{})
	if err != nil {
		panic(fmt.Sprintf("cannot get logger: %v", err))
	}

	dnr, err := dnr.TestnetDNR(nil)
	if err != nil {
		panic(fmt.Sprintf("cannot get mock DNR: %v", err))
	}

	ocean, err := dark.NewOcean(log, 5, dnr)
	if err != nil {
		panic(fmt.Sprintf("cannot get dark ocean: %v", err))
	}
	return ocean.GetPools()
}

var _ = Describe("HTTP handlers", func() {
	Context("when posting orders", func() {

		It("should return 500 for empty request bodies", func() {
			pools := GetPools()

			r := httptest.NewRequest("POST", "http://localhost/orders", nil)
			w := httptest.NewRecorder()

			trader, err := identity.NewMultiAddressFromString("/ip4/127.0.0.1/tcp/80/republic/8MGfbzAMS59Gb4cSjpm34soGNYsM2f")
			Ω(err).ShouldNot(HaveOccurred())

			handler := relay.RecoveryHandler(relay.PostOrdersHandler(&trader, pools))
			handler.ServeHTTP(w, r)

			Ω(w.Code).Should(Equal(http.StatusInternalServerError))
			Expect(w.Body.String()).To(Equal("cannot decode json into an order or a list of order fragments: EOF EOF"))
		})

		It("should return 200 for full orders", func() {

			pools := GetPools()

			fullOrder := order.Order{}

			defaultStackZero := stackint.Zero()

			fullOrder.ID = []byte("vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQyV6ryi1wDSM=")
			fullOrder.Type = 2
			fullOrder.Parity = 1
			fullOrder.Expiry = time.Time{}
			fullOrder.FstCode = order.CurrencyCodeETH
			fullOrder.SndCode = order.CurrencyCodeBTC
			fullOrder.Price = &defaultStackZero
			fullOrder.MaxVolume = &defaultStackZero
			fullOrder.MinVolume = &defaultStackZero
			fullOrder.Nonce = &defaultStackZero

			s, _ := json.Marshal(fullOrder)
			body := bytes.NewBuffer(s)
			r := httptest.NewRequest("POST", "http://localhost/orders", body)
			w := httptest.NewRecorder()

			trader, err := identity.NewMultiAddressFromString("/ip4/127.0.0.1/tcp/80/republic/8MGfbzAMS59Gb4cSjpm34soGNYsM2f")
			Ω(err).ShouldNot(HaveOccurred())

			handler := relay.RecoveryHandler(relay.PostOrdersHandler(&trader, pools))
			handler.ServeHTTP(w, r)

			// Check the status code is what we expect.
			Ω(w.Code).Should(Equal(http.StatusOK))
		})

		It("should return 200 for fragmented orders", func() {

			pools := GetPools()

			defaultStackZero := stackint.Zero()
			prime, _ := stackint.FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111")

			fragmentedOrder := relay.OrderFragments{}
			fragmentSet := make([]relay.Fragments, 1)
			fragments := []*order.Fragment{}

			var err error
			fragments, err = order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, &defaultStackZero, &defaultStackZero, &defaultStackZero, &defaultStackZero).Split(2, 1, &prime)
			Ω(err).ShouldNot(HaveOccurred())

			fragmentSet[0].Fragment = fragments
			fragmentSet[0].DarkPool = []byte("vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQyV6ryi1wDSM=")

			fragmentedOrder.ID = []byte("vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQyV6ryi1wDSM=")
			fragmentedOrder.Type = 2
			fragmentedOrder.Parity = 1
			fragmentedOrder.Expiry = time.Time{}
			fragmentedOrder.FragmentSet = fragmentSet
			s, _ := json.Marshal(fragmentedOrder)
			body := bytes.NewBuffer(s)

			r := httptest.NewRequest("POST", "http://localhost/orders", body)
			w := httptest.NewRecorder()

			trader, err := identity.NewMultiAddressFromString("/ip4/127.0.0.1/tcp/80/republic/8MGfbzAMS59Gb4cSjpm34soGNYsM2f")
			Ω(err).ShouldNot(HaveOccurred())

			handler := relay.RecoveryHandler(relay.PostOrdersHandler(&trader, pools))
			handler.ServeHTTP(w, r)

			Ω(w.Code).Should(Equal(http.StatusOK))
		})

		It("should return 500 for malformed orders", func() {

			pools := GetPools()

			fullOrder := []byte("this is not an order or an order fragment")
			s, _ := json.Marshal(fullOrder)
			body := bytes.NewBuffer(s)

			r := httptest.NewRequest("POST", "http://localhost/orders", body)
			w := httptest.NewRecorder()

			trader, err := identity.NewMultiAddressFromString("/ip4/127.0.0.1/tcp/80/republic/8MGfbzAMS59Gb4cSjpm34soGNYsM2f")
			Ω(err).ShouldNot(HaveOccurred())

			handler := relay.RecoveryHandler(relay.PostOrdersHandler(&trader, pools))
			handler.ServeHTTP(w, r)

			Ω(w.Code).Should(Equal(http.StatusInternalServerError))
			Expect(w.Body.String()).To(Equal("cannot decode json into an order or a list of order fragments: json: cannot unmarshal string into Go value of type order.Order EOF"))
		})
	})

	Context("when cancelling orders", func() {
		It("should return 200 for cancel order requests", func() {

			pools := GetPools()
			cancelRequest := relay.HTTPDelete{}
			cancelRequest.ID = []byte("vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQyV6ryi1wDSM=")
			s, _ := json.Marshal(cancelRequest)
			body := bytes.NewBuffer(s)

			trader, err := identity.NewMultiAddressFromString("/ip4/127.0.0.1/tcp/80/republic/8MGfbzAMS59Gb4cSjpm34soGNYsM2f")
			Ω(err).ShouldNot(HaveOccurred())
			
			r := httptest.NewRequest("POST", "http://localhost/orders", body)
			w := httptest.NewRecorder()

			handler := relay.RecoveryHandler(relay.DeleteOrderHandler(&trader, pools))
			handler.ServeHTTP(w, r)

			Ω(w.Code).Should(Equal(http.StatusOK))
			Expect(w.Body.String()).To(Equal(""))
		})

		It("should return 500 for malformed cancel order requests", func() {

			pools := GetPools()

			cancelRequest := []byte("this is not an order or an order fragment")
			s, _ := json.Marshal(cancelRequest)
			body := bytes.NewBuffer(s)

			trader, err := identity.NewMultiAddressFromString("/ip4/127.0.0.1/tcp/80/republic/8MGfbzAMS59Gb4cSjpm34soGNYsM2f")
			Ω(err).ShouldNot(HaveOccurred())
			
			r := httptest.NewRequest("POST", "http://localhost/orders/{23213}", body)
			w := httptest.NewRecorder()

			handler := relay.RecoveryHandler(relay.DeleteOrderHandler(&trader, pools))
			handler.ServeHTTP(w, r)

			Ω(w.Code).Should(Equal(http.StatusInternalServerError))
			Expect(w.Body.String()).To(Equal("cannot decode json: json: cannot unmarshal string into Go value of type relay.HTTPDelete"))
		})
	})
})
