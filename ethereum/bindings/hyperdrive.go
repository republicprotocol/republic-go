// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// HyperdriveABI is the input ABI used to generate the binding from.
const HyperdriveABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"tx\",\"type\":\"bytes32[]\"}],\"name\":\"sendTx\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"dnr\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// HyperdriveBin is the compiled bytecode used for deploying new contracts.
const HyperdriveBin = `0x608060405234801561001057600080fd5b5060405160208061028c833981016040525160018054600160a060020a031916600160a060020a0390921691909117905561023c806100506000396000f3006080604052600436106100405763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166387a0c5288114610045575b600080fd5b34801561005157600080fd5b506040805160206004803580820135838102808601850190965280855261009a9536959394602494938501929182918501908490808284375094975061009c9650505050505050565b005b600154604080517f4f5550fc000000000000000000000000000000000000000000000000000000008152336c0100000000000000000000000081026bffffffffffffffffffffffff19166004830152915160009373ffffffffffffffffffffffffffffffffffffffff1691634f5550fc91602480830192602092919082900301818887803b15801561012d57600080fd5b505af1158015610141573d6000803e3d6000fd5b505050506040513d602081101561015757600080fd5b5051151561016457600080fd5b600091505b82518210156101b557600080848481518110151561018357fe5b602090810290910181015182528101919091526040016000205460ff16156101aa57600080fd5b600190910190610169565b600091505b825182101561020b57600160008085858151811015156101d657fe5b6020908102919091018101518252810191909152604001600020805460ff1916911515919091179055600191909101906101ba565b5050505600a165627a7a72305820704b51dbefe1c7ba0df3a7f1f74334b60fbe4e1e3e33f7b030c03584360a0bec0029`

// DeployHyperdrive deploys a new Ethereum contract, binding an instance of Hyperdrive to it.
func DeployHyperdrive(auth *bind.TransactOpts, backend bind.ContractBackend, dnr common.Address) (common.Address, *types.Transaction, *Hyperdrive, error) {
	parsed, err := abi.JSON(strings.NewReader(HyperdriveABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(HyperdriveBin), backend, dnr)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Hyperdrive{HyperdriveCaller: HyperdriveCaller{contract: contract}, HyperdriveTransactor: HyperdriveTransactor{contract: contract}, HyperdriveFilterer: HyperdriveFilterer{contract: contract}}, nil
}

// Hyperdrive is an auto generated Go binding around an Ethereum contract.
type Hyperdrive struct {
	HyperdriveCaller     // Read-only binding to the contract
	HyperdriveTransactor // Write-only binding to the contract
	HyperdriveFilterer   // Log filterer for contract events
}

// HyperdriveCaller is an auto generated read-only Go binding around an Ethereum contract.
type HyperdriveCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HyperdriveTransactor is an auto generated write-only Go binding around an Ethereum contract.
type HyperdriveTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HyperdriveFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type HyperdriveFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HyperdriveSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type HyperdriveSession struct {
	Contract     *Hyperdrive       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// HyperdriveCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type HyperdriveCallerSession struct {
	Contract *HyperdriveCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// HyperdriveTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type HyperdriveTransactorSession struct {
	Contract     *HyperdriveTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// HyperdriveRaw is an auto generated low-level Go binding around an Ethereum contract.
type HyperdriveRaw struct {
	Contract *Hyperdrive // Generic contract binding to access the raw methods on
}

// HyperdriveCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type HyperdriveCallerRaw struct {
	Contract *HyperdriveCaller // Generic read-only contract binding to access the raw methods on
}

// HyperdriveTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type HyperdriveTransactorRaw struct {
	Contract *HyperdriveTransactor // Generic write-only contract binding to access the raw methods on
}

// NewHyperdrive creates a new instance of Hyperdrive, bound to a specific deployed contract.
func NewHyperdrive(address common.Address, backend bind.ContractBackend) (*Hyperdrive, error) {
	contract, err := bindHyperdrive(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Hyperdrive{HyperdriveCaller: HyperdriveCaller{contract: contract}, HyperdriveTransactor: HyperdriveTransactor{contract: contract}, HyperdriveFilterer: HyperdriveFilterer{contract: contract}}, nil
}

// NewHyperdriveCaller creates a new read-only instance of Hyperdrive, bound to a specific deployed contract.
func NewHyperdriveCaller(address common.Address, caller bind.ContractCaller) (*HyperdriveCaller, error) {
	contract, err := bindHyperdrive(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &HyperdriveCaller{contract: contract}, nil
}

// NewHyperdriveTransactor creates a new write-only instance of Hyperdrive, bound to a specific deployed contract.
func NewHyperdriveTransactor(address common.Address, transactor bind.ContractTransactor) (*HyperdriveTransactor, error) {
	contract, err := bindHyperdrive(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &HyperdriveTransactor{contract: contract}, nil
}

// NewHyperdriveFilterer creates a new log filterer instance of Hyperdrive, bound to a specific deployed contract.
func NewHyperdriveFilterer(address common.Address, filterer bind.ContractFilterer) (*HyperdriveFilterer, error) {
	contract, err := bindHyperdrive(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &HyperdriveFilterer{contract: contract}, nil
}

// bindHyperdrive binds a generic wrapper to an already deployed contract.
func bindHyperdrive(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(HyperdriveABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Hyperdrive *HyperdriveRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Hyperdrive.Contract.HyperdriveCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Hyperdrive *HyperdriveRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Hyperdrive.Contract.HyperdriveTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Hyperdrive *HyperdriveRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Hyperdrive.Contract.HyperdriveTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Hyperdrive *HyperdriveCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Hyperdrive.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Hyperdrive *HyperdriveTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Hyperdrive.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Hyperdrive *HyperdriveTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Hyperdrive.Contract.contract.Transact(opts, method, params...)
}

// SendTx is a paid mutator transaction binding the contract method 0x87a0c528.
//
// Solidity: function sendTx(tx bytes32[]) returns()
func (_Hyperdrive *HyperdriveTransactor) SendTx(opts *bind.TransactOpts, tx [][32]byte) (*types.Transaction, error) {
	return _Hyperdrive.contract.Transact(opts, "sendTx", tx)
}

// SendTx is a paid mutator transaction binding the contract method 0x87a0c528.
//
// Solidity: function sendTx(tx bytes32[]) returns()
func (_Hyperdrive *HyperdriveSession) SendTx(tx [][32]byte) (*types.Transaction, error) {
	return _Hyperdrive.Contract.SendTx(&_Hyperdrive.TransactOpts, tx)
}

// SendTx is a paid mutator transaction binding the contract method 0x87a0c528.
//
// Solidity: function sendTx(tx bytes32[]) returns()
func (_Hyperdrive *HyperdriveTransactorSession) SendTx(tx [][32]byte) (*types.Transaction, error) {
	return _Hyperdrive.Contract.SendTx(&_Hyperdrive.TransactOpts, tx)
}
