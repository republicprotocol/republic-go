package hyperdrive

import (
	"bytes"
	"context"
	"encoding/binary"

	"github.com/republicprotocol/republic-go/identity"
	"golang.org/x/crypto/sha3"
)

func ProcessProposal(ctx context.Context, proposalChIn <-chan Proposal, signer identity.Signer, verifier identity.Verifier, capacity int) (<-chan Prepare, <-chan Fault, <-chan error) {
	prepareCh := make(chan Prepare, capacity)
	faultCh := make(chan Fault, capacity)
	errCh := make(chan error, capacity)

	go func() {
		defer close(prepareCh)
		defer close(faultCh)
		defer close(errCh)

		store := NewMessageMapStore()

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return

			case proposal, ok := <-proposalChIn:
				if !ok {
					return
				}

				message, err := VerifyAndSignMessage(&proposal, &store, signer, verifier, 0)

				if err != nil {
					errCh <- err
					continue
				}

				// After verifying and signing the message check for Faults
				switch message := message.(type) {
				case *Proposal:
					prepare := Prepare{
						Proposal: *message,
					}
					signature, err := signer.Sign(prepare.Hash())
					if err != nil {
						errCh <- err
						continue
					}
					prepare.Signatures = identity.Signatures{signature}

					select {
					case <-ctx.Done():
						errCh <- ctx.Err()
						return
					case prepareCh <- prepare:
					}
				case *Fault:
					select {
					case <-ctx.Done():
						errCh <- ctx.Err()
						return
					case faultCh <- *message:
					}
				default:
					// Gracefully ignore invalid messages
					continue
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
// Block for preparation, commitment, and finalization.
type Proposal struct {
	Block

	// Signature of the Replica that produced this Proposal
	identity.Signature
}

// Hash implements the Hasher interface.
func (proposal *Proposal) Hash() identity.Hash {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, ProposalHeader)
	binary.Write(&buf, binary.BigEndian, proposal.Block)

	return sha3.Sum256(buf.Bytes())
}

// Fault implements the Message interface.
func (proposal *Proposal) Fault() *Fault {
	return &Fault{
		Rank:   proposal.Block.Rank,
		Height: proposal.Block.Height,
	}
}

// Verify implements the Message interface.
func (proposal *Proposal) Verify(verifier identity.Verifier) error {

	// TODO: Complete verification
	if err := proposal.Block.VerifyBlock(verifier); err != nil {
		return err
	}
	return verifier.VerifySignature(proposal.Signature)
}

// SetSignatures implements the Message interface.
func (proposal *Proposal) SetSignatures(signatures identity.Signatures) {
	return
}

// GetSignatures implements the Message interface.
func (proposal *Proposal) GetSignatures() identity.Signatures {
	return identity.Signatures{proposal.Signature}
}
