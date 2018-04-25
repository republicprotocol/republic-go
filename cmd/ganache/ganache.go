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

	"github.com/republicprotocol/republic-go/blockchain/test/ganache"
)

const reset = "\x1b[0m"
const green = "\x1b[32;1m"
const yellow = "\x1b[33;1m"

func main() {
	argSleep := flag.Int("sleep", 10, "Seconds to sleep after starting ganache")
	flag.Parse()

	fmt.Printf("Ganache is listening on %shttp://localhost:8545%s...\n", green, reset)

	cmd := ganache.Start()
	go killAtExit(cmd)

	time.Sleep(time.Duration(*argSleep) * time.Second)

	conn, err := ganache.Connect("http://localhost:8545")
	if err != nil {
		log.Fatalf("cannot connect to ganache: %v", err)
		return
	}

	if err := ganache.DeployContracts(conn); err != nil {
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
		cmd.Process.Kill()
		os.Exit(0)
	}()
}
