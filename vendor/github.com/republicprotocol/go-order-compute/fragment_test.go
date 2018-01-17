package compute_test

import (
	cryptoRand "crypto/rand"
	"math/big"
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/republicprotocol/go-order-compute"
	. "github.com/republicprotocol/go-sss"
)

var _ = Describe("Computations over fragments", func() {

	It("should correctly reconstruct an addition", func() {
		// Shamir parameters.
		N := int64(99)
		K := int64(50)
		prime, _ := big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)

		// Create the secrets.
		a := big.NewInt(400)
		b := big.NewInt(20)

		// Split the secrets.
		aShares, err := Split(N, K, prime, a)
		Ω(err).ShouldNot(HaveOccurred())
		bShares, err := Split(N, K, prime, b)
		Ω(err).ShouldNot(HaveOccurred())

		// Add the shares.
		outShares := make(Shares, N)
		for i := int64(0); i < N; i++ {
			out, err := Add(prime, aShares[i], bShares[i])
			Ω(err).ShouldNot(HaveOccurred())
			outShares[i] = out
		}

		// For all K greater than, or equal to, 50 attempt to decode the secret.
		for k := K; k <= N; k++ {
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
				kShares[index] = outShares[index]
			}
			// Join the output shares.
			ab, err := Join(prime, kShares)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(ab.Cmp(big.NewInt(0).Add(a, b))).Should(Equal(0))
		}
	})

	It("should correctly reconstruct a random number", func() {
		// Shamir parameters.
		N := int64(99)
		K := int64(50)
		prime, _ := big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)

		// Generate N random numbers, one for each participant
		randomMax := big.NewInt(0).Sub(prime, big.NewInt(1))
		randomNumbers := make([]*big.Int, N)
		globalRandomNumber := big.NewInt(0)
		for i := int64(0); i < N; i++ {
			r, err := cryptoRand.Int(cryptoRand.Reader, randomMax)
			Ω(err).ShouldNot(HaveOccurred())
			randomNumbers[i] = r
			globalRandomNumber.Add(globalRandomNumber, r)
			globalRandomNumber.Mod(globalRandomNumber, prime)
		}

		// Create shares for each random number.
		randomNumberShares := make([]Shares, N)
		for i := int64(0); i < N; i++ {
			shares, err := Split(N, K, prime, randomNumbers[i])
			Ω(err).ShouldNot(HaveOccurred())
			randomNumberShares[i] = shares
		}

		// For each participant, collect the relevant random shares from other
		// participants.
		outShares := make(Shares, N)
		for i := int64(0); i < N; i++ {
			randomNumberShare := randomNumberShares[i][i]
			for j := int64(0); j < N; j++ {
				if i == j {
					continue
				}
				r, err := Add(prime, randomNumberShare, randomNumberShares[j][i])
				Ω(err).ShouldNot(HaveOccurred())
				randomNumberShare = r
			}
			outShares[i] = randomNumberShare
		}

		// For all K greater than, or equal to, 50 attempt to decode the secret.
		for k := K; k <= N; k++ {
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
				kShares[index] = outShares[index]
			}
			// Join the output shares.
			randomNumber, err := Join(prime, kShares)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(randomNumber.Cmp(globalRandomNumber)).Should(Equal(0))
		}
	})

})
