package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/republicprotocol/republic-go/darknode"
)

func main() {
}

func NewLocalDarkNodes(numberOfDarkNodes, numberOfBootstrapDarkNodes int) (darknode.DarkNodes, []context.Context, []context.CancelFunc) {
	darkNodes := make(darknode.DarkNodes, numberOfDarkNodes)
	ctxs := make([]context.Context, numberOfDarkNodes)
	cancels := make([]context.CancelFunc, numberOfDarkNodes)
	for i := 0; i < numberOfDarkNodes; i++ {
		key := keystore.NewKeyForDirectICAP(rand.Reader)
		darkNodes[i] = NewLocalDarkNode(key, "127.0.0.1", fmt.Sprintf("%d", 3000+i))
		ctxs[i], cancels[i] = context.WithCancel(context.Background())
	}
	return darkNodes, ctxs, cancels
}

func NewLocalDarkNode(key *keystore.Key, host, port string) darknode.DarkNode {
	config := darknode.NewLocalConfig(key, host, port)
	node , err := darknode.NewDarkNode(config)
	if err != nil {
		log.Fatal("fail to create new dark node,", err)
	}
	return node
}
