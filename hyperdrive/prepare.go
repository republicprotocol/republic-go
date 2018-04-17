package hyperdrive

import (
	"bytes"
	"context"
	"encoding/binary"

	"golang.org/x/crypto/sha3"
)

// PrepareHeader distinguishes Prepare from other message types that have the
// same content.
const PrepareHeader = byte(2)

// A Prepare messages signals that a Replica has received a valid Proposal.
type Prepare struct {
	Proposal

	// Signatures of the Replicas that signed this Prepare
	Signatures
}

// Hash implements the Hasher interface.
func (prepare *Prepare) Hash() Hash {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, PrepareHeader)
	binary.Write(&buf, binary.BigEndian, prepare.Proposal.Hash())
	return sha3.Sum256(buf.Bytes())
}

func (prepare *Prepare) Fault() Fault {
	return Fault{
		Rank:   prepare.Proposal.Block.Rank,
		Height: prepare.Proposal.Block.Height,
	}
}

// Verify the Prepare message. Returns an error if the message is invalid,
// otherwise nil.
func (prepare *Prepare) Verify() error {
	return nil
}

func (prepare *Prepare) SetSignatures(signatures Signatures) {
	prepare.Signatures = signatures
}

func (prepare *Prepare) GetSignatures() Signatures {
	return prepare.Signatures
}

func ProcessPreparation(ctx context.Context, prepareChIn <-chan Prepare, signer Signer, capacity, threshold int) (<-chan Commit, <-chan Fault, <-chan error) {
	commitCh := make(chan Commit, threshold)
	faultCh := make(chan Fault, threshold)
	errCh := make(chan error, threshold)

	go func() {
		defer close(commitCh)
		defer close(faultCh)
		defer close(errCh)

		store := NewMessageMapStore()

		for {
			select {

			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case prepare, ok := <-prepareChIn:
				if !ok {
					return
				}

				message, err := VerifyAndSignMessage(&prepare, &store, signer, threshold)
				if err != nil {
					errCh <- err
					continue
				}
				// After verifying and signing the message check for Faults
				switch message := message.(type) {

				case Prepare:
					commit := Commit{
						Prepare:    message,
						Signatures: message.Signatures,
					}
					select {
					case <-ctx.Done():
						errCh <- ctx.Err()
						return
					case commitCh <- commit:
					}
				case Fault:
					select {
					case <-ctx.Done():
						errCh <- ctx.Err()
						return
					case faultCh <- message:
					}
				default:
					// Gracefully ignore invalid messages
					continue
				}
			}
		}
	}()

	return commitCh, faultCh, errCh
}
