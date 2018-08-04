package testutils

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/leveldb"
	"github.com/republicprotocol/republic-go/registry"
	"github.com/republicprotocol/republic-go/swarm"
)

// Swarmer is a mock implementation of the swarm.Swarmer interface.
type Swarmer struct {
	multiAddrsMu *sync.Mutex
	multiAddrs   map[identity.Address]identity.MultiAddress
}

func NewMockSwarmer() Swarmer {
	return Swarmer{
		multiAddrsMu: new(sync.Mutex),
		multiAddrs:   make(map[identity.Address]identity.MultiAddress),
	}
}

func (swarmer *Swarmer) Ping(ctx context.Context) error {
	return nil
}

func (swarmer *Swarmer) BroadcastMultiAddress(ctx context.Context, multiAddr identity.MultiAddress) error {
	return nil
}

func (swarmer *Swarmer) Query(ctx context.Context, query identity.Address) (identity.MultiAddress, error) {
	return identity.MultiAddress{}, nil
}

func (swarmer *Swarmer) MultiAddress() identity.MultiAddress {
	return identity.MultiAddress{}
}

func (swarmer *Swarmer) Peers() (identity.MultiAddresses, error) {
	return make([]identity.MultiAddress, len(swarmer.multiAddrs)), nil
}

func (swarmer *Swarmer) Pong(ctx context.Context, to identity.MultiAddress) error {
	return nil
}

func (swarmer *Swarmer) PutMultiAddress(multiAddr identity.MultiAddress) {
	swarmer.multiAddrsMu.Lock()
	defer swarmer.multiAddrsMu.Unlock()
	swarmer.multiAddrs[multiAddr.Address()] = multiAddr
}

func (swarmer *Swarmer) RemoveMultiAddress(multiAddr identity.MultiAddress) {
	swarmer.multiAddrsMu.Lock()
	defer swarmer.multiAddrsMu.Unlock()
	if _, ok := swarmer.multiAddrs[multiAddr.Address()]; ok {
		delete(swarmer.multiAddrs, multiAddr.Address())
	}
}

// MockServerHub will store all Servers that Clients use to Query and Ping
type MockServerHub struct {
	ConnsMu *sync.Mutex
	Conns   map[identity.Address]swarm.Server
	Active  map[identity.Address]bool
}

func (serverHub *MockServerHub) Register(serverAddr identity.Address, server swarm.Server) {
	serverHub.ConnsMu.Lock()
	defer serverHub.ConnsMu.Unlock()

	serverHub.Conns[serverAddr] = server
	serverHub.Active[serverAddr] = true
}

func (serverHub *MockServerHub) Deregister(serverAddr identity.Address) bool {
	serverHub.ConnsMu.Lock()
	defer serverHub.ConnsMu.Unlock()

	isActive, _ := serverHub.Active[serverAddr]
	if isActive {
		serverHub.Active[serverAddr] = false
	}
	return isActive
}

func (serverHub *MockServerHub) IsRegistered(serverAddr identity.Address) bool {
	serverHub.ConnsMu.Lock()
	defer serverHub.ConnsMu.Unlock()

	isActive, _ := serverHub.Active[serverAddr]

	return isActive
}

// ClientType defines different clients in terms of behaviour in
// the gossip network.
type ClientType int

const (
	Honest             = ClientType(iota)
	AlwaysDrop         // AlwaysDrop will drop multi-addresses it receives
	BroadcastStale     // Adversary2 will broadcast stale multi-addresses of other nodes
	BroadcastIncorrect // Adversary3 will broadcast incorrect multi-addresses
)

// ClientTypes contains an array of all possible ClientType values.
var ClientTypes = []ClientType{
	Honest,
	AlwaysDrop,
	BroadcastStale,
	BroadcastIncorrect,
}

type MockSwarmClient struct {
	addr       identity.Address
	store      swarm.MultiAddressStorer
	serverHub  *MockServerHub
	clientType ClientType
}

func NewMockSwarmClient(MockServerHub *MockServerHub, clientType ClientType, verifier *registry.Crypter) (MockSwarmClient, swarm.MultiAddressStorer, error) {
	multiAddr, err := RandomMultiAddress()
	if err != nil {
		return MockSwarmClient{}, nil, err
	}
	multiAddr.Nonce = 1
	signature, err := verifier.Sign(multiAddr.Hash())
	if err != nil {
		return MockSwarmClient{}, nil, err
	}
	multiAddr.Signature = signature

	// Create leveldb store and store own multiAddress.
	db, err := leveldb.NewStore(fmt.Sprintf("./tmp/swarmer-%v.out", multiAddr.Address()), 72*time.Hour)
	if err != nil {
		return MockSwarmClient{}, nil, err
	}
	store := db.SwarmMultiAddressStore()
	if err = store.PutMultiAddress(multiAddr); err != nil {
		return MockSwarmClient{}, nil, err
	}

	return MockSwarmClient{
		addr:       identity.Address(multiAddr.Address()),
		store:      store,
		serverHub:  MockServerHub,
		clientType: clientType,
	}, store, nil
}

func (client *MockSwarmClient) Ping(ctx context.Context, to identity.MultiAddress, multiAddr identity.MultiAddress) error {
	client.serverHub.ConnsMu.Lock()
	if isActive, _ := client.serverHub.Active[to.Address()]; !isActive {
		client.serverHub.ConnsMu.Unlock()
		return errors.New("server is not active")
	}
	server := client.serverHub.Conns[to.Address()]
	client.serverHub.ConnsMu.Unlock()

	randomSleep()

	switch client.clientType {
	case Honest, AlwaysDrop:
		return server.Ping(ctx, multiAddr)
	case BroadcastStale:
		multi, err := client.store.MultiAddress(multiAddr.Address())
		if err != nil {
			return err
		}
		return server.Ping(ctx, multi)
	case BroadcastIncorrect:
		multiAddr.Signature = []byte{}
		return server.Ping(ctx, multiAddr)
	default:
		return nil
	}
}

func (client *MockSwarmClient) Pong(ctx context.Context, multiAddr identity.MultiAddress) error {
	client.serverHub.ConnsMu.Lock()
	if isActive, _ := client.serverHub.Active[multiAddr.Address()]; !isActive {
		client.serverHub.ConnsMu.Unlock()
		return errors.New("server is not active")
	}
	server := client.serverHub.Conns[multiAddr.Address()]
	client.serverHub.ConnsMu.Unlock()

	multi, err := client.store.MultiAddress(client.addr)
	if err != nil {
		return err
	}

	randomSleep()

	switch client.clientType {
	case Honest, BroadcastStale:
		return server.Pong(ctx, multi)
	case BroadcastIncorrect:
		multi.Signature = []byte{}
		return server.Pong(ctx, multi)
	default:
		return nil
	}
}

func (client *MockSwarmClient) Query(ctx context.Context, to identity.MultiAddress, query identity.Address) (identity.MultiAddresses, error) {
	client.serverHub.ConnsMu.Lock()
	if isActive, _ := client.serverHub.Active[to.Address()]; !isActive {
		client.serverHub.ConnsMu.Unlock()
		return identity.MultiAddresses{}, errors.New("server is not active")
	}
	server := client.serverHub.Conns[to.Address()]
	client.serverHub.ConnsMu.Unlock()

	randomSleep()

	switch client.clientType {
	case Honest, BroadcastStale, BroadcastIncorrect:
		return server.Query(ctx, query)
	default:
		return identity.MultiAddresses{}, nil
	}
}

func (client *MockSwarmClient) MultiAddress() identity.MultiAddress {
	multi, err := client.store.MultiAddress(client.addr)
	if err != nil {
		log.Println("error retrieving multiAddress from store", err)
		return identity.MultiAddress{}
	}
	return multi
}

func randomSleep() {
	r := rand.Intn(120)
	time.Sleep(time.Duration(r) * time.Millisecond)
}

type MockSwarmBinder struct {
}

// NewMockSwarmBinder returns a MockSwarmBinder
func NewMockSwarmBinder() *MockSwarmBinder {
	return &MockSwarmBinder{}
}

func (binder *MockSwarmBinder) IsRegistered(darknodeAddr identity.Address) (bool, error) {
	return true, nil
}

func (binder *MockSwarmBinder) PublicKey(addr identity.Address) (rsa.PublicKey, error) {
	return rsa.PublicKey{}, nil
}
