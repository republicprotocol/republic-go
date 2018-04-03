package hyper

import "context"

type Proposal struct {
}

func ProcessProposal(ctx context.Context, proposalChIn chan Proposal) (chan Prepare, chan Fault, chan error) {
	prepareCh := make(chan Prepare)
	faultCh := make(chan Fault)
	errCh := make(chan error)

	go func() {
		defer close(prepareCh)
		defer close(faultCh)
		defer close(errCh)

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case _ = <-proposalChIn:
				// TODO: process the proposal
			}
		}
	}()

	return prepareCh, faultCh, errCh
}
