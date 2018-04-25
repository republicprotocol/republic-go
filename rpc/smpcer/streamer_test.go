package smpcer_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/rpc/smpcer"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/rpc/client"
)

var _ = Describe("Streamer", func() {

	Context("NewStreamer method", func() {

		It("should return a Streamer object", func() {
			addr, _, err := identity.NewAddress()
			Expect(err).ShouldNot(HaveOccurred())
			multiaddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/80/republic/%v", addr))
			Expect(err).ShouldNot(HaveOccurred())
			connPool := client.NewConnPool(5)

			Expect(NewStreamer(multiaddress, &connPool)).ShouldNot(BeNil())
		})
	})
})
