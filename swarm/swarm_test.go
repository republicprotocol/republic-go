package swarm_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/swarm"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/leveldb"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var _ = Describe("Swarm", func() {

	var (
		numberOfClients          = 100
		numberOfBootstrapClients = 6
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
			swarmers := make([]Swarmer, numberOfClients)
			ecdsaKeys := make([]crypto.EcdsaKey, numberOfClients)

			// Creating a common server hub for all clients to use.
			serverHub := &mockServerHub{
				connsMu: new(sync.Mutex),
				conns:   map[identity.Address]Server{},
				active:  map[identity.Address]bool{},
			}

			var err error
			// Initialize all the clients
			for i := 0; i < numberOfClients; i++ {
				ecdsaKeys[i], err = crypto.RandomEcdsaKey()
				Expect(err).ShouldNot(HaveOccurred())
				client, store, err := newMockSwarmClient(serverHub, i, &ecdsaKeys[i])
				Expect(err).ShouldNot(HaveOccurred())
				clients[i] = &client
				stores[i] = store

				swarmers[i] = NewSwarmer(clients[i], stores[i], α, &ecdsaKeys[i])
			}

			ctx, cancelCtx := context.WithCancel(context.Background())
			defer cancelCtx()

			// Bootstrapping created clients
			dispatch.CoForAll(numberOfClients, func(i int) {
				defer GinkgoRecover()

				for j := 0; j < numberOfBootstrapClients; j++ {
					if _, err := stores[i].PutMultiAddress(clients[j].MultiAddress()); err != nil {
						Expect(err).ShouldNot(HaveOccurred())
					}
				}
			})

			dispatch.CoForAll(numberOfClients, func(i int) {
				defer GinkgoRecover()

				server := NewServer(swarmers[i], stores[i], α)
				serverHub.Register(clients[i].MultiAddress().Address(), server)
			})

			dispatch.CoForAll(numberOfClients, func(i int) {
				defer GinkgoRecover()

				err := swarmers[i].Ping(ctx)
				Expect(err).ShouldNot(HaveOccurred())

				for j := 0; j < numberOfClients; j++ {
					if i == j {
						continue
					}
					if serverHub.IsRegistered(clients[j].MultiAddress().Address()) {
						multiAddr, err := swarmers[i].Query(ctx, clients[j].MultiAddress().Address())
						if err != nil {
							continue
						}
						Expect(multiAddr.String()).To(Equal(clients[j].MultiAddress().String()))
					}
				}
			})

			dispatch.CoForAll(numberOfClients, func(i int) {
				defer GinkgoRecover()

				peers, err := swarmers[i].Peers()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(peers)).To(BeNumerically(">=", numberOfClients*9/10))
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

func (serverHub *mockServerHub) Deregister(serverAddr identity.Address) {
	serverHub.connsMu.Lock()
	defer serverHub.connsMu.Unlock()

	delete(serverHub.conns, serverAddr)
	serverHub.active[serverAddr] = false
}

func (serverHub *mockServerHub) IsRegistered(serverAddr identity.Address) bool {
	serverHub.connsMu.Lock()
	defer serverHub.connsMu.Unlock()

	isActive, _ := serverHub.active[serverAddr]

	return isActive
}

type mockSwarmClient struct {
	store     MultiAddressStorer
	multiAddr identity.MultiAddress
	serverHub *mockServerHub
}

func newMockSwarmClient(mockServerHub *mockServerHub, i int, key *crypto.EcdsaKey) (mockSwarmClient, MultiAddressStorer, error) {

	multiAddr, err := identity.Address(key.Address()).MultiAddress()
	if err != nil {
		return mockSwarmClient{}, nil, err
	}
	multiAddr.Nonce = 1
	signature, err := key.Sign(multiAddr.Hash())
	Expect(err).ShouldNot(HaveOccurred())
	multiAddr.Signature = signature

	// Create leveldb store and store own multiAddress.
	db, err := leveldb.NewStore(fmt.Sprintf("./tmp/swarmer.%v.out", i+1), 72*time.Hour)
	Expect(err).ShouldNot(HaveOccurred())
	store := db.SwarmMultiAddressStore()
	_, err = store.PutMultiAddress(multiAddr)
	Expect(err).ShouldNot(HaveOccurred())

	return mockSwarmClient{
		store:     store,
		multiAddr: multiAddr,
		serverHub: mockServerHub,
	}, store, nil
}

func (client *mockSwarmClient) Ping(ctx context.Context, to identity.MultiAddress, multiAddr identity.MultiAddress) error {
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

func (client *mockSwarmClient) Pong(ctx context.Context, multiAddr identity.MultiAddress) error {
	var server Server
	isActive := false

	client.serverHub.connsMu.Lock()
	if isActive, _ = client.serverHub.active[multiAddr.Address()]; isActive {
		server = client.serverHub.conns[multiAddr.Address()]
	}
	client.serverHub.connsMu.Unlock()

	multi, err := client.store.MultiAddress(client.multiAddr.Address())
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
	var server Server
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
	multi, err := client.store.MultiAddress(client.multiAddr.Address())
	if err != nil {
		log.Println("err in getting the multiaddress in store")
	}

	return multi
}

func randomSleep() {
	r := rand.Intn(120)
	time.Sleep(time.Duration(r) * time.Millisecond)
}
