package main

import (
	"fmt"
	"net/http"

	"github.com/republicprotocol/republic-go/relay"
)

func main() {
	r := relay.NewRouter()
	if err := http.ListenAndServe("0.0.0.0:8000", r); err != nil {
		fmt.Sprintf("could not start router: %v", err)
	}
}
