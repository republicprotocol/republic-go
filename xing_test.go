package xing_test

import (
	"bytes"
	"sort"

	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"

	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-identity"
	xing "github.com/republicprotocol/go-xing"
)

var _ = Describe("X", func() {

	Context("when assigning the X overlay", func() {

		var miners []xing.Miner
		var overlayMiners []xing.Miner

		BeforeEach(func() {
			epoch, err := generateEpoch()
			Ω(err).ShouldNot(HaveOccurred())

			// Generate overlay miners.
			overlayMiners, err = generateMiners()
			Ω(err).ShouldNot(HaveOccurred())

			// Compute the X Overlay.
			numberOfMNetworks := xing.NumberOfMNetworks(len(overlayMiners))
			xing.AssignXOverlay(overlayMiners, epoch, numberOfMNetworks)

			// Clone miners into a clean slice.
			miners = make([]xing.Miner, len(overlayMiners))
			for i := range overlayMiners {
				miners[i] = xing.NewMiner(overlayMiners[i].ID, overlayMiners[i].Commitment)
			}

			// Compute the individual components of the Miners.
			xing.AssignXHash(miners, epoch)
			xing.AssignClass(miners, numberOfMNetworks)
			xing.AssignMNetwork(miners, numberOfMNetworks)
		})

		It("should assign X hashes", func() {
			for i := range overlayMiners {
				Ω(bytes.Equal(overlayMiners[i].X, miners[i].X)).Should(Equal(true))
			}
		})
		It("should assign classes", func() {
			for i := range overlayMiners {
				Ω(overlayMiners[i].Class).Should(Equal(miners[i].Class))
			}
		})
		It("should assign M networks", func() {
			for i := range overlayMiners {
				Ω(overlayMiners[i].MNetwork).Should(Equal(miners[i].MNetwork))
			}
		})
	})

	Context("when assigning X hashes", func() {
		It("should generate the correct X hashes", func() {
			epoch, err := generateEpoch()
			Ω(err).ShouldNot(HaveOccurred())
			miners, err := generateMiners()
			Ω(err).ShouldNot(HaveOccurred())
			xing.AssignXHash(miners, epoch)
			for _, miner := range miners {
				Ω(bytes.Equal(miner.X, crypto.Keccak256(epoch[:], miner.Commitment[:]))).Should(Equal(true))
			}
		})

		It("should pass the require X hashes check", func() {
			epoch, err := generateEpoch()
			Ω(err).ShouldNot(HaveOccurred())
			miners, err := generateMiners()
			Ω(err).ShouldNot(HaveOccurred())
			xing.AssignXHash(miners, epoch)
			Ω(xing.RequireXHashes(miners)).Should(Equal(true))
		})
	})

	Context("when calculating the number of classes", func() {
		It("should always be odd", func() {
			for n := 7; n < 1000; n++ {
				c := xing.NumberOfClasses(n)
				Ω(c%2 == 1).Should(Equal(true))
			}
		})
	})

	Context("when assigning classes", func() {
		It("should generate the correct classes", func() {
			numberOfMiners := 1000
			numberOfMNetworks := xing.NumberOfMNetworks(numberOfMiners)
			epoch, err := generateEpoch()
			Ω(err).ShouldNot(HaveOccurred())
			miners, err := generateMiners()
			Ω(err).ShouldNot(HaveOccurred())
			xing.AssignXHash(miners, epoch)
			xing.AssignClass(miners, numberOfMNetworks)
		})
		It("should sort the miners", func() {
			numberOfMiners := 1000
			numberOfMNetworks := xing.NumberOfMNetworks(numberOfMiners)
			epoch, err := generateEpoch()
			Ω(err).ShouldNot(HaveOccurred())
			miners, err := generateMiners()
			Ω(err).ShouldNot(HaveOccurred())
			xing.AssignXHash(miners, epoch)
			// Purposefully unsort the list.
			m := miners[0]
			miners[0] = miners[len(miners)-1]
			miners[len(miners)-1] = m
			// Calculate classes and ensure it is still sorted.
			xing.AssignClass(miners, numberOfMNetworks)
			Ω(sort.SliceIsSorted(miners, func(i, j int) bool {
				return miners[i].X.LessThan(miners[j].X)
			})).Should(Equal(true))
		})
	})

	Context("when assigning M networks", func() {
		It("should generate the correct M networks", func() {
			numberOfMiners := 1000
			numberOfMNetworks := xing.NumberOfMNetworks(numberOfMiners)
			epoch, err := generateEpoch()
			Ω(err).ShouldNot(HaveOccurred())
			miners, err := generateMiners()
			Ω(err).ShouldNot(HaveOccurred())
			xing.AssignXHash(miners, epoch)
			xing.AssignMNetwork(miners, numberOfMNetworks)
		})
		It("should sort the miners", func() {
			numberOfMiners := 1000
			numberOfMNetworks := xing.NumberOfMNetworks(numberOfMiners)
			epoch, err := generateEpoch()
			Ω(err).ShouldNot(HaveOccurred())
			miners, err := generateMiners()
			Ω(err).ShouldNot(HaveOccurred())
			xing.AssignXHash(miners, epoch)
			// Purposefully unsort the list.
			m := miners[0]
			miners[0] = miners[len(miners)-1]
			miners[len(miners)-1] = m
			// Calculate classes and ensure it is still sorted.
			xing.AssignMNetwork(miners, numberOfMNetworks)
			Ω(sort.SliceIsSorted(miners, func(i, j int) bool {
				return miners[i].X.LessThan(miners[j].X)
			})).Should(Equal(true))
		})
	})
})

func generateEpoch() (xing.Hash, error) {
	id, _, err := identity.NewID()
	if err != nil {
		return nil, err
	}
	return crypto.Keccak256([]byte(id.String())), nil
}

func generateMiners() ([]xing.Miner, error) {
	miners := make([]xing.Miner, 100)
	for i := 0; i < len(miners); i++ {
		id, _, err := identity.NewID()
		if err != nil {
			return nil, err
		}
		miners[i] = xing.NewMiner(id, crypto.Keccak256([]byte(id.String())))
	}
	return miners, nil
}
