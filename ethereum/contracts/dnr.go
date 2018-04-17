package contracts

import (
	"context"
	"errors"
	"time"

	"github.com/republicprotocol/republic-go/ethereum/contracts"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/republicprotocol/republic-go/ethereum/bindings"
	"github.com/republicprotocol/republic-go/ethereum/client"
	"github.com/republicprotocol/republic-go/stackint"
)

// Epoch contains a blockhash and a timestamp
type Epoch struct {
	Blockhash [32]byte
	Timestamp stackint.Int1024
}

// DarkNodeRegistry is the dark node interface
type DarkNodeRegistry struct {
	network                 client.Network
	context                 context.Context
	client                  client.Connection
	transactOpts            *bind.TransactOpts
	callOpts                *bind.CallOpts
	binding                 *bindings.DarkNodeRegistry
	tokenBinding            *bindings.RepublicToken
	darkNodeRegistryAddress common.Address
}

// NewDarkNodeRegistry returns a Dark node registrar
func NewDarkNodeRegistry(context context.Context, clientDetails client.Connection, transactOpts *bind.TransactOpts, callOpts *bind.CallOpts) (DarkNodeRegistry, error) {
	contract, err := bindings.NewDarkNodeRegistry(clientDetails.DNRAddress, bind.ContractBackend(clientDetails.Client))
	if err != nil {
		return DarkNodeRegistry{}, err
	}
	renContract, err := bindings.NewRepublicToken(clientDetails.RenAddress, bind.ContractBackend(clientDetails.Client))
	if err != nil {
		return DarkNodeRegistry{}, err
	}
	return DarkNodeRegistry{
		network:                 clientDetails.Network,
		context:                 context,
		client:                  clientDetails,
		transactOpts:            transactOpts,
		callOpts:                callOpts,
		binding:                 contract,
		tokenBinding:            renContract,
		darkNodeRegistryAddress: clientDetails.DNRAddress,
	}, nil
}

// Register a new dark node
func (darkNodeRegistry *DarkNodeRegistry) Register(darkNodeID []byte, publicKey []byte, bond *stackint.Int1024) (*types.Transaction, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return &types.Transaction{}, err
	}

	txn, err := darkNodeRegistry.binding.Register(darkNodeRegistry.transactOpts, darkNodeIDByte, publicKey, bond.ToBigInt())
	if err != nil {
		return nil, err
	}
	_, err = darkNodeRegistry.client.PatchedWaitMined(darkNodeRegistry.context, txn)
	return txn, err
}

// Deregister an existing dark node
func (darkNodeRegistry *DarkNodeRegistry) Deregister(darkNodeID []byte) (*types.Transaction, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return &types.Transaction{}, err
	}
	tx, err := darkNodeRegistry.binding.Deregister(darkNodeRegistry.transactOpts, darkNodeIDByte)
	if err != nil {
		return tx, err
	}
	_, err = darkNodeRegistry.client.PatchedWaitMined(darkNodeRegistry.context, tx)
	return tx, err
}

// Refund withdraws the bond. Must be called before reregistering.
func (darkNodeRegistry *DarkNodeRegistry) Refund(darkNodeID []byte) (*types.Transaction, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return &types.Transaction{}, err
	}
	tx, err := darkNodeRegistry.binding.Refund(darkNodeRegistry.transactOpts, darkNodeIDByte)
	if err != nil {
		return tx, err
	}
	_, err = darkNodeRegistry.client.PatchedWaitMined(darkNodeRegistry.context, tx)
	return tx, err
}

// GetBond retrieves the bond of an existing dark node
func (darkNodeRegistry *DarkNodeRegistry) GetBond(darkNodeID []byte) (stackint.Int1024, error) {
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

// IsRegistered returns true if the node is registered
func (darkNodeRegistry *DarkNodeRegistry) IsRegistered(darkNodeID []byte) (bool, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return false, err
	}
	return darkNodeRegistry.binding.IsRegistered(darkNodeRegistry.callOpts, darkNodeIDByte)
}

// IsDeregistered returns true if the node is deregistered
func (darkNodeRegistry *DarkNodeRegistry) IsDeregistered(darkNodeID []byte) (bool, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return false, err
	}
	return darkNodeRegistry.binding.IsDeregistered(darkNodeRegistry.callOpts, darkNodeIDByte)
}

// ApproveRen doesn't actually talk to the DNR - instead it approved Ren to it
func (darkNodeRegistry *DarkNodeRegistry) ApproveRen(value *stackint.Int1024) error {
	txn, err := darkNodeRegistry.tokenBinding.Approve(darkNodeRegistry.transactOpts, darkNodeRegistry.client.DNRAddress, value.ToBigInt())
	if err != nil {
		return err
	}
	_, err = darkNodeRegistry.client.PatchedWaitMined(darkNodeRegistry.context, txn)
	return err
}

// CurrentEpoch returns the current epoch
func (darkNodeRegistry *DarkNodeRegistry) CurrentEpoch() (Epoch, error) {
	epoch, err := darkNodeRegistry.binding.CurrentEpoch(darkNodeRegistry.callOpts)
	if err != nil {
		return Epoch{}, err
	}
	timestamp, err := stackint.FromBigInt(epoch.Timestamp)
	if err != nil {
		return Epoch{}, err
	}

	var blockhash [32]byte
	for i, b := range epoch.Blockhash.Bytes() {
		blockhash[i] = b
	}

	return Epoch{
		Blockhash: blockhash,
		Timestamp: timestamp,
	}, nil
}

// Epoch updates the current Epoch if the Minimum Epoch Interval has passed since the previous Epoch
func (darkNodeRegistry *DarkNodeRegistry) Epoch() (*types.Transaction, error) {
	tx, err := darkNodeRegistry.binding.Epoch(darkNodeRegistry.transactOpts)
	if err != nil {
		return nil, err
	}
	_, err = darkNodeRegistry.client.PatchedWaitMined(darkNodeRegistry.context, tx)
	return tx, err
}

// TimeUntilEpoch calculates the time remaining until the next Epoch can be called
func (darkNodeRegistry *DarkNodeRegistry) TimeUntilEpoch() (time.Duration, error) {
	epoch, err := darkNodeRegistry.CurrentEpoch()
	if err != nil {
		return 0, err
	}

	minInterval, err := darkNodeRegistry.MinimumEpochInterval()

	nextTime := epoch.Timestamp.Add(&minInterval)
	unix, err := nextTime.ToUint()
	if err != nil {
		// Either minInterval is really big, or unix epoch time has overflowed uint64s.
		return 0, err
	}

	toWait := time.Second * time.Duration(int64(unix)-time.Now().Unix())

	// Ensure toWait is at least 1
	if toWait < 1*time.Second {
		toWait = 1 * time.Second
	}

	// Try again within a minute
	if toWait > time.Minute {
		toWait = time.Minute
	}

	return toWait, nil

}

// WaitForEpoch guarantees that an Epoch as passed (and calls Epoch if connected to Ganache)
func (darkNodeRegistry *DarkNodeRegistry) WaitForEpoch() error {

	previousEpoch, err := darkNodeRegistry.CurrentEpoch()
	if err != nil {
		return err
	}

	currentEpoch := previousEpoch

	for currentEpoch.Blockhash == previousEpoch.Blockhash {

		// Calculate how much time to sleep for
		// If epoch can already be called, returns 1 second
		toWait, err := darkNodeRegistry.TimeUntilEpoch()
		if err != nil {
			return err
		}

		time.Sleep(toWait)

		// If on Ganache, have to call epoch manually
		if darkNodeRegistry.Chain == contracts.ChainGanache {
			darkNodeRegistry.SetGasLimit(300000)
			tx, err := darkNodeRegistry.binding.Epoch(darkNodeRegistry.transactOpts)
			darkNodeRegistry.SetGasLimit(0)
			if err != nil {
				return err
			}
			darkNodeRegistry.client.PatchedWaitMined(darkNodeRegistry.context, tx)
		}

		currentEpoch, err = darkNodeRegistry.CurrentEpoch()
		if err != nil {
			return err
		}
	}

	return nil
}

// GetOwner gets the owner of the given dark node
func (darkNodeRegistry *DarkNodeRegistry) GetOwner(darkNodeID []byte) (common.Address, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return common.Address{}, err
	}
	return darkNodeRegistry.binding.GetOwner(darkNodeRegistry.callOpts, darkNodeIDByte)
}

// GetPublicKey gets the public key of the goven dark node
func (darkNodeRegistry *DarkNodeRegistry) GetPublicKey(darkNodeID []byte) ([]byte, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return []byte{}, err
	}
	return darkNodeRegistry.binding.GetPublicKey(darkNodeRegistry.callOpts, darkNodeIDByte)
}

// GetAllNodes gets all dark nodes
func (darkNodeRegistry *DarkNodeRegistry) GetAllNodes() ([][]byte, error) {
	ret, err := darkNodeRegistry.binding.GetDarkNodes(darkNodeRegistry.callOpts)
	if err != nil {
		return nil, err
	}
	arr := make([][]byte, len(ret))
	for i := range ret {
		arr[i] = ret[i][:]
	}
	return arr, nil
}

// MinimumBond gets the minimum viable bond amount
func (darkNodeRegistry *DarkNodeRegistry) MinimumBond() (stackint.Int1024, error) {
	bond, err := darkNodeRegistry.binding.MinimumBond(darkNodeRegistry.callOpts)
	if err != nil {
		return stackint.Int1024{}, err
	}
	return stackint.FromBigInt(bond)
}

// MinimumEpochInterval gets the minimum epoch interval
func (darkNodeRegistry *DarkNodeRegistry) MinimumEpochInterval() (stackint.Int1024, error) {
	interval, err := darkNodeRegistry.binding.MinimumEpochInterval(darkNodeRegistry.callOpts)
	if err != nil {
		return stackint.Int1024{}, err
	}
	return stackint.FromBigInt(interval)
}

// MinimumDarkPoolSize gets the minumum dark pool size
func (darkNodeRegistry *DarkNodeRegistry) MinimumDarkPoolSize() (stackint.Int1024, error) {
	interval, err := darkNodeRegistry.binding.MinimumDarkPoolSize(darkNodeRegistry.callOpts)
	if err != nil {
		return stackint.Int1024{}, err
	}
	return stackint.FromBigInt(interval)
}

// SetGasLimit sets the gas limit to use for transactions
func (darkNodeRegistry *DarkNodeRegistry) SetGasLimit(limit uint64) {
	darkNodeRegistry.transactOpts.GasLimit = limit
}

// WaitUntilRegistration waits until the registration is successful
func (darkNodeRegistry *DarkNodeRegistry) WaitUntilRegistration(darkNodeID []byte) error {
	isRegistered := false
	for !isRegistered {
		var err error
		isRegistered, err = darkNodeRegistry.IsRegistered(darkNodeID)
		if err != nil {
			return err
		}
		darkNodeRegistry.WaitForEpoch()

	}
	return nil
}

func toByte(id []byte) ([20]byte, error) {
	twentyByte := [20]byte{}
	if len(id) != 20 {
		return twentyByte, errors.New("length mismatch")
	}
	for i := range id {
		twentyByte[i] = id[i]
	}
	return twentyByte, nil
}
