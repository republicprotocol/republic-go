package hyperdrive

import (
	"bytes"
	"context"
	"encoding/binary"

	"golang.org/x/crypto/sha3"
)

func ProcessProposal(ctx context.Context, proposalChIn <-chan Proposal, signer Signer, verifier Verifier, capacity int) (<-chan Prepare, <-chan Fault, <-chan error) {
	prepareCh := make(chan Prepare, capacity)
	faultCh := make(chan Fault, capacity)
	errCh := make(chan error, capacity)

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

				if err := proposal.Verify(verifier); err != nil {
					fault := Fault{
						Rank:   proposal.Block.Rank,
						Height: proposal.Block.Height,
					}
					signature, err := signer.Sign(fault.Hash())
					if err != nil {
						errCh <- err
						continue
					}
					fault.Signatures = Signatures{signature}

					select {
					case <-ctx.Done():
						errCh <- ctx.Err()
						return
					case faultCh <- fault:
					}
					continue
				}

				prepare := Prepare{
					Proposal: proposal,
				}
				signature, err := signer.Sign(prepare.Hash())
				if err != nil {
					errCh <- err
					continue
				}
				prepare.Signatures = Signatures{signature}

				select {
				case <-ctx.Done():
					errCh <- ctx.Err()
					return
				case prepareCh <- prepare:
				}
			}
		}
	}()

	return prepareCh, faultCh, errCh
}

// ProposalHeader distinguishes Proposal from other message types that have the
// same content.
const ProposalHeader = byte(2)

// A Proposal message is sent by the Commander Replica to propose the next
// Block for preparation, committment, and finalization.
type Proposal struct {
	Block

	// Signature of the Replica that produced this Proposal
	Signature
}

// Hash implements the Hasher interface.
func (proposal *Proposal) Hash() Hash {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, ProposalHeader)
	binary.Write(&buf, binary.BigEndian, proposal.Block.Hash())
	return sha3.Sum256(buf.Bytes())
}

// Fault implements the Message interface.
func (proposal *Proposal) Fault() Fault {
	return Fault{
		Rank:   proposal.Block.Rank,
		Height: proposal.Block.Height,
	}
}

// Verify implements the Message interface.
func (proposal *Proposal) Verify(verifier Verifier) error {
	// TODO: Complete verification
	if err := proposal.Block.Verify(verifier); err != nil {
		return err
	}
	return verifier.VerifyProposer(proposal.Signature)
}

// SetSignatures implements the Message interface.
func (proposal *Proposal) SetSignatures(signatures Signatures) {
	return
}

// GetSignatures implements the Message interface.
func (proposal *Proposal) GetSignatures() Signatures {
	return Signatures{proposal.Signature}
}
