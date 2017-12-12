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

	// Setup big numbers so that we do not have to keep recreating them in each
	// loop.
	accum := big.NewInt(0)
	co := big.NewInt(0)
	base := big.NewInt(0)
	exp := big.NewInt(0)
	expMod := big.NewInt(0)

	// Create N shares.
	shares := make(Shares, shamir.N)
	for x := int64(1); x <= shamir.N; x++ {
		accum.Set(coefficients[0])
		base.SetInt64(x)
		exp.Set(base)
		expMod.Mod(exp, shamir.Prime)
		// Evaluate the polyomial at x.
		for _, coefficient := range coefficients[1:] {
			co.Set(coefficient)
			co.Mul(co, expMod)
			co.Mod(co, shamir.Prime)
			accum.Add(accum, co)
			accum.Mod(accum, shamir.Prime)
			exp.Mul(exp, base)
			expMod.Mod(exp, shamir.Prime)
		}
		shares[x-1] = Share{
			Key:   x,
			Value: big.NewInt(0).Set(accum),
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

	// Setup big numbers so that we do not have to keep recreating them in each
	// loop.
	value := big.NewInt(0)
	num := big.NewInt(1)
	den := big.NewInt(1)
	start := big.NewInt(0)
	next := big.NewInt(0)
	nextNeg := big.NewInt(0)
	nextDiff := big.NewInt(0)

	// Compute the Lagrange basic polynomial interpolation.
	for i := 0; i < len(shares); i++ {
		num.SetInt64(1)
		den.SetInt64(1)
		for j := 0; j < len(shares); j++ {
			if i == j {
				continue
			}
			start.SetInt64(int64(shares[i].Key))
			next.SetInt64(int64(shares[j].Key))
			nextNeg.SetInt64(0)
			nextNeg.Sub(nextNeg, next)
			num.Mul(num, nextNeg)
			num.Mod(num, shamir.Prime)
			nextDiff.Sub(start, next)
			den.Mul(den, nextDiff)
			den.Mod(den, shamir.Prime)
		}
		den.ModInverse(den, shamir.Prime)
		value.Mul(shares[i].Value, num)
		value.Mul(value, den)
		secret.Add(secret, value)
		secret.Mod(secret, shamir.Prime)
	}

	return secret, nil
}
