package shamir

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"

	"github.com/republicprotocol/republic-go/stackint"
)

// A Share struct represents some share of a secret after the secret has been
// encoded.
type Share struct {
	Key   int64
	Value stackint.Int1024
}

// Shares are a slice of Share structs.
type Shares []Share

// Split a secret into Shares. N represents the number of Shares that the
// secret will be split into, and K represents the number of Share required to
// reconstruct the secret. Prime is used to define the finite field from which
// secrets can be selected. A slice of Shares, or an error, is returned.
func Split(n int64, k int64, prime, secret *stackint.Int1024) (Shares, error) {
	// Validate the encoding by checking that N is greater than K, and that the
	// secret is within the finite field.
	if n < k {
		return nil, NewNKError(n, k)
	}
	if prime.LessThanOrEqual(secret) {
		return nil, NewFiniteFieldError(secret)
	}

	// Generate K polynomial coefficients, where the first coefficient is the
	// secret.
	one := stackint.One()
	max := prime.Sub(&one)
	coefficients := make([]*stackint.Int1024, k)
	coefficients[0] = secret

	for i := int64(1); i < k; i++ {
		coefficient, err := stackint.Random(rand.Reader, &max)
		if err != nil {
			return nil, err
		}
		coefficients[i] = &coefficient
	}

	// Create N shares.
	shares := make(Shares, n)
	for x := int64(1); x <= n; x++ {

		// accum := coefficients[0]
		accum := (*coefficients[0]).Clone()
		// base := x
		base := stackint.FromUint64(uint64(x))
		// expMod := base % prime
		exp := base.Clone()
		expMod := exp.Mod(prime)

		// Evaluate the polyomial at x.
		for j := range coefficients[1:] {

			// co := (coefficients * expoMod) % prime
			coefficient := coefficients[j].Clone()
			co := coefficient.MulModulo(&expMod, prime)

			// accum = (accum + co) % prime
			accum = accum.AddModulo(&co, prime)

			// expMod = (expMod * base ) % prime
			exp = exp.Mul(&base)
			expMod = exp.Mod(prime)
		}
		shares[x-1] = Share{
			Key:   x,
			Value: accum,
		}
	}
	return shares, nil
}

// Join Shares into a secret. Prime is used to define the finite field from
// which the secret was selected. The reconstructed secret, or an error, is
// returned.
func Join(prime *stackint.Int1024, shares Shares) *stackint.Int1024 {
	secret := stackint.Zero()

	// Compute the Lagrange basic polynomial interpolation.
	for i := 0; i < len(shares); i++ {
		num := stackint.One()
		den := stackint.One()

		for j := 0; j < len(shares); j++ {
			if i == j {
				continue
			}
			// startposition = shares[formula][0];
			start := stackint.FromUint64(uint64(shares[i].Key))

			// nextposition = shares[count][0];
			next := stackint.FromUint64(uint64(shares[j].Key))

			// numerator = (numerator * -nextposition) % prime;
			nextGen := num.MulModulo(&next, prime)
			num = prime.Sub(&nextGen)

			// denominator = (denominator * (startposition - nextposition)) % prime;
			nextDiff := start.SubModulo(&next, prime)
			den = den.MulModulo(&nextDiff, prime)
		}

		den = den.ModInverse(prime)
		value := shares[i].Value.MulModulo(&num, prime)
		value = value.MulModulo(&den, prime)
		secret = secret.AddModulo(&value, prime)
	}

	return &secret
}

// ToBytes encodes the Share into a slice of bytes.
func ToBytes(share Share) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, share.Key)
	binary.Write(buf, binary.LittleEndian, share.Value.ToBytes())
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
		Value: stackint.FromBytes(data),
	}, nil
}
