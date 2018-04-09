package relay

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/republicprotocol/republic-go/order"
)

// handleHTTPRequests will handle POST, DELETE and GET requests
func handleHTTPRequests(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		// TODO: Handle GET requests here
	case "POST":
		ord := order.Order{}
		err := json.NewDecoder(r.Body).Decode(&ord)
		if err != nil {
			fmt.Fprintf("cannot decode json: %v", err)
			return
		}

		SendOrderToDarkOcean(ord)
	case "DELETE":
		// Handle cancel requests here

	default:
		fmt.Fprintf(w, "Only GET, POST and DELETE methods are supported.")
	}
}
