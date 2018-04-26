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

	mu          *sync.Mutex
	senders     map[identity.Address]chan *ComputeMessage
	receivers   map[identity.Address]*dispatch.Splitter
	receiverChs map[identity.Address]chan *ComputeMessage
	errs        map[identity.Address]*dispatch.Splitter
	errChs      map[identity.Address]chan error
	rcs         map[identity.Address]int
}

func NewStreamer(multiAddress identity.MultiAddress, connPool *client.ConnPool) Streamer {
	return Streamer{
		multiAddress: multiAddress,
		connPool:     connPool,

		mu:          new(sync.Mutex),
		senders:     map[identity.Address]chan *ComputeMessage{},
		receivers:   map[identity.Address]*dispatch.Splitter{},
		receiverChs: map[identity.Address]chan *ComputeMessage{},
		errs:        map[identity.Address]*dispatch.Splitter{},
		errChs:      map[identity.Address]chan error{},
		rcs:         map[identity.Address]int{},
	}
}

func (streamer *Streamer) connect(multiAddr identity.MultiAddress, done <-chan struct{}, sender <-chan *ComputeMessage) (<-chan *ComputeMessage, <-chan error) {
	receiver := make(chan *ComputeMessage)
	errs := make(chan error, 1)
	go func() {
		defer close(receiver)
		defer close(errs)

		addr := multiAddr.Address()

		streamer.acquireConn(multiAddr)
		defer streamer.releaseConn(addr)

		if err := streamer.receivers[addr].Subscribe(receiver); err != nil {
			errs <- err
			return
		}
		defer streamer.receivers[addr].Unsubscribe(receiver)

		if err := streamer.errs[addr].Subscribe(errs); err != nil {
			errs <- err
			return
		}
		defer streamer.errs[addr].Unsubscribe(errs)

		dispatch.Pipe(done, sender, streamer.senders[addr])
	}()
	return receiver, errs
}

func (streamer *Streamer) acquireConn(multiAddr identity.MultiAddress) {
	streamer.mu.Lock()
	defer streamer.mu.Unlock()

	addr := multiAddr.Address()

	if streamer.rcs[addr] == 0 {
		streamer.rcs[addr]++
		streamer.senders[addr] = make(chan *ComputeMessage)
		streamer.receivers[addr] = dispatch.NewSplitter(MaxConnections)
		streamer.receiverChs[addr] = make(chan *ComputeMessage)
		go streamer.receivers[addr].Split(streamer.receiverChs[addr])
		streamer.errs[addr] = dispatch.NewSplitter(MaxConnections)
		streamer.errChs[addr] = make(chan error)
		go streamer.errs[addr].Split(streamer.stream(multiAddr))
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
		delete(streamer.receiverChs, addr)
		delete(streamer.rcs, addr)
		delete(streamer.errs, addr)
		delete(streamer.errChs, addr)
	}
}

func (streamer *Streamer) stream(multiAddress identity.MultiAddress) <-chan error {
	errs := make(chan error, 1)
	go func() {
		defer close(errs)

		stream, err := streamer.compute(context.Background(), multiAddress)
		if err != nil {
			errs <- err
			return
		}

		addr := multiAddress.Address()

		done := make(chan struct{})
		defer close(done)

		// Read all messages from the sender channel and write them to the gRPC
		// stream
		senderErrs := make(chan error, 1)
		go func() {
			defer close(senderErrs)
			for {
				select {
				case <-done:
					// The receiver has terminated
					return
				case <-stream.Context().Done():
					// Writing to the error channel will cause the receiver to
					// terminate
					senderErrs <- stream.Context().Err()
					return
				case message, ok := <-streamer.senders[addr]:
					if !ok {
						return
					}
					if err := stream.Send(message); err != nil {
						senderErrs <- err
						return
					}
				}
			}
		}()

		// Read all messages from the gRPC stream and write them to the receiver
		// channel
		for {
			message, err := stream.Recv()
			if err != nil {
				errs <- err
				return
			}
			select {
			case <-stream.Context().Done():
				errs <- stream.Context().Err()
				return
			case err, ok := <-errs:
				// The sender has terminated, possibly with an error that should be
				// returned
				if !ok {
					return
				}
				errs <- err
				return
			case streamer.receiverChs[addr] <- message:
			}
		}
	}()
	return errs
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
