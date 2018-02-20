package compute

import (
	"bytes"
	"encoding/binary"
	"math/big"
	"time"

	"github.com/republicprotocol/go-identity"

	"github.com/republicprotocol/go-sss"

	"github.com/ethereum/go-ethereum/crypto"
)

// A CurrencyCode is a numerical representation of the currencies supported by
// Orders.
type CurrencyCode int64

// The possible values for a CurrencyCode.
const (
	CurrencyCodeBTC CurrencyCode = 1
	CurrencyCodeETH CurrencyCode = 2
	CurrencyCodeREN CurrencyCode = 3
	CurrencyCodeDGD CurrencyCode = 4
)

// An OrderType is a publicly bit of information that determines the type of
// trade that an Order is representing.
type OrderType int64

// The possible values for an OrderType.
const (
	OrderTypeIBBO  OrderType = 1
	OrderTypeLimit OrderType = 2
)

// OrderParity determines whether an Order is a buy or a sell.
type OrderParity int64

// The possible values for OrderParity.
const (
	OrderParityBuy  OrderParity = 1
	OrderParitySell OrderParity = 2
)

// An OrderID is the Keccak256 hash of an Order.
type OrderID []byte

// Equals checks if two OrderIDs are equal in value.
func (id OrderID) Equals(other OrderID) bool {
	return bytes.Equal(id, other)
}

// String returns the OrderID as a string.
func (id OrderID) String() string {
	return string(id)
}

// An Order represents the want to perform a trade of assets. Public data in
// the Order must be exposed for computation, but private data should not be
// exposed to anyone other than the trader that wants to execute the Order.
type Order struct {
	// Signature
	Owner     identity.ID `json:"owner"`
	Signature []byte      `json:"signature"`

	// Public
	ID     OrderID     `json:"id"`
	Type   OrderType   `json:"type"`
	Parity OrderParity `json:"parity"`
	Expiry time.Time   `json:"expiry"`

	// Secure
	FstCode   CurrencyCode `json:"fstCode"`
	SndCode   CurrencyCode `json:"sndCode"`
	Price     *big.Int     `json:"price"`
	MaxVolume *big.Int     `json:"maxVolume"`
	MinVolume *big.Int     `json:"minVolume"`

	// Private
	Nonce *big.Int `json:"nonce"`
}

// NewOrder returns a new Order and computes the OrderID for the Order.
func NewOrder(ty OrderType, parity OrderParity, expiry time.Time, fstCode, sndCode CurrencyCode, price, maxVolume, minVolume, nonce *big.Int) *Order {
	order := &Order{
		Type:      ty,
		Parity:    parity,
		Expiry:    expiry,
		FstCode:   fstCode,
		SndCode:   sndCode,
		Price:     price,
		MaxVolume: maxVolume,
		MinVolume: minVolume,
		Nonce:     nonce,
	}
	order.GenerateID()
	return order
}

// Split the Order into n OrderFragments, where k OrderFragments are needed to
// reconstruct the Order. Returns a slice of all n OrderFragments, or an error.
func (order *Order) Split(n, k int64, prime *big.Int) ([]*OrderFragment, error) {
	fstCodeShares, err := sss.Split(n, k, prime, big.NewInt(int64(order.FstCode)))
	if err != nil {
		return nil, err
	}
	sndCodeShares, err := sss.Split(n, k, prime, big.NewInt(int64(order.SndCode)))
	if err != nil {
		return nil, err
	}
	priceShares, err := sss.Split(n, k, prime, order.Price)
	if err != nil {
		return nil, err
	}
	maxVolumeShares, err := sss.Split(n, k, prime, order.MaxVolume)
	if err != nil {
		return nil, err
	}
	minVolumeShares, err := sss.Split(n, k, prime, order.MinVolume)
	if err != nil {
		return nil, err
	}
	orderFragments := make([]*OrderFragment, n)
	for i := int64(0); i < n; i++ {
		orderFragments[i] = NewOrderFragment(
			order.ID,
			order.Type,
			order.Parity,
			fstCodeShares[i],
			sndCodeShares[i],
			priceShares[i],
			maxVolumeShares[i],
			minVolumeShares[i],
		)
	}
	return orderFragments, nil
}

// Bytes returns an Order serialized into a bytes.
func (order *Order) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, order.Type)
	binary.Write(buf, binary.LittleEndian, order.Parity)
	binary.Write(buf, binary.LittleEndian, order.Expiry)
	binary.Write(buf, binary.LittleEndian, order.FstCode)
	binary.Write(buf, binary.LittleEndian, order.SndCode)
	binary.Write(buf, binary.LittleEndian, order.Price.Bytes())
	binary.Write(buf, binary.LittleEndian, order.MaxVolume.Bytes())
	binary.Write(buf, binary.LittleEndian, order.MinVolume.Bytes())
	binary.Write(buf, binary.LittleEndian, order.Nonce.Bytes())
	return buf.Bytes()
}

// Equals checks if two Orders are equal in value.
func (order *Order) Equals(other *Order) bool {
	return order.ID.Equals(other.ID) &&
		order.Type == other.Type &&
		order.Parity == other.Parity &&
		order.Expiry.Equal(other.Expiry) &&
		order.FstCode == other.FstCode &&
		order.SndCode == other.SndCode &&
		order.Price.Cmp(other.Price) == 0 &&
		order.MaxVolume.Cmp(other.MaxVolume) == 0 &&
		order.MinVolume.Cmp(other.MinVolume) == 0 &&
		order.Nonce.Cmp(other.Nonce) == 0
}

// GenerateID of the Order.
func (order *Order) GenerateID() {
	order.ID = OrderID(crypto.Keccak256(order.Bytes()))
}

// An OrderFragmentID is the Keccak256 hash of an OrderFragment.
type OrderFragmentID []byte

// Equals checks if two OrderFragmentIDs are equal in value.
func (id OrderFragmentID) Equals(other OrderFragmentID) bool {
	return bytes.Equal(id, other)
}

// An OrderFragment is a secret share of an Order. Is is created using Shamir
// secret sharing where the secret is an Order encoded as a big.Int.
type OrderFragment struct {
	// Signature
	Owner     identity.ID
	Signature []byte

	// Public
	ID          OrderFragmentID
	OrderID     OrderID
	OrderType   OrderType
	OrderParity OrderParity
	OrderExpiry time.Time

	// Secure
	FstCodeShare   sss.Share
	SndCodeShare   sss.Share
	PriceShare     sss.Share
	MaxVolumeShare sss.Share
	MinVolumeShare sss.Share
}

// NewOrderFragment returns a new OrderFragment and computes the
// OrderFragmentID for the OrderFragment.
func NewOrderFragment(orderID OrderID, orderType OrderType, orderParity OrderParity, fstCodeShare, sndCodeShare, priceShare, maxVolumeShare, minVolumeShare sss.Share) *OrderFragment {
	orderFragment := &OrderFragment{
		OrderID:        orderID,
		OrderType:      orderType,
		OrderParity:    orderParity,
		FstCodeShare:   fstCodeShare,
		SndCodeShare:   sndCodeShare,
		PriceShare:     priceShare,
		MaxVolumeShare: maxVolumeShare,
		MinVolumeShare: minVolumeShare,
	}
	orderFragment.ID = OrderFragmentID(crypto.Keccak256(orderFragment.Bytes()))
	return orderFragment
}

// Add two OrderFragments together and return the resulting output
// OrderFragment. The output OrderFragment will have its ID computed.
func (orderFragment *OrderFragment) Add(other *OrderFragment, prime *big.Int) (*ResultFragment, error) {
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
		Value: big.NewInt(0).Add(buyOrderFragment.FstCodeShare.Value, sellOrderFragment.FstCodeShare.Value),
	}
	sndCodeShare := sss.Share{
		Key:   buyOrderFragment.SndCodeShare.Key,
		Value: big.NewInt(0).Add(buyOrderFragment.SndCodeShare.Value, sellOrderFragment.SndCodeShare.Value),
	}
	priceShare := sss.Share{
		Key:   buyOrderFragment.PriceShare.Key,
		Value: big.NewInt(0).Add(buyOrderFragment.PriceShare.Value, sellOrderFragment.PriceShare.Value),
	}
	maxVolumeShare := sss.Share{
		Key:   buyOrderFragment.MaxVolumeShare.Key,
		Value: big.NewInt(0).Add(buyOrderFragment.MaxVolumeShare.Value, sellOrderFragment.MaxVolumeShare.Value),
	}
	minVolumeShare := sss.Share{
		Key:   buyOrderFragment.MinVolumeShare.Key,
		Value: big.NewInt(0).Add(sellOrderFragment.MinVolumeShare.Value, buyOrderFragment.MinVolumeShare.Value),
	}
	fstCodeShare.Value.Mod(fstCodeShare.Value, prime)
	sndCodeShare.Value.Mod(sndCodeShare.Value, prime)
	priceShare.Value.Mod(priceShare.Value, prime)
	maxVolumeShare.Value.Mod(maxVolumeShare.Value, prime)
	minVolumeShare.Value.Mod(minVolumeShare.Value, prime)
	resultFragment := NewResultFragment(
		buyOrderFragment.OrderID,
		sellOrderFragment.OrderID,
		buyOrderFragment.ID,
		sellOrderFragment.ID,
		fstCodeShare,
		sndCodeShare,
		priceShare,
		maxVolumeShare,
		minVolumeShare,
	)
	return resultFragment, nil
}

// Sub two OrderFragments from one another and return the resulting output
// OrderFragment. The output OrderFragment will have its ID computed.
func (orderFragment *OrderFragment) Sub(other *OrderFragment, prime *big.Int) (*ResultFragment, error) {
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
	resultFragment := NewResultFragment(
		buyOrderFragment.OrderID,
		sellOrderFragment.OrderID,
		buyOrderFragment.ID,
		sellOrderFragment.ID,
		fstCodeShare,
		sndCodeShare,
		priceShare,
		maxVolumeShare,
		minVolumeShare,
	)
	return resultFragment, nil
}

// Bytes returns an Order serialized into a bytes.
func (orderFragment *OrderFragment) Bytes() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, orderFragment.OrderID)
	binary.Write(buf, binary.LittleEndian, orderFragment.OrderType)
	binary.Write(buf, binary.LittleEndian, orderFragment.OrderParity)

	binary.Write(buf, binary.LittleEndian, orderFragment.FstCodeShare.Key)
	binary.Write(buf, binary.LittleEndian, orderFragment.FstCodeShare.Value.Bytes())
	binary.Write(buf, binary.LittleEndian, orderFragment.SndCodeShare.Key)
	binary.Write(buf, binary.LittleEndian, orderFragment.SndCodeShare.Value.Bytes())
	binary.Write(buf, binary.LittleEndian, orderFragment.PriceShare.Key)
	binary.Write(buf, binary.LittleEndian, orderFragment.PriceShare.Value.Bytes())
	binary.Write(buf, binary.LittleEndian, orderFragment.MaxVolumeShare.Key)
	binary.Write(buf, binary.LittleEndian, orderFragment.MaxVolumeShare.Value.Bytes())
	binary.Write(buf, binary.LittleEndian, orderFragment.MinVolumeShare.Key)
	binary.Write(buf, binary.LittleEndian, orderFragment.MinVolumeShare.Value.Bytes())

	return buf.Bytes()
}

// IsCompatible returns an error when the two OrderFragments do not have
// the same share indices.
func (orderFragment *OrderFragment) IsCompatible(rhs *OrderFragment) error {
	if orderFragment.OrderParity == rhs.OrderParity {
		return NewOrderParityError(orderFragment.OrderParity)
	}
	if orderFragment.FstCodeShare.Key != rhs.FstCodeShare.Key {
		return NewOrderFragmentationError(orderFragment.FstCodeShare.Key, rhs.FstCodeShare.Key)
	}
	if orderFragment.SndCodeShare.Key != rhs.SndCodeShare.Key {
		return NewOrderFragmentationError(orderFragment.SndCodeShare.Key, rhs.SndCodeShare.Key)
	}
	if orderFragment.PriceShare.Key != rhs.PriceShare.Key {
		return NewOrderFragmentationError(orderFragment.PriceShare.Key, rhs.PriceShare.Key)
	}
	if orderFragment.MaxVolumeShare.Key != rhs.MaxVolumeShare.Key {
		return NewOrderFragmentationError(orderFragment.MaxVolumeShare.Key, rhs.MaxVolumeShare.Key)
	}
	if orderFragment.MinVolumeShare.Key != rhs.MinVolumeShare.Key {
		return NewOrderFragmentationError(orderFragment.MinVolumeShare.Key, rhs.MinVolumeShare.Key)
	}
	return nil
}
