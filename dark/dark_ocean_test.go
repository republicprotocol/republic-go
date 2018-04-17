package dark_test

import (
	"context"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/contracts/connection"
	"github.com/republicprotocol/republic-go/contracts/dnr"
	"github.com/republicprotocol/republic-go/dark"
	"github.com/republicprotocol/republic-go/logger"
)

var _ = Describe("Dark Oceans", func() {
	Context("testrpc", func() {
		It("should send a message to the channel", func() {
			log, err := logger.NewLogger(logger.Options{})
			Ω(err).ShouldNot(HaveOccurred())

			dnr, err := dnr.TestnetDNR(nil)
			Ω(err).ShouldNot(HaveOccurred())

			ocean, err := dark.NewOcean(log, 5, dnr)
			Ω(err).ShouldNot(HaveOccurred())

			channel := make(chan struct{}, 1)
			go ocean.Watch(channel)
			Eventually(channel).Should(Receive())

			Eventually(channel).Should(Receive())
			dnr.WaitForEpoch()

			Ω(nil).Should(BeNil())
		})
	})

	Context("ropsten", func() {
		const key = `{"version":3,"id":"7844982f-abe7-4690-8c15-34f75f847c66","address":"db205ea9d35d8c01652263d58351af75cfbcbf07","Crypto":{"ciphertext":"378dce3c1279b36b071e1c7e2540ac1271581bff0bbe36b94f919cb73c491d3a","cipherparams":{"iv":"2eb92da55cc2aa62b7ffddba891f5d35"},"cipher":"aes-128-ctr","kdf":"scrypt","kdfparams":{"dklen":32,"salt":"80d3341678f83a14024ba9c3edab072e6bd2eea6aa0fbc9e0a33bae27ffa3d6d","n":8192,"r":8,"p":1},"mac":"3d07502ea6cd6b96a508138d8b8cd2e46c3966240ff276ce288059ba4235cb0d"}}`

		It("should send a message to the channel", func() {
			mockLogger, err := logger.NewLogger(logger.Options{})
			Ω(err).ShouldNot(HaveOccurred())

			auth, err := bind.NewTransactor(strings.NewReader(key), "password1")
			Ω(err).ShouldNot(HaveOccurred())

			client, err := connection.FromURI("https://ropsten.infura.io/", "ropsten")
			Ω(err).ShouldNot(HaveOccurred())

			dnr, err := dnr.NewDarkNodeRegistry(context.Background(), &client, auth, &bind.CallOpts{})
			Ω(err).ShouldNot(HaveOccurred())

			ocean, err := dark.NewOcean(mockLogger, 5, dnr)
			Ω(err).ShouldNot(HaveOccurred())

			channel := make(chan struct{}, 1)
			go ocean.Watch(channel)
			Eventually(channel, 2*time.Second).Should(Receive())

			// Would have to wait for an epoch, will slow test down too much
			// Eventually(channel).Should(Receive())
			// dnr.WaitForEpoch()

			Ω(nil).Should(BeNil())
		})
	})
})
