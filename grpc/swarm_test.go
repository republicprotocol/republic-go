package grpc_test

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/grpc"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/dht"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/swarm"
	"github.com/republicprotocol/republic-go/testutils"
)

var _ = Describe("Swarming", func() {

	var server *Server
	var service SwarmService
	var serviceDHT dht.DHT
	var serviceMultiAddr identity.MultiAddress
	var serviceClient swarm.Client
	var client swarm.Client

	BeforeEach(func() {
		var err error

		serviceClient, err = newSwarmClient()
		Expect(err).ShouldNot(HaveOccurred())

		serviceDHT = dht.NewDHT(serviceClient.MultiAddress().Address(), 20)
		service = NewSwarmService(swarm.NewServer(testutils.NewCrypter(), serviceClient, &serviceDHT))
		serviceMultiAddr = serviceClient.MultiAddress()
		server = NewServer()
		service.Register(server)

		client, err = newSwarmClient()
		Expect(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		server.Stop()
	})

	Context("when pinging a service", func() {

		It("should return the multiaddress of the service", func(done Done) {
			defer close(done)

			go func() {
				defer GinkgoRecover()

				err := server.Start("0.0.0.0:18514")
				Expect(err).ShouldNot(HaveOccurred())
			}()
			time.Sleep(time.Millisecond)

			multiAddr, err := client.Ping(context.Background(), serviceMultiAddr)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(multiAddr.String()).Should(Equal(serviceMultiAddr.String()))
		})

		It("should add the client to the service DHT", func(done Done) {
			defer close(done)

			go func() {
				defer GinkgoRecover()

				err := server.Start("0.0.0.0:18514")
				Expect(err).ShouldNot(HaveOccurred())
			}()
			time.Sleep(time.Millisecond)

			_, err := client.Ping(context.Background(), serviceMultiAddr)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(serviceDHT.MultiAddresses()).Should(HaveLen(1))
		})

	})

	Context("when querying a service", func() {

		It("should return the multiaddress of the service close to the query", func(done Done) {
			defer close(done)

			go func() {
				defer GinkgoRecover()

				err := server.Start("0.0.0.0:18514")
				Expect(err).ShouldNot(HaveOccurred())
			}()
			time.Sleep(time.Millisecond)

			_, err := client.Ping(context.Background(), serviceMultiAddr)
			Expect(err).ShouldNot(HaveOccurred())

			multiAddrs, err := client.Query(context.Background(), serviceMultiAddr, client.MultiAddress().Address(), [65]byte{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(multiAddrs).Should(HaveLen(1))
		})

	})
})

func newSwarmClient() (swarm.Client, error) {
	ecdsaKey, err := crypto.RandomEcdsaKey()
	if err != nil {
		return nil, err
	}
	addr := identity.Address(ecdsaKey.Address())
	multiAddr, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/0.0.0.0/tcp/18514/republic/%v", addr))
	if err != nil {
		return nil, err
	}
	client := NewSwarmClient(multiAddr)
	return client, nil
}
