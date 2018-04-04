package hyper

import (
	"bytes"
	"context"
	"encoding/binary"
)

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
					signedProposal, err := signProposal(proposal, signer)
					if err != nil {
						errCh <- err
					} else {
						prepareCh <- Prepare{
							signedProposal.Signature,
							signedProposal.Block,
							signedProposal.Rank,
							signedProposal.Height,
						}
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
	// TODO: validate rank, make sure that the current rank and the rank of the proposer is same
	// TODO: validate height, make sure that the height is correct
	return valid
}

func signProposal(p Proposal, signer Signer) (Proposal, error) {
	b, err := signBlock(p.Block, signer)
	if err != nil {
		return Proposal{}, err
	}
	p.Block = b
	var proposalBuf bytes.Buffer
	binary.Write(&proposalBuf, binary.BigEndian, p)
	sig, err := signer.Sign(proposalBuf.Bytes())
	return Proposal{
		sig,
		p.Block,
		p.Rank,
		p.Height,
	}, nil
}

func signBlock(b Block, signer Signer) (Block, error) {
	var blockBuf bytes.Buffer
	binary.Write(&blockBuf, binary.BigEndian, b)
	sig, err := signer.Sign(blockBuf.Bytes())
	if err != nil {
		return Block{}, err
	}
	return Block{
		b.tuples,
		sig,
	}, nil
}
