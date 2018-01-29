package rpc_test

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-rpc"
	"google.golang.org/grpc"
	"net"
	"time"
)

type mockServer struct {
	Multi identity.MultiAddress
}

func (s *mockServer) Ping(ctx context.Context, address *rpc.MultiAddress) (*rpc.Nothing, error) {
	return &rpc.Nothing{}, nil
}

func (s *mockServer) Peers(ctx context.Context, address *rpc.MultiAddress) (*rpc.MultiAddresses, error) {
	return &rpc.MultiAddresses{}, nil
}

func (s *mockServer) QueryCloserPeers(ctx context.Context, address *rpc.Query) (*rpc.MultiAddresses, error) {
	return &rpc.MultiAddresses{}, nil
}

type mockClient struct {
	Multi identity.MultiAddress
}

var _ = Describe("Swarm node", func() {
	var server *grpc.Server
	var rpcServer mockServer
	var rpcClient mockClient
	var badMulti1 identity.MultiAddress
	var badMulti2 identity.MultiAddress
	var err error
	var defaultTimeout = time.Second * 5

	createServe := func() {
		keyPair, err := identity.NewKeyPair()
		Ω(err).ShouldNot(HaveOccurred())
		multi, err := identity.NewMultiAddressFromString(
			fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/republic/%s", 3000, keyPair.Address()))
		Ω(err).ShouldNot(HaveOccurred())

		server = grpc.NewServer()
		rpcServer = mockServer{Multi: multi}
		rpc.RegisterSwarmNodeServer(server, &rpcServer)
	}

	createClient := func() {
		keyPair, err := identity.NewKeyPair()
		Ω(err).ShouldNot(HaveOccurred())
		multi, err := identity.NewMultiAddressFromString(
			fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/republic/%s", 4000, keyPair.Address()))
		Ω(err).ShouldNot(HaveOccurred())
		rpcClient = mockClient{Multi: multi}
	}

	BeforeEach(func() {
		createClient()
		createServe()
		badMulti1, err = identity.NewMultiAddressFromString("/ip4/192.168.0.1/republic/8MHzQ7ZQDvvT8Nqo3HLQQDZvfcHJYB")
		Ω(err).ShouldNot(HaveOccurred())
		badMulti2, err = identity.NewMultiAddressFromString("/tcp/80/republic/8MHzQ7ZQDvvT8Nqo3HLQQDZvfcHJYB")
		Ω(err).ShouldNot(HaveOccurred())
	})

	Context("Ping target", func() {
		It("should be able to ping target through grpc", func() {

			lis, err := net.Listen("tcp", ":3000")
			Ω(err).ShouldNot(HaveOccurred())
			go func(server *grpc.Server) {
				defer GinkgoRecover()
				Ω(server.Serve(lis)).ShouldNot(HaveOccurred())
			}(server)
			defer server.Stop()

			err = rpc.PingTarget(rpcServer.Multi, rpcClient.Multi, defaultTimeout)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should return an error if pinging a wrong address", func() {
			err = rpc.PingTarget(badMulti1, rpcClient.Multi, defaultTimeout)
			Ω(err).Should(HaveOccurred())
			err = rpc.PingTarget(badMulti2, rpcClient.Multi, defaultTimeout)
			Ω(err).Should(HaveOccurred())
		})

		It("should return an timeout error if getting no response in a minute", func() {
			err := rpc.PingTarget(rpcServer.Multi, rpcClient.Multi, defaultTimeout)
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("Get peers from target", func() {
		It("should be able to get peers through grpc", func() {
			lis, err := net.Listen("tcp", ":3000")
			Ω(err).ShouldNot(HaveOccurred())
			go func(server *grpc.Server) {
				defer GinkgoRecover()
				Ω(server.Serve(lis)).ShouldNot(HaveOccurred())
			}(server)
			defer server.Stop()

			multi, err := rpc.GetPeersFromTarget(rpcServer.Multi, rpcClient.Multi, defaultTimeout)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(multi).Should(Equal(identity.MultiAddresses{}))
		})

		It("should return an error if calling a bad address", func() {
			_, err = rpc.GetPeersFromTarget(badMulti1, rpcClient.Multi, defaultTimeout)
			Ω(err).Should(HaveOccurred())
			_, err = rpc.GetPeersFromTarget(badMulti2, rpcClient.Multi, defaultTimeout)
			Ω(err).Should(HaveOccurred())
		})

		It("should return an timeout error if getting no response in a minute", func() {
			_, err := rpc.GetPeersFromTarget(rpcServer.Multi, rpcClient.Multi, defaultTimeout)
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("Query close peers", func() {
		It("should be able to query close peers through grpc", func() {
			lis, err := net.Listen("tcp", ":3000")
			Ω(err).ShouldNot(HaveOccurred())
			go func(server *grpc.Server) {
				defer GinkgoRecover()
				Ω(server.Serve(lis)).ShouldNot(HaveOccurred())
			}(server)
			defer server.Stop()

			target := identity.Address("8MHzQ7ZQDvvT8Nqo3HLQQDZvfcHJYB")
			multi, err := rpc.QueryCloserPeersFromTarget(rpcServer.Multi, rpcClient.Multi, target, true, defaultTimeout)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(multi).Should(Equal(identity.MultiAddresses{}))
		})

		It("should return an error if quering a bad address", func() {
			target := identity.Address("8MHzQ7ZQDvvT8Nqo3HLQQDZvfcHJYB")
			_, err = rpc.QueryCloserPeersFromTarget(badMulti1, rpcClient.Multi, target, true, defaultTimeout)
			Ω(err).Should(HaveOccurred())
			_, err = rpc.QueryCloserPeersFromTarget(badMulti2, rpcClient.Multi, target, true, defaultTimeout)
			Ω(err).Should(HaveOccurred())
		})

		It("should return an timeout error if getting no response in a minute", func() {
			target := identity.Address("8MHzQ7ZQDvvT8Nqo3HLQQDZvfcHJYB")
			_, err := rpc.QueryCloserPeersFromTarget(rpcServer.Multi, rpcClient.Multi, target, true, defaultTimeout)
			Ω(err).Should(HaveOccurred())
		})
	})
})
