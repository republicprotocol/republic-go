package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strings"
	"syscall"

	"github.com/republicprotocol/republic-go/identity"

	"github.com/republicprotocol/republic-go/darknode"
)

var configFile *string

func main() {

	// Parse command-line arguments
	configFile = flag.String("config", path.Join(os.Getenv("HOME"), ".darknode/config.json"), "JSON configuration file")
	flag.Parse()

	// Load configuration file
	config, err := darknode.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	// Get IP-address
	ipAddr, err := getIPAddress()
	if err != nil {
		log.Fatalf("cannot get ip-address: %v", err)
	}

	// Get multi-address
	multiAddr, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/%s/tcp/18514/republic/%s", ipAddr, config.Address))
	if err != nil {
		log.Fatalf("cannot get multiaddress: %v", err)
	}

	// Create the Darknode
	node, err := darknode.NewDarknode(multiAddr, config)
	if err != nil {
		log.Fatalf("cannot create darknode: %v", err)
	}

	// Run the Darknode until the process is terminated
	done := make(chan struct{})
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	go func() {
		defer close(done)
		<-sig
	}()
	for err := range node.Run(done) {
		log.Println(err)
	}
}

func getIPAddress() (string, error) {

	out, err := exec.Command("curl", "https://ipinfo.io/ip").Output()
	out = []byte(strings.Trim(string(out), "\n "))
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}

	return string(out), nil
}
