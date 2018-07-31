package swarm_test

import (
	"context"
	"math/rand"
	"os"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/swarm"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/testutils"
)

var _ = Describe("Swarm", func() {

	var (
		numberOfClients          = 100
		numberOfBootstrapClients = 5
		α                        = 4
	)

	registerClientsAndBootstrap := func(ctx context.Context, honest bool) ([]Client, []Swarmer, *testutils.MockServerHub, error) {
		// Creating clients.
		stores := make([]MultiAddressStorer, numberOfClients)
		clients := make([]Client, numberOfClients)
		swarmers := make([]Swarmer, numberOfClients)
		ecdsaKeys := make([]crypto.EcdsaKey, numberOfClients)

		serverHub := &testutils.MockServerHub{
			ConnsMu: new(sync.Mutex),
			Conns:   map[identity.Address]Server{},
			Active:  map[identity.Address]bool{},
		}

		for i := 0; i < numberOfClients; i++ {
			key, err := crypto.RandomEcdsaKey()
			if err != nil {
				return nil, nil, serverHub, err
			}
			clientType := testutils.Honest
			if !honest {
				clientType = testutils.ClientTypes[rand.Intn(len(testutils.ClientTypes))]
			}
			client, store, err := testutils.NewMockSwarmClient(serverHub, &key, clientType)
			if err != nil {
				return nil, nil, serverHub, err
			}
			clients[i] = &client
			ecdsaKeys[i] = key
			stores[i] = store
			swarmers[i] = NewSwarmer(clients[i], stores[i], α, &ecdsaKeys[i])
		}

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
		return clients, swarmers, serverHub, nil
	}

	Context("when connecting to the network", func() {

		AfterEach(func() {
			os.RemoveAll("./tmp")
		})

		Context("when traders are honest, they", func() {

			It("should be able to query most peers", func() {
				ctx, cancelCtx := context.WithCancel(context.Background())
				defer cancelCtx()

				clients, swarmers, serverHub, err := registerClientsAndBootstrap(ctx, true)
				Expect(err).ShouldNot(HaveOccurred())

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

		Context("when some traders are dishonest, they", func() {

			It("should be able to connect to atleast the bootstrap nodes", func() {
				ctx, cancelCtx := context.WithCancel(context.Background())
				defer cancelCtx()

				_, swarmers, _, err := registerClientsAndBootstrap(ctx, false)
				Expect(err).ShouldNot(HaveOccurred())

				dispatch.CoForAll(numberOfClients, func(i int) {
					defer GinkgoRecover()

					err := swarmers[i].Ping(ctx)
					Expect(err).ShouldNot(HaveOccurred())
				})

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
