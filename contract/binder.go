package contract

import (
	"go/types"
	"context"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math/big"
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
	context      context.Context
	conn         Conn
	transactOpts *bind.TransactOpts
	callOpts     *bind.CallOpts

	republicToken    bindings.RepublicToken
	darknodeRegistry bindings.DarknodeRegistry
	orderbook        bindings.Orderbook

	renExSettlement bindings.RenExSettlement
}

// NewBinder returns a Binder to communicate with contracts
func NewBinder(ctx context.Context, auth *bind.TransactOpts, conn Conn) (Binder, error) {
	transactOpts := auth
	transactOpts.GasPrice = big.NewInt(20000000000)

	nonce, err := conn.Client.PendingNonceAt(context.Background(), transactOpts.From)
	if err != nil {
		return Binder{}, err
	}
	transactOpts.Nonce = nonce	

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

	renExSettlement, err := bindings.NewRenExSettlement(common.HexToAddress(conn.Config.RenExSettlementAddress), bind.ContractBackend(conn.Client))
	if err != nil {
		fmt.Println(fmt.Errorf("cannot bind to RenExSettlement: %v", err))
		return Binder{}, err
	}

	return Binder{
		mu:           new(sync.RWMutex),
		network:      conn.Config.Network,
		context:      ctx,
		conn:         conn,
		transactOpts: transactOpts,
		callOpts:     &bind.CallOpts{},

		republicToken:    *republicToken,
		darknodeRegistry: *darknodeRegistry,
		renExSettlement:  *renExSettlement,
		orderbook:        *orderbook}, nil
}

// SendTx locks binder resources to execute function f (handling nonces explicitly)
// and will wait until the block has been mined on the blockchain. This will allow
// parallel requests to the blockchain since the binder will be unlocked before
// waiting for transaction to complete execution on the blockchain.
func (binder *Binder) SendTx(f func () (*types.Transaction, error)) error {
	tx, err := func() (*types.Transaction, error) {
		binder.mu.Lock()
		defer binder.mu.Unlock()

		return sendTx(f)
	}()
	_, err = binder.conn.PatchedWaitMined(context.Background(), tx)
	return err
}

func (binder *Binder) sendTx(f func () (*types.Transaction, error)) (*types.Transaction, error) {
	tx, err := f()
	if err == core.ErrNonceTooLow || err core.ErrReplaceUnderpriced {
		binder.nonce.Add(binder.nonce, big.NewInt(1))
		return sendTx(f)
	}
	if err == core.ErrNonceTooHigh {
		binder.nonce.Sub(binder.nonce, big.NewInt(1))
		return sendTx(f)
	}
	if err == nil {
		binder.nonce.Add(binder.nonce, big.NewInt(1))
		return tx, nil
	}
	return tx, err
}

// SubmitOrder to the RenEx accounts
func (binder *Binder) SubmitOrder(ord order.Order) error {
	return binder.SendTx(func() (*types.Transaction, error) {
		return binder.submitOrder(ord)
	})
}

func (binder *Binder) submitOrder(ord order.Order) (*types.Transaction, error) {
	nonceHash := big.NewInt(0).SetBytes(ord.BytesFromNonce())
	log.Printf("[submit order] id: %v,tokens:%d, priceCo:%v, priceExp:%v, volumeCo:%v, volumeExp:%v, minVol:%v, minVolExp:%v", base64.StdEncoding.EncodeToString(ord.ID[:]), uint64(ord.Tokens), uint16(ord.Price.Co), uint16(ord.Price.Exp), uint16(ord.Volume.Co), uint16(ord.Volume.Exp), uint16(ord.MinimumVolume.Co), uint16(ord.MinimumVolume.Exp))
	return binder.renExSettlement.SubmitOrder(binder.transactOpts, uint8(ord.Type), uint8(ord.Parity), uint64(ord.Expiry.Unix()), uint64(ord.Tokens), uint16(ord.Price.Co), uint16(ord.Price.Exp), uint16(ord.Volume.Co), uint16(ord.Volume.Exp), uint16(ord.MinimumVolume.Co), uint16(ord.MinimumVolume.Exp), nonceHash)
}

// SubmitMatch will submit a matched order pair to the RenEx accounts
func (binder *Binder) SubmitMatch(buy, sell order.ID) error {
	return binder.SendTx(func() (*types.Transaction, error) {
		return binder.renExSettlement.SubmitMatch(binder.transactOpts, buy, sell)
	})
}

// Settle the order pair which gets confirmed by the Orderbook
func (binder *Binder) Settle(buy order.Order, sell order.Order) error {
	if err := binder.SendTx(func() (*types.Transaction, error) {
		return binder.submitOrder(buy)
	}); err != nil {
		return err
	}
	if err := binder.SendTx(func() (*types.Transaction, error) {
		return binder.submitOrder(sell)
	}); err != nil {
		return err
	}
	return binder.SendTx(func() (*types.Transaction, error) {
		return binder.renExSettlement.SubmitMatch(binder.transactOpts, buy, sell)
	})
}

// SettlementDetail will return settlement details from the smart contract
func (binder *Binder) SettlementDetail(buy, sell order.ID) (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	return binder.settlementDetail(buy, sell)
}

func (binder *Binder) settlementDetail(buy, sell order.ID) (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, error) {
	price, lowVolume, highVolume, lowFee, highFee, err := binder.renExSettlement.GetSettlementDetails(binder.callOpts, buy, sell)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	return price, lowVolume, highVolume, lowFee, highFee, nil
}

// Register a new dark node with the dark node registrar
func (binder *Binder) Register(darknodeID []byte, publicKey []byte, bond *stackint.Int1024) (error) {
	darknodeIDByte, err := toByte(darknodeID)
	if err != nil {
		return err
	}
	return binder.SendTx(func() (*types.Transaction, error) {
		return binder.darknodeRegistry.Register(binder.transactOpts, darknodeIDByte, publicKey, bond.ToBigInt())
	})
}

// Deregister an existing dark node.
func (binder *Binder) Deregister(darknodeID []byte) (error) {
	darknodeIDByte, err := toByte(darknodeID)
	if err != nil {
		return err
	}
	return binder.SendTx(func() (*types.Transaction, error) {
		return binder.darknodeRegistry.Deregister(binder.transactOpts, darknodeIDByte)
	})
}

// Refund withdraws the bond. Must be called before reregistering.
func (binder *Binder) Refund(darknodeID []byte) (*types.Transaction, error) {
	darknodeIDByte, err := toByte(darknodeID)
	if err != nil {
		return err
	}
	return binder.SendTx(func() (*types.Transaction, error) {
		return binder.darknodeRegistry.Refund(binder.transactOpts, darknodeIDByte)
	})
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
	return binder.darknodeRegistry.IsDeregistered(binder.callOpts, darknodeIDByte)
}

// ApproveRen doesn't actually talk to the DNR - instead it approves Ren to it
func (binder *Binder) ApproveRen(value *stackint.Int1024) (error) {
	return binder.SendTx(func() (*types.Transaction, error) {
		return binder.republicToken.Approve(binder.transactOpts, common.HexToAddress(binder.conn.Config.DarknodeRegistryAddress), value.ToBigInt())
	})
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
	return binder.darknodeRegistry.GetOwner(binder.callOpts, darknodeIDByte)
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
	interval, err := binder.darknodeRegistry.MinimumDarkPoolSize(binder.callOpts)
	if err != nil {
		return stackint.Int1024{}, err
	}
	return stackint.FromBigInt(interval)
}

// Pods returns the Pod configuration for the current Epoch.
func (binder *Binder) Pods() ([]registry.Pod, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	return binder.pods()
}

func (binder *Binder) pods() ([]registry.Pod, error) {
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
	epoch, err := binder.darknodeRegistry.CurrentEpoch(binder.callOpts)
	if err != nil {
		return []registry.Pod{}, err
	}
	epochVal := epoch.Epochhash
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

	return binder.epoch()
}

func (binder *Binder) epoch() (registry.Epoch, error) {
	epoch, err := binder.darknodeRegistry.CurrentEpoch(binder.callOpts)
	if err != nil {
		return registry.Epoch{}, err
	}

	epochBlocknumber, err := stackint.FromBigInt(epoch.Blocknumber)
	if err != nil {
		return registry.Epoch{}, err
	}

	var blockhash [32]byte
	for i, b := range epoch.Epochhash.Bytes() {
		blockhash[i] = b
	}

	pods, err := binder.pods()
	if err != nil {
		return registry.Epoch{}, err
	}

	darknodes, err := binder.darknodes()
	if err != nil {
		return registry.Epoch{}, err
	}

	blocknumber, err := epochBlocknumber.ToUint()
	if err != nil {
		return registry.Epoch{}, err
	}

	return registry.Epoch{
		Hash:        blockhash,
		Pods:        pods,
		Darknodes:   darknodes,
		BlockNumber: blocknumber,
	}, nil
}

// NextEpoch will try to turn the Epoch and returns the resulting Epoch. If
// the turning of the Epoch failed, the current Epoch is returned.
func (binder *Binder) NextEpoch() (registry.Epoch, error) {
	if err := binder.SendTx(func() (*types.Transaction, error) {
		return binder.darknodeRegistry.Epoch(binder.transactOpts)
	}); err != nil {
		return nil, err
	}

	return binder.epoch()
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
	pods, err := binder.pods()
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
	if err := binder.SendTx(func() (*types.Transaction, error) {
		return binder.orderbook.OpenBuyOrder(binder.transactOpts, signature[:], id)
	}); err != nil {
		return err
	}
	
	return binder.waitForOrderDepth(id)
}

// OpenSellOrder on the Orderbook. The signature will be used to identify
// the trader that owns the order. The order must be in an undefined state
// to be opened.
func (binder *Binder) OpenSellOrder(signature [65]byte, id order.ID) error {
	if err := binder.SendTx(func() (*types.Transaction, error) {
		return binder.orderbook.OpenSellOrder(binder.transactOpts, signature[:], id)
	}); err != nil {
		return err
	}
	
	return binder.waitForOrderDepth(id)
}

// CancelOrder on the Orderbook. The signature will be used to verify that
// the request was created by the trader that owns the order. The order
// must be in the opened state to be canceled.
func (binder *Binder) CancelOrder(signature [65]byte, id order.ID) error {
	return binder.SendTx(func() (*types.Transaction, error) {
		return binder.orderbook.CancelOrder(binder.transactOpts, signature[:], id)
	})
}

// ConfirmOrder match on the Orderbook.
func (binder *Binder) ConfirmOrder(id order.ID, match order.ID) error {
	if err := binder.SendTx(func() (*types.Transaction, error) {
		return binder.confirmOrder(id, match)
	}); err != nil {
		return err
	}

	for {
		depth, err := binder.orderbook.OrderDepth(binder.callOpts, id)
		if err != nil {
			return err
		}

		if depth.Uint64()-before.Uint64() >= BlocksForConfirmation {
			return nil
		}
	}
}

func (binder *Binder) confirmOrder(id order.ID, match order.ID) (*types.Transaction, error) {
	orderMatches := [][32]byte{match}
	before, err := binder.orderbook.OrderDepth(binder.callOpts, id)
	if err != nil {
		return nil, err
	}
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

// BuyOrders in the Orderbook starting at an offset and returning limited
// numbers of buy orders.
func (binder *Binder) BuyOrders(offset, limit int) ([]order.ID, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	return binder.buyOrders(offset, limit)
}

func (binder *Binder) buyOrders(offset, limit int) ([]order.ID, error) {
	orders := make([]order.ID, 0, limit)
	for i := 0; i < limit; i++ {
		ordID, ok, err := binder.orderbook.BuyOrder(binder.callOpts, big.NewInt(int64(offset+i)))
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
func (binder *Binder) BlockNumber(orderID order.ID) (uint, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	return binder.blockNumber(orderID)
}

func (binder *Binder) blockNumber(orderID order.ID) (uint, error) {
	blockNumber, err := binder.orderbook.OrderBlockNumber(binder.callOpts, orderID)
	if err != nil {
		return 0, err
	}

	return uint(blockNumber.Uint64()), nil
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

func (binder *Binder) waitForOrderDepth(id order.ID) error {
	for {
		depth, err := binder.orderbook.OrderDepth(binder.callOpts, id)
		if err != nil {
			return err
		}

		if depth.Uint64() >= BlocksForConfirmation {
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
