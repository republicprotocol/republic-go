package client_test

import (
	"context"
	"fmt"
	"net"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/rpc/client"

	"github.com/republicprotocol/republic-go/identity"
	"google.golang.org/grpc"
)

var _ = Describe("Client connections", func() {

	var server *grpc.Server

	BeforeSuite(func(done Done) {
		defer close(done)

		server = grpc.NewServer()
		listener, err := net.Listen("tcp", "127.0.0.1:3000")
		Expect(err).ShouldNot(HaveOccurred())
		go func() {
			defer GinkgoRecover()
			err := server.Serve(listener)
			Expect(err).ShouldNot(HaveOccurred())
		}()
	})

	AfterSuite(func() {
		server.Stop()
	})

	Context("Dial method", func() {

		It("should return a connection object if valid multiaddress is provided", func() {

			multiAddr, err := createConnMultiAddress("127.0.0.1", "3000")
			Expect(err).ShouldNot(HaveOccurred())

			conn, err := Dial(context.Background(), multiAddr)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(conn).ShouldNot(BeNil())

		})
	})

	Context("Clone method", func() {

		It("should return a clone of the connection object", func() {
			conn := createConn()
			connUpdate := conn.Clone()
			Expect(connUpdate).To(Equal(conn))
		})
	})

	Context("Close method", func() {

		It("should shutdown connection object if all its references have been closed", func() {
			conn := createConn()
			err := conn.Close()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(conn.ClientConn.GetState().String()).Should(Equal("SHUTDOWN"))
		})

		It("should not SHUTDOWN connection object if it has atleast one reference alive", func() {
			conn := createConn()
			conn = conn.Clone()
			err := conn.Close()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(conn.ClientConn.GetState().String()).ShouldNot((Equal("SHUTDOWN")))
		})
	})
})

func createConn() *Conn {
	connMultiAddress, err := createConnMultiAddress("127.0.0.1", "3000")
	Expect(err).ShouldNot(HaveOccurred())
	conn, err := Dial(context.Background(), connMultiAddress)
	Expect(err).ShouldNot(HaveOccurred())
	return conn
}

func createConnMultiAddress(host, port string) (identity.MultiAddress, error) {
	addr, _, err := identity.NewAddress()
	if err != nil {
		return identity.MultiAddress{}, err
	}
	return identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/%s/tcp/%s/republic/%v", host, port, addr))
}
