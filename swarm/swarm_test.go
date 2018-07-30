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
	"github.com/republicprotocol/republic-go/crypto"
	. "github.com/republicprotocol/republic-go/swarm"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/leveldb"
	"github.com/republicprotocol/republic-go/testutils"
)

var _ = Describe("Swarm", func() {

	var (
		numberOfClients          = 100
		numberOfBootstrapClients = 5
		α                        = 4
	)

	Context("when querying after bootstrapping", func() {

		AfterEach(func() {
			os.RemoveAll("./tmp")
		})

		It("should be able to find most peers in the network", func() {
			// Creating clients.
			stores := make([]MultiAddressStorer, numberOfClients)
			clients := make([]Client, numberOfClients)
			multiAddresses := make(identity.MultiAddresses, numberOfClients)
			swarmers := make([]Swarmer, numberOfClients)

			// Creating a common server hub for all clients to use.
			serverHub := &mockServerHub{
				connsMu: new(sync.Mutex),
				conns:   map[identity.Address]Server{},
				active:  map[identity.Address]bool{},
			}

			for i := 0; i < numberOfClients; i++ {
				client, store, err := newMockClientToServer(serverHub, i, i < numberOfClients/4)
				Expect(err).ShouldNot(HaveOccurred())
				clients[i] = &client
				multiAddresses[i] = clients[i].MultiAddress()
				stores[i] = store

				ecdsaKey, err := crypto.RandomEcdsaKey()
				// Creating swarmer for the client.
				swarmers[i], err = NewSwarmer(clients[i], stores[i], α, &ecdsaKey)
				Expect(err).ShouldNot(HaveOccurred())
			}

			ctx, cancelCtx := context.WithCancel(context.Background())
			defer cancelCtx()

			// Bootstrapping created clients
			dispatch.CoForAll(numberOfClients, func(i int) {
				defer GinkgoRecover()

				for j := 0; j < numberOfBootstrapClients; j++ {
					if i == j {
						continue
					}
					if _, err := stores[i].PutMultiAddress(multiAddresses[j]); err != nil {
						Expect(err).ShouldNot(HaveOccurred())
					}
				}
			})

			dispatch.CoForAll(numberOfClients, func(i int) {
				defer GinkgoRecover()

				server := NewServer(swarmers[i], stores[i], α)
				serverHub.Register(multiAddresses[i].Address(), server)

				err := swarmers[i].Ping(ctx)
				Expect(err).ShouldNot(HaveOccurred())

				for j := 0; j < numberOfClients; j++ {
					if i == j {
						continue
					}
					if serverHub.IsRegistered(multiAddresses[j].Address()) {
						multiAddr, err := swarmers[i].Query(ctx, multiAddresses[j].Address())
						if err != nil {
							continue
						}
						Expect(multiAddr.String()).To(Equal(multiAddresses[j].String()))
					}
				}
			})

			dispatch.CoForAll(numberOfClients, func(i int) {
				defer GinkgoRecover()

				peers, err := swarmers[i].Peers()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(peers)).To(BeNumerically(">=", (numberOfClients - (numberOfClients / 10))))
			})
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
	store       MultiAddressStorer
	multiAddr   identity.MultiAddress
	isAdversary bool
	serverHub   *mockServerHub
}

func newMockClientToServer(mockServerHub *mockServerHub, i int, isAdversary bool) (mockClientToServer, MultiAddressStorer, error) {
	multiAddr, err := testutils.RandomMultiAddress()
	if err != nil {
		return mockClientToServer{}, nil, err
	}

	// Create leveldb store and store own multiaddress.
	db, err := leveldb.NewStore(fmt.Sprintf("./tmp/swarmer.%v.out", i+1), 72*time.Hour)
	Expect(err).ShouldNot(HaveOccurred())
	store := db.SwarmMultiAddressStore()
	_, err = store.PutMultiAddress(multiAddr)
	Expect(err).ShouldNot(HaveOccurred())

	return mockClientToServer{
		store:       store,
		multiAddr:   multiAddr,
		isAdversary: isAdversary,
		serverHub:   mockServerHub,
	}, store, nil
}

func (client *mockClientToServer) Ping(ctx context.Context, to identity.MultiAddress, multiAddr identity.MultiAddress) error {
	var server Server
	isActive := false

	client.serverHub.connsMu.Lock()
	if isActive, _ = client.serverHub.active[to.Address()]; isActive {
		server = client.serverHub.conns[to.Address()]
	}
	client.serverHub.connsMu.Unlock()

	if isActive {
		randomSleep()
		return server.Ping(ctx, multiAddr)
	}
	return errors.New("address not active")
}

func (client *mockClientToServer) Pong(ctx context.Context, multiAddr identity.MultiAddress) error {
	var server Server
	isActive := false

	client.serverHub.connsMu.Lock()
	if isActive, _ = client.serverHub.active[multiAddr.Address()]; isActive {
		server = client.serverHub.conns[multiAddr.Address()]
	}
	client.serverHub.connsMu.Unlock()

	_, err := client.store.MultiAddress(client.multiAddr.Address())
	if err != nil {
		return err
	}

	if isActive {
		randomSleep()
		return server.Pong(ctx, client.multiAddr)
	}
	return errors.New("pong address not active")
}

func (client *mockClientToServer) Query(ctx context.Context, to identity.MultiAddress, query identity.Address) (identity.MultiAddresses, error) {
	var server Server
	isActive := false

	if client.isAdversary && rand.Uint64()%2 == 0 {
		return identity.MultiAddresses{}, nil
	}

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

func (client *mockClientToServer) MultiAddress() identity.MultiAddress {
	return client.multiAddr
}

func randomSleep() {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(120)
	time.Sleep(time.Duration(r) * time.Millisecond)
}
