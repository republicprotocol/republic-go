package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/republicprotocol/republic-go/http/adapter"
	"github.com/republicprotocol/republic-go/order"
)

const reset = "\x1b[0m"

type OpenOrderRequest struct {
	Signature            string                       `json:"signature"`
	OrderFragmentMapping adapter.OrderFragmentMapping `json:"orderFragmentMapping"`
}

type CancelOrderRequest struct {
	Signature string   `json:"signature"`
	ID        order.ID `json:"id"`
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
func AuthorizationHandler(authProvider adapter.AuthProvider, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if authProvider.RequireAuth() {
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
			if coms[0] != "Bearer" {
				if err := authProvider.Verify(coms[1]); err != nil {
					writeError(w, http.StatusUnauthorized, err.Error())
					return
				}
			}
		}
		h.ServeHTTP(w, r)
	})
}

// OpenOrderHandler handles all HTTP open order requests
func OpenOrderHandler(adapter adapter.OpenOrderAdapter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		openOrderRequest := OpenOrderRequest{}
		if err := json.NewDecoder(r.Body).Decode(&openOrderRequest); err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot decode json into an order or a list of order fragments: %v", err))
			return
		}
		if err := adapter.OpenOrder(openOrderRequest.Signature, openOrderRequest.OrderFragmentMapping); err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot open order: %v", err))
			return
		}
		w.WriteHeader(http.StatusCreated)
	})
}

// CancelOrderHandler handles HTTP Delete Requests
func CancelOrderHandler(adapter adapter.CancelOrderAdapter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// FIXME: Check that cancelOrder.ID matches mux.Vars(r)["orderID"]

		cancelOrderRequest := CancelOrderRequest{}
		if err := json.NewDecoder(r.Body).Decode(&cancelOrderRequest); err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot decode json: %v", err))
			return
		}
		if err := adapter.CancelOrder(cancelOrderRequest.Signature, cancelOrderRequest.ID); err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot cancel order: %v", err))
			return
		}
		w.WriteHeader(http.StatusGone)
	})
}

// GetOrderHandler handles all HTTP GET requests.
/* func (relay *Relay) GetOrderHandler() http.Handler {
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
} */

func writeError(w http.ResponseWriter, statusCode int, err string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(err))
	return
}
