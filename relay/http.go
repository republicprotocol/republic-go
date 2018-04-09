package relay

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/republicprotocol/republic-go/order"
)

func HandleHTTPRequests() {
	http.HandleFunc("/", requestHandler)
}

// Handles POST, DELETE and GET requests.
func requestHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// To-do: Add authentication + get status from ID.
		slices := strings.Split(r.URL.Path, "/")
		id := slices[len(slices) - 1]
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{} {
			"id": id,
			"status": "..",
		})
	case "POST":
		ord := order.Order{}
		err := json.NewDecoder(r.Body).Decode(&ord)
		if err != nil {
			log.Println(err.Error())
			return
		}
		SendOrderToDarkOcean(ord)
	case "DELETE":
		// Handle cancel requests here
	default:
		log.Println(w, "Invalid request")
	}
}
