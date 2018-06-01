package ledger

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/bindings"
	"github.com/republicprotocol/republic-go/order"
)

const BlocksForConfirmation = 1

// RenLedgerContract
type RenLedgerContract struct {
	network      ethereum.Network
	context      context.Context
	conn         ethereum.Conn
	transactOpts *bind.TransactOpts
	callOpts     *bind.CallOpts
	binding      *bindings.RenLedger
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

func (ledger *RenLedgerContract) OpenBuyOrder(signature [65]byte, id order.ID) error {

	ledger.transactOpts.GasPrice = big.NewInt(int64(5000000000))
	tx, err := ledger.binding.OpenBuyOrder(ledger.transactOpts, signature[:], id)

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

		if depth.Uint64() >= BlocksForConfirmation {
			return nil
		}
	}
}

func (ledger *RenLedgerContract) OpenSellOrder(signature [65]byte, id order.ID) error {

	ledger.transactOpts.GasPrice = big.NewInt(int64(5000000000))
	tx, err := ledger.binding.OpenSellOrder(ledger.transactOpts, signature[:], id)

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

		if depth.Uint64() >= BlocksForConfirmation {
			return nil
		}
	}
}

func (ledger *RenLedgerContract) CancelOrder(signature [65]byte, id order.ID) error {

	before, err := ledger.binding.OrderDepth(ledger.callOpts, id)
	if err != nil {
		return err
	}
	tx, err := ledger.binding.CancelOrder(ledger.transactOpts, signature[:], id)
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

func (ledger *RenLedgerContract) Priority(id order.ID) (uint64, error) {
	priority, err := ledger.binding.OrderPriority(ledger.callOpts, id)
	if err != nil {
		return 0, err
	}

	return priority.Uint64(), nil
}

func (ledger *RenLedgerContract) Status(id order.ID) (order.Status, error) {
	var orderID [32]byte
	copy(orderID[:], id[:])
	state, err := ledger.binding.OrderState(ledger.callOpts, orderID)
	if err != nil {
		return order.Nil, err
	}

	return order.Status(state), nil
}

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

func (ledger *RenLedgerContract) Matches(id order.ID) ([]order.ID, error) {
	matches, err := ledger.binding.OrderMatch(ledger.callOpts, id)
	if err != nil {
		return nil, err
	}
	orderIDs := make([]order.ID, len(matches))
	for i := range matches {
		orderIDs[i] = matches[i]
	}

	return orderIDs, nil
}

func (ledger *RenLedgerContract) BuyOrders(offset, limit int) ([]order.ID, error) {
	orders := make([]order.ID, 0, limit)
	for i := 0; i < limit; i++ {
		ordId, ok, err := ledger.binding.BuyOrder(ledger.callOpts, big.NewInt(int64(offset+i)))
		if !ok {
			return orders, nil
		}
		if err != nil {
			return nil, err
		}

		orders = append(orders, ordId)
	}
	return orders, nil
}

func (ledger *RenLedgerContract) SellOrders(offset, limit int) ([]order.ID, error) {
	orders := make([]order.ID, 0, limit)
	for i := 0; i < limit; i++ {
		ordId, ok, err := ledger.binding.SellOrder(ledger.callOpts, big.NewInt(int64(offset+i)))
		if !ok {
			return orders, nil
		}
		if err != nil {
			return nil, err
		}

		orders = append(orders, ordId)
	}

	return orders, nil
}

func (ledger *RenLedgerContract) Trader(id order.ID) (string, error) {
	address, err := ledger.binding.OrderTrader(ledger.callOpts, id)
	if err != nil {
		return "", err
	}

	return address.String(), nil
}

func (ledger *RenLedgerContract) Broker(id order.ID) (common.Address, error) {
	address, err := ledger.binding.OrderBroker(ledger.callOpts, id)
	if err != nil {
		return common.Address{}, err
	}

	return address, nil
}

func (ledger *RenLedgerContract) Confirmer(id order.ID) (common.Address, error) {
	address, err := ledger.binding.OrderConfirmer(ledger.callOpts, id)
	if err != nil {
		return common.Address{}, err
	}

	return address, nil
}

func (ledger *RenLedgerContract) Fee() (*big.Int, error) {
	return big.NewInt(0), nil
}

// GetBlockNumberOfTx gets the block number of the transaction hash.
// It calls infura using json-RPC.
func (ledger *RenLedgerContract) GetBlockNumberOfTx(transaction common.Hash) (uint64, error) {
	switch ledger.conn.Config.Network {
	// According to this https://github.com/ethereum/go-ethereum/issues/15210
	// we have to use json-rpc to get the block number.
	case ethereum.NetworkRopsten:
		hash := []byte(`{"jsonrpc":"2.0","method":"eth_getTransactionByHash","params":["` + transaction.Hex() + `"],"id":1}`)
		resp, err := http.Post("https://ropsten.infura.io", "application/json", bytes.NewBuffer(hash))
		if err != nil {
			return 0, err
		}

		// Read the response status
		if resp.StatusCode != http.StatusOK {
			return 0, fmt.Errorf("request failed with status code %d", resp.StatusCode)
		}
		// Get the result
		var data map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return 0, err
		}
		if result, ok := data["result"]; ok {
			if blockNumber, ok := result.(map[string]interface{})["blockNumber"]; ok {
				if blockNumberStr, ok := blockNumber.(string); ok {
					return strconv.ParseUint(blockNumberStr[2:], 16, 64)
				}
			}
		}
		return 0, errors.New("fail to unmarshal the json response")
	}

	return 0, nil
}

// CurrentBlock gets the latest block.
func (ledger *RenLedgerContract) CurrentBlock() (*types.Block, error) {
	return ledger.conn.Client.BlockByNumber(ledger.context, nil)
}

// Depth will return depth of confirmation blocks
func (ledger *RenLedgerContract) Depth(orderID order.ID) (uint, error) {

	depth, err := ledger.binding.OrderDepth(ledger.callOpts, orderID)
	if err != nil {
		return 0, err
	}

	return uint(depth.Uint64()), nil
}
