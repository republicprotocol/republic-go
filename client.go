package rpc

import (
	"fmt"
	"time"

	identity "github.com/republicprotocol/go-identity"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Client struct {
	*grpc.ClientConn

	Options ClientOptions
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

func NewClient(multiAddress identity.MultiAddress) (Client, error) {
	client := Client{
		Options: DefaultClientOptions(),
	}

	host, err := multiAddress.ValueForProtocol(identity.IP4Code)
	if err != nil {
		return client, err
	}
	port, err := multiAddress.ValueForProtocol(identity.TCPCode)
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
	}

	return client, nil
}

func (client Client) BroadcastDeltaFragment(deltaFragment *DeltaFragment) (*DeltaFragment, error) {

	var resp *DeltaFragment
	var err error

	serializedDeltaFragment = SerializeDeltaFragment(deltaFragment)
	for i := 0; i < client.Options.TimeoutRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), client.Options.Timeout+(client.Options.TimeoutBackoff*time.Duration(i)))
		defer cancel()

		resp, err = client.BroadcastDeltaFragment()
	}

	serializedOrderFragment := SerializeOrderFragment(orderFragment)
	serializedOrderFragment.To = SerializeAddress(to)
	serializedOrderFragment.From = SerializeMultiAddress(from)
	_, err = client.SendOrderFragment(ctx, serializedOrderFragment, grpc.FailFast(false))
	return nil, err
}
