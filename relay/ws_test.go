package relay_test

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	. "github.com/republicprotocol/republic-go/relay"
	"github.com/republicprotocol/republic-go/stackint"
)

var _ = Describe("WebSocket streaming", func() {
	Context("when connecting to the socket", func() {
		It("should error for missing parameters", func() {
			orderBook := orderbook.NewOrderBook(100)
			server := httptest.NewServer(RecoveryHandler(GetOrdersHandler(orderBook)))
			u, _ := url.Parse(server.URL)
			u.Scheme = "ws"
			u.Path = "orders"
			// Note that we don't specify any query parameters.

			conn, _, _ := websocket.DefaultDialer.Dial(u.String(), nil)
			messageType, _, err := conn.ReadMessage()

			// In this case when we attempt to read, the server is already closed, so
			// we should have an unexpected close error.
			Ω(messageType).Should(Equal(-1))
			Ω(websocket.IsUnexpectedCloseError(err)).Should(Equal(true))
		})

		It("should be able to successfully connect to the socket with valid parameters", func() {
			orderBook := orderbook.NewOrderBook(100)
			server := httptest.NewServer(RecoveryHandler(GetOrdersHandler(orderBook)))
			u, _ := url.Parse(server.URL)
			u.Scheme = "ws"
			u.Path = "orders"
			u.RawQuery = "id=test"

			conn, _, _ := websocket.DefaultDialer.Dial(u.String(), nil)
			conn.SetReadDeadline(time.Now().Add(time.Second))
			messageType, _, err := conn.ReadMessage()

			// In this case the server is still open when we read, but the deadline
			// times out due to not receiving a message, so we should not not have an
			// unexpected close error.
			Ω(messageType).Should(Equal(-1))
			Ω(websocket.IsUnexpectedCloseError(err)).Should(Equal(false))
		})

		It("should retrieve information about an order", func() {
			var wg sync.WaitGroup

			orderBook := orderbook.NewOrderBook(100)

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

			wg.Add(1)
			go func() {
				defer wg.Done()
				// Open an order with the specified ID.
				orderBook.Open(orderMessage)
			}()

			server := httptest.NewServer(RecoveryHandler(GetOrdersHandler(orderBook)))
			u, _ := url.Parse(server.URL)
			u.Scheme = "ws"
			u.Path = "orders"
			u.RawQuery = "id=vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQ"

			// Connect to the socket.
			conn, _, _ := websocket.DefaultDialer.Dial(u.String(), nil)

			// Check if there is a message to be read.
			_, message, err := conn.ReadMessage()
			socketMessage := new(orderbook.Message)
			if err := json.Unmarshal(message, socketMessage); err != nil {
				fmt.Println(err)
			}
			Ω(orderMessage.Ord.ID).Should(Equal(socketMessage.Ord.ID))

			// Update the status of the order and check if there is another
			// message to be read.
			orderBook.Settle(orderMessage)
			messageType, message, err := conn.ReadMessage()
			fmt.Println(string(message))

			// TODO: Check the status of this message.

			Ω(messageType).ShouldNot(Equal(-1))
			Ω(websocket.IsUnexpectedCloseError(err)).Should(Equal(false))

			wg.Wait()
		})

		It("should not retrieve information about unspecified orders", func() {
			var wg sync.WaitGroup

			orderBook := orderbook.NewOrderBook(100)

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

			wg.Add(1)
			go func() {
				defer wg.Done()
				orderBook.Open(orderMessage)
			}()

			server := httptest.NewServer(RecoveryHandler(GetOrdersHandler(orderBook)))
			u, _ := url.Parse(server.URL)
			u.Scheme = "ws"
			u.Path = "orders"
			u.RawQuery = "id=test"

			// Connect to the socket.
			conn, _, _ := websocket.DefaultDialer.Dial(u.String(), nil)

			// Check if there is a message to be read.
			conn.SetReadDeadline(time.Now().Add(time.Second))
			messageType, _, err := conn.ReadMessage()
			Ω(messageType).Should(Equal(-1))
			Ω(websocket.IsUnexpectedCloseError(err)).Should(Equal(false))

			wg.Wait()
		})

		// TODO: Test the handling of statuses.
	})
})
