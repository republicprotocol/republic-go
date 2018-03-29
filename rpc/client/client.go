package client

import (
	"fmt"
	"io"
	"runtime"
	"time"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// A Client is used to create and mange a gRPC connection. It provides methods
// for all RPCs and handles all timeouts and retries.
type Client struct {
	Connection *grpc.ClientConn
	To         *rpc.MultiAddress
	From       *rpc.MultiAddress

	Options ClientOptions
	rpc.SwarmClient
	rpc.SmpcClient
	rpc.RelayClient
	rpc.SyncerClient
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
		To:      rpc.MarshalMultiAddress(&to),
		From:    rpc.MarshalMultiAddress(&from),
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

	client.SwarmClient = rpc.NewSwarmClient(client.Connection)
	client.RelayClient = rpc.NewRelayClient(client.Connection)
	client.SmpcClient = rpc.NewSmpcClient(client.Connection)
	client.SyncerClient = rpc.NewSyncerClient(client.Connection)

	return client, nil
}

// TimeoutFunc uses the timeout options of the Client to call a function. It
// returns the last error that occured, or nil.
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
// It returns the last error that occured, or nil. If the RPC is successful it
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
func (client *Client) QueryPeers(target *rpc.Address) (chan *rpc.MultiAddress, error) {

	ch := make(chan *rpc.MultiAddress)
	err := client.StreamTimeoutFunc(func(ctx context.Context) error {
		stream, err := client.SwarmClient.QueryPeers(ctx, &rpc.Query{
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
func (client *Client) QueryPeersDeep(target *rpc.Address) (chan *rpc.MultiAddress, error) {
	ch := make(chan *rpc.MultiAddress)
	err := client.StreamTimeoutFunc(func(ctx context.Context) error {
		stream, err := client.SwarmClient.QueryPeersDeep(ctx, &rpc.Query{
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
func (client *Client) Sync() (chan *rpc.SyncBlock, error) {
	ch := make(chan *rpc.SyncBlock)
	err := client.TimeoutFunc(func(ctx context.Context) error {
		stream, err := client.SyncerClient.Sync(ctx, &rpc.SyncRequest{
			From: client.From,
		}, grpc.FailFast(false))
		if err != nil {
			return err
		}
		go func() {
			defer func() { recover() }()
			for {
				syncBlock, err := stream.Recv()
				if err == io.EOF {
					close(ch)
					return
				}
				if err != nil {
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
func (client *Client) SignOrderFragment(orderFragmentSignature *rpc.OrderFragmentId) (*rpc.OrderFragmentId, error) {
	var val *rpc.OrderFragmentId
	var err error
	err = client.TimeoutFunc(func(ctx context.Context) error {
		val, err = client.RelayClient.SignOrderFragment(ctx, &rpc.OrderFragmentId{
			Signature: orderFragmentSignature.Signature,
			OrderFragmentId: orderFragmentSignature.OrderFragmentId,
		}, grpc.FailFast(false))
		return err
	})
	return val, err
}

// OpenOrder RPC.
func (client *Client) OpenOrder(openOrderRequest *rpc.OpenOrderRequest) error {
	return client.TimeoutFunc(func(ctx context.Context) error {
		_, err := client.RelayClient.OpenOrder(ctx, &rpc.OpenOrderRequest{
			From:           client.From,
			OrderFragment:  openOrderRequest.OrderFragment,
		}, grpc.FailFast(false))
		return err
	})
}

// CancelOrder RPC.
func (client *Client) CancelOrder(cancelOrderRequest *rpc.CancelOrderRequest) error {
	return client.TimeoutFunc(func(ctx context.Context) error {
		_, err := client.RelayClient.CancelOrder(ctx, &rpc.CancelOrderRequest{
			From:           client.From,
			OrderFragmentId: cancelOrderRequest.OrderFragmentId,
		}, grpc.FailFast(false))
		return err
	})
}

//// RandomFragmentShares RPC.
//func (client *Client) RandomFragmentShares() (*RandomFragments, error) {
//	var val *RandomFragments
//	var err error
//	err = client.TimeoutFunc(func(ctx context.Context) error {
//		val, err = client.DarkClient.RandomFragmentShares(ctx, &RandomFragmentSharesRequest{
//			From: client.From,
//		}, grpc.FailFast(false))
//		return err
//	})
//	return val, err
//}

//// ResidueFragmentShares RPC.
//func (client *Client) ResidueFragmentShares(randomFragments *RandomFragments) (*ResidueFragments, error) {
//	var val *ResidueFragments
//	var err error
//	err = client.TimeoutFunc(func(ctx context.Context) error {
//		val, err = client.DarkClient.ResidueFragmentShares(ctx, &ResidueFragmentSharesRequest{
//			From:            client.From,
//			RandomFragments: randomFragments,
//		}, grpc.FailFast(false))
//		return err
//	})
//	return val, err
//}

//// ComputeResidueFragment RPC.
//func (client *Client) ComputeResidueFragment(residueFragments *ResidueFragments) error {
//	return client.TimeoutFunc(func(ctx context.Context) error {
//		_, err := client.DarkClient.ComputeResidueFragment(ctx, &ComputeResidueFragmentRequest{
//			From:             client.From,
//			ResidueFragments: residueFragments,
//		}, grpc.FailFast(false))
//		return err
//	})
//}

//// BroadcastAlphaBetaFragment RPC.
//func (client *Client) BroadcastAlphaBetaFragment(alphaBetaFragment *AlphaBetaFragment) (*AlphaBetaFragment, error) {
//	var val *AlphaBetaFragment
//	var err error
//	err = client.TimeoutFunc(func(ctx context.Context) error {
//		val, err = client.DarkClient.BroadcastAlphaBetaFragment(ctx, &BroadcastAlphaBetaFragmentRequest{
//			From:              client.From,
//			AlphaBetaFragment: alphaBetaFragment,
//		}, grpc.FailFast(false))
//		return err
//	})
//	return val, err
//}

//// BroadcastDeltaFragment RPC.
//func (client *Client) BroadcastDeltaFragment(deltaFragment *DeltaFragment) (*DeltaFragment, error) {
//	var val *DeltaFragment
//	var err error
//	err = client.TimeoutFunc(func(ctx context.Context) error {
//		val, err = client.DarkClient.BroadcastDeltaFragment(ctx, &BroadcastDeltaFragmentRequest{
//			From:          client.From,
//			DeltaFragment: deltaFragment,
//		}, grpc.FailFast(false))
//		return err
//	})
//	return val, err
//}
