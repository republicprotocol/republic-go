package swarm_test

import (
	"context"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/swarm"

	"github.com/republicprotocol/republic-go/dht"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/testutils"
)

var _ = Describe("Swarm", func() {

	Context("when bootstrapping", func() {

		It("should be able to query any peer after bootstrapping", func() {
			numberOfClients := 50
			numberOfBootstrapClients := 5

			// Creating clients
			dhts := make([]dht.DHT, numberOfClients)
			clients := make([]Client, numberOfClients)
			multiAddrs := make(identity.MultiAddresses, numberOfClients)
			bootstrapMultiaddrs := make(identity.MultiAddresses, numberOfBootstrapClients)

			// Creating a common server hub for all clients to use
			serverHub := &mockServerHub{
				connsMu: new(sync.Mutex),
				conns:   map[identity.Address]Server{},
			}

			for i := 0; i < numberOfClients; i++ {
				client, err := newMockClientToServer(serverHub)
				Expect(err).ShouldNot(HaveOccurred())
				clients[i] = &client
				multiAddrs[i] = clients[i].MultiAddress()

				// Store bootstrap multiAddresses
				if i < numberOfBootstrapClients {
					bootstrapMultiaddrs[i] = multiAddrs[i]
				}

				dhts[i] = dht.NewDHT(multiAddrs[i].Address(), numberOfClients)
				server := NewServer(clients[i], &dhts[i])
				serverHub.Register(multiAddrs[i].Address(), server)
			}

			ctx, cancelCtx := context.WithCancel(context.Background())
			defer cancelCtx()

			// Bootstrapping created clients
			dispatch.CoForAll(numberOfClients, func(i int) {
				defer GinkgoRecover()

				// Creating swarmer for the client
				swarmer := NewSwarmer(clients[i], &dhts[i])
				err := swarmer.Bootstrap(ctx, bootstrapMultiaddrs)
				Expect(err).ShouldNot(HaveOccurred())
			})

			// Query for clients
			for i := 0; i < numberOfClients; i++ {
				for j := 0; j < numberOfClients; j++ {
					if i == j {
						continue
					}
					multiAddr, err := NewSwarmer(clients[i], &dhts[i]).Query(ctx, multiAddrs[j].Address(), -1)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(multiAddr).To(Equal(multiAddrs[j]))
				}
			}
		})
	})
})

// mockServerHub will store all Servers that Clients use to Query and Ping
type mockServerHub struct {
	connsMu *sync.Mutex
	conns   map[identity.Address]Server
}

func (serverHub *mockServerHub) Register(serverAddr identity.Address, server Server) {
	serverHub.connsMu.Lock()
	defer serverHub.connsMu.Unlock()

	// If server is not already present, add it to the serverHub
	if _, ok := serverHub.conns[serverAddr]; !ok {
		serverHub.conns[serverAddr] = server
	}
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

func (client *mockClientToServer) Ping(ctx context.Context, to identity.MultiAddress) (identity.MultiAddress, error) {
	return client.serverHub.conns[to.Address()].Ping(ctx, client.multiAddr)
}

func (client *mockClientToServer) Query(ctx context.Context, to identity.MultiAddress, query identity.Address, querySig [65]byte) (identity.MultiAddresses, error) {
	return client.serverHub.conns[to.Address()].Query(ctx, query, querySig)
}

func (client *mockClientToServer) MultiAddress() identity.MultiAddress {
	return client.multiAddr
}

type mockClient struct {
	multiAddr   identity.MultiAddress
	storedAddrs identity.MultiAddresses
}

func newMockClient(multiAddrs identity.MultiAddresses) (Client, error) {
	multiAddr, err := testutils.RandomMultiAddress()
	if err != nil {
		return nil, err
	}

	return &mockClient{
		multiAddr:   multiAddr,
		storedAddrs: multiAddrs,
	}, nil
}

func (client *mockClient) Ping(ctx context.Context, to identity.MultiAddress) (identity.MultiAddress, error) {
	return to, nil
}

func (client *mockClient) Query(ctx context.Context, to identity.MultiAddress, query identity.Address, querySig [65]byte) (identity.MultiAddresses, error) {
	return client.storedAddrs, nil
}

func (client *mockClient) MultiAddress() identity.MultiAddress {
	return client.multiAddr
}
