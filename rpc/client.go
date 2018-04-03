package rpc

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/republicprotocol/republic-go/identity"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// A Client is used to create and mange a gRPC connection. It provides methods
// for all RPCs and handles all timeouts and retries.
type Client struct {
	SmpcClient
	SwarmClient
	RelayClient
	SyncerClient

	Connection *grpc.ClientConn
	To         *MultiAddress
	From       *MultiAddress
	Options    ClientOptions
}

// NewClient returns a Client that is connected to the given MultiAddress and
// will always identify itself from the given MultiAddress. The connection will
// be closed when the Client is garbage collected.
func NewClient(to, from identity.MultiAddress) (*Client, error) {
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
		To:      MarshalMultiAddress(&to),
		From:    MarshalMultiAddress(&from),
	}

	if err := client.TimeoutFunc(func(ctx context.Context) error {
		connection, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%s", host, port), grpc.WithInsecure())
		if err != nil {
			return err
		}
		client.Connection = connection
		runtime.SetFinalizer(client, func(client *Client) {
			client.Connection.Close()
		})
		return nil
	}); err != nil {
		return client, err
	}

	client.SmpcClient = NewSmpcClient(client.Connection)
	client.SwarmClient = NewSwarmClient(client.Connection)
	client.RelayClient = NewRelayClient(client.Connection)
	client.SyncerClient = NewSyncerClient(client.Connection)

	return client, nil
}

// TimeoutFunc uses the timeout options of the Client to call a function. It
// returns the last error that occurred, or nil.
func (client *Client) TimeoutFunc(f func(ctx context.Context) error) error {
	var err error
	for i := 0; i < client.Options.TimeoutRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), client.Options.Timeout+(client.Options.TimeoutBackoff*time.Duration(i)))
		defer cancel()
		if err = f(ctx); err != nil {
			continue
		}
		return nil
	}
	return err
}

// StreamTimeoutFunc uses the timeout options of the Client to call a function.
// It returns the last error that occurred, or nil. If the RPC is successful it
// will not cancel the context.
func (client *Client) StreamTimeoutFunc(f func(ctx context.Context) error) error {
	var err error
	for i := 0; i < client.Options.TimeoutRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), client.Options.Timeout+(client.Options.TimeoutBackoff*time.Duration(i)))
		defer func() {
			if err := recover(); err != nil {
				cancel()
			}
		}()
		if err = f(ctx); err != nil {
			cancel()
			continue
		}
		return nil
	}
	return err
}

// Ping RPC.
func (client *Client) Ping() error {
	return client.TimeoutFunc(func(ctx context.Context) error {
		_, err := client.SwarmClient.Ping(ctx, client.To, grpc.FailFast(false))
		return err
	})
}

// QueryPeers RPC.
func (client *Client) QueryPeers(target *Address) (chan *MultiAddress, error) {

	ch := make(chan *MultiAddress)
	err := client.StreamTimeoutFunc(func(ctx context.Context) error {
		stream, err := client.SwarmClient.QueryPeers(ctx, &Query{
			From:   client.From,
			Target: target,
		}, grpc.FailFast(false))
		if err != nil {
			return err
		}
		go func() {
			defer func() { recover() }()
			for {
				multiAddress, err := stream.Recv()
				if err != nil {
					close(ch)
					return
				}
				ch <- multiAddress
			}
		}()
		return nil
	})
	return ch, err
}

// QueryPeersDeep RPC.
func (client *Client) QueryPeersDeep(target *Address) (chan *MultiAddress, error) {
	ch := make(chan *MultiAddress)
	err := client.StreamTimeoutFunc(func(ctx context.Context) error {
		stream, err := client.SwarmClient.QueryPeersDeep(ctx, &Query{
			From:   client.From,
			Target: target,
		}, grpc.FailFast(false))
		if err != nil {
			return err
		}
		go func() {
			defer func() { recover() }()
			for {
				multiAddress, err := stream.Recv()
				if err != nil {
					close(ch)
					return
				}
				ch <- multiAddress
			}
		}()
		return nil
	})
	return ch, err
}

// Sync RPC.
func (client *Client) Sync() (chan *SyncBlock, error) {
	ch := make(chan *SyncBlock)
	err := client.TimeoutFunc(func(ctx context.Context) error {
		stream, err := client.SyncerClient.Sync(ctx, &SyncRequest{
			From: client.From,
		}, grpc.FailFast(false))
		if err != nil {
			return err
		}
		go func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Recovered in f", r)
				}
			}()
			for {
				syncBlock, err := stream.Recv()
				if err != nil {
					log.Println("err receiving from the stream", err)
					close(ch)
					return
				}
				ch <- syncBlock
			}
		}()
		return nil
	})
	return ch, err
}

// SignOrderFragment RPC.
func (client *Client) SignOrderFragment(orderFragmentSignature *OrderFragmentId) (*OrderFragmentId, error) {
	var val *OrderFragmentId
	var err error
	err = client.TimeoutFunc(func(ctx context.Context) error {
		val, err = client.RelayClient.SignOrderFragment(ctx, &OrderFragmentId{
			Signature:       orderFragmentSignature.Signature,
			OrderFragmentId: orderFragmentSignature.OrderFragmentId,
		}, grpc.FailFast(false))
		return err
	})
	return val, err
}

// OpenOrder RPC.
func (client *Client) OpenOrder(openOrderRequest *OpenOrderRequest) error {
	return client.TimeoutFunc(func(ctx context.Context) error {
		_, err := client.RelayClient.OpenOrder(ctx, &OpenOrderRequest{
			From:          client.From,
			OrderFragment: openOrderRequest.OrderFragment,
		}, grpc.FailFast(false))
		return err
	})
}

// CancelOrder RPC.
func (client *Client) CancelOrder(cancelOrderRequest *CancelOrderRequest) error {
	return client.TimeoutFunc(func(ctx context.Context) error {
		_, err := client.RelayClient.CancelOrder(ctx, &CancelOrderRequest{
			From:            client.From,
			OrderFragmentId: cancelOrderRequest.OrderFragmentId,
		}, grpc.FailFast(false))
		return err
	})
}
