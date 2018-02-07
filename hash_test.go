package xing_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"math/rand"

	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-x"
)

var _ = Describe("Hashes", func() {

	Context("when comparing hashes", func() {
		It("should get different result if you swap the variable", func() {
			hash1, hash2 := randomHash(), randomHash()
			Ω(hash1.LessThan(hash2)).ShouldNot(Equal(hash2.LessThan(hash1)))
		})
		It("should produce the correct result regarding two hashes", func() {
			hash := randomHash()
			shouldBeGreaterHash, shouldBeLessHash := xing.Hash(make([]byte, 32)), xing.Hash(make([]byte, 32))
			copy(shouldBeGreaterHash, hash)
			copy(shouldBeLessHash, hash)

			randomIndex := rand.Intn(32)
			for {
				if shouldBeGreaterHash[randomIndex] == uint8(255) || shouldBeLessHash[randomIndex] == uint8(0) {
					randomIndex = rand.Intn(32)
				} else {
					break
				}
			}

			shouldBeLessHash[randomIndex] -= 1
			shouldBeGreaterHash[randomIndex] += 1
			Ω(shouldBeLessHash.LessThan(hash)).Should(BeTrue())
			Ω(shouldBeGreaterHash.LessThan(hash)).Should(BeFalse())
		})
		It("should not less than itself", func() {
			hash := randomHash()
			Ω(hash.LessThan(hash)).Should(BeFalse())
		})
	})

	Context("Creating new miner", func() {
		It("should create a new miner with the given id and hash", func() {
			id, _, err := identity.NewID()
			Ω(err).ShouldNot(HaveOccurred())
			commitment := make([]byte, 32)
			for i, _ := range commitment {
				commitment[i] = uint8(rand.Intn(256))
			}
			miner := xing.NewMiner(id, commitment)
			Ω(miner.ID).Should(Equal(id))
			for i, _ := range miner.Commitment {
				Ω(miner.Commitment[i]).Should(Equal(commitment[i]))
			}
		})
	})
})

func randomHash() xing.Hash {
	hash := make([]byte, 32)
	for i, _ := range hash {
		hash[i] = uint8(rand.Intn(256))
	}
	return xing.Hash(hash)

}
