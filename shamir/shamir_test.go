package shamir_test

import (
	"fmt"
	"math/big"
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/shamir"
	"github.com/republicprotocol/republic-go/stackint"
)

var _ = Describe("Shamir's secret sharing", func() {

	primeStr := "179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111"
	prime := stackint.FromString(primeStr)

	// It("should correctly encode integers greater than 2^1024", func() {
	// 	// The maximum 1024 bit integer.

	// 	max, ok := big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137215", 10)
	// 	Ω(ok).Should(Equal(true))
	// 	max.Sub(max, big.NewInt(1))
	// 	// The first prime above 1024 bits.
	// 	prime, ok := big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)
	// 	Ω(ok).Should(Equal(true))
	// 	Ω(prime.Cmp(max) > 0).Should(Equal(true))
	// })

	Context("serialization", func() {
		It("should be able to serialize and deserialize shares", func() {
			// prime, ok := big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)
			// Ω(ok).Should(Equal(true))

			for i := int64(0); i < 1000; i++ {
				bytes := ToBytes(Share{
					Key:   i,
					Value: prime,
				})
				share, err := FromBytes(bytes)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(share.Key).Should(Equal(i))
				Ω(share.Value.Cmp(&prime)).Should(Equal(0))
			}
		})
		It("should return an error when deserializing an empty byte slice", func() {
			_, err := FromBytes([]byte{})
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("splitting", func() {
		It("should return the correct number of shares", func() {
			// Shamir parameters.
			n := int64(100)
			k := int64(50)
			primeBig, ok := big.NewInt(0).SetString(primeStr, 10)
			Ω(ok).Should(Equal(true))
			secret := stackint.FromUint64(1234)
			secretBig := big.NewInt(1234)
			// Split the secret.
			shares, err := Split(n, k, &prime, &secret, primeBig, secretBig)
			Ω(err).Should(BeNil())
			Ω(int64(len(shares))).Should(Equal(n))
		})
	})

	Context("joining", func() {
		It("should return the correct secret from K shares", func() {
			// Shamir parameters.
			N := int64(100)
			K := int64(50)
			primeBig, ok := big.NewInt(0).SetString(primeStr, 10)
			Ω(ok).Should(Equal(true))
			secret := stackint.FromUint64(1234)
			secretBig := big.NewInt(1234)
			// Split the secret.
			shares, err := Split(N, K, &prime, &secret, primeBig, secretBig)
			Ω(err).Should(BeNil())
			Ω(int64(len(shares))).Should(Equal(N))
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
				decodedSecret, decodedSecretBig := Join(&prime, primeBig, kShares)
				fmt.Println("!!!")
				fmt.Println(decodedSecret.String())
				fmt.Println(decodedSecretBig.String())
				Ω(decodedSecret.Cmp(&secret)).Should(Equal(0))
			}
		})

		// 	It("should return an incorrect secret from less than K shares", func() {
		// 		// Shamir parameters.
		// 		N := int64(100)
		// 		K := int64(50)
		// 		prime, ok := big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)
		// 		Ω(ok).Should(Equal(true))
		// 		secret := big.NewInt(1234)
		// 		// Split the secret.
		// 		shares, err := Split(N, K, prime, secret)
		// 		Ω(err).Should(BeNil())
		// 		Ω(int64(len(shares))).Should(Equal(N))
		// 		// For all K less than 50 attempt to decode the secret.
		// 		for k := int64(1); k < 50; k++ {
		// 			// Pick K unique indices in the range [0, k).
		// 			indices := map[int]struct{}{}
		// 			for i := 0; i < int(k); i++ {
		// 				for {
		// 					index := rand.Intn(int(k))
		// 					if _, ok := indices[index]; !ok {
		// 						indices[index] = struct{}{}
		// 						break
		// 					}
		// 				}
		// 			}
		// 			// Use K shares to reconstruct the secret.
		// 			kShares := make(Shares, k)
		// 			for index := range indices {
		// 				kShares[index] = shares[index]
		// 			}
		// 			decodedSecret := Join(prime, kShares)
		// 			Ω(decodedSecret.Cmp(secret)).ShouldNot(Equal(0))
		// 		}
		// 	})
	})
})
