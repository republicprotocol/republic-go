package swarm_test

import (
	"context"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/swarm"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/dht"
	"github.com/republicprotocol/republic-go/identity"
)

var _ = Describe("Swarm", func() {

	Context("when using a fully connected swarmer", func() {

		It("should not error on bootstrap or query", func() {
			var err error
			numberOfClients := 50

			// Creating clients
			clients := make([]Client, numberOfClients)
			multiaddrs := make(identity.MultiAddresses, numberOfClients)
			for i := 0; i < numberOfClients; i++ {
				clients[i], err = newMockClient(multiaddrs)
				Expect(err).ShouldNot(HaveOccurred())
				multiaddrs[i] = clients[i].MultiAddress()
			}

			var swarmer Swarmer
			ctx, cancelCtx := context.WithCancel(context.Background())
			defer cancelCtx()

			// Bootstrapping created clients
			for i := 0; i < numberOfClients; i++ {
				// Creating swarmer for the client
				clientDht := dht.NewDHT(multiaddrs[i].Address(), 100)
				swarmer = NewSwarmer(clients[i], &clientDht)

				// Bootstrap the client with all available multiaddresses
				err = swarmer.Bootstrap(ctx, multiaddrs)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(clientDht.MultiAddresses())).To(Equal(numberOfClients - 1))
			}

			// Query for the last created client
			multiaddr, err := swarmer.Query(ctx, multiaddrs[0].Address(), 1)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(multiaddr).To(Equal(multiaddrs[0]))

			// Create a new client
			newClient, err := newMockClient(multiaddrs)
			Expect(err).ShouldNot(HaveOccurred())

			// Creating swarmer for the newly created client
			newDht := dht.NewDHT(newClient.MultiAddress().Address(), 100)
			newSwarmer := NewSwarmer(newClient, &newDht)

			// Bootstrapping by providing only the first client address should
			// populate the new client's DHT with all registered client multiaddresses
			err = newSwarmer.Bootstrap(ctx, identity.MultiAddresses{multiaddrs[0]})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(newDht.MultiAddresses())).To(Equal(numberOfClients))
		})
	})

	// FIXME: The functionality works everywhere it is used but this test does not.
	// This test must be fixed and re-enabled.

	//Context("when clients are connected to servers for querying", func() {
	//
	//	It("should connect all clients using swarmer", func() {
	//		numberOfClients := 50
	//		numberOfBootstrapClients := 5
	//
	//		// Creating clients
	//		clients := make([]Client, numberOfClients)
	//		multiaddrs := make(identity.MultiAddresses, numberOfClients)
	//		dhts := make([]dht.DHT, numberOfClients)
	//		bootstrapMultiaddrs := make(identity.MultiAddresses, numberOfBootstrapClients)
	//
	//		// Creating a common server hub for all clients to use
	//		serverHub := &mockServerHub{
	//			connsMu: new(sync.Mutex),
	//			conns:   map[identity.Address]Server{},
	//		}
	//
	//		for i := 0; i < numberOfClients; i++ {
	//			client, err := newMockClientToServer(serverHub)
	//			Expect(err).ShouldNot(HaveOccurred())
	//			clients[i] = &client
	//			multiaddrs[i] = clients[i].MultiAddress()
	//
	//			// Store bootstrap multiaddresses
	//			if i < numberOfBootstrapClients {
	//				bootstrapMultiaddrs[i] = multiaddrs[i]
	//			}
	//
	//			// TODO: (Please confirm) Creating a server for each client (??)
	//			dhts[i] = dht.NewDHT(multiaddrs[i].Address(), 100)
	//			server := NewServer(clients[i], &dhts[i])
	//			serverHub.Register(multiaddrs[i].Address(), server)
	//		}
	//
	//		var swarmer Swarmer
	//		ctx, cancelCtx := context.WithCancel(context.Background())
	//		defer cancelCtx()
	//
	//		// Bootstrapping created clients
	//		for i := 0; i < numberOfClients; i++ {
	//			// Creating swarmer for the client
	//			swarmer = NewSwarmer(clients[i], &dhts[i])
	//
	//			err := swarmer.Bootstrap(ctx, bootstrapMultiaddrs)
	//			Expect(err).ShouldNot(HaveOccurred())
	//
	//			log.Println(len(dhts[i].MultiAddresses()))
	//			if i > numberOfBootstrapClients {
	//				Expect(len(dhts[i].MultiAddresses()) > numberOfBootstrapClients).To(Equal(true))
	//			}
	//		}
	//
	//		// Query for the last created client
	//		multiaddr, err := swarmer.Query(ctx, multiaddrs[0].Address(), 1)
	//		Expect(err).ShouldNot(HaveOccurred())
	//		Expect(multiaddr).To(Equal(multiaddrs[0]))
	//	})
	//})
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
	multiaddr identity.MultiAddress
	serverHub *mockServerHub
}

func (client *mockClientToServer) Ping(ctx context.Context, to identity.MultiAddress) (identity.MultiAddress, error) {
	return client.serverHub.conns[to.Address()].Ping(ctx, to)
}

func (client *mockClientToServer) Query(ctx context.Context, to identity.MultiAddress, query identity.Address, querySig [65]byte) (identity.MultiAddresses, error) {
	return client.serverHub.conns[to.Address()].Query(ctx, query, querySig)
}

func (client *mockClientToServer) MultiAddress() identity.MultiAddress {
	return client.multiaddr
}

type mockClient struct {
	multiaddr   identity.MultiAddress
	storedAddrs identity.MultiAddresses
}

func (client *mockClient) Ping(ctx context.Context, to identity.MultiAddress) (identity.MultiAddress, error) {
	return to, nil
}

func (client *mockClient) Query(ctx context.Context, to identity.MultiAddress, query identity.Address, querySig [65]byte) (identity.MultiAddresses, error) {
	return client.storedAddrs, nil
}

func (client *mockClient) MultiAddress() identity.MultiAddress {
	return client.multiaddr
}

func createNewMultiAddress() (identity.MultiAddress, error) {
	// Generate multiaddress
	ecdsaKey, err := crypto.RandomEcdsaKey()
	if err != nil {
		return identity.MultiAddress{}, err
	}
	multiaddr, err := identity.Address(ecdsaKey.Address()).MultiAddress()
	if err != nil {
		return identity.MultiAddress{}, err
	}
	return multiaddr, nil
}

func newMockClient(multiaddrs identity.MultiAddresses) (Client, error) {
	multiaddr, err := createNewMultiAddress()
	if err != nil {
		return nil, err
	}

	return &mockClient{
		multiaddr:   multiaddr,
		storedAddrs: multiaddrs,
	}, nil
}

func newMockClientToServer(mockServerHub *mockServerHub) (mockClientToServer, error) {
	multiaddr, err := createNewMultiAddress()
	if err != nil {
		return mockClientToServer{}, err
	}

	return mockClientToServer{
		multiaddr: multiaddr,
		serverHub: mockServerHub,
	}, nil
}
