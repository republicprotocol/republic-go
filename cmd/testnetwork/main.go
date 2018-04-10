package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/republicprotocol/republic-go/contracts/connection"
)

var debug = false
var sleep = 5

const reset = "\x1b[0m"
const green = "\x1b[32;1m"
const yellow = "\x1b[33;1m"

func main() {
	parseCommandLineFlags()

	log := fmt.Sprintf("Started Ganache server on port %s8545%s.", green, reset)
	if !debug {
		log = fmt.Sprintf("%s Run with `-debug` to show output.", log)
	}
	fmt.Printf("%s\n", log)

	var wg sync.WaitGroup
	cmd := connection.StartTestnet(debug, &wg)
	go killAtExit(cmd)

	time.Sleep(time.Duration(sleep) * time.Second)

	err := connection.DeployContractsToGanache("http://localhost:8545")
	if err != nil {
		log.Fatal(err)
	}

	cmd.Wait()
}

func parseCommandLineFlags() error {
	debugPtr := flag.Bool("debug", false, "Print output to stdout")
	sleepPtr := flag.Int("sleep", 5, "Time to wait for ganache to start-up")

	flag.Parse()

	debug = *debugPtr
	sleep = *sleepPtr

	return nil
}

func killAtExit(cmd *exec.Cmd) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Printf("%sShutting down Ganache...%s\n", yellow, reset)
		cmd.Process.Kill()
		os.Exit(0)
	}()
}
