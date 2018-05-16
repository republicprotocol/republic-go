package smpc

import (
	"context"
	"time"

	"github.com/republicprotocol/republic-go/shamir"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/smpc/delta"
	"github.com/republicprotocol/republic-go/stream"
	"github.com/republicprotocol/republic-go/swarm"
)

type Smpcer interface {
	Input() chan<- Inst
	Output() <-chan Compute
}

type Smpc struct {
	client  stream.Client
	server  stream.Server
	swarmer swarm.Swarmer
}

func NewSmpc(client stream.Client, server stream.Server, swarmer swarm.Swarmer) Smpc {
	return Smpc{
		client:  client,
		server:  server,
		swarmer: swarmer,
	}
}

func (smpc *Smpc) Connect(addr identity.Address) (stream.Stream, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	multiAddr, err := smpc.swarmer.Query(ctx, addr, 3)
	if err != nil {
		return nil, err
	}

	// Compare address and decide to connect to or listen for connection.
	var stream stream.Stream
	if multiAddr.Address() > smpc.Address() {
		stream, err = smpc.client.Connect(ctx, multiAddr)
	} else {
		stream, err = smpc.server.Listen(ctx, multiAddr.Address())
	}

	return stream, nil
}
