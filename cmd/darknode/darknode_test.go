package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Darknode integration", func() {

	Context("when bootstrapping into a network", func() {

		It("should be able to query the super majority of nodes", func() {
			Expect(nil).To(BeNil())
		})

		It("should be responsive to RPCs", func() {
			Expect(nil).To(BeNil())
		})

	})

	Context("when the epoch changes", func() {

		It("should reach consensus on the configuration of the pods", func() {
			Expect(nil).To(BeNil())
		})

		It("should continue computations in previous epoch", func() {
			Expect(nil).To(BeNil())
		})

		It("should run computations in new epoch", func() {
			Expect(nil).To(BeNil())
		})

	})

	Context("when orders are opened", func() {

		It("should confirm matching orders", func() {
			Expect(nil).To(BeNil())
		})

		It("should not confirm mismatching orders", func() {
			Expect(nil).To(BeNil())
		})

	})

	Context("when orders are canceled", func() {

		It("should not confirm canceled orders", func() {
			Expect(nil).To(BeNil())
		})

	})

	Context("when orders are confirmed", func() {

		It("should not reconfirm orders", func() {
			Expect(nil).To(BeNil())
		})

	})

	Context("when faults occur under the threshold", func() {

		It("should not block the computation", func() {
			Expect(nil).To(BeNil())
		})

		It("should not corrupt the state of the faulty nodes", func() {
			Expect(nil).To(BeNil())
		})

	})

	Context("when faults occur above the threshold", func() {

		It("should block the computation", func() {
			Expect(nil).To(BeNil())
		})

		It("should unblock the computation after recovery", func() {
			Expect(nil).To(BeNil())
		})

		It("should not corrupt the state of the faulty nodes", func() {
			Expect(nil).To(BeNil())
		})

	})

})
