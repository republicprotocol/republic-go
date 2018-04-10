package relay

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	OrderID string `json:"orderID"`
	Status  string `json:"status"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func GetOrdersHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("cannot open websocket connection: %v", err), http.StatusBadRequest)
		}
		streamOrders(r, conn)
	})
}

// Notifies client if status of specified order has changed.
func streamOrders(r *http.Request, conn *websocket.Conn) {
	traderID := r.FormValue("trader")
	orderID := r.FormValue("order")
	statuses := strings.Split(r.FormValue("status"), ",")
	if (traderID == "" || orderID == "") {
		return
	}
	log.Println("Client subscribed")
	m := Message{}
	orders := make(map[string]string)
	for {
		// status := getOrderStatus(orderID)
		status := "confirmed"

		if (status != orders[orderID]) {
			for _, item := range statuses {
				// If the user specified the current status, notify them.
				if item == status {
					orders[orderID] = status
					m.OrderID = orderID
					m.Status = status
					if err := conn.WriteJSON(m); err != nil {
						fmt.Sprintf("cannot write json: %v", err)
						break
					}
				}
			}
		}
		time.Sleep(time.Second)
	}
	log.Println("Client unsubscribed")
}
