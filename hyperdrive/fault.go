package hyper

import "context"

type Fault struct {
	rank   Rank
	height Height
}

type FaultHash [32]byte
type Faults map[FaultHash]uint8

func ProcessFault(ctx context.Context, faultChIn chan Fault) (chan Fault, chan error) {
	faultCh := make(chan Fault)
	errCh := make(chan error)
	faults := make(Faults)

	go func() {
		defer close(faultCh)
		defer close(errCh)

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case fault := <-faultChIn:
				h := getFaultHash(fault)
				if faults[h] >= THRESHOLD-1 {
					updateRank(fault)
				}
				faults[h]++
			}
		}
	}()

	return faultCh, errCh
}

func getFaultHash(f Fault) FaultHash {
	return FaultHash{}
}

func updateRank(f Fault) {}
