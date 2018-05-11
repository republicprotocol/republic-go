package http

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/republicprotocol/republic-go/http/adapter"
	"github.com/republicprotocol/republic-go/order"
)

const reset = "\x1b[0m"

// OpenOrderRequest is an JSON object sent to the HTTP handlers to request the
// opening of an order.
type OpenOrderRequest struct {
	Signature            string                       `json:"signature"`
	OrderFragmentMapping adapter.OrderFragmentMapping `json:"orderFragmentMapping"`
}

// CancelOrderRequest is an JSON object sent to the HTTP handlers to request
// the cancelation of an order.
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

func ListenAndServe(bind, port string, openOrderAdapter adapter.OpenOrderAdapter, cancelOrderAdapter adapter.CancelOrderAdapter) error {
	r := mux.NewRouter().StrictSlash(true)
	r.Methods("POST").Path("/orders").Handler(RecoveryHandler(OpenOrderHandler(openOrderAdapter)))
	r.Methods("DELETE").Path("/orders/{orderID}").Handler(RecoveryHandler(CancelOrderHandler(cancelOrderAdapter)))
	return http.ListenAndServe(fmt.Sprintf("%v:%v", bind, port), r)
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
		orderIDString, ok := mux.Vars(r)["orderID"]
		if !ok {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot cancel order: nil id"))
			return
		}
		orderID, err := base64.StdEncoding.DecodeString(orderIDString)
		if err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot decode order id: %v", err))
			return
		}

		cancelOrderRequest := CancelOrderRequest{}
		if err := json.NewDecoder(r.Body).Decode(&cancelOrderRequest); err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot decode json: %v", err))
			return
		}
		if !bytes.Equal(cancelOrderRequest.ID, orderID) {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot cancel order: invalid id"))
			return
		}
		if err := adapter.CancelOrder(cancelOrderRequest.Signature, cancelOrderRequest.ID); err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot cancel order: %v", err))
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
