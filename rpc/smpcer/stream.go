package smpcer

import (
	fmt "fmt"
	"sync"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/rpc/client"
	"golang.org/x/net/context"
)

type Streamer struct {
	multiAddress identity.MultiAddress
	connPool     *client.ConnPool

	mu        *sync.Mutex
	senders   map[identity.Address]chan *ComputeMessage
	receivers map[identity.Address]*dispatch.Splitter
	errs      map[identity.Address]*dispatch.Splitter
	rcs       map[identity.Address]int
}

func NewStreamer(multiAddress identity.MultiAddress, connPool *client.ConnPool) Streamer {
	return Streamer{
		multiAddress: multiAddress,
		connPool:     connPool,

		mu:        new(sync.Mutex),
		senders:   map[identity.Address]chan *ComputeMessage{},
		receivers: map[identity.Address]*dispatch.Splitter{},
		errs:      map[identity.Address]*dispatch.Splitter{},
		rcs:       map[identity.Address]int{},
	}
}

func (streamer *Streamer) connect(ctx context.Context, multiAddr identity.MultiAddress, sender <-chan *ComputeMessage) (<-chan *ComputeMessage, <-chan error) {
	panic("unimplemented")
}

func (streamer *Streamer) acquireConn(multiAddr identity.MultiAddress) {
	streamer.mu.Lock()
	defer streamer.mu.Unlock()

	addr := multiAddr.Address()

	if streamer.rcs[addr] == 0 {
		streamer.rcs[addr]++
		sender := make(chan *ComputeMessage)
		streamer.senders[addr] = sender
		streamer.receivers[addr] = dispatch.NewSplitter(MaxConnections)
		go streamer.stream(multiAddr)
	}
	streamer.rcs[addr]++
}

func (streamer *Streamer) releaseConn(addr identity.Address) {
	streamer.mu.Lock()
	defer streamer.mu.Unlock()

	streamer.rcs[addr]--
	if streamer.rcs[addr] == 0 {
		close(streamer.senders[addr])
		delete(streamer.senders, addr)
		delete(streamer.receivers, addr)
		delete(streamer.errs, addr)
	}
}

func (streamer *Streamer) stream(multiAddress identity.MultiAddress) {
	panic("unimplemented")
}

func (streamer *Streamer) compute(ctx context.Context, multiAddress identity.MultiAddress) (Smpc_ComputeClient, error) {

	// Dial the client.ConnPool for a client.Conn
	conn, err := streamer.connPool.Dial(ctx, multiAddress)
	if err != nil {
		return nil, fmt.Errorf("cannot dial %v: %v", multiAddress, err)
	}
	defer conn.Close()

	// Create an SmpcClient and call the Compute RPC
	smpcClient := NewSmpcClient(conn.ClientConn)
	stream, err := smpcClient.Compute(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot open stream to %v: %v", multiAddress, err)
	}

	// Send an authentication message to the Smpc service
	auth := &ComputeMessage{
		Signature: []byte{}, // FIXME: Provide verifiable signature
		Value: &ComputeMessage_Address{
			Address: streamer.multiAddress.Address().String(),
		},
	}
	if err := stream.Send(auth); err != nil {
		return nil, fmt.Errorf("cannot authenticate with %v: %v", multiAddress, err)
	}

	return stream, nil
}
