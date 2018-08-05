package swarm_test

import (
	"context"
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
	"github.com/republicprotocol/republic-go/registry"
	"github.com/republicprotocol/republic-go/testutils"
)

var _ = Describe("Swarm", func() {

	var (
		numberOfClients          = 100
		numberOfBootstrapClients = 5
		α                        = 10
	)

	registerClientsAndBootstrap := func(ctx context.Context, honest bool) ([]Client, []Swarmer, *testutils.MockServerHub, error) {
		// Creating clients.
		stores := make([]MultiAddressStorer, numberOfClients)
		clients := make([]Client, numberOfClients)
		swarmers := make([]Swarmer, numberOfClients)
		verifiers := make([]registry.Crypter, numberOfClients)

		serverHub := &testutils.MockServerHub{
			ConnsMu: new(sync.Mutex),
			Conns:   map[identity.Address]Server{},
			Active:  map[identity.Address]bool{},
		}

		for i := 0; i < numberOfClients; i++ {
			key, err := crypto.RandomKeystore()
			if err != nil {
				return nil, nil, serverHub, err
			}
			verifiers[i] = registry.NewCrypter(key, testutils.NewMockSwarmBinder(), numberOfClients, time.Hour)

			clientType := testutils.Honest
			if !honest && i > numberOfBootstrapClients {
				clientType = testutils.ClientTypes[rand.Intn(len(testutils.ClientTypes))]
			}

			client, store, err := testutils.NewMockSwarmClient(serverHub, clientType, &verifiers[i])
			if err != nil {
				return nil, nil, serverHub, err
			}
			clients[i] = &client
			stores[i] = store
			alpha := rand.Intn(α-2) + 2
			swarmers[i] = NewSwarmer(clients[i], stores[i], alpha, &verifiers[i])

			server := NewServer(swarmers[i], stores[i], alpha, &verifiers[i])
			serverHub.Register(clients[i].MultiAddress().Address(), server)
		}

		// Bootstrapping created clients
		dispatch.CoForAll(numberOfClients, func(i int) {
			defer GinkgoRecover()

			for j := 0; j < numberOfBootstrapClients; j++ {
				if err := stores[i].InsertMultiAddress(clients[j].MultiAddress()); err != nil {
					Expect(err).ShouldNot(HaveOccurred())
				}
			}
		})
		return clients, swarmers, serverHub, nil
	}

	Context("when connecting to the network", func() {

		Context("when traders are honest, they", func() {

			AfterEach(func() {
				os.RemoveAll("./tmp")
			})

			It("should be able to query most peers", func() {
				ctx, cancelCtx := context.WithCancel(context.Background())
				defer cancelCtx()

				clients, swarmers, serverHub, err := registerClientsAndBootstrap(ctx, true)
				Expect(err).ShouldNot(HaveOccurred())

				log.Println("Ping self-address to join the network")
				dispatch.CoForAll(numberOfClients, func(i int) {
					defer GinkgoRecover()

					err := swarmers[i].Ping(ctx)
					Expect(err).ShouldNot(HaveOccurred())
				})

				log.Println("Query other peers address")
				dispatch.CoForAll(numberOfClients, func(i int) {
					defer GinkgoRecover()

					for j := 0; j < numberOfClients; j++ {
						if i == j {
							continue
						}
						if serverHub.IsRegistered(clients[j].MultiAddress().Address()) {
							multiAddr, err := swarmers[i].Query(ctx, clients[j].MultiAddress().Address())
							if err != nil {
								log.Printf("cannot query %v, %v", clients[j].MultiAddress().Address(), err)
								continue
							}
							Expect(multiAddr.String()).To(Equal(clients[j].MultiAddress().String()))
						}
					}
				})

				log.Println("Check number of connected nodes")
				dispatch.CoForAll(numberOfClients, func(i int) {
					defer GinkgoRecover()

					peers, err := swarmers[i].Peers()
					Expect(err).ShouldNot(HaveOccurred())
					Expect(len(peers)).To(BeNumerically(">=", numberOfClients*9/10))
				})
			})
		})

		Context("when some traders are dishonest, they", func() {

			AfterEach(func() {
				os.RemoveAll("./tmp")
			})

			It("should be able to connect to at least the bootstrap nodes", func() {
				ctx, cancelCtx := context.WithCancel(context.Background())
				defer cancelCtx()

				_, swarmers, _, err := registerClientsAndBootstrap(ctx, false)
				Expect(err).ShouldNot(HaveOccurred())

				dispatch.CoForAll(numberOfClients, func(i int) {
					defer GinkgoRecover()

					err := swarmers[i].Ping(ctx)
					Expect(err).ShouldNot(HaveOccurred())
				})

				time.Sleep(2 * time.Second)
				dispatch.CoForAll(numberOfClients, func(i int) {
					defer GinkgoRecover()

					peers, err := swarmers[i].Peers()
					Expect(err).ShouldNot(HaveOccurred())
					Expect(len(peers)).To(BeNumerically(">", numberOfBootstrapClients))
				})
			})
		})
	})
})
