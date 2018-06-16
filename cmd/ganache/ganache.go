// package main

// import (
// 	"flag"
// 	"fmt"
// 	"log"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"time"

// 	"github.com/republicprotocol/republic-go/blockchain/test/ganache"
// )

// const reset = "\x1b[0m"
// const green = "\x1b[32;1m"
// const yellow = "\x1b[33;1m"

// func main() {
// 	argSleep := flag.Int("sleep", 10, "Seconds to sleep after starting ganache")
// 	flag.Parse()

// 	fmt.Printf("Ganache is listening on %shttp://localhost:8545%s...\n", green, reset)

// 	ganache.Start()
// 	time.Sleep(time.Duration(*argSleep) * time.Second)

// 	conn, err := ganache.Connect("http://localhost:8545")
// 	if err != nil {
// 		log.Fatalf("cannot connect to ganache: %v", err)
// 	}

// 	if err := ganache.DeployContracts(conn); err != nil {
// 		log.Fatalf("cannot deploy contracts to ganache: %v", err)
// 	}

// 	signals := make(chan os.Signal, 1)
// 	defer close(signals)

// 	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
// 	<-signals
// }
