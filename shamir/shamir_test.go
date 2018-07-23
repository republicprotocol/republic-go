package shamir_test

import (
	"encoding/json"
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/shamir"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/stackint"
)

var _ = Describe("Shamir's secret sharing", func() {

	Context("when checking equality", func() {

		It("should return true for equal shares", func() {
			for i := uint64(0); i < 100; i++ {
				index := uint64(rand.Int63())
				value := uint64(rand.Int63()) % Prime
				share := Share{
					Index: index,
					Value: value,
				}
				shareOther := Share{
					Index: index,
					Value: value,
				}
				Expect(share.Equal(&shareOther)).Should(BeTrue())
			}
		})

		It("should return false for unequal shares", func() {
			for i := uint64(0); i < 100; i++ {
				share := Share{
					Index: uint64(rand.Int63()),
					Value: uint64(rand.Int63()) % Prime,
				}
				shareOther := Share{
					Index: uint64(rand.Int63()),
					Value: uint64(rand.Int63()) % Prime,
				}
				Expect(share.Equal(&shareOther)).Should(BeFalse())
			}
		})

	})

	Context("when performing arithmetic", func() {

		It("should equal subtraction on the secrets when done on shares", func() {
			for i := uint64(0); i < 100; i++ {

				secret := ((uint64(rand.Int63()) % Prime) / 2) + (Prime / 2)
				secretOther := (uint64(rand.Int63()) % Prime) / 2

				shares, err := Split(72, 48, secret)
				Expect(err).ShouldNot(HaveOccurred())
				sharesOther, err := Split(72, 48, secretOther)
				Expect(err).ShouldNot(HaveOccurred())
				sharesResult := make(Shares, 72)
				for j := 0; j < 72; j++ {
					sharesResult[j] = shares[j].Sub(&sharesOther[j])
				}

				secretResult := Join(sharesResult)
				Expect(secretResult).Should(Equal(secret - secretOther))
			}
		})

	})

	Context("when marshaling and unmarshaling", func() {

		It("should equal itself after marshaling them unmarshaling in binary", func() {
			for i := uint64(0); i < 100; i++ {
				share := Share{
					Index: uint64(rand.Int63()),
					Value: uint64(rand.Int63()) % Prime,
				}
				data, err := share.MarshalBinary()
				Expect(err).ShouldNot(HaveOccurred())
				unmarshaledShare := Share{}
				err = unmarshaledShare.UnmarshalBinary(data)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(share.Index).Should(Equal(unmarshaledShare.Index))
				Expect(share.Value).Should(Equal(unmarshaledShare.Value))
			}
		})

		It("should equal itself after marshaling them unmarshaling in JSON", func() {
			for i := uint64(0); i < 100; i++ {
				share := Share{
					Index: uint64(rand.Int63()),
					Value: uint64(rand.Int63()) % Prime,
				}
				data, err := share.MarshalJSON()
				Expect(err).ShouldNot(HaveOccurred())
				unmarshaledShare := Share{}
				err = unmarshaledShare.UnmarshalJSON(data)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(share.Index).Should(Equal(unmarshaledShare.Index))
				Expect(share.Value).Should(Equal(unmarshaledShare.Value))
			}
		})

		It("should return an error when unmarshaling an empty data as binary", func() {
			share := Share{}
			err := share.UnmarshalBinary([]byte{})
			Expect(err).Should(HaveOccurred())
		})

		It("should return an error when unmarshaling an empty data as JSON", func() {
			share := Share{}
			err := share.UnmarshalJSON([]byte{})
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("when encrypting and decrypting", func() {

		It("should equal itself after an encryption then decryption", func() {
			rsaKey, err := crypto.RandomRsaKey()
			Expect(err).ShouldNot(HaveOccurred())
			for i := uint64(0); i < 100; i++ {
				share := Share{
					Index: uint64(rand.Int63()),
					Value: uint64(rand.Int63()) % Prime,
				}
				cipherText, err := share.Encrypt(rsaKey.PublicKey)
				Expect(err).ShouldNot(HaveOccurred())
				decryptedShare := Share{}
				err = decryptedShare.Decrypt(rsaKey.PrivateKey, cipherText)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(share.Index).Should(Equal(decryptedShare.Index))
				Expect(share.Value).Should(Equal(decryptedShare.Value))
			}
		})

		It("should return an error after encryption then decryption with different keys", func() {
			rsaKey, err := crypto.RandomRsaKey()
			Expect(err).ShouldNot(HaveOccurred())
			rsaKeyOther, err := crypto.RandomRsaKey()
			Expect(err).ShouldNot(HaveOccurred())
			for i := uint64(0); i < 100; i++ {
				share := Share{
					Index: uint64(rand.Int63()),
					Value: uint64(rand.Int63()) % Prime,
				}
				cipherText, err := share.Encrypt(rsaKey.PublicKey)
				Expect(err).ShouldNot(HaveOccurred())
				decryptedShare := Share{}
				err = decryptedShare.Decrypt(rsaKeyOther.PrivateKey, cipherText)
				Expect(err).Should(HaveOccurred())
			}
		})
	})

	Context("when splitting", func() {

		It("should return the required number of shares", func() {
			n := int64(100)
			k := int64(50)
			secret := uint64(1234)
			shares, err := Split(n, k, secret)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(int64(len(shares))).Should(Equal(n))
		})

		It("should return an error when k is greater than n", func() {
			n := int64(50)
			k := int64(100)
			secret := uint64(1234)
			_, err := Split(n, k, secret)
			Expect(err).Should(Equal(ErrNKError))
		})

		It("should return an error when the secret is greater than the prime", func() {
			n := int64(100)
			k := int64(50)
			secret := Prime + 1
			_, err := Split(n, k, secret)
			Expect(err).Should(Equal(ErrFiniteField))
		})

	})

	Context("when joining", func() {

		It("should rejoin shares unmarshalled by json", func() {
			// Shamir parameters.
			N := int64(8)
			K := int64(6)

			js := []byte(
				"[[0,0,0,0,0,0,0,1,92,83,98,143,101,106,148,77],[0,0,0,0,0,0,0,2,226,79,128,237,26,163,39,111],[0,0,0,0,0,0,0,3,176,18,63,8,204,25,116,185],[0,0,0,0,0,0,0,4,88,16,175,44,168,0,238,177],[0,0,0,0,0,0,0,5,19,183,191,129,190,183,148,255],[0,0,0,0,0,0,0,6,186,85,23,6,213,110,27,129],[0,0,0,0,0,0,0,7,27,135,163,115,37,13,223,52],[0,0,0,0,0,0,0,8,106,208,228,226,226,36,32,81]]",
			)
			secretStackInt := stackint.FromUint(40)

			var shares Shares
			err := json.Unmarshal(js, &shares)
			Expect(err).Should(BeNil())

			Expect(int64(len(shares))).Should(Equal(N))

			// For all K greater than, or equal to, 50 attempt to decode the secret.
			// Pick K unique indices in the range [0, k).
			indices := map[int]struct{}{}
			for i := 0; i < int(K); i++ {
				for {
					index := rand.Intn(int(K))
					if _, ok := indices[index]; !ok {
						indices[index] = struct{}{}
						break
					}
				}
			}
			// Use K shares to reconstruct the secret.
			kShares := make(Shares, K)
			for index := range indices {
				kShares[index] = shares[index]
			}
			decodedSecret := stackint.FromUint(uint(Join(kShares)))
			Expect(decodedSecret.Cmp(&secretStackInt)).Should(Equal(0))
		})

		It("should return the required secret from the threshold of shares", func() {
			// Shamir parameters.
			N := int64(100)
			K := int64(50)
			secret := uint64(1234)
			// Split the secret.
			shares, err := Split(N, K, secret)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(int64(len(shares))).Should(Equal(N))
			// For all K greater than, or equal to, 50 attempt to decode the secret.
			for k := int64(50); k < 101; k++ {
				// Pick K unique indices in the range [0, k).
				indices := map[int]struct{}{}
				for i := 0; i < int(k); i++ {
					for {
						index := rand.Intn(int(k))
						if _, ok := indices[index]; !ok {
							indices[index] = struct{}{}
							break
						}
					}
				}
				// Use K shares to reconstruct the secret.
				kShares := make(Shares, k)
				for index := range indices {
					kShares[index] = shares[index]
				}
				decodedSecret := stackint.FromUint(uint(Join(kShares)))
				secretStackInt := stackint.FromUint(uint(secret))
				Expect(decodedSecret.Cmp(&secretStackInt)).Should(Equal(0))
			}
		})

		It("should not return the required secret from less than the threshold of shares", func() {
			// Shamir parameters.
			N := int64(100)
			K := int64(50)
			secret := uint64(1234)
			// Split the secret.
			shares, err := Split(N, K, secret)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(int64(len(shares))).Should(Equal(N))
			// For all K less than 50 attempt to decode the secret.
			for k := int64(1); k < 50; k++ {
				// Pick K unique indices in the range [0, k).
				indices := map[int]struct{}{}
				for i := 0; i < int(k); i++ {
					for {
						index := rand.Intn(int(k))
						if _, ok := indices[index]; !ok {
							indices[index] = struct{}{}
							break
						}
					}
				}
				// Use K shares to reconstruct the secret.
				kShares := make(Shares, k)
				for index := range indices {
					kShares[index] = shares[index]
				}
				decodedSecret := Join(kShares)
				Expect(decodedSecret).Should(Not(Equal(&secret)))
			}
		})
	})
})
