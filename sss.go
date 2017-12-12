package sss

import (
	"crypto/rand"
	"math/big"
)

type Share struct {
	Key   int64
	Value *big.Int
}

type Shares []Share

type Shamir struct {
	N     int64
	K     int64
	Prime *big.Int
}

func NewShamir(n int64, k int64, prime *big.Int) *Shamir {
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
	for i := int64(1); i < shamir.K; i++ {
		coefficient, err := rand.Int(rand.Reader, max)
		if err != nil {
			return nil, err
		}
		coefficients[i] = coefficient
	}

	// Create N shares.
	shares := make(Shares, shamir.N)
	for x := int64(1); x <= shamir.N; x++ {
		accum := big.NewInt(0).Set(coefficients[0])
		for exp := int64(1); exp < shamir.K; exp++ {
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

	// If we have more shares than necessary, take the first K shares.
	if int64(len(shares)) > shamir.K {
		shares = shares[:shamir.K]
	}

	for i := 0; i < len(shares); i++ {

		num := big.NewInt(1)
		den := big.NewInt(1)
		for j := 0; j < len(shares); j++ {
			if i == j {
				continue
			}

			start := big.NewInt(int64(shares[i].Key))
			next := big.NewInt(int64(shares[j].Key))

			negNext := big.NewInt(0)
			negNext.Sub(negNext, next)
			num.Mul(num, negNext)
			num.Mod(num, shamir.Prime)

			startDiffNext := big.NewInt(0).Set(start)
			startDiffNext.Sub(startDiffNext, next)
			den.Mul(den, startDiffNext)
			den.Mod(den, shamir.Prime)
		}

		value := big.NewInt(0).Set(shares[i].Value)
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
