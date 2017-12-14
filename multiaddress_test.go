package identity_test

import (
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/republicprotocol/go-identity"
)

var _ = Describe("Address", func() {

	It("should be able to convert to Bas64", func() {
		id, err := NewPublicID()
		Ω(err).Should(BeNil())
		log.Println(id.Base64())
	})

})
