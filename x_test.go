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
				Ω(c <= n).Should(Equal(true))
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
