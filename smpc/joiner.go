package smpc

import (
	"errors"
	"sync"

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
type JoinID [32]byte

// A Join is used to join a set of shamir.Shares. The shamir.Shares within a
// Join are all associated with different shared values. All shamir.Shares must
// have the same index value.
type Join struct {
	ID     JoinID
	Index  JoinIndex
	Shares shamir.Shares
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
func NewJoiner(k int64) Joiner {
	return Joiner{
		k:     k,
		cache: make(shamir.Shares, k),

		joinSetsMu: new(sync.Mutex),
		joinSets:   map[JoinID]Joins{},
	}
}

// InsertJoinAndSetCallback for a JoinID. Previously callbacks for the same
// JoinID will be replaced.
func (joiner *Joiner) InsertJoinAndSetCallback(join Join, callback func(id JoinID, values []uint64)) error {
	if len(join.Shares) > MaxJoinLength {
		return ErrJoinLengthExceedsMax
	}

	maybeCallback := Callback(nil)
	maybeValues := [MaxJoinLength]uint64{}
	maybeValuesLen := 0

	err := func() error {
		joiner.joinSetsMu.Lock()
		defer joiner.joinSetsMu.Unlock()

		// Initialize the JoinSet for this JoinID if it has not been initialized
		if _, ok := joiner.joinSets[join.ID]; !ok {
			joiner.joinSets[join.ID] = JoinSet{
				Set:       map[JoinIndex]Join{},
				Values:    [MaxJoinLength]uint64{},
				ValuesLen: len(join.Shares),
			}
		}

		// Insert this join, if it is needed, and set the callback
		if len(join.Shares) != joiner.joinSets[join.ID].ValuesLen {
			return ErrJoinLengthUnequal
		}
		if !joiner.joins[join.ID].ValueOk {
			joiner.joins[join.ID].Set[join.Index] = join
		}
		joiner.joins[join.ID].Callback = callback

		// Short circuit if there are not enough Joins to successfully perform a
		// reconstruction
		if len(joiner.joins[join.ID].Set) < joiner.k {
			return nil
		}

		// If the reconstruction has not happened, perform the reconstruction
		if !joiner.joins[join.ID].ValueOk {
			for i := 0; i < joiner.joins[join.ID].ValuesLen; i++ {
				k := 0
				for _, join := range joiner.joins[join.ID].Set {
					joiner.cache[k] = join.Shares[i]
					k++
					if k >= joiner.k {
						break
					}
				}
				joiner.joins[join.ID].Values[i] = shamir.Join(joiner.kCache)
			}
			joiner.joins[join.ID].ValuesOk = true
		}

		// Copy values to ensure that future mutations do not interfere
		// with the callback (which happens outside of the mutex to
		// encourage liveness)
		maybeCallback = joiner.joins[join.ID].Callback
		maybeValues = joiner.joins[join.ID].Values
		maybeValuesLen = joiner.joins[join.ID].ValuesLen
		return nil
	}()
	if err != nil {
		return err
	}

	if maybeCallback != nil {
		maybeCallback(join.ID, maybeValues[:maybeValuesLen])
	}
	return nil
}
