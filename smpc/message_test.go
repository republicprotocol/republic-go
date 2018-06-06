package smpc_test

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/smpc"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/republicprotocol/republic-go/crypto"
)

var numberOfMessages = 24

var _ = Describe("Messages", func() {

	Context("Marshal and unmarshal Message", func() {

		It("should equal itself after marshaling and unmarshaling in binary if the message is of type MessageTypeJoin", func() {
			messageJoins := generateMessageJoin(numberOfMessages)
			messages := make([]Message, len(messageJoins))
			for i := range messages {
				messages[i] = Message{
					MessageType:         MessageTypeJoin,
					MessageJoin:         &messageJoins[i],
					MessageJoinResponse: nil,
				}
			}

			for i := range messages {
				data, err := messages[i].MarshalBinary()
				Ω(err).ShouldNot(HaveOccurred())

				var message Message
				Ω(message.UnmarshalBinary(data)).ShouldNot(HaveOccurred())
				Ω(message.MessageType).Should(Equal(MessageTypeJoin))
				Ω(message.MessageJoinResponse).Should(BeNil())
				Ω(bytes.Compare(messages[i].MessageJoin.NetworkID[:], message.MessageJoin.NetworkID[:])).Should(Equal(0))
				Ω(bytes.Compare(messages[i].MessageJoin.Join.ID[:], message.MessageJoin.Join.ID[:])).Should(Equal(0))
				Ω(messages[i].MessageJoin.Join.Index).Should(Equal(message.MessageJoin.Join.Index))
				Ω(len(messages[i].MessageJoin.Join.Shares)).Should(Equal(len(message.MessageJoin.Join.Shares)))
				for j := range messages[i].MessageJoin.Join.Shares {
					Ω(messages[i].MessageJoin.Join.Shares[j].Equal(&message.MessageJoin.Join.Shares[j]))
				}
			}
		})

		It("should equal itself after marshaling and unmarshaling in binary if the message is of type MessageJoinResponse", func() {
			messageJoinResponses := generateMessageJoinResponse(numberOfMessages)
			messages := make([]Message, len(messageJoinResponses))
			for i := range messages {
				messages[i] = Message{
					MessageType:         MessageTypeJoinResponse,
					MessageJoin:         nil,
					MessageJoinResponse: &messageJoinResponses[i],
				}
			}

			for i := range messages {
				data, err := messages[i].MarshalBinary()
				Ω(err).ShouldNot(HaveOccurred())

				var message Message
				Ω(message.UnmarshalBinary(data)).ShouldNot(HaveOccurred())
				Ω(message.MessageType).Should(Equal(MessageTypeJoinResponse))
				Ω(message.MessageJoin).Should(BeNil())
				Ω(bytes.Compare(messages[i].MessageJoinResponse.NetworkID[:], message.MessageJoinResponse.NetworkID[:])).Should(Equal(0))
				Ω(bytes.Compare(messages[i].MessageJoinResponse.Join.ID[:], message.MessageJoinResponse.Join.ID[:])).Should(Equal(0))
				Ω(messages[i].MessageJoinResponse.Join.Index).Should(Equal(message.MessageJoinResponse.Join.Index))
				Ω(len(messages[i].MessageJoinResponse.Join.Shares)).Should(Equal(len(message.MessageJoinResponse.Join.Shares)))
				for j := range messages[i].MessageJoinResponse.Join.Shares {
					Ω(messages[i].MessageJoinResponse.Join.Shares[j].Equal(&message.MessageJoinResponse.Join.Shares[j]))
				}
			}
		})
	})

	Context("Marshal and unmarshal MessageJoin", func() {

		It("should equal itself after marshaling and unmarshaling in binary", func() {
			messages := generateMessageJoin(numberOfMessages)

			for i := range messages {
				data, err := messages[i].MarshalBinary()
				Ω(err).ShouldNot(HaveOccurred())

				var message MessageJoin
				Ω(message.UnmarshalBinary(data)).ShouldNot(HaveOccurred())

				Ω(bytes.Compare(messages[i].NetworkID[:], message.NetworkID[:])).Should(Equal(0))
				Ω(bytes.Compare(messages[i].Join.ID[:], message.Join.ID[:])).Should(Equal(0))
				Ω(messages[i].Join.Index).Should(Equal(message.Join.Index))
				Ω(len(messages[i].Join.Shares)).Should(Equal(len(message.Join.Shares)))
				for j := range messages[i].Join.Shares {
					Ω(messages[i].Join.Shares[j].Equal(&message.Join.Shares[j]))
				}
			}
		})
	})

	Context("Marshal and unmarshal MessageJoinResponse", func() {

		It("should equal itself after marshaling and unmarshaling in binary", func() {
			messages := generateMessageJoinResponse(numberOfMessages)

			for i := range messages {
				data, err := messages[i].MarshalBinary()
				Ω(err).ShouldNot(HaveOccurred())

				var message MessageJoinResponse
				Ω(message.UnmarshalBinary(data)).ShouldNot(HaveOccurred())

				Ω(bytes.Compare(messages[i].NetworkID[:], message.NetworkID[:])).Should(Equal(0))
				Ω(bytes.Compare(messages[i].Join.ID[:], message.Join.ID[:])).Should(Equal(0))
				Ω(messages[i].Join.Index).Should(Equal(message.Join.Index))
				Ω(len(messages[i].Join.Shares)).Should(Equal(len(message.Join.Shares)))
				for j := range messages[i].Join.Shares {
					Ω(messages[i].Join.Shares[j].Equal(&message.Join.Shares[j]))
				}
			}
		})
	})

})

func generateMessageJoin(k int) []MessageJoin {
	messages := make([]MessageJoin, k)
	_, joins := generateJoins()
	var networkID [32]byte
	copy(networkID[:], crypto.Keccak256([]byte{uint8(math.MaxUint8)}))
	for i := range messages {
		messages[i] = MessageJoin{
			NetworkID: networkID,
			Join:      joins[i],
		}
	}

	return messages
}

func generateMessageJoinResponse(k int) []MessageJoinResponse {
	messages := make([]MessageJoinResponse, k)
	_, joins := generateJoins()
	var networkID [32]byte
	copy(networkID[:], crypto.Keccak256([]byte{uint8(math.MaxUint8)}))
	for i := range messages {
		messages[i] = MessageJoinResponse{
			NetworkID: networkID,
			Join:      joins[i],
		}
	}

	return messages
}
