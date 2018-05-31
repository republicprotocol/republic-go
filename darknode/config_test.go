package darknode_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/darknode"
)

var _ = Describe("Config", func() {

	Context("when unmarshalling config from file", func() {

		It("should not return error", func() {
			config, err := NewConfigFromJSONFile("../cmd/darknode/config.json")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(config).ShouldNot(BeNil())
		})

	})
})
