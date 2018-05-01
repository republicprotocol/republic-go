package relay_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/republicprotocol/republic-go/blockchain/ethereum/dnr"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/relay"
	"github.com/republicprotocol/republic-go/stackint"
)

var _ = PDescribe("HTTP handlers", func() {

	Context("when handling authentication", func() {

		It("should return 401 for unauthorized tokens", func() {

			r := httptest.NewRequest("POST", "http://localhost/orders", nil)
			r.Header.Set("Authorization", "Bearer token test")
			w := httptest.NewRecorder()

			relayNode := relay.Relay{}
			handler := relay.RecoveryHandler(relayNode.AuthorizationHandler(relayNode.OpenOrdersHandler()))
			handler.ServeHTTP(w, r)

			Ω(w.Code).Should(Equal(http.StatusUnauthorized))
		})

		It("should return 401 for authorization headers that are not Bearer type", func() {

			r := httptest.NewRequest("POST", "http://localhost/orders", nil)
			r.Header.Set("Authorization", "Not-Bearer token")
			w := httptest.NewRecorder()

			relayNode := relay.Relay{}
			handler := relay.RecoveryHandler(relayNode.AuthorizationHandler(relayNode.OpenOrdersHandler()))
			handler.ServeHTTP(w, r)

			Ω(w.Code).Should(Equal(http.StatusUnauthorized))
		})

		It("should return 401 for requests without headers", func() {

			r := httptest.NewRequest("POST", "http://localhost/orders", nil)
			w := httptest.NewRecorder()

			relayNode := relay.Relay{}
			handler := relay.RecoveryHandler(relayNode.AuthorizationHandler(relayNode.OpenOrdersHandler()))
			handler.ServeHTTP(w, r)

			Ω(w.Code).Should(Equal(http.StatusUnauthorized))
		})

		// It("should return 200 for authorized tokens", func() {

		// 	r := httptest.NewRequest("POST", "http://localhost/orders", nil)
		// 	r.Header.Set("Authorization", "Bearer token")
		// 	w := httptest.NewRecorder()

		// 	handler := relay.RecoveryHandler(relay.AuthorizationHandler(relay.OpenOrdersHandler((relay.Relay{})), "token"))
		// 	handler.ServeHTTP(w, r)

		// 	Ω(w.Code).Should(Equal(http.StatusOK))
		// })
	})

	Context("when posting orders", func() {

		It("should return 400 for empty request bodies", func() {
			// pools, trader := getPoolsAndTrader()

			r := httptest.NewRequest("POST", "http://localhost/orders", nil)
			w := httptest.NewRecorder()

			relayNode := relay.Relay{}
			handler := relay.RecoveryHandler(relayNode.OpenOrdersHandler())
			handler.ServeHTTP(w, r)

			Ω(w.Code).Should(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("cannot decode json into an order or a list of order fragments:"))
		})

		// It("should return 201 for full orders", func() {
		// 	// pools, trader := getPoolsAndTrader()

		// 	fullOrder := getFullOrder()

		// 	sendOrder := relay.OpenOrderRequest{}
		// 	sendOrder.Order = fullOrder
		// 	sendOrder.OrderFragments = relay.OrderFragments{}

		// 	s, _ := json.Marshal(sendOrder)
		// 	body := bytes.NewBuffer(s)
		// 	r := httptest.NewRequest("POST", "http://localhost/orders", body)
		// 	w := httptest.NewRecorder()

		// 	relayNode := relay.Relay{}
		// 	handler := relay.RecoveryHandler(relay.OpenOrdersHandler(relayNode))
		// 	handler.ServeHTTP(w, r)

		// 	Ω(w.Code).Should(Equal(http.StatusCreated))
		// })

		// It("should return 201 for fragmented orders", func() {
		// 	pools, trader := getPoolsAndTrader()

		// 	fragmentedOrder, err := generateFragmentedOrderForDarkPool(pools[0])
		// 	Ω(err).ShouldNot(HaveOccurred())

		// 	sendOrder := relay.OpenOrderRequest{}
		// 	sendOrder.Order = order.Order{}
		// 	sendOrder.OrderFragments = fragmentedOrder

		// 	data, err := json.Marshal(sendOrder)
		// 	Ω(err).ShouldNot(HaveOccurred())
		// 	body := bytes.NewBuffer(data)

		// 	r := httptest.NewRequest("POST", "http://localhost/orders", body)
		// 	w := httptest.NewRecorder()

		// 	handler := relay.RecoveryHandler(relay.OpenOrdersHandler(trader, pools))
		// 	handler.ServeHTTP(w, r)

		// 	Ω(w.Code).Should(Equal(http.StatusCreated))
		// })

		It("should return 400 for malformed orders", func() {
			// pools, trader := getPoolsAndTrader()

			incorrectOrder := []byte("this is not an order or an order fragment")
			s, _ := json.Marshal(incorrectOrder)
			body := bytes.NewBuffer(s)

			r := httptest.NewRequest("POST", "http://localhost/orders", body)
			w := httptest.NewRecorder()

			relayNode := relay.Relay{}
			handler := relay.RecoveryHandler(relayNode.OpenOrdersHandler())
			handler.ServeHTTP(w, r)

			Ω(w.Code).Should(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("cannot decode json into an order or a list of order fragments:"))
		})

		It("should return 400 for empty order constructs", func() {
			// pools, trader := getPoolsAndTrader()

			sendOrder := relay.OpenOrderRequest{}
			sendOrder.Order = order.Order{}
			sendOrder.OrderFragments = relay.OrderFragments{}

			s, err := json.Marshal(sendOrder)
			Ω(err).ShouldNot(HaveOccurred())
			body := bytes.NewBuffer(s)

			r := httptest.NewRequest("POST", "http://localhost/orders", body)
			w := httptest.NewRecorder()

			relayNode := relay.Relay{}
			handler := relay.RecoveryHandler(relayNode.OpenOrdersHandler())
			handler.ServeHTTP(w, r)

			Ω(w.Code).Should(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("cannot decode json into an order or a list of order fragments: empty object"))
		})
	})

	Context("when getting orders", func() {
		It("should return the correct information when given a valid ID", func() {
			maxConnections := 3
			book := orderbook.NewOrderbook(maxConnections)

			defaultStackVal, _ := stackint.FromString("179761232312312")
			ord := order.Order{}
			ord.ID = []byte("vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQ")
			ord.Type = 2
			ord.Parity = 1
			ord.Expiry = time.Time{}
			ord.FstCode = order.CurrencyCodeETH
			ord.SndCode = order.CurrencyCodeBTC
			ord.Price = defaultStackVal
			ord.MaxVolume = defaultStackVal
			ord.MinVolume = defaultStackVal
			ord.Nonce = defaultStackVal

			orderMessage := orderbook.NewEntry(ord, order.Open)
			book.Open(orderMessage)

			r := httptest.NewRequest("GET", "http://localhost/orders/vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQ", nil)
			w := httptest.NewRecorder()

			relayNode := relay.NewRelay(relay.Config{}, dnr.DarknodeRegistry{}, &book, nil, nil, nil)
			// string(ord.ID)
			handler := relay.RecoveryHandler(relayNode.GetOrderHandler())
			handler.ServeHTTP(w, r)

			message := new(order.Order)
			if err := json.Unmarshal(w.Body.Bytes(), message); err != nil {
				fmt.Println(err)
			}
			Ω(message.ID).Should(Equal(ord.ID))
			Ω(w.Code).Should(Equal(http.StatusOK))
		})

		It("should error when when given an invalid ID", func() {
			maxConnections := 3
			book := orderbook.NewOrderbook(maxConnections)

			r := httptest.NewRequest("GET", "http://localhost/orders/test", nil)
			w := httptest.NewRecorder()

			relayNode := relay.NewRelay(relay.Config{}, dnr.DarknodeRegistry{}, &book, nil, nil, nil)
			handler := relay.RecoveryHandler(relayNode.GetOrderHandler())
			handler.ServeHTTP(w, r)

			Expect(w.Body.String()).To(ContainSubstring("order id is invalid"))
			Ω(w.Code).Should(Equal(http.StatusBadRequest))
		})
	})

	Context("when cancelling orders", func() {
		// It("should return 410 for cancel order requests", func() {
		// 	// pools, trader := getPoolsAndTrader()

		// 	cancelRequest := relay.CancelOrderRequest{}
		// 	cancelRequest.ID = []byte("vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQyV6ryi1wDSM=")

		// 	s, _ := json.Marshal(cancelRequest)
		// 	body := bytes.NewBuffer(s)

		// 	r := httptest.NewRequest("POST", "http://localhost/orders", body)
		// 	w := httptest.NewRecorder()

		// 	relayNode := relay.Relay{}
		// 	handler := relay.RecoveryHandler(relay.CancelOrderHandler(relayNode))
		// 	handler.ServeHTTP(w, r)

		// 	Ω(w.Code).Should(Equal(http.StatusGone))
		// })

		It("should return 400 for malformed cancel order requests", func() {
			// pools, trader := getPoolsAndTrader()

			cancelRequest := []byte("this is not an order or an order fragment")

			s, _ := json.Marshal(cancelRequest)
			body := bytes.NewBuffer(s)

			r := httptest.NewRequest("POST", "http://localhost/orders/23213", body)
			w := httptest.NewRecorder()

			relayNode := relay.Relay{}
			handler := relay.RecoveryHandler(relayNode.CancelOrderHandler())
			handler.ServeHTTP(w, r)

			Ω(w.Code).Should(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("cannot decode json:"))
		})
	})
})
