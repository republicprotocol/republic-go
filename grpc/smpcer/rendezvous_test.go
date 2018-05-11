package smpcer_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/grpc/smpcer"
)

var _ = Describe("Rendezvous", func() {

	Context("NewRendezvous method", func() {

		It("should return a Rendezvous object", func() {
			Expect(NewRendezvous()).ShouldNot(BeNil())
		})
	})
})
