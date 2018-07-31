package smpc

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"sync"

	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/shamir"
)

// MaxJoinLength restricts the maximum number of shamir.Shares that can be
// stored in a Join, and therefore the maximum number of values that can be
// reconstructed by a JoinSet.
const MaxJoinLength = 16

// ErrJoinLengthUnequal is returned when two Joins with the same JoinID have a
// different number of shamir.Shares.
var ErrJoinLengthUnequal = errors.New("join length unequal")

// ErrJoinLengthExceedsMax is returned when a Join has too many shamir.Shares
// compared to the MaxJoinLength.
var ErrJoinLengthExceedsMax = errors.New("join length exceeds max")

// A JoinID is used to identify an SMPC join over a network of SMPC nodes. For
// a value to be joined, all SMPC nodes must use the same JoinID for that
// value.
type JoinID [33]byte

// A Join is used to join a set of shamir.Shares. The shamir.Shares within a
// Join are all associated with different shared values. All shamir.Shares must
// have the same index value.
type Join struct {
	ID     JoinID
	Index  JoinIndex
	Shares shamir.Shares
}

// MarshalBinary implements the encoding.Marshaler interface.
func (join *Join) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, join.ID); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, join.Index); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, int64(len(join.Shares))); err != nil {
		return nil, err
	}
	for _, share := range join.Shares {
		shareData, err := share.MarshalBinary()
		if err != nil {
			return nil, err
		}
		if err := binary.Write(buf, binary.BigEndian, shareData); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.Unmarshaler interface.
func (join *Join) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, &join.ID); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &join.Index); err != nil {
		return err
	}
	numShares := int64(0)
	if err := binary.Read(buf, binary.BigEndian, &numShares); err != nil {
		return err
	}
	join.Shares = make(shamir.Shares, numShares)
	for i := int64(0); i < numShares; i++ {
		shareData := [16]byte{}
		if _, err := buf.Read(shareData[:]); err != nil {
			return err
		}
		if err := join.Shares[i].UnmarshalBinary(shareData[:]); err != nil {
			return err
		}
	}
	return nil
}

// JoinIndex is the index of all shamir.Shares in a Join.
type JoinIndex uint64

// Callback for a successful reconstruction of a JoinSet. Callbacks are not
// guaranteed to be called mutually exclusively, and must ensure their own
// concurrent safety.
type Callback func(JoinID, []uint64)

// JoinSet is a set of Joins, with different JoinIndices, but with the same
// JoinID.
type JoinSet struct {

	// Set of Joins with different JoinIndices.
	Set map[JoinIndex]Join

	// Values to be reconstructed once enough Joins have been added to the Set
	// and a boolean to reflect whether or not the reconstruction has already
	// happened.
	Values    [MaxJoinLength]uint64
	ValuesOk  bool
	ValuesLen int

	// Callback for when the reconstruction happens.
	Callback Callback
}

// A Joiner received Joins and groups them together based on their JoinID. Once
// a sufficient number of Joins have been collected, the shamir.Shares are
// zipped across all Joins, and each zip is reconstructed into a value.
type Joiner struct {
	k     int64
	cache shamir.Shares

	joinSetsMu *sync.Mutex
	joinSets   map[JoinID]JoinSet
}

// NewJoiner returns an empty Joiner that needs k shamir.Shares before it can
// reconstruct a value.
func NewJoiner(k int64) *Joiner {
	return &Joiner{
		k:     k,
		cache: make(shamir.Shares, k),

		joinSetsMu: new(sync.Mutex),
		joinSets:   map[JoinID]JoinSet{},
	}
}

// InsertJoinAndSetCallback for a JoinID. If a Callback has been set for this
// JoinID it will be replaced. The Callback will only be called once per call
// to Joiner.InsertJoinAndSetCallback; it will be automatically set to nil
// after it is called. Passing a nil Callback will remove any existing
// Callback.
func (joiner *Joiner) InsertJoinAndSetCallback(join Join, callback Callback) error {
	return joiner.insertJoin(join, callback, true)
}

// InsertJoin for a JoinID. If a Callback has been set for this JoinID it will
// be called if the insertion results in a successful reconstruction of values.
func (joiner *Joiner) InsertJoin(join Join) error {
	return joiner.insertJoin(join, nil, false)
}

func (joiner *Joiner) insertJoin(join Join, callback Callback, overrideCallback bool) error {
	if len(join.Shares) > MaxJoinLength {
		return ErrJoinLengthExceedsMax
	}

	maybeCallback := Callback(nil)
	maybeValues := [MaxJoinLength]uint64{}
	maybeValuesLen := 0

	err := func() error {
		joiner.joinSetsMu.Lock()
		defer joiner.joinSetsMu.Unlock()

		// Load the JoinSet and store any mutations when this function returns
		joinSet, ok := joiner.joinSets[join.ID]
		defer func() {
			joiner.joinSets[join.ID] = joinSet
		}()

		// Initialize the JoinSet for this JoinID if it has not been initialized
		if !ok {
			joinSet = JoinSet{
				Set:       map[JoinIndex]Join{},
				Values:    [MaxJoinLength]uint64{},
				ValuesLen: len(join.Shares),
			}
		}

		// Insert this join, if it is needed, and set the callback
		if len(join.Shares) != joinSet.ValuesLen {
			logger.Error(fmt.Sprintf("%v: expected %v, got %v", ErrJoinLengthUnequal, joinSet.ValuesLen, len(join.Shares)))
			return ErrJoinLengthUnequal
		}
		if !joinSet.ValuesOk {
			joinSet.Set[join.Index] = join
		}
		if overrideCallback {
			joinSet.Callback = callback
		}

		// Short circuit if there are not enough Joins to successfully perform a
		// reconstruction
		if int64(len(joinSet.Set)) < joiner.k {
			return nil
		}

		// If the reconstruction has not happened, perform the reconstruction
		if !joinSet.ValuesOk {
			for i := 0; i < joinSet.ValuesLen; i++ {
				k := int64(0)
				for _, join := range joinSet.Set {
					joiner.cache[k] = join.Shares[i]
					k++
					if k >= joiner.k {
						break
					}
				}
				joinSet.Values[i] = shamir.Join(joiner.cache)
			}
			joinSet.ValuesOk = true
		}

		// Copy values to ensure that future mutations do not interfere
		// with the callback (which happens outside of the mutex to
		// encourage liveness)
		maybeCallback = joinSet.Callback
		maybeValues = joinSet.Values
		maybeValuesLen = joinSet.ValuesLen

		// Ensure that the Callback is only ever called once
		joinSet.Callback = nil

		return nil
	}()
	if err != nil {
		return err
	}

	if maybeCallback != nil {
		maybeCallback(join.ID, maybeValues[:maybeValuesLen])
		return nil
	}
	return nil
}
