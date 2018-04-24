package client_test

import (
	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"
	// . "github.com/republicprotocol/republic-go/rpc/client"
)

var _ = Describe("Client connection pools", func() {

	Context("NewConnPool method", func() {

		It("should return an empty connection pool object with cacheLimit set", func() {

		})
	})

	Context("Dial method", func() {

		It("should return existing connection with updated timestamp when it is present in cache", func() {

		})

		Context("when it is not present in cache", func() {

			It("should return new connection", func() {

			})

			Context("and when cache is full", func() {

				It("should delete the oldest connection in cache and return new connection", func() {

				})
			})
		})
	})

	Context("on calling Close method", func() {

		It("should close all connections in the pool", func() {

		})
	})
})
