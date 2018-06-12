package ome_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/ome"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/smpc"
	"github.com/republicprotocol/republic-go/testutils"
)

var _ = Describe("Settler", func() {
	var storer Storer
	var smpcer smpc.Smpcer
	var accounts cal.DarkpoolAccounts

	BeforeEach(func() {
		storer = testutils.NewStorer()
		smpcer = testutils.NewAlwaysMatchSmpc()
		accounts = testutils.NewDarkpoolAccounts()
	})

	Context("when a computation has been resolved to a match and been confirmed ", func() {
		It("should be able to reconstruct the order and settle it.", func() {
			// todo

			Ω(true).Should(BeTrue())
		})
	})
})
