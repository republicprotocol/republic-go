package connection_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/contracts/connection"
)

var _ = Describe("connection", func() {

	Context("FromURI", func() {
		It("should panic for unimplemented mainnet", func() {
			Ω(func() { connection.FromURI("", "mainnet") }).Should(Panic())
		})

		It("has default connection for ropsten", func() {
			conn, err := connection.FromURI("", "ropsten")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(conn).ShouldNot(BeNil())
		})
	})

})
