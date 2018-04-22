package relay

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
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
func GetOrdersHandler(book *orderbook.Orderbook) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		defer conn.Close()
		if err != nil {
			fmt.Printf("cannot open websocket connection: %v", err)
			return
		}
		writeUpdatesToWebSocket(w, r, conn, book)
	})
}

// writeUpdatesToWebSocket to notify the client of status changes to specific
// orders.
func writeUpdatesToWebSocket(w http.ResponseWriter, r *http.Request, conn *websocket.Conn, book *orderbook.Orderbook) {
	// Retrieve ID and statuses from URL.
	orderID := r.FormValue("id")
	statuses := strings.Split(r.FormValue("status"), ",")
	orderStatuses := []int{}

	if orderID == "" {
		return
	}

	for _, item := range statuses {
		switch item {
		case "open":
			orderStatuses = append(orderStatuses, 0)
		case "unconfirmed":
			orderStatuses = append(orderStatuses, 1)
		case "canceled":
			orderStatuses = append(orderStatuses, 2)
		case "confirmed":
			orderStatuses = append(orderStatuses, 3)
		case "settled":
			orderStatuses = append(orderStatuses, 4)
		}
	}

	// Handle ping/pong.
	writeDeadline := 10 * time.Second
	pingInterval := 30 * time.Second
	pongInterval := 60 * time.Second
	ping := time.NewTicker(pingInterval)
	defer ping.Stop()
	conn.SetReadDeadline(time.Now().Add(pongInterval))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongInterval))
		return nil
	})

	messages := make(chan orderbook.Entry, 100)
	defer close(messages)

	go func() {
		if err := book.Subscribe(messages); err != nil {
			fmt.Printf("unable to subscribe to order book: %v", err)
		}
	}()
	defer book.Unsubscribe(messages)

	for {
		select {
		case message, ok := <-messages:
			if !ok {
				return
			}
			if !bytes.Equal(message.Order.ID, []byte(orderID)) {
				break
			}
			// Loop through specified statuses.
			for _, status := range orderStatuses {
				if status == int(message.Status) {
					conn.SetWriteDeadline(time.Now().Add(writeDeadline))
					if err := conn.WriteJSON(message); err != nil {
						fmt.Printf("cannot send json: %v", err) // FIXME: Use a logger
						return
					}
				}
			}
			// If the user hasn't specified a status, send them all messages.
			if len(orderStatuses) == 0 {
				conn.SetWriteDeadline(time.Now().Add(writeDeadline))
				if err := conn.WriteJSON(message); err != nil {
					fmt.Printf("cannot send json: %v", err) // FIXME: Use a logger
					return
				}
			}
		case <-ping.C:
			conn.SetWriteDeadline(time.Now().Add(writeDeadline))
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				fmt.Printf("cannot ping websocket: %s", err) // FIXME: Use a logger
				return
			}
		}
	}
}
