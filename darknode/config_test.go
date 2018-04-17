package darknode_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/darknode"
)

var _ = Describe("Configs", func() {

	Context("when loading malformed files", func() {

		It("should return an error when no file is found", func() {
			_, err := LoadConfig("notfound.json")
			Ω(err).Should(HaveOccurred())
		})

	})

})
