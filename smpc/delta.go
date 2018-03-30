package smpc

import (
	"bytes"
	"fmt"
	"math/big"
	"sync"

	"github.com/republicprotocol/go-do"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jbenet/go-base58"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/shamir"
)

type DeltaFragmentMatrix struct {
	do.GuardedObject

	prime *big.Int

	buyOrderFragments  map[string]*order.Fragment
	sellOrderFragments map[string]*order.Fragment
	deltaFragments     map[string]map[string]DeltaFragment

	deltaFragmentsQueue         DeltaFragments
	deltaFragmentsQueueNotEmpty *do.Guard
}

func NewDeltaFragmentMatrix(prime *big.Int) *DeltaFragmentMatrix {
	deltaFragmentMatrix := new(DeltaFragmentMatrix)
	deltaFragmentMatrix.GuardedObject = do.NewGuardedObject()
	deltaFragmentMatrix.prime = prime
	deltaFragmentMatrix.buyOrderFragments = map[string]*order.Fragment{}
	deltaFragmentMatrix.sellOrderFragments = map[string]*order.Fragment{}
	deltaFragmentMatrix.deltaFragments = map[string]map[string]DeltaFragment{}
	deltaFragmentMatrix.deltaFragmentsQueue = DeltaFragments{}
	deltaFragmentMatrix.deltaFragmentsQueueNotEmpty = deltaFragmentMatrix.Guard(func() bool {
		return len(deltaFragmentMatrix.deltaFragmentsQueue) > 0
	})
	return deltaFragmentMatrix
}

func (matrix *DeltaFragmentMatrix) ComputeBuyOrder(buyOrderFragment *order.Fragment) {
	matrix.Enter(nil)
	defer matrix.Exit()

	matrix.buyOrderFragments[string(buyOrderFragment.OrderID)] = buyOrderFragment

	for _, sellOrderFragment := range matrix.sellOrderFragments {
		if !buyOrderFragment.IsCompatible(sellOrderFragment) {
			continue
		}
		if _, ok := matrix.deltaFragments[string(buyOrderFragment.OrderID)]; !ok {
			matrix.deltaFragments[string(buyOrderFragment.OrderID)] = map[string]DeltaFragment{}
		}
		deltaFragment := NewDeltaFragment(buyOrderFragment, sellOrderFragment, matrix.prime)
		matrix.deltaFragments[string(buyOrderFragment.OrderID)][string(sellOrderFragment.OrderID)] = deltaFragment
		matrix.deltaFragmentsQueue = append(matrix.deltaFragmentsQueue, deltaFragment)
	}
}

func (matrix *DeltaFragmentMatrix) ComputeSellOrder(sellOrderFragment *order.Fragment) {
	matrix.Enter(nil)
	defer matrix.Exit()

	matrix.sellOrderFragments[string(sellOrderFragment.OrderID)] = sellOrderFragment

	for _, buyOrderFragment := range matrix.buyOrderFragments {
		if !buyOrderFragment.IsCompatible(sellOrderFragment) {
			continue
		}
		if _, ok := matrix.deltaFragments[string(buyOrderFragment.OrderID)]; !ok {
			matrix.deltaFragments[string(buyOrderFragment.OrderID)] = map[string]DeltaFragment{}
		}
		deltaFragment := NewDeltaFragment(buyOrderFragment, sellOrderFragment, matrix.prime)
		matrix.deltaFragments[string(buyOrderFragment.OrderID)][string(sellOrderFragment.OrderID)] = deltaFragment
		matrix.deltaFragmentsQueue = append(matrix.deltaFragmentsQueue, deltaFragment)
	}

}

func (matrix *DeltaFragmentMatrix) RemoveBuyOrder(id order.ID) {
	matrix.Enter(nil)
	defer matrix.Exit()

	delete(matrix.deltaFragments, string(id))
	delete(matrix.buyOrderFragments, string(id))
}

func (matrix *DeltaFragmentMatrix) RemoveSellOrder(id order.ID) {
	matrix.Enter(nil)
	defer matrix.Exit()

	delete(matrix.sellOrderFragments, string(id))
	for buyOrderID := range matrix.deltaFragments {
		delete(matrix.deltaFragments[buyOrderID], string(id))
	}
}

func (matrix *DeltaFragmentMatrix) WaitForDeltaFragments(deltaFragments DeltaFragments) int {
	matrix.Enter(matrix.deltaFragmentsQueueNotEmpty)
	defer matrix.Exit()

	n := 0
	for i := 0; i < len(deltaFragments) && i < len(matrix.deltaFragmentsQueue); i++ {
		deltaFragments[i] = matrix.deltaFragmentsQueue[i]
		n++
	}

	if n >= len(matrix.deltaFragmentsQueue) {
		matrix.deltaFragmentsQueue = matrix.deltaFragmentsQueue[0:0]
	} else {
		matrix.deltaFragmentsQueue = matrix.deltaFragmentsQueue[n:]
	}
	return n
}

type DeltaBuilder struct {
	do.GuardedObject

	k     int64
	prime *big.Int

	deltas                 map[string]Delta
	deltaFragments         map[string]DeltaFragment
	deltasToDeltaFragments map[string]DeltaFragments

	deltasQueue         Deltas
	deltasQueueNotEmpty *do.Guard
}

func NewDeltaBuilder(k int64, prime *big.Int) *DeltaBuilder {
	builder := new(DeltaBuilder)
	builder.GuardedObject = do.NewGuardedObject()
	builder.k = k
	builder.prime = prime
	builder.deltas = map[string]Delta{}
	builder.deltaFragments = map[string]DeltaFragment{}
	builder.deltasToDeltaFragments = map[string]DeltaFragments{}
	builder.deltasQueue = Deltas{}
	builder.deltasQueueNotEmpty = builder.Guard(func() bool {
		return len(builder.deltasQueue) > 0
	})
	return builder
}

func (builder *DeltaBuilder) ComputeDelta(deltaFragments DeltaFragments) {
	builder.Enter(nil)
	defer builder.Exit()

	for _, deltaFragment := range deltaFragments {
		// Store the DeltaFragment if it has not been seen before
		if builder.hasDeltaFragment(deltaFragment.ID) {
			continue
		}
		builder.deltaFragments[string(deltaFragment.ID)] = deltaFragment

		// Associate the DeltaFragment with its respective Delta if the Delta
		// has not been built yet
		if builder.hasDelta(deltaFragment.DeltaID) {
			continue
		}
		if _, ok := builder.deltasToDeltaFragments[string(deltaFragment.DeltaID)]; ok {
			builder.deltasToDeltaFragments[string(deltaFragment.DeltaID)] =
				append(builder.deltasToDeltaFragments[string(deltaFragment.DeltaID)], deltaFragment)
		} else {
			builder.deltasToDeltaFragments[string(deltaFragment.DeltaID)] = DeltaFragments{deltaFragment}
		}

		// Build the Delta
		deltaFragments := builder.deltasToDeltaFragments[string(deltaFragment.DeltaID)]
		if int64(len(deltaFragments)) >= builder.k {
			if !IsCompatible(deltaFragments) {
				continue
			}
			delta := NewDelta(deltaFragments, builder.prime)
			builder.deltas[string(delta.ID)] = delta
			builder.deltasQueue = append(builder.deltasQueue, delta)
		}
	}
}

func (builder *DeltaBuilder) WaitForDeltas(deltas Deltas) int {
	builder.Enter(builder.deltasQueueNotEmpty)
	defer builder.Exit()

	n := 0
	for i := 0; i < len(deltas) && i < len(builder.deltasQueue); i++ {
		deltas[i] = builder.deltasQueue[i]
		n++
	}

	if n >= len(builder.deltasQueue) {
		builder.deltasQueue = builder.deltasQueue[0:0]
	} else {
		builder.deltasQueue = builder.deltasQueue[n:]
	}
	return n
}

func (builder *DeltaBuilder) hasDelta(deltaID DeltaID) bool {
	_, ok := builder.deltas[string(deltaID)]
	return ok
}

func (builder *DeltaBuilder) hasDeltaFragment(deltaFragmentID DeltaFragmentID) bool {
	_, ok := builder.deltaFragments[string(deltaFragmentID)]
	return ok
}

type DeltaHandler func(deltas Deltas)

// A DeltaID is the Keccak256 hash of the order IDs that were used to compute
// the associated Delta.
type DeltaID []byte

// Equal returns an equality check between two DeltaIDs.
func (id DeltaID) Equal(other DeltaID) bool {
	return bytes.Equal(id, other)
}

// String returns a DeltaID as a Base58 encoded string.
func (id DeltaID) String() string {
	return base58.Encode(id)
}

type Deltas []Delta

type Delta struct {
	ID          DeltaID
	BuyOrderID  order.ID
	SellOrderID order.ID
	FstCode     *big.Int
	SndCode     *big.Int
	Price       *big.Int
	MaxVolume   *big.Int
	MinVolume   *big.Int
}

func NewDelta(deltaFragments DeltaFragments, prime *big.Int) Delta {

	// Collect Shares across all DeltaFragments.
	k := len(deltaFragments)
	fstCodeShares := make(shamir.Shares, k)
	sndCodeShares := make(shamir.Shares, k)
	priceShares := make(shamir.Shares, k)
	maxVolumeShares := make(shamir.Shares, k)
	minVolumeShares := make(shamir.Shares, k)
	for i, deltaFragment := range deltaFragments {
		fstCodeShares[i] = deltaFragment.FstCodeShare
		sndCodeShares[i] = deltaFragment.SndCodeShare
		priceShares[i] = deltaFragment.PriceShare
		maxVolumeShares[i] = deltaFragment.MaxVolumeShare
		minVolumeShares[i] = deltaFragment.MinVolumeShare
	}

	// Join the Shares into a Result.
	delta := Delta{
		BuyOrderID:  deltaFragments[0].BuyOrderID,
		SellOrderID: deltaFragments[0].SellOrderID,
	}
	delta.FstCode = shamir.Join(prime, fstCodeShares)
	delta.SndCode = shamir.Join(prime, sndCodeShares)
	delta.Price = shamir.Join(prime, priceShares)
	delta.MaxVolume = shamir.Join(prime, maxVolumeShares)
	delta.MinVolume = shamir.Join(prime, minVolumeShares)

	// Compute the ResultID and return the Result.
	delta.ID = DeltaID(crypto.Keccak256(delta.BuyOrderID[:], delta.SellOrderID[:]))
	return delta
}

func (delta *Delta) IsMatch(prime *big.Int) bool {
	zeroThreshold := big.NewInt(0).Div(prime, big.NewInt(2))
	if delta.FstCode.Cmp(big.NewInt(0)) != 0 {
		return false
	}
	if delta.SndCode.Cmp(big.NewInt(0)) != 0 {
		return false
	}
	if delta.Price.Cmp(zeroThreshold) == 1 {
		return false
	}
	if delta.MaxVolume.Cmp(zeroThreshold) == 1 {
		return false
	}
	if delta.MinVolume.Cmp(zeroThreshold) == 1 {
		return false
	}
	return true
}

// A DeltaFragmentID is the Keccak256 hash of the order IDs that were used to
// compute the associated DeltaFragment.
type DeltaFragmentID []byte

// Equal returns an equality check between two DeltaFragmentIDs.
func (id DeltaFragmentID) Equal(other DeltaFragmentID) bool {
	return bytes.Equal(id, other)
}

// String returns a DeltaFragmentID as a Base58 encoded string.
func (id DeltaFragmentID) String() string {
	return base58.Encode(id)
}

type DeltaFragments []DeltaFragment

// A DeltaFragment is a secret share of a Final. Is is performing a
// computation over two OrderFragments.
type DeltaFragment struct {
	ID                  DeltaFragmentID
	DeltaID             DeltaID
	BuyOrderID          order.ID
	SellOrderID         order.ID
	BuyOrderFragmentID  order.FragmentID
	SellOrderFragmentID order.FragmentID

	FstCodeShare   shamir.Share
	SndCodeShare   shamir.Share
	PriceShare     shamir.Share
	MaxVolumeShare shamir.Share
	MinVolumeShare shamir.Share
}

func NewDeltaFragment(left *order.Fragment, right *order.Fragment, prime *big.Int) DeltaFragment {
	var buyOrderFragment *order.Fragment
	var sellOrderFragment *order.Fragment
	if left.OrderParity == order.ParityBuy {
		buyOrderFragment = left
		sellOrderFragment = right
	} else {
		buyOrderFragment = right
		sellOrderFragment = left
	}

	fstCodeShare := shamir.Share{
		Key:   buyOrderFragment.FstCodeShare.Key,
		Value: big.NewInt(0).Add(buyOrderFragment.FstCodeShare.Value, big.NewInt(0).Sub(prime, sellOrderFragment.FstCodeShare.Value)),
	}
	sndCodeShare := shamir.Share{
		Key:   buyOrderFragment.SndCodeShare.Key,
		Value: big.NewInt(0).Add(buyOrderFragment.SndCodeShare.Value, big.NewInt(0).Sub(prime, sellOrderFragment.SndCodeShare.Value)),
	}
	priceShare := shamir.Share{
		Key:   buyOrderFragment.PriceShare.Key,
		Value: big.NewInt(0).Add(buyOrderFragment.PriceShare.Value, big.NewInt(0).Sub(prime, sellOrderFragment.PriceShare.Value)),
	}
	maxVolumeShare := shamir.Share{
		Key:   buyOrderFragment.MaxVolumeShare.Key,
		Value: big.NewInt(0).Add(buyOrderFragment.MaxVolumeShare.Value, big.NewInt(0).Sub(prime, sellOrderFragment.MinVolumeShare.Value)),
	}
	minVolumeShare := shamir.Share{
		Key:   buyOrderFragment.MinVolumeShare.Key,
		Value: big.NewInt(0).Add(sellOrderFragment.MaxVolumeShare.Value, big.NewInt(0).Sub(prime, buyOrderFragment.MinVolumeShare.Value)),
	}
	fstCodeShare.Value.Mod(fstCodeShare.Value, prime)
	sndCodeShare.Value.Mod(sndCodeShare.Value, prime)
	priceShare.Value.Mod(priceShare.Value, prime)
	maxVolumeShare.Value.Mod(maxVolumeShare.Value, prime)
	minVolumeShare.Value.Mod(minVolumeShare.Value, prime)

	return DeltaFragment{
		ID:                  DeltaFragmentID(crypto.Keccak256([]byte(buyOrderFragment.ID), []byte(sellOrderFragment.ID))),
		DeltaID:             DeltaID(crypto.Keccak256([]byte(buyOrderFragment.OrderID), []byte(sellOrderFragment.OrderID))),
		BuyOrderID:          buyOrderFragment.OrderID,
		SellOrderID:         sellOrderFragment.OrderID,
		BuyOrderFragmentID:  buyOrderFragment.ID,
		SellOrderFragmentID: sellOrderFragment.ID,
		FstCodeShare:        fstCodeShare,
		SndCodeShare:        sndCodeShare,
		PriceShare:          priceShare,
		MaxVolumeShare:      maxVolumeShare,
		MinVolumeShare:      minVolumeShare,
	}
}

// Equals checks if two DeltaFragments are equal in value.
func (deltaFragment *DeltaFragment) Equals(other *DeltaFragment) bool {
	return deltaFragment.ID.Equal(other.ID) &&
		deltaFragment.DeltaID.Equal(other.DeltaID) &&
		deltaFragment.BuyOrderID.Equal(other.BuyOrderID) &&
		deltaFragment.SellOrderID.Equal(other.SellOrderID) &&
		deltaFragment.BuyOrderFragmentID.Equal(other.BuyOrderFragmentID) &&
		deltaFragment.SellOrderFragmentID.Equal(other.SellOrderFragmentID) &&
		deltaFragment.FstCodeShare.Key == other.FstCodeShare.Key &&
		deltaFragment.FstCodeShare.Value.Cmp(other.FstCodeShare.Value) == 0 &&
		deltaFragment.SndCodeShare.Key == other.SndCodeShare.Key &&
		deltaFragment.SndCodeShare.Value.Cmp(other.SndCodeShare.Value) == 0 &&
		deltaFragment.PriceShare.Key == other.PriceShare.Key &&
		deltaFragment.PriceShare.Value.Cmp(other.PriceShare.Value) == 0 &&
		deltaFragment.MaxVolumeShare.Key == other.MaxVolumeShare.Key &&
		deltaFragment.MaxVolumeShare.Value.Cmp(other.MaxVolumeShare.Value) == 0 &&
		deltaFragment.MinVolumeShare.Key == other.MinVolumeShare.Key &&
		deltaFragment.MinVolumeShare.Value.Cmp(other.MinVolumeShare.Value) == 0
}

// IsCompatible returns true if all DeltaFragments are fragments of the same
// Delta, otherwise it returns false.
func IsCompatible(deltaFragments DeltaFragments) bool {
	if len(deltaFragments) == 0 {
		return false
	}
	for i := range deltaFragments {
		if !deltaFragments[i].DeltaID.Equal(deltaFragments[0].DeltaID) {
			return false
		}
	}
	return true
}

// DeltaQueues is a slice of DeltaQueue components.
type DeltaQueues []DeltaQueue

// A DeltaQueue owns a channel of Delta components.
type DeltaQueue struct {
	chMu   *sync.RWMutex
	chOpen bool
	ch     chan Delta
}

// NewDeltaQueue returns a MessageQueue interface that channels Delta
//components.
func NewDeltaQueue(messageQueueLimit int) DeltaQueue {
	return DeltaQueue{
		chMu:   new(sync.RWMutex),
		chOpen: true,
		ch:     make(chan Delta, messageQueueLimit),
	}
}

// Run the DeltaQueue. The DeltaQueue is an abstraction over a channel of Delta
// components and does not need to be run. This method does nothing.
func (queue *DeltaQueue) Run() error {
	return nil
}

// Shutdown the DeltaQueue. If it has already been Shutdown, an error will be
// returned.
func (queue *DeltaQueue) Shutdown() error {
	queue.chMu.Lock()
	defer queue.chMu.Unlock()

	queue.chOpen = false
	close(queue.ch)
	return nil
}

// Send a message to the DeltaQueue. The Message must be a Delta component,
// otherwise an error is returned.
func (queue *DeltaQueue) Send(message dispatch.Message) error {
	queue.chMu.RLock()
	defer queue.chMu.RUnlock()

	if !queue.chOpen {
		return nil
	}

	switch message := message.(type) {
	case Delta:
		queue.ch <- message
	default:
		return fmt.Errorf("cannot send message: unrecognized type %T", message)
	}
	return nil
}

// Recv a message from the DeltaQueue. All Messages returned will be Delta
// components.
func (queue *DeltaQueue) Recv() (dispatch.Message, bool) {
	message, ok := <-queue.ch
	return message, ok
}
