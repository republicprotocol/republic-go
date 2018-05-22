package stream_test

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/stream"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
)

var _ = Describe("streamer channel", func() {

	Context("on registering a client - server pair with a channel hub", func() {

		It("should register a new stream with channel hub", func() {
			channelHub := NewChannelHub()
			clientMultiAddr, err := createNewMultiAddress()
			Expect(err).ShouldNot(HaveOccurred())
			serverMultiAddr, err := createNewMultiAddress()
			Expect(err).ShouldNot(HaveOccurred())
			channelClient := NewChannelClient(clientMultiAddr.Address(), &channelHub)
			ctx, cancelCtx := context.WithCancel(context.Background())
			defer cancelCtx()
			closeStream, err := channelClient.Connect(ctx, serverMultiAddr)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(closeStream).ToNot(BeNil())
		})

		It("should send and receive messages", func() {
			channelHub := NewChannelHub()
			clientMultiAddr, err := createNewMultiAddress()
			Expect(err).ShouldNot(HaveOccurred())
			serverMultiAddr, err := createNewMultiAddress()
			Expect(err).ShouldNot(HaveOccurred())
			channelClient := NewChannelClient(clientMultiAddr.Address(), &channelHub)
			ctx, cancelCtx := context.WithCancel(context.Background())
			defer cancelCtx()
			closeStream, err := channelClient.Connect(ctx, serverMultiAddr)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(closeStream).ToNot(BeNil())
			channelServer := NewChannelServer(serverMultiAddr.Address(), &channelHub)
			newCloseStream, err := channelServer.Listen(ctx, clientMultiAddr.Address())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newCloseStream).ToNot(BeNil())
			var receivedMessage, message *mockByteMessage
			go func() {
				receivedMessage = &mockByteMessage{}
				err = closeStream.Recv(receivedMessage)
				Expect(err).ShouldNot(HaveOccurred())
			}()
			go func() {
				message = &mockByteMessage{
					value: []byte{2},
				}
				err = newCloseStream.Send(message)
				Expect(err).ShouldNot(HaveOccurred())
			}()
			Expect(receivedMessage).To(Equal(message))
		})
	})
})

type mockByteMessage struct {
	value []byte
}

func (message *mockByteMessage) IsMessage() {
	return
}

func (message *mockByteMessage) MarshalBinary() (data []byte, err error) {
	return message.value, nil
}

func (message *mockByteMessage) UnmarshalBinary(data []byte) error {
	message.value = data
	return nil
}

func createNewMultiAddress() (identity.MultiAddress, error) {
	keystore, err := crypto.RandomKeystore()
	if err != nil {
		return identity.MultiAddress{}, err
	}
	return identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/80/republic/%s", keystore.Address()))
}
