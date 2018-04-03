package hyper

import "context"

type Signer interface {
	Sign([]byte) (Signature, error)
}

type Proposal struct {
	Signature
	Block
	Rank
	Height
}

func ProcessProposal(ctx context.Context, proposalChIn chan Proposal, signer Signer, sharedBlocks SharedBlocks) (chan Prepare, chan Fault, chan error) {
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
			case proposal, ok := <-proposalChIn:
				if !ok {
					return
				}
				valid := validateProposal(proposal, sharedBlocks)
				if valid {
					prepareCh <- Prepare{
						proposal.Signature,
						proposal.Block,
						proposal.Rank,
						proposal.Height,
					}
				} else {
					faultCh <- Fault{
						proposal.Rank,
						proposal.Height,
					}
				}
			}
		}
	}()

	return prepareCh, faultCh, errCh
}

func validateProposal(p Proposal, sb SharedBlocks) bool {
	valid := validateBlock(p.Block, sb)
	// TODO: check whether the signer inside block, and signer in the proposal are same
	// TODO: validate rank, make sure that the current rank and the rank of the proposer is same
	// TODO: validate height, make sure that the height is correct
	return valid
}
