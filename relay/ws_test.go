package relay_test

import (
	"log"
	"net/http/httptest"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/orderbook"
	. "github.com/republicprotocol/republic-go/relay"
)

var _ = Describe("WebSocket streaming", func() {
	Context("when connecting to the socket", func() {
		It("should error for invalid parameters", func() {
			orderBook := orderbook.NewOrderBook(100)
			server := httptest.NewServer(RecoveryHandler(GetOrdersHandler(orderBook)))
			u, _ := url.Parse(server.URL)
			u.Scheme = "ws"
			u.Path = "orders"
			// Note that we don't specify any query parameters.

			conn, _, _ := websocket.DefaultDialer.Dial(u.String(), nil)
			_, _, err := conn.ReadMessage()

			// In this case when we attempt to read, the server is already closed, so
			// we should have an unexpected close error.
			Ω(websocket.IsUnexpectedCloseError(err)).Should(Equal(true))
		})

		It("should be able to successfully read from socket with valid parameters", func() {
			orderBook := orderbook.NewOrderBook(100)
			server := httptest.NewServer(RecoveryHandler(GetOrdersHandler(orderBook)))
			u, _ := url.Parse(server.URL)
			u.Scheme = "ws"
			u.Path = "orders"
			u.RawQuery = "id=test"

			conn, _, _ := websocket.DefaultDialer.Dial(u.String(), nil)
			conn.SetReadDeadline(time.Now().Add(time.Second))
			_, _, err := conn.ReadMessage()

			// In this case the server is still open when we read, but the deadline
			// times out due to not receiving a message, so we should not not have an
			// unexpected close error.
			Ω(websocket.IsUnexpectedCloseError(err)).Should(Equal(false))
		})

		// It("should routinely send a ping message", func() {
		// 	orderBook := orderbook.NewOrderBook(100)
		// 	server := httptest.NewServer(RecoveryHandler(GetOrdersHandler(orderBook)))
		// 	u, _ := url.Parse(server.URL)
		// 	u.Scheme = "ws"
		// 	u.Path = "orders"
		// 	u.RawQuery = "id=test"
		//
		// 	conn, _, _ := websocket.DefaultDialer.Dial(u.String(), nil)
		// 	// conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		// 	var err error
		// 	var messageType int
		// 	log.Println("waiting for message")
		// 	messageType, _, err = conn.ReadMessage()
		// 	log.Println(messageType)
		//
		// 	// In this case we have increased the read deadline, so we are able to
		// 	// receive a ping message.
		// 	Ω(websocket.IsUnexpectedCloseError(err)).Should(Equal(false))
		// })
	})
})
