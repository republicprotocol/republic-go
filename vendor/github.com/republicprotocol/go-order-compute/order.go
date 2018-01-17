package compute

import (
	"bytes"
	"encoding/binary"

	"github.com/ethereum/go-ethereum/crypto"
)

// An OrderType is a publicly bit of information that determines the type of
// trade that an Order is representing.
type OrderType int64

// The possible values for an OrderType.
const (
	OrderTypeIBBO  = 1
	OrderTypeLimit = 2
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

// An OrderID is the Keccak256 hash of an Order.
type OrderID []byte

// OrderIDs is an alias.
type OrderIDs []OrderID

// An Order represents the want to perform a trade of assets. Public data in
// the Order must be exposed for computation, but private data should not be
// exposed to anyone other than the trader that wants to execute the Order.
type Order struct {
	// Public data.
	ID   OrderID
	Type OrderType
	Buy  int64

	// Private data.
	LittleCode CurrencyCode
	BigCode    CurrencyCode
	Price      int64
	MaxVolume  int64
	MinVolume  int64
	Nonce      int64
}

// NewOrder returns a new Order and computes the OrderID for the Order.
func NewOrder(ty OrderType, buy int64, littleCode CurrencyCode, bigCode CurrencyCode, price int64, maxVolume int64, minVolume int64, nonce int64) *Order {
	order := &Order{
		Type: ty,
		Buy:  buy,

		LittleCode: littleCode,
		BigCode:    bigCode,
		Price:      price,
		MaxVolume:  maxVolume,
		MinVolume:  minVolume,
		Nonce:      nonce,
	}
	order.ID = OrderID(crypto.Keccak256(order.Bytes()))
	return order
}

// Bytes returns an Order serialized into a bytes.
func (order *Order) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, order.Type)
	binary.Write(buf, binary.LittleEndian, order.Buy)
	binary.Write(buf, binary.LittleEndian, order.LittleCode)
	binary.Write(buf, binary.LittleEndian, order.BigCode)
	binary.Write(buf, binary.LittleEndian, order.Price)
	binary.Write(buf, binary.LittleEndian, order.MaxVolume)
	binary.Write(buf, binary.LittleEndian, order.MinVolume)
	binary.Write(buf, binary.LittleEndian, order.Nonce)
	return buf.Bytes()
}
