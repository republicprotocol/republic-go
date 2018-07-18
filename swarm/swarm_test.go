package swarm_test

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/leveldb"
	. "github.com/republicprotocol/republic-go/swarm"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/testutils"
)

var _ = Describe("Swarm", func() {

	AfterEach(func() {
		os.RemoveAll("./tmp")
	})

	Context("when bootstrapping", func() {

		It("should be able to query any peer after bootstrapping", func() {
			numberOfClients := 100
			numberOfBootstrapClients := 5
			α := 10

			// Creating clients.
			stores := make([]MultiAddressStorer, numberOfClients)
			clients := make([]Client, numberOfClients)
			multiAddrs := make(identity.MultiAddresses, numberOfClients)
			swarmers := make([]Swarmer, numberOfClients)

			// Creating a common server hub for all clients to use.
			serverHub := &mockServerHub{
				connsMu: new(sync.Mutex),
				conns:   map[identity.Address]Server{},
				active:  map[identity.Address]bool{},
			}

			for i := 0; i < numberOfClients; i++ {
				client, err := newMockClientToServer(serverHub)
				Expect(err).ShouldNot(HaveOccurred())
				clients[i] = &client
				multiAddrs[i] = clients[i].MultiAddress()

				// Create leveldb store and store own multiaddress.
				db, err := leveldb.NewStore(fmt.Sprintf("./tmp/swarmer.%v.out", i+1), 72*time.Hour)
				Expect(err).ShouldNot(HaveOccurred())
				stores[i] = db.MultiAddressStore()
				stores[i].PutMultiAddress(multiAddrs[i], 1)

				// Creating swarmer for the client.
				swarmers[i] = NewSwarmer(clients[i], stores[i], α, 5*time.Second)
				server := NewServer(swarmers[i], stores[i], α)
				serverHub.Register(multiAddrs[i].Address(), server)
			}

			ctx, cancelCtx := context.WithCancel(context.Background())
			defer cancelCtx()

			// Bootstrapping created clients
			dispatch.CoForAll(numberOfClients, func(i int) {
				defer GinkgoRecover()

				for j := 0; j < numberOfBootstrapClients; j++ {
					stores[i].PutMultiAddress(multiAddrs[j], 1)
				}
				err := swarmers[i].Ping(ctx)
				if err != nil {
					Expect(err).ShouldNot(HaveOccurred())
				}
			})

			// Query for clients
			for i := 0; i < numberOfClients; i++ {
				for j := 0; j < numberOfClients; j++ {
					if i == j {
						continue
					}
					multiAddr, err := swarmers[i].Query(ctx, multiAddrs[j].Address())
					if err != nil && serverHub.IsRegistered(multiAddrs[j].Address()) {
						Expect(err).ShouldNot(HaveOccurred())
					}
					Expect(multiAddr.String()).To(Equal(multiAddrs[j].String()))
				}
			}
		})
	})
})

// mockServerHub will store all Servers that Clients use to Query and Ping
type mockServerHub struct {
	connsMu *sync.Mutex
	conns   map[identity.Address]Server
	active  map[identity.Address]bool
}

func (serverHub *mockServerHub) Register(serverAddr identity.Address, server Server) {
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

type mockClientToServer struct {
	multiAddr identity.MultiAddress
	serverHub *mockServerHub
}

func newMockClientToServer(mockServerHub *mockServerHub) (mockClientToServer, error) {
	multiAddr, err := testutils.RandomMultiAddress()
	if err != nil {
		return mockClientToServer{}, err
	}

	return mockClientToServer{
		multiAddr: multiAddr,
		serverHub: mockServerHub,
	}, nil
}

func (client *mockClientToServer) Ping(ctx context.Context, to identity.MultiAddress, multiAddr identity.MultiAddress, nonce uint64) error {
	var clientAddr Server
	isActive := false

	client.serverHub.connsMu.Lock()
	if isActive, _ = client.serverHub.active[to.Address()]; isActive {
		clientAddr = client.serverHub.conns[to.Address()]
	}
	client.serverHub.connsMu.Unlock()

	if isActive {
		randomSleep()
		return clientAddr.Ping(ctx, multiAddr, nonce)
	}
	return errors.New("address not active")
}

func (client *mockClientToServer) Pong(ctx context.Context, multiAddr identity.MultiAddress, nonce uint64) error {
	var clientAddr Server
	isActive := false

	client.serverHub.connsMu.Lock()
	if isActive, _ = client.serverHub.active[multiAddr.Address()]; isActive {
		clientAddr = client.serverHub.conns[multiAddr.Address()]
	}
	client.serverHub.connsMu.Unlock()

	if isActive {
		randomSleep()
		return clientAddr.Pong(ctx, client.multiAddr, nonce)
	}
	return errors.New("pong address not active")
}

func (client *mockClientToServer) Query(ctx context.Context, to identity.MultiAddress, query identity.Address, querySig [65]byte) (identity.MultiAddresses, error) {
	var clientAddr Server
	isActive := false

	client.serverHub.connsMu.Lock()
	if isActive, _ = client.serverHub.active[to.Address()]; isActive {
		clientAddr = client.serverHub.conns[to.Address()]
	}
	client.serverHub.connsMu.Unlock()

	if isActive {
		randomSleep()
		return clientAddr.Query(ctx, query, querySig)
	}
	return identity.MultiAddresses{}, errors.New("server not active")
}

func (client *mockClientToServer) MultiAddress() identity.MultiAddress {
	return client.multiAddr
}

func randomSleep() {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(120)
	time.Sleep(time.Duration(r) * time.Millisecond)
}
