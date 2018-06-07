package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/republicprotocol/republic-go/http/adapter"
	"github.com/rs/cors"
)

func NewStatusServer(statusAdapter adapter.StatusAdapter) http.Handler {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/status", statusHandler(statusAdapter)).Methods("GET")
	r.Use(RecoveryHandler)

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET"},
	}).Handler(r)

	return handler
}

// statusHandler
func statusHandler(statusAdapter adapter.StatusAdapter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		str, err := json.Marshal(statusAdapter.Status())
		if err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("cannot convert status object into json: %v", err))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(str)
	}
}
