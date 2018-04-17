package hyper

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

// Verify the Commit message. Returns an error if the message is invalid,
// otherwise nil.
func (commit *Commit) Verify() error {
	return nil
}

// Hash implements the Hasher interface.
func (commit *Commit) Hash() Hash {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, CommitHeader)
	binary.Write(&buf, binary.BigEndian, commit.Prepare.Hash())
	return sha3.Sum256(buf.Bytes())
}

// ProcessCommits by collecting Commits. Once a threshold of Commits has been
// reached for a Block, the Block is certified and produced to the Block
// channel. The incrementing of height must be done by reading Blocks produced
// by this process, and comparing it to the current height.
func ProcessCommits(ctx context.Context, commitChIn <-chan Commit, signer Signer, capacity int) (<-chan Commit, <-chan Fault, <-chan error) {
	commitCh := make(chan Commit, capacity)
	faultCh := make(chan Fault, capacity)
	errCh := make(chan error, capacity)

	go func() {
		defer close(commitCh)
		defer close(blockCh)
		defer close(faultCh)
		defer close(errCh)

		threshold := validator.Threshold()

		commits := map[[32]byte]Commit{}
		certifications := map[[32]byte]struct{}{}

		for {
			select {

			case <-ctx.Done():
				errCh <- ctx.Err()
				return

			case commit, ok := <-commitChIn:
				if !ok {
					return
				}

				res, err := Foo(commit, signer)
				if err != nil {
					errCh <- err
					continue
				}

				switch res.(type) {
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
					continue
				}

				hash := commit.Hash()
				if _, ok := certifications[hash]; ok {
					continue
				}

				if err := commit.Verify(); err != nil {
					fault, err := FaultFromCommit(&commit, signer)
					if err != nil {
						errCh <- err
						continue
					}
					select {
					case <-ctx.Done():
						errCh <- ctx.Err()
						return
					case faultCh <- fault:
					}
					continue
				}

				if _, ok := commits[hash]; !ok {
					signature, err := signer.Sign(commit)
					if err != nil {
						errCh <- err
						continue
					}
					commit.Signatures = commit.Signatures.Merge(Signatures{signature})
					commits[hash] = commit

					select {
					case <-ctx.Done():
						errCh <- ctx.Err()
						return
					case commitCh <- commits[hash]:
					}
				} else {
					commits[hash].Signatures = commits[hash].Signatures.Merge(commit.Signatures)
				}

				if len(commits[hash].Signatures) >= threshold {
					certifications[hash] = struct{}{}
					select {
					case <-ctx.Done():
						errCh <- ctx.Err()
						return
					case commitCh <- commits[hash]:
					}
				}
			}
		}
	}()

	return commitCh, faultCh, errCh
}
