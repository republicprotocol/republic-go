package swarmer

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/republicprotocol/republic-go/crypto"
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
	crypter      crypto.Crypter
	multiAddress identity.MultiAddress
	dht          *dht.DHT
	connPool     *client.ConnPool
	bootstrapped bool
}

// NewClient returns a Client that identifies itself using an
// identity.MultiAddress and uses a dht.DHT and client.ConnPool when
// interacting with the swarm network.
func NewClient(crypter crypto.Crypter, multiAddress identity.MultiAddress, dht *dht.DHT, connPool *client.ConnPool) Client {
	return Client{
		crypter:      crypter,
		multiAddress: multiAddress,
		dht:          dht,
		connPool:     connPool,
		bootstrapped: false,
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
	errs := make(chan error, 2*len(bootstrapMultiAddrs)+1)
	go func() {
		defer close(errs)
		defer func() {
			client.bootstrapped = true
		}()
		dispatch.CoForAll(bootstrapMultiAddrs, func(i int) {
			if err := client.Ping(ctx, bootstrapMultiAddrs[i]); err != nil {
				errs <- fmt.Errorf("cannot ping bootstrap node %v: %v", bootstrapMultiAddrs[i], err)
			}
		})
		_, err := client.query(ctx, client.Address(), depth, true)
		if err != nil {
			errs <- fmt.Errorf("error while bootstrapping: %v", err)
		}
	}()
	return errs
}

// Ping a peer in the swarm network. A context can be used to cancel or expire
// the ping. Once this function returns, the cancellation and expiration of the
// Context will do nothing.
func (client *Client) Ping(ctx context.Context, peer identity.MultiAddress) error {
	conn, err := client.connPool.Dial(ctx, peer)
	if err != nil {
		return fmt.Errorf("cannot dial %v: %v", peer, err)
	}
	defer conn.Close()

	swarmClient := NewSwarmClient(conn.ClientConn)
	requestSignature, err := client.crypter.Sign(client.MultiAddress())
	if err != nil {
		return fmt.Errorf("cannot sign request: %v", err)
	}
	request := &PingRequest{
		Signature:    requestSignature,
		MultiAddress: client.MultiAddress().String(),
	}

	response, err := swarmClient.Ping(ctx, request)
	if err != nil && err != io.EOF {
		return err
	}
	multiAddr, err := identity.NewMultiAddressFromString(response.GetMultiAddress())
	if err != nil {
		return fmt.Errorf("cannot parse %v: %v", response.GetMultiAddress(), err)
	}
	multiAddr.Signature = response.GetSignature()
	if err := client.crypter.Verify(multiAddr, multiAddr.Signature); err != nil {
		return fmt.Errorf("cannot verify signature of %v: %v", response.GetMultiAddress(), err)
	}
	if err := client.UpdateDHT(multiAddr); err != nil {
		return fmt.Errorf("cannot store %v in dht: %v", response.GetMultiAddress(), err)
	}
	return nil
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
	return client.query(ctx, query, depth, false)
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

	swarmClient := NewSwarmClient(conn.ClientConn)
	requestSignature, err := client.crypter.Sign(query)
	if err != nil {
		return identity.MultiAddresses{}, fmt.Errorf("cannot sign request: %v", err)
	}
	request := &QueryRequest{
		Signature: requestSignature,
		Address:   query.String(),
	}

	// Regardless of the success of the ping, continue querying the peer
	_ = client.Ping(ctx, peer)

	stream, err := swarmClient.Query(ctx, request)
	if err != nil && err != io.EOF {
		return identity.MultiAddresses{}, err
	}
	if err := client.UpdateDHT(peer); err != nil {
		return identity.MultiAddresses{}, fmt.Errorf("cannot update dht: %v", err)
	}

	multiAddrs := identity.MultiAddresses{}
	for {
		message, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return multiAddrs, nil
			}
			return multiAddrs, err
		}

		multiAddr, err := identity.NewMultiAddressFromString(message.GetMultiAddress())
		if err != nil {
			log.Printf("cannot parse %v: %v", message.GetMultiAddress(), err)
			continue
		}
		multiAddr.Signature = message.GetSignature()
		if err := client.crypter.Verify(multiAddr, multiAddr.Signature); err != nil {
			log.Printf("cannot verify signature of %v: %v", message.GetMultiAddress(), err)
			continue
		}
		if err := client.UpdateDHT(multiAddr); err != nil {
			log.Printf("cannot store %v in dht: %v", message.GetMultiAddress(), err)
			continue
		}
		multiAddrs = append(multiAddrs, multiAddr)
	}
}

// Address of the Client.
func (client *Client) Address() identity.Address {
	return client.multiAddress.Address()
}

// Bootstrapped shows whether the node has finished bootstrapping.
func (client *Client) Bootstrapped() bool {
	return client.bootstrapped
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

func (client *Client) query(ctx context.Context, query identity.Address, depth int, isBootstrapping bool) (identity.MultiAddress, error) {
	whitelist := identity.MultiAddresses{}
	blacklist := map[identity.Address]struct{}{}

	// Build a list of identity.MultiAddresses that are closer to the query
	// than the Swarm service
	multiAddrs := client.dht.MultiAddresses()
	for _, multiAddr := range multiAddrs {
		// Short circuit if the Swarm service is directly connected to the
		// query
		if query == multiAddr.Address() && !isBootstrapping {
			return multiAddr, nil
		}

		isPeerCloser, err := identity.Closer(multiAddr.Address(), client.Address(), query)
		if err != nil {
			return identity.MultiAddress{}, fmt.Errorf("cannot compare address distances %v and %v: %v", multiAddr.Address(), client.Address(), err)
		}
		if (isPeerCloser || isBootstrapping) && query != client.Address() {
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
		if err != nil && err != io.EOF {
			return identity.MultiAddress{}, fmt.Errorf("cannot send query to %v: %v", peer, err)
		}

		// Add the peer to the DHT after a successful query and ignore the
		// error
		for _, multiAddr := range multiAddrs {
			if multiAddr.Address() == query && !isBootstrapping {
				return multiAddr, nil
			}
			if _, ok := blacklist[multiAddr.Address()]; ok {
				continue
			}
			whitelist = append(whitelist, multiAddr)
		}
	}
	return identity.MultiAddress{}, ErrNotFound
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
	requestSignature, err := client.crypter.Sign(client.multiAddress)
	if err != nil {
		return false
	}
	request := &PingRequest{
		Signature:    requestSignature,
		MultiAddress: client.multiAddress.String(),
	}
	if _, err := swarmClient.Ping(ctx, request); err != nil {
		client.dht.RemoveMultiAddress(multiAddress)
		return true
	}
	client.dht.UpdateMultiAddress(multiAddress)
	return false
}
