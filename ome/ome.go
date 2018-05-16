package ome

import (
	"sync"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/shamir"
	"github.com/republicprotocol/republic-go/smpc"
	"github.com/republicprotocol/republic-go/smpc/delta"
)

type Omer interface {
	OnChangeEpoch(ξ cal.Epoch) error
}

type Ome struct {
	done     <-chan struct{}
	mu       *sync.Mutex
	epoch    cal.Epoch
	delegate Delegate
	storer   Storer
	syncer   Syncer
	ranker   Ranker
	smpcer   smpc.Smpcer
	deltas   map[delta.ID][]delta.Fragment
}

func NewOme(done <-chan struct{}, pod <-chan cal.Epoch, delegate Delegate, storer Storer, syncer Syncer, ranker Ranker, smpcer smpc.Smpcer) Ome {
	ome := Ome{}
	ome.done = done
	ome.mu = new(sync.Mutex)
	ome.delegate = delegate
	ome.storer = storer
	ome.syncer = syncer
	ome.ranker = ranker
	ome.smpcer = smpcer
	ome.deltas = map[delta.ID][]delta.Fragment{}
	subChan := make(chan smpc.InstConnect)

	// Listen for network change and update the pod infor
	go func() {
		defer close(subchan)
		for {
			select {
			case <-done:
				return
			case pod, ok := <-pod:
				if !ok {
					return
				}
				ome.mu.Lock()
				ome.pod = pod
				ome.mu.Unlock()


			}
		}
	}()

	return ome
}



func (engine *Ome) OpenOrder(fragment order.Fragment) error {
	return engine.storer.Insert(fragment)
}

func (engine *Ome) CancelOrder(fragment order.Fragment) error {
	return engine.storer.Delete(fragment.OrderID)
}

func (engine *Ome) SyncRenLedger() error {
	return engine.syncer.SyncRenLedger(engine.done, engine.ranker)
}

func (engine *Ome) SyncRenLedger() error {
	done := make(chan struct{})
	defer close(done)

	go func() {
		// todo : handle error

	}()

	go func() {
		pairs := []OrderPair{}
		for {
			select {
			case <-done:
				return
			default:
				// todo : check if we have the order Fragments
				// if so, call smpc to generate the delta fragment
				// store the delta fragemtn , check if we have enough thredshold
				// if so check if its a match or not ,
				// confirm order if it's a match .

				orderPairs := engine.ranker.Get(50)
				for _, pair := range orderPairs {
					fragment, err := engine.storer.Get(pair.orderID)
					if err != nil {
						pairs = append(pairs, pair)
						continue
					}
					// We'll support 1 to multi matches in the future
					if len(pair.matches) != 1 {
						return
					}
					matchFragment, err := engine.storer.Get(pair.matches[0])
					if err != nil {
						pairs = append(pairs, pair)
						continue
					}

					delta, err := engine.smpcer.LessThan(fragment, matchFragment)
					if err != nil {
						return
					}

					engine.deltas[delta.DeltaID] = append(engine.deltas[delta.DeltaID], delta)
					threshold := engine.syncer.Threshold()
					if threshold == -1 {
						return
					}

					if len(engine.deltas[delta.DeltaID]) > threshold {
						delta, err := engine.smpcer.Join(engine.deltas[delta.DeltaID]...)
						if err != nil {
							return
						}
						// Fixme : delta.IsMatch shoudl take an uint64
						if delta.IsMatch(shamir.Prime) {
							err = engine.ConfirmOrder(delta.BuyOrderID, []order.ID{delta.SellOrderID})
							if err != nil {
								return
							}
							err = engine.delegate.ConfirmOrder(delta.BuyOrderID, []order.ID{delta.SellOrderID})
							if err != nil {
								return
							}

							engine.ranker.Remove(delta.BuyOrderID)
							engine.ranker.Remove(delta.SellOrderID)
						}
					}
				}
			}
		}
	}()

	return nil
}


func (engine *Ome) Compute(done <-chan struct{}) error {

}

func (engine *Ome) ConfirmOrder(id order.ID, matches []order.ID) error {
	if err := engine.syncer.ConfirmOrder(id, matches); err != nil {
		return err
	}

	// Remove the order from the ranker
	engine.ranker.Remove(id)
	for _, match := range matches {
		engine.ranker.Remove(match)
	}

	return nil
}

type Delegate interface {
	OnConfirmOrder(order order.Order, matches []order.Order) error
}
