package orderbook_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/orderbook"
)

var _ = Describe("orderBookDB", func() {
	Context("opening new DB connection", func() {
		It("should error for invalid filepath", func() {
			db, err := orderbook.NewDatabase("")
			Ω(db).ShouldNot(Equal(nil))
			Ω(err).Should(HaveOccurred())
		})
	})
})
