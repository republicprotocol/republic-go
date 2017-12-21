package main

import (
	"flag"
	"google.golang.org/grpc/reflection"

	"bufio"
	"fmt"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-swarm"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"strings"
)

// Declare command line arguments.
var port = flag.Int("port", 8080, "RPC listening port")

func main() {

	// Parse command line arguments.
	flag.Parse()

	// Generate identity and address for this node.
	keyPair, err := identity.NewKeyPair()
	if err != nil {
		log.Fatal(err)
	}
	multiAddress, err := identity.NewMultiAddress(fmt.Sprintf("/ip4/0.0.0.0/tcp/%v/republic/%s", *port, keyPair.PublicAddress()))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Republic address:", keyPair.PublicAddress())

	// listen to the tcp port
	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	// Create gRPC services.

	peer := swarm.NewPeer(&swarm.Config{
		KeyPair:      keyPair,
		MultiAddress: multiAddress,
		Peers:        make([]identity.MultiAddress, 0, 100),
	})

	// Register reflection service on gRPC server.
	reflection.Register(s)
	log.Printf("Listening for connections on %d...\n", *port)
	// Start the server and listen to request in background
	go s.Serve(lis)

	// handle request
	scanner := bufio.NewScanner(os.Stdin)
	var text []string
	for {
		fmt.Print("Enter command: ")
		scanner.Scan()
		text = strings.Split(scanner.Text(), " ")
		if len(text) == 0 {
			continue
		} else if text[0] == "quit" {
			break
		} else if len(text) != 2 {
			log.Println("Please enter a valid command")
			continue
		}

		// command handler
		switch text[0] {
		case "ping":
			log.Println("Ping: " + text[1])
			pong, err := peer.PingPeer(text[1])
			if err != nil {
				log.Fatal(err)
			}
			log.Println("Pong: " + pong.Multi)
		case "peers":
			if len(text) == 1 {
				log.Println("Peers from the routing table:")
				for _, multi := range peer.DHT.MultiAddresses() {
					log.Println(multi)
				}
			} else {
				log.Printf("Ask all peers from : %s \n", text[1])
				peers, _ := peer.AskPeers(text[1])
				for _, j := range peers.Multis {
					log.Printf("Get peer from server : %s \n", j)
				}
			}
		case "find":
			log.Printf("Try to find node : %s \n", text[1])
			peers, err := peer.FindPeer(identity.Address(text[1]))
			if err != nil {
				log.Fatal(err)
			}
			for _, j := range peers.Multis {
				log.Printf("Peer found : %s \n", j)
			}
		}

	}

}
