package testutils

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/leveldb"
	"github.com/republicprotocol/republic-go/swarm"
)

// mockServerHub will store all Servers that Clients use to Query and Ping
type mockServerHub struct {
	connsMu *sync.Mutex
	conns   map[identity.Address]swarm.Server
	active  map[identity.Address]bool
}

func (serverHub *mockServerHub) Register(serverAddr identity.Address, server swarm.Server) {
	serverHub.connsMu.Lock()
	defer serverHub.connsMu.Unlock()

	serverHub.conns[serverAddr] = server
	serverHub.active[serverAddr] = true
}

func (serverHub *mockServerHub) Deregister(serverAddr identity.Address) bool {
	serverHub.connsMu.Lock()
	defer serverHub.connsMu.Unlock()

	isActive, _ := serverHub.active[serverAddr]
	if isActive {
		serverHub.active[serverAddr] = false
	}
	return isActive
}

func (serverHub *mockServerHub) IsRegistered(serverAddr identity.Address) bool {
	serverHub.connsMu.Lock()
	defer serverHub.connsMu.Unlock()

	isActive, _ := serverHub.active[serverAddr]

	return isActive
}

// ClientType defines different clients in terms of behaviour in
// the gossip network
type ClientType int

const (
	Honest     = ClientType(iota)
	Adversary1 // Adversary1 will drop multi-addresses it receives
	Adversary2 // Adversary2 will broadcast stale multi-addresses of other nodes
	Adversary3 // Adversary3 will broadcast incorrect multi-addresses
)

type mockSwarmClient struct {
	addr       identity.Address
	store      swarm.MultiAddressStorer
	serverHub  *mockServerHub
	clientType ClientType
}

func newMockSwarmClient(mockServerHub *mockServerHub, key *crypto.EcdsaKey, clientType ClientType) (mockSwarmClient, error) {
	multiAddr, err := identity.Address(key.Address()).MultiAddress()
	if err != nil {
		return mockSwarmClient{}, err
	}
	multiAddr.Nonce = 1
	signature, err := key.Sign(multiAddr.Hash())
	if err != nil {
		return mockSwarmClient{}, err
	}
	multiAddr.Signature = signature

	// Create leveldb store and store own multiAddress.
	db, err := leveldb.NewStore(fmt.Sprintf("./tmp/swarmer-%v.out", key.Address()), 72*time.Hour)
	if err != nil {
		return mockSwarmClient{}, err
	}
	store := db.SwarmMultiAddressStore()
	_, err = store.PutMultiAddress(multiAddr)
	if err != nil {
		return mockSwarmClient{}, err
	}

	return mockSwarmClient{
		addr:       identity.Address(key.Address()),
		store:      store,
		serverHub:  mockServerHub,
		clientType: clientType,
	}, nil
}

func (client *mockSwarmClient) Ping(ctx context.Context, to identity.MultiAddress, multiAddr identity.MultiAddress) error {
	server, ok := client.serverHub.conns[to.Address()]
	if !ok {
		return errors.New("address not active")
	}
	randomSleep()

	// TODO :
	switch client.clientType {
	case Honest:
		return server.Ping(ctx, multiAddr)
	case Adversary1:

	case Adversary2:

	case Adversary3:

	}

	return nil
}

func (client *mockSwarmClient) Pong(ctx context.Context, multiAddr identity.MultiAddress) error {
	var server swarm.Server
	isActive := false

	client.serverHub.connsMu.Lock()
	if isActive, _ = client.serverHub.active[multiAddr.Address()]; isActive {
		server = client.serverHub.conns[multiAddr.Address()]
	}
	client.serverHub.connsMu.Unlock()

	multi, err := client.store.MultiAddress(client.addr)
	if err != nil {
		return err
	}

	if isActive {
		randomSleep()
		return server.Pong(ctx, multi)
	}
	return errors.New("pong address not active")
}

func (client *mockSwarmClient) Query(ctx context.Context, to identity.MultiAddress, query identity.Address) (identity.MultiAddresses, error) {
	var server swarm.Server
	isActive := false

	// if client.isAdversary && rand.Uint64()%2 == 0 {
	// 	return identity.MultiAddresses{}, nil
	// }

	client.serverHub.connsMu.Lock()
	if isActive, _ = client.serverHub.active[to.Address()]; isActive {
		server = client.serverHub.conns[to.Address()]
	}
	client.serverHub.connsMu.Unlock()

	if isActive {
		randomSleep()
		return server.Query(ctx, query)
	}
	return identity.MultiAddresses{}, errors.New("server not active")
}

func (client *mockSwarmClient) MultiAddress() identity.MultiAddress {
	multi, err := client.store.MultiAddress(client.addr)
	if err != nil {
		log.Println("err in getting the multiaddress in store")
	}

	return multi
}

func randomSleep() {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(120)
	time.Sleep(time.Duration(r) * time.Millisecond)
}
