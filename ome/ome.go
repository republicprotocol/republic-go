package ome

import (
	"context"
	"errors"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
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
	ledger    cal.RenLedger
	orderbook orderbook.Orderbooker
	ranker    Ranker
	smpcer    smpc.Smpcer
	deltas    map[delta.ID][]delta.Fragment
}

func NewOme(done <-chan struct{}, delegate Delegate, ledger cal.RenLedger, orderbook orderbook.Orderbooker, ranker Ranker, smpcer smpc.Smpcer) Ome {
	ome := Ome{}
	ome.done = done
	ome.epoch = make(chan cal.Epoch)
	ome.ledger = ledger
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
	var currentEpoch cal.Epoch
	errs := make(chan error)
	input, output := ome.smpcer.Instructions(), ome.smpcer.Results()

	err := ome.smpcer.Start()
	if err != nil {
		errs <- err
		close(errs)
		return errs
	}

	go func() {
		defer close(errs)
		defer ome.smpcer.Shutdown()

		deltaWaitingList := []delta.Delta{}

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
			case result, ok := <-output:
				if !ok {
					return
				}

				for i := 0; i < len(deltaWaitingList); i++ {
					buyOrd, sellOrd, err := getOrderDetails(deltaWaitingList[i], ome.orderbook)
					if err != nil {
						continue
					}
					err = ome.ledger.ConfirmOrder(result.Delta.BuyOrderID, []order.ID{result.Delta.SellOrderID})
					if err != nil {
						continue
					}
					err = ome.delegate.OnConfirmOrder(buyOrd, []order.Order{sellOrd})
					if err != nil {
						continue
					}
					deltaWaitingList = append(deltaWaitingList[:i], deltaWaitingList[i+1:]...)
					i--
				}

				if result.Delta.IsMatch() {
					buyOrd, sellOrd, err := getOrderDetails(result.Delta, ome.orderbook)
					if err != nil {
						deltaWaitingList = append(deltaWaitingList, result.Delta)
						continue
					}
					err = ome.ledger.ConfirmOrder(result.Delta.BuyOrderID, []order.ID{result.Delta.SellOrderID})
					if err != nil {
						deltaWaitingList = append(deltaWaitingList, result.Delta)
						continue
					}
					err = ome.delegate.OnConfirmOrder(buyOrd, []order.Order{sellOrd})
					if err != nil {
						deltaWaitingList = append(deltaWaitingList, result.Delta)
						continue
					}
				}
			default:
				orderPairs := ome.ranker.OrderPairs(50 )

				for _, pair := range orderPairs {
					inst, err := generateInstruction(pair, ome.orderbook, currentEpoch)
					if err != nil {
						continue
					}
					input <- inst
					ome.ranker.Remove(pair.SellOrder , pair.BuyOrder)
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

func getOrderDetails(delta delta.Delta, orderbooker orderbook.Orderbooker) (order.Order, order.Order, error) {
	buyOrd, err := orderbooker.Order(delta.BuyOrderID)
	if err != nil {
		return order.Order{}, order.Order{}, err
	}
	sellOrd, err := orderbooker.Order(delta.SellOrderID)
	if err != nil {
		return order.Order{}, order.Order{}, err
	}

	return buyOrd, sellOrd, nil
}
