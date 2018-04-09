package dark_test

import (
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
			dnr, err := dnr.TestnetDNR(nil)
			if err != nil {
				panic(err)
			}

			ocean, err := dark.NewOcean(log, 5, dnr)
			if err != nil {
				panic(err)
			}

			channel := make(chan struct{}, 1)
			go ocean.Watch(channel)
			Eventually(channel).Should(Receive())

			dnr.WaitForEpoch()

			Ω(nil).Should(BeNil())
		})
	})
})
