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
func GetPools(dnr dnr.DarkNodeRegistry) dark.Pools {
	log, err := logger.NewLogger(logger.Options{})
	if err != nil {
		panic(fmt.Sprintf("cannot get logger: %v", err))
	}

	ocean, err := dark.NewOcean(log, 5, dnr)
	if err != nil {
		panic(fmt.Sprintf("cannot get dark ocean: %v", err))
	}
	return ocean.GetPools()
}

func GetFullOrder() order.Order {
	fullOrder := order.Order{}

	defaultStackVal, _ := stackint.FromString("179761232312312")

	fullOrder.ID = []byte("vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQyV6ryi1wDSM=")
	fullOrder.Type = 2
	fullOrder.Parity = 1
	fullOrder.Expiry = time.Time{}
	fullOrder.FstCode = order.CurrencyCodeETH
	fullOrder.SndCode = order.CurrencyCodeBTC
	fullOrder.Price = &defaultStackVal
	fullOrder.MaxVolume = &defaultStackVal
	fullOrder.MinVolume = &defaultStackVal
	fullOrder.Nonce = &defaultStackVal
	return fullOrder
}

func GetFragmentedOrder() relay.Fragments {
	defaultStackVal, _ := stackint.FromString("179761232312312")

	fragmentedOrder := relay.Fragments{}
	fragmentSet := map[string][]*order.Fragment{}
	fragments := []*order.Fragment{}

	var err error
	fragments, err = order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, &defaultStackVal, &defaultStackVal, &defaultStackVal, &defaultStackVal).Split(2, 1, &Prime)
	Ω(err).ShouldNot(HaveOccurred())

	fragmentSet["vrZhWU3VV9LRIM="] = fragments
	// fragmentSet[0].DarkPool = []byte("vrZhWU3VV9LRIM=")

	fragmentedOrder.ID = []byte("vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQyV6ryi1wDSM=")
	fragmentedOrder.Type = 2
	fragmentedOrder.Parity = 1
	fragmentedOrder.Expiry = time.Time{}
	fragmentedOrder.DarkPools = fragmentSet

	return fragmentedOrder
}

var _ = Describe("HTTP handlers", func() {
	Context("when posting orders", func() {

		It("should return 400 for empty request bodies", func() {
			pools := GetPools(epochDNR)

			r := httptest.NewRequest("POST", "http://localhost/orders", nil)
			w := httptest.NewRecorder()

			trader, err := identity.NewMultiAddressFromString("/ip4/127.0.0.1/tcp/80/republic/8MGfbzAMS59Gb4cSjpm34soGNYsM2f")
			Ω(err).ShouldNot(HaveOccurred())

			handler := relay.RecoveryHandler(relay.PostOrdersHandler(&trader, pools))
			handler.ServeHTTP(w, r)

			Ω(w.Code).Should(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("cannot decode json into an order or a list of order fragments:"))
			// log.Println("2.1")
		})

		It("should return 201 for full orders", func() {

			pools := GetPools(epochDNR)

			fullOrder := GetFullOrder()

			sendOrder := relay.HTTPPost{}
			sendOrder.Order = fullOrder
			sendOrder.OrderFragments = relay.Fragments{}

			s, _ := json.Marshal(sendOrder)
			body := bytes.NewBuffer(s)
			r := httptest.NewRequest("POST", "http://localhost/orders", body)
			w := httptest.NewRecorder()

			trader, err := identity.NewMultiAddressFromString("/ip4/127.0.0.1/tcp/80/republic/8MGfbzAMS59Gb4cSjpm34soGNYsM2f")
			Ω(err).ShouldNot(HaveOccurred())

			handler := relay.RecoveryHandler(relay.PostOrdersHandler(&trader, pools))
			handler.ServeHTTP(w, r)

			Ω(w.Code).Should(Equal(http.StatusCreated))
			// log.Println("2.2")
		})

		It("should return 201 for fragmented orders", func() {

			pools := GetPools(epochDNR)

			// defaultStackZero := stackint.Zero()

			fragmentedOrder, err := generateFragmentedOrderForDarkPool(pools[0])
			Ω(err).ShouldNot(HaveOccurred())

			sendOrder := relay.HTTPPost{}
			sendOrder.Order = order.Order{}
			sendOrder.OrderFragments = fragmentedOrder

			s, err := json.Marshal(sendOrder)
			Ω(err).ShouldNot(HaveOccurred())
			body := bytes.NewBuffer(s)

			r := httptest.NewRequest("POST", "http://localhost/orders", body)
			w := httptest.NewRecorder()

			trader, err := identity.NewMultiAddressFromString("/ip4/127.0.0.1/tcp/80/republic/8MGfbzAMS59Gb4cSjpm34soGNYsM2f")
			Ω(err).ShouldNot(HaveOccurred())

			handler := relay.RecoveryHandler(relay.PostOrdersHandler(&trader, pools))
			handler.ServeHTTP(w, r)

			Ω(w.Code).Should(Equal(http.StatusCreated))

			// log.Println("2.3")
		})

		It("should return 400 for malformed orders", func() {

			pools := GetPools(epochDNR)

			incorrectOrder := []byte("this is not an order or an order fragment")
			s, _ := json.Marshal(incorrectOrder)
			body := bytes.NewBuffer(s)

			r := httptest.NewRequest("POST", "http://localhost/orders", body)
			w := httptest.NewRecorder()

			trader, err := identity.NewMultiAddressFromString("/ip4/127.0.0.1/tcp/80/republic/8MGfbzAMS59Gb4cSjpm34soGNYsM2f")
			Ω(err).ShouldNot(HaveOccurred())

			handler := relay.RecoveryHandler(relay.PostOrdersHandler(&trader, pools))
			handler.ServeHTTP(w, r)

			Ω(w.Code).Should(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("cannot decode json into an order or a list of order fragments:"))
			// log.Println("2.4")
		})

		It("should return 400 for empty order constructs", func() {

			pools := GetPools(epochDNR)

			sendOrder := relay.HTTPPost{}
			sendOrder.Order = order.Order{}
			sendOrder.OrderFragments = relay.Fragments{}

			s, err := json.Marshal(sendOrder)
			Ω(err).ShouldNot(HaveOccurred())
			body := bytes.NewBuffer(s)

			r := httptest.NewRequest("POST", "http://localhost/orders", body)
			w := httptest.NewRecorder()

			trader, err := identity.NewMultiAddressFromString("/ip4/127.0.0.1/tcp/80/republic/8MGfbzAMS59Gb4cSjpm34soGNYsM2f")
			Ω(err).ShouldNot(HaveOccurred())

			handler := relay.RecoveryHandler(relay.PostOrdersHandler(&trader, pools))
			handler.ServeHTTP(w, r)

			Ω(w.Code).Should(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("cannot decode json into an order or a list of order fragments: empty object"))
			// log.Println("2.5")
		})
	})

	Context("when cancelling orders", func() {
		It("should return 410 for cancel order requests", func() {

			pools := GetPools(epochDNR)
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

			Ω(w.Code).Should(Equal(http.StatusGone))
			// log.Println("2.6")
		})

		It("should return 400 for malformed cancel order requests", func() {

			pools := GetPools(epochDNR)

			cancelRequest := []byte("this is not an order or an order fragment")
			s, _ := json.Marshal(cancelRequest)
			body := bytes.NewBuffer(s)

			trader, err := identity.NewMultiAddressFromString("/ip4/127.0.0.1/tcp/80/republic/8MGfbzAMS59Gb4cSjpm34soGNYsM2f")
			Ω(err).ShouldNot(HaveOccurred())

			r := httptest.NewRequest("POST", "http://localhost/orders/{23213}", body)
			w := httptest.NewRecorder()

			handler := relay.RecoveryHandler(relay.DeleteOrderHandler(&trader, pools))
			handler.ServeHTTP(w, r)

			Ω(w.Code).Should(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("cannot decode json:"))
			// log.Println("2.7")
		})
	})
})
