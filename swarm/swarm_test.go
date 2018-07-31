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
			// Creating a common server hub for all clients to use.
			serverHub := &mockServerHub{
				connsMu: new(sync.Mutex),
				conns:   map[identity.Address]Server{},
				active:  map[identity.Address]bool{},
			}

			clients, keys, err := createClients(numberOfClients, serverHub)
			Expect(err).ShouldNot(HaveOccurred())
			swarmers, err := createClients()

		})
	})

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

func createClients(numberOfClients int, hub *mockServerHub) ([]Client, []crypto.EcdsaKey, error) {
	clients := make([]Client, numberOfClients)
	keys := make([]crypto.EcdsaKey, numberOfClients)
	for i := 0; i < numberOfClients; i++ {
		key, err := crypto.RandomEcdsaKey()
		if err != nil {
			return nil, nil, err
		}
		client, err := newMockSwarmClient(hub, i, &key)
		if err != nil {
			return nil, nil, err
		}
		clients[i] = &client
		keys[i] = key
	}
}

func createSwarmers(clients []Client, keys []crypto.EcdsaKey) ([]Swarmer, error) {
	numberOfClients := len(clients)
	swarmers := make([]Client, numberOfClients)
	for i := 0; i < numberOfClients; i++ {

		swarmers[i] = NewSwarmer(clients[i], stores[i], α, &ecdsaKeys[i])
	}
}
