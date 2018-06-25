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
	runes := []rune(base64.StdEncoding.EncodeToString(id[:]))
	return string(runes[:4])
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
	TokenABC Token = 65537
	TokenXYZ Token = 65538
)

// String returns a human-readable representation of a Token.
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
	case TokenABC:
		return "ABC"
	case TokenXYZ:
		return "XYZ"
	default:
		return "unexpected token"
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
	TokensETHABC Tokens = Tokens((uint64(TokenETH) << 32) | uint64(TokenABC))
	TokensETHXYZ Tokens = Tokens((uint64(TokenETH) << 32) | uint64(TokenXYZ))
	TokensDGXREN Tokens = Tokens((uint64(TokenDGX) << 32) | uint64(TokenREN))
	TokensDGXABC Tokens = Tokens((uint64(TokenDGX) << 32) | uint64(TokenABC))
	TokensDGXXYZ Tokens = Tokens((uint64(TokenDGX) << 32) | uint64(TokenXYZ))
)

// PriorityToken returns the priority token of a token pair.
func (tokens Tokens) PriorityToken() Token {
	return Token(tokens & 0x00000000FFFFFFFF)
}

// NonPriorityToken returns the non-priority token of a token pair.
func (tokens Tokens) NonPriorityToken() Token {
	return Token(tokens >> 32)
}

// String returns a human-readable representation of Tokens.
func (tokens Tokens) String() string {
	switch tokens {
	case TokensBTCETH:
		return "BTC-ETH"
	case TokensBTCDGX:
		return "BTC-DGX"
	case TokensBTCREN:
		return "BTC-REN"
	case TokensETHDGX:
		return "ETH-DGX"
	case TokensETHREN:
		return "ETH-REN"
	case TokensDGXREN:
		return "DGX-REN"
	default:
		return "unexpected tokens"
	}
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

// String returns a human-readable representation of the Parity.
func (parity Parity) String() string {
	switch parity {
	case ParityBuy:
		return "buy"
	case ParitySell:
		return "sell"
	default:
		return "unexpected parity"
	}
}

// Settlement is a unique identifier for the settlement layer used by the
// Order.
type Settlement uint32

// Values for Settlementt.
const (
	SettlementNil (Settlement) = iota
	SettlementRenEx
	SettlementRenExAtomic
)

// String implements the Stringer interface.
func (settlement Settlement) String() string {
	switch settlement {
	case SettlementRenEx:
		return "RenEx"
	case SettlementRenExAtomic:
		return "RenEx Atomic"
	default:
		return "unexpected order settlement"
	}
}

// The Status shows what status the order is in.
type Status uint8

// Status values.
const (
	Nil = Status(iota)
	Open
	Confirmed
	Canceled
)

// String implements the Stringer interface.
func (status Status) String() string {
	switch status {
	case Nil:
		return "nil"
	case Open:
		return "open"
	case Confirmed:
		return "confirmed"
	case Canceled:
		return "canceled"
	default:
		return "unexpected order status"
	}
}

// An Order represents the want to perform a trade of assets.
type Order struct {
	Signature  Signature  `json:"signature"`
	ID         ID         `json:"id"`
	Type       Type       `json:"type"`
	Parity     Parity     `json:"parity"`
	Settlement Settlement `json:"settlement"`
	Expiry     time.Time  `json:"expiry"`

	Tokens        Tokens `json:"tokens"`
	Price         CoExp  `json:"price"`
	Volume        CoExp  `json:"volume"`
	MinimumVolume CoExp  `json:"minimumVolume"`
	Nonce         uint64 `json:"nonce"`
}

// NewOrder returns a new Order and computes the ID.
func NewOrder(ty Type, parity Parity, settlement Settlement, expiry time.Time, tokens Tokens, price, volume, minimumVolume CoExp, nonce uint64) Order {
	order := Order{
		Type:       ty,
		Parity:     parity,
		Settlement: settlement,
		Expiry:     expiry,

		Tokens:        tokens,
		Price:         price,
		Volume:        volume,
		MinimumVolume: minimumVolume,
		Nonce:         nonce % shamir.Prime,
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
	return json.NewEncoder(file).Encode(&orders)
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
	nonces, err := shamir.Split(n, k, order.Nonce)
	if err != nil {
		return nil, err
	}
	fragments := make([]Fragment, n)
	for i := range fragments {
		fragments[i], err = NewFragment(
			order.ID,
			order.Type,
			order.Parity,
			order.Settlement,
			order.Expiry,
			tokens[i],
			CoExpShare{Co: priceCos[i], Exp: priceExps[i]},
			CoExpShare{Co: volumeCos[i], Exp: volumeExps[i]},
			CoExpShare{Co: minimumVolumeCos[i], Exp: minimumVolumeExps[i]},
			nonces[i],
		)
		if err != nil {
			return nil, err
		}
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
	binary.Write(buf, binary.BigEndian, order.Settlement)
	binary.Write(buf, binary.BigEndian, order.Expiry.Unix())
	binary.Write(buf, binary.BigEndian, order.Tokens)
	binary.Write(buf, binary.BigEndian, order.Price)
	binary.Write(buf, binary.BigEndian, order.Volume)
	binary.Write(buf, binary.BigEndian, order.MinimumVolume)
	binary.Write(buf, binary.BigEndian, order.BytesFromNonce())
	return buf.Bytes()
}

// Equal returns an equality check between two Orders.
func (order *Order) Equal(other *Order) bool {
	return bytes.Equal(order.ID[:], other.ID[:]) &&
		order.Type == other.Type &&
		order.Parity == other.Parity &&
		order.Settlement == other.Settlement &&
		order.Expiry.Equal(other.Expiry) &&
		order.Tokens == other.Tokens &&
		order.Price == other.Price &&
		order.Volume == other.Volume &&
		order.MinimumVolume == other.MinimumVolume &&
		order.Nonce == other.Nonce
}

// BytesFromNonce returns the uint64 nonce as a slice of bytes.
func (order *Order) BytesFromNonce() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, order.Nonce)
	return crypto.Keccak256(buf.Bytes())
}
