package stream_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/stream"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
)

var _ = Describe("streamer channel", func() {

	channelHub := ChannelHub{}
	clientMultiAddr := identity.MultiAddress{}
	serverMultiAddr := identity.MultiAddress{}

	BeforeEach(func() {
		var err error

		// Create a new channel hub
		channelHub = NewChannelHub()

		//Create client and server multiaddresses
		clientMultiAddr, err = createNewMultiAddress()
		Expect(err).ShouldNot(HaveOccurred())
		serverMultiAddr, err = createNewMultiAddress()
		Expect(err).ShouldNot(HaveOccurred())
	})

	Context("on registering a client - server pair with a channel hub", func() {

		It("should register a new stream with channel hub", func() {
			// Create a new channel client with the client address
			channelClient := NewChannelClient(clientMultiAddr.Address(), &channelHub)

			//Connect client to the server
			ctx, cancelCtx := context.WithCancel(context.Background())
			defer cancelCtx()
			closeStream, err := channelClient.Connect(ctx, serverMultiAddr)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(closeStream).ToNot(BeNil())
		})

		It("should send and receive messages", func() {
			// Create a new channel client with the client address and connect it to a server
			channelClient := NewChannelClient(clientMultiAddr.Address(), &channelHub)
			ctx, cancelCtx := context.WithCancel(context.Background())
			defer cancelCtx()
			clientStream, err := channelClient.Connect(ctx, serverMultiAddr)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(clientStream).ToNot(BeNil())

			// Create a new server channel and start listening to client
			channelServer := NewChannelServer(serverMultiAddr.Address(), &channelHub)
			serverStream, err := channelServer.Listen(ctx, clientMultiAddr.Address())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(serverStream).ToNot(BeNil())

			// Send 5 messages to client from server
			var receivedMessage, message *mockByteMessage
			
			go func() {
				for i := 0; i < 5; i++ {
					message = &mockByteMessage{
						value: []byte{byte(i)},
					}
					err = serverStream.Send(message)
					Expect(err).ShouldNot(HaveOccurred())
				}
			}()
			for i := 0; i < 5; i++ {
				receivedMessage = &mockByteMessage{}
				err = clientStream.Recv(receivedMessage)
				Expect(err).ShouldNot(HaveOccurred())
			}
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
	ecdsaKey, err := crypto.RandomEcdsaKey()
	if err != nil {
		return identity.MultiAddress{}, err
	}
	return identity.Address(ecdsaKey.Address()).MultiAddress()
}
