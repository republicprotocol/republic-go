package test

import (
	"flag"

	"github.com/republicprotocol/go-dark-node"
)

func main() {

	numberOfNodes := flag.Int("nodes", 20, "number of nodes")

	nodes := generateNodes(numberOfNodes)

	err := deployNodes(nodes)

	sendingOrders()

	collectLogs()

}

func generateNodes(numberOfNodes int) []*node.DarkNode {
	return nil
}

func deployNodes(nodes []*node.DarkNode) error {
	return nil
}

func sendingOrders() {

}

func collectLogs() {

}
