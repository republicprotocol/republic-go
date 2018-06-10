package orderbook

import (
	"context"
	"encoding/base64"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/order"
)

// Client for invoking the Server.OpenOrder ROC on a remote Server.
type Client interface {

	// OpenOrder by sending an order.EncryptedFragment to an
	// identity.MultiAddress. The order.EncryptedFragment will be stored by the
	// Server hosted at the identity.MultiAddress.
	OpenOrder(context.Context, identity.MultiAddress, order.EncryptedFragment) error
}

// Server for opening order.EncryptedFragments. This RPC should only be called
// after the respective order.Order has been opened on the Ethereum blockchain
// otherwise it will be ignored by the Server.
type Server interface {
	OpenOrder(context.Context, order.EncryptedFragment) error
}

// An Orderbook is responsible for receiving orders. It reads order.Order
// states from the Syncer, and reads order.EncryptedFragemnts from the Server.
// By combining these two interfaces into a single unified interface all data
// required for processing orders is exposed by the Orderbook interface.
type Orderbook interface {
	Server
	Syncer
}

type orderbook struct {
	crypto.RsaKey

	syncer Syncer
	storer Storer
}

// NewOrderbook returns an Orderbok that uses a crypto.RsaKey to decrypt the
// order.EncryptedFragments that it receives, and stores them in a Storer.
func NewOrderbook(key crypto.RsaKey, syncer Syncer, storer Storer) Orderbook {
	return &orderbook{
		RsaKey: key,

		syncer: syncer,
		storer: storer,
	}
}

// OpenOrder implements the Server interface.
func (book *orderbook) OpenOrder(ctx context.Context, orderFragment order.EncryptedFragment) error {
	fragment, err := orderFragment.Decrypt(*book.RsaKey.PrivateKey)
	if err != nil {
		return err
	}
	if fragment.OrderParity == order.ParityBuy {
		logger.BuyOrderReceived(logger.LevelDebugLow, base64.StdEncoding.EncodeToString(fragment.OrderID[:8]), base64.StdEncoding.EncodeToString(fragment.ID[:8]))
	} else {
		logger.SellOrderReceived(logger.LevelDebugLow, base64.StdEncoding.EncodeToString(fragment.OrderID[:8]), base64.StdEncoding.EncodeToString(fragment.ID[:8]))
	}

	return book.storer.InsertOrderFragment(fragment)
}

// Sync implements the Syncer interface.
func (book *orderbook) Sync() (ChangeSet, error) {
	return book.syncer.Sync()
}
