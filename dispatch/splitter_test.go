package dispatch_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/dispatch"
)

var _ = Describe("splitter", func() {
	Context("run message queue", func() {
		splitter := dispatch.NewSplitter()
		Ω(true).Should(BeTrue())
	})
})
