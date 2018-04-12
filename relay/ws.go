package relay

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// GetOrdersHandler handles WebSocket requests.
func GetOrdersHandler(orderBook *orderbook.OrderBook) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		defer conn.Close()
		if err != nil {
			fmt.Sprintf("cannot open websocket connection: %v", err)
			return
		}
		streamOrders(w, r, conn, orderBook)
	})
}

// streamOrders notifies client if status of specified order has changed.
func streamOrders(w http.ResponseWriter, r *http.Request, conn *websocket.Conn, orderBook *orderbook.OrderBook) {
	// Retrieve ID from URL.
	orderID := r.FormValue("id")
	if orderID == "" {
		return
	}

	// Handle ping/pong.
	writeDeadline := 10 * time.Second
	pingInterval := 5 * time.Second
	pongInterval := 60 * time.Second
	ping := time.NewTicker(pingInterval)
	defer ping.Stop()
	conn.SetReadDeadline(time.Now().Add(pongInterval))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongInterval))
		return nil
	})

	orders := make(map[string]order.Status)
	messages := make(chan *orderbook.Message, 100)
	queue := NewWriteOnlyChannelQueue(messages, 100)
	go func() {
		if err := orderBook.Subscribe("id", queue); err != nil {
			fmt.Sprintf("unable to subscribe to order book: %v", err)
		}
	}()
	orders[orderID] = order.Open

	for {
		select {
		case message, ok := <-messages:
			if !ok {
				return
			}
			if message.Status >= orders[string(message.Ord.ID)] {
				conn.SetWriteDeadline(time.Now().Add(writeDeadline))
				if err := conn.WriteJSON(message); err != nil {
					fmt.Sprintf("cannot send json: %v", err)
					return
				}
			}
		case <-ping.C:
			conn.SetWriteDeadline(time.Now().Add(writeDeadline))
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				fmt.Sprintf("cannot ping websocket: %s", err)
				return
			}
		}
	}
}
