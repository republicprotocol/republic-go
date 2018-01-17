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

// OrderIDs is an alias.
type OrderIDs []OrderID

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
