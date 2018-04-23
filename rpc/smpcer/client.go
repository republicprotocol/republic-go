package smpcer

import (
	"errors"
	fmt "fmt"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/rpc/client"
	"golang.org/x/net/context"
)

var ErrConnectToSelf = errors.New("connect to self")

type Client struct {
	multiAddress identity.MultiAddress
	router       *RendezvousRouter
	connPool     *client.ConnPool
}

func NewClient(multiAddress identity.MultiAddress, router *RendezvousRouter, connPool *client.ConnPool) Client {
	return Client{
		multiAddress: multiAddress,
		router:       router,
		connPool:     connPool,
	}
}

func (client *Client) Compute(ctx context.Context, multiAddress identity.MultiAddress, sender <-chan *ComputeMessage) (<-chan *ComputeMessage, <-chan error) {
	if client.Address() == multiAddress.Address() {
		receiver := make(chan *ComputeMessage)
		errs := make(chan error, 1)
		defer close(receiver)
		defer close(errs)
		errs <- ErrConnectToSelf
		return receiver, errs
	}
	if client.Address() < multiAddress.Address() {
		return client.connect(ctx, multiAddress, sender)
	}
	return client.waitForConnect(ctx, multiAddress, sender)
}

// Address of the Client.
func (client *Client) Address() identity.Address {
	return client.multiAddress.Address()
}

// MultiAddress of the Client.
func (client *Client) MultiAddress() identity.MultiAddress {
	return client.multiAddress
}

// RendezvousRouter used by the Client to connect synchronously with Smpc
// services.
func (client *Client) RendezvousRouter() *RendezvousRouter {
	return client.router
}

func (client *Client) connect(ctx context.Context, multiAddress identity.MultiAddress, sender <-chan *ComputeMessage) (<-chan *ComputeMessage, <-chan error) {
	receiver := make(chan *ComputeMessage)
	errs := make(chan error, 1)
	go func() {
		defer close(receiver)
		defer close(errs)

		conn, err := client.connPool.Dial(ctx, multiAddress)
		if err != nil {
			errs <- fmt.Errorf("cannot dial %v: %v", multiAddress, err)
			return
		}
		defer conn.Close()

		// FIXME: Provide verifiable signature
		smpcClient := NewSmpcClient(conn.ClientConn)
		stream, err := smpcClient.Compute(ctx)
		if err != nil {
			errs <- fmt.Errorf("cannot open stream to %v: %v", multiAddress, err)
			return
		}
		auth := &ComputeMessage{
			Signature: []byte{},
			Value: &ComputeMessage_Address{
				Address: client.Address().String(),
			},
		}
		if err := stream.Send(auth); err != nil {
			errs <- fmt.Errorf("cannot authenticate with %v: %v", multiAddress, err)
			return
		}

		rendezvous := client.router.Acquire(multiAddress.Address().String())
		defer client.router.Release(multiAddress.Address().String())

		if err := client.mergeStreamAndRendezvous(stream, rendezvous); err != nil {
			errs <- err
			return
		}
	}()
	return receiver, errs
}

func (client *Client) waitForConnect(ctx context.Context, multiAddress identity.MultiAddress, sender <-chan *ComputeMessage) (<-chan *ComputeMessage, <-chan error) {
	receiver := make(chan *ComputeMessage)
	errs := make(chan error, 1)
	go func() {
		defer close(receiver)
		defer close(errs)

		rendezvous := client.router.Acquire(multiAddress.Address().String())
		defer client.router.Release(multiAddress.Address().String())

		for {
			select {
			case <-ctx.Done():
				errs <- ctx.Err()
				return
			case message, ok := <-sender:
				if !ok {
					return
				}
				select {
				case <-ctx.Done():
					errs <- ctx.Err()
					return
				case rendezvous.Sender <- message:
				}
			case message, ok := <-rendezvous.Receiver:
				if !ok {
					return
				}
				select {
				case <-ctx.Done():
					errs <- ctx.Err()
					return
				case receiver <- message:
				}
			}
		}
	}()
	return receiver, errs
}
