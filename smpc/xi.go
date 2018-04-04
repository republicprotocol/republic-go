package smpc

import (
	"context"
	"crypto/rand"
	"log"
	"time"

	"github.com/republicprotocol/republic-go/stackint"

	"github.com/republicprotocol/republic-go/shamir"
)

func ProduceXiFragmentGenerators(ctx context.Context, n, k int64) (chan XiFragmentGenerator, chan error) {
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

func ProcessXiFragmentGenerators(ctx context.Context, xiFragmentGeneratorChIn chan XiFragmentGenerator) (chan XiFragmentAdditiveShares, chan error) {
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

func ProcessXiFragmentAdditiveShares(ctx context.Context, xiFragmentAdditiveSharesChIn chan XiFragmentAdditiveShares) (chan XiFragmentMultiplicativeShares, chan error) {
	xiFragmentMultiplicativeSharesCh := make(chan XiFragmentMultiplicativeShares)
	errCh := make(chan error)

	go func() {
		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case xiFragmentAdditiveShares, ok := <-xiFragmentAdditiveSharesChIn:
				if !ok {
					return
				}

				a := SumShares(xiFragmentAdditiveShares.A)
				b := SumShares(xiFragmentAdditiveShares.B)
				r := SumShares(xiFragmentAdditiveShares.R)

				ab := shamir.Share{
					Key:   a.Key,
					Value: a.Value.MulModulo(&b.Value, &Prime),
				}
				rSquared := shamir.Share{
					Key:   r.Key,
					Value: r.Value.MulModulo(&r.Value, &Prime),
				}

				xiFragmentMultiplicativeShares := XiFragmentMultiplicativeShares{
					ID: xiFragmentAdditiveShares.ID,
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

func ConsumeXiFragmentMultiplicativeShares(ctx context.Context, xiFragmentMultiplicativeSharesChIn chan XiFragmentMultiplicativeShares) chan error {
	errCh := make(chan error)

	go func() {
		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case xiFragmentMultiplicativeShares, ok := <-xiFragmentMultiplicativeSharesChIn:
				if !ok {
					return
				}
				log.Println(xiFragmentMultiplicativeShares)
			}
		}
	}()

	return errCh
}

type XiFragmentGenerator struct {
	ID stackint.Int1024
	N  int64
	K  int64
}

type XiFragmentAdditiveShares struct {
	ID   stackint.Int1024
	A, B shamir.Shares
	R    shamir.Shares
}

type XiFragmentAdditive struct {
	ID   stackint.Int1024
	A, B shamir.Share
	R    shamir.Share
}

type XiFragmentMultiplicativeShares struct {
	ID       stackint.Int1024
	AB       shamir.Shares
	RSquared shamir.Shares
}

func RandomXiFragmentAdditiveShares(xiFragmentGenerator *XiFragmentGenerator, prime *stackint.Int1024) (XiFragmentAdditiveShares, error) {
	var err error
	xiFragmentAdditiveShares := XiFragmentAdditiveShares{
		ID: xiFragmentGenerator.ID,
	}
	xiFragmentAdditiveShares.A, err = RandomShares(xiFragmentGenerator.N, xiFragmentGenerator.K, prime)
	if err != nil {
		return xiFragmentAdditiveShares, err
	}
	xiFragmentAdditiveShares.B, err = RandomShares(xiFragmentGenerator.N, xiFragmentGenerator.K, prime)
	if err != nil {
		return xiFragmentAdditiveShares, err
	}
	xiFragmentAdditiveShares.R, err = RandomShares(xiFragmentGenerator.N, xiFragmentGenerator.K, prime)
	if err != nil {
		return xiFragmentAdditiveShares, err
	}
	return xiFragmentAdditiveShares, nil
}

func RandomShares(n, k int64, prime *stackint.Int1024) (shamir.Shares, error) {
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

func SumShares(shares shamir.Shares) shamir.Share {
	if len(shares) == 0 {
		return shamir.Share{
			Key:   0,
			Value: stackint.Zero(),
		}
	}
	key := shares[0].Key
	value := shares[0].Value
	for i := 1; i < len(shares); i++ {
		value = value.Add(&shares[i].Value)
	}
	return shamir.Share{
		Key:   key,
		Value: value,
	}
}
