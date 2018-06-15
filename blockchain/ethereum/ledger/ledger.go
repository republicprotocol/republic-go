package ledger

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/bindings"
	"github.com/republicprotocol/republic-go/order"
)

// ErrMismatchedOrderLengths is returned when there isn't an equal number of order IDs,
// signatures and order parities
var ErrMismatchedOrderLengths = errors.New("mismatched order lengths")

// BlocksForConfirmation is the number of Ethereum blocks required to consider
// changes to an order's status (Open, Canceled or Confirmed) in the Ledger to
// be confirmed. The functions `OpenBuyOrder`, `OpenSellOrder`, `CancelOrder`
// and `ConfirmOrder` return only after the required number of confirmations has
// been reached.
const BlocksForConfirmation = 1

// FIXME: Protect me with a mutex!
// RenLedgerContract implements the cal.RenLedger interface
type RenLedgerContract struct {
	network      ethereum.Network
	context      context.Context
	conn         ethereum.Conn
	transactOpts *bind.TransactOpts
	callOpts     *bind.CallOpts
	binding      *bindings.Orderbook
	address      common.Address
}

// NewRenLedgerContract creates a new NewRenLedgerContract with given parameters.
func NewRenLedgerContract(ctx context.Context, conn ethereum.Conn, transactOpts *bind.TransactOpts, callOpts *bind.CallOpts) (RenLedgerContract, error) {
	contract, err := bindings.NewOrderbook(common.HexToAddress(conn.Config.RenLedgerAddress), bind.ContractBackend(conn.Client))
	if err != nil {
		return RenLedgerContract{}, err
	}

	return RenLedgerContract{
		network:      conn.Config.Network,
		context:      ctx,
		conn:         conn,
		transactOpts: transactOpts,
		callOpts:     callOpts,
		binding:      contract,
		address:      common.HexToAddress(conn.Config.RenLedgerAddress),
	}, nil
}

// OpenOrders implements the cal.RenLedger interface
func (ledger *RenLedgerContract) OpenOrders(signatures [][65]byte, orderIDs []order.ID, orderParities []order.Parity) (int, error) {
	if len(signatures) != len(orderIDs) || len(signatures) != len(orderParities) {
		return 0, ErrMismatchedOrderLengths
	}

	nonce, err := ledger.conn.Client.PendingNonceAt(context.Background(), ledger.transactOpts.From)
	if err != nil {
		return 0, err
	}

	txs := make([]*types.Transaction, 0, len(signatures))
	for i := range signatures {
		ledger.transactOpts.GasPrice = big.NewInt(int64(5000000000))
		ledger.transactOpts.Nonce.Add(big.NewInt(0).SetUint64(nonce), big.NewInt(int64(i)))

		var tx *types.Transaction
		if orderParities[i] == order.ParityBuy {
			tx, err = ledger.binding.OpenBuyOrder(ledger.transactOpts, signatures[i][:], orderIDs[i])
		} else {
			tx, err = ledger.binding.OpenSellOrder(ledger.transactOpts, signatures[i][:], orderIDs[i])
		}
		if err != nil {
			break
		}
		txs = append(txs, tx)
	}

	for i := range txs {
		err := ledger.waitForOrderDepth(txs[i], orderIDs[i])
		if err != nil {
			return i, err
		}
	}

	return len(txs), err
}

// OpenBuyOrder implements the cal.RenLedger interface.
func (ledger *RenLedgerContract) OpenBuyOrder(signature [65]byte, id order.ID) error {
	ledger.transactOpts.GasPrice = big.NewInt(int64(20000000000))

	tx, err := ledger.binding.OpenBuyOrder(ledger.transactOpts, signature[:], id)
	if err != nil {
		return err
	}

	return ledger.waitForOrderDepth(tx, id)
}

// OpenSellOrder implements the cal.RenLedger interface.
func (ledger *RenLedgerContract) OpenSellOrder(signature [65]byte, id order.ID) error {
	ledger.transactOpts.GasPrice = big.NewInt(int64(20000000000))
	ledger.transactOpts.GasLimit = 500000

	tx, err := ledger.binding.OpenSellOrder(ledger.transactOpts, signature[:], id)
	if err != nil {
		return err
	}
	ledger.transactOpts.GasLimit = 0

	return ledger.waitForOrderDepth(tx, id)
}

// CancelOrder implements the cal.RenLedger interface.
func (ledger *RenLedgerContract) CancelOrder(signature [65]byte, id order.ID) error {
	tx, err := ledger.binding.CancelOrder(ledger.transactOpts, signature[:], id)
	if err != nil {
		return err
	}
	_, err = ledger.conn.PatchedWaitMined(ledger.context, tx)
	if err != nil {
		return err
	}
	return nil
}

// ConfirmOrder implements the cal.RenLedger interface
func (ledger *RenLedgerContract) ConfirmOrder(id order.ID, match order.ID) error {
	orderMatches := [][32]byte{match}

	before, err := ledger.binding.OrderDepth(ledger.callOpts, id)
	if err != nil {
		return err
	}
	tx, err := ledger.binding.ConfirmOrder(ledger.transactOpts, [32]byte(id), orderMatches)
	if err != nil {
		return err
	}
	_, err = ledger.conn.PatchedWaitMined(ledger.context, tx)
	if err != nil {
		return err
	}

	for {
		depth, err := ledger.binding.OrderDepth(ledger.callOpts, id)
		if err != nil {
			return err
		}

		if depth.Uint64()-before.Uint64() >= BlocksForConfirmation {
			return nil
		}
	}
}

// Priority implements the cal.RenLedger interface
func (ledger *RenLedgerContract) Priority(id order.ID) (uint64, error) {
	priority, err := ledger.binding.OrderPriority(ledger.callOpts, id)
	if err != nil {
		return 0, err
	}

	return priority.Uint64(), nil
}

// Status implements the cal.RenLedger interface
func (ledger *RenLedgerContract) Status(id order.ID) (order.Status, error) {
	var orderID [32]byte
	copy(orderID[:], id[:])
	state, err := ledger.binding.OrderState(ledger.callOpts, orderID)
	if err != nil {
		return order.Nil, err
	}

	return order.Status(state), nil
}

// OrderMatch implements the cal.RenLedger interface
func (ledger *RenLedgerContract) OrderMatch(id order.ID) (order.ID, error) {

	matches, err := ledger.binding.OrderMatch(ledger.callOpts, [32]byte(id))
	if err != nil {
		return order.ID{}, err
	}
	orderIDs := make([]order.ID, len(matches))
	for i := range matches {
		orderIDs[i] = matches[i]
	}
	if len(orderIDs) != 1 {
		return order.ID{}, errors.New("no matches found for the order")
	}

	return orderIDs[0], nil
}

// BuyOrders implements the cal.RenLedger interface
func (ledger *RenLedgerContract) BuyOrders(offset, limit int) ([]order.ID, error) {
	orders := make([]order.ID, 0, limit)
	for i := 0; i < limit; i++ {
		ordID, ok, err := ledger.binding.BuyOrder(ledger.callOpts, big.NewInt(int64(offset+i)))
		if !ok {
			return orders, nil
		}
		if err != nil {
			return nil, err
		}

		orders = append(orders, ordID)
	}
	return orders, nil
}

// SellOrders implements the cal.RenLedger interface
func (ledger *RenLedgerContract) SellOrders(offset, limit int) ([]order.ID, error) {
	orders := make([]order.ID, 0, limit)
	for i := 0; i < limit; i++ {
		ordID, ok, err := ledger.binding.SellOrder(ledger.callOpts, big.NewInt(int64(offset+i)))
		if !ok {
			return orders, nil
		}
		if err != nil {
			return nil, err
		}

		orders = append(orders, ordID)
	}

	return orders, nil
}

// Trader implements the cal.RenLedger interface
func (ledger *RenLedgerContract) Trader(id order.ID) (string, error) {
	address, err := ledger.binding.OrderTrader(ledger.callOpts, id)
	if err != nil {
		return "", err
	}

	return address.String(), nil
}

// Broker returns the address of the broker who submitted the order
func (ledger *RenLedgerContract) Broker(id order.ID) (common.Address, error) {
	address, err := ledger.binding.OrderBroker(ledger.callOpts, id)
	if err != nil {
		return common.Address{}, err
	}

	return address, nil
}

// Confirmer returns the address of the confirmer who submitted the order
func (ledger *RenLedgerContract) Confirmer(id order.ID) (common.Address, error) {
	address, err := ledger.binding.OrderConfirmer(ledger.callOpts, id)
	if err != nil {
		return common.Address{}, err
	}

	return address, nil
}

// Fee implements the cal.RenLedger interface
func (ledger *RenLedgerContract) Fee() (*big.Int, error) {
	return big.NewInt(0), nil
}

// Depth implements the cal.RenLedger interface
func (ledger *RenLedgerContract) Depth(orderID order.ID) (uint, error) {

	depth, err := ledger.binding.OrderDepth(ledger.callOpts, orderID)
	if err != nil {
		return 0, err
	}

	return uint(depth.Uint64()), nil
}

// BlockNumber implements the cal.RenLedger interface
func (ledger *RenLedgerContract) BlockNumber(orderID order.ID) (uint, error) {
	blockNumber, err := ledger.binding.OrderBlockNumber(ledger.callOpts, orderID)
	if err != nil {
		return 0, err
	}

	return uint(blockNumber.Uint64()), nil
}

// OrderCounts returns the total number of orders in the ledger
func (ledger *RenLedgerContract) OrderCounts() (uint64, error) {
	counts, err := ledger.binding.GetOrdersCount(ledger.callOpts)
	if err != nil {
		return 0, err
	}

	return counts.Uint64(), nil
}

// OrderID returns the order at a given index in the ledger
func (ledger *RenLedgerContract) OrderID(index int) ([32]byte, error) {
	i := big.NewInt(int64(index))
	id, exist, err := ledger.binding.GetOrder(ledger.callOpts, i)
	if !exist {
		return [32]byte{}, errors.New("order not exist")
	}
	if err != nil {
		return [32]byte{}, err
	}

	return id, nil
}

func (ledger *RenLedgerContract) waitForOrderDepth(tx *types.Transaction, id order.ID) error {
	_, err := ledger.conn.PatchedWaitMined(ledger.context, tx)
	if err != nil {
		return err
	}

	for {
		depth, err := ledger.binding.OrderDepth(ledger.callOpts, id)
		if err != nil {
			return err
		}

		if depth.Uint64() >= BlocksForConfirmation {
			return nil
		}
		time.Sleep(time.Second * 14)
	}
}
