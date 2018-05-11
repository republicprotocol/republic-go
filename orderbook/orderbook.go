package orderbook

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/republicprotocol/republic-go/delta"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/order"
	"github.com/syndtr/goleveldb/leveldb"
)

// ErrWriteToClosedOrderbook is returned when an attempt to updated the
// Orderbook is made after a call to Orderbook.Close.
var ErrWriteToClosedOrderbook = errors.New("write to closed orderbook")

//type Syncer interface {
//	Open(order order.Order) error
//	Match(order order.Order) error
//	Confirm(order order.Order) error
//	Release(order order.Order) error
//	Settle(order order.Order) error
//	Cancel(order order.Order) error
//	Status(id order.ID) (order.Status, error)
//	Clear()
//}

// An Orderbook is responsible for store the historical orders.
// It also broadcasts the newly received orders to its subscribers.
type Orderbook struct {
	status         *leveldb.DB
	orderFragments *leveldb.DB
	orders         *leveldb.DB
	atoms          *leveldb.DB
	matches        *leveldb.DB

	broadcaster       *dispatch.Broadcaster
	broadcasterChDone chan struct{}
	broadcasterCh     chan interface{}
}

// NewOrderbook creates a new Orderbook.
func NewOrderbook() (Orderbook, error) {
	dbPath := path.Join(os.Getenv("HOME"), ".darknode", "orderbook")
	status, err := leveldb.OpenFile(path.Join(dbPath, "status"), nil)
	if err != nil {
		return Orderbook{}, err
	}
	orderFragments, err := leveldb.OpenFile(path.Join(dbPath, "fragments"), nil)
	if err != nil {
		return Orderbook{}, err
	}
	orders, err := leveldb.OpenFile(path.Join(dbPath, "orders"), nil)
	if err != nil {
		return Orderbook{}, err
	}
	atoms, err := leveldb.OpenFile(path.Join(dbPath, "atoms"), nil)
	if err != nil {
		return Orderbook{}, err
	}
	matches, err := leveldb.OpenFile(path.Join(dbPath, "matches"), nil)
	if err != nil {
		return Orderbook{}, err
	}

	broadcaster := dispatch.NewBroadcaster(dispatch.MaxListeners)
	broadcasterChDone := make(chan struct{})
	broadcasterCh := make(chan interface{})
	go broadcaster.Broadcast(broadcasterChDone, broadcasterCh)

	return Orderbook{
		status:         status,
		orderFragments: orderFragments,
		orders:         orders,
		atoms:          atoms,
		matches:        matches,

		broadcaster:       broadcaster,
		broadcasterCh:     broadcasterCh,
		broadcasterChDone: broadcasterChDone,
	}, nil
}

// Close the Orderbook. All listeners will eventually be closed and no more
// listeners will be accepted.
func (orderbook *Orderbook) Close() {
	// Close all the levelDB files
	orderbook.status.Close()
	orderbook.orderFragments.Close()
	orderbook.orders.Close()
	orderbook.matches.Close()
	orderbook.atoms.Close()

	// Close the broadcaster
	orderbook.broadcaster.Close()
	close(orderbook.broadcasterChDone)
}

// Listen to the orderbook for updates. Calls to Orderbook.Listen are
// non-blocking, and the background worker is terminated when the done
// channel is closed. A read-only channel of updates is returned, and
// will be closed when no more data will be written to it.
func (orderbook *Orderbook) Listen(done <-chan struct{}) <-chan Update {
	listener := orderbook.broadcaster.Listen(done)
	subscriber := make(chan Update)

	go func() {
		defer close(subscriber)
		dispatch.CoBegin(func() {
			for {
				select {
				case <-done:
					return
				case <-orderbook.broadcasterChDone:
					return
				case msg, ok := <-listener:
					if !ok {
						return
					}
					if msg, ok := msg.(Update); ok {
						select {
						case <-done:
							return
						case <-orderbook.broadcasterChDone:
							return
						case subscriber <- msg:
						}
					}
				}
			}
		}, func() {
			// Stream historical data
			iter := orderbook.status.NewIterator(nil, nil)
			for iter.Next() {
				key, value := iter.Key(), iter.Value()
				if len(value) != 1 {
					return
				}
				update := NewUpdate(order.ID(key), value[0])
				select {
				case <-done:
					return
				case <-orderbook.broadcasterChDone:
					return
				case subscriber <- update:
				}
			}

			iter.Release()
		})
	}()

	return subscriber
}

// Open is called when we first receive the order fragment.
func (orderbook *Orderbook) Open(fragment order.Fragment) error {
	key := fragment.OrderID

	// Check existence of the order
	_ , err := orderbook.status.Get(key, nil )
	if err != leveldb.ErrNotFound{
		return err
	}

	encoded, err := json.Marshal(fragment)
	if err != nil {
		return err
	}
	if err := orderbook.orderFragments.Put([]byte(key), encoded, nil); err != nil {
		return err
	}
	if err := orderbook.status.Put([]byte(key), []byte{order.Open}, nil); err != nil {
		return err
	}

	update := NewUpdate(fragment.OrderID, order.Open)
	select {
	case <-orderbook.broadcasterChDone:
		return ErrWriteToClosedOrderbook
	case orderbook.broadcasterCh <- update:
		return nil
	}
}

// Match is called when we discover a match for the order.
func (orderbook *Orderbook) Match(dlt delta.Delta) error {
	// Save buy order in the orderbook
	buyKey := dlt.BuyOrderID

	// Check previous status of the order
	buyStatus , err := orderbook.status.Get(buyKey, nil )
	if err != nil {
		return err
	}
	if bytes.Compare(buyStatus , []byte{order.Open}) != 0{
		return fmt.Errorf("cannot matched order with status %v", buyStatus)
	}

	//todo : do we care about the expiry and nonce ( can get from the fragments table)
	fstCode, err := dlt.FstCode.ToUint()
	if err != nil {
		return err
	}
	sndCode, err := dlt.SndCode.ToUint()
	if err != nil {
		return err
	}
	buyOrder := order.Order{
		ID:        dlt.SellOrderID,
		Type:      order.TypeLimit,
		Parity:    order.ParityBuy,
		FstCode:   int64(fstCode),
		SndCode:   int64(sndCode),
		Price:     dlt.Price,
		MaxVolume: dlt.MaxVolume,
		MinVolume: dlt.MinVolume,
	}
	encodedBuy, err := json.Marshal(buyOrder)
	if err != nil {
		return err
	}
	if err := orderbook.orders.Put([]byte(buyKey), encodedBuy, nil); err != nil {
		return err
	}
	if err := orderbook.status.Put([]byte(buyKey), []byte{order.Unconfirmed}, nil); err != nil {
		return err
	}
	buyMatches, err := json.Marshal([]order.ID{dlt.SellOrderID})
	if err != nil {
		return err
	}
	if err := orderbook.matches.Put([]byte(buyStatus), buyMatches, nil); err != nil {
		return err
	}

	// Save sell order in the orderbook
	sellKey := dlt.SellOrderID

	// Check previous status of the order
	sellStatus , err := orderbook.status.Get(sellKey, nil )
	if err != nil {
		return err
	}
	if bytes.Compare(sellStatus , []byte{order.Open}) != 0{
		return fmt.Errorf("cannot matched order with status %v", sellStatus)
	}
	//todo : do we care about the expiry and nonce ( can get from the fragments table)
	fstCode, err = dlt.FstCode.ToUint()
	if err != nil {
		return err
	}
	sndCode, err = dlt.SndCode.ToUint()
	if err != nil {
		return err
	}
	sellOrder := order.Order{
		ID:        dlt.SellOrderID,
		Type:      order.TypeLimit,
		Parity:    order.ParitySell,
		FstCode:   int64(fstCode),
		SndCode:   int64(sndCode),
		Price:     dlt.Price,
		MaxVolume: dlt.MaxVolume,
		MinVolume: dlt.MinVolume,
	}
	encodedSell, err := json.Marshal(sellOrder)
	if err := orderbook.orders.Put([]byte(sellKey), encodedSell, nil); err != nil {
		return err
	}
	if err := orderbook.status.Put([]byte(sellKey), []byte{order.Unconfirmed}, nil); err != nil {
		return err
	}
	sellMatches, err := json.Marshal([]order.ID{dlt.SellOrderID})
	if err != nil {
		return err
	}
	if err := orderbook.matches.Put([]byte(sellKey), sellMatches, nil); err != nil {
		return err
	}

	buyUpdate := NewUpdate(dlt.BuyOrderID, order.Unconfirmed)
	sellUpdate := NewUpdate(dlt.SellOrderID, order.Unconfirmed)

	select {
	case <-orderbook.broadcasterChDone:
		return ErrWriteToClosedOrderbook
	case orderbook.broadcasterCh <- buyUpdate:
		select {
		case <-orderbook.broadcasterChDone:
			return ErrWriteToClosedOrderbook
		case orderbook.broadcasterCh <- sellUpdate:
			return nil
		}
	}
}

// Confirm is called when the order has been confirmed by the hyperdrive.
func (orderbook *Orderbook) Confirm(id order.ID) error {
	// Check existence of the order
	status, err := orderbook.status.Get(id, nil )
	if err != nil{
		return err
	}
	if bytes.Compare(status , []byte{order.Unconfirmed}) != 0{
		return fmt.Errorf("cannot confirm order with status %v", status)
	}

	if err := orderbook.status.Put(id, []byte{order.Confirmed}, nil); err != nil {
		return err
	}

	update := NewUpdate(id, order.Confirmed)
	select {
	case <-orderbook.broadcasterChDone:
		return ErrWriteToClosedOrderbook
	case orderbook.broadcasterCh <- update:
		return nil
	}
}

// Release is called when the order has been denied by the hyperdrive.
func (orderbook *Orderbook) Release(id order.ID) error {
	// Check existence of the order
	status, err := orderbook.status.Get(id, nil )
	if err != nil{
		return err
	}
	if bytes.Compare(status , []byte{order.Unconfirmed}) != 0{
		return fmt.Errorf("cannot release order with status %v", status)
	}

	if err := orderbook.status.Put(id, []byte{order.Open}, nil); err != nil {
		return err
	}
	if err := orderbook.matches.Delete(id, nil); err != nil {
		return err
	}

	update := NewUpdate(id, order.Open)
	select {
	case <-orderbook.broadcasterChDone:
		return ErrWriteToClosedOrderbook
	case orderbook.broadcasterCh <- update:
		return nil
	}
}

// Settle is called when the order is settled.
func (orderbook *Orderbook) Settle(id order.ID) error {
	// Check previous status of the order
	status, err := orderbook.status.Get(id, nil )
	if err != nil{
		return err
	}
	if bytes.Compare(status , []byte{order.Confirmed}) != 0{
		return fmt.Errorf("cannot settle order with status %v", status)
	}

	if err := orderbook.status.Put(id, []byte{order.Settled}, nil); err != nil {
		return err
	}

	update := NewUpdate(id, order.Settled)
	select {
	case <-orderbook.broadcasterChDone:
		return ErrWriteToClosedOrderbook
	case orderbook.broadcasterCh <- update:
		return nil
	}
}

// Cancel is called when the order is canceled.
func (orderbook *Orderbook) Cancel(id order.ID) error {
	// Check previous status of the order
	status, err := orderbook.status.Get(id, nil )
	if err != nil{
		return err
	}
	if bytes.Compare(status , []byte{order.Open}) != 0 || bytes.Compare(status , []byte{order.Unconfirmed}) != 0{
		return fmt.Errorf("too late too cancel the order")
	}
	if err := orderbook.status.Put(id, []byte{order.Canceled}, nil); err != nil {
		return err
	}

	update := NewUpdate(id, order.Canceled)
	select {
	case <-orderbook.broadcasterChDone:
		return ErrWriteToClosedOrderbook
	case orderbook.broadcasterCh <- update:
		return nil
	}
}

// Order retrieves information regarding an order.
func (orderbook *Orderbook) Order(id order.ID) (order.Order, error) {
	ord := order.Order{}
	ordBytes, err := orderbook.orders.Get(id, nil)
	if err != nil {
		return order.Order{}, err
	}
	err = json.Unmarshal(ordBytes, ord)
	if err != nil {
		return order.Order{}, err
	}

	return ord, nil
}

// OrderFragment retrieves information regarding an orderFragment.
func (orderbook *Orderbook) OrderFragment(id order.ID) (order.Fragment, error) {
	fragment := order.Fragment{}
	fragmentBytes, err := orderbook.orderFragments.Get(id, nil)
	if err != nil {
		return order.Fragment{}, err
	}
	err = json.Unmarshal(fragmentBytes, fragment)
	if err != nil {
		return order.Fragment{}, err
	}

	return fragment, nil
}

// Status retrieves status regarding an orderID.
func (orderbook *Orderbook) Status(id order.ID) (order.Status, error) {
	status, err := orderbook.status.Get(id, nil)
	if err != nil {
		return order.Nil, err
	}

	if len(status) != 1 {
		return order.Nil, errors.New("internal error : wrong status length ")
	}
	return status[0], nil
}

// CounterOrders retrieves the matched orders regarding an orderID.
func (orderbook *Orderbook) CounterOrders(id order.ID) ([]order.ID, error) {
	var matches []order.ID
	matchesBytes, err := orderbook.matches.Get(id, nil)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(matchesBytes, &matches); err != nil {
		return nil, err
	}

	return matches, nil
}
