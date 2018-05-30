package stream_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/dispatch"
	. "github.com/republicprotocol/republic-go/stream"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
)

var _ = Describe("Channel streams", func() {

	var hub ChannelHub

	BeforeEach(func() {
		hub = NewChannelHub()
	})

	Context("when sending and receiving messages", func() {

		var clientMultiAddr identity.MultiAddress
		var client Streamer
		var clientStream Stream
		var clientCancel context.CancelFunc

		var serverMultiAddr identity.MultiAddress
		var server Streamer
		var serverStream Stream
		var serverCancel context.CancelFunc

		BeforeEach(func() {
			var err error

			clientMultiAddr, err = createNewMultiAddress()
			Expect(err).ShouldNot(HaveOccurred())
			serverMultiAddr, err = createNewMultiAddress()
			Expect(err).ShouldNot(HaveOccurred())

			client = NewChannelStreamer(clientMultiAddr.Address(), &hub)
			server = NewChannelStreamer(serverMultiAddr.Address(), &hub)

			ctx, cancel := context.WithCancel(context.Background())
			clientStream, err = client.Open(ctx, serverMultiAddr)
			clientCancel = cancel
			Expect(err).ShouldNot(HaveOccurred())

			ctx, cancel = context.WithCancel(context.Background())
			serverStream, err = server.Open(ctx, clientMultiAddr)
			serverCancel = cancel
			Expect(err).ShouldNot(HaveOccurred())
		})

		AfterEach(func() {
			clientCancel()
			serverCancel()
		})

		It("should receive messages sent by the client", func() {
			dispatch.CoBegin(func() {
				defer GinkgoRecover()
				for i := 0; i < 256; i++ {
					message := mockMessage([]byte{byte(i)})
					err := clientStream.Send(&message)
					Expect(err).ShouldNot(HaveOccurred())
				}
			}, func() {
				defer GinkgoRecover()
				for i := 0; i < 256; i++ {
					message := mockMessage{}
					err := serverStream.Recv(&message)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(message).Should(Equal(mockMessage([]byte{byte(i)})))
				}
			})
		})

		It("should receive messages sent by the server", func() {
			dispatch.CoBegin(func() {
				defer GinkgoRecover()
				for i := 0; i < 256; i++ {
					message := mockMessage{}
					err := clientStream.Recv(&message)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(message).Should(Equal(mockMessage([]byte{byte(i)})))
				}
			}, func() {
				defer GinkgoRecover()
				for i := 0; i < 256; i++ {
					message := mockMessage([]byte{byte(i)})
					err := serverStream.Send(&message)
					Expect(err).ShouldNot(HaveOccurred())
				}
			})
		})

	})
})

func createNewMultiAddress() (identity.MultiAddress, error) {
	ecdsaKey, err := crypto.RandomEcdsaKey()
	if err != nil {
		return identity.MultiAddress{}, err
	}
	return identity.Address(ecdsaKey.Address()).MultiAddress()
}
