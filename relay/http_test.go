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
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/relay"
	"github.com/republicprotocol/republic-go/stackint"
)

var _ = Describe("HTTP handlers", func() {
	var err error
	epochDNR, err = dnr.TestnetDNR(nil)
	if err != nil {
		panic(err)
	}
	Context("when posting orders", func() {

		It("should return 400 for empty request bodies", func() {
			pools, trader := getPoolsAndTrader()

			r := httptest.NewRequest("POST", "http://localhost/orders", nil)
			w := httptest.NewRecorder()

			handler := relay.RecoveryHandler(relay.OpenOrdersHandler(trader, pools))
			handler.ServeHTTP(w, r)

			Ω(w.Code).Should(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("cannot decode json into an order or a list of order fragments:"))
		})

		It("should return 201 for full orders", func() {
			pools, trader := getPoolsAndTrader()

			fullOrder := getFullOrder()

			sendOrder := relay.OpenOrderRequest{}
			sendOrder.Order = fullOrder
			sendOrder.OrderFragments = relay.OrderFragments{}

			s, _ := json.Marshal(sendOrder)
			body := bytes.NewBuffer(s)
			r := httptest.NewRequest("POST", "http://localhost/orders", body)
			w := httptest.NewRecorder()

			handler := relay.RecoveryHandler(relay.OpenOrdersHandler(trader, pools))
			handler.ServeHTTP(w, r)

			Ω(w.Code).Should(Equal(http.StatusCreated))
		})

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
			pools, trader := getPoolsAndTrader()

			incorrectOrder := []byte("this is not an order or an order fragment")
			s, _ := json.Marshal(incorrectOrder)
			body := bytes.NewBuffer(s)

			r := httptest.NewRequest("POST", "http://localhost/orders", body)
			w := httptest.NewRecorder()

			handler := relay.RecoveryHandler(relay.OpenOrdersHandler(trader, pools))
			handler.ServeHTTP(w, r)

			Ω(w.Code).Should(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("cannot decode json into an order or a list of order fragments:"))
		})

		It("should return 400 for empty order constructs", func() {
			pools, trader := getPoolsAndTrader()

			sendOrder := relay.OpenOrderRequest{}
			sendOrder.Order = order.Order{}
			sendOrder.OrderFragments = relay.OrderFragments{}

			s, err := json.Marshal(sendOrder)
			Ω(err).ShouldNot(HaveOccurred())
			body := bytes.NewBuffer(s)

			r := httptest.NewRequest("POST", "http://localhost/orders", body)
			w := httptest.NewRecorder()

			handler := relay.RecoveryHandler(relay.OpenOrdersHandler(trader, pools))
			handler.ServeHTTP(w, r)

			Ω(w.Code).Should(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("cannot decode json into an order or a list of order fragments: empty object"))
		})
	})

	Context("when getting orders", func() {
		It("should return the correct information when given a valid ID", func() {
			maxConnections := 3
			orderBook := orderbook.NewOrderBook(maxConnections)

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

			var hash [32]byte
			orderMessage := orderbook.NewMessage(ord, order.Open, hash)
			orderBook.Open(orderMessage)

			r := httptest.NewRequest("GET", "http://localhost/orders/vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQ", nil)
			w := httptest.NewRecorder()

			handler := relay.RecoveryHandler(relay.GetOrderHandler(orderBook, string(ord.ID)))
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
			orderBook := orderbook.NewOrderBook(maxConnections)

			r := httptest.NewRequest("GET", "http://localhost/orders/test", nil)
			w := httptest.NewRecorder()

			handler := relay.RecoveryHandler(relay.GetOrderHandler(orderBook, ""))
			handler.ServeHTTP(w, r)

			Expect(w.Body.String()).To(ContainSubstring("order id is invalid"))
			Ω(w.Code).Should(Equal(http.StatusBadRequest))
		})
	})

	Context("when cancelling orders", func() {
		It("should return 410 for cancel order requests", func() {
			pools, trader := getPoolsAndTrader()

			cancelRequest := relay.CancelOrderRequest{}
			cancelRequest.ID = []byte("vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQyV6ryi1wDSM=")

			s, _ := json.Marshal(cancelRequest)
			body := bytes.NewBuffer(s)

			r := httptest.NewRequest("POST", "http://localhost/orders", body)
			w := httptest.NewRecorder()

			handler := relay.RecoveryHandler(relay.CancelOrderHandler(trader, pools))
			handler.ServeHTTP(w, r)

			Ω(w.Code).Should(Equal(http.StatusGone))
		})

		It("should return 400 for malformed cancel order requests", func() {
			pools, trader := getPoolsAndTrader()

			cancelRequest := []byte("this is not an order or an order fragment")

			s, _ := json.Marshal(cancelRequest)
			body := bytes.NewBuffer(s)

			r := httptest.NewRequest("POST", "http://localhost/orders/23213", body)
			w := httptest.NewRecorder()

			handler := relay.RecoveryHandler(relay.CancelOrderHandler(trader, pools))
			handler.ServeHTTP(w, r)

			Ω(w.Code).Should(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("cannot decode json:"))
		})
	})
})
