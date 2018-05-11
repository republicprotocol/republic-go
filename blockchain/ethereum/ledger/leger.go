package ledger

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/bindings"
	"github.com/republicprotocol/republic-go/order"
	"math/big"
)

// RenLedgerContract
type RenLedgerContract struct {
	network      ethereum.Network
	context      context.Context
	conn         ethereum.Conn
	transactOpts *bind.TransactOpts
	callOpts     *bind.CallOpts
	binding      *bindings.RenLedger
	tokenBinding *bindings.RepublicToken
	address      common.Address
}

// NewRenLedgerContract creates a new NewRenLedgerContract with given parameters.
func NewRenLedgerContract(ctx context.Context, conn ethereum.Conn, transactOpts *bind.TransactOpts, callOpts *bind.CallOpts) (RenLedgerContract, error) {
	contract, err := bindings.NewRenLedger(common.HexToAddress(conn.Config.RenLedgerAddress), bind.ContractBackend(conn.Client))
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

func (ledger *RenLedgerContract) OpenOrder(signature []byte, id order.ID) error {
	var orderID [32]byte
	copy(orderID[:], id[:])

	tx, err := ledger.binding.OpenOrder(ledger.transactOpts, signature, orderID)
	if err != nil {
		return err
	}
	_, err = ledger.conn.PatchedWaitMined(ledger.context, tx)
	return err
}

func (ledger *RenLedgerContract) CancelOrder(signature []byte, id order.ID) error {
	var orderID [32]byte
	copy(orderID[:], id[:])

	tx, err := ledger.binding.CancelOrder(ledger.transactOpts, signature, orderID)
	if err != nil {
		return err
	}
	_, err = ledger.conn.PatchedWaitMined(ledger.context, tx)
	return err
}

func (ledger *RenLedgerContract) ConfirmOrder(id order.ID, matches []order.ID) error {
	orderMatches := make ([][32]byte, len(matches))
	for i := range orderMatches {
		copy (orderMatches[i][:], matches[i][:])
	}
	var orderID [32]byte
	copy(orderID[:], id[:])

	tx, err := ledger.binding.ConfirmOrder(ledger.transactOpts, orderID, orderMatches)
	if err != nil {
		return err
	}
	_, err = ledger.conn.PatchedWaitMined(ledger.context, tx)
	return err
}

func (ledger *RenLedgerContract) Priority(id order.ID) (uint64, error) {
	var orderID [32]byte
	copy(orderID[:], id[:])
	priority , err := ledger.binding.OrderPriority(ledger.callOpts, orderID)
	if err != nil {
		return 0, err
	}

	return priority.Uint64(), nil
}

func (ledger *RenLedgerContract) Status(id order.ID) (order.Status, error) {
	var orderID [32]byte
	copy(orderID[:], id[:])
	state , err := ledger.binding.OrderState(ledger.callOpts, orderID)
	if err != nil {
		return order.Nil, err
	}

	return order.Status(state), nil
}

func (ledger *RenLedgerContract) Matches(id order.ID) ([]order.ID, error) {
	var orderID [32]byte
	copy(orderID[:], id[:])
	matches , err := ledger.binding.OrderMatch(ledger.callOpts, orderID)
	if err != nil {
		return nil, err
	}
	orderIDs := make ([]order.ID, len(matches))
	for i := range matches{
		orderIDs[i] = matches[i][:]
	}

	return orderIDs, nil
}

func (ledger *RenLedgerContract) Orderbook(index int) (order.ID, error) {
	i := big.NewInt(int64(index))
	id , err := ledger.binding.Orderbook(ledger.callOpts, i)
	if err != nil {
		return nil, err
	}

	return order.ID(id[:]), nil
}

func (ledger *RenLedgerContract) Trader(id order.ID) (common.Address, error) {
	var orderID [32]byte
	copy(orderID[:], id[:])
	address , err := ledger.binding.OrderTrader(ledger.callOpts, orderID)
	if err != nil {
		return common.Address{}, err
	}

	return address, nil
}

func (ledger *RenLedgerContract) Broker(id order.ID) (common.Address, error) {
	var orderID [32]byte
	copy(orderID[:], id[:])
	address , err := ledger.binding.OrderBroker(ledger.callOpts, orderID)
	if err != nil {
		return common.Address{}, err
	}

	return address, nil
}

func (ledger *RenLedgerContract) Confirmer(id order.ID) (common.Address, error) {
	var orderID [32]byte
	copy(orderID[:], id[:])
	address , err := ledger.binding.OrderConfirmer(ledger.callOpts, orderID)
	if err != nil {
		return common.Address{}, err
	}

	return address, nil
}