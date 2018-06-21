package contract_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/stackint"
	"github.com/republicprotocol/republic-go/testutils"
)

var _ = Describe("Contract Binder", func() {

	binder, _ := testutils.GanacheBeforeSuite(func() {

	})

	testutils.GanacheAfterSuite(func() {

	})
	Context("when interacting with ethereum smart contracts", func() {

		It("should not return a nonce error", func() {

			// Trigger new Epoch
			epoch, err := binder.NextEpoch()
			Expect(err).ShouldNot(HaveOccurred())

			// Bond for the darknode
			bond, err := stackint.FromString("0")
			Expect(err).ShouldNot(HaveOccurred())

			// Darknode address
			darknodeAddr := identity.Address("8MK6bwP1ADVPaMQ4Gxfm85KYbEdJ6Y")

			// Keystore
			keystore, err := crypto.RandomRsaKey()
			Expect(err).ShouldNot(HaveOccurred())
			publicKey, err := crypto.BytesFromRsaPublicKey(&(keystore.PublicKey))

			// Register darknode with the darknodeRegistry
			_, err = binder.Register(darknodeAddr.ID(), publicKey, &bond)
			Expect(err).ShouldNot(HaveOccurred())

			time.Sleep(14 * time.Second)

			// Check if darknode has registered
			tx, err := binder.IsDeregistered(darknodeAddr.ID())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(tx).To(BeTrue())

		})
	})
})
