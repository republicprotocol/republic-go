package main

import (
	"github.com/republicprotocol/republic/crypto"
	"log"
	"github.com/republicprotocol/republic"
	"github.com/republicprotocol/republic/dht"
)

const (
	address = "localhost:8080"
)


func main() {
	//// Set up a connection to the server.
	//conn, err := grpc.Dial(address, grpc.WithInsecure())
	//if err != nil {
	//	log.Fatalf("did not connect: %v", err)
	//}
	//defer conn.Close()
	//c := rpc.NewNodeClient(conn)
	//
	//// Generating identity for the client node
	//secp, err := crypto.NewSECP256K1()
	//if err != nil {
	//	log.Fatalf("failed to identify self: %v", err)
	//}
	//id := secp.PublicAddress()
	//grpcID := &rpc.ID{Address: id}
	//log.Println("Republic address:", id)
	//
	//// Ping the server
	//log.Println("Ping: " + address)
	//rID, err := c.Ping(context.Background(), grpcID)
	//if err != nil {
	//	log.Fatalf("could not ping: %v", err)
	//}
	//log.Printf("Pong: %s \n", rID.Address)
	//
	//// Get all peers from the server
	//log.Printf("Ask all peers from : %s \n", rID.Address)
	//rMultiAddresses, err := c.Peers(context.Background(), grpcID)
	//if err != nil {
	//	log.Fatalf("could not get peers: %v", err)
	//}
	//
	//for _, j := range rMultiAddresses.Multis {
	//	log.Printf("Get peer from server : %s \n", j)
	//}
	//
	//// Get peers from a node that are closer to a target than the node itself
	//log.Printf("Ask peers close to: %s \n", id)
	//rMultiAddresses, err = c.CloserPeers(context.Background(), &rpc.Path{To: grpcID, From: grpcID})
	//if err != nil {
	//	log.Fatalf("could not get peers: %v", err)
	//}
	//
	//for _, j := range rMultiAddresses.Multis {
	//	log.Printf("Closer peer : %s \n", j)
	//}

	// Generating identity for the client node
	secp, err := crypto.NewSECP256K1()
	if err != nil {
		log.Fatalf("failed to identify self: %v", err)
	}
	id := secp.PublicAddress()
	log.Println("Republic address:", id)

	node := swarm.NewNode(dht.ID(id))

	// ping the server
	log.Println("Ping: " + address)
	pong, err := node.PingNode(address)
	log.Println("Pong: " + pong.Address)

	log.Printf("Ask all peers from : %s \n", pong.Address)
	peers,_ := node.PeersNode(address)
	for _, j := range peers.Multis {
		log.Printf("Get peer from server : %s \n", j)
	}

	// todo : has bug, need to fix it
	//log.Println("Find node Ugz7rk5EjOSyhIenhKWH6RRW07k= : ")
	//target, err := node.FindNode("Ugz7rk5EjOSyhIenhKWH6RRW07k=")
	//if err != nil {
	//	log.Fatal("err:", err)
	//}
	//log.Println("Find target: "+target)
}
