package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/republicprotocol/republic-go/http/adapter"
	"github.com/rs/cors"
)

// OpenOrderRequest is an JSON object sent to the HTTP handlers to request the
// opening of an order.
type OpenOrderRequest struct {
	Signature             string                        `json:"signature"`
	OrderFragmentMappings adapter.OrderFragmentMappings `json:"orderFragmentMappings"`
}

func NewServer(openOrderAdapter adapter.OpenOrderAdapter, cancelOrderAdapter adapter.CancelOrderAdapter) http.Handler {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/orders", OpenOrderHandler(openOrderAdapter)).Methods("POST")
	r.HandleFunc("/orders", CancelOrderHandler(cancelOrderAdapter)).Methods("DELETE")
	r.Use(RecoveryHandler)

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	}).Handler(r)

	return handler
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

// OpenOrderHandler handles all HTTP open order requests
func OpenOrderHandler(adapter adapter.OpenOrderAdapter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		openOrderRequest := OpenOrderRequest{}
		if err := json.NewDecoder(r.Body).Decode(&openOrderRequest); err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot decode json into an order or a list of order fragments: %v", err))
			return
		}
		if err := adapter.OpenOrder(openOrderRequest.Signature, openOrderRequest.OrderFragmentMappings); err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot open order: %v", err))
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

// CancelOrderHandler handles HTTP Delete Requests
func CancelOrderHandler(adapter adapter.CancelOrderAdapter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderID := r.URL.Query().Get("id")
		if orderID == "" {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot cancel order: nil id"))
			return
		}
		signature := r.URL.Query().Get("signature")
		if signature == "" {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot cancel order: nil signature"))
			return
		}
		if err := adapter.CancelOrder(signature, orderID); err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("cannot cancel order: %v", err))
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func writeError(w http.ResponseWriter, statusCode int, err string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(err))
	return
}
