package contracts

import (
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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/republicprotocol/republic-go/contracts/bindings"
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
// changes to an order's status (Open, Canceled or Confirmed) in the Ledger to
// be confirmed. The functions `OpenBuyOrder`, `OpenSellOrder`, `CancelOrder`
// and `ConfirmOrder` return only after the required number of confirmations has
// been reached.
const BlocksForConfirmation = 1

// Epoch consists of the blockhash and number for the epoch
type Epoch struct {
	Blockhash   [32]byte
	BlockNumber stackint.Int1024
}

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
	renExSettlement  bindings.RenExSettlement
	ledger           bindings.Orderbook
	rewardsVault     bindings.RewardVault
}

// GetContractBindings returns a Binder to communicate with contracts
func GetContractBindings(ctx context.Context, keystore crypto.Keystore, conf Config) (Binder, error) {
	conn, err := Connect(conf)
	if err != nil {
		return Binder{}, fmt.Errorf("cannot connect to ethereum: %v", err)
	}
	auth := bind.NewKeyedTransactor(keystore.EcdsaKey.PrivateKey)
	auth.GasPrice = big.NewInt(1000000000)

	darknodeRegistry, err := bindings.NewDarknodeRegistry(common.HexToAddress(conn.Config.DarknodeRegistryAddress), bind.ContractBackend(conn.Client))
	if err != nil {
		fmt.Println(fmt.Errorf("cannot bind to darkpool: %v", err))
		return Binder{}, err
	}

	republicToken, err := bindings.NewRepublicToken(common.HexToAddress(conn.Config.RepublicTokenAddress), bind.ContractBackend(conn.Client))
	if err != nil {
		fmt.Println(fmt.Errorf("cannot bind to darkpool: %v", err))
		return Binder{}, err
	}

	ledger, err := bindings.NewOrderbook(common.HexToAddress(conn.Config.RenLedgerAddress), bind.ContractBackend(conn.Client))
	if err != nil {
		fmt.Println(fmt.Errorf("cannot bind to ren ledger: %v", err))
		return Binder{}, err
	}

	renExSettlement, err := bindings.NewRenExSettlement(common.HexToAddress(conn.Config.RenExAccountsAddress), bind.ContractBackend(conn.Client))
	if err != nil {
		fmt.Println(fmt.Errorf("cannot bind to RenEx accounts: %v", err))
		return Binder{}, err
	}

	return Binder{
		mu:           new(sync.RWMutex),
		network:      conn.Config.Network,
		context:      ctx,
		conn:         conn,
		transactOpts: auth,
		callOpts:     &bind.CallOpts{},

		republicToken:    *republicToken,
		darknodeRegistry: *darknodeRegistry,
		renExSettlement:  *renExSettlement,
		ledger:           *ledger,
		rewardsVault:     bindings.RewardVault{}}, nil
}

// SubmitOrder to the RenEx accounts
func (binder *Binder) SubmitOrder(ord order.Order) error {
	binder.mu.Lock()
	defer binder.mu.Unlock()

	return binder.submitOrder(ord)
}

// SubmitMatch will submit a matched order pair to the RenEx accounts
func (binder *Binder) SubmitMatch(buy, sell order.ID) error {
	binder.mu.Lock()
	defer binder.mu.Unlock()

	return binder.submitMatch(buy, sell)
}

// Settle the order pair which gets confirmed by the RenLedger
func (binder *Binder) Settle(buy order.Order, sell order.Order) error {
	binder.mu.Lock()
	defer binder.mu.Unlock()

	err := binder.submitOrder(buy)
	if err != nil {
		return err
	}
	err = binder.submitOrder(sell)
	if err != nil {
		return err
	}

	return binder.submitMatch(buy.ID, sell.ID)
}

// SettlementDetail will return settlement details from the smart contract
func (binder *Binder) SettlementDetail(buy, sell order.ID) (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	price, lowVolume, highVolume, lowFee, highFee, err := binder.renExSettlement.GetSettlementDetails(binder.callOpts, buy, sell)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	return price, lowVolume, highVolume, lowFee, highFee, nil
}

func (binder *Binder) submitMatch(buy, sell order.ID) error {
	binder.transactOpts.GasLimit = 500000
	tx, err := binder.renExSettlement.SubmitMatch(binder.transactOpts, buy, sell)
	if err != nil {
		return err
	}
	_, err = binder.conn.PatchedWaitMined(binder.context, tx)
	return err
}

func (binder *Binder) submitOrder(ord order.Order) error {
	nonceHash := big.NewInt(0).SetBytes(ord.BytesFromNonce())
	log.Printf("[submit order] id: %v,tokens:%d, priceCo:%v, priceExp:%v, volumeCo:%v, volumeExp:%v, minVol:%v, minVolExp:%v", base64.StdEncoding.EncodeToString(ord.ID[:]), uint64(ord.Tokens), uint16(ord.Price.Co), uint16(ord.Price.Exp), uint16(ord.Volume.Co), uint16(ord.Volume.Exp), uint16(ord.MinimumVolume.Co), uint16(ord.MinimumVolume.Exp))
	tx, err := binder.renExSettlement.SubmitOrder(binder.transactOpts, uint8(ord.Type), uint8(ord.Parity), uint64(ord.Expiry.Unix()), uint64(ord.Tokens), uint16(ord.Price.Co), uint16(ord.Price.Exp), uint16(ord.Volume.Co), uint16(ord.Volume.Exp), uint16(ord.MinimumVolume.Co), uint16(ord.MinimumVolume.Exp), nonceHash)
	if err != nil {
		return err
	}
	_, err = binder.conn.PatchedWaitMined(binder.context, tx)
	return err
}

// Register a new dark node with the dark node registrar
func (binder *Binder) Register(darkNodeID []byte, publicKey []byte, bond *stackint.Int1024) (*types.Transaction, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return &types.Transaction{}, err
	}

	txn, err := binder.darknodeRegistry.Register(binder.transactOpts, darkNodeIDByte, publicKey, bond.ToBigInt())
	if err != nil {
		return nil, err
	}
	_, err = binder.conn.PatchedWaitMined(binder.context, txn)
	return txn, err
}

// Deregister an existing dark node
func (binder *Binder) Deregister(darkNodeID []byte) (*types.Transaction, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return &types.Transaction{}, err
	}
	tx, err := binder.darknodeRegistry.Deregister(binder.transactOpts, darkNodeIDByte)
	if err != nil {
		return nil, err
	}
	_, err = binder.conn.PatchedWaitMined(binder.context, tx)
	return tx, err
}

// Refund withdraws the bond. Must be called before reregistering.
func (binder *Binder) Refund(darkNodeID []byte) (*types.Transaction, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return &types.Transaction{}, err
	}
	tx, err := binder.darknodeRegistry.Refund(binder.transactOpts, darkNodeIDByte)
	if err != nil {
		return nil, err
	}
	_, err = binder.conn.PatchedWaitMined(binder.context, tx)
	return tx, err
}

// GetBond retrieves the bond of an existing dark node
func (binder *Binder) GetBond(darkNodeID []byte) (stackint.Int1024, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return stackint.Int1024{}, err
	}
	bond, err := binder.darknodeRegistry.GetBond(binder.callOpts, darkNodeIDByte)
	if err != nil {
		return stackint.Int1024{}, err
	}
	return stackint.FromBigInt(bond)
}

// IsRegistered returns true if the identity.Address is a current
// registered Darknode. Otherwise, it returns false.
func (binder *Binder) IsRegistered(darknodeAddr identity.Address) (bool, error) {
	darkNodeIDByte, err := toByte(darknodeAddr.ID())
	if err != nil {
		return false, err
	}
	return binder.darknodeRegistry.IsRegistered(binder.callOpts, darkNodeIDByte)
}

// IsDeregistered returns true if the node is deregistered
func (binder *Binder) IsDeregistered(darkNodeID []byte) (bool, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return false, err
	}
	return binder.darknodeRegistry.IsDeregistered(binder.callOpts, darkNodeIDByte)
}

// ApproveRen doesn't actually talk to the DNR - instead it approves Ren to it
func (binder *Binder) ApproveRen(value *stackint.Int1024) (*types.Transaction, error) {
	txn, err := binder.republicToken.Approve(binder.transactOpts, common.HexToAddress(binder.conn.Config.DarknodeRegistryAddress), value.ToBigInt())
	if err != nil {
		return nil, err
	}
	_, err = binder.conn.PatchedWaitMined(binder.context, txn)
	return txn, err
}

// CurrentEpoch returns the current epoch
func (binder *Binder) CurrentEpoch() (Epoch, error) {
	epoch, err := binder.darknodeRegistry.CurrentEpoch(binder.callOpts)
	if err != nil {
		return Epoch{}, err
	}
	blocknumber, err := stackint.FromBigInt(epoch.Blocknumber)
	if err != nil {
		return Epoch{}, err
	}

	var blockhash [32]byte
	for i, b := range epoch.Epochhash.Bytes() {
		blockhash[i] = b
	}

	return Epoch{
		Blockhash:   blockhash,
		BlockNumber: blocknumber,
	}, nil
}

// NextEpoch will try to turn the Epoch and returns the resulting Epoch. If
// the turning of the Epoch failed, the current Epoch is returned.
func (binder *Binder) NextEpoch() (registry.Epoch, error) {
	binder.TriggerEpoch()
	return binder.Epoch()
}

// TriggerEpoch updates the current Epoch if the Minimum Epoch Interval has
// passed since the previous Epoch
func (binder *Binder) TriggerEpoch() (*types.Transaction, error) {
	tx, err := binder.darknodeRegistry.Epoch(binder.transactOpts)
	if err != nil {
		return nil, err
	}
	_, err = binder.conn.PatchedWaitMined(binder.context, tx)
	return tx, err
}

// GetOwner gets the owner of the given dark node
func (binder *Binder) GetOwner(darkNodeID []byte) (common.Address, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return common.Address{}, err
	}
	return binder.darknodeRegistry.GetOwner(binder.callOpts, darkNodeIDByte)
}

// PublicKey returns the RSA public key of the Darknode registered with the
// given identity.Address.
func (binder *Binder) PublicKey(darknodeAddr identity.Address) (rsa.PublicKey, error) {
	darkNodeIDByte, err := toByte(darknodeAddr.ID())
	if err != nil {
		return rsa.PublicKey{}, err
	}
	pubKeyBytes, err := binder.darknodeRegistry.GetPublicKey(binder.callOpts, darkNodeIDByte)
	if err != nil {
		return rsa.PublicKey{}, err
	}
	return crypto.RsaPublicKeyFromBytes(pubKeyBytes)
}

// Darknodes registered in the pod.
func (binder *Binder) Darknodes() (identity.Addresses, error) {
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
	bond, err := binder.darknodeRegistry.MinimumBond(binder.callOpts)
	if err != nil {
		return stackint.Int1024{}, err
	}
	return stackint.FromBigInt(bond)
}

// MinimumEpochInterval returns the minimum number of seconds between
// epochs.
func (binder *Binder) MinimumEpochInterval() (*big.Int, error) {
	return binder.darknodeRegistry.MinimumEpochInterval(binder.callOpts)
}

// MinimumPodSize gets the minimum pod size
func (binder *Binder) MinimumPodSize() (stackint.Int1024, error) {
	interval, err := binder.darknodeRegistry.MinimumDarkPoolSize(binder.callOpts)
	if err != nil {
		return stackint.Int1024{}, err
	}
	return stackint.FromBigInt(interval)
}

// SetGasLimit sets the gas limit to use for transactions
func (binder *Binder) SetGasLimit(limit uint64) {
	binder.transactOpts.GasLimit = limit
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

// Pods returns the Pod configuration for the current Epoch.
func (binder *Binder) Pods() ([]registry.Pod, error) {
	darknodeAddrs, err := binder.Darknodes()

	numberOfNodesInPod, err := binder.MinimumPodSize()
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
	numberOfDarkNodes := big.NewInt(int64(len(darknodeAddrs)))
	x := big.NewInt(0).Mod(epochVal, numberOfDarkNodes)
	positionInOcean := make([]int, len(darknodeAddrs))
	for i := 0; i < len(darknodeAddrs); i++ {
		positionInOcean[i] = -1
	}
	pods := make([]registry.Pod, (len(darknodeAddrs) / int(numberOfNodesInPod.ToBigInt().Int64())))
	for i := 0; i < len(darknodeAddrs); i++ {
		isRegistered, err := binder.IsRegistered(darknodeAddrs[x.Int64()])
		if err != nil {
			return []registry.Pod{}, err
		}
		for !isRegistered || positionInOcean[x.Int64()] != -1 {
			x.Add(x, big.NewInt(1))
			x.Mod(x, numberOfDarkNodes)
			isRegistered, err = binder.IsRegistered(darknodeAddrs[x.Int64()])
			if err != nil {
				return []registry.Pod{}, err
			}
		}
		positionInOcean[x.Int64()] = i
		podID := i % (len(darknodeAddrs) / int(numberOfNodesInPod.ToBigInt().Int64()))
		pods[podID].Darknodes = append(pods[podID].Darknodes, darknodeAddrs[x.Int64()])
		x.Mod(x.Add(x, epochVal), numberOfDarkNodes)
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
	epoch, err := binder.CurrentEpoch()
	if err != nil {
		return registry.Epoch{}, err
	}

	pods, err := binder.Pods()
	if err != nil {
		return registry.Epoch{}, err
	}

	darknodes, err := binder.Darknodes()
	if err != nil {
		return registry.Epoch{}, err
	}

	blocknumber, err := epoch.BlockNumber.ToUint()
	if err != nil {
		return registry.Epoch{}, err
	}

	return registry.Epoch{
		Hash:        epoch.Blockhash,
		Pods:        pods,
		Darknodes:   darknodes,
		BlockNumber: blocknumber,
	}, nil
}

// Pod returns the Pod that contains the given identity.Address in the
// current Epoch. It returns ErrPodNotFound if the identity.Address is not
// registered in the current Epoch.
func (binder *Binder) Pod(addr identity.Address) (registry.Pod, error) {
	pods, err := binder.Pods()
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

// OpenOrders on the Ren Ledger. Returns the number of orders successfully
// opened.
func (binder *Binder) OpenOrders(signatures [][65]byte, orderIDs []order.ID, orderParities []order.Parity) (int, error) {
	binder.mu.Lock()
	defer binder.mu.Unlock()

	if len(signatures) != len(orderIDs) || len(signatures) != len(orderParities) {
		return 0, ErrMismatchedOrderLengths
	}

	nonce, err := binder.conn.Client.PendingNonceAt(context.Background(), binder.transactOpts.From)
	if err != nil {
		return 0, err
	}

	txs := make([]*types.Transaction, 0, len(signatures))
	for i := range signatures {
		binder.transactOpts.GasPrice = big.NewInt(int64(5000000000))
		binder.transactOpts.Nonce.Add(big.NewInt(0).SetUint64(nonce), big.NewInt(int64(i)))

		var tx *types.Transaction
		if orderParities[i] == order.ParityBuy {
			tx, err = binder.ledger.OpenBuyOrder(binder.transactOpts, signatures[i][:], orderIDs[i])
		} else {
			tx, err = binder.ledger.OpenSellOrder(binder.transactOpts, signatures[i][:], orderIDs[i])
		}
		if err != nil {
			break
		}
		txs = append(txs, tx)
	}

	for i := range txs {
		err := binder.waitForOrderDepth(txs[i], orderIDs[i])
		if err != nil {
			return i, err
		}
	}

	return len(txs), err
}

// OpenBuyOrder on the Ren Ledger. The signature will be used to identify
// the trader that owns the order. The order must be in an undefined state
// to be opened.
func (binder *Binder) OpenBuyOrder(signature [65]byte, id order.ID) error {
	binder.mu.Lock()
	defer binder.mu.Unlock()

	binder.transactOpts.GasPrice = big.NewInt(int64(20000000000))
	binder.transactOpts.GasLimit = 3000000

	tx, err := binder.ledger.OpenBuyOrder(binder.transactOpts, signature[:], id)
	if err != nil {
		return err
	}

	return binder.waitForOrderDepth(tx, id)
}

// OpenSellOrder on the Ren Ledger. The signature will be used to identify
// the trader that owns the order. The order must be in an undefined state
// to be opened.
func (binder *Binder) OpenSellOrder(signature [65]byte, id order.ID) error {
	binder.mu.Lock()
	defer binder.mu.Unlock()

	binder.transactOpts.GasPrice = big.NewInt(int64(20000000000))
	binder.transactOpts.GasLimit = 3000000

	tx, err := binder.ledger.OpenSellOrder(binder.transactOpts, signature[:], id)
	if err != nil {
		return err
	}
	binder.transactOpts.GasLimit = 0

	return binder.waitForOrderDepth(tx, id)
}

// CancelOrder on the Ren Ledger. The signature will be used to verify that
// the request was created by the trader that owns the order. The order
// must be in the opened state to be canceled.
func (binder *Binder) CancelOrder(signature [65]byte, id order.ID) error {
	binder.mu.Lock()
	defer binder.mu.Unlock()

	tx, err := binder.ledger.CancelOrder(binder.transactOpts, signature[:], id)
	if err != nil {
		return err
	}
	_, err = binder.conn.PatchedWaitMined(binder.context, tx)
	if err != nil {
		return err
	}
	return nil
}

// ConfirmOrder match on the Ren Ledger.
func (binder *Binder) ConfirmOrder(id order.ID, match order.ID) error {
	binder.mu.Lock()
	defer binder.mu.Unlock()

	orderMatches := [][32]byte{match}
	before, err := binder.ledger.OrderDepth(binder.callOpts, id)
	if err != nil {
		return err
	}
	tx, err := binder.ledger.ConfirmOrder(binder.transactOpts, [32]byte(id), orderMatches)
	if err != nil {
		return err
	}
	_, err = binder.conn.PatchedWaitMined(binder.context, tx)
	if err != nil {
		return err
	}

	for {
		depth, err := binder.ledger.OrderDepth(binder.callOpts, id)
		if err != nil {
			return err
		}

		if depth.Uint64()-before.Uint64() >= BlocksForConfirmation {
			return nil
		}
	}
}

// Priority will return the priority of the order
func (binder *Binder) Priority(id order.ID) (uint64, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	priority, err := binder.ledger.OrderPriority(binder.callOpts, id)
	if err != nil {
		return 0, err
	}

	return priority.Uint64(), nil
}

// Status will return the status of the order
func (binder *Binder) Status(id order.ID) (order.Status, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	var orderID [32]byte
	copy(orderID[:], id[:])
	state, err := binder.ledger.OrderState(binder.callOpts, orderID)
	if err != nil {
		return order.Nil, err
	}

	return order.Status(state), nil
}

// OrderMatch of an order, if any.
func (binder *Binder) OrderMatch(id order.ID) (order.ID, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	matches, err := binder.ledger.OrderMatch(binder.callOpts, [32]byte(id))
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

// BuyOrders in the Ren Ledger starting at an offset and returning limited
// numbers of buy orders.
func (binder *Binder) BuyOrders(offset, limit int) ([]order.ID, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	orders := make([]order.ID, 0, limit)
	for i := 0; i < limit; i++ {
		ordID, ok, err := binder.ledger.BuyOrder(binder.callOpts, big.NewInt(int64(offset+i)))
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

// SellOrders in the Ren Ledger starting at an offset and returning limited
// numbers of sell orders.
func (binder *Binder) SellOrders(offset, limit int) ([]order.ID, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	orders := make([]order.ID, 0, limit)
	for i := 0; i < limit; i++ {
		ordID, ok, err := binder.ledger.SellOrder(binder.callOpts, big.NewInt(int64(offset+i)))
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

	address, err := binder.ledger.OrderTrader(binder.callOpts, id)
	if err != nil {
		return "", err
	}

	return address.String(), nil
}

// Broker returns the address of the broker who submitted the order
func (binder *Binder) Broker(id order.ID) (common.Address, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	address, err := binder.ledger.OrderBroker(binder.callOpts, id)
	if err != nil {
		return common.Address{}, err
	}

	return address, nil
}

// Confirmer returns the address of the confirmer who submitted the order
func (binder *Binder) Confirmer(id order.ID) (common.Address, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	address, err := binder.ledger.OrderConfirmer(binder.callOpts, id)
	if err != nil {
		return common.Address{}, err
	}

	return address, nil
}

// Fee required to open an order.
func (binder *Binder) Fee() (*big.Int, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	return big.NewInt(0), nil
}

// Depth will return depth of confirmation blocks
func (binder *Binder) Depth(orderID order.ID) (uint, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	depth, err := binder.ledger.OrderDepth(binder.callOpts, orderID)
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

	blockNumber, err := binder.ledger.OrderBlockNumber(binder.callOpts, orderID)
	if err != nil {
		return 0, err
	}

	return uint(blockNumber.Uint64()), nil
}

// OrderCounts returns the total number of orders in the ledger
func (binder *Binder) OrderCounts() (uint64, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	counts, err := binder.ledger.GetOrdersCount(binder.callOpts)
	if err != nil {
		return 0, err
	}

	return counts.Uint64(), nil
}

// OrderID returns the order at a given index in the ledger
func (binder *Binder) OrderID(index int) ([32]byte, error) {
	binder.mu.RLock()
	defer binder.mu.RUnlock()

	i := big.NewInt(int64(index))
	id, exist, err := binder.ledger.GetOrder(binder.callOpts, i)
	if !exist {
		return [32]byte{}, errors.New("order not exist")
	}
	if err != nil {
		return [32]byte{}, err
	}

	return id, nil
}

func (binder *Binder) waitForOrderDepth(tx *types.Transaction, id order.ID) error {
	_, err := binder.conn.PatchedWaitMined(binder.context, tx)
	if err != nil {
		return err
	}

	for {
		depth, err := binder.ledger.OrderDepth(binder.callOpts, id)
		if err != nil {
			return err
		}

		if depth.Uint64() >= BlocksForConfirmation {
			return nil
		}
		time.Sleep(time.Second * 14)
	}
}
