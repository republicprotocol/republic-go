package hd

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
	"github.com/republicprotocol/republic-go/hyperdrive"
)

// HyperdriveContract defines the golang interface for interacting with the
// hyperdrive smart contract on Ethereum.
type HyperdriveContract struct {
	network           ethereum.Network
	context           context.Context
	client            ethereum.Conn
	callOpts          *bind.CallOpts
	transactOpts      *bind.TransactOpts
	binding           *bindings.Hyperdrive
	hyperdriveAddress common.Address
}

// NewHyperdriveContract creates a new HyperdriveContract with given parameters.
func NewHyperdriveContract(ctx context.Context, conn ethereum.Conn, transactOpts *bind.TransactOpts, callOpts *bind.CallOpts) (HyperdriveContract, error) {
	contract, err := bindings.NewHyperdrive(conn.HyperdriveAddress, bind.ContractBackend(conn.Client))
	if err != nil {
		return HyperdriveContract{}, err
	}

	return HyperdriveContract{
		network:           conn.Network,
		context:           ctx,
		client:            conn,
		callOpts:          callOpts,
		transactOpts:      transactOpts,
		binding:           contract,
		hyperdriveAddress: conn.HyperdriveAddress,
	}, nil
}

// SendTx sends an tx to the contract. It will block until the transaction has
// been mined. It returns an error if there is an conflict with previous txs.
// You need to register with the darkNodeRegistry to send the tx.
func (hyper HyperdriveContract) SendTx(tx hyperdrive.Tx) (*types.Transaction, error) {
	var hash [32]byte
	copy(hash[:], tx.Hash)

	nonces := make([][32]byte, len(tx.Nonces))
	for i := range nonces {
		copy(nonces[i][:], tx.Nonces[i])
	}

	transaction, err := hyper.binding.SendOrderMatch(hyper.transactOpts, hash, nonces)
	if err != nil {
		return nil, err
	}

	_, err = hyper.client.PatchedWaitMined(hyper.context, transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (hyper *HyperdriveContract) CheckOrders(nonce hyperdrive.Nonce) (uint64, error) {
	var nonceIn32Bytes [32]byte
	copy(nonceIn32Bytes[:], nonce)
	bn, err := hyper.binding.ConfirmedOrders(hyper.callOpts, nonceIn32Bytes)
	if err != nil {
		return 0, err
	}
	return bn.Uint64(), nil
}

// GetDepth read the depth of the nonce from the hyperdrive contract.
func (hyper *HyperdriveContract) GetDepth(nonce hyperdrive.Nonce) (uint64, error) {
	var nonceIn32Bytes [32]byte
	copy(nonceIn32Bytes[:], nonce)
	depth, err := hyper.binding.Depth(hyper.callOpts, nonceIn32Bytes)
	if err != nil {
		return 0, err
	}
	return depth.Uint64(), nil
}

// GetBlockNumberOfTx gets the block number of the transaction hash.
// It calls infura using json-RPC.
func (hyper HyperdriveContract) GetBlockNumberOfTx(transaction common.Hash) (uint64, error) {
	switch hyper.client.Network {
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
func (hyper HyperdriveContract) CurrentBlock() (*types.Block, error) {
	return hyper.client.Client.BlockByNumber(hyper.context, nil)
}

// BlockByNumber gets the block of given number.
func (hyper *HyperdriveContract) BlockByNumber(blockNumber *big.Int) (*types.Block, error) {
	return hyper.client.Client.BlockByNumber(hyper.context, blockNumber)
}

// SetGasLimit sets the gas limit to use for transactions.
func (hyper *HyperdriveContract) SetGasLimit(limit uint64) {
	hyper.transactOpts.GasLimit = limit
}
