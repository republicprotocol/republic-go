package order

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"os"
	"time"

	"github.com/jbenet/go-base58"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/shamir"
	"github.com/republicprotocol/republic-go/stackint"
)

// A CurrencyCode is a numerical representation of the currencies supported by
// Republic Protocol.
type CurrencyCode = int64

// CurrencyCode values.
const (
	CurrencyCodeBTC CurrencyCode = 1
	CurrencyCodeETH CurrencyCode = 2
	CurrencyCodeREN CurrencyCode = 3
	CurrencyCodeDGD CurrencyCode = 4
)

// A Type is a publicly bit of information that determines the type of
// trade that an Order is representing.
type Type int64

// Type values.
const (
	TypeIBBO  Type = 1
	TypeLimit Type = 2
)

// The Parity of an Order determines whether it is buy or a sell.
type Parity int64

// Parity values.
const (
	ParityBuy  Parity = 1
	ParitySell Parity = 2
)

// An ID is the Keccak256 hash of an Order.
type ID []byte

// Equal returns an equality check between two IDs.
func (id ID) Equal(other ID) bool {
	return bytes.Equal(id, other)
}

// String returns an ID as a Base58 encoded string.
func (id ID) String() string {
	return base58.Encode(id)
}

func (id ID) Hash() []byte {
	return crypto.Keccak256([]byte(id.String()))
}

// The Status shows what status the order is in.
type Status = uint8

// Status values.
const (
	Nil = Status(iota)
	Open
	Unconfirmed
	Canceled
	Confirmed
	Settled
)

// An Order represents the want to perform a trade of assets.
type Order struct {
	Signature []byte `json:"signature"`
	ID        ID     `json:"id"`

	Type   Type      `json:"type"`
	Parity Parity    `json:"parity"`
	Expiry time.Time `json:"expiry"`

	FstCode   CurrencyCode     `json:"fstCode"`
	SndCode   CurrencyCode     `json:"sndCode"`
	Price     stackint.Int1024 `json:"price"`
	MaxVolume stackint.Int1024 `json:"maxVolume"`
	MinVolume stackint.Int1024 `json:"minVolume"`

	Nonce stackint.Int1024 `json:"nonce"`
}

// NewOrder returns a new Order and computes the ID.
func NewOrder(ty Type, parity Parity, expiry time.Time, fstCode, sndCode CurrencyCode, price, maxVolume, minVolume, nonce stackint.Int1024) *Order {
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
	order.ID = ID(order.Hash())
	return order
}

// NewOrderFromJSONFile returns an order that is unmarshaled from a JSON file.
func NewOrderFromJSONFile(fileName string) (Order, error) {
	order := Order{}
	file, err := os.Open(fileName)
	if err != nil {
		return order, err
	}
	defer file.Close()
	if err := json.NewDecoder(file).Decode(&order); err != nil {
		return order, err
	}
	return order, nil
}

// NewOrdersFromJSONFile returns an array of orders that is unmarshaled from a JSON file.
func NewOrdersFromJSONFile(fileName string) ([]Order, error) {
	orders := []Order{}
	file, err := os.Open(fileName)
	if err != nil {
		return orders, err
	}
	defer file.Close()
	if err := json.NewDecoder(file).Decode(&orders); err != nil {
		return orders, err
	}
	return orders, nil
}

// WriteOrdersToJSONFile writes an array of orders into a JSON file.
func WriteOrdersToJSONFile(fileName string, orders []*Order) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := json.NewEncoder(file).Encode(&orders); err != nil {
		return err
	}
	return nil
}

// Split the Order into n OrderFragments, where k OrderFragments are needed to
// reconstruct the Order. Returns a slice of all n OrderFragments, or an error.
func (order *Order) Split(n, k int64, prime *stackint.Int1024) ([]*Fragment, error) {
	fstCode := stackint.FromUint(uint(order.FstCode))
	fstCodeShares, err := shamir.Split(n, k, prime, &fstCode)
	if err != nil {
		return nil, err
	}
	sndCode := stackint.FromUint(uint(order.SndCode))
	sndCodeShares, err := shamir.Split(n, k, prime, &sndCode)
	if err != nil {
		return nil, err
	}
	priceShares, err := shamir.Split(n, k, prime, &order.Price)
	if err != nil {
		return nil, err
	}
	maxVolumeShares, err := shamir.Split(n, k, prime, &order.MaxVolume)
	if err != nil {
		return nil, err
	}
	minVolumeShares, err := shamir.Split(n, k, prime, &order.MinVolume)
	if err != nil {
		return nil, err
	}
	fragments := make([]*Fragment, n)
	for i := range fragments {
		fragments[i] = NewFragment(
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
	return fragments, nil
}

// Hash returns the Keccak256 hash of an Order. This hash is used to create the
// ID and signature for an Order.
func (order *Order) Hash() []byte {
	return crypto.Keccak256(order.Bytes())
}

// Sign signs the order using the provided keypair, and assigns it the the order's
// Signature field.
func (order *Order) Sign(signer crypto.Signer) error {
	var err error
	order.Signature, err = signer.Sign(order)
	return err
}

// VerifySignature verifies that the Signature field has been signed by the provided
// ID's private key, returning an error if the signature is invalid
func (order *Order) VerifySignature(addr string) error {
	return crypto.VerifySignature(order, order.Signature, addr)
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

// Equal returns an equality check between two Orders.
func (order *Order) Equal(other *Order) bool {
	return order.ID.Equal(other.ID) &&
		order.Type == other.Type &&
		order.Parity == other.Parity &&
		order.Expiry.Equal(other.Expiry) &&
		order.FstCode == other.FstCode &&
		order.SndCode == other.SndCode &&
		order.Price.Cmp(&other.Price) == 0 &&
		order.MaxVolume.Cmp(&other.MaxVolume) == 0 &&
		order.MinVolume.Cmp(&other.MinVolume) == 0 &&
		order.Nonce.Cmp(&other.Nonce) == 0
}
