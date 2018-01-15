package x_test

import (
	"bytes"

	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"

	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-x"
)

var _ = Describe("X", func() {

	Context("when assigning the X overlay", func() {

		var miners []x.Miner
		var overlayMiners []x.Miner

		BeforeEach(func() {
			epoch, err := generateEpoch()
			Ω(err).ShouldNot(HaveOccurred())

			// Generate overlay miners.
			overlayMiners, err = generateMiners()
			Ω(err).ShouldNot(HaveOccurred())

			// Compute the X Overlay.
			numberOfMNetworks := x.NumberOfMNetworks(len(overlayMiners))
			x.AssignXOverlay(overlayMiners, epoch, numberOfMNetworks)

			// Clone miners into a clean slice.
			miners = make([]x.Miner, len(overlayMiners))
			for i := range overlayMiners {
				miners[i] = x.NewMiner(overlayMiners[i].ID, overlayMiners[i].Commitment)
			}

			// Compute the individual components of the Miners.
			x.AssignXHash(miners, epoch)
			x.AssignClass(miners, numberOfMNetworks)
			x.AssignMNetwork(miners, numberOfMNetworks)
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
			x.AssignXHash(miners, epoch)
			for _, miner := range miners {
				Ω(bytes.Equal(miner.X, crypto.Keccak256(epoch[:], miner.Commitment[:]))).Should(Equal(true))
			}
		})

		It("should pass the require X hashes check", func() {
			epoch, err := generateEpoch()
			Ω(err).ShouldNot(HaveOccurred())
			miners, err := generateMiners()
			Ω(err).ShouldNot(HaveOccurred())
			x.AssignXHash(miners, epoch)
			Ω(x.RequireXHashes(miners)).Should(Equal(true))
		})
	})

	Context("when calculating the number of classes", func() {
		It("should always be odd", func() {
			for n := 7; n < 1000; n++ {
				c := x.NumberOfClasses(n)
				Ω(c%2 == 1).Should(Equal(true))
			}
		})
	})

	Context("when assigning classes", func() {
		It("should generate the correct classes", func() {
			numberOfMiners := 1000
			numberOfMNetworks := x.NumberOfMNetworks(numberOfMiners)
			epoch, err := generateEpoch()
			Ω(err).ShouldNot(HaveOccurred())
			miners, err := generateMiners()
			Ω(err).ShouldNot(HaveOccurred())
			x.AssignXHash(miners, epoch)
			x.AssignClass(miners, numberOfMNetworks)
		})
	})

	Context("when assigning M networks", func() {
		It("should generate the correct M networks", func() {
			numberOfMiners := 1000
			numberOfMNetworks := x.NumberOfMNetworks(numberOfMiners)
			epoch, err := generateEpoch()
			Ω(err).ShouldNot(HaveOccurred())
			miners, err := generateMiners()
			Ω(err).ShouldNot(HaveOccurred())
			x.AssignXHash(miners, epoch)
			x.AssignMNetwork(miners, numberOfMNetworks)
		})
	})
})

func generateEpoch() (x.Hash, error) {
	id, _, err := identity.NewID()
	if err != nil {
		return nil, err
	}
	return crypto.Keccak256([]byte(id.String())), nil
}

func generateMiners() ([]x.Miner, error) {
	miners := make([]x.Miner, 100)
	for i := 0; i < len(miners); i++ {
		id, _, err := identity.NewID()
		if err != nil {
			return nil, err
		}
		miners[i] = x.NewMiner(id, crypto.Keccak256([]byte(id.String())))
	}
	return miners, nil
}
