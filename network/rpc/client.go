package rpc

import (
	"fmt"
	"time"

	"github.com/republicprotocol/go-order-compute"
	"github.com/republicprotocol/republic-go/identity"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// A Client is used to send RPCs to services.
type Client struct {
	Connection *grpc.ClientConn
	To         *Multiaddress
	From       *Multiaddress

	Options   ClientOptions
	Swarm     SwarmClient
	DarkOcean DarkOceanClient
}

func NewClient(to, from identity.Multiaddress) (*Client, error) {
	host, err := to.ValueForProtocol(identity.IP4Code)
	if err != nil {
		return nil, err
	}
	port, err := to.ValueForProtocol(identity.TCPCode)
	if err != nil {
		return nil, err
	}

	client := &Client{
		Options: DefaultClientOptions(),
		To:      SerializeMultiaddress(to),
		From:    SerializeMultiaddress(from),
	}
	if err := client.timeout(func(ctx context.Context) error {
		connection, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%s", host, port), grpc.WithInsecure())
		if err != nil {
			return err
		}
		client.Connection = connection
		return nil
	}); err != nil {
		return client, err
	}

	client.Swarm = NewSwarmClient(client.Connection)
	client.DarkOcean = NewDarkOceanClient(client.Connection)
	return client, nil
}

func (client *Client) BroadcastDeltaFragment(deltaFragment *compute.DeltaFragment) (*DeltaFragment, error) {
	var response *DeltaFragment
	var err error

	serializedDeltaFragment := SerializeDeltaFragment(deltaFragment)
	err = client.timeout(func(ctx context.Context) error {
		response, err = client.DarkOcean.BroadcastDeltaFragment(ctx, &BroadcastDeltaFragmentRequest{
			From:          client.From,
			DeltaFragment: serializedDeltaFragment,
		}, grpc.FailFast(false))
		return err
	})

	return response, err
}

func (client *Client) timeout(f func(ctx context.Context) error) error {
	var err error
	for i := 0; i < client.Options.TimeoutRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), client.Options.Timeout+(client.Options.TimeoutBackoff*time.Duration(i)))
		defer cancel()
		if err = f(ctx); err == nil {
			break
		}
	}
	return err
}
