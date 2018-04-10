package relay

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

// Checks to see if status of specified order has changed.
// To-do: Check using real data. Add support for multiple statuses/orders.
func streamOrders(r *http.Request, conn *websocket.Conn) {
	traderID := strings.Replace(r.FormValue("trader"), "\"", "", -1)
	orderID := strings.Replace(r.FormValue("order"), "\"", "", -1)
	status := strings.Replace(r.FormValue("status"), "\"", "", -1)
	if traderID == "" || orderID == "" || status == "" {
		return
	}
	log.Println("Client subscribed")
	for {
		if err := conn.WriteMessage(websocket.TextMessage, []byte("Here is a status update")); err != nil {
			log.Println(err.Error())
			break
		}
		time.Sleep(time.Second)
	}
	log.Println("Client unsubscribed")
}
