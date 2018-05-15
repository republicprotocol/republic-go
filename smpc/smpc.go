package smpc

import (
	"context"
	"time"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/rpc"
	"github.com/republicprotocol/republic-go/swarm"
)

type Smpcer interface {
	Connect(addr identity.Address) (rpc.Stream, error)
}

type Smpc struct {
	streamer rpc.Streamer
	swarmer  swarm.Swarmer
}

func NewSmpc(streamer rpc.Streamer, swarmer swarm.Swarmer) Smpc {
	return Smpc{
		streamer: streamer,
		swarmer:  swarmer,
	}
}

func (smpc *Smpc) Connect(addr identity.Address) (rpc.Stream, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	multiAddr, err := smpc.swarmer.Query(ctx, addr)
	if err != nil {
		return nil, err
	}

	return smpc.streamer.Connect(ctx, multiAddr)
}
