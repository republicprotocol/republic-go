package rpc

import (
	"fmt"
	"io"
	"runtime"
	"time"

	"github.com/republicprotocol/republic-go/identity"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// A Client is used to create and mange a gRPC connection. It provides methods
// for all RPCs and handles all timeouts and retries.
type Client struct {
	Connection *grpc.ClientConn
	To         *MultiAddress
	From       *MultiAddress

	Options ClientOptions
	SwarmClient
	DarkClient
	GossipClient
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
		To:      SerializeMultiAddress(to),
		From:    SerializeMultiAddress(from),
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

	client.SwarmClient = NewSwarmClient(client.Connection)
	client.DarkClient = NewDarkClient(client.Connection)
	client.GossipClient = NewGossipClient(client.Connection)

	return client, nil
}

// TimeoutFunc uses the timeout options of the Client to call a function. It
// returns the last error that occured, or nil.
func (client *Client) TimeoutFunc(f func(ctx context.Context) error) error {
	var err error
	for i := 0; i < client.Options.TimeoutRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), client.Options.Timeout+(client.Options.TimeoutBackoff*time.Duration(i)))
		defer cancel()
		if err = f(ctx); err == nil {
			return err
		}
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
	err := client.TimeoutFunc(func(ctx context.Context) error {
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
				if err == io.EOF {
					close(ch)
					return
				}
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

// QueryPeersBlocking calls the QueryPeers RPC and blocks until the RPC is
// finished.
func (client *Client) QueryPeersBlocking(target *Address) ([]*MultiAddress, error) {
	multiAddresses := make([]*MultiAddress, 0)
	err := client.TimeoutFunc(func(ctx context.Context) error {
		stream, err := client.SwarmClient.QueryPeers(ctx, &Query{
			From:   client.From,
			Target: target,
		}, grpc.FailFast(false))
		if err != nil {
			return err
		}
		for {
			multiAddress, err := stream.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
			multiAddresses = append(multiAddresses, multiAddress)
		}
	})
	return multiAddresses, err
}

// QueryPeersDeep RPC.
func (client *Client) QueryPeersDeep(target *Address) (chan *MultiAddress, error) {
	ch := make(chan *MultiAddress)
	err := client.TimeoutFunc(func(ctx context.Context) error {
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
				if err == io.EOF {
					close(ch)
					return
				}
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

// QueryPeersDeepBlocking calls the QueryPeersDeep RPC and blocks until the RPC
// is finished.
func (client *Client) QueryPeersDeepBlocking(target *Address) ([]*MultiAddress, error) {
	multiAddresses := make([]*MultiAddress, 0)
	err := client.TimeoutFunc(func(ctx context.Context) error {
		stream, err := client.SwarmClient.QueryPeersDeep(ctx, &Query{
			From:   client.From,
			Target: target,
		}, grpc.FailFast(false))
		if err != nil {
			return err
		}
		for {
			multiAddress, err := stream.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
			multiAddresses = append(multiAddresses, multiAddress)
		}
	})
	return multiAddresses, err
}

// Sync RPC.
func (client *Client) Sync() (chan *SyncBlock, error) {
	ch := make(chan *SyncBlock)
	err := client.TimeoutFunc(func(ctx context.Context) error {
		stream, err := client.DarkClient.Sync(ctx, &SyncRequest{
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
func (client *Client) SignOrderFragment(orderFragmentSignature *OrderFragmentSignature) (*OrderFragmentSignature, error) {
	var val *OrderFragmentSignature
	var err error
	err = client.TimeoutFunc(func(ctx context.Context) error {
		val, err = client.DarkClient.SignOrderFragment(ctx, &SignOrderFragmentRequest{
			From: client.From,
			OrderFragmentSignature: orderFragmentSignature,
		}, grpc.FailFast(false))
		return err
	})
	return val, err
}

// OpenOrder RPC.
func (client *Client) OpenOrder(orderSignature *OrderSignature, orderFragment *OrderFragment) error {
	return client.TimeoutFunc(func(ctx context.Context) error {
		_, err := client.DarkClient.OpenOrder(ctx, &OpenOrderRequest{
			From:           client.From,
			OrderSignature: orderSignature,
			OrderFragment:  orderFragment,
		}, grpc.FailFast(false))
		return err
	})
}

// CancelOrder RPC.
func (client *Client) CancelOrder(orderSignature *OrderSignature) error {
	return client.TimeoutFunc(func(ctx context.Context) error {
		_, err := client.DarkClient.CancelOrder(ctx, &CancelOrderRequest{
			From:           client.From,
			OrderSignature: orderSignature,
		}, grpc.FailFast(false))
		return err
	})
}

// RandomFragmentShares RPC.
func (client *Client) RandomFragmentShares() (*RandomFragments, error) {
	var val *RandomFragments
	var err error
	err = client.TimeoutFunc(func(ctx context.Context) error {
		val, err = client.DarkClient.RandomFragmentShares(ctx, &RandomFragmentSharesRequest{
			From: client.From,
		}, grpc.FailFast(false))
		return err
	})
	return val, err
}

// ResidueFragmentShares RPC.
func (client *Client) ResidueFragmentShares(randomFragments *RandomFragments) (*ResidueFragments, error) {
	var val *ResidueFragments
	var err error
	err = client.TimeoutFunc(func(ctx context.Context) error {
		val, err = client.DarkClient.ResidueFragmentShares(ctx, &ResidueFragmentSharesRequest{
			From:            client.From,
			RandomFragments: randomFragments,
		}, grpc.FailFast(false))
		return err
	})
	return val, err
}

// ComputeResidueFragment RPC.
func (client *Client) ComputeResidueFragment(residueFragments *ResidueFragments) error {
	return client.TimeoutFunc(func(ctx context.Context) error {
		_, err := client.DarkClient.ComputeResidueFragment(ctx, &ComputeResidueFragmentRequest{
			From:             client.From,
			ResidueFragments: residueFragments,
		}, grpc.FailFast(false))
		return err
	})
}

// BroadcastAlphaBetaFragment RPC.
func (client *Client) BroadcastAlphaBetaFragment(alphaBetaFragment *AlphaBetaFragment) (*AlphaBetaFragment, error) {
	var val *AlphaBetaFragment
	var err error
	err = client.TimeoutFunc(func(ctx context.Context) error {
		val, err = client.DarkClient.BroadcastAlphaBetaFragment(ctx, &BroadcastAlphaBetaFragmentRequest{
			From:              client.From,
			AlphaBetaFragment: alphaBetaFragment,
		}, grpc.FailFast(false))
		return err
	})
	return val, err
}

// BroadcastDeltaFragment RPC.
func (client *Client) BroadcastDeltaFragment(deltaFragment *DeltaFragment) (*DeltaFragment, error) {
	var val *DeltaFragment
	var err error
	err = client.TimeoutFunc(func(ctx context.Context) error {
		val, err = client.DarkClient.BroadcastDeltaFragment(ctx, &BroadcastDeltaFragmentRequest{
			From:          client.From,
			DeltaFragment: deltaFragment,
		}, grpc.FailFast(false))
		return err
	})
	return val, err
}

// Gossip RPC.
func (client *Client) Gossip(rumor *Rumor) (*Rumor, error) {
	var val *Rumor
	var err error
	err = client.TimeoutFunc(func(ctx context.Context) error {
		val, err = client.GossipClient.Gossip(ctx, &GossipRequest{
			From:  client.From,
			Rumor: rumor,
		}, grpc.FailFast(false))
		return err
	})
	return val, err
}

// Finalize RPC.
func (client *Client) Finalize(rumor *Rumor) (*Rumor, error) {
	var val *Rumor
	var err error
	err = client.TimeoutFunc(func(ctx context.Context) error {
		val, err = client.GossipClient.Finalize(ctx, &FinalizeRequest{
			From:  client.From,
			Rumor: rumor,
		}, grpc.FailFast(false))
		return err
	})
	return val, err
}
