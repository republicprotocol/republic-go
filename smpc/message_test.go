package smpc_test

import (
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/smpc"

	"github.com/republicprotocol/republic-go/shamir"
)

var _ = Describe("Messages", func() {

	Context("when marshaling and unmarshaling messages", func() {

		It("should equal itself after marshaling and unmarshaling in binary", func() {
			for i := uint64(0); i < 100; i++ {
				messageJ := MessageJ{
					Share: shamir.Share{
						Index: uint64(rand.Int63()),
						Value: uint64(rand.Int63()),
					},
				}
				message := Message{
					MessageType: MessageTypeJ,
					MessageJ:    &messageJ,
				}
				_, err := rand.Read(messageJ.InstID[:])
				Expect(err).ShouldNot(HaveOccurred())
				_, err = rand.Read(messageJ.NetworkID[:])
				Expect(err).ShouldNot(HaveOccurred())

				data, err := message.MarshalBinary()
				Expect(err).ShouldNot(HaveOccurred())
				unmarshaledMessage := Message{MessageJ: &MessageJ{}}
				err = unmarshaledMessage.UnmarshalBinary(data)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(message).To(Equal(unmarshaledMessage))
			}
		})

		It("should return an error when marshaling an unknown message", func() {
			message := Message{
				MessageType: -69,
			}
			_, err := message.MarshalBinary()
			Expect(err).Should(Equal(ErrUnexpectedMessageType))
		})

		It("should return an error when marshaling an unknown message", func() {
			message := Message{}
			err := message.UnmarshalBinary([]byte{byte(255)})
			Expect(err).Should(Equal(ErrUnexpectedMessageType))
		})

		It("should return an error when unmarshaling an empty data as binary", func() {
			message := Message{}
			err := message.UnmarshalBinary([]byte{})
			Expect(err).Should(HaveOccurred())
		})

	})

	Context("when marshaling and unmarshaling messageJs", func() {

		It("should equal itself after marshaling and unmarshaling in binary", func() {
			for i := uint64(0); i < 100; i++ {
				messageJ := MessageJ{
					Share: shamir.Share{
						Index: uint64(rand.Int63()),
						Value: uint64(rand.Int63()),
					},
				}
				_, err := rand.Read(messageJ.InstID[:])
				Expect(err).ShouldNot(HaveOccurred())
				_, err = rand.Read(messageJ.NetworkID[:])
				Expect(err).ShouldNot(HaveOccurred())

				data, err := messageJ.MarshalBinary()
				Expect(err).ShouldNot(HaveOccurred())
				unmarshaledMessageJ := MessageJ{}
				err = unmarshaledMessageJ.UnmarshalBinary(data)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(messageJ).To(Equal(unmarshaledMessageJ))
			}
		})

		It("should return an error when unmarshaling an empty data as binary", func() {
			message := Message{}
			err := message.UnmarshalBinary([]byte{})
			Expect(err).Should(HaveOccurred())
		})

	})

})
