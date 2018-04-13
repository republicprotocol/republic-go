package relay

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/republicprotocol/republic-go/dark"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
)

const reset = "\x1b[0m"

// The HTTPPost object
type HTTPPost struct {
	Order               order.Order                   `json:"order"`
	OrderFragments      Fragments                     `json:"orderFragments"`
}

// The HTTPDelete object
type HTTPDelete struct {
	signature           []byte                        `json:"signature"`
	ID                  order.ID                      `json:"id"`
}

// Fragments will store a list of Fragment Sets with their order details
type Fragments struct {
	Signature           []byte                         `json:"signature"`
	ID                  order.ID                       `json:"id"`

	Type                order.Type                     `json:"type"`
	Parity              order.Parity                   `json:"parity"`
	Expiry              time.Time                      `json:"expiry"`

	DarkPools           map[string][]*order.Fragment   `json:"darkPools"`
}

// RecoveryHandler handles errors while processing the requests and populates the errors in the response
func RecoveryHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				writeError(w, http.StatusInternalServerError, fmt.Sprintf("%v", r))
			}
		}()
		h.ServeHTTP(w, r)
	})
}

// PostOrdersHandler handles all HTTP Post requests
func PostOrdersHandler(multiAddress *identity.MultiAddress, darkPools dark.Pools) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postOrder := HTTPPost{}
		
		if err := json.NewDecoder(r.Body).Decode(&postOrder); err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot decode json into an order or a list of order fragments: %v", err))
			return
		}
		
		if len(postOrder.OrderFragments.DarkPools) > 0 {
			if err := SendOrderFragmentsToDarkOcean(postOrder.OrderFragments, multiAddress, darkPools); err != nil {
				writeError(w, http.StatusInternalServerError, fmt.Sprintf("error sending order fragments : %v", err))
				return
			}
		} else if postOrder.Order.ID.String() != "" {
			if err := SendOrderToDarkOcean(postOrder.Order, multiAddress, darkPools); err != nil {
				writeError(w, http.StatusInternalServerError, fmt.Sprintf("error sending orders : %v", err))
				return
			}
		} else {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot decode json into an order or a list of order fragments: empty object"))
			return
		}
		
		w.WriteHeader(http.StatusCreated)
	})
}

func HandleGetOrder() http.Handler {
	// To-do: Add authentication.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["orderID"]
		// status := getOrderStatus(id)
		status := "confirmed"
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]interface{} {
			"id":     id,
			"status": status,
		}); err != nil {
			fmt.Sprintf("cannot encode json object: %v", err)
		}
	})
}

// DeleteOrderHandler handles HTTP Delete Requests
func DeleteOrderHandler(multiAddress *identity.MultiAddress, darkPools dark.Pools) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cancelOrder := HTTPDelete{}
		if err := json.NewDecoder(r.Body).Decode(&cancelOrder); err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot decode json: %v", err))
			return
		}
		if err := CancelOrder(cancelOrder.ID, multiAddress, darkPools); err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("error canceling orders : %v", err))
			return
		}
		w.WriteHeader(http.StatusGone)
	})
}

func writeError (w http.ResponseWriter, httpCode int, err string) {
	w.WriteHeader(httpCode)
	w.Write([]byte(err))
	return
}
