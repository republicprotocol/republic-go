package main

import (
	"flag"
	"google.golang.org/grpc/reflection"

	"github.com/republicprotocol/go-swarm/rpc"
	"google.golang.org/grpc"
	"log"
	"net"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-swarm"
)

// Declare command line arguments.
var port = flag.String("port", "8080", "RPC listening port")

func main() {

	// Parse command line arguments.
	flag.Parse()

	// Generate identity for this node.
	secp, err := identity.NewKeyPair()
	if err != nil {
		log.Fatalf("failed to identify self: %v", err)
	}
	address := secp.PublicAddress()
	log.Println("Republic address:", address)

	// listen to the tcp port
	lis, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	// Create gRPC services.
	node := swarm.NewNode(address)
	rpc.RegisterDHTServer(s, node)

	// Register reflection service on gRPC server.
	reflection.Register(s)
	log.Printf("Listening for connections on %s...\n", *port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
