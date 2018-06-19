package ome_test

import (
	"bytes"
	"encoding/base64"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/ome"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/testutils"
)

var _ = Describe("Computations", func() {
	var buy, sell order.Order

	BeforeEach(func() {
		buy, sell = testutils.RandomBuyOrder(), testutils.RandomSellOrder()
	})

	Context("computations id ", func() {
		It("should be able to compare two computation ID", func() {
			computationID := NewComputationID(buy.ID, sell.ID)
			expectedID := ComputationID{}
			copy(expectedID[:], crypto.Keccak256(buy.ID[:], sell.ID[:]))
			Ω(bytes.Equal(computationID[:], expectedID[:]))
		})

		It("should implement the Stringer interface ", func() {
			computationID := NewComputationID(buy.ID, sell.ID)
			idInString := fmt.Sprintf("%v", computationID)

			expectedString := base64.StdEncoding.EncodeToString(computationID[:8])
			Ω(idInString).Should(Equal(expectedString))
		})
	})

	Context("computations state ", func() {
		It("should implement the Stringer interface ", func() {
			Ω(fmt.Sprintf("%v", ComputationStateNil)).Should(Equal("nil"))
			Ω(fmt.Sprintf("%v", ComputationStateMatched)).Should(Equal("matched"))
			Ω(fmt.Sprintf("%v", ComputationStateMismatched)).Should(Equal("mismatched"))
			Ω(fmt.Sprintf("%v", ComputationStateAccepted)).Should(Equal("accepted"))
			Ω(fmt.Sprintf("%v", ComputationStateRejected)).Should(Equal("rejected"))
			Ω(fmt.Sprintf("%v", ComputationStateSettled)).Should(Equal("settled"))
			Ω(fmt.Sprintf("%v", ComputationState(100))).Should(Equal("unsupported state"))
		})
	})

	Context("computation", func() {
		It("should be comparable", func() {
			comp := testutils.RandomComputation()
			another := testutils.RandomComputation()
			Ω(comp.Equal(&another)).Should(BeFalse())
			Ω(comp.Equal(&comp)).Should(BeTrue())
			Ω(another.Equal(&another)).Should(BeTrue())
		})
	})
})
