package smpcer

import (
	"fmt"
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
	rcs       map[identity.Address]int
	dones     map[identity.Address]chan struct{}
	senders   map[identity.Address]*dispatch.Broadcaster
	receivers map[identity.Address]*dispatch.Broadcaster
	errs      map[identity.Address]*dispatch.Broadcaster
}

func NewStreamer(multiAddress identity.MultiAddress, connPool *client.ConnPool) Streamer {
	return Streamer{
		multiAddress: multiAddress,
		connPool:     connPool,

		mu:        new(sync.Mutex),
		rcs:       map[identity.Address]int{},
		dones:     map[identity.Address]chan struct{}{},
		senders:   map[identity.Address]*dispatch.Broadcaster{},
		receivers: map[identity.Address]*dispatch.Broadcaster{},
		errs:      map[identity.Address]*dispatch.Broadcaster{},
	}
}

func (streamer *Streamer) connect(multiAddr identity.MultiAddress, done <-chan struct{}, sender <-chan interface{}) (<-chan interface{}, <-chan interface{}) {
	receiver, errs := streamer.acquireConn(multiAddr, sender)
	go func() {
		defer streamer.releaseConn(multiAddr.Address())
		<-done
	}()
	return receiver, errs
}

func (streamer *Streamer) acquireConn(multiAddr identity.MultiAddress, sender <-chan interface{}) (<-chan interface{}, <-chan interface{}) {
	streamer.mu.Lock()
	defer streamer.mu.Unlock()

	addr := multiAddr.Address()

	streamer.rcs[addr]++
	if streamer.rcs[addr] == 1 {
		streamer.dones[addr] = make(chan struct{})
		streamer.senders[addr] = dispatch.NewBroadcaster()
		streamer.receivers[addr] = dispatch.NewBroadcaster()
		streamer.errs[addr] = dispatch.NewBroadcaster()

		sender := streamer.senders[addr].Listen(streamer.dones[addr])
		receiver, errs := streamer.openGrpcStream(multiAddr, streamer.dones[addr], sender)
		go streamer.receivers[addr].Broadcast(streamer.dones[addr], receiver)
		go streamer.errs[addr].Broadcast(streamer.dones[addr], errs)
	}

	go streamer.senders[addr].Broadcast(streamer.dones[addr], sender)
	receiver := streamer.receivers[addr].Listen(streamer.dones[addr])
	errs := streamer.errs[addr].Listen(streamer.dones[addr])

	return receiver, errs
}

func (streamer *Streamer) releaseConn(addr identity.Address) {
	streamer.mu.Lock()
	defer streamer.mu.Unlock()

	streamer.rcs[addr]--
	if streamer.rcs[addr] == 0 {
		close(streamer.dones[addr])
		streamer.senders[addr].Close()
		streamer.receivers[addr].Close()
		streamer.errs[addr].Close()
		delete(streamer.dones, addr)
		delete(streamer.senders, addr)
		delete(streamer.receivers, addr)
		delete(streamer.errs, addr)
	}
}

func (streamer *Streamer) openGrpcStream(multiAddr identity.MultiAddress, done <-chan struct{}, sender <-chan interface{}) (<-chan interface{}, <-chan interface{}) {
	receiver := make(chan interface{})
	errs := make(chan interface{}, 2)

	go func() {
		defer close(receiver)
		defer close(errs)

		grpcStream, err := streamer.openGrpcCompute(context.Background(), multiAddr)
		if err != nil {
			errs <- err
			return
		}

		dispatch.CoBegin(func() {
			// Writing messages from the senders to the grpc client
			for {
				select {
				case <-done:
					return
				case <-grpcStream.Context().Done():
					errs <- grpcStream.Context().Err()
					return
				case msg, ok := <-sender:
					if !ok {
						return
					}
					if msg, ok := msg.(*ComputeMessage); ok {
						if err := grpcStream.Send(msg); err != nil {
							errs <- err
							return
						}
					}
				}
			}
		}, func() {
			// Writing messages from the grpc client to the receivers
			for {
				msg, err := grpcStream.Recv()
				if err != nil {
					errs <- err
					return
				}
				select {
				case <-done:
					return
				case <-grpcStream.Context().Done():
					errs <- grpcStream.Context().Err()
					return
				case receiver <- msg:
				}
			}
		})
	}()

	return receiver, errs
}

func (streamer *Streamer) openGrpcCompute(ctx context.Context, multiAddress identity.MultiAddress) (Smpc_ComputeClient, error) {

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
