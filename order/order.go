package order

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"os"
	"time"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/shamir"
)

// A Signature is the ECDSA signature of an order ID.
type Signature [65]byte

// An ID is the keccak256 hash of an Order.
type ID [32]byte

// Equal returns an equality check between two DeltaFragmentIDs.
func (id ID) Equal(other ID) bool {
	return bytes.Equal(id[:], other[:])
}

// String returns a truncated base64 encoding of the ID.
func (id ID) String() string {
	return base64.StdEncoding.EncodeToString(id[:8])
}

// Token is a numerical representation of a token supported by Republic
// Protocol.
type Token uint32

// Token values.
const (
	TokenBTC Token = 0
	TokenETH Token = 1
	TokenDGX Token = 256
	TokenREN Token = 65536
)

func (token Token) String() string {
	switch token {
	case TokenBTC:
		return "BTC"
	case TokenETH:
		return "ETH"
	case TokenDGX:
		return "DGX"
	case TokenREN:
		return "REN"
	default:
		return "unrecognized token "
	}
}

// Tokens are a numerical representation of the token pairings supported by
// Republic Protocol.
type Tokens uint64

// Tokens values.
const (
	TokensBTCETH Tokens = Tokens((uint64(TokenBTC) << 32) | uint64(TokenETH))
	TokensBTCDGX Tokens = Tokens((uint64(TokenBTC) << 32) | uint64(TokenDGX))
	TokensBTCREN Tokens = Tokens((uint64(TokenBTC) << 32) | uint64(TokenREN))
	TokensETHDGX Tokens = Tokens((uint64(TokenETH) << 32) | uint64(TokenDGX))
	TokensETHREN Tokens = Tokens((uint64(TokenETH) << 32) | uint64(TokenREN))
	TokensDGXREN Tokens = Tokens((uint64(TokenDGX) << 32) | uint64(TokenREN))
)

// PriorityToken returns the priority token of a token pair.
func (tokens Tokens) PriorityToken() Token {
	return Token(tokens & 0x00000000FFFFFFFF)
}

// NonPriorityToken returns the non-priority token of a token pair.
func (tokens Tokens) NonPriorityToken() Token {
	return Token(tokens >> 32)
}

// A Type is a publicly bit of information that determines the type of
// trade that an Order is representing.
type Type int8

// Type values.
const (
	TypeMidpoint Type = 0
	TypeLimit    Type = 1
)

// The Parity of an Order determines whether it is buy or a sell.
type Parity int8

// Parity values.
const (
	ParityBuy  Parity = 0
	ParitySell Parity = 1
)

// The Status shows what status the order is in.
type Status uint8

// Status values.
const (
	Nil = Status(iota)
	Open
	Confirmed
	Canceled
)

// An Order represents the want to perform a trade of assets.
type Order struct {
	Signature     Signature `json:"signature"`
	ID            ID        `json:"id"`
	Type          Type      `json:"type"`
	Parity        Parity    `json:"parity"`
	Expiry        time.Time `json:"expiry"`
	Tokens        Tokens    `json:"tokens"`
	Price         CoExp     `json:"price"`
	Volume        CoExp     `json:"volume"`
	MinimumVolume CoExp     `json:"minimumVolume"`
	Nonce         int64     `json:"nonce"`
}

// NewOrder returns a new Order and computes the ID.
func NewOrder(ty Type, parity Parity, expiry time.Time, tokens Tokens, price, volume, minimumVolume CoExp, nonce int64) Order {
	order := Order{
		Type:          ty,
		Parity:        parity,
		Expiry:        expiry,
		Tokens:        tokens,
		Price:         price,
		Volume:        volume,
		MinimumVolume: minimumVolume,
		Nonce:         nonce,
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
func (order *Order) Split(n, k int64) ([]Fragment, error) {
	tokens, err := shamir.Split(n, k, uint64(order.Tokens))
	if err != nil {
		return nil, err
	}
	priceCos, err := shamir.Split(n, k, uint64(order.Price.Co))
	if err != nil {
		return nil, err
	}
	priceExps, err := shamir.Split(n, k, uint64(order.Price.Exp))
	if err != nil {
		return nil, err
	}
	volumeCos, err := shamir.Split(n, k, uint64(order.Volume.Co))
	if err != nil {
		return nil, err
	}
	volumeExps, err := shamir.Split(n, k, uint64(order.Volume.Exp))
	if err != nil {
		return nil, err
	}
	minimumVolumeCos, err := shamir.Split(n, k, uint64(order.MinimumVolume.Co))
	if err != nil {
		return nil, err
	}
	minimumVolumeExps, err := shamir.Split(n, k, uint64(order.MinimumVolume.Exp))
	if err != nil {
		return nil, err
	}
	fragments := make([]Fragment, n)
	for i := range fragments {
		fragments[i] = NewFragment(
			order.ID,
			order.Type,
			order.Parity,
			tokens[i],
			CoExpShare{Co: priceCos[i], Exp: priceExps[i]},
			CoExpShare{Co: volumeCos[i], Exp: volumeExps[i]},
			CoExpShare{Co: minimumVolumeCos[i], Exp: minimumVolumeExps[i]},
		)
	}
	return fragments, nil
}

// Hash returns the Keccak256 hash of an Order. This hash is used to create the
// ID and signature for an Order.
func (order *Order) Hash() [32]byte {
	hash := crypto.Keccak256(order.Bytes())
	hash32 := [32]byte{}
	copy(hash32[:], hash)
	return hash32
}

// Bytes returns an Order serialized into a bytes.
// TODO: This function should return an error.
func (order *Order) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, order.Type)
	binary.Write(buf, binary.BigEndian, order.Parity)
	binary.Write(buf, binary.BigEndian, order.Expiry.Unix())
	binary.Write(buf, binary.BigEndian, order.Tokens)
	binary.Write(buf, binary.BigEndian, order.Price)
	binary.Write(buf, binary.BigEndian, order.Volume)
	binary.Write(buf, binary.BigEndian, order.MinimumVolume)
	binary.Write(buf, binary.BigEndian, order.bytesFromNonce())
	return buf.Bytes()
}

// Equal returns an equality check between two Orders.
func (order *Order) Equal(other *Order) bool {
	return bytes.Equal(order.ID[:], other.ID[:]) &&
		order.Type == other.Type &&
		order.Parity == other.Parity &&
		order.Expiry.Equal(other.Expiry) &&
		order.Tokens == other.Tokens &&
		order.Price == other.Price &&
		order.Volume == other.Volume &&
		order.MinimumVolume == other.MinimumVolume &&
		order.Nonce == other.Nonce
}

func (order *Order) bytesFromNonce() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, order.Nonce)
	return crypto.Keccak256(buf.Bytes())
}
