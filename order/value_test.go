package order_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/order"

	"github.com/republicprotocol/republic-go/shamir"
)

var _ = Describe("Values", func() {

	Context("when marshaling and unmarshaling", func() {

		It("should be able to marshal and unmarshal CoExp as JSON", func() {
			value := NewCoExp(uint64(1), uint64(2))

			data, err := value.MarshalJSON()
			Expect(err).ShouldNot(HaveOccurred())

			decodedValue := CoExp{}
			err = decodedValue.UnmarshalJSON(data)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(decodedValue.Exp).Should(Equal(value.Exp))
			Expect(decodedValue.Co).Should(Equal(value.Co))
		})

		It("should return an error when unmarshalling invalid JSON", func() {
			decodedValue := CoExp{}
			err := decodedValue.UnmarshalJSON([]byte{byte(1)})
			Expect(err).Should(HaveOccurred())

			Expect(decodedValue.Exp).Should(Equal(uint64(0)))
		})

		It("should be able to marshal and unmarshal CoExpShare as JSON", func() {
			coeff := shamir.Share{
				Index: uint64(1),
				Value: uint64(5),
			}
			exp := shamir.Share{
				Index: uint64(10),
				Value: uint64(50),
			}
			value := createFragmentValue(coeff, exp)

			data, err := value.MarshalJSON()
			Expect(err).ShouldNot(HaveOccurred())

			decodedValue := CoExpShare{}
			err = decodedValue.UnmarshalJSON(data)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(decodedValue.Exp).Should(Equal(value.Exp))
			Expect(decodedValue.Co).Should(Equal(value.Co))
		})

		It("should be able to marshal and unmarshal EncryptedCoExp as JSON", func() {
			encryptedFragmentValue := EncryptedCoExpShare{
				Co:  []byte{1},
				Exp: []byte{2},
			}

			data, err := encryptedFragmentValue.MarshalJSON()
			Expect(err).ShouldNot(HaveOccurred())

			decodedValue := EncryptedCoExpShare{}
			err = decodedValue.UnmarshalJSON(data)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(decodedValue.Exp).Should(Equal(encryptedFragmentValue.Exp))
			Expect(decodedValue.Co).Should(Equal(encryptedFragmentValue.Co))
		})

		Context("when marshalling using big endian encoding", func() {

			It("should not return an error when CoExp is marshalled", func() {
				value := NewCoExp(uint64(1), uint64(2))

				data, err := value.MarshalBinary()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(data).ShouldNot(BeNil())
			})

			It("should not return an error when CoExpShare is marshalled", func() {
				coeff := shamir.Share{
					Index: uint64(1),
					Value: uint64(5),
				}
				exp := shamir.Share{
					Index: uint64(10),
					Value: uint64(50),
				}
				value := createFragmentValue(coeff, exp)

				data, err := value.MarshalBinary()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(data).ShouldNot(BeNil())
			})

			It("should not return an error when EncryptedCoExpShare is marshalled", func() {
				encryptedFragmentValue := EncryptedCoExpShare{
					Co:  []byte{1},
					Exp: []byte{2},
				}

				data, err := encryptedFragmentValue.MarshalBinary()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(data).ShouldNot(BeNil())
			})
		})
	})

	Context("when testing for equality", func() {

		It("should return true for similar fragment values", func() {
			value1 := CoExpShare{}
			value2 := CoExpShare{}

			Expect(value1.Equal(&value2)).Should(Equal(true))
		})

		It("should return false for different fragment values", func() {
			coeff := shamir.Share{
				Index: uint64(1),
				Value: uint64(5),
			}
			value := shamir.Share{
				Index: uint64(10),
				Value: uint64(50),
			}
			value1 := createFragmentValue(coeff, shamir.Share{})
			value2 := createFragmentValue(shamir.Share{}, value)

			Expect(value1.Equal(&value2)).Should(Equal(false))
		})

		It("should return true for similar encrypted fragment values", func() {
			encryptedFragmentValue1 := EncryptedCoExpShare{
				Co:  []byte{1},
				Exp: []byte{2},
			}
			encryptedFragmentValue2 := EncryptedCoExpShare{
				Co:  []byte{1},
				Exp: []byte{2},
			}

			Expect(encryptedFragmentValue1.Equal(&encryptedFragmentValue2)).Should(Equal(true))
		})

		It("should return false for different encrypted fragment values", func() {
			encryptedFragmentValue1 := EncryptedCoExpShare{
				Co:  []byte{2},
				Exp: []byte{1},
			}
			encryptedFragmentValue2 := EncryptedCoExpShare{
				Co:  []byte{1},
				Exp: []byte{2},
			}

			Expect(encryptedFragmentValue1.Equal(&encryptedFragmentValue2)).Should(Equal(false))
		})
	})
})

func createFragmentValue(coeff, exp shamir.Share) CoExpShare {
	return CoExpShare{
		Co:  coeff,
		Exp: exp,
	}
}
