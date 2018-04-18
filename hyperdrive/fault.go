package hyperdrive

import (
	"bytes"
	"context"
	"encoding/binary"

	"golang.org/x/crypto/sha3"
)

func ProcessFault(ctx context.Context, faultChIn chan Fault, signer Signer, verifier Verifier, capacity int) (chan Fault, chan error) {
	faultCh := make(chan Fault, capacity)
	errCh := make(chan error, capacity)

	go func() {
		defer close(faultCh)
		defer close(errCh)

		store := NewMessageMapStore()

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case fault := <-faultChIn:
				message, err := VerifyAndSignMessage(&fault, &store, signer, verifier, 0)
				if err != nil {
					errCh <- err
					continue
				}
				// After verifying and signing the message check for Faults
				switch message := message.(type) {
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

	return faultCh, errCh
}

const FaultHeader = byte(5)

type Fault struct {
	Rank
	Height

	// Signatures of the Replicas that signed this Fault
	Signatures
}

// Hash implements the Hasher interface.
func (fault *Fault) Hash() Hash {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, FaultHeader)
	binary.Write(&buf, binary.BigEndian, fault.Rank)
	binary.Write(&buf, binary.BigEndian, fault.Height)

	return sha3.Sum256(buf.Bytes())
}

func (fault *Fault) Fault() *Fault {
	return fault
}

func (fault *Fault) Verify(verifier Verifier) error {
	return nil
}

func (fault *Fault) SetSignatures(signatures Signatures) {
	fault.Signatures = signatures
}

func (fault *Fault) GetSignatures() Signatures {
	return fault.Signatures
}
