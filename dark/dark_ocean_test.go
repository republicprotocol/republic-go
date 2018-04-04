package dark_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/contracts/dnr"
	"github.com/republicprotocol/republic-go/dark"
	"github.com/republicprotocol/republic-go/logger"
)

var _ = Describe("Dark Oceans", func() {
	Context("Watch", func() {
		It("should send a message to the channel", func() {
			log, err := logger.NewLogger(logger.Options{})
			if err != nil {
				panic(err)
			}
			dnr, err := dnr.NewMockDarkNodeRegistrar()

			ocean, err := dark.NewOcean(log, dnr)
			if err != nil {
				panic(err)
			}

			channel := make(chan struct{}, 1)
			go ocean.Watch(time.Second, channel)
			Eventually(channel).Should(Receive())

			dnr.Epoch()

			Ω(nil).Should(BeNil())
		})
	})
})
