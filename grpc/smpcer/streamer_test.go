package smpcer_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/grpc/smpcer"
)

var _ = Describe("Streamer", func() {

	Context("NewStreamer method", func() {

		It("should return a Streamer object", func() {
			multiaddress, connPool, err := createMultiAddrAndConnPool("127.0.0.1", "3000")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(NewStreamer(multiaddress, &connPool)).ShouldNot(BeNil())
		})
	})
})
