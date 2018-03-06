package rpc_test

import (
	"context"
	"fmt"
	"math/big"
	"net"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-order-compute"
	sss "github.com/republicprotocol/go-sss"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/network/rpc"
	"google.golang.org/grpc"
)

func (s *mockServer) BroadcastDeltaFragment(ctx context.Context, broadcastDeltaFragmentRequest *rpc.BroadcastDeltaFragmentRequest) (*rpc.DeltaFragment, error) {
	return &rpc.DeltaFragment{}, nil
}

var _ = Describe("Custom Client", func() {
	Context("Broadcast delta fragments ", func() {
		var from, to identity.MultiAddress
		var sssShare = sss.Share{Key: 1, Value: &big.Int{}}
		deltaFragment := &compute.DeltaFragment{
			ID:                  []byte("byte"),
			BuyOrderID:          []byte("byte"),
			SellOrderID:         []byte("byte"),
			BuyOrderFragmentID:  []byte("byte"),
			SellOrderFragmentID: []byte("byte"),
			FstCodeShare:        sssShare,
			SndCodeShare:        sssShare,
			PriceShare:          sssShare,
			MaxVolumeShare:      sssShare,
			MinVolumeShare:      sssShare,
		}

		BeforeEach(func() {
			address, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			from, err = identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/3000/republic/%v", address.String()))
			Ω(err).ShouldNot(HaveOccurred())

			address, _, err = identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			to, err = identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/4000/republic/%v", address.String()))
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should not error", func() {
			server := grpc.NewServer()
			rpcServer := mockServer{MultiAddress: to}
			rpc.RegisterDarkNodeServer(server, &rpcServer)
			lis, err := net.Listen("tcp", ":4000")
			Ω(err).ShouldNot(HaveOccurred())
			go func(server *grpc.Server) {
				defer GinkgoRecover()
				Ω(server.Serve(lis)).ShouldNot(HaveOccurred())
			}(server)
			defer server.Stop()

			client, err := rpc.NewClient(to, from)
			Ω(err).ShouldNot(HaveOccurred())

			_, err = client.BroadcastDeltaFragment(deltaFragment)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should error with bad-formatted multiaddress", func() {
			multiWithouIP, err := identity.NewMultiAddressFromString(fmt.Sprintf("/tcp/3000/republic/%v", from.Address().String()))
			Ω(err).ShouldNot(HaveOccurred())
			multiWithouPort, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/republic/%v", from.Address().String()))
			Ω(err).ShouldNot(HaveOccurred())

			_, err = rpc.NewClient(multiWithouIP, from)
			Ω(err).Should(HaveOccurred())
			_, err = rpc.NewClient(multiWithouPort, from)
			Ω(err).Should(HaveOccurred())

		})

		It("should error when the server if offline", func() {
			client, err := rpc.NewClient(to, from, rpc.ClientOptions{
				Timeout:        5 * time.Second,
				TimeoutBackoff: 0 * time.Second,
				TimeoutRetries: 3,
			})
			Ω(err).ShouldNot(HaveOccurred())
			_, err = client.BroadcastDeltaFragment(deltaFragment)
			Ω(err).Should(HaveOccurred())
		})
	})
})

var _ = Describe("Swarm node", func() {

	var err error
	var server *grpc.Server
	var rpcServer mockServer
	var rpcClient mockClient
	var multiAddressMissingHost identity.MultiAddress
	var multiAddressMissingPort identity.MultiAddress
	var defaultTimeout = time.Second

	createServe := func() {
		keyPair, err := identity.NewKeyPair()
		Ω(err).ShouldNot(HaveOccurred())
		multiAddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/republic/%s", 3000, keyPair.Address()))
		Ω(err).ShouldNot(HaveOccurred())
		server = grpc.NewServer()
		rpcServer = mockServer{MultiAddress: multiAddress}
		rpc.RegisterSwarmNodeServer(server, &rpcServer)
	}

	createClient := func() {
		keyPair, err := identity.NewKeyPair()
		Ω(err).ShouldNot(HaveOccurred())
		multiAddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/republic/%s", 4000, keyPair.Address()))
		Ω(err).ShouldNot(HaveOccurred())
		rpcClient = mockClient{MultiAddress: multiAddress}
	}

	createBadMultiAddresses := func() {
		multiAddressMissingPort, err = identity.NewMultiAddressFromString("/ip4/192.168.0.1/republic/8MHzQ7ZQDvvT8Nqo3HLQQDZvfcHJYB")
		Ω(err).ShouldNot(HaveOccurred())
		multiAddressMissingHost, err = identity.NewMultiAddressFromString("/tcp/80/republic/8MHzQ7ZQDvvT8Nqo3HLQQDZvfcHJYB")
		Ω(err).ShouldNot(HaveOccurred())
	}

	BeforeEach(func() {
		createClient()
		createServe()
		createBadMultiAddresses()
	})

	Context("pinging a target", func() {
		It("should return nothing", func() {
			lis, err := net.Listen("tcp", ":3000")
			Ω(err).ShouldNot(HaveOccurred())
			go func(server *grpc.Server) {
				defer GinkgoRecover()
				Ω(server.Serve(lis)).ShouldNot(HaveOccurred())
			}(server)
			defer server.Stop()

			err = rpc.PingTarget(rpcServer.MultiAddress, rpcClient.MultiAddress, defaultTimeout)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should return an error for bad multi-addresses", func() {
			err = rpc.PingTarget(multiAddressMissingPort, rpcClient.MultiAddress, defaultTimeout)
			Ω(err).Should(HaveOccurred())
			err = rpc.PingTarget(multiAddressMissingHost, rpcClient.MultiAddress, defaultTimeout)
			Ω(err).Should(HaveOccurred())
		})

		It("should return a timeout error when there is no response within the timeout duration", func() {
			err := rpc.PingTarget(rpcServer.MultiAddress, rpcClient.MultiAddress, defaultTimeout)
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("query close peers from a target", func() {
		target := identity.Address("8MHzQ7ZQDvvT8Nqo3HLQQDZvfcHJYB")

		It("should return all peers", func() {
			lis, err := net.Listen("tcp", ":3000")
			Ω(err).ShouldNot(HaveOccurred())
			go func(server *grpc.Server) {
				defer GinkgoRecover()
				Ω(server.Serve(lis)).ShouldNot(HaveOccurred())
			}(server)
			defer server.Stop()

			closePeers, err := rpc.QueryCloserPeersFromTarget(rpcServer.MultiAddress, rpcClient.MultiAddress, target, defaultTimeout)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(closePeers).Should(Equal(identity.MultiAddresses{}))
		})

		It("should return an error for bad multi-addresses", func() {
			_, err = rpc.QueryCloserPeersFromTarget(multiAddressMissingPort, rpcClient.MultiAddress, target, defaultTimeout)
			Ω(err).Should(HaveOccurred())
			_, err = rpc.QueryCloserPeersFromTarget(multiAddressMissingHost, rpcClient.MultiAddress, target, defaultTimeout)
			Ω(err).Should(HaveOccurred())
		})

		It("should return a timeout error when there is no response within the timeout duration", func() {
			_, err := rpc.QueryCloserPeersFromTarget(rpcServer.MultiAddress, rpcClient.MultiAddress, target, defaultTimeout)
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("query closer peers on frontier ", func() {
		It("should return multi-addresses", func() {
			lis, err := net.Listen("tcp", ":3000")
			Ω(err).ShouldNot(HaveOccurred())
			go func(server *grpc.Server) {
				defer GinkgoRecover()
				Ω(server.Serve(lis)).ShouldNot(HaveOccurred())
			}(server)
			defer server.Stop()

			target := identity.Address("8MHzQ7ZQDvvT8Nqo3HLQQDZvfcHJYB")
			multiAddresses, err := rpc.QueryCloserPeersOnFrontierFromTarget(rpcServer.MultiAddress, rpcClient.MultiAddress, target, defaultTimeout)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(len(multiAddresses)).Should(Equal(1))
		})

		It("should return an error for bad multi-addresses", func() {
			target := identity.Address("8MHzQ7ZQDvvT8Nqo3HLQQDZvfcHJYB")
			_, err = rpc.QueryCloserPeersOnFrontierFromTarget(multiAddressMissingPort, rpcClient.MultiAddress, target, defaultTimeout)
			Ω(err).Should(HaveOccurred())
			_, err = rpc.QueryCloserPeersOnFrontierFromTarget(multiAddressMissingHost, rpcClient.MultiAddress, target, defaultTimeout)
			Ω(err).Should(HaveOccurred())
		})

		It("should return a timeout error when there is no response within the timeout duration", func() {
			target := identity.Address("8MHzQ7ZQDvvT8Nqo3HLQQDZvfcHJYB")
			_, err := rpc.QueryCloserPeersOnFrontierFromTarget(rpcServer.MultiAddress, rpcClient.MultiAddress, target, defaultTimeout)
			Ω(err).Should(HaveOccurred())
		})
	})
})
