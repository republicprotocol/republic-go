package shamir

import (
	"bytes"
	"crypto/rsa"
	"encoding/binary"
	"encoding/json"
	"errors"
	"math/big"
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

// Constants used for Pedersen commitments.
var (
	CommitG, _ = big.NewInt(0).SetString("55770337541865622645762792995310051207397982365281550224398537343173859535105139234635957625862144230699334804049712012240627186833575368075584941905877907874721866477877240648297575495799055309964125528850827031430744090495932459174225210240271018890912546873876168921012222806248779191616184071117514042939", 10)
	CommitH, _ = big.NewInt(0).SetString("115473249134132086626276466548821312109767217648495447434241772177665837404180106795894544927519633232260996631113705689399974939497093162658781255907718160986084723710881782554865045754197429739205723247463738713971412477405153816377087610084931720375701175784456548744957930969662291742829878205406745940998", 10)
	CommitP, _ = big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430638301075261977215604623753829209190123552924403809454455467971378126933547351328256473517128194204856274749623411850633650752318731262510900221876378083808528817", 10)
)

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

		// Evaluate the polynomial at x.
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

type Blindings []Blinding

type Blinding struct {
	*big.Int
}

func (b *Blinding) Encrypt(pubKey rsa.PublicKey) ([]byte, error) {
	rsaKey := crypto.RsaKey{PrivateKey: &rsa.PrivateKey{PublicKey: pubKey}}
	if b.Int == nil {
		return []byte{}, nil
	}
	data := b.Int.Bytes()
	return rsaKey.Encrypt(data)
}

func (b *Blinding) Decrypt(privKey *rsa.PrivateKey, cipherText []byte) error {
	if cipherText == nil || len(cipherText) == 0 {
		return nil
	}
	rsaKey := crypto.RsaKey{PrivateKey: privKey}
	bs, err := rsaKey.Decrypt(cipherText)
	if err != nil {
		return err
	}
	if b.Int == nil {
		b.Int = big.NewInt(0)
	}
	b.Int.SetBytes(bs)
	return nil
}

func (b Blinding) MarshalJSON() ([]byte, error) {
	if b.Int == nil {
		return json.Marshal([]byte{0})
	}
	return json.Marshal(b.Int.Bytes())
}

func (b *Blinding) UnmarshalJSON(data []byte) error {
	bs := []byte{}
	if err := json.Unmarshal(data, &bs); err != nil {
		return err
	}

	if b.Int == nil {
		b.Int = big.NewInt(0)
	}
	b.Int.SetBytes(bs)

	return nil
}

type Commitment struct {
	*big.Int
}

func NewCommitment(x Share, s Blinding) Commitment {
	gˣ := big.NewInt(0).Exp(CommitG, big.NewInt(0).SetUint64(x.Value), CommitP)
	hˢ := big.NewInt(0).Exp(CommitH, s.Int, CommitP)
	gˣhˢ := big.NewInt(0).Mul(gˣ, hˢ)
	return Commitment{gˣhˢ.Mod(gˣhˢ, CommitP)}
}
