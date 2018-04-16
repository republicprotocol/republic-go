package hyper

import (
	"context"

	"github.com/republicprotocol/republic-go/identity"
)

type Signer interface {
	Sign() Signature
	ID() identity.ID
}

type Signature string

type Proposal struct {
	Block

	// Signature of the Replica that produced this Proposal along with all of
	// the other Replicas that have signed it
	Signature
	Signatures []Signature
}

func ProcessProposal(ctx context.Context, proposalChIn <-chan Proposal, validator Validator) (<-chan Prepare, <-chan Fault, <-chan error) {
	prepareCh := make(chan Prepare, validator.Threshold())
	faultCh := make(chan Fault, validator.Threshold())
	errCh := make(chan error, validator.Threshold())
	counter := uint64(0)

	go func() {
		defer close(prepareCh)
		defer close(faultCh)
		defer close(errCh)

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case proposal, ok := <-proposalChIn:
				if !ok {
					return
				}
				counter++
				if validator.ValidateProposal(proposal) {
					prepare := Prepare{
						validator.Sign(),
						proposal.Block,
						proposal.Rank,
						proposal.Height,
					}
					// log.Println("Validated proposal on", validator.Sign())
					prepareCh <- prepare
				} else {
					fault := Fault{
						proposal.Rank,
						proposal.Height,
						validator.Sign(),
					}
					select {
					case <-ctx.Done():
						errCh <- ctx.Err()
						return
					case faultCh <- fault:
					}
				}
			}
		}
	}()

	return prepareCh, faultCh, errCh
}
