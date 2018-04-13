package relay

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
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

	messages := make(chan *orderbook.Message, 100)
	queue := NewWriteOnlyChannelQueue(messages, 100)
	go func() {
		if err := orderBook.Subscribe("id", queue); err != nil {
			fmt.Printf("unable to subscribe to order book: %v", err)
		}
	}()

	for {
		fmt.Println("waiting for message")
		select {
		case message, ok := <-messages:
			fmt.Printf("received a message")
			if !ok {
				return
			}
			// TODO: Check status.
			if string(message.Ord.ID) == orderID {
				conn.SetWriteDeadline(time.Now().Add(writeDeadline))
				if err := conn.WriteJSON(message); err != nil {
					fmt.Printf("cannot send json: %v", err)
					return
				}
			}
		case <-ping.C:
			conn.SetWriteDeadline(time.Now().Add(writeDeadline))
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				fmt.Printf("cannot ping websocket: %s", err)
				return
			}
		}
	}
}
