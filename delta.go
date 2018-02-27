package compute

import (
	"bytes"
	"math/big"

	"github.com/republicprotocol/go-do"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/go-sss"
)

type DeltaShard struct {
	Signature      []byte
	DeltaFragments []*DeltaFragment
}

func NewDeltaShard(deltaFragments []*DeltaFragment) DeltaShard {
	return DeltaShard{
		DeltaFragments: deltaFragments,
	}
}

// A DeltaID is the Keccak256 hash of the OrderIDs that were used to compute
// the respective Result.
type DeltaID []byte

// Equals checks if two ResultIDs are equal in value.
func (id DeltaID) Equals(other DeltaID) bool {
	return bytes.Equal(id, other)
}

// String returns the ResultID as a string.
func (id DeltaID) String() string {
	return string(id)
}

// A Delta is the publicly computed value of comparing two Orders.
type Delta struct {
	ID          DeltaID
	BuyOrderID  OrderID
	SellOrderID OrderID
	FstCode     *big.Int
	SndCode     *big.Int
	Price       *big.Int
	MaxVolume   *big.Int
	MinVolume   *big.Int
}

func NewDelta(deltaFragments []*DeltaFragment, prime *big.Int) (*Delta, error) {
	// Check that all ResultFragments are compatible with each other.
	err := isCompatible(deltaFragments)
	if err != nil {
		return nil, err
	}

	// Collect sss.Shares across all DeltaFragments.
	k := len(deltaFragments)
	fstCodeShares := make(sss.Shares, k)
	sndCodeShares := make(sss.Shares, k)
	priceShares := make(sss.Shares, k)
	maxVolumeShares := make(sss.Shares, k)
	minVolumeShares := make(sss.Shares, k)
	for i, resultFragment := range deltaFragments {
		fstCodeShares[i] = resultFragment.FstCodeShare
		sndCodeShares[i] = resultFragment.SndCodeShare
		priceShares[i] = resultFragment.PriceShare
		maxVolumeShares[i] = resultFragment.MaxVolumeShare
		minVolumeShares[i] = resultFragment.MinVolumeShare
	}

	// Join the sss.Shares into a Result.
	delta := &Delta{
		BuyOrderID:  deltaFragments[0].BuyOrderID,
		SellOrderID: deltaFragments[0].SellOrderID,
	}
	delta.FstCode = sss.Join(prime, fstCodeShares)
	delta.SndCode = sss.Join(prime, sndCodeShares)
	delta.Price = sss.Join(prime, priceShares)
	delta.MaxVolume = sss.Join(prime, maxVolumeShares)
	delta.MinVolume = sss.Join(prime, minVolumeShares)

	// Compute the ResultID and return the Result.
	delta.ID = DeltaID(crypto.Keccak256(delta.BuyOrderID[:], delta.SellOrderID[:]))
	return delta, nil
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

// A DeltaFragmentID is the Keccak256 hash of its OrderFragmentIDs.
type DeltaFragmentID []byte

// Equals checks if two DeltaFragmentIDs are equal in value.
func (id DeltaFragmentID) Equals(other DeltaFragmentID) bool {
	return bytes.Equal(id, other)
}

// String returns the DeltaFragmentID as a string.
func (id DeltaFragmentID) String() string {
	return string(id)
}

// A DeltaFragment is a secret share of a Final. Is is performing a
// computation over two OrderFragments.
type DeltaFragment struct {
	// Public data.
	ID                  DeltaFragmentID
	BuyOrderID          OrderID
	SellOrderID         OrderID
	BuyOrderFragmentID  OrderFragmentID
	SellOrderFragmentID OrderFragmentID

	// Private data.
	FstCodeShare   sss.Share
	SndCodeShare   sss.Share
	PriceShare     sss.Share
	MaxVolumeShare sss.Share
	MinVolumeShare sss.Share
}

func NewDeltaFragment(left *OrderFragment, right *OrderFragment, prime *big.Int) (*DeltaFragment, error) {
	if err := left.IsCompatible(right); err != nil {
		return nil, err
	}
	var buyOrderFragment *OrderFragment
	var sellOrderFragment *OrderFragment
	if left.OrderParity == OrderParityBuy {
		buyOrderFragment = left
		sellOrderFragment = right
	} else {
		buyOrderFragment = right
		sellOrderFragment = left
	}

	deltaFragment, err := buyOrderFragment.Sub(sellOrderFragment, prime)
	if err != nil {
		return nil, err
	}
	return deltaFragment, nil
}

// Equals checks if two DeltaFragments are equal in value.
func (deltaFragment *DeltaFragment) Equals(other *DeltaFragment) bool {
	return deltaFragment.ID.Equals(other.ID) &&
		deltaFragment.BuyOrderID.Equals(other.BuyOrderID) &&
		deltaFragment.SellOrderID.Equals(other.SellOrderID) &&
		deltaFragment.BuyOrderFragmentID.Equals(other.BuyOrderFragmentID) &&
		deltaFragment.SellOrderFragmentID.Equals(other.SellOrderFragmentID) &&
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

// DeltaID returns the ID of the Delta to which this DeltaFragment will
// eventually reconstruct.
func (deltaFragment *DeltaFragment) DeltaID() DeltaID {
	return DeltaID(crypto.Keccak256(deltaFragment.BuyOrderID[:], deltaFragment.SellOrderID[:]))
}

// IsCompatible returns an error when the two deltaFragments do not have
// the same share indices.
func isCompatible(deltaFragments []*DeltaFragment) error {
	if len(deltaFragments) == 0 {
		return NewEmptySliceError("result fragments")
	}
	buyOrderID := deltaFragments[0].BuyOrderID
	sellOrderID := deltaFragments[0].SellOrderID
	for i := range deltaFragments {
		if !deltaFragments[i].BuyOrderID.Equals(buyOrderID) {
			return NewOrderFragmentationError(0, int64(i))
		}
		if !deltaFragments[i].SellOrderID.Equals(sellOrderID) {
			return NewOrderFragmentationError(0, int64(i))
		}
	}
	return nil
}

type DeltaBuilder struct {
	do.GuardedObject

	k                      int64
	prime                  *big.Int
	deltas                 map[string]*Delta
	deltaFragments         map[string]*DeltaFragment
	deltasToDeltaFragments map[string][]*DeltaFragment
}

func NewDeltaBuilder(k int64, prime *big.Int) *DeltaBuilder {
	return &DeltaBuilder{
		GuardedObject:          do.NewGuardedObject(),
		k:                      k,
		prime:                  prime,
		deltas:                 map[string]*Delta{},
		deltaFragments:         map[string]*DeltaFragment{},
		deltasToDeltaFragments: map[string][]*DeltaFragment{},
	}
}

func (builder *DeltaBuilder) InsertDeltaFragment(deltaFragment *DeltaFragment) (*Delta, error) {
	builder.Enter(nil)
	defer builder.Exit()
	return builder.insertDeltaFragment(deltaFragment)
}

func (builder *DeltaBuilder) insertDeltaFragment(deltaFragment *DeltaFragment) (*Delta, error) {
	// Is the delta already built, or are we adding a delta fragment that we
	// have already seen
	if builder.hasDelta(deltaFragment.DeltaID()) {
		return builder.deltas[deltaFragment.DeltaID().String()], nil
	}
	if builder.hasDeltaFragment(deltaFragment.ID) {
		return builder.deltas[deltaFragment.DeltaID().String()], nil
	}

	// Add the delta fragment to the builder and attach it to the appropriate
	// delta
	builder.deltaFragments[deltaFragment.ID.String()] = deltaFragment
	if deltaFragments, ok := builder.deltasToDeltaFragments[deltaFragment.DeltaID().String()]; ok {
		deltaFragments = append(deltaFragments, deltaFragment)
		builder.deltasToDeltaFragments[deltaFragment.DeltaID().String()] = deltaFragments
	} else {
		builder.deltasToDeltaFragments[deltaFragment.DeltaID().String()] = []*DeltaFragment{deltaFragment}
	}

	// Build the delta if possible and return it
	deltaFragments := builder.deltasToDeltaFragments[deltaFragment.DeltaID().String()]
	if int64(len(deltaFragments)) >= builder.k {
		delta, err := NewDelta(deltaFragments, builder.prime)
		if err != nil {
			return delta, err
		}
		builder.deltas[delta.ID.String()] = delta
		return delta, nil
	}

	return nil, nil
}

func (builder *DeltaBuilder) HasDelta(deltaID DeltaID) bool {
	builder.EnterReadOnly(nil)
	defer builder.Exit()
	return builder.hasDelta(deltaID)
}

func (builder *DeltaBuilder) hasDelta(deltaID DeltaID) bool {
	_, ok := builder.deltas[deltaID.String()]
	return ok
}

func (builder *DeltaBuilder) HasDeltaFragment(deltaFragmentID DeltaFragmentID) bool {
	builder.EnterReadOnly(nil)
	defer builder.Exit()
	return builder.hasDeltaFragment(deltaFragmentID)
}

func (builder *DeltaBuilder) hasDeltaFragment(deltaFragmentID DeltaFragmentID) bool {
	_, ok := builder.deltaFragments[deltaFragmentID.String()]
	return ok
}

type DeltaFragmentMatrix struct {
	do.GuardedObject

	prime                  *big.Int
	buyOrderFragments      map[string]*OrderFragment
	sellOrderFragments     map[string]*OrderFragment
	buySellDeltaFragments  map[string]map[string]*DeltaFragment
	completeOrderFragments map[string]*OrderFragment
}

func NewDeltaFragmentMatrix(prime *big.Int) *DeltaFragmentMatrix {
	return &DeltaFragmentMatrix{
		GuardedObject:          do.NewGuardedObject(),
		prime:                  prime,
		buyOrderFragments:      map[string]*OrderFragment{},
		sellOrderFragments:     map[string]*OrderFragment{},
		buySellDeltaFragments:  map[string]map[string]*DeltaFragment{},
		completeOrderFragments: map[string]*OrderFragment{},
	}
}

func (matrix *DeltaFragmentMatrix) InsertOrderFragment(orderFragment *OrderFragment) ([]*DeltaFragment, error) {
	matrix.Enter(nil)
	defer matrix.Exit()
	if orderFragment.OrderParity == OrderParityBuy {
		return matrix.insertBuyOrderFragment(orderFragment)
	}
	return matrix.insertSellOrderFragment(orderFragment)
}

func (matrix *DeltaFragmentMatrix) insertBuyOrderFragment(buyOrderFragment *OrderFragment) ([]*DeltaFragment, error) {
	if _, ok := matrix.buyOrderFragments[buyOrderFragment.ID.String()]; ok {
		return []*DeltaFragment{}, nil
	}
	if _, ok := matrix.completeOrderFragments[buyOrderFragment.ID.String()]; ok {
		return []*DeltaFragment{}, nil
	}

	deltaFragments := make([]*DeltaFragment, 0, len(matrix.sellOrderFragments))
	deltaFragmentsMap := map[string]*DeltaFragment{}
	for i := range matrix.sellOrderFragments {
		deltaFragment, err := buyOrderFragment.Sub(matrix.sellOrderFragments[i], matrix.prime)
		if err != nil {
			return deltaFragments, err
		}
		deltaFragments = append(deltaFragments, deltaFragment)
		deltaFragmentsMap[matrix.sellOrderFragments[i].ID.String()] = deltaFragment
	}

	matrix.buySellDeltaFragments[buyOrderFragment.ID.String()] = deltaFragmentsMap
	return deltaFragments, nil
}

func (matrix *DeltaFragmentMatrix) insertSellOrderFragment(sellOrderFragment *OrderFragment) ([]*DeltaFragment, error) {
	if _, ok := matrix.sellOrderFragments[sellOrderFragment.ID.String()]; ok {
		return []*DeltaFragment{}, nil
	}
	if _, ok := matrix.completeOrderFragments[sellOrderFragment.ID.String()]; ok {
		return []*DeltaFragment{}, nil
	}

	deltaFragments := make([]*DeltaFragment, 0, len(matrix.buyOrderFragments))
	for i := range matrix.buyOrderFragments {
		deltaFragment, err := matrix.buyOrderFragments[i].Sub(sellOrderFragment, matrix.prime)
		if err != nil {
			return deltaFragments, err
		}
		if _, ok := matrix.buySellDeltaFragments[matrix.buyOrderFragments[i].ID.String()]; ok {
			deltaFragments = append(deltaFragments, deltaFragment)
			matrix.buySellDeltaFragments[matrix.buyOrderFragments[i].ID.String()][sellOrderFragment.ID.String()] = deltaFragment
		}
	}
	return deltaFragments, nil
}

func (matrix *DeltaFragmentMatrix) RemoveOrderFragment(orderFragment *OrderFragment) error {
	matrix.Enter(nil)
	defer matrix.Exit()
	if orderFragment.OrderParity == OrderParityBuy {
		return matrix.removeBuyOrderFragment(orderFragment)
	}
	return matrix.removeSellOrderFragment(orderFragment)
}

func (matrix *DeltaFragmentMatrix) removeBuyOrderFragment(buyOrderFragment *OrderFragment) error {
	if _, ok := matrix.buyOrderFragments[buyOrderFragment.ID.String()]; !ok {
		return nil
	}

	delete(matrix.buyOrderFragments, buyOrderFragment.ID.String())
	delete(matrix.buySellDeltaFragments, buyOrderFragment.ID.String())

	matrix.completeOrderFragments[buyOrderFragment.ID.String()] = buyOrderFragment
	return nil
}

func (matrix *DeltaFragmentMatrix) removeSellOrderFragment(sellOrderFragment *OrderFragment) error {
	if _, ok := matrix.sellOrderFragments[sellOrderFragment.ID.String()]; !ok {
		return nil
	}

	for i := range matrix.buySellDeltaFragments {
		delete(matrix.buySellDeltaFragments[i], sellOrderFragment.ID.String())
	}

	matrix.completeOrderFragments[sellOrderFragment.ID.String()] = sellOrderFragment
	return nil
}

func (matrix *DeltaFragmentMatrix) DeltaFragment(buyOrderFragmentID, sellOrderFragmentID OrderFragmentID) *DeltaFragment {
	matrix.EnterReadOnly(nil)
	defer matrix.ExitReadOnly()
	return matrix.deltaFragment(buyOrderFragmentID, sellOrderFragmentID)
}

func (matrix *DeltaFragmentMatrix) deltaFragment(buyOrderFragmentID, sellOrderFragmentID OrderFragmentID) *DeltaFragment {
	if deltaFragments, ok := matrix.buySellDeltaFragments[buyOrderFragmentID.String()]; ok {
		if deltaFragment, ok := deltaFragments[sellOrderFragmentID.String()]; ok {
			return deltaFragment
		}
	}
	return nil
}
