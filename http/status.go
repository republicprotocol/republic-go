package http

import (
	"encoding/json"
	"fmt"
	https "net/http"

	"github.com/gorilla/mux"
	"github.com/republicprotocol/renex-ingress-api-go/http"
	"github.com/republicprotocol/republic-go/http/adapter"
	"github.com/rs/cors"
)

// NewStatusServer returns a new http.Handler for serving darknode status
func NewStatusServer(statusAdapter adapter.StatusAdapter) https.Handler {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/status", statusHandler(statusAdapter)).Methods("GET")
	r.Use(http.RecoveryHandler)

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET"},
	}).Handler(r)

	return handler
}

// statusHandler
func statusHandler(statusAdapter adapter.StatusAdapter) https.HandlerFunc {
	return func(w https.ResponseWriter, r *https.Request) {
		status, err := statusAdapter.Status()
		if err != nil {
			http.WriteError(w, https.StatusBadRequest, fmt.Sprintf("cannot retrieve status object: %v", err))
			return
		}
		str, err := json.Marshal(status)
		if err != nil {
			http.WriteError(w, https.StatusBadRequest, fmt.Sprintf("cannot convert status object into json: %v", err))
			return
		}
		// Set content type to JSON before StatusOK or it will be ignored
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(https.StatusOK)
		w.Write(str)
	}
}
