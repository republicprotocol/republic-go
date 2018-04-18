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

// HyperdriveEpochABI is the input ABI used to generate the binding from.
const HyperdriveEpochABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"minimumEpochInterval\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentEpoch\",\"outputs\":[{\"name\":\"blockhash\",\"type\":\"uint256\"},{\"name\":\"timestamp\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"epoch\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_minimumEpochInterval\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"NewEpoch\",\"type\":\"event\"}]"

// HyperdriveEpochBin is the compiled bytecode used for deploying new contracts.
const HyperdriveEpochBin = `0x608060405234801561001057600080fd5b506040516020806101c2833981016040818152915160025581810190915260001943014080825242602090920182905260005560015561016d806100556000396000f3006080604052600436106100565763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166355cacda5811461005b5780637667180814610082578063900cf0cf146100b0575b600080fd5b34801561006757600080fd5b506100706100c7565b60408051918252519081900360200190f35b34801561008e57600080fd5b506100976100cd565b6040805192835260208301919091528051918290030190f35b3480156100bc57600080fd5b506100c56100d6565b005b60025481565b60005460015482565b6002546001546000910142116100eb57600080fd5b50604080518082018252600019430140808252600254600180549091016020909301839052600082815592905591517fe358419ca0dd7928a310d787a606dfae5d869f5071249efa6107105e7afc40bc9190a1505600a165627a7a7230582098390aed1f40bd6cc9257ac65082b7bc3bc766b9e113589a146240e642e277da0029`

// DeployHyperdriveEpoch deploys a new Ethereum contract, binding an instance of HyperdriveEpoch to it.
func DeployHyperdriveEpoch(auth *bind.TransactOpts, backend bind.ContractBackend, _minimumEpochInterval *big.Int) (common.Address, *types.Transaction, *HyperdriveEpoch, error) {
	parsed, err := abi.JSON(strings.NewReader(HyperdriveEpochABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(HyperdriveEpochBin), backend, _minimumEpochInterval)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &HyperdriveEpoch{HyperdriveEpochCaller: HyperdriveEpochCaller{contract: contract}, HyperdriveEpochTransactor: HyperdriveEpochTransactor{contract: contract}, HyperdriveEpochFilterer: HyperdriveEpochFilterer{contract: contract}}, nil
}

// HyperdriveEpoch is an auto generated Go binding around an Ethereum contract.
type HyperdriveEpoch struct {
	HyperdriveEpochCaller     // Read-only binding to the contract
	HyperdriveEpochTransactor // Write-only binding to the contract
	HyperdriveEpochFilterer   // Log filterer for contract events
}

// HyperdriveEpochCaller is an auto generated read-only Go binding around an Ethereum contract.
type HyperdriveEpochCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HyperdriveEpochTransactor is an auto generated write-only Go binding around an Ethereum contract.
type HyperdriveEpochTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HyperdriveEpochFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type HyperdriveEpochFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HyperdriveEpochSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type HyperdriveEpochSession struct {
	Contract     *HyperdriveEpoch  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// HyperdriveEpochCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type HyperdriveEpochCallerSession struct {
	Contract *HyperdriveEpochCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// HyperdriveEpochTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type HyperdriveEpochTransactorSession struct {
	Contract     *HyperdriveEpochTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// HyperdriveEpochRaw is an auto generated low-level Go binding around an Ethereum contract.
type HyperdriveEpochRaw struct {
	Contract *HyperdriveEpoch // Generic contract binding to access the raw methods on
}

// HyperdriveEpochCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type HyperdriveEpochCallerRaw struct {
	Contract *HyperdriveEpochCaller // Generic read-only contract binding to access the raw methods on
}

// HyperdriveEpochTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type HyperdriveEpochTransactorRaw struct {
	Contract *HyperdriveEpochTransactor // Generic write-only contract binding to access the raw methods on
}

// NewHyperdriveEpoch creates a new instance of HyperdriveEpoch, bound to a specific deployed contract.
func NewHyperdriveEpoch(address common.Address, backend bind.ContractBackend) (*HyperdriveEpoch, error) {
	contract, err := bindHyperdriveEpoch(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &HyperdriveEpoch{HyperdriveEpochCaller: HyperdriveEpochCaller{contract: contract}, HyperdriveEpochTransactor: HyperdriveEpochTransactor{contract: contract}, HyperdriveEpochFilterer: HyperdriveEpochFilterer{contract: contract}}, nil
}

// NewHyperdriveEpochCaller creates a new read-only instance of HyperdriveEpoch, bound to a specific deployed contract.
func NewHyperdriveEpochCaller(address common.Address, caller bind.ContractCaller) (*HyperdriveEpochCaller, error) {
	contract, err := bindHyperdriveEpoch(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &HyperdriveEpochCaller{contract: contract}, nil
}

// NewHyperdriveEpochTransactor creates a new write-only instance of HyperdriveEpoch, bound to a specific deployed contract.
func NewHyperdriveEpochTransactor(address common.Address, transactor bind.ContractTransactor) (*HyperdriveEpochTransactor, error) {
	contract, err := bindHyperdriveEpoch(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &HyperdriveEpochTransactor{contract: contract}, nil
}

// NewHyperdriveEpochFilterer creates a new log filterer instance of HyperdriveEpoch, bound to a specific deployed contract.
func NewHyperdriveEpochFilterer(address common.Address, filterer bind.ContractFilterer) (*HyperdriveEpochFilterer, error) {
	contract, err := bindHyperdriveEpoch(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &HyperdriveEpochFilterer{contract: contract}, nil
}

// bindHyperdriveEpoch binds a generic wrapper to an already deployed contract.
func bindHyperdriveEpoch(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(HyperdriveEpochABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_HyperdriveEpoch *HyperdriveEpochRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _HyperdriveEpoch.Contract.HyperdriveEpochCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_HyperdriveEpoch *HyperdriveEpochRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HyperdriveEpoch.Contract.HyperdriveEpochTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_HyperdriveEpoch *HyperdriveEpochRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _HyperdriveEpoch.Contract.HyperdriveEpochTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_HyperdriveEpoch *HyperdriveEpochCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _HyperdriveEpoch.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_HyperdriveEpoch *HyperdriveEpochTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HyperdriveEpoch.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_HyperdriveEpoch *HyperdriveEpochTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _HyperdriveEpoch.Contract.contract.Transact(opts, method, params...)
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() constant returns(blockhash uint256, timestamp uint256)
func (_HyperdriveEpoch *HyperdriveEpochCaller) CurrentEpoch(opts *bind.CallOpts) (struct {
	Blockhash *big.Int
	Timestamp *big.Int
}, error) {
	ret := new(struct {
		Blockhash *big.Int
		Timestamp *big.Int
	})
	out := ret
	err := _HyperdriveEpoch.contract.Call(opts, out, "currentEpoch")
	return *ret, err
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() constant returns(blockhash uint256, timestamp uint256)
func (_HyperdriveEpoch *HyperdriveEpochSession) CurrentEpoch() (struct {
	Blockhash *big.Int
	Timestamp *big.Int
}, error) {
	return _HyperdriveEpoch.Contract.CurrentEpoch(&_HyperdriveEpoch.CallOpts)
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() constant returns(blockhash uint256, timestamp uint256)
func (_HyperdriveEpoch *HyperdriveEpochCallerSession) CurrentEpoch() (struct {
	Blockhash *big.Int
	Timestamp *big.Int
}, error) {
	return _HyperdriveEpoch.Contract.CurrentEpoch(&_HyperdriveEpoch.CallOpts)
}

// MinimumEpochInterval is a free data retrieval call binding the contract method 0x55cacda5.
//
// Solidity: function minimumEpochInterval() constant returns(uint256)
func (_HyperdriveEpoch *HyperdriveEpochCaller) MinimumEpochInterval(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _HyperdriveEpoch.contract.Call(opts, out, "minimumEpochInterval")
	return *ret0, err
}

// MinimumEpochInterval is a free data retrieval call binding the contract method 0x55cacda5.
//
// Solidity: function minimumEpochInterval() constant returns(uint256)
func (_HyperdriveEpoch *HyperdriveEpochSession) MinimumEpochInterval() (*big.Int, error) {
	return _HyperdriveEpoch.Contract.MinimumEpochInterval(&_HyperdriveEpoch.CallOpts)
}

// MinimumEpochInterval is a free data retrieval call binding the contract method 0x55cacda5.
//
// Solidity: function minimumEpochInterval() constant returns(uint256)
func (_HyperdriveEpoch *HyperdriveEpochCallerSession) MinimumEpochInterval() (*big.Int, error) {
	return _HyperdriveEpoch.Contract.MinimumEpochInterval(&_HyperdriveEpoch.CallOpts)
}

// Epoch is a paid mutator transaction binding the contract method 0x900cf0cf.
//
// Solidity: function epoch() returns()
func (_HyperdriveEpoch *HyperdriveEpochTransactor) Epoch(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HyperdriveEpoch.contract.Transact(opts, "epoch")
}

// Epoch is a paid mutator transaction binding the contract method 0x900cf0cf.
//
// Solidity: function epoch() returns()
func (_HyperdriveEpoch *HyperdriveEpochSession) Epoch() (*types.Transaction, error) {
	return _HyperdriveEpoch.Contract.Epoch(&_HyperdriveEpoch.TransactOpts)
}

// Epoch is a paid mutator transaction binding the contract method 0x900cf0cf.
//
// Solidity: function epoch() returns()
func (_HyperdriveEpoch *HyperdriveEpochTransactorSession) Epoch() (*types.Transaction, error) {
	return _HyperdriveEpoch.Contract.Epoch(&_HyperdriveEpoch.TransactOpts)
}

// HyperdriveEpochNewEpochIterator is returned from FilterNewEpoch and is used to iterate over the raw logs and unpacked data for NewEpoch events raised by the HyperdriveEpoch contract.
type HyperdriveEpochNewEpochIterator struct {
	Event *HyperdriveEpochNewEpoch // Event containing the contract specifics and raw log

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
func (it *HyperdriveEpochNewEpochIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HyperdriveEpochNewEpoch)
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
		it.Event = new(HyperdriveEpochNewEpoch)
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
func (it *HyperdriveEpochNewEpochIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HyperdriveEpochNewEpochIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HyperdriveEpochNewEpoch represents a NewEpoch event raised by the HyperdriveEpoch contract.
type HyperdriveEpochNewEpoch struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterNewEpoch is a free log retrieval operation binding the contract event 0xe358419ca0dd7928a310d787a606dfae5d869f5071249efa6107105e7afc40bc.
//
// Solidity: event NewEpoch()
func (_HyperdriveEpoch *HyperdriveEpochFilterer) FilterNewEpoch(opts *bind.FilterOpts) (*HyperdriveEpochNewEpochIterator, error) {

	logs, sub, err := _HyperdriveEpoch.contract.FilterLogs(opts, "NewEpoch")
	if err != nil {
		return nil, err
	}
	return &HyperdriveEpochNewEpochIterator{contract: _HyperdriveEpoch.contract, event: "NewEpoch", logs: logs, sub: sub}, nil
}

// WatchNewEpoch is a free log subscription operation binding the contract event 0xe358419ca0dd7928a310d787a606dfae5d869f5071249efa6107105e7afc40bc.
//
// Solidity: event NewEpoch()
func (_HyperdriveEpoch *HyperdriveEpochFilterer) WatchNewEpoch(opts *bind.WatchOpts, sink chan<- *HyperdriveEpochNewEpoch) (event.Subscription, error) {

	logs, sub, err := _HyperdriveEpoch.contract.WatchLogs(opts, "NewEpoch")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HyperdriveEpochNewEpoch)
				if err := _HyperdriveEpoch.contract.UnpackLog(event, "NewEpoch", log); err != nil {
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

// LinkedListABI is the input ABI used to generate the binding from.
const LinkedListABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"NULL\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes20\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// LinkedListBin is the compiled bytecode used for deploying new contracts.
const LinkedListBin = `0x60b361002f600b82828239805160001a6073146000811461001f57610021565bfe5b5030600052607381538281f300730000000000000000000000000000000000000000301460806040526004361060555763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663f26be3fc8114605a575b600080fd5b60606082565b604080516bffffffffffffffffffffffff199092168252519081900360200190f35b6000815600a165627a7a72305820985e6b41267e59699a4ca854b54875525abc8abb42f03d3cd9ed5be2433983660029`

// DeployLinkedList deploys a new Ethereum contract, binding an instance of LinkedList to it.
func DeployLinkedList(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *LinkedList, error) {
	parsed, err := abi.JSON(strings.NewReader(LinkedListABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(LinkedListBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LinkedList{LinkedListCaller: LinkedListCaller{contract: contract}, LinkedListTransactor: LinkedListTransactor{contract: contract}, LinkedListFilterer: LinkedListFilterer{contract: contract}}, nil
}

// LinkedList is an auto generated Go binding around an Ethereum contract.
type LinkedList struct {
	LinkedListCaller     // Read-only binding to the contract
	LinkedListTransactor // Write-only binding to the contract
	LinkedListFilterer   // Log filterer for contract events
}

// LinkedListCaller is an auto generated read-only Go binding around an Ethereum contract.
type LinkedListCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LinkedListTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LinkedListTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LinkedListFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LinkedListFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LinkedListSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LinkedListSession struct {
	Contract     *LinkedList       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LinkedListCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LinkedListCallerSession struct {
	Contract *LinkedListCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// LinkedListTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LinkedListTransactorSession struct {
	Contract     *LinkedListTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// LinkedListRaw is an auto generated low-level Go binding around an Ethereum contract.
type LinkedListRaw struct {
	Contract *LinkedList // Generic contract binding to access the raw methods on
}

// LinkedListCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LinkedListCallerRaw struct {
	Contract *LinkedListCaller // Generic read-only contract binding to access the raw methods on
}

// LinkedListTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LinkedListTransactorRaw struct {
	Contract *LinkedListTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLinkedList creates a new instance of LinkedList, bound to a specific deployed contract.
func NewLinkedList(address common.Address, backend bind.ContractBackend) (*LinkedList, error) {
	contract, err := bindLinkedList(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LinkedList{LinkedListCaller: LinkedListCaller{contract: contract}, LinkedListTransactor: LinkedListTransactor{contract: contract}, LinkedListFilterer: LinkedListFilterer{contract: contract}}, nil
}

// NewLinkedListCaller creates a new read-only instance of LinkedList, bound to a specific deployed contract.
func NewLinkedListCaller(address common.Address, caller bind.ContractCaller) (*LinkedListCaller, error) {
	contract, err := bindLinkedList(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LinkedListCaller{contract: contract}, nil
}

// NewLinkedListTransactor creates a new write-only instance of LinkedList, bound to a specific deployed contract.
func NewLinkedListTransactor(address common.Address, transactor bind.ContractTransactor) (*LinkedListTransactor, error) {
	contract, err := bindLinkedList(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LinkedListTransactor{contract: contract}, nil
}

// NewLinkedListFilterer creates a new log filterer instance of LinkedList, bound to a specific deployed contract.
func NewLinkedListFilterer(address common.Address, filterer bind.ContractFilterer) (*LinkedListFilterer, error) {
	contract, err := bindLinkedList(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LinkedListFilterer{contract: contract}, nil
}

// bindLinkedList binds a generic wrapper to an already deployed contract.
func bindLinkedList(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LinkedListABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LinkedList *LinkedListRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LinkedList.Contract.LinkedListCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LinkedList *LinkedListRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LinkedList.Contract.LinkedListTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LinkedList *LinkedListRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LinkedList.Contract.LinkedListTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LinkedList *LinkedListCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LinkedList.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LinkedList *LinkedListTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LinkedList.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LinkedList *LinkedListTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LinkedList.Contract.contract.Transact(opts, method, params...)
}

// NULL is a free data retrieval call binding the contract method 0xf26be3fc.
//
// Solidity: function NULL() constant returns(bytes20)
func (_LinkedList *LinkedListCaller) NULL(opts *bind.CallOpts) ([20]byte, error) {
	var (
		ret0 = new([20]byte)
	)
	out := ret0
	err := _LinkedList.contract.Call(opts, out, "NULL")
	return *ret0, err
}

// NULL is a free data retrieval call binding the contract method 0xf26be3fc.
//
// Solidity: function NULL() constant returns(bytes20)
func (_LinkedList *LinkedListSession) NULL() ([20]byte, error) {
	return _LinkedList.Contract.NULL(&_LinkedList.CallOpts)
}

// NULL is a free data retrieval call binding the contract method 0xf26be3fc.
//
// Solidity: function NULL() constant returns(bytes20)
func (_LinkedList *LinkedListCallerSession) NULL() ([20]byte, error) {
	return _LinkedList.Contract.NULL(&_LinkedList.CallOpts)
}
