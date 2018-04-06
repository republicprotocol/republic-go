package smpc

import (
	"context"
	"crypto/rand"
	"time"

	"github.com/republicprotocol/republic-go/shamir"
	"github.com/republicprotocol/republic-go/stackint"
)

// ProduceObscureRngs by periodically writing to an output channel.
func ProduceObscureRngs(ctx context.Context, n, k int64) (<-chan ObscureRng, <-chan error) {
	obscureRngCh := make(chan ObscureRng)
	errCh := make(chan error)

	go func() {
		defer close(obscureRngCh)
		defer close(errCh)

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			random, err := stackint.Random(rand.Reader, &Prime)
			if err != nil {
				errCh <- err
				continue
			}
			obscureRng := ObscureRng{
				N:  n,
				K:  k,
				ID: random.Bytes(),
			}
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case <-ticker.C:
				select {
				case <-ctx.Done():
					errCh <- ctx.Err()
					return
				case obscureRngCh <- obscureRng:
				}
			}
		}
	}()

	return obscureRngCh, errCh
}

// ProcessObscureRngs by reading from an input channel, generating Shamir
// secret shares for a set of random numbers, and writing ObscureRngShares to
// an output channel.
func ProcessObscureRngs(ctx context.Context, obscureRngChIn <-chan ObscureRng) (<-chan ObscureRngShares, <-chan error) {
	obscureRngSharesCh := make(chan ObscureRngShares)
	errCh := make(chan error)

	go func() {
		defer close(obscureRngSharesCh)
		defer close(errCh)

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case obscureRng, ok := <-obscureRngChIn:
				if !ok {
					return
				}
				obscureRngShares, err := NewObscureRngShares(&obscureRng, &Prime)
				if err != nil {
					errCh <- err
					continue
				}
				select {
				case <-ctx.Done():
					errCh <- ctx.Err()
					return
				case obscureRngSharesCh <- obscureRngShares:
				}
			}
		}
	}()

	return obscureRngSharesCh, errCh
}

// ProcessObscureRngSharesIndexed by reading from an input channel, summating
// the indexed shares, and multiplying them, and producing ObscureMulShares by
// writing to an output channel.
func ProcessObscureRngSharesIndexed(ctx context.Context, obscureRngSharesIndexedChIn <-chan ObscureRngSharesIndexed) (<-chan ObscureMulShares, <-chan error) {
	obscureMulSharesCh := make(chan ObscureMulShares)
	errCh := make(chan error)

	go func() {
		defer close(obscureMulSharesCh)
		defer close(errCh)

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case obscureRngSharesIndexed, ok := <-obscureRngSharesIndexedChIn:
				if !ok {
					return
				}

				a := SummateShares(obscureRngSharesIndexed.A, &Prime)
				b := SummateShares(obscureRngSharesIndexed.B, &Prime)
				r := SummateShares(obscureRngSharesIndexed.R, &Prime)

				ab := shamir.Share{
					Key:   a.Key,
					Value: a.Value.MulModulo(&b.Value, &Prime),
				}
				rSq := shamir.Share{
					Key:   r.Key,
					Value: r.Value.MulModulo(&r.Value, &Prime),
				}
				abShares, err := shamir.Split(obscureRngSharesIndexed.N, obscureRngSharesIndexed.K, &Prime, &ab.Value)
				if err != nil {
					errCh <- err
					continue
				}
				rSqShares, err := shamir.Split(obscureRngSharesIndexed.N, obscureRngSharesIndexed.K, &Prime, &rSq.Value)
				if err != nil {
					errCh <- err
					continue
				}

				obscureMulShares := ObscureMulShares{
					ID:  obscureRngSharesIndexed.ID,
					AB:  abShares,
					RSq: rSqShares,
				}
				select {
				case <-ctx.Done():
					errCh <- ctx.Err()
					return
				case obscureMulSharesCh <- obscureMulShares:
				}
			}
		}
	}()

	return obscureMulSharesCh, errCh
}

// ProcessObscureMulSharesIndexed by reading from an input channel, using the
// indexed shares to produce a share of each multiplication, and finally
// producing an ObscureResidueFragment by writing to an output channel.
func ProcessObscureMulSharesIndexed(ctx context.Context, obscureMulSharesIndexedChIn <-chan ObscureMulSharesIndexed) (<-chan ObscureResidueFragment, <-chan error) {
	obscureResidueFragmentCh := make(chan ObscureResidueFragment)
	errCh := make(chan error)

	go func() {
		defer close(obscureResidueFragmentCh)
		defer close(errCh)

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case obscureMulSharesIndexed, ok := <-obscureMulSharesIndexedChIn:
				if !ok {
					return
				}

				t := obscureMulSharesIndexed.N

				// Reconstruct share of AB
				Hj := stackint.Zero()
				for i := range obscureMulSharesIndexed.AB {
					isNegative := λ[t][i] < 0
					λi := stackint.Zero()
					if isNegative {
						λi = stackint.FromUint(uint(-λ[t][i]))
					} else {
						λi = stackint.FromUint(uint(λ[t][i]))
					}
					λᵢhᵢj := λi.MulModulo(&obscureMulSharesIndexed.AB[i].Value, &Prime)
					if isNegative {
						Hj = Hj.SubModulo(&λᵢhᵢj, &Prime)
					} else {
						Hj = Hj.AddModulo(&λᵢhᵢj, &Prime)
					}
				}
				AB := shamir.Share{
					// FIXME: Do not assume the existence of an element.
					// We should probably look at a better way of knowing
					// which key a node is associated with.
					Key:   obscureMulSharesIndexed.AB[0].Key,
					Value: Hj,
				}

				// Reconstruct share of RSquared
				Hj = stackint.Zero()
				for i := range obscureMulSharesIndexed.RSq {
					isNegative := λ[t][i] < 0
					λi := stackint.Zero()
					if isNegative {
						λi = stackint.FromUint(uint(-λ[t][i]))
					} else {
						λi = stackint.FromUint(uint(λ[t][i]))
					}
					λᵢhᵢj := λi.MulModulo(&obscureMulSharesIndexed.RSq[i].Value, &Prime)
					if isNegative {
						Hj = Hj.SubModulo(&λᵢhᵢj, &Prime)
					} else {
						Hj = Hj.AddModulo(&λᵢhᵢj, &Prime)
					}
				}
				RSq := shamir.Share{
					// FIXME: Do not assume the existence of an element.
					// We should probably look at a better way of knowing
					// which key a node is associated with.
					Key:   obscureMulSharesIndexed.RSq[0].Key,
					Value: Hj,
				}

				// Output the ObscureResidueFragment
				obscureResidueFragment := ObscureResidueFragment{
					N:   obscureMulSharesIndexed.N,
					K:   obscureMulSharesIndexed.K,
					ID:  obscureMulSharesIndexed.ID,
					AB:  AB,
					RSq: RSq,
				}
				select {
				case <-ctx.Done():
					errCh <- ctx.Err()
					return
				case obscureResidueFragmentCh <- obscureResidueFragment:
				}
			}
		}
	}()

	return obscureResidueFragmentCh, errCh
}

// An ObscureResidueFragment is used to multiply a finite field element by a
// random quadratic residue. Iff the result is a quadratic residue, then the
// finite field element is also a quadratic residue.
type ObscureResidueFragment struct {
	// Used to parameterize the number of shares generated, and the number of
	// shares needed for reconstruction
	N, K int64

	// A unique identifier that links this ObscureResidueFragment to the
	// ObscureRng request that was used to generate it
	ID []byte

	// Beaver triple shares used for fast multiplication
	A, B, AB shamir.Share

	// A share of the random quadratic residue
	RSq shamir.Share
}

// An ObscureRng is used to request a round of obscure random number
// generation. This process results in all parties having shares of a random
// number without any party knowing the random number.
type ObscureRng struct {
	// Used to parameterize the number of shares generated, and the number of
	// shares needed for reconstruction
	N, K int64

	// A random unique identifier for this round of ObscureRng
	ID []byte
}

// ObscureRngShares are Shamir secret shares of a random number. All parties
// generate one of these during a round of ObscureRng, with one share for each
// party.
type ObscureRngShares struct {
	// Used to parameterize the number of shares generated, and the number of
	// shares needed for reconstruction
	N, K int64

	// The round of ObscureRng that these ObscureRngShares are genereted from
	ID []byte

	// Shamir secret shares for three random numbers
	A, B, R shamir.Shares
}

// NewObscureRngShares from an ObscureRng.
func NewObscureRngShares(obscureRng *ObscureRng, prime *stackint.Int1024) (ObscureRngShares, error) {
	var err error
	obscureRngShares := ObscureRngShares{
		N:  obscureRng.N,
		K:  obscureRng.K,
		ID: obscureRng.ID,
	}
	obscureRngShares.A, err = GenerateRandomShares(obscureRng.N, obscureRng.K, prime)
	if err != nil {
		return obscureRngShares, err
	}
	obscureRngShares.B, err = GenerateRandomShares(obscureRng.N, obscureRng.K, prime)
	if err != nil {
		return obscureRngShares, err
	}
	obscureRngShares.R, err = GenerateRandomShares(obscureRng.N, obscureRng.K, prime)
	if err != nil {
		return obscureRngShares, err
	}
	return obscureRngShares, nil
}

// ObscureRngSharesIndexed are all Shamir secret shares, from all parties, that
// have a specific index. Each Shamir secret share in an
// ObscureRngSharesIndexed is a share with a specific index, taken from
// ObscureRngShares.
type ObscureRngSharesIndexed struct {
	// Used to parameterize the number of shares generated, and the number of
	// shares needed for reconstruction
	N, K int64

	// The round of ObscureRng that these ObscureRngSharesIndexed are genereted
	// from
	ID []byte

	// Shamir secret shares, all with the same index, from different
	// ObscureRngShares
	A, B, R shamir.Shares
}

// ObscureMulShares are Shamir secret shares of a multiplication between two
// Shamir secret shares. All parties generate one of these, for their
// respective shares, during a round of multiplication.
type ObscureMulShares struct {
	// Used to parameterize the number of shares generated, and the number of
	// shares needed for reconstruction
	N, K int64

	// The round of ObscureRng that these ObscureMulShares are genereted from
	ID []byte

	// Shamir secret shares for the multiplication of shares summed from an
	// ObscureRngSharesIndexed
	AB, RSq shamir.Shares
}

// ObscureMulSharesIndexed are the same as ObscureRngSharesIndexed, but from
// ObscureMulShares.
type ObscureMulSharesIndexed struct {
	// Used to parameterize the number of shares generated, and the number of
	// shares needed for reconstruction
	N, K int64

	// The round of ObscureRng that these ObscureMulSharesIndexed are genereted
	// from
	ID []byte

	// Shamir secret shares, all with the same index, from different
	// ObscureMulShares
	AB, RSq shamir.Shares
}

// GenerateRandomShares in a finite field defined by the given prime.
func GenerateRandomShares(n, k int64, prime *stackint.Int1024) (shamir.Shares, error) {
	r, err := stackint.Random(rand.Reader, prime)
	if err != nil {
		return nil, err
	}
	rs, err := shamir.Split(n, k, prime, &r)
	if err != nil {
		return nil, err
	}
	return rs, nil
}

// SummateShares in a finite field defined by the given prime. All indices be
// the same.
func SummateShares(shares shamir.Shares, prime *stackint.Int1024) shamir.Share {
	if len(shares) == 0 {
		return shamir.Share{
			Key:   0,
			Value: stackint.Zero(),
		}
	}
	key := shares[0].Key
	value := shares[0].Value
	for i := 1; i < len(shares); i++ {
		value = value.AddModulo(&shares[i].Value, prime)
	}
	return shamir.Share{
		Key:   key,
		Value: value,
	}
}

// λ defines the first row of inverted square Van Der Monde matrices. The index
// defines the size of the Van Der Monde matrix, and the elements are defined by
// A(i,j) = i^(j-1).
var λ = map[int64][]int64{
	3:  {3, -3, 1},
	7:  {7, -21, 35, -35, 21, -7, 1},
	11: {11, -55, 165, -330, 462, -462, 330, -165, 55, -11, 1},
	15: {15, -105, 455, -1365, 3003, -5005, 6435, -6435, 5005, -3003, 1365, -455, 105, -15, 1},
	19: {19, -171, 969, -3876, 11628, -27132, 50388, -75582, 92378, -92378, 75582, -50388, 27132, -11628, 3876, -969, 171, -19, 1},
	23: {23, -253, 1771, -8855, 33649, -100947, 245157, -490314, 817190, -1144066, 1352078, -1352078, 1144066, -817190, 490314, -245157, 100947, -33649, 8855, -1771, 253, -23, 1},
	27: {27, -351, 2925, -17550, 80730, -296010, 888030, -2220075, 4686825, -8436285, 13037895, -17383860, 20058300, -20058300, 17383860, -13037895, 8436285, -4686825, 2220075, -888030, 296010, -80730, 17550, -2925, 351, -27, 1},
	31: {31, -465, 4495, -31465, 169911, -736281, 2629575, -7888725, 20160075, -44352165, 84672315, -141120525, 206253075, -265182525, 300540195, -300540195, 265182525, -206253075, 141120525, -84672315, 44352165, -20160075, 7888725, -2629575, 736281, -169911, 31465, -4495, 465, -31, 1},
}
