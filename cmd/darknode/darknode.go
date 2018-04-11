package main

import (
	"context"

	"github.com/republicprotocol/republic-go/darknode"
)

func main() {
}

func NewLocalDarknodes(numberOfDarknodes, numberOfBootstrapDarknodes int) (darknode.Darknodes, []context.Context, []context.CancelFunc) {
	darknodes := make(Darknodes, numberOfDarknodes)
	ctxs := make([]context.Context, numberOfDarknodes)
	cancels := make([]context.CancelFunc, numberOfDarknodes)
	for i := 0; i < numberOfDarknodes; i++ {
		var config Config
		if i < numberOfBootstrapDarknodes {
			config = darknode.NewLocalConfig(3000 + i)
		} else {
			config = darknode.NewLocalConfig(4000 + i)
		}
		darknodes[i] = darknode.NewDarknode(config)
		ctx[i], cancels[i] = context.WithCancel(context.Background())
	}
	return darknodes, ctxs, cancels
}
