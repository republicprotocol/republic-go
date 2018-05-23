package swarm_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/swarm"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/dht"
	"github.com/republicprotocol/republic-go/identity"
)

var _ = Describe("swarm ", func() {

	Context("a completely functioning swarmer", func() {

		It("should not error on bootstrap or query", func() {
			var err error
			numberOfClients := 50

			// Creating clients
			clients := make([]Client, numberOfClients)
			multiaddrs := make(identity.MultiAddresses, numberOfClients)
			for i := 0; i < numberOfClients; i++ {
				clients[i], err = createNewClient(multiaddrs)
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
			newClient, err := createNewClient(multiaddrs)
			Expect(err).ShouldNot(HaveOccurred())

			// Creating swarmer for the newly crsated client
			newDht := dht.NewDHT(newClient.MultiAddress().Address(), 100)
			newSwarmer := NewSwarmer(newClient, &newDht)

			// Bootstrapping by providing only the first client address should
			// populate the new client's DHT with all registered client multiaddresses
			err = newSwarmer.Bootstrap(ctx, identity.MultiAddresses{multiaddrs[0]})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(newDht.MultiAddresses())).To(Equal(numberOfClients))
		})
	})
})

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

func createNewClient(multiaddrs identity.MultiAddresses) (Client, error) {
	multiaddr, err := createNewMultiAddress()
	if err != nil {
		return nil, err
	}

	return &mockClient{
		multiaddr:   multiaddr,
		storedAddrs: multiaddrs,
	}, nil
}
