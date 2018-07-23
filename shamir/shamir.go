package shamir

import (
	"bytes"
	"crypto/rsa"
	"encoding/binary"
	"encoding/json"
	"errors"
	"math/rand"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/stackint"
)

// ErrNKError is returned when the numbers of shared required to reconstruct a
// secret is greater than the number of shares the secret is split into.
var ErrNKError = errors.New("expected n to be greater than or equal to k")

// ErrFiniteField is returned when a secret is not in the finite field.
var ErrFiniteField = errors.New("expected secret to be in the finite field")

// ErrUnmarshalNilBytes is returned when trying to unmarshal a nil, or empty,
// byte slice.
var ErrUnmarshalNilBytes = errors.New("unmarshal nil bytes")

// Prime is the prime number used to define the finite field.
const Prime uint64 = 17012364981921935471

// A Share struct represents some share of a secret after the secret has been
// encoded.
type Share struct {
	Index uint64
	Value uint64
}

// Sub one share from another within the finite field and return the result.
// The index of the result will always be set to the receiver index.
func (share *Share) Sub(arg *Share) Share {
	return Share{
		Index: share.Index,
		Value: addMod(share.Value, subMod(Prime, arg.Value, Prime), Prime),
	}
}

// MarshalJSON implements the json.Marshaler interface.
func (share Share) MarshalJSON() ([]byte, error) {
	bytes, err := share.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return json.Marshal(bytes)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (share *Share) UnmarshalJSON(data []byte) error {
	var bytes []byte
	if err := json.Unmarshal(data, &bytes); err != nil {
		return err
	}
	return share.UnmarshalBinary(bytes)
}

// Equal returns true when both Shares are equal. Otherwise, it returns false.
func (share *Share) Equal(other *Share) bool {
	return share.Index == other.Index &&
		share.Value == other.Value
}

// MarshalBinary implements the encoding.BinaryMarshaler interface. The uint64
// index is encoded using binary.BigEndian and then the uint64 value is encoded
// using binary.BigEndian.
func (share Share) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, share.Index); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, share.Value); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface. The
// uint64 index is decoded using binary.BigEndian and then the uint64 value is
// decoded using binary.BigEndian.
func (share *Share) UnmarshalBinary(data []byte) error {
	if data == nil || len(data) == 0 {
		return ErrUnmarshalNilBytes
	}
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, &share.Index); err != nil {
		return err
	}
	return binary.Read(buf, binary.BigEndian, &share.Value)
}

// Encrypt a Share using an rsa.PublicKey.
func (share *Share) Encrypt(pubKey rsa.PublicKey) ([]byte, error) {
	rsaKey := crypto.RsaKey{PrivateKey: &rsa.PrivateKey{PublicKey: pubKey}}
	data, err := share.MarshalBinary()
	if err != nil {
		return []byte{}, err
	}
	return rsaKey.Encrypt(data)
}

// Decrypt cipher text into Share using an crypto.RsaKey.
func (share *Share) Decrypt(privKey *rsa.PrivateKey, cipherText []byte) error {
	rsaKey := crypto.RsaKey{PrivateKey: privKey}
	plainText, err := rsaKey.Decrypt(cipherText)
	if err != nil {
		return err
	}
	return share.UnmarshalBinary(plainText)
}

// Shares are a slice of Share structs.
type Shares []Share

// Split a secret into Shares. N represents the number of Shares that the
// secret will be split into, and K represents the number of Share required to
// reconstruct the secret. A slice of Shares, or an error, is returned.
func Split(n, k int64, secret uint64) (Shares, error) {
	// Validate the encoding by checking that N is greater than K, and that the
	// secret is within the finite field.
	if n < k {
		return nil, ErrNKError
	}
	if Prime <= secret {
		return nil, ErrFiniteField
	}

	// Generate K polynomial coefficients, where the first coefficient is the
	// secret.
	coefficients := make([]uint64, k)
	coefficients[0] = secret

	for i := int64(1); i < k; i++ {
		coefficients[i] = rand.Uint64()
		for coefficients[i] >= Prime {
			coefficients[i] = rand.Uint64()
		}
	}

	// Create N shares.
	shares := make(Shares, n)
	for x := int64(1); x <= n; x++ {

		accum := coefficients[0]
		base := uint64(x)
		exp := base % Prime

		// Evaluate the polyomial at x.
		for j := range coefficients[1:] {

			// co := (coefficients * expoMod) % prime
			coefficient := coefficients[j]
			co := mulMod(coefficient, exp, Prime)

			accum = addMod(accum, co, Prime)

			// exp = (exp * base ) % prime
			exp = mulMod(exp, base, Prime)
		}
		shares[x-1] = Share{
			Index: uint64(x),
			Value: accum,
		}
	}
	return shares, nil
}

// Join Shares into a secret. Prime is used to define the finite field from
// which the secret was selected. The reconstructed secret, or an error, is
// returned.
func Join(shares Shares) uint64 {
	secret := uint64(0)

	// Compute the Lagrange basic polynomial interpolation.
	for i := 0; i < len(shares); i++ {
		num := uint64(1)
		den := uint64(1)

		for j := 0; j < len(shares); j++ {
			if i == j {
				continue
			}
			// startposition = shares[formula][0];
			start := shares[i].Index

			// nextposition = shares[count][0];
			next := shares[j].Index

			// numerator = (numerator * -nextposition) % prime;
			nextGen := mulMod(num, next, Prime)
			num = Prime - nextGen

			// denominator = (denominator * (startposition - nextposition)) % prime;
			nextDiff := subMod(start, next, Prime)
			den = mulMod(den, nextDiff, Prime)
		}

		den = invMod(den, Prime)
		value := mulMod(shares[i].Value, num, Prime)
		value = mulMod(value, den, Prime)
		secret = addMod(secret, value, Prime)
	}

	return secret
}

func addMod(x uint64, y uint64, mod uint64) uint64 {
	stackX := stackint.FromUint(uint(x))
	stackY := stackint.FromUint(uint(y))
	stackM := stackint.FromUint(uint(mod))
	stackR := stackX.AddModulo(&stackY, &stackM)
	r, err := stackR.ToUint()
	if err != nil {
		panic(err)
	}
	return uint64(r)
}

func subMod(x uint64, y uint64, mod uint64) uint64 {
	stackX := stackint.FromUint(uint(x))
	stackY := stackint.FromUint(uint(y))
	stackM := stackint.FromUint(uint(mod))
	stackR := stackX.SubModulo(&stackY, &stackM)
	r, err := stackR.ToUint()
	if err != nil {
		panic(err)
	}
	return uint64(r)
}

func mulMod(x uint64, y uint64, mod uint64) uint64 {
	stackX := stackint.FromUint(uint(x))
	stackY := stackint.FromUint(uint(y))
	stackM := stackint.FromUint(uint(mod))
	stackR := stackX.MulModulo(&stackY, &stackM)
	r, err := stackR.ToUint()
	if err != nil {
		panic(err)
	}
	return uint64(r)
}

func invMod(x uint64, mod uint64) uint64 {
	stackX := stackint.FromUint(uint(x))
	stackM := stackint.FromUint(uint(mod))
	stackR := stackX.ModInverse(&stackM)
	r, err := stackR.ToUint()
	if err != nil {
		panic(err)
	}
	return uint64(r)
}
