package ome

import (
	"context"
	"errors"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/shamir"
	"github.com/republicprotocol/republic-go/smpc"
	"github.com/republicprotocol/republic-go/smpc/delta"
)

var ErrNotFoundInStore = errors.New("not found in store")

type Delegate interface {
	OnConfirmOrder(order order.Order, matches []order.Order) error
}

type OmeClient interface {
	OpenOrder(ctx context.Context, multiAddr identity.MultiAddress, orderFragment order.Fragment) error
}

type OmeServer interface {
	OpenOrder(ctx context.Context, orderFragment order.Fragment) error
}

type Omer interface {
	OnChangeEpoch(ξ cal.Epoch)
}

type Ome struct {
	done      <-chan struct{}
	epoch     <-chan cal.Epoch
	delegate  Delegate
	orderbook orderbook.Orderbooker
	ranker    Ranker
	smpcer    smpc.Smpcer
	deltas    map[delta.ID][]delta.Fragment
}

func NewOme(done <-chan struct{}, delegate Delegate, orderbook orderbook.Orderbooker, ranker Ranker, smpcer smpc.Smpcer) Ome {
	ome := Ome{}
	ome.done = done
	ome.epoch = make(chan cal.Epoch)
	ome.delegate = delegate
	ome.orderbook = orderbook
	ome.ranker = ranker
	ome.smpcer = smpcer
	ome.deltas = map[delta.ID][]delta.Fragment{}
	return ome
}

func (ome *Ome) OnChangeEpoch(ξ cal.Epoch) {
	ome.epoch <- ξ
}

func (ome *Ome) Sync() error {
	updates, err := ome.orderbook.Sync()
	if err != nil {
		return err
	}
	for _, update := range updates {
		switch update.Status {
		case order.Open:
			ome.ranker.Insert(update.ID, update.Parity, update.Priority)
		case order.Canceled, order.Settled, order.Confirmed:
			ome.ranker.Remove(update.ID)
		}
	}

	return nil
}

func (ome *Ome) Compute() chan error {
	errs := make(chan error)
	var currentEpoch cal.Epoch

	orderPairs := ome.ranker.OrderPairs(ome.done)
	input, output := ome.smpcer.Input(), ome.smpcer.Output()

	go func() {
		defer close(errs)

		waitingList := []OrderPair{}
		for {
			select {
			case <-ome.done:
				return
			case epoch, ok := <-ome.epoch:
				if !ok {
					return
				}
				currentEpoch = epoch
				input <- smpc.Inst{
					InstConnect: &smpc.InstConnect{
						PeersID: epoch.Hash[:],
						Peers:   epoch.Darknodes,
						N:       len(epoch.Darknodes),
						K:       (len(epoch.Darknodes) + 1) * 2 / 3,
					},
				}
			case pair, ok := <-orderPairs:
				// Try again
				if !ok {
					orderPairs = ome.ranker.OrderPairs(ome.done)
					continue
				}

				for i := 0; i < len(waitingList); i++ {
					//todo : need to have a way for element to expire
					inst, err := generateInstruction(waitingList[i], ome.orderbook, currentEpoch)
					if err != nil {
						continue
					}

					input <- inst
					waitingList = append(waitingList[:i], waitingList[i+1:]...)
					i--
				}

				inst, err := generateInstruction(pair, ome.orderbook, currentEpoch)
				if err != nil {
					waitingList = append(waitingList, pair)
					continue
				}
				input <- inst
			case result := <-output:
				if result.Delta.IsMatch(shamir.Prime) {
					// todo :talk to the orderbook ?
					err := ome.delegate.OnConfirmOrder()
					if err != nil {

					}
				}
			}
		}
	}()

	return errs
}

func generateInstruction(pair OrderPair, orderbook orderbook.Orderbooker, epoch cal.Epoch) (smpc.Inst, error) {
	// Get orderFragment from the orderbook
	buyOrderFragment, err := orderbook.OrderFragment(pair.BuyOrder)
	if err != nil {
		return smpc.Inst{}, err
	}
	sellOrderFragment, err := orderbook.OrderFragment(pair.SellOrder)
	if err != nil {
		return smpc.Inst{}, err
	}

	// generate instruction message for smpc
	return smpc.Inst{
		InstCompute: &smpc.InstCompute{
			PeersID: epoch.Hash[:],
			Buy:     buyOrderFragment,
			Sell:    sellOrderFragment,
		},
	}, nil
}
