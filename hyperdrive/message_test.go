package hyperdrive_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/hyperdrive"
)

var _ = Describe("Messages", func() {

	Context("when verifying and signing", func() {

		It("should return an error when signing returns an error", func(done Done) {
			defer close(done)

			prepare := Prepare{}
			messageStore := NewMessageMapStore()
			signer := NewErrorSigner()
			threshold := 0

			message, err := VerifyAndSignMessage(&prepare, &messageStore, &signer, threshold)
			Ω(message).Should(BeNil())
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("when verifying and signing", func() {
		It("should return an error when signing returns an error", func() {

			prepare := Prepare{}
			messageStore := NewMessageMapStore()
			signer := NewErrorSigner()
			threshold := 0

			message, err := VerifyAndSignMessage(&prepare, &messageStore, &signer, threshold)
			Ω(message).Should(BeNil())
			Ω(err).Should(HaveOccurred())
		})
	})



})
