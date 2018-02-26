package rpc_test

import (
	"context"
	"fmt"
	"net"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-order-compute"
	"github.com/republicprotocol/go-rpc"
	"google.golang.org/grpc"
)

func (s *mockServer) Sync(syncRequest *rpc.SyncRequest, stream rpc.DarkNode_SyncServer) error {
	stream.Send(&rpc.SyncBlock{})
	return nil
}

func (s *mockServer) ElectShard(ctx context.Context, electShardRequest *rpc.ElectShardRequest) (*rpc.Shard, error) {
	return &rpc.Shard{}, nil
}

func (s *mockServer) ComputeShard(ctx context.Context, computeShardRequest *rpc.ComputeShardRequest) (*rpc.Nothing, error) {
	return &rpc.Nothing{}, nil
}

func (s *mockServer) FinalizeShard(ctx context.Context, finalizeShardRequest *rpc.FinalizeShardRequest) (*rpc.Nothing, error) {
	return &rpc.Nothing{}, nil
}

var _ = Describe("Dark Node", func() {
	var server *grpc.Server
	var rpcServer mockServer
	var rpcClient mockClient

	createServe := func() {
		keyPair, err := identity.NewKeyPair()
		Ω(err).ShouldNot(HaveOccurred())
		multiAddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/republic/%s", 3000, keyPair.Address()))
		Ω(err).ShouldNot(HaveOccurred())
		server = grpc.NewServer()
		rpcServer = mockServer{MultiAddress: multiAddress}
		rpc.RegisterDarkNodeServer(server, &rpcServer)
	}

	createClient := func() {
		keyPair, err := identity.NewKeyPair()
		Ω(err).ShouldNot(HaveOccurred())
		multiAddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/republic/%s", 4000, keyPair.Address()))
		Ω(err).ShouldNot(HaveOccurred())
		rpcClient = mockClient{MultiAddress: multiAddress}
	}

	BeforeEach(func() {
		createClient()
		createServe()
	})

	Context("sync order book", func() {
		It("should stream back results by using a channel", func() {
			lis, err := net.Listen("tcp", ":3000")
			Ω(err).ShouldNot(HaveOccurred())
			go func(server *grpc.Server) {
				defer GinkgoRecover()
				Ω(server.Serve(lis)).ShouldNot(HaveOccurred())
			}(server)
			defer server.Stop()

			resultChan, _ := rpc.SyncWithTarget(rpcServer.MultiAddress, rpcClient.MultiAddress, defaultTimeout)
			syncBlock := <-resultChan
			Ω(syncBlock.Err).Should(BeNil())
		})
	})

	Context("elect shard", func() {
		It("should return no error", func() {
			lis, err := net.Listen("tcp", ":3000")
			Ω(err).ShouldNot(HaveOccurred())
			go func(server *grpc.Server) {
				defer GinkgoRecover()
				Ω(server.Serve(lis)).ShouldNot(HaveOccurred())
			}(server)
			defer server.Stop()

			inputShard := compute.NewShard([]*compute.DeltaFragment{}, []*compute.ResidueFragment{})
			shard, err := rpc.StartElectShard(rpcServer.MultiAddress, rpcClient.MultiAddress, inputShard, defaultTimeout)
			Ω(*shard).Should(Equal(rpc.Shard{}))
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	Context("compute shard", func() {
		It("should return no error", func() {
			lis, err := net.Listen("tcp", ":3000")
			Ω(err).ShouldNot(HaveOccurred())
			go func(server *grpc.Server) {
				defer GinkgoRecover()
				Ω(server.Serve(lis)).ShouldNot(HaveOccurred())
			}(server)
			defer server.Stop()

			inputShard := compute.NewShard([]*compute.DeltaFragment{}, []*compute.ResidueFragment{})
			err = rpc.AskToComputeShard(rpcServer.MultiAddress, rpcClient.MultiAddress, inputShard, defaultTimeout)
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	Context("compute shard", func() {
		It("should return no error", func() {
			lis, err := net.Listen("tcp", ":3000")
			Ω(err).ShouldNot(HaveOccurred())
			go func(server *grpc.Server) {
				defer GinkgoRecover()
				Ω(server.Serve(lis)).ShouldNot(HaveOccurred())
			}(server)
			defer server.Stop()

			inputShard := compute.NewShard([]*compute.DeltaFragment{}, []*compute.ResidueFragment{})
			err = rpc.AskToComputeShard(rpcServer.MultiAddress, rpcClient.MultiAddress, inputShard, defaultTimeout)
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	Context("delta shard", func() {
		It("should return no error", func() {
			lis, err := net.Listen("tcp", ":3000")
			Ω(err).ShouldNot(HaveOccurred())
			go func(server *grpc.Server) {
				defer GinkgoRecover()
				Ω(server.Serve(lis)).ShouldNot(HaveOccurred())
			}(server)
			defer server.Stop()

			inputShard := compute.NewDeltaShard()
			err = rpc.FinalizeShard(rpcServer.MultiAddress, rpcClient.MultiAddress, inputShard, defaultTimeout)
			Ω(err).ShouldNot(HaveOccurred())
		})
	})
})
