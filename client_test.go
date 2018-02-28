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
