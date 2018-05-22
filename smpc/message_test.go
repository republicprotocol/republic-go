package smpc_test

import (
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/smpc"

	"github.com/republicprotocol/republic-go/shamir"
)

var _ = Describe("smpc messages ", func() {

	Context("when marshaling and unmarshaling standard messages", func() {

		It("should equal itself after marshaling and unmarshaling in binary", func() {
			for i := uint64(0); i < 100; i++ {
				messageJ := MessageJ{
					InstID:    [32]byte{byte(i)},
					NetworkID: [32]byte{byte(i)},
					Share: shamir.Share{
						Index: uint64(rand.Int63()),
						Value: uint64(rand.Int63()),
					},
				}
				message := Message{
					MessageType: 1,
					MessageJ:    &messageJ,
				}
				data, err := message.MarshalBinary()
				Expect(err).ShouldNot(HaveOccurred())
				unmarshaledMessage := Message{MessageJ: &MessageJ{}}
				err = unmarshaledMessage.UnmarshalBinary(data)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(message).To(Equal(unmarshaledMessage))
			}
		})

		It("should return an error when unmarshaling an empty data as binary", func() {
			message := Message{}
			err := message.UnmarshalBinary([]byte{})
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("when marshaling and unmarshaling join messages to the SMPC", func() {

		It("should equal itself after marshaling and unmarshaling in binary", func() {
			for i := uint64(0); i < 100; i++ {
				messageJ := MessageJ{
					InstID:    [32]byte{byte(i)},
					NetworkID: [32]byte{byte(i)},
					Share: shamir.Share{
						Index: uint64(rand.Int63()),
						Value: uint64(rand.Int63()),
					},
				}
				data, err := messageJ.MarshalBinary()
				Expect(err).ShouldNot(HaveOccurred())
				unmarshaledMessage := MessageJ{}
				err = unmarshaledMessage.UnmarshalBinary(data)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(messageJ).To(Equal(unmarshaledMessage))
			}
		})

		It("should return an error when unmarshaling an empty data as binary", func() {
			message := MessageJ{}
			err := message.UnmarshalBinary([]byte{})
			Expect(err).Should(HaveOccurred())
		})
	})
})
