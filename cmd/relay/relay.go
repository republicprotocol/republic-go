package main

import (
	"flag"
	"fmt"
	// "net/http"

	"github.com/republicprotocol/republic-go/relay"
)

func main() {
	// Parse the command-line arguments
	keystore := flag.String("keystore", "", "path of keystore file")
	passphrase := flag.String("passphrase", "", "passphrase to decrypt keystore")
	bindAddress := flag.String("bind", "", "bind address")
	port := flag.Int("port", 80, "port to bind API")
	token := flag.String("token", "", "optional token")
	flag.Parse()

	if flag.Parsed() {
		if *keystore == "" || *passphrase == "" || *bindAddress == "" {
			// fmt.Println("usage: relay <--keystore> <--passphrase> <--bind> <--port> [--token]")
			flag.PrintDefaults()
			fmt.Sprintf("%v %v", *port, *token)
			return
		}
	}

	r := relay.NewRouter()
	if err := http.ListenAndServe("0.0.0.0:8000", r); err != nil {
		fmt.Sprintf("could not start router: %v", err)
	}
}
