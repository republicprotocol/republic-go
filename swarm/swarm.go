package swarm

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/republicprotocol/republic-go/identity"
)

// A Client exposes methods for invoking RPCs on a remote server.
type Client interface {

	// Ping a remote server to propagate a multi-address throughout the
	// network.
	Ping(ctx context.Context, to identity.MultiAddress, multiAddr identity.MultiAddress, nonce uint64) error

	// Pong a remote server with own multi-address in response to a Ping.
	Pong(ctx context.Context, to identity.MultiAddress) error

	Query(ctx context.Context, to identity.MultiAddress, query identity.Address) (identity.MultiAddresses, error)

	// MultiAddress used when invoking the Pong RPC.
	MultiAddress() identity.MultiAddress
}

type Swarmer interface {

	// Ping will increment multiaddress nonce by 1 and send this information
	// to α randomly selected nodes. Ping must be called initially to connect
	// to the network. For this to work, there must be atleast one multiaddress
	// of a node in the network available in the storer.
	Ping(ctx context.Context) error

	Pong(ctx context.Context, to identity.MultiAddress) error

	BroadcastMultiAddress(ctx context.Context, multiAddr identity.MultiAddress, nonce uint64) error

	Query(ctx context.Context, query identity.Address) (identity.MultiAddress, error)

	// MultiAddress used when pinging and ponging.
	MultiAddress() identity.MultiAddress

	// Peers will return the latest version of all known multi-addresses. These
	// multi-addresses are not guaranteed to be connected.
	Peers() (identity.MultiAddresses, error)
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

func (swarmer *swarmer) Pong(ctx context.Context, to identity.MultiAddress) error {
	return swarmer.client.Pong(ctx, to)
}

// BroadcastMultiAddress implements the Swarmer interface.
func (swarmer *swarmer) BroadcastMultiAddress(ctx context.Context, multiAddr identity.MultiAddress, nonce uint64) error {
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

func (swarmer *swarmer) Peers() (identity.MultiAddresses, error) {
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
	randomMultiAddrs, err := randomMultiAddrs(swarmer.storer, swarmer.MultiAddress().Address(), query, swarmer.α)
	if err != nil {
		return identity.MultiAddress{}, err
	}

	multiAddrsIt, err := swarmer.storer.MultiAddresses()
	if err != nil {
		log.Println(err)
		// return identity.MultiAddress{}, err
	}
	mas, _, err := multiAddrsIt.Collect()
	if err != nil {
		log.Println(err)
		// return identity.MultiAddress{}, err
	}

	log.Printf("got %v random multiaddrs for a total of %v in store", len(randomMultiAddrs), len(mas))

	seenAddrs := map[identity.Address]struct{}{}

	// Query α closer multiaddrs until the node is reached or there are no
	// more newer multiaddresses in the store.
	for len(randomMultiAddrs) > 0 {
		multiAddr := randomMultiAddrs[0]
		randomMultiAddrs = randomMultiAddrs[1:]
		if _, ok := seenAddrs[multiAddr.Address()]; ok {
			continue
		}
		seenAddrs[multiAddr.Address()] = struct{}{}

		// Check if same address is returned
		if multiAddr.Address() == swarmer.MultiAddress().Address() {
			continue
		}

		if multiAddr.Address() == query {
			return multiAddr, nil
		}

		// Query for more random multiaddresses addresses.
		multiAddrs, err := swarmer.client.Query(ctx, multiAddr, query)
		if err != nil {
			log.Printf("cannot query %v: %v", multiAddr.Address(), err)
			continue
		}

		for _, multi := range multiAddrs {
			if multi.Address() == query {
				return multi, nil
			}

			if _, ok := seenAddrs[multi.Address()]; ok {
				continue
			}

			// Get nonce of multiaddress if present in store and
			// store it back to the store.
			_, nonce, err := swarmer.storer.MultiAddress(multi.Address())
			if err != nil && err != ErrMultiAddressNotFound {
				continue
			}

			_, err = swarmer.storer.PutMultiAddress(multi, nonce)
			if err != nil {
				log.Printf("cannot store %v: %v", multi.Address(), err)
				continue
			}

			if _, ok := seenAddrs[multi.Address()]; !ok {
				randomMultiAddrs = append(randomMultiAddrs, multi)
			}
		}
		log.Printf("new %v random multiaddrs", len(randomMultiAddrs))
	}

	for key := range seenAddrs {
		log.Printf("saw %v", key)
	}
	return identity.MultiAddress{}, ErrMultiAddressNotFound
}

// pingNodes will ping α random nodes in the storer using the client to gossip about the multiaddress and nonce seen.
func (swarmer *swarmer) pingNodes(ctx context.Context, multiAddr identity.MultiAddress, nonce uint64) error {
	multiAddrs, err := swarmer.Peers()
	if err != nil {
		return err
	}

	if len(multiAddrs) <= swarmer.α {
		for _, multi := range multiAddrs {
			if multi.Address() == multiAddr.Address() {
				continue
			}

			log.Printf("multiaddress: %v", multiAddrs[i])
			if err := swarmer.client.Ping(ctx, multi, multiAddr, nonce); err != nil {
				log.Printf("cannot ping node with address %v: %v", multi, err)
			}
		}
		return nil
	}

	seenAddrs := map[int]struct{}{}
	for len(seenAddrs) < swarmer.α {
		i := rand.Intn(len(multiAddrs))
		if _, ok := seenAddrs[i]; ok {
			continue
		}
		seenAddrs[i] = struct{}{}

		if multiAddrs[i].Address() == multiAddr.Address() {
			continue
		}

		log.Printf("multiaddress: %v", multiAddrs[i])
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

	Query(ctx context.Context, query identity.Address) (identity.MultiAddresses, error)
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
	if err := server.swarmer.Pong(ctx, multiAddr); err != nil {
		return err
	}
	if !changed {
		return nil
	}

	// If the multiaddress is new or has modifications, gossip the new
	// information to α random nodes in the network.
	return server.swarmer.BroadcastMultiAddress(ctx, multiAddr, nonce)
}

// Pong will store unseen multiaddresses in the storer.
func (server *server) Pong(ctx context.Context, from identity.MultiAddress, nonce uint64) error {
	_, err := server.multiAddrStore.PutMultiAddress(from, nonce)
	return err
}

func (server *server) Query(ctx context.Context, query identity.Address) (identity.MultiAddresses, error) {
	return randomMultiAddrs(server.multiAddrStore, server.swarmer.MultiAddress().Address(), query, server.α)
}

func randomMultiAddrs(storer MultiAddressStorer, self, query identity.Address, α int) (identity.MultiAddresses, error) {
	log.Printf("Alpha: %v", α)
	multiAddr, _, err := storer.MultiAddress(query)
	if err == nil {
		log.Printf("got multiaddress: %v", multiAddr.Address())
		return []identity.MultiAddress{multiAddr}, nil
	}

	multiAddrsIter, err := storer.MultiAddresses()
	if err != nil {
		log.Printf("error at getting multiaddresses: %v", err)
		return identity.MultiAddresses{}, err
	}
	multiAddrs, _, err := multiAddrsIter.Collect()
	if err != nil {
		log.Printf("error at collecting multiaddresses: %v", err)
		return identity.MultiAddresses{}, err
	}
	if len(multiAddrs) <= α {
		log.Println("here")
		for _, m := range multiAddrs {
			log.Printf("got %v", m.Address())
		}
		return multiAddrs, nil
	}

	rand.Seed(time.Now().UnixNano())

	results := identity.MultiAddresses{}
	for len(results) < α {
		// Randomly select a multi-address and make sure it is not selected
		// more than once
		i := rand.Intn(len(multiAddrs))
		multiAddr := multiAddrs[i]

		multiAddrs[i] = multiAddrs[len(multiAddrs)-1]
		multiAddrs = multiAddrs[:len(multiAddrs)-1]

		// Do not return own multi-address
		if multiAddr.Address() == self {
			continue
		}
		results = append(results, multiAddr)
	}

	return results, nil
}
