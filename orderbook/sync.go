package orderbook

import (
	"fmt"
	"log"

	"github.com/republicprotocol/republic-go/order"
)

type Syncer interface {
	Sync() (Notifications, error)
}

type syncer struct {
	// Stores for storing and loading data
	pointerStore PointerStorer
	orderStore   OrderStorer

	// ContractBinder exposes methods to pull changes from Ethereum and the
	// control parameters for customising when to pull changes
	contractBinder ContractBinder
	limit          int
	resyncPointer  int
}

func NewSyncer(pointerStore PointerStorer, orderStore OrderStorer, contractBinder ContractBinder, limit int) Syncer {
	return &syncer{
		pointerStore: pointerStore,
		orderStore:   orderStore,

		contractBinder: contractBinder,
		limit:          limit,
		resyncPointer:  0,
	}
}

// Sync implements the Syncer interface.
func (syncer *syncer) Sync() (Notifications, error) {
	notifications := make(Notifications, 0, syncer.limit)
	if err := syncer.sync(&notifications); err != nil {
		return notifications, err
	}
	if err := syncer.resync(&notifications); err != nil {
		return notifications, err
	}
	return notifications, nil
}

func (syncer *syncer) sync(notifications *Notifications) error {
	// Load the current pointer
	pointer, err := syncer.pointerStore.Pointer()
	if err != nil {
		return fmt.Errorf("cannot load pointer: %v", err)
	}

	// Synchronise new orders from the ContractBinder
	orderIDs, orderStatuses, traders, err := syncer.contractBinder.Orders(int(pointer), syncer.limit)
	if err != nil {
		return fmt.Errorf("cannot load orders from contract binder: %v", err)
	}
	if len(orderIDs) > 0 {
		log.Printf("[info] (sync) changed = %v", len(orderIDs))
	}

	// Store the resulting pointer so that we do not re-sync orders next time
	if err := syncer.pointerStore.PutPointer(pointer + Pointer(len(orderIDs))); err != nil {
		log.Printf("[error] (sync) cannot store pointer: %v", err)
	}

	// Logging data
	numOpenOrders := 0
	numConfirmedOrders := 0
	numCanceledOrders := 0
	numUnknownOrders := 0
	defer func() {
		if numOpenOrders > 0 {
			log.Printf("[info] (sync) opened = %v", numOpenOrders)
		}
		if numConfirmedOrders > 0 {
			log.Printf("[info] (sync) confirmed = %v", numConfirmedOrders)
		}
		if numCanceledOrders > 0 {
			log.Printf("[info] (sync) canceled = %v", numCanceledOrders)
		}
		if numUnknownOrders > 0 {
			log.Printf("[info] (sync) unknown = %v", numUnknownOrders)
		}
	}()

	for i, orderID := range orderIDs {
		switch orderStatuses[i] {
		case order.Open:
			numOpenOrders++
			notification := NotificationOpenOrder{OrderID: orderID, Trader: traders[i]}
			*notifications = append(*notifications, notification)
		case order.Confirmed:
			numConfirmedOrders++
			notification := NotificationConfirmOrder{OrderID: orderID}
			*notifications = append(*notifications, notification)
		case order.Canceled:
			numCanceledOrders++
			notification := NotificationCancelOrder{OrderID: orderID}
			*notifications = append(*notifications, notification)
		default:
			numUnknownOrders++
		}
	}
	return nil
}

func (syncer *syncer) resync(notifications *Notifications) error {
	orderIter, err := syncer.orderStore.Orders()
	if err != nil {
		return fmt.Errorf("cannot load pointer: %v", err)
	}
	defer orderIter.Release()

	orders, orderStatuses, _, err := orderIter.Collect()
	if err != nil {
		return fmt.Errorf("cannot collect orders: %v", err)
	}
	if len(orders) == 0 {
		return nil
	}

	// Log information about the resync at the end of the function
	numClosedOrders := 0
	defer func() {
		if numClosedOrders > 0 {
			log.Printf("[info] (sync) closed = %v", numClosedOrders)
		}
	}()

	// Function for deleting order IDs from storage
	deleteOrder := func(orderID order.ID, orderStatus order.Status) {
		numClosedOrders++
		if err := syncer.orderStore.DeleteOrder(orderID); err != nil {
			log.Printf("[error] (sync) cannot delete order: %v", err)
			return
		}

		switch orderStatus {
		case order.Confirmed:
			notification := NotificationConfirmOrder{OrderID: orderID}
			*notifications = append(*notifications, notification)
		case order.Canceled:
			notification := NotificationCancelOrder{OrderID: orderID}
			*notifications = append(*notifications, notification)
		default:
			return
		}
	}

	offset := syncer.resyncPointer
	limit := 2 * syncer.limit
	if limit > len(orders) {
		limit = len(orders)
	}
	for i := 0; i < 2*syncer.limit; i++ {
		syncer.resyncPointer = (offset + i) % len(orders)

		orderID, orderStatus := orders[syncer.resyncPointer], orderStatuses[syncer.resyncPointer]
		if orderStatus != order.Open {
			deleteOrder(orderID, orderStatus)
			continue
		}

		orderStatus, err = syncer.contractBinder.Status(orderID)
		if err != nil {
			log.Printf("[error] (sync) cannot load order status: %v", err)
			continue
		} else if orderStatus != order.Open {
			deleteOrder(orderID, orderStatus)
		}

		orderDepth, err := syncer.contractBinder.Depth(orderID)
		if err != nil {
			log.Printf("[error] (sync) cannot load order status: %v", err)
			continue
		}
		if orderDepth > 10000 {
			deleteOrder(orderID, orderStatus)
		}
	}
	return nil
}
