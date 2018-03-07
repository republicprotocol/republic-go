package shamir

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"math/big"
)

// A Share struct represents some share of a secret after the secret has been
// encoded.
type Share struct {
	Key   int64
	Value *big.Int
}

// Shares are a slice of Share structs.
type Shares []Share

// Split a secret into Shares. N represents the number of Shares that the
// secret will be split into, and K represents the number of Share required to
// reconstruct the secret. Prime is used to define the finite field from which
// secrets can be selected. A slice of Shares, or an error, is returned.
func Split(n int64, k int64, prime *big.Int, secret *big.Int) (Shares, error) {
	// Validate the encoding by checking that N is greater than K, and that the
	// secret is within the finite field.
	if n < k {
		return nil, NewNKError(n, k)
	}
	if prime.Cmp(secret) <= 0 {
		return nil, NewFiniteFieldError(secret)
	}

	// Generate K polynomial coefficients, where the first coefficient is the
	// secret.
	max := big.NewInt(0).Sub(prime, big.NewInt(1))
	coefficients := make([]*big.Int, k)
	coefficients[0] = secret
	for i := int64(1); i < k; i++ {
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
	shares := make(Shares, n)
	for x := int64(1); x <= n; x++ {
		accum.Set(coefficients[0])
		base.SetInt64(x)
		exp.Set(base)
		expMod.Mod(exp, prime)
		// Evaluate the polyomial at x.
		for _, coefficient := range coefficients[1:] {
			co.Set(coefficient)
			co.Mul(co, expMod)
			co.Mod(co, prime)
			accum.Add(accum, co)
			accum.Mod(accum, prime)
			exp.Mul(exp, base)
			expMod.Mod(exp, prime)
		}
		shares[x-1] = Share{
			Key:   x,
			Value: big.NewInt(0).Set(accum),
		}
	}
	return shares, nil
}

// Join Shares into a secret. Prime is used to define the finite field from
// which the secret was selected. The reconstructed secret, or an error, is
// returned.
func Join(prime *big.Int, shares Shares) *big.Int {
	secret := big.NewInt(0)

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
			num.Mod(num, prime)
			nextDiff.Sub(start, next)
			den.Mul(den, nextDiff)
			den.Mod(den, prime)
		}
		den.ModInverse(den, prime)
		value.Mul(shares[i].Value, num)
		value.Mul(value, den)
		secret.Add(secret, value)
		secret.Mod(secret, prime)
	}

	return secret
}

// ToBytes encodes the Share into a slice of bytes.
func ToBytes(share Share) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, share.Key)
	binary.Write(buf, binary.LittleEndian, share.Value.Bytes())
	return buf.Bytes()
}

// FromBytes decodes a slice of bytes into a Share.
func FromBytes(bs []byte) (Share, error) {
	key := int64(0)
	buf := bytes.NewReader(bs)
	if err := binary.Read(buf, binary.LittleEndian, &key); err != nil {
		return Share{}, err
	}
	data := make([]byte, buf.Len())
	if err := binary.Read(buf, binary.LittleEndian, data); err != nil {
		return Share{}, err
	}
	return Share{
		Key:   key,
		Value: big.NewInt(0).SetBytes(data),
	}, nil
}
