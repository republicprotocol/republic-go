package order_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/order"

	"github.com/republicprotocol/republic-go/shamir"
)

var _ = Describe("Value", func() {

	Context("when marshaling and unmarshaling", func() {

		It("should be able to marshal and unmarshal Value as JSON", func() {
			value := createValue(uint32(1), uint32(2))

			data, err := value.MarshalJSON()
			Expect(err).ShouldNot(HaveOccurred())

			decodedValue := Value{}
			err = decodedValue.UnmarshalJSON(data)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(decodedValue.Exp).Should(Equal(value.Exp))
			Expect(decodedValue.Co).Should(Equal(value.Co))
		})

		It("should be able to marshal and unmarshal FragmentValue as JSON", func() {
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

			decodedValue := FragmentValue{}
			err = decodedValue.UnmarshalJSON(data)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(decodedValue.Exp).Should(Equal(value.Exp))
			Expect(decodedValue.Co).Should(Equal(value.Co))
		})

		It("should be able to marshal and unmarshal EncryptedFragmentValue as JSON", func() {
			encryptedFragmentValue := EncryptedFragmentValue{
				Co:  []byte{1},
				Exp: []byte{2},
			}

			data, err := encryptedFragmentValue.MarshalJSON()
			Expect(err).ShouldNot(HaveOccurred())

			decodedValue := EncryptedFragmentValue{}
			err = decodedValue.UnmarshalJSON(data)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(decodedValue.Exp).Should(Equal(encryptedFragmentValue.Exp))
			Expect(decodedValue.Co).Should(Equal(encryptedFragmentValue.Co))
		})

		Context("when marshalling using bigEndian encoding", func() {

			It("should not return an error when Value is marshalled", func() {
				value := createValue(uint32(1), uint32(2))

				data, err := value.MarshalBinary()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(data).ShouldNot(BeNil())
			})

			It("should not return an error when FragmentValue is marshalled", func() {
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

			It("should not return an error when EncryptedFragmentValue is marshalled", func() {
				encryptedFragmentValue := EncryptedFragmentValue{
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
			value1 := FragmentValue{}
			value2 := FragmentValue{}

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
			encryptedFragmentValue1 := EncryptedFragmentValue{
				Co:  []byte{1},
				Exp: []byte{2},
			}
			encryptedFragmentValue2 := EncryptedFragmentValue{
				Co:  []byte{1},
				Exp: []byte{2},
			}

			Expect(encryptedFragmentValue1.Equal(&encryptedFragmentValue2)).Should(Equal(true))
		})

		It("should return false for different encrypted fragment values", func() {
			encryptedFragmentValue1 := EncryptedFragmentValue{
				Co:  []byte{2},
				Exp: []byte{1},
			}
			encryptedFragmentValue2 := EncryptedFragmentValue{
				Co:  []byte{1},
				Exp: []byte{2},
			}

			Expect(encryptedFragmentValue1.Equal(&encryptedFragmentValue2)).Should(Equal(false))
		})
	})
})

func createFragmentValue(coeff, exp shamir.Share) FragmentValue {
	return FragmentValue{
		Co:  coeff,
		Exp: exp,
	}
}
