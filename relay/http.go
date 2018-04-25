package relay

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/republicprotocol/republic-go/order"
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

// AuthorizationHandler handles errors while processing the requests and populates the errors in the response
func (relay *Relay) AuthorizationHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if relay.Config.Token != "" {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeError(w, http.StatusUnauthorized, "")
				return
			}
			coms := strings.Split(authHeader, " ")
			if len(coms) != 2 {
				writeError(w, http.StatusUnauthorized, "")
				return
			}
			if coms[0] != "Bearer" || coms[1] != relay.Config.Token {
				writeError(w, http.StatusUnauthorized, "")
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

// OpenOrdersHandler handles all HTTP open order requests
func (relay *Relay) OpenOrdersHandler() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		openOrder := OpenOrderRequest{}

		if err := json.NewDecoder(r.Body).Decode(&openOrder); err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot decode json into an order or a list of order fragments: %v", err))
			return
		}

		if len(openOrder.OrderFragments.DarkPools) > 0 {
			if err := relay.SendOrderFragmentsToDarkOcean(openOrder.OrderFragments); err != nil {
				writeError(w, http.StatusInternalServerError, fmt.Sprintf("error sending order fragments: %v", err))
				return
			}
		} else if openOrder.Order.ID.String() != "" {
			if err := relay.SendOrderToDarkOcean(openOrder.Order); err != nil {
				writeError(w, http.StatusInternalServerError, fmt.Sprintf("error sending orders: %v", err))
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
func (relay *Relay) GetOrderHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orderID := vars["orderID"]
		if orderID == "" {
			writeError(w, http.StatusBadRequest, "order id is invalid")
			return
		}

		// Check if there exists an item in the order book with the given ID.
		message := relay.orderbook.Order([]byte(orderID))
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
func (relay *Relay) CancelOrderHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// FIXME: Check that cancelOrder.ID matches mux.Vars(r)["orderID"]

		cancelOrder := CancelOrderRequest{}
		if err := json.NewDecoder(r.Body).Decode(&cancelOrder); err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot decode json: %v", err))
			return
		}
		if err := relay.CancelOrder(cancelOrder.ID); err != nil {
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
