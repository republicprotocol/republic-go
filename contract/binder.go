package contract

import (
	"bytes"
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/republicprotocol/republic-go/contract/bindings"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/registry"
	"github.com/republicprotocol/republic-go/stackint"
)

const EthereumAddress = "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"

// ErrPodNotFound is returned when dark node address was not found in any pod
var ErrPodNotFound = errors.New("cannot find node in any pod")

// ErrLengthMismatch is returned when ID is not an expected 20 byte value
var ErrLengthMismatch = errors.New("length mismatch")

// ErrMismatchedOrderLengths is returned when there isn't an equal number of order IDs,
// signatures and order parities
var ErrMismatchedOrderLengths = errors.New("mismatched order lengths")

// BlocksForConfirmation is the number of Ethereum blocks required to consider
// changes to an order's status (Open, Canceled or Confirmed) in the orderbook to
// be confirmed. The functions `OpenBuyOrder`, `OpenSellOrder`, `CancelOrder`
// and `ConfirmOrder` return only after the required number of confirmations has
// been reached.
const BlocksForConfirmation = 1

// Binder implements all methods that will communicate with the smart contracts
type Binder struct {
	mu           *sync.RWMutex
	network      Network
	conn         Conn
	transactOpts *bind.TransactOpts
	callOpts     *bind.CallOpts

	republicToken    *bindings.RepublicToken
	darknodeRegistry *bindings.DarknodeRegistry
	orderbook        *bindings.Orderbook

	settlementRegistry *bindings.SettlementRegistry
	renExSettlement    *bindings.Settlement
}

// NewBinder returns a Binder to communicate with contracts
func NewBinder(auth *bind.TransactOpts, conn Conn) (Binder, error) {
	transactOpts := *auth
	transactOpts.GasPrice = big.NewInt(5000000000)

	nonce, err := conn.Client.PendingNonceAt(context.Background(), transactOpts.From)
	if err != nil {
		return Binder{}, err
	}
	transactOpts.Nonce = big.NewInt(int64(nonce))

	darknodeRegistry, err := bindings.NewDarknodeRegistry(common.HexToAddress(conn.Config.DarknodeRegistryAddress), bind.ContractBackend(conn.Client))
	if err != nil {
		fmt.Println(fmt.Errorf("cannot bind to DarknodeRegistry: %v", err))
		return Binder{}, err
	}

	republicToken, err := bindings.NewRepublicToken(common.HexToAddress(conn.Config.RepublicTokenAddress), bind.ContractBackend(conn.Client))
	if err != nil {
		fmt.Println(fmt.Errorf("cannot bind to RepublicToken: %v", err))
		return Binder{}, err
	}

	orderbook, err := bindings.NewOrderbook(common.HexToAddress(conn.Config.OrderbookAddress), bind.ContractBackend(conn.Client))
	if err != nil {
		fmt.Println(fmt.Errorf("cannot bind to Orderbook: %v", err))
		return Binder{}, err
	}

	settlementRegistry, err := bindings.NewSettlementRegistry(common.HexToAddress(conn.Config.SettlementRegistryAddress), bind.ContractBackend(conn.Client))
	if err != nil {
		fmt.Println(fmt.Errorf("cannot bind to SettlementRegistry: %v", err))
		return Binder{}, err
	}

	renExSettlementAddress, err := settlementRegistry.SettlementContract(&bind.CallOpts{}, uint64(order.SettlementRenEx))
	if err != nil {
		fmt.Println(fmt.Errorf("cannot bind to RenExSettlementAddress: %v", err))
		return Binder{}, err
	}

	renExSettlement, err := bindings.NewSettlement(renExSettlementAddress, bind.ContractBackend(conn.Client))
	if err != nil {
		fmt.Println(fmt.Errorf("cannot bind to RenExSettlement: %v", err))
		return Binder{}, err
	}

	return Binder{
		mu:           new(sync.RWMutex),
		network:      conn.Config.Network,
		conn:         conn,
		transactOpts: &transactOpts,
		callOpts:     &bind.CallOpts{},

		republicToken:    republicToken,
		darknodeRegistry: darknodeRegistry,
		orderbook:        orderbook,

		settlementRegistry: settlementRegistry,
		renExSettlement:    renExSettlement,
	}, nil
}

// From returns the common.Address used to submit transactions.
func (binder *Binder) From() common.Address {
	binder.mu.RLock()
	defer binder.mu.RUnlock()
	return binder.transactOpts.From
}

// SendTx locks binder resources to execute function f (handling nonces explicitly)
// and will wait until the block has been mined on the blockchain. This will allow
// parallel requests to the blockchain since the binder will be unlocked before
// waiting for transaction to complete execution on the blockchain.
func (binder *Binder) SendTx(f func() (*types.Transaction, error)) (*types.Transaction, error) {
	binder.mu.Lock()
	defer binder.mu.Unlock()

	return binder.sendTx(f)
}

func (binder *Binder) sendTx(f func() (*types.Transaction, error)) (*types.Transaction, error) {
	tx, err := f()
	if err == nil {
		binder.transactOpts.Nonce.Add(binder.transactOpts.Nonce, big.NewInt(1))
		return tx, nil
	}
	if err == core.ErrNonceTooLow || err == core.ErrReplaceUnderpriced || strings.Contains(err.Error(), "nonce is too low") {
		binder.transactOpts.Nonce.Add(binder.transactOpts.Nonce, big.NewInt(1))
		return binder.sendTx(f)
	}
	if err == core.ErrNonceTooHigh {
		binder.transactOpts.Nonce.Sub(binder.transactOpts.Nonce, big.NewInt(1))
		return binder.sendTx(f)
	}

	// If any other type of nonce error occurs we will refresh the nonce and
	// try again for up to 1 minute
	var nonce uint64
	for try := 0; try < 60 && strings.Contains(err.Error(), "nonce"); try++ {
		time.Sleep(time.Second)
		nonce, err = binder.conn.Client.PendingNonceAt(context.Background(), binder.transactOpts.From)
		if err != nil {
			continue
		}
		binder.transactOpts.Nonce = big.NewInt(int64(nonce))
		if tx, err = f(); err == nil {
			binder.transactOpts.Nonce.Add(binder.transactOpts.Nonce, big.NewInt(1))
			return tx, nil
		}
	}

	return tx, err
}

// SubmitOrder to the RenEx accounts
func (binder *Binder) SubmitOrder(ord order.Order) error {
	tx, err := binder.SendTx(func() (*types.Transaction, error) {
		return binder.submitOrder(ord)
	})
	if err != nil {
		return err
	}

	_, err = binder.conn.PatchedWaitMined(context.Background(), tx)
	return err
}

func (binder *Binder) submitOrder(ord order.Order) (*types.Transaction, error) {
	// If the gas price is greater than the gas price limit, temporarily lower
	// the gas price for this request
	lastGasPrice := binder.transactOpts.GasPrice
	submitOrderGasPriceLimit, err := binder.renExSettlement.SubmitOrderGasPriceLimit(binder.callOpts)
	if err == nil {
		// Set gas price to the appropriate limit
		binder.transactOpts.GasPrice = submitOrderGasPriceLimit
		// Reset gas price
		defer func() {
			binder.transactOpts.GasPrice = lastGasPrice
		}()
	}

	log.Printf("[info] (submit order) order = %v, tokens = %v", ord.ID, ord.Tokens)

	return binder.renExSettlement.SubmitOrder(binder.transactOpts, ord.PrefixHash(), uint64(ord.Settlement), uint64(ord.Tokens), big.NewInt(0).SetUint64(ord.Price), big.NewInt(0).SetUint64(ord.Volume), big.NewInt(0).SetUint64(ord.MinimumVolume))
}

// SubmitMatch will submit a matched order pair to the RenEx accounts
func (binder *Binder) SubmitMatch(buy, sell order.ID) error {
	tx, err := binder.SendTx(func() (*types.Transaction, error) {
		return binder.submitMatch(buy, sell)
	})
	if err != nil {
		return err
	}

	_, err = binder.conn.PatchedWaitMined(context.Background(), tx)
	return err
}

func (binder *Binder) submitMatch(buy, sell order.ID) (*types.Transaction, error) {
	log.Printf("[info] (submit match) buy = %v, sell = %v", buy, sell)
	return binder.renExSettlement.Settle(binder.transactOpts, buy, sell)
}

// Settle the order pair which gets confirmed by the Orderbook
func (binder *Binder) Settle(buy order.Order, sell order.Order) error {
	binder.mu.Lock()
	defer binder.mu.Unlock()

	// TODO: Do we need to be able to check the Settlement contract for the
	// order status, or can we rely on Infura to block transactions that are
	// known to fail?

	var wg sync.WaitGroup

	// Submit buy order
	if sendTx, sendTxErr := binder.sendTx(func() (*types.Transaction, error) {
		return binder.submitOrder(buy)
	}); sendTxErr == nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, waitErr := binder.conn.PatchedWaitMined(context.Background(), sendTx)
			if waitErr != nil {
				log.Printf("[error] (settle) cannot wait to submit buy = %v: %v", buy.ID, waitErr)
			}
		}()
	} else {
		log.Printf("[error] (settle) cannot submit buy = %v: %v", buy.ID, sendTxErr)
	}

	// Submit sell order
	if sendTx, sendTxErr := binder.sendTx(func() (*types.Transaction, error) {
		return binder.submitOrder(sell)
	}); sendTxErr == nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, waitErr := binder.conn.PatchedWaitMined(context.Background(), sendTx)
			if waitErr != nil {
				log.Printf("[error] (settle) cannot wait to submit sell = %v: %v", sell.ID, waitErr)
			}
		}()
	} else {
		log.Printf("[error] (settle) cannot submit sell = %v: %v", sell.ID, sendTxErr)
	}

	// Wait for both submitOrder calls to be mined
	wg.Wait()

	// Submit the match and wait for it to be mined
	tx, sendTxErr := binder.sendTx(func() (*types.Transaction, error) {
		return binder.submitMatch(buy.ID, sell.ID)
	})
	if sendTxErr != nil {
		return fmt.Errorf("cannot settle buy = %v, sell = %v: %v", buy.ID, sell.ID, sendTxErr)
	}
	_, waitErr := binder.conn.PatchedWaitMined(context.Background(), tx)
	if waitErr != nil {
		return fmt.Errorf("cannot wait to settle buy = %v, sell = %v: %v", buy.ID, sell.ID, waitErr)
	}
	return nil
}

// Register a new dark node with the dark node registrar
func (binder *Binder) Register(darknodeID []byte, publicKey []byte, bond *stackint.Int1024) error {
	darknodeIDByte, err := toByte(darknodeID)
	if err != nil {
		return err
	}
	tx, err := binder.SendTx(func() (*types.Transaction, error) {
		return binder.register(darknodeIDByte, publicKey, bond)
	})
	if err != nil {
		return err
	}

	_, err = binder.conn.PatchedWaitMined(context.Background(), tx)
	return err
}

func (binder *Binder) register(darknodeIDByte [20]byte, publicKey []byte, bond *stackint.Int1024) (*types.Transaction, error) {
	return binder.darknodeRegistry.Register(binder.transactOpts, darknodeIDByte, publicKey, bond.ToBigInt())
}

// Deregister an existing dark node.
func (binder *Binder) Deregister(darknodeID []byte) error {
	darknodeIDByte, err := toByte(darknodeID)
	if err != nil {
		return err
	}
	tx, err := binder.SendTx(func() (*types.Transaction, error) {
		return binder.deregister(darknodeIDByte)
	})
	if err != nil {
		return err
	}

	_, err = binder.conn.PatchedWaitMined(context.Background(), tx)
	return err
}

func (binder *Binder) deregister(darknodeIDByte [20]byte) (*types.Transaction, error) {
	return binder.darknodeRegistry.Deregister(binder.transactOpts, darknodeIDByte)
}

// Refund withdraws the bond. Must be called before reregistering.
func (binder *Binder) Refund(darknodeID []byte) error {
	darknodeIDByte, err := toByte(darknodeID)
	if err != nil {
		return err
	}
	tx, err := binder.SendTx(func() (*types.Transaction, error) {
		return binder.refund(darknodeIDByte)
	})
	if err != nil {
		return err
	}

	_, err = binder.conn.PatchedWaitMined(context.Background(), tx)
	return err
}

func (binder *Binder) refund(darknodeIDByte [20]byte) (*types.Transaction, error) {
	return binder.darknodeRegistry.Refund(binder.transactOpts, darknodeIDByte)
}

// GetBond retrieves the bond of an existing dark node
func (binder *Binder) GetBond(darknodeID []byte) (stackint.Int1024, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	return binder.getBond(darknodeID)
}

func (binder *Binder) getBond(darknodeID []byte) (stackint.Int1024, error) {
	darknodeIDByte, err := toByte(darknodeID)
	if err != nil {
		return stackint.Int1024{}, err
	}
	bond, err := binder.darknodeRegistry.GetDarknodeBond(binder.callOpts, darknodeIDByte)
	if err != nil {
		return stackint.Int1024{}, err
	}
	return stackint.FromBigInt(bond)
}

// IsRegistered returns true if the identity.Address is a current
// registered Darknode. Otherwise, it returns false.
func (binder *Binder) IsRegistered(darknodeAddr identity.Address) (bool, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	return binder.isRegistered(darknodeAddr)
}

func (binder *Binder) isRegistered(darknodeAddr identity.Address) (bool, error) {
	darknodeIDByte, err := toByte(darknodeAddr.ID())
	if err != nil {
		return false, err
	}
	// log.Println("registering", string(darknodeIDByte[:]))
	return binder.darknodeRegistry.IsRegistered(binder.callOpts, darknodeIDByte)
}

// IsDeregistered returns true if the node is deregistered
func (binder *Binder) IsDeregistered(darknodeID []byte) (bool, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	return binder.isDeregistered(darknodeID)
}

func (binder *Binder) isDeregistered(darknodeID []byte) (bool, error) {
	darknodeIDByte, err := toByte(darknodeID)
	if err != nil {
		return false, err
	}
	// log.Println("deregistering", string(darknodeIDByte[:]))
	return binder.darknodeRegistry.IsDeregistered(binder.callOpts, darknodeIDByte)
}

// ApproveRen doesn't actually talk to the DNR - instead it approves Ren to it
func (binder *Binder) ApproveRen(value *stackint.Int1024) error {
	tx, err := binder.SendTx(func() (*types.Transaction, error) {
		return binder.approveRen(value)
	})
	if err != nil {
		return err
	}

	_, err = binder.conn.PatchedWaitMined(context.Background(), tx)
	return err
}

func (binder *Binder) approveRen(value *stackint.Int1024) (*types.Transaction, error) {
	return binder.republicToken.Approve(binder.transactOpts, common.HexToAddress(binder.conn.Config.DarknodeRegistryAddress), value.ToBigInt())
}

// GetOwner gets the owner of the given dark node
func (binder *Binder) GetOwner(darknodeID []byte) (common.Address, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	return binder.getOwner(darknodeID)
}

func (binder *Binder) getOwner(darknodeID []byte) (common.Address, error) {
	darknodeIDByte, err := toByte(darknodeID)
	if err != nil {
		return common.Address{}, err
	}
	return binder.darknodeRegistry.GetDarknodeOwner(binder.callOpts, darknodeIDByte)
}

// PublicKey returns the RSA public key of the Darknode registered with the
// given identity.Address.
func (binder *Binder) PublicKey(darknodeAddr identity.Address) (rsa.PublicKey, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	return binder.publicKey(darknodeAddr)
}

func (binder *Binder) publicKey(darknodeAddr identity.Address) (rsa.PublicKey, error) {
	darknodeIDByte, err := toByte(darknodeAddr.ID())
	if err != nil {
		return rsa.PublicKey{}, err
	}
	pubKeyBytes, err := binder.darknodeRegistry.GetDarknodePublicKey(binder.callOpts, darknodeIDByte)
	if err != nil {
		return rsa.PublicKey{}, err
	}
	return crypto.RsaPublicKeyFromBytes(pubKeyBytes)
}

// Darknodes registered in the pod.
func (binder *Binder) Darknodes() (identity.Addresses, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	return binder.darknodes()
}

func (binder *Binder) darknodes() (identity.Addresses, error) {
	numDarknodesBig, err := binder.darknodeRegistry.NumDarknodes(binder.callOpts)
	if err != nil {
		return nil, err
	}
	numDarknodes := numDarknodesBig.Int64()
	darknodes := make(identity.Addresses, 0, numDarknodes)

	// Get the first 20 pods worth of darknodes
	nilValue := common.HexToAddress("0x0000000000000000000000000000000000000000")
	values, err := binder.darknodeRegistry.GetDarknodes(binder.callOpts, nilValue, big.NewInt(480))

	// Loop until all darknode have been loaded
	for {
		if err != nil {
			return nil, err
		}
		for _, value := range values {
			if bytes.Equal(value.Bytes(), nilValue.Bytes()) {
				// We are finished when a nil address is returned
				return darknodes, nil
			}
			darknodes = append(darknodes, identity.ID(value.Bytes()).Address())
		}
		lastValue := values[len(values)-1]
		values, err = binder.darknodeRegistry.GetDarknodes(binder.callOpts, lastValue, big.NewInt(480))
		if err != nil {
			return nil, err
		}
		// Skip the first value returned so that we do not duplicate values
		values = values[1:]
	}
}

func (binder *Binder) previousDarknodes() (identity.Addresses, error) {
	numDarknodesBig, err := binder.darknodeRegistry.NumDarknodesPreviousEpoch(binder.callOpts)
	if err != nil {
		return nil, err
	}
	numDarknodes := numDarknodesBig.Int64()
	darknodes := make(identity.Addresses, 0, numDarknodes)

	// Get the first 20 pods worth of darknodes
	nilValue := common.HexToAddress("0x0000000000000000000000000000000000000000")
	values, err := binder.darknodeRegistry.GetPreviousDarknodes(binder.callOpts, nilValue, big.NewInt(480))

	// Loop until all darknode have been loaded
	for {
		if err != nil {
			return nil, err
		}
		for _, value := range values {
			if bytes.Equal(value.Bytes(), nilValue.Bytes()) {
				// We are finished when a nil address is returned
				return darknodes, nil
			}
			darknodes = append(darknodes, identity.ID(value.Bytes()).Address())
		}
		lastValue := values[len(values)-1]
		values, err = binder.darknodeRegistry.GetPreviousDarknodes(binder.callOpts, lastValue, big.NewInt(480))
		if err != nil {
			return nil, err
		}
		// Skip the first value returned so that we do not duplicate values
		values = values[1:]
	}
}

// MinimumBond gets the minimum viable bond amount
func (binder *Binder) MinimumBond() (stackint.Int1024, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	return binder.minimumBond()
}

func (binder *Binder) minimumBond() (stackint.Int1024, error) {
	bond, err := binder.darknodeRegistry.MinimumBond(binder.callOpts)
	if err != nil {
		return stackint.Int1024{}, err
	}
	return stackint.FromBigInt(bond)
}

// MinimumEpochInterval returns the minimum number of seconds between
// epochs.
func (binder *Binder) MinimumEpochInterval() (*big.Int, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	return binder.minimumEpochInterval()
}

func (binder *Binder) minimumEpochInterval() (*big.Int, error) {
	return binder.darknodeRegistry.MinimumEpochInterval(binder.callOpts)
}

// MinimumPodSize gets the minimum pod size
func (binder *Binder) MinimumPodSize() (stackint.Int1024, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	return binder.minimumPodSize()
}

func (binder *Binder) minimumPodSize() (stackint.Int1024, error) {
	interval, err := binder.darknodeRegistry.MinimumPodSize(binder.callOpts)
	if err != nil {
		return stackint.Int1024{}, err
	}
	return stackint.FromBigInt(interval)
}

// Pods returns the Pod configuration for the current Epoch.
func (binder *Binder) Pods() ([]registry.Pod, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	epoch, err := binder.darknodeRegistry.CurrentEpoch(binder.callOpts)
	if err != nil {
		return []registry.Pod{}, err
	}
	darknodes, err := binder.darknodes()
	if err != nil {
		return []registry.Pod{}, err
	}

	return binder.pods(epoch.Epochhash, darknodes)
}

// PreviousPods returns the Pod configuration for the previous Epoch.
func (binder *Binder) PreviousPods() ([]registry.Pod, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	previousEpoch, err := binder.darknodeRegistry.PreviousEpoch(binder.callOpts)
	if err != nil {
		return []registry.Pod{}, err
	}
	previousDarknodes, err := binder.previousDarknodes()
	if err != nil {
		return []registry.Pod{}, err
	}

	return binder.pods(previousEpoch.Epochhash, previousDarknodes)
}

func (binder *Binder) pods(epochVal *big.Int, darknodeAddrs identity.Addresses) ([]registry.Pod, error) {

	numberOfNodesInPod, err := binder.minimumPodSize()
	if err != nil {
		return []registry.Pod{}, err
	}
	if len(darknodeAddrs) < int(numberOfNodesInPod.ToBigInt().Int64()) {
		return []registry.Pod{}, fmt.Errorf("degraded pod: expected at least %v addresses, got %v", int(numberOfNodesInPod.ToBigInt().Int64()), len(darknodeAddrs))
	}

	numberOfDarknodes := big.NewInt(int64(len(darknodeAddrs)))
	x := big.NewInt(0).Mod(epochVal, numberOfDarknodes)
	positionInOcean := make([]int, len(darknodeAddrs))
	for i := 0; i < len(darknodeAddrs); i++ {
		positionInOcean[i] = -1
	}
	pods := make([]registry.Pod, (len(darknodeAddrs) / int(numberOfNodesInPod.ToBigInt().Int64())))
	for i := 0; i < len(darknodeAddrs); i++ {
		isRegistered, err := binder.isRegistered(darknodeAddrs[x.Int64()])
		if err != nil {
			return []registry.Pod{}, err
		}
		for !isRegistered || positionInOcean[x.Int64()] != -1 {
			x.Add(x, big.NewInt(1))
			x.Mod(x, numberOfDarknodes)
			isRegistered, err = binder.isRegistered(darknodeAddrs[x.Int64()])
			if err != nil {
				return []registry.Pod{}, err
			}
		}
		positionInOcean[x.Int64()] = i
		podID := i % (len(darknodeAddrs) / int(numberOfNodesInPod.ToBigInt().Int64()))
		pods[podID].Darknodes = append(pods[podID].Darknodes, darknodeAddrs[x.Int64()])
		x.Mod(x.Add(x, epochVal), numberOfDarknodes)
	}

	for i := range pods {
		hashData := [][]byte{}
		for _, darknodeAddr := range pods[i].Darknodes {
			hashData = append(hashData, darknodeAddr.ID())
		}
		copy(pods[i].Hash[:], crypto.Keccak256(hashData...))
		pods[i].Position = i
	}
	return pods, nil
}

// Epoch returns the current Epoch which includes the Pod configuration.
func (binder *Binder) Epoch() (registry.Epoch, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	epoch, err := binder.darknodeRegistry.CurrentEpoch(binder.callOpts)
	if err != nil {
		return registry.Epoch{}, err
	}
	darknodeAddrs, err := binder.darknodes()
	if err != nil {
		return registry.Epoch{}, err
	}

	return binder.epoch(epoch, darknodeAddrs)
}

func (binder *Binder) EpochHash() ([32]byte, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	epoch, err := binder.darknodeRegistry.CurrentEpoch(binder.callOpts)
	if err != nil {
		return [32]byte{}, err
	}
	var res [32]byte
	copy(res[:], epoch.Epochhash.Bytes())

	return res, nil
}

// PreviousEpoch returns the previous Epoch which includes the Pod configuration.
func (binder *Binder) PreviousEpoch() (registry.Epoch, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	previousEpoch, err := binder.darknodeRegistry.PreviousEpoch(binder.callOpts)
	if err != nil {
		return registry.Epoch{}, err
	}
	previousDarknodes, err := binder.previousDarknodes()
	if err != nil {
		return registry.Epoch{}, err
	}

	return binder.epoch(previousEpoch, previousDarknodes)
}

func (binder *Binder) epoch(epoch struct {
	Epochhash   *big.Int
	Blocknumber *big.Int
}, darknodeAddrs identity.Addresses) (registry.Epoch, error) {
	blockInterval, err := binder.darknodeRegistry.MinimumEpochInterval(binder.callOpts)
	if err != nil {
		return registry.Epoch{}, err
	}

	var blockhash [32]byte
	copy(blockhash[:], epoch.Epochhash.Bytes())

	pods, err := binder.pods(epoch.Epochhash, darknodeAddrs)
	if err != nil {
		return registry.Epoch{}, err
	}

	darknodes, err := binder.darknodes()
	if err != nil {
		return registry.Epoch{}, err
	}

	return registry.Epoch{
		Hash:          blockhash,
		Pods:          pods,
		Darknodes:     darknodes,
		BlockNumber:   epoch.Blocknumber,
		BlockInterval: blockInterval,
	}, nil
}

// NextEpoch will try to turn the Epoch and returns the resulting Epoch. If
// the turning of the Epoch failed, the current Epoch is returned.
func (binder *Binder) NextEpoch() (registry.Epoch, error) {
	tx, err := binder.SendTx(func() (*types.Transaction, error) {
		return binder.nextEpoch()
	})
	if err != nil {
		return registry.Epoch{}, err
	}

	_, err = binder.conn.PatchedWaitMined(context.Background(), tx)
	if err != nil {
		return registry.Epoch{}, err
	}

	epoch, err := binder.darknodeRegistry.CurrentEpoch(binder.callOpts)
	if err != nil {
		return registry.Epoch{}, err
	}
	darknodeAddrs, err := binder.darknodes()
	if err != nil {
		return registry.Epoch{}, err
	}

	return binder.epoch(epoch, darknodeAddrs)
}

func (binder *Binder) nextEpoch() (*types.Transaction, error) {
	return binder.darknodeRegistry.Epoch(binder.transactOpts)
}

// Pod returns the Pod that contains the given identity.Address in the
// current Epoch. It returns ErrPodNotFound if the identity.Address is not
// registered in the current Epoch.
func (binder *Binder) Pod(addr identity.Address) (registry.Pod, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	return binder.pod(addr)
}

func (binder *Binder) pod(addr identity.Address) (registry.Pod, error) {
	epoch, err := binder.darknodeRegistry.CurrentEpoch(binder.callOpts)
	if err != nil {
		return registry.Pod{}, err
	}
	darknodeAddrs, err := binder.darknodes()
	if err != nil {
		return registry.Pod{}, err
	}

	pods, err := binder.pods(epoch.Epochhash, darknodeAddrs)
	if err != nil {
		return registry.Pod{}, err
	}

	for i := range pods {
		for j := range pods[i].Darknodes {
			if pods[i].Darknodes[j] == addr {
				return pods[i], nil
			}
		}
	}

	return registry.Pod{}, ErrPodNotFound
}

// OpenOrder on the Orderbook. The signature will be used to identify the broker
// that verifies the order. The order must be in an undefined state to be
// opened.
func (binder *Binder) OpenOrder(settlement order.Settlement, signature [65]byte, id order.ID) error {
	tx, err := binder.SendTx(func() (*types.Transaction, error) {
		return binder.openOrder(settlement, signature, id)
	})
	if err != nil {
		return err
	}

	return binder.waitForOrderDepth(tx, id, 0)
}

func (binder *Binder) openOrder(settlement order.Settlement, signature [65]byte, id order.ID) (*types.Transaction, error) {
	return binder.orderbook.OpenOrder(binder.transactOpts, uint64(settlement), signature[:], id)
}

// CancelOrder on the Orderbook. The signature will be used to verify that
// the request was created by the trader that owns the order. The order
// must be in the opened state to be canceled.
func (binder *Binder) CancelOrder(id order.ID) error {
	tx, err := binder.SendTx(func() (*types.Transaction, error) {
		return binder.cancelOrder(id)
	})
	if err != nil {
		return err
	}

	_, err = binder.conn.PatchedWaitMined(context.Background(), tx)
	return err
}

func (binder *Binder) cancelOrder(id order.ID) (*types.Transaction, error) {
	return binder.orderbook.CancelOrder(binder.transactOpts, id)
}

// ConfirmOrder match on the Orderbook.
func (binder *Binder) ConfirmOrder(id order.ID, match order.ID) error {
	before, err := binder.orderDepth(id)
	if err != nil {
		return err
	}

	tx, err := binder.SendTx(func() (*types.Transaction, error) {
		return binder.confirmOrder(id, match)
	})
	if err != nil {
		return err
	}

	return binder.waitForOrderDepth(tx, id, before.Uint64())
}

func (binder *Binder) orderDepth(id order.ID) (*big.Int, error) {
	return binder.orderbook.OrderDepth(binder.callOpts, id)
}

func (binder *Binder) confirmOrder(id order.ID, match order.ID) (*types.Transaction, error) {
	return binder.orderbook.ConfirmOrder(binder.transactOpts, [32]byte(id), [32]byte(match))
}

// Priority will return the priority of the order
func (binder *Binder) Priority(id order.ID) (uint64, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	return binder.priority(id)
}

func (binder *Binder) priority(id order.ID) (uint64, error) {
	priority, err := binder.orderbook.OrderPriority(binder.callOpts, id)
	if err != nil {
		return 0, err
	}

	return priority.Uint64(), nil
}

// Status will return the status of the order
func (binder *Binder) Status(id order.ID) (order.Status, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	return binder.status(id)
}

func (binder *Binder) status(id order.ID) (order.Status, error) {
	var orderID [32]byte
	copy(orderID[:], id[:])
	state, err := binder.orderbook.OrderState(binder.callOpts, orderID)
	if err != nil {
		return order.Nil, err
	}

	return order.Status(state), nil
}

// OrderMatch of an order, if any.
func (binder *Binder) OrderMatch(id order.ID) (order.ID, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	return binder.orderMatch(id)
}

func (binder *Binder) orderMatch(id order.ID) (order.ID, error) {
	match, err := binder.orderbook.OrderMatch(binder.callOpts, [32]byte(id))
	if err != nil {
		return order.ID{}, err
	}
	return order.ID(match), nil
}

// Orders in the Orderbook starting at an offset and returning limited
// numbers of orders.
func (binder *Binder) Orders(offset, limit int) ([]order.ID, []order.Status, []string, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	return binder.orders(offset, limit)
}

func (binder *Binder) orders(offset, limit int) ([]order.ID, []order.Status, []string, error) {
	orderIDsBytes, tradersAddrs, orderStatusesUInt8, err := binder.orderbook.GetOrders(binder.callOpts, big.NewInt(int64(offset)), big.NewInt(int64(limit)))
	if err != nil {
		return nil, nil, nil, err
	}

	orderIDs := make([]order.ID, len(orderIDsBytes))
	orderStatuses := make([]order.Status, len(orderStatusesUInt8))
	traders := make([]string, len(tradersAddrs))

	for i := range orderIDs {
		orderIDs[i] = orderIDsBytes[i]
	}
	for i := range orderStatuses {
		orderStatuses[i] = order.Status(orderStatusesUInt8[i])
	}
	for i := range traders {
		traders[i] = tradersAddrs[i].String()
	}

	return orderIDs, orderStatuses, traders, nil
}

// Trader returns the trader who submits the order
func (binder *Binder) Trader(id order.ID) (string, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	return binder.trader(id)
}

func (binder *Binder) trader(id order.ID) (string, error) {
	address, err := binder.orderbook.OrderTrader(binder.callOpts, id)
	if err != nil {
		return "", err
	}

	return address.String(), nil
}

// Confirmer returns the address of the confirmer who submitted the order
func (binder *Binder) Confirmer(id order.ID) (common.Address, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	return binder.confirmer(id)
}

func (binder *Binder) confirmer(id order.ID) (common.Address, error) {
	address, err := binder.orderbook.OrderConfirmer(binder.callOpts, id)
	if err != nil {
		return common.Address{}, err
	}

	return address, nil
}

// Fee required to open an order.
func (binder *Binder) Fee() (*big.Int, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	return binder.fee()
}

func (binder *Binder) fee() (*big.Int, error) {
	return big.NewInt(0), nil
}

// Depth will return depth of confirmation blocks
func (binder *Binder) Depth(orderID order.ID) (uint, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	return binder.depth(orderID)
}

func (binder *Binder) depth(orderID order.ID) (uint, error) {
	depth, err := binder.orderbook.OrderDepth(binder.callOpts, orderID)
	if err != nil {
		return 0, err
	}

	return uint(depth.Uint64()), nil
}

// BlockNumber will return the block number when the order status
// last mode modified
func (binder *Binder) BlockNumber(orderID order.ID) (*big.Int, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	return binder.blockNumber(orderID)
}

func (binder *Binder) blockNumber(orderID order.ID) (*big.Int, error) {
	return binder.orderbook.OrderBlockNumber(binder.callOpts, orderID)
}

// OrderCounts returns the total number of orders in the orderbook
func (binder *Binder) OrderCounts() (uint64, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	return binder.orderCounts()
}

func (binder *Binder) orderCounts() (uint64, error) {
	counts, err := binder.orderbook.OrdersCount(binder.callOpts)
	if err != nil {
		return 0, err
	}

	return counts.Uint64(), nil
}

func (binder *Binder) waitForOrderDepth(tx *types.Transaction, id order.ID, before uint64) error {
	_, err := binder.conn.PatchedWaitMined(context.Background(), tx)
	if err != nil {
		return err
	}

	for {
		depth, err := binder.orderbook.OrderDepth(binder.callOpts, id)
		if err != nil {
			return err
		}

		if depth.Uint64()-before >= BlocksForConfirmation {
			return nil
		}
		time.Sleep(time.Second * 14)
	}
}

func toByte(id []byte) ([20]byte, error) {
	twentyByte := [20]byte{}
	if len(id) != 20 {
		return twentyByte, ErrLengthMismatch
	}
	for i := range id {
		twentyByte[i] = id[i]
	}
	return twentyByte, nil
}
