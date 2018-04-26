package dnr

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/bindings"
	"github.com/republicprotocol/republic-go/stackint"
)

// Epoch contains a blockhash and a timestamp
type Epoch struct {
	Blockhash [32]byte
	Timestamp stackint.Int1024
}

// DarknodeRegistry is the dark node interface
type DarknodeRegistry struct {
<<<<<<< HEAD:blockchain/ethereum/dnr/dnr.go
	network                 ethereum.Network
=======
	network                 client.Network
>>>>>>> hyperdrive:ethereum/contracts/dnr.go
	context                 context.Context
	conn                    ethereum.Conn
	transactOpts            *bind.TransactOpts
	callOpts                *bind.CallOpts
	binding                 *bindings.DarknodeRegistry
	tokenBinding            *bindings.RepublicToken
	DarknodeRegistryAddress common.Address
}

// NewDarknodeRegistry returns a Dark node registrar
<<<<<<< HEAD:blockchain/ethereum/dnr/dnr.go
func NewDarknodeRegistry(context context.Context, conn ethereum.Conn, transactOpts *bind.TransactOpts, callOpts *bind.CallOpts) (DarknodeRegistry, error) {
	contract, err := bindings.NewDarkNodeRegistry(conn.DarknodeRegistryAddress, bind.ContractBackend(conn.Client))
=======
func NewDarknodeRegistry(context context.Context, clientDetails client.Connection, transactOpts *bind.TransactOpts, callOpts *bind.CallOpts) (DarknodeRegistry, error) {
	contract, err := bindings.NewDarknodeRegistry(clientDetails.DNRAddress, bind.ContractBackend(clientDetails.Client))
>>>>>>> hyperdrive:ethereum/contracts/dnr.go
	if err != nil {
		return DarknodeRegistry{}, err
	}
	renContract, err := bindings.NewRepublicToken(conn.RepublicTokenAddress, bind.ContractBackend(conn.Client))
	if err != nil {
		return DarknodeRegistry{}, err
	}
	return DarknodeRegistry{
<<<<<<< HEAD:blockchain/ethereum/dnr/dnr.go
		network:                 conn.Network,
=======
		network:                 clientDetails.Network,
>>>>>>> hyperdrive:ethereum/contracts/dnr.go
		context:                 context,
		conn:                    conn,
		transactOpts:            transactOpts,
		callOpts:                callOpts,
		binding:                 contract,
		tokenBinding:            renContract,
<<<<<<< HEAD:blockchain/ethereum/dnr/dnr.go
		darkNodeRegistryAddress: conn.DarknodeRegistryAddress,
=======
		DarknodeRegistryAddress: clientDetails.DNRAddress,
>>>>>>> hyperdrive:ethereum/contracts/dnr.go
	}, nil
}

// Register a new dark node
<<<<<<< HEAD:blockchain/ethereum/dnr/dnr.go
func (darkNodeRegistry *DarknodeRegistry) Register(darkNodeID []byte, publicKey []byte, bond *stackint.Int1024) (*types.Transaction, error) {
=======
func (DarknodeRegistry *DarknodeRegistry) Register(darkNodeID []byte, publicKey []byte, bond *stackint.Int1024) (*types.Transaction, error) {
>>>>>>> hyperdrive:ethereum/contracts/dnr.go
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return &types.Transaction{}, err
	}

	txn, err := DarknodeRegistry.binding.Register(DarknodeRegistry.transactOpts, darkNodeIDByte, publicKey, bond.ToBigInt())
	if err != nil {
		fmt.Println(darkNodeRegistry.transactOpts.GasLimit)
		panic(err)
	}
<<<<<<< HEAD:blockchain/ethereum/dnr/dnr.go
	_, err = darkNodeRegistry.conn.PatchedWaitMined(darkNodeRegistry.context, txn)
=======
	_, err = DarknodeRegistry.client.PatchedWaitMined(DarknodeRegistry.context, txn)
>>>>>>> hyperdrive:ethereum/contracts/dnr.go
	return txn, err
}

// Deregister an existing dark node
<<<<<<< HEAD:blockchain/ethereum/dnr/dnr.go
func (darkNodeRegistry *DarknodeRegistry) Deregister(darkNodeID []byte) (*types.Transaction, error) {
=======
func (DarknodeRegistry *DarknodeRegistry) Deregister(darkNodeID []byte) (*types.Transaction, error) {
>>>>>>> hyperdrive:ethereum/contracts/dnr.go
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return &types.Transaction{}, err
	}
	tx, err := DarknodeRegistry.binding.Deregister(DarknodeRegistry.transactOpts, darkNodeIDByte)
	if err != nil {
		fmt.Println(darkNodeRegistry.transactOpts.GasLimit)
		panic(err)
	}
<<<<<<< HEAD:blockchain/ethereum/dnr/dnr.go
	_, err = darkNodeRegistry.conn.PatchedWaitMined(darkNodeRegistry.context, tx)
=======
	_, err = DarknodeRegistry.client.PatchedWaitMined(DarknodeRegistry.context, tx)
>>>>>>> hyperdrive:ethereum/contracts/dnr.go
	return tx, err
}

// Refund withdraws the bond. Must be called before reregistering.
<<<<<<< HEAD:blockchain/ethereum/dnr/dnr.go
func (darkNodeRegistry *DarknodeRegistry) Refund(darkNodeID []byte) (*types.Transaction, error) {
=======
func (DarknodeRegistry *DarknodeRegistry) Refund(darkNodeID []byte) (*types.Transaction, error) {
>>>>>>> hyperdrive:ethereum/contracts/dnr.go
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return &types.Transaction{}, err
	}
	tx, err := DarknodeRegistry.binding.Refund(DarknodeRegistry.transactOpts, darkNodeIDByte)
	if err != nil {
		fmt.Println(darkNodeRegistry.transactOpts.GasLimit)
		panic(err)
	}
<<<<<<< HEAD:blockchain/ethereum/dnr/dnr.go
	_, err = darkNodeRegistry.conn.PatchedWaitMined(darkNodeRegistry.context, tx)
=======
	_, err = DarknodeRegistry.client.PatchedWaitMined(DarknodeRegistry.context, tx)
>>>>>>> hyperdrive:ethereum/contracts/dnr.go
	return tx, err
}

// GetBond retrieves the bond of an existing dark node
<<<<<<< HEAD:blockchain/ethereum/dnr/dnr.go
func (darkNodeRegistry *DarknodeRegistry) GetBond(darkNodeID []byte) (stackint.Int1024, error) {
=======
func (DarknodeRegistry *DarknodeRegistry) GetBond(darkNodeID []byte) (stackint.Int1024, error) {
>>>>>>> hyperdrive:ethereum/contracts/dnr.go
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
<<<<<<< HEAD:blockchain/ethereum/dnr/dnr.go
func (darkNodeRegistry *DarknodeRegistry) IsRegistered(darkNodeID []byte) (bool, error) {
=======
func (DarknodeRegistry *DarknodeRegistry) IsRegistered(darkNodeID []byte) (bool, error) {
>>>>>>> hyperdrive:ethereum/contracts/dnr.go
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return false, err
	}
	return DarknodeRegistry.binding.IsRegistered(DarknodeRegistry.callOpts, darkNodeIDByte)
}

// IsDeregistered returns true if the node is deregistered
<<<<<<< HEAD:blockchain/ethereum/dnr/dnr.go
func (darkNodeRegistry *DarknodeRegistry) IsDeregistered(darkNodeID []byte) (bool, error) {
=======
func (DarknodeRegistry *DarknodeRegistry) IsDeregistered(darkNodeID []byte) (bool, error) {
>>>>>>> hyperdrive:ethereum/contracts/dnr.go
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return false, err
	}
	return DarknodeRegistry.binding.IsDeregistered(DarknodeRegistry.callOpts, darkNodeIDByte)
}

// ApproveRen doesn't actually talk to the DNR - instead it approved Ren to it
<<<<<<< HEAD:blockchain/ethereum/dnr/dnr.go
func (darkNodeRegistry *DarknodeRegistry) ApproveRen(value *stackint.Int1024) (*types.Transaction, error) {
	txn, err := darkNodeRegistry.tokenBinding.Approve(darkNodeRegistry.transactOpts, darkNodeRegistry.conn.DarknodeRegistryAddress, value.ToBigInt())
=======
func (DarknodeRegistry *DarknodeRegistry) ApproveRen(value *stackint.Int1024) error {
	txn, err := DarknodeRegistry.tokenBinding.Approve(DarknodeRegistry.transactOpts, DarknodeRegistry.client.DNRAddress, value.ToBigInt())
>>>>>>> hyperdrive:ethereum/contracts/dnr.go
	if err != nil {
		return nil, err
	}
<<<<<<< HEAD:blockchain/ethereum/dnr/dnr.go
	_, err = darkNodeRegistry.conn.PatchedWaitMined(darkNodeRegistry.context, txn)
	return txn, err
}

// CurrentEpoch returns the current epoch
func (darkNodeRegistry *DarknodeRegistry) CurrentEpoch() (Epoch, error) {
	epoch, err := darkNodeRegistry.binding.CurrentEpoch(darkNodeRegistry.callOpts)
=======
	_, err = DarknodeRegistry.client.PatchedWaitMined(DarknodeRegistry.context, txn)
	return err
}

// CurrentEpoch returns the current epoch
func (DarknodeRegistry *DarknodeRegistry) CurrentEpoch() (Epoch, error) {
	epoch, err := DarknodeRegistry.binding.CurrentEpoch(DarknodeRegistry.callOpts)
>>>>>>> hyperdrive:ethereum/contracts/dnr.go
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
<<<<<<< HEAD:blockchain/ethereum/dnr/dnr.go
func (darkNodeRegistry *DarknodeRegistry) Epoch() (*types.Transaction, error) {
	tx, err := darkNodeRegistry.binding.Epoch(darkNodeRegistry.transactOpts)
	if err != nil {
		return nil, err
	}
	_, err = darkNodeRegistry.conn.PatchedWaitMined(darkNodeRegistry.context, tx)
=======
func (DarknodeRegistry *DarknodeRegistry) Epoch() (*types.Transaction, error) {
	tx, err := DarknodeRegistry.binding.Epoch(DarknodeRegistry.transactOpts)
	if err != nil {
		return nil, err
	}
	_, err = DarknodeRegistry.client.PatchedWaitMined(DarknodeRegistry.context, tx)
>>>>>>> hyperdrive:ethereum/contracts/dnr.go
	return tx, err
}

// TimeUntilEpoch calculates the time remaining until the next Epoch can be called
<<<<<<< HEAD:blockchain/ethereum/dnr/dnr.go
func (darkNodeRegistry *DarknodeRegistry) TimeUntilEpoch() (time.Duration, error) {
	epoch, err := darkNodeRegistry.CurrentEpoch()
=======
func (DarknodeRegistry *DarknodeRegistry) TimeUntilEpoch() (time.Duration, error) {
	epoch, err := DarknodeRegistry.CurrentEpoch()
>>>>>>> hyperdrive:ethereum/contracts/dnr.go
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
<<<<<<< HEAD:blockchain/ethereum/dnr/dnr.go
func (darkNodeRegistry *DarknodeRegistry) WaitForEpoch() error {
=======
func (DarknodeRegistry *DarknodeRegistry) WaitForEpoch() error {
>>>>>>> hyperdrive:ethereum/contracts/dnr.go

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
<<<<<<< HEAD:blockchain/ethereum/dnr/dnr.go
		if darkNodeRegistry.network == ethereum.NetworkGanache {
			tx, err := darkNodeRegistry.binding.Epoch(darkNodeRegistry.transactOpts)
			if err != nil {
				return err
			}
			darkNodeRegistry.conn.PatchedWaitMined(darkNodeRegistry.context, tx)
=======
		if DarknodeRegistry.network == client.NetworkGanache {
			DarknodeRegistry.SetGasLimit(300000)
			tx, err := DarknodeRegistry.binding.Epoch(DarknodeRegistry.transactOpts)
			DarknodeRegistry.SetGasLimit(0)
			if err != nil {
				return err
			}
			DarknodeRegistry.client.PatchedWaitMined(DarknodeRegistry.context, tx)
>>>>>>> hyperdrive:ethereum/contracts/dnr.go
		}

		currentEpoch, err = DarknodeRegistry.CurrentEpoch()
		if err != nil {
			return err
		}
	}

	return nil
}

// GetOwner gets the owner of the given dark node
<<<<<<< HEAD:blockchain/ethereum/dnr/dnr.go
func (darkNodeRegistry *DarknodeRegistry) GetOwner(darkNodeID []byte) (common.Address, error) {
=======
func (DarknodeRegistry *DarknodeRegistry) GetOwner(darkNodeID []byte) (common.Address, error) {
>>>>>>> hyperdrive:ethereum/contracts/dnr.go
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return common.Address{}, err
	}
	return DarknodeRegistry.binding.GetOwner(DarknodeRegistry.callOpts, darkNodeIDByte)
}

// GetPublicKey gets the public key of the goven dark node
<<<<<<< HEAD:blockchain/ethereum/dnr/dnr.go
func (darkNodeRegistry *DarknodeRegistry) GetPublicKey(darkNodeID []byte) ([]byte, error) {
=======
func (DarknodeRegistry *DarknodeRegistry) GetPublicKey(darkNodeID []byte) ([]byte, error) {
>>>>>>> hyperdrive:ethereum/contracts/dnr.go
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return []byte{}, err
	}
	return DarknodeRegistry.binding.GetPublicKey(DarknodeRegistry.callOpts, darkNodeIDByte)
}

// GetAllNodes gets all dark nodes
<<<<<<< HEAD:blockchain/ethereum/dnr/dnr.go
func (darkNodeRegistry *DarknodeRegistry) GetAllNodes() ([][]byte, error) {
	ret, err := darkNodeRegistry.binding.GetDarkNodes(darkNodeRegistry.callOpts)
=======
func (DarknodeRegistry *DarknodeRegistry) GetAllNodes() ([][]byte, error) {
	ret, err := DarknodeRegistry.binding.GetDarkNodes(DarknodeRegistry.callOpts)
>>>>>>> hyperdrive:ethereum/contracts/dnr.go
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
<<<<<<< HEAD:blockchain/ethereum/dnr/dnr.go
func (darkNodeRegistry *DarknodeRegistry) MinimumBond() (stackint.Int1024, error) {
	bond, err := darkNodeRegistry.binding.MinimumBond(darkNodeRegistry.callOpts)
=======
func (DarknodeRegistry *DarknodeRegistry) MinimumBond() (stackint.Int1024, error) {
	bond, err := DarknodeRegistry.binding.MinimumBond(DarknodeRegistry.callOpts)
>>>>>>> hyperdrive:ethereum/contracts/dnr.go
	if err != nil {
		return stackint.Int1024{}, err
	}
	return stackint.FromBigInt(bond)
}

// MinimumEpochInterval gets the minimum epoch interval
<<<<<<< HEAD:blockchain/ethereum/dnr/dnr.go
func (darkNodeRegistry *DarknodeRegistry) MinimumEpochInterval() (stackint.Int1024, error) {
	interval, err := darkNodeRegistry.binding.MinimumEpochInterval(darkNodeRegistry.callOpts)
=======
func (DarknodeRegistry *DarknodeRegistry) MinimumEpochInterval() (stackint.Int1024, error) {
	interval, err := DarknodeRegistry.binding.MinimumEpochInterval(DarknodeRegistry.callOpts)
>>>>>>> hyperdrive:ethereum/contracts/dnr.go
	if err != nil {
		return stackint.Int1024{}, err
	}
	return stackint.FromBigInt(interval)
}

// MinimumDarkPoolSize gets the minumum dark pool size
<<<<<<< HEAD:blockchain/ethereum/dnr/dnr.go
func (darkNodeRegistry *DarknodeRegistry) MinimumDarkPoolSize() (stackint.Int1024, error) {
	interval, err := darkNodeRegistry.binding.MinimumDarkPoolSize(darkNodeRegistry.callOpts)
=======
func (DarknodeRegistry *DarknodeRegistry) MinimumDarkPoolSize() (stackint.Int1024, error) {
	interval, err := DarknodeRegistry.binding.MinimumDarkPoolSize(DarknodeRegistry.callOpts)
>>>>>>> hyperdrive:ethereum/contracts/dnr.go
	if err != nil {
		return stackint.Int1024{}, err
	}
	return stackint.FromBigInt(interval)
}

// SetGasLimit sets the gas limit to use for transactions
<<<<<<< HEAD:blockchain/ethereum/dnr/dnr.go
func (darkNodeRegistry *DarknodeRegistry) SetGasLimit(limit uint64) {
	darkNodeRegistry.transactOpts.GasLimit = limit
}

// WaitUntilRegistration waits until the registration is successful
func (darkNodeRegistry *DarknodeRegistry) WaitUntilRegistration(darkNodeID []byte) error {
=======
func (DarknodeRegistry *DarknodeRegistry) SetGasLimit(limit uint64) {
	DarknodeRegistry.transactOpts.GasLimit = limit
}

// WaitUntilRegistration waits until the registration is successful
func (DarknodeRegistry *DarknodeRegistry) WaitUntilRegistration(darkNodeID []byte) error {
>>>>>>> hyperdrive:ethereum/contracts/dnr.go
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
