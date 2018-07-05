package orderbook

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"time"

	"github.com/republicprotocol/republic-go/logger"

	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/registry"
)

type syncer struct {
	// Epoch at the beginning of the syncer range
	epoch registry.Epoch

	// Stores for storing and loading data
	pointerStore       PointerStorer
	orderStore         OrderStorer
	orderFragmentStore OrderFragmentStorer

	// ContractBinder exposes methods to pull changes from Ethereum and the
	// control parameters for customising when to pull changes
	contractBinder ContractBinder
	interval       time.Duration
	limit          int
}

// newSyncer returns a new Syncer that will synchronise changes from a
// ContractBinder for a specific registry.Epoch. At each time interval, the
// Syncer will synchronise all changes up to some maximum limit.
func newSyncer(epoch registry.Epoch, pointerStore PointerStorer, orderStore OrderStorer, orderFragmentStore OrderFragmentStorer, contractBinder ContractBinder, interval time.Duration, limit int) *syncer {
	return &syncer{
		epoch: epoch,

		pointerStore:       pointerStore,
		orderStore:         orderStore,
		orderFragmentStore: orderFragmentStore,

		contractBinder: contractBinder,
		interval:       interval,
		limit:          limit,
	}
}

// sync orders and order.Status changes until the done channel is closed.
// All changes are produced as Notifications. Notifications of the
// NotificationOpenOrder type will not have an associated order.Fragment.
func (syncer *syncer) sync(done <-chan struct{}, orderFragments <-chan order.Fragment) (<-chan Notification, <-chan error) {
	notifications := make(chan Notification)
	errs := make(chan error)

	go func() {
		defer close(notifications)
		defer close(errs)

		ticker := time.NewTicker(syncer.interval)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				logger.Network(logger.LevelDebug, fmt.Sprintf("synchronising orderbook in epoch %v", base64.StdEncoding.EncodeToString(syncer.epoch.Hash[:8])))
				syncer.syncClosures(done, notifications, errs)
				syncer.syncOpens(done, notifications, errs)
			case orderFragment, ok := <-orderFragments:
				if !ok {
					return
				}
				logger.Network(logger.LevelDebug, fmt.Sprintf("synchronising order fragment in epoch %v", base64.StdEncoding.EncodeToString(syncer.epoch.Hash[:8])))
				syncer.insertOrderFragment(orderFragment, done, notifications, errs)
			}
		}
	}()

	return notifications, errs
}

// syncClosures iterates through all orders and deletes those that are no
// longer open.
func (syncer *syncer) syncClosures(done <-chan struct{}, notifications chan<- Notification, errs chan<- error) {
	orderIter, err := syncer.orderStore.Orders()
	if err != nil {
		select {
		case <-done:
		case errs <- fmt.Errorf("cannot load order iterator: %v", err):
		}
		return
	}
	defer orderIter.Release()

	// Function for deleting order IDs from storage
	deleteOrder := func(orderID order.ID) {
		if err := syncer.orderStore.DeleteOrder(orderID); err != nil {
			select {
			case <-done:
				return
			case errs <- fmt.Errorf("cannot delete order: %v", err):
			}
		}
	}

	for orderIter.Next() {
		// Get the next order, and its status, and mark it for deltion if it is
		// not open
		orderID, orderStatus, err := orderIter.Cursor()
		if err != nil {
			select {
			case <-done:
				return
			case errs <- fmt.Errorf("cannot load order iterator cursor: %v", err):
				continue
			}
		}
		if orderStatus != order.Open {
			deleteOrder(orderID)
			continue
		}

		// Refresh the status and mark it for deltion if it is not open
		orderStatus, err = syncer.contractBinder.Status(orderID)
		if err != nil {
			select {
			case <-done:
				return
			case errs <- fmt.Errorf("cannot sync order status: %v", err):
				continue
			}
		}
		if orderStatus != order.Open {
			deleteOrder(orderID)
			continue
		}
	}
}

// syncOpens attempts to synchronise all new orders since the last time a
// synchronisation happened. Usually, these orders will be open, but if the
// interval between synchronisations is large enough then it is possible that
// they are already closed.
func (syncer *syncer) syncOpens(done <-chan struct{}, notifications chan<- Notification, errs chan<- error) {

	// Load the current pointer
	pointer, err := syncer.pointerStore.Pointer()
	if err != nil {
		select {
		case <-done:
		case errs <- fmt.Errorf("cannot load pointer: %v", err):
		}
		return
	}

	// Synchronise new orders from the ContractBinder
	orderIDs, orderStatuses, traders, err := syncer.contractBinder.Orders(int(pointer), syncer.limit)
	if err != nil {
		select {
		case <-done:
		case errs <- fmt.Errorf("cannot sync orders: %v", err):
		}
		return
	}

	// Store the resulting pointer so that we do not re-sync orders next time
	if err := syncer.pointerStore.PutPointer(pointer + Pointer(len(orderIDs))); err != nil {
		select {
		case <-done:
		case errs <- fmt.Errorf("cannot store pointer: %v", err):
		}
	}

	blockInterval := big.NewInt(0).Mul(big.NewInt(2), syncer.epoch.BlockInterval)
	for i, orderID := range orderIDs {

		// Ignore orders that are outside the considered block range
		blockNumber, err := syncer.contractBinder.BlockNumber(orderID)
		if err != nil {
			select {
			case <-done:
				return
			case errs <- fmt.Errorf("cannot sync order block number: %v", err):
				continue
			}
		}
		if blockNumber.Cmp(syncer.epoch.BlockNumber) == -1 {
			continue
		}
		if blockNumber.Sub(blockNumber, blockInterval).Cmp(syncer.epoch.BlockNumber) == 1 {
			continue
		}

		// Synchronise the status of this order and generate the appropriate
		// notification
		switch orderStatuses[i] {

		// Open orders need to check for the respective order fragment before a
		// notification can be generated
		case order.Open:
			syncer.insertOrder(orderID, orderStatuses[i], traders[i], done, notifications, errs)

		// Other statuses can generate notifications immediately
		case order.Confirmed:
			notification := NotificationConfirmOrder{OrderID: orderID}
			select {
			case <-done:
				return
			case notifications <- notification:
			}
		case order.Canceled:
			notification := NotificationCancelOrder{OrderID: orderID}
			select {
			case <-done:
				return
			case notifications <- notification:
			}
		}
	}
}

func (syncer *syncer) insertOrder(orderID order.ID, orderStatus order.Status, trader string, done <-chan struct{}, notifications chan<- Notification, errs chan<- error) {

	// Store the order
	if err := syncer.orderStore.PutOrder(orderID, orderStatus); err != nil {
		select {
		case <-done:
			return
		case errs <- fmt.Errorf("cannot store order: %v", err):
			// Continue even if there was an error storing the order
		}
	}

	// Check for the respective order fragment
	orderFragment, err := syncer.orderFragmentStore.OrderFragment(syncer.epoch, orderID)
	if err != nil {
		if err == ErrOrderFragmentNotFound {
			// No order fragment received yet
			return
		}
		select {
		case <-done:
		case errs <- err:
		}
		return
	}

	notification := NotificationOpenOrder{
		OrderID:       orderID,
		OrderFragment: orderFragment,
		Trader:        trader,
	}
	select {
	case <-done:
	case notifications <- notification:
	}
}

func (syncer *syncer) insertOrderFragment(orderFragment order.Fragment, done <-chan struct{}, notifications chan<- Notification, errs chan<- error) {
	// Store the order fragment
	if err := syncer.orderFragmentStore.PutOrderFragment(syncer.epoch, orderFragment); err != nil {
		select {
		case <-done:
			return
		case errs <- fmt.Errorf("cannot store order fragment: %v", err):
			// Continue even if there was an error storing the order fragment
		}
	}

	// Check for the respective order
	orderStatus, err := syncer.orderStore.Order(orderFragment.OrderID)
	if err != nil {
		if err == ErrOrderNotFound {
			// No order synchronised yet
			return
		}
		select {
		case <-done:
		case errs <- err:
		}
		return
	}
	if orderStatus != order.Open {
		// Order is synchronised but it is not open
		return
	}

	notification := NotificationOpenOrder{
		OrderID:       orderFragment.OrderID,
		OrderFragment: orderFragment,
	}
	select {
	case <-done:
	case notifications <- notification:
	}
}
