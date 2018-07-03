package orderbook

import (
	"fmt"
	"math/big"
	"sync/atomic"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/order"
)

type syncer struct {
	contract          ContractBinder
	blockNumberOffset *big.Int
	blockNumberLimit  *big.Int

	limit int

	purgeOffset int
	purgeCache  []int

	ordersOffset int
	orders       []order.ID
}

// newSyncer returns a new Syncer that will synchronise a bounded number of
// changes from a ContractBinder. The Syncer will only considered order.IDs
// that were opened in a specific block number range.
func newSyncer(contract ContractBinder, blockNumberOffset, blockNumberLimit *big.Int, offset, limit int) *syncer {
	return &syncer{
		contract:          contract,
		blockNumberOffset: blockNumberOffset,
		blockNumberLimit:  blockNumberLimit,

		limit: limit,

		purgeOffset: 0,
		purgeCache:  make([]int, limit),

		ordersOffset: offset,
		orders:       make([]order.ID, limit),
	}
}

// sync orders and order.Status changes until the done channel is closed.
// All changes are produced as Notifications. Notifications of the
// NotificationOpenOrder type will not have an associated order.Fragment.
func (syncer *syncer) sync(done <-chan struct{}, notifications chan<- Notification, errs chan<- error) {
	syncer.syncConfirmsAndCancels(done, notifications, errs)
	syncer.syncOpens(done, notifications, errs)
}

func (syncer *syncer) syncConfirmsAndCancels(done <-chan struct{}, notifications chan<- Notification, errs chan<- error) {
	if len(syncer.orders) == 0 {
		return
	}

	// Compute the bounds of the orders that will be inspected for purging
	offset := syncer.purgeOffset % len(syncer.orders)
	limit := offset + syncer.limit
	if limit > len(syncer.orders) {
		limit = len(syncer.orders)
	}

	// j indexes into the syncer.purgeCache so that we can store indices that
	// need to be purged without using a heavy mutex, or runtime allocation
	j := int64(0)

	// Iterate over the bounded window concurrently and inspect each order.ID
	// to see if it needs to be purged
	dispatch.ForAll(syncer.orders[offset:limit], func(i int) {
		orderID := syncer.orders[i]

		// Produce notifications for changes in status
		notification := Notification(nil)
		status, err := syncer.contract.Status(orderID)
		if err != nil {
			select {
			case <-done:
			case errs <- fmt.Errorf("cannot sync order block status: %v", err):
			}
			return
		}
		switch status {
		case order.Confirmed:
			notification = NotificationConfirmOrder{OrderID: orderID}
		case order.Canceled:
			notification = NotificationCancelOrder{OrderID: orderID}
		default:
			return
		}
		select {
		case <-done:
			return
		case notifications <- notification:
		}

		// Store the index of the current order.ID so that it can be purged
		j := atomic.AddInt64(&j, 1)
		syncer.purgeCache[j-1] = offset + i
	})

	// Purge all stored indices by removing them from the syncer.orders slice
	for i := int64(0); i < j; i++ {
		purge := syncer.purgeCache[i]
		syncer.orders = append(syncer.orders[:purge], syncer.orders[purge+1:]...)
	}
	syncer.purgeOffset += syncer.limit
}

func (syncer *syncer) syncOpens(done <-chan struct{}, notifications chan<- Notification, errs chan<- error) {

	orderIDs, err := syncer.contract.Orders(syncer.ordersOffset, syncer.limit)
	if err != nil {
		select {
		case <-done:
		case errs <- err:
		}
		return
	}
	syncer.ordersOffset += len(orderIDs)

	// Produce a notification for each order that has been synchronised
	for _, orderID := range orderIDs {

		// Ignore orders that are outside the considered block range
		blockNumber, err := syncer.contract.BlockNumber(orderID)
		if err != nil {
			select {
			case <-done:
				return
			case errs <- fmt.Errorf("cannot sync order block number: %v", err):
				continue
			}
		}
		if blockNumber.Cmp(syncer.blockNumberOffset) == -1 {
			continue
		}
		if blockNumber.Sub(blockNumber, syncer.blockNumberLimit).Cmp(syncer.blockNumberOffset) == 1 {
			continue
		}

		// Create and produce the notification
		notification := Notification(nil)
		status, err := syncer.contract.Status(orderID)
		if err != nil {
			select {
			case <-done:
				return
			case errs <- fmt.Errorf("cannot sync order block status: %v", err):
				continue
			}
		}
		switch status {
		case order.Open:
			notification = NotificationOpenOrder{OrderID: orderID}
			syncer.orders = append(syncer.orders, orderID)
		case order.Confirmed:
			notification = NotificationConfirmOrder{OrderID: orderID}
		case order.Canceled:
			notification = NotificationCancelOrder{OrderID: orderID}
		default:
			continue
		}

		select {
		case <-done:
			return
		case notifications <- notification:
		}
	}
}
