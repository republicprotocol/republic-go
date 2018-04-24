package smpcer

import (
	"errors"
	fmt "fmt"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/rpc/client"
	"golang.org/x/net/context"
)

var ErrConnectToSelf = errors.New("connect to self")

type Client struct {
	multiAddress identity.MultiAddress
	rendezvous   Rendezvous
	streamer     Streamer
	connPool     *client.ConnPool
}

func NewClient(multiAddress identity.MultiAddress, connPool *client.ConnPool) Client {
	return Client{
		multiAddress: multiAddress,
		rendezvous:   NewRendezvous(),
		streamer:     NewStreamer(multiAddress, connPool),
		connPool:     connPool,
	}
}

func (client *Client) OpenOrder(ctx context.Context, multiAddr identity.MultiAddress, orderFragment order.Fragment) error {
	conn, err := client.connPool.Dial(ctx, multiAddr)
	if err != nil {
		return fmt.Errorf("cannot dial %v: %v", multiAddr, err)
	}
	defer conn.Close()

	smpcerClient := NewSmpcClient(conn.ClientConn)
	request := &OpenOrderRequest{
		Signature: []byte{}, // FIXME: Provide verifiable signature
		OrderFragment: &OrderFragment{
			OrderFragmenId: orderFragment.ID,
			OrderId:        orderFragment.OrderID,
			Expiry:         orderFragment.OrderExpiry.Unix(),
			Type:           int32(orderFragment.OrderType), // FIXME: Put actual token pairing here
			Tokens:         int32(0),                       // FIXME: Put actual token pairing here
			PriceShare: &Share{
				Index: orderFragment.PriceShare.Key,
				Value: []byte{}, // RSA encrption of orderFragment.PriceShare.Value
			},
			VolumeShare: &Share{
				Index: orderFragment.MaxVolumeShare.Key, // FIXME: Unify volumes
				Value: []byte{},                         // RSA encrption of orderFragment.PriceShare.Value
			},
		},
	}
	_, err = smpcerClient.OpenOrder(ctx, request)
	return err
}

func (client *Client) CloseOrder(ctx context.Context, multiAddr identity.MultiAddress orderId []byte) error {
	conn, err := client.connPool.Dial(ctx, multiAddr)
	if err != nil {
		return fmt.Errorf("cannot dial %v:%v", multiAddr, err)
	}
	defer conn.Close()

	smpcerClient := NewSmpcClient(conn.ClientConn)
	request := &CancelOrderRequest {
		Signature: []byte{}, // FIXME: Provide verifiable signature
	}
}

func (client *Client) Compute(ctx context.Context, multiAddress identity.MultiAddress, sender <-chan *ComputeMessage) (<-chan *ComputeMessage, <-chan error) {
	if client.Address() == multiAddress.Address() {
		// The Client is attempting to connect to itself
		receiver := make(chan *ComputeMessage)
		defer close(receiver)
		errs := make(chan error, 1)
		defer close(errs)
		errs <- ErrConnectToSelf
		return receiver, errs
	}

	if client.Address() < multiAddress.Address() {
		// The Client should open a gRPC stream
		return client.connect(ctx, multiAddress, sender)
	}

	// The Client must wait for the Smpc service to accept a gRPC stream from
	// a Client on another machine
	return client.wait(ctx, multiAddress, sender)
}

// Address of the Client.
func (client *Client) Address() identity.Address {
	return client.multiAddress.Address()
}

// MultiAddress of the Client.
func (client *Client) MultiAddress() identity.MultiAddress {
	return client.multiAddress
}

func (client *Client) connect(ctx context.Context, multiAddress identity.MultiAddress, sender <-chan *ComputeMessage) (<-chan *ComputeMessage, <-chan error) {
	return client.streamer.connect(ctx, multiAddress, sender)
}

func (client *Client) wait(ctx context.Context, multiAddress identity.MultiAddress, sender <-chan *ComputeMessage) (<-chan *ComputeMessage, <-chan error) {
	return client.rendezvous.wait(multiAddress.Address(), ctx.Done(), sender)
}
