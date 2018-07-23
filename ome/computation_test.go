package ome_test

import (
	"bytes"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/ome"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/testutils"
)

var _ = Describe("Computations", func() {
	var buyFragment, sellFragment order.Fragment

	BeforeEach(func() {
		buyFragments, err := testutils.RandomBuyOrderFragments(6, 4)
		Expect(err).ShouldNot(HaveOccurred())
		buyFragment = buyFragments[0]
		sellFragments, err := testutils.RandomSellOrderFragments(6, 4)
		Expect(err).ShouldNot(HaveOccurred())
		sellFragment = sellFragments[0]
	})

	Context("when checking for equality", func() {
		It("should return true for equal computation IDs", func() {
			computationID := NewComputationID(buyFragment.OrderID, sellFragment.OrderID)
			expectedID := ComputationID{}
			copy(expectedID[:], crypto.Keccak256(buyFragment.OrderID[:], sellFragment.OrderID[:]))
			Ω(bytes.Equal(computationID[:], expectedID[:]))
		})

		It("should return true for the same computations compared against itself", func() {
			computation := NewComputation([32]byte{}, buyFragment, sellFragment, ComputationStateNil, false)
			Expect(computation.Equal(&computation)).Should(BeTrue())
			Expect(computation.ID.String()).Should(Equal(computation.ID.String()))
		})

		It("should return true for equal computations compared against each other", func() {
			lhs := NewComputation([32]byte{}, buyFragment, sellFragment, ComputationStateNil, false)
			rhs := NewComputation([32]byte{}, buyFragment, sellFragment, ComputationStateNil, false)
			lhs.Timestamp = rhs.Timestamp
			Expect(lhs.Equal(&rhs)).Should(BeTrue())
			Expect(lhs.ID.String()).Should(Equal(rhs.ID.String()))
		})

		It("should return false for unequal computations compared against each other", func() {
			lhs := NewComputation([32]byte{}, buyFragment, sellFragment, ComputationStateNil, false)
			rhs := NewComputation([32]byte{}, sellFragment, buyFragment, ComputationStateNil, false)
			Expect(lhs.Equal(&rhs)).Should(BeFalse())
			Expect(lhs.ID.String()).ShouldNot(Equal(rhs.ID.String()))
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
})
