package grpc_test

import (
	"context"
	"fmt"
	"net"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/dht"
	. "github.com/republicprotocol/republic-go/grpc"

	"github.com/republicprotocol/republic-go/identity"
)

var _ = Describe("Status", func() {

	var serviceMultiAddr identity.MultiAddress
	var service StatusService
	var server *Server

	BeforeEach(func() {
		var err error

		keystore, err := crypto.RandomKeystore()
		Expect(err).ShouldNot(HaveOccurred())

		serviceMultiAddr, err = identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/0.0.0.0/tcp/18514/republic/%v", keystore.Address()))
		Expect(err).ShouldNot(HaveOccurred())

		dht := dht.NewDHT(serviceMultiAddr.Address(), 128)
		service = NewStatusService(&dht)
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
			conn, err := Dial(context.Background(), serviceMultiAddr)
			Expect(err).ShouldNot(HaveOccurred())

			client := NewStatusServiceClient(conn)
			status, err := client.Status(context.Background(), &StatusRequest{})
			Expect(err).ShouldNot(HaveOccurred())

			Expect(status.Address).Should(Equal(serviceMultiAddr.Address().String()))
		})
	})
})
