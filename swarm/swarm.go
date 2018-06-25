package swarm

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/dht"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
)

// ErrMultiAddressNotFound is returned from a query when no
// identity.MultiAddress can be found for the identity.Address.
var ErrMultiAddressNotFound = errors.New("multiaddress not found")

type Server interface {
	Ping(ctx context.Context, from identity.MultiAddress) (identity.MultiAddress, error)
	Query(ctx context.Context, query identity.Address, querySig [65]byte) (identity.MultiAddresses, error)
}

type server struct {
	verifier   crypto.Verifier
	dhtManager dhtManager
}

func NewServer(verifier crypto.Verifier, client Client, dht *dht.DHT) Server {
	return &server{
		verifier: verifier,
		dhtManager: dhtManager{
			client: client,
			dht:    dht,
		},
	}
}

func (server *server) Ping(ctx context.Context, from identity.MultiAddress) (identity.MultiAddress, error) {
	if err := server.verifier.Verify(from.Hash(), from.Signature); err != nil {
		return server.dhtManager.client.MultiAddress(), nil
	}
	return server.dhtManager.client.MultiAddress(), server.dhtManager.updateDHT(from)
}

func (server *server) Query(ctx context.Context, query identity.Address, querySig [65]byte) (identity.MultiAddresses, error) {
	addr := server.dhtManager.client.MultiAddress().Address()
	multiAddrs := server.dhtManager.dht.MultiAddresses()
	multiAddrsCloser := make(identity.MultiAddresses, 0, len(multiAddrs)/2)
	for _, multiAddr := range multiAddrs {
		isPeerCloser, err := identity.Closer(multiAddr.Address(), addr, query)
		if err != nil {
			return multiAddrs, err
		}
		if isPeerCloser {
			multiAddrsCloser = append(multiAddrsCloser, multiAddr)
		}
	}
	return multiAddrsCloser, nil
}

type Client interface {

	// Ping a node. Returns the identity.MultiAddress of the node. An
	// implementation of Client should pass its own identity.MultiAddress to
	// the node during the ping.
	Ping(ctx context.Context, to identity.MultiAddress) (identity.MultiAddress, error)

	// Query a node for the identity.MultiAddress of an identity.Address.
	// Returns a list of identity.MultiAddresses that are closer to the query
	// than the node that was queried.
	Query(ctx context.Context, to identity.MultiAddress, query identity.Address, querySig [65]byte) (identity.MultiAddresses, error)

	// MultiAddress of the Client.
	MultiAddress() identity.MultiAddress
}

type Swarmer interface {

	// Bootstrap into the network. Starting from a list of known
	// identity.MultiAddresses, the Swarmer will query for itself throughout
	// the network. Doing so will connect the Swarmer to nodes in the network
	// that have identity.Addresses close to its own.
	Bootstrap(ctx context.Context, multiAddrs identity.MultiAddresses) error

	// Query for the identity.MultiAddress of an identity.Address using a BFS
	// algorithm. The depth parameters limits the BFS, however a depth below
	// zero will perform an exhaustive search. Returns ErrMultiAddressNotFound
	// if no matching results are found.
	Query(ctx context.Context, query identity.Address, depth int) (identity.MultiAddress, error)

	// MultiAddress of the Swarmer.
	MultiAddress() identity.MultiAddress
}

type swarmer struct {
	client     Client
	dhtManager dhtManager
}

func NewSwarmer(client Client, dht *dht.DHT) Swarmer {
	return &swarmer{
		client: client,
		dhtManager: dhtManager{
			client,
			dht,
		},
	}
}

// Bootstrap implements the Swarmer interface.
func (swarmer *swarmer) Bootstrap(ctx context.Context, multiAddrs identity.MultiAddresses) error {
	errs := make(chan error, len(multiAddrs)+1)

	go func() {
		defer close(errs)
		dispatch.CoForAll(multiAddrs, func(i int) {
			if multiAddrs[i].Address() == swarmer.client.MultiAddress().Address() {
				return
			}
			multiAddr, err := swarmer.client.Ping(ctx, multiAddrs[i])
			if err != nil {
				errs <- fmt.Errorf("cannot ping bootstrap node %v: %v", multiAddrs[i], err)
				return
			}
			if err := swarmer.dhtManager.updateDHT(multiAddr); err != nil {
				errs <- fmt.Errorf("cannot update dht with bootstrap node %v: %v", multiAddrs[i], err)
				return
			}
		})
		if _, err := swarmer.query(ctx, swarmer.client.MultiAddress().Address(), -1, true); err != nil {
			errs <- fmt.Errorf("error while bootstrapping: %v", err)
			return
		}
		logger.Network(logger.LevelInfo, fmt.Sprintf("connected to %v peers after bootstrapping", len(swarmer.dhtManager.dht.MultiAddresses())))
	}()

	for err := range errs {
		return err
	}
	return nil
}

// Query implements the Swarmer interface.
func (swarmer *swarmer) Query(ctx context.Context, query identity.Address, depth int) (identity.MultiAddress, error) {
	return swarmer.query(ctx, query, depth, false)
}

// MultiAddress implements the Swarmer interface.
func (swarmer *swarmer) MultiAddress() identity.MultiAddress {
	return swarmer.client.MultiAddress()
}

func (swarmer *swarmer) query(ctx context.Context, query identity.Address, depth int, isBootstrapping bool) (identity.MultiAddress, error) {
	whitelist := identity.MultiAddresses{}
	blacklist := map[identity.Address]struct{}{}
	blacklist[swarmer.client.MultiAddress().Address()] = struct{}{}

	// Build a list of identity.MultiAddresses that are closer to the query
	// than the Swarm service
	multiAddrs := swarmer.dhtManager.dht.MultiAddresses()
	for _, multiAddr := range multiAddrs {
		if isBootstrapping {
			if query != multiAddr.Address() {
				whitelist = append(whitelist, multiAddr)
			}
			continue
		}

		// Short circuit if the Swarm service is directly connected to the
		// query
		if query == multiAddr.Address() {
			return multiAddr, nil
		}
		isPeerCloser, err := identity.Closer(multiAddr.Address(), swarmer.client.MultiAddress().Address(), query)
		if err != nil {
			return identity.MultiAddress{}, fmt.Errorf("cannot compare address distances %v and %v: %v", multiAddr.Address(), swarmer.client.MultiAddress().Address(), err)
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

		if isBootstrapping {
			if _, err := swarmer.client.Ping(ctx, peer); err != nil {
				continue
			}
		}

		// Query for identity.MultiAddresses that are closer to the query
		// target than the peer itself, and add them to the whitelist
		multiAddrs, err := swarmer.client.Query(ctx, peer, query, [65]byte{})
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
			if err := swarmer.dhtManager.updateDHT(multiAddr); err != nil {
				logger.Network(logger.LevelInfo, fmt.Sprintf("cannot update dht with %v: %v", multiAddr, err))
			}
			whitelist = append(whitelist, multiAddr)
		}
	}
	if isBootstrapping {
		return identity.MultiAddress{}, nil
	}
	return identity.MultiAddress{}, ErrMultiAddressNotFound
}

type dhtManager struct {
	client Client
	dht    *dht.DHT
}

func (dhtManager *dhtManager) updateDHT(multiAddr identity.MultiAddress) error {
	if dhtManager.client.MultiAddress().Address() == multiAddr.Address() {
		return nil
	}
	if err := dhtManager.dht.UpdateMultiAddress(multiAddr); err != nil {
		if err == dht.ErrFullBucket {
			if dhtManager.pruneDHT(multiAddr.Address()) {
				return dhtManager.dht.UpdateMultiAddress(multiAddr)
			}
		}
		return err
	}
	return nil
}

func (dhtManager *dhtManager) pruneDHT(addr identity.Address) bool {
	bucket, err := dhtManager.dht.FindBucket(addr)
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
	multiAddr := bucket.MultiAddresses[0]
	multiAddrUpdated, err := dhtManager.client.Ping(ctx, multiAddr)
	if err != nil {
		dhtManager.dht.RemoveMultiAddress(multiAddr)
		return true
	}
	dhtManager.dht.UpdateMultiAddress(multiAddrUpdated)
	return false
}
