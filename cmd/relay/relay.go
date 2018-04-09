package main

import (
	"log"
	"net/http"

	"github.com/republicprotocol/republic-go/relay"
)

func main() {
	relay.HandleSocketRequests()
	relay.HandleHTTPRequests()
	if err := http.ListenAndServe("0.0.0.0:8000", nil); err != nil {
		log.Fatal(err.Error())
	}
}
