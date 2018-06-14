package dnr

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/bindings"
	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/stackint"
)

// ErrPodNotFound is returned when dark node address was not found in any pod
var ErrPodNotFound = errors.New("cannot find node in any pod")

// ErrLengthMismatch is returned when ID is not an expected 20 byte value
var ErrLengthMismatch = errors.New("length mismatch")

// Epoch contains a blockhash and a timestamp
type Epoch struct {
	Blockhash   [32]byte
	BlockNumber stackint.Int1024
}

// DarknodeRegistry is the dark node interface
type DarknodeRegistry struct {
	network                 ethereum.Network
	context                 context.Context
	conn                    ethereum.Conn
	transactOpts            *bind.TransactOpts
	callOpts                *bind.CallOpts
	binding                 *bindings.DarknodeRegistry
	tokenBinding            *bindings.RepublicToken
	DarknodeRegistryAddress common.Address
}

// NewDarknodeRegistry returns a Dark node registrar
func NewDarknodeRegistry(context context.Context, conn ethereum.Conn, transactOpts *bind.TransactOpts, callOpts *bind.CallOpts) (DarknodeRegistry, error) {
	contract, err := bindings.NewDarknodeRegistry(common.HexToAddress(conn.Config.DarknodeRegistryAddress), bind.ContractBackend(conn.Client))
	if err != nil {
		return DarknodeRegistry{}, err
	}
	renContract, err := bindings.NewRepublicToken(common.HexToAddress(conn.Config.RepublicTokenAddress), bind.ContractBackend(conn.Client))
	if err != nil {
		return DarknodeRegistry{}, err
	}
	return DarknodeRegistry{
		network:                 conn.Config.Network,
		context:                 context,
		conn:                    conn,
		transactOpts:            transactOpts,
		callOpts:                callOpts,
		binding:                 contract,
		tokenBinding:            renContract,
		DarknodeRegistryAddress: common.HexToAddress(conn.Config.DarknodeRegistryAddress),
	}, nil
}

// Register a new dark node with the dark node registrar
func (darkNodeRegistry *DarknodeRegistry) Register(darkNodeID []byte, publicKey []byte, bond *stackint.Int1024) (*types.Transaction, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return &types.Transaction{}, err
	}

	txn, err := darkNodeRegistry.binding.Register(darkNodeRegistry.transactOpts, darkNodeIDByte, publicKey, bond.ToBigInt())
	if err != nil {
		panic(err)
	}
	_, err = darkNodeRegistry.conn.PatchedWaitMined(darkNodeRegistry.context, txn)
	return txn, err
}

// Deregister an existing dark node
func (darkNodeRegistry *DarknodeRegistry) Deregister(darkNodeID []byte) (*types.Transaction, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return &types.Transaction{}, err
	}
	tx, err := darkNodeRegistry.binding.Deregister(darkNodeRegistry.transactOpts, darkNodeIDByte)
	if err != nil {
		panic(err)
	}
	_, err = darkNodeRegistry.conn.PatchedWaitMined(darkNodeRegistry.context, tx)
	return tx, err
}

// Refund withdraws the bond. Must be called before reregistering.
func (darkNodeRegistry *DarknodeRegistry) Refund(darkNodeID []byte) (*types.Transaction, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return &types.Transaction{}, err
	}
	tx, err := darkNodeRegistry.binding.Refund(darkNodeRegistry.transactOpts, darkNodeIDByte)
	if err != nil {
		panic(err)
	}
	_, err = darkNodeRegistry.conn.PatchedWaitMined(darkNodeRegistry.context, tx)
	return tx, err
}

// GetBond retrieves the bond of an existing dark node
func (darkNodeRegistry *DarknodeRegistry) GetBond(darkNodeID []byte) (stackint.Int1024, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return stackint.Int1024{}, err
	}
	bond, err := darkNodeRegistry.binding.GetBond(darkNodeRegistry.callOpts, darkNodeIDByte)
	if err != nil {
		return stackint.Int1024{}, err
	}
	return stackint.FromBigInt(bond)
}

// IsRegistered implements the cal.Darkpool interface
func (darkNodeRegistry *DarknodeRegistry) IsRegistered(darknodeAddr identity.Address) (bool, error) {
	darkNodeIDByte, err := toByte(darknodeAddr.ID())
	if err != nil {
		return false, err
	}
	return darkNodeRegistry.binding.IsRegistered(darkNodeRegistry.callOpts, darkNodeIDByte)
}

// IsDeregistered returns true if the node is deregistered
func (darkNodeRegistry *DarknodeRegistry) IsDeregistered(darkNodeID []byte) (bool, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return false, err
	}
	return darkNodeRegistry.binding.IsDeregistered(darkNodeRegistry.callOpts, darkNodeIDByte)
}

// ApproveRen doesn't actually talk to the DNR - instead it approves Ren to it
func (darkNodeRegistry *DarknodeRegistry) ApproveRen(value *stackint.Int1024) (*types.Transaction, error) {
	txn, err := darkNodeRegistry.tokenBinding.Approve(darkNodeRegistry.transactOpts, darkNodeRegistry.DarknodeRegistryAddress, value.ToBigInt())
	if err != nil {
		return nil, err
	}
	_, err = darkNodeRegistry.conn.PatchedWaitMined(darkNodeRegistry.context, txn)
	return txn, err
}

// CurrentEpoch returns the current epoch
func (darkNodeRegistry *DarknodeRegistry) CurrentEpoch() (Epoch, error) {
	epoch, err := darkNodeRegistry.binding.CurrentEpoch(darkNodeRegistry.callOpts)
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

// NextEpoch implements the cal.Darkpool interfacce.
func (darkNodeRegistry *DarknodeRegistry) NextEpoch() (cal.Epoch, error) {
	darkNodeRegistry.TriggerEpoch()
	return darkNodeRegistry.Epoch()
}

// TriggerEpoch updates the current Epoch if the Minimum Epoch Interval has
// passed since the previous Epoch
func (darkNodeRegistry *DarknodeRegistry) TriggerEpoch() (*types.Transaction, error) {
	tx, err := darkNodeRegistry.binding.Epoch(darkNodeRegistry.transactOpts)
	if err != nil {
		return nil, err
	}
	_, err = darkNodeRegistry.conn.PatchedWaitMined(darkNodeRegistry.context, tx)
	return tx, err
}

// GetOwner gets the owner of the given dark node
func (darkNodeRegistry *DarknodeRegistry) GetOwner(darkNodeID []byte) (common.Address, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return common.Address{}, err
	}
	return darkNodeRegistry.binding.GetOwner(darkNodeRegistry.callOpts, darkNodeIDByte)
}

// PublicKey implements the cal.Darkpool interface
func (darkNodeRegistry *DarknodeRegistry) PublicKey(darknodeAddr identity.Address) (rsa.PublicKey, error) {
	darkNodeIDByte, err := toByte(darknodeAddr.ID())
	if err != nil {
		return rsa.PublicKey{}, err
	}
	pubKeyBytes, err := darkNodeRegistry.binding.GetPublicKey(darkNodeRegistry.callOpts, darkNodeIDByte)
	if err != nil {
		return rsa.PublicKey{}, err
	}
	return crypto.RsaPublicKeyFromBytes(pubKeyBytes)
}

// Darknodes implements the cal.Darkpool interface
func (darkNodeRegistry *DarknodeRegistry) Darknodes() (identity.Addresses, error) {
	ret, err := darkNodeRegistry.binding.GetDarknodes(darkNodeRegistry.callOpts)
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
func (darkNodeRegistry *DarknodeRegistry) MinimumBond() (stackint.Int1024, error) {
	bond, err := darkNodeRegistry.binding.MinimumBond(darkNodeRegistry.callOpts)
	if err != nil {
		return stackint.Int1024{}, err
	}
	return stackint.FromBigInt(bond)
}

// MinimumEpochInterval implements the cal.Darkpool interface
func (darkNodeRegistry *DarknodeRegistry) MinimumEpochInterval() (*big.Int, error) {
	return darkNodeRegistry.binding.MinimumEpochInterval(darkNodeRegistry.callOpts)
}

// MinimumPodSize gets the minimum pod size
func (darkNodeRegistry *DarknodeRegistry) MinimumPodSize() (stackint.Int1024, error) {
	interval, err := darkNodeRegistry.binding.MinimumDarkPoolSize(darkNodeRegistry.callOpts)
	if err != nil {
		return stackint.Int1024{}, err
	}
	return stackint.FromBigInt(interval)
}

// SetGasLimit sets the gas limit to use for transactions
func (darkNodeRegistry *DarknodeRegistry) SetGasLimit(limit uint64) {
	darkNodeRegistry.transactOpts.GasLimit = limit
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

// Pods implements the cal.Darkpool interface
func (darkNodeRegistry *DarknodeRegistry) Pods() ([]cal.Pod, error) {
	darknodeAddrs, err := darkNodeRegistry.Darknodes()

	numberOfNodesInPod, err := darkNodeRegistry.MinimumPodSize()
	if err != nil {
		return []cal.Pod{}, err
	}
	if len(darknodeAddrs) < int(numberOfNodesInPod.ToBigInt().Int64()) {
		return []cal.Pod{}, fmt.Errorf("degraded pod: expected at least %v addresses, got %v", int(numberOfNodesInPod.ToBigInt().Int64()), len(darknodeAddrs))
	}
	epoch, err := darkNodeRegistry.binding.CurrentEpoch(darkNodeRegistry.callOpts)
	if err != nil {
		return []cal.Pod{}, err
	}
	epochVal := epoch.Epochhash
	numberOfDarkNodes := big.NewInt(int64(len(darknodeAddrs)))
	x := big.NewInt(0).Mod(epochVal, numberOfDarkNodes)
	positionInOcean := make([]int, len(darknodeAddrs))
	for i := 0; i < len(darknodeAddrs); i++ {
		positionInOcean[i] = -1
	}
	pods := make([]cal.Pod, (len(darknodeAddrs) / int(numberOfNodesInPod.ToBigInt().Int64())))
	for i := 0; i < len(darknodeAddrs); i++ {
		isRegistered, err := darkNodeRegistry.IsRegistered(darknodeAddrs[x.Int64()])
		if err != nil {
			return []cal.Pod{}, err
		}
		for !isRegistered || positionInOcean[x.Int64()] != -1 {
			x.Add(x, big.NewInt(1))
			x.Mod(x, numberOfDarkNodes)
			isRegistered, err = darkNodeRegistry.IsRegistered(darknodeAddrs[x.Int64()])
			if err != nil {
				return []cal.Pod{}, err
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
	}
	return pods, nil
}

// Epoch implements the cal.Darkpool interface
func (darkNodeRegistry *DarknodeRegistry) Epoch() (cal.Epoch, error) {
	epoch, err := darkNodeRegistry.CurrentEpoch()
	if err != nil {
		return cal.Epoch{}, err
	}

	pods, err := darkNodeRegistry.Pods()
	if err != nil {
		return cal.Epoch{}, err
	}

	darknodes, err := darkNodeRegistry.Darknodes()
	if err != nil {
		return cal.Epoch{}, err
	}

	blocknumber, err := epoch.BlockNumber.ToUint()
	if err != nil {
		return cal.Epoch{}, err
	}

	return cal.Epoch{
		Hash:        epoch.Blockhash,
		Pods:        pods,
		Darknodes:   darknodes,
		BlockNumber: blocknumber,
	}, nil
}

// Pod implements the cal.Darkpool interface
func (darkNodeRegistry *DarknodeRegistry) Pod(addr identity.Address) (cal.Pod, error) {
	pods, err := darkNodeRegistry.Pods()
	if err != nil {
		return cal.Pod{}, err
	}

	for i := range pods {
		for j := range pods[i].Darknodes {
			if pods[i].Darknodes[j] == addr {
				return pods[i], nil
			}
		}
	}

	return cal.Pod{}, ErrPodNotFound
}
