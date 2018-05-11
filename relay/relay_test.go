package relay_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/relay"
)

var _ = Describe("Relay", func() {

	Context("when opening orders", func() {

		It("should open orders with a sufficient number of order fragments", func() {
			Expect("implemented").To(Equal("unimplemented"))
		})

		It("should not open orders with an insufficient number of order fragments", func() {
			Expect("implemented").To(Equal("unimplemented"))
		})

		It("should not open orders with malformed order fragments", func() {
			Expect("implemented").To(Equal("unimplemented"))
		})

		It("should not open orders that have not been signed correctly", func() {
			Expect("implemented").To(Equal("unimplemented"))
		})

	})

	Context("when canceling orders", func() {

		It("should cancel orders that are open", func() {
			Expect("implemented").To(Equal("unimplemented"))
		})

		It("should not cancel orders that are not open", func() {
			Expect("implemented").To(Equal("unimplemented"))
		})

		It("should not cancel orders that have not been signed correctly", func() {
			Expect("implemented").To(Equal("unimplemented"))
		})

	})

	Context("when getting orders", func() {

		It("should return an order for verified traders", func() {
			Expect("implemented").To(Equal("unimplemented"))
		})

		It("should not return an order for unverified traders", func() {
			Expect("implemented").To(Equal("unimplemented"))
		})

		It("should return orders for verified traders", func() {
			Expect("implemented").To(Equal("unimplemented"))
		})

		It("should not return orders for unverified traders", func() {
			Expect("implemented").To(Equal("unimplemented"))
		})

	})

})
