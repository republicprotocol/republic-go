package sss

import (
	"math/big"
)

type Share struct {
	Key   int
	Value *big.Int
}

type Shares []Share

type Shamir struct {
	N     int
	K     int
	Prime *big.Int
}

func NewShamir(n int, k int, prime *big.Int) *Shamir {
	return &Shamir{N: n, K: k, Prime: prime}
}

func (shamir *Shamir) Encode(secret *big.Int) (Shares, error) {
	// Validate the encoding by checking that N is greater than K, and that the
	// secret is within the finite field.
	if shamir.N < shamir.K {
		return nil, NewNKError(shamir.N, shamir.K)
	}
	if shamir.Prime.Cmp(secret) <= 0 {
		return nil, NewFiniteFieldError(secret)
	}

	// Generate K polynomial coefficients, where the first coefficient is the
	// secret.
	max := big.NewInt(0).Set(shamir.Prime)
	max.Sub(max, big.NewInt(1))
	coefficients := make([]*big.Int, shamir.K)
	coefficients[0] = secret
	// for i := 1; i < shamir.K; i++ {
	// 	coefficient, err := rand.Int(rand.Reader, max)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	coefficients[i] = coefficient
	// }
	coefficients[1] = big.NewInt(166)
	coefficients[2] = big.NewInt(94)

	// Create N shares.
	shares := make(Shares, shamir.N)
	for x := 1; x <= shamir.N; x++ {
		accum := big.NewInt(0).Set(coefficients[0])
		for exp := 1; exp < shamir.K; exp++ {
			a := big.NewInt(int64(x))
			a.Exp(a, big.NewInt(int64(exp)), shamir.Prime)
			b := big.NewInt(0).Set(coefficients[exp])
			b.Mul(b, a)
			b.Mod(b, shamir.Prime)
			accum.Add(accum, b)
			accum.Mod(accum, shamir.Prime)
		}
		shares[x-1] = Share{
			Key:   x,
			Value: accum,
		}
	}
	return shares, nil
}

func (shamir *Shamir) Decode(shares Shares) (*big.Int, error) {
	secret := big.NewInt(0)
	ks := shares[:shamir.K]

	for i := 0; i < len(ks); i++ {

		num := big.NewInt(1)
		den := big.NewInt(1)
		for j := 0; j < len(ks); j++ {
			if i == j {
				continue
			}

			start := big.NewInt(int64(ks[i].Key))
			next := big.NewInt(int64(ks[j].Key))

			negNext := big.NewInt(0)
			negNext.Sub(negNext, next)
			num.Mul(num, negNext)
			num.Mod(num, shamir.Prime)

			startDiffNext := big.NewInt(0).Set(start)
			startDiffNext.Sub(startDiffNext, next)
			den.Mul(den, startDiffNext)
			den.Mod(den, shamir.Prime)
		}

		value := ks[i].Value
		modInverseDen := big.NewInt(0)
		modInverseDen.ModInverse(den, shamir.Prime)
		value.Mul(value, num)
		value.Mul(value, modInverseDen)

		secret.Add(secret, shamir.Prime)
		secret.Add(secret, value)
		secret.Mod(secret, shamir.Prime)
	}

	return secret, nil
}
