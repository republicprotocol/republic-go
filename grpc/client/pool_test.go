package client_test

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/crypto"
	. "github.com/republicprotocol/republic-go/grpc/client"

	"github.com/republicprotocol/republic-go/identity"
)

var _ = Describe("Client connection pools", func() {

	var connPool ConnPool

	Context("NewConnPool method", func() {

		It("should return an empty connection pool object with cacheLimit set", func() {
			Expect(NewConnPool(5)).ShouldNot(BeNil())
		})
	})

	Context("Dial method", func() {

		Context("when it is not present in cache", func() {
			connPool = NewConnPool(5)
			It("should return new connection", func() {
				keystore, err := crypto.RandomKeystore()
				Expect(err).ShouldNot(HaveOccurred())
				multiaddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/80/republic/%s", keystore.Address()))
				Expect(err).ShouldNot(HaveOccurred())
				conn, err := connPool.Dial(context.Background(), multiaddress)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(conn).ShouldNot(BeNil())
			})

			Context("and when cache is full", func() {

				It("should delete the oldest connection in cache and return new connection", func() {

				})
			})
		})

		It("should return existing connection with updated timestamp when it is present in cache", func() {

		})

	})

	Context("on calling Close method", func() {

		It("should close all connections in the pool", func() {

		})
	})
})
