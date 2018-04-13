package relay_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/relay"
)

var _ = Describe("HTTP handlers", func() {
	Context("when posting orders", func() {

		It("should return 400 for empty request bodies", func() {
			pools, trader := getPoolsAndTrader()

			r := httptest.NewRequest("POST", "http://localhost/orders", nil)
			w := httptest.NewRecorder()

			handler := relay.RecoveryHandler(relay.PostOrdersHandler(&trader, pools))
			handler.ServeHTTP(w, r)

			Ω(w.Code).Should(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("cannot decode json into an order or a list of order fragments:"))
		})

		It("should return 201 for full orders", func() {
			pools, trader := getPoolsAndTrader()

			fullOrder := getFullOrder()

			sendOrder := relay.HTTPPost{}
			sendOrder.Order = fullOrder
			sendOrder.OrderFragments = relay.Fragments{}

			s, _ := json.Marshal(sendOrder)
			body := bytes.NewBuffer(s)
			r := httptest.NewRequest("POST", "http://localhost/orders", body)
			w := httptest.NewRecorder()

			handler := relay.RecoveryHandler(relay.PostOrdersHandler(&trader, pools))
			handler.ServeHTTP(w, r)

			Ω(w.Code).Should(Equal(http.StatusCreated))
		})

		It("should return 201 for fragmented orders", func() {
			pools, trader := getPoolsAndTrader()

			fragmentedOrder, err := generateFragmentedOrderForDarkPool(pools[0])
			Ω(err).ShouldNot(HaveOccurred())

			sendOrder := relay.HTTPPost{}
			sendOrder.Order = order.Order{}
			sendOrder.OrderFragments = fragmentedOrder

			data, err := json.Marshal(sendOrder)
			Ω(err).ShouldNot(HaveOccurred())
			body := bytes.NewBuffer(data)

			r := httptest.NewRequest("POST", "http://localhost/orders", body)
			w := httptest.NewRecorder()

			handler := relay.RecoveryHandler(relay.PostOrdersHandler(&trader, pools))
			handler.ServeHTTP(w, r)

			Ω(w.Code).Should(Equal(http.StatusCreated))
		})

		It("should return 400 for malformed orders", func() {
			pools, trader := getPoolsAndTrader()

			incorrectOrder := []byte("this is not an order or an order fragment")
			s, _ := json.Marshal(incorrectOrder)
			body := bytes.NewBuffer(s)

			r := httptest.NewRequest("POST", "http://localhost/orders", body)
			w := httptest.NewRecorder()

			handler := relay.RecoveryHandler(relay.PostOrdersHandler(&trader, pools))
			handler.ServeHTTP(w, r)

			Ω(w.Code).Should(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("cannot decode json into an order or a list of order fragments:"))
		})

		It("should return 400 for empty order constructs", func() {
			pools, trader := getPoolsAndTrader()

			sendOrder := relay.HTTPPost{}
			sendOrder.Order = order.Order{}
			sendOrder.OrderFragments = relay.Fragments{}

			s, err := json.Marshal(sendOrder)
			Ω(err).ShouldNot(HaveOccurred())
			body := bytes.NewBuffer(s)

			r := httptest.NewRequest("POST", "http://localhost/orders", body)
			w := httptest.NewRecorder()

			handler := relay.RecoveryHandler(relay.PostOrdersHandler(&trader, pools))
			handler.ServeHTTP(w, r)

			Ω(w.Code).Should(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("cannot decode json into an order or a list of order fragments: empty object"))
		})
	})

	Context("when cancelling orders", func() {
		It("should return 410 for cancel order requests", func() {
			pools, trader := getPoolsAndTrader()

			cancelRequest := relay.HTTPDelete{}
			cancelRequest.ID = []byte("vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQyV6ryi1wDSM=")

			s, _ := json.Marshal(cancelRequest)
			body := bytes.NewBuffer(s)

			r := httptest.NewRequest("POST", "http://localhost/orders", body)
			w := httptest.NewRecorder()

			handler := relay.RecoveryHandler(relay.DeleteOrderHandler(&trader, pools))
			handler.ServeHTTP(w, r)

			Ω(w.Code).Should(Equal(http.StatusGone))
		})

		It("should return 400 for malformed cancel order requests", func() {
			pools, trader := getPoolsAndTrader()

			cancelRequest := []byte("this is not an order or an order fragment")

			s, _ := json.Marshal(cancelRequest)
			body := bytes.NewBuffer(s)

			r := httptest.NewRequest("POST", "http://localhost/orders/{23213}", body)
			w := httptest.NewRecorder()

			handler := relay.RecoveryHandler(relay.DeleteOrderHandler(&trader, pools))
			handler.ServeHTTP(w, r)

			Ω(w.Code).Should(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("cannot decode json:"))
		})
	})
})
