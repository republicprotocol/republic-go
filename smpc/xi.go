package smpc

import (
	"context"
	"crypto/rand"
	"time"

	"github.com/republicprotocol/republic-go/stackint"

	"github.com/republicprotocol/republic-go/shamir"
)

func ProduceXiFragmentGenerators(ctx context.Context, n, k int64) (<-chan XiFragmentGenerator, <-chan error) {
	xiFragmentGeneratorCh := make(chan XiFragmentGenerator)
	errCh := make(chan error)

	go func() {
		defer close(xiFragmentGeneratorCh)
		defer close(errCh)

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			random, err := stackint.Random(rand.Reader, &Prime)
			if err != nil {
				errCh <- err
				continue
			}
			xiFragmentGenerator := XiFragmentGenerator{
				ID: random,
				N:  n,
				K:  k,
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
				case xiFragmentGeneratorCh <- xiFragmentGenerator:
				}
			}
		}
	}()

	return xiFragmentGeneratorCh, errCh
}

func ProcessXiFragmentGenerators(ctx context.Context, xiFragmentGeneratorChIn <-chan XiFragmentGenerator) (<-chan XiFragmentAdditiveShares, <-chan error) {
	xiFragmentAdditiveSharesCh := make(chan XiFragmentAdditiveShares)
	errCh := make(chan error)

	go func() {
		defer close(xiFragmentAdditiveSharesCh)
		defer close(errCh)

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case xiFragmentGenerator, ok := <-xiFragmentGeneratorChIn:
				if !ok {
					return
				}
				xiFragmentAdditiveShares, err := RandomXiFragmentAdditiveShares(&xiFragmentGenerator, &Prime)
				if err != nil {
					errCh <- err
					continue
				}
				select {
				case <-ctx.Done():
					errCh <- ctx.Err()
					return
				case xiFragmentAdditiveSharesCh <- xiFragmentAdditiveShares:
				}
			}
		}
	}()

	return xiFragmentAdditiveSharesCh, errCh
}

func ProcessXiFragmentAdditiveShares(ctx context.Context, xiFragmentAdditiveSharesChIn <-chan XiFragmentAdditiveShares) (<-chan XiFragmentMultiplicativeShares, <-chan error) {
	xiFragmentMultiplicativeSharesCh := make(chan XiFragmentMultiplicativeShares)
	errCh := make(chan error)

	go func() {
		defer close(xiFragmentMultiplicativeSharesCh)
		defer close(errCh)

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case xiFragmentAdditiveShares, ok := <-xiFragmentAdditiveSharesChIn:
				if !ok {
					return
				}

				a := SummateShares(xiFragmentAdditiveShares.A, &Prime)
				b := SummateShares(xiFragmentAdditiveShares.B, &Prime)
				r := SummateShares(xiFragmentAdditiveShares.R, &Prime)

				ab := shamir.Share{
					Key:   a.Key,
					Value: a.Value.MulModulo(&b.Value, &Prime),
				}
				rSquared := shamir.Share{
					Key:   r.Key,
					Value: r.Value.MulModulo(&r.Value, &Prime),
				}
				abShares, err := shamir.Split(xiFragmentAdditiveShares.N, xiFragmentAdditiveShares.K, &Prime, &ab.Value)
				if err != nil {
					errCh <- err
					continue
				}
				rSquaredShares, err := shamir.Split(xiFragmentAdditiveShares.N, xiFragmentAdditiveShares.K, &Prime, &rSquared.Value)
				if err != nil {
					errCh <- err
					continue
				}

				xiFragmentMultiplicativeShares := XiFragmentMultiplicativeShares{
					ID:       xiFragmentAdditiveShares.ID,
					AB:       abShares,
					RSquared: rSquaredShares,
				}
				select {
				case <-ctx.Done():
					errCh <- ctx.Err()
					return
				case xiFragmentMultiplicativeSharesCh <- xiFragmentMultiplicativeShares:
				}
			}
		}
	}()

	return xiFragmentMultiplicativeSharesCh, errCh
}

func ProcessXiFragmentMultiplicativeShares(ctx context.Context, xiFragmentMultiplicativeSharesChIn <-chan XiFragmentMultiplicativeShares) (<-chan XiFragment, <-chan error) {
	xiFragmentCh := make(chan XiFragment)
	errCh := make(chan error)

	go func() {
		defer close(xiFragmentCh)
		defer close(errCh)

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case xiFragmentMultiplicativeShares, ok := <-xiFragmentMultiplicativeSharesChIn:
				if !ok {
					return
				}

				t := xiFragmentMultiplicativeShares.N

				// Reconstruct share of AB
				ABHj := stackint.Zero()
				for i := range xiFragmentMultiplicativeShares.AB {
					isNegative := false
					λi := λ[t][i]
					if λi < 0 {
						isNegative = true
						λi = -λi
					}
					bigλi := stackint.FromUint(uint(λi))
					innerProduct := bigλi.MulModulo(&xiFragmentMultiplicativeShares.AB[i].Value, &Prime)
					if isNegative {
						ABHj = ABHj.SubModulo(&innerProduct, &Prime)
					} else {
						ABHj = ABHj.AddModulo(&innerProduct, &Prime)
					}
				}

				// Reconstruct share of RSquared
				RSquaredHj := stackint.Zero()
				for i := range xiFragmentMultiplicativeShares.RSquared {
					isNegative := false
					λi := λ[t][i]
					if λi < 0 {
						isNegative = true
						λi = -λi
					}
					bigλi := stackint.FromUint(uint(λi))
					innerProduct := bigλi.MulModulo(&xiFragmentMultiplicativeShares.RSquared[i].Value, &Prime)
					if isNegative {
						RSquaredHj = RSquaredHj.SubModulo(&innerProduct, &Prime)
					} else {
						RSquaredHj = RSquaredHj.AddModulo(&innerProduct, &Prime)
					}
				}

				// Output the XiFragment
				xiFragment := XiFragment{
					ID: xiFragmentMultiplicativeShares.ID,
					N:  xiFragmentMultiplicativeShares.N,
					K:  xiFragmentMultiplicativeShares.K,
					AB: shamir.Share{
						// FIXME: Do not assume the existence of an element.
						// We should probably look at a better way of knowing
						// which key a node is associated with.
						Key:   xiFragmentMultiplicativeShares.AB[0].Key,
						Value: ABHj,
					},
					RSquared: shamir.Share{
						// FIXME: Same as above.
						Key:   xiFragmentMultiplicativeShares.RSquared[0].Key,
						Value: RSquaredHj,
					},
				}
				select {
				case <-ctx.Done():
					errCh <- ctx.Err()
					return
				case xiFragmentCh <- xiFragment:
				}
			}
		}
	}()

	return xiFragmentCh, errCh
}

type XiFragmentGenerator struct {
	ID stackint.Int1024
	N  int64
	K  int64
}

type XiFragmentAdditiveShares struct {
	ID stackint.Int1024
	N  int64
	K  int64
	A  shamir.Shares
	B  shamir.Shares
	R  shamir.Shares
}

type XiFragmentAdditive struct {
	ID stackint.Int1024
	N  int64
	K  int64
	A  shamir.Share
	B  shamir.Share
	R  shamir.Share
}

type XiFragmentMultiplicativeShares struct {
	ID       stackint.Int1024
	N        int64
	K        int64
	AB       shamir.Shares
	RSquared shamir.Shares
}

type XiFragment struct {
	ID       stackint.Int1024
	N        int64
	K        int64
	AB       shamir.Share
	RSquared shamir.Share
}

func RandomXiFragmentAdditiveShares(xiFragmentGenerator *XiFragmentGenerator, prime *stackint.Int1024) (XiFragmentAdditiveShares, error) {
	var err error
	xiFragmentAdditiveShares := XiFragmentAdditiveShares{
		ID: xiFragmentGenerator.ID,
	}
	xiFragmentAdditiveShares.A, err = GenerateRandomShares(xiFragmentGenerator.N, xiFragmentGenerator.K, prime)
	if err != nil {
		return xiFragmentAdditiveShares, err
	}
	xiFragmentAdditiveShares.B, err = GenerateRandomShares(xiFragmentGenerator.N, xiFragmentGenerator.K, prime)
	if err != nil {
		return xiFragmentAdditiveShares, err
	}
	xiFragmentAdditiveShares.R, err = GenerateRandomShares(xiFragmentGenerator.N, xiFragmentGenerator.K, prime)
	if err != nil {
		return xiFragmentAdditiveShares, err
	}
	return xiFragmentAdditiveShares, nil
}

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
