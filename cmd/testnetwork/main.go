package main

import (
	"flag"
	"fmt"

	"github.com/republicprotocol/republic-go/contracts/connection"
)

var debug = false

const reset = "\x1b[0m"
const green = "\x1b[32;1m"

func main() {
	parseCommandLineFlags()

	log := fmt.Sprintf("Started Ganache server on port %s8545%s.", green, reset)
	if !debug {
		log = fmt.Sprintf("%s Run with `-debug` to show output.", log)
	}
	fmt.Printf("%s\n", log)

	connection.StartTestnet(debug)
}

func parseCommandLineFlags() error {
	debugPtr := flag.Bool("debug", false, "Print output to stdout")

	flag.Parse()

	debug = *debugPtr

	return nil
}
