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

func GetOrdersHandler(orderBook *orderbook.OrderBook) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("cannot open websocket connection: %v", err), http.StatusBadRequest)
		}
		streamOrders(r, conn, orderBook)
	})
}

// Notifies client if status of specified order has changed.
func streamOrders(r *http.Request, conn *websocket.Conn, orderBook *orderbook.OrderBook) {
	// Retrieve parameters from URL.
	traderID := r.FormValue("trader")
	orderID := r.FormValue("order")
	// statuses := strings.Split(r.FormValue("status"), ",")
	if traderID == "" || orderID == "" {
		return
	}

	// Handle ping/pong.
	writeDeadline := 10 * time.Second
	pingInterval := 30 * time.Second
	pongInterval := 60 * time.Second
	ping := time.NewTicker(pingInterval)
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongInterval))
		return nil
	})

	defer func() {
		ping.Stop()
		conn.Close()
	}()

	orders := make(map[string]order.Status)
	messages := make(chan *orderbook.Message, 100)
	queue := NewWriteOnlyChannelQueue(messages, 100)
	orderBook.Subscribe("id", queue)

	for {
		select {
		case message, ok := <-messages:
			if !ok {
				return
			}
			if message.Status > orders[string(message.Ord.ID)] {
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
