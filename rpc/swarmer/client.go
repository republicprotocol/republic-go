package swarmer

import (
	"context"
	"fmt"
	"time"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/rpc/client"
	"github.com/republicprotocol/republic-go/rpc/dht"
)

// Client is an abstraction over the gRPC SwarmClient RPCs that implements high
// level client functionality. The Client uses a dht.DHT to store all
// identity.MultiAddresses during its interactions with Swarm services, and a
// client.ConnPool to reuse gRPC connections. The Client is used by Swarm
// services to interact with other Swarm services in the network.
type Client struct {
	multiAddress identity.MultiAddress
	dht          *dht.DHT
	connPool     *client.ConnPool
}

// NewClient returns a Client that identifies itself using an
// identity.MultiAddress and uses a dht.DHT and client.ConnPool when
// interacting with the swarm network.
func NewClient(multiAddress identity.MultiAddress, dht *dht.DHT, connPool *client.ConnPool) Client {
	return Client{
		multiAddress: multiAddress,
		dht:          dht,
		connPool:     connPool,
	}
}

// Bootstrap into the swarm network by calling Client.QueryTo for all bootstrap
// identity.MultiAddress concurrently, looking for itself. A context can be
// used to cancel or expire the bootstrap. Once all calls to Client.QueryTo
// returns, the cancellation and expiration of the Context will do nothing. A
// depth can be used to limit how far into the swarm network the Client will
// query for itself before stopping, with a negative value defining that the
// Client should not stop until all search paths are fully exhausted.
func (client *Client) Bootstrap(ctx context.Context, bootstrapMultiAddrs identity.MultiAddresses, depth int) <-chan error {
	errs := make(chan error, len(bootstrapMultiAddrs))
	go func() {
		defer close(errs)
		dispatch.CoForAll(bootstrapMultiAddrs, func(i int) {
			_, err := client.QueryTo(ctx, bootstrapMultiAddrs[i], client.Address())
			if err != nil {
				errs <- fmt.Errorf("error while bootstrapping: %v", err)
			}
		})
	}()
	return errs
}

// Query the swarm network for the identity.MultiAddress of an
// identity.Address. A context can be used to cancel or expire the query. Once
// this function returns, the cancellation and expiration of the Context will
// do nothing. The Client will use its dht.DHT to begin the query. If it has no
// known peers that are closer to the queried identity.Address than itself, the
// query will terminate immediately. A depth can be used to limit how far into
// the swarm network the Client will query before stopping, with a negative
// value defining that the Client should not stop until all search paths are
// fully exhausted.
func (client *Client) Query(ctx context.Context, query identity.Address, depth int) (identity.MultiAddress, error) {

	whitelist := identity.MultiAddresses{}
	blacklist := map[identity.Address]struct{}{}

	// Build a list of identity.MultiAddresses that are closer to the query
	// than the Swarm service
	multiAddrs := client.dht.MultiAddresses()
	for _, multiAddr := range multiAddrs {
		// Short circuit if the Swarm service is directly connected to the
		// query
		if query == multiAddr.Address() {
			return multiAddr, nil
		}

		isPeerCloser, err := identity.Closer(multiAddr.Address(), client.Address(), query)
		if err != nil {
			return identity.MultiAddress{}, fmt.Errorf("cannot compare address distances %v and %v: %v", multiAddr.Address(), client.Address(), err)
		}
		if isPeerCloser {
			whitelist = append(whitelist, multiAddr)
		}
	}

	// Search all peers for identity.MultiAddresses that are closer to the
	// query until the depth limit is reach or there are no more peers left to
	// search
	for i := 0; (i < depth || depth < 0) && len(whitelist) > 0; i++ {

		peer := whitelist[0]
		whitelist = whitelist[1:]
		if _, ok := blacklist[peer.Address()]; ok {
			continue
		}
		blacklist[peer.Address()] = struct{}{}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Query for identity.MultiAddresses that are closer to the query
		// target than the peer itself, and add them to the whitelist
		multiAddrs, err := client.QueryTo(ctx, peer, query)
		if err != nil {
			return identity.MultiAddress{}, fmt.Errorf("cannot send query to %v: %v", peer, err)
		}
		for _, multiAddr := range multiAddrs {
			whitelist = append(whitelist, multiAddr)
		}
	}

	return identity.MultiAddress{}, ErrNotFound
}

// QueryTo a peer for the identity.MultiAddress of an identity.Address. A
// context can be used to cancel or expire the query. Once this function
// returns, the cancellation and expiration of the Context will do nothing. The
// Client will use the given peer to begin the query. A depth can be used to
// limit how far into the swarm network the Client will query before stopping,
// with a negative value defining that the Client should not stop until all
// search paths are fully exhausted.
func (client *Client) QueryTo(ctx context.Context, peer identity.MultiAddress, query identity.Address) (identity.MultiAddresses, error) {
	conn, err := client.connPool.Dial(ctx, peer)
	if err != nil {
		return identity.MultiAddresses{}, fmt.Errorf("cannot dial %v: %v", peer, err)
	}
	defer conn.Close()

	// FIXME: Provide verifiable signature
	swarmClient := NewSwarmClient(conn.ClientConn)
	request := &QueryRequest{
		Signature: []byte{},
		Address:   query.String(),
	}
	stream, err := swarmClient.Query(ctx, request)
	if err != nil {
		return identity.MultiAddresses{}, err
	}
	if err := client.UpdateDHT(peer); err != nil {
		return identity.MultiAddresses{}, fmt.Errorf("cannot update dht: %v", err)
	}

	multiAddrs := identity.MultiAddresses{}
	for {
		message, err := stream.Recv()
		if err != nil {
			return multiAddrs, err
		}

		// FIXME: Verify the message signature
		signature := message.GetSignature()
		multiAddr, err := identity.NewMultiAddressFromString(message.GetMultiAddress())
		if err != nil {
			continue
		}
		if err := multiAddr.VerifySignature(signature); err != nil {
			continue
		}
		if err := client.UpdateDHT(multiAddr); err != nil {
			continue
		}
		multiAddrs = append(multiAddrs, multiAddr)
	}
}

// Address of the Client.
func (client *Client) Address() identity.Address {
	return client.multiAddress.Address()
}

// MultiAddress of the Client.
func (client *Client) MultiAddress() identity.MultiAddress {
	return client.multiAddress
}

// DHT used by the Client for storing the identity.MultiAddresses of its peers.
func (client *Client) DHT() *dht.DHT {
	return client.dht
}

// UpdateDHT with an identity.MultiAddress. If the respective dht.Bucket is
// full, then the Client will ping the oldest peer in the dht.Bucket and if
// the peer is unresponsive it will be removed to make room for the new
// identity.MultiAddress.
func (client *Client) UpdateDHT(multiAddress identity.MultiAddress) error {
	if client.multiAddress.Address() == multiAddress.Address() {
		return nil
	}
	if err := client.dht.UpdateMultiAddress(multiAddress); err != nil {
		if err == dht.ErrFullBucket {
			if client.pruneDHT(multiAddress.Address()) {
				return client.dht.UpdateMultiAddress(multiAddress)
			}
		}
		return err
	}
	return nil
}

func (client *Client) pruneDHT(addr identity.Address) bool {
	bucket, err := client.dht.FindBucket(addr)
	if err != nil {
		return false
	}
	if bucket == nil || bucket.Length() == 0 {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Ping the oldest identity.MultiAddress in the bucket and see if the
	// service is still responsive
	multiAddress := bucket.MultiAddresses[0]
	conn, err := client.connPool.Dial(ctx, multiAddress)
	if err != nil {
		client.dht.RemoveMultiAddress(multiAddress)
		return true
	}
	defer conn.Close()

	// FIXME: Provide verifiable signature
	swarmClient := NewSwarmClient(conn.ClientConn)
	request := &PingRequest{
		Signature:    []byte{},
		MultiAddress: client.multiAddress.String(),
	}
	if _, err := swarmClient.Ping(ctx, request); err != nil {
		client.dht.RemoveMultiAddress(multiAddress)
		return true
	}
	client.dht.UpdateMultiAddress(multiAddress)
	return false
}
