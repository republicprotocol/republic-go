package rpc_test

import (
	"context"
	"fmt"
	"math/big"
	"net"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-order-compute"
	"github.com/republicprotocol/go-rpc"
	"github.com/republicprotocol/go-sss"
	"google.golang.org/grpc"
)

//var result = &compute.Final{
//	ID:          []byte("resultID"),
//	BuyOrderID:  []byte("BuyOrderID"),
//	SellOrderID: []byte("SellOrderID"),
//	FstCode:     big.NewInt(0),
//	SndCode:     big.NewInt(0),
//	Price:       big.NewInt(0),
//	MaxVolume:   big.NewInt(0),
//	MinVolume:   big.NewInt(0),
//}

func (s *mockServer) SendOrderFragmentCommitment(ctx context.Context, orderFragmentCommitment *rpc.OrderFragmentCommitment) (*rpc.OrderFragmentCommitment, error) {
	return &rpc.OrderFragmentCommitment{}, nil
}

func (s *mockServer) SendOrderFragment(ctx context.Context, orderFragment *rpc.OrderFragment) (*rpc.Nothing, error) {
	return &rpc.Nothing{}, nil
}

//func (s *mockServer) Notifications(multiAddress *rpc.MultiAddress, stream rpc.DarkNode_NotificationsServer) error {
//	stream.Send(rpc.SerializeFinal(result))
//	return nil
//}
//
//func (s *mockServer) GetFinals(multiAddress *rpc.MultiAddress, stream rpc.DarkNode_GetFinalsServer) error {
//	stream.Send(rpc.SerializeFinal(result))
//	return nil
//}

var _ = Describe("Dark Network", func() {
	var server *grpc.Server
	var rpcServer mockServer
	var rpcClient mockClient

	var orderFragment *compute.OrderFragment
	var badServerAddress identity.MultiAddress
	var err error

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

	createFragments := func() {
		sssShare := sss.Share{Key: 1, Value: &big.Int{}}
		orderFragment = compute.NewOrderFragment([]byte("orderID"), compute.OrderTypeIBBO, compute.OrderParityBuy,
			sssShare, sssShare, sssShare, sssShare, sssShare)
		badServerAddress, err = identity.NewMultiAddressFromString("/ip4/192.168.0.1/republic/8MHzQ7ZQDvvT8Nqo3HLQQDZvfcHJYB")
		Ω(err).ShouldNot(HaveOccurred())
	}

	BeforeEach(func() {
		createClient()
		createServe()
		createFragments()
	})

	Context("sending order fragments commitments", func() {

		It("should return nothing", func() {
			lis, err := net.Listen("tcp", ":3000")
			Ω(err).ShouldNot(HaveOccurred())
			go func(server *grpc.Server) {
				defer GinkgoRecover()
				Ω(server.Serve(lis)).ShouldNot(HaveOccurred())
			}(server)
			defer server.Stop()

			err = rpc.SendOrderFragmentCommitmentToTarget(rpcServer.MultiAddress, rpcClient.MultiAddress, defaultTimeout)
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	Context("sending order fragments", func() {
		keyPair, _ := identity.NewKeyPair()
		to := keyPair.Address()
		from := rpcClient.MultiAddress

		It("should return nothing", func() {
			lis, err := net.Listen("tcp", ":3000")
			Ω(err).ShouldNot(HaveOccurred())
			go func(server *grpc.Server) {
				defer GinkgoRecover()
				Ω(server.Serve(lis)).ShouldNot(HaveOccurred())
			}(server)
			defer server.Stop()

			err = rpc.SendOrderFragmentToTarget(rpcServer.MultiAddress, to, from, orderFragment, defaultTimeout)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should return an error for bad multi-addresses", func() {
			err = rpc.SendOrderFragmentToTarget(badServerAddress, to, from, orderFragment, defaultTimeout)
			Ω(err).Should(HaveOccurred())
		})

		It("should return a timeout error when there is no response within the timeout duration", func() {
			err := rpc.SendOrderFragmentToTarget(rpcServer.MultiAddress, to, from, orderFragment, defaultTimeout)
			Ω(err).Should(HaveOccurred())
		})
	})

	//Context("getting notifications from target", func() {
	//	It("should be able to query notifications from the server", func() {
	//		lis, err := net.Listen("tcp", ":3000")
	//		Ω(err).ShouldNot(HaveOccurred())
	//		go func(server *grpc.Server) {
	//			defer GinkgoRecover()
	//			Ω(server.Serve(lis)).ShouldNot(HaveOccurred())
	//		}(server)
	//		defer server.Stop()
	//		resultChan, _ := rpc.NotificationsFromTarget(rpcServer.MultiAddress, rpcClient.MultiAddress, defaultTimeout)
	//		res := <-resultChan
	//
	//		Ω(res.Ok.(*compute.Final)).Should(Equal(result))
	//	})
	//
	//	It("should return an error when dialing an offline server", func() {
	//		resultChan, _ := rpc.NotificationsFromTarget(rpcServer.MultiAddress, rpcClient.MultiAddress, defaultTimeout)
	//		res := <-resultChan
	//		Ω(res.Err).ShouldNot(BeNil())
	//	})
	//})
	//
	//Context("getting results from the target", func() {
	//	It("should return nothing", func() {
	//		lis, err := net.Listen("tcp", ":3000")
	//		Ω(err).ShouldNot(HaveOccurred())
	//		go func(server *grpc.Server) {
	//			defer GinkgoRecover()
	//			Ω(server.Serve(lis)).ShouldNot(HaveOccurred())
	//		}(server)
	//		defer server.Stop()
	//
	//		results, err := rpc.GetFinalsFromTarget(rpcServer.MultiAddress, rpcClient.MultiAddress, defaultTimeout)
	//		Ω(err).ShouldNot(HaveOccurred())
	//		Ω(len(results)).Should(BeNumerically(">", 0))
	//		Ω(results[0]).Should(Equal(result))
	//	})
	//
	//	It("should return an error when dialing an offline server", func() {
	//		_, err := rpc.GetFinalsFromTarget(rpcServer.MultiAddress, rpcClient.MultiAddress, defaultTimeout)
	//		Ω(err).Should(HaveOccurred())
	//	})
	//})

})
