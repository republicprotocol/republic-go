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
	pointerStore       PointerStorer
	orderStore         OrderStorer
	orderFragmentStore OrderFragmentStorer

	// ContractBinder exposes methods to pull changes from Ethereum and the
	// control parameters for customizing when to pull changes
	contractBinder ContractBinder
	limit          int
	resyncPointer  int
}

func NewSyncer(pointerStore PointerStorer, orderStore OrderStorer, orderFragmentStore OrderFragmentStorer, contractBinder ContractBinder, limit int) Syncer {
	return &syncer{
		pointerStore:       pointerStore,
		orderStore:         orderStore,
		orderFragmentStore: orderFragmentStore,

		contractBinder: contractBinder,
		limit:          limit,
		resyncPointer:  0,
	}
}

// Sync implements the Syncer interface.
func (syncer *syncer) Sync() (Notifications, error) {
	notifications := make(Notifications, 0, syncer.limit)
	log.Println("[info] (sync) started")
	if err := syncer.sync(&notifications); err != nil {
		log.Printf("[info] (sync) errored = %v", err)
		return notifications, err
	}
	log.Println("[info] (resync) started")
	if err := syncer.resync(&notifications); err != nil {
		log.Printf("[info] (resync) errored = %v", err)
		return notifications, err
	}
	log.Printf("[info] (sync) completed = %v", len(notifications))
	return notifications, nil
}

func (syncer *syncer) sync(notifications *Notifications) error {
	// Load the current pointer
	pointer, err := syncer.pointerStore.Pointer()
	if err != nil {
		return fmt.Errorf("cannot load pointer: %v", err)
	}

	// Synchronize new orders from the ContractBinder
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
			notification := NotificationOpenOrder{OrderID: orderID, Trader: traders[i], Priority: uint(pointer) + uint(i)}
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

	orders, _, traders, priorities, err := orderIter.Collect()
	if err != nil {
		log.Printf("[error] (resync) cannot collect orders: %v", err)
	}
	if len(orders) == 0 {
		return nil
	}

	log.Printf("[info] (resync) %v orders, resync ptr = %v", len(orders), syncer.resyncPointer)

	// Log information about the resync at the end of the function
	numClosedOrders := 0
	defer func() {
		if numClosedOrders > 0 {
			log.Printf("[info] (resync) closed = %v", numClosedOrders)
		}
	}()

	// Function for deleting order IDs from storage
	deleteOrder := func(orderID order.ID, orderStatus order.Status) {
		numClosedOrders++
		if _, err := syncer.orderFragmentStore.OrderFragment(orderID); err == nil {
			if err := syncer.orderFragmentStore.DeleteOrderFragment(orderID); err != nil {
				log.Printf("[error] (resync) cannot delete order fragment: %v", err)
			}
		}
		if err := syncer.orderStore.DeleteOrder(orderID); err != nil {
			log.Printf("[error] (resync) cannot delete order: %v", err)
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
	if offset >= limit {
		return nil
	}
	for i := offset; i < limit; i++ {
		pointer := (i) % len(orders)

		orderID := orders[pointer]
		trader := traders[pointer]
		priority := priorities[pointer]
		orderStatus, err := syncer.contractBinder.Status(orderID)
		if err != nil {
			log.Printf("[error] (resync) cannot load order status: %v", err)
			continue
		}

		switch orderStatus {
		case order.Canceled:
			deleteOrder(orderID, orderStatus)
		case order.Confirmed:
			settleStatus, err := syncer.contractBinder.SettlementStatus(orderID)
			if err != nil {
				log.Printf("[error] (resync) cannot load order settlement status: %v", err)
				continue
			}
			if settleStatus > 1 {
				deleteOrder(orderID, order.Confirmed)
			}
		case order.Open:
			if fragment, err := syncer.orderFragmentStore.OrderFragment(orderID); err == nil {
				log.Printf("[info] (resync) generating new notification %v, resync ptr = %v", orderID, syncer.resyncPointer+pointer)
				notification := NotificationOpenOrder{OrderID: orderID, OrderFragment: fragment, Priority: priority, Trader: trader}
				*notifications = append(*notifications, notification)
			}
		}
	}
	syncer.resyncPointer += limit - offset
	return nil
}
