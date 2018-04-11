package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/republicprotocol/republic-go/cmd/darknode"
)

const ganacheRPC = "http://localhost:8545"
const numberOfDarknodes = 48
const numberOfBootstrapDarknodes = 5

var _ = Describe("Darknode", func() {

	Context("when watching the ocean", func() {

		It("should update local views of the ocean", func() {

			conn := ganache.Connect(ganacheRPC)

			darknodes, ctxs, cancels := NewLocalDarknodes()
			for i := range darknodes {
				darknodes[i].Run(ctx[i])
			}

		})

		It("should converge on a global view of the ocean", func() {

		})

		It("should persist computations from recent epochs", func() {

		})

		It("should not persist computations from distant epochs", func() {

		})
	})

	Context("when computing order matches", func() {

		It("should process the distribute order table in parallel with other pools", func() {

		})

		It("should update the order book after computing an order match", func() {

		})

	})

	Context("when confirming order matches", func() {

		It("should update the order book after confirming an order match", func() {

		})

		It("should update the order book after releasing an order match", func() {

		})
	})
})
