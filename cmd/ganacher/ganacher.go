package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/republicprotocol/republic-go/ethereum/ganache"
)

const reset = "\x1b[0m"
const green = "\x1b[32;1m"
const yellow = "\x1b[33;1m"

func main() {
	argSleep := flag.Int("sleep", 10, "Time to wait for ganache to start-up")
	flag.Parse()

	fmt.Printf("Started Ganache server on port %s8545%s...\n", green, reset)

	cmd := ganache.Start()
	go killAtExit(cmd)

	time.Sleep(time.Duration(*argSleep) * time.Second)

	conn, err := ganache.Connect("http://localhost:8545")
	if err != nil {
		log.Fatalf("cannot connect to ganache: %v", err)
		return
	}

	if err := ganache.DeployContracts(*conn); err != nil {
		log.Fatalf("cannot deploy contracts to ganache: %v", err)
		return
	}

	cmd.Wait()
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
