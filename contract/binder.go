package contract

import (
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
	renExSettlement  *bindings.Settlement
	renExBalance     *bindings.RenExBalances
	erc20            *bindings.ERC20
}

// NewBinder returns a Binder to communicate with contracts
func NewBinder(auth *bind.TransactOpts, conn Conn) (Binder, error) {
	transactOpts := *auth
	transactOpts.GasPrice = big.NewInt(20000000000)

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

	renExSettlement, err := bindings.NewSettlement(common.HexToAddress(conn.Config.RenExSettlementAddress), bind.ContractBackend(conn.Client))
	if err != nil {
		fmt.Println(fmt.Errorf("cannot bind to RenExSettlement: %v", err))
		return Binder{}, err
	}

	renExBalance, err := bindings.NewRenExBalances(common.HexToAddress(conn.Config.RenExBalancesAddress), bind.ContractBackend(conn.Client))
	if err != nil {
		fmt.Println(fmt.Errorf("cannot bind to RenExBalance: %v", err))
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
		renExSettlement:  renExSettlement,
		renExBalance:     renExBalance,
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
	submissionGasPriceLimit, err := binder.renExSettlement.SubmissionGasPriceLimit(binder.callOpts)
	if err == nil {
		// Set gas price to the appropriate limit
		binder.transactOpts.GasPrice = submissionGasPriceLimit
		// Reset gas price
		defer func() {
			binder.transactOpts.GasPrice = lastGasPrice
		}()
	}

	nonceHash := big.NewInt(0).SetBytes(ord.BytesFromNonce())
	log.Printf("[info] (submit order) order = %v, tokens = %v", ord.ID, ord.Tokens)
	return binder.renExSettlement.SubmitOrder(binder.transactOpts, uint32(ord.Settlement), uint8(ord.Type), uint8(ord.Parity), uint64(ord.Expiry.Unix()), uint64(ord.Tokens), uint16(ord.Price.Co), uint16(ord.Price.Exp), uint16(ord.Volume.Co), uint16(ord.Volume.Exp), uint16(ord.MinimumVolume.Co), uint16(ord.MinimumVolume.Exp), nonceHash)
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
	return binder.renExSettlement.SubmitMatch(binder.transactOpts, buy, sell)
}

// Settle the order pair which gets confirmed by the Orderbook
func (binder *Binder) Settle(buy order.Order, sell order.Order) (err error) {
	binder.mu.Lock()
	defer binder.mu.Unlock()

	// Submit orders
	if _, sendTxErr := binder.sendTx(func() (*types.Transaction, error) {
		return binder.submitOrder(buy)
	}); sendTxErr != nil {
		err = fmt.Errorf("cannot settle buy = %v: %v", buy.ID, sendTxErr)
	}
	if _, sendTxErr := binder.sendTx(func() (*types.Transaction, error) {
		return binder.submitOrder(sell)
	}); sendTxErr != nil {
		err = fmt.Errorf("cannot settle sell = %v: %v", sell.ID, sendTxErr)
	}

	// Submit match
	tx, sendTxErr := binder.sendTx(func() (*types.Transaction, error) {
		return binder.submitMatch(buy.ID, sell.ID)
	})
	if sendTxErr != nil {
		err = fmt.Errorf("cannot settle buy = %v, sell = %v: %v", buy.ID, sell.ID, sendTxErr)
		return err
	}

	// Wait for last transaction
	if _, waitErr := binder.conn.PatchedWaitMined(context.Background(), tx); waitErr != nil {
		err = fmt.Errorf("cannot wait to settle buy = %v, sell = %v: %v", buy.ID, sell.ID, waitErr)
		return err
	}

	return err
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
	bond, err := binder.darknodeRegistry.GetBond(binder.callOpts, darknodeIDByte)
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
	pubKeyBytes, err := binder.darknodeRegistry.GetPublicKey(binder.callOpts, darknodeIDByte)
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
	ret, err := binder.darknodeRegistry.GetDarknodes(binder.callOpts)
	if err != nil {
		return nil, err
	}
	arr := make(identity.Addresses, len(ret))
	for i := range ret {
		arr[i] = identity.ID(ret[i][:]).Address()
	}
	return arr, nil
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

	return binder.pods(epoch.Epochhash)
}

// PreviousPods returns the Pod configuration for the previous Epoch.
func (binder *Binder) PreviousPods() ([]registry.Pod, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	previousEpoch, err := binder.darknodeRegistry.PreviousEpoch(binder.callOpts)
	if err != nil {
		return []registry.Pod{}, err
	}

	return binder.pods(previousEpoch.Epochhash)
}

func (binder *Binder) pods(epochVal *big.Int) ([]registry.Pod, error) {
	darknodeAddrs, err := binder.darknodes()
	if err != nil {
		return []registry.Pod{}, err
	}

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

	return binder.epoch(epoch)
}

// PreviousEpoch returns the previous Epoch which includes the Pod configuration.
func (binder *Binder) PreviousEpoch() (registry.Epoch, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	previousEpoch, err := binder.darknodeRegistry.PreviousEpoch(binder.callOpts)
	if err != nil {
		return registry.Epoch{}, err
	}

	return binder.epoch(previousEpoch)
}

func (binder *Binder) epoch(epoch struct {
	Epochhash   *big.Int
	Blocknumber *big.Int
}) (registry.Epoch, error) {
	blockInterval, err := binder.darknodeRegistry.MinimumEpochInterval(binder.callOpts)
	if err != nil {
		return registry.Epoch{}, err
	}

	var blockhash [32]byte
	copy(blockhash[:], epoch.Epochhash.Bytes())

	pods, err := binder.pods(epoch.Epochhash)
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

		// FIXME: Such a low gas price is only appropriate during testnet.
		previousGasPrice := binder.transactOpts.GasPrice
		binder.transactOpts.GasPrice = big.NewInt(1000000000)
		defer func() {
			binder.transactOpts.GasPrice = previousGasPrice
		}()

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

	return binder.epoch(epoch)
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

	pods, err := binder.pods(epoch.Epochhash)
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

// OpenBuyOrder on the Orderbook. The signature will be used to identify
// the trader that owns the order. The order must be in an undefined state
// to be opened.
func (binder *Binder) OpenBuyOrder(signature [65]byte, id order.ID) error {
	tx, err := binder.SendTx(func() (*types.Transaction, error) {
		return binder.openBuyOrder(signature, id)
	})
	if err != nil {
		return err
	}

	return binder.waitForOrderDepth(tx, id, 0)
}

func (binder *Binder) openBuyOrder(signature [65]byte, id order.ID) (*types.Transaction, error) {
	return binder.orderbook.OpenBuyOrder(binder.transactOpts, signature[:], id)
}

// OpenSellOrder on the Orderbook. The signature will be used to identify
// the trader that owns the order. The order must be in an undefined state
// to be opened.
func (binder *Binder) OpenSellOrder(signature [65]byte, id order.ID) error {
	tx, err := binder.SendTx(func() (*types.Transaction, error) {
		return binder.openSellOrder(signature, id)
	})
	if err != nil {
		return err
	}

	return binder.waitForOrderDepth(tx, id, 0)
}

func (binder *Binder) openSellOrder(signature [65]byte, id order.ID) (*types.Transaction, error) {
	return binder.orderbook.OpenSellOrder(binder.transactOpts, signature[:], id)
}

// CancelOrder on the Orderbook. The signature will be used to verify that
// the request was created by the trader that owns the order. The order
// must be in the opened state to be canceled.
func (binder *Binder) CancelOrder(signature [65]byte, id order.ID) error {
	tx, err := binder.SendTx(func() (*types.Transaction, error) {
		return binder.cancelOrder(signature, id)
	})
	if err != nil {
		return err
	}

	_, err = binder.conn.PatchedWaitMined(context.Background(), tx)
	return err
}

func (binder *Binder) cancelOrder(signature [65]byte, id order.ID) (*types.Transaction, error) {
	return binder.orderbook.CancelOrder(binder.transactOpts, signature[:], id)
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
	orderMatches := [][32]byte{match}
	return binder.orderbook.ConfirmOrder(binder.transactOpts, [32]byte(id), orderMatches)
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
	matches, err := binder.orderbook.OrderMatch(binder.callOpts, [32]byte(id))
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

// SellOrders in the Orderbook starting at an offset and returning limited
// numbers of sell orders.
func (binder *Binder) SellOrders(offset, limit int) ([]order.ID, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	return binder.sellOrders(offset, limit)
}

func (binder *Binder) sellOrders(offset, limit int) ([]order.ID, error) {
	orders := make([]order.ID, 0, limit)
	for i := 0; i < limit; i++ {
		ordID, ok, err := binder.orderbook.SellOrder(binder.callOpts, big.NewInt(int64(offset+i)))
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

// Broker returns the address of the broker who submitted the order
func (binder *Binder) Broker(id order.ID) (common.Address, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	return binder.broker(id)
}

func (binder *Binder) broker(id order.ID) (common.Address, error) {
	address, err := binder.orderbook.OrderBroker(binder.callOpts, id)
	if err != nil {
		return common.Address{}, err
	}

	return address, nil
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
	counts, err := binder.orderbook.GetOrdersCount(binder.callOpts)
	if err != nil {
		return 0, err
	}

	return counts.Uint64(), nil
}

// OrderID returns the order at a given index in the orderbook
func (binder *Binder) OrderID(index int) ([32]byte, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	return binder.orderID(index)
}

func (binder *Binder) orderID(index int) ([32]byte, error) {
	i := big.NewInt(int64(index))
	id, exist, err := binder.orderbook.GetOrder(binder.callOpts, i)
	if !exist {
		return [32]byte{}, errors.New("order not exist")
	}
	if err != nil {
		return [32]byte{}, err
	}

	return id, nil
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

func (binder *Binder) Deposit(tokenAddress common.Address, value *big.Int) error {
	tx, err := binder.SendTx(func() (*types.Transaction, error) {
		oldValue := binder.transactOpts.Value
		defer func() {
			binder.transactOpts.Value = oldValue
		}()
		if tokenAddress.Hex() == EthereumAddress {
			binder.transactOpts.Value = value
		}

		return binder.renExBalance.Deposit(binder.transactOpts, tokenAddress, value)
	})
	if err != nil {
		return err
	}

	_, err = binder.conn.PatchedWaitMined(context.Background(), tx)
	return err
}

func (binder *Binder) GetBalance(traderAddress common.Address) ([]common.Address, []*big.Int, error) {
	return binder.renExBalance.GetBalances(binder.callOpts, traderAddress)
}

func (binder *Binder) Withdraw(tokenAddress common.Address, value *big.Int) error {
	tx, err := binder.SendTx(func() (*types.Transaction, error) {
		return binder.renExBalance.Withdraw(binder.transactOpts, tokenAddress, value)
	})
	if err != nil {
		return err
	}

	_, err = binder.conn.PatchedWaitMined(context.Background(), tx)
	return err
}

func (binder *Binder) CurrentBlockNumber() (*big.Int, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	block, err := binder.conn.Client.BlockByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return block.Number(), err
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
