package rpc_test

import (
	"context"
	"fmt"
	"net"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-rpc"
	"google.golang.org/grpc"
)

type mockServer struct {
	identity.MultiAddress
}

func (s *mockServer) Ping(ctx context.Context, address *rpc.MultiAddress) (*rpc.Nothing, error) {
	return &rpc.Nothing{}, nil
}

func (s *mockServer) QueryCloserPeers(ctx context.Context, query *rpc.Query) (*rpc.MultiAddresses, error) {
	return &rpc.MultiAddresses{}, nil
}

func (s *mockServer) QueryCloserPeersOnFrontier(query *rpc.Query, stream rpc.SwarmNode_QueryCloserPeersOnFrontierServer) error {
	multiAddress, err := identity.NewMultiAddressFromString("/ip4/127.0.0.1/tcp/80/republic/8MGfbzAMS59Gb4cSjpm34soGNYsM2f")
	if err != nil {
		return err

	}
	return  stream.Send(rpc.SerializeMultiAddress(multiAddress))
}

type mockClient struct {
	identity.MultiAddress
}

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

	Context("Query close peers from a target", func() {
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
