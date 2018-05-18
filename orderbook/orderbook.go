package orderbook

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/shamir"
)

var ErrOrderNotFound = errors.New("order not found")
var ErrOrderFragmentNotFound = errors.New("order fragment not found")

type Client interface {

	// OpenOrder by sending an order.EncryptedFragment to an
	// identity.MultiAddress. The order.EncryptedFragment will be stored by the
	// Server hosted at the identity.MultiAddress.
	OpenOrder(context.Context, identity.MultiAddress, order.EncryptedFragment) error
}

type Server interface {
	OpenOrder(context.Context, order.EncryptedFragment) error
}

type Listener interface {
	OnConfirmOrderMatch(order.Order, order.Order)
}

type Orderbooker interface {
	Server

	// Sync orders and order states from the Ren Ledger to this local
	// Orderbooker. Returns a list of changes that were made to this local
	// Orderbooker during the synchronization.
	Sync() ([]OrderUpdate, error)

	// OrderFragment stored in this local Orderbooker. These are received from
	// other Orderbookers calling Orderbooker.OpenOrder to send an
	// order.EncryptedFragment to this local Orderbooker.
	OrderFragment(order.ID) (order.Fragment, error)

	// Order that has been reconstructed and stored in this local Orderbooker.
	// This only happens for orders that have been matched and confirmed.
	Order(order.ID) (order.Order, error)

	// Confirm an order match.
	ConfirmOrderMatch(order.ID, order.ID) error
}

type orderbook struct {
	crypto.RsaKey
	storer Storer
	syncer Syncer
}

func NewOrderbook(key crypto.RsaKey, storer Storer, syncer Syncer) Orderbooker {
	return &orderbook{
		RsaKey: key,
		storer: storer,
		syncer: syncer,
	}
}

func (book *orderbook) OpenOrder(ctx context.Context, orderFragment order.EncryptedFragment) error {
	tokens, err := book.RsaKey.Decrypt(orderFragment.Tokens)
	if err != nil {
		return err
	}
	var tokenShare shamir.Share
	err = json.Unmarshal(tokens, tokenShare)
	if err != nil {
		return fmt.Errorf("cannot decrypt tokens: %v", err)
	}

	// Decrypt price
	decryptedPriceCo, err := book.RsaKey.Decrypt(orderFragment.Price.Co)
	if err != nil {
		return fmt.Errorf("cannot decrypt price co: %v", err)
	}
	decryptedPriceExp, err := book.RsaKey.Decrypt(orderFragment.Price.Exp)
	if err != nil {
		return fmt.Errorf("cannot decrypt price exp: %v", err)
	}
	price := order.CoExpShare{
		Co: shamir.Share{}, Exp: shamir.Share{},
	}
	if err := price.Co.UnmarshalBinary(decryptedPriceCo); err != nil {
		return err
	}
	if err := price.Exp.UnmarshalBinary(decryptedPriceExp); err != nil {
		return err
	}

	// Decrypt volume
	decryptedVolumeCo, err := book.RsaKey.Decrypt(orderFragment.Volume.Co)
	if err != nil {
		return err
	}
	decryptedVolumeExp, err := book.RsaKey.Decrypt(orderFragment.Volume.Exp)
	if err != nil {
		return err
	}
	volume := order.CoExpShare{
		Co: shamir.Share{}, Exp: shamir.Share{},
	}
	if err := volume.Co.UnmarshalBinary(decryptedVolumeCo); err != nil {
		return err
	}
	if err := volume.Exp.UnmarshalBinary(decryptedVolumeExp); err != nil {
		return err
	}

	// Decrypt minVolume
	decryptedMinVolumeCo, err := book.RsaKey.Decrypt(orderFragment.Volume.Co)
	if err != nil {
		return err
	}
	decryptedMinVolumeExp, err := book.RsaKey.Decrypt(orderFragment.Volume.Exp)
	if err != nil {
		return err
	}
	minVolume := order.CoExpShare{
		Co: shamir.Share{}, Exp: shamir.Share{},
	}
	if err := minVolume.Co.UnmarshalBinary(decryptedMinVolumeCo); err != nil {
		return err
	}
	if err := minVolume.Exp.UnmarshalBinary(decryptedMinVolumeExp); err != nil {
		return err
	}

	fragment := order.Fragment{
		OrderID:       orderFragment.OrderID,
		OrderType:     orderFragment.OrderType,
		OrderParity:   orderFragment.OrderParity,
		OrderExpiry:   orderFragment.OrderExpiry,
		ID:            orderFragment.ID,
		Tokens:        tokenShare,
		Price:         price,
		Volume:        volume,
		MinimumVolume: minVolume,
	}

	return book.storer.InsertOrderFragment(fragment)
}

func (book *orderbook) Sync() ([]OrderUpdate, error) {
	return book.syncer.Sync()
}

func (book *orderbook) OrderFragment(id order.ID) (order.Fragment, error) {
	return book.storer.OrderFragment(id)
}

func (book *orderbook) Order(id order.ID) (order.Order, error) {
	return book.storer.Order(id)
}

func (book *orderbook) ConfirmOrderMatch(buy order.ID, sell order.ID) error {
	return book.syncer.ConfirmOrderMatch(buy, sell)
}
