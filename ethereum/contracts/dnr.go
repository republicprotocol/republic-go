package contracts

import (
	"context"
	"errors"
	"time"

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

// DarknodeRegistry is the dark node interface
type DarknodeRegistry struct {
	network                 client.Network
	context                 context.Context
	client                  client.Connection
	transactOpts            *bind.TransactOpts
	callOpts                *bind.CallOpts
	binding                 *bindings.DarknodeRegistry
	tokenBinding            *bindings.RepublicToken
	DarknodeRegistryAddress common.Address
}

// NewDarknodeRegistry returns a Dark node registrar
func NewDarknodeRegistry(context context.Context, clientDetails client.Connection, transactOpts *bind.TransactOpts, callOpts *bind.CallOpts) (DarknodeRegistry, error) {
	contract, err := bindings.NewDarknodeRegistry(clientDetails.DNRAddress, bind.ContractBackend(clientDetails.Client))
	if err != nil {
		return DarknodeRegistry{}, err
	}
	renContract, err := bindings.NewRepublicToken(clientDetails.RenAddress, bind.ContractBackend(clientDetails.Client))
	if err != nil {
		return DarknodeRegistry{}, err
	}
	return DarknodeRegistry{
		network:                 clientDetails.Network,
		context:                 context,
		client:                  clientDetails,
		transactOpts:            transactOpts,
		callOpts:                callOpts,
		binding:                 contract,
		tokenBinding:            renContract,
		DarknodeRegistryAddress: clientDetails.DNRAddress,
	}, nil
}

// Register a new dark node
func (DarknodeRegistry *DarknodeRegistry) Register(darkNodeID []byte, publicKey []byte, bond *stackint.Int1024) (*types.Transaction, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return &types.Transaction{}, err
	}

	txn, err := DarknodeRegistry.binding.Register(DarknodeRegistry.transactOpts, darkNodeIDByte, publicKey, bond.ToBigInt())
	if err != nil {
		return nil, err
	}
	_, err = DarknodeRegistry.client.PatchedWaitMined(DarknodeRegistry.context, txn)
	return txn, err
}

// Deregister an existing dark node
func (DarknodeRegistry *DarknodeRegistry) Deregister(darkNodeID []byte) (*types.Transaction, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return &types.Transaction{}, err
	}
	tx, err := DarknodeRegistry.binding.Deregister(DarknodeRegistry.transactOpts, darkNodeIDByte)
	if err != nil {
		return tx, err
	}
	_, err = DarknodeRegistry.client.PatchedWaitMined(DarknodeRegistry.context, tx)
	return tx, err
}

// Refund withdraws the bond. Must be called before reregistering.
func (DarknodeRegistry *DarknodeRegistry) Refund(darkNodeID []byte) (*types.Transaction, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return &types.Transaction{}, err
	}
	tx, err := DarknodeRegistry.binding.Refund(DarknodeRegistry.transactOpts, darkNodeIDByte)
	if err != nil {
		return tx, err
	}
	_, err = DarknodeRegistry.client.PatchedWaitMined(DarknodeRegistry.context, tx)
	return tx, err
}

// GetBond retrieves the bond of an existing dark node
func (DarknodeRegistry *DarknodeRegistry) GetBond(darkNodeID []byte) (stackint.Int1024, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return stackint.Int1024{}, err
	}
	bond, err := DarknodeRegistry.binding.GetBond(DarknodeRegistry.callOpts, darkNodeIDByte)
	if err != nil {
		return stackint.Int1024{}, err
	}
	return stackint.FromBigInt(bond)
}

// IsRegistered returns true if the node is registered
func (DarknodeRegistry *DarknodeRegistry) IsRegistered(darkNodeID []byte) (bool, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return false, err
	}
	return DarknodeRegistry.binding.IsRegistered(DarknodeRegistry.callOpts, darkNodeIDByte)
}

// IsDeregistered returns true if the node is deregistered
func (DarknodeRegistry *DarknodeRegistry) IsDeregistered(darkNodeID []byte) (bool, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return false, err
	}
	return DarknodeRegistry.binding.IsDeregistered(DarknodeRegistry.callOpts, darkNodeIDByte)
}

// ApproveRen doesn't actually talk to the DNR - instead it approved Ren to it
func (DarknodeRegistry *DarknodeRegistry) ApproveRen(value *stackint.Int1024) error {
	txn, err := DarknodeRegistry.tokenBinding.Approve(DarknodeRegistry.transactOpts, DarknodeRegistry.client.DNRAddress, value.ToBigInt())
	if err != nil {
		return err
	}
	_, err = DarknodeRegistry.client.PatchedWaitMined(DarknodeRegistry.context, txn)
	return err
}

// CurrentEpoch returns the current epoch
func (DarknodeRegistry *DarknodeRegistry) CurrentEpoch() (Epoch, error) {
	epoch, err := DarknodeRegistry.binding.CurrentEpoch(DarknodeRegistry.callOpts)
	if err != nil {
		return Epoch{}, err
	}
	timestamp, err := stackint.FromBigInt(epoch.Timestamp)
	if err != nil {
		return Epoch{}, err
	}

	var blockhash [32]byte
	for i, b := range epoch.Epochhash.Bytes() {
		blockhash[i] = b
	}

	return Epoch{
		Blockhash: blockhash,
		Timestamp: timestamp,
	}, nil
}

// Epoch updates the current Epoch if the Minimum Epoch Interval has passed since the previous Epoch
func (DarknodeRegistry *DarknodeRegistry) Epoch() (*types.Transaction, error) {
	tx, err := DarknodeRegistry.binding.Epoch(DarknodeRegistry.transactOpts)
	if err != nil {
		return nil, err
	}
	_, err = DarknodeRegistry.client.PatchedWaitMined(DarknodeRegistry.context, tx)
	return tx, err
}

// TimeUntilEpoch calculates the time remaining until the next Epoch can be called
func (DarknodeRegistry *DarknodeRegistry) TimeUntilEpoch() (time.Duration, error) {
	epoch, err := DarknodeRegistry.CurrentEpoch()
	if err != nil {
		return 0, err
	}

	minInterval, err := DarknodeRegistry.MinimumEpochInterval()

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
func (DarknodeRegistry *DarknodeRegistry) WaitForEpoch() error {

	previousEpoch, err := DarknodeRegistry.CurrentEpoch()
	if err != nil {
		return err
	}

	currentEpoch := previousEpoch

	for currentEpoch.Blockhash == previousEpoch.Blockhash {

		// Calculate how much time to sleep for
		// If epoch can already be called, returns 1 second
		toWait, err := DarknodeRegistry.TimeUntilEpoch()
		if err != nil {
			return err
		}

		time.Sleep(toWait)

		// If on Ganache, have to call epoch manually
		if DarknodeRegistry.network == client.NetworkGanache {
			DarknodeRegistry.SetGasLimit(300000)
			tx, err := DarknodeRegistry.binding.Epoch(DarknodeRegistry.transactOpts)
			DarknodeRegistry.SetGasLimit(0)
			if err != nil {
				return err
			}
			DarknodeRegistry.client.PatchedWaitMined(DarknodeRegistry.context, tx)
		}

		currentEpoch, err = DarknodeRegistry.CurrentEpoch()
		if err != nil {
			return err
		}
	}

	return nil
}

// GetOwner gets the owner of the given dark node
func (DarknodeRegistry *DarknodeRegistry) GetOwner(darkNodeID []byte) (common.Address, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return common.Address{}, err
	}
	return DarknodeRegistry.binding.GetOwner(DarknodeRegistry.callOpts, darkNodeIDByte)
}

// GetPublicKey gets the public key of the goven dark node
func (DarknodeRegistry *DarknodeRegistry) GetPublicKey(darkNodeID []byte) ([]byte, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return []byte{}, err
	}
	return DarknodeRegistry.binding.GetPublicKey(DarknodeRegistry.callOpts, darkNodeIDByte)
}

// GetAllNodes gets all dark nodes
func (DarknodeRegistry *DarknodeRegistry) GetAllNodes() ([][]byte, error) {
	ret, err := DarknodeRegistry.binding.GetDarkNodes(DarknodeRegistry.callOpts)
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
func (DarknodeRegistry *DarknodeRegistry) MinimumBond() (stackint.Int1024, error) {
	bond, err := DarknodeRegistry.binding.MinimumBond(DarknodeRegistry.callOpts)
	if err != nil {
		return stackint.Int1024{}, err
	}
	return stackint.FromBigInt(bond)
}

// MinimumEpochInterval gets the minimum epoch interval
func (DarknodeRegistry *DarknodeRegistry) MinimumEpochInterval() (stackint.Int1024, error) {
	interval, err := DarknodeRegistry.binding.MinimumEpochInterval(DarknodeRegistry.callOpts)
	if err != nil {
		return stackint.Int1024{}, err
	}
	return stackint.FromBigInt(interval)
}

// MinimumDarkPoolSize gets the minumum dark pool size
func (DarknodeRegistry *DarknodeRegistry) MinimumDarkPoolSize() (stackint.Int1024, error) {
	interval, err := DarknodeRegistry.binding.MinimumDarkPoolSize(DarknodeRegistry.callOpts)
	if err != nil {
		return stackint.Int1024{}, err
	}
	return stackint.FromBigInt(interval)
}

// SetGasLimit sets the gas limit to use for transactions
func (DarknodeRegistry *DarknodeRegistry) SetGasLimit(limit uint64) {
	DarknodeRegistry.transactOpts.GasLimit = limit
}

// WaitUntilRegistration waits until the registration is successful
func (DarknodeRegistry *DarknodeRegistry) WaitUntilRegistration(darkNodeID []byte) error {
	isRegistered := false
	for !isRegistered {
		var err error
		isRegistered, err = DarknodeRegistry.IsRegistered(darkNodeID)
		if err != nil {
			return err
		}
		DarknodeRegistry.WaitForEpoch()

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
