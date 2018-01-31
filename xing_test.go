package rpc_test

import (
	"context"
	"fmt"
	"math/big"
	"net"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-order-compute"
	"github.com/republicprotocol/go-rpc"
	"github.com/republicprotocol/go-sss"
	"google.golang.org/grpc"
)

func (s *mockServer) SendOrderFragment(ctx context.Context, orderFragment *rpc.OrderFragment) (*rpc.Nothing, error) {
	return &rpc.Nothing{}, nil
}

func (s *mockServer) SendResultFragment(ctx context.Context, resultFragment *rpc.ResultFragment) (*rpc.Nothing, error) {
	return &rpc.Nothing{}, nil
}

var _ = Describe("Xing Overlay Network", func() {
	var server *grpc.Server
	var rpcServer mockServer
	var rpcClient mockClient
	var defaultTimeout = time.Second
	var orderFragment *compute.OrderFragment
	var resultFragment *compute.ResultFragment
	var badTargetAddress identity.MultiAddress
	var err error

	createServe := func() {
		keyPair, err := identity.NewKeyPair()
		Ω(err).ShouldNot(HaveOccurred())
		multiAddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/republic/%s", 3000, keyPair.Address()))
		Ω(err).ShouldNot(HaveOccurred())
		server = grpc.NewServer()
		rpcServer = mockServer{MultiAddress: multiAddress}
		rpc.RegisterXingNodeServer(server, &rpcServer)
	}

	createClient := func() {
		keyPair, err := identity.NewKeyPair()
		Ω(err).ShouldNot(HaveOccurred())
		multiAddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/republic/%s", 4000, keyPair.Address()))
		Ω(err).ShouldNot(HaveOccurred())
		rpcClient = mockClient{MultiAddress: multiAddress}
	}

	createFragments := func() {
		sssShare := sss.Share{Key: 1, Value: &big.Int{}}
		orderFragment = compute.NewOrderFragment([]byte("orderID"), compute.OrderTypeIBBO, compute.OrderParityBuy,
			sssShare, sssShare, sssShare, sssShare, sssShare)
		resultFragment = compute.NewResultFragment([]byte("butOrderID"), []byte("sellOrderID"),
			[]byte("butOrderFragmentID"), []byte("sellOrderFragmentID"),
			sssShare, sssShare, sssShare, sssShare, sssShare)
		badTargetAddress, err = identity.NewMultiAddressFromString("/ip4/192.168.0.1/republic/8MHzQ7ZQDvvT8Nqo3HLQQDZvfcHJYB")
		Ω(err).ShouldNot(HaveOccurred())
	}

	BeforeEach(func() {
		createClient()
		createServe()
		createFragments()
	})

	Context("sending order fragments", func() {
		It("should return nothing", func() {
			lis, err := net.Listen("tcp", ":3000")
			Ω(err).ShouldNot(HaveOccurred())
			go func(server *grpc.Server) {
				defer GinkgoRecover()
				Ω(server.Serve(lis)).ShouldNot(HaveOccurred())
			}(server)
			defer server.Stop()

			err = rpc.SendOrderFragmentToTarget(rpcServer.MultiAddress, *orderFragment, defaultTimeout)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should return an error for bad multi-addresses", func() {
			err = rpc.SendOrderFragmentToTarget(badTargetAddress, *orderFragment, defaultTimeout)
			Ω(err).Should(HaveOccurred())
		})

		It("should return a timeout error when there is no response within the timeout duration", func() {
			err := rpc.SendOrderFragmentToTarget(rpcServer.MultiAddress, *orderFragment, defaultTimeout)
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("sending result fragments", func() {
		It("should return nothing", func() {
			lis, err := net.Listen("tcp", ":3000")
			Ω(err).ShouldNot(HaveOccurred())
			go func(server *grpc.Server) {
				defer GinkgoRecover()
				Ω(server.Serve(lis)).ShouldNot(HaveOccurred())
			}(server)
			defer server.Stop()

			err = rpc.SendResultFragmentToTarget(rpcServer.MultiAddress, *resultFragment, defaultTimeout)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should return an error for bad multi-addresses", func() {
			err = rpc.SendResultFragmentToTarget(badTargetAddress, *resultFragment, defaultTimeout)
			Ω(err).Should(HaveOccurred())
		})

		It("should return a timeout error when there is no response within the timeout duration", func() {
			err := rpc.SendResultFragmentToTarget(rpcServer.MultiAddress, *resultFragment, defaultTimeout)
			Ω(err).Should(HaveOccurred())
		})
	})
})
