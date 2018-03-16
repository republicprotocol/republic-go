package shamir

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/republicprotocol/republic-go/stackint"
)

// A Share struct represents some share of a secret after the secret has been
// encoded.
type Share struct {
	Key      int64
	Value    stackint.Int1024
	ValueBig *big.Int
}

// Shares are a slice of Share structs.
type Shares []Share

// Split a secret into Shares. N represents the number of Shares that the
// secret will be split into, and K represents the number of Share required to
// reconstruct the secret. Prime is used to define the finite field from which
// secrets can be selected. A slice of Shares, or an error, is returned.
func Split(n int64, k int64, prime, secret *stackint.Int1024, primeBig, secretBig *big.Int) (Shares, error) {
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
	maxBig := big.NewInt(0).Sub(primeBig, big.NewInt(1))
	coefficients := make([]*stackint.Int1024, k)
	coefficients[0] = secret
	coefficientsBig := make([]*big.Int, k)
	coefficientsBig[0] = secretBig

	for i := int64(1); i < k; i++ {
		coefficient, err := stackint.Random(rand.Reader, &max)
		if err != nil {
			return nil, err
		}
		coefficients[i] = &coefficient
		coefficientBig, err := rand.Int(rand.Reader, maxBig)
		if err != nil {
			return nil, err
		}
		coefficientsBig[i] = coefficientBig
	}

	// Setup big numbers so that we do not have to keep recreating them in each
	// loop.
	accumBig := big.NewInt(0)
	coBig := big.NewInt(0)
	baseBig := big.NewInt(0)
	expBig := big.NewInt(0)
	expModBig := big.NewInt(0)

	// Create N shares.
	shares := make(Shares, n)
	for x := int64(1); x <= n; x++ {
		accum := (*coefficients[0]).Clone()
		accumBig.Set(coefficientsBig[0])
		base := stackint.FromUint64(uint64(x))
		baseBig.SetInt64(x)
		exp := base
		expBig.Set(baseBig)
		expMod := exp.Mod(prime)
		expModBig.Mod(expBig, primeBig)
		// Evaluate the polyomial at x.
		for _, coefficient := range coefficients[1:] {
			co := coefficient.Mul(&expMod)
			coBig.Mul(coBig, expModBig)
			co = co.Mod(prime)
			coBig.Mod(coBig, primeBig)
			accum.Inc(&co)
			accumBig.Add(accumBig, coBig)
			accum = accum.Mod(prime)
			accumBig.Mod(accumBig, primeBig)
			exp = exp.Mul(&base)
			expBig.Mul(expBig, baseBig)
			expMod = exp.Mod(prime)
			expModBig.Mod(expBig, primeBig)
		}
		expBig, _ := big.NewInt(0).SetString(accum.String(), 0)
		if expBig.Cmp(accumBig) != 0 {
			fmt.Println(accum.String())
			fmt.Println(accumBig.String())
			// panic("!!!")
		}
		shares[x-1] = Share{
			Key:      x,
			Value:    accum,
			ValueBig: accumBig,
		}
	}
	return shares, nil
}

// Join Shares into a secret. Prime is used to define the finite field from
// which the secret was selected. The reconstructed secret, or an error, is
// returned.
func Join(prime *stackint.Int1024, primeBig *big.Int, shares Shares) (*stackint.Int1024, *big.Int) {
	secret := stackint.Zero()
	secretBig := big.NewInt(0)

	// Setup big numbers so that we do not have to keep recreating them in each
	// loop.
	valueBig := big.NewInt(0)
	numBig := big.NewInt(1)
	denBig := big.NewInt(1)
	startBig := big.NewInt(0)
	nextBig := big.NewInt(0)
	nextNegBig := big.NewInt(0)
	nextDiffBig := big.NewInt(0)

	// Compute the Lagrange basic polynomial interpolation.
	for i := 0; i < len(shares); i++ {
		num := stackint.One()
		den := stackint.One()

		numBig.SetInt64(1)
		denBig.SetInt64(1)

		for j := 0; j < len(shares); j++ {
			if i == j {
				continue
			}
			// startposition = shares[formula][0];
			startBig.SetInt64(int64(shares[i].Key))
			start := stackint.FromUint64(uint64(shares[i].Key))

			// nextposition = shares[count][0];
			nextBig.SetInt64(int64(shares[j].Key))
			next := stackint.FromUint64(uint64(shares[j].Key))

			// numerator = (numerator * -nextposition) % prime;
			nextNegBig.SetInt64(0)
			nextNegBig.Sub(nextNegBig, nextBig)
			numBig.Mul(numBig, nextNegBig)
			numBig.Mod(numBig, primeBig)
			nextGen := num.Mul(&next)
			nextGen = nextGen.Mod(prime)
			num = prime.Sub(&nextGen)

			// denominator = (denominator * (startposition - nextposition)) % prime;
			nextDiffBig.Sub(startBig, nextBig)
			denBig.Mul(denBig, nextDiffBig)
			denBig.Mod(denBig, primeBig)
			nextDiff := start.SubModulo(&next, prime)
			den = den.Mul(&nextDiff)
			den = den.Mod(prime)
		}

		// valueBig = shares[formula][1]
		// accumBig = (primeBig + accumBig + (valueBig * numeratorBig * modInverse(denominatorBig))) % primeBig
		denBig.ModInverse(denBig, primeBig)
		valueBig.Mul(shares[i].ValueBig, numBig)
		valueBig.Mul(valueBig, denBig)
		secretBig.Add(secretBig, valueBig)
		secretBig.Mod(secretBig, primeBig)

		den = den.ModInverse(prime)
		value := shares[i].Value.Mul(&num)
		value = value.Mul(&den)
		secret.Inc(&value)
		secret = secret.Mod(prime)
	}

	return &secret, secretBig
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
		Key:      key,
		Value:    stackint.FromBytes(data),
		ValueBig: big.NewInt(0).SetBytes(data),
	}, nil
}
