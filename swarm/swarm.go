package swarm

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"sync"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
)

// ErrNonceTooLow is returned if the nonce of the multiAddress is lower than the
// one present in the local store.
var ErrNonceTooLow = errors.New("nonce too low")

// A Client exposes methods for invoking RPCs on a remote server.
type Client interface {

	// Ping a remote server to propagate a multi-address throughout the
	// network.
	Ping(ctx context.Context, to, multi identity.MultiAddress) error

	// Pong a remote server with own multi-address in response to a Ping.
	Pong(ctx context.Context, to identity.MultiAddress) error

	// Query a node for the identity.MultiAddress of an identity.Address.
	// Returns a list of random identity.MultiAddresses from the node that
	// was queried.
	Query(ctx context.Context, to identity.MultiAddress, query identity.Address) (identity.MultiAddresses, error)

	// MultiAddress used when invoking the Pong RPC.
	MultiAddress() identity.MultiAddress
}

type Swarmer interface {

	// Ping will increment nonce of the multi-address by 1 and send this information
	// to α randomly selected nodes. Ping must be called initially to connect
	// to the network. For this to work, there must be at least one multiAddress
	// of a node in the network available in the storer.
	Ping(ctx context.Context) error

	// Pong is used to inform a target node after it's multiaddress is received.
	Pong(ctx context.Context, to identity.MultiAddress) error

	// BroadcastMultiAddress to a maximum of α randomly selected nodes.
	BroadcastMultiAddress(ctx context.Context, multiAddr identity.MultiAddress) error

	// Query a node for the identity.MultiAddress of an identity.Address.
	// Returns a list of random identity.MultiAddresses from the node that
	// was queried.
	Query(ctx context.Context, query identity.Address) (identity.MultiAddress, error)

	// MultiAddress used when pinging and ponging.
	MultiAddress() identity.MultiAddress

	// Peers will return the latest version of all known multi-addresses. These
	// multi-addresses are not guaranteed to be connected.
	Peers() (identity.MultiAddresses, error)
}

type swarmer struct {
	client Client
	key    *crypto.EcdsaKey
	storer MultiAddressStorer
	α      int
}

// NewSwarmer will return an object that implements the Swarmer interface.
func NewSwarmer(client Client, storer MultiAddressStorer, α int, key *crypto.EcdsaKey) Swarmer {
	return &swarmer{
		client: client,
		key:    key,
		storer: storer,
		α:      α,
	}
}

// Ping will update the multi-address and nonce in the storer and send
// the swarmer's multi-address to α randomly selected nodes.
func (swarmer *swarmer) Ping(ctx context.Context) error {
	multi, err := swarmer.storer.MultiAddress(swarmer.MultiAddress().Address())
	if err != nil {
		return err
	}
	multi.Nonce++
	signature, err := swarmer.key.Sign(multi.Hash())
	if err != nil {
		return err
	}
	multi.Signature = signature
	if err := swarmer.storer.PutMultiAddress(multi); err != nil {
		return err
	}

	return swarmer.pingNodes(ctx, swarmer.MultiAddress())
}

func (swarmer *swarmer) Pong(ctx context.Context, to identity.MultiAddress) error {
	return swarmer.client.Pong(ctx, to)
}

// BroadcastMultiAddress implements the Swarmer interface.
func (swarmer *swarmer) BroadcastMultiAddress(ctx context.Context, multiAddr identity.MultiAddress) error {
	return swarmer.pingNodes(ctx, multiAddr)
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
	multiAddressesIterator, err := swarmer.storer.MultiAddresses()
	if err != nil {
		return nil, err
	}
	defer multiAddressesIterator.Release()

	return multiAddressesIterator.Collect()
}

func (swarmer *swarmer) query(ctx context.Context, query identity.Address) (identity.MultiAddress, error) {
	// Is the multi-address same as the swarmer's multiAddress?
	if swarmer.MultiAddress().Address() == query {
		return swarmer.MultiAddress(), nil
	}
	// Is the multi-address present in the storer?
	multiAddr, err := swarmer.storer.MultiAddress(query)
	if err == nil {
		return multiAddr, nil
	}
	if err != ErrMultiAddressNotFound {
		return identity.MultiAddress{}, err
	}

	// If multi-address is not present in the store, query for a maximum of α random nodes.
	randomMultiAddrs, err := randomMultiAddrs(swarmer.storer, swarmer.MultiAddress().Address(), swarmer.α)
	if err != nil {
		return identity.MultiAddress{}, err
	}

	// Create two maps to records the addrs we have seen and queried
	seenMu := new(sync.Mutex)
	seenAddrs := map[identity.Address]struct{}{
		swarmer.MultiAddress().Address(): {},
	}

	// Query α multiAddresses until the node is reached or there are no
	// more newer multi-addresses can be queried.
	for {
		if len(randomMultiAddrs) == 0 {
			break
		}
		if _, ok := seenAddrs[query]; ok {
			target, err := swarmer.storer.MultiAddress(query)
			if err != nil {
				logger.Error("cannot get multiAddress from the storer")
				return identity.MultiAddress{}, err
			}
			log.Println("found !!!!")
			return target, nil
		}

		// Pick at most α multiAddresses
		length := swarmer.α
		if len(randomMultiAddrs) < swarmer.α {
			length = len(randomMultiAddrs)
		}
		peersThisRound := randomMultiAddrs[:length]
		randomMultiAddrs = randomMultiAddrs[length:]

		// Query the α multiAddresses simultaneously
		dispatch.CoForAll(peersThisRound, func(i int) {
			multiAddr := peersThisRound[i]
			seenMu.Lock()
			_, ok := seenAddrs[multiAddr.Address()]
			seenAddrs[multiAddr.Address()] = struct{}{}
			seenMu.Unlock()
			if ok {
				return
			}
			multiAddrs, err := swarmer.client.Query(ctx, multiAddr, query)
			if err != nil {
				log.Printf("cannot query %v: %v", multiAddr.Address(), err)
				return
			}

			dispatch.CoForAll(multiAddrs, func(j int) {

				// Verify the multi address
				multi := multiAddrs[j]
				log.Println("querying", multi.Address().String())
				verifier := crypto.NewEcdsaVerifier(multi.Address().String())
				if err := verifier.Verify(multi.Hash(), multi.Signature); err != nil {
					return
				}

				// Mark the new multi as seen and add to the query backlog.
				seenMu.Lock()
				if _, ok := seenAddrs[multi.Address()]; ok {
					seenMu.Unlock()
					return
				}
				seenAddrs[multi.Address()] = struct{}{}
				randomMultiAddrs = append(randomMultiAddrs, multi)
				seenMu.Unlock()

				// Put the new multi in our storer if it has a higher nonce
				oldMulti, err := swarmer.storer.MultiAddress(multi.Address())
				if err != nil && err != ErrMultiAddressNotFound {
					log.Printf("cannot get nonce of %v : %v", multi.Address(), err)
					return
				}
				if err == ErrMultiAddressNotFound || oldMulti.Nonce < multi.Nonce {
					err = swarmer.storer.PutMultiAddress(multi)
					if err != nil {
						log.Printf("cannot store %v: %v", multi.Address(), err)
						return
					}
				}
			})
		})
	}

	return identity.MultiAddress{}, ErrMultiAddressNotFound
}

// pingNodes will ping α random nodes in the storer using the client to gossip about the multi-address and nonce seen.
func (swarmer *swarmer) pingNodes(ctx context.Context, multiAddr identity.MultiAddress) error {
	multiAddrs, err := swarmer.Peers()
	if err != nil {
		return err
	}

	if len(multiAddrs) <= swarmer.α {
		for _, multi := range multiAddrs {
			if multi.Address() == multiAddr.Address() {
				continue
			}

			if err := swarmer.client.Ping(ctx, multi, multiAddr); err != nil {
				log.Printf("cannot ping node with address %v: %v", multi, err)
				continue
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

		if err := swarmer.client.Ping(ctx, multiAddrs[i], multiAddr); err != nil {
			log.Printf("cannot ping node with address %v: %v", multiAddrs[i], err)
		}
	}

	return nil
}

type Server interface {

	// Ping will register the multi-address and nonce into a storer and
	// broadcast this information to the network.
	Ping(ctx context.Context, from identity.MultiAddress) error

	// Pong will handle responses from unseen nodes and register their
	// multi-addresses in the storer.
	Pong(ctx context.Context, from identity.MultiAddress) error

	// Query will return the multi-address of the query, if available in
	// the storer. Otherwise, it will return α random multi-addresses from
	// the storer.
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

func (server *server) Ping(ctx context.Context, multiAddr identity.MultiAddress) error {
	// Verifies the signature
	verifier := crypto.NewEcdsaVerifier(multiAddr.Address().String())
	if err := verifier.Verify(multiAddr.Hash(), multiAddr.Signature); err != nil {
		return err
	}

	// Pong back
	if err := server.swarmer.Pong(ctx, multiAddr); err != nil {
		return err
	}
	// Compare the nonce and see if we need to gossip the ping.
	oldMulti, err := server.multiAddrStore.MultiAddress(multiAddr.Address())
	if err != nil {
		return err
	}
	if oldMulti.Nonce < multiAddr.Nonce {
		err := server.multiAddrStore.PutMultiAddress(multiAddr)
		if err != nil {
			return err
		}
		return server.swarmer.BroadcastMultiAddress(ctx, multiAddr)
	}

	return nil
}

// Pong will store unseen multi-addresses in the storer.
func (server *server) Pong(ctx context.Context, from identity.MultiAddress) error {
	// Verifies the signature
	verifier := crypto.NewEcdsaVerifier(from.Address().String())
	if err := verifier.Verify(from.Hash(), from.Signature); err != nil {
		return err
	}

	// Compare the nonce and see if we need to gossip the ping.
	oldMulti, err := server.multiAddrStore.MultiAddress(from.Address())
	if err != nil {
		return err
	}
	if oldMulti.Nonce < from.Nonce {
		err := server.multiAddrStore.PutMultiAddress(from)
		if err != nil {
			return err
		}
		return server.swarmer.BroadcastMultiAddress(ctx, from)
	}

	return nil
}

func (server *server) Query(ctx context.Context, query identity.Address) (identity.MultiAddresses, error) {
	multiAddr, err := server.multiAddrStore.MultiAddress(query)
	if err == nil {
		return []identity.MultiAddress{multiAddr}, nil
	}
	return randomMultiAddrs(server.multiAddrStore, server.swarmer.MultiAddress().Address(), server.α)
}

// randomMultiAddrs returns maximum α random multi-addresses from the storer.
func randomMultiAddrs(storer MultiAddressStorer, self identity.Address, α int) (identity.MultiAddresses, error) {
	// Get all known multi-addresses from the storer.
	multiAddrsIter, err := storer.MultiAddresses()
	if err != nil {
		log.Printf("error getting multiaddresses: %v", err)
		return identity.MultiAddresses{}, err
	}
	defer multiAddrsIter.Release()

	multiAddrs, err := multiAddrsIter.Collect()
	if err != nil {
		log.Printf("error collecting multiaddresses: %v", err)
		return identity.MultiAddresses{}, err
	}

	if len(multiAddrs) <= α {
		return multiAddrs, nil
	}

	// Randomly select α multi-addresses.
	results := identity.MultiAddresses{}
	for len(results) < α {
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
