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

// BasicTokenABI is the input ABI used to generate the binding from.
const BasicTokenABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]"

// BasicTokenBin is the compiled bytecode used for deploying new contracts.
const BasicTokenBin = `0x608060405234801561001057600080fd5b5061027c806100206000396000f3006080604052600436106100565763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166318160ddd811461005b57806370a0823114610082578063a9059cbb146100b0575b600080fd5b34801561006757600080fd5b506100706100f5565b60408051918252519081900360200190f35b34801561008e57600080fd5b5061007073ffffffffffffffffffffffffffffffffffffffff600435166100fb565b3480156100bc57600080fd5b506100e173ffffffffffffffffffffffffffffffffffffffff60043516602435610123565b604080519115158252519081900360200190f35b60015490565b73ffffffffffffffffffffffffffffffffffffffff1660009081526020819052604090205490565b600073ffffffffffffffffffffffffffffffffffffffff8316151561014757600080fd5b3360009081526020819052604090205482111561016357600080fd5b33600090815260208190526040902054610183908363ffffffff61022b16565b336000908152602081905260408082209290925573ffffffffffffffffffffffffffffffffffffffff8516815220546101c2908363ffffffff61023d16565b73ffffffffffffffffffffffffffffffffffffffff8416600081815260208181526040918290209390935580518581529051919233927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9281900390910190a350600192915050565b60008282111561023757fe5b50900390565b8181018281101561024a57fe5b929150505600a165627a7a723058203380cc6307c525bb245b03e37aa15c5d6e8f025e3bfc8d1c227bfd08aad6de020029`

// DeployBasicToken deploys a new Ethereum contract, binding an instance of BasicToken to it.
func DeployBasicToken(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BasicToken, error) {
	parsed, err := abi.JSON(strings.NewReader(BasicTokenABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(BasicTokenBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BasicToken{BasicTokenCaller: BasicTokenCaller{contract: contract}, BasicTokenTransactor: BasicTokenTransactor{contract: contract}, BasicTokenFilterer: BasicTokenFilterer{contract: contract}}, nil
}

// BasicToken is an auto generated Go binding around an Ethereum contract.
type BasicToken struct {
	BasicTokenCaller     // Read-only binding to the contract
	BasicTokenTransactor // Write-only binding to the contract
	BasicTokenFilterer   // Log filterer for contract events
}

// BasicTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type BasicTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BasicTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BasicTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BasicTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BasicTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BasicTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BasicTokenSession struct {
	Contract     *BasicToken       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BasicTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BasicTokenCallerSession struct {
	Contract *BasicTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// BasicTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BasicTokenTransactorSession struct {
	Contract     *BasicTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// BasicTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type BasicTokenRaw struct {
	Contract *BasicToken // Generic contract binding to access the raw methods on
}

// BasicTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BasicTokenCallerRaw struct {
	Contract *BasicTokenCaller // Generic read-only contract binding to access the raw methods on
}

// BasicTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BasicTokenTransactorRaw struct {
	Contract *BasicTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBasicToken creates a new instance of BasicToken, bound to a specific deployed contract.
func NewBasicToken(address common.Address, backend bind.ContractBackend) (*BasicToken, error) {
	contract, err := bindBasicToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BasicToken{BasicTokenCaller: BasicTokenCaller{contract: contract}, BasicTokenTransactor: BasicTokenTransactor{contract: contract}, BasicTokenFilterer: BasicTokenFilterer{contract: contract}}, nil
}

// NewBasicTokenCaller creates a new read-only instance of BasicToken, bound to a specific deployed contract.
func NewBasicTokenCaller(address common.Address, caller bind.ContractCaller) (*BasicTokenCaller, error) {
	contract, err := bindBasicToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BasicTokenCaller{contract: contract}, nil
}

// NewBasicTokenTransactor creates a new write-only instance of BasicToken, bound to a specific deployed contract.
func NewBasicTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*BasicTokenTransactor, error) {
	contract, err := bindBasicToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BasicTokenTransactor{contract: contract}, nil
}

// NewBasicTokenFilterer creates a new log filterer instance of BasicToken, bound to a specific deployed contract.
func NewBasicTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*BasicTokenFilterer, error) {
	contract, err := bindBasicToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BasicTokenFilterer{contract: contract}, nil
}

// bindBasicToken binds a generic wrapper to an already deployed contract.
func bindBasicToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BasicTokenABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BasicToken *BasicTokenRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _BasicToken.Contract.BasicTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BasicToken *BasicTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BasicToken.Contract.BasicTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BasicToken *BasicTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BasicToken.Contract.BasicTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BasicToken *BasicTokenCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _BasicToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BasicToken *BasicTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BasicToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BasicToken *BasicTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BasicToken.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(uint256)
func (_BasicToken *BasicTokenCaller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _BasicToken.contract.Call(opts, out, "balanceOf", _owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(uint256)
func (_BasicToken *BasicTokenSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _BasicToken.Contract.BalanceOf(&_BasicToken.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(uint256)
func (_BasicToken *BasicTokenCallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _BasicToken.Contract.BalanceOf(&_BasicToken.CallOpts, _owner)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_BasicToken *BasicTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _BasicToken.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_BasicToken *BasicTokenSession) TotalSupply() (*big.Int, error) {
	return _BasicToken.Contract.TotalSupply(&_BasicToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_BasicToken *BasicTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _BasicToken.Contract.TotalSupply(&_BasicToken.CallOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_BasicToken *BasicTokenTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _BasicToken.contract.Transact(opts, "transfer", _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_BasicToken *BasicTokenSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _BasicToken.Contract.Transfer(&_BasicToken.TransactOpts, _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_BasicToken *BasicTokenTransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _BasicToken.Contract.Transfer(&_BasicToken.TransactOpts, _to, _value)
}

// BasicTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the BasicToken contract.
type BasicTokenTransferIterator struct {
	Event *BasicTokenTransfer // Event containing the contract specifics and raw log

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
func (it *BasicTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BasicTokenTransfer)
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
		it.Event = new(BasicTokenTransfer)
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
func (it *BasicTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BasicTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BasicTokenTransfer represents a Transfer event raised by the BasicToken contract.
type BasicTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
func (_BasicToken *BasicTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BasicTokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BasicToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BasicTokenTransferIterator{contract: _BasicToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
func (_BasicToken *BasicTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *BasicTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BasicToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BasicTokenTransfer)
				if err := _BasicToken.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// BindingsABI is the input ABI used to generate the binding from.
const BindingsABI = "[]"

// BindingsBin is the compiled bytecode used for deploying new contracts.
const BindingsBin = `0x6080604052348015600f57600080fd5b50603580601d6000396000f3006080604052600080fd00a165627a7a72305820cc45a792f39fb10a8e6abe01aa6054f086c6c31d46bc68080a5b27970cf76c140029`

// DeployBindings deploys a new Ethereum contract, binding an instance of Bindings to it.
func DeployBindings(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Bindings, error) {
	parsed, err := abi.JSON(strings.NewReader(BindingsABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(BindingsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Bindings{BindingsCaller: BindingsCaller{contract: contract}, BindingsTransactor: BindingsTransactor{contract: contract}, BindingsFilterer: BindingsFilterer{contract: contract}}, nil
}

// Bindings is an auto generated Go binding around an Ethereum contract.
type Bindings struct {
	BindingsCaller     // Read-only binding to the contract
	BindingsTransactor // Write-only binding to the contract
	BindingsFilterer   // Log filterer for contract events
}

// BindingsCaller is an auto generated read-only Go binding around an Ethereum contract.
type BindingsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BindingsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BindingsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BindingsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BindingsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BindingsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BindingsSession struct {
	Contract     *Bindings         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BindingsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BindingsCallerSession struct {
	Contract *BindingsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// BindingsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BindingsTransactorSession struct {
	Contract     *BindingsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// BindingsRaw is an auto generated low-level Go binding around an Ethereum contract.
type BindingsRaw struct {
	Contract *Bindings // Generic contract binding to access the raw methods on
}

// BindingsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BindingsCallerRaw struct {
	Contract *BindingsCaller // Generic read-only contract binding to access the raw methods on
}

// BindingsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BindingsTransactorRaw struct {
	Contract *BindingsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBindings creates a new instance of Bindings, bound to a specific deployed contract.
func NewBindings(address common.Address, backend bind.ContractBackend) (*Bindings, error) {
	contract, err := bindBindings(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Bindings{BindingsCaller: BindingsCaller{contract: contract}, BindingsTransactor: BindingsTransactor{contract: contract}, BindingsFilterer: BindingsFilterer{contract: contract}}, nil
}

// NewBindingsCaller creates a new read-only instance of Bindings, bound to a specific deployed contract.
func NewBindingsCaller(address common.Address, caller bind.ContractCaller) (*BindingsCaller, error) {
	contract, err := bindBindings(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BindingsCaller{contract: contract}, nil
}

// NewBindingsTransactor creates a new write-only instance of Bindings, bound to a specific deployed contract.
func NewBindingsTransactor(address common.Address, transactor bind.ContractTransactor) (*BindingsTransactor, error) {
	contract, err := bindBindings(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BindingsTransactor{contract: contract}, nil
}

// NewBindingsFilterer creates a new log filterer instance of Bindings, bound to a specific deployed contract.
func NewBindingsFilterer(address common.Address, filterer bind.ContractFilterer) (*BindingsFilterer, error) {
	contract, err := bindBindings(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BindingsFilterer{contract: contract}, nil
}

// bindBindings binds a generic wrapper to an already deployed contract.
func bindBindings(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BindingsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bindings *BindingsRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Bindings.Contract.BindingsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bindings *BindingsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bindings.Contract.BindingsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bindings *BindingsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bindings.Contract.BindingsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bindings *BindingsCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Bindings.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bindings *BindingsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bindings.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bindings *BindingsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bindings.Contract.contract.Transact(opts, method, params...)
}

// BrokerVerifierABI is the input ABI used to generate the binding from.
const BrokerVerifierABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_trader\",\"type\":\"address\"},{\"name\":\"_signature\",\"type\":\"bytes\"},{\"name\":\"_orderID\",\"type\":\"bytes32\"}],\"name\":\"verifyOpenSignature\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// BrokerVerifierBin is the compiled bytecode used for deploying new contracts.
const BrokerVerifierBin = `0x`

// DeployBrokerVerifier deploys a new Ethereum contract, binding an instance of BrokerVerifier to it.
func DeployBrokerVerifier(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BrokerVerifier, error) {
	parsed, err := abi.JSON(strings.NewReader(BrokerVerifierABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(BrokerVerifierBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BrokerVerifier{BrokerVerifierCaller: BrokerVerifierCaller{contract: contract}, BrokerVerifierTransactor: BrokerVerifierTransactor{contract: contract}, BrokerVerifierFilterer: BrokerVerifierFilterer{contract: contract}}, nil
}

// BrokerVerifier is an auto generated Go binding around an Ethereum contract.
type BrokerVerifier struct {
	BrokerVerifierCaller     // Read-only binding to the contract
	BrokerVerifierTransactor // Write-only binding to the contract
	BrokerVerifierFilterer   // Log filterer for contract events
}

// BrokerVerifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type BrokerVerifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BrokerVerifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BrokerVerifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BrokerVerifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BrokerVerifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BrokerVerifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BrokerVerifierSession struct {
	Contract     *BrokerVerifier   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BrokerVerifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BrokerVerifierCallerSession struct {
	Contract *BrokerVerifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// BrokerVerifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BrokerVerifierTransactorSession struct {
	Contract     *BrokerVerifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// BrokerVerifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type BrokerVerifierRaw struct {
	Contract *BrokerVerifier // Generic contract binding to access the raw methods on
}

// BrokerVerifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BrokerVerifierCallerRaw struct {
	Contract *BrokerVerifierCaller // Generic read-only contract binding to access the raw methods on
}

// BrokerVerifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BrokerVerifierTransactorRaw struct {
	Contract *BrokerVerifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBrokerVerifier creates a new instance of BrokerVerifier, bound to a specific deployed contract.
func NewBrokerVerifier(address common.Address, backend bind.ContractBackend) (*BrokerVerifier, error) {
	contract, err := bindBrokerVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BrokerVerifier{BrokerVerifierCaller: BrokerVerifierCaller{contract: contract}, BrokerVerifierTransactor: BrokerVerifierTransactor{contract: contract}, BrokerVerifierFilterer: BrokerVerifierFilterer{contract: contract}}, nil
}

// NewBrokerVerifierCaller creates a new read-only instance of BrokerVerifier, bound to a specific deployed contract.
func NewBrokerVerifierCaller(address common.Address, caller bind.ContractCaller) (*BrokerVerifierCaller, error) {
	contract, err := bindBrokerVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BrokerVerifierCaller{contract: contract}, nil
}

// NewBrokerVerifierTransactor creates a new write-only instance of BrokerVerifier, bound to a specific deployed contract.
func NewBrokerVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*BrokerVerifierTransactor, error) {
	contract, err := bindBrokerVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BrokerVerifierTransactor{contract: contract}, nil
}

// NewBrokerVerifierFilterer creates a new log filterer instance of BrokerVerifier, bound to a specific deployed contract.
func NewBrokerVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*BrokerVerifierFilterer, error) {
	contract, err := bindBrokerVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BrokerVerifierFilterer{contract: contract}, nil
}

// bindBrokerVerifier binds a generic wrapper to an already deployed contract.
func bindBrokerVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BrokerVerifierABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BrokerVerifier *BrokerVerifierRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _BrokerVerifier.Contract.BrokerVerifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BrokerVerifier *BrokerVerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BrokerVerifier.Contract.BrokerVerifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BrokerVerifier *BrokerVerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BrokerVerifier.Contract.BrokerVerifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BrokerVerifier *BrokerVerifierCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _BrokerVerifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BrokerVerifier *BrokerVerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BrokerVerifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BrokerVerifier *BrokerVerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BrokerVerifier.Contract.contract.Transact(opts, method, params...)
}

// VerifyOpenSignature is a paid mutator transaction binding the contract method 0x472e1910.
//
// Solidity: function verifyOpenSignature(_trader address, _signature bytes, _orderID bytes32) returns(bool)
func (_BrokerVerifier *BrokerVerifierTransactor) VerifyOpenSignature(opts *bind.TransactOpts, _trader common.Address, _signature []byte, _orderID [32]byte) (*types.Transaction, error) {
	return _BrokerVerifier.contract.Transact(opts, "verifyOpenSignature", _trader, _signature, _orderID)
}

// VerifyOpenSignature is a paid mutator transaction binding the contract method 0x472e1910.
//
// Solidity: function verifyOpenSignature(_trader address, _signature bytes, _orderID bytes32) returns(bool)
func (_BrokerVerifier *BrokerVerifierSession) VerifyOpenSignature(_trader common.Address, _signature []byte, _orderID [32]byte) (*types.Transaction, error) {
	return _BrokerVerifier.Contract.VerifyOpenSignature(&_BrokerVerifier.TransactOpts, _trader, _signature, _orderID)
}

// VerifyOpenSignature is a paid mutator transaction binding the contract method 0x472e1910.
//
// Solidity: function verifyOpenSignature(_trader address, _signature bytes, _orderID bytes32) returns(bool)
func (_BrokerVerifier *BrokerVerifierTransactorSession) VerifyOpenSignature(_trader common.Address, _signature []byte, _orderID [32]byte) (*types.Transaction, error) {
	return _BrokerVerifier.Contract.VerifyOpenSignature(&_BrokerVerifier.TransactOpts, _trader, _signature, _orderID)
}

// BurnableTokenABI is the input ABI used to generate the binding from.
const BurnableTokenABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"burner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Burn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]"

// BurnableTokenBin is the compiled bytecode used for deploying new contracts.
const BurnableTokenBin = `0x608060405234801561001057600080fd5b50610361806100206000396000f3006080604052600436106100615763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166318160ddd811461006657806342966c681461008d57806370a08231146100a7578063a9059cbb146100c8575b600080fd5b34801561007257600080fd5b5061007b610100565b60408051918252519081900360200190f35b34801561009957600080fd5b506100a5600435610106565b005b3480156100b357600080fd5b5061007b600160a060020a0360043516610113565b3480156100d457600080fd5b506100ec600160a060020a036004351660243561012e565b604080519115158252519081900360200190f35b60015490565b610110338261020f565b50565b600160a060020a031660009081526020819052604090205490565b6000600160a060020a038316151561014557600080fd5b3360009081526020819052604090205482111561016157600080fd5b33600090815260208190526040902054610181908363ffffffff61031016565b3360009081526020819052604080822092909255600160a060020a038516815220546101b3908363ffffffff61032216565b600160a060020a038416600081815260208181526040918290209390935580518581529051919233927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9281900390910190a350600192915050565b600160a060020a03821660009081526020819052604090205481111561023457600080fd5b600160a060020a03821660009081526020819052604090205461025d908263ffffffff61031016565b600160a060020a038316600090815260208190526040902055600154610289908263ffffffff61031016565b600155604080518281529051600160a060020a038416917fcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5919081900360200190a2604080518281529051600091600160a060020a038516917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9181900360200190a35050565b60008282111561031c57fe5b50900390565b8181018281101561032f57fe5b929150505600a165627a7a723058208c240b47ec8c59dccb593c7831125d41919c48ee2b0da7340749261d0e0bf5bf0029`

// DeployBurnableToken deploys a new Ethereum contract, binding an instance of BurnableToken to it.
func DeployBurnableToken(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BurnableToken, error) {
	parsed, err := abi.JSON(strings.NewReader(BurnableTokenABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(BurnableTokenBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BurnableToken{BurnableTokenCaller: BurnableTokenCaller{contract: contract}, BurnableTokenTransactor: BurnableTokenTransactor{contract: contract}, BurnableTokenFilterer: BurnableTokenFilterer{contract: contract}}, nil
}

// BurnableToken is an auto generated Go binding around an Ethereum contract.
type BurnableToken struct {
	BurnableTokenCaller     // Read-only binding to the contract
	BurnableTokenTransactor // Write-only binding to the contract
	BurnableTokenFilterer   // Log filterer for contract events
}

// BurnableTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type BurnableTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BurnableTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BurnableTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BurnableTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BurnableTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BurnableTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BurnableTokenSession struct {
	Contract     *BurnableToken    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BurnableTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BurnableTokenCallerSession struct {
	Contract *BurnableTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// BurnableTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BurnableTokenTransactorSession struct {
	Contract     *BurnableTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// BurnableTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type BurnableTokenRaw struct {
	Contract *BurnableToken // Generic contract binding to access the raw methods on
}

// BurnableTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BurnableTokenCallerRaw struct {
	Contract *BurnableTokenCaller // Generic read-only contract binding to access the raw methods on
}

// BurnableTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BurnableTokenTransactorRaw struct {
	Contract *BurnableTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBurnableToken creates a new instance of BurnableToken, bound to a specific deployed contract.
func NewBurnableToken(address common.Address, backend bind.ContractBackend) (*BurnableToken, error) {
	contract, err := bindBurnableToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BurnableToken{BurnableTokenCaller: BurnableTokenCaller{contract: contract}, BurnableTokenTransactor: BurnableTokenTransactor{contract: contract}, BurnableTokenFilterer: BurnableTokenFilterer{contract: contract}}, nil
}

// NewBurnableTokenCaller creates a new read-only instance of BurnableToken, bound to a specific deployed contract.
func NewBurnableTokenCaller(address common.Address, caller bind.ContractCaller) (*BurnableTokenCaller, error) {
	contract, err := bindBurnableToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BurnableTokenCaller{contract: contract}, nil
}

// NewBurnableTokenTransactor creates a new write-only instance of BurnableToken, bound to a specific deployed contract.
func NewBurnableTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*BurnableTokenTransactor, error) {
	contract, err := bindBurnableToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BurnableTokenTransactor{contract: contract}, nil
}

// NewBurnableTokenFilterer creates a new log filterer instance of BurnableToken, bound to a specific deployed contract.
func NewBurnableTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*BurnableTokenFilterer, error) {
	contract, err := bindBurnableToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BurnableTokenFilterer{contract: contract}, nil
}

// bindBurnableToken binds a generic wrapper to an already deployed contract.
func bindBurnableToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BurnableTokenABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BurnableToken *BurnableTokenRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _BurnableToken.Contract.BurnableTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BurnableToken *BurnableTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnableToken.Contract.BurnableTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BurnableToken *BurnableTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnableToken.Contract.BurnableTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BurnableToken *BurnableTokenCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _BurnableToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BurnableToken *BurnableTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnableToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BurnableToken *BurnableTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnableToken.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(uint256)
func (_BurnableToken *BurnableTokenCaller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _BurnableToken.contract.Call(opts, out, "balanceOf", _owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(uint256)
func (_BurnableToken *BurnableTokenSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _BurnableToken.Contract.BalanceOf(&_BurnableToken.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(uint256)
func (_BurnableToken *BurnableTokenCallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _BurnableToken.Contract.BalanceOf(&_BurnableToken.CallOpts, _owner)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_BurnableToken *BurnableTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _BurnableToken.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_BurnableToken *BurnableTokenSession) TotalSupply() (*big.Int, error) {
	return _BurnableToken.Contract.TotalSupply(&_BurnableToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_BurnableToken *BurnableTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _BurnableToken.Contract.TotalSupply(&_BurnableToken.CallOpts)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(_value uint256) returns()
func (_BurnableToken *BurnableTokenTransactor) Burn(opts *bind.TransactOpts, _value *big.Int) (*types.Transaction, error) {
	return _BurnableToken.contract.Transact(opts, "burn", _value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(_value uint256) returns()
func (_BurnableToken *BurnableTokenSession) Burn(_value *big.Int) (*types.Transaction, error) {
	return _BurnableToken.Contract.Burn(&_BurnableToken.TransactOpts, _value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(_value uint256) returns()
func (_BurnableToken *BurnableTokenTransactorSession) Burn(_value *big.Int) (*types.Transaction, error) {
	return _BurnableToken.Contract.Burn(&_BurnableToken.TransactOpts, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_BurnableToken *BurnableTokenTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _BurnableToken.contract.Transact(opts, "transfer", _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_BurnableToken *BurnableTokenSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _BurnableToken.Contract.Transfer(&_BurnableToken.TransactOpts, _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_BurnableToken *BurnableTokenTransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _BurnableToken.Contract.Transfer(&_BurnableToken.TransactOpts, _to, _value)
}

// BurnableTokenBurnIterator is returned from FilterBurn and is used to iterate over the raw logs and unpacked data for Burn events raised by the BurnableToken contract.
type BurnableTokenBurnIterator struct {
	Event *BurnableTokenBurn // Event containing the contract specifics and raw log

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
func (it *BurnableTokenBurnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnableTokenBurn)
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
		it.Event = new(BurnableTokenBurn)
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
func (it *BurnableTokenBurnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BurnableTokenBurnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BurnableTokenBurn represents a Burn event raised by the BurnableToken contract.
type BurnableTokenBurn struct {
	Burner common.Address
	Value  *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBurn is a free log retrieval operation binding the contract event 0xcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5.
//
// Solidity: e Burn(burner indexed address, value uint256)
func (_BurnableToken *BurnableTokenFilterer) FilterBurn(opts *bind.FilterOpts, burner []common.Address) (*BurnableTokenBurnIterator, error) {

	var burnerRule []interface{}
	for _, burnerItem := range burner {
		burnerRule = append(burnerRule, burnerItem)
	}

	logs, sub, err := _BurnableToken.contract.FilterLogs(opts, "Burn", burnerRule)
	if err != nil {
		return nil, err
	}
	return &BurnableTokenBurnIterator{contract: _BurnableToken.contract, event: "Burn", logs: logs, sub: sub}, nil
}

// WatchBurn is a free log subscription operation binding the contract event 0xcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5.
//
// Solidity: e Burn(burner indexed address, value uint256)
func (_BurnableToken *BurnableTokenFilterer) WatchBurn(opts *bind.WatchOpts, sink chan<- *BurnableTokenBurn, burner []common.Address) (event.Subscription, error) {

	var burnerRule []interface{}
	for _, burnerItem := range burner {
		burnerRule = append(burnerRule, burnerItem)
	}

	logs, sub, err := _BurnableToken.contract.WatchLogs(opts, "Burn", burnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BurnableTokenBurn)
				if err := _BurnableToken.contract.UnpackLog(event, "Burn", log); err != nil {
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

// BurnableTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the BurnableToken contract.
type BurnableTokenTransferIterator struct {
	Event *BurnableTokenTransfer // Event containing the contract specifics and raw log

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
func (it *BurnableTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnableTokenTransfer)
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
		it.Event = new(BurnableTokenTransfer)
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
func (it *BurnableTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BurnableTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BurnableTokenTransfer represents a Transfer event raised by the BurnableToken contract.
type BurnableTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
func (_BurnableToken *BurnableTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnableTokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnableToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BurnableTokenTransferIterator{contract: _BurnableToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
func (_BurnableToken *BurnableTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *BurnableTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnableToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BurnableTokenTransfer)
				if err := _BurnableToken.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// DarknodeRegistryABI is the input ABI used to generate the binding from.
const DarknodeRegistryABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"_darknodeID\",\"type\":\"address\"}],\"name\":\"isPendingRegistration\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numDarknodesNextEpoch\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_darknodeID\",\"type\":\"address\"},{\"name\":\"_publicKey\",\"type\":\"bytes\"},{\"name\":\"_bond\",\"type\":\"uint256\"}],\"name\":\"register\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_nextMinimumBond\",\"type\":\"uint256\"}],\"name\":\"updateMinimumBond\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numDarknodes\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_darknodeID\",\"type\":\"address\"}],\"name\":\"getDarknodeOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"nextSlasher\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_darknodeID\",\"type\":\"address\"}],\"name\":\"isPendingDeregistration\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_start\",\"type\":\"address\"},{\"name\":\"_count\",\"type\":\"uint256\"}],\"name\":\"getPreviousDarknodes\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"nextMinimumEpochInterval\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minimumEpochInterval\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_prover\",\"type\":\"address\"},{\"name\":\"_challenger1\",\"type\":\"address\"},{\"name\":\"_challenger2\",\"type\":\"address\"}],\"name\":\"slash\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_darknodeID\",\"type\":\"address\"}],\"name\":\"isRefundable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"previousEpoch\",\"outputs\":[{\"name\":\"epochhash\",\"type\":\"uint256\"},{\"name\":\"blocknumber\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"nextMinimumBond\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_nextMinimumEpochInterval\",\"type\":\"uint256\"}],\"name\":\"updateMinimumEpochInterval\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"nextMinimumPodSize\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numDarknodesPreviousEpoch\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentEpoch\",\"outputs\":[{\"name\":\"epochhash\",\"type\":\"uint256\"},{\"name\":\"blocknumber\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_darknodeID\",\"type\":\"address\"}],\"name\":\"isRegisteredInPreviousEpoch\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_darknodeID\",\"type\":\"address\"}],\"name\":\"isDeregistered\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_nextMinimumPodSize\",\"type\":\"uint256\"}],\"name\":\"updateMinimumPodSize\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_darknodeID\",\"type\":\"address\"}],\"name\":\"deregister\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_darknodeID\",\"type\":\"address\"}],\"name\":\"getDarknodePublicKey\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"ren\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"epoch\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"store\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minimumBond\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"slasher\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_slasher\",\"type\":\"address\"}],\"name\":\"updateSlasher\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_darknodeID\",\"type\":\"address\"}],\"name\":\"getDarknodeBond\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferStoreOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_darknodeID\",\"type\":\"address\"}],\"name\":\"isRegistered\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minimumPodSize\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_darknodeID\",\"type\":\"address\"}],\"name\":\"isDeregisterable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_start\",\"type\":\"address\"},{\"name\":\"_count\",\"type\":\"uint256\"}],\"name\":\"getDarknodes\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_darknodeID\",\"type\":\"address\"}],\"name\":\"refund\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_darknodeID\",\"type\":\"address\"}],\"name\":\"isRefunded\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_VERSION\",\"type\":\"string\"},{\"name\":\"_renAddress\",\"type\":\"address\"},{\"name\":\"_storeAddress\",\"type\":\"address\"},{\"name\":\"_minimumBond\",\"type\":\"uint256\"},{\"name\":\"_minimumPodSize\",\"type\":\"uint256\"},{\"name\":\"_minimumEpochInterval\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_darknodeID\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_bond\",\"type\":\"uint256\"}],\"name\":\"LogDarknodeRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_darknodeID\",\"type\":\"address\"}],\"name\":\"LogDarknodeDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"LogDarknodeOwnerRefunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"LogNewEpoch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"previousMinimumBond\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"nextMinimumBond\",\"type\":\"uint256\"}],\"name\":\"LogMinimumBondUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"previousMinimumPodSize\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"nextMinimumPodSize\",\"type\":\"uint256\"}],\"name\":\"LogMinimumPodSizeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"previousMinimumEpochInterval\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"nextMinimumEpochInterval\",\"type\":\"uint256\"}],\"name\":\"LogMinimumEpochIntervalUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"previousSlasher\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"nextSlasher\",\"type\":\"address\"}],\"name\":\"LogSlasherUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"}],\"name\":\"OwnershipRenounced\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"

// DarknodeRegistryBin is the compiled bytecode used for deploying new contracts.
const DarknodeRegistryBin = `0x60806040523480156200001157600080fd5b5060405162002701380380620027018339810160409081528151602080840151928401516060850151608086015160a087015160008054600160a060020a031916331790559490960180519096929491936200007391600191890190620000f7565b5060128054600160a060020a03958616600160a060020a031991821617909155601180549690951695169490941790925560058190556009556006819055600a556007819055600b555060408051808201909152436000198101408083526020909201819052600d91909155600e556000600281905560038190556004556200019c565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106200013a57805160ff19168380011785556200016a565b828001600101855582156200016a579182015b828111156200016a5782518255916020019190600101906200014d565b50620001789291506200017c565b5090565b6200019991905b8082111562000178576000815560010162000183565b90565b61255580620001ac6000396000f3006080604052600436106101ea5763ffffffff60e060020a600035041663040fa05181146101ef5780630847e9fa146102245780630aeb6b401461024b5780630ff9aafe1461027d5780631460e603146102955780631cedf8a3146102aa57806321a2ad3a146102e7578063303ee989146102fc578063438460741461031d578063455dc46d1461039157806355cacda5146103a6578063563bf264146103bb5780635aebd1cb146103e85780635cdaab481461040957806360a22fe41461043757806363b851b91461044c578063702c25ee14610464578063715018a61461047957806371740d161461048e57806376671808146104a35780637be266da146104b85780638020fc1f146104d957806380a0c461146104fa57806384ac33ec1461051257806384d2688c146105335780638a9b4067146105c95780638da5cb5b146105de578063900cf0cf146105f3578063975057e714610608578063aa7517e11461061d578063b134427114610632578063b3139d3814610647578063ba0f5b2014610668578063c2250a9914610689578063c3c5a547146106aa578063c7dbc2be146106cb578063e1878925146106e0578063ec5325c114610701578063f2fde38b14610725578063fa89401a14610746578063ffa1ad7414610767578063ffc9152e1461077c575b600080fd5b3480156101fb57600080fd5b50610210600160a060020a036004351661079d565b604080519115158252519081900360200190f35b34801561023057600080fd5b5061023961084b565b60408051918252519081900360200190f35b34801561025757600080fd5b5061027b60048035600160a060020a03169060248035908101910135604435610851565b005b34801561028957600080fd5b5061027b600435610bf0565b3480156102a157600080fd5b50610239610c0c565b3480156102b657600080fd5b506102cb600160a060020a0360043516610c12565b60408051600160a060020a039092168252519081900360200190f35b3480156102f357600080fd5b506102cb610c99565b34801561030857600080fd5b50610210600160a060020a0360043516610ca8565b34801561032957600080fd5b50610341600160a060020a0360043516602435610cfa565b60408051602080825283518183015283519192839290830191858101910280838360005b8381101561037d578181015183820152602001610365565b505050509050019250505060405180910390f35b34801561039d57600080fd5b50610239610d1d565b3480156103b257600080fd5b50610239610d23565b3480156103c757600080fd5b5061027b600160a060020a0360043581169060243581169060443516610d29565b3480156103f457600080fd5b50610210600160a060020a03600435166111af565b34801561041557600080fd5b5061041e61124a565b6040805192835260208301919091528051918290030190f35b34801561044357600080fd5b50610239611253565b34801561045857600080fd5b5061027b600435611259565b34801561047057600080fd5b50610239611275565b34801561048557600080fd5b5061027b61127b565b34801561049a57600080fd5b506102396112e7565b3480156104af57600080fd5b5061041e6112ed565b3480156104c457600080fd5b50610210600160a060020a03600435166112f6565b3480156104e557600080fd5b50610210600160a060020a036004351661131b565b34801561050657600080fd5b5061027b6004356113b0565b34801561051e57600080fd5b5061027b600160a060020a03600435166113cc565b34801561053f57600080fd5b50610554600160a060020a03600435166115e0565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561058e578181015183820152602001610576565b50505050905090810190601f1680156105bb5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156105d557600080fd5b506102cb6116d9565b3480156105ea57600080fd5b506102cb6116e8565b3480156105ff57600080fd5b5061027b6116f7565b34801561061457600080fd5b506102cb61199f565b34801561062957600080fd5b506102396119ae565b34801561063e57600080fd5b506102cb6119b4565b34801561065357600080fd5b5061027b600160a060020a03600435166119c3565b34801561067457600080fd5b50610239600160a060020a0360043516611a09565b34801561069557600080fd5b5061027b600160a060020a0360043516611a74565b3480156106b657600080fd5b50610210600160a060020a0360043516611b0d565b3480156106d757600080fd5b50610239611b32565b3480156106ec57600080fd5b50610210600160a060020a0360043516611b38565b34801561070d57600080fd5b50610341600160a060020a0360043516602435611bcf565b34801561073157600080fd5b5061027b600160a060020a0360043516611bea565b34801561075257600080fd5b5061027b600160a060020a0360043516611c0d565b34801561077357600080fd5b50610554611f10565b34801561078857600080fd5b50610210600160a060020a0360043516611f9d565b601254604080517fe2d7c64c000000000000000000000000000000000000000000000000000000008152600160a060020a03848116600483015291516000938493169163e2d7c64c91602480830192602092919082900301818787803b15801561080657600080fd5b505af115801561081a573d6000803e3d6000fd5b505050506040513d602081101561083057600080fd5b5051905080158015906108445750600e5481115b9392505050565b60035481565b8361085b81611f9d565b15156108d6576040805160e560020a62461bcd028152602060048201526024808201527f6d75737420626520726566756e646564206f72206e657665722072656769737460448201527f6572656400000000000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b600554821015610930576040805160e560020a62461bcd02815260206004820152601160248201527f696e73756666696369656e7420626f6e64000000000000000000000000000000604482015290519081900360640190fd5b601154604080517f23b872dd000000000000000000000000000000000000000000000000000000008152336004820152306024820152604481018590529051600160a060020a03909216916323b872dd916064808201926020929091908290030181600087803b1580156109a357600080fd5b505af11580156109b7573d6000803e3d6000fd5b505050506040513d60208110156109cd57600080fd5b50511515610a25576040805160e560020a62461bcd02815260206004820152601460248201527f626f6e64207472616e73666572206661696c6564000000000000000000000000604482015290519081900360640190fd5b601154601254604080517fa9059cbb000000000000000000000000000000000000000000000000000000008152600160a060020a039283166004820152602481018690529051919092169163a9059cbb9160448083019260209291908290030181600087803b158015610a9757600080fd5b505af1158015610aab573d6000803e3d6000fd5b505050506040513d6020811015610ac157600080fd5b5050601254600754600e546040517fa85ef579000000000000000000000000000000000000000000000000000000008152600160a060020a03898116600483019081523360248401819052604484018990529390940160848301819052600060a4840181905260c06064850190815260c485018b9052929096169563a85ef579958c95948a948d948d949093909290919060e401868680828437820191505098505050505050505050600060405180830381600087803b158015610b8457600080fd5b505af1158015610b98573d6000803e3d6000fd5b5050600380546001019055505060408051600160a060020a03871681526020810184905281517fd2819ba4c736158371edf0be38fd8d1fc435609832e392f118c4c79160e5bd7b929181900390910190a15050505050565b600054600160a060020a03163314610c0757600080fd5b600955565b60025481565b6012546040805160e060020a63a3078815028152600160a060020a0384811660048301529151600093929092169163a30788159160248082019260209290919082900301818787803b158015610c6757600080fd5b505af1158015610c7b573d6000803e3d6000fd5b505050506040513d6020811015610c9157600080fd5b505192915050565b600c54600160a060020a031681565b6012546040805160e160020a62c985b7028152600160a060020a0384811660048301529151600093849316916301930b6e91602480830192602092919082900301818787803b15801561080657600080fd5b606081801515610d0957506004545b610d15848260016120c8565b949350505050565b600b5481565b60075481565b6008546000908190600160a060020a03163314610d90576040805160e560020a62461bcd02815260206004820152600f60248201527f6d75737420626520736c61736865720000000000000000000000000000000000604482015290519081900360640190fd5b601254604080517fcad41357000000000000000000000000000000000000000000000000000000008152600160a060020a0388811660048301529151600293929092169163cad41357916024808201926020929091908290030181600087803b158015610dfc57600080fd5b505af1158015610e10573d6000803e3d6000fd5b505050506040513d6020811015610e2657600080fd5b5051811515610e3157fe5b049150600482601254604080517ffbc402fc000000000000000000000000000000000000000000000000000000008152600160a060020a038a81166004830152602482018890529151949093049450169163fbc402fc9160448082019260009290919082900301818387803b158015610ea957600080fd5b505af1158015610ebd573d6000803e3d6000fd5b50505050610eca85611b38565b15610fa257601254600754600e54604080517f3ac39d4b000000000000000000000000000000000000000000000000000000008152600160a060020a038a8116600483015292909301602484015251921691633ac39d4b9160448082019260009290919082900301818387803b158015610f4357600080fd5b505af1158015610f57573d6000803e3d6000fd5b505060038054600019019055505060408051600160a060020a038716815290517f2dc89de5703d2c341a22ebfc7c4d3f197e5e1f0c19bc2e1135f387163cb927e49181900360200190a15b6011546012546040805160e060020a63a3078815028152600160a060020a03888116600483015291519382169363a9059cbb939092169163a3078815916024808201926020929091908290030181600087803b15801561100157600080fd5b505af1158015611015573d6000803e3d6000fd5b505050506040513d602081101561102b57600080fd5b50516040805160e060020a63ffffffff8516028152600160a060020a039092166004830152602482018590525160448083019260209291908290030181600087803b15801561107957600080fd5b505af115801561108d573d6000803e3d6000fd5b505050506040513d60208110156110a357600080fd5b50506011546012546040805160e060020a63a3078815028152600160a060020a03878116600483015291519382169363a9059cbb939092169163a3078815916024808201926020929091908290030181600087803b15801561110457600080fd5b505af1158015611118573d6000803e3d6000fd5b505050506040513d602081101561112e57600080fd5b50516040805160e060020a63ffffffff8516028152600160a060020a039092166004830152602482018590525160448083019260209291908290030181600087803b15801561117c57600080fd5b505af1158015611190573d6000803e3d6000fd5b505050506040513d60208110156111a657600080fd5b50505050505050565b60006111ba8261131b565b801561124457506010546012546040805160e160020a62c985b7028152600160a060020a038681166004830152915191909216916301930b6e9160248083019260209291908290030181600087803b15801561121557600080fd5b505af1158015611229573d6000803e3d6000fd5b505050506040513d602081101561123f57600080fd5b505111155b92915050565b600f5460105482565b60095481565b600054600160a060020a0316331461127057600080fd5b600b55565b600a5481565b600054600160a060020a0316331461129257600080fd5b60008054604051600160a060020a03909116917ff8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c6482091a26000805473ffffffffffffffffffffffffffffffffffffffff19169055565b60045481565b600d54600e5482565b60408051808201909152600f5481526010546020820152600090611244908390612350565b6012546040805160e160020a62c985b7028152600160a060020a0384811660048301529151600093849316916301930b6e91602480830192602092919082900301818787803b15801561136d57600080fd5b505af1158015611381573d6000803e3d6000fd5b505050506040513d602081101561139757600080fd5b5051905080158015906108445750600e54101592915050565b600054600160a060020a031633146113c757600080fd5b600a55565b806113d681611b38565b151561142c576040805160e560020a62461bcd02815260206004820152601660248201527f6d757374206265206465726567697374657261626c6500000000000000000000604482015290519081900360640190fd5b6012546040805160e060020a63a3078815028152600160a060020a038086166004830152915185933393169163a30788159160248083019260209291908290030181600087803b15801561147f57600080fd5b505af1158015611493573d6000803e3d6000fd5b505050506040513d60208110156114a957600080fd5b5051600160a060020a031614611509576040805160e560020a62461bcd02815260206004820152601660248201527f6d757374206265206461726b6e6f6465206f776e657200000000000000000000604482015290519081900360640190fd5b601254600754600e54604080517f3ac39d4b000000000000000000000000000000000000000000000000000000008152600160a060020a03888116600483015292909301602484015251921691633ac39d4b9160448082019260009290919082900301818387803b15801561157d57600080fd5b505af1158015611591573d6000803e3d6000fd5b505060038054600019019055505060408051600160a060020a038516815290517f2dc89de5703d2c341a22ebfc7c4d3f197e5e1f0c19bc2e1135f387163cb927e49181900360200190a1505050565b601254604080517fee594a50000000000000000000000000000000000000000000000000000000008152600160a060020a0384811660048301529151606093929092169163ee594a509160248082019260009290919082900301818387803b15801561164b57600080fd5b505af115801561165f573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052602081101561168857600080fd5b8101908080516401000000008111156116a057600080fd5b820160208101848111156116b357600080fd5b81516401000000008111828201871017156116cd57600080fd5b50909695505050505050565b601154600160a060020a031681565b600054600160a060020a031681565b601054600090151561176557600054600160a060020a03163314611765576040805160e560020a62461bcd02815260206004820152601d60248201527f6e6f7420617574686f72697a6564202866697273742065706f63687329000000604482015290519081900360640190fd5b600754600e54014310156117c3576040805160e560020a62461bcd02815260206004820152601d60248201527f65706f636820696e74657276616c20686173206e6f7420706173736564000000604482015290519081900360640190fd5b50600d8054600f55600e805460105560408051808201909152436000198101408083526020909201819052928190559190556002805460045560035490556005546009541461185057600954600581905560408051828152602081019290925280517f7c77c94944e9e4e5b0d46f1297127d060020792687cd743401d782346c68f6559281900390910190a15b600654600a541461189f57600a54600681905560408051828152602081019290925280517f6d520e46e5714982ddf8cb6216bcb3e1c1d5b79d337afc305335f819394f5d6a9281900390910190a15b600754600b54146118ee57600b54600781905560408051828152602081019290925280517fb218cde2730b79a0667ddf869466ee66a12ef56fe65fa4986a590f8a7108c9de9281900390910190a15b600854600c54600160a060020a0390811691161461197357600c546008805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a039283169081179182905560408051929093168252602082015281517f933228a1c3ba8fadd3ce47a9db5b898be647f89af99ba7c1b9a655f59ea306c8929181900390910190a15b6040517feff7e281fe3b4211ed1f0a5e6419bcc40f4552974f771357e66926421f0a58e890600090a150565b601254600160a060020a031681565b60055481565b600854600160a060020a031681565b600054600160a060020a031633146119da57600080fd5b600c805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b601254604080517fcad41357000000000000000000000000000000000000000000000000000000008152600160a060020a0384811660048301529151600093929092169163cad413579160248082019260209290919082900301818787803b158015610c6757600080fd5b600054600160a060020a03163314611a8b57600080fd5b601254604080517ff2fde38b000000000000000000000000000000000000000000000000000000008152600160a060020a0384811660048301529151919092169163f2fde38b91602480830192600092919082900301818387803b158015611af257600080fd5b505af1158015611b06573d6000803e3d6000fd5b5050505050565b60408051808201909152600d548152600e546020820152600090611244908390612350565b60065481565b6012546040805160e160020a62c985b7028152600160a060020a0384811660048301529151600093849316916301930b6e91602480830192602092919082900301818787803b158015611b8a57600080fd5b505af1158015611b9e573d6000803e3d6000fd5b505050506040513d6020811015611bb457600080fd5b50519050611bc183611b0d565b801561084457501592915050565b606081801515611bde57506002545b610d15848260006120c8565b600054600160a060020a03163314611c0157600080fd5b611c0a816124ac565b50565b60008082611c1a816111af565b1515611c96576040805160e560020a62461bcd02815260206004820152602b60248201527f6d7573742062652064657265676973746572656420666f72206174206c65617360448201527f74206f6e652065706f6368000000000000000000000000000000000000000000606482015290519081900360840190fd5b6012546040805160e060020a63a3078815028152600160a060020a0387811660048301529151919092169163a30788159160248083019260209291908290030181600087803b158015611ce857600080fd5b505af1158015611cfc573d6000803e3d6000fd5b505050506040513d6020811015611d1257600080fd5b5051601254604080517fcad41357000000000000000000000000000000000000000000000000000000008152600160a060020a038881166004830152915193965091169163cad41357916024808201926020929091908290030181600087803b158015611d7e57600080fd5b505af1158015611d92573d6000803e3d6000fd5b505050506040513d6020811015611da857600080fd5b5051601254604080517f41b44392000000000000000000000000000000000000000000000000000000008152600160a060020a03888116600483015291519395509116916341b443929160248082019260009290919082900301818387803b158015611e1357600080fd5b505af1158015611e27573d6000803e3d6000fd5b5050601154604080517fa9059cbb000000000000000000000000000000000000000000000000000000008152600160a060020a03888116600483015260248201889052915191909216935063a9059cbb925060448083019260209291908290030181600087803b158015611e9a57600080fd5b505af1158015611eae573d6000803e3d6000fd5b505050506040513d6020811015611ec457600080fd5b505060408051600160a060020a03851681526020810184905281517f96ab9e56c79eee4a72db6e2879cbfbecdba5c65b411f4861824e66b89df19764929181900390910190a150505050565b60018054604080516020600284861615610100026000190190941693909304601f81018490048402820184019092528181529291830182828015611f955780601f10611f6a57610100808354040283529160200191611f95565b820191906000526020600020905b815481529060010190602001808311611f7857829003601f168201915b505050505081565b601254604080517fe2d7c64c000000000000000000000000000000000000000000000000000000008152600160a060020a03848116600483015291516000938493849391169163e2d7c64c9160248082019260209290919082900301818787803b15801561200a57600080fd5b505af115801561201e573d6000803e3d6000fd5b505050506040513d602081101561203457600080fd5b50516012546040805160e160020a62c985b7028152600160a060020a03888116600483015291519395509116916301930b6e916024808201926020929091908290030181600087803b15801561208957600080fd5b505af115801561209d573d6000803e3d6000fd5b505050506040513d60208110156120b357600080fd5b5051905081158015610d155750159392505050565b60608281600080808415156120dd5760025494505b84604051908082528060200260200182016040528015612107578160200160208202803883390190505b50935060009250889150600160a060020a03821615156121a357601260009054906101000a9004600160a060020a0316600160a060020a0316631bce6ff36040518163ffffffff1660e060020a028152600401602060405180830381600087803b15801561217457600080fd5b505af1158015612188573d6000803e3d6000fd5b505050506040513d602081101561219e57600080fd5b505191505b8483101561234357600160a060020a03821615156121c057612343565b86156121d6576121cf826112f6565b90506121e2565b6121df82611b0d565b90505b80151561228457601254604080517fab73e316000000000000000000000000000000000000000000000000000000008152600160a060020a0385811660048301529151919092169163ab73e3169160248083019260209291908290030181600087803b15801561225157600080fd5b505af1158015612265573d6000803e3d6000fd5b505050506040513d602081101561227b57600080fd5b505191506121a3565b81848481518110151561229357fe5b600160a060020a039283166020918202909201810191909152601254604080517fab73e31600000000000000000000000000000000000000000000000000000000815286851660048201529051919093169263ab73e3169260248083019391928290030181600087803b15801561230957600080fd5b505af115801561231d573d6000803e3d6000fd5b505050506040513d602081101561233357600080fd5b50516001939093019291506121a3565b5091979650505050505050565b601254604080517fe2d7c64c000000000000000000000000000000000000000000000000000000008152600160a060020a03858116600483015291516000938493849384938493929092169163e2d7c64c9160248082019260209290919082900301818787803b1580156123c357600080fd5b505af11580156123d7573d6000803e3d6000fd5b505050506040513d60208110156123ed57600080fd5b50516012546040805160e160020a62c985b7028152600160a060020a038b8116600483015291519397509116916301930b6e916024808201926020929091908290030181600087803b15801561244257600080fd5b505af1158015612456573d6000803e3d6000fd5b505050506040513d602081101561246c57600080fd5b505192508315801590612483575085602001518411155b91508215806124955750856020015183115b90508180156124a15750805b979650505050505050565b600160a060020a03811615156124c157600080fd5b60008054604051600160a060020a03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a36000805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03929092169190911790555600a165627a7a723058203ca0ffa3522020513e86bfc48d54dd91c755b8938e2a045ce0b50270ec95eba10029`

// DeployDarknodeRegistry deploys a new Ethereum contract, binding an instance of DarknodeRegistry to it.
func DeployDarknodeRegistry(auth *bind.TransactOpts, backend bind.ContractBackend, _VERSION string, _renAddress common.Address, _storeAddress common.Address, _minimumBond *big.Int, _minimumPodSize *big.Int, _minimumEpochInterval *big.Int) (common.Address, *types.Transaction, *DarknodeRegistry, error) {
	parsed, err := abi.JSON(strings.NewReader(DarknodeRegistryABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(DarknodeRegistryBin), backend, _VERSION, _renAddress, _storeAddress, _minimumBond, _minimumPodSize, _minimumEpochInterval)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DarknodeRegistry{DarknodeRegistryCaller: DarknodeRegistryCaller{contract: contract}, DarknodeRegistryTransactor: DarknodeRegistryTransactor{contract: contract}, DarknodeRegistryFilterer: DarknodeRegistryFilterer{contract: contract}}, nil
}

// DarknodeRegistry is an auto generated Go binding around an Ethereum contract.
type DarknodeRegistry struct {
	DarknodeRegistryCaller     // Read-only binding to the contract
	DarknodeRegistryTransactor // Write-only binding to the contract
	DarknodeRegistryFilterer   // Log filterer for contract events
}

// DarknodeRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type DarknodeRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DarknodeRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DarknodeRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DarknodeRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DarknodeRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DarknodeRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DarknodeRegistrySession struct {
	Contract     *DarknodeRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DarknodeRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DarknodeRegistryCallerSession struct {
	Contract *DarknodeRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// DarknodeRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DarknodeRegistryTransactorSession struct {
	Contract     *DarknodeRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// DarknodeRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type DarknodeRegistryRaw struct {
	Contract *DarknodeRegistry // Generic contract binding to access the raw methods on
}

// DarknodeRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DarknodeRegistryCallerRaw struct {
	Contract *DarknodeRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// DarknodeRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DarknodeRegistryTransactorRaw struct {
	Contract *DarknodeRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDarknodeRegistry creates a new instance of DarknodeRegistry, bound to a specific deployed contract.
func NewDarknodeRegistry(address common.Address, backend bind.ContractBackend) (*DarknodeRegistry, error) {
	contract, err := bindDarknodeRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DarknodeRegistry{DarknodeRegistryCaller: DarknodeRegistryCaller{contract: contract}, DarknodeRegistryTransactor: DarknodeRegistryTransactor{contract: contract}, DarknodeRegistryFilterer: DarknodeRegistryFilterer{contract: contract}}, nil
}

// NewDarknodeRegistryCaller creates a new read-only instance of DarknodeRegistry, bound to a specific deployed contract.
func NewDarknodeRegistryCaller(address common.Address, caller bind.ContractCaller) (*DarknodeRegistryCaller, error) {
	contract, err := bindDarknodeRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DarknodeRegistryCaller{contract: contract}, nil
}

// NewDarknodeRegistryTransactor creates a new write-only instance of DarknodeRegistry, bound to a specific deployed contract.
func NewDarknodeRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*DarknodeRegistryTransactor, error) {
	contract, err := bindDarknodeRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DarknodeRegistryTransactor{contract: contract}, nil
}

// NewDarknodeRegistryFilterer creates a new log filterer instance of DarknodeRegistry, bound to a specific deployed contract.
func NewDarknodeRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*DarknodeRegistryFilterer, error) {
	contract, err := bindDarknodeRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DarknodeRegistryFilterer{contract: contract}, nil
}

// bindDarknodeRegistry binds a generic wrapper to an already deployed contract.
func bindDarknodeRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DarknodeRegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DarknodeRegistry *DarknodeRegistryRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _DarknodeRegistry.Contract.DarknodeRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DarknodeRegistry *DarknodeRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.DarknodeRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DarknodeRegistry *DarknodeRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.DarknodeRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DarknodeRegistry *DarknodeRegistryCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _DarknodeRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DarknodeRegistry *DarknodeRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DarknodeRegistry *DarknodeRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.contract.Transact(opts, method, params...)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_DarknodeRegistry *DarknodeRegistryCaller) VERSION(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "VERSION")
	return *ret0, err
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_DarknodeRegistry *DarknodeRegistrySession) VERSION() (string, error) {
	return _DarknodeRegistry.Contract.VERSION(&_DarknodeRegistry.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) VERSION() (string, error) {
	return _DarknodeRegistry.Contract.VERSION(&_DarknodeRegistry.CallOpts)
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() constant returns(epochhash uint256, blocknumber uint256)
func (_DarknodeRegistry *DarknodeRegistryCaller) CurrentEpoch(opts *bind.CallOpts) (struct {
	Epochhash   *big.Int
	Blocknumber *big.Int
}, error) {
	ret := new(struct {
		Epochhash   *big.Int
		Blocknumber *big.Int
	})
	out := ret
	err := _DarknodeRegistry.contract.Call(opts, out, "currentEpoch")
	return *ret, err
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() constant returns(epochhash uint256, blocknumber uint256)
func (_DarknodeRegistry *DarknodeRegistrySession) CurrentEpoch() (struct {
	Epochhash   *big.Int
	Blocknumber *big.Int
}, error) {
	return _DarknodeRegistry.Contract.CurrentEpoch(&_DarknodeRegistry.CallOpts)
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() constant returns(epochhash uint256, blocknumber uint256)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) CurrentEpoch() (struct {
	Epochhash   *big.Int
	Blocknumber *big.Int
}, error) {
	return _DarknodeRegistry.Contract.CurrentEpoch(&_DarknodeRegistry.CallOpts)
}

// GetDarknodeBond is a free data retrieval call binding the contract method 0xba0f5b20.
//
// Solidity: function getDarknodeBond(_darknodeID address) constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistryCaller) GetDarknodeBond(opts *bind.CallOpts, _darknodeID common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "getDarknodeBond", _darknodeID)
	return *ret0, err
}

// GetDarknodeBond is a free data retrieval call binding the contract method 0xba0f5b20.
//
// Solidity: function getDarknodeBond(_darknodeID address) constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistrySession) GetDarknodeBond(_darknodeID common.Address) (*big.Int, error) {
	return _DarknodeRegistry.Contract.GetDarknodeBond(&_DarknodeRegistry.CallOpts, _darknodeID)
}

// GetDarknodeBond is a free data retrieval call binding the contract method 0xba0f5b20.
//
// Solidity: function getDarknodeBond(_darknodeID address) constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) GetDarknodeBond(_darknodeID common.Address) (*big.Int, error) {
	return _DarknodeRegistry.Contract.GetDarknodeBond(&_DarknodeRegistry.CallOpts, _darknodeID)
}

// GetDarknodeOwner is a free data retrieval call binding the contract method 0x1cedf8a3.
//
// Solidity: function getDarknodeOwner(_darknodeID address) constant returns(address)
func (_DarknodeRegistry *DarknodeRegistryCaller) GetDarknodeOwner(opts *bind.CallOpts, _darknodeID common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "getDarknodeOwner", _darknodeID)
	return *ret0, err
}

// GetDarknodeOwner is a free data retrieval call binding the contract method 0x1cedf8a3.
//
// Solidity: function getDarknodeOwner(_darknodeID address) constant returns(address)
func (_DarknodeRegistry *DarknodeRegistrySession) GetDarknodeOwner(_darknodeID common.Address) (common.Address, error) {
	return _DarknodeRegistry.Contract.GetDarknodeOwner(&_DarknodeRegistry.CallOpts, _darknodeID)
}

// GetDarknodeOwner is a free data retrieval call binding the contract method 0x1cedf8a3.
//
// Solidity: function getDarknodeOwner(_darknodeID address) constant returns(address)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) GetDarknodeOwner(_darknodeID common.Address) (common.Address, error) {
	return _DarknodeRegistry.Contract.GetDarknodeOwner(&_DarknodeRegistry.CallOpts, _darknodeID)
}

// GetDarknodePublicKey is a free data retrieval call binding the contract method 0x84d2688c.
//
// Solidity: function getDarknodePublicKey(_darknodeID address) constant returns(bytes)
func (_DarknodeRegistry *DarknodeRegistryCaller) GetDarknodePublicKey(opts *bind.CallOpts, _darknodeID common.Address) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "getDarknodePublicKey", _darknodeID)
	return *ret0, err
}

// GetDarknodePublicKey is a free data retrieval call binding the contract method 0x84d2688c.
//
// Solidity: function getDarknodePublicKey(_darknodeID address) constant returns(bytes)
func (_DarknodeRegistry *DarknodeRegistrySession) GetDarknodePublicKey(_darknodeID common.Address) ([]byte, error) {
	return _DarknodeRegistry.Contract.GetDarknodePublicKey(&_DarknodeRegistry.CallOpts, _darknodeID)
}

// GetDarknodePublicKey is a free data retrieval call binding the contract method 0x84d2688c.
//
// Solidity: function getDarknodePublicKey(_darknodeID address) constant returns(bytes)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) GetDarknodePublicKey(_darknodeID common.Address) ([]byte, error) {
	return _DarknodeRegistry.Contract.GetDarknodePublicKey(&_DarknodeRegistry.CallOpts, _darknodeID)
}

// GetDarknodes is a free data retrieval call binding the contract method 0xec5325c1.
//
// Solidity: function getDarknodes(_start address, _count uint256) constant returns(address[])
func (_DarknodeRegistry *DarknodeRegistryCaller) GetDarknodes(opts *bind.CallOpts, _start common.Address, _count *big.Int) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "getDarknodes", _start, _count)
	return *ret0, err
}

// GetDarknodes is a free data retrieval call binding the contract method 0xec5325c1.
//
// Solidity: function getDarknodes(_start address, _count uint256) constant returns(address[])
func (_DarknodeRegistry *DarknodeRegistrySession) GetDarknodes(_start common.Address, _count *big.Int) ([]common.Address, error) {
	return _DarknodeRegistry.Contract.GetDarknodes(&_DarknodeRegistry.CallOpts, _start, _count)
}

// GetDarknodes is a free data retrieval call binding the contract method 0xec5325c1.
//
// Solidity: function getDarknodes(_start address, _count uint256) constant returns(address[])
func (_DarknodeRegistry *DarknodeRegistryCallerSession) GetDarknodes(_start common.Address, _count *big.Int) ([]common.Address, error) {
	return _DarknodeRegistry.Contract.GetDarknodes(&_DarknodeRegistry.CallOpts, _start, _count)
}

// GetPreviousDarknodes is a free data retrieval call binding the contract method 0x43846074.
//
// Solidity: function getPreviousDarknodes(_start address, _count uint256) constant returns(address[])
func (_DarknodeRegistry *DarknodeRegistryCaller) GetPreviousDarknodes(opts *bind.CallOpts, _start common.Address, _count *big.Int) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "getPreviousDarknodes", _start, _count)
	return *ret0, err
}

// GetPreviousDarknodes is a free data retrieval call binding the contract method 0x43846074.
//
// Solidity: function getPreviousDarknodes(_start address, _count uint256) constant returns(address[])
func (_DarknodeRegistry *DarknodeRegistrySession) GetPreviousDarknodes(_start common.Address, _count *big.Int) ([]common.Address, error) {
	return _DarknodeRegistry.Contract.GetPreviousDarknodes(&_DarknodeRegistry.CallOpts, _start, _count)
}

// GetPreviousDarknodes is a free data retrieval call binding the contract method 0x43846074.
//
// Solidity: function getPreviousDarknodes(_start address, _count uint256) constant returns(address[])
func (_DarknodeRegistry *DarknodeRegistryCallerSession) GetPreviousDarknodes(_start common.Address, _count *big.Int) ([]common.Address, error) {
	return _DarknodeRegistry.Contract.GetPreviousDarknodes(&_DarknodeRegistry.CallOpts, _start, _count)
}

// IsDeregisterable is a free data retrieval call binding the contract method 0xe1878925.
//
// Solidity: function isDeregisterable(_darknodeID address) constant returns(bool)
func (_DarknodeRegistry *DarknodeRegistryCaller) IsDeregisterable(opts *bind.CallOpts, _darknodeID common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "isDeregisterable", _darknodeID)
	return *ret0, err
}

// IsDeregisterable is a free data retrieval call binding the contract method 0xe1878925.
//
// Solidity: function isDeregisterable(_darknodeID address) constant returns(bool)
func (_DarknodeRegistry *DarknodeRegistrySession) IsDeregisterable(_darknodeID common.Address) (bool, error) {
	return _DarknodeRegistry.Contract.IsDeregisterable(&_DarknodeRegistry.CallOpts, _darknodeID)
}

// IsDeregisterable is a free data retrieval call binding the contract method 0xe1878925.
//
// Solidity: function isDeregisterable(_darknodeID address) constant returns(bool)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) IsDeregisterable(_darknodeID common.Address) (bool, error) {
	return _DarknodeRegistry.Contract.IsDeregisterable(&_DarknodeRegistry.CallOpts, _darknodeID)
}

// IsDeregistered is a free data retrieval call binding the contract method 0x8020fc1f.
//
// Solidity: function isDeregistered(_darknodeID address) constant returns(bool)
func (_DarknodeRegistry *DarknodeRegistryCaller) IsDeregistered(opts *bind.CallOpts, _darknodeID common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "isDeregistered", _darknodeID)
	return *ret0, err
}

// IsDeregistered is a free data retrieval call binding the contract method 0x8020fc1f.
//
// Solidity: function isDeregistered(_darknodeID address) constant returns(bool)
func (_DarknodeRegistry *DarknodeRegistrySession) IsDeregistered(_darknodeID common.Address) (bool, error) {
	return _DarknodeRegistry.Contract.IsDeregistered(&_DarknodeRegistry.CallOpts, _darknodeID)
}

// IsDeregistered is a free data retrieval call binding the contract method 0x8020fc1f.
//
// Solidity: function isDeregistered(_darknodeID address) constant returns(bool)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) IsDeregistered(_darknodeID common.Address) (bool, error) {
	return _DarknodeRegistry.Contract.IsDeregistered(&_DarknodeRegistry.CallOpts, _darknodeID)
}

// IsPendingDeregistration is a free data retrieval call binding the contract method 0x303ee989.
//
// Solidity: function isPendingDeregistration(_darknodeID address) constant returns(bool)
func (_DarknodeRegistry *DarknodeRegistryCaller) IsPendingDeregistration(opts *bind.CallOpts, _darknodeID common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "isPendingDeregistration", _darknodeID)
	return *ret0, err
}

// IsPendingDeregistration is a free data retrieval call binding the contract method 0x303ee989.
//
// Solidity: function isPendingDeregistration(_darknodeID address) constant returns(bool)
func (_DarknodeRegistry *DarknodeRegistrySession) IsPendingDeregistration(_darknodeID common.Address) (bool, error) {
	return _DarknodeRegistry.Contract.IsPendingDeregistration(&_DarknodeRegistry.CallOpts, _darknodeID)
}

// IsPendingDeregistration is a free data retrieval call binding the contract method 0x303ee989.
//
// Solidity: function isPendingDeregistration(_darknodeID address) constant returns(bool)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) IsPendingDeregistration(_darknodeID common.Address) (bool, error) {
	return _DarknodeRegistry.Contract.IsPendingDeregistration(&_DarknodeRegistry.CallOpts, _darknodeID)
}

// IsPendingRegistration is a free data retrieval call binding the contract method 0x040fa051.
//
// Solidity: function isPendingRegistration(_darknodeID address) constant returns(bool)
func (_DarknodeRegistry *DarknodeRegistryCaller) IsPendingRegistration(opts *bind.CallOpts, _darknodeID common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "isPendingRegistration", _darknodeID)
	return *ret0, err
}

// IsPendingRegistration is a free data retrieval call binding the contract method 0x040fa051.
//
// Solidity: function isPendingRegistration(_darknodeID address) constant returns(bool)
func (_DarknodeRegistry *DarknodeRegistrySession) IsPendingRegistration(_darknodeID common.Address) (bool, error) {
	return _DarknodeRegistry.Contract.IsPendingRegistration(&_DarknodeRegistry.CallOpts, _darknodeID)
}

// IsPendingRegistration is a free data retrieval call binding the contract method 0x040fa051.
//
// Solidity: function isPendingRegistration(_darknodeID address) constant returns(bool)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) IsPendingRegistration(_darknodeID common.Address) (bool, error) {
	return _DarknodeRegistry.Contract.IsPendingRegistration(&_DarknodeRegistry.CallOpts, _darknodeID)
}

// IsRefundable is a free data retrieval call binding the contract method 0x5aebd1cb.
//
// Solidity: function isRefundable(_darknodeID address) constant returns(bool)
func (_DarknodeRegistry *DarknodeRegistryCaller) IsRefundable(opts *bind.CallOpts, _darknodeID common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "isRefundable", _darknodeID)
	return *ret0, err
}

// IsRefundable is a free data retrieval call binding the contract method 0x5aebd1cb.
//
// Solidity: function isRefundable(_darknodeID address) constant returns(bool)
func (_DarknodeRegistry *DarknodeRegistrySession) IsRefundable(_darknodeID common.Address) (bool, error) {
	return _DarknodeRegistry.Contract.IsRefundable(&_DarknodeRegistry.CallOpts, _darknodeID)
}

// IsRefundable is a free data retrieval call binding the contract method 0x5aebd1cb.
//
// Solidity: function isRefundable(_darknodeID address) constant returns(bool)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) IsRefundable(_darknodeID common.Address) (bool, error) {
	return _DarknodeRegistry.Contract.IsRefundable(&_DarknodeRegistry.CallOpts, _darknodeID)
}

// IsRefunded is a free data retrieval call binding the contract method 0xffc9152e.
//
// Solidity: function isRefunded(_darknodeID address) constant returns(bool)
func (_DarknodeRegistry *DarknodeRegistryCaller) IsRefunded(opts *bind.CallOpts, _darknodeID common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "isRefunded", _darknodeID)
	return *ret0, err
}

// IsRefunded is a free data retrieval call binding the contract method 0xffc9152e.
//
// Solidity: function isRefunded(_darknodeID address) constant returns(bool)
func (_DarknodeRegistry *DarknodeRegistrySession) IsRefunded(_darknodeID common.Address) (bool, error) {
	return _DarknodeRegistry.Contract.IsRefunded(&_DarknodeRegistry.CallOpts, _darknodeID)
}

// IsRefunded is a free data retrieval call binding the contract method 0xffc9152e.
//
// Solidity: function isRefunded(_darknodeID address) constant returns(bool)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) IsRefunded(_darknodeID common.Address) (bool, error) {
	return _DarknodeRegistry.Contract.IsRefunded(&_DarknodeRegistry.CallOpts, _darknodeID)
}

// IsRegistered is a free data retrieval call binding the contract method 0xc3c5a547.
//
// Solidity: function isRegistered(_darknodeID address) constant returns(bool)
func (_DarknodeRegistry *DarknodeRegistryCaller) IsRegistered(opts *bind.CallOpts, _darknodeID common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "isRegistered", _darknodeID)
	return *ret0, err
}

// IsRegistered is a free data retrieval call binding the contract method 0xc3c5a547.
//
// Solidity: function isRegistered(_darknodeID address) constant returns(bool)
func (_DarknodeRegistry *DarknodeRegistrySession) IsRegistered(_darknodeID common.Address) (bool, error) {
	return _DarknodeRegistry.Contract.IsRegistered(&_DarknodeRegistry.CallOpts, _darknodeID)
}

// IsRegistered is a free data retrieval call binding the contract method 0xc3c5a547.
//
// Solidity: function isRegistered(_darknodeID address) constant returns(bool)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) IsRegistered(_darknodeID common.Address) (bool, error) {
	return _DarknodeRegistry.Contract.IsRegistered(&_DarknodeRegistry.CallOpts, _darknodeID)
}

// IsRegisteredInPreviousEpoch is a free data retrieval call binding the contract method 0x7be266da.
//
// Solidity: function isRegisteredInPreviousEpoch(_darknodeID address) constant returns(bool)
func (_DarknodeRegistry *DarknodeRegistryCaller) IsRegisteredInPreviousEpoch(opts *bind.CallOpts, _darknodeID common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "isRegisteredInPreviousEpoch", _darknodeID)
	return *ret0, err
}

// IsRegisteredInPreviousEpoch is a free data retrieval call binding the contract method 0x7be266da.
//
// Solidity: function isRegisteredInPreviousEpoch(_darknodeID address) constant returns(bool)
func (_DarknodeRegistry *DarknodeRegistrySession) IsRegisteredInPreviousEpoch(_darknodeID common.Address) (bool, error) {
	return _DarknodeRegistry.Contract.IsRegisteredInPreviousEpoch(&_DarknodeRegistry.CallOpts, _darknodeID)
}

// IsRegisteredInPreviousEpoch is a free data retrieval call binding the contract method 0x7be266da.
//
// Solidity: function isRegisteredInPreviousEpoch(_darknodeID address) constant returns(bool)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) IsRegisteredInPreviousEpoch(_darknodeID common.Address) (bool, error) {
	return _DarknodeRegistry.Contract.IsRegisteredInPreviousEpoch(&_DarknodeRegistry.CallOpts, _darknodeID)
}

// MinimumBond is a free data retrieval call binding the contract method 0xaa7517e1.
//
// Solidity: function minimumBond() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistryCaller) MinimumBond(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "minimumBond")
	return *ret0, err
}

// MinimumBond is a free data retrieval call binding the contract method 0xaa7517e1.
//
// Solidity: function minimumBond() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistrySession) MinimumBond() (*big.Int, error) {
	return _DarknodeRegistry.Contract.MinimumBond(&_DarknodeRegistry.CallOpts)
}

// MinimumBond is a free data retrieval call binding the contract method 0xaa7517e1.
//
// Solidity: function minimumBond() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) MinimumBond() (*big.Int, error) {
	return _DarknodeRegistry.Contract.MinimumBond(&_DarknodeRegistry.CallOpts)
}

// MinimumEpochInterval is a free data retrieval call binding the contract method 0x55cacda5.
//
// Solidity: function minimumEpochInterval() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistryCaller) MinimumEpochInterval(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "minimumEpochInterval")
	return *ret0, err
}

// MinimumEpochInterval is a free data retrieval call binding the contract method 0x55cacda5.
//
// Solidity: function minimumEpochInterval() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistrySession) MinimumEpochInterval() (*big.Int, error) {
	return _DarknodeRegistry.Contract.MinimumEpochInterval(&_DarknodeRegistry.CallOpts)
}

// MinimumEpochInterval is a free data retrieval call binding the contract method 0x55cacda5.
//
// Solidity: function minimumEpochInterval() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) MinimumEpochInterval() (*big.Int, error) {
	return _DarknodeRegistry.Contract.MinimumEpochInterval(&_DarknodeRegistry.CallOpts)
}

// MinimumPodSize is a free data retrieval call binding the contract method 0xc7dbc2be.
//
// Solidity: function minimumPodSize() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistryCaller) MinimumPodSize(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "minimumPodSize")
	return *ret0, err
}

// MinimumPodSize is a free data retrieval call binding the contract method 0xc7dbc2be.
//
// Solidity: function minimumPodSize() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistrySession) MinimumPodSize() (*big.Int, error) {
	return _DarknodeRegistry.Contract.MinimumPodSize(&_DarknodeRegistry.CallOpts)
}

// MinimumPodSize is a free data retrieval call binding the contract method 0xc7dbc2be.
//
// Solidity: function minimumPodSize() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) MinimumPodSize() (*big.Int, error) {
	return _DarknodeRegistry.Contract.MinimumPodSize(&_DarknodeRegistry.CallOpts)
}

// NextMinimumBond is a free data retrieval call binding the contract method 0x60a22fe4.
//
// Solidity: function nextMinimumBond() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistryCaller) NextMinimumBond(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "nextMinimumBond")
	return *ret0, err
}

// NextMinimumBond is a free data retrieval call binding the contract method 0x60a22fe4.
//
// Solidity: function nextMinimumBond() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistrySession) NextMinimumBond() (*big.Int, error) {
	return _DarknodeRegistry.Contract.NextMinimumBond(&_DarknodeRegistry.CallOpts)
}

// NextMinimumBond is a free data retrieval call binding the contract method 0x60a22fe4.
//
// Solidity: function nextMinimumBond() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) NextMinimumBond() (*big.Int, error) {
	return _DarknodeRegistry.Contract.NextMinimumBond(&_DarknodeRegistry.CallOpts)
}

// NextMinimumEpochInterval is a free data retrieval call binding the contract method 0x455dc46d.
//
// Solidity: function nextMinimumEpochInterval() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistryCaller) NextMinimumEpochInterval(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "nextMinimumEpochInterval")
	return *ret0, err
}

// NextMinimumEpochInterval is a free data retrieval call binding the contract method 0x455dc46d.
//
// Solidity: function nextMinimumEpochInterval() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistrySession) NextMinimumEpochInterval() (*big.Int, error) {
	return _DarknodeRegistry.Contract.NextMinimumEpochInterval(&_DarknodeRegistry.CallOpts)
}

// NextMinimumEpochInterval is a free data retrieval call binding the contract method 0x455dc46d.
//
// Solidity: function nextMinimumEpochInterval() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) NextMinimumEpochInterval() (*big.Int, error) {
	return _DarknodeRegistry.Contract.NextMinimumEpochInterval(&_DarknodeRegistry.CallOpts)
}

// NextMinimumPodSize is a free data retrieval call binding the contract method 0x702c25ee.
//
// Solidity: function nextMinimumPodSize() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistryCaller) NextMinimumPodSize(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "nextMinimumPodSize")
	return *ret0, err
}

// NextMinimumPodSize is a free data retrieval call binding the contract method 0x702c25ee.
//
// Solidity: function nextMinimumPodSize() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistrySession) NextMinimumPodSize() (*big.Int, error) {
	return _DarknodeRegistry.Contract.NextMinimumPodSize(&_DarknodeRegistry.CallOpts)
}

// NextMinimumPodSize is a free data retrieval call binding the contract method 0x702c25ee.
//
// Solidity: function nextMinimumPodSize() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) NextMinimumPodSize() (*big.Int, error) {
	return _DarknodeRegistry.Contract.NextMinimumPodSize(&_DarknodeRegistry.CallOpts)
}

// NextSlasher is a free data retrieval call binding the contract method 0x21a2ad3a.
//
// Solidity: function nextSlasher() constant returns(address)
func (_DarknodeRegistry *DarknodeRegistryCaller) NextSlasher(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "nextSlasher")
	return *ret0, err
}

// NextSlasher is a free data retrieval call binding the contract method 0x21a2ad3a.
//
// Solidity: function nextSlasher() constant returns(address)
func (_DarknodeRegistry *DarknodeRegistrySession) NextSlasher() (common.Address, error) {
	return _DarknodeRegistry.Contract.NextSlasher(&_DarknodeRegistry.CallOpts)
}

// NextSlasher is a free data retrieval call binding the contract method 0x21a2ad3a.
//
// Solidity: function nextSlasher() constant returns(address)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) NextSlasher() (common.Address, error) {
	return _DarknodeRegistry.Contract.NextSlasher(&_DarknodeRegistry.CallOpts)
}

// NumDarknodes is a free data retrieval call binding the contract method 0x1460e603.
//
// Solidity: function numDarknodes() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistryCaller) NumDarknodes(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "numDarknodes")
	return *ret0, err
}

// NumDarknodes is a free data retrieval call binding the contract method 0x1460e603.
//
// Solidity: function numDarknodes() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistrySession) NumDarknodes() (*big.Int, error) {
	return _DarknodeRegistry.Contract.NumDarknodes(&_DarknodeRegistry.CallOpts)
}

// NumDarknodes is a free data retrieval call binding the contract method 0x1460e603.
//
// Solidity: function numDarknodes() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) NumDarknodes() (*big.Int, error) {
	return _DarknodeRegistry.Contract.NumDarknodes(&_DarknodeRegistry.CallOpts)
}

// NumDarknodesNextEpoch is a free data retrieval call binding the contract method 0x0847e9fa.
//
// Solidity: function numDarknodesNextEpoch() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistryCaller) NumDarknodesNextEpoch(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "numDarknodesNextEpoch")
	return *ret0, err
}

// NumDarknodesNextEpoch is a free data retrieval call binding the contract method 0x0847e9fa.
//
// Solidity: function numDarknodesNextEpoch() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistrySession) NumDarknodesNextEpoch() (*big.Int, error) {
	return _DarknodeRegistry.Contract.NumDarknodesNextEpoch(&_DarknodeRegistry.CallOpts)
}

// NumDarknodesNextEpoch is a free data retrieval call binding the contract method 0x0847e9fa.
//
// Solidity: function numDarknodesNextEpoch() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) NumDarknodesNextEpoch() (*big.Int, error) {
	return _DarknodeRegistry.Contract.NumDarknodesNextEpoch(&_DarknodeRegistry.CallOpts)
}

// NumDarknodesPreviousEpoch is a free data retrieval call binding the contract method 0x71740d16.
//
// Solidity: function numDarknodesPreviousEpoch() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistryCaller) NumDarknodesPreviousEpoch(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "numDarknodesPreviousEpoch")
	return *ret0, err
}

// NumDarknodesPreviousEpoch is a free data retrieval call binding the contract method 0x71740d16.
//
// Solidity: function numDarknodesPreviousEpoch() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistrySession) NumDarknodesPreviousEpoch() (*big.Int, error) {
	return _DarknodeRegistry.Contract.NumDarknodesPreviousEpoch(&_DarknodeRegistry.CallOpts)
}

// NumDarknodesPreviousEpoch is a free data retrieval call binding the contract method 0x71740d16.
//
// Solidity: function numDarknodesPreviousEpoch() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) NumDarknodesPreviousEpoch() (*big.Int, error) {
	return _DarknodeRegistry.Contract.NumDarknodesPreviousEpoch(&_DarknodeRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_DarknodeRegistry *DarknodeRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_DarknodeRegistry *DarknodeRegistrySession) Owner() (common.Address, error) {
	return _DarknodeRegistry.Contract.Owner(&_DarknodeRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) Owner() (common.Address, error) {
	return _DarknodeRegistry.Contract.Owner(&_DarknodeRegistry.CallOpts)
}

// PreviousEpoch is a free data retrieval call binding the contract method 0x5cdaab48.
//
// Solidity: function previousEpoch() constant returns(epochhash uint256, blocknumber uint256)
func (_DarknodeRegistry *DarknodeRegistryCaller) PreviousEpoch(opts *bind.CallOpts) (struct {
	Epochhash   *big.Int
	Blocknumber *big.Int
}, error) {
	ret := new(struct {
		Epochhash   *big.Int
		Blocknumber *big.Int
	})
	out := ret
	err := _DarknodeRegistry.contract.Call(opts, out, "previousEpoch")
	return *ret, err
}

// PreviousEpoch is a free data retrieval call binding the contract method 0x5cdaab48.
//
// Solidity: function previousEpoch() constant returns(epochhash uint256, blocknumber uint256)
func (_DarknodeRegistry *DarknodeRegistrySession) PreviousEpoch() (struct {
	Epochhash   *big.Int
	Blocknumber *big.Int
}, error) {
	return _DarknodeRegistry.Contract.PreviousEpoch(&_DarknodeRegistry.CallOpts)
}

// PreviousEpoch is a free data retrieval call binding the contract method 0x5cdaab48.
//
// Solidity: function previousEpoch() constant returns(epochhash uint256, blocknumber uint256)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) PreviousEpoch() (struct {
	Epochhash   *big.Int
	Blocknumber *big.Int
}, error) {
	return _DarknodeRegistry.Contract.PreviousEpoch(&_DarknodeRegistry.CallOpts)
}

// Ren is a free data retrieval call binding the contract method 0x8a9b4067.
//
// Solidity: function ren() constant returns(address)
func (_DarknodeRegistry *DarknodeRegistryCaller) Ren(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "ren")
	return *ret0, err
}

// Ren is a free data retrieval call binding the contract method 0x8a9b4067.
//
// Solidity: function ren() constant returns(address)
func (_DarknodeRegistry *DarknodeRegistrySession) Ren() (common.Address, error) {
	return _DarknodeRegistry.Contract.Ren(&_DarknodeRegistry.CallOpts)
}

// Ren is a free data retrieval call binding the contract method 0x8a9b4067.
//
// Solidity: function ren() constant returns(address)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) Ren() (common.Address, error) {
	return _DarknodeRegistry.Contract.Ren(&_DarknodeRegistry.CallOpts)
}

// Slasher is a free data retrieval call binding the contract method 0xb1344271.
//
// Solidity: function slasher() constant returns(address)
func (_DarknodeRegistry *DarknodeRegistryCaller) Slasher(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "slasher")
	return *ret0, err
}

// Slasher is a free data retrieval call binding the contract method 0xb1344271.
//
// Solidity: function slasher() constant returns(address)
func (_DarknodeRegistry *DarknodeRegistrySession) Slasher() (common.Address, error) {
	return _DarknodeRegistry.Contract.Slasher(&_DarknodeRegistry.CallOpts)
}

// Slasher is a free data retrieval call binding the contract method 0xb1344271.
//
// Solidity: function slasher() constant returns(address)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) Slasher() (common.Address, error) {
	return _DarknodeRegistry.Contract.Slasher(&_DarknodeRegistry.CallOpts)
}

// Store is a free data retrieval call binding the contract method 0x975057e7.
//
// Solidity: function store() constant returns(address)
func (_DarknodeRegistry *DarknodeRegistryCaller) Store(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "store")
	return *ret0, err
}

// Store is a free data retrieval call binding the contract method 0x975057e7.
//
// Solidity: function store() constant returns(address)
func (_DarknodeRegistry *DarknodeRegistrySession) Store() (common.Address, error) {
	return _DarknodeRegistry.Contract.Store(&_DarknodeRegistry.CallOpts)
}

// Store is a free data retrieval call binding the contract method 0x975057e7.
//
// Solidity: function store() constant returns(address)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) Store() (common.Address, error) {
	return _DarknodeRegistry.Contract.Store(&_DarknodeRegistry.CallOpts)
}

// Deregister is a paid mutator transaction binding the contract method 0x84ac33ec.
//
// Solidity: function deregister(_darknodeID address) returns()
func (_DarknodeRegistry *DarknodeRegistryTransactor) Deregister(opts *bind.TransactOpts, _darknodeID common.Address) (*types.Transaction, error) {
	return _DarknodeRegistry.contract.Transact(opts, "deregister", _darknodeID)
}

// Deregister is a paid mutator transaction binding the contract method 0x84ac33ec.
//
// Solidity: function deregister(_darknodeID address) returns()
func (_DarknodeRegistry *DarknodeRegistrySession) Deregister(_darknodeID common.Address) (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.Deregister(&_DarknodeRegistry.TransactOpts, _darknodeID)
}

// Deregister is a paid mutator transaction binding the contract method 0x84ac33ec.
//
// Solidity: function deregister(_darknodeID address) returns()
func (_DarknodeRegistry *DarknodeRegistryTransactorSession) Deregister(_darknodeID common.Address) (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.Deregister(&_DarknodeRegistry.TransactOpts, _darknodeID)
}

// Epoch is a paid mutator transaction binding the contract method 0x900cf0cf.
//
// Solidity: function epoch() returns()
func (_DarknodeRegistry *DarknodeRegistryTransactor) Epoch(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DarknodeRegistry.contract.Transact(opts, "epoch")
}

// Epoch is a paid mutator transaction binding the contract method 0x900cf0cf.
//
// Solidity: function epoch() returns()
func (_DarknodeRegistry *DarknodeRegistrySession) Epoch() (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.Epoch(&_DarknodeRegistry.TransactOpts)
}

// Epoch is a paid mutator transaction binding the contract method 0x900cf0cf.
//
// Solidity: function epoch() returns()
func (_DarknodeRegistry *DarknodeRegistryTransactorSession) Epoch() (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.Epoch(&_DarknodeRegistry.TransactOpts)
}

// Refund is a paid mutator transaction binding the contract method 0xfa89401a.
//
// Solidity: function refund(_darknodeID address) returns()
func (_DarknodeRegistry *DarknodeRegistryTransactor) Refund(opts *bind.TransactOpts, _darknodeID common.Address) (*types.Transaction, error) {
	return _DarknodeRegistry.contract.Transact(opts, "refund", _darknodeID)
}

// Refund is a paid mutator transaction binding the contract method 0xfa89401a.
//
// Solidity: function refund(_darknodeID address) returns()
func (_DarknodeRegistry *DarknodeRegistrySession) Refund(_darknodeID common.Address) (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.Refund(&_DarknodeRegistry.TransactOpts, _darknodeID)
}

// Refund is a paid mutator transaction binding the contract method 0xfa89401a.
//
// Solidity: function refund(_darknodeID address) returns()
func (_DarknodeRegistry *DarknodeRegistryTransactorSession) Refund(_darknodeID common.Address) (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.Refund(&_DarknodeRegistry.TransactOpts, _darknodeID)
}

// Register is a paid mutator transaction binding the contract method 0x0aeb6b40.
//
// Solidity: function register(_darknodeID address, _publicKey bytes, _bond uint256) returns()
func (_DarknodeRegistry *DarknodeRegistryTransactor) Register(opts *bind.TransactOpts, _darknodeID common.Address, _publicKey []byte, _bond *big.Int) (*types.Transaction, error) {
	return _DarknodeRegistry.contract.Transact(opts, "register", _darknodeID, _publicKey, _bond)
}

// Register is a paid mutator transaction binding the contract method 0x0aeb6b40.
//
// Solidity: function register(_darknodeID address, _publicKey bytes, _bond uint256) returns()
func (_DarknodeRegistry *DarknodeRegistrySession) Register(_darknodeID common.Address, _publicKey []byte, _bond *big.Int) (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.Register(&_DarknodeRegistry.TransactOpts, _darknodeID, _publicKey, _bond)
}

// Register is a paid mutator transaction binding the contract method 0x0aeb6b40.
//
// Solidity: function register(_darknodeID address, _publicKey bytes, _bond uint256) returns()
func (_DarknodeRegistry *DarknodeRegistryTransactorSession) Register(_darknodeID common.Address, _publicKey []byte, _bond *big.Int) (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.Register(&_DarknodeRegistry.TransactOpts, _darknodeID, _publicKey, _bond)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DarknodeRegistry *DarknodeRegistryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DarknodeRegistry.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DarknodeRegistry *DarknodeRegistrySession) RenounceOwnership() (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.RenounceOwnership(&_DarknodeRegistry.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DarknodeRegistry *DarknodeRegistryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.RenounceOwnership(&_DarknodeRegistry.TransactOpts)
}

// Slash is a paid mutator transaction binding the contract method 0x563bf264.
//
// Solidity: function slash(_prover address, _challenger1 address, _challenger2 address) returns()
func (_DarknodeRegistry *DarknodeRegistryTransactor) Slash(opts *bind.TransactOpts, _prover common.Address, _challenger1 common.Address, _challenger2 common.Address) (*types.Transaction, error) {
	return _DarknodeRegistry.contract.Transact(opts, "slash", _prover, _challenger1, _challenger2)
}

// Slash is a paid mutator transaction binding the contract method 0x563bf264.
//
// Solidity: function slash(_prover address, _challenger1 address, _challenger2 address) returns()
func (_DarknodeRegistry *DarknodeRegistrySession) Slash(_prover common.Address, _challenger1 common.Address, _challenger2 common.Address) (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.Slash(&_DarknodeRegistry.TransactOpts, _prover, _challenger1, _challenger2)
}

// Slash is a paid mutator transaction binding the contract method 0x563bf264.
//
// Solidity: function slash(_prover address, _challenger1 address, _challenger2 address) returns()
func (_DarknodeRegistry *DarknodeRegistryTransactorSession) Slash(_prover common.Address, _challenger1 common.Address, _challenger2 common.Address) (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.Slash(&_DarknodeRegistry.TransactOpts, _prover, _challenger1, _challenger2)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_DarknodeRegistry *DarknodeRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _DarknodeRegistry.contract.Transact(opts, "transferOwnership", _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_DarknodeRegistry *DarknodeRegistrySession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.TransferOwnership(&_DarknodeRegistry.TransactOpts, _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_DarknodeRegistry *DarknodeRegistryTransactorSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.TransferOwnership(&_DarknodeRegistry.TransactOpts, _newOwner)
}

// TransferStoreOwnership is a paid mutator transaction binding the contract method 0xc2250a99.
//
// Solidity: function transferStoreOwnership(_newOwner address) returns()
func (_DarknodeRegistry *DarknodeRegistryTransactor) TransferStoreOwnership(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _DarknodeRegistry.contract.Transact(opts, "transferStoreOwnership", _newOwner)
}

// TransferStoreOwnership is a paid mutator transaction binding the contract method 0xc2250a99.
//
// Solidity: function transferStoreOwnership(_newOwner address) returns()
func (_DarknodeRegistry *DarknodeRegistrySession) TransferStoreOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.TransferStoreOwnership(&_DarknodeRegistry.TransactOpts, _newOwner)
}

// TransferStoreOwnership is a paid mutator transaction binding the contract method 0xc2250a99.
//
// Solidity: function transferStoreOwnership(_newOwner address) returns()
func (_DarknodeRegistry *DarknodeRegistryTransactorSession) TransferStoreOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.TransferStoreOwnership(&_DarknodeRegistry.TransactOpts, _newOwner)
}

// UpdateMinimumBond is a paid mutator transaction binding the contract method 0x0ff9aafe.
//
// Solidity: function updateMinimumBond(_nextMinimumBond uint256) returns()
func (_DarknodeRegistry *DarknodeRegistryTransactor) UpdateMinimumBond(opts *bind.TransactOpts, _nextMinimumBond *big.Int) (*types.Transaction, error) {
	return _DarknodeRegistry.contract.Transact(opts, "updateMinimumBond", _nextMinimumBond)
}

// UpdateMinimumBond is a paid mutator transaction binding the contract method 0x0ff9aafe.
//
// Solidity: function updateMinimumBond(_nextMinimumBond uint256) returns()
func (_DarknodeRegistry *DarknodeRegistrySession) UpdateMinimumBond(_nextMinimumBond *big.Int) (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.UpdateMinimumBond(&_DarknodeRegistry.TransactOpts, _nextMinimumBond)
}

// UpdateMinimumBond is a paid mutator transaction binding the contract method 0x0ff9aafe.
//
// Solidity: function updateMinimumBond(_nextMinimumBond uint256) returns()
func (_DarknodeRegistry *DarknodeRegistryTransactorSession) UpdateMinimumBond(_nextMinimumBond *big.Int) (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.UpdateMinimumBond(&_DarknodeRegistry.TransactOpts, _nextMinimumBond)
}

// UpdateMinimumEpochInterval is a paid mutator transaction binding the contract method 0x63b851b9.
//
// Solidity: function updateMinimumEpochInterval(_nextMinimumEpochInterval uint256) returns()
func (_DarknodeRegistry *DarknodeRegistryTransactor) UpdateMinimumEpochInterval(opts *bind.TransactOpts, _nextMinimumEpochInterval *big.Int) (*types.Transaction, error) {
	return _DarknodeRegistry.contract.Transact(opts, "updateMinimumEpochInterval", _nextMinimumEpochInterval)
}

// UpdateMinimumEpochInterval is a paid mutator transaction binding the contract method 0x63b851b9.
//
// Solidity: function updateMinimumEpochInterval(_nextMinimumEpochInterval uint256) returns()
func (_DarknodeRegistry *DarknodeRegistrySession) UpdateMinimumEpochInterval(_nextMinimumEpochInterval *big.Int) (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.UpdateMinimumEpochInterval(&_DarknodeRegistry.TransactOpts, _nextMinimumEpochInterval)
}

// UpdateMinimumEpochInterval is a paid mutator transaction binding the contract method 0x63b851b9.
//
// Solidity: function updateMinimumEpochInterval(_nextMinimumEpochInterval uint256) returns()
func (_DarknodeRegistry *DarknodeRegistryTransactorSession) UpdateMinimumEpochInterval(_nextMinimumEpochInterval *big.Int) (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.UpdateMinimumEpochInterval(&_DarknodeRegistry.TransactOpts, _nextMinimumEpochInterval)
}

// UpdateMinimumPodSize is a paid mutator transaction binding the contract method 0x80a0c461.
//
// Solidity: function updateMinimumPodSize(_nextMinimumPodSize uint256) returns()
func (_DarknodeRegistry *DarknodeRegistryTransactor) UpdateMinimumPodSize(opts *bind.TransactOpts, _nextMinimumPodSize *big.Int) (*types.Transaction, error) {
	return _DarknodeRegistry.contract.Transact(opts, "updateMinimumPodSize", _nextMinimumPodSize)
}

// UpdateMinimumPodSize is a paid mutator transaction binding the contract method 0x80a0c461.
//
// Solidity: function updateMinimumPodSize(_nextMinimumPodSize uint256) returns()
func (_DarknodeRegistry *DarknodeRegistrySession) UpdateMinimumPodSize(_nextMinimumPodSize *big.Int) (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.UpdateMinimumPodSize(&_DarknodeRegistry.TransactOpts, _nextMinimumPodSize)
}

// UpdateMinimumPodSize is a paid mutator transaction binding the contract method 0x80a0c461.
//
// Solidity: function updateMinimumPodSize(_nextMinimumPodSize uint256) returns()
func (_DarknodeRegistry *DarknodeRegistryTransactorSession) UpdateMinimumPodSize(_nextMinimumPodSize *big.Int) (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.UpdateMinimumPodSize(&_DarknodeRegistry.TransactOpts, _nextMinimumPodSize)
}

// UpdateSlasher is a paid mutator transaction binding the contract method 0xb3139d38.
//
// Solidity: function updateSlasher(_slasher address) returns()
func (_DarknodeRegistry *DarknodeRegistryTransactor) UpdateSlasher(opts *bind.TransactOpts, _slasher common.Address) (*types.Transaction, error) {
	return _DarknodeRegistry.contract.Transact(opts, "updateSlasher", _slasher)
}

// UpdateSlasher is a paid mutator transaction binding the contract method 0xb3139d38.
//
// Solidity: function updateSlasher(_slasher address) returns()
func (_DarknodeRegistry *DarknodeRegistrySession) UpdateSlasher(_slasher common.Address) (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.UpdateSlasher(&_DarknodeRegistry.TransactOpts, _slasher)
}

// UpdateSlasher is a paid mutator transaction binding the contract method 0xb3139d38.
//
// Solidity: function updateSlasher(_slasher address) returns()
func (_DarknodeRegistry *DarknodeRegistryTransactorSession) UpdateSlasher(_slasher common.Address) (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.UpdateSlasher(&_DarknodeRegistry.TransactOpts, _slasher)
}

// DarknodeRegistryLogDarknodeDeregisteredIterator is returned from FilterLogDarknodeDeregistered and is used to iterate over the raw logs and unpacked data for LogDarknodeDeregistered events raised by the DarknodeRegistry contract.
type DarknodeRegistryLogDarknodeDeregisteredIterator struct {
	Event *DarknodeRegistryLogDarknodeDeregistered // Event containing the contract specifics and raw log

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
func (it *DarknodeRegistryLogDarknodeDeregisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DarknodeRegistryLogDarknodeDeregistered)
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
		it.Event = new(DarknodeRegistryLogDarknodeDeregistered)
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
func (it *DarknodeRegistryLogDarknodeDeregisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DarknodeRegistryLogDarknodeDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DarknodeRegistryLogDarknodeDeregistered represents a LogDarknodeDeregistered event raised by the DarknodeRegistry contract.
type DarknodeRegistryLogDarknodeDeregistered struct {
	DarknodeID common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogDarknodeDeregistered is a free log retrieval operation binding the contract event 0x2dc89de5703d2c341a22ebfc7c4d3f197e5e1f0c19bc2e1135f387163cb927e4.
//
// Solidity: e LogDarknodeDeregistered(_darknodeID address)
func (_DarknodeRegistry *DarknodeRegistryFilterer) FilterLogDarknodeDeregistered(opts *bind.FilterOpts) (*DarknodeRegistryLogDarknodeDeregisteredIterator, error) {

	logs, sub, err := _DarknodeRegistry.contract.FilterLogs(opts, "LogDarknodeDeregistered")
	if err != nil {
		return nil, err
	}
	return &DarknodeRegistryLogDarknodeDeregisteredIterator{contract: _DarknodeRegistry.contract, event: "LogDarknodeDeregistered", logs: logs, sub: sub}, nil
}

// WatchLogDarknodeDeregistered is a free log subscription operation binding the contract event 0x2dc89de5703d2c341a22ebfc7c4d3f197e5e1f0c19bc2e1135f387163cb927e4.
//
// Solidity: e LogDarknodeDeregistered(_darknodeID address)
func (_DarknodeRegistry *DarknodeRegistryFilterer) WatchLogDarknodeDeregistered(opts *bind.WatchOpts, sink chan<- *DarknodeRegistryLogDarknodeDeregistered) (event.Subscription, error) {

	logs, sub, err := _DarknodeRegistry.contract.WatchLogs(opts, "LogDarknodeDeregistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DarknodeRegistryLogDarknodeDeregistered)
				if err := _DarknodeRegistry.contract.UnpackLog(event, "LogDarknodeDeregistered", log); err != nil {
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

// DarknodeRegistryLogDarknodeOwnerRefundedIterator is returned from FilterLogDarknodeOwnerRefunded and is used to iterate over the raw logs and unpacked data for LogDarknodeOwnerRefunded events raised by the DarknodeRegistry contract.
type DarknodeRegistryLogDarknodeOwnerRefundedIterator struct {
	Event *DarknodeRegistryLogDarknodeOwnerRefunded // Event containing the contract specifics and raw log

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
func (it *DarknodeRegistryLogDarknodeOwnerRefundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DarknodeRegistryLogDarknodeOwnerRefunded)
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
		it.Event = new(DarknodeRegistryLogDarknodeOwnerRefunded)
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
func (it *DarknodeRegistryLogDarknodeOwnerRefundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DarknodeRegistryLogDarknodeOwnerRefundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DarknodeRegistryLogDarknodeOwnerRefunded represents a LogDarknodeOwnerRefunded event raised by the DarknodeRegistry contract.
type DarknodeRegistryLogDarknodeOwnerRefunded struct {
	Owner  common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterLogDarknodeOwnerRefunded is a free log retrieval operation binding the contract event 0x96ab9e56c79eee4a72db6e2879cbfbecdba5c65b411f4861824e66b89df19764.
//
// Solidity: e LogDarknodeOwnerRefunded(_owner address, _amount uint256)
func (_DarknodeRegistry *DarknodeRegistryFilterer) FilterLogDarknodeOwnerRefunded(opts *bind.FilterOpts) (*DarknodeRegistryLogDarknodeOwnerRefundedIterator, error) {

	logs, sub, err := _DarknodeRegistry.contract.FilterLogs(opts, "LogDarknodeOwnerRefunded")
	if err != nil {
		return nil, err
	}
	return &DarknodeRegistryLogDarknodeOwnerRefundedIterator{contract: _DarknodeRegistry.contract, event: "LogDarknodeOwnerRefunded", logs: logs, sub: sub}, nil
}

// WatchLogDarknodeOwnerRefunded is a free log subscription operation binding the contract event 0x96ab9e56c79eee4a72db6e2879cbfbecdba5c65b411f4861824e66b89df19764.
//
// Solidity: e LogDarknodeOwnerRefunded(_owner address, _amount uint256)
func (_DarknodeRegistry *DarknodeRegistryFilterer) WatchLogDarknodeOwnerRefunded(opts *bind.WatchOpts, sink chan<- *DarknodeRegistryLogDarknodeOwnerRefunded) (event.Subscription, error) {

	logs, sub, err := _DarknodeRegistry.contract.WatchLogs(opts, "LogDarknodeOwnerRefunded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DarknodeRegistryLogDarknodeOwnerRefunded)
				if err := _DarknodeRegistry.contract.UnpackLog(event, "LogDarknodeOwnerRefunded", log); err != nil {
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

// DarknodeRegistryLogDarknodeRegisteredIterator is returned from FilterLogDarknodeRegistered and is used to iterate over the raw logs and unpacked data for LogDarknodeRegistered events raised by the DarknodeRegistry contract.
type DarknodeRegistryLogDarknodeRegisteredIterator struct {
	Event *DarknodeRegistryLogDarknodeRegistered // Event containing the contract specifics and raw log

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
func (it *DarknodeRegistryLogDarknodeRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DarknodeRegistryLogDarknodeRegistered)
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
		it.Event = new(DarknodeRegistryLogDarknodeRegistered)
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
func (it *DarknodeRegistryLogDarknodeRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DarknodeRegistryLogDarknodeRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DarknodeRegistryLogDarknodeRegistered represents a LogDarknodeRegistered event raised by the DarknodeRegistry contract.
type DarknodeRegistryLogDarknodeRegistered struct {
	DarknodeID common.Address
	Bond       *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogDarknodeRegistered is a free log retrieval operation binding the contract event 0xd2819ba4c736158371edf0be38fd8d1fc435609832e392f118c4c79160e5bd7b.
//
// Solidity: e LogDarknodeRegistered(_darknodeID address, _bond uint256)
func (_DarknodeRegistry *DarknodeRegistryFilterer) FilterLogDarknodeRegistered(opts *bind.FilterOpts) (*DarknodeRegistryLogDarknodeRegisteredIterator, error) {

	logs, sub, err := _DarknodeRegistry.contract.FilterLogs(opts, "LogDarknodeRegistered")
	if err != nil {
		return nil, err
	}
	return &DarknodeRegistryLogDarknodeRegisteredIterator{contract: _DarknodeRegistry.contract, event: "LogDarknodeRegistered", logs: logs, sub: sub}, nil
}

// WatchLogDarknodeRegistered is a free log subscription operation binding the contract event 0xd2819ba4c736158371edf0be38fd8d1fc435609832e392f118c4c79160e5bd7b.
//
// Solidity: e LogDarknodeRegistered(_darknodeID address, _bond uint256)
func (_DarknodeRegistry *DarknodeRegistryFilterer) WatchLogDarknodeRegistered(opts *bind.WatchOpts, sink chan<- *DarknodeRegistryLogDarknodeRegistered) (event.Subscription, error) {

	logs, sub, err := _DarknodeRegistry.contract.WatchLogs(opts, "LogDarknodeRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DarknodeRegistryLogDarknodeRegistered)
				if err := _DarknodeRegistry.contract.UnpackLog(event, "LogDarknodeRegistered", log); err != nil {
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

// DarknodeRegistryLogMinimumBondUpdatedIterator is returned from FilterLogMinimumBondUpdated and is used to iterate over the raw logs and unpacked data for LogMinimumBondUpdated events raised by the DarknodeRegistry contract.
type DarknodeRegistryLogMinimumBondUpdatedIterator struct {
	Event *DarknodeRegistryLogMinimumBondUpdated // Event containing the contract specifics and raw log

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
func (it *DarknodeRegistryLogMinimumBondUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DarknodeRegistryLogMinimumBondUpdated)
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
		it.Event = new(DarknodeRegistryLogMinimumBondUpdated)
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
func (it *DarknodeRegistryLogMinimumBondUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DarknodeRegistryLogMinimumBondUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DarknodeRegistryLogMinimumBondUpdated represents a LogMinimumBondUpdated event raised by the DarknodeRegistry contract.
type DarknodeRegistryLogMinimumBondUpdated struct {
	PreviousMinimumBond *big.Int
	NextMinimumBond     *big.Int
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterLogMinimumBondUpdated is a free log retrieval operation binding the contract event 0x7c77c94944e9e4e5b0d46f1297127d060020792687cd743401d782346c68f655.
//
// Solidity: e LogMinimumBondUpdated(previousMinimumBond uint256, nextMinimumBond uint256)
func (_DarknodeRegistry *DarknodeRegistryFilterer) FilterLogMinimumBondUpdated(opts *bind.FilterOpts) (*DarknodeRegistryLogMinimumBondUpdatedIterator, error) {

	logs, sub, err := _DarknodeRegistry.contract.FilterLogs(opts, "LogMinimumBondUpdated")
	if err != nil {
		return nil, err
	}
	return &DarknodeRegistryLogMinimumBondUpdatedIterator{contract: _DarknodeRegistry.contract, event: "LogMinimumBondUpdated", logs: logs, sub: sub}, nil
}

// WatchLogMinimumBondUpdated is a free log subscription operation binding the contract event 0x7c77c94944e9e4e5b0d46f1297127d060020792687cd743401d782346c68f655.
//
// Solidity: e LogMinimumBondUpdated(previousMinimumBond uint256, nextMinimumBond uint256)
func (_DarknodeRegistry *DarknodeRegistryFilterer) WatchLogMinimumBondUpdated(opts *bind.WatchOpts, sink chan<- *DarknodeRegistryLogMinimumBondUpdated) (event.Subscription, error) {

	logs, sub, err := _DarknodeRegistry.contract.WatchLogs(opts, "LogMinimumBondUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DarknodeRegistryLogMinimumBondUpdated)
				if err := _DarknodeRegistry.contract.UnpackLog(event, "LogMinimumBondUpdated", log); err != nil {
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

// DarknodeRegistryLogMinimumEpochIntervalUpdatedIterator is returned from FilterLogMinimumEpochIntervalUpdated and is used to iterate over the raw logs and unpacked data for LogMinimumEpochIntervalUpdated events raised by the DarknodeRegistry contract.
type DarknodeRegistryLogMinimumEpochIntervalUpdatedIterator struct {
	Event *DarknodeRegistryLogMinimumEpochIntervalUpdated // Event containing the contract specifics and raw log

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
func (it *DarknodeRegistryLogMinimumEpochIntervalUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DarknodeRegistryLogMinimumEpochIntervalUpdated)
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
		it.Event = new(DarknodeRegistryLogMinimumEpochIntervalUpdated)
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
func (it *DarknodeRegistryLogMinimumEpochIntervalUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DarknodeRegistryLogMinimumEpochIntervalUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DarknodeRegistryLogMinimumEpochIntervalUpdated represents a LogMinimumEpochIntervalUpdated event raised by the DarknodeRegistry contract.
type DarknodeRegistryLogMinimumEpochIntervalUpdated struct {
	PreviousMinimumEpochInterval *big.Int
	NextMinimumEpochInterval     *big.Int
	Raw                          types.Log // Blockchain specific contextual infos
}

// FilterLogMinimumEpochIntervalUpdated is a free log retrieval operation binding the contract event 0xb218cde2730b79a0667ddf869466ee66a12ef56fe65fa4986a590f8a7108c9de.
//
// Solidity: e LogMinimumEpochIntervalUpdated(previousMinimumEpochInterval uint256, nextMinimumEpochInterval uint256)
func (_DarknodeRegistry *DarknodeRegistryFilterer) FilterLogMinimumEpochIntervalUpdated(opts *bind.FilterOpts) (*DarknodeRegistryLogMinimumEpochIntervalUpdatedIterator, error) {

	logs, sub, err := _DarknodeRegistry.contract.FilterLogs(opts, "LogMinimumEpochIntervalUpdated")
	if err != nil {
		return nil, err
	}
	return &DarknodeRegistryLogMinimumEpochIntervalUpdatedIterator{contract: _DarknodeRegistry.contract, event: "LogMinimumEpochIntervalUpdated", logs: logs, sub: sub}, nil
}

// WatchLogMinimumEpochIntervalUpdated is a free log subscription operation binding the contract event 0xb218cde2730b79a0667ddf869466ee66a12ef56fe65fa4986a590f8a7108c9de.
//
// Solidity: e LogMinimumEpochIntervalUpdated(previousMinimumEpochInterval uint256, nextMinimumEpochInterval uint256)
func (_DarknodeRegistry *DarknodeRegistryFilterer) WatchLogMinimumEpochIntervalUpdated(opts *bind.WatchOpts, sink chan<- *DarknodeRegistryLogMinimumEpochIntervalUpdated) (event.Subscription, error) {

	logs, sub, err := _DarknodeRegistry.contract.WatchLogs(opts, "LogMinimumEpochIntervalUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DarknodeRegistryLogMinimumEpochIntervalUpdated)
				if err := _DarknodeRegistry.contract.UnpackLog(event, "LogMinimumEpochIntervalUpdated", log); err != nil {
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

// DarknodeRegistryLogMinimumPodSizeUpdatedIterator is returned from FilterLogMinimumPodSizeUpdated and is used to iterate over the raw logs and unpacked data for LogMinimumPodSizeUpdated events raised by the DarknodeRegistry contract.
type DarknodeRegistryLogMinimumPodSizeUpdatedIterator struct {
	Event *DarknodeRegistryLogMinimumPodSizeUpdated // Event containing the contract specifics and raw log

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
func (it *DarknodeRegistryLogMinimumPodSizeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DarknodeRegistryLogMinimumPodSizeUpdated)
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
		it.Event = new(DarknodeRegistryLogMinimumPodSizeUpdated)
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
func (it *DarknodeRegistryLogMinimumPodSizeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DarknodeRegistryLogMinimumPodSizeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DarknodeRegistryLogMinimumPodSizeUpdated represents a LogMinimumPodSizeUpdated event raised by the DarknodeRegistry contract.
type DarknodeRegistryLogMinimumPodSizeUpdated struct {
	PreviousMinimumPodSize *big.Int
	NextMinimumPodSize     *big.Int
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterLogMinimumPodSizeUpdated is a free log retrieval operation binding the contract event 0x6d520e46e5714982ddf8cb6216bcb3e1c1d5b79d337afc305335f819394f5d6a.
//
// Solidity: e LogMinimumPodSizeUpdated(previousMinimumPodSize uint256, nextMinimumPodSize uint256)
func (_DarknodeRegistry *DarknodeRegistryFilterer) FilterLogMinimumPodSizeUpdated(opts *bind.FilterOpts) (*DarknodeRegistryLogMinimumPodSizeUpdatedIterator, error) {

	logs, sub, err := _DarknodeRegistry.contract.FilterLogs(opts, "LogMinimumPodSizeUpdated")
	if err != nil {
		return nil, err
	}
	return &DarknodeRegistryLogMinimumPodSizeUpdatedIterator{contract: _DarknodeRegistry.contract, event: "LogMinimumPodSizeUpdated", logs: logs, sub: sub}, nil
}

// WatchLogMinimumPodSizeUpdated is a free log subscription operation binding the contract event 0x6d520e46e5714982ddf8cb6216bcb3e1c1d5b79d337afc305335f819394f5d6a.
//
// Solidity: e LogMinimumPodSizeUpdated(previousMinimumPodSize uint256, nextMinimumPodSize uint256)
func (_DarknodeRegistry *DarknodeRegistryFilterer) WatchLogMinimumPodSizeUpdated(opts *bind.WatchOpts, sink chan<- *DarknodeRegistryLogMinimumPodSizeUpdated) (event.Subscription, error) {

	logs, sub, err := _DarknodeRegistry.contract.WatchLogs(opts, "LogMinimumPodSizeUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DarknodeRegistryLogMinimumPodSizeUpdated)
				if err := _DarknodeRegistry.contract.UnpackLog(event, "LogMinimumPodSizeUpdated", log); err != nil {
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

// DarknodeRegistryLogNewEpochIterator is returned from FilterLogNewEpoch and is used to iterate over the raw logs and unpacked data for LogNewEpoch events raised by the DarknodeRegistry contract.
type DarknodeRegistryLogNewEpochIterator struct {
	Event *DarknodeRegistryLogNewEpoch // Event containing the contract specifics and raw log

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
func (it *DarknodeRegistryLogNewEpochIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DarknodeRegistryLogNewEpoch)
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
		it.Event = new(DarknodeRegistryLogNewEpoch)
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
func (it *DarknodeRegistryLogNewEpochIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DarknodeRegistryLogNewEpochIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DarknodeRegistryLogNewEpoch represents a LogNewEpoch event raised by the DarknodeRegistry contract.
type DarknodeRegistryLogNewEpoch struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNewEpoch is a free log retrieval operation binding the contract event 0xeff7e281fe3b4211ed1f0a5e6419bcc40f4552974f771357e66926421f0a58e8.
//
// Solidity: e LogNewEpoch()
func (_DarknodeRegistry *DarknodeRegistryFilterer) FilterLogNewEpoch(opts *bind.FilterOpts) (*DarknodeRegistryLogNewEpochIterator, error) {

	logs, sub, err := _DarknodeRegistry.contract.FilterLogs(opts, "LogNewEpoch")
	if err != nil {
		return nil, err
	}
	return &DarknodeRegistryLogNewEpochIterator{contract: _DarknodeRegistry.contract, event: "LogNewEpoch", logs: logs, sub: sub}, nil
}

// WatchLogNewEpoch is a free log subscription operation binding the contract event 0xeff7e281fe3b4211ed1f0a5e6419bcc40f4552974f771357e66926421f0a58e8.
//
// Solidity: e LogNewEpoch()
func (_DarknodeRegistry *DarknodeRegistryFilterer) WatchLogNewEpoch(opts *bind.WatchOpts, sink chan<- *DarknodeRegistryLogNewEpoch) (event.Subscription, error) {

	logs, sub, err := _DarknodeRegistry.contract.WatchLogs(opts, "LogNewEpoch")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DarknodeRegistryLogNewEpoch)
				if err := _DarknodeRegistry.contract.UnpackLog(event, "LogNewEpoch", log); err != nil {
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

// DarknodeRegistryLogSlasherUpdatedIterator is returned from FilterLogSlasherUpdated and is used to iterate over the raw logs and unpacked data for LogSlasherUpdated events raised by the DarknodeRegistry contract.
type DarknodeRegistryLogSlasherUpdatedIterator struct {
	Event *DarknodeRegistryLogSlasherUpdated // Event containing the contract specifics and raw log

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
func (it *DarknodeRegistryLogSlasherUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DarknodeRegistryLogSlasherUpdated)
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
		it.Event = new(DarknodeRegistryLogSlasherUpdated)
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
func (it *DarknodeRegistryLogSlasherUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DarknodeRegistryLogSlasherUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DarknodeRegistryLogSlasherUpdated represents a LogSlasherUpdated event raised by the DarknodeRegistry contract.
type DarknodeRegistryLogSlasherUpdated struct {
	PreviousSlasher common.Address
	NextSlasher     common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterLogSlasherUpdated is a free log retrieval operation binding the contract event 0x933228a1c3ba8fadd3ce47a9db5b898be647f89af99ba7c1b9a655f59ea306c8.
//
// Solidity: e LogSlasherUpdated(previousSlasher address, nextSlasher address)
func (_DarknodeRegistry *DarknodeRegistryFilterer) FilterLogSlasherUpdated(opts *bind.FilterOpts) (*DarknodeRegistryLogSlasherUpdatedIterator, error) {

	logs, sub, err := _DarknodeRegistry.contract.FilterLogs(opts, "LogSlasherUpdated")
	if err != nil {
		return nil, err
	}
	return &DarknodeRegistryLogSlasherUpdatedIterator{contract: _DarknodeRegistry.contract, event: "LogSlasherUpdated", logs: logs, sub: sub}, nil
}

// WatchLogSlasherUpdated is a free log subscription operation binding the contract event 0x933228a1c3ba8fadd3ce47a9db5b898be647f89af99ba7c1b9a655f59ea306c8.
//
// Solidity: e LogSlasherUpdated(previousSlasher address, nextSlasher address)
func (_DarknodeRegistry *DarknodeRegistryFilterer) WatchLogSlasherUpdated(opts *bind.WatchOpts, sink chan<- *DarknodeRegistryLogSlasherUpdated) (event.Subscription, error) {

	logs, sub, err := _DarknodeRegistry.contract.WatchLogs(opts, "LogSlasherUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DarknodeRegistryLogSlasherUpdated)
				if err := _DarknodeRegistry.contract.UnpackLog(event, "LogSlasherUpdated", log); err != nil {
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

// DarknodeRegistryOwnershipRenouncedIterator is returned from FilterOwnershipRenounced and is used to iterate over the raw logs and unpacked data for OwnershipRenounced events raised by the DarknodeRegistry contract.
type DarknodeRegistryOwnershipRenouncedIterator struct {
	Event *DarknodeRegistryOwnershipRenounced // Event containing the contract specifics and raw log

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
func (it *DarknodeRegistryOwnershipRenouncedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DarknodeRegistryOwnershipRenounced)
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
		it.Event = new(DarknodeRegistryOwnershipRenounced)
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
func (it *DarknodeRegistryOwnershipRenouncedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DarknodeRegistryOwnershipRenouncedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DarknodeRegistryOwnershipRenounced represents a OwnershipRenounced event raised by the DarknodeRegistry contract.
type DarknodeRegistryOwnershipRenounced struct {
	PreviousOwner common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipRenounced is a free log retrieval operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_DarknodeRegistry *DarknodeRegistryFilterer) FilterOwnershipRenounced(opts *bind.FilterOpts, previousOwner []common.Address) (*DarknodeRegistryOwnershipRenouncedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _DarknodeRegistry.contract.FilterLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DarknodeRegistryOwnershipRenouncedIterator{contract: _DarknodeRegistry.contract, event: "OwnershipRenounced", logs: logs, sub: sub}, nil
}

// WatchOwnershipRenounced is a free log subscription operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_DarknodeRegistry *DarknodeRegistryFilterer) WatchOwnershipRenounced(opts *bind.WatchOpts, sink chan<- *DarknodeRegistryOwnershipRenounced, previousOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _DarknodeRegistry.contract.WatchLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DarknodeRegistryOwnershipRenounced)
				if err := _DarknodeRegistry.contract.UnpackLog(event, "OwnershipRenounced", log); err != nil {
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

// DarknodeRegistryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the DarknodeRegistry contract.
type DarknodeRegistryOwnershipTransferredIterator struct {
	Event *DarknodeRegistryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *DarknodeRegistryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DarknodeRegistryOwnershipTransferred)
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
		it.Event = new(DarknodeRegistryOwnershipTransferred)
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
func (it *DarknodeRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DarknodeRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DarknodeRegistryOwnershipTransferred represents a OwnershipTransferred event raised by the DarknodeRegistry contract.
type DarknodeRegistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_DarknodeRegistry *DarknodeRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*DarknodeRegistryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DarknodeRegistry.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DarknodeRegistryOwnershipTransferredIterator{contract: _DarknodeRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_DarknodeRegistry *DarknodeRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DarknodeRegistryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DarknodeRegistry.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DarknodeRegistryOwnershipTransferred)
				if err := _DarknodeRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// DarknodeRegistryStoreABI is the input ABI used to generate the binding from.
const DarknodeRegistryStoreABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"darknodeID\",\"type\":\"address\"}],\"name\":\"darknodeDeregisteredAt\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"begin\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"darknodeID\",\"type\":\"address\"},{\"name\":\"deregisteredAt\",\"type\":\"uint256\"}],\"name\":\"updateDarknodeDeregisteredAt\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"darknodeID\",\"type\":\"address\"}],\"name\":\"removeDarknode\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"ren\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"darknodeID\",\"type\":\"address\"}],\"name\":\"darknodeOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_darknodeID\",\"type\":\"address\"},{\"name\":\"_darknodeOwner\",\"type\":\"address\"},{\"name\":\"_bond\",\"type\":\"uint256\"},{\"name\":\"_publicKey\",\"type\":\"bytes\"},{\"name\":\"_registeredAt\",\"type\":\"uint256\"},{\"name\":\"_deregisteredAt\",\"type\":\"uint256\"}],\"name\":\"appendDarknode\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"darknodeID\",\"type\":\"address\"}],\"name\":\"next\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"darknodeID\",\"type\":\"address\"}],\"name\":\"darknodeBond\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"darknodeID\",\"type\":\"address\"}],\"name\":\"darknodeRegisteredAt\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"darknodeID\",\"type\":\"address\"}],\"name\":\"darknodePublicKey\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"darknodeID\",\"type\":\"address\"},{\"name\":\"bond\",\"type\":\"uint256\"}],\"name\":\"updateDarknodeBond\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_VERSION\",\"type\":\"string\"},{\"name\":\"_ren\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"}],\"name\":\"OwnershipRenounced\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"

// DarknodeRegistryStoreBin is the compiled bytecode used for deploying new contracts.
const DarknodeRegistryStoreBin = `0x608060405234801561001057600080fd5b5060405161113038038061113083398101604052805160208083015160008054600160a060020a0319163317905591909201805190926100559160019185019061007c565b5060048054600160a060020a031916600160a060020a039290921691909117905550610117565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106100bd57805160ff19168380011785556100ea565b828001600101855582156100ea579182015b828111156100ea5782518255916020019190600101906100cf565b506100f69291506100fa565b5090565b61011491905b808211156100f65760008155600101610100565b90565b61100a806101266000396000f3006080604052600436106100e55763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166301930b6e81146100ea5780631bce6ff31461011d5780633ac39d4b1461014e57806341b4439214610174578063715018a6146101955780638a9b4067146101aa5780638da5cb5b146101bf578063a3078815146101d4578063a85ef579146101f5578063ab73e31614610234578063cad4135714610255578063e2d7c64c14610276578063ee594a5014610297578063f2fde38b1461032d578063fbc402fc1461034e578063ffa1ad7414610372575b600080fd5b3480156100f657600080fd5b5061010b600160a060020a0360043516610387565b60408051918252519081900360200190f35b34801561012957600080fd5b506101326103be565b60408051600160a060020a039092168252519081900360200190f35b34801561015a57600080fd5b50610172600160a060020a03600435166024356103e6565b005b34801561018057600080fd5b50610172600160a060020a036004351661041c565b3480156101a157600080fd5b50610172610596565b3480156101b657600080fd5b50610132610602565b3480156101cb57600080fd5b50610132610611565b3480156101e057600080fd5b50610132600160a060020a0360043516610620565b34801561020157600080fd5b50610172600160a060020a036004803582169160248035909116916044359160643590810191013560843560a435610657565b34801561024057600080fd5b50610132600160a060020a0360043516610767565b34801561026157600080fd5b5061010b600160a060020a0360043516610790565b34801561028257600080fd5b5061010b600160a060020a03600435166107c7565b3480156102a357600080fd5b506102b8600160a060020a03600435166107ff565b6040805160208082528351818301528351919283929083019185019080838360005b838110156102f25781810151838201526020016102da565b50505050905090810190601f16801561031f5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561033957600080fd5b50610172600160a060020a03600435166108c3565b34801561035a57600080fd5b50610172600160a060020a03600435166024356108e6565b34801561037e57600080fd5b506102b8610a23565b60008054600160a060020a0316331461039f57600080fd5b50600160a060020a031660009081526002602052604090206003015490565b60008054600160a060020a031633146103d657600080fd5b6103e06003610ab0565b90505b90565b600054600160a060020a031633146103fd57600080fd5b600160a060020a03909116600090815260026020526040902060030155565b60008054600160a060020a0316331461043457600080fd5b50600160a060020a03811660009081526002602081905260408220600181018054825473ffffffffffffffffffffffffffffffffffffffff19168355908490559181018390556003810183905590916104906004830182610ec9565b505061049d600383610acf565b6004805460008054604080517fa9059cbb000000000000000000000000000000000000000000000000000000008152600160a060020a0392831695810195909552602485018690525192169263a9059cbb9260448083019360209383900390910190829087803b15801561051057600080fd5b505af1158015610524573d6000803e3d6000fd5b505050506040513d602081101561053a57600080fd5b50511515610592576040805160e560020a62461bcd02815260206004820152601460248201527f626f6e64207472616e73666572206661696c6564000000000000000000000000604482015290519081900360640190fd5b5050565b600054600160a060020a031633146105ad57600080fd5b60008054604051600160a060020a03909116917ff8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c6482091a26000805473ffffffffffffffffffffffffffffffffffffffff19169055565b600454600160a060020a031681565b600054600160a060020a031681565b60008054600160a060020a0316331461063857600080fd5b50600160a060020a039081166000908152600260205260409020541690565b61065f610f0d565b600054600160a060020a0316331461067657600080fd5b60a06040519081016040528088600160a060020a0316815260200187815260200184815260200183815260200186868080601f01602080910402602001604051908101604052809392919081815260200183838082843750505092909352505050600160a060020a038981166000908152600260208181526040928390208551815473ffffffffffffffffffffffffffffffffffffffff1916951694909417845584810151600185015591840151908301556060830151600383015560808301518051939450849361074e9260048501920190610f46565b5090505061075d600389610bf9565b5050505050505050565b60008054600160a060020a0316331461077f57600080fd5b61078a600383610c0c565b92915050565b60008054600160a060020a031633146107a857600080fd5b50600160a060020a031660009081526002602052604090206001015490565b60008054600160a060020a031633146107df57600080fd5b50600160a060020a03166000908152600260208190526040909120015490565b600054606090600160a060020a0316331461081957600080fd5b600160a060020a038216600090815260026020818152604092839020600401805484516001821615610100026000190190911693909304601f81018390048302840183019094528383529192908301828280156108b75780601f1061088c576101008083540402835291602001916108b7565b820191906000526020600020905b81548152906001019060200180831161089a57829003601f168201915b50505050509050919050565b600054600160a060020a031633146108da57600080fd5b6108e381610c92565b50565b60008054600160a060020a031633146108fe57600080fd5b50600160a060020a038216600090815260026020526040902060010180549082905581811115610a1e576004805460008054604080517fa9059cbb000000000000000000000000000000000000000000000000000000008152600160a060020a039283169581019590955286860360248601525192169263a9059cbb9260448083019360209383900390910190829087803b15801561099c57600080fd5b505af11580156109b0573d6000803e3d6000fd5b505050506040513d60208110156109c657600080fd5b50511515610a1e576040805160e560020a62461bcd02815260206004820152601460248201527f63616e6e6f74207472616e7366657220626f6e64000000000000000000000000604482015290519081900360640190fd5b505050565b60018054604080516020600284861615610100026000190190941693909304601f81018490048402820184019092528181529291830182828015610aa85780601f10610a7d57610100808354040283529160200191610aa8565b820191906000526020600020905b815481529060010190602001808311610a8b57829003601f168201915b505050505081565b60008080526020919091526040902060010154600160a060020a031690565b600080610adc8484610d0f565b1515610b32576040805160e560020a62461bcd02815260206004820152600b60248201527f6e6f7420696e206c697374000000000000000000000000000000000000000000604482015290519081900360640190fd5b600160a060020a0383161515610b4757610bf3565b5050600160a060020a0381811660008181526020859052604080822080546001808301805461010093849004891680885286882090930180549190991673ffffffffffffffffffffffffffffffffffffffff199182168117909955888752948620805474ffffffffffffffffffffffffffffffffffffffff0019169383029390931790925594909352805474ffffffffffffffffffffffffffffffffffffffffff191690558154169055905b50505050565b61059282610c0684610d2e565b83610d50565b6000610c188383610d0f565b1515610c6e576040805160e560020a62461bcd02815260206004820152600b60248201527f6e6f7420696e206c697374000000000000000000000000000000000000000000604482015290519081900360640190fd5b50600160a060020a0390811660009081526020929092526040909120600101541690565b600160a060020a0381161515610ca757600080fd5b60008054604051600160a060020a03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a36000805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b600160a060020a03166000908152602091909152604090205460ff1690565b6000808052602082905260409020546101009004600160a060020a0316919050565b6000610d5c8483610d0f565b15610db1576040805160e560020a62461bcd02815260206004820152600f60248201527f616c726561647920696e206c6973740000000000000000000000000000000000604482015290519081900360640190fd5b610dbb8484610d0f565b80610dcd5750600160a060020a038316155b1515610e23576040805160e560020a62461bcd02815260206004820152600b60248201527f6e6f7420696e206c697374000000000000000000000000000000000000000000604482015290519081900360640190fd5b50600160a060020a039182166000818152602094909452604080852060019081018054948616808852838820805474ffffffffffffffffffffffffffffffffffffffff001990811661010097880217825581850180549890991673ffffffffffffffffffffffffffffffffffffffff1998891681179099558354909716821790925595875291862080549094169285029290921790925591909252815460ff1916179055565b50805460018160011615610100020316600290046000825580601f10610eef57506108e3565b601f0160209004906000526020600020908101906108e39190610fc4565b60a0604051908101604052806000600160a060020a03168152602001600081526020016000815260200160008152602001606081525090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10610f8757805160ff1916838001178555610fb4565b82800160010185558215610fb4579182015b82811115610fb4578251825591602001919060010190610f99565b50610fc0929150610fc4565b5090565b6103e391905b80821115610fc05760008155600101610fca5600a165627a7a72305820015437ac611d5af2e1f7198c13f70eb87782ea2ab53e3ef00e887140465ff0b70029`

// DeployDarknodeRegistryStore deploys a new Ethereum contract, binding an instance of DarknodeRegistryStore to it.
func DeployDarknodeRegistryStore(auth *bind.TransactOpts, backend bind.ContractBackend, _VERSION string, _ren common.Address) (common.Address, *types.Transaction, *DarknodeRegistryStore, error) {
	parsed, err := abi.JSON(strings.NewReader(DarknodeRegistryStoreABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(DarknodeRegistryStoreBin), backend, _VERSION, _ren)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DarknodeRegistryStore{DarknodeRegistryStoreCaller: DarknodeRegistryStoreCaller{contract: contract}, DarknodeRegistryStoreTransactor: DarknodeRegistryStoreTransactor{contract: contract}, DarknodeRegistryStoreFilterer: DarknodeRegistryStoreFilterer{contract: contract}}, nil
}

// DarknodeRegistryStore is an auto generated Go binding around an Ethereum contract.
type DarknodeRegistryStore struct {
	DarknodeRegistryStoreCaller     // Read-only binding to the contract
	DarknodeRegistryStoreTransactor // Write-only binding to the contract
	DarknodeRegistryStoreFilterer   // Log filterer for contract events
}

// DarknodeRegistryStoreCaller is an auto generated read-only Go binding around an Ethereum contract.
type DarknodeRegistryStoreCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DarknodeRegistryStoreTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DarknodeRegistryStoreTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DarknodeRegistryStoreFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DarknodeRegistryStoreFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DarknodeRegistryStoreSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DarknodeRegistryStoreSession struct {
	Contract     *DarknodeRegistryStore // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// DarknodeRegistryStoreCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DarknodeRegistryStoreCallerSession struct {
	Contract *DarknodeRegistryStoreCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// DarknodeRegistryStoreTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DarknodeRegistryStoreTransactorSession struct {
	Contract     *DarknodeRegistryStoreTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// DarknodeRegistryStoreRaw is an auto generated low-level Go binding around an Ethereum contract.
type DarknodeRegistryStoreRaw struct {
	Contract *DarknodeRegistryStore // Generic contract binding to access the raw methods on
}

// DarknodeRegistryStoreCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DarknodeRegistryStoreCallerRaw struct {
	Contract *DarknodeRegistryStoreCaller // Generic read-only contract binding to access the raw methods on
}

// DarknodeRegistryStoreTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DarknodeRegistryStoreTransactorRaw struct {
	Contract *DarknodeRegistryStoreTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDarknodeRegistryStore creates a new instance of DarknodeRegistryStore, bound to a specific deployed contract.
func NewDarknodeRegistryStore(address common.Address, backend bind.ContractBackend) (*DarknodeRegistryStore, error) {
	contract, err := bindDarknodeRegistryStore(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DarknodeRegistryStore{DarknodeRegistryStoreCaller: DarknodeRegistryStoreCaller{contract: contract}, DarknodeRegistryStoreTransactor: DarknodeRegistryStoreTransactor{contract: contract}, DarknodeRegistryStoreFilterer: DarknodeRegistryStoreFilterer{contract: contract}}, nil
}

// NewDarknodeRegistryStoreCaller creates a new read-only instance of DarknodeRegistryStore, bound to a specific deployed contract.
func NewDarknodeRegistryStoreCaller(address common.Address, caller bind.ContractCaller) (*DarknodeRegistryStoreCaller, error) {
	contract, err := bindDarknodeRegistryStore(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DarknodeRegistryStoreCaller{contract: contract}, nil
}

// NewDarknodeRegistryStoreTransactor creates a new write-only instance of DarknodeRegistryStore, bound to a specific deployed contract.
func NewDarknodeRegistryStoreTransactor(address common.Address, transactor bind.ContractTransactor) (*DarknodeRegistryStoreTransactor, error) {
	contract, err := bindDarknodeRegistryStore(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DarknodeRegistryStoreTransactor{contract: contract}, nil
}

// NewDarknodeRegistryStoreFilterer creates a new log filterer instance of DarknodeRegistryStore, bound to a specific deployed contract.
func NewDarknodeRegistryStoreFilterer(address common.Address, filterer bind.ContractFilterer) (*DarknodeRegistryStoreFilterer, error) {
	contract, err := bindDarknodeRegistryStore(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DarknodeRegistryStoreFilterer{contract: contract}, nil
}

// bindDarknodeRegistryStore binds a generic wrapper to an already deployed contract.
func bindDarknodeRegistryStore(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DarknodeRegistryStoreABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DarknodeRegistryStore *DarknodeRegistryStoreRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _DarknodeRegistryStore.Contract.DarknodeRegistryStoreCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DarknodeRegistryStore *DarknodeRegistryStoreRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DarknodeRegistryStore.Contract.DarknodeRegistryStoreTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DarknodeRegistryStore *DarknodeRegistryStoreRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DarknodeRegistryStore.Contract.DarknodeRegistryStoreTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DarknodeRegistryStore *DarknodeRegistryStoreCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _DarknodeRegistryStore.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DarknodeRegistryStore *DarknodeRegistryStoreTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DarknodeRegistryStore.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DarknodeRegistryStore *DarknodeRegistryStoreTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DarknodeRegistryStore.Contract.contract.Transact(opts, method, params...)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_DarknodeRegistryStore *DarknodeRegistryStoreCaller) VERSION(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _DarknodeRegistryStore.contract.Call(opts, out, "VERSION")
	return *ret0, err
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_DarknodeRegistryStore *DarknodeRegistryStoreSession) VERSION() (string, error) {
	return _DarknodeRegistryStore.Contract.VERSION(&_DarknodeRegistryStore.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_DarknodeRegistryStore *DarknodeRegistryStoreCallerSession) VERSION() (string, error) {
	return _DarknodeRegistryStore.Contract.VERSION(&_DarknodeRegistryStore.CallOpts)
}

// Begin is a free data retrieval call binding the contract method 0x1bce6ff3.
//
// Solidity: function begin() constant returns(address)
func (_DarknodeRegistryStore *DarknodeRegistryStoreCaller) Begin(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DarknodeRegistryStore.contract.Call(opts, out, "begin")
	return *ret0, err
}

// Begin is a free data retrieval call binding the contract method 0x1bce6ff3.
//
// Solidity: function begin() constant returns(address)
func (_DarknodeRegistryStore *DarknodeRegistryStoreSession) Begin() (common.Address, error) {
	return _DarknodeRegistryStore.Contract.Begin(&_DarknodeRegistryStore.CallOpts)
}

// Begin is a free data retrieval call binding the contract method 0x1bce6ff3.
//
// Solidity: function begin() constant returns(address)
func (_DarknodeRegistryStore *DarknodeRegistryStoreCallerSession) Begin() (common.Address, error) {
	return _DarknodeRegistryStore.Contract.Begin(&_DarknodeRegistryStore.CallOpts)
}

// DarknodeBond is a free data retrieval call binding the contract method 0xcad41357.
//
// Solidity: function darknodeBond(darknodeID address) constant returns(uint256)
func (_DarknodeRegistryStore *DarknodeRegistryStoreCaller) DarknodeBond(opts *bind.CallOpts, darknodeID common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DarknodeRegistryStore.contract.Call(opts, out, "darknodeBond", darknodeID)
	return *ret0, err
}

// DarknodeBond is a free data retrieval call binding the contract method 0xcad41357.
//
// Solidity: function darknodeBond(darknodeID address) constant returns(uint256)
func (_DarknodeRegistryStore *DarknodeRegistryStoreSession) DarknodeBond(darknodeID common.Address) (*big.Int, error) {
	return _DarknodeRegistryStore.Contract.DarknodeBond(&_DarknodeRegistryStore.CallOpts, darknodeID)
}

// DarknodeBond is a free data retrieval call binding the contract method 0xcad41357.
//
// Solidity: function darknodeBond(darknodeID address) constant returns(uint256)
func (_DarknodeRegistryStore *DarknodeRegistryStoreCallerSession) DarknodeBond(darknodeID common.Address) (*big.Int, error) {
	return _DarknodeRegistryStore.Contract.DarknodeBond(&_DarknodeRegistryStore.CallOpts, darknodeID)
}

// DarknodeDeregisteredAt is a free data retrieval call binding the contract method 0x01930b6e.
//
// Solidity: function darknodeDeregisteredAt(darknodeID address) constant returns(uint256)
func (_DarknodeRegistryStore *DarknodeRegistryStoreCaller) DarknodeDeregisteredAt(opts *bind.CallOpts, darknodeID common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DarknodeRegistryStore.contract.Call(opts, out, "darknodeDeregisteredAt", darknodeID)
	return *ret0, err
}

// DarknodeDeregisteredAt is a free data retrieval call binding the contract method 0x01930b6e.
//
// Solidity: function darknodeDeregisteredAt(darknodeID address) constant returns(uint256)
func (_DarknodeRegistryStore *DarknodeRegistryStoreSession) DarknodeDeregisteredAt(darknodeID common.Address) (*big.Int, error) {
	return _DarknodeRegistryStore.Contract.DarknodeDeregisteredAt(&_DarknodeRegistryStore.CallOpts, darknodeID)
}

// DarknodeDeregisteredAt is a free data retrieval call binding the contract method 0x01930b6e.
//
// Solidity: function darknodeDeregisteredAt(darknodeID address) constant returns(uint256)
func (_DarknodeRegistryStore *DarknodeRegistryStoreCallerSession) DarknodeDeregisteredAt(darknodeID common.Address) (*big.Int, error) {
	return _DarknodeRegistryStore.Contract.DarknodeDeregisteredAt(&_DarknodeRegistryStore.CallOpts, darknodeID)
}

// DarknodeOwner is a free data retrieval call binding the contract method 0xa3078815.
//
// Solidity: function darknodeOwner(darknodeID address) constant returns(address)
func (_DarknodeRegistryStore *DarknodeRegistryStoreCaller) DarknodeOwner(opts *bind.CallOpts, darknodeID common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DarknodeRegistryStore.contract.Call(opts, out, "darknodeOwner", darknodeID)
	return *ret0, err
}

// DarknodeOwner is a free data retrieval call binding the contract method 0xa3078815.
//
// Solidity: function darknodeOwner(darknodeID address) constant returns(address)
func (_DarknodeRegistryStore *DarknodeRegistryStoreSession) DarknodeOwner(darknodeID common.Address) (common.Address, error) {
	return _DarknodeRegistryStore.Contract.DarknodeOwner(&_DarknodeRegistryStore.CallOpts, darknodeID)
}

// DarknodeOwner is a free data retrieval call binding the contract method 0xa3078815.
//
// Solidity: function darknodeOwner(darknodeID address) constant returns(address)
func (_DarknodeRegistryStore *DarknodeRegistryStoreCallerSession) DarknodeOwner(darknodeID common.Address) (common.Address, error) {
	return _DarknodeRegistryStore.Contract.DarknodeOwner(&_DarknodeRegistryStore.CallOpts, darknodeID)
}

// DarknodePublicKey is a free data retrieval call binding the contract method 0xee594a50.
//
// Solidity: function darknodePublicKey(darknodeID address) constant returns(bytes)
func (_DarknodeRegistryStore *DarknodeRegistryStoreCaller) DarknodePublicKey(opts *bind.CallOpts, darknodeID common.Address) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _DarknodeRegistryStore.contract.Call(opts, out, "darknodePublicKey", darknodeID)
	return *ret0, err
}

// DarknodePublicKey is a free data retrieval call binding the contract method 0xee594a50.
//
// Solidity: function darknodePublicKey(darknodeID address) constant returns(bytes)
func (_DarknodeRegistryStore *DarknodeRegistryStoreSession) DarknodePublicKey(darknodeID common.Address) ([]byte, error) {
	return _DarknodeRegistryStore.Contract.DarknodePublicKey(&_DarknodeRegistryStore.CallOpts, darknodeID)
}

// DarknodePublicKey is a free data retrieval call binding the contract method 0xee594a50.
//
// Solidity: function darknodePublicKey(darknodeID address) constant returns(bytes)
func (_DarknodeRegistryStore *DarknodeRegistryStoreCallerSession) DarknodePublicKey(darknodeID common.Address) ([]byte, error) {
	return _DarknodeRegistryStore.Contract.DarknodePublicKey(&_DarknodeRegistryStore.CallOpts, darknodeID)
}

// DarknodeRegisteredAt is a free data retrieval call binding the contract method 0xe2d7c64c.
//
// Solidity: function darknodeRegisteredAt(darknodeID address) constant returns(uint256)
func (_DarknodeRegistryStore *DarknodeRegistryStoreCaller) DarknodeRegisteredAt(opts *bind.CallOpts, darknodeID common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DarknodeRegistryStore.contract.Call(opts, out, "darknodeRegisteredAt", darknodeID)
	return *ret0, err
}

// DarknodeRegisteredAt is a free data retrieval call binding the contract method 0xe2d7c64c.
//
// Solidity: function darknodeRegisteredAt(darknodeID address) constant returns(uint256)
func (_DarknodeRegistryStore *DarknodeRegistryStoreSession) DarknodeRegisteredAt(darknodeID common.Address) (*big.Int, error) {
	return _DarknodeRegistryStore.Contract.DarknodeRegisteredAt(&_DarknodeRegistryStore.CallOpts, darknodeID)
}

// DarknodeRegisteredAt is a free data retrieval call binding the contract method 0xe2d7c64c.
//
// Solidity: function darknodeRegisteredAt(darknodeID address) constant returns(uint256)
func (_DarknodeRegistryStore *DarknodeRegistryStoreCallerSession) DarknodeRegisteredAt(darknodeID common.Address) (*big.Int, error) {
	return _DarknodeRegistryStore.Contract.DarknodeRegisteredAt(&_DarknodeRegistryStore.CallOpts, darknodeID)
}

// Next is a free data retrieval call binding the contract method 0xab73e316.
//
// Solidity: function next(darknodeID address) constant returns(address)
func (_DarknodeRegistryStore *DarknodeRegistryStoreCaller) Next(opts *bind.CallOpts, darknodeID common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DarknodeRegistryStore.contract.Call(opts, out, "next", darknodeID)
	return *ret0, err
}

// Next is a free data retrieval call binding the contract method 0xab73e316.
//
// Solidity: function next(darknodeID address) constant returns(address)
func (_DarknodeRegistryStore *DarknodeRegistryStoreSession) Next(darknodeID common.Address) (common.Address, error) {
	return _DarknodeRegistryStore.Contract.Next(&_DarknodeRegistryStore.CallOpts, darknodeID)
}

// Next is a free data retrieval call binding the contract method 0xab73e316.
//
// Solidity: function next(darknodeID address) constant returns(address)
func (_DarknodeRegistryStore *DarknodeRegistryStoreCallerSession) Next(darknodeID common.Address) (common.Address, error) {
	return _DarknodeRegistryStore.Contract.Next(&_DarknodeRegistryStore.CallOpts, darknodeID)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_DarknodeRegistryStore *DarknodeRegistryStoreCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DarknodeRegistryStore.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_DarknodeRegistryStore *DarknodeRegistryStoreSession) Owner() (common.Address, error) {
	return _DarknodeRegistryStore.Contract.Owner(&_DarknodeRegistryStore.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_DarknodeRegistryStore *DarknodeRegistryStoreCallerSession) Owner() (common.Address, error) {
	return _DarknodeRegistryStore.Contract.Owner(&_DarknodeRegistryStore.CallOpts)
}

// Ren is a free data retrieval call binding the contract method 0x8a9b4067.
//
// Solidity: function ren() constant returns(address)
func (_DarknodeRegistryStore *DarknodeRegistryStoreCaller) Ren(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DarknodeRegistryStore.contract.Call(opts, out, "ren")
	return *ret0, err
}

// Ren is a free data retrieval call binding the contract method 0x8a9b4067.
//
// Solidity: function ren() constant returns(address)
func (_DarknodeRegistryStore *DarknodeRegistryStoreSession) Ren() (common.Address, error) {
	return _DarknodeRegistryStore.Contract.Ren(&_DarknodeRegistryStore.CallOpts)
}

// Ren is a free data retrieval call binding the contract method 0x8a9b4067.
//
// Solidity: function ren() constant returns(address)
func (_DarknodeRegistryStore *DarknodeRegistryStoreCallerSession) Ren() (common.Address, error) {
	return _DarknodeRegistryStore.Contract.Ren(&_DarknodeRegistryStore.CallOpts)
}

// AppendDarknode is a paid mutator transaction binding the contract method 0xa85ef579.
//
// Solidity: function appendDarknode(_darknodeID address, _darknodeOwner address, _bond uint256, _publicKey bytes, _registeredAt uint256, _deregisteredAt uint256) returns()
func (_DarknodeRegistryStore *DarknodeRegistryStoreTransactor) AppendDarknode(opts *bind.TransactOpts, _darknodeID common.Address, _darknodeOwner common.Address, _bond *big.Int, _publicKey []byte, _registeredAt *big.Int, _deregisteredAt *big.Int) (*types.Transaction, error) {
	return _DarknodeRegistryStore.contract.Transact(opts, "appendDarknode", _darknodeID, _darknodeOwner, _bond, _publicKey, _registeredAt, _deregisteredAt)
}

// AppendDarknode is a paid mutator transaction binding the contract method 0xa85ef579.
//
// Solidity: function appendDarknode(_darknodeID address, _darknodeOwner address, _bond uint256, _publicKey bytes, _registeredAt uint256, _deregisteredAt uint256) returns()
func (_DarknodeRegistryStore *DarknodeRegistryStoreSession) AppendDarknode(_darknodeID common.Address, _darknodeOwner common.Address, _bond *big.Int, _publicKey []byte, _registeredAt *big.Int, _deregisteredAt *big.Int) (*types.Transaction, error) {
	return _DarknodeRegistryStore.Contract.AppendDarknode(&_DarknodeRegistryStore.TransactOpts, _darknodeID, _darknodeOwner, _bond, _publicKey, _registeredAt, _deregisteredAt)
}

// AppendDarknode is a paid mutator transaction binding the contract method 0xa85ef579.
//
// Solidity: function appendDarknode(_darknodeID address, _darknodeOwner address, _bond uint256, _publicKey bytes, _registeredAt uint256, _deregisteredAt uint256) returns()
func (_DarknodeRegistryStore *DarknodeRegistryStoreTransactorSession) AppendDarknode(_darknodeID common.Address, _darknodeOwner common.Address, _bond *big.Int, _publicKey []byte, _registeredAt *big.Int, _deregisteredAt *big.Int) (*types.Transaction, error) {
	return _DarknodeRegistryStore.Contract.AppendDarknode(&_DarknodeRegistryStore.TransactOpts, _darknodeID, _darknodeOwner, _bond, _publicKey, _registeredAt, _deregisteredAt)
}

// RemoveDarknode is a paid mutator transaction binding the contract method 0x41b44392.
//
// Solidity: function removeDarknode(darknodeID address) returns()
func (_DarknodeRegistryStore *DarknodeRegistryStoreTransactor) RemoveDarknode(opts *bind.TransactOpts, darknodeID common.Address) (*types.Transaction, error) {
	return _DarknodeRegistryStore.contract.Transact(opts, "removeDarknode", darknodeID)
}

// RemoveDarknode is a paid mutator transaction binding the contract method 0x41b44392.
//
// Solidity: function removeDarknode(darknodeID address) returns()
func (_DarknodeRegistryStore *DarknodeRegistryStoreSession) RemoveDarknode(darknodeID common.Address) (*types.Transaction, error) {
	return _DarknodeRegistryStore.Contract.RemoveDarknode(&_DarknodeRegistryStore.TransactOpts, darknodeID)
}

// RemoveDarknode is a paid mutator transaction binding the contract method 0x41b44392.
//
// Solidity: function removeDarknode(darknodeID address) returns()
func (_DarknodeRegistryStore *DarknodeRegistryStoreTransactorSession) RemoveDarknode(darknodeID common.Address) (*types.Transaction, error) {
	return _DarknodeRegistryStore.Contract.RemoveDarknode(&_DarknodeRegistryStore.TransactOpts, darknodeID)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DarknodeRegistryStore *DarknodeRegistryStoreTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DarknodeRegistryStore.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DarknodeRegistryStore *DarknodeRegistryStoreSession) RenounceOwnership() (*types.Transaction, error) {
	return _DarknodeRegistryStore.Contract.RenounceOwnership(&_DarknodeRegistryStore.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DarknodeRegistryStore *DarknodeRegistryStoreTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _DarknodeRegistryStore.Contract.RenounceOwnership(&_DarknodeRegistryStore.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_DarknodeRegistryStore *DarknodeRegistryStoreTransactor) TransferOwnership(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _DarknodeRegistryStore.contract.Transact(opts, "transferOwnership", _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_DarknodeRegistryStore *DarknodeRegistryStoreSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _DarknodeRegistryStore.Contract.TransferOwnership(&_DarknodeRegistryStore.TransactOpts, _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_DarknodeRegistryStore *DarknodeRegistryStoreTransactorSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _DarknodeRegistryStore.Contract.TransferOwnership(&_DarknodeRegistryStore.TransactOpts, _newOwner)
}

// UpdateDarknodeBond is a paid mutator transaction binding the contract method 0xfbc402fc.
//
// Solidity: function updateDarknodeBond(darknodeID address, bond uint256) returns()
func (_DarknodeRegistryStore *DarknodeRegistryStoreTransactor) UpdateDarknodeBond(opts *bind.TransactOpts, darknodeID common.Address, bond *big.Int) (*types.Transaction, error) {
	return _DarknodeRegistryStore.contract.Transact(opts, "updateDarknodeBond", darknodeID, bond)
}

// UpdateDarknodeBond is a paid mutator transaction binding the contract method 0xfbc402fc.
//
// Solidity: function updateDarknodeBond(darknodeID address, bond uint256) returns()
func (_DarknodeRegistryStore *DarknodeRegistryStoreSession) UpdateDarknodeBond(darknodeID common.Address, bond *big.Int) (*types.Transaction, error) {
	return _DarknodeRegistryStore.Contract.UpdateDarknodeBond(&_DarknodeRegistryStore.TransactOpts, darknodeID, bond)
}

// UpdateDarknodeBond is a paid mutator transaction binding the contract method 0xfbc402fc.
//
// Solidity: function updateDarknodeBond(darknodeID address, bond uint256) returns()
func (_DarknodeRegistryStore *DarknodeRegistryStoreTransactorSession) UpdateDarknodeBond(darknodeID common.Address, bond *big.Int) (*types.Transaction, error) {
	return _DarknodeRegistryStore.Contract.UpdateDarknodeBond(&_DarknodeRegistryStore.TransactOpts, darknodeID, bond)
}

// UpdateDarknodeDeregisteredAt is a paid mutator transaction binding the contract method 0x3ac39d4b.
//
// Solidity: function updateDarknodeDeregisteredAt(darknodeID address, deregisteredAt uint256) returns()
func (_DarknodeRegistryStore *DarknodeRegistryStoreTransactor) UpdateDarknodeDeregisteredAt(opts *bind.TransactOpts, darknodeID common.Address, deregisteredAt *big.Int) (*types.Transaction, error) {
	return _DarknodeRegistryStore.contract.Transact(opts, "updateDarknodeDeregisteredAt", darknodeID, deregisteredAt)
}

// UpdateDarknodeDeregisteredAt is a paid mutator transaction binding the contract method 0x3ac39d4b.
//
// Solidity: function updateDarknodeDeregisteredAt(darknodeID address, deregisteredAt uint256) returns()
func (_DarknodeRegistryStore *DarknodeRegistryStoreSession) UpdateDarknodeDeregisteredAt(darknodeID common.Address, deregisteredAt *big.Int) (*types.Transaction, error) {
	return _DarknodeRegistryStore.Contract.UpdateDarknodeDeregisteredAt(&_DarknodeRegistryStore.TransactOpts, darknodeID, deregisteredAt)
}

// UpdateDarknodeDeregisteredAt is a paid mutator transaction binding the contract method 0x3ac39d4b.
//
// Solidity: function updateDarknodeDeregisteredAt(darknodeID address, deregisteredAt uint256) returns()
func (_DarknodeRegistryStore *DarknodeRegistryStoreTransactorSession) UpdateDarknodeDeregisteredAt(darknodeID common.Address, deregisteredAt *big.Int) (*types.Transaction, error) {
	return _DarknodeRegistryStore.Contract.UpdateDarknodeDeregisteredAt(&_DarknodeRegistryStore.TransactOpts, darknodeID, deregisteredAt)
}

// DarknodeRegistryStoreOwnershipRenouncedIterator is returned from FilterOwnershipRenounced and is used to iterate over the raw logs and unpacked data for OwnershipRenounced events raised by the DarknodeRegistryStore contract.
type DarknodeRegistryStoreOwnershipRenouncedIterator struct {
	Event *DarknodeRegistryStoreOwnershipRenounced // Event containing the contract specifics and raw log

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
func (it *DarknodeRegistryStoreOwnershipRenouncedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DarknodeRegistryStoreOwnershipRenounced)
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
		it.Event = new(DarknodeRegistryStoreOwnershipRenounced)
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
func (it *DarknodeRegistryStoreOwnershipRenouncedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DarknodeRegistryStoreOwnershipRenouncedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DarknodeRegistryStoreOwnershipRenounced represents a OwnershipRenounced event raised by the DarknodeRegistryStore contract.
type DarknodeRegistryStoreOwnershipRenounced struct {
	PreviousOwner common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipRenounced is a free log retrieval operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_DarknodeRegistryStore *DarknodeRegistryStoreFilterer) FilterOwnershipRenounced(opts *bind.FilterOpts, previousOwner []common.Address) (*DarknodeRegistryStoreOwnershipRenouncedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _DarknodeRegistryStore.contract.FilterLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DarknodeRegistryStoreOwnershipRenouncedIterator{contract: _DarknodeRegistryStore.contract, event: "OwnershipRenounced", logs: logs, sub: sub}, nil
}

// WatchOwnershipRenounced is a free log subscription operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_DarknodeRegistryStore *DarknodeRegistryStoreFilterer) WatchOwnershipRenounced(opts *bind.WatchOpts, sink chan<- *DarknodeRegistryStoreOwnershipRenounced, previousOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _DarknodeRegistryStore.contract.WatchLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DarknodeRegistryStoreOwnershipRenounced)
				if err := _DarknodeRegistryStore.contract.UnpackLog(event, "OwnershipRenounced", log); err != nil {
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

// DarknodeRegistryStoreOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the DarknodeRegistryStore contract.
type DarknodeRegistryStoreOwnershipTransferredIterator struct {
	Event *DarknodeRegistryStoreOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *DarknodeRegistryStoreOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DarknodeRegistryStoreOwnershipTransferred)
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
		it.Event = new(DarknodeRegistryStoreOwnershipTransferred)
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
func (it *DarknodeRegistryStoreOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DarknodeRegistryStoreOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DarknodeRegistryStoreOwnershipTransferred represents a OwnershipTransferred event raised by the DarknodeRegistryStore contract.
type DarknodeRegistryStoreOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_DarknodeRegistryStore *DarknodeRegistryStoreFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*DarknodeRegistryStoreOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DarknodeRegistryStore.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DarknodeRegistryStoreOwnershipTransferredIterator{contract: _DarknodeRegistryStore.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_DarknodeRegistryStore *DarknodeRegistryStoreFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DarknodeRegistryStoreOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DarknodeRegistryStore.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DarknodeRegistryStoreOwnershipTransferred)
				if err := _DarknodeRegistryStore.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// DarknodeRewardVaultABI is the input ABI used to generate the binding from.
const DarknodeRewardVaultABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"darknodeBalances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_darknode\",\"type\":\"address\"},{\"name\":\"_token\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"darknodeRegistry\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newDarknodeRegistry\",\"type\":\"address\"}],\"name\":\"updateDarknodeRegistry\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"ETHEREUM\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_darknode\",\"type\":\"address\"},{\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_VERSION\",\"type\":\"string\"},{\"name\":\"_darknodeRegistry\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"previousDarknodeRegistry\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"nextDarknodeRegistry\",\"type\":\"address\"}],\"name\":\"LogDarknodeRegistryUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"}],\"name\":\"OwnershipRenounced\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"

// DarknodeRewardVaultBin is the compiled bytecode used for deploying new contracts.
const DarknodeRewardVaultBin = `0x608060405234801561001057600080fd5b50604051610af1380380610af183398101604052805160208083015160008054600160a060020a0319163317905591909201805190926100559160019185019061007c565b5060028054600160a060020a031916600160a060020a039290921691909117905550610117565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106100bd57805160ff19168380011785556100ea565b828001600101855582156100ea579182015b828111156100ea5782518255916020019190600101906100cf565b506100f69291506100fa565b5090565b61011491905b808211156100f65760008155600101610100565b90565b6109cb806101266000396000f3006080604052600436106100a35763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166370324b7781146100a8578063715018a6146100e15780638340f549146100f85780638da5cb5b146101155780639e45e0d014610146578063aaff096d1461015b578063f2fde38b1461017c578063f7cdf47c1461019d578063f940e385146101b2578063ffa1ad74146101d9575b600080fd5b3480156100b457600080fd5b506100cf600160a060020a0360043581169060243516610263565b60408051918252519081900360200190f35b3480156100ed57600080fd5b506100f6610280565b005b6100f6600160a060020a03600435811690602435166044356102ec565b34801561012157600080fd5b5061012a61051a565b60408051600160a060020a039092168252519081900360200190f35b34801561015257600080fd5b5061012a610529565b34801561016757600080fd5b506100f6600160a060020a0360043516610538565b34801561018857600080fd5b506100f6600160a060020a03600435166105c6565b3480156101a957600080fd5b5061012a6105e9565b3480156101be57600080fd5b506100f6600160a060020a0360043581169060243516610601565b3480156101e557600080fd5b506101ee610882565b6040805160208082528351818301528351919283929083019185019080838360005b83811015610228578181015183820152602001610210565b50505050905090810190601f1680156102555780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b600360209081526000928352604080842090915290825290205481565b600054600160a060020a0316331461029757600080fd5b60008054604051600160a060020a03909116917ff8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c6482091a26000805473ffffffffffffffffffffffffffffffffffffffff19169055565b600160a060020a03821673eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee141561036d57348114610368576040805160e560020a62461bcd02815260206004820152601360248201527f6d69736d6174636865642074782076616c756500000000000000000000000000604482015290519081900360640190fd5b6104b3565b34156103c3576040805160e560020a62461bcd02815260206004820152601960248201527f756e6578706563746564206574686572207472616e7366657200000000000000604482015290519081900360640190fd5b604080517f23b872dd000000000000000000000000000000000000000000000000000000008152336004820152306024820152604481018390529051600160a060020a038416916323b872dd9160648083019260209291908290030181600087803b15801561043157600080fd5b505af1158015610445573d6000803e3d6000fd5b505050506040513d602081101561045b57600080fd5b505115156104b3576040805160e560020a62461bcd02815260206004820152601560248201527f746f6b656e207472616e73666572206661696c65640000000000000000000000604482015290519081900360640190fd5b600160a060020a038084166000908152600360209081526040808320938616835292905220546104e9908263ffffffff61090f16565b600160a060020a03938416600090815260036020908152604080832095909616825293909352929091209190915550565b600054600160a060020a031681565b600254600160a060020a031681565b600054600160a060020a0316331461054f57600080fd5b60025460408051600160a060020a039283168152918316602083015280517ff9f6dd5c784f63cc27c1079c73574a73485a6c2e7f7e2181c5eb2be8c693cfb79281900390910190a16002805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b600054600160a060020a031633146105dd57600080fd5b6105e681610922565b50565b73eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee81565b600254604080517f1cedf8a3000000000000000000000000000000000000000000000000000000008152600160a060020a038581166004830152915160009384931691631cedf8a391602480830192602092919082900301818787803b15801561066a57600080fd5b505af115801561067e573d6000803e3d6000fd5b505050506040513d602081101561069457600080fd5b50519150600160a060020a03821615156106f8576040805160e560020a62461bcd02815260206004820152601660248201527f696e76616c6964206461726b6e6f6465206f776e657200000000000000000000604482015290519081900360640190fd5b50600160a060020a03838116600090815260036020908152604080832093861680845293909152812080549190559073eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee141561077e57604051600160a060020a0383169082156108fc029083906000818181858888f19350505050158015610778573d6000803e3d6000fd5b5061087c565b82600160a060020a031663a9059cbb83836040518363ffffffff167c01000000000000000000000000000000000000000000000000000000000281526004018083600160a060020a0316600160a060020a0316815260200182815260200192505050602060405180830381600087803b1580156107fa57600080fd5b505af115801561080e573d6000803e3d6000fd5b505050506040513d602081101561082457600080fd5b5051151561087c576040805160e560020a62461bcd02815260206004820152601560248201527f746f6b656e207472616e73666572206661696c65640000000000000000000000604482015290519081900360640190fd5b50505050565b60018054604080516020600284861615610100026000190190941693909304601f810184900484028201840190925281815292918301828280156109075780601f106108dc57610100808354040283529160200191610907565b820191906000526020600020905b8154815290600101906020018083116108ea57829003601f168201915b505050505081565b8181018281101561091c57fe5b92915050565b600160a060020a038116151561093757600080fd5b60008054604051600160a060020a03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a36000805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03929092169190911790555600a165627a7a72305820c7bf0c933ccfd92393207f834b67c799b18c299cd7b03e75c6fe909fde2823180029`

// DeployDarknodeRewardVault deploys a new Ethereum contract, binding an instance of DarknodeRewardVault to it.
func DeployDarknodeRewardVault(auth *bind.TransactOpts, backend bind.ContractBackend, _VERSION string, _darknodeRegistry common.Address) (common.Address, *types.Transaction, *DarknodeRewardVault, error) {
	parsed, err := abi.JSON(strings.NewReader(DarknodeRewardVaultABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(DarknodeRewardVaultBin), backend, _VERSION, _darknodeRegistry)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DarknodeRewardVault{DarknodeRewardVaultCaller: DarknodeRewardVaultCaller{contract: contract}, DarknodeRewardVaultTransactor: DarknodeRewardVaultTransactor{contract: contract}, DarknodeRewardVaultFilterer: DarknodeRewardVaultFilterer{contract: contract}}, nil
}

// DarknodeRewardVault is an auto generated Go binding around an Ethereum contract.
type DarknodeRewardVault struct {
	DarknodeRewardVaultCaller     // Read-only binding to the contract
	DarknodeRewardVaultTransactor // Write-only binding to the contract
	DarknodeRewardVaultFilterer   // Log filterer for contract events
}

// DarknodeRewardVaultCaller is an auto generated read-only Go binding around an Ethereum contract.
type DarknodeRewardVaultCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DarknodeRewardVaultTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DarknodeRewardVaultTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DarknodeRewardVaultFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DarknodeRewardVaultFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DarknodeRewardVaultSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DarknodeRewardVaultSession struct {
	Contract     *DarknodeRewardVault // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// DarknodeRewardVaultCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DarknodeRewardVaultCallerSession struct {
	Contract *DarknodeRewardVaultCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// DarknodeRewardVaultTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DarknodeRewardVaultTransactorSession struct {
	Contract     *DarknodeRewardVaultTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// DarknodeRewardVaultRaw is an auto generated low-level Go binding around an Ethereum contract.
type DarknodeRewardVaultRaw struct {
	Contract *DarknodeRewardVault // Generic contract binding to access the raw methods on
}

// DarknodeRewardVaultCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DarknodeRewardVaultCallerRaw struct {
	Contract *DarknodeRewardVaultCaller // Generic read-only contract binding to access the raw methods on
}

// DarknodeRewardVaultTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DarknodeRewardVaultTransactorRaw struct {
	Contract *DarknodeRewardVaultTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDarknodeRewardVault creates a new instance of DarknodeRewardVault, bound to a specific deployed contract.
func NewDarknodeRewardVault(address common.Address, backend bind.ContractBackend) (*DarknodeRewardVault, error) {
	contract, err := bindDarknodeRewardVault(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DarknodeRewardVault{DarknodeRewardVaultCaller: DarknodeRewardVaultCaller{contract: contract}, DarknodeRewardVaultTransactor: DarknodeRewardVaultTransactor{contract: contract}, DarknodeRewardVaultFilterer: DarknodeRewardVaultFilterer{contract: contract}}, nil
}

// NewDarknodeRewardVaultCaller creates a new read-only instance of DarknodeRewardVault, bound to a specific deployed contract.
func NewDarknodeRewardVaultCaller(address common.Address, caller bind.ContractCaller) (*DarknodeRewardVaultCaller, error) {
	contract, err := bindDarknodeRewardVault(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DarknodeRewardVaultCaller{contract: contract}, nil
}

// NewDarknodeRewardVaultTransactor creates a new write-only instance of DarknodeRewardVault, bound to a specific deployed contract.
func NewDarknodeRewardVaultTransactor(address common.Address, transactor bind.ContractTransactor) (*DarknodeRewardVaultTransactor, error) {
	contract, err := bindDarknodeRewardVault(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DarknodeRewardVaultTransactor{contract: contract}, nil
}

// NewDarknodeRewardVaultFilterer creates a new log filterer instance of DarknodeRewardVault, bound to a specific deployed contract.
func NewDarknodeRewardVaultFilterer(address common.Address, filterer bind.ContractFilterer) (*DarknodeRewardVaultFilterer, error) {
	contract, err := bindDarknodeRewardVault(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DarknodeRewardVaultFilterer{contract: contract}, nil
}

// bindDarknodeRewardVault binds a generic wrapper to an already deployed contract.
func bindDarknodeRewardVault(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DarknodeRewardVaultABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DarknodeRewardVault *DarknodeRewardVaultRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _DarknodeRewardVault.Contract.DarknodeRewardVaultCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DarknodeRewardVault *DarknodeRewardVaultRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DarknodeRewardVault.Contract.DarknodeRewardVaultTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DarknodeRewardVault *DarknodeRewardVaultRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DarknodeRewardVault.Contract.DarknodeRewardVaultTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DarknodeRewardVault *DarknodeRewardVaultCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _DarknodeRewardVault.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DarknodeRewardVault *DarknodeRewardVaultTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DarknodeRewardVault.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DarknodeRewardVault *DarknodeRewardVaultTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DarknodeRewardVault.Contract.contract.Transact(opts, method, params...)
}

// ETHEREUM is a free data retrieval call binding the contract method 0xf7cdf47c.
//
// Solidity: function ETHEREUM() constant returns(address)
func (_DarknodeRewardVault *DarknodeRewardVaultCaller) ETHEREUM(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DarknodeRewardVault.contract.Call(opts, out, "ETHEREUM")
	return *ret0, err
}

// ETHEREUM is a free data retrieval call binding the contract method 0xf7cdf47c.
//
// Solidity: function ETHEREUM() constant returns(address)
func (_DarknodeRewardVault *DarknodeRewardVaultSession) ETHEREUM() (common.Address, error) {
	return _DarknodeRewardVault.Contract.ETHEREUM(&_DarknodeRewardVault.CallOpts)
}

// ETHEREUM is a free data retrieval call binding the contract method 0xf7cdf47c.
//
// Solidity: function ETHEREUM() constant returns(address)
func (_DarknodeRewardVault *DarknodeRewardVaultCallerSession) ETHEREUM() (common.Address, error) {
	return _DarknodeRewardVault.Contract.ETHEREUM(&_DarknodeRewardVault.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_DarknodeRewardVault *DarknodeRewardVaultCaller) VERSION(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _DarknodeRewardVault.contract.Call(opts, out, "VERSION")
	return *ret0, err
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_DarknodeRewardVault *DarknodeRewardVaultSession) VERSION() (string, error) {
	return _DarknodeRewardVault.Contract.VERSION(&_DarknodeRewardVault.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_DarknodeRewardVault *DarknodeRewardVaultCallerSession) VERSION() (string, error) {
	return _DarknodeRewardVault.Contract.VERSION(&_DarknodeRewardVault.CallOpts)
}

// DarknodeBalances is a free data retrieval call binding the contract method 0x70324b77.
//
// Solidity: function darknodeBalances( address,  address) constant returns(uint256)
func (_DarknodeRewardVault *DarknodeRewardVaultCaller) DarknodeBalances(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DarknodeRewardVault.contract.Call(opts, out, "darknodeBalances", arg0, arg1)
	return *ret0, err
}

// DarknodeBalances is a free data retrieval call binding the contract method 0x70324b77.
//
// Solidity: function darknodeBalances( address,  address) constant returns(uint256)
func (_DarknodeRewardVault *DarknodeRewardVaultSession) DarknodeBalances(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _DarknodeRewardVault.Contract.DarknodeBalances(&_DarknodeRewardVault.CallOpts, arg0, arg1)
}

// DarknodeBalances is a free data retrieval call binding the contract method 0x70324b77.
//
// Solidity: function darknodeBalances( address,  address) constant returns(uint256)
func (_DarknodeRewardVault *DarknodeRewardVaultCallerSession) DarknodeBalances(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _DarknodeRewardVault.Contract.DarknodeBalances(&_DarknodeRewardVault.CallOpts, arg0, arg1)
}

// DarknodeRegistry is a free data retrieval call binding the contract method 0x9e45e0d0.
//
// Solidity: function darknodeRegistry() constant returns(address)
func (_DarknodeRewardVault *DarknodeRewardVaultCaller) DarknodeRegistry(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DarknodeRewardVault.contract.Call(opts, out, "darknodeRegistry")
	return *ret0, err
}

// DarknodeRegistry is a free data retrieval call binding the contract method 0x9e45e0d0.
//
// Solidity: function darknodeRegistry() constant returns(address)
func (_DarknodeRewardVault *DarknodeRewardVaultSession) DarknodeRegistry() (common.Address, error) {
	return _DarknodeRewardVault.Contract.DarknodeRegistry(&_DarknodeRewardVault.CallOpts)
}

// DarknodeRegistry is a free data retrieval call binding the contract method 0x9e45e0d0.
//
// Solidity: function darknodeRegistry() constant returns(address)
func (_DarknodeRewardVault *DarknodeRewardVaultCallerSession) DarknodeRegistry() (common.Address, error) {
	return _DarknodeRewardVault.Contract.DarknodeRegistry(&_DarknodeRewardVault.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_DarknodeRewardVault *DarknodeRewardVaultCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DarknodeRewardVault.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_DarknodeRewardVault *DarknodeRewardVaultSession) Owner() (common.Address, error) {
	return _DarknodeRewardVault.Contract.Owner(&_DarknodeRewardVault.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_DarknodeRewardVault *DarknodeRewardVaultCallerSession) Owner() (common.Address, error) {
	return _DarknodeRewardVault.Contract.Owner(&_DarknodeRewardVault.CallOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0x8340f549.
//
// Solidity: function deposit(_darknode address, _token address, _value uint256) returns()
func (_DarknodeRewardVault *DarknodeRewardVaultTransactor) Deposit(opts *bind.TransactOpts, _darknode common.Address, _token common.Address, _value *big.Int) (*types.Transaction, error) {
	return _DarknodeRewardVault.contract.Transact(opts, "deposit", _darknode, _token, _value)
}

// Deposit is a paid mutator transaction binding the contract method 0x8340f549.
//
// Solidity: function deposit(_darknode address, _token address, _value uint256) returns()
func (_DarknodeRewardVault *DarknodeRewardVaultSession) Deposit(_darknode common.Address, _token common.Address, _value *big.Int) (*types.Transaction, error) {
	return _DarknodeRewardVault.Contract.Deposit(&_DarknodeRewardVault.TransactOpts, _darknode, _token, _value)
}

// Deposit is a paid mutator transaction binding the contract method 0x8340f549.
//
// Solidity: function deposit(_darknode address, _token address, _value uint256) returns()
func (_DarknodeRewardVault *DarknodeRewardVaultTransactorSession) Deposit(_darknode common.Address, _token common.Address, _value *big.Int) (*types.Transaction, error) {
	return _DarknodeRewardVault.Contract.Deposit(&_DarknodeRewardVault.TransactOpts, _darknode, _token, _value)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DarknodeRewardVault *DarknodeRewardVaultTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DarknodeRewardVault.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DarknodeRewardVault *DarknodeRewardVaultSession) RenounceOwnership() (*types.Transaction, error) {
	return _DarknodeRewardVault.Contract.RenounceOwnership(&_DarknodeRewardVault.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DarknodeRewardVault *DarknodeRewardVaultTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _DarknodeRewardVault.Contract.RenounceOwnership(&_DarknodeRewardVault.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_DarknodeRewardVault *DarknodeRewardVaultTransactor) TransferOwnership(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _DarknodeRewardVault.contract.Transact(opts, "transferOwnership", _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_DarknodeRewardVault *DarknodeRewardVaultSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _DarknodeRewardVault.Contract.TransferOwnership(&_DarknodeRewardVault.TransactOpts, _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_DarknodeRewardVault *DarknodeRewardVaultTransactorSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _DarknodeRewardVault.Contract.TransferOwnership(&_DarknodeRewardVault.TransactOpts, _newOwner)
}

// UpdateDarknodeRegistry is a paid mutator transaction binding the contract method 0xaaff096d.
//
// Solidity: function updateDarknodeRegistry(_newDarknodeRegistry address) returns()
func (_DarknodeRewardVault *DarknodeRewardVaultTransactor) UpdateDarknodeRegistry(opts *bind.TransactOpts, _newDarknodeRegistry common.Address) (*types.Transaction, error) {
	return _DarknodeRewardVault.contract.Transact(opts, "updateDarknodeRegistry", _newDarknodeRegistry)
}

// UpdateDarknodeRegistry is a paid mutator transaction binding the contract method 0xaaff096d.
//
// Solidity: function updateDarknodeRegistry(_newDarknodeRegistry address) returns()
func (_DarknodeRewardVault *DarknodeRewardVaultSession) UpdateDarknodeRegistry(_newDarknodeRegistry common.Address) (*types.Transaction, error) {
	return _DarknodeRewardVault.Contract.UpdateDarknodeRegistry(&_DarknodeRewardVault.TransactOpts, _newDarknodeRegistry)
}

// UpdateDarknodeRegistry is a paid mutator transaction binding the contract method 0xaaff096d.
//
// Solidity: function updateDarknodeRegistry(_newDarknodeRegistry address) returns()
func (_DarknodeRewardVault *DarknodeRewardVaultTransactorSession) UpdateDarknodeRegistry(_newDarknodeRegistry common.Address) (*types.Transaction, error) {
	return _DarknodeRewardVault.Contract.UpdateDarknodeRegistry(&_DarknodeRewardVault.TransactOpts, _newDarknodeRegistry)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf940e385.
//
// Solidity: function withdraw(_darknode address, _token address) returns()
func (_DarknodeRewardVault *DarknodeRewardVaultTransactor) Withdraw(opts *bind.TransactOpts, _darknode common.Address, _token common.Address) (*types.Transaction, error) {
	return _DarknodeRewardVault.contract.Transact(opts, "withdraw", _darknode, _token)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf940e385.
//
// Solidity: function withdraw(_darknode address, _token address) returns()
func (_DarknodeRewardVault *DarknodeRewardVaultSession) Withdraw(_darknode common.Address, _token common.Address) (*types.Transaction, error) {
	return _DarknodeRewardVault.Contract.Withdraw(&_DarknodeRewardVault.TransactOpts, _darknode, _token)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf940e385.
//
// Solidity: function withdraw(_darknode address, _token address) returns()
func (_DarknodeRewardVault *DarknodeRewardVaultTransactorSession) Withdraw(_darknode common.Address, _token common.Address) (*types.Transaction, error) {
	return _DarknodeRewardVault.Contract.Withdraw(&_DarknodeRewardVault.TransactOpts, _darknode, _token)
}

// DarknodeRewardVaultLogDarknodeRegistryUpdatedIterator is returned from FilterLogDarknodeRegistryUpdated and is used to iterate over the raw logs and unpacked data for LogDarknodeRegistryUpdated events raised by the DarknodeRewardVault contract.
type DarknodeRewardVaultLogDarknodeRegistryUpdatedIterator struct {
	Event *DarknodeRewardVaultLogDarknodeRegistryUpdated // Event containing the contract specifics and raw log

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
func (it *DarknodeRewardVaultLogDarknodeRegistryUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DarknodeRewardVaultLogDarknodeRegistryUpdated)
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
		it.Event = new(DarknodeRewardVaultLogDarknodeRegistryUpdated)
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
func (it *DarknodeRewardVaultLogDarknodeRegistryUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DarknodeRewardVaultLogDarknodeRegistryUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DarknodeRewardVaultLogDarknodeRegistryUpdated represents a LogDarknodeRegistryUpdated event raised by the DarknodeRewardVault contract.
type DarknodeRewardVaultLogDarknodeRegistryUpdated struct {
	PreviousDarknodeRegistry common.Address
	NextDarknodeRegistry     common.Address
	Raw                      types.Log // Blockchain specific contextual infos
}

// FilterLogDarknodeRegistryUpdated is a free log retrieval operation binding the contract event 0xf9f6dd5c784f63cc27c1079c73574a73485a6c2e7f7e2181c5eb2be8c693cfb7.
//
// Solidity: e LogDarknodeRegistryUpdated(previousDarknodeRegistry address, nextDarknodeRegistry address)
func (_DarknodeRewardVault *DarknodeRewardVaultFilterer) FilterLogDarknodeRegistryUpdated(opts *bind.FilterOpts) (*DarknodeRewardVaultLogDarknodeRegistryUpdatedIterator, error) {

	logs, sub, err := _DarknodeRewardVault.contract.FilterLogs(opts, "LogDarknodeRegistryUpdated")
	if err != nil {
		return nil, err
	}
	return &DarknodeRewardVaultLogDarknodeRegistryUpdatedIterator{contract: _DarknodeRewardVault.contract, event: "LogDarknodeRegistryUpdated", logs: logs, sub: sub}, nil
}

// WatchLogDarknodeRegistryUpdated is a free log subscription operation binding the contract event 0xf9f6dd5c784f63cc27c1079c73574a73485a6c2e7f7e2181c5eb2be8c693cfb7.
//
// Solidity: e LogDarknodeRegistryUpdated(previousDarknodeRegistry address, nextDarknodeRegistry address)
func (_DarknodeRewardVault *DarknodeRewardVaultFilterer) WatchLogDarknodeRegistryUpdated(opts *bind.WatchOpts, sink chan<- *DarknodeRewardVaultLogDarknodeRegistryUpdated) (event.Subscription, error) {

	logs, sub, err := _DarknodeRewardVault.contract.WatchLogs(opts, "LogDarknodeRegistryUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DarknodeRewardVaultLogDarknodeRegistryUpdated)
				if err := _DarknodeRewardVault.contract.UnpackLog(event, "LogDarknodeRegistryUpdated", log); err != nil {
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

// DarknodeRewardVaultOwnershipRenouncedIterator is returned from FilterOwnershipRenounced and is used to iterate over the raw logs and unpacked data for OwnershipRenounced events raised by the DarknodeRewardVault contract.
type DarknodeRewardVaultOwnershipRenouncedIterator struct {
	Event *DarknodeRewardVaultOwnershipRenounced // Event containing the contract specifics and raw log

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
func (it *DarknodeRewardVaultOwnershipRenouncedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DarknodeRewardVaultOwnershipRenounced)
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
		it.Event = new(DarknodeRewardVaultOwnershipRenounced)
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
func (it *DarknodeRewardVaultOwnershipRenouncedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DarknodeRewardVaultOwnershipRenouncedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DarknodeRewardVaultOwnershipRenounced represents a OwnershipRenounced event raised by the DarknodeRewardVault contract.
type DarknodeRewardVaultOwnershipRenounced struct {
	PreviousOwner common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipRenounced is a free log retrieval operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_DarknodeRewardVault *DarknodeRewardVaultFilterer) FilterOwnershipRenounced(opts *bind.FilterOpts, previousOwner []common.Address) (*DarknodeRewardVaultOwnershipRenouncedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _DarknodeRewardVault.contract.FilterLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DarknodeRewardVaultOwnershipRenouncedIterator{contract: _DarknodeRewardVault.contract, event: "OwnershipRenounced", logs: logs, sub: sub}, nil
}

// WatchOwnershipRenounced is a free log subscription operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_DarknodeRewardVault *DarknodeRewardVaultFilterer) WatchOwnershipRenounced(opts *bind.WatchOpts, sink chan<- *DarknodeRewardVaultOwnershipRenounced, previousOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _DarknodeRewardVault.contract.WatchLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DarknodeRewardVaultOwnershipRenounced)
				if err := _DarknodeRewardVault.contract.UnpackLog(event, "OwnershipRenounced", log); err != nil {
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

// DarknodeRewardVaultOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the DarknodeRewardVault contract.
type DarknodeRewardVaultOwnershipTransferredIterator struct {
	Event *DarknodeRewardVaultOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *DarknodeRewardVaultOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DarknodeRewardVaultOwnershipTransferred)
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
		it.Event = new(DarknodeRewardVaultOwnershipTransferred)
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
func (it *DarknodeRewardVaultOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DarknodeRewardVaultOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DarknodeRewardVaultOwnershipTransferred represents a OwnershipTransferred event raised by the DarknodeRewardVault contract.
type DarknodeRewardVaultOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_DarknodeRewardVault *DarknodeRewardVaultFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*DarknodeRewardVaultOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DarknodeRewardVault.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DarknodeRewardVaultOwnershipTransferredIterator{contract: _DarknodeRewardVault.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_DarknodeRewardVault *DarknodeRewardVaultFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DarknodeRewardVaultOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DarknodeRewardVault.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DarknodeRewardVaultOwnershipTransferred)
				if err := _DarknodeRewardVault.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// DarknodeSlasherABI is the input ABI used to generate the binding from.
const DarknodeSlasherABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"trustedOrderbook\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"orderSubmitted\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_buyID\",\"type\":\"bytes32\"},{\"name\":\"_sellID\",\"type\":\"bytes32\"}],\"name\":\"submitChallenge\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"challengeSubmitted\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"orderDetails\",\"outputs\":[{\"name\":\"settlementID\",\"type\":\"uint64\"},{\"name\":\"tokens\",\"type\":\"uint64\"},{\"name\":\"price\",\"type\":\"uint256\"},{\"name\":\"volume\",\"type\":\"uint256\"},{\"name\":\"minimumVolume\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"details\",\"type\":\"bytes\"},{\"name\":\"settlementID\",\"type\":\"uint64\"},{\"name\":\"tokens\",\"type\":\"uint64\"},{\"name\":\"price\",\"type\":\"uint256\"},{\"name\":\"volume\",\"type\":\"uint256\"},{\"name\":\"minimumVolume\",\"type\":\"uint256\"}],\"name\":\"submitChallengeOrder\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"challengers\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"trustedDarknodeRegistry\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_VERSION\",\"type\":\"string\"},{\"name\":\"_darknodeRegistry\",\"type\":\"address\"},{\"name\":\"_orderbook\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"}],\"name\":\"OwnershipRenounced\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"

// DarknodeSlasherBin is the compiled bytecode used for deploying new contracts.
const DarknodeSlasherBin = `0x60806040523480156200001157600080fd5b506040516200131a3803806200131a83398101604090815281516020808401519284015160008054600160a060020a031916331790559190930180519093620000609160019186019062000094565b5060028054600160a060020a03938416600160a060020a031991821617909155600380549290931691161790555062000139565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10620000d757805160ff191683800117855562000107565b8280016001018555821562000107579182015b8281111562000107578251825591602001919060010190620000ea565b506200011592915062000119565b5090565b6200013691905b8082111562000115576000815560010162000120565b90565b6111d180620001496000396000f3006080604052600436106100b95763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166326cf660d81146100be5780634778a5be146100ef5780634d4582a11461011b578063715018a6146101385780638c8b6fc51461014d5780638da5cb5b14610168578063a3c50b321461017d578063c41c33af146101cd578063f2fde38b1461020b578063fc322d731461022c578063fc722b2b14610244578063ffa1ad7414610259575b600080fd5b3480156100ca57600080fd5b506100d36102e3565b60408051600160a060020a039092168252519081900360200190f35b3480156100fb57600080fd5b506101076004356102f2565b604080519115158252519081900360200190f35b34801561012757600080fd5b50610136600435602435610307565b005b34801561014457600080fd5b5061013661099b565b34801561015957600080fd5b50610107600435602435610a07565b34801561017457600080fd5b506100d3610a27565b34801561018957600080fd5b50610195600435610a36565b6040805167ffffffffffffffff9687168152949095166020850152838501929092526060830152608082015290519081900360a00190f35b3480156101d957600080fd5b50610136602460048035828101929101359067ffffffffffffffff90358116906044351660643560843560a435610a76565b34801561021757600080fd5b50610136600160a060020a0360043516610d9e565b34801561023857600080fd5b506100d3600435610dc1565b34801561025057600080fd5b506100d3610ddc565b34801561026557600080fd5b5061026e610deb565b6040805160208082528351818301528351919283929083019185019080838360005b838110156102a8578181015183820152602001610290565b50505050905090810190601f1680156102d55780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b600354600160a060020a031681565b60046020526000908152604090205460ff1681565b60008281526005602090815260408083208484529091528120548190819060ff161561037d576040805160e560020a62461bcd02815260206004820152601260248201527f616c7265616479206368616c6c656e6765640000000000000000000000000000604482015290519081900360640190fd5b60008581526004602052604090205460ff1615156103e5576040805160e560020a62461bcd02815260206004820152601360248201527f64657461696c7320756e617661696c61626c6500000000000000000000000000604482015290519081900360640190fd5b60008481526004602052604090205460ff16151561044d576040805160e560020a62461bcd02815260206004820152601360248201527f64657461696c7320756e617661696c61626c6500000000000000000000000000604482015290519081900360640190fd5b600354604080517faf3e8a400000000000000000000000000000000000000000000000000000000081526004810188905290518692600160a060020a03169163af3e8a409160248083019260209291908290030181600087803b1580156104b357600080fd5b505af11580156104c7573d6000803e3d6000fd5b505050506040513d60208110156104dd57600080fd5b505114610534576040805160e560020a62461bcd02815260206004820152601260248201527f756e636f6e6669726d6564206f72646572730000000000000000000000000000604482015290519081900360640190fd5b61068e60066000876000191660001916815260200190815260200160002060a060405190810160405290816000820160009054906101000a900467ffffffffffffffff1667ffffffffffffffff1667ffffffffffffffff1681526020016000820160089054906101000a900467ffffffffffffffff1667ffffffffffffffff1667ffffffffffffffff168152602001600182015481526020016002820154815260200160038201548152505060066000876000191660001916815260200190815260200160002060a060405190810160405290816000820160009054906101000a900467ffffffffffffffff1667ffffffffffffffff1667ffffffffffffffff1681526020016000820160089054906101000a900467ffffffffffffffff1667ffffffffffffffff1667ffffffffffffffff1681526020016001820154815260200160028201548152602001600382015481525050610e78565b600354604080517fb1a0801000000000000000000000000000000000000000000000000000000000815260048101889052905192159550600160a060020a039091169163b1a08010916024808201926020929091908290030181600087803b1580156106f957600080fd5b505af115801561070d573d6000803e3d6000fd5b505050506040513d602081101561072357600080fd5b5051600354604080517fb1a08010000000000000000000000000000000000000000000000000000000008152600481018990529051600160a060020a03938416939092169163b1a08010916024808201926020929091908290030181600087803b15801561079057600080fd5b505af11580156107a4573d6000803e3d6000fd5b505050506040513d60208110156107ba57600080fd5b5051600160a060020a031614915082806107d15750815b1515610827576040805160e560020a62461bcd02815260206004820152601160248201527f696e76616c6964206368616c6c656e6765000000000000000000000000000000604482015290519081900360640190fd5b600354604080517f1107c3f7000000000000000000000000000000000000000000000000000000008152600481018890529051600160a060020a0390921691631107c3f7916024808201926020929091908290030181600087803b15801561088e57600080fd5b505af11580156108a2573d6000803e3d6000fd5b505050506040513d60208110156108b857600080fd5b5051600086815260056020818152604080842089855282528084208054600160ff1991821681179092559383528185208b865283528185208054909416179092556002546007909152818320548884528284205483517f563bf264000000000000000000000000000000000000000000000000000000008152600160a060020a038088166004830152928316602482015290821660448201529251949550169263563bf2649260648084019391929182900301818387803b15801561097c57600080fd5b505af1158015610990573d6000803e3d6000fd5b505050505050505050565b600054600160a060020a031633146109b257600080fd5b60008054604051600160a060020a03909116917ff8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c6482091a26000805473ffffffffffffffffffffffffffffffffffffffff19169055565b600560209081526000928352604080842090915290825290205460ff1681565b600054600160a060020a031681565b600660205260009081526040902080546001820154600283015460039093015467ffffffffffffffff808416946801000000000000000090940416929085565b610a7e611161565b600254604080517fc3c5a5470000000000000000000000000000000000000000000000000000000081523360048201529051600092600160a060020a03169163c3c5a54791602480830192602092919082900301818787803b158015610ae357600080fd5b505af1158015610af7573d6000803e3d6000fd5b505050506040513d6020811015610b0d57600080fd5b505180610ba85750600254604080517f8020fc1f0000000000000000000000000000000000000000000000000000000081523360048201529051600160a060020a0390921691638020fc1f916024808201926020929091908290030181600087803b158015610b7b57600080fd5b505af1158015610b8f573d6000803e3d6000fd5b505050506040513d6020811015610ba557600080fd5b50515b1515610bfe576040805160e560020a62461bcd02815260206004820152601060248201527f6d757374206265206461726b6e6f646500000000000000000000000000000000604482015290519081900360640190fd5b60a0604051908101604052808867ffffffffffffffff1681526020018767ffffffffffffffff168152602001868152602001858152602001848152509150610c7789898080601f01602080910402602001604051908101604052809392919081815260200183838082843750889450610f0b9350505050565b60008181526004602052604090205490915060ff1615610ce1576040805160e560020a62461bcd02815260206004820152601160248201527f616c7265616479207375626d6974746564000000000000000000000000000000604482015290519081900360640190fd5b6000908152600660209081526040808320845181548487015167ffffffffffffffff90811668010000000000000000026fffffffffffffffff0000000000000000199190931667ffffffffffffffff199092169190911716178155818501516001808301919091556060860151600283015560809095015160039091015560078252808320805473ffffffffffffffffffffffffffffffffffffffff19163317905560049091529020805460ff1916909117905550505050505050565b600054600160a060020a03163314610db557600080fd5b610dbe81611070565b50565b600760205260009081526040902054600160a060020a031681565b600254600160a060020a031681565b60018054604080516020600284861615610100026000190190941693909304601f81018490048402820184019092528181529291830182828015610e705780601f10610e4557610100808354040283529160200191610e70565b820191906000526020600020905b815481529060010190602001808311610e5357829003601f168201915b505050505081565b6000610e8c836020015183602001516110ed565b1515610e9a57506000610f05565b816040015183604001511015610eb257506000610f05565b816080015183606001511015610eca57506000610f05565b826080015182606001511015610ee257506000610f05565b8151835167ffffffffffffffff908116911614610f0157506000610f05565b5060015b92915050565b600082826000015183602001518460400151856060015186608001516040516020018087805190602001908083835b60208310610f595780518252601f199092019160209182019101610f3a565b6001836020036101000a0380198251168184511680821785525050505050509050018667ffffffffffffffff1667ffffffffffffffff1678010000000000000000000000000000000000000000000000000281526008018567ffffffffffffffff1667ffffffffffffffff16780100000000000000000000000000000000000000000000000002815260080184815260200183815260200182815260200196505050505050506040516020818303038152906040526040518082805190602001908083835b6020831061103d5780518252601f19909201916020918201910161101e565b5181516020939093036101000a600019018019909116921691909117905260405192018290039091209695505050505050565b600160a060020a038116151561108557600080fd5b60008054604051600160a060020a03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a36000805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b600064010000000067ffffffffffffffff83160463ffffffff908116908416148015611133575064010000000067ffffffffffffffff84160463ffffffff908116908316145b801561115a575063ffffffff80841664010000000067ffffffffffffffff86160490911611155b9392505050565b60a060405190810160405280600067ffffffffffffffff168152602001600067ffffffffffffffff16815260200160008152602001600081526020016000815250905600a165627a7a7230582056d93db5bda848f78fec922018e917ac51c2942c5728619cec54549b81e3c0650029`

// DeployDarknodeSlasher deploys a new Ethereum contract, binding an instance of DarknodeSlasher to it.
func DeployDarknodeSlasher(auth *bind.TransactOpts, backend bind.ContractBackend, _VERSION string, _darknodeRegistry common.Address, _orderbook common.Address) (common.Address, *types.Transaction, *DarknodeSlasher, error) {
	parsed, err := abi.JSON(strings.NewReader(DarknodeSlasherABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(DarknodeSlasherBin), backend, _VERSION, _darknodeRegistry, _orderbook)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DarknodeSlasher{DarknodeSlasherCaller: DarknodeSlasherCaller{contract: contract}, DarknodeSlasherTransactor: DarknodeSlasherTransactor{contract: contract}, DarknodeSlasherFilterer: DarknodeSlasherFilterer{contract: contract}}, nil
}

// DarknodeSlasher is an auto generated Go binding around an Ethereum contract.
type DarknodeSlasher struct {
	DarknodeSlasherCaller     // Read-only binding to the contract
	DarknodeSlasherTransactor // Write-only binding to the contract
	DarknodeSlasherFilterer   // Log filterer for contract events
}

// DarknodeSlasherCaller is an auto generated read-only Go binding around an Ethereum contract.
type DarknodeSlasherCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DarknodeSlasherTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DarknodeSlasherTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DarknodeSlasherFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DarknodeSlasherFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DarknodeSlasherSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DarknodeSlasherSession struct {
	Contract     *DarknodeSlasher  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DarknodeSlasherCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DarknodeSlasherCallerSession struct {
	Contract *DarknodeSlasherCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// DarknodeSlasherTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DarknodeSlasherTransactorSession struct {
	Contract     *DarknodeSlasherTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// DarknodeSlasherRaw is an auto generated low-level Go binding around an Ethereum contract.
type DarknodeSlasherRaw struct {
	Contract *DarknodeSlasher // Generic contract binding to access the raw methods on
}

// DarknodeSlasherCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DarknodeSlasherCallerRaw struct {
	Contract *DarknodeSlasherCaller // Generic read-only contract binding to access the raw methods on
}

// DarknodeSlasherTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DarknodeSlasherTransactorRaw struct {
	Contract *DarknodeSlasherTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDarknodeSlasher creates a new instance of DarknodeSlasher, bound to a specific deployed contract.
func NewDarknodeSlasher(address common.Address, backend bind.ContractBackend) (*DarknodeSlasher, error) {
	contract, err := bindDarknodeSlasher(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DarknodeSlasher{DarknodeSlasherCaller: DarknodeSlasherCaller{contract: contract}, DarknodeSlasherTransactor: DarknodeSlasherTransactor{contract: contract}, DarknodeSlasherFilterer: DarknodeSlasherFilterer{contract: contract}}, nil
}

// NewDarknodeSlasherCaller creates a new read-only instance of DarknodeSlasher, bound to a specific deployed contract.
func NewDarknodeSlasherCaller(address common.Address, caller bind.ContractCaller) (*DarknodeSlasherCaller, error) {
	contract, err := bindDarknodeSlasher(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DarknodeSlasherCaller{contract: contract}, nil
}

// NewDarknodeSlasherTransactor creates a new write-only instance of DarknodeSlasher, bound to a specific deployed contract.
func NewDarknodeSlasherTransactor(address common.Address, transactor bind.ContractTransactor) (*DarknodeSlasherTransactor, error) {
	contract, err := bindDarknodeSlasher(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DarknodeSlasherTransactor{contract: contract}, nil
}

// NewDarknodeSlasherFilterer creates a new log filterer instance of DarknodeSlasher, bound to a specific deployed contract.
func NewDarknodeSlasherFilterer(address common.Address, filterer bind.ContractFilterer) (*DarknodeSlasherFilterer, error) {
	contract, err := bindDarknodeSlasher(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DarknodeSlasherFilterer{contract: contract}, nil
}

// bindDarknodeSlasher binds a generic wrapper to an already deployed contract.
func bindDarknodeSlasher(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DarknodeSlasherABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DarknodeSlasher *DarknodeSlasherRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _DarknodeSlasher.Contract.DarknodeSlasherCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DarknodeSlasher *DarknodeSlasherRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DarknodeSlasher.Contract.DarknodeSlasherTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DarknodeSlasher *DarknodeSlasherRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DarknodeSlasher.Contract.DarknodeSlasherTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DarknodeSlasher *DarknodeSlasherCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _DarknodeSlasher.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DarknodeSlasher *DarknodeSlasherTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DarknodeSlasher.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DarknodeSlasher *DarknodeSlasherTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DarknodeSlasher.Contract.contract.Transact(opts, method, params...)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_DarknodeSlasher *DarknodeSlasherCaller) VERSION(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _DarknodeSlasher.contract.Call(opts, out, "VERSION")
	return *ret0, err
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_DarknodeSlasher *DarknodeSlasherSession) VERSION() (string, error) {
	return _DarknodeSlasher.Contract.VERSION(&_DarknodeSlasher.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_DarknodeSlasher *DarknodeSlasherCallerSession) VERSION() (string, error) {
	return _DarknodeSlasher.Contract.VERSION(&_DarknodeSlasher.CallOpts)
}

// ChallengeSubmitted is a free data retrieval call binding the contract method 0x8c8b6fc5.
//
// Solidity: function challengeSubmitted( bytes32,  bytes32) constant returns(bool)
func (_DarknodeSlasher *DarknodeSlasherCaller) ChallengeSubmitted(opts *bind.CallOpts, arg0 [32]byte, arg1 [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _DarknodeSlasher.contract.Call(opts, out, "challengeSubmitted", arg0, arg1)
	return *ret0, err
}

// ChallengeSubmitted is a free data retrieval call binding the contract method 0x8c8b6fc5.
//
// Solidity: function challengeSubmitted( bytes32,  bytes32) constant returns(bool)
func (_DarknodeSlasher *DarknodeSlasherSession) ChallengeSubmitted(arg0 [32]byte, arg1 [32]byte) (bool, error) {
	return _DarknodeSlasher.Contract.ChallengeSubmitted(&_DarknodeSlasher.CallOpts, arg0, arg1)
}

// ChallengeSubmitted is a free data retrieval call binding the contract method 0x8c8b6fc5.
//
// Solidity: function challengeSubmitted( bytes32,  bytes32) constant returns(bool)
func (_DarknodeSlasher *DarknodeSlasherCallerSession) ChallengeSubmitted(arg0 [32]byte, arg1 [32]byte) (bool, error) {
	return _DarknodeSlasher.Contract.ChallengeSubmitted(&_DarknodeSlasher.CallOpts, arg0, arg1)
}

// Challengers is a free data retrieval call binding the contract method 0xfc322d73.
//
// Solidity: function challengers( bytes32) constant returns(address)
func (_DarknodeSlasher *DarknodeSlasherCaller) Challengers(opts *bind.CallOpts, arg0 [32]byte) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DarknodeSlasher.contract.Call(opts, out, "challengers", arg0)
	return *ret0, err
}

// Challengers is a free data retrieval call binding the contract method 0xfc322d73.
//
// Solidity: function challengers( bytes32) constant returns(address)
func (_DarknodeSlasher *DarknodeSlasherSession) Challengers(arg0 [32]byte) (common.Address, error) {
	return _DarknodeSlasher.Contract.Challengers(&_DarknodeSlasher.CallOpts, arg0)
}

// Challengers is a free data retrieval call binding the contract method 0xfc322d73.
//
// Solidity: function challengers( bytes32) constant returns(address)
func (_DarknodeSlasher *DarknodeSlasherCallerSession) Challengers(arg0 [32]byte) (common.Address, error) {
	return _DarknodeSlasher.Contract.Challengers(&_DarknodeSlasher.CallOpts, arg0)
}

// OrderDetails is a free data retrieval call binding the contract method 0xa3c50b32.
//
// Solidity: function orderDetails( bytes32) constant returns(settlementID uint64, tokens uint64, price uint256, volume uint256, minimumVolume uint256)
func (_DarknodeSlasher *DarknodeSlasherCaller) OrderDetails(opts *bind.CallOpts, arg0 [32]byte) (struct {
	SettlementID  uint64
	Tokens        uint64
	Price         *big.Int
	Volume        *big.Int
	MinimumVolume *big.Int
}, error) {
	ret := new(struct {
		SettlementID  uint64
		Tokens        uint64
		Price         *big.Int
		Volume        *big.Int
		MinimumVolume *big.Int
	})
	out := ret
	err := _DarknodeSlasher.contract.Call(opts, out, "orderDetails", arg0)
	return *ret, err
}

// OrderDetails is a free data retrieval call binding the contract method 0xa3c50b32.
//
// Solidity: function orderDetails( bytes32) constant returns(settlementID uint64, tokens uint64, price uint256, volume uint256, minimumVolume uint256)
func (_DarknodeSlasher *DarknodeSlasherSession) OrderDetails(arg0 [32]byte) (struct {
	SettlementID  uint64
	Tokens        uint64
	Price         *big.Int
	Volume        *big.Int
	MinimumVolume *big.Int
}, error) {
	return _DarknodeSlasher.Contract.OrderDetails(&_DarknodeSlasher.CallOpts, arg0)
}

// OrderDetails is a free data retrieval call binding the contract method 0xa3c50b32.
//
// Solidity: function orderDetails( bytes32) constant returns(settlementID uint64, tokens uint64, price uint256, volume uint256, minimumVolume uint256)
func (_DarknodeSlasher *DarknodeSlasherCallerSession) OrderDetails(arg0 [32]byte) (struct {
	SettlementID  uint64
	Tokens        uint64
	Price         *big.Int
	Volume        *big.Int
	MinimumVolume *big.Int
}, error) {
	return _DarknodeSlasher.Contract.OrderDetails(&_DarknodeSlasher.CallOpts, arg0)
}

// OrderSubmitted is a free data retrieval call binding the contract method 0x4778a5be.
//
// Solidity: function orderSubmitted( bytes32) constant returns(bool)
func (_DarknodeSlasher *DarknodeSlasherCaller) OrderSubmitted(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _DarknodeSlasher.contract.Call(opts, out, "orderSubmitted", arg0)
	return *ret0, err
}

// OrderSubmitted is a free data retrieval call binding the contract method 0x4778a5be.
//
// Solidity: function orderSubmitted( bytes32) constant returns(bool)
func (_DarknodeSlasher *DarknodeSlasherSession) OrderSubmitted(arg0 [32]byte) (bool, error) {
	return _DarknodeSlasher.Contract.OrderSubmitted(&_DarknodeSlasher.CallOpts, arg0)
}

// OrderSubmitted is a free data retrieval call binding the contract method 0x4778a5be.
//
// Solidity: function orderSubmitted( bytes32) constant returns(bool)
func (_DarknodeSlasher *DarknodeSlasherCallerSession) OrderSubmitted(arg0 [32]byte) (bool, error) {
	return _DarknodeSlasher.Contract.OrderSubmitted(&_DarknodeSlasher.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_DarknodeSlasher *DarknodeSlasherCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DarknodeSlasher.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_DarknodeSlasher *DarknodeSlasherSession) Owner() (common.Address, error) {
	return _DarknodeSlasher.Contract.Owner(&_DarknodeSlasher.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_DarknodeSlasher *DarknodeSlasherCallerSession) Owner() (common.Address, error) {
	return _DarknodeSlasher.Contract.Owner(&_DarknodeSlasher.CallOpts)
}

// TrustedDarknodeRegistry is a free data retrieval call binding the contract method 0xfc722b2b.
//
// Solidity: function trustedDarknodeRegistry() constant returns(address)
func (_DarknodeSlasher *DarknodeSlasherCaller) TrustedDarknodeRegistry(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DarknodeSlasher.contract.Call(opts, out, "trustedDarknodeRegistry")
	return *ret0, err
}

// TrustedDarknodeRegistry is a free data retrieval call binding the contract method 0xfc722b2b.
//
// Solidity: function trustedDarknodeRegistry() constant returns(address)
func (_DarknodeSlasher *DarknodeSlasherSession) TrustedDarknodeRegistry() (common.Address, error) {
	return _DarknodeSlasher.Contract.TrustedDarknodeRegistry(&_DarknodeSlasher.CallOpts)
}

// TrustedDarknodeRegistry is a free data retrieval call binding the contract method 0xfc722b2b.
//
// Solidity: function trustedDarknodeRegistry() constant returns(address)
func (_DarknodeSlasher *DarknodeSlasherCallerSession) TrustedDarknodeRegistry() (common.Address, error) {
	return _DarknodeSlasher.Contract.TrustedDarknodeRegistry(&_DarknodeSlasher.CallOpts)
}

// TrustedOrderbook is a free data retrieval call binding the contract method 0x26cf660d.
//
// Solidity: function trustedOrderbook() constant returns(address)
func (_DarknodeSlasher *DarknodeSlasherCaller) TrustedOrderbook(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DarknodeSlasher.contract.Call(opts, out, "trustedOrderbook")
	return *ret0, err
}

// TrustedOrderbook is a free data retrieval call binding the contract method 0x26cf660d.
//
// Solidity: function trustedOrderbook() constant returns(address)
func (_DarknodeSlasher *DarknodeSlasherSession) TrustedOrderbook() (common.Address, error) {
	return _DarknodeSlasher.Contract.TrustedOrderbook(&_DarknodeSlasher.CallOpts)
}

// TrustedOrderbook is a free data retrieval call binding the contract method 0x26cf660d.
//
// Solidity: function trustedOrderbook() constant returns(address)
func (_DarknodeSlasher *DarknodeSlasherCallerSession) TrustedOrderbook() (common.Address, error) {
	return _DarknodeSlasher.Contract.TrustedOrderbook(&_DarknodeSlasher.CallOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DarknodeSlasher *DarknodeSlasherTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DarknodeSlasher.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DarknodeSlasher *DarknodeSlasherSession) RenounceOwnership() (*types.Transaction, error) {
	return _DarknodeSlasher.Contract.RenounceOwnership(&_DarknodeSlasher.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DarknodeSlasher *DarknodeSlasherTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _DarknodeSlasher.Contract.RenounceOwnership(&_DarknodeSlasher.TransactOpts)
}

// SubmitChallenge is a paid mutator transaction binding the contract method 0x4d4582a1.
//
// Solidity: function submitChallenge(_buyID bytes32, _sellID bytes32) returns()
func (_DarknodeSlasher *DarknodeSlasherTransactor) SubmitChallenge(opts *bind.TransactOpts, _buyID [32]byte, _sellID [32]byte) (*types.Transaction, error) {
	return _DarknodeSlasher.contract.Transact(opts, "submitChallenge", _buyID, _sellID)
}

// SubmitChallenge is a paid mutator transaction binding the contract method 0x4d4582a1.
//
// Solidity: function submitChallenge(_buyID bytes32, _sellID bytes32) returns()
func (_DarknodeSlasher *DarknodeSlasherSession) SubmitChallenge(_buyID [32]byte, _sellID [32]byte) (*types.Transaction, error) {
	return _DarknodeSlasher.Contract.SubmitChallenge(&_DarknodeSlasher.TransactOpts, _buyID, _sellID)
}

// SubmitChallenge is a paid mutator transaction binding the contract method 0x4d4582a1.
//
// Solidity: function submitChallenge(_buyID bytes32, _sellID bytes32) returns()
func (_DarknodeSlasher *DarknodeSlasherTransactorSession) SubmitChallenge(_buyID [32]byte, _sellID [32]byte) (*types.Transaction, error) {
	return _DarknodeSlasher.Contract.SubmitChallenge(&_DarknodeSlasher.TransactOpts, _buyID, _sellID)
}

// SubmitChallengeOrder is a paid mutator transaction binding the contract method 0xc41c33af.
//
// Solidity: function submitChallengeOrder(details bytes, settlementID uint64, tokens uint64, price uint256, volume uint256, minimumVolume uint256) returns()
func (_DarknodeSlasher *DarknodeSlasherTransactor) SubmitChallengeOrder(opts *bind.TransactOpts, details []byte, settlementID uint64, tokens uint64, price *big.Int, volume *big.Int, minimumVolume *big.Int) (*types.Transaction, error) {
	return _DarknodeSlasher.contract.Transact(opts, "submitChallengeOrder", details, settlementID, tokens, price, volume, minimumVolume)
}

// SubmitChallengeOrder is a paid mutator transaction binding the contract method 0xc41c33af.
//
// Solidity: function submitChallengeOrder(details bytes, settlementID uint64, tokens uint64, price uint256, volume uint256, minimumVolume uint256) returns()
func (_DarknodeSlasher *DarknodeSlasherSession) SubmitChallengeOrder(details []byte, settlementID uint64, tokens uint64, price *big.Int, volume *big.Int, minimumVolume *big.Int) (*types.Transaction, error) {
	return _DarknodeSlasher.Contract.SubmitChallengeOrder(&_DarknodeSlasher.TransactOpts, details, settlementID, tokens, price, volume, minimumVolume)
}

// SubmitChallengeOrder is a paid mutator transaction binding the contract method 0xc41c33af.
//
// Solidity: function submitChallengeOrder(details bytes, settlementID uint64, tokens uint64, price uint256, volume uint256, minimumVolume uint256) returns()
func (_DarknodeSlasher *DarknodeSlasherTransactorSession) SubmitChallengeOrder(details []byte, settlementID uint64, tokens uint64, price *big.Int, volume *big.Int, minimumVolume *big.Int) (*types.Transaction, error) {
	return _DarknodeSlasher.Contract.SubmitChallengeOrder(&_DarknodeSlasher.TransactOpts, details, settlementID, tokens, price, volume, minimumVolume)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_DarknodeSlasher *DarknodeSlasherTransactor) TransferOwnership(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _DarknodeSlasher.contract.Transact(opts, "transferOwnership", _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_DarknodeSlasher *DarknodeSlasherSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _DarknodeSlasher.Contract.TransferOwnership(&_DarknodeSlasher.TransactOpts, _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_DarknodeSlasher *DarknodeSlasherTransactorSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _DarknodeSlasher.Contract.TransferOwnership(&_DarknodeSlasher.TransactOpts, _newOwner)
}

// DarknodeSlasherOwnershipRenouncedIterator is returned from FilterOwnershipRenounced and is used to iterate over the raw logs and unpacked data for OwnershipRenounced events raised by the DarknodeSlasher contract.
type DarknodeSlasherOwnershipRenouncedIterator struct {
	Event *DarknodeSlasherOwnershipRenounced // Event containing the contract specifics and raw log

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
func (it *DarknodeSlasherOwnershipRenouncedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DarknodeSlasherOwnershipRenounced)
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
		it.Event = new(DarknodeSlasherOwnershipRenounced)
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
func (it *DarknodeSlasherOwnershipRenouncedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DarknodeSlasherOwnershipRenouncedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DarknodeSlasherOwnershipRenounced represents a OwnershipRenounced event raised by the DarknodeSlasher contract.
type DarknodeSlasherOwnershipRenounced struct {
	PreviousOwner common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipRenounced is a free log retrieval operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_DarknodeSlasher *DarknodeSlasherFilterer) FilterOwnershipRenounced(opts *bind.FilterOpts, previousOwner []common.Address) (*DarknodeSlasherOwnershipRenouncedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _DarknodeSlasher.contract.FilterLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DarknodeSlasherOwnershipRenouncedIterator{contract: _DarknodeSlasher.contract, event: "OwnershipRenounced", logs: logs, sub: sub}, nil
}

// WatchOwnershipRenounced is a free log subscription operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_DarknodeSlasher *DarknodeSlasherFilterer) WatchOwnershipRenounced(opts *bind.WatchOpts, sink chan<- *DarknodeSlasherOwnershipRenounced, previousOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _DarknodeSlasher.contract.WatchLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DarknodeSlasherOwnershipRenounced)
				if err := _DarknodeSlasher.contract.UnpackLog(event, "OwnershipRenounced", log); err != nil {
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

// DarknodeSlasherOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the DarknodeSlasher contract.
type DarknodeSlasherOwnershipTransferredIterator struct {
	Event *DarknodeSlasherOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *DarknodeSlasherOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DarknodeSlasherOwnershipTransferred)
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
		it.Event = new(DarknodeSlasherOwnershipTransferred)
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
func (it *DarknodeSlasherOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DarknodeSlasherOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DarknodeSlasherOwnershipTransferred represents a OwnershipTransferred event raised by the DarknodeSlasher contract.
type DarknodeSlasherOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_DarknodeSlasher *DarknodeSlasherFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*DarknodeSlasherOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DarknodeSlasher.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DarknodeSlasherOwnershipTransferredIterator{contract: _DarknodeSlasher.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_DarknodeSlasher *DarknodeSlasherFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DarknodeSlasherOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DarknodeSlasher.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DarknodeSlasherOwnershipTransferred)
				if err := _DarknodeSlasher.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ECRecoveryABI is the input ABI used to generate the binding from.
const ECRecoveryABI = "[]"

// ECRecoveryBin is the compiled bytecode used for deploying new contracts.
const ECRecoveryBin = `0x604c602c600b82828239805160001a60731460008114601c57601e565bfe5b5030600052607381538281f30073000000000000000000000000000000000000000030146080604052600080fd00a165627a7a7230582051ce5400e9db2ec99f60568a05f33dba86f40cf2f0c53aa9337bccedec093bd10029`

// DeployECRecovery deploys a new Ethereum contract, binding an instance of ECRecovery to it.
func DeployECRecovery(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ECRecovery, error) {
	parsed, err := abi.JSON(strings.NewReader(ECRecoveryABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ECRecoveryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ECRecovery{ECRecoveryCaller: ECRecoveryCaller{contract: contract}, ECRecoveryTransactor: ECRecoveryTransactor{contract: contract}, ECRecoveryFilterer: ECRecoveryFilterer{contract: contract}}, nil
}

// ECRecovery is an auto generated Go binding around an Ethereum contract.
type ECRecovery struct {
	ECRecoveryCaller     // Read-only binding to the contract
	ECRecoveryTransactor // Write-only binding to the contract
	ECRecoveryFilterer   // Log filterer for contract events
}

// ECRecoveryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ECRecoveryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ECRecoveryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ECRecoveryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ECRecoveryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ECRecoveryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ECRecoverySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ECRecoverySession struct {
	Contract     *ECRecovery       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ECRecoveryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ECRecoveryCallerSession struct {
	Contract *ECRecoveryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// ECRecoveryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ECRecoveryTransactorSession struct {
	Contract     *ECRecoveryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// ECRecoveryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ECRecoveryRaw struct {
	Contract *ECRecovery // Generic contract binding to access the raw methods on
}

// ECRecoveryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ECRecoveryCallerRaw struct {
	Contract *ECRecoveryCaller // Generic read-only contract binding to access the raw methods on
}

// ECRecoveryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ECRecoveryTransactorRaw struct {
	Contract *ECRecoveryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewECRecovery creates a new instance of ECRecovery, bound to a specific deployed contract.
func NewECRecovery(address common.Address, backend bind.ContractBackend) (*ECRecovery, error) {
	contract, err := bindECRecovery(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ECRecovery{ECRecoveryCaller: ECRecoveryCaller{contract: contract}, ECRecoveryTransactor: ECRecoveryTransactor{contract: contract}, ECRecoveryFilterer: ECRecoveryFilterer{contract: contract}}, nil
}

// NewECRecoveryCaller creates a new read-only instance of ECRecovery, bound to a specific deployed contract.
func NewECRecoveryCaller(address common.Address, caller bind.ContractCaller) (*ECRecoveryCaller, error) {
	contract, err := bindECRecovery(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ECRecoveryCaller{contract: contract}, nil
}

// NewECRecoveryTransactor creates a new write-only instance of ECRecovery, bound to a specific deployed contract.
func NewECRecoveryTransactor(address common.Address, transactor bind.ContractTransactor) (*ECRecoveryTransactor, error) {
	contract, err := bindECRecovery(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ECRecoveryTransactor{contract: contract}, nil
}

// NewECRecoveryFilterer creates a new log filterer instance of ECRecovery, bound to a specific deployed contract.
func NewECRecoveryFilterer(address common.Address, filterer bind.ContractFilterer) (*ECRecoveryFilterer, error) {
	contract, err := bindECRecovery(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ECRecoveryFilterer{contract: contract}, nil
}

// bindECRecovery binds a generic wrapper to an already deployed contract.
func bindECRecovery(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ECRecoveryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ECRecovery *ECRecoveryRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ECRecovery.Contract.ECRecoveryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ECRecovery *ECRecoveryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ECRecovery.Contract.ECRecoveryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ECRecovery *ECRecoveryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ECRecovery.Contract.ECRecoveryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ECRecovery *ECRecoveryCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ECRecovery.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ECRecovery *ECRecoveryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ECRecovery.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ECRecovery *ECRecoveryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ECRecovery.Contract.contract.Transact(opts, method, params...)
}

// ERC20ABI is the input ABI used to generate the binding from.
const ERC20ABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"who\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]"

// ERC20Bin is the compiled bytecode used for deploying new contracts.
const ERC20Bin = `0x`

// DeployERC20 deploys a new Ethereum contract, binding an instance of ERC20 to it.
func DeployERC20(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ERC20, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ERC20Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ERC20{ERC20Caller: ERC20Caller{contract: contract}, ERC20Transactor: ERC20Transactor{contract: contract}, ERC20Filterer: ERC20Filterer{contract: contract}}, nil
}

// ERC20 is an auto generated Go binding around an Ethereum contract.
type ERC20 struct {
	ERC20Caller     // Read-only binding to the contract
	ERC20Transactor // Write-only binding to the contract
	ERC20Filterer   // Log filterer for contract events
}

// ERC20Caller is an auto generated read-only Go binding around an Ethereum contract.
type ERC20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type ERC20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ERC20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ERC20Session struct {
	Contract     *ERC20            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ERC20CallerSession struct {
	Contract *ERC20Caller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ERC20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ERC20TransactorSession struct {
	Contract     *ERC20Transactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC20Raw is an auto generated low-level Go binding around an Ethereum contract.
type ERC20Raw struct {
	Contract *ERC20 // Generic contract binding to access the raw methods on
}

// ERC20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ERC20CallerRaw struct {
	Contract *ERC20Caller // Generic read-only contract binding to access the raw methods on
}

// ERC20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ERC20TransactorRaw struct {
	Contract *ERC20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewERC20 creates a new instance of ERC20, bound to a specific deployed contract.
func NewERC20(address common.Address, backend bind.ContractBackend) (*ERC20, error) {
	contract, err := bindERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC20{ERC20Caller: ERC20Caller{contract: contract}, ERC20Transactor: ERC20Transactor{contract: contract}, ERC20Filterer: ERC20Filterer{contract: contract}}, nil
}

// NewERC20Caller creates a new read-only instance of ERC20, bound to a specific deployed contract.
func NewERC20Caller(address common.Address, caller bind.ContractCaller) (*ERC20Caller, error) {
	contract, err := bindERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20Caller{contract: contract}, nil
}

// NewERC20Transactor creates a new write-only instance of ERC20, bound to a specific deployed contract.
func NewERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*ERC20Transactor, error) {
	contract, err := bindERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20Transactor{contract: contract}, nil
}

// NewERC20Filterer creates a new log filterer instance of ERC20, bound to a specific deployed contract.
func NewERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*ERC20Filterer, error) {
	contract, err := bindERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC20Filterer{contract: contract}, nil
}

// bindERC20 binds a generic wrapper to an already deployed contract.
func bindERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20 *ERC20Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC20.Contract.ERC20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20 *ERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20.Contract.ERC20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20 *ERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20.Contract.ERC20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20 *ERC20CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20 *ERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20 *ERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(owner address, spender address) constant returns(uint256)
func (_ERC20 *ERC20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20.contract.Call(opts, out, "allowance", owner, spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(owner address, spender address) constant returns(uint256)
func (_ERC20 *ERC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ERC20.Contract.Allowance(&_ERC20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(owner address, spender address) constant returns(uint256)
func (_ERC20 *ERC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ERC20.Contract.Allowance(&_ERC20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(who address) constant returns(uint256)
func (_ERC20 *ERC20Caller) BalanceOf(opts *bind.CallOpts, who common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20.contract.Call(opts, out, "balanceOf", who)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(who address) constant returns(uint256)
func (_ERC20 *ERC20Session) BalanceOf(who common.Address) (*big.Int, error) {
	return _ERC20.Contract.BalanceOf(&_ERC20.CallOpts, who)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(who address) constant returns(uint256)
func (_ERC20 *ERC20CallerSession) BalanceOf(who common.Address) (*big.Int, error) {
	return _ERC20.Contract.BalanceOf(&_ERC20.CallOpts, who)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20 *ERC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20 *ERC20Session) TotalSupply() (*big.Int, error) {
	return _ERC20.Contract.TotalSupply(&_ERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20 *ERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _ERC20.Contract.TotalSupply(&_ERC20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns(bool)
func (_ERC20 *ERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns(bool)
func (_ERC20 *ERC20Session) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Approve(&_ERC20.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns(bool)
func (_ERC20 *ERC20TransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Approve(&_ERC20.TransactOpts, spender, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns(bool)
func (_ERC20 *ERC20Transactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns(bool)
func (_ERC20 *ERC20Session) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Transfer(&_ERC20.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns(bool)
func (_ERC20 *ERC20TransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Transfer(&_ERC20.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns(bool)
func (_ERC20 *ERC20Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns(bool)
func (_ERC20 *ERC20Session) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.TransferFrom(&_ERC20.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns(bool)
func (_ERC20 *ERC20TransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.TransferFrom(&_ERC20.TransactOpts, from, to, value)
}

// ERC20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ERC20 contract.
type ERC20ApprovalIterator struct {
	Event *ERC20Approval // Event containing the contract specifics and raw log

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
func (it *ERC20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20Approval)
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
		it.Event = new(ERC20Approval)
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
func (it *ERC20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20Approval represents a Approval event raised by the ERC20 contract.
type ERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(owner indexed address, spender indexed address, value uint256)
func (_ERC20 *ERC20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ERC20ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ERC20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ERC20ApprovalIterator{contract: _ERC20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(owner indexed address, spender indexed address, value uint256)
func (_ERC20 *ERC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ERC20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20Approval)
				if err := _ERC20.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ERC20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ERC20 contract.
type ERC20TransferIterator struct {
	Event *ERC20Transfer // Event containing the contract specifics and raw log

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
func (it *ERC20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20Transfer)
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
		it.Event = new(ERC20Transfer)
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
func (it *ERC20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20Transfer represents a Transfer event raised by the ERC20 contract.
type ERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
func (_ERC20 *ERC20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ERC20TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ERC20TransferIterator{contract: _ERC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
func (_ERC20 *ERC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20Transfer)
				if err := _ERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ERC20BasicABI is the input ABI used to generate the binding from.
const ERC20BasicABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"who\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]"

// ERC20BasicBin is the compiled bytecode used for deploying new contracts.
const ERC20BasicBin = `0x`

// DeployERC20Basic deploys a new Ethereum contract, binding an instance of ERC20Basic to it.
func DeployERC20Basic(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ERC20Basic, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20BasicABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ERC20BasicBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ERC20Basic{ERC20BasicCaller: ERC20BasicCaller{contract: contract}, ERC20BasicTransactor: ERC20BasicTransactor{contract: contract}, ERC20BasicFilterer: ERC20BasicFilterer{contract: contract}}, nil
}

// ERC20Basic is an auto generated Go binding around an Ethereum contract.
type ERC20Basic struct {
	ERC20BasicCaller     // Read-only binding to the contract
	ERC20BasicTransactor // Write-only binding to the contract
	ERC20BasicFilterer   // Log filterer for contract events
}

// ERC20BasicCaller is an auto generated read-only Go binding around an Ethereum contract.
type ERC20BasicCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20BasicTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ERC20BasicTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20BasicFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ERC20BasicFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20BasicSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ERC20BasicSession struct {
	Contract     *ERC20Basic       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC20BasicCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ERC20BasicCallerSession struct {
	Contract *ERC20BasicCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// ERC20BasicTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ERC20BasicTransactorSession struct {
	Contract     *ERC20BasicTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// ERC20BasicRaw is an auto generated low-level Go binding around an Ethereum contract.
type ERC20BasicRaw struct {
	Contract *ERC20Basic // Generic contract binding to access the raw methods on
}

// ERC20BasicCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ERC20BasicCallerRaw struct {
	Contract *ERC20BasicCaller // Generic read-only contract binding to access the raw methods on
}

// ERC20BasicTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ERC20BasicTransactorRaw struct {
	Contract *ERC20BasicTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC20Basic creates a new instance of ERC20Basic, bound to a specific deployed contract.
func NewERC20Basic(address common.Address, backend bind.ContractBackend) (*ERC20Basic, error) {
	contract, err := bindERC20Basic(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC20Basic{ERC20BasicCaller: ERC20BasicCaller{contract: contract}, ERC20BasicTransactor: ERC20BasicTransactor{contract: contract}, ERC20BasicFilterer: ERC20BasicFilterer{contract: contract}}, nil
}

// NewERC20BasicCaller creates a new read-only instance of ERC20Basic, bound to a specific deployed contract.
func NewERC20BasicCaller(address common.Address, caller bind.ContractCaller) (*ERC20BasicCaller, error) {
	contract, err := bindERC20Basic(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20BasicCaller{contract: contract}, nil
}

// NewERC20BasicTransactor creates a new write-only instance of ERC20Basic, bound to a specific deployed contract.
func NewERC20BasicTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC20BasicTransactor, error) {
	contract, err := bindERC20Basic(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20BasicTransactor{contract: contract}, nil
}

// NewERC20BasicFilterer creates a new log filterer instance of ERC20Basic, bound to a specific deployed contract.
func NewERC20BasicFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC20BasicFilterer, error) {
	contract, err := bindERC20Basic(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC20BasicFilterer{contract: contract}, nil
}

// bindERC20Basic binds a generic wrapper to an already deployed contract.
func bindERC20Basic(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20BasicABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20Basic *ERC20BasicRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC20Basic.Contract.ERC20BasicCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20Basic *ERC20BasicRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20Basic.Contract.ERC20BasicTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20Basic *ERC20BasicRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20Basic.Contract.ERC20BasicTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20Basic *ERC20BasicCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC20Basic.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20Basic *ERC20BasicTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20Basic.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20Basic *ERC20BasicTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20Basic.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(who address) constant returns(uint256)
func (_ERC20Basic *ERC20BasicCaller) BalanceOf(opts *bind.CallOpts, who common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20Basic.contract.Call(opts, out, "balanceOf", who)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(who address) constant returns(uint256)
func (_ERC20Basic *ERC20BasicSession) BalanceOf(who common.Address) (*big.Int, error) {
	return _ERC20Basic.Contract.BalanceOf(&_ERC20Basic.CallOpts, who)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(who address) constant returns(uint256)
func (_ERC20Basic *ERC20BasicCallerSession) BalanceOf(who common.Address) (*big.Int, error) {
	return _ERC20Basic.Contract.BalanceOf(&_ERC20Basic.CallOpts, who)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20Basic *ERC20BasicCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20Basic.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20Basic *ERC20BasicSession) TotalSupply() (*big.Int, error) {
	return _ERC20Basic.Contract.TotalSupply(&_ERC20Basic.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20Basic *ERC20BasicCallerSession) TotalSupply() (*big.Int, error) {
	return _ERC20Basic.Contract.TotalSupply(&_ERC20Basic.CallOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns(bool)
func (_ERC20Basic *ERC20BasicTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Basic.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns(bool)
func (_ERC20Basic *ERC20BasicSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Basic.Contract.Transfer(&_ERC20Basic.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns(bool)
func (_ERC20Basic *ERC20BasicTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Basic.Contract.Transfer(&_ERC20Basic.TransactOpts, to, value)
}

// ERC20BasicTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ERC20Basic contract.
type ERC20BasicTransferIterator struct {
	Event *ERC20BasicTransfer // Event containing the contract specifics and raw log

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
func (it *ERC20BasicTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20BasicTransfer)
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
		it.Event = new(ERC20BasicTransfer)
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
func (it *ERC20BasicTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20BasicTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20BasicTransfer represents a Transfer event raised by the ERC20Basic contract.
type ERC20BasicTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
func (_ERC20Basic *ERC20BasicFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ERC20BasicTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC20Basic.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ERC20BasicTransferIterator{contract: _ERC20Basic.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
func (_ERC20Basic *ERC20BasicFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ERC20BasicTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC20Basic.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20BasicTransfer)
				if err := _ERC20Basic.contract.UnpackLog(event, "Transfer", log); err != nil {
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
const LinkedListABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"NULL\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// LinkedListBin is the compiled bytecode used for deploying new contracts.
const LinkedListBin = `0x60ba61002f600b82828239805160001a6073146000811461001f57610021565bfe5b5030600052607381538281f300730000000000000000000000000000000000000000301460806040526004361060555763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663f26be3fc8114605a575b600080fd5b60606089565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b6000815600a165627a7a7230582075f4a127afe163172cfcf00bd1b47c2f6d6b63ac967564c118278247bcb405220029`

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
// Solidity: function NULL() constant returns(address)
func (_LinkedList *LinkedListCaller) NULL(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _LinkedList.contract.Call(opts, out, "NULL")
	return *ret0, err
}

// NULL is a free data retrieval call binding the contract method 0xf26be3fc.
//
// Solidity: function NULL() constant returns(address)
func (_LinkedList *LinkedListSession) NULL() (common.Address, error) {
	return _LinkedList.Contract.NULL(&_LinkedList.CallOpts)
}

// NULL is a free data retrieval call binding the contract method 0xf26be3fc.
//
// Solidity: function NULL() constant returns(address)
func (_LinkedList *LinkedListCallerSession) NULL() (common.Address, error) {
	return _LinkedList.Contract.NULL(&_LinkedList.CallOpts)
}

// OrderbookABI is the input ABI used to generate the binding from.
const OrderbookABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"_orderID\",\"type\":\"bytes32\"}],\"name\":\"orderConfirmer\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_orderID\",\"type\":\"bytes32\"},{\"name\":\"_matchedOrderID\",\"type\":\"bytes32\"}],\"name\":\"confirmOrder\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"ordersCount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_orderID\",\"type\":\"bytes32\"}],\"name\":\"cancelOrder\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_orderID\",\"type\":\"bytes32\"}],\"name\":\"orderBlockNumber\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"ren\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_offset\",\"type\":\"uint256\"},{\"name\":\"_limit\",\"type\":\"uint256\"}],\"name\":\"getOrders\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32[]\"},{\"name\":\"\",\"type\":\"address[]\"},{\"name\":\"\",\"type\":\"uint8[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"orders\",\"outputs\":[{\"name\":\"state\",\"type\":\"uint8\"},{\"name\":\"trader\",\"type\":\"address\"},{\"name\":\"confirmer\",\"type\":\"address\"},{\"name\":\"settlementID\",\"type\":\"uint64\"},{\"name\":\"priority\",\"type\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"name\":\"matchedOrder\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"darknodeRegistry\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_orderID\",\"type\":\"bytes32\"}],\"name\":\"orderDepth\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_settlementID\",\"type\":\"uint64\"},{\"name\":\"_signature\",\"type\":\"bytes\"},{\"name\":\"_orderID\",\"type\":\"bytes32\"}],\"name\":\"openOrder\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_orderID\",\"type\":\"bytes32\"}],\"name\":\"orderState\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newDarknodeRegistry\",\"type\":\"address\"}],\"name\":\"updateDarknodeRegistry\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_orderID\",\"type\":\"bytes32\"}],\"name\":\"orderMatch\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_orderID\",\"type\":\"bytes32\"}],\"name\":\"orderTrader\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_orderID\",\"type\":\"bytes32\"}],\"name\":\"orderPriority\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"settlementRegistry\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_VERSION\",\"type\":\"string\"},{\"name\":\"_renAddress\",\"type\":\"address\"},{\"name\":\"_darknodeRegistry\",\"type\":\"address\"},{\"name\":\"_settlementRegistry\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"previousFee\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"nextFee\",\"type\":\"uint256\"}],\"name\":\"LogFeeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"previousDarknodeRegistry\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"nextDarknodeRegistry\",\"type\":\"address\"}],\"name\":\"LogDarknodeRegistryUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"}],\"name\":\"OwnershipRenounced\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"

// OrderbookBin is the compiled bytecode used for deploying new contracts.
const OrderbookBin = `0x60806040523480156200001157600080fd5b50604051620014f1380380620014f1833981016040908152815160208084015192840151606085015160008054600160a060020a03191633179055929094018051909492916200006791600191870190620000ac565b5060028054600160a060020a03948516600160a060020a0319918216179091556003805493851693821693909317909255600480549190931691161790555062000151565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10620000ef57805160ff19168380011785556200011f565b828001600101855582156200011f579182015b828111156200011f57825182559160200191906001019062000102565b506200012d92915062000131565b5090565b6200014e91905b808211156200012d576000815560010162000138565b90565b61139080620001616000396000f30060806040526004361061011c5763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416631107c3f781146101215780632a0401f01461015557806335daa73114610172578063715018a6146101995780637489ec23146101ae57806389895d53146101c65780638a9b4067146101de5780638da5cb5b146101f35780638f72fc77146102085780639c3f1e90146103015780639e45e0d01461037c578063a188fcb814610391578063a1ce7e03146103a9578063aab14d04146103da578063aaff096d14610416578063af3e8a4014610437578063b1a080101461044f578063b248e4e114610467578063d09812e11461047f578063f2fde38b14610494578063ffa1ad74146104b5575b600080fd5b34801561012d57600080fd5b5061013960043561053f565b60408051600160a060020a039092168252519081900360200190f35b34801561016157600080fd5b50610170600435602435610560565b005b34801561017e57600080fd5b506101876107b1565b60408051918252519081900360200190f35b3480156101a557600080fd5b506101706107b7565b3480156101ba57600080fd5b50610170600435610823565b3480156101d257600080fd5b50610187600435610a0b565b3480156101ea57600080fd5b50610139610a20565b3480156101ff57600080fd5b50610139610a2f565b34801561021457600080fd5b50610223600435602435610a3e565b60405180806020018060200180602001848103845287818151815260200191508051906020019060200280838360005b8381101561026b578181015183820152602001610253565b50505050905001848103835286818151815260200191508051906020019060200280838360005b838110156102aa578181015183820152602001610292565b50505050905001848103825285818151815260200191508051906020019060200280838360005b838110156102e95781810151838201526020016102d1565b50505050905001965050505050505060405180910390f35b34801561030d57600080fd5b50610319600435610be8565b6040518088600381111561032957fe5b60ff168152600160a060020a0397881660208201529590961660408087019190915267ffffffffffffffff9094166060860152608085019290925260a084015260c0830152519081900360e00192509050f35b34801561038857600080fd5b50610139610c4e565b34801561039d57600080fd5b50610187600435610c5d565b3480156103b557600080fd5b506101706004803567ffffffffffffffff169060248035908101910135604435610c95565b3480156103e657600080fd5b506103f260043561113b565b6040518082600381111561040257fe5b60ff16815260200191505060405180910390f35b34801561042257600080fd5b50610170600160a060020a0360043516611150565b34801561044357600080fd5b506101876004356111de565b34801561045b57600080fd5b506101396004356111f3565b34801561047357600080fd5b50610187600435611213565b34801561048b57600080fd5b50610139611228565b3480156104a057600080fd5b50610170600160a060020a0360043516611237565b3480156104c157600080fd5b506104ca61125a565b6040805160208082528351818301528351919283929083019185019080838360005b838110156105045781810151838201526020016104ec565b50505050905090810190601f1680156105315780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b600081815260066020526040902060010154600160a060020a03165b919050565b600354604080517fc3c5a547000000000000000000000000000000000000000000000000000000008152336004820181905291519192600160a060020a03169163c3c5a547916024808201926020929091908290030181600087803b1580156105c857600080fd5b505af11580156105dc573d6000803e3d6000fd5b505050506040513d60208110156105f257600080fd5b5051151561064a576040805160e560020a62461bcd02815260206004820152601b60248201527f6d7573742062652072656769737465726564206461726b6e6f64650000000000604482015290519081900360640190fd5b600160008481526006602052604090205460ff16600381111561066957fe5b146106be576040805160e560020a62461bcd02815260206004820152601460248201527f696e76616c6964206f7264657220737461747573000000000000000000000000604482015290519081900360640190fd5b600160008381526006602052604090205460ff1660038111156106dd57fe5b14610732576040805160e560020a62461bcd02815260206004820152601460248201527f696e76616c6964206f7264657220737461747573000000000000000000000000604482015290519081900360640190fd5b50600082815260066020526040808220805460ff1990811660029081178355600180840180543373ffffffffffffffffffffffffffffffffffffffff19918216811790925560048087018a9055436003978801819055998952969097208054909416909217835582018054909516179093559082019390935590910155565b60055490565b600054600160a060020a031633146107ce57600080fd5b60008054604051600160a060020a03909116917ff8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c6482091a26000805473ffffffffffffffffffffffffffffffffffffffff19169055565b6000600160008381526006602052604090205460ff16600381111561084457fe5b14610899576040805160e560020a62461bcd02815260206004820152601360248201527f696e76616c6964206f7264657220737461746500000000000000000000000000604482015290519081900360640190fd5b6004805460008481526006602090815260408083206001015481517fa6065c960000000000000000000000000000000000000000000000000000000081527401000000000000000000000000000000000000000090910467ffffffffffffffff169581019590955251600160a060020a039093169363a6065c969360248083019491928390030190829087803b15801561093257600080fd5b505af1158015610946573d6000803e3d6000fd5b505050506040513d602081101561095c57600080fd5b50516000838152600660205260409020549091506101009004600160a060020a0316331480610993575033600160a060020a038216145b15156109e9576040805160e560020a62461bcd02815260206004820152600e60248201527f6e6f7420617574686f72697a6564000000000000000000000000000000000000604482015290519081900360640190fd5b506000908152600660205260409020805460ff19166003908117825543910155565b60009081526006602052604090206003015490565b600254600160a060020a031681565b600054600160a060020a031681565b6060806060600060608060606000806005805490508b101515610a6057610bdb565b6005548a96508b87011115610a78576005548b900395505b85604051908082528060200260200182016040528015610aa2578160200160208202803883390190505b50945085604051908082528060200260200182016040528015610acf578160200160208202803883390190505b50935085604051908082528060200260200182016040528015610afc578160200160208202803883390190505b509250600091505b85821015610bd15760058054838d01908110610b1c57fe5b90600052602060002001549050808583815181101515610b3857fe5b6020908102919091018101919091526000828152600690915260409020548451610100909104600160a060020a031690859084908110610b7457fe5b600160a060020a03909216602092830290910182015260008281526006909152604090205460ff166003811115610ba757fe5b8383815181101515610bb557fe5b60ff909216602092830290910190910152600190910190610b04565b8484849850985098505b5050505050509250925092565b6006602052600090815260409020805460018201546002830154600384015460049094015460ff841694600160a060020a0361010090950485169484169367ffffffffffffffff7401000000000000000000000000000000000000000090910416929187565b600354600160a060020a031681565b6000818152600660205260408120600301541515610c7d5750600061055b565b50600090815260066020526040902060030154430390565b6000808060008481526006602052604090205460ff166003811115610cb657fe5b14610d0b576040805160e560020a62461bcd02815260206004820152601460248201527f696e76616c6964206f7264657220737461747573000000000000000000000000604482015290519081900360640190fd5b60048054604080517ff8f9be3600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8a169381019390935251339450600160a060020a039091169163f8f9be369160248083019260209291908290030181600087803b158015610d7f57600080fd5b505af1158015610d93573d6000803e3d6000fd5b505050506040513d6020811015610da957600080fd5b50511515610e01576040805160e560020a62461bcd02815260206004820152601960248201527f736574746c656d656e74206e6f74207265676973746572656400000000000000604482015290519081900360640190fd5b60048054604080517fa6065c9600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8a169381019390935251600160a060020a039091169163a6065c969160248083019260209291908290030181600087803b158015610e7257600080fd5b505af1158015610e86573d6000803e3d6000fd5b505050506040513d6020811015610e9c57600080fd5b50516040517f472e1910000000000000000000000000000000000000000000000000000000008152600160a060020a038481166004830190815260448301879052606060248401908152606484018990529394509084169263472e19109286928a928a928a929190608401858580828437820191505095505050505050602060405180830381600087803b158015610f3357600080fd5b505af1158015610f47573d6000803e3d6000fd5b505050506040513d6020811015610f5d57600080fd5b50511515610fb5576040805160e560020a62461bcd02815260206004820152601860248201527f696e76616c69642062726f6b6572207369676e61747572650000000000000000604482015290519081900360640190fd5b6040805160e0810182526001808252600160a060020a038516602080840191909152600083850181905267ffffffffffffffff8b166060850152600554830160808501524360a085015260c0840181905287815260069091529290922081518154929391929091839160ff19169083600381111561102f57fe5b021790555060208201518154600160a060020a039182166101000274ffffffffffffffffffffffffffffffffffffffff0019909116178255604083015160018084018054606087015167ffffffffffffffff1674010000000000000000000000000000000000000000027fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff9490951673ffffffffffffffffffffffffffffffffffffffff1990911617929092169290921790556080830151600283015560a0830151600383015560c0909201516004909101556005805491820181556000527f036b6384b5eca791c62761152d0c79bb0604c104a5fb6f4eb0703f3154bb3db001929092555050505050565b60009081526006602052604090205460ff1690565b600054600160a060020a0316331461116757600080fd5b60035460408051600160a060020a039283168152918316602083015280517ff9f6dd5c784f63cc27c1079c73574a73485a6c2e7f7e2181c5eb2be8c693cfb79281900390910190a16003805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b60009081526006602052604090206004015490565b6000908152600660205260409020546101009004600160a060020a031690565b60009081526006602052604090206002015490565b600454600160a060020a031681565b600054600160a060020a0316331461124e57600080fd5b611257816112e7565b50565b60018054604080516020600284861615610100026000190190941693909304601f810184900484028201840190925281815292918301828280156112df5780601f106112b4576101008083540402835291602001916112df565b820191906000526020600020905b8154815290600101906020018083116112c257829003601f168201915b505050505081565b600160a060020a03811615156112fc57600080fd5b60008054604051600160a060020a03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a36000805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03929092169190911790555600a165627a7a7230582051b56e26f7dae3c9e51492ee15b23bcf52a37bffb0601256a30cde6c0213be8b0029`

// DeployOrderbook deploys a new Ethereum contract, binding an instance of Orderbook to it.
func DeployOrderbook(auth *bind.TransactOpts, backend bind.ContractBackend, _VERSION string, _renAddress common.Address, _darknodeRegistry common.Address, _settlementRegistry common.Address) (common.Address, *types.Transaction, *Orderbook, error) {
	parsed, err := abi.JSON(strings.NewReader(OrderbookABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(OrderbookBin), backend, _VERSION, _renAddress, _darknodeRegistry, _settlementRegistry)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Orderbook{OrderbookCaller: OrderbookCaller{contract: contract}, OrderbookTransactor: OrderbookTransactor{contract: contract}, OrderbookFilterer: OrderbookFilterer{contract: contract}}, nil
}

// Orderbook is an auto generated Go binding around an Ethereum contract.
type Orderbook struct {
	OrderbookCaller     // Read-only binding to the contract
	OrderbookTransactor // Write-only binding to the contract
	OrderbookFilterer   // Log filterer for contract events
}

// OrderbookCaller is an auto generated read-only Go binding around an Ethereum contract.
type OrderbookCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OrderbookTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OrderbookTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OrderbookFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OrderbookFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OrderbookSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OrderbookSession struct {
	Contract     *Orderbook        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OrderbookCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OrderbookCallerSession struct {
	Contract *OrderbookCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// OrderbookTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OrderbookTransactorSession struct {
	Contract     *OrderbookTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// OrderbookRaw is an auto generated low-level Go binding around an Ethereum contract.
type OrderbookRaw struct {
	Contract *Orderbook // Generic contract binding to access the raw methods on
}

// OrderbookCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OrderbookCallerRaw struct {
	Contract *OrderbookCaller // Generic read-only contract binding to access the raw methods on
}

// OrderbookTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OrderbookTransactorRaw struct {
	Contract *OrderbookTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOrderbook creates a new instance of Orderbook, bound to a specific deployed contract.
func NewOrderbook(address common.Address, backend bind.ContractBackend) (*Orderbook, error) {
	contract, err := bindOrderbook(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Orderbook{OrderbookCaller: OrderbookCaller{contract: contract}, OrderbookTransactor: OrderbookTransactor{contract: contract}, OrderbookFilterer: OrderbookFilterer{contract: contract}}, nil
}

// NewOrderbookCaller creates a new read-only instance of Orderbook, bound to a specific deployed contract.
func NewOrderbookCaller(address common.Address, caller bind.ContractCaller) (*OrderbookCaller, error) {
	contract, err := bindOrderbook(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OrderbookCaller{contract: contract}, nil
}

// NewOrderbookTransactor creates a new write-only instance of Orderbook, bound to a specific deployed contract.
func NewOrderbookTransactor(address common.Address, transactor bind.ContractTransactor) (*OrderbookTransactor, error) {
	contract, err := bindOrderbook(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OrderbookTransactor{contract: contract}, nil
}

// NewOrderbookFilterer creates a new log filterer instance of Orderbook, bound to a specific deployed contract.
func NewOrderbookFilterer(address common.Address, filterer bind.ContractFilterer) (*OrderbookFilterer, error) {
	contract, err := bindOrderbook(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OrderbookFilterer{contract: contract}, nil
}

// bindOrderbook binds a generic wrapper to an already deployed contract.
func bindOrderbook(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OrderbookABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Orderbook *OrderbookRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Orderbook.Contract.OrderbookCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Orderbook *OrderbookRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Orderbook.Contract.OrderbookTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Orderbook *OrderbookRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Orderbook.Contract.OrderbookTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Orderbook *OrderbookCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Orderbook.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Orderbook *OrderbookTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Orderbook.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Orderbook *OrderbookTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Orderbook.Contract.contract.Transact(opts, method, params...)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_Orderbook *OrderbookCaller) VERSION(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Orderbook.contract.Call(opts, out, "VERSION")
	return *ret0, err
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_Orderbook *OrderbookSession) VERSION() (string, error) {
	return _Orderbook.Contract.VERSION(&_Orderbook.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_Orderbook *OrderbookCallerSession) VERSION() (string, error) {
	return _Orderbook.Contract.VERSION(&_Orderbook.CallOpts)
}

// DarknodeRegistry is a free data retrieval call binding the contract method 0x9e45e0d0.
//
// Solidity: function darknodeRegistry() constant returns(address)
func (_Orderbook *OrderbookCaller) DarknodeRegistry(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Orderbook.contract.Call(opts, out, "darknodeRegistry")
	return *ret0, err
}

// DarknodeRegistry is a free data retrieval call binding the contract method 0x9e45e0d0.
//
// Solidity: function darknodeRegistry() constant returns(address)
func (_Orderbook *OrderbookSession) DarknodeRegistry() (common.Address, error) {
	return _Orderbook.Contract.DarknodeRegistry(&_Orderbook.CallOpts)
}

// DarknodeRegistry is a free data retrieval call binding the contract method 0x9e45e0d0.
//
// Solidity: function darknodeRegistry() constant returns(address)
func (_Orderbook *OrderbookCallerSession) DarknodeRegistry() (common.Address, error) {
	return _Orderbook.Contract.DarknodeRegistry(&_Orderbook.CallOpts)
}

// GetOrders is a free data retrieval call binding the contract method 0x8f72fc77.
//
// Solidity: function getOrders(_offset uint256, _limit uint256) constant returns(bytes32[], address[], uint8[])
func (_Orderbook *OrderbookCaller) GetOrders(opts *bind.CallOpts, _offset *big.Int, _limit *big.Int) ([][32]byte, []common.Address, []uint8, error) {
	var (
		ret0 = new([][32]byte)
		ret1 = new([]common.Address)
		ret2 = new([]uint8)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
	}
	err := _Orderbook.contract.Call(opts, out, "getOrders", _offset, _limit)
	return *ret0, *ret1, *ret2, err
}

// GetOrders is a free data retrieval call binding the contract method 0x8f72fc77.
//
// Solidity: function getOrders(_offset uint256, _limit uint256) constant returns(bytes32[], address[], uint8[])
func (_Orderbook *OrderbookSession) GetOrders(_offset *big.Int, _limit *big.Int) ([][32]byte, []common.Address, []uint8, error) {
	return _Orderbook.Contract.GetOrders(&_Orderbook.CallOpts, _offset, _limit)
}

// GetOrders is a free data retrieval call binding the contract method 0x8f72fc77.
//
// Solidity: function getOrders(_offset uint256, _limit uint256) constant returns(bytes32[], address[], uint8[])
func (_Orderbook *OrderbookCallerSession) GetOrders(_offset *big.Int, _limit *big.Int) ([][32]byte, []common.Address, []uint8, error) {
	return _Orderbook.Contract.GetOrders(&_Orderbook.CallOpts, _offset, _limit)
}

// OrderBlockNumber is a free data retrieval call binding the contract method 0x89895d53.
//
// Solidity: function orderBlockNumber(_orderID bytes32) constant returns(uint256)
func (_Orderbook *OrderbookCaller) OrderBlockNumber(opts *bind.CallOpts, _orderID [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Orderbook.contract.Call(opts, out, "orderBlockNumber", _orderID)
	return *ret0, err
}

// OrderBlockNumber is a free data retrieval call binding the contract method 0x89895d53.
//
// Solidity: function orderBlockNumber(_orderID bytes32) constant returns(uint256)
func (_Orderbook *OrderbookSession) OrderBlockNumber(_orderID [32]byte) (*big.Int, error) {
	return _Orderbook.Contract.OrderBlockNumber(&_Orderbook.CallOpts, _orderID)
}

// OrderBlockNumber is a free data retrieval call binding the contract method 0x89895d53.
//
// Solidity: function orderBlockNumber(_orderID bytes32) constant returns(uint256)
func (_Orderbook *OrderbookCallerSession) OrderBlockNumber(_orderID [32]byte) (*big.Int, error) {
	return _Orderbook.Contract.OrderBlockNumber(&_Orderbook.CallOpts, _orderID)
}

// OrderConfirmer is a free data retrieval call binding the contract method 0x1107c3f7.
//
// Solidity: function orderConfirmer(_orderID bytes32) constant returns(address)
func (_Orderbook *OrderbookCaller) OrderConfirmer(opts *bind.CallOpts, _orderID [32]byte) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Orderbook.contract.Call(opts, out, "orderConfirmer", _orderID)
	return *ret0, err
}

// OrderConfirmer is a free data retrieval call binding the contract method 0x1107c3f7.
//
// Solidity: function orderConfirmer(_orderID bytes32) constant returns(address)
func (_Orderbook *OrderbookSession) OrderConfirmer(_orderID [32]byte) (common.Address, error) {
	return _Orderbook.Contract.OrderConfirmer(&_Orderbook.CallOpts, _orderID)
}

// OrderConfirmer is a free data retrieval call binding the contract method 0x1107c3f7.
//
// Solidity: function orderConfirmer(_orderID bytes32) constant returns(address)
func (_Orderbook *OrderbookCallerSession) OrderConfirmer(_orderID [32]byte) (common.Address, error) {
	return _Orderbook.Contract.OrderConfirmer(&_Orderbook.CallOpts, _orderID)
}

// OrderDepth is a free data retrieval call binding the contract method 0xa188fcb8.
//
// Solidity: function orderDepth(_orderID bytes32) constant returns(uint256)
func (_Orderbook *OrderbookCaller) OrderDepth(opts *bind.CallOpts, _orderID [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Orderbook.contract.Call(opts, out, "orderDepth", _orderID)
	return *ret0, err
}

// OrderDepth is a free data retrieval call binding the contract method 0xa188fcb8.
//
// Solidity: function orderDepth(_orderID bytes32) constant returns(uint256)
func (_Orderbook *OrderbookSession) OrderDepth(_orderID [32]byte) (*big.Int, error) {
	return _Orderbook.Contract.OrderDepth(&_Orderbook.CallOpts, _orderID)
}

// OrderDepth is a free data retrieval call binding the contract method 0xa188fcb8.
//
// Solidity: function orderDepth(_orderID bytes32) constant returns(uint256)
func (_Orderbook *OrderbookCallerSession) OrderDepth(_orderID [32]byte) (*big.Int, error) {
	return _Orderbook.Contract.OrderDepth(&_Orderbook.CallOpts, _orderID)
}

// OrderMatch is a free data retrieval call binding the contract method 0xaf3e8a40.
//
// Solidity: function orderMatch(_orderID bytes32) constant returns(bytes32)
func (_Orderbook *OrderbookCaller) OrderMatch(opts *bind.CallOpts, _orderID [32]byte) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Orderbook.contract.Call(opts, out, "orderMatch", _orderID)
	return *ret0, err
}

// OrderMatch is a free data retrieval call binding the contract method 0xaf3e8a40.
//
// Solidity: function orderMatch(_orderID bytes32) constant returns(bytes32)
func (_Orderbook *OrderbookSession) OrderMatch(_orderID [32]byte) ([32]byte, error) {
	return _Orderbook.Contract.OrderMatch(&_Orderbook.CallOpts, _orderID)
}

// OrderMatch is a free data retrieval call binding the contract method 0xaf3e8a40.
//
// Solidity: function orderMatch(_orderID bytes32) constant returns(bytes32)
func (_Orderbook *OrderbookCallerSession) OrderMatch(_orderID [32]byte) ([32]byte, error) {
	return _Orderbook.Contract.OrderMatch(&_Orderbook.CallOpts, _orderID)
}

// OrderPriority is a free data retrieval call binding the contract method 0xb248e4e1.
//
// Solidity: function orderPriority(_orderID bytes32) constant returns(uint256)
func (_Orderbook *OrderbookCaller) OrderPriority(opts *bind.CallOpts, _orderID [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Orderbook.contract.Call(opts, out, "orderPriority", _orderID)
	return *ret0, err
}

// OrderPriority is a free data retrieval call binding the contract method 0xb248e4e1.
//
// Solidity: function orderPriority(_orderID bytes32) constant returns(uint256)
func (_Orderbook *OrderbookSession) OrderPriority(_orderID [32]byte) (*big.Int, error) {
	return _Orderbook.Contract.OrderPriority(&_Orderbook.CallOpts, _orderID)
}

// OrderPriority is a free data retrieval call binding the contract method 0xb248e4e1.
//
// Solidity: function orderPriority(_orderID bytes32) constant returns(uint256)
func (_Orderbook *OrderbookCallerSession) OrderPriority(_orderID [32]byte) (*big.Int, error) {
	return _Orderbook.Contract.OrderPriority(&_Orderbook.CallOpts, _orderID)
}

// OrderState is a free data retrieval call binding the contract method 0xaab14d04.
//
// Solidity: function orderState(_orderID bytes32) constant returns(uint8)
func (_Orderbook *OrderbookCaller) OrderState(opts *bind.CallOpts, _orderID [32]byte) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Orderbook.contract.Call(opts, out, "orderState", _orderID)
	return *ret0, err
}

// OrderState is a free data retrieval call binding the contract method 0xaab14d04.
//
// Solidity: function orderState(_orderID bytes32) constant returns(uint8)
func (_Orderbook *OrderbookSession) OrderState(_orderID [32]byte) (uint8, error) {
	return _Orderbook.Contract.OrderState(&_Orderbook.CallOpts, _orderID)
}

// OrderState is a free data retrieval call binding the contract method 0xaab14d04.
//
// Solidity: function orderState(_orderID bytes32) constant returns(uint8)
func (_Orderbook *OrderbookCallerSession) OrderState(_orderID [32]byte) (uint8, error) {
	return _Orderbook.Contract.OrderState(&_Orderbook.CallOpts, _orderID)
}

// OrderTrader is a free data retrieval call binding the contract method 0xb1a08010.
//
// Solidity: function orderTrader(_orderID bytes32) constant returns(address)
func (_Orderbook *OrderbookCaller) OrderTrader(opts *bind.CallOpts, _orderID [32]byte) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Orderbook.contract.Call(opts, out, "orderTrader", _orderID)
	return *ret0, err
}

// OrderTrader is a free data retrieval call binding the contract method 0xb1a08010.
//
// Solidity: function orderTrader(_orderID bytes32) constant returns(address)
func (_Orderbook *OrderbookSession) OrderTrader(_orderID [32]byte) (common.Address, error) {
	return _Orderbook.Contract.OrderTrader(&_Orderbook.CallOpts, _orderID)
}

// OrderTrader is a free data retrieval call binding the contract method 0xb1a08010.
//
// Solidity: function orderTrader(_orderID bytes32) constant returns(address)
func (_Orderbook *OrderbookCallerSession) OrderTrader(_orderID [32]byte) (common.Address, error) {
	return _Orderbook.Contract.OrderTrader(&_Orderbook.CallOpts, _orderID)
}

// Orders is a free data retrieval call binding the contract method 0x9c3f1e90.
//
// Solidity: function orders( bytes32) constant returns(state uint8, trader address, confirmer address, settlementID uint64, priority uint256, blockNumber uint256, matchedOrder bytes32)
func (_Orderbook *OrderbookCaller) Orders(opts *bind.CallOpts, arg0 [32]byte) (struct {
	State        uint8
	Trader       common.Address
	Confirmer    common.Address
	SettlementID uint64
	Priority     *big.Int
	BlockNumber  *big.Int
	MatchedOrder [32]byte
}, error) {
	ret := new(struct {
		State        uint8
		Trader       common.Address
		Confirmer    common.Address
		SettlementID uint64
		Priority     *big.Int
		BlockNumber  *big.Int
		MatchedOrder [32]byte
	})
	out := ret
	err := _Orderbook.contract.Call(opts, out, "orders", arg0)
	return *ret, err
}

// Orders is a free data retrieval call binding the contract method 0x9c3f1e90.
//
// Solidity: function orders( bytes32) constant returns(state uint8, trader address, confirmer address, settlementID uint64, priority uint256, blockNumber uint256, matchedOrder bytes32)
func (_Orderbook *OrderbookSession) Orders(arg0 [32]byte) (struct {
	State        uint8
	Trader       common.Address
	Confirmer    common.Address
	SettlementID uint64
	Priority     *big.Int
	BlockNumber  *big.Int
	MatchedOrder [32]byte
}, error) {
	return _Orderbook.Contract.Orders(&_Orderbook.CallOpts, arg0)
}

// Orders is a free data retrieval call binding the contract method 0x9c3f1e90.
//
// Solidity: function orders( bytes32) constant returns(state uint8, trader address, confirmer address, settlementID uint64, priority uint256, blockNumber uint256, matchedOrder bytes32)
func (_Orderbook *OrderbookCallerSession) Orders(arg0 [32]byte) (struct {
	State        uint8
	Trader       common.Address
	Confirmer    common.Address
	SettlementID uint64
	Priority     *big.Int
	BlockNumber  *big.Int
	MatchedOrder [32]byte
}, error) {
	return _Orderbook.Contract.Orders(&_Orderbook.CallOpts, arg0)
}

// OrdersCount is a free data retrieval call binding the contract method 0x35daa731.
//
// Solidity: function ordersCount() constant returns(uint256)
func (_Orderbook *OrderbookCaller) OrdersCount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Orderbook.contract.Call(opts, out, "ordersCount")
	return *ret0, err
}

// OrdersCount is a free data retrieval call binding the contract method 0x35daa731.
//
// Solidity: function ordersCount() constant returns(uint256)
func (_Orderbook *OrderbookSession) OrdersCount() (*big.Int, error) {
	return _Orderbook.Contract.OrdersCount(&_Orderbook.CallOpts)
}

// OrdersCount is a free data retrieval call binding the contract method 0x35daa731.
//
// Solidity: function ordersCount() constant returns(uint256)
func (_Orderbook *OrderbookCallerSession) OrdersCount() (*big.Int, error) {
	return _Orderbook.Contract.OrdersCount(&_Orderbook.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Orderbook *OrderbookCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Orderbook.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Orderbook *OrderbookSession) Owner() (common.Address, error) {
	return _Orderbook.Contract.Owner(&_Orderbook.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Orderbook *OrderbookCallerSession) Owner() (common.Address, error) {
	return _Orderbook.Contract.Owner(&_Orderbook.CallOpts)
}

// Ren is a free data retrieval call binding the contract method 0x8a9b4067.
//
// Solidity: function ren() constant returns(address)
func (_Orderbook *OrderbookCaller) Ren(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Orderbook.contract.Call(opts, out, "ren")
	return *ret0, err
}

// Ren is a free data retrieval call binding the contract method 0x8a9b4067.
//
// Solidity: function ren() constant returns(address)
func (_Orderbook *OrderbookSession) Ren() (common.Address, error) {
	return _Orderbook.Contract.Ren(&_Orderbook.CallOpts)
}

// Ren is a free data retrieval call binding the contract method 0x8a9b4067.
//
// Solidity: function ren() constant returns(address)
func (_Orderbook *OrderbookCallerSession) Ren() (common.Address, error) {
	return _Orderbook.Contract.Ren(&_Orderbook.CallOpts)
}

// SettlementRegistry is a free data retrieval call binding the contract method 0xd09812e1.
//
// Solidity: function settlementRegistry() constant returns(address)
func (_Orderbook *OrderbookCaller) SettlementRegistry(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Orderbook.contract.Call(opts, out, "settlementRegistry")
	return *ret0, err
}

// SettlementRegistry is a free data retrieval call binding the contract method 0xd09812e1.
//
// Solidity: function settlementRegistry() constant returns(address)
func (_Orderbook *OrderbookSession) SettlementRegistry() (common.Address, error) {
	return _Orderbook.Contract.SettlementRegistry(&_Orderbook.CallOpts)
}

// SettlementRegistry is a free data retrieval call binding the contract method 0xd09812e1.
//
// Solidity: function settlementRegistry() constant returns(address)
func (_Orderbook *OrderbookCallerSession) SettlementRegistry() (common.Address, error) {
	return _Orderbook.Contract.SettlementRegistry(&_Orderbook.CallOpts)
}

// CancelOrder is a paid mutator transaction binding the contract method 0x7489ec23.
//
// Solidity: function cancelOrder(_orderID bytes32) returns()
func (_Orderbook *OrderbookTransactor) CancelOrder(opts *bind.TransactOpts, _orderID [32]byte) (*types.Transaction, error) {
	return _Orderbook.contract.Transact(opts, "cancelOrder", _orderID)
}

// CancelOrder is a paid mutator transaction binding the contract method 0x7489ec23.
//
// Solidity: function cancelOrder(_orderID bytes32) returns()
func (_Orderbook *OrderbookSession) CancelOrder(_orderID [32]byte) (*types.Transaction, error) {
	return _Orderbook.Contract.CancelOrder(&_Orderbook.TransactOpts, _orderID)
}

// CancelOrder is a paid mutator transaction binding the contract method 0x7489ec23.
//
// Solidity: function cancelOrder(_orderID bytes32) returns()
func (_Orderbook *OrderbookTransactorSession) CancelOrder(_orderID [32]byte) (*types.Transaction, error) {
	return _Orderbook.Contract.CancelOrder(&_Orderbook.TransactOpts, _orderID)
}

// ConfirmOrder is a paid mutator transaction binding the contract method 0x2a0401f0.
//
// Solidity: function confirmOrder(_orderID bytes32, _matchedOrderID bytes32) returns()
func (_Orderbook *OrderbookTransactor) ConfirmOrder(opts *bind.TransactOpts, _orderID [32]byte, _matchedOrderID [32]byte) (*types.Transaction, error) {
	return _Orderbook.contract.Transact(opts, "confirmOrder", _orderID, _matchedOrderID)
}

// ConfirmOrder is a paid mutator transaction binding the contract method 0x2a0401f0.
//
// Solidity: function confirmOrder(_orderID bytes32, _matchedOrderID bytes32) returns()
func (_Orderbook *OrderbookSession) ConfirmOrder(_orderID [32]byte, _matchedOrderID [32]byte) (*types.Transaction, error) {
	return _Orderbook.Contract.ConfirmOrder(&_Orderbook.TransactOpts, _orderID, _matchedOrderID)
}

// ConfirmOrder is a paid mutator transaction binding the contract method 0x2a0401f0.
//
// Solidity: function confirmOrder(_orderID bytes32, _matchedOrderID bytes32) returns()
func (_Orderbook *OrderbookTransactorSession) ConfirmOrder(_orderID [32]byte, _matchedOrderID [32]byte) (*types.Transaction, error) {
	return _Orderbook.Contract.ConfirmOrder(&_Orderbook.TransactOpts, _orderID, _matchedOrderID)
}

// OpenOrder is a paid mutator transaction binding the contract method 0xa1ce7e03.
//
// Solidity: function openOrder(_settlementID uint64, _signature bytes, _orderID bytes32) returns()
func (_Orderbook *OrderbookTransactor) OpenOrder(opts *bind.TransactOpts, _settlementID uint64, _signature []byte, _orderID [32]byte) (*types.Transaction, error) {
	return _Orderbook.contract.Transact(opts, "openOrder", _settlementID, _signature, _orderID)
}

// OpenOrder is a paid mutator transaction binding the contract method 0xa1ce7e03.
//
// Solidity: function openOrder(_settlementID uint64, _signature bytes, _orderID bytes32) returns()
func (_Orderbook *OrderbookSession) OpenOrder(_settlementID uint64, _signature []byte, _orderID [32]byte) (*types.Transaction, error) {
	return _Orderbook.Contract.OpenOrder(&_Orderbook.TransactOpts, _settlementID, _signature, _orderID)
}

// OpenOrder is a paid mutator transaction binding the contract method 0xa1ce7e03.
//
// Solidity: function openOrder(_settlementID uint64, _signature bytes, _orderID bytes32) returns()
func (_Orderbook *OrderbookTransactorSession) OpenOrder(_settlementID uint64, _signature []byte, _orderID [32]byte) (*types.Transaction, error) {
	return _Orderbook.Contract.OpenOrder(&_Orderbook.TransactOpts, _settlementID, _signature, _orderID)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Orderbook *OrderbookTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Orderbook.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Orderbook *OrderbookSession) RenounceOwnership() (*types.Transaction, error) {
	return _Orderbook.Contract.RenounceOwnership(&_Orderbook.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Orderbook *OrderbookTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Orderbook.Contract.RenounceOwnership(&_Orderbook.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_Orderbook *OrderbookTransactor) TransferOwnership(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _Orderbook.contract.Transact(opts, "transferOwnership", _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_Orderbook *OrderbookSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _Orderbook.Contract.TransferOwnership(&_Orderbook.TransactOpts, _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_Orderbook *OrderbookTransactorSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _Orderbook.Contract.TransferOwnership(&_Orderbook.TransactOpts, _newOwner)
}

// UpdateDarknodeRegistry is a paid mutator transaction binding the contract method 0xaaff096d.
//
// Solidity: function updateDarknodeRegistry(_newDarknodeRegistry address) returns()
func (_Orderbook *OrderbookTransactor) UpdateDarknodeRegistry(opts *bind.TransactOpts, _newDarknodeRegistry common.Address) (*types.Transaction, error) {
	return _Orderbook.contract.Transact(opts, "updateDarknodeRegistry", _newDarknodeRegistry)
}

// UpdateDarknodeRegistry is a paid mutator transaction binding the contract method 0xaaff096d.
//
// Solidity: function updateDarknodeRegistry(_newDarknodeRegistry address) returns()
func (_Orderbook *OrderbookSession) UpdateDarknodeRegistry(_newDarknodeRegistry common.Address) (*types.Transaction, error) {
	return _Orderbook.Contract.UpdateDarknodeRegistry(&_Orderbook.TransactOpts, _newDarknodeRegistry)
}

// UpdateDarknodeRegistry is a paid mutator transaction binding the contract method 0xaaff096d.
//
// Solidity: function updateDarknodeRegistry(_newDarknodeRegistry address) returns()
func (_Orderbook *OrderbookTransactorSession) UpdateDarknodeRegistry(_newDarknodeRegistry common.Address) (*types.Transaction, error) {
	return _Orderbook.Contract.UpdateDarknodeRegistry(&_Orderbook.TransactOpts, _newDarknodeRegistry)
}

// OrderbookLogDarknodeRegistryUpdatedIterator is returned from FilterLogDarknodeRegistryUpdated and is used to iterate over the raw logs and unpacked data for LogDarknodeRegistryUpdated events raised by the Orderbook contract.
type OrderbookLogDarknodeRegistryUpdatedIterator struct {
	Event *OrderbookLogDarknodeRegistryUpdated // Event containing the contract specifics and raw log

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
func (it *OrderbookLogDarknodeRegistryUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OrderbookLogDarknodeRegistryUpdated)
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
		it.Event = new(OrderbookLogDarknodeRegistryUpdated)
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
func (it *OrderbookLogDarknodeRegistryUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OrderbookLogDarknodeRegistryUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OrderbookLogDarknodeRegistryUpdated represents a LogDarknodeRegistryUpdated event raised by the Orderbook contract.
type OrderbookLogDarknodeRegistryUpdated struct {
	PreviousDarknodeRegistry common.Address
	NextDarknodeRegistry     common.Address
	Raw                      types.Log // Blockchain specific contextual infos
}

// FilterLogDarknodeRegistryUpdated is a free log retrieval operation binding the contract event 0xf9f6dd5c784f63cc27c1079c73574a73485a6c2e7f7e2181c5eb2be8c693cfb7.
//
// Solidity: e LogDarknodeRegistryUpdated(previousDarknodeRegistry address, nextDarknodeRegistry address)
func (_Orderbook *OrderbookFilterer) FilterLogDarknodeRegistryUpdated(opts *bind.FilterOpts) (*OrderbookLogDarknodeRegistryUpdatedIterator, error) {

	logs, sub, err := _Orderbook.contract.FilterLogs(opts, "LogDarknodeRegistryUpdated")
	if err != nil {
		return nil, err
	}
	return &OrderbookLogDarknodeRegistryUpdatedIterator{contract: _Orderbook.contract, event: "LogDarknodeRegistryUpdated", logs: logs, sub: sub}, nil
}

// WatchLogDarknodeRegistryUpdated is a free log subscription operation binding the contract event 0xf9f6dd5c784f63cc27c1079c73574a73485a6c2e7f7e2181c5eb2be8c693cfb7.
//
// Solidity: e LogDarknodeRegistryUpdated(previousDarknodeRegistry address, nextDarknodeRegistry address)
func (_Orderbook *OrderbookFilterer) WatchLogDarknodeRegistryUpdated(opts *bind.WatchOpts, sink chan<- *OrderbookLogDarknodeRegistryUpdated) (event.Subscription, error) {

	logs, sub, err := _Orderbook.contract.WatchLogs(opts, "LogDarknodeRegistryUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OrderbookLogDarknodeRegistryUpdated)
				if err := _Orderbook.contract.UnpackLog(event, "LogDarknodeRegistryUpdated", log); err != nil {
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

// OrderbookLogFeeUpdatedIterator is returned from FilterLogFeeUpdated and is used to iterate over the raw logs and unpacked data for LogFeeUpdated events raised by the Orderbook contract.
type OrderbookLogFeeUpdatedIterator struct {
	Event *OrderbookLogFeeUpdated // Event containing the contract specifics and raw log

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
func (it *OrderbookLogFeeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OrderbookLogFeeUpdated)
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
		it.Event = new(OrderbookLogFeeUpdated)
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
func (it *OrderbookLogFeeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OrderbookLogFeeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OrderbookLogFeeUpdated represents a LogFeeUpdated event raised by the Orderbook contract.
type OrderbookLogFeeUpdated struct {
	PreviousFee *big.Int
	NextFee     *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterLogFeeUpdated is a free log retrieval operation binding the contract event 0x6308fd23d4b466bc07d9c9ef31631f5b96c7bac2efcb3d214fe3cc678e7ae00a.
//
// Solidity: e LogFeeUpdated(previousFee uint256, nextFee uint256)
func (_Orderbook *OrderbookFilterer) FilterLogFeeUpdated(opts *bind.FilterOpts) (*OrderbookLogFeeUpdatedIterator, error) {

	logs, sub, err := _Orderbook.contract.FilterLogs(opts, "LogFeeUpdated")
	if err != nil {
		return nil, err
	}
	return &OrderbookLogFeeUpdatedIterator{contract: _Orderbook.contract, event: "LogFeeUpdated", logs: logs, sub: sub}, nil
}

// WatchLogFeeUpdated is a free log subscription operation binding the contract event 0x6308fd23d4b466bc07d9c9ef31631f5b96c7bac2efcb3d214fe3cc678e7ae00a.
//
// Solidity: e LogFeeUpdated(previousFee uint256, nextFee uint256)
func (_Orderbook *OrderbookFilterer) WatchLogFeeUpdated(opts *bind.WatchOpts, sink chan<- *OrderbookLogFeeUpdated) (event.Subscription, error) {

	logs, sub, err := _Orderbook.contract.WatchLogs(opts, "LogFeeUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OrderbookLogFeeUpdated)
				if err := _Orderbook.contract.UnpackLog(event, "LogFeeUpdated", log); err != nil {
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

// OrderbookOwnershipRenouncedIterator is returned from FilterOwnershipRenounced and is used to iterate over the raw logs and unpacked data for OwnershipRenounced events raised by the Orderbook contract.
type OrderbookOwnershipRenouncedIterator struct {
	Event *OrderbookOwnershipRenounced // Event containing the contract specifics and raw log

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
func (it *OrderbookOwnershipRenouncedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OrderbookOwnershipRenounced)
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
		it.Event = new(OrderbookOwnershipRenounced)
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
func (it *OrderbookOwnershipRenouncedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OrderbookOwnershipRenouncedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OrderbookOwnershipRenounced represents a OwnershipRenounced event raised by the Orderbook contract.
type OrderbookOwnershipRenounced struct {
	PreviousOwner common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipRenounced is a free log retrieval operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_Orderbook *OrderbookFilterer) FilterOwnershipRenounced(opts *bind.FilterOpts, previousOwner []common.Address) (*OrderbookOwnershipRenouncedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _Orderbook.contract.FilterLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return &OrderbookOwnershipRenouncedIterator{contract: _Orderbook.contract, event: "OwnershipRenounced", logs: logs, sub: sub}, nil
}

// WatchOwnershipRenounced is a free log subscription operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_Orderbook *OrderbookFilterer) WatchOwnershipRenounced(opts *bind.WatchOpts, sink chan<- *OrderbookOwnershipRenounced, previousOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _Orderbook.contract.WatchLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OrderbookOwnershipRenounced)
				if err := _Orderbook.contract.UnpackLog(event, "OwnershipRenounced", log); err != nil {
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

// OrderbookOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Orderbook contract.
type OrderbookOwnershipTransferredIterator struct {
	Event *OrderbookOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *OrderbookOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OrderbookOwnershipTransferred)
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
		it.Event = new(OrderbookOwnershipTransferred)
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
func (it *OrderbookOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OrderbookOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OrderbookOwnershipTransferred represents a OwnershipTransferred event raised by the Orderbook contract.
type OrderbookOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_Orderbook *OrderbookFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*OrderbookOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Orderbook.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &OrderbookOwnershipTransferredIterator{contract: _Orderbook.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_Orderbook *OrderbookFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OrderbookOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Orderbook.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OrderbookOwnershipTransferred)
				if err := _Orderbook.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// OwnableABI is the input ABI used to generate the binding from.
const OwnableABI = "[{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"}],\"name\":\"OwnershipRenounced\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"

// OwnableBin is the compiled bytecode used for deploying new contracts.
const OwnableBin = `0x608060405234801561001057600080fd5b5060008054600160a060020a0319163317905561020b806100326000396000f3006080604052600436106100565763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663715018a6811461005b5780638da5cb5b14610072578063f2fde38b146100a3575b600080fd5b34801561006757600080fd5b506100706100c4565b005b34801561007e57600080fd5b50610087610130565b60408051600160a060020a039092168252519081900360200190f35b3480156100af57600080fd5b50610070600160a060020a036004351661013f565b600054600160a060020a031633146100db57600080fd5b60008054604051600160a060020a03909116917ff8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c6482091a26000805473ffffffffffffffffffffffffffffffffffffffff19169055565b600054600160a060020a031681565b600054600160a060020a0316331461015657600080fd5b61015f81610162565b50565b600160a060020a038116151561017757600080fd5b60008054604051600160a060020a03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a36000805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03929092169190911790555600a165627a7a72305820d07f4d5221378e0ccfc5230ace30fcd5568360922dc2baa242decb0ba9afd4100029`

// DeployOwnable deploys a new Ethereum contract, binding an instance of Ownable to it.
func DeployOwnable(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Ownable, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnableABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(OwnableBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Ownable{OwnableCaller: OwnableCaller{contract: contract}, OwnableTransactor: OwnableTransactor{contract: contract}, OwnableFilterer: OwnableFilterer{contract: contract}}, nil
}

// Ownable is an auto generated Go binding around an Ethereum contract.
type Ownable struct {
	OwnableCaller     // Read-only binding to the contract
	OwnableTransactor // Write-only binding to the contract
	OwnableFilterer   // Log filterer for contract events
}

// OwnableCaller is an auto generated read-only Go binding around an Ethereum contract.
type OwnableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OwnableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OwnableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OwnableSession struct {
	Contract     *Ownable          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OwnableCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OwnableCallerSession struct {
	Contract *OwnableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// OwnableTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OwnableTransactorSession struct {
	Contract     *OwnableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// OwnableRaw is an auto generated low-level Go binding around an Ethereum contract.
type OwnableRaw struct {
	Contract *Ownable // Generic contract binding to access the raw methods on
}

// OwnableCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OwnableCallerRaw struct {
	Contract *OwnableCaller // Generic read-only contract binding to access the raw methods on
}

// OwnableTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OwnableTransactorRaw struct {
	Contract *OwnableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOwnable creates a new instance of Ownable, bound to a specific deployed contract.
func NewOwnable(address common.Address, backend bind.ContractBackend) (*Ownable, error) {
	contract, err := bindOwnable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Ownable{OwnableCaller: OwnableCaller{contract: contract}, OwnableTransactor: OwnableTransactor{contract: contract}, OwnableFilterer: OwnableFilterer{contract: contract}}, nil
}

// NewOwnableCaller creates a new read-only instance of Ownable, bound to a specific deployed contract.
func NewOwnableCaller(address common.Address, caller bind.ContractCaller) (*OwnableCaller, error) {
	contract, err := bindOwnable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OwnableCaller{contract: contract}, nil
}

// NewOwnableTransactor creates a new write-only instance of Ownable, bound to a specific deployed contract.
func NewOwnableTransactor(address common.Address, transactor bind.ContractTransactor) (*OwnableTransactor, error) {
	contract, err := bindOwnable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OwnableTransactor{contract: contract}, nil
}

// NewOwnableFilterer creates a new log filterer instance of Ownable, bound to a specific deployed contract.
func NewOwnableFilterer(address common.Address, filterer bind.ContractFilterer) (*OwnableFilterer, error) {
	contract, err := bindOwnable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OwnableFilterer{contract: contract}, nil
}

// bindOwnable binds a generic wrapper to an already deployed contract.
func bindOwnable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnableABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ownable *OwnableRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Ownable.Contract.OwnableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ownable *OwnableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ownable.Contract.OwnableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ownable *OwnableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ownable.Contract.OwnableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ownable *OwnableCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Ownable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ownable *OwnableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ownable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ownable *OwnableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ownable.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Ownable *OwnableCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Ownable.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Ownable *OwnableSession) Owner() (common.Address, error) {
	return _Ownable.Contract.Owner(&_Ownable.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Ownable *OwnableCallerSession) Owner() (common.Address, error) {
	return _Ownable.Contract.Owner(&_Ownable.CallOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ownable *OwnableTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ownable.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ownable *OwnableSession) RenounceOwnership() (*types.Transaction, error) {
	return _Ownable.Contract.RenounceOwnership(&_Ownable.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ownable *OwnableTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Ownable.Contract.RenounceOwnership(&_Ownable.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_Ownable *OwnableTransactor) TransferOwnership(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _Ownable.contract.Transact(opts, "transferOwnership", _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_Ownable *OwnableSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _Ownable.Contract.TransferOwnership(&_Ownable.TransactOpts, _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_Ownable *OwnableTransactorSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _Ownable.Contract.TransferOwnership(&_Ownable.TransactOpts, _newOwner)
}

// OwnableOwnershipRenouncedIterator is returned from FilterOwnershipRenounced and is used to iterate over the raw logs and unpacked data for OwnershipRenounced events raised by the Ownable contract.
type OwnableOwnershipRenouncedIterator struct {
	Event *OwnableOwnershipRenounced // Event containing the contract specifics and raw log

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
func (it *OwnableOwnershipRenouncedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnableOwnershipRenounced)
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
		it.Event = new(OwnableOwnershipRenounced)
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
func (it *OwnableOwnershipRenouncedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnableOwnershipRenouncedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnableOwnershipRenounced represents a OwnershipRenounced event raised by the Ownable contract.
type OwnableOwnershipRenounced struct {
	PreviousOwner common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipRenounced is a free log retrieval operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_Ownable *OwnableFilterer) FilterOwnershipRenounced(opts *bind.FilterOpts, previousOwner []common.Address) (*OwnableOwnershipRenouncedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _Ownable.contract.FilterLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return &OwnableOwnershipRenouncedIterator{contract: _Ownable.contract, event: "OwnershipRenounced", logs: logs, sub: sub}, nil
}

// WatchOwnershipRenounced is a free log subscription operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_Ownable *OwnableFilterer) WatchOwnershipRenounced(opts *bind.WatchOpts, sink chan<- *OwnableOwnershipRenounced, previousOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _Ownable.contract.WatchLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnableOwnershipRenounced)
				if err := _Ownable.contract.UnpackLog(event, "OwnershipRenounced", log); err != nil {
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

// OwnableOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Ownable contract.
type OwnableOwnershipTransferredIterator struct {
	Event *OwnableOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *OwnableOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnableOwnershipTransferred)
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
		it.Event = new(OwnableOwnershipTransferred)
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
func (it *OwnableOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnableOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnableOwnershipTransferred represents a OwnershipTransferred event raised by the Ownable contract.
type OwnableOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_Ownable *OwnableFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*OwnableOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Ownable.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &OwnableOwnershipTransferredIterator{contract: _Ownable.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_Ownable *OwnableFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OwnableOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Ownable.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnableOwnershipTransferred)
				if err := _Ownable.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// PausableABI is the input ABI used to generate the binding from.
const PausableABI = "[{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"}],\"name\":\"OwnershipRenounced\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"

// PausableBin is the compiled bytecode used for deploying new contracts.
const PausableBin = `0x608060405260008054600160a860020a031916331790556103c4806100256000396000f3006080604052600436106100775763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416633f4ba83a811461007c5780635c975abb14610093578063715018a6146100bc5780638456cb59146100d15780638da5cb5b146100e6578063f2fde38b14610117575b600080fd5b34801561008857600080fd5b50610091610138565b005b34801561009f57600080fd5b506100a86101bf565b604080519115158252519081900360200190f35b3480156100c857600080fd5b506100916101e0565b3480156100dd57600080fd5b5061009161024c565b3480156100f257600080fd5b506100fb6102e9565b60408051600160a060020a039092168252519081900360200190f35b34801561012357600080fd5b50610091600160a060020a03600435166102f8565b600054600160a060020a0316331461014f57600080fd5b60005474010000000000000000000000000000000000000000900460ff16151561017857600080fd5b6000805474ff0000000000000000000000000000000000000000191681556040517f7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b339190a1565b60005474010000000000000000000000000000000000000000900460ff1681565b600054600160a060020a031633146101f757600080fd5b60008054604051600160a060020a03909116917ff8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c6482091a26000805473ffffffffffffffffffffffffffffffffffffffff19169055565b600054600160a060020a0316331461026357600080fd5b60005474010000000000000000000000000000000000000000900460ff161561028b57600080fd5b6000805474ff00000000000000000000000000000000000000001916740100000000000000000000000000000000000000001781556040517f6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff6259190a1565b600054600160a060020a031681565b600054600160a060020a0316331461030f57600080fd5b6103188161031b565b50565b600160a060020a038116151561033057600080fd5b60008054604051600160a060020a03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a36000805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03929092169190911790555600a165627a7a723058207d0a7a42de164bf3e33f827ae05b6303ba2a10e703b3c9ca4bead9f8150273bf0029`

// DeployPausable deploys a new Ethereum contract, binding an instance of Pausable to it.
func DeployPausable(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Pausable, error) {
	parsed, err := abi.JSON(strings.NewReader(PausableABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(PausableBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Pausable{PausableCaller: PausableCaller{contract: contract}, PausableTransactor: PausableTransactor{contract: contract}, PausableFilterer: PausableFilterer{contract: contract}}, nil
}

// Pausable is an auto generated Go binding around an Ethereum contract.
type Pausable struct {
	PausableCaller     // Read-only binding to the contract
	PausableTransactor // Write-only binding to the contract
	PausableFilterer   // Log filterer for contract events
}

// PausableCaller is an auto generated read-only Go binding around an Ethereum contract.
type PausableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PausableTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PausableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PausableFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PausableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PausableSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PausableSession struct {
	Contract     *Pausable         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PausableCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PausableCallerSession struct {
	Contract *PausableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// PausableTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PausableTransactorSession struct {
	Contract     *PausableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// PausableRaw is an auto generated low-level Go binding around an Ethereum contract.
type PausableRaw struct {
	Contract *Pausable // Generic contract binding to access the raw methods on
}

// PausableCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PausableCallerRaw struct {
	Contract *PausableCaller // Generic read-only contract binding to access the raw methods on
}

// PausableTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PausableTransactorRaw struct {
	Contract *PausableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPausable creates a new instance of Pausable, bound to a specific deployed contract.
func NewPausable(address common.Address, backend bind.ContractBackend) (*Pausable, error) {
	contract, err := bindPausable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Pausable{PausableCaller: PausableCaller{contract: contract}, PausableTransactor: PausableTransactor{contract: contract}, PausableFilterer: PausableFilterer{contract: contract}}, nil
}

// NewPausableCaller creates a new read-only instance of Pausable, bound to a specific deployed contract.
func NewPausableCaller(address common.Address, caller bind.ContractCaller) (*PausableCaller, error) {
	contract, err := bindPausable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PausableCaller{contract: contract}, nil
}

// NewPausableTransactor creates a new write-only instance of Pausable, bound to a specific deployed contract.
func NewPausableTransactor(address common.Address, transactor bind.ContractTransactor) (*PausableTransactor, error) {
	contract, err := bindPausable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PausableTransactor{contract: contract}, nil
}

// NewPausableFilterer creates a new log filterer instance of Pausable, bound to a specific deployed contract.
func NewPausableFilterer(address common.Address, filterer bind.ContractFilterer) (*PausableFilterer, error) {
	contract, err := bindPausable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PausableFilterer{contract: contract}, nil
}

// bindPausable binds a generic wrapper to an already deployed contract.
func bindPausable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PausableABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Pausable *PausableRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Pausable.Contract.PausableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Pausable *PausableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Pausable.Contract.PausableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Pausable *PausableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Pausable.Contract.PausableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Pausable *PausableCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Pausable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Pausable *PausableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Pausable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Pausable *PausableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Pausable.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Pausable *PausableCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Pausable.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Pausable *PausableSession) Owner() (common.Address, error) {
	return _Pausable.Contract.Owner(&_Pausable.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Pausable *PausableCallerSession) Owner() (common.Address, error) {
	return _Pausable.Contract.Owner(&_Pausable.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() constant returns(bool)
func (_Pausable *PausableCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Pausable.contract.Call(opts, out, "paused")
	return *ret0, err
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() constant returns(bool)
func (_Pausable *PausableSession) Paused() (bool, error) {
	return _Pausable.Contract.Paused(&_Pausable.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() constant returns(bool)
func (_Pausable *PausableCallerSession) Paused() (bool, error) {
	return _Pausable.Contract.Paused(&_Pausable.CallOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Pausable *PausableTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Pausable.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Pausable *PausableSession) Pause() (*types.Transaction, error) {
	return _Pausable.Contract.Pause(&_Pausable.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Pausable *PausableTransactorSession) Pause() (*types.Transaction, error) {
	return _Pausable.Contract.Pause(&_Pausable.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Pausable *PausableTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Pausable.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Pausable *PausableSession) RenounceOwnership() (*types.Transaction, error) {
	return _Pausable.Contract.RenounceOwnership(&_Pausable.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Pausable *PausableTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Pausable.Contract.RenounceOwnership(&_Pausable.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_Pausable *PausableTransactor) TransferOwnership(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _Pausable.contract.Transact(opts, "transferOwnership", _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_Pausable *PausableSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _Pausable.Contract.TransferOwnership(&_Pausable.TransactOpts, _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_Pausable *PausableTransactorSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _Pausable.Contract.TransferOwnership(&_Pausable.TransactOpts, _newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Pausable *PausableTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Pausable.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Pausable *PausableSession) Unpause() (*types.Transaction, error) {
	return _Pausable.Contract.Unpause(&_Pausable.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Pausable *PausableTransactorSession) Unpause() (*types.Transaction, error) {
	return _Pausable.Contract.Unpause(&_Pausable.TransactOpts)
}

// PausableOwnershipRenouncedIterator is returned from FilterOwnershipRenounced and is used to iterate over the raw logs and unpacked data for OwnershipRenounced events raised by the Pausable contract.
type PausableOwnershipRenouncedIterator struct {
	Event *PausableOwnershipRenounced // Event containing the contract specifics and raw log

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
func (it *PausableOwnershipRenouncedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PausableOwnershipRenounced)
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
		it.Event = new(PausableOwnershipRenounced)
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
func (it *PausableOwnershipRenouncedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PausableOwnershipRenouncedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PausableOwnershipRenounced represents a OwnershipRenounced event raised by the Pausable contract.
type PausableOwnershipRenounced struct {
	PreviousOwner common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipRenounced is a free log retrieval operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_Pausable *PausableFilterer) FilterOwnershipRenounced(opts *bind.FilterOpts, previousOwner []common.Address) (*PausableOwnershipRenouncedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _Pausable.contract.FilterLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return &PausableOwnershipRenouncedIterator{contract: _Pausable.contract, event: "OwnershipRenounced", logs: logs, sub: sub}, nil
}

// WatchOwnershipRenounced is a free log subscription operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_Pausable *PausableFilterer) WatchOwnershipRenounced(opts *bind.WatchOpts, sink chan<- *PausableOwnershipRenounced, previousOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _Pausable.contract.WatchLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PausableOwnershipRenounced)
				if err := _Pausable.contract.UnpackLog(event, "OwnershipRenounced", log); err != nil {
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

// PausableOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Pausable contract.
type PausableOwnershipTransferredIterator struct {
	Event *PausableOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *PausableOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PausableOwnershipTransferred)
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
		it.Event = new(PausableOwnershipTransferred)
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
func (it *PausableOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PausableOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PausableOwnershipTransferred represents a OwnershipTransferred event raised by the Pausable contract.
type PausableOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_Pausable *PausableFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*PausableOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Pausable.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &PausableOwnershipTransferredIterator{contract: _Pausable.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_Pausable *PausableFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PausableOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Pausable.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PausableOwnershipTransferred)
				if err := _Pausable.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// PausablePauseIterator is returned from FilterPause and is used to iterate over the raw logs and unpacked data for Pause events raised by the Pausable contract.
type PausablePauseIterator struct {
	Event *PausablePause // Event containing the contract specifics and raw log

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
func (it *PausablePauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PausablePause)
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
		it.Event = new(PausablePause)
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
func (it *PausablePauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PausablePauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PausablePause represents a Pause event raised by the Pausable contract.
type PausablePause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterPause is a free log retrieval operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: e Pause()
func (_Pausable *PausableFilterer) FilterPause(opts *bind.FilterOpts) (*PausablePauseIterator, error) {

	logs, sub, err := _Pausable.contract.FilterLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return &PausablePauseIterator{contract: _Pausable.contract, event: "Pause", logs: logs, sub: sub}, nil
}

// WatchPause is a free log subscription operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: e Pause()
func (_Pausable *PausableFilterer) WatchPause(opts *bind.WatchOpts, sink chan<- *PausablePause) (event.Subscription, error) {

	logs, sub, err := _Pausable.contract.WatchLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PausablePause)
				if err := _Pausable.contract.UnpackLog(event, "Pause", log); err != nil {
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

// PausableUnpauseIterator is returned from FilterUnpause and is used to iterate over the raw logs and unpacked data for Unpause events raised by the Pausable contract.
type PausableUnpauseIterator struct {
	Event *PausableUnpause // Event containing the contract specifics and raw log

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
func (it *PausableUnpauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PausableUnpause)
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
		it.Event = new(PausableUnpause)
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
func (it *PausableUnpauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PausableUnpauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PausableUnpause represents a Unpause event raised by the Pausable contract.
type PausableUnpause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUnpause is a free log retrieval operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: e Unpause()
func (_Pausable *PausableFilterer) FilterUnpause(opts *bind.FilterOpts) (*PausableUnpauseIterator, error) {

	logs, sub, err := _Pausable.contract.FilterLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return &PausableUnpauseIterator{contract: _Pausable.contract, event: "Unpause", logs: logs, sub: sub}, nil
}

// WatchUnpause is a free log subscription operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: e Unpause()
func (_Pausable *PausableFilterer) WatchUnpause(opts *bind.WatchOpts, sink chan<- *PausableUnpause) (event.Subscription, error) {

	logs, sub, err := _Pausable.contract.WatchLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PausableUnpause)
				if err := _Pausable.contract.UnpackLog(event, "Unpause", log); err != nil {
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

// PausableTokenABI is the input ABI used to generate the binding from.
const PausableTokenABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseApproval\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseApproval\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"}],\"name\":\"OwnershipRenounced\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]"

// PausableTokenBin is the compiled bytecode used for deploying new contracts.
const PausableTokenBin = `0x608060405260038054600160a860020a03191633179055610a84806100256000396000f3006080604052600436106100cf5763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663095ea7b381146100d457806318160ddd1461010c57806323b872dd146101335780633f4ba83a1461015d5780635c975abb14610174578063661884631461018957806370a08231146101ad578063715018a6146101ce5780638456cb59146101e35780638da5cb5b146101f8578063a9059cbb14610229578063d73dd6231461024d578063dd62ed3e14610271578063f2fde38b14610298575b600080fd5b3480156100e057600080fd5b506100f8600160a060020a03600435166024356102b9565b604080519115158252519081900360200190f35b34801561011857600080fd5b506101216102e4565b60408051918252519081900360200190f35b34801561013f57600080fd5b506100f8600160a060020a03600435811690602435166044356102ea565b34801561016957600080fd5b50610172610317565b005b34801561018057600080fd5b506100f861038f565b34801561019557600080fd5b506100f8600160a060020a036004351660243561039f565b3480156101b957600080fd5b50610121600160a060020a03600435166103c3565b3480156101da57600080fd5b506101726103de565b3480156101ef57600080fd5b5061017261044c565b34801561020457600080fd5b5061020d6104c9565b60408051600160a060020a039092168252519081900360200190f35b34801561023557600080fd5b506100f8600160a060020a03600435166024356104d8565b34801561025957600080fd5b506100f8600160a060020a03600435166024356104fc565b34801561027d57600080fd5b50610121600160a060020a0360043581169060243516610520565b3480156102a457600080fd5b50610172600160a060020a036004351661054b565b60035460009060a060020a900460ff16156102d357600080fd5b6102dd838361056e565b9392505050565b60015490565b60035460009060a060020a900460ff161561030457600080fd5b61030f8484846105d4565b949350505050565b600354600160a060020a0316331461032e57600080fd5b60035460a060020a900460ff16151561034657600080fd5b6003805474ff0000000000000000000000000000000000000000191690556040517f7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b3390600090a1565b60035460a060020a900460ff1681565b60035460009060a060020a900460ff16156103b957600080fd5b6102dd838361074b565b600160a060020a031660009081526020819052604090205490565b600354600160a060020a031633146103f557600080fd5b600354604051600160a060020a03909116907ff8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c6482090600090a26003805473ffffffffffffffffffffffffffffffffffffffff19169055565b600354600160a060020a0316331461046357600080fd5b60035460a060020a900460ff161561047a57600080fd5b6003805474ff0000000000000000000000000000000000000000191660a060020a1790556040517f6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff62590600090a1565b600354600160a060020a031681565b60035460009060a060020a900460ff16156104f257600080fd5b6102dd838361083b565b60035460009060a060020a900460ff161561051657600080fd5b6102dd838361091c565b600160a060020a03918216600090815260026020908152604080832093909416825291909152205490565b600354600160a060020a0316331461056257600080fd5b61056b816109b5565b50565b336000818152600260209081526040808320600160a060020a038716808552908352818420869055815186815291519394909390927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925928290030190a350600192915050565b6000600160a060020a03831615156105eb57600080fd5b600160a060020a03841660009081526020819052604090205482111561061057600080fd5b600160a060020a038416600090815260026020908152604080832033845290915290205482111561064057600080fd5b600160a060020a038416600090815260208190526040902054610669908363ffffffff610a3316565b600160a060020a03808616600090815260208190526040808220939093559085168152205461069e908363ffffffff610a4516565b600160a060020a038085166000908152602081815260408083209490945591871681526002825282812033825290915220546106e0908363ffffffff610a3316565b600160a060020a03808616600081815260026020908152604080832033845282529182902094909455805186815290519287169391927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef929181900390910190a35060019392505050565b336000908152600260209081526040808320600160a060020a0386168452909152812054808311156107a057336000908152600260209081526040808320600160a060020a03881684529091528120556107d5565b6107b0818463ffffffff610a3316565b336000908152600260209081526040808320600160a060020a03891684529091529020555b336000818152600260209081526040808320600160a060020a0389168085529083529281902054815190815290519293927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925929181900390910190a35060019392505050565b6000600160a060020a038316151561085257600080fd5b3360009081526020819052604090205482111561086e57600080fd5b3360009081526020819052604090205461088e908363ffffffff610a3316565b3360009081526020819052604080822092909255600160a060020a038516815220546108c0908363ffffffff610a4516565b600160a060020a038416600081815260208181526040918290209390935580518581529051919233927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9281900390910190a350600192915050565b336000908152600260209081526040808320600160a060020a0386168452909152812054610950908363ffffffff610a4516565b336000818152600260209081526040808320600160a060020a0389168085529083529281902085905580519485525191937f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925929081900390910190a350600192915050565b600160a060020a03811615156109ca57600080fd5b600354604051600160a060020a038084169216907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a36003805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b600082821115610a3f57fe5b50900390565b81810182811015610a5257fe5b929150505600a165627a7a72305820b74990bbc2c31603ee33b72f704b7d064e6449244ddb06fbfd184f6df591549f0029`

// DeployPausableToken deploys a new Ethereum contract, binding an instance of PausableToken to it.
func DeployPausableToken(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *PausableToken, error) {
	parsed, err := abi.JSON(strings.NewReader(PausableTokenABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(PausableTokenBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &PausableToken{PausableTokenCaller: PausableTokenCaller{contract: contract}, PausableTokenTransactor: PausableTokenTransactor{contract: contract}, PausableTokenFilterer: PausableTokenFilterer{contract: contract}}, nil
}

// PausableToken is an auto generated Go binding around an Ethereum contract.
type PausableToken struct {
	PausableTokenCaller     // Read-only binding to the contract
	PausableTokenTransactor // Write-only binding to the contract
	PausableTokenFilterer   // Log filterer for contract events
}

// PausableTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type PausableTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PausableTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PausableTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PausableTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PausableTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PausableTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PausableTokenSession struct {
	Contract     *PausableToken    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PausableTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PausableTokenCallerSession struct {
	Contract *PausableTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// PausableTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PausableTokenTransactorSession struct {
	Contract     *PausableTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// PausableTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type PausableTokenRaw struct {
	Contract *PausableToken // Generic contract binding to access the raw methods on
}

// PausableTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PausableTokenCallerRaw struct {
	Contract *PausableTokenCaller // Generic read-only contract binding to access the raw methods on
}

// PausableTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PausableTokenTransactorRaw struct {
	Contract *PausableTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPausableToken creates a new instance of PausableToken, bound to a specific deployed contract.
func NewPausableToken(address common.Address, backend bind.ContractBackend) (*PausableToken, error) {
	contract, err := bindPausableToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PausableToken{PausableTokenCaller: PausableTokenCaller{contract: contract}, PausableTokenTransactor: PausableTokenTransactor{contract: contract}, PausableTokenFilterer: PausableTokenFilterer{contract: contract}}, nil
}

// NewPausableTokenCaller creates a new read-only instance of PausableToken, bound to a specific deployed contract.
func NewPausableTokenCaller(address common.Address, caller bind.ContractCaller) (*PausableTokenCaller, error) {
	contract, err := bindPausableToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PausableTokenCaller{contract: contract}, nil
}

// NewPausableTokenTransactor creates a new write-only instance of PausableToken, bound to a specific deployed contract.
func NewPausableTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*PausableTokenTransactor, error) {
	contract, err := bindPausableToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PausableTokenTransactor{contract: contract}, nil
}

// NewPausableTokenFilterer creates a new log filterer instance of PausableToken, bound to a specific deployed contract.
func NewPausableTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*PausableTokenFilterer, error) {
	contract, err := bindPausableToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PausableTokenFilterer{contract: contract}, nil
}

// bindPausableToken binds a generic wrapper to an already deployed contract.
func bindPausableToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PausableTokenABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PausableToken *PausableTokenRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _PausableToken.Contract.PausableTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PausableToken *PausableTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PausableToken.Contract.PausableTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PausableToken *PausableTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PausableToken.Contract.PausableTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PausableToken *PausableTokenCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _PausableToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PausableToken *PausableTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PausableToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PausableToken *PausableTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PausableToken.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(uint256)
func (_PausableToken *PausableTokenCaller) Allowance(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _PausableToken.contract.Call(opts, out, "allowance", _owner, _spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(uint256)
func (_PausableToken *PausableTokenSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _PausableToken.Contract.Allowance(&_PausableToken.CallOpts, _owner, _spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(uint256)
func (_PausableToken *PausableTokenCallerSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _PausableToken.Contract.Allowance(&_PausableToken.CallOpts, _owner, _spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(uint256)
func (_PausableToken *PausableTokenCaller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _PausableToken.contract.Call(opts, out, "balanceOf", _owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(uint256)
func (_PausableToken *PausableTokenSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _PausableToken.Contract.BalanceOf(&_PausableToken.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(uint256)
func (_PausableToken *PausableTokenCallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _PausableToken.Contract.BalanceOf(&_PausableToken.CallOpts, _owner)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_PausableToken *PausableTokenCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _PausableToken.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_PausableToken *PausableTokenSession) Owner() (common.Address, error) {
	return _PausableToken.Contract.Owner(&_PausableToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_PausableToken *PausableTokenCallerSession) Owner() (common.Address, error) {
	return _PausableToken.Contract.Owner(&_PausableToken.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() constant returns(bool)
func (_PausableToken *PausableTokenCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _PausableToken.contract.Call(opts, out, "paused")
	return *ret0, err
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() constant returns(bool)
func (_PausableToken *PausableTokenSession) Paused() (bool, error) {
	return _PausableToken.Contract.Paused(&_PausableToken.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() constant returns(bool)
func (_PausableToken *PausableTokenCallerSession) Paused() (bool, error) {
	return _PausableToken.Contract.Paused(&_PausableToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_PausableToken *PausableTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _PausableToken.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_PausableToken *PausableTokenSession) TotalSupply() (*big.Int, error) {
	return _PausableToken.Contract.TotalSupply(&_PausableToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_PausableToken *PausableTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _PausableToken.Contract.TotalSupply(&_PausableToken.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(bool)
func (_PausableToken *PausableTokenTransactor) Approve(opts *bind.TransactOpts, _spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _PausableToken.contract.Transact(opts, "approve", _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(bool)
func (_PausableToken *PausableTokenSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _PausableToken.Contract.Approve(&_PausableToken.TransactOpts, _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(bool)
func (_PausableToken *PausableTokenTransactorSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _PausableToken.Contract.Approve(&_PausableToken.TransactOpts, _spender, _value)
}

// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
//
// Solidity: function decreaseApproval(_spender address, _subtractedValue uint256) returns(success bool)
func (_PausableToken *PausableTokenTransactor) DecreaseApproval(opts *bind.TransactOpts, _spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
	return _PausableToken.contract.Transact(opts, "decreaseApproval", _spender, _subtractedValue)
}

// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
//
// Solidity: function decreaseApproval(_spender address, _subtractedValue uint256) returns(success bool)
func (_PausableToken *PausableTokenSession) DecreaseApproval(_spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
	return _PausableToken.Contract.DecreaseApproval(&_PausableToken.TransactOpts, _spender, _subtractedValue)
}

// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
//
// Solidity: function decreaseApproval(_spender address, _subtractedValue uint256) returns(success bool)
func (_PausableToken *PausableTokenTransactorSession) DecreaseApproval(_spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
	return _PausableToken.Contract.DecreaseApproval(&_PausableToken.TransactOpts, _spender, _subtractedValue)
}

// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
//
// Solidity: function increaseApproval(_spender address, _addedValue uint256) returns(success bool)
func (_PausableToken *PausableTokenTransactor) IncreaseApproval(opts *bind.TransactOpts, _spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
	return _PausableToken.contract.Transact(opts, "increaseApproval", _spender, _addedValue)
}

// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
//
// Solidity: function increaseApproval(_spender address, _addedValue uint256) returns(success bool)
func (_PausableToken *PausableTokenSession) IncreaseApproval(_spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
	return _PausableToken.Contract.IncreaseApproval(&_PausableToken.TransactOpts, _spender, _addedValue)
}

// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
//
// Solidity: function increaseApproval(_spender address, _addedValue uint256) returns(success bool)
func (_PausableToken *PausableTokenTransactorSession) IncreaseApproval(_spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
	return _PausableToken.Contract.IncreaseApproval(&_PausableToken.TransactOpts, _spender, _addedValue)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_PausableToken *PausableTokenTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PausableToken.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_PausableToken *PausableTokenSession) Pause() (*types.Transaction, error) {
	return _PausableToken.Contract.Pause(&_PausableToken.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_PausableToken *PausableTokenTransactorSession) Pause() (*types.Transaction, error) {
	return _PausableToken.Contract.Pause(&_PausableToken.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_PausableToken *PausableTokenTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PausableToken.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_PausableToken *PausableTokenSession) RenounceOwnership() (*types.Transaction, error) {
	return _PausableToken.Contract.RenounceOwnership(&_PausableToken.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_PausableToken *PausableTokenTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _PausableToken.Contract.RenounceOwnership(&_PausableToken.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_PausableToken *PausableTokenTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _PausableToken.contract.Transact(opts, "transfer", _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_PausableToken *PausableTokenSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _PausableToken.Contract.Transfer(&_PausableToken.TransactOpts, _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_PausableToken *PausableTokenTransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _PausableToken.Contract.Transfer(&_PausableToken.TransactOpts, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
func (_PausableToken *PausableTokenTransactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _PausableToken.contract.Transact(opts, "transferFrom", _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
func (_PausableToken *PausableTokenSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _PausableToken.Contract.TransferFrom(&_PausableToken.TransactOpts, _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
func (_PausableToken *PausableTokenTransactorSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _PausableToken.Contract.TransferFrom(&_PausableToken.TransactOpts, _from, _to, _value)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_PausableToken *PausableTokenTransactor) TransferOwnership(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _PausableToken.contract.Transact(opts, "transferOwnership", _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_PausableToken *PausableTokenSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _PausableToken.Contract.TransferOwnership(&_PausableToken.TransactOpts, _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_PausableToken *PausableTokenTransactorSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _PausableToken.Contract.TransferOwnership(&_PausableToken.TransactOpts, _newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_PausableToken *PausableTokenTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PausableToken.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_PausableToken *PausableTokenSession) Unpause() (*types.Transaction, error) {
	return _PausableToken.Contract.Unpause(&_PausableToken.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_PausableToken *PausableTokenTransactorSession) Unpause() (*types.Transaction, error) {
	return _PausableToken.Contract.Unpause(&_PausableToken.TransactOpts)
}

// PausableTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the PausableToken contract.
type PausableTokenApprovalIterator struct {
	Event *PausableTokenApproval // Event containing the contract specifics and raw log

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
func (it *PausableTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PausableTokenApproval)
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
		it.Event = new(PausableTokenApproval)
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
func (it *PausableTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PausableTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PausableTokenApproval represents a Approval event raised by the PausableToken contract.
type PausableTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(owner indexed address, spender indexed address, value uint256)
func (_PausableToken *PausableTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*PausableTokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _PausableToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &PausableTokenApprovalIterator{contract: _PausableToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(owner indexed address, spender indexed address, value uint256)
func (_PausableToken *PausableTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *PausableTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _PausableToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PausableTokenApproval)
				if err := _PausableToken.contract.UnpackLog(event, "Approval", log); err != nil {
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

// PausableTokenOwnershipRenouncedIterator is returned from FilterOwnershipRenounced and is used to iterate over the raw logs and unpacked data for OwnershipRenounced events raised by the PausableToken contract.
type PausableTokenOwnershipRenouncedIterator struct {
	Event *PausableTokenOwnershipRenounced // Event containing the contract specifics and raw log

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
func (it *PausableTokenOwnershipRenouncedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PausableTokenOwnershipRenounced)
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
		it.Event = new(PausableTokenOwnershipRenounced)
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
func (it *PausableTokenOwnershipRenouncedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PausableTokenOwnershipRenouncedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PausableTokenOwnershipRenounced represents a OwnershipRenounced event raised by the PausableToken contract.
type PausableTokenOwnershipRenounced struct {
	PreviousOwner common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipRenounced is a free log retrieval operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_PausableToken *PausableTokenFilterer) FilterOwnershipRenounced(opts *bind.FilterOpts, previousOwner []common.Address) (*PausableTokenOwnershipRenouncedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _PausableToken.contract.FilterLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return &PausableTokenOwnershipRenouncedIterator{contract: _PausableToken.contract, event: "OwnershipRenounced", logs: logs, sub: sub}, nil
}

// WatchOwnershipRenounced is a free log subscription operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_PausableToken *PausableTokenFilterer) WatchOwnershipRenounced(opts *bind.WatchOpts, sink chan<- *PausableTokenOwnershipRenounced, previousOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _PausableToken.contract.WatchLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PausableTokenOwnershipRenounced)
				if err := _PausableToken.contract.UnpackLog(event, "OwnershipRenounced", log); err != nil {
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

// PausableTokenOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the PausableToken contract.
type PausableTokenOwnershipTransferredIterator struct {
	Event *PausableTokenOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *PausableTokenOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PausableTokenOwnershipTransferred)
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
		it.Event = new(PausableTokenOwnershipTransferred)
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
func (it *PausableTokenOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PausableTokenOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PausableTokenOwnershipTransferred represents a OwnershipTransferred event raised by the PausableToken contract.
type PausableTokenOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_PausableToken *PausableTokenFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*PausableTokenOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _PausableToken.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &PausableTokenOwnershipTransferredIterator{contract: _PausableToken.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_PausableToken *PausableTokenFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PausableTokenOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _PausableToken.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PausableTokenOwnershipTransferred)
				if err := _PausableToken.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// PausableTokenPauseIterator is returned from FilterPause and is used to iterate over the raw logs and unpacked data for Pause events raised by the PausableToken contract.
type PausableTokenPauseIterator struct {
	Event *PausableTokenPause // Event containing the contract specifics and raw log

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
func (it *PausableTokenPauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PausableTokenPause)
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
		it.Event = new(PausableTokenPause)
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
func (it *PausableTokenPauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PausableTokenPauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PausableTokenPause represents a Pause event raised by the PausableToken contract.
type PausableTokenPause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterPause is a free log retrieval operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: e Pause()
func (_PausableToken *PausableTokenFilterer) FilterPause(opts *bind.FilterOpts) (*PausableTokenPauseIterator, error) {

	logs, sub, err := _PausableToken.contract.FilterLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return &PausableTokenPauseIterator{contract: _PausableToken.contract, event: "Pause", logs: logs, sub: sub}, nil
}

// WatchPause is a free log subscription operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: e Pause()
func (_PausableToken *PausableTokenFilterer) WatchPause(opts *bind.WatchOpts, sink chan<- *PausableTokenPause) (event.Subscription, error) {

	logs, sub, err := _PausableToken.contract.WatchLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PausableTokenPause)
				if err := _PausableToken.contract.UnpackLog(event, "Pause", log); err != nil {
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

// PausableTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the PausableToken contract.
type PausableTokenTransferIterator struct {
	Event *PausableTokenTransfer // Event containing the contract specifics and raw log

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
func (it *PausableTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PausableTokenTransfer)
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
		it.Event = new(PausableTokenTransfer)
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
func (it *PausableTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PausableTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PausableTokenTransfer represents a Transfer event raised by the PausableToken contract.
type PausableTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
func (_PausableToken *PausableTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PausableTokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PausableToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &PausableTokenTransferIterator{contract: _PausableToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
func (_PausableToken *PausableTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *PausableTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PausableToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PausableTokenTransfer)
				if err := _PausableToken.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// PausableTokenUnpauseIterator is returned from FilterUnpause and is used to iterate over the raw logs and unpacked data for Unpause events raised by the PausableToken contract.
type PausableTokenUnpauseIterator struct {
	Event *PausableTokenUnpause // Event containing the contract specifics and raw log

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
func (it *PausableTokenUnpauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PausableTokenUnpause)
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
		it.Event = new(PausableTokenUnpause)
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
func (it *PausableTokenUnpauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PausableTokenUnpauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PausableTokenUnpause represents a Unpause event raised by the PausableToken contract.
type PausableTokenUnpause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUnpause is a free log retrieval operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: e Unpause()
func (_PausableToken *PausableTokenFilterer) FilterUnpause(opts *bind.FilterOpts) (*PausableTokenUnpauseIterator, error) {

	logs, sub, err := _PausableToken.contract.FilterLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return &PausableTokenUnpauseIterator{contract: _PausableToken.contract, event: "Unpause", logs: logs, sub: sub}, nil
}

// WatchUnpause is a free log subscription operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: e Unpause()
func (_PausableToken *PausableTokenFilterer) WatchUnpause(opts *bind.WatchOpts, sink chan<- *PausableTokenUnpause) (event.Subscription, error) {

	logs, sub, err := _PausableToken.contract.WatchLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PausableTokenUnpause)
				if err := _PausableToken.contract.UnpackLog(event, "Unpause", log); err != nil {
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

// RepublicTokenABI is the input ABI used to generate the binding from.
const RepublicTokenABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"INITIAL_SUPPLY\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseApproval\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"beneficiary\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferTokens\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseApproval\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"burner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Burn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"}],\"name\":\"OwnershipRenounced\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]"

// RepublicTokenBin is the compiled bytecode used for deploying new contracts.
const RepublicTokenBin = `0x60806040526003805460a060020a60ff021916905534801561002057600080fd5b5060038054600160a060020a031916339081179091556b033b2e3c9fd0803ce8000000600181905560009182526020829052604090912055610e31806100676000396000f3006080604052600436106101115763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166306fdde038114610116578063095ea7b3146101a057806318160ddd146101d857806323b872dd146101ff5780632ff2e9dc14610229578063313ce5671461023e5780633f4ba83a1461026957806342966c68146102805780635c975abb1461029857806366188463146102ad57806370a08231146102d1578063715018a6146102f25780638456cb59146103075780638da5cb5b1461031c57806395d89b411461034d578063a9059cbb14610362578063bec3fa1714610386578063d73dd623146103aa578063dd62ed3e146103ce578063f2fde38b146103f5575b600080fd5b34801561012257600080fd5b5061012b610416565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561016557818101518382015260200161014d565b50505050905090810190601f1680156101925780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156101ac57600080fd5b506101c4600160a060020a036004351660243561044d565b604080519115158252519081900360200190f35b3480156101e457600080fd5b506101ed610478565b60408051918252519081900360200190f35b34801561020b57600080fd5b506101c4600160a060020a036004358116906024351660443561047e565b34801561023557600080fd5b506101ed6104ab565b34801561024a57600080fd5b506102536104bb565b6040805160ff9092168252519081900360200190f35b34801561027557600080fd5b5061027e6104c0565b005b34801561028c57600080fd5b5061027e600435610538565b3480156102a457600080fd5b506101c4610545565b3480156102b957600080fd5b506101c4600160a060020a0360043516602435610555565b3480156102dd57600080fd5b506101ed600160a060020a0360043516610579565b3480156102fe57600080fd5b5061027e610594565b34801561031357600080fd5b5061027e610602565b34801561032857600080fd5b5061033161067f565b60408051600160a060020a039092168252519081900360200190f35b34801561035957600080fd5b5061012b61068e565b34801561036e57600080fd5b506101c4600160a060020a03600435166024356106c5565b34801561039257600080fd5b506101c4600160a060020a03600435166024356106e9565b3480156103b657600080fd5b506101c4600160a060020a03600435166024356107c1565b3480156103da57600080fd5b506101ed600160a060020a03600435811690602435166107e5565b34801561040157600080fd5b5061027e600160a060020a0360043516610810565b60408051808201909152600e81527f52657075626c696320546f6b656e000000000000000000000000000000000000602082015281565b60035460009060a060020a900460ff161561046757600080fd5b6104718383610830565b9392505050565b60015490565b60035460009060a060020a900460ff161561049857600080fd5b6104a3848484610896565b949350505050565b6b033b2e3c9fd0803ce800000081565b601281565b600354600160a060020a031633146104d757600080fd5b60035460a060020a900460ff1615156104ef57600080fd5b6003805474ff0000000000000000000000000000000000000000191690556040517f7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b3390600090a1565b61054233826109fb565b50565b60035460a060020a900460ff1681565b60035460009060a060020a900460ff161561056f57600080fd5b6104718383610aea565b600160a060020a031660009081526020819052604090205490565b600354600160a060020a031633146105ab57600080fd5b600354604051600160a060020a03909116907ff8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c6482090600090a26003805473ffffffffffffffffffffffffffffffffffffffff19169055565b600354600160a060020a0316331461061957600080fd5b60035460a060020a900460ff161561063057600080fd5b6003805474ff0000000000000000000000000000000000000000191660a060020a1790556040517f6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff62590600090a1565b600354600160a060020a031681565b60408051808201909152600381527f52454e0000000000000000000000000000000000000000000000000000000000602082015281565b60035460009060a060020a900460ff16156106df57600080fd5b6104718383610bda565b600354600090600160a060020a0316331461070357600080fd5b6000821161071057600080fd5b600354600160a060020a031660009081526020819052604090205461073b908363ffffffff610ca916565b600354600160a060020a039081166000908152602081905260408082209390935590851681522054610773908363ffffffff610cbb16565b600160a060020a038085166000818152602081815260409182902094909455600354815187815291519294931692600080516020610de683398151915292918290030190a350600192915050565b60035460009060a060020a900460ff16156107db57600080fd5b6104718383610cce565b600160a060020a03918216600090815260026020908152604080832093909416825291909152205490565b600354600160a060020a0316331461082757600080fd5b61054281610d67565b336000818152600260209081526040808320600160a060020a038716808552908352818420869055815186815291519394909390927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925928290030190a350600192915050565b6000600160a060020a03831615156108ad57600080fd5b600160a060020a0384166000908152602081905260409020548211156108d257600080fd5b600160a060020a038416600090815260026020908152604080832033845290915290205482111561090257600080fd5b600160a060020a03841660009081526020819052604090205461092b908363ffffffff610ca916565b600160a060020a038086166000908152602081905260408082209390935590851681522054610960908363ffffffff610cbb16565b600160a060020a038085166000908152602081815260408083209490945591871681526002825282812033825290915220546109a2908363ffffffff610ca916565b600160a060020a0380861660008181526002602090815260408083203384528252918290209490945580518681529051928716939192600080516020610de6833981519152929181900390910190a35060019392505050565b600160a060020a038216600090815260208190526040902054811115610a2057600080fd5b600160a060020a038216600090815260208190526040902054610a49908263ffffffff610ca916565b600160a060020a038316600090815260208190526040902055600154610a75908263ffffffff610ca916565b600155604080518281529051600160a060020a038416917fcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5919081900360200190a2604080518281529051600091600160a060020a03851691600080516020610de68339815191529181900360200190a35050565b336000908152600260209081526040808320600160a060020a038616845290915281205480831115610b3f57336000908152600260209081526040808320600160a060020a0388168452909152812055610b74565b610b4f818463ffffffff610ca916565b336000908152600260209081526040808320600160a060020a03891684529091529020555b336000818152600260209081526040808320600160a060020a0389168085529083529281902054815190815290519293927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925929181900390910190a35060019392505050565b6000600160a060020a0383161515610bf157600080fd5b33600090815260208190526040902054821115610c0d57600080fd5b33600090815260208190526040902054610c2d908363ffffffff610ca916565b3360009081526020819052604080822092909255600160a060020a03851681522054610c5f908363ffffffff610cbb16565b600160a060020a03841660008181526020818152604091829020939093558051858152905191923392600080516020610de68339815191529281900390910190a350600192915050565b600082821115610cb557fe5b50900390565b81810182811015610cc857fe5b92915050565b336000908152600260209081526040808320600160a060020a0386168452909152812054610d02908363ffffffff610cbb16565b336000818152600260209081526040808320600160a060020a0389168085529083529281902085905580519485525191937f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925929081900390910190a350600192915050565b600160a060020a0381161515610d7c57600080fd5b600354604051600160a060020a038084169216907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a36003805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03929092169190911790555600ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa165627a7a723058205f984844e88552dfc75d11b0e9d7a4c2e8d59e8b3f97fa47e2fe0e89cd586ef10029`

// DeployRepublicToken deploys a new Ethereum contract, binding an instance of RepublicToken to it.
func DeployRepublicToken(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *RepublicToken, error) {
	parsed, err := abi.JSON(strings.NewReader(RepublicTokenABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(RepublicTokenBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &RepublicToken{RepublicTokenCaller: RepublicTokenCaller{contract: contract}, RepublicTokenTransactor: RepublicTokenTransactor{contract: contract}, RepublicTokenFilterer: RepublicTokenFilterer{contract: contract}}, nil
}

// RepublicToken is an auto generated Go binding around an Ethereum contract.
type RepublicToken struct {
	RepublicTokenCaller     // Read-only binding to the contract
	RepublicTokenTransactor // Write-only binding to the contract
	RepublicTokenFilterer   // Log filterer for contract events
}

// RepublicTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type RepublicTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RepublicTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RepublicTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RepublicTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RepublicTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RepublicTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RepublicTokenSession struct {
	Contract     *RepublicToken    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RepublicTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RepublicTokenCallerSession struct {
	Contract *RepublicTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// RepublicTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RepublicTokenTransactorSession struct {
	Contract     *RepublicTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// RepublicTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type RepublicTokenRaw struct {
	Contract *RepublicToken // Generic contract binding to access the raw methods on
}

// RepublicTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RepublicTokenCallerRaw struct {
	Contract *RepublicTokenCaller // Generic read-only contract binding to access the raw methods on
}

// RepublicTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RepublicTokenTransactorRaw struct {
	Contract *RepublicTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRepublicToken creates a new instance of RepublicToken, bound to a specific deployed contract.
func NewRepublicToken(address common.Address, backend bind.ContractBackend) (*RepublicToken, error) {
	contract, err := bindRepublicToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RepublicToken{RepublicTokenCaller: RepublicTokenCaller{contract: contract}, RepublicTokenTransactor: RepublicTokenTransactor{contract: contract}, RepublicTokenFilterer: RepublicTokenFilterer{contract: contract}}, nil
}

// NewRepublicTokenCaller creates a new read-only instance of RepublicToken, bound to a specific deployed contract.
func NewRepublicTokenCaller(address common.Address, caller bind.ContractCaller) (*RepublicTokenCaller, error) {
	contract, err := bindRepublicToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RepublicTokenCaller{contract: contract}, nil
}

// NewRepublicTokenTransactor creates a new write-only instance of RepublicToken, bound to a specific deployed contract.
func NewRepublicTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*RepublicTokenTransactor, error) {
	contract, err := bindRepublicToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RepublicTokenTransactor{contract: contract}, nil
}

// NewRepublicTokenFilterer creates a new log filterer instance of RepublicToken, bound to a specific deployed contract.
func NewRepublicTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*RepublicTokenFilterer, error) {
	contract, err := bindRepublicToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RepublicTokenFilterer{contract: contract}, nil
}

// bindRepublicToken binds a generic wrapper to an already deployed contract.
func bindRepublicToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RepublicTokenABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RepublicToken *RepublicTokenRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _RepublicToken.Contract.RepublicTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RepublicToken *RepublicTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RepublicToken.Contract.RepublicTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RepublicToken *RepublicTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RepublicToken.Contract.RepublicTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RepublicToken *RepublicTokenCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _RepublicToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RepublicToken *RepublicTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RepublicToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RepublicToken *RepublicTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RepublicToken.Contract.contract.Transact(opts, method, params...)
}

// INITIALSUPPLY is a free data retrieval call binding the contract method 0x2ff2e9dc.
//
// Solidity: function INITIAL_SUPPLY() constant returns(uint256)
func (_RepublicToken *RepublicTokenCaller) INITIALSUPPLY(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RepublicToken.contract.Call(opts, out, "INITIAL_SUPPLY")
	return *ret0, err
}

// INITIALSUPPLY is a free data retrieval call binding the contract method 0x2ff2e9dc.
//
// Solidity: function INITIAL_SUPPLY() constant returns(uint256)
func (_RepublicToken *RepublicTokenSession) INITIALSUPPLY() (*big.Int, error) {
	return _RepublicToken.Contract.INITIALSUPPLY(&_RepublicToken.CallOpts)
}

// INITIALSUPPLY is a free data retrieval call binding the contract method 0x2ff2e9dc.
//
// Solidity: function INITIAL_SUPPLY() constant returns(uint256)
func (_RepublicToken *RepublicTokenCallerSession) INITIALSUPPLY() (*big.Int, error) {
	return _RepublicToken.Contract.INITIALSUPPLY(&_RepublicToken.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(uint256)
func (_RepublicToken *RepublicTokenCaller) Allowance(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RepublicToken.contract.Call(opts, out, "allowance", _owner, _spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(uint256)
func (_RepublicToken *RepublicTokenSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _RepublicToken.Contract.Allowance(&_RepublicToken.CallOpts, _owner, _spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(uint256)
func (_RepublicToken *RepublicTokenCallerSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _RepublicToken.Contract.Allowance(&_RepublicToken.CallOpts, _owner, _spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(uint256)
func (_RepublicToken *RepublicTokenCaller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RepublicToken.contract.Call(opts, out, "balanceOf", _owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(uint256)
func (_RepublicToken *RepublicTokenSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _RepublicToken.Contract.BalanceOf(&_RepublicToken.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(uint256)
func (_RepublicToken *RepublicTokenCallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _RepublicToken.Contract.BalanceOf(&_RepublicToken.CallOpts, _owner)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_RepublicToken *RepublicTokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _RepublicToken.contract.Call(opts, out, "decimals")
	return *ret0, err
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_RepublicToken *RepublicTokenSession) Decimals() (uint8, error) {
	return _RepublicToken.Contract.Decimals(&_RepublicToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_RepublicToken *RepublicTokenCallerSession) Decimals() (uint8, error) {
	return _RepublicToken.Contract.Decimals(&_RepublicToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_RepublicToken *RepublicTokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _RepublicToken.contract.Call(opts, out, "name")
	return *ret0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_RepublicToken *RepublicTokenSession) Name() (string, error) {
	return _RepublicToken.Contract.Name(&_RepublicToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_RepublicToken *RepublicTokenCallerSession) Name() (string, error) {
	return _RepublicToken.Contract.Name(&_RepublicToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_RepublicToken *RepublicTokenCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _RepublicToken.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_RepublicToken *RepublicTokenSession) Owner() (common.Address, error) {
	return _RepublicToken.Contract.Owner(&_RepublicToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_RepublicToken *RepublicTokenCallerSession) Owner() (common.Address, error) {
	return _RepublicToken.Contract.Owner(&_RepublicToken.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() constant returns(bool)
func (_RepublicToken *RepublicTokenCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _RepublicToken.contract.Call(opts, out, "paused")
	return *ret0, err
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() constant returns(bool)
func (_RepublicToken *RepublicTokenSession) Paused() (bool, error) {
	return _RepublicToken.Contract.Paused(&_RepublicToken.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() constant returns(bool)
func (_RepublicToken *RepublicTokenCallerSession) Paused() (bool, error) {
	return _RepublicToken.Contract.Paused(&_RepublicToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_RepublicToken *RepublicTokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _RepublicToken.contract.Call(opts, out, "symbol")
	return *ret0, err
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_RepublicToken *RepublicTokenSession) Symbol() (string, error) {
	return _RepublicToken.Contract.Symbol(&_RepublicToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_RepublicToken *RepublicTokenCallerSession) Symbol() (string, error) {
	return _RepublicToken.Contract.Symbol(&_RepublicToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_RepublicToken *RepublicTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RepublicToken.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_RepublicToken *RepublicTokenSession) TotalSupply() (*big.Int, error) {
	return _RepublicToken.Contract.TotalSupply(&_RepublicToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_RepublicToken *RepublicTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _RepublicToken.Contract.TotalSupply(&_RepublicToken.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(bool)
func (_RepublicToken *RepublicTokenTransactor) Approve(opts *bind.TransactOpts, _spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _RepublicToken.contract.Transact(opts, "approve", _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(bool)
func (_RepublicToken *RepublicTokenSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _RepublicToken.Contract.Approve(&_RepublicToken.TransactOpts, _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(bool)
func (_RepublicToken *RepublicTokenTransactorSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _RepublicToken.Contract.Approve(&_RepublicToken.TransactOpts, _spender, _value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(_value uint256) returns()
func (_RepublicToken *RepublicTokenTransactor) Burn(opts *bind.TransactOpts, _value *big.Int) (*types.Transaction, error) {
	return _RepublicToken.contract.Transact(opts, "burn", _value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(_value uint256) returns()
func (_RepublicToken *RepublicTokenSession) Burn(_value *big.Int) (*types.Transaction, error) {
	return _RepublicToken.Contract.Burn(&_RepublicToken.TransactOpts, _value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(_value uint256) returns()
func (_RepublicToken *RepublicTokenTransactorSession) Burn(_value *big.Int) (*types.Transaction, error) {
	return _RepublicToken.Contract.Burn(&_RepublicToken.TransactOpts, _value)
}

// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
//
// Solidity: function decreaseApproval(_spender address, _subtractedValue uint256) returns(success bool)
func (_RepublicToken *RepublicTokenTransactor) DecreaseApproval(opts *bind.TransactOpts, _spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
	return _RepublicToken.contract.Transact(opts, "decreaseApproval", _spender, _subtractedValue)
}

// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
//
// Solidity: function decreaseApproval(_spender address, _subtractedValue uint256) returns(success bool)
func (_RepublicToken *RepublicTokenSession) DecreaseApproval(_spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
	return _RepublicToken.Contract.DecreaseApproval(&_RepublicToken.TransactOpts, _spender, _subtractedValue)
}

// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
//
// Solidity: function decreaseApproval(_spender address, _subtractedValue uint256) returns(success bool)
func (_RepublicToken *RepublicTokenTransactorSession) DecreaseApproval(_spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
	return _RepublicToken.Contract.DecreaseApproval(&_RepublicToken.TransactOpts, _spender, _subtractedValue)
}

// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
//
// Solidity: function increaseApproval(_spender address, _addedValue uint256) returns(success bool)
func (_RepublicToken *RepublicTokenTransactor) IncreaseApproval(opts *bind.TransactOpts, _spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
	return _RepublicToken.contract.Transact(opts, "increaseApproval", _spender, _addedValue)
}

// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
//
// Solidity: function increaseApproval(_spender address, _addedValue uint256) returns(success bool)
func (_RepublicToken *RepublicTokenSession) IncreaseApproval(_spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
	return _RepublicToken.Contract.IncreaseApproval(&_RepublicToken.TransactOpts, _spender, _addedValue)
}

// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
//
// Solidity: function increaseApproval(_spender address, _addedValue uint256) returns(success bool)
func (_RepublicToken *RepublicTokenTransactorSession) IncreaseApproval(_spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
	return _RepublicToken.Contract.IncreaseApproval(&_RepublicToken.TransactOpts, _spender, _addedValue)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_RepublicToken *RepublicTokenTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RepublicToken.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_RepublicToken *RepublicTokenSession) Pause() (*types.Transaction, error) {
	return _RepublicToken.Contract.Pause(&_RepublicToken.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_RepublicToken *RepublicTokenTransactorSession) Pause() (*types.Transaction, error) {
	return _RepublicToken.Contract.Pause(&_RepublicToken.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_RepublicToken *RepublicTokenTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RepublicToken.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_RepublicToken *RepublicTokenSession) RenounceOwnership() (*types.Transaction, error) {
	return _RepublicToken.Contract.RenounceOwnership(&_RepublicToken.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_RepublicToken *RepublicTokenTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _RepublicToken.Contract.RenounceOwnership(&_RepublicToken.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_RepublicToken *RepublicTokenTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _RepublicToken.contract.Transact(opts, "transfer", _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_RepublicToken *RepublicTokenSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _RepublicToken.Contract.Transfer(&_RepublicToken.TransactOpts, _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_RepublicToken *RepublicTokenTransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _RepublicToken.Contract.Transfer(&_RepublicToken.TransactOpts, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
func (_RepublicToken *RepublicTokenTransactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _RepublicToken.contract.Transact(opts, "transferFrom", _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
func (_RepublicToken *RepublicTokenSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _RepublicToken.Contract.TransferFrom(&_RepublicToken.TransactOpts, _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
func (_RepublicToken *RepublicTokenTransactorSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _RepublicToken.Contract.TransferFrom(&_RepublicToken.TransactOpts, _from, _to, _value)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_RepublicToken *RepublicTokenTransactor) TransferOwnership(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _RepublicToken.contract.Transact(opts, "transferOwnership", _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_RepublicToken *RepublicTokenSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _RepublicToken.Contract.TransferOwnership(&_RepublicToken.TransactOpts, _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_RepublicToken *RepublicTokenTransactorSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _RepublicToken.Contract.TransferOwnership(&_RepublicToken.TransactOpts, _newOwner)
}

// TransferTokens is a paid mutator transaction binding the contract method 0xbec3fa17.
//
// Solidity: function transferTokens(beneficiary address, amount uint256) returns(bool)
func (_RepublicToken *RepublicTokenTransactor) TransferTokens(opts *bind.TransactOpts, beneficiary common.Address, amount *big.Int) (*types.Transaction, error) {
	return _RepublicToken.contract.Transact(opts, "transferTokens", beneficiary, amount)
}

// TransferTokens is a paid mutator transaction binding the contract method 0xbec3fa17.
//
// Solidity: function transferTokens(beneficiary address, amount uint256) returns(bool)
func (_RepublicToken *RepublicTokenSession) TransferTokens(beneficiary common.Address, amount *big.Int) (*types.Transaction, error) {
	return _RepublicToken.Contract.TransferTokens(&_RepublicToken.TransactOpts, beneficiary, amount)
}

// TransferTokens is a paid mutator transaction binding the contract method 0xbec3fa17.
//
// Solidity: function transferTokens(beneficiary address, amount uint256) returns(bool)
func (_RepublicToken *RepublicTokenTransactorSession) TransferTokens(beneficiary common.Address, amount *big.Int) (*types.Transaction, error) {
	return _RepublicToken.Contract.TransferTokens(&_RepublicToken.TransactOpts, beneficiary, amount)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_RepublicToken *RepublicTokenTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RepublicToken.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_RepublicToken *RepublicTokenSession) Unpause() (*types.Transaction, error) {
	return _RepublicToken.Contract.Unpause(&_RepublicToken.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_RepublicToken *RepublicTokenTransactorSession) Unpause() (*types.Transaction, error) {
	return _RepublicToken.Contract.Unpause(&_RepublicToken.TransactOpts)
}

// RepublicTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the RepublicToken contract.
type RepublicTokenApprovalIterator struct {
	Event *RepublicTokenApproval // Event containing the contract specifics and raw log

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
func (it *RepublicTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RepublicTokenApproval)
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
		it.Event = new(RepublicTokenApproval)
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
func (it *RepublicTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RepublicTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RepublicTokenApproval represents a Approval event raised by the RepublicToken contract.
type RepublicTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(owner indexed address, spender indexed address, value uint256)
func (_RepublicToken *RepublicTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*RepublicTokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _RepublicToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &RepublicTokenApprovalIterator{contract: _RepublicToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(owner indexed address, spender indexed address, value uint256)
func (_RepublicToken *RepublicTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *RepublicTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _RepublicToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RepublicTokenApproval)
				if err := _RepublicToken.contract.UnpackLog(event, "Approval", log); err != nil {
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

// RepublicTokenBurnIterator is returned from FilterBurn and is used to iterate over the raw logs and unpacked data for Burn events raised by the RepublicToken contract.
type RepublicTokenBurnIterator struct {
	Event *RepublicTokenBurn // Event containing the contract specifics and raw log

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
func (it *RepublicTokenBurnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RepublicTokenBurn)
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
		it.Event = new(RepublicTokenBurn)
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
func (it *RepublicTokenBurnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RepublicTokenBurnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RepublicTokenBurn represents a Burn event raised by the RepublicToken contract.
type RepublicTokenBurn struct {
	Burner common.Address
	Value  *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBurn is a free log retrieval operation binding the contract event 0xcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5.
//
// Solidity: e Burn(burner indexed address, value uint256)
func (_RepublicToken *RepublicTokenFilterer) FilterBurn(opts *bind.FilterOpts, burner []common.Address) (*RepublicTokenBurnIterator, error) {

	var burnerRule []interface{}
	for _, burnerItem := range burner {
		burnerRule = append(burnerRule, burnerItem)
	}

	logs, sub, err := _RepublicToken.contract.FilterLogs(opts, "Burn", burnerRule)
	if err != nil {
		return nil, err
	}
	return &RepublicTokenBurnIterator{contract: _RepublicToken.contract, event: "Burn", logs: logs, sub: sub}, nil
}

// WatchBurn is a free log subscription operation binding the contract event 0xcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5.
//
// Solidity: e Burn(burner indexed address, value uint256)
func (_RepublicToken *RepublicTokenFilterer) WatchBurn(opts *bind.WatchOpts, sink chan<- *RepublicTokenBurn, burner []common.Address) (event.Subscription, error) {

	var burnerRule []interface{}
	for _, burnerItem := range burner {
		burnerRule = append(burnerRule, burnerItem)
	}

	logs, sub, err := _RepublicToken.contract.WatchLogs(opts, "Burn", burnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RepublicTokenBurn)
				if err := _RepublicToken.contract.UnpackLog(event, "Burn", log); err != nil {
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

// RepublicTokenOwnershipRenouncedIterator is returned from FilterOwnershipRenounced and is used to iterate over the raw logs and unpacked data for OwnershipRenounced events raised by the RepublicToken contract.
type RepublicTokenOwnershipRenouncedIterator struct {
	Event *RepublicTokenOwnershipRenounced // Event containing the contract specifics and raw log

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
func (it *RepublicTokenOwnershipRenouncedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RepublicTokenOwnershipRenounced)
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
		it.Event = new(RepublicTokenOwnershipRenounced)
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
func (it *RepublicTokenOwnershipRenouncedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RepublicTokenOwnershipRenouncedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RepublicTokenOwnershipRenounced represents a OwnershipRenounced event raised by the RepublicToken contract.
type RepublicTokenOwnershipRenounced struct {
	PreviousOwner common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipRenounced is a free log retrieval operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_RepublicToken *RepublicTokenFilterer) FilterOwnershipRenounced(opts *bind.FilterOpts, previousOwner []common.Address) (*RepublicTokenOwnershipRenouncedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _RepublicToken.contract.FilterLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return &RepublicTokenOwnershipRenouncedIterator{contract: _RepublicToken.contract, event: "OwnershipRenounced", logs: logs, sub: sub}, nil
}

// WatchOwnershipRenounced is a free log subscription operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_RepublicToken *RepublicTokenFilterer) WatchOwnershipRenounced(opts *bind.WatchOpts, sink chan<- *RepublicTokenOwnershipRenounced, previousOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _RepublicToken.contract.WatchLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RepublicTokenOwnershipRenounced)
				if err := _RepublicToken.contract.UnpackLog(event, "OwnershipRenounced", log); err != nil {
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

// RepublicTokenOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the RepublicToken contract.
type RepublicTokenOwnershipTransferredIterator struct {
	Event *RepublicTokenOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *RepublicTokenOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RepublicTokenOwnershipTransferred)
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
		it.Event = new(RepublicTokenOwnershipTransferred)
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
func (it *RepublicTokenOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RepublicTokenOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RepublicTokenOwnershipTransferred represents a OwnershipTransferred event raised by the RepublicToken contract.
type RepublicTokenOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_RepublicToken *RepublicTokenFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*RepublicTokenOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _RepublicToken.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &RepublicTokenOwnershipTransferredIterator{contract: _RepublicToken.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_RepublicToken *RepublicTokenFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RepublicTokenOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _RepublicToken.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RepublicTokenOwnershipTransferred)
				if err := _RepublicToken.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// RepublicTokenPauseIterator is returned from FilterPause and is used to iterate over the raw logs and unpacked data for Pause events raised by the RepublicToken contract.
type RepublicTokenPauseIterator struct {
	Event *RepublicTokenPause // Event containing the contract specifics and raw log

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
func (it *RepublicTokenPauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RepublicTokenPause)
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
		it.Event = new(RepublicTokenPause)
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
func (it *RepublicTokenPauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RepublicTokenPauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RepublicTokenPause represents a Pause event raised by the RepublicToken contract.
type RepublicTokenPause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterPause is a free log retrieval operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: e Pause()
func (_RepublicToken *RepublicTokenFilterer) FilterPause(opts *bind.FilterOpts) (*RepublicTokenPauseIterator, error) {

	logs, sub, err := _RepublicToken.contract.FilterLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return &RepublicTokenPauseIterator{contract: _RepublicToken.contract, event: "Pause", logs: logs, sub: sub}, nil
}

// WatchPause is a free log subscription operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: e Pause()
func (_RepublicToken *RepublicTokenFilterer) WatchPause(opts *bind.WatchOpts, sink chan<- *RepublicTokenPause) (event.Subscription, error) {

	logs, sub, err := _RepublicToken.contract.WatchLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RepublicTokenPause)
				if err := _RepublicToken.contract.UnpackLog(event, "Pause", log); err != nil {
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

// RepublicTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the RepublicToken contract.
type RepublicTokenTransferIterator struct {
	Event *RepublicTokenTransfer // Event containing the contract specifics and raw log

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
func (it *RepublicTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RepublicTokenTransfer)
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
		it.Event = new(RepublicTokenTransfer)
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
func (it *RepublicTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RepublicTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RepublicTokenTransfer represents a Transfer event raised by the RepublicToken contract.
type RepublicTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
func (_RepublicToken *RepublicTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RepublicTokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RepublicToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &RepublicTokenTransferIterator{contract: _RepublicToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
func (_RepublicToken *RepublicTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *RepublicTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RepublicToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RepublicTokenTransfer)
				if err := _RepublicToken.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// RepublicTokenUnpauseIterator is returned from FilterUnpause and is used to iterate over the raw logs and unpacked data for Unpause events raised by the RepublicToken contract.
type RepublicTokenUnpauseIterator struct {
	Event *RepublicTokenUnpause // Event containing the contract specifics and raw log

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
func (it *RepublicTokenUnpauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RepublicTokenUnpause)
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
		it.Event = new(RepublicTokenUnpause)
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
func (it *RepublicTokenUnpauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RepublicTokenUnpauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RepublicTokenUnpause represents a Unpause event raised by the RepublicToken contract.
type RepublicTokenUnpause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUnpause is a free log retrieval operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: e Unpause()
func (_RepublicToken *RepublicTokenFilterer) FilterUnpause(opts *bind.FilterOpts) (*RepublicTokenUnpauseIterator, error) {

	logs, sub, err := _RepublicToken.contract.FilterLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return &RepublicTokenUnpauseIterator{contract: _RepublicToken.contract, event: "Unpause", logs: logs, sub: sub}, nil
}

// WatchUnpause is a free log subscription operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: e Unpause()
func (_RepublicToken *RepublicTokenFilterer) WatchUnpause(opts *bind.WatchOpts, sink chan<- *RepublicTokenUnpause) (event.Subscription, error) {

	logs, sub, err := _RepublicToken.contract.WatchLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RepublicTokenUnpause)
				if err := _RepublicToken.contract.UnpackLog(event, "Unpause", log); err != nil {
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

// SafeMathABI is the input ABI used to generate the binding from.
const SafeMathABI = "[]"

// SafeMathBin is the compiled bytecode used for deploying new contracts.
const SafeMathBin = `0x604c602c600b82828239805160001a60731460008114601c57601e565bfe5b5030600052607381538281f30073000000000000000000000000000000000000000030146080604052600080fd00a165627a7a72305820a02a84ac4f3af3b8d9db41a87667de0285507cc2e99581d2d45fbe62a74f203d0029`

// DeploySafeMath deploys a new Ethereum contract, binding an instance of SafeMath to it.
func DeploySafeMath(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SafeMath, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMathABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SafeMathBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// SafeMath is an auto generated Go binding around an Ethereum contract.
type SafeMath struct {
	SafeMathCaller     // Read-only binding to the contract
	SafeMathTransactor // Write-only binding to the contract
	SafeMathFilterer   // Log filterer for contract events
}

// SafeMathCaller is an auto generated read-only Go binding around an Ethereum contract.
type SafeMathCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SafeMathTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SafeMathFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SafeMathSession struct {
	Contract     *SafeMath         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SafeMathCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SafeMathCallerSession struct {
	Contract *SafeMathCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// SafeMathTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SafeMathTransactorSession struct {
	Contract     *SafeMathTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// SafeMathRaw is an auto generated low-level Go binding around an Ethereum contract.
type SafeMathRaw struct {
	Contract *SafeMath // Generic contract binding to access the raw methods on
}

// SafeMathCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SafeMathCallerRaw struct {
	Contract *SafeMathCaller // Generic read-only contract binding to access the raw methods on
}

// SafeMathTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SafeMathTransactorRaw struct {
	Contract *SafeMathTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSafeMath creates a new instance of SafeMath, bound to a specific deployed contract.
func NewSafeMath(address common.Address, backend bind.ContractBackend) (*SafeMath, error) {
	contract, err := bindSafeMath(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// NewSafeMathCaller creates a new read-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathCaller(address common.Address, caller bind.ContractCaller) (*SafeMathCaller, error) {
	contract, err := bindSafeMath(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathCaller{contract: contract}, nil
}

// NewSafeMathTransactor creates a new write-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathTransactor(address common.Address, transactor bind.ContractTransactor) (*SafeMathTransactor, error) {
	contract, err := bindSafeMath(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathTransactor{contract: contract}, nil
}

// NewSafeMathFilterer creates a new log filterer instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathFilterer(address common.Address, filterer bind.ContractFilterer) (*SafeMathFilterer, error) {
	contract, err := bindSafeMath(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SafeMathFilterer{contract: contract}, nil
}

// bindSafeMath binds a generic wrapper to an already deployed contract.
func bindSafeMath(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMathABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SafeMath.Contract.SafeMathCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SafeMath.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transact(opts, method, params...)
}

// SettlementABI is the input ABI used to generate the binding from.
const SettlementABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"submitOrderGasPriceLimit\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_details\",\"type\":\"bytes\"},{\"name\":\"_settlementID\",\"type\":\"uint64\"},{\"name\":\"_tokens\",\"type\":\"uint64\"},{\"name\":\"_price\",\"type\":\"uint256\"},{\"name\":\"_volume\",\"type\":\"uint256\"},{\"name\":\"_minimumVolume\",\"type\":\"uint256\"}],\"name\":\"submitOrder\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_buyID\",\"type\":\"bytes32\"},{\"name\":\"_sellID\",\"type\":\"bytes32\"}],\"name\":\"settle\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// SettlementBin is the compiled bytecode used for deploying new contracts.
const SettlementBin = `0x`

// DeploySettlement deploys a new Ethereum contract, binding an instance of Settlement to it.
func DeploySettlement(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Settlement, error) {
	parsed, err := abi.JSON(strings.NewReader(SettlementABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SettlementBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Settlement{SettlementCaller: SettlementCaller{contract: contract}, SettlementTransactor: SettlementTransactor{contract: contract}, SettlementFilterer: SettlementFilterer{contract: contract}}, nil
}

// Settlement is an auto generated Go binding around an Ethereum contract.
type Settlement struct {
	SettlementCaller     // Read-only binding to the contract
	SettlementTransactor // Write-only binding to the contract
	SettlementFilterer   // Log filterer for contract events
}

// SettlementCaller is an auto generated read-only Go binding around an Ethereum contract.
type SettlementCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SettlementTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SettlementTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SettlementFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SettlementFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SettlementSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SettlementSession struct {
	Contract     *Settlement       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SettlementCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SettlementCallerSession struct {
	Contract *SettlementCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// SettlementTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SettlementTransactorSession struct {
	Contract     *SettlementTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// SettlementRaw is an auto generated low-level Go binding around an Ethereum contract.
type SettlementRaw struct {
	Contract *Settlement // Generic contract binding to access the raw methods on
}

// SettlementCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SettlementCallerRaw struct {
	Contract *SettlementCaller // Generic read-only contract binding to access the raw methods on
}

// SettlementTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SettlementTransactorRaw struct {
	Contract *SettlementTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSettlement creates a new instance of Settlement, bound to a specific deployed contract.
func NewSettlement(address common.Address, backend bind.ContractBackend) (*Settlement, error) {
	contract, err := bindSettlement(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Settlement{SettlementCaller: SettlementCaller{contract: contract}, SettlementTransactor: SettlementTransactor{contract: contract}, SettlementFilterer: SettlementFilterer{contract: contract}}, nil
}

// NewSettlementCaller creates a new read-only instance of Settlement, bound to a specific deployed contract.
func NewSettlementCaller(address common.Address, caller bind.ContractCaller) (*SettlementCaller, error) {
	contract, err := bindSettlement(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SettlementCaller{contract: contract}, nil
}

// NewSettlementTransactor creates a new write-only instance of Settlement, bound to a specific deployed contract.
func NewSettlementTransactor(address common.Address, transactor bind.ContractTransactor) (*SettlementTransactor, error) {
	contract, err := bindSettlement(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SettlementTransactor{contract: contract}, nil
}

// NewSettlementFilterer creates a new log filterer instance of Settlement, bound to a specific deployed contract.
func NewSettlementFilterer(address common.Address, filterer bind.ContractFilterer) (*SettlementFilterer, error) {
	contract, err := bindSettlement(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SettlementFilterer{contract: contract}, nil
}

// bindSettlement binds a generic wrapper to an already deployed contract.
func bindSettlement(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SettlementABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Settlement *SettlementRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Settlement.Contract.SettlementCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Settlement *SettlementRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Settlement.Contract.SettlementTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Settlement *SettlementRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Settlement.Contract.SettlementTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Settlement *SettlementCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Settlement.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Settlement *SettlementTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Settlement.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Settlement *SettlementTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Settlement.Contract.contract.Transact(opts, method, params...)
}

// SubmitOrderGasPriceLimit is a free data retrieval call binding the contract method 0x60556aa6.
//
// Solidity: function submitOrderGasPriceLimit() constant returns(uint256)
func (_Settlement *SettlementCaller) SubmitOrderGasPriceLimit(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Settlement.contract.Call(opts, out, "submitOrderGasPriceLimit")
	return *ret0, err
}

// SubmitOrderGasPriceLimit is a free data retrieval call binding the contract method 0x60556aa6.
//
// Solidity: function submitOrderGasPriceLimit() constant returns(uint256)
func (_Settlement *SettlementSession) SubmitOrderGasPriceLimit() (*big.Int, error) {
	return _Settlement.Contract.SubmitOrderGasPriceLimit(&_Settlement.CallOpts)
}

// SubmitOrderGasPriceLimit is a free data retrieval call binding the contract method 0x60556aa6.
//
// Solidity: function submitOrderGasPriceLimit() constant returns(uint256)
func (_Settlement *SettlementCallerSession) SubmitOrderGasPriceLimit() (*big.Int, error) {
	return _Settlement.Contract.SubmitOrderGasPriceLimit(&_Settlement.CallOpts)
}

// Settle is a paid mutator transaction binding the contract method 0xd32a9cd9.
//
// Solidity: function settle(_buyID bytes32, _sellID bytes32) returns()
func (_Settlement *SettlementTransactor) Settle(opts *bind.TransactOpts, _buyID [32]byte, _sellID [32]byte) (*types.Transaction, error) {
	return _Settlement.contract.Transact(opts, "settle", _buyID, _sellID)
}

// Settle is a paid mutator transaction binding the contract method 0xd32a9cd9.
//
// Solidity: function settle(_buyID bytes32, _sellID bytes32) returns()
func (_Settlement *SettlementSession) Settle(_buyID [32]byte, _sellID [32]byte) (*types.Transaction, error) {
	return _Settlement.Contract.Settle(&_Settlement.TransactOpts, _buyID, _sellID)
}

// Settle is a paid mutator transaction binding the contract method 0xd32a9cd9.
//
// Solidity: function settle(_buyID bytes32, _sellID bytes32) returns()
func (_Settlement *SettlementTransactorSession) Settle(_buyID [32]byte, _sellID [32]byte) (*types.Transaction, error) {
	return _Settlement.Contract.Settle(&_Settlement.TransactOpts, _buyID, _sellID)
}

// SubmitOrder is a paid mutator transaction binding the contract method 0xb86f602c.
//
// Solidity: function submitOrder(_details bytes, _settlementID uint64, _tokens uint64, _price uint256, _volume uint256, _minimumVolume uint256) returns()
func (_Settlement *SettlementTransactor) SubmitOrder(opts *bind.TransactOpts, _details []byte, _settlementID uint64, _tokens uint64, _price *big.Int, _volume *big.Int, _minimumVolume *big.Int) (*types.Transaction, error) {
	return _Settlement.contract.Transact(opts, "submitOrder", _details, _settlementID, _tokens, _price, _volume, _minimumVolume)
}

// SubmitOrder is a paid mutator transaction binding the contract method 0xb86f602c.
//
// Solidity: function submitOrder(_details bytes, _settlementID uint64, _tokens uint64, _price uint256, _volume uint256, _minimumVolume uint256) returns()
func (_Settlement *SettlementSession) SubmitOrder(_details []byte, _settlementID uint64, _tokens uint64, _price *big.Int, _volume *big.Int, _minimumVolume *big.Int) (*types.Transaction, error) {
	return _Settlement.Contract.SubmitOrder(&_Settlement.TransactOpts, _details, _settlementID, _tokens, _price, _volume, _minimumVolume)
}

// SubmitOrder is a paid mutator transaction binding the contract method 0xb86f602c.
//
// Solidity: function submitOrder(_details bytes, _settlementID uint64, _tokens uint64, _price uint256, _volume uint256, _minimumVolume uint256) returns()
func (_Settlement *SettlementTransactorSession) SubmitOrder(_details []byte, _settlementID uint64, _tokens uint64, _price *big.Int, _volume *big.Int, _minimumVolume *big.Int) (*types.Transaction, error) {
	return _Settlement.Contract.SubmitOrder(&_Settlement.TransactOpts, _details, _settlementID, _tokens, _price, _volume, _minimumVolume)
}

// SettlementRegistryABI is the input ABI used to generate the binding from.
const SettlementRegistryABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_settlementID\",\"type\":\"uint64\"}],\"name\":\"deregisterSettlement\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_settlementID\",\"type\":\"uint64\"}],\"name\":\"settlementContract\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_settlementID\",\"type\":\"uint64\"}],\"name\":\"brokerVerifierContract\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_settlementID\",\"type\":\"uint64\"},{\"name\":\"_settlementContract\",\"type\":\"address\"},{\"name\":\"_brokerVerifierContract\",\"type\":\"address\"}],\"name\":\"registerSettlement\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"settlementDetails\",\"outputs\":[{\"name\":\"registered\",\"type\":\"bool\"},{\"name\":\"settlementContract\",\"type\":\"address\"},{\"name\":\"brokerVerifierContract\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_settlementID\",\"type\":\"uint64\"}],\"name\":\"settlementRegistration\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_VERSION\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"settlementID\",\"type\":\"uint64\"},{\"indexed\":false,\"name\":\"settlementContract\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"brokerVerifierContract\",\"type\":\"address\"}],\"name\":\"LogSettlementRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"settlementID\",\"type\":\"uint64\"},{\"indexed\":false,\"name\":\"settlementContract\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"brokerVerifierContract\",\"type\":\"address\"}],\"name\":\"LogSettlementUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"settlementID\",\"type\":\"uint64\"}],\"name\":\"LogSettlementDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"}],\"name\":\"OwnershipRenounced\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"

// SettlementRegistryBin is the compiled bytecode used for deploying new contracts.
const SettlementRegistryBin = `0x608060405234801561001057600080fd5b506040516108d13803806108d183398101604052805160008054600160a060020a0319163317905501805161004c906001906020840190610053565b50506100ee565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061009457805160ff19168380011785556100c1565b828001600101855582156100c1579182015b828111156100c15782518255916020019190600101906100a6565b506100cd9291506100d1565b5090565b6100eb91905b808211156100cd57600081556001016100d7565b90565b6107d4806100fd6000396000f3006080604052600436106100a35763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166306b82a7d81146100a85780633809b240146100cc578063715018a61461010a5780638da5cb5b1461011f578063a6065c9614610134578063c2190c9414610156578063f2fde38b1461018a578063f8a09cd0146101ab578063f8f9be36146101f7578063ffa1ad741461022d575b600080fd5b3480156100b457600080fd5b506100ca67ffffffffffffffff600435166102b7565b005b3480156100d857600080fd5b506100ee67ffffffffffffffff600435166103e3565b60408051600160a060020a039092168252519081900360200190f35b34801561011657600080fd5b506100ca61040d565b34801561012b57600080fd5b506100ee610479565b34801561014057600080fd5b506100ee67ffffffffffffffff60043516610488565b34801561016257600080fd5b506100ca67ffffffffffffffff60043516600160a060020a03602435811690604435166104b0565b34801561019657600080fd5b506100ca600160a060020a036004351661062c565b3480156101b757600080fd5b506101cd67ffffffffffffffff6004351661064f565b604080519315158452600160a060020a039283166020850152911682820152519081900360600190f35b34801561020357600080fd5b5061021967ffffffffffffffff6004351661067f565b604080519115158252519081900360200190f35b34801561023957600080fd5b5061024261069e565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561027c578181015183820152602001610264565b50505050905090810190601f1680156102a95780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b600054600160a060020a031633146102ce57600080fd5b67ffffffffffffffff811660009081526002602052604090205460ff16151561035857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600e60248201527f6e6f742072656769737465726564000000000000000000000000000000000000604482015290519081900360640190fd5b67ffffffffffffffff8116600081815260026020908152604091829020805474ffffffffffffffffffffffffffffffffffffffffff19168155600101805473ffffffffffffffffffffffffffffffffffffffff19169055815192835290517f9e153628c2e49d8816ef9c792fac73678b5846d086e651aa7a6ea95f1b7f66759281900390910190a150565b67ffffffffffffffff166000908152600260205260409020546101009004600160a060020a031690565b600054600160a060020a0316331461042457600080fd5b60008054604051600160a060020a03909116917ff8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c6482091a26000805473ffffffffffffffffffffffffffffffffffffffff19169055565b600054600160a060020a031681565b67ffffffffffffffff16600090815260026020526040902060010154600160a060020a031690565b60008054600160a060020a031633146104c857600080fd5b5067ffffffffffffffff83166000818152600260208181526040808420805482516060810184526001808252600160a060020a03808c168388019081528b821696840196875299909852959094529251955185166101000274ffffffffffffffffffffffffffffffffffffffff001996151560ff1985161796909616959095178555519190930180549190921673ffffffffffffffffffffffffffffffffffffffff1990911617905560ff1680156105d2576040805167ffffffffffffffff86168152600160a060020a03808616602083015284168183015290517ffa11b2b19489638a86cff6fcb6aa85783c0a8503fdf36fda4ae42bb4156cad519181900360600190a1610626565b6040805167ffffffffffffffff86168152600160a060020a03808616602083015284168183015290517ff84f53f26b0e14dd3102f23ddd5488d7d4f6b05a7fa0881e1e030392c289000d9181900360600190a15b50505050565b600054600160a060020a0316331461064357600080fd5b61064c8161072b565b50565b6002602052600090815260409020805460019091015460ff821691600160a060020a036101009091048116911683565b67ffffffffffffffff1660009081526002602052604090205460ff1690565b60018054604080516020600284861615610100026000190190941693909304601f810184900484028201840190925281815292918301828280156107235780601f106106f857610100808354040283529160200191610723565b820191906000526020600020905b81548152906001019060200180831161070657829003601f168201915b505050505081565b600160a060020a038116151561074057600080fd5b60008054604051600160a060020a03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a36000805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03929092169190911790555600a165627a7a723058206e221f1fa4498e6a0716d27124f8921ec61b882e3c82466e6b9e7f9ab865c18a0029`

// DeploySettlementRegistry deploys a new Ethereum contract, binding an instance of SettlementRegistry to it.
func DeploySettlementRegistry(auth *bind.TransactOpts, backend bind.ContractBackend, _VERSION string) (common.Address, *types.Transaction, *SettlementRegistry, error) {
	parsed, err := abi.JSON(strings.NewReader(SettlementRegistryABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SettlementRegistryBin), backend, _VERSION)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SettlementRegistry{SettlementRegistryCaller: SettlementRegistryCaller{contract: contract}, SettlementRegistryTransactor: SettlementRegistryTransactor{contract: contract}, SettlementRegistryFilterer: SettlementRegistryFilterer{contract: contract}}, nil
}

// SettlementRegistry is an auto generated Go binding around an Ethereum contract.
type SettlementRegistry struct {
	SettlementRegistryCaller     // Read-only binding to the contract
	SettlementRegistryTransactor // Write-only binding to the contract
	SettlementRegistryFilterer   // Log filterer for contract events
}

// SettlementRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type SettlementRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SettlementRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SettlementRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SettlementRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SettlementRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SettlementRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SettlementRegistrySession struct {
	Contract     *SettlementRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// SettlementRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SettlementRegistryCallerSession struct {
	Contract *SettlementRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// SettlementRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SettlementRegistryTransactorSession struct {
	Contract     *SettlementRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// SettlementRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type SettlementRegistryRaw struct {
	Contract *SettlementRegistry // Generic contract binding to access the raw methods on
}

// SettlementRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SettlementRegistryCallerRaw struct {
	Contract *SettlementRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// SettlementRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SettlementRegistryTransactorRaw struct {
	Contract *SettlementRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSettlementRegistry creates a new instance of SettlementRegistry, bound to a specific deployed contract.
func NewSettlementRegistry(address common.Address, backend bind.ContractBackend) (*SettlementRegistry, error) {
	contract, err := bindSettlementRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SettlementRegistry{SettlementRegistryCaller: SettlementRegistryCaller{contract: contract}, SettlementRegistryTransactor: SettlementRegistryTransactor{contract: contract}, SettlementRegistryFilterer: SettlementRegistryFilterer{contract: contract}}, nil
}

// NewSettlementRegistryCaller creates a new read-only instance of SettlementRegistry, bound to a specific deployed contract.
func NewSettlementRegistryCaller(address common.Address, caller bind.ContractCaller) (*SettlementRegistryCaller, error) {
	contract, err := bindSettlementRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SettlementRegistryCaller{contract: contract}, nil
}

// NewSettlementRegistryTransactor creates a new write-only instance of SettlementRegistry, bound to a specific deployed contract.
func NewSettlementRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*SettlementRegistryTransactor, error) {
	contract, err := bindSettlementRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SettlementRegistryTransactor{contract: contract}, nil
}

// NewSettlementRegistryFilterer creates a new log filterer instance of SettlementRegistry, bound to a specific deployed contract.
func NewSettlementRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*SettlementRegistryFilterer, error) {
	contract, err := bindSettlementRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SettlementRegistryFilterer{contract: contract}, nil
}

// bindSettlementRegistry binds a generic wrapper to an already deployed contract.
func bindSettlementRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SettlementRegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SettlementRegistry *SettlementRegistryRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SettlementRegistry.Contract.SettlementRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SettlementRegistry *SettlementRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SettlementRegistry.Contract.SettlementRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SettlementRegistry *SettlementRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SettlementRegistry.Contract.SettlementRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SettlementRegistry *SettlementRegistryCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SettlementRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SettlementRegistry *SettlementRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SettlementRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SettlementRegistry *SettlementRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SettlementRegistry.Contract.contract.Transact(opts, method, params...)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_SettlementRegistry *SettlementRegistryCaller) VERSION(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _SettlementRegistry.contract.Call(opts, out, "VERSION")
	return *ret0, err
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_SettlementRegistry *SettlementRegistrySession) VERSION() (string, error) {
	return _SettlementRegistry.Contract.VERSION(&_SettlementRegistry.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_SettlementRegistry *SettlementRegistryCallerSession) VERSION() (string, error) {
	return _SettlementRegistry.Contract.VERSION(&_SettlementRegistry.CallOpts)
}

// BrokerVerifierContract is a free data retrieval call binding the contract method 0xa6065c96.
//
// Solidity: function brokerVerifierContract(_settlementID uint64) constant returns(address)
func (_SettlementRegistry *SettlementRegistryCaller) BrokerVerifierContract(opts *bind.CallOpts, _settlementID uint64) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _SettlementRegistry.contract.Call(opts, out, "brokerVerifierContract", _settlementID)
	return *ret0, err
}

// BrokerVerifierContract is a free data retrieval call binding the contract method 0xa6065c96.
//
// Solidity: function brokerVerifierContract(_settlementID uint64) constant returns(address)
func (_SettlementRegistry *SettlementRegistrySession) BrokerVerifierContract(_settlementID uint64) (common.Address, error) {
	return _SettlementRegistry.Contract.BrokerVerifierContract(&_SettlementRegistry.CallOpts, _settlementID)
}

// BrokerVerifierContract is a free data retrieval call binding the contract method 0xa6065c96.
//
// Solidity: function brokerVerifierContract(_settlementID uint64) constant returns(address)
func (_SettlementRegistry *SettlementRegistryCallerSession) BrokerVerifierContract(_settlementID uint64) (common.Address, error) {
	return _SettlementRegistry.Contract.BrokerVerifierContract(&_SettlementRegistry.CallOpts, _settlementID)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_SettlementRegistry *SettlementRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _SettlementRegistry.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_SettlementRegistry *SettlementRegistrySession) Owner() (common.Address, error) {
	return _SettlementRegistry.Contract.Owner(&_SettlementRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_SettlementRegistry *SettlementRegistryCallerSession) Owner() (common.Address, error) {
	return _SettlementRegistry.Contract.Owner(&_SettlementRegistry.CallOpts)
}

// SettlementContract is a free data retrieval call binding the contract method 0x3809b240.
//
// Solidity: function settlementContract(_settlementID uint64) constant returns(address)
func (_SettlementRegistry *SettlementRegistryCaller) SettlementContract(opts *bind.CallOpts, _settlementID uint64) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _SettlementRegistry.contract.Call(opts, out, "settlementContract", _settlementID)
	return *ret0, err
}

// SettlementContract is a free data retrieval call binding the contract method 0x3809b240.
//
// Solidity: function settlementContract(_settlementID uint64) constant returns(address)
func (_SettlementRegistry *SettlementRegistrySession) SettlementContract(_settlementID uint64) (common.Address, error) {
	return _SettlementRegistry.Contract.SettlementContract(&_SettlementRegistry.CallOpts, _settlementID)
}

// SettlementContract is a free data retrieval call binding the contract method 0x3809b240.
//
// Solidity: function settlementContract(_settlementID uint64) constant returns(address)
func (_SettlementRegistry *SettlementRegistryCallerSession) SettlementContract(_settlementID uint64) (common.Address, error) {
	return _SettlementRegistry.Contract.SettlementContract(&_SettlementRegistry.CallOpts, _settlementID)
}

// SettlementDetails is a free data retrieval call binding the contract method 0xf8a09cd0.
//
// Solidity: function settlementDetails( uint64) constant returns(registered bool, settlementContract address, brokerVerifierContract address)
func (_SettlementRegistry *SettlementRegistryCaller) SettlementDetails(opts *bind.CallOpts, arg0 uint64) (struct {
	Registered             bool
	SettlementContract     common.Address
	BrokerVerifierContract common.Address
}, error) {
	ret := new(struct {
		Registered             bool
		SettlementContract     common.Address
		BrokerVerifierContract common.Address
	})
	out := ret
	err := _SettlementRegistry.contract.Call(opts, out, "settlementDetails", arg0)
	return *ret, err
}

// SettlementDetails is a free data retrieval call binding the contract method 0xf8a09cd0.
//
// Solidity: function settlementDetails( uint64) constant returns(registered bool, settlementContract address, brokerVerifierContract address)
func (_SettlementRegistry *SettlementRegistrySession) SettlementDetails(arg0 uint64) (struct {
	Registered             bool
	SettlementContract     common.Address
	BrokerVerifierContract common.Address
}, error) {
	return _SettlementRegistry.Contract.SettlementDetails(&_SettlementRegistry.CallOpts, arg0)
}

// SettlementDetails is a free data retrieval call binding the contract method 0xf8a09cd0.
//
// Solidity: function settlementDetails( uint64) constant returns(registered bool, settlementContract address, brokerVerifierContract address)
func (_SettlementRegistry *SettlementRegistryCallerSession) SettlementDetails(arg0 uint64) (struct {
	Registered             bool
	SettlementContract     common.Address
	BrokerVerifierContract common.Address
}, error) {
	return _SettlementRegistry.Contract.SettlementDetails(&_SettlementRegistry.CallOpts, arg0)
}

// SettlementRegistration is a free data retrieval call binding the contract method 0xf8f9be36.
//
// Solidity: function settlementRegistration(_settlementID uint64) constant returns(bool)
func (_SettlementRegistry *SettlementRegistryCaller) SettlementRegistration(opts *bind.CallOpts, _settlementID uint64) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _SettlementRegistry.contract.Call(opts, out, "settlementRegistration", _settlementID)
	return *ret0, err
}

// SettlementRegistration is a free data retrieval call binding the contract method 0xf8f9be36.
//
// Solidity: function settlementRegistration(_settlementID uint64) constant returns(bool)
func (_SettlementRegistry *SettlementRegistrySession) SettlementRegistration(_settlementID uint64) (bool, error) {
	return _SettlementRegistry.Contract.SettlementRegistration(&_SettlementRegistry.CallOpts, _settlementID)
}

// SettlementRegistration is a free data retrieval call binding the contract method 0xf8f9be36.
//
// Solidity: function settlementRegistration(_settlementID uint64) constant returns(bool)
func (_SettlementRegistry *SettlementRegistryCallerSession) SettlementRegistration(_settlementID uint64) (bool, error) {
	return _SettlementRegistry.Contract.SettlementRegistration(&_SettlementRegistry.CallOpts, _settlementID)
}

// DeregisterSettlement is a paid mutator transaction binding the contract method 0x06b82a7d.
//
// Solidity: function deregisterSettlement(_settlementID uint64) returns()
func (_SettlementRegistry *SettlementRegistryTransactor) DeregisterSettlement(opts *bind.TransactOpts, _settlementID uint64) (*types.Transaction, error) {
	return _SettlementRegistry.contract.Transact(opts, "deregisterSettlement", _settlementID)
}

// DeregisterSettlement is a paid mutator transaction binding the contract method 0x06b82a7d.
//
// Solidity: function deregisterSettlement(_settlementID uint64) returns()
func (_SettlementRegistry *SettlementRegistrySession) DeregisterSettlement(_settlementID uint64) (*types.Transaction, error) {
	return _SettlementRegistry.Contract.DeregisterSettlement(&_SettlementRegistry.TransactOpts, _settlementID)
}

// DeregisterSettlement is a paid mutator transaction binding the contract method 0x06b82a7d.
//
// Solidity: function deregisterSettlement(_settlementID uint64) returns()
func (_SettlementRegistry *SettlementRegistryTransactorSession) DeregisterSettlement(_settlementID uint64) (*types.Transaction, error) {
	return _SettlementRegistry.Contract.DeregisterSettlement(&_SettlementRegistry.TransactOpts, _settlementID)
}

// RegisterSettlement is a paid mutator transaction binding the contract method 0xc2190c94.
//
// Solidity: function registerSettlement(_settlementID uint64, _settlementContract address, _brokerVerifierContract address) returns()
func (_SettlementRegistry *SettlementRegistryTransactor) RegisterSettlement(opts *bind.TransactOpts, _settlementID uint64, _settlementContract common.Address, _brokerVerifierContract common.Address) (*types.Transaction, error) {
	return _SettlementRegistry.contract.Transact(opts, "registerSettlement", _settlementID, _settlementContract, _brokerVerifierContract)
}

// RegisterSettlement is a paid mutator transaction binding the contract method 0xc2190c94.
//
// Solidity: function registerSettlement(_settlementID uint64, _settlementContract address, _brokerVerifierContract address) returns()
func (_SettlementRegistry *SettlementRegistrySession) RegisterSettlement(_settlementID uint64, _settlementContract common.Address, _brokerVerifierContract common.Address) (*types.Transaction, error) {
	return _SettlementRegistry.Contract.RegisterSettlement(&_SettlementRegistry.TransactOpts, _settlementID, _settlementContract, _brokerVerifierContract)
}

// RegisterSettlement is a paid mutator transaction binding the contract method 0xc2190c94.
//
// Solidity: function registerSettlement(_settlementID uint64, _settlementContract address, _brokerVerifierContract address) returns()
func (_SettlementRegistry *SettlementRegistryTransactorSession) RegisterSettlement(_settlementID uint64, _settlementContract common.Address, _brokerVerifierContract common.Address) (*types.Transaction, error) {
	return _SettlementRegistry.Contract.RegisterSettlement(&_SettlementRegistry.TransactOpts, _settlementID, _settlementContract, _brokerVerifierContract)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SettlementRegistry *SettlementRegistryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SettlementRegistry.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SettlementRegistry *SettlementRegistrySession) RenounceOwnership() (*types.Transaction, error) {
	return _SettlementRegistry.Contract.RenounceOwnership(&_SettlementRegistry.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SettlementRegistry *SettlementRegistryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _SettlementRegistry.Contract.RenounceOwnership(&_SettlementRegistry.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_SettlementRegistry *SettlementRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _SettlementRegistry.contract.Transact(opts, "transferOwnership", _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_SettlementRegistry *SettlementRegistrySession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _SettlementRegistry.Contract.TransferOwnership(&_SettlementRegistry.TransactOpts, _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_SettlementRegistry *SettlementRegistryTransactorSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _SettlementRegistry.Contract.TransferOwnership(&_SettlementRegistry.TransactOpts, _newOwner)
}

// SettlementRegistryLogSettlementDeregisteredIterator is returned from FilterLogSettlementDeregistered and is used to iterate over the raw logs and unpacked data for LogSettlementDeregistered events raised by the SettlementRegistry contract.
type SettlementRegistryLogSettlementDeregisteredIterator struct {
	Event *SettlementRegistryLogSettlementDeregistered // Event containing the contract specifics and raw log

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
func (it *SettlementRegistryLogSettlementDeregisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementRegistryLogSettlementDeregistered)
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
		it.Event = new(SettlementRegistryLogSettlementDeregistered)
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
func (it *SettlementRegistryLogSettlementDeregisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementRegistryLogSettlementDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementRegistryLogSettlementDeregistered represents a LogSettlementDeregistered event raised by the SettlementRegistry contract.
type SettlementRegistryLogSettlementDeregistered struct {
	SettlementID uint64
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterLogSettlementDeregistered is a free log retrieval operation binding the contract event 0x9e153628c2e49d8816ef9c792fac73678b5846d086e651aa7a6ea95f1b7f6675.
//
// Solidity: e LogSettlementDeregistered(settlementID uint64)
func (_SettlementRegistry *SettlementRegistryFilterer) FilterLogSettlementDeregistered(opts *bind.FilterOpts) (*SettlementRegistryLogSettlementDeregisteredIterator, error) {

	logs, sub, err := _SettlementRegistry.contract.FilterLogs(opts, "LogSettlementDeregistered")
	if err != nil {
		return nil, err
	}
	return &SettlementRegistryLogSettlementDeregisteredIterator{contract: _SettlementRegistry.contract, event: "LogSettlementDeregistered", logs: logs, sub: sub}, nil
}

// WatchLogSettlementDeregistered is a free log subscription operation binding the contract event 0x9e153628c2e49d8816ef9c792fac73678b5846d086e651aa7a6ea95f1b7f6675.
//
// Solidity: e LogSettlementDeregistered(settlementID uint64)
func (_SettlementRegistry *SettlementRegistryFilterer) WatchLogSettlementDeregistered(opts *bind.WatchOpts, sink chan<- *SettlementRegistryLogSettlementDeregistered) (event.Subscription, error) {

	logs, sub, err := _SettlementRegistry.contract.WatchLogs(opts, "LogSettlementDeregistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementRegistryLogSettlementDeregistered)
				if err := _SettlementRegistry.contract.UnpackLog(event, "LogSettlementDeregistered", log); err != nil {
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

// SettlementRegistryLogSettlementRegisteredIterator is returned from FilterLogSettlementRegistered and is used to iterate over the raw logs and unpacked data for LogSettlementRegistered events raised by the SettlementRegistry contract.
type SettlementRegistryLogSettlementRegisteredIterator struct {
	Event *SettlementRegistryLogSettlementRegistered // Event containing the contract specifics and raw log

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
func (it *SettlementRegistryLogSettlementRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementRegistryLogSettlementRegistered)
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
		it.Event = new(SettlementRegistryLogSettlementRegistered)
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
func (it *SettlementRegistryLogSettlementRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementRegistryLogSettlementRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementRegistryLogSettlementRegistered represents a LogSettlementRegistered event raised by the SettlementRegistry contract.
type SettlementRegistryLogSettlementRegistered struct {
	SettlementID           uint64
	SettlementContract     common.Address
	BrokerVerifierContract common.Address
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterLogSettlementRegistered is a free log retrieval operation binding the contract event 0xf84f53f26b0e14dd3102f23ddd5488d7d4f6b05a7fa0881e1e030392c289000d.
//
// Solidity: e LogSettlementRegistered(settlementID uint64, settlementContract address, brokerVerifierContract address)
func (_SettlementRegistry *SettlementRegistryFilterer) FilterLogSettlementRegistered(opts *bind.FilterOpts) (*SettlementRegistryLogSettlementRegisteredIterator, error) {

	logs, sub, err := _SettlementRegistry.contract.FilterLogs(opts, "LogSettlementRegistered")
	if err != nil {
		return nil, err
	}
	return &SettlementRegistryLogSettlementRegisteredIterator{contract: _SettlementRegistry.contract, event: "LogSettlementRegistered", logs: logs, sub: sub}, nil
}

// WatchLogSettlementRegistered is a free log subscription operation binding the contract event 0xf84f53f26b0e14dd3102f23ddd5488d7d4f6b05a7fa0881e1e030392c289000d.
//
// Solidity: e LogSettlementRegistered(settlementID uint64, settlementContract address, brokerVerifierContract address)
func (_SettlementRegistry *SettlementRegistryFilterer) WatchLogSettlementRegistered(opts *bind.WatchOpts, sink chan<- *SettlementRegistryLogSettlementRegistered) (event.Subscription, error) {

	logs, sub, err := _SettlementRegistry.contract.WatchLogs(opts, "LogSettlementRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementRegistryLogSettlementRegistered)
				if err := _SettlementRegistry.contract.UnpackLog(event, "LogSettlementRegistered", log); err != nil {
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

// SettlementRegistryLogSettlementUpdatedIterator is returned from FilterLogSettlementUpdated and is used to iterate over the raw logs and unpacked data for LogSettlementUpdated events raised by the SettlementRegistry contract.
type SettlementRegistryLogSettlementUpdatedIterator struct {
	Event *SettlementRegistryLogSettlementUpdated // Event containing the contract specifics and raw log

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
func (it *SettlementRegistryLogSettlementUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementRegistryLogSettlementUpdated)
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
		it.Event = new(SettlementRegistryLogSettlementUpdated)
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
func (it *SettlementRegistryLogSettlementUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementRegistryLogSettlementUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementRegistryLogSettlementUpdated represents a LogSettlementUpdated event raised by the SettlementRegistry contract.
type SettlementRegistryLogSettlementUpdated struct {
	SettlementID           uint64
	SettlementContract     common.Address
	BrokerVerifierContract common.Address
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterLogSettlementUpdated is a free log retrieval operation binding the contract event 0xfa11b2b19489638a86cff6fcb6aa85783c0a8503fdf36fda4ae42bb4156cad51.
//
// Solidity: e LogSettlementUpdated(settlementID uint64, settlementContract address, brokerVerifierContract address)
func (_SettlementRegistry *SettlementRegistryFilterer) FilterLogSettlementUpdated(opts *bind.FilterOpts) (*SettlementRegistryLogSettlementUpdatedIterator, error) {

	logs, sub, err := _SettlementRegistry.contract.FilterLogs(opts, "LogSettlementUpdated")
	if err != nil {
		return nil, err
	}
	return &SettlementRegistryLogSettlementUpdatedIterator{contract: _SettlementRegistry.contract, event: "LogSettlementUpdated", logs: logs, sub: sub}, nil
}

// WatchLogSettlementUpdated is a free log subscription operation binding the contract event 0xfa11b2b19489638a86cff6fcb6aa85783c0a8503fdf36fda4ae42bb4156cad51.
//
// Solidity: e LogSettlementUpdated(settlementID uint64, settlementContract address, brokerVerifierContract address)
func (_SettlementRegistry *SettlementRegistryFilterer) WatchLogSettlementUpdated(opts *bind.WatchOpts, sink chan<- *SettlementRegistryLogSettlementUpdated) (event.Subscription, error) {

	logs, sub, err := _SettlementRegistry.contract.WatchLogs(opts, "LogSettlementUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementRegistryLogSettlementUpdated)
				if err := _SettlementRegistry.contract.UnpackLog(event, "LogSettlementUpdated", log); err != nil {
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

// SettlementRegistryOwnershipRenouncedIterator is returned from FilterOwnershipRenounced and is used to iterate over the raw logs and unpacked data for OwnershipRenounced events raised by the SettlementRegistry contract.
type SettlementRegistryOwnershipRenouncedIterator struct {
	Event *SettlementRegistryOwnershipRenounced // Event containing the contract specifics and raw log

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
func (it *SettlementRegistryOwnershipRenouncedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementRegistryOwnershipRenounced)
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
		it.Event = new(SettlementRegistryOwnershipRenounced)
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
func (it *SettlementRegistryOwnershipRenouncedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementRegistryOwnershipRenouncedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementRegistryOwnershipRenounced represents a OwnershipRenounced event raised by the SettlementRegistry contract.
type SettlementRegistryOwnershipRenounced struct {
	PreviousOwner common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipRenounced is a free log retrieval operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_SettlementRegistry *SettlementRegistryFilterer) FilterOwnershipRenounced(opts *bind.FilterOpts, previousOwner []common.Address) (*SettlementRegistryOwnershipRenouncedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _SettlementRegistry.contract.FilterLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SettlementRegistryOwnershipRenouncedIterator{contract: _SettlementRegistry.contract, event: "OwnershipRenounced", logs: logs, sub: sub}, nil
}

// WatchOwnershipRenounced is a free log subscription operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_SettlementRegistry *SettlementRegistryFilterer) WatchOwnershipRenounced(opts *bind.WatchOpts, sink chan<- *SettlementRegistryOwnershipRenounced, previousOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _SettlementRegistry.contract.WatchLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementRegistryOwnershipRenounced)
				if err := _SettlementRegistry.contract.UnpackLog(event, "OwnershipRenounced", log); err != nil {
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

// SettlementRegistryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the SettlementRegistry contract.
type SettlementRegistryOwnershipTransferredIterator struct {
	Event *SettlementRegistryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *SettlementRegistryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementRegistryOwnershipTransferred)
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
		it.Event = new(SettlementRegistryOwnershipTransferred)
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
func (it *SettlementRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementRegistryOwnershipTransferred represents a OwnershipTransferred event raised by the SettlementRegistry contract.
type SettlementRegistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_SettlementRegistry *SettlementRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*SettlementRegistryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SettlementRegistry.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SettlementRegistryOwnershipTransferredIterator{contract: _SettlementRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_SettlementRegistry *SettlementRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SettlementRegistryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SettlementRegistry.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementRegistryOwnershipTransferred)
				if err := _SettlementRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// SettlementUtilsABI is the input ABI used to generate the binding from.
const SettlementUtilsABI = "[]"

// SettlementUtilsBin is the compiled bytecode used for deploying new contracts.
const SettlementUtilsBin = `0x604c602c600b82828239805160001a60731460008114601c57601e565bfe5b5030600052607381538281f30073000000000000000000000000000000000000000030146080604052600080fd00a165627a7a723058201d14e882b4e9e2061c75107ba9740a5361b4214d7d24504574e070a582e2d9de0029`

// DeploySettlementUtils deploys a new Ethereum contract, binding an instance of SettlementUtils to it.
func DeploySettlementUtils(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SettlementUtils, error) {
	parsed, err := abi.JSON(strings.NewReader(SettlementUtilsABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SettlementUtilsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SettlementUtils{SettlementUtilsCaller: SettlementUtilsCaller{contract: contract}, SettlementUtilsTransactor: SettlementUtilsTransactor{contract: contract}, SettlementUtilsFilterer: SettlementUtilsFilterer{contract: contract}}, nil
}

// SettlementUtils is an auto generated Go binding around an Ethereum contract.
type SettlementUtils struct {
	SettlementUtilsCaller     // Read-only binding to the contract
	SettlementUtilsTransactor // Write-only binding to the contract
	SettlementUtilsFilterer   // Log filterer for contract events
}

// SettlementUtilsCaller is an auto generated read-only Go binding around an Ethereum contract.
type SettlementUtilsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SettlementUtilsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SettlementUtilsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SettlementUtilsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SettlementUtilsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SettlementUtilsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SettlementUtilsSession struct {
	Contract     *SettlementUtils  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SettlementUtilsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SettlementUtilsCallerSession struct {
	Contract *SettlementUtilsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// SettlementUtilsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SettlementUtilsTransactorSession struct {
	Contract     *SettlementUtilsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// SettlementUtilsRaw is an auto generated low-level Go binding around an Ethereum contract.
type SettlementUtilsRaw struct {
	Contract *SettlementUtils // Generic contract binding to access the raw methods on
}

// SettlementUtilsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SettlementUtilsCallerRaw struct {
	Contract *SettlementUtilsCaller // Generic read-only contract binding to access the raw methods on
}

// SettlementUtilsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SettlementUtilsTransactorRaw struct {
	Contract *SettlementUtilsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSettlementUtils creates a new instance of SettlementUtils, bound to a specific deployed contract.
func NewSettlementUtils(address common.Address, backend bind.ContractBackend) (*SettlementUtils, error) {
	contract, err := bindSettlementUtils(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SettlementUtils{SettlementUtilsCaller: SettlementUtilsCaller{contract: contract}, SettlementUtilsTransactor: SettlementUtilsTransactor{contract: contract}, SettlementUtilsFilterer: SettlementUtilsFilterer{contract: contract}}, nil
}

// NewSettlementUtilsCaller creates a new read-only instance of SettlementUtils, bound to a specific deployed contract.
func NewSettlementUtilsCaller(address common.Address, caller bind.ContractCaller) (*SettlementUtilsCaller, error) {
	contract, err := bindSettlementUtils(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SettlementUtilsCaller{contract: contract}, nil
}

// NewSettlementUtilsTransactor creates a new write-only instance of SettlementUtils, bound to a specific deployed contract.
func NewSettlementUtilsTransactor(address common.Address, transactor bind.ContractTransactor) (*SettlementUtilsTransactor, error) {
	contract, err := bindSettlementUtils(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SettlementUtilsTransactor{contract: contract}, nil
}

// NewSettlementUtilsFilterer creates a new log filterer instance of SettlementUtils, bound to a specific deployed contract.
func NewSettlementUtilsFilterer(address common.Address, filterer bind.ContractFilterer) (*SettlementUtilsFilterer, error) {
	contract, err := bindSettlementUtils(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SettlementUtilsFilterer{contract: contract}, nil
}

// bindSettlementUtils binds a generic wrapper to an already deployed contract.
func bindSettlementUtils(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SettlementUtilsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SettlementUtils *SettlementUtilsRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SettlementUtils.Contract.SettlementUtilsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SettlementUtils *SettlementUtilsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SettlementUtils.Contract.SettlementUtilsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SettlementUtils *SettlementUtilsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SettlementUtils.Contract.SettlementUtilsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SettlementUtils *SettlementUtilsCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SettlementUtils.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SettlementUtils *SettlementUtilsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SettlementUtils.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SettlementUtils *SettlementUtilsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SettlementUtils.Contract.contract.Transact(opts, method, params...)
}

// StandardTokenABI is the input ABI used to generate the binding from.
const StandardTokenABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseApproval\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseApproval\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]"

// StandardTokenBin is the compiled bytecode used for deploying new contracts.
const StandardTokenBin = `0x608060405234801561001057600080fd5b506106b3806100206000396000f30060806040526004361061008d5763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663095ea7b3811461009257806318160ddd146100ca57806323b872dd146100f1578063661884631461011b57806370a082311461013f578063a9059cbb14610160578063d73dd62314610184578063dd62ed3e146101a8575b600080fd5b34801561009e57600080fd5b506100b6600160a060020a03600435166024356101cf565b604080519115158252519081900360200190f35b3480156100d657600080fd5b506100df610235565b60408051918252519081900360200190f35b3480156100fd57600080fd5b506100b6600160a060020a036004358116906024351660443561023b565b34801561012757600080fd5b506100b6600160a060020a03600435166024356103b2565b34801561014b57600080fd5b506100df600160a060020a03600435166104a2565b34801561016c57600080fd5b506100b6600160a060020a03600435166024356104bd565b34801561019057600080fd5b506100b6600160a060020a036004351660243561059e565b3480156101b457600080fd5b506100df600160a060020a0360043581169060243516610637565b336000818152600260209081526040808320600160a060020a038716808552908352818420869055815186815291519394909390927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925928290030190a350600192915050565b60015490565b6000600160a060020a038316151561025257600080fd5b600160a060020a03841660009081526020819052604090205482111561027757600080fd5b600160a060020a03841660009081526002602090815260408083203384529091529020548211156102a757600080fd5b600160a060020a0384166000908152602081905260409020546102d0908363ffffffff61066216565b600160a060020a038086166000908152602081905260408082209390935590851681522054610305908363ffffffff61067416565b600160a060020a03808516600090815260208181526040808320949094559187168152600282528281203382529091522054610347908363ffffffff61066216565b600160a060020a03808616600081815260026020908152604080832033845282529182902094909455805186815290519287169391927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef929181900390910190a35060019392505050565b336000908152600260209081526040808320600160a060020a03861684529091528120548083111561040757336000908152600260209081526040808320600160a060020a038816845290915281205561043c565b610417818463ffffffff61066216565b336000908152600260209081526040808320600160a060020a03891684529091529020555b336000818152600260209081526040808320600160a060020a0389168085529083529281902054815190815290519293927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925929181900390910190a35060019392505050565b600160a060020a031660009081526020819052604090205490565b6000600160a060020a03831615156104d457600080fd5b336000908152602081905260409020548211156104f057600080fd5b33600090815260208190526040902054610510908363ffffffff61066216565b3360009081526020819052604080822092909255600160a060020a03851681522054610542908363ffffffff61067416565b600160a060020a038416600081815260208181526040918290209390935580518581529051919233927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9281900390910190a350600192915050565b336000908152600260209081526040808320600160a060020a03861684529091528120546105d2908363ffffffff61067416565b336000818152600260209081526040808320600160a060020a0389168085529083529281902085905580519485525191937f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925929081900390910190a350600192915050565b600160a060020a03918216600090815260026020908152604080832093909416825291909152205490565b60008282111561066e57fe5b50900390565b8181018281101561068157fe5b929150505600a165627a7a7230582002a8ac0e1fb78e81f525672fea77888ca221a22f23d1d03ea17a5f794f4c07950029`

// DeployStandardToken deploys a new Ethereum contract, binding an instance of StandardToken to it.
func DeployStandardToken(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *StandardToken, error) {
	parsed, err := abi.JSON(strings.NewReader(StandardTokenABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(StandardTokenBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &StandardToken{StandardTokenCaller: StandardTokenCaller{contract: contract}, StandardTokenTransactor: StandardTokenTransactor{contract: contract}, StandardTokenFilterer: StandardTokenFilterer{contract: contract}}, nil
}

// StandardToken is an auto generated Go binding around an Ethereum contract.
type StandardToken struct {
	StandardTokenCaller     // Read-only binding to the contract
	StandardTokenTransactor // Write-only binding to the contract
	StandardTokenFilterer   // Log filterer for contract events
}

// StandardTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type StandardTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StandardTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StandardTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StandardTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StandardTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StandardTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StandardTokenSession struct {
	Contract     *StandardToken    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StandardTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StandardTokenCallerSession struct {
	Contract *StandardTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// StandardTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StandardTokenTransactorSession struct {
	Contract     *StandardTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// StandardTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type StandardTokenRaw struct {
	Contract *StandardToken // Generic contract binding to access the raw methods on
}

// StandardTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StandardTokenCallerRaw struct {
	Contract *StandardTokenCaller // Generic read-only contract binding to access the raw methods on
}

// StandardTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StandardTokenTransactorRaw struct {
	Contract *StandardTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStandardToken creates a new instance of StandardToken, bound to a specific deployed contract.
func NewStandardToken(address common.Address, backend bind.ContractBackend) (*StandardToken, error) {
	contract, err := bindStandardToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StandardToken{StandardTokenCaller: StandardTokenCaller{contract: contract}, StandardTokenTransactor: StandardTokenTransactor{contract: contract}, StandardTokenFilterer: StandardTokenFilterer{contract: contract}}, nil
}

// NewStandardTokenCaller creates a new read-only instance of StandardToken, bound to a specific deployed contract.
func NewStandardTokenCaller(address common.Address, caller bind.ContractCaller) (*StandardTokenCaller, error) {
	contract, err := bindStandardToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StandardTokenCaller{contract: contract}, nil
}

// NewStandardTokenTransactor creates a new write-only instance of StandardToken, bound to a specific deployed contract.
func NewStandardTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*StandardTokenTransactor, error) {
	contract, err := bindStandardToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StandardTokenTransactor{contract: contract}, nil
}

// NewStandardTokenFilterer creates a new log filterer instance of StandardToken, bound to a specific deployed contract.
func NewStandardTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*StandardTokenFilterer, error) {
	contract, err := bindStandardToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StandardTokenFilterer{contract: contract}, nil
}

// bindStandardToken binds a generic wrapper to an already deployed contract.
func bindStandardToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StandardTokenABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StandardToken *StandardTokenRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _StandardToken.Contract.StandardTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StandardToken *StandardTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StandardToken.Contract.StandardTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StandardToken *StandardTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StandardToken.Contract.StandardTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StandardToken *StandardTokenCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _StandardToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StandardToken *StandardTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StandardToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StandardToken *StandardTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StandardToken.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(uint256)
func (_StandardToken *StandardTokenCaller) Allowance(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _StandardToken.contract.Call(opts, out, "allowance", _owner, _spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(uint256)
func (_StandardToken *StandardTokenSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _StandardToken.Contract.Allowance(&_StandardToken.CallOpts, _owner, _spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(uint256)
func (_StandardToken *StandardTokenCallerSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _StandardToken.Contract.Allowance(&_StandardToken.CallOpts, _owner, _spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(uint256)
func (_StandardToken *StandardTokenCaller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _StandardToken.contract.Call(opts, out, "balanceOf", _owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(uint256)
func (_StandardToken *StandardTokenSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _StandardToken.Contract.BalanceOf(&_StandardToken.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(uint256)
func (_StandardToken *StandardTokenCallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _StandardToken.Contract.BalanceOf(&_StandardToken.CallOpts, _owner)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_StandardToken *StandardTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _StandardToken.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_StandardToken *StandardTokenSession) TotalSupply() (*big.Int, error) {
	return _StandardToken.Contract.TotalSupply(&_StandardToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_StandardToken *StandardTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _StandardToken.Contract.TotalSupply(&_StandardToken.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(bool)
func (_StandardToken *StandardTokenTransactor) Approve(opts *bind.TransactOpts, _spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _StandardToken.contract.Transact(opts, "approve", _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(bool)
func (_StandardToken *StandardTokenSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _StandardToken.Contract.Approve(&_StandardToken.TransactOpts, _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(bool)
func (_StandardToken *StandardTokenTransactorSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _StandardToken.Contract.Approve(&_StandardToken.TransactOpts, _spender, _value)
}

// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
//
// Solidity: function decreaseApproval(_spender address, _subtractedValue uint256) returns(bool)
func (_StandardToken *StandardTokenTransactor) DecreaseApproval(opts *bind.TransactOpts, _spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
	return _StandardToken.contract.Transact(opts, "decreaseApproval", _spender, _subtractedValue)
}

// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
//
// Solidity: function decreaseApproval(_spender address, _subtractedValue uint256) returns(bool)
func (_StandardToken *StandardTokenSession) DecreaseApproval(_spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
	return _StandardToken.Contract.DecreaseApproval(&_StandardToken.TransactOpts, _spender, _subtractedValue)
}

// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
//
// Solidity: function decreaseApproval(_spender address, _subtractedValue uint256) returns(bool)
func (_StandardToken *StandardTokenTransactorSession) DecreaseApproval(_spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
	return _StandardToken.Contract.DecreaseApproval(&_StandardToken.TransactOpts, _spender, _subtractedValue)
}

// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
//
// Solidity: function increaseApproval(_spender address, _addedValue uint256) returns(bool)
func (_StandardToken *StandardTokenTransactor) IncreaseApproval(opts *bind.TransactOpts, _spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
	return _StandardToken.contract.Transact(opts, "increaseApproval", _spender, _addedValue)
}

// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
//
// Solidity: function increaseApproval(_spender address, _addedValue uint256) returns(bool)
func (_StandardToken *StandardTokenSession) IncreaseApproval(_spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
	return _StandardToken.Contract.IncreaseApproval(&_StandardToken.TransactOpts, _spender, _addedValue)
}

// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
//
// Solidity: function increaseApproval(_spender address, _addedValue uint256) returns(bool)
func (_StandardToken *StandardTokenTransactorSession) IncreaseApproval(_spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
	return _StandardToken.Contract.IncreaseApproval(&_StandardToken.TransactOpts, _spender, _addedValue)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_StandardToken *StandardTokenTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _StandardToken.contract.Transact(opts, "transfer", _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_StandardToken *StandardTokenSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _StandardToken.Contract.Transfer(&_StandardToken.TransactOpts, _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_StandardToken *StandardTokenTransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _StandardToken.Contract.Transfer(&_StandardToken.TransactOpts, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
func (_StandardToken *StandardTokenTransactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _StandardToken.contract.Transact(opts, "transferFrom", _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
func (_StandardToken *StandardTokenSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _StandardToken.Contract.TransferFrom(&_StandardToken.TransactOpts, _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
func (_StandardToken *StandardTokenTransactorSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _StandardToken.Contract.TransferFrom(&_StandardToken.TransactOpts, _from, _to, _value)
}

// StandardTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the StandardToken contract.
type StandardTokenApprovalIterator struct {
	Event *StandardTokenApproval // Event containing the contract specifics and raw log

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
func (it *StandardTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StandardTokenApproval)
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
		it.Event = new(StandardTokenApproval)
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
func (it *StandardTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StandardTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StandardTokenApproval represents a Approval event raised by the StandardToken contract.
type StandardTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(owner indexed address, spender indexed address, value uint256)
func (_StandardToken *StandardTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*StandardTokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _StandardToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &StandardTokenApprovalIterator{contract: _StandardToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(owner indexed address, spender indexed address, value uint256)
func (_StandardToken *StandardTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *StandardTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _StandardToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StandardTokenApproval)
				if err := _StandardToken.contract.UnpackLog(event, "Approval", log); err != nil {
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

// StandardTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the StandardToken contract.
type StandardTokenTransferIterator struct {
	Event *StandardTokenTransfer // Event containing the contract specifics and raw log

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
func (it *StandardTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StandardTokenTransfer)
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
		it.Event = new(StandardTokenTransfer)
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
func (it *StandardTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StandardTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StandardTokenTransfer represents a Transfer event raised by the StandardToken contract.
type StandardTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
func (_StandardToken *StandardTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*StandardTokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _StandardToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &StandardTokenTransferIterator{contract: _StandardToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
func (_StandardToken *StandardTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *StandardTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _StandardToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StandardTokenTransfer)
				if err := _StandardToken.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// UtilsABI is the input ABI used to generate the binding from.
const UtilsABI = "[]"

// UtilsBin is the compiled bytecode used for deploying new contracts.
const UtilsBin = `0x604c602c600b82828239805160001a60731460008114601c57601e565bfe5b5030600052607381538281f30073000000000000000000000000000000000000000030146080604052600080fd00a165627a7a72305820fb1e7b496982848637c59ff3c93b47b3603fd2d38d7116ede15380ad621103bf0029`

// DeployUtils deploys a new Ethereum contract, binding an instance of Utils to it.
func DeployUtils(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Utils, error) {
	parsed, err := abi.JSON(strings.NewReader(UtilsABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(UtilsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Utils{UtilsCaller: UtilsCaller{contract: contract}, UtilsTransactor: UtilsTransactor{contract: contract}, UtilsFilterer: UtilsFilterer{contract: contract}}, nil
}

// Utils is an auto generated Go binding around an Ethereum contract.
type Utils struct {
	UtilsCaller     // Read-only binding to the contract
	UtilsTransactor // Write-only binding to the contract
	UtilsFilterer   // Log filterer for contract events
}

// UtilsCaller is an auto generated read-only Go binding around an Ethereum contract.
type UtilsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UtilsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UtilsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UtilsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UtilsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UtilsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UtilsSession struct {
	Contract     *Utils            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UtilsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UtilsCallerSession struct {
	Contract *UtilsCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// UtilsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UtilsTransactorSession struct {
	Contract     *UtilsTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UtilsRaw is an auto generated low-level Go binding around an Ethereum contract.
type UtilsRaw struct {
	Contract *Utils // Generic contract binding to access the raw methods on
}

// UtilsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UtilsCallerRaw struct {
	Contract *UtilsCaller // Generic read-only contract binding to access the raw methods on
}

// UtilsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UtilsTransactorRaw struct {
	Contract *UtilsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUtils creates a new instance of Utils, bound to a specific deployed contract.
func NewUtils(address common.Address, backend bind.ContractBackend) (*Utils, error) {
	contract, err := bindUtils(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Utils{UtilsCaller: UtilsCaller{contract: contract}, UtilsTransactor: UtilsTransactor{contract: contract}, UtilsFilterer: UtilsFilterer{contract: contract}}, nil
}

// NewUtilsCaller creates a new read-only instance of Utils, bound to a specific deployed contract.
func NewUtilsCaller(address common.Address, caller bind.ContractCaller) (*UtilsCaller, error) {
	contract, err := bindUtils(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UtilsCaller{contract: contract}, nil
}

// NewUtilsTransactor creates a new write-only instance of Utils, bound to a specific deployed contract.
func NewUtilsTransactor(address common.Address, transactor bind.ContractTransactor) (*UtilsTransactor, error) {
	contract, err := bindUtils(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UtilsTransactor{contract: contract}, nil
}

// NewUtilsFilterer creates a new log filterer instance of Utils, bound to a specific deployed contract.
func NewUtilsFilterer(address common.Address, filterer bind.ContractFilterer) (*UtilsFilterer, error) {
	contract, err := bindUtils(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UtilsFilterer{contract: contract}, nil
}

// bindUtils binds a generic wrapper to an already deployed contract.
func bindUtils(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(UtilsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Utils *UtilsRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Utils.Contract.UtilsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Utils *UtilsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Utils.Contract.UtilsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Utils *UtilsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Utils.Contract.UtilsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Utils *UtilsCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Utils.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Utils *UtilsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Utils.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Utils *UtilsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Utils.Contract.contract.Transact(opts, method, params...)
}
