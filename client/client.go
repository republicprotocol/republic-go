package main

import (
	"github.com/republicprotocol/republic/crypto"
	"github.com/republicprotocol/republic/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

const (
	address = "localhost:8080"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := rpc.NewNodeClient(conn)

	// Generating idendity for the client node
	secp, err := crypto.NewSECP256K1()
	if err != nil {
		log.Fatalf("failed to identify self: %v", err)
	}
	id := secp.PublicAddress()
	log.Println("Republic address:", id)

	// Ping the server
	rID, err := c.Ping(context.Background(), &rpc.ID{Address: id})
	log.Println("Ping: " + address)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Pong: %s \n", rID.Address)
}
