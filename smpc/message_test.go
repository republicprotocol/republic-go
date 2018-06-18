package smpc_test

import (
	"bytes"
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/smpc"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/republicprotocol/republic-go/crypto"
)

var (
	n = int64(24)
	k = 2 * (n + 1) / 3
)

var _ = Describe("Messages", func() {

	Context("when marshaling and unmarshaling message of type MessageTypeJoin", func() {
		var messageJoins []MessageJoin
		var messages []Message

		BeforeEach(func() {
			messageJoins = generateMessageJoin(n, k)
			messages = make([]Message, len(messageJoins))
			for i := range messages {
				messages[i] = Message{
					MessageType:         MessageTypeJoin,
					MessageJoin:         &messageJoins[i],
					MessageJoinResponse: nil,
				}
			}
		})

		It("should equal itself after marshaling and unmarshaling to binary", func() {
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

		It("should error if the messageType is wrong", func() {
			for i := range messages {
				messages[i].MessageType = MessageType(3)
				_, err := messages[i].MarshalBinary()
				log.Println(err)
				Ω(err).Should(HaveOccurred())
			}
		})

		It("should implements the stream.Message interface", func() {
			for i := range messages {
				messages[i].IsMessage()
			}
		})
	})

	Context("when marshaling and unmarshaling message of type MessageJoinResponse", func() {
		var messageJoinResponses []MessageJoinResponse
		var messages []Message

		BeforeEach(func() {
			messageJoinResponses = generateMessageJoinResponse(n, k)
			messages := make([]Message, len(messageJoinResponses))
			for i := range messages {
				messages[i] = Message{
					MessageType:         MessageTypeJoinResponse,
					MessageJoin:         nil,
					MessageJoinResponse: &messageJoinResponses[i],
				}
			}
		})

		It("should equal itself after marshaling and unmarshaling to binary", func() {
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

		It("should error if the messageType is wrong", func() {
			for i := range messages {
				messages[i].MessageType = MessageType(3)
				_, err := messages[i].MarshalBinary()
				log.Println(err)
				Ω(err).Should(HaveOccurred())
			}
		})

		It("should implements the stream.Message interface", func() {
			for i := range messages {
				messages[i].IsMessage()
			}
		})
	})
})

func generateMessageJoin(n, k int64) []MessageJoin {
	messages := make([]MessageJoin, n)
	_, joins := generateJoins(n, k)
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

func generateMessageJoinResponse(n, k int64) []MessageJoinResponse {
	messages := make([]MessageJoinResponse, n)
	_, joins := generateJoins(n, k)
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
