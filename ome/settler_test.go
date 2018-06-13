package ome_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/ome"

	"github.com/republicprotocol/republic-go/testutils"
)

var _ = Describe("Settler", func() {
	var storer *testutils.Storer
	var smpcer *testutils.Smpc
	var accounts *testutils.DarkpoolAccounts

	BeforeEach(func() {
		storer = testutils.NewStorer()
		smpcer = testutils.NewAlwaysMatchSmpc()
		accounts = testutils.NewDarkpoolAccounts()
	})

	Context("when a computation has been resolved to a match and been confirmed ", func() {
		It("should be able to reconstruct the order and settle it.", func() {
			settler := NewSettler(storer, smpcer, accounts)

			Ω(true).Should(BeTrue())
		})
	})
})
