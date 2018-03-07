package compute

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	base58 "github.com/jbenet/go-base58"
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/shamir"
)

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

func NewDelta(deltaFragments []*DeltaFragment, prime *big.Int) *Delta {
	// Check that all ResultFragments are compatible with each other.
	if !IsCompatible(deltaFragments) {
		return nil
	}

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
	delta := &Delta{
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

func NewDeltaFragment(left *order.Fragment, right *order.Fragment, prime *big.Int) *DeltaFragment {
	if !left.IsCompatible(right) {
		return nil
	}

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

	return &DeltaFragment{
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
func IsCompatible(deltaFragments []*DeltaFragment) bool {
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

func (builder *DeltaBuilder) InsertDeltaFragment(deltaFragment *DeltaFragment) *Delta {
	builder.Enter(nil)
	defer builder.Exit()
	return builder.insertDeltaFragment(deltaFragment)
}

func (builder *DeltaBuilder) insertDeltaFragment(deltaFragment *DeltaFragment) *Delta {
	// Is the delta already built, or are we adding a delta fragment that we
	// have already seen
	if builder.hasDelta(deltaFragment.DeltaID) {
		return nil // Only return new deltas
	}
	if builder.hasDeltaFragment(deltaFragment.ID) {
		return nil // Only return new deltas
	}

	// Add the delta fragment to the builder and attach it to the appropriate
	// delta
	builder.deltaFragments[string(deltaFragment.ID)] = deltaFragment
	if deltaFragments, ok := builder.deltasToDeltaFragments[string(deltaFragment.DeltaID)]; ok {
		deltaFragments = append(deltaFragments, deltaFragment)
		builder.deltasToDeltaFragments[string(deltaFragment.DeltaID)] = deltaFragments
	} else {
		builder.deltasToDeltaFragments[string(deltaFragment.DeltaID)] = []*DeltaFragment{deltaFragment}
	}

	// Build the delta if possible and return it
	deltaFragments := builder.deltasToDeltaFragments[string(deltaFragment.DeltaID)]
	if int64(len(deltaFragments)) >= builder.k {
		delta := NewDelta(deltaFragments, builder.prime)
		if delta == nil {
			return nil
		}
		builder.deltas[string(delta.ID)] = delta
		return delta
	}

	return nil
}

func (builder *DeltaBuilder) HasDelta(deltaID DeltaID) bool {
	builder.EnterReadOnly(nil)
	defer builder.Exit()
	return builder.hasDelta(deltaID)
}

func (builder *DeltaBuilder) hasDelta(deltaID DeltaID) bool {
	_, ok := builder.deltas[string(deltaID)]
	return ok
}

func (builder *DeltaBuilder) HasDeltaFragment(deltaFragmentID DeltaFragmentID) bool {
	builder.EnterReadOnly(nil)
	defer builder.Exit()
	return builder.hasDeltaFragment(deltaFragmentID)
}

func (builder *DeltaBuilder) hasDeltaFragment(deltaFragmentID DeltaFragmentID) bool {
	_, ok := builder.deltaFragments[string(deltaFragmentID)]
	return ok
}

type DeltaFragmentMatrix struct {
	do.GuardedObject

	prime                  *big.Int
	buyOrderFragments      map[string]*order.Fragment
	sellOrderFragments     map[string]*order.Fragment
	buySellDeltaFragments  map[string]map[string]*DeltaFragment
	completeOrderFragments map[string]*order.Fragment
}

func NewDeltaFragmentMatrix(prime *big.Int) *DeltaFragmentMatrix {
	return &DeltaFragmentMatrix{
		GuardedObject:          do.NewGuardedObject(),
		prime:                  prime,
		buyOrderFragments:      map[string]*order.Fragment{},
		sellOrderFragments:     map[string]*order.Fragment{},
		buySellDeltaFragments:  map[string]map[string]*DeltaFragment{},
		completeOrderFragments: map[string]*order.Fragment{},
	}
}

func (matrix *DeltaFragmentMatrix) InsertOrderFragment(orderFragment *order.Fragment) ([]*DeltaFragment, error) {
	matrix.Enter(nil)
	defer matrix.Exit()
	if orderFragment.OrderParity == order.ParityBuy {
		return matrix.insertBuyOrderFragment(orderFragment)
	}
	return matrix.insertSellOrderFragment(orderFragment)
}

func (matrix *DeltaFragmentMatrix) insertBuyOrderFragment(buyOrderFragment *order.Fragment) ([]*DeltaFragment, error) {
	if _, ok := matrix.buyOrderFragments[string(buyOrderFragment.ID)]; ok {
		return []*DeltaFragment{}, nil
	}
	if _, ok := matrix.completeOrderFragments[string(buyOrderFragment.ID)]; ok {
		return []*DeltaFragment{}, nil
	}

	deltaFragments := make([]*DeltaFragment, 0, len(matrix.sellOrderFragments))
	deltaFragmentsMap := map[string]*DeltaFragment{}
	for i := range matrix.sellOrderFragments {
		deltaFragment := NewDeltaFragment(buyOrderFragment, matrix.sellOrderFragments[i], matrix.prime)
		if deltaFragment == nil {
			continue
		}
		deltaFragments = append(deltaFragments, deltaFragment)
		deltaFragmentsMap[string(matrix.sellOrderFragments[i].ID)] = deltaFragment
	}

	matrix.buyOrderFragments[string(buyOrderFragment.ID)] = buyOrderFragment
	matrix.buySellDeltaFragments[string(buyOrderFragment.ID)] = deltaFragmentsMap
	return deltaFragments, nil
}

func (matrix *DeltaFragmentMatrix) insertSellOrderFragment(sellOrderFragment *order.Fragment) ([]*DeltaFragment, error) {
	if _, ok := matrix.sellOrderFragments[string(sellOrderFragment.ID)]; ok {
		return []*DeltaFragment{}, nil
	}
	if _, ok := matrix.completeOrderFragments[string(sellOrderFragment.ID)]; ok {
		return []*DeltaFragment{}, nil
	}

	deltaFragments := make([]*DeltaFragment, 0, len(matrix.buyOrderFragments))
	for i := range matrix.buyOrderFragments {
		deltaFragment := NewDeltaFragment(matrix.buyOrderFragments[i], sellOrderFragment, matrix.prime)
		if deltaFragment == nil {
			continue
		}
		if _, ok := matrix.buySellDeltaFragments[string(matrix.buyOrderFragments[i].ID)]; ok {
			deltaFragments = append(deltaFragments, deltaFragment)
			matrix.buySellDeltaFragments[string(matrix.buyOrderFragments[i].ID)][string(sellOrderFragment.ID)] = deltaFragment
		}
	}

	matrix.sellOrderFragments[string(sellOrderFragment.ID)] = sellOrderFragment
	return deltaFragments, nil
}

func (matrix *DeltaFragmentMatrix) RemoveOrderFragment(orderFragment *order.Fragment) error {
	matrix.Enter(nil)
	defer matrix.Exit()
	if orderFragment.OrderParity == order.ParityBuy {
		return matrix.removeBuyOrderFragment(orderFragment)
	}
	return matrix.removeSellOrderFragment(orderFragment)
}

func (matrix *DeltaFragmentMatrix) removeBuyOrderFragment(buyOrderFragment *order.Fragment) error {
	if _, ok := matrix.buyOrderFragments[string(buyOrderFragment.ID)]; !ok {
		return nil
	}

	delete(matrix.buyOrderFragments, string(buyOrderFragment.ID))
	delete(matrix.buySellDeltaFragments, string(buyOrderFragment.ID))

	matrix.completeOrderFragments[string(buyOrderFragment.ID)] = buyOrderFragment
	return nil
}

func (matrix *DeltaFragmentMatrix) removeSellOrderFragment(sellOrderFragment *order.Fragment) error {
	if _, ok := matrix.sellOrderFragments[string(sellOrderFragment.ID)]; !ok {
		return nil
	}

	for i := range matrix.buySellDeltaFragments {
		delete(matrix.buySellDeltaFragments[i], string(sellOrderFragment.ID))
	}

	matrix.completeOrderFragments[string(sellOrderFragment.ID)] = sellOrderFragment
	return nil
}

func (matrix *DeltaFragmentMatrix) DeltaFragment(buyOrderFragmentID, sellOrderFragmentID order.FragmentID) *DeltaFragment {
	matrix.EnterReadOnly(nil)
	defer matrix.ExitReadOnly()
	return matrix.deltaFragment(buyOrderFragmentID, sellOrderFragmentID)
}

func (matrix *DeltaFragmentMatrix) deltaFragment(buyOrderFragmentID, sellOrderFragmentID order.FragmentID) *DeltaFragment {
	if deltaFragments, ok := matrix.buySellDeltaFragments[string(buyOrderFragmentID)]; ok {
		if deltaFragment, ok := deltaFragments[string(sellOrderFragmentID)]; ok {
			return deltaFragment
		}
	}
	return nil
}
