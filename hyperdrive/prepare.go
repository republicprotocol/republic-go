package hyper

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
	Signatures []Signature
}

// Verify the Prepare message. Returns an error if the message is invalid,
// otherwise nil.
func (prepare *Prepare) Verify() error {
	return nil
}

// Hash implements the Hasher interface.
func (prepare *Prepare) Hash() Hash {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, PrepareHeader)
	binary.Write(&buf, binary.BigEndian, prepare.Proposal.Hash())
	return sha3.Sum256(buf.Bytes())
}

func ProcessPreparation(ctx context.Context, prepareChIn <-chan Prepare, validator Validator) (<-chan Commit, <-chan Fault, <-chan error) {
	threshold := validator.Threshold()
	commitCh := make(chan Commit, threshold)
	faultCh := make(chan Fault, threshold)
	errCh := make(chan error, threshold)
	prepares := map[[32]byte]int{}
	commited := map[[32]byte]bool{}
	counter := uint64(0)

	go func() {
		defer close(commitCh)
		defer close(faultCh)
		defer close(errCh)

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case prepare, ok := <-prepareChIn:
				// log.Println("Successfully reading prepares from channels")
				if !ok {
					return
				}
				counter++
				h := BlockHash(prepare.Block)
				if commited[h] {
					continue
				}

				if !validator.ValidatePrepare(prepare) {
					faultCh <- Fault{
						Rank:      prepare.Rank,
						Height:    prepare.Height,
						Signature: validator.Sign(),
					}
					continue
				}

				// log.Println("Validated prepare on", validator.Sign())
				// log.Println("Counting prepares on", validator.Sign(), prepares[h], "with threshold", threshold)
				if prepares[h] >= threshold-1 {
					// log.Println("Commited on ", validator.Sign())
					commitCh <- Commit{
						Rank:               prepare.Rank,
						Height:             prepare.Height,
						Block:              prepare.Block,
						ThresholdSignature: ThresholdSignature("Threshold_BLS"),
						Signature:          validator.Sign(),
					}
					commited[h] = true
				} else {
					prepares[h]++
				}
			}
		}
	}()

	return commitCh, faultCh, errCh
}
