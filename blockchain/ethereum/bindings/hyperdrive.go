// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// HyperdriveABI is the input ABI used to generate the binding from.
const HyperdriveABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"hash\",\"type\":\"bytes32\"},{\"name\":\"orderHashes\",\"type\":\"bytes32[]\"}],\"name\":\"sendOrderMatch\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"getOrderMatch\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"bytes32[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"orderHash\",\"type\":\"bytes32\"}],\"name\":\"depth\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"confirmedOrders\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"dnr\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"orderHash\",\"type\":\"bytes32\"}],\"name\":\"Confirm\",\"type\":\"event\"}]"

// HyperdriveBin is the compiled bytecode used for deploying new contracts.
const HyperdriveBin = `0x608060405234801561001057600080fd5b50604051602080610526833981016040525160028054600160a060020a031916600160a060020a039092169190911790556104d6806100506000396000f3006080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416632e7c8ba581146100665780636a5ba3ce146100c25780639e77a29814610134578063e58d75111461015e575b600080fd5b34801561007257600080fd5b506040805160206004602480358281013584810280870186019097528086526100c0968435963696604495919490910192918291850190849080828437509497506101769650505050505050565b005b3480156100ce57600080fd5b506100da600435610388565b60408051838152602080820183815284519383019390935283519192916060840191858101910280838360005b8381101561011f578181015183820152602001610107565b50505050905001935050505060405180910390f35b34801561014057600080fd5b5061014c6004356103f9565b60408051918252519081900360200190f35b34801561016a57600080fd5b5061014c60043561042e565b600254604080517f4f5550fc000000000000000000000000000000000000000000000000000000008152336c0100000000000000000000000081026bffffffffffffffffffffffff19166004830152915160009373ffffffffffffffffffffffffffffffffffffffff1691634f5550fc91602480830192602092919082900301818887803b15801561020757600080fd5b505af115801561021b573d6000803e3d6000fd5b505050506040513d602081101561023157600080fd5b5051151561023e57600080fd5b600091505b825182101561028d5760016000848481518110151561025e57fe5b60209081029091018101518252810191909152604001600020541561028257600080fd5b600190910190610243565b600091505b8251821015610382576040805180820190915284815260208101849052835160009081908690869081106102c257fe5b60209081029091018101518252818101929092526040016000208251815582820151805191926102fa92600185019290910190610440565b509050504360016000858581518110151561031157fe5b602090810290910181015182528101919091526040016000205582517f94d94affe11eeecd469f8d64826f21e58723e70a2d6f0a69d3475ec3f6ff6f689084908490811061035b57fe5b602090810290910181015160408051918252519081900390910190a1600190910190610292565b50505050565b600081815260208181526040808320805460019091018054835181860281018601909452808452606094929391928391908301828280156103e957602002820191906000526020600020905b815481526001909101906020018083116103d4575b5050505050905091509150915091565b600081815260016020526040812054151561041657506000610429565b5060008181526001602052604090205443035b919050565b60016020526000908152604090205481565b82805482825590600052602060002090810192821561047d579160200282015b8281111561047d5782518255602090920191600190910190610460565b5061048992915061048d565b5090565b6104a791905b808211156104895760008155600101610493565b905600a165627a7a72305820df404005bdd8118ecf8e97d45154c110da72abc0beb631286bf4f44f4ce450ac0029`

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

// ConfirmedOrders is a free data retrieval call binding the contract method 0xe58d7511.
//
// Solidity: function confirmedOrders( bytes32) constant returns(uint256)
func (_Hyperdrive *HyperdriveCaller) ConfirmedOrders(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Hyperdrive.contract.Call(opts, out, "confirmedOrders", arg0)
	return *ret0, err
}

// ConfirmedOrders is a free data retrieval call binding the contract method 0xe58d7511.
//
// Solidity: function confirmedOrders( bytes32) constant returns(uint256)
func (_Hyperdrive *HyperdriveSession) ConfirmedOrders(arg0 [32]byte) (*big.Int, error) {
	return _Hyperdrive.Contract.ConfirmedOrders(&_Hyperdrive.CallOpts, arg0)
}

// ConfirmedOrders is a free data retrieval call binding the contract method 0xe58d7511.
//
// Solidity: function confirmedOrders( bytes32) constant returns(uint256)
func (_Hyperdrive *HyperdriveCallerSession) ConfirmedOrders(arg0 [32]byte) (*big.Int, error) {
	return _Hyperdrive.Contract.ConfirmedOrders(&_Hyperdrive.CallOpts, arg0)
}

// Depth is a free data retrieval call binding the contract method 0x9e77a298.
//
// Solidity: function depth(orderHash bytes32) constant returns(uint256)
func (_Hyperdrive *HyperdriveCaller) Depth(opts *bind.CallOpts, orderHash [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Hyperdrive.contract.Call(opts, out, "depth", orderHash)
	return *ret0, err
}

// Depth is a free data retrieval call binding the contract method 0x9e77a298.
//
// Solidity: function depth(orderHash bytes32) constant returns(uint256)
func (_Hyperdrive *HyperdriveSession) Depth(orderHash [32]byte) (*big.Int, error) {
	return _Hyperdrive.Contract.Depth(&_Hyperdrive.CallOpts, orderHash)
}

// Depth is a free data retrieval call binding the contract method 0x9e77a298.
//
// Solidity: function depth(orderHash bytes32) constant returns(uint256)
func (_Hyperdrive *HyperdriveCallerSession) Depth(orderHash [32]byte) (*big.Int, error) {
	return _Hyperdrive.Contract.Depth(&_Hyperdrive.CallOpts, orderHash)
}

// GetOrderMatch is a free data retrieval call binding the contract method 0x6a5ba3ce.
//
// Solidity: function getOrderMatch(hash bytes32) constant returns(bytes32, bytes32[])
func (_Hyperdrive *HyperdriveCaller) GetOrderMatch(opts *bind.CallOpts, hash [32]byte) ([32]byte, [][32]byte, error) {
	var (
		ret0 = new([32]byte)
		ret1 = new([][32]byte)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := _Hyperdrive.contract.Call(opts, out, "getOrderMatch", hash)
	return *ret0, *ret1, err
}

// GetOrderMatch is a free data retrieval call binding the contract method 0x6a5ba3ce.
//
// Solidity: function getOrderMatch(hash bytes32) constant returns(bytes32, bytes32[])
func (_Hyperdrive *HyperdriveSession) GetOrderMatch(hash [32]byte) ([32]byte, [][32]byte, error) {
	return _Hyperdrive.Contract.GetOrderMatch(&_Hyperdrive.CallOpts, hash)
}

// GetOrderMatch is a free data retrieval call binding the contract method 0x6a5ba3ce.
//
// Solidity: function getOrderMatch(hash bytes32) constant returns(bytes32, bytes32[])
func (_Hyperdrive *HyperdriveCallerSession) GetOrderMatch(hash [32]byte) ([32]byte, [][32]byte, error) {
	return _Hyperdrive.Contract.GetOrderMatch(&_Hyperdrive.CallOpts, hash)
}

// SendOrderMatch is a paid mutator transaction binding the contract method 0x2e7c8ba5.
//
// Solidity: function sendOrderMatch(hash bytes32, orderHashes bytes32[]) returns()
func (_Hyperdrive *HyperdriveTransactor) SendOrderMatch(opts *bind.TransactOpts, hash [32]byte, orderHashes [][32]byte) (*types.Transaction, error) {
	return _Hyperdrive.contract.Transact(opts, "sendOrderMatch", hash, orderHashes)
}

// SendOrderMatch is a paid mutator transaction binding the contract method 0x2e7c8ba5.
//
// Solidity: function sendOrderMatch(hash bytes32, orderHashes bytes32[]) returns()
func (_Hyperdrive *HyperdriveSession) SendOrderMatch(hash [32]byte, orderHashes [][32]byte) (*types.Transaction, error) {
	return _Hyperdrive.Contract.SendOrderMatch(&_Hyperdrive.TransactOpts, hash, orderHashes)
}

// SendOrderMatch is a paid mutator transaction binding the contract method 0x2e7c8ba5.
//
// Solidity: function sendOrderMatch(hash bytes32, orderHashes bytes32[]) returns()
func (_Hyperdrive *HyperdriveTransactorSession) SendOrderMatch(hash [32]byte, orderHashes [][32]byte) (*types.Transaction, error) {
	return _Hyperdrive.Contract.SendOrderMatch(&_Hyperdrive.TransactOpts, hash, orderHashes)
}

// HyperdriveConfirmIterator is returned from FilterConfirm and is used to iterate over the raw logs and unpacked data for Confirm events raised by the Hyperdrive contract.
type HyperdriveConfirmIterator struct {
	Event *HyperdriveConfirm // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *HyperdriveConfirmIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HyperdriveConfirm)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(HyperdriveConfirm)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *HyperdriveConfirmIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HyperdriveConfirmIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HyperdriveConfirm represents a Confirm event raised by the Hyperdrive contract.
type HyperdriveConfirm struct {
	OrderHash [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterConfirm is a free log retrieval operation binding the contract event 0x94d94affe11eeecd469f8d64826f21e58723e70a2d6f0a69d3475ec3f6ff6f68.
//
// Solidity: event Confirm(orderHash bytes32)
func (_Hyperdrive *HyperdriveFilterer) FilterConfirm(opts *bind.FilterOpts) (*HyperdriveConfirmIterator, error) {

	logs, sub, err := _Hyperdrive.contract.FilterLogs(opts, "Confirm")
	if err != nil {
		return nil, err
	}
	return &HyperdriveConfirmIterator{contract: _Hyperdrive.contract, event: "Confirm", logs: logs, sub: sub}, nil
}

// WatchConfirm is a free log subscription operation binding the contract event 0x94d94affe11eeecd469f8d64826f21e58723e70a2d6f0a69d3475ec3f6ff6f68.
//
// Solidity: event Confirm(orderHash bytes32)
func (_Hyperdrive *HyperdriveFilterer) WatchConfirm(opts *bind.WatchOpts, sink chan<- *HyperdriveConfirm) (event.Subscription, error) {

	logs, sub, err := _Hyperdrive.contract.WatchLogs(opts, "Confirm")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HyperdriveConfirm)
				if err := _Hyperdrive.contract.UnpackLog(event, "Confirm", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}
