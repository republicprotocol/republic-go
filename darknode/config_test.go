package darknode_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Configurations", func() {

	Context("negative tests", func() {
		It("should return an error when trying to open an non-existent file", func() {
			_, err := LoadConfig("non-existent.json")
			Ω(err).Should(HaveOccurred())
		})
	})

})
