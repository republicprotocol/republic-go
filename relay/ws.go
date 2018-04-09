package relay

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func HandleSocketRequests() {
	http.HandleFunc("/orders", socketHandler)
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}
	go streamOrders(r, conn)
}

func streamOrders(r *http.Request, conn *websocket.Conn) {
	traderID := strings.Replace(r.FormValue("trader"), "\"", "", -1)
	orderID := strings.Replace(r.FormValue("order"), "\"", "", -1)
	status := strings.Replace(r.FormValue("status"), "\"", "", -1)
	if (traderID == "" || orderID == "" || status == "") {
		return
	}
	log.Println("Client subscribed")
	for {
		if err := conn.WriteMessage(websocket.TextMessage, []byte("Here is a status update")); err != nil {
			log.Println(err)
			break
		}
		time.Sleep(time.Second)
	}
	log.Println("Client unsubscribed")
}
