package relay_test

import (
	"encoding/json"
	"net/http/httptest"
	"net/url"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/relay"

	"github.com/gorilla/websocket"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/dnr"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/stackint"
)

var _ = Describe("WebSocket streaming", func() {

	FContext("when connecting to the socket", func() {
		It("should error for missing parameters", func() {
			book := orderbook.NewOrderbook()
			relay := NewRelay(Config{}, dnr.DarknodeRegistry{}, &book, nil, nil, nil)

			server := httptest.NewServer(RecoveryHandler(relay.GetOrdersHandler()))
			u, _ := url.Parse(server.URL)
			u.Scheme = "ws"
			u.Path = "orders"
			// Note that we don't specify any query parameters.

			conn, _, _ := websocket.DefaultDialer.Dial(u.String(), nil)
			messageType, _, err := conn.ReadMessage()

			// In this case when we attempt to read, the server is already closed, so
			// we should have an unexpected close error.
			Ω(messageType).Should(Equal(-1))
			// TODO: Check for websocket error responses.
			Ω(websocket.IsUnexpectedCloseError(err)).Should(Equal(true))
		})

		It("should be able to successfully connect to the socket with valid parameters", func() {
			book := orderbook.NewOrderbook()
			relay := NewRelay(Config{}, dnr.DarknodeRegistry{}, &book, nil, nil, nil)

			server := httptest.NewServer(RecoveryHandler(relay.GetOrdersHandler()))
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

			book := orderbook.NewOrderbook()
			relay := NewRelay(Config{}, dnr.DarknodeRegistry{}, &book, nil, nil, nil)

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

			wg.Add(1)
			go func() {
				defer wg.Done()
				// Open an order with the specified ID.
				book.Open(ord)
			}()

			// Connect to the socket.
			server := httptest.NewServer(RecoveryHandler(relay.GetOrdersHandler()))
			u, _ := url.Parse(server.URL)
			u.Scheme = "ws"
			u.Path = "orders"
			u.RawQuery = "id=vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQ"
			conn, _, _ := websocket.DefaultDialer.Dial(u.String(), nil)

			// We should be able to read the initial message.
			messageType, message, err := conn.ReadMessage()
			var socketMessage orderbook.Entry
			err = json.Unmarshal(message, &socketMessage)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(socketMessage.Order.ID).Should(Equal(ord.ID))
			Ω(messageType).Should(Equal(websocket.TextMessage))

			// Update the status of the order and check if there is another
			// message to be read.
			book.Settle(ord)
			messageType, message, err = conn.ReadMessage()
			err = json.Unmarshal(message, &socketMessage)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(socketMessage.Status).Should(Equal(order.Settled))
			Ω(messageType).Should(Equal(websocket.TextMessage))
			Ω(websocket.IsUnexpectedCloseError(err)).Should(Equal(false))

			wg.Wait()
		})

		It("should not retrieve information about unspecified orders", func() {
			var wg sync.WaitGroup

			book := orderbook.NewOrderbook()
			relay := NewRelay(Config{}, dnr.DarknodeRegistry{}, &book, nil, nil, nil)

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

			wg.Add(1)
			go func() {
				defer wg.Done()
				book.Open(ord)
			}()

			// Connect to the socket.
			server := httptest.NewServer(RecoveryHandler(relay.GetOrdersHandler()))
			u, _ := url.Parse(server.URL)
			u.Scheme = "ws"
			u.Path = "orders"
			u.RawQuery = "id=test"
			conn, _, _ := websocket.DefaultDialer.Dial(u.String(), nil)

			// We should not receive any messages.
			conn.SetReadDeadline(time.Now().Add(time.Second))
			messageType, _, err := conn.ReadMessage()
			Ω(messageType).Should(Equal(-1))
			Ω(websocket.IsUnexpectedCloseError(err)).Should(Equal(false))

			wg.Wait()
		})

		FIt("should provide information about specified statuses", func() {
			var wg sync.WaitGroup

			book := orderbook.NewOrderbook()
			relay := NewRelay(Config{}, dnr.DarknodeRegistry{}, &book, nil, nil, nil)

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

			// Connect to the socket.
			server := httptest.NewServer(RecoveryHandler(relay.GetOrdersHandler()))
			u, _ := url.Parse(server.URL)
			u.Scheme = "ws"
			u.Path = "orders"
			u.RawQuery = "id=vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQ&status=unconfirmed,confirmed"
			conn, _, _ := websocket.DefaultDialer.Dial(u.String(), nil)

			// We should be able to read the initial message.
			err := book.Open(ord)
			Ω(err).ShouldNot(HaveOccurred())

			// Update the status of the order and check if there is another
			// message to be read.
			err = book.Match(ord)
			Ω(err).ShouldNot(HaveOccurred())
			messageType, message, err := conn.ReadMessage()
			var socketMessage orderbook.Entry
			err = json.Unmarshal(message, &socketMessage)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(socketMessage.Status).Should(Equal(order.Unconfirmed))

			err = book.Confirm(ord)
			Ω(err).ShouldNot(HaveOccurred())
			messageType, message, err = conn.ReadMessage()
			err = json.Unmarshal(message, &socketMessage)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(socketMessage.Status).Should(Equal(order.Confirmed))
			Ω(messageType).Should(Equal(websocket.TextMessage))
			Ω(websocket.IsUnexpectedCloseError(err)).Should(Equal(false))

			wg.Wait()
		})

		It("should not provide information about unspecified statuses", func() {
			var wg sync.WaitGroup

			book := orderbook.NewOrderbook()
			relay := NewRelay(Config{}, dnr.DarknodeRegistry{}, &book, nil, nil, nil)

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

			wg.Add(1)
			go func() {
				defer wg.Done()
				// Open an order with the specified ID.
				book.Open(ord)
			}()

			// Connect to the socket.
			server := httptest.NewServer(RecoveryHandler(relay.GetOrdersHandler()))
			u, _ := url.Parse(server.URL)
			u.Scheme = "ws"
			u.Path = "orders"
			u.RawQuery = "id=vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQ&status=unconfirmed,confirmed"
			conn, _, _ := websocket.DefaultDialer.Dial(u.String(), nil)

			// We should be able to read the initial message.
			_, message, err := conn.ReadMessage()
			var socketMessage orderbook.Entry
			err = json.Unmarshal(message, &socketMessage)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(socketMessage.Order.ID).Should(Equal(ord.ID))

			// We should not receive the following message, as we have not
			// included the status as a parameter.
			book.Settle(ord)
			messageType, _, err := conn.ReadMessage()
			Ω(messageType).Should(Equal(-1))
			// TODO: Check for websocket error responses.
			Ω(websocket.IsUnexpectedCloseError(err)).Should(Equal(true))

			wg.Wait()
		})
	})
})
