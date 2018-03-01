package rpc

import (
	"fmt"
	"log"
	"time"

	identity "github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-order-compute"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Client struct {
	*grpc.ClientConn

	DarkNode DarkNodeClient
	Options  ClientOptions
	From     *MultiAddress
	To       *MultiAddress
}

type ClientOptions struct {
	Timeout        time.Duration
	TimeoutBackoff time.Duration
	TimeoutRetries int
}

func DefaultClientOptions() ClientOptions {
	return ClientOptions{
		Timeout:        30 * time.Second,
		TimeoutBackoff: 0 * time.Second,
		TimeoutRetries: 3,
	}
}

func NewClient(to, from identity.MultiAddress, options ...ClientOptions) (Client, error) {
	client := Client{
		Options: buildClientOptions(options...),
		From:    SerializeMultiAddress(from),
		To:      SerializeMultiAddress(to),
	}

	host, err := to.ValueForProtocol(identity.IP4Code)
	if err != nil {
		return client, err
	}
	port, err := to.ValueForProtocol(identity.TCPCode)
	if err != nil {
		return client, err
	}

	for i := 0; i < client.Options.TimeoutRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), client.Options.Timeout+(client.Options.TimeoutBackoff*time.Duration(i)))
		defer cancel()
		client.ClientConn, err = grpc.DialContext(ctx, fmt.Sprintf("%s:%s", host, port), grpc.WithInsecure())
		if err == nil {
			break
		}
		log.Println(err)
	}
	if err != nil {
		return client, err
	}
	client.DarkNode = NewDarkNodeClient(client.ClientConn)

	return client, nil
}

func (client Client) BroadcastDeltaFragment(deltaFragment *compute.DeltaFragment) (*DeltaFragment, error) {
	var resp *DeltaFragment
	var err error

	serializedDeltaFragment := SerializeDeltaFragment(deltaFragment)
	for i := 0; i < client.Options.TimeoutRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), client.Options.Timeout+(client.Options.TimeoutBackoff*time.Duration(i)))
		defer cancel()

		resp, err = client.DarkNode.BroadcastDeltaFragment(ctx, &BroadcastDeltaFragmentRequest{
			From:          client.From,
			DeltaFragment: serializedDeltaFragment,
		}, grpc.FailFast(false))
		if err == nil {
			break
		}
		log.Printf("broadcastDelaFragment %s to %s", err.Error(), client.To.Multi)
	}

	return resp, err
}

func buildClientOptions(options ...ClientOptions) ClientOptions {
	opts := DefaultClientOptions()
	for _, opt := range options {
		opts.Timeout = opt.Timeout
		opts.TimeoutBackoff = opt.TimeoutBackoff
		opts.TimeoutRetries = opt.TimeoutRetries
	}
	return opts
}
