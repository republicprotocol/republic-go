package smpc

import (
	"context"
	"time"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/stream"
	"github.com/republicprotocol/republic-go/swarm"
)

type Smpcer interface {
	Connect(addr identity.Address) (stream.Stream, error)
}

type Smpc struct {
	client  stream.Client
	swarmer swarm.Swarmer
}

func NewSmpc(client stream.Client, swarmer swarm.Swarmer) Smpc {
	return Smpc{
		client:  client,
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
	stream, err := smpc.client.Connect(ctx, multiAddr)
	if err != nil {
		return nil, err
	}

	return stream, nil
}
