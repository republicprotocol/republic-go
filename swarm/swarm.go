package swarm

import (
	"context"
	"log"
	"math"
	"math/rand"

	"github.com/republicprotocol/republic-go/identity"
)

type Client interface {

	// Ping a node to inform it of the existence of a new multiaddress
	Ping(ctx context.Context, to identity.MultiAddress, multiAddr identity.MultiAddress, nonce uint64) error

	// An implementation of Client's Pong should pass its own identity.MultiAddress
	// to the node during the pong.
	Pong(ctx context.Context, to identity.MultiAddress) error

	// Query a node for the identity.MultiAddress of an identity.Address.
	// Returns a list of identity.MultiAddresses that are closer to the query
	// than the node that was queried.
	Query(ctx context.Context, to identity.MultiAddress, query identity.Address, querySig [65]byte) (identity.MultiAddresses, error)

	// MultiAddress of the Client.
	MultiAddress() identity.MultiAddress
}

type Swarmer interface {

	// Ping will increment multiaddress nonce by 1 and send this information
	// to α randomly selected nodes. Ping must be called initially to connect
	// to the network. For this to work, there must be atleast one multiaddress
	// of a node in the network available in the storer.
	Ping(ctx context.Context) error

	// Broadcast will send the multiaddress and nonce to a randomly selected
	// α nodes in the network. Before this, a Pong request will be initiated
	// to inform the node at the multiaddress about the existence of this node.
	Broadcast(ctx context.Context, multiAddr identity.MultiAddress, nonce uint64) error

	// Query for the identity.MultiAddress of an identity.Address using a BFS
	// algorithm. The depth parameters limits the BFS, however a depth below
	// zero will perform an exhaustive search. Returns ErrMultiAddressNotFound
	// if no matching results are found.
	Query(ctx context.Context, query identity.Address) (identity.MultiAddress, error)

	// MultiAddress of the Swarmer.
	MultiAddress() identity.MultiAddress

	// GetConnectedPeers will return the multiaddresses of all the darknodes that
	// are connected to the swarmer.
	GetConnectedPeers() (identity.MultiAddresses, error)
}

type swarmer struct {
	client Client
	storer MultiAddressStorer
	α      int
	nonce  uint64
}

// NewSwarmer will return an object that implements the Swarmer interface.
func NewSwarmer(client Client, storer MultiAddressStorer, α int) (Swarmer, error) {
	swarmer := &swarmer{
		client: client,
		storer: storer,
		α:      α,
	}
	if err := swarmer.retrieveAndUpdateSelf(); err != nil {
		return nil, err
	}
	return swarmer, nil
}

// Ping will update the multiaddress and nonce in the storer and send
// the swarmer's multiaddress to α randomly selected nodes.
func (swarmer *swarmer) Ping(ctx context.Context) error {
	if err := swarmer.retrieveAndUpdateSelf(); err != nil {
		return err
	}

	return swarmer.pingNodes(ctx, swarmer.MultiAddress(), swarmer.nonce)
}

// Broadcast implements the Swarmer interface.
func (swarmer *swarmer) Broadcast(ctx context.Context, multiAddr identity.MultiAddress, nonce uint64) error {
	if err := swarmer.client.Pong(ctx, multiAddr); err != nil {
		return err
	}

	if err := swarmer.pingNodes(ctx, multiAddr, nonce); err != nil {
		return err
	}

	return nil
}

// Query implements the Swarmer interface.
func (swarmer *swarmer) Query(ctx context.Context, query identity.Address) (identity.MultiAddress, error) {
	return swarmer.query(ctx, query)
}

// MultiAddress implements the Swarmer interface.
func (swarmer *swarmer) MultiAddress() identity.MultiAddress {
	return swarmer.client.MultiAddress()
}

func (swarmer *swarmer) GetConnectedPeers() (identity.MultiAddresses, error) {
	multiaddressesIterator, err := swarmer.storer.MultiAddresses()
	if err != nil {
		return nil, err
	}
	multiAddrs, _, err := multiaddressesIterator.Collect()
	if err != nil {
		return nil, err
	}
	return multiAddrs, nil
}

func (swarmer *swarmer) query(ctx context.Context, query identity.Address) (identity.MultiAddress, error) {

	// Is the multiaddress present in the storer?
	multiAddr, _, err := swarmer.storer.MultiAddress(query)
	if err == nil {
		return multiAddr, nil
	}
	if err != ErrMultiAddressNotFound {
		return identity.MultiAddress{}, err
	}

	if swarmer.MultiAddress().Address() == query {
		return swarmer.MultiAddress(), nil
	}

	// If multiaddress is not present in the store, query for closer nodes.
	multiAddrsCloser, err := closerMultiAddrs(swarmer.storer, swarmer.MultiAddress().Address(), query, swarmer.α)
	if err != nil {
		return identity.MultiAddress{}, err
	}

	keys := map[identity.Address]struct{}{}

	// Query α closer multiaddrs until the node is reached or there are no
	// more newer multiaddresses in the store.
	for i := 0; len(multiAddrsCloser) > 0; i++ {
		multiAddr := multiAddrsCloser[0]
		multiAddrsCloser = multiAddrsCloser[1:]
		if _, ok := keys[multiAddr.Address()]; ok {
			continue
		}
		keys[multiAddr.Address()] = struct{}{}

		// Check if same address is returned
		if multiAddr.Address() == swarmer.MultiAddress().Address() {
			continue
		}

		if multiAddr.Address() == query {
			return multiAddr, nil
		}

		// Query for closer addresses.
		closer, err := swarmer.client.Query(ctx, multiAddr, query, [65]byte{})
		if err != nil {
			log.Printf("cannot query %v: %v", multiAddr.Address(), err)
			continue
		}

		for _, multi := range closer {
			_, err := swarmer.storer.PutMultiAddress(multi, 1)
			if err != nil {
				log.Printf("cannot store %v: %v", multi.Address(), err)
				continue
			}

			if multi.Address() == query {
				return multi, nil
			}

			if _, ok := keys[multi.Address()]; !ok {
				multiAddrsCloser = append(multiAddrsCloser, multi)
			}
		}
	}
	return identity.MultiAddress{}, ErrMultiAddressNotFound
}

// pingNodes will ping α random nodes in the storer using the client to gossip about the multiaddress and nonce seen.
func (swarmer *swarmer) pingNodes(ctx context.Context, multiAddr identity.MultiAddress, nonce uint64) error {
	multiAddrs, err := swarmer.GetConnectedPeers()
	if err != nil {
		return err
	}

	keys := map[int]struct{}{}
	for len(keys) < int(math.Min(float64(swarmer.α), float64(len(multiAddrs)))) {
		i := rand.Intn(len(multiAddrs))
		if _, ok := keys[i]; ok {
			continue
		}
		keys[i] = struct{}{}

		if err := swarmer.client.Ping(ctx, multiAddrs[i], multiAddr, nonce); err != nil {
			log.Printf("cannot ping node with address %v: %v", multiAddrs[i], err)
		}
	}

	return nil
}

func (swarmer *swarmer) retrieveAndUpdateSelf() error {
	_, nonce, err := swarmer.storer.MultiAddress(swarmer.MultiAddress().Address())
	if err != nil && err != ErrMultiAddressNotFound {
		return err
	}
	_, err = swarmer.storer.PutMultiAddress(swarmer.MultiAddress(), nonce)
	if err != nil {
		return err
	}
	swarmer.nonce = nonce
	return nil
}

type Server interface {

	// Ping will register the multiaddress and nonce into a storer and
	// broadcast this information to the network.
	Ping(ctx context.Context, from identity.MultiAddress, nonce uint64) error

	// Pong will handle responses from unseen nodes and register their
	// multiaddresses in the storer.
	Pong(ctx context.Context, from identity.MultiAddress, nonce uint64) error

	// Query will return a list of multiaddresses that are closest to the
	// query node. This list is limited by α.
	Query(ctx context.Context, query identity.Address, querySig [65]byte) (identity.MultiAddresses, error)
}

type server struct {
	swarmer        Swarmer
	multiAddrStore MultiAddressStorer
	α              int
}

func NewServer(swarmer Swarmer, multiAddrStore MultiAddressStorer, α int) Server {
	return &server{
		swarmer:        swarmer,
		multiAddrStore: multiAddrStore,
		α:              α,
	}
}

func (server *server) Ping(ctx context.Context, multiAddr identity.MultiAddress, nonce uint64) error {
	// FIXME: Verify multi address signature

	changed, err := server.multiAddrStore.PutMultiAddress(multiAddr, nonce)
	if err != nil {
		return err
	}
	if !changed {
		return nil
	}

	// If the multiaddress is new or has modifications, gossip the new
	// information to α random nodes in the network.
	return server.swarmer.Broadcast(ctx, multiAddr, nonce)
}

// Pong will store unseen multiaddresses in the storer.
func (server *server) Pong(ctx context.Context, from identity.MultiAddress, nonce uint64) error {
	_, err := server.multiAddrStore.PutMultiAddress(from, nonce)
	return err
}

func (server *server) Query(ctx context.Context, query identity.Address, querySig [65]byte) (identity.MultiAddresses, error) {
	return closerMultiAddrs(server.multiAddrStore, server.swarmer.MultiAddress().Address(), query, server.α)
}

func closerMultiAddrs(storer MultiAddressStorer, addr, query identity.Address, α int) (identity.MultiAddresses, error) {
	multi, _, err := storer.MultiAddress(query)
	if err == nil {
		return []identity.MultiAddress{multi}, nil
	}

	multiAddrsCloser := identity.MultiAddresses{}
	multiaddressesIterator, err := storer.MultiAddresses()
	if err != nil {
		return identity.MultiAddresses{}, err
	}
	multiAddrs, _, err := multiaddressesIterator.Collect()
	if err != nil {
		return identity.MultiAddresses{}, err
	}

	keys := map[int]struct{}{}
	for len(keys) < len(multiAddrs) && len(multiAddrsCloser) < α {
		i := rand.Intn(len(multiAddrs))
		if _, ok := keys[i]; ok {
			continue
		}
		keys[i] = struct{}{}

		if multiAddrs[i].Address() == query {
			multiAddrsCloser = append(multiAddrsCloser, multiAddrs[i])
			return multiAddrsCloser, nil
		}
		if multiAddrs[i].Address() == addr {
			continue
		}

		multiAddrsCloser = append(multiAddrsCloser, multiAddrs[i])
	}

	return multiAddrsCloser, nil
}
