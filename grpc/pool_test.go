package grpc_test

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/grpc"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
)

var _ = Describe("Client connection pools", func() {

	var connPool ConnPool

	Context("when dialing connections", func() {

		It("should return new connection if not already present in cache or an existing connection", func() {
			n := 10
			connPool = NewConnPool(n)
			Expect(connPool).ShouldNot(BeNil())

			// Keep adding connections until the cache limit has reached
			var multiaddress identity.MultiAddress
			var err error

			for i := 0; i < n; i++ {
				multiaddress, err = generateMultiaddress()
				Expect(err).ShouldNot(HaveOccurred())

				conn, err := connPool.Dial(context.Background(), multiaddress)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(conn).ShouldNot(BeNil())
			}

			// Store a copy of connPool for testing
			connPoolBeforeClose := connPool

			// Testing an existing connection is returned when the multiaddress is provided
			conn, err := connPool.Dial(context.Background(), multiaddress)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(conn).ShouldNot(BeNil())

			// Confirm that a new connection has not been added to the connection pool
			Expect(connPoolBeforeClose).Should(Equal(connPool))

			// Testing a new connection request is fulfilled on an full cache
			multiaddress, err = generateMultiaddress()
			Expect(err).ShouldNot(HaveOccurred())

			conn, err = connPool.Dial(context.Background(), multiaddress)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(conn).ShouldNot(BeNil())

			// Testing that Close clears all connections.
			connPoolBeforeClose = connPool
			err = connPool.Close()
			Expect(err).ShouldNot(HaveOccurred())

			// Confirm that all connections have been cleared
			Expect(connPoolBeforeClose).ShouldNot(Equal(connPool))
		})
	})
})

func generateMultiaddress() (identity.MultiAddress, error) {
	keystore, err := crypto.RandomKeystore()
	if err != nil {
		return identity.MultiAddress{}, err
	}
	multiaddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/80/republic/%s", keystore.Address()))
	if err != nil {
		return identity.MultiAddress{}, err
	}
	return multiaddress, nil
}
