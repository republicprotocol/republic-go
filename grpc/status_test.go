package grpc_test

import (
	"context"
	"fmt"
	"net"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/grpc"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/testutils"
)

var _ = Describe("Status", func() {

	var serviceMultiAddr identity.MultiAddress
	var service StatusService
	var server *Server
	var swarmer testutils.Swarmer

	BeforeEach(func() {
		var err error

		keystore, err := crypto.RandomKeystore()
		Expect(err).ShouldNot(HaveOccurred())

		serviceMultiAddr, err = identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/0.0.0.0/tcp/18514/republic/%v", keystore.Address()))
		Expect(err).ShouldNot(HaveOccurred())

		swarmer = testutils.NewMockSwarmer(serviceMultiAddr)
		service = NewStatusService(&swarmer)
		server = NewServer()
		service.Register(server)

		lis, err := net.Listen("tcp", "0.0.0.0:18514")
		Expect(err).ShouldNot(HaveOccurred())

		go server.Serve(lis)
	})

	AfterEach(func() {
		server.Stop()
	})

	Context("when getting the status", func() {
		It("should return the expected status information", func() {
			NumberOfPeers := 5
			conn, err := Dial(context.Background(), serviceMultiAddr)
			Expect(err).ShouldNot(HaveOccurred())

			for i := 0; i < NumberOfPeers; i++ {
				multi, err := testutils.RandomMultiAddress()
				Expect(err).ShouldNot(HaveOccurred())
				swarmer.PutMultiAddress(multi, 1)
			}

			client := NewStatusServiceClient(conn)
			status, err := client.Status(context.Background(), &StatusRequest{})
			Expect(err).ShouldNot(HaveOccurred())

			Expect(status.Address).Should(Equal(serviceMultiAddr.Address().String()))
			Expect(status.Peers).Should(Equal(int64(NumberOfPeers)))
		})
	})
})
