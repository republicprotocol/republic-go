package compute

import (
	"bytes"
	"encoding/binary"
	"math/big"

	"github.com/republicprotocol/go-sss"

	"github.com/ethereum/go-ethereum/crypto"
)

// A CurrencyCode is a numerical representation of the currencies supported by
// Orders.
type CurrencyCode int64

// The possible values for a CurrencyCode.
const (
	CurrencyCodeBTC = 1
	CurrencyCodeETH = 2
	CurrencyCodeREN = 3
)

// An OrderType is a publicly bit of information that determines the type of
// trade that an Order is representing.
type OrderType int64

// The possible values for an OrderType.
const (
	OrderTypeIBBO  = 1
	OrderTypeLimit = 2
)

type OrderBuySell int64

const (
	OrderBuy  = 1
	OrderSell = 0
)

// An OrderID is the Keccak256 hash of an Order.
type OrderID []byte

// An Order represents the want to perform a trade of assets. Public data in
// the Order must be exposed for computation, but private data should not be
// exposed to anyone other than the trader that wants to execute the Order.
type Order struct {
	// Public data.
	ID      OrderID
	Type    OrderType
	BuySell OrderBuySell

	// Private data.
	FstCode   CurrencyCode
	SndCode   CurrencyCode
	Price     int64
	MaxVolume int64
	MinVolume int64
	Nonce     int64
}

// NewOrder returns a new Order and computes the OrderID for the Order.
func NewOrder(ty OrderType, buySell OrderBuySell, fstCode CurrencyCode, sndCode CurrencyCode, price int64, maxVolume int64, minVolume int64, nonce int64) *Order {
	order := &Order{
		Type:    ty,
		BuySell: buySell,

		FstCode:   fstCode,
		SndCode:   sndCode,
		Price:     price,
		MaxVolume: maxVolume,
		MinVolume: minVolume,
		Nonce:     nonce,
	}
	order.ID = OrderID(crypto.Keccak256(order.Bytes()))
	return order
}

// Bytes returns an Order serialized into a bytes.
func (order *Order) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, order.Type)
	binary.Write(buf, binary.LittleEndian, order.BuySell)
	binary.Write(buf, binary.LittleEndian, order.FstCode)
	binary.Write(buf, binary.LittleEndian, order.SndCode)
	binary.Write(buf, binary.LittleEndian, order.Price)
	binary.Write(buf, binary.LittleEndian, order.MaxVolume)
	binary.Write(buf, binary.LittleEndian, order.MinVolume)
	binary.Write(buf, binary.LittleEndian, order.Nonce)
	return buf.Bytes()
}

func (order *Order) Split(n, k int64, prime *big.Int) ([]*OrderFragment, error) {
	fstCodeShares, err := sss.Split(n, k, prime, big.NewInt(int64(order.FstCode)))
	if err != nil {
		return nil, err
	}
	sndCodeShares, err := sss.Split(n, k, prime, big.NewInt(int64(order.SndCode)))
	if err != nil {
		return nil, err
	}
	priceShares, err := sss.Split(n, k, prime, big.NewInt(order.Price))
	if err != nil {
		return nil, err
	}
	maxVolumeShares, err := sss.Split(n, k, prime, big.NewInt(order.MaxVolume))
	if err != nil {
		return nil, err
	}
	minVolumeShares, err := sss.Split(n, k, prime, big.NewInt(order.MinVolume))
	if err != nil {
		return nil, err
	}
	orderFragments := make([]*OrderFragment, n)
	for i := int64(0); i < n; i++ {
		orderFragments[i] = NewOrderFragment(
			order.ID,
			order.Type,
			order.BuySell,
			fstCodeShares[i],
			sndCodeShares[i],
			priceShares[i],
			maxVolumeShares[i],
			minVolumeShares[i],
		)
	}
	return orderFragments, nil
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
	// Public data.
	ID           OrderFragmentID
	OrderID      OrderID
	OrderType    OrderType
	OrderBuySell OrderBuySell

	// Private data.
	FstCodeShare   sss.Share
	SndCodeShare   sss.Share
	PriceShare     sss.Share
	MaxVolumeShare sss.Share
	MinVolumeShare sss.Share
}

// NewOrderFragment returns a new OrderFragment and computes the
// OrderFragmentID for the OrderFragment.
func NewOrderFragment(orderID OrderID, orderType OrderType, orderBuySell OrderBuySell, fstCodeShare, sndCodeShare, priceShare, maxVolumeShare, minVolumeShare sss.Share) *OrderFragment {
	orderFragment := &OrderFragment{
		OrderID:        orderID,
		OrderType:      orderType,
		OrderBuySell:   orderBuySell,
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
	if orderFragment.OrderBuySell == OrderBuy {
		buyOrderFragment = orderFragment
		sellOrderFragment = other
	} else {
		buyOrderFragment = other
		sellOrderFragment = orderFragment
	}

	// Perform the addition using the buyOrderFragment as the LHS and the
	// sellOrderFragment as the RHS.
	computed := &ComputedOrderFragment{
		BuyOrderID:          buyOrderFragment.OrderID,
		BuyOrderFragmentID:  buyOrderFragment.ID,
		SellOrderID:         sellOrderFragment.OrderID,
		SellOrderFragmentID: sellOrderFragment.ID,

		FstCodeShare: sss.Share{
			Key:   buyOrderFragment.FstCodeShare.Key,
			Value: big.NewInt(0).Add(buyOrderFragment.FstCodeShare.Value, sellOrderFragment.FstCodeShare.Value),
		},
		SndCodeShare: sss.Share{
			Key:   buyOrderFragment.SndCodeShare.Key,
			Value: big.NewInt(0).Add(buyOrderFragment.SndCodeShare.Value, sellOrderFragment.SndCodeShare.Value),
		},
		PriceShare: sss.Share{
			Key:   buyOrderFragment.PriceShare.Key,
			Value: big.NewInt(0).Add(buyOrderFragment.PriceShare.Value, sellOrderFragment.PriceShare.Value),
		},
		MaxVolumeShare: sss.Share{
			Key:   buyOrderFragment.MaxVolumeShare.Key,
			Value: big.NewInt(0).Add(buyOrderFragment.MaxVolumeShare.Value, sellOrderFragment.MaxVolumeShare.Value),
		},
		MinVolumeShare: sss.Share{
			Key:   buyOrderFragment.MinVolumeShare.Key,
			Value: big.NewInt(0).Add(sellOrderFragment.MinVolumeShare.Value, buyOrderFragment.MinVolumeShare.Value),
		},
	}
	computed.FstCodeShare.Value.Mod(computed.FstCodeShare.Value, prime)
	computed.SndCodeShare.Value.Mod(computed.SndCodeShare.Value, prime)
	computed.PriceShare.Value.Mod(computed.PriceShare.Value, prime)
	computed.MaxVolumeShare.Value.Mod(computed.MaxVolumeShare.Value, prime)
	computed.MinVolumeShare.Value.Mod(computed.MinVolumeShare.Value, prime)
	computed.ID = OrderFragmentID(crypto.Keccak256(computed.Bytes()))
	return computed, nil
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
	if orderFragment.OrderBuySell == OrderBuy {
		buyOrderFragment = orderFragment
		sellOrderFragment = other
	} else {
		buyOrderFragment = other
		sellOrderFragment = orderFragment
	}

	// Perform the addition using the buyOrderFragment as the LHS and the
	// sellOrderFragment as the RHS.
	computed := &ComputedOrderFragment{
		BuyOrderID:          buyOrderFragment.OrderID,
		BuyOrderFragmentID:  buyOrderFragment.ID,
		SellOrderID:         sellOrderFragment.OrderID,
		SellOrderFragmentID: sellOrderFragment.ID,

		FstCodeShare: sss.Share{
			Key:   buyOrderFragment.FstCodeShare.Key,
			Value: big.NewInt(0).Add(buyOrderFragment.FstCodeShare.Value, big.NewInt(0).Set(prime).Sub(prime, sellOrderFragment.FstCodeShare.Value)),
		},
		SndCodeShare: sss.Share{
			Key:   buyOrderFragment.SndCodeShare.Key,
			Value: big.NewInt(0).Add(buyOrderFragment.SndCodeShare.Value, big.NewInt(0).Set(prime).Sub(prime, sellOrderFragment.SndCodeShare.Value)),
		},
		PriceShare: sss.Share{
			Key:   buyOrderFragment.PriceShare.Key,
			Value: big.NewInt(0).Add(buyOrderFragment.PriceShare.Value, big.NewInt(0).Set(prime).Sub(prime, sellOrderFragment.PriceShare.Value)),
		},
		MaxVolumeShare: sss.Share{
			Key:   buyOrderFragment.MaxVolumeShare.Key,
			Value: big.NewInt(0).Add(buyOrderFragment.MaxVolumeShare.Value, big.NewInt(0).Set(prime).Sub(prime, sellOrderFragment.MaxVolumeShare.Value)),
		},
		MinVolumeShare: sss.Share{
			Key:   buyOrderFragment.MinVolumeShare.Key,
			Value: big.NewInt(0).Add(sellOrderFragment.MinVolumeShare.Value, big.NewInt(0).Set(prime).Sub(prime, buyOrderFragment.MinVolumeShare.Value)),
		},
	}
	computed.FstCodeShare.Value.Mod(computed.FstCodeShare.Value, prime)
	computed.SndCodeShare.Value.Mod(computed.SndCodeShare.Value, prime)
	computed.PriceShare.Value.Mod(computed.PriceShare.Value, prime)
	computed.MaxVolumeShare.Value.Mod(computed.MaxVolumeShare.Value, prime)
	computed.MinVolumeShare.Value.Mod(computed.MinVolumeShare.Value, prime)
	computed.ID = OrderFragmentID(crypto.Keccak256(computed.Bytes()))
	return computed, nil
}

// Bytes returns an Order serialized into a bytes.
func (orderFragment *OrderFragment) Bytes() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, orderFragment.OrderID)
	binary.Write(buf, binary.LittleEndian, orderFragment.OrderType)
	binary.Write(buf, binary.LittleEndian, orderFragment.OrderBuySell)

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
	if orderFragment.OrderBuySell == rhs.OrderBuySell {
		return NewOrderComputationError(orderFragment.OrderBuySell)
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
