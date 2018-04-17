package hyperdrive

import (
	"bytes"
	"context"
	"encoding/binary"

	"golang.org/x/crypto/sha3"
)

// CommitHeader distinguishes Commits from other message types that have the
// same content.
const CommitHeader = byte(3)

// A Commit messages signals that a Replica wants to commit to the finalization
// of a Block.
type Commit struct {
	Prepare

	// Signatures of the Replicas that have signed this Commit
	Signatures
}

// Hash implements the Hasher interface.
func (commit *Commit) Hash() Hash {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, CommitHeader)
	binary.Write(&buf, binary.BigEndian, commit.Prepare.Hash())
	return sha3.Sum256(buf.Bytes())
}

func (commit *Commit) Fault() Fault {
	return Fault{
		Rank:   commit.Prepare.Block.Rank,
		Height: commit.Prepare.Block.Height,
	}
}

// Verify the Commit message. Returns an error if the message is invalid,
// otherwise nil.
func (commit *Commit) Verify(verifier Verifier) error {
	// TODO: Complete verification
	if err := commit.Prepare.Verify(verifier); err != nil {
		return err
	}
	return verifier.VerifySignatures(commit.Signatures)
}

func (commit *Commit) SetSignatures(signatures Signatures) {
	commit.Signatures = signatures
}

func (commit *Commit) GetSignatures() Signatures {
	return commit.Signatures
}

// ProcessCommits by collecting Commits. Once a threshold of Commits has been
// reached for a Block, the Block is certified and produced to the Block
// channel. The incrementing of height must be done by reading Blocks produced
// by this process, and comparing it to the current height.
func ProcessCommits(ctx context.Context, commitChIn <-chan Commit, signer Signer, verifier Verifier, capacity, threshold int) (<-chan Commit, <-chan Fault, <-chan error) {
	commitCh := make(chan Commit, capacity)
	faultCh := make(chan Fault, capacity)
	errCh := make(chan error, capacity)

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

			case commit, ok := <-commitChIn:
				if !ok {
					return
				}

				message, err := VerifyAndSignMessage(&commit, &store, signer, verifier, threshold)
				if err != nil {
					errCh <- err
					continue
				}
				// After verifying and signing the message check for Faults
				switch message := message.(type) {

				case Commit:
					select {
					case <-ctx.Done():
						errCh <- ctx.Err()
						return
					case commitCh <- message:
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
