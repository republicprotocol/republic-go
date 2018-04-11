package main

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/republicprotocol/republic-go/darknode"
)

func main() {
}

func NewLocalDarknodes(numberOfDarknodes, numberOfBootstrapDarknodes int) (darknode.Darknodes, []context.Context, []context.CancelFunc) {
	darknodes := make(Darknodes, numberOfDarknodes)
	ctxs := make([]context.Context, numberOfDarknodes)
	cancels := make([]context.CancelFunc, numberOfDarknodes)
	for i := 0; i < numberOfDarknodes; i++ {
		key := keystore.NewKeyForDirectICAP(rand.Reader)
		darknodes[i] = NewLocalDarknode(key, "127.0.0.1", fmt.Sprintf("%d", 3000+i))
		ctx[i], cancels[i] = context.WithCancel(context.Background())
	}
	return darknodes, ctxs, cancels
}

func NewLocalDarknode(key *keystore.Key, host, port string) Darknode {
	config := darknode.NewLocalConfig(key, host, port)
	return darknode.NewDarknode(config)
}
