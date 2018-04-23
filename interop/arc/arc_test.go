package arc_test

import (
	"crypto/sha256"
	"math/big"
	"time"

	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"
	// . "github.com/republicprotocol/republic-go/interop/arc"
)

var _ = Describe("Arc", func() {
	Context("BTC-ETH swap", func() {
		It("can be fulfilled by both parties", func() {
			secret := randomBytes32()
			hashLock := sha256.Sum256(secret)
			BTCAtom := arc_bitcoin.NewBitcoinArc(RPC_USERNAME, RPC_PASSWORD, CHAIN)
			err := BTCAtom.Initiate(hashLock, []byte(aliceAddr), []byte(bobAddr), big.NewInt(3000000), time.Now().Unix()+10000)
			Ω(err).Should(BeNil())
			err = BTCAtom.Redeem(secret)
			Ω(err).Should(BeNil())
		})
	})
})
