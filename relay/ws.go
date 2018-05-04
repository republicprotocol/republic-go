package relay

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/republicprotocol/republic-go/order"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// GetOrdersHandler handles WebSocket requests.
func (relay *Relay) GetOrdersHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		defer conn.Close()
		if err != nil {
			fmt.Printf("cannot open websocket connection: %v", err)
			return
		}
		relay.writeUpdatesToWebSocket(w, r, conn)
	})
}

// writeUpdatesToWebSocket to notify the client of status changes to specific
// orders.
func (relay *Relay) writeUpdatesToWebSocket(w http.ResponseWriter, r *http.Request, conn *websocket.Conn) {
	// Retrieve ID and statuses from URL.
	orderID := r.FormValue("id")
	statuses := strings.Split(r.FormValue("status"), ",")
	orderStatuses := []order.Status{}

	if orderID == "" {
		return
	}

	for _, item := range statuses {
		switch item {
		case "open":
			orderStatuses = append(orderStatuses, order.Open)
		case "unconfirmed":
			orderStatuses = append(orderStatuses, order.Unconfirmed)
		case "canceled":
			orderStatuses = append(orderStatuses, order.Canceled)
		case "confirmed":
			orderStatuses = append(orderStatuses, order.Confirmed)
		case "settled":
			orderStatuses = append(orderStatuses, order.Settled)
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

	done := make(chan struct{})
	messages := relay.orderbook.Listen(done)
	defer close(done)

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
				if status == message.Status {
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
