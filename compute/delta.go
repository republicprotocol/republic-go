package compute

import (
	"bytes"
	"math/big"

	base58 "github.com/jbenet/go-base58"
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/order"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/go-sss"
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

	FstCodeShare   sss.Share
	SndCodeShare   sss.Share
	PriceShare     sss.Share
	MaxVolumeShare sss.Share
	MinVolumeShare sss.Share
}

func NewDeltaFragment(left *order.Fragment, right *order.Fragment, prime *big.Int) (*DeltaFragment, error) {
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

	deltaFragment := &DeltaFragment{
		ID:                  DeltaFragmentID(crypto.Keccak256([]byte(buyOrderFragment.ID), []byte(sellOrderFragment.ID))),
		DeltaID:             DeltaID(crypto.Keccak256([]byte(buyOrderFragment.OrderID), []byte(sellOrderFragment.OrderID))),
		BuyOrderID:          buyOrderFragment.OrderID,
		SellOrderID:         sellOrderFragment.OrderID,
		BuyOrderFragmentID:  buyOrderFragment.ID,
		SellOrderFragmentID: sellOrderFragment.ID,
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
		return nil, nil // Only return new deltas
	}
	if builder.hasDeltaFragment(deltaFragment.ID) {
		return nil, nil // Only return new deltas
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

	matrix.buyOrderFragments[buyOrderFragment.ID.String()] = buyOrderFragment
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

	matrix.sellOrderFragments[sellOrderFragment.ID.String()] = sellOrderFragment
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

// Sub two OrderFragments from one another and return the resulting output
// ResultFragment. The output ResultFragment will have its ID computed.
func (orderFragment *OrderFragment) Sub(other *OrderFragment, prime *big.Int) (*DeltaFragment, error) {
	// Check that the OrderFragments have compatible sss.Shares, and that one
	// of them is an OrderBuy and the other is an OrderSell.
	if err := orderFragment.IsCompatible(other); err != nil {
		return nil, err
	}

	// Label the OrderFragments appropriately.
	var buyOrderFragment, sellOrderFragment *OrderFragment
	if orderFragment.OrderParity == OrderParityBuy {
		buyOrderFragment = orderFragment
		sellOrderFragment = other
	} else {
		buyOrderFragment = other
		sellOrderFragment = orderFragment
	}

	// Perform the addition using the buyOrderFragment as the LHS and the
	// sellOrderFragment as the RHS.
	fstCodeShare := sss.Share{
		Key:   buyOrderFragment.FstCodeShare.Key,
		Value: big.NewInt(0).Add(buyOrderFragment.FstCodeShare.Value, big.NewInt(0).Sub(prime, sellOrderFragment.FstCodeShare.Value)),
	}
	sndCodeShare := sss.Share{
		Key:   buyOrderFragment.SndCodeShare.Key,
		Value: big.NewInt(0).Add(buyOrderFragment.SndCodeShare.Value, big.NewInt(0).Sub(prime, sellOrderFragment.SndCodeShare.Value)),
	}
	priceShare := sss.Share{
		Key:   buyOrderFragment.PriceShare.Key,
		Value: big.NewInt(0).Add(buyOrderFragment.PriceShare.Value, big.NewInt(0).Sub(prime, sellOrderFragment.PriceShare.Value)),
	}
	maxVolumeShare := sss.Share{
		Key:   buyOrderFragment.MaxVolumeShare.Key,
		Value: big.NewInt(0).Add(buyOrderFragment.MaxVolumeShare.Value, big.NewInt(0).Sub(prime, sellOrderFragment.MinVolumeShare.Value)),
	}
	minVolumeShare := sss.Share{
		Key:   buyOrderFragment.MinVolumeShare.Key,
		Value: big.NewInt(0).Add(sellOrderFragment.MaxVolumeShare.Value, big.NewInt(0).Sub(prime, buyOrderFragment.MinVolumeShare.Value)),
	}
	fstCodeShare.Value.Mod(fstCodeShare.Value, prime)
	sndCodeShare.Value.Mod(sndCodeShare.Value, prime)
	priceShare.Value.Mod(priceShare.Value, prime)
	maxVolumeShare.Value.Mod(maxVolumeShare.Value, prime)
	minVolumeShare.Value.Mod(minVolumeShare.Value, prime)
	deltaFragment := &DeltaFragment{
		nil,
		buyOrderFragment.OrderID,
		sellOrderFragment.OrderID,
		buyOrderFragment.ID,
		sellOrderFragment.ID,
		fstCodeShare,
		sndCodeShare,
		priceShare,
		maxVolumeShare,
		minVolumeShare,
	}
	deltaFragment.ID = DeltaFragmentID(crypto.Keccak256(deltaFragment.BuyOrderFragmentID[:], deltaFragment.SellOrderFragmentID[:]))
	return deltaFragment, nil
}
