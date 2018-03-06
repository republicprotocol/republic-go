package rpc

import (
	"fmt"
	"io"
	"time"

	"github.com/republicprotocol/go-order-compute"
	"github.com/republicprotocol/republic-go/identity"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

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
	if err := client.TimeoutFunc(func(ctx context.Context) error {
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

func (client *Client) Ping() error {
	return client.TimeoutFunc(func(ctx context.Context) error {
		_, err := client.Swarm.Ping(ctx, client.To, grpc.FailFast(false))
		return err
	})
}

func (client *Client) Query(query identity.Address) (identity.Multiaddresses, error) {
	multiaddresses := make(identity.Multiaddresses, 0)
	err := client.TimeoutFunc(func(ctx context.Context) error {
		stream, err := client.Swarm.Query(ctx, SerializeAddress(query), grpc.FailFast(false))
		if err != nil {
			return err
		}
		for {
			multiaddress, err := stream.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
			deserializedMultiaddress, err := DeserializeMultiaddress(multiaddress)
			if err != nil {
				return err
			}
			multiaddresses = append(multiaddresses, deserializedMultiaddress)
		}
	})
	return multiaddresses, err
}

func (client *Client) QueryDeep(query identity.Address) (identity.Multiaddresses, error) {
	multiaddresses := make(identity.Multiaddresses, 0)
	err := client.TimeoutFunc(func(ctx context.Context) error {
		stream, err := client.Swarm.QueryDeep(ctx, SerializeAddress(query), grpc.FailFast(false))
		if err != nil {
			return err
		}
		for {
			multiaddress, err := stream.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
			deserializedMultiaddress, err := DeserializeMultiaddress(multiaddress)
			if err != nil {
				return err
			}
			multiaddresses = append(multiaddresses, deserializedMultiaddress)
		}
	})
	return multiaddresses, err
}

func (client *Client) BroadcastDeltaFragment(deltaFragment *compute.DeltaFragment) (*DeltaFragment, error) {
	var response *DeltaFragment
	var err error

	serializedDeltaFragment := SerializeDeltaFragment(deltaFragment)
	err = client.TimeoutFunc(func(ctx context.Context) error {
		response, err = client.DarkOcean.BroadcastDeltaFragment(ctx, &BroadcastDeltaFragmentRequest{
			From:          client.From,
			DeltaFragment: serializedDeltaFragment,
		}, grpc.FailFast(false))
		return err
	})

	return response, err
}

// TimeoutFunc uses the timeout options of the Client to call a function. It
// returns the last error that occured, or nil.
func (client *Client) TimeoutFunc(f func(ctx context.Context) error) (err error) {
	for i := 0; i < client.Options.TimeoutRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), client.Options.Timeout+(client.Options.TimeoutBackoff*time.Duration(i)))
		defer cancel()
		if err = f(ctx); err == nil {
			return
		}
	}
}
