package relay

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
)

const reset = "\x1b[0m"

// OpenOrderRequest is a JSON request to open an order in the Darkpool that is
// optionally split into order fragments, and optionally signed.
type OpenOrderRequest struct {
	Order          order.Order    `json:"order"`
	OrderFragments OrderFragments `json:"orderFragments"`
}

// CancelOrderRequest is a JSON request to cancel an order that is optionally
// signed.
type CancelOrderRequest struct {
	Signature []byte   `json:"signature"`
	ID        order.ID `json:"id"`
}

// OrderFragments is a JSON representation of order fragments that have been
// split for the different pools.
type OrderFragments struct {
	Signature []byte   `json:"signature"`
	ID        order.ID `json:"id"`

	Type   order.Type   `json:"type"`
	Parity order.Parity `json:"parity"`
	Expiry time.Time    `json:"expiry"`

	DarkPools map[string][]*order.Fragment `json:"darkPools"`
}

// AuthorizationHandler handles errors while processing the requests and populates the errors in the response
func AuthorizationHandler(h http.Handler, token string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token != "" {
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				if strings.Split(authHeader, " ")[1] != token {
					writeError(w, http.StatusUnauthorized, fmt.Sprintln("Unauthorized token"))
					return
				}
			}
		}
		h.ServeHTTP(w, r)
	})
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

// OpenOrdersHandler handles all HTTP open order requests
func OpenOrdersHandler(relayConfig Relay) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		openOrder := OpenOrderRequest{}

		if err := json.NewDecoder(r.Body).Decode(&openOrder); err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot decode json into an order or a list of order fragments: %v", err))
			return
		}

		if len(openOrder.OrderFragments.DarkPools) > 0 {
			if err := SendOrderFragmentsToDarkOcean(openOrder.OrderFragments, relayConfig); err != nil {
				writeError(w, http.StatusInternalServerError, fmt.Sprintf("error sending order fragments : %v", err))
				return
			}
		} else if openOrder.Order.ID.String() != "" {
			if err := SendOrderToDarkOcean(openOrder.Order, relayConfig); err != nil {
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

// GetOrderHandler handles all HTTP GET requests.
func GetOrderHandler(book *orderbook.Orderbook, id string) http.Handler {
	// TODO: Add authentication.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orderID := vars["orderID"]
		if orderID == "" {
			if id == "" {
				writeError(w, http.StatusBadRequest, "order id is invalid")
				return
			}
			orderID = id
		}

		// Check if there exists an item in the order book with the given ID.
		message := book.Order([]byte(orderID))
		if message.Order.ID == nil {
			writeError(w, http.StatusBadRequest, "order id is invalid")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(message.Order); err != nil {
			fmt.Printf("cannot encode object as json: %v", err)
		}
	})
}

// CancelOrderHandler handles HTTP Delete Requests
func CancelOrderHandler(relayConfig Relay) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cancelOrder := CancelOrderRequest{}
		if err := json.NewDecoder(r.Body).Decode(&cancelOrder); err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot decode json: %v", err))
			return
		}
		if err := CancelOrder(cancelOrder.ID, relayConfig); err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("error canceling orders : %v", err))
			return
		}
		w.WriteHeader(http.StatusGone)
	})
}

func writeError(w http.ResponseWriter, statusCode int, err string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(err))
	return
}
