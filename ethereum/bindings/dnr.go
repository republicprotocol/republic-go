package bindings

//
//import (
//	"math/big"
//	"strings"
//
//	"github.com/ethereum/go-ethereum"
//	"github.com/ethereum/go-ethereum/accounts/abi"
//	"github.com/ethereum/go-ethereum/accounts/abi/bind"
//	"github.com/ethereum/go-ethereum/common"
//	"github.com/ethereum/go-ethereum/core/types"
//	"github.com/ethereum/go-ethereum/event"
//)
//
//// BasicTokenABI is the input ABI used to generate the binding from.
//const BasicTokenABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]"
//
//// BasicTokenBin is the compiled bytecode used for deploying new contracts.
//const BasicTokenBin = `0x6060604052341561000f57600080fd5b6102818061001e6000396000f3006060604052600436106100615763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166318160ddd811461006657806327e235e31461008b57806370a08231146100aa578063a9059cbb146100c9575b600080fd5b341561007157600080fd5b6100796100ff565b60405190815260200160405180910390f35b341561009657600080fd5b610079600160a060020a0360043516610105565b34156100b557600080fd5b610079600160a060020a0360043516610117565b34156100d457600080fd5b6100eb600160a060020a0360043516602435610132565b604051901515815260200160405180910390f35b60005481565b60016020526000908152604090205481565b600160a060020a031660009081526001602052604090205490565b6000600160a060020a038316151561014957600080fd5b600160a060020a03331660009081526001602052604090205482111561016e57600080fd5b600160a060020a033316600090815260016020526040902054610197908363ffffffff61022d16565b600160a060020a0333811660009081526001602052604080822093909355908516815220546101cc908363ffffffff61023f16565b600160a060020a0380851660008181526001602052604090819020939093559133909116907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9085905190815260200160405180910390a350600192915050565b60008282111561023957fe5b50900390565b60008282018381101561024e57fe5b93925050505600a165627a7a723058203009a2682080b7122335e277a04de8b3d09fbd179da95e113b1220e25409b8fa0029`
//
//// DeployBasicToken deploys a new Ethereum contract, binding an instance of BasicToken to it.
//func DeployBasicToken(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BasicToken, error) {
//	parsed, err := abi.JSON(strings.NewReader(BasicTokenABI))
//	if err != nil {
//		return common.Address{}, nil, nil, err
//	}
//	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(BasicTokenBin), backend)
//	if err != nil {
//		return common.Address{}, nil, nil, err
//	}
//	return address, tx, &BasicToken{BasicTokenCaller: BasicTokenCaller{contract: contract}, BasicTokenTransactor: BasicTokenTransactor{contract: contract}, BasicTokenFilterer: BasicTokenFilterer{contract: contract}}, nil
//}
//
//// BasicToken is an auto generated Go binding around an Ethereum contract.
//type BasicToken struct {
//	BasicTokenCaller     // Read-only binding to the contract
//	BasicTokenTransactor // Write-only binding to the contract
//	BasicTokenFilterer   // Log filterer for contract events
//}
//
//// BasicTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
//type BasicTokenCaller struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// BasicTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
//type BasicTokenTransactor struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// BasicTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
//type BasicTokenFilterer struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// BasicTokenSession is an auto generated Go binding around an Ethereum contract,
//// with pre-set call and transact options.
//type BasicTokenSession struct {
//	Contract     *BasicToken       // Generic contract binding to set the session for
//	CallOpts     bind.CallOpts     // Call options to use throughout this session
//	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
//}
//
//// BasicTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
//// with pre-set call options.
//type BasicTokenCallerSession struct {
//	Contract *BasicTokenCaller // Generic contract caller binding to set the session for
//	CallOpts bind.CallOpts     // Call options to use throughout this session
//}
//
//// BasicTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
//// with pre-set transact options.
//type BasicTokenTransactorSession struct {
//	Contract     *BasicTokenTransactor // Generic contract transactor binding to set the session for
//	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
//}
//
//// BasicTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
//type BasicTokenRaw struct {
//	Contract *BasicToken // Generic contract binding to access the raw methods on
//}
//
//// BasicTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
//type BasicTokenCallerRaw struct {
//	Contract *BasicTokenCaller // Generic read-only contract binding to access the raw methods on
//}
//
//// BasicTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
//type BasicTokenTransactorRaw struct {
//	Contract *BasicTokenTransactor // Generic write-only contract binding to access the raw methods on
//}
//
//// NewBasicToken creates a new instance of BasicToken, bound to a specific deployed contract.
//func NewBasicToken(address common.Address, backend bind.ContractBackend) (*BasicToken, error) {
//	contract, err := bindBasicToken(address, backend, backend, backend)
//	if err != nil {
//		return nil, err
//	}
//	return &BasicToken{BasicTokenCaller: BasicTokenCaller{contract: contract}, BasicTokenTransactor: BasicTokenTransactor{contract: contract}, BasicTokenFilterer: BasicTokenFilterer{contract: contract}}, nil
//}
//
//// NewBasicTokenCaller creates a new read-only instance of BasicToken, bound to a specific deployed contract.
//func NewBasicTokenCaller(address common.Address, caller bind.ContractCaller) (*BasicTokenCaller, error) {
//	contract, err := bindBasicToken(address, caller, nil, nil)
//	if err != nil {
//		return nil, err
//	}
//	return &BasicTokenCaller{contract: contract}, nil
//}
//
//// NewBasicTokenTransactor creates a new write-only instance of BasicToken, bound to a specific deployed contract.
//func NewBasicTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*BasicTokenTransactor, error) {
//	contract, err := bindBasicToken(address, nil, transactor, nil)
//	if err != nil {
//		return nil, err
//	}
//	return &BasicTokenTransactor{contract: contract}, nil
//}
//
//// NewBasicTokenFilterer creates a new log filterer instance of BasicToken, bound to a specific deployed contract.
//func NewBasicTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*BasicTokenFilterer, error) {
//	contract, err := bindBasicToken(address, nil, nil, filterer)
//	if err != nil {
//		return nil, err
//	}
//	return &BasicTokenFilterer{contract: contract}, nil
//}
//
//// bindBasicToken binds a generic wrapper to an already deployed contract.
//func bindBasicToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
//	parsed, err := abi.JSON(strings.NewReader(BasicTokenABI))
//	if err != nil {
//		return nil, err
//	}
//	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
//}
//
//// Call invokes the (constant) contract method with params as input values and
//// sets the output to result. The result type might be a single field for simple
//// returns, a slice of interfaces for anonymous returns and a struct for named
//// returns.
//func (_BasicToken *BasicTokenRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
//	return _BasicToken.Contract.BasicTokenCaller.contract.Call(opts, result, method, params...)
//}
//
//// Transfer initiates a plain transaction to move funds to the contract, calling
//// its default method if one is available.
//func (_BasicToken *BasicTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _BasicToken.Contract.BasicTokenTransactor.contract.Transfer(opts)
//}
//
//// Transact invokes the (paid) contract method with params as input values.
//func (_BasicToken *BasicTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
//	return _BasicToken.Contract.BasicTokenTransactor.contract.Transact(opts, method, params...)
//}
//
//// Call invokes the (constant) contract method with params as input values and
//// sets the output to result. The result type might be a single field for simple
//// returns, a slice of interfaces for anonymous returns and a struct for named
//// returns.
//func (_BasicToken *BasicTokenCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
//	return _BasicToken.Contract.contract.Call(opts, result, method, params...)
//}
//
//// Transfer initiates a plain transaction to move funds to the contract, calling
//// its default method if one is available.
//func (_BasicToken *BasicTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _BasicToken.Contract.contract.Transfer(opts)
//}
//
//// Transact invokes the (paid) contract method with params as input values.
//func (_BasicToken *BasicTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
//	return _BasicToken.Contract.contract.Transact(opts, method, params...)
//}
//
//// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
////
//// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
//func (_BasicToken *BasicTokenCaller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
//	var (
//		ret0 = new(*big.Int)
//	)
//	out := ret0
//	err := _BasicToken.contract.Call(opts, out, "balanceOf", _owner)
//	return *ret0, err
//}
//
//// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
////
//// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
//func (_BasicToken *BasicTokenSession) BalanceOf(_owner common.Address) (*big.Int, error) {
//	return _BasicToken.Contract.BalanceOf(&_BasicToken.CallOpts, _owner)
//}
//
//// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
////
//// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
//func (_BasicToken *BasicTokenCallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
//	return _BasicToken.Contract.BalanceOf(&_BasicToken.CallOpts, _owner)
//}
//
//// Balances is a free data retrieval call binding the contract method 0x27e235e3.
////
//// Solidity: function balances( address) constant returns(uint256)
//func (_BasicToken *BasicTokenCaller) Balances(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
//	var (
//		ret0 = new(*big.Int)
//	)
//	out := ret0
//	err := _BasicToken.contract.Call(opts, out, "balances", arg0)
//	return *ret0, err
//}
//
//// Balances is a free data retrieval call binding the contract method 0x27e235e3.
////
//// Solidity: function balances( address) constant returns(uint256)
//func (_BasicToken *BasicTokenSession) Balances(arg0 common.Address) (*big.Int, error) {
//	return _BasicToken.Contract.Balances(&_BasicToken.CallOpts, arg0)
//}
//
//// Balances is a free data retrieval call binding the contract method 0x27e235e3.
////
//// Solidity: function balances( address) constant returns(uint256)
//func (_BasicToken *BasicTokenCallerSession) Balances(arg0 common.Address) (*big.Int, error) {
//	return _BasicToken.Contract.Balances(&_BasicToken.CallOpts, arg0)
//}
//
//// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
////
//// Solidity: function totalSupply() constant returns(uint256)
//func (_BasicToken *BasicTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
//	var (
//		ret0 = new(*big.Int)
//	)
//	out := ret0
//	err := _BasicToken.contract.Call(opts, out, "totalSupply")
//	return *ret0, err
//}
//
//// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
////
//// Solidity: function totalSupply() constant returns(uint256)
//func (_BasicToken *BasicTokenSession) TotalSupply() (*big.Int, error) {
//	return _BasicToken.Contract.TotalSupply(&_BasicToken.CallOpts)
//}
//
//// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
////
//// Solidity: function totalSupply() constant returns(uint256)
//func (_BasicToken *BasicTokenCallerSession) TotalSupply() (*big.Int, error) {
//	return _BasicToken.Contract.TotalSupply(&_BasicToken.CallOpts)
//}
//
//// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
////
//// Solidity: function transfer(_to address, _value uint256) returns(bool)
//func (_BasicToken *BasicTokenTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _BasicToken.contract.Transact(opts, "transfer", _to, _value)
//}
//
//// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
////
//// Solidity: function transfer(_to address, _value uint256) returns(bool)
//func (_BasicToken *BasicTokenSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _BasicToken.Contract.Transfer(&_BasicToken.TransactOpts, _to, _value)
//}
//
//// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
////
//// Solidity: function transfer(_to address, _value uint256) returns(bool)
//func (_BasicToken *BasicTokenTransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _BasicToken.Contract.Transfer(&_BasicToken.TransactOpts, _to, _value)
//}
//
//// BasicTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the BasicToken contract.
//type BasicTokenTransferIterator struct {
//	Event *BasicTokenTransfer // Event containing the contract specifics and raw log
//
//	contract *bind.BoundContract // Generic contract to use for unpacking event data
//	event    string              // Event name to use for unpacking event data
//
//	logs chan types.Log        // Log channel receiving the found contract events
//	sub  ethereum.Subscription // Subscription for errors, completion and termination
//	done bool                  // Whether the subscription completed delivering logs
//	fail error                 // Occurred error to stop iteration
//}
//
//// Next advances the iterator to the subsequent event, returning whether there
//// are any more events found. In case of a retrieval or parsing error, false is
//// returned and Error() can be queried for the exact failure.
//func (it *BasicTokenTransferIterator) Next() bool {
//	// If the iterator failed, stop iterating
//	if it.fail != nil {
//		return false
//	}
//	// If the iterator completed, deliver directly whatever's available
//	if it.done {
//		select {
//		case log := <-it.logs:
//			it.Event = new(BasicTokenTransfer)
//			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//				it.fail = err
//				return false
//			}
//			it.Event.Raw = log
//			return true
//
//		default:
//			return false
//		}
//	}
//	// Iterator still in progress, wait for either a data or an error event
//	select {
//	case log := <-it.logs:
//		it.Event = new(BasicTokenTransfer)
//		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//			it.fail = err
//			return false
//		}
//		it.Event.Raw = log
//		return true
//
//	case err := <-it.sub.Err():
//		it.done = true
//		it.fail = err
//		return it.Next()
//	}
//}
//
//// Error returns any retrieval or parsing error occurred during filtering.
//func (it *BasicTokenTransferIterator) Error() error {
//	return it.fail
//}
//
//// Close terminates the iteration process, releasing any pending underlying
//// resources.
//func (it *BasicTokenTransferIterator) Close() error {
//	it.sub.Unsubscribe()
//	return nil
//}
//
//// BasicTokenTransfer represents a Transfer event raised by the BasicToken contract.
//type BasicTokenTransfer struct {
//	From  common.Address
//	To    common.Address
//	Value *big.Int
//	Raw   types.Log // Blockchain specific contextual infos
//}
//
//// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
////
//// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
//func (_BasicToken *BasicTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BasicTokenTransferIterator, error) {
//
//	var fromRule []interface{}
//	for _, fromItem := range from {
//		fromRule = append(fromRule, fromItem)
//	}
//	var toRule []interface{}
//	for _, toItem := range to {
//		toRule = append(toRule, toItem)
//	}
//
//	logs, sub, err := _BasicToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
//	if err != nil {
//		return nil, err
//	}
//	return &BasicTokenTransferIterator{contract: _BasicToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
//}
//
//// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
////
//// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
//func (_BasicToken *BasicTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *BasicTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {
//
//	var fromRule []interface{}
//	for _, fromItem := range from {
//		fromRule = append(fromRule, fromItem)
//	}
//	var toRule []interface{}
//	for _, toItem := range to {
//		toRule = append(toRule, toItem)
//	}
//
//	logs, sub, err := _BasicToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
//	if err != nil {
//		return nil, err
//	}
//	return event.NewSubscription(func(quit <-chan struct{}) error {
//		defer sub.Unsubscribe()
//		for {
//			select {
//			case log := <-logs:
//				// New log arrived, parse the event and forward to the user
//				event := new(BasicTokenTransfer)
//				if err := _BasicToken.contract.UnpackLog(event, "Transfer", log); err != nil {
//					return err
//				}
//				event.Raw = log
//
//				select {
//				case sink <- event:
//				case err := <-sub.Err():
//					return err
//				case <-quit:
//					return nil
//				}
//			case err := <-sub.Err():
//				return err
//			case <-quit:
//				return nil
//			}
//		}
//	}), nil
//}
//
//// BurnableTokenABI is the input ABI used to generate the binding from.
//const BurnableTokenABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseApproval\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseApproval\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"burner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Burn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]"
//
//// BurnableTokenBin is the compiled bytecode used for deploying new contracts.
//const BurnableTokenBin = `0x6060604052341561000f57600080fd5b61080e8061001e6000396000f3006060604052600436106100a35763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663095ea7b381146100a857806318160ddd146100de57806323b872dd1461010357806327e235e31461012b57806342966c681461014a578063661884631461016257806370a0823114610184578063a9059cbb146101a3578063d73dd623146101c5578063dd62ed3e146101e7575b600080fd5b34156100b357600080fd5b6100ca600160a060020a036004351660243561020c565b604051901515815260200160405180910390f35b34156100e957600080fd5b6100f1610278565b60405190815260200160405180910390f35b341561010e57600080fd5b6100ca600160a060020a036004358116906024351660443561027e565b341561013657600080fd5b6100f1600160a060020a0360043516610400565b341561015557600080fd5b610160600435610412565b005b341561016d57600080fd5b6100ca600160a060020a03600435166024356104db565b341561018f57600080fd5b6100f1600160a060020a03600435166105d5565b34156101ae57600080fd5b6100ca600160a060020a03600435166024356105f0565b34156101d057600080fd5b6100ca600160a060020a03600435166024356106eb565b34156101f257600080fd5b6100f1600160a060020a036004358116906024351661078f565b600160a060020a03338116600081815260026020908152604080832094871680845294909152808220859055909291907f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259085905190815260200160405180910390a350600192915050565b60005481565b6000600160a060020a038316151561029557600080fd5b600160a060020a0384166000908152600160205260409020548211156102ba57600080fd5b600160a060020a03808516600090815260026020908152604080832033909416835292905220548211156102ed57600080fd5b600160a060020a038416600090815260016020526040902054610316908363ffffffff6107ba16565b600160a060020a03808616600090815260016020526040808220939093559085168152205461034b908363ffffffff6107cc16565b600160a060020a03808516600090815260016020908152604080832094909455878316825260028152838220339093168252919091522054610393908363ffffffff6107ba16565b600160a060020a03808616600081815260026020908152604080832033861684529091529081902093909355908516917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9085905190815260200160405180910390a35060019392505050565b60016020526000908152604090205481565b600080821161042057600080fd5b600160a060020a03331660009081526001602052604090205482111561044557600080fd5b5033600160a060020a03811660009081526001602052604090205461046a90836107ba565b600160a060020a03821660009081526001602052604081209190915554610497908363ffffffff6107ba16565b600055600160a060020a0381167fcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca58360405190815260200160405180910390a25050565b600160a060020a0333811660009081526002602090815260408083209386168352929052908120548083111561053857600160a060020a03338116600090815260026020908152604080832093881683529290529081205561056f565b610548818463ffffffff6107ba16565b600160a060020a033381166000908152600260209081526040808320938916835292905220555b600160a060020a0333811660008181526002602090815260408083209489168084529490915290819020547f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925915190815260200160405180910390a35060019392505050565b600160a060020a031660009081526001602052604090205490565b6000600160a060020a038316151561060757600080fd5b600160a060020a03331660009081526001602052604090205482111561062c57600080fd5b600160a060020a033316600090815260016020526040902054610655908363ffffffff6107ba16565b600160a060020a03338116600090815260016020526040808220939093559085168152205461068a908363ffffffff6107cc16565b600160a060020a0380851660008181526001602052604090819020939093559133909116907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9085905190815260200160405180910390a350600192915050565b600160a060020a033381166000908152600260209081526040808320938616835292905290812054610723908363ffffffff6107cc16565b600160a060020a0333811660008181526002602090815260408083209489168084529490915290819020849055919290917f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92591905190815260200160405180910390a350600192915050565b600160a060020a03918216600090815260026020908152604080832093909416825291909152205490565b6000828211156107c657fe5b50900390565b6000828201838110156107db57fe5b93925050505600a165627a7a72305820919c8940b8573ebf62bc3e59f8191350c89d4a1e8b214abb413482a221da38540029`
//
//// DeployBurnableToken deploys a new Ethereum contract, binding an instance of BurnableToken to it.
//func DeployBurnableToken(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BurnableToken, error) {
//	parsed, err := abi.JSON(strings.NewReader(BurnableTokenABI))
//	if err != nil {
//		return common.Address{}, nil, nil, err
//	}
//	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(BurnableTokenBin), backend)
//	if err != nil {
//		return common.Address{}, nil, nil, err
//	}
//	return address, tx, &BurnableToken{BurnableTokenCaller: BurnableTokenCaller{contract: contract}, BurnableTokenTransactor: BurnableTokenTransactor{contract: contract}, BurnableTokenFilterer: BurnableTokenFilterer{contract: contract}}, nil
//}
//
//// BurnableToken is an auto generated Go binding around an Ethereum contract.
//type BurnableToken struct {
//	BurnableTokenCaller     // Read-only binding to the contract
//	BurnableTokenTransactor // Write-only binding to the contract
//	BurnableTokenFilterer   // Log filterer for contract events
//}
//
//// BurnableTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
//type BurnableTokenCaller struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// BurnableTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
//type BurnableTokenTransactor struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// BurnableTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
//type BurnableTokenFilterer struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// BurnableTokenSession is an auto generated Go binding around an Ethereum contract,
//// with pre-set call and transact options.
//type BurnableTokenSession struct {
//	Contract     *BurnableToken    // Generic contract binding to set the session for
//	CallOpts     bind.CallOpts     // Call options to use throughout this session
//	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
//}
//
//// BurnableTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
//// with pre-set call options.
//type BurnableTokenCallerSession struct {
//	Contract *BurnableTokenCaller // Generic contract caller binding to set the session for
//	CallOpts bind.CallOpts        // Call options to use throughout this session
//}
//
//// BurnableTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
//// with pre-set transact options.
//type BurnableTokenTransactorSession struct {
//	Contract     *BurnableTokenTransactor // Generic contract transactor binding to set the session for
//	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
//}
//
//// BurnableTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
//type BurnableTokenRaw struct {
//	Contract *BurnableToken // Generic contract binding to access the raw methods on
//}
//
//// BurnableTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
//type BurnableTokenCallerRaw struct {
//	Contract *BurnableTokenCaller // Generic read-only contract binding to access the raw methods on
//}
//
//// BurnableTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
//type BurnableTokenTransactorRaw struct {
//	Contract *BurnableTokenTransactor // Generic write-only contract binding to access the raw methods on
//}
//
//// NewBurnableToken creates a new instance of BurnableToken, bound to a specific deployed contract.
//func NewBurnableToken(address common.Address, backend bind.ContractBackend) (*BurnableToken, error) {
//	contract, err := bindBurnableToken(address, backend, backend, backend)
//	if err != nil {
//		return nil, err
//	}
//	return &BurnableToken{BurnableTokenCaller: BurnableTokenCaller{contract: contract}, BurnableTokenTransactor: BurnableTokenTransactor{contract: contract}, BurnableTokenFilterer: BurnableTokenFilterer{contract: contract}}, nil
//}
//
//// NewBurnableTokenCaller creates a new read-only instance of BurnableToken, bound to a specific deployed contract.
//func NewBurnableTokenCaller(address common.Address, caller bind.ContractCaller) (*BurnableTokenCaller, error) {
//	contract, err := bindBurnableToken(address, caller, nil, nil)
//	if err != nil {
//		return nil, err
//	}
//	return &BurnableTokenCaller{contract: contract}, nil
//}
//
//// NewBurnableTokenTransactor creates a new write-only instance of BurnableToken, bound to a specific deployed contract.
//func NewBurnableTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*BurnableTokenTransactor, error) {
//	contract, err := bindBurnableToken(address, nil, transactor, nil)
//	if err != nil {
//		return nil, err
//	}
//	return &BurnableTokenTransactor{contract: contract}, nil
//}
//
//// NewBurnableTokenFilterer creates a new log filterer instance of BurnableToken, bound to a specific deployed contract.
//func NewBurnableTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*BurnableTokenFilterer, error) {
//	contract, err := bindBurnableToken(address, nil, nil, filterer)
//	if err != nil {
//		return nil, err
//	}
//	return &BurnableTokenFilterer{contract: contract}, nil
//}
//
//// bindBurnableToken binds a generic wrapper to an already deployed contract.
//func bindBurnableToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
//	parsed, err := abi.JSON(strings.NewReader(BurnableTokenABI))
//	if err != nil {
//		return nil, err
//	}
//	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
//}
//
//// Call invokes the (constant) contract method with params as input values and
//// sets the output to result. The result type might be a single field for simple
//// returns, a slice of interfaces for anonymous returns and a struct for named
//// returns.
//func (_BurnableToken *BurnableTokenRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
//	return _BurnableToken.Contract.BurnableTokenCaller.contract.Call(opts, result, method, params...)
//}
//
//// Transfer initiates a plain transaction to move funds to the contract, calling
//// its default method if one is available.
//func (_BurnableToken *BurnableTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _BurnableToken.Contract.BurnableTokenTransactor.contract.Transfer(opts)
//}
//
//// Transact invokes the (paid) contract method with params as input values.
//func (_BurnableToken *BurnableTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
//	return _BurnableToken.Contract.BurnableTokenTransactor.contract.Transact(opts, method, params...)
//}
//
//// Call invokes the (constant) contract method with params as input values and
//// sets the output to result. The result type might be a single field for simple
//// returns, a slice of interfaces for anonymous returns and a struct for named
//// returns.
//func (_BurnableToken *BurnableTokenCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
//	return _BurnableToken.Contract.contract.Call(opts, result, method, params...)
//}
//
//// Transfer initiates a plain transaction to move funds to the contract, calling
//// its default method if one is available.
//func (_BurnableToken *BurnableTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _BurnableToken.Contract.contract.Transfer(opts)
//}
//
//// Transact invokes the (paid) contract method with params as input values.
//func (_BurnableToken *BurnableTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
//	return _BurnableToken.Contract.contract.Transact(opts, method, params...)
//}
//
//// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
////
//// Solidity: function allowance(_owner address, _spender address) constant returns(uint256)
//func (_BurnableToken *BurnableTokenCaller) Allowance(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (*big.Int, error) {
//	var (
//		ret0 = new(*big.Int)
//	)
//	out := ret0
//	err := _BurnableToken.contract.Call(opts, out, "allowance", _owner, _spender)
//	return *ret0, err
//}
//
//// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
////
//// Solidity: function allowance(_owner address, _spender address) constant returns(uint256)
//func (_BurnableToken *BurnableTokenSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
//	return _BurnableToken.Contract.Allowance(&_BurnableToken.CallOpts, _owner, _spender)
//}
//
//// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
////
//// Solidity: function allowance(_owner address, _spender address) constant returns(uint256)
//func (_BurnableToken *BurnableTokenCallerSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
//	return _BurnableToken.Contract.Allowance(&_BurnableToken.CallOpts, _owner, _spender)
//}
//
//// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
////
//// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
//func (_BurnableToken *BurnableTokenCaller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
//	var (
//		ret0 = new(*big.Int)
//	)
//	out := ret0
//	err := _BurnableToken.contract.Call(opts, out, "balanceOf", _owner)
//	return *ret0, err
//}
//
//// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
////
//// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
//func (_BurnableToken *BurnableTokenSession) BalanceOf(_owner common.Address) (*big.Int, error) {
//	return _BurnableToken.Contract.BalanceOf(&_BurnableToken.CallOpts, _owner)
//}
//
//// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
////
//// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
//func (_BurnableToken *BurnableTokenCallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
//	return _BurnableToken.Contract.BalanceOf(&_BurnableToken.CallOpts, _owner)
//}
//
//// Balances is a free data retrieval call binding the contract method 0x27e235e3.
////
//// Solidity: function balances( address) constant returns(uint256)
//func (_BurnableToken *BurnableTokenCaller) Balances(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
//	var (
//		ret0 = new(*big.Int)
//	)
//	out := ret0
//	err := _BurnableToken.contract.Call(opts, out, "balances", arg0)
//	return *ret0, err
//}
//
//// Balances is a free data retrieval call binding the contract method 0x27e235e3.
////
//// Solidity: function balances( address) constant returns(uint256)
//func (_BurnableToken *BurnableTokenSession) Balances(arg0 common.Address) (*big.Int, error) {
//	return _BurnableToken.Contract.Balances(&_BurnableToken.CallOpts, arg0)
//}
//
//// Balances is a free data retrieval call binding the contract method 0x27e235e3.
////
//// Solidity: function balances( address) constant returns(uint256)
//func (_BurnableToken *BurnableTokenCallerSession) Balances(arg0 common.Address) (*big.Int, error) {
//	return _BurnableToken.Contract.Balances(&_BurnableToken.CallOpts, arg0)
//}
//
//// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
////
//// Solidity: function totalSupply() constant returns(uint256)
//func (_BurnableToken *BurnableTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
//	var (
//		ret0 = new(*big.Int)
//	)
//	out := ret0
//	err := _BurnableToken.contract.Call(opts, out, "totalSupply")
//	return *ret0, err
//}
//
//// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
////
//// Solidity: function totalSupply() constant returns(uint256)
//func (_BurnableToken *BurnableTokenSession) TotalSupply() (*big.Int, error) {
//	return _BurnableToken.Contract.TotalSupply(&_BurnableToken.CallOpts)
//}
//
//// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
////
//// Solidity: function totalSupply() constant returns(uint256)
//func (_BurnableToken *BurnableTokenCallerSession) TotalSupply() (*big.Int, error) {
//	return _BurnableToken.Contract.TotalSupply(&_BurnableToken.CallOpts)
//}
//
//// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
////
//// Solidity: function approve(_spender address, _value uint256) returns(bool)
//func (_BurnableToken *BurnableTokenTransactor) Approve(opts *bind.TransactOpts, _spender common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _BurnableToken.contract.Transact(opts, "approve", _spender, _value)
//}
//
//// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
////
//// Solidity: function approve(_spender address, _value uint256) returns(bool)
//func (_BurnableToken *BurnableTokenSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _BurnableToken.Contract.Approve(&_BurnableToken.TransactOpts, _spender, _value)
//}
//
//// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
////
//// Solidity: function approve(_spender address, _value uint256) returns(bool)
//func (_BurnableToken *BurnableTokenTransactorSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _BurnableToken.Contract.Approve(&_BurnableToken.TransactOpts, _spender, _value)
//}
//
//// Burn is a paid mutator transaction binding the contract method 0x42966c68.
////
//// Solidity: function burn(_value uint256) returns()
//func (_BurnableToken *BurnableTokenTransactor) Burn(opts *bind.TransactOpts, _value *big.Int) (*types.Transaction, error) {
//	return _BurnableToken.contract.Transact(opts, "burn", _value)
//}
//
//// Burn is a paid mutator transaction binding the contract method 0x42966c68.
////
//// Solidity: function burn(_value uint256) returns()
//func (_BurnableToken *BurnableTokenSession) Burn(_value *big.Int) (*types.Transaction, error) {
//	return _BurnableToken.Contract.Burn(&_BurnableToken.TransactOpts, _value)
//}
//
//// Burn is a paid mutator transaction binding the contract method 0x42966c68.
////
//// Solidity: function burn(_value uint256) returns()
//func (_BurnableToken *BurnableTokenTransactorSession) Burn(_value *big.Int) (*types.Transaction, error) {
//	return _BurnableToken.Contract.Burn(&_BurnableToken.TransactOpts, _value)
//}
//
//// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
////
//// Solidity: function decreaseApproval(_spender address, _subtractedValue uint256) returns(bool)
//func (_BurnableToken *BurnableTokenTransactor) DecreaseApproval(opts *bind.TransactOpts, _spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
//	return _BurnableToken.contract.Transact(opts, "decreaseApproval", _spender, _subtractedValue)
//}
//
//// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
////
//// Solidity: function decreaseApproval(_spender address, _subtractedValue uint256) returns(bool)
//func (_BurnableToken *BurnableTokenSession) DecreaseApproval(_spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
//	return _BurnableToken.Contract.DecreaseApproval(&_BurnableToken.TransactOpts, _spender, _subtractedValue)
//}
//
//// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
////
//// Solidity: function decreaseApproval(_spender address, _subtractedValue uint256) returns(bool)
//func (_BurnableToken *BurnableTokenTransactorSession) DecreaseApproval(_spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
//	return _BurnableToken.Contract.DecreaseApproval(&_BurnableToken.TransactOpts, _spender, _subtractedValue)
//}
//
//// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
////
//// Solidity: function increaseApproval(_spender address, _addedValue uint256) returns(bool)
//func (_BurnableToken *BurnableTokenTransactor) IncreaseApproval(opts *bind.TransactOpts, _spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
//	return _BurnableToken.contract.Transact(opts, "increaseApproval", _spender, _addedValue)
//}
//
//// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
////
//// Solidity: function increaseApproval(_spender address, _addedValue uint256) returns(bool)
//func (_BurnableToken *BurnableTokenSession) IncreaseApproval(_spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
//	return _BurnableToken.Contract.IncreaseApproval(&_BurnableToken.TransactOpts, _spender, _addedValue)
//}
//
//// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
////
//// Solidity: function increaseApproval(_spender address, _addedValue uint256) returns(bool)
//func (_BurnableToken *BurnableTokenTransactorSession) IncreaseApproval(_spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
//	return _BurnableToken.Contract.IncreaseApproval(&_BurnableToken.TransactOpts, _spender, _addedValue)
//}
//
//// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
////
//// Solidity: function transfer(_to address, _value uint256) returns(bool)
//func (_BurnableToken *BurnableTokenTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _BurnableToken.contract.Transact(opts, "transfer", _to, _value)
//}
//
//// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
////
//// Solidity: function transfer(_to address, _value uint256) returns(bool)
//func (_BurnableToken *BurnableTokenSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _BurnableToken.Contract.Transfer(&_BurnableToken.TransactOpts, _to, _value)
//}
//
//// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
////
//// Solidity: function transfer(_to address, _value uint256) returns(bool)
//func (_BurnableToken *BurnableTokenTransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _BurnableToken.Contract.Transfer(&_BurnableToken.TransactOpts, _to, _value)
//}
//
//// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
////
//// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
//func (_BurnableToken *BurnableTokenTransactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _BurnableToken.contract.Transact(opts, "transferFrom", _from, _to, _value)
//}
//
//// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
////
//// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
//func (_BurnableToken *BurnableTokenSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _BurnableToken.Contract.TransferFrom(&_BurnableToken.TransactOpts, _from, _to, _value)
//}
//
//// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
////
//// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
//func (_BurnableToken *BurnableTokenTransactorSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _BurnableToken.Contract.TransferFrom(&_BurnableToken.TransactOpts, _from, _to, _value)
//}
//
//// BurnableTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the BurnableToken contract.
//type BurnableTokenApprovalIterator struct {
//	Event *BurnableTokenApproval // Event containing the contract specifics and raw log
//
//	contract *bind.BoundContract // Generic contract to use for unpacking event data
//	event    string              // Event name to use for unpacking event data
//
//	logs chan types.Log        // Log channel receiving the found contract events
//	sub  ethereum.Subscription // Subscription for errors, completion and termination
//	done bool                  // Whether the subscription completed delivering logs
//	fail error                 // Occurred error to stop iteration
//}
//
//// Next advances the iterator to the subsequent event, returning whether there
//// are any more events found. In case of a retrieval or parsing error, false is
//// returned and Error() can be queried for the exact failure.
//func (it *BurnableTokenApprovalIterator) Next() bool {
//	// If the iterator failed, stop iterating
//	if it.fail != nil {
//		return false
//	}
//	// If the iterator completed, deliver directly whatever's available
//	if it.done {
//		select {
//		case log := <-it.logs:
//			it.Event = new(BurnableTokenApproval)
//			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//				it.fail = err
//				return false
//			}
//			it.Event.Raw = log
//			return true
//
//		default:
//			return false
//		}
//	}
//	// Iterator still in progress, wait for either a data or an error event
//	select {
//	case log := <-it.logs:
//		it.Event = new(BurnableTokenApproval)
//		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//			it.fail = err
//			return false
//		}
//		it.Event.Raw = log
//		return true
//
//	case err := <-it.sub.Err():
//		it.done = true
//		it.fail = err
//		return it.Next()
//	}
//}
//
//// Error returns any retrieval or parsing error occurred during filtering.
//func (it *BurnableTokenApprovalIterator) Error() error {
//	return it.fail
//}
//
//// Close terminates the iteration process, releasing any pending underlying
//// resources.
//func (it *BurnableTokenApprovalIterator) Close() error {
//	it.sub.Unsubscribe()
//	return nil
//}
//
//// BurnableTokenApproval represents a Approval event raised by the BurnableToken contract.
//type BurnableTokenApproval struct {
//	Owner   common.Address
//	Spender common.Address
//	Value   *big.Int
//	Raw     types.Log // Blockchain specific contextual infos
//}
//
//// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
////
//// Solidity: event Approval(owner indexed address, spender indexed address, value uint256)
//func (_BurnableToken *BurnableTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*BurnableTokenApprovalIterator, error) {
//
//	var ownerRule []interface{}
//	for _, ownerItem := range owner {
//		ownerRule = append(ownerRule, ownerItem)
//	}
//	var spenderRule []interface{}
//	for _, spenderItem := range spender {
//		spenderRule = append(spenderRule, spenderItem)
//	}
//
//	logs, sub, err := _BurnableToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
//	if err != nil {
//		return nil, err
//	}
//	return &BurnableTokenApprovalIterator{contract: _BurnableToken.contract, event: "Approval", logs: logs, sub: sub}, nil
//}
//
//// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
////
//// Solidity: event Approval(owner indexed address, spender indexed address, value uint256)
//func (_BurnableToken *BurnableTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *BurnableTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {
//
//	var ownerRule []interface{}
//	for _, ownerItem := range owner {
//		ownerRule = append(ownerRule, ownerItem)
//	}
//	var spenderRule []interface{}
//	for _, spenderItem := range spender {
//		spenderRule = append(spenderRule, spenderItem)
//	}
//
//	logs, sub, err := _BurnableToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
//	if err != nil {
//		return nil, err
//	}
//	return event.NewSubscription(func(quit <-chan struct{}) error {
//		defer sub.Unsubscribe()
//		for {
//			select {
//			case log := <-logs:
//				// New log arrived, parse the event and forward to the user
//				event := new(BurnableTokenApproval)
//				if err := _BurnableToken.contract.UnpackLog(event, "Approval", log); err != nil {
//					return err
//				}
//				event.Raw = log
//
//				select {
//				case sink <- event:
//				case err := <-sub.Err():
//					return err
//				case <-quit:
//					return nil
//				}
//			case err := <-sub.Err():
//				return err
//			case <-quit:
//				return nil
//			}
//		}
//	}), nil
//}
//
//// BurnableTokenBurnIterator is returned from FilterBurn and is used to iterate over the raw logs and unpacked data for Burn events raised by the BurnableToken contract.
//type BurnableTokenBurnIterator struct {
//	Event *BurnableTokenBurn // Event containing the contract specifics and raw log
//
//	contract *bind.BoundContract // Generic contract to use for unpacking event data
//	event    string              // Event name to use for unpacking event data
//
//	logs chan types.Log        // Log channel receiving the found contract events
//	sub  ethereum.Subscription // Subscription for errors, completion and termination
//	done bool                  // Whether the subscription completed delivering logs
//	fail error                 // Occurred error to stop iteration
//}
//
//// Next advances the iterator to the subsequent event, returning whether there
//// are any more events found. In case of a retrieval or parsing error, false is
//// returned and Error() can be queried for the exact failure.
//func (it *BurnableTokenBurnIterator) Next() bool {
//	// If the iterator failed, stop iterating
//	if it.fail != nil {
//		return false
//	}
//	// If the iterator completed, deliver directly whatever's available
//	if it.done {
//		select {
//		case log := <-it.logs:
//			it.Event = new(BurnableTokenBurn)
//			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//				it.fail = err
//				return false
//			}
//			it.Event.Raw = log
//			return true
//
//		default:
//			return false
//		}
//	}
//	// Iterator still in progress, wait for either a data or an error event
//	select {
//	case log := <-it.logs:
//		it.Event = new(BurnableTokenBurn)
//		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//			it.fail = err
//			return false
//		}
//		it.Event.Raw = log
//		return true
//
//	case err := <-it.sub.Err():
//		it.done = true
//		it.fail = err
//		return it.Next()
//	}
//}
//
//// Error returns any retrieval or parsing error occurred during filtering.
//func (it *BurnableTokenBurnIterator) Error() error {
//	return it.fail
//}
//
//// Close terminates the iteration process, releasing any pending underlying
//// resources.
//func (it *BurnableTokenBurnIterator) Close() error {
//	it.sub.Unsubscribe()
//	return nil
//}
//
//// BurnableTokenBurn represents a Burn event raised by the BurnableToken contract.
//type BurnableTokenBurn struct {
//	Burner common.Address
//	Value  *big.Int
//	Raw    types.Log // Blockchain specific contextual infos
//}
//
//// FilterBurn is a free log retrieval operation binding the contract event 0xcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5.
////
//// Solidity: event Burn(burner indexed address, value uint256)
//func (_BurnableToken *BurnableTokenFilterer) FilterBurn(opts *bind.FilterOpts, burner []common.Address) (*BurnableTokenBurnIterator, error) {
//
//	var burnerRule []interface{}
//	for _, burnerItem := range burner {
//		burnerRule = append(burnerRule, burnerItem)
//	}
//
//	logs, sub, err := _BurnableToken.contract.FilterLogs(opts, "Burn", burnerRule)
//	if err != nil {
//		return nil, err
//	}
//	return &BurnableTokenBurnIterator{contract: _BurnableToken.contract, event: "Burn", logs: logs, sub: sub}, nil
//}
//
//// WatchBurn is a free log subscription operation binding the contract event 0xcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5.
////
//// Solidity: event Burn(burner indexed address, value uint256)
//func (_BurnableToken *BurnableTokenFilterer) WatchBurn(opts *bind.WatchOpts, sink chan<- *BurnableTokenBurn, burner []common.Address) (event.Subscription, error) {
//
//	var burnerRule []interface{}
//	for _, burnerItem := range burner {
//		burnerRule = append(burnerRule, burnerItem)
//	}
//
//	logs, sub, err := _BurnableToken.contract.WatchLogs(opts, "Burn", burnerRule)
//	if err != nil {
//		return nil, err
//	}
//	return event.NewSubscription(func(quit <-chan struct{}) error {
//		defer sub.Unsubscribe()
//		for {
//			select {
//			case log := <-logs:
//				// New log arrived, parse the event and forward to the user
//				event := new(BurnableTokenBurn)
//				if err := _BurnableToken.contract.UnpackLog(event, "Burn", log); err != nil {
//					return err
//				}
//				event.Raw = log
//
//				select {
//				case sink <- event:
//				case err := <-sub.Err():
//					return err
//				case <-quit:
//					return nil
//				}
//			case err := <-sub.Err():
//				return err
//			case <-quit:
//				return nil
//			}
//		}
//	}), nil
//}
//
//// BurnableTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the BurnableToken contract.
//type BurnableTokenTransferIterator struct {
//	Event *BurnableTokenTransfer // Event containing the contract specifics and raw log
//
//	contract *bind.BoundContract // Generic contract to use for unpacking event data
//	event    string              // Event name to use for unpacking event data
//
//	logs chan types.Log        // Log channel receiving the found contract events
//	sub  ethereum.Subscription // Subscription for errors, completion and termination
//	done bool                  // Whether the subscription completed delivering logs
//	fail error                 // Occurred error to stop iteration
//}
//
//// Next advances the iterator to the subsequent event, returning whether there
//// are any more events found. In case of a retrieval or parsing error, false is
//// returned and Error() can be queried for the exact failure.
//func (it *BurnableTokenTransferIterator) Next() bool {
//	// If the iterator failed, stop iterating
//	if it.fail != nil {
//		return false
//	}
//	// If the iterator completed, deliver directly whatever's available
//	if it.done {
//		select {
//		case log := <-it.logs:
//			it.Event = new(BurnableTokenTransfer)
//			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//				it.fail = err
//				return false
//			}
//			it.Event.Raw = log
//			return true
//
//		default:
//			return false
//		}
//	}
//	// Iterator still in progress, wait for either a data or an error event
//	select {
//	case log := <-it.logs:
//		it.Event = new(BurnableTokenTransfer)
//		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//			it.fail = err
//			return false
//		}
//		it.Event.Raw = log
//		return true
//
//	case err := <-it.sub.Err():
//		it.done = true
//		it.fail = err
//		return it.Next()
//	}
//}
//
//// Error returns any retrieval or parsing error occurred during filtering.
//func (it *BurnableTokenTransferIterator) Error() error {
//	return it.fail
//}
//
//// Close terminates the iteration process, releasing any pending underlying
//// resources.
//func (it *BurnableTokenTransferIterator) Close() error {
//	it.sub.Unsubscribe()
//	return nil
//}
//
//// BurnableTokenTransfer represents a Transfer event raised by the BurnableToken contract.
//type BurnableTokenTransfer struct {
//	From  common.Address
//	To    common.Address
//	Value *big.Int
//	Raw   types.Log // Blockchain specific contextual infos
//}
//
//// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
////
//// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
//func (_BurnableToken *BurnableTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnableTokenTransferIterator, error) {
//
//	var fromRule []interface{}
//	for _, fromItem := range from {
//		fromRule = append(fromRule, fromItem)
//	}
//	var toRule []interface{}
//	for _, toItem := range to {
//		toRule = append(toRule, toItem)
//	}
//
//	logs, sub, err := _BurnableToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
//	if err != nil {
//		return nil, err
//	}
//	return &BurnableTokenTransferIterator{contract: _BurnableToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
//}
//
//// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
////
//// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
//func (_BurnableToken *BurnableTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *BurnableTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {
//
//	var fromRule []interface{}
//	for _, fromItem := range from {
//		fromRule = append(fromRule, fromItem)
//	}
//	var toRule []interface{}
//	for _, toItem := range to {
//		toRule = append(toRule, toItem)
//	}
//
//	logs, sub, err := _BurnableToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
//	if err != nil {
//		return nil, err
//	}
//	return event.NewSubscription(func(quit <-chan struct{}) error {
//		defer sub.Unsubscribe()
//		for {
//			select {
//			case log := <-logs:
//				// New log arrived, parse the event and forward to the user
//				event := new(BurnableTokenTransfer)
//				if err := _BurnableToken.contract.UnpackLog(event, "Transfer", log); err != nil {
//					return err
//				}
//				event.Raw = log
//
//				select {
//				case sink <- event:
//				case err := <-sub.Err():
//					return err
//				case <-quit:
//					return nil
//				}
//			case err := <-sub.Err():
//				return err
//			case <-quit:
//				return nil
//			}
//		}
//	}), nil
//}
//
//// DarkNodeRegistryABI is the input ABI used to generate the binding from.
//const DarkNodeRegistryABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"numDarkNodes\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numDarkNodesNextEpoch\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_darkNodeID\",\"type\":\"bytes20\"}],\"name\":\"isDeregistered\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_darkNodeID\",\"type\":\"bytes20\"}],\"name\":\"getPublicKey\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_darkNodeID\",\"type\":\"bytes20\"},{\"name\":\"_publicKey\",\"type\":\"bytes\"},{\"name\":\"_bond\",\"type\":\"uint256\"}],\"name\":\"register\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_darkNodeID\",\"type\":\"bytes20\"}],\"name\":\"isRegistered\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minimumEpochInterval\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_darkNodeID\",\"type\":\"bytes20\"}],\"name\":\"refund\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_darkNodeID\",\"type\":\"bytes20\"}],\"name\":\"getBond\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentEpoch\",\"outputs\":[{\"name\":\"blockhash\",\"type\":\"uint256\"},{\"name\":\"timestamp\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getDarkNodes\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes20[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"epoch\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minimumBond\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minimumDarkPoolSize\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_darkNodeID\",\"type\":\"bytes20\"}],\"name\":\"isUnregistered\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_darkNodeID\",\"type\":\"bytes20\"}],\"name\":\"deregister\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_darkNodeID\",\"type\":\"bytes20\"}],\"name\":\"getOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_token\",\"type\":\"address\"},{\"name\":\"_minimumBond\",\"type\":\"uint256\"},{\"name\":\"_minimumDarkPoolSize\",\"type\":\"uint256\"},{\"name\":\"_minimumEpochInterval\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_darkNodeID\",\"type\":\"bytes20\"},{\"indexed\":false,\"name\":\"_bond\",\"type\":\"uint256\"}],\"name\":\"DarkNodeRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_darkNodeID\",\"type\":\"bytes20\"}],\"name\":\"DarkNodeDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"OwnerRefunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"NewEpoch\",\"type\":\"event\"}]"
//
//// DarkNodeRegistryBin is the compiled bytecode used for deploying new contracts.
//const DarkNodeRegistryBin = `0x6060604052341561000f57600080fd5b6040516080806110978339810160405280805191906020018051919060200180519190602001805160008054600160a060020a031916600160a060020a038816179055600585905560068490556007819055915060409050805190810160405260001943014081524260208201526008815181556020820151600190910155505060006003819055600455505050610feb806100ac6000396000f3006060604052600436106100d75763ffffffff60e060020a600035041663060a2cfd81146100dc5780630620eb9214610101578063171f6ea81461011457806332ccd52f14610148578063375a8be3146101df5780634f5550fc1461024357806355cacda5146102635780635a8f9b811461027657806368f209eb1461029657806376671808146102b6578063879ae084146102e1578063900cf0cf14610347578063aa7517e11461035a578063b31575d51461036d578063d3841c2514610380578063e08b4c8a146103a0578063e487eb58146103c0575b600080fd5b34156100e757600080fd5b6100ef6103fc565b60405190815260200160405180910390f35b341561010c57600080fd5b6100ef610402565b341561011f57600080fd5b6101346001606060020a031960043516610408565b604051901515815260200160405180910390f35b341561015357600080fd5b6101686001606060020a031960043516610457565b60405160208082528190810183818151815260200191508051906020019080838360005b838110156101a457808201518382015260200161018c565b50505050905090810190601f1680156101d15780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34156101ea57600080fd5b610241600480356001606060020a0319169060446024803590810190830135806020601f82018190048102016040519081016040528181529291906020840183838082843750949650509335935061053392505050565b005b341561024e57600080fd5b6101346001606060020a03196004351661076b565b341561026e57600080fd5b6100ef6107cb565b341561028157600080fd5b6102416001606060020a0319600435166107d1565b34156102a157600080fd5b6100ef6001606060020a0319600435166109ce565b34156102c157600080fd5b6102c96109ee565b60405191825260208201526040908101905180910390f35b34156102ec57600080fd5b6102f46109f7565b60405160208082528190810183818151815260200191508051906020019060200280838360005b8381101561033357808201518382015260200161031b565b505050509050019250505060405180910390f35b341561035257600080fd5b610241610aae565b341561036557600080fd5b6100ef610b2a565b341561037857600080fd5b6100ef610b30565b341561038b57600080fd5b6101346001606060020a031960043516610b36565b34156103ab57600080fd5b6102416001606060020a031960043516610b56565b34156103cb57600080fd5b6103e06001606060020a031960043516610c17565b604051600160a060020a03909116815260200160405180910390f35b60035481565b60045481565b6001606060020a031981166000908152600160205260408120600301541580159061045157506009546001606060020a0319831660009081526001602052604090206003015411155b92915050565b61045f610f12565b60016000836bffffffffffffffffffffffff19166bffffffffffffffffffffffff191681526020019081526020016000206004018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156105275780601f106104fc57610100808354040283529160200191610527565b820191906000526020600020905b81548152906001019060200180831161050a57829003601f168201915b50505050509050919050565b8261053d81610b36565b151561054857600080fd5b60055482101561055757600080fd5b600054600160a060020a031663dd62ed3e333060405160e060020a63ffffffff8516028152600160a060020a03928316600482015291166024820152604401602060405180830381600087803b15156105af57600080fd5b5af115156105bc57600080fd5b505050604051805183111590506105d257600080fd5b600054600160a060020a03166323b872dd33308560405160e060020a63ffffffff8616028152600160a060020a0393841660048201529190921660248201526044810191909152606401602060405180830381600087803b151561063557600080fd5b5af1151561064257600080fd5b50505060405180519050151561065757600080fd5b60a06040519081016040908152600160a060020a033316825260208083018590526007546009540182840152600060608401819052608084018790526001606060020a0319881681526001909152208151815473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0391909116178155602082015181600101556040820151816002015560608201518160030155608082015181600401908051610708929160200190610f24565b50905050610717600285610c3c565b6004805460010190557fcde427a7822769e806a704726602a3a8a44458da1352d9b5a46bcfb95142ad3884836040516001606060020a0319909216825260208201526040908101905180910390a150505050565b6001606060020a03198116600090815260016020526040812060020154158015906107b457506009546001606060020a0319831660009081526001602052604090206002015411155b801561045157506107c482610408565b1592915050565b60075481565b6001606060020a03198116600090815260016020526040812054829033600160a060020a0390811691161461080557600080fd5b8261080f81610408565b151561081a57600080fd5b6001606060020a031984166000908152600160208190526040822001549350831161084157fe5b61084c600285610c67565b60a0604051908101604052806000600160a060020a03168152602001600081526020016000815260200160008152602001602060405190810160409081526000808352919092526001606060020a0319871681526001602052208151815473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0391909116178155602082015181600101556040820151816002015560608201518160030155608082015181600401908051610908929160200190610f24565b5050600054600160a060020a0316905063a9059cbb338560405160e060020a63ffffffff8516028152600160a060020a0390921660048301526024820152604401602060405180830381600087803b151561096257600080fd5b5af1151561096f57600080fd5b50505060405180519050151561098457600080fd5b7f8dce8f4eb4097fbdf948a109703513f52a7fabcc7e328ba7f29ad8165033e9533384604051600160a060020a03909216825260208201526040908101905180910390a150505050565b6001606060020a0319166000908152600160208190526040909120015490565b60085460095482565b6109ff610f12565b610a07610f12565b600080600354604051805910610a1a5750595b9080825280602002602001820160405250925060009150610a3b6002610d60565b90505b600354821015610aa657610a518161076b565b1515610a6957610a62600282610d85565b9050610a3e565b80838381518110610a7657fe5b6001606060020a0319909216602092830290910190910152610a99600282610d85565b6001909201919050610a3e565b509092915050565b600754600954600091014211610ac357600080fd5b50600019430140604080519081016040528181526007546009540160208201526008815181556020820151600190910155506004546003557fe358419ca0dd7928a310d787a606dfae5d869f5071249efa6107105e7afc40bc60405160405180910390a150565b60055481565b60065481565b6001606060020a0319166000908152600160205260409020600201541590565b80610b608161076b565b1515610b6b57600080fd5b6001606060020a03198216600090815260016020526040902054829033600160a060020a03908116911614610b9f57600080fd5b6007546009546001606060020a031985166000908152600160205260409081902091909201600390910155600480546000190190557fe3365373f442312a7d66aa1003b60fd4a0f9ed1c85fd8bb384240f8eb4fe916c908490516001606060020a0319909116815260200160405180910390a1505050565b6001606060020a031916600090815260016020526040902054600160a060020a031690565b610c468282610dcc565b15610c5057600080fd5b610c6382610c5d84610dec565b83610e14565b5050565b600080610c748484610dcc565b1515610c7f57600080fd5b6001606060020a031983161515610c9557610d5a565b50506001606060020a03198082166000818152602085905260408082208054600180830180546c01000000000000000000000000610100948590048102808b168952878920909401805492820282810473ffffffffffffffffffffffffffffffffffffffff1994851617909155998a168852958720805496840490940274ffffffffffffffffffffffffffffffffffffffff00199096169590951790925594909352805474ffffffffffffffffffffffffffffffffffffffffff191690558154169055905b50505050565b600080805260209190915260409020600101546c010000000000000000000000000290565b6000610d918383610dcc565b1515610d9c57600080fd5b506001606060020a031916600090815260209190915260409020600101546c010000000000000000000000000290565b6001606060020a0319166000908152602091909152604090205460ff1690565b60008080526020829052604090205461010090046c0100000000000000000000000002919050565b6000610e208483610dcc565b15610e2a57600080fd5b610e348484610dcc565b80610e4757506001606060020a03198316155b1515610e5257600080fd5b506001606060020a0319808316600090815260209490945260408085206001908101805485851680895284892080546c01000000000000000000000000998a900461010090810274ffffffffffffffffffffffffffffffffffffffff00199283161783558287018054958c028c810473ffffffffffffffffffffffffffffffffffffffff199788161790915586549b909a049a9094168a179094559690951688529287208054969093029516949094179055909252815460ff1916179055565b60206040519081016040526000815290565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10610f6557805160ff1916838001178555610f92565b82800160010185558215610f92579182015b82811115610f92578251825591602001919060010190610f77565b50610f9e929150610fa2565b5090565b610fbc91905b80821115610f9e5760008155600101610fa8565b905600a165627a7a72305820b482c337f6c6c9a753d538b1a4a450388ef42e25a444b74f0575414670a9c1bc0029`
//
//// DeployDarkNodeRegistry deploys a new Ethereum contract, binding an instance of DarkNodeRegistry to it.
//func DeployDarkNodeRegistry(auth *bind.TransactOpts, backend bind.ContractBackend, _token common.Address, _minimumBond *big.Int, _minimumDarkPoolSize *big.Int, _minimumEpochInterval *big.Int) (common.Address, *types.Transaction, *DarkNodeRegistry, error) {
//	parsed, err := abi.JSON(strings.NewReader(DarkNodeRegistryABI))
//	if err != nil {
//		return common.Address{}, nil, nil, err
//	}
//	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(DarkNodeRegistryBin), backend, _token, _minimumBond, _minimumDarkPoolSize, _minimumEpochInterval)
//	if err != nil {
//		return common.Address{}, nil, nil, err
//	}
//	return address, tx, &DarkNodeRegistry{DarkNodeRegistryCaller: DarkNodeRegistryCaller{contract: contract}, DarkNodeRegistryTransactor: DarkNodeRegistryTransactor{contract: contract}, DarkNodeRegistryFilterer: DarkNodeRegistryFilterer{contract: contract}}, nil
//}
//
//// DarkNodeRegistry is an auto generated Go binding around an Ethereum contract.
//type DarkNodeRegistry struct {
//	DarkNodeRegistryCaller     // Read-only binding to the contract
//	DarkNodeRegistryTransactor // Write-only binding to the contract
//	DarkNodeRegistryFilterer   // Log filterer for contract events
//}
//
//// DarkNodeRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
//type DarkNodeRegistryCaller struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// DarkNodeRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
//type DarkNodeRegistryTransactor struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// DarkNodeRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
//type DarkNodeRegistryFilterer struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// DarkNodeRegistrySession is an auto generated Go binding around an Ethereum contract,
//// with pre-set call and transact options.
//type DarkNodeRegistrySession struct {
//	Contract     *DarkNodeRegistry // Generic contract binding to set the session for
//	CallOpts     bind.CallOpts     // Call options to use throughout this session
//	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
//}
//
//// DarkNodeRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
//// with pre-set call options.
//type DarkNodeRegistryCallerSession struct {
//	Contract *DarkNodeRegistryCaller // Generic contract caller binding to set the session for
//	CallOpts bind.CallOpts           // Call options to use throughout this session
//}
//
//// DarkNodeRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
//// with pre-set transact options.
//type DarkNodeRegistryTransactorSession struct {
//	Contract     *DarkNodeRegistryTransactor // Generic contract transactor binding to set the session for
//	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
//}
//
//// DarkNodeRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
//type DarkNodeRegistryRaw struct {
//	Contract *DarkNodeRegistry // Generic contract binding to access the raw methods on
//}
//
//// DarkNodeRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
//type DarkNodeRegistryCallerRaw struct {
//	Contract *DarkNodeRegistryCaller // Generic read-only contract binding to access the raw methods on
//}
//
//// DarkNodeRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
//type DarkNodeRegistryTransactorRaw struct {
//	Contract *DarkNodeRegistryTransactor // Generic write-only contract binding to access the raw methods on
//}
//
//// NewDarkNodeRegistry creates a new instance of DarkNodeRegistry, bound to a specific deployed contract.
//func NewDarkNodeRegistry(address common.Address, backend bind.ContractBackend) (*DarkNodeRegistry, error) {
//	contract, err := bindDarkNodeRegistry(address, backend, backend, backend)
//	if err != nil {
//		return nil, err
//	}
//	return &DarkNodeRegistry{DarkNodeRegistryCaller: DarkNodeRegistryCaller{contract: contract}, DarkNodeRegistryTransactor: DarkNodeRegistryTransactor{contract: contract}, DarkNodeRegistryFilterer: DarkNodeRegistryFilterer{contract: contract}}, nil
//}
//
//// NewDarkNodeRegistryCaller creates a new read-only instance of DarkNodeRegistry, bound to a specific deployed contract.
//func NewDarkNodeRegistryCaller(address common.Address, caller bind.ContractCaller) (*DarkNodeRegistryCaller, error) {
//	contract, err := bindDarkNodeRegistry(address, caller, nil, nil)
//	if err != nil {
//		return nil, err
//	}
//	return &DarkNodeRegistryCaller{contract: contract}, nil
//}
//
//// NewDarkNodeRegistryTransactor creates a new write-only instance of DarkNodeRegistry, bound to a specific deployed contract.
//func NewDarkNodeRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*DarkNodeRegistryTransactor, error) {
//	contract, err := bindDarkNodeRegistry(address, nil, transactor, nil)
//	if err != nil {
//		return nil, err
//	}
//	return &DarkNodeRegistryTransactor{contract: contract}, nil
//}
//
//// NewDarkNodeRegistryFilterer creates a new log filterer instance of DarkNodeRegistry, bound to a specific deployed contract.
//func NewDarkNodeRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*DarkNodeRegistryFilterer, error) {
//	contract, err := bindDarkNodeRegistry(address, nil, nil, filterer)
//	if err != nil {
//		return nil, err
//	}
//	return &DarkNodeRegistryFilterer{contract: contract}, nil
//}
//
//// bindDarkNodeRegistry binds a generic wrapper to an already deployed contract.
//func bindDarkNodeRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
//	parsed, err := abi.JSON(strings.NewReader(DarkNodeRegistryABI))
//	if err != nil {
//		return nil, err
//	}
//	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
//}
//
//// Call invokes the (constant) contract method with params as input values and
//// sets the output to result. The result type might be a single field for simple
//// returns, a slice of interfaces for anonymous returns and a struct for named
//// returns.
//func (_DarkNodeRegistry *DarkNodeRegistryRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
//	return _DarkNodeRegistry.Contract.DarkNodeRegistryCaller.contract.Call(opts, result, method, params...)
//}
//
//// Transfer initiates a plain transaction to move funds to the contract, calling
//// its default method if one is available.
//func (_DarkNodeRegistry *DarkNodeRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _DarkNodeRegistry.Contract.DarkNodeRegistryTransactor.contract.Transfer(opts)
//}
//
//// Transact invokes the (paid) contract method with params as input values.
//func (_DarkNodeRegistry *DarkNodeRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
//	return _DarkNodeRegistry.Contract.DarkNodeRegistryTransactor.contract.Transact(opts, method, params...)
//}
//
//// Call invokes the (constant) contract method with params as input values and
//// sets the output to result. The result type might be a single field for simple
//// returns, a slice of interfaces for anonymous returns and a struct for named
//// returns.
//func (_DarkNodeRegistry *DarkNodeRegistryCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
//	return _DarkNodeRegistry.Contract.contract.Call(opts, result, method, params...)
//}
//
//// Transfer initiates a plain transaction to move funds to the contract, calling
//// its default method if one is available.
//func (_DarkNodeRegistry *DarkNodeRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _DarkNodeRegistry.Contract.contract.Transfer(opts)
//}
//
//// Transact invokes the (paid) contract method with params as input values.
//func (_DarkNodeRegistry *DarkNodeRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
//	return _DarkNodeRegistry.Contract.contract.Transact(opts, method, params...)
//}
//
//// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
////
//// Solidity: function currentEpoch() constant returns(blockhash uint256, timestamp uint256)
//func (_DarkNodeRegistry *DarkNodeRegistryCaller) CurrentEpoch(opts *bind.CallOpts) (struct {
//	Blockhash *big.Int
//	Timestamp *big.Int
//}, error) {
//	ret := new(struct {
//		Blockhash *big.Int
//		Timestamp *big.Int
//	})
//	out := ret
//	err := _DarkNodeRegistry.contract.Call(opts, out, "currentEpoch")
//	return *ret, err
//}
//
//// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
////
//// Solidity: function currentEpoch() constant returns(blockhash uint256, timestamp uint256)
//func (_DarkNodeRegistry *DarkNodeRegistrySession) CurrentEpoch() (struct {
//	Blockhash *big.Int
//	Timestamp *big.Int
//}, error) {
//	return _DarkNodeRegistry.Contract.CurrentEpoch(&_DarkNodeRegistry.CallOpts)
//}
//
//// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
////
//// Solidity: function currentEpoch() constant returns(blockhash uint256, timestamp uint256)
//func (_DarkNodeRegistry *DarkNodeRegistryCallerSession) CurrentEpoch() (struct {
//	Blockhash *big.Int
//	Timestamp *big.Int
//}, error) {
//	return _DarkNodeRegistry.Contract.CurrentEpoch(&_DarkNodeRegistry.CallOpts)
//}
//
//// GetBond is a free data retrieval call binding the contract method 0x68f209eb.
////
//// Solidity: function getBond(_darkNodeID bytes20) constant returns(uint256)
//func (_DarkNodeRegistry *DarkNodeRegistryCaller) GetBond(opts *bind.CallOpts, _darkNodeID [20]byte) (*big.Int, error) {
//	var (
//		ret0 = new(*big.Int)
//	)
//	out := ret0
//	err := _DarkNodeRegistry.contract.Call(opts, out, "getBond", _darkNodeID)
//	return *ret0, err
//}
//
//// GetBond is a free data retrieval call binding the contract method 0x68f209eb.
////
//// Solidity: function getBond(_darkNodeID bytes20) constant returns(uint256)
//func (_DarkNodeRegistry *DarkNodeRegistrySession) GetBond(_darkNodeID [20]byte) (*big.Int, error) {
//	return _DarkNodeRegistry.Contract.GetBond(&_DarkNodeRegistry.CallOpts, _darkNodeID)
//}
//
//// GetBond is a free data retrieval call binding the contract method 0x68f209eb.
////
//// Solidity: function getBond(_darkNodeID bytes20) constant returns(uint256)
//func (_DarkNodeRegistry *DarkNodeRegistryCallerSession) GetBond(_darkNodeID [20]byte) (*big.Int, error) {
//	return _DarkNodeRegistry.Contract.GetBond(&_DarkNodeRegistry.CallOpts, _darkNodeID)
//}
//
//// GetDarkNodes is a free data retrieval call binding the contract method 0x879ae084.
////
//// Solidity: function getDarkNodes() constant returns(bytes20[])
//func (_DarkNodeRegistry *DarkNodeRegistryCaller) GetDarkNodes(opts *bind.CallOpts) ([][20]byte, error) {
//	var (
//		ret0 = new([][20]byte)
//	)
//	out := ret0
//	err := _DarkNodeRegistry.contract.Call(opts, out, "getDarkNodes")
//	return *ret0, err
//}
//
//// GetDarkNodes is a free data retrieval call binding the contract method 0x879ae084.
////
//// Solidity: function getDarkNodes() constant returns(bytes20[])
//func (_DarkNodeRegistry *DarkNodeRegistrySession) GetDarkNodes() ([][20]byte, error) {
//	return _DarkNodeRegistry.Contract.GetDarkNodes(&_DarkNodeRegistry.CallOpts)
//}
//
//// GetDarkNodes is a free data retrieval call binding the contract method 0x879ae084.
////
//// Solidity: function getDarkNodes() constant returns(bytes20[])
//func (_DarkNodeRegistry *DarkNodeRegistryCallerSession) GetDarkNodes() ([][20]byte, error) {
//	return _DarkNodeRegistry.Contract.GetDarkNodes(&_DarkNodeRegistry.CallOpts)
//}
//
//// GetOwner is a free data retrieval call binding the contract method 0xe487eb58.
////
//// Solidity: function getOwner(_darkNodeID bytes20) constant returns(address)
//func (_DarkNodeRegistry *DarkNodeRegistryCaller) GetOwner(opts *bind.CallOpts, _darkNodeID [20]byte) (common.Address, error) {
//	var (
//		ret0 = new(common.Address)
//	)
//	out := ret0
//	err := _DarkNodeRegistry.contract.Call(opts, out, "getOwner", _darkNodeID)
//	return *ret0, err
//}
//
//// GetOwner is a free data retrieval call binding the contract method 0xe487eb58.
////
//// Solidity: function getOwner(_darkNodeID bytes20) constant returns(address)
//func (_DarkNodeRegistry *DarkNodeRegistrySession) GetOwner(_darkNodeID [20]byte) (common.Address, error) {
//	return _DarkNodeRegistry.Contract.GetOwner(&_DarkNodeRegistry.CallOpts, _darkNodeID)
//}
//
//// GetOwner is a free data retrieval call binding the contract method 0xe487eb58.
////
//// Solidity: function getOwner(_darkNodeID bytes20) constant returns(address)
//func (_DarkNodeRegistry *DarkNodeRegistryCallerSession) GetOwner(_darkNodeID [20]byte) (common.Address, error) {
//	return _DarkNodeRegistry.Contract.GetOwner(&_DarkNodeRegistry.CallOpts, _darkNodeID)
//}
//
//// GetPublicKey is a free data retrieval call binding the contract method 0x32ccd52f.
////
//// Solidity: function getPublicKey(_darkNodeID bytes20) constant returns(bytes)
//func (_DarkNodeRegistry *DarkNodeRegistryCaller) GetPublicKey(opts *bind.CallOpts, _darkNodeID [20]byte) ([]byte, error) {
//	var (
//		ret0 = new([]byte)
//	)
//	out := ret0
//	err := _DarkNodeRegistry.contract.Call(opts, out, "getPublicKey", _darkNodeID)
//	return *ret0, err
//}
//
//// GetPublicKey is a free data retrieval call binding the contract method 0x32ccd52f.
////
//// Solidity: function getPublicKey(_darkNodeID bytes20) constant returns(bytes)
//func (_DarkNodeRegistry *DarkNodeRegistrySession) GetPublicKey(_darkNodeID [20]byte) ([]byte, error) {
//	return _DarkNodeRegistry.Contract.GetPublicKey(&_DarkNodeRegistry.CallOpts, _darkNodeID)
//}
//
//// GetPublicKey is a free data retrieval call binding the contract method 0x32ccd52f.
////
//// Solidity: function getPublicKey(_darkNodeID bytes20) constant returns(bytes)
//func (_DarkNodeRegistry *DarkNodeRegistryCallerSession) GetPublicKey(_darkNodeID [20]byte) ([]byte, error) {
//	return _DarkNodeRegistry.Contract.GetPublicKey(&_DarkNodeRegistry.CallOpts, _darkNodeID)
//}
//
//// IsDeregistered is a free data retrieval call binding the contract method 0x171f6ea8.
////
//// Solidity: function isDeregistered(_darkNodeID bytes20) constant returns(bool)
//func (_DarkNodeRegistry *DarkNodeRegistryCaller) IsDeregistered(opts *bind.CallOpts, _darkNodeID [20]byte) (bool, error) {
//	var (
//		ret0 = new(bool)
//	)
//	out := ret0
//	err := _DarkNodeRegistry.contract.Call(opts, out, "isDeregistered", _darkNodeID)
//	return *ret0, err
//}
//
//// IsDeregistered is a free data retrieval call binding the contract method 0x171f6ea8.
////
//// Solidity: function isDeregistered(_darkNodeID bytes20) constant returns(bool)
//func (_DarkNodeRegistry *DarkNodeRegistrySession) IsDeregistered(_darkNodeID [20]byte) (bool, error) {
//	return _DarkNodeRegistry.Contract.IsDeregistered(&_DarkNodeRegistry.CallOpts, _darkNodeID)
//}
//
//// IsDeregistered is a free data retrieval call binding the contract method 0x171f6ea8.
////
//// Solidity: function isDeregistered(_darkNodeID bytes20) constant returns(bool)
//func (_DarkNodeRegistry *DarkNodeRegistryCallerSession) IsDeregistered(_darkNodeID [20]byte) (bool, error) {
//	return _DarkNodeRegistry.Contract.IsDeregistered(&_DarkNodeRegistry.CallOpts, _darkNodeID)
//}
//
//// IsRegistered is a free data retrieval call binding the contract method 0x4f5550fc.
////
//// Solidity: function isRegistered(_darkNodeID bytes20) constant returns(bool)
//func (_DarkNodeRegistry *DarkNodeRegistryCaller) IsRegistered(opts *bind.CallOpts, _darkNodeID [20]byte) (bool, error) {
//	var (
//		ret0 = new(bool)
//	)
//	out := ret0
//	err := _DarkNodeRegistry.contract.Call(opts, out, "isRegistered", _darkNodeID)
//	return *ret0, err
//}
//
//// IsRegistered is a free data retrieval call binding the contract method 0x4f5550fc.
////
//// Solidity: function isRegistered(_darkNodeID bytes20) constant returns(bool)
//func (_DarkNodeRegistry *DarkNodeRegistrySession) IsRegistered(_darkNodeID [20]byte) (bool, error) {
//	return _DarkNodeRegistry.Contract.IsRegistered(&_DarkNodeRegistry.CallOpts, _darkNodeID)
//}
//
//// IsRegistered is a free data retrieval call binding the contract method 0x4f5550fc.
////
//// Solidity: function isRegistered(_darkNodeID bytes20) constant returns(bool)
//func (_DarkNodeRegistry *DarkNodeRegistryCallerSession) IsRegistered(_darkNodeID [20]byte) (bool, error) {
//	return _DarkNodeRegistry.Contract.IsRegistered(&_DarkNodeRegistry.CallOpts, _darkNodeID)
//}
//
//// IsUnregistered is a free data retrieval call binding the contract method 0xd3841c25.
////
//// Solidity: function isUnregistered(_darkNodeID bytes20) constant returns(bool)
//func (_DarkNodeRegistry *DarkNodeRegistryCaller) IsUnregistered(opts *bind.CallOpts, _darkNodeID [20]byte) (bool, error) {
//	var (
//		ret0 = new(bool)
//	)
//	out := ret0
//	err := _DarkNodeRegistry.contract.Call(opts, out, "isUnregistered", _darkNodeID)
//	return *ret0, err
//}
//
//// IsUnregistered is a free data retrieval call binding the contract method 0xd3841c25.
////
//// Solidity: function isUnregistered(_darkNodeID bytes20) constant returns(bool)
//func (_DarkNodeRegistry *DarkNodeRegistrySession) IsUnregistered(_darkNodeID [20]byte) (bool, error) {
//	return _DarkNodeRegistry.Contract.IsUnregistered(&_DarkNodeRegistry.CallOpts, _darkNodeID)
//}
//
//// IsUnregistered is a free data retrieval call binding the contract method 0xd3841c25.
////
//// Solidity: function isUnregistered(_darkNodeID bytes20) constant returns(bool)
//func (_DarkNodeRegistry *DarkNodeRegistryCallerSession) IsUnregistered(_darkNodeID [20]byte) (bool, error) {
//	return _DarkNodeRegistry.Contract.IsUnregistered(&_DarkNodeRegistry.CallOpts, _darkNodeID)
//}
//
//// MinimumBond is a free data retrieval call binding the contract method 0xaa7517e1.
////
//// Solidity: function minimumBond() constant returns(uint256)
//func (_DarkNodeRegistry *DarkNodeRegistryCaller) MinimumBond(opts *bind.CallOpts) (*big.Int, error) {
//	var (
//		ret0 = new(*big.Int)
//	)
//	out := ret0
//	err := _DarkNodeRegistry.contract.Call(opts, out, "minimumBond")
//	return *ret0, err
//}
//
//// MinimumBond is a free data retrieval call binding the contract method 0xaa7517e1.
////
//// Solidity: function minimumBond() constant returns(uint256)
//func (_DarkNodeRegistry *DarkNodeRegistrySession) MinimumBond() (*big.Int, error) {
//	return _DarkNodeRegistry.Contract.MinimumBond(&_DarkNodeRegistry.CallOpts)
//}
//
//// MinimumBond is a free data retrieval call binding the contract method 0xaa7517e1.
////
//// Solidity: function minimumBond() constant returns(uint256)
//func (_DarkNodeRegistry *DarkNodeRegistryCallerSession) MinimumBond() (*big.Int, error) {
//	return _DarkNodeRegistry.Contract.MinimumBond(&_DarkNodeRegistry.CallOpts)
//}
//
//// MinimumDarkPoolSize is a free data retrieval call binding the contract method 0xb31575d5.
////
//// Solidity: function minimumDarkPoolSize() constant returns(uint256)
//func (_DarkNodeRegistry *DarkNodeRegistryCaller) MinimumDarkPoolSize(opts *bind.CallOpts) (*big.Int, error) {
//	var (
//		ret0 = new(*big.Int)
//	)
//	out := ret0
//	err := _DarkNodeRegistry.contract.Call(opts, out, "minimumDarkPoolSize")
//	return *ret0, err
//}
//
//// MinimumDarkPoolSize is a free data retrieval call binding the contract method 0xb31575d5.
////
//// Solidity: function minimumDarkPoolSize() constant returns(uint256)
//func (_DarkNodeRegistry *DarkNodeRegistrySession) MinimumDarkPoolSize() (*big.Int, error) {
//	return _DarkNodeRegistry.Contract.MinimumDarkPoolSize(&_DarkNodeRegistry.CallOpts)
//}
//
//// MinimumDarkPoolSize is a free data retrieval call binding the contract method 0xb31575d5.
////
//// Solidity: function minimumDarkPoolSize() constant returns(uint256)
//func (_DarkNodeRegistry *DarkNodeRegistryCallerSession) MinimumDarkPoolSize() (*big.Int, error) {
//	return _DarkNodeRegistry.Contract.MinimumDarkPoolSize(&_DarkNodeRegistry.CallOpts)
//}
//
//// MinimumEpochInterval is a free data retrieval call binding the contract method 0x55cacda5.
////
//// Solidity: function minimumEpochInterval() constant returns(uint256)
//func (_DarkNodeRegistry *DarkNodeRegistryCaller) MinimumEpochInterval(opts *bind.CallOpts) (*big.Int, error) {
//	var (
//		ret0 = new(*big.Int)
//	)
//	out := ret0
//	err := _DarkNodeRegistry.contract.Call(opts, out, "minimumEpochInterval")
//	return *ret0, err
//}
//
//// MinimumEpochInterval is a free data retrieval call binding the contract method 0x55cacda5.
////
//// Solidity: function minimumEpochInterval() constant returns(uint256)
//func (_DarkNodeRegistry *DarkNodeRegistrySession) MinimumEpochInterval() (*big.Int, error) {
//	return _DarkNodeRegistry.Contract.MinimumEpochInterval(&_DarkNodeRegistry.CallOpts)
//}
//
//// MinimumEpochInterval is a free data retrieval call binding the contract method 0x55cacda5.
////
//// Solidity: function minimumEpochInterval() constant returns(uint256)
//func (_DarkNodeRegistry *DarkNodeRegistryCallerSession) MinimumEpochInterval() (*big.Int, error) {
//	return _DarkNodeRegistry.Contract.MinimumEpochInterval(&_DarkNodeRegistry.CallOpts)
//}
//
//// NumDarkNodes is a free data retrieval call binding the contract method 0x060a2cfd.
////
//// Solidity: function numDarkNodes() constant returns(uint256)
//func (_DarkNodeRegistry *DarkNodeRegistryCaller) NumDarkNodes(opts *bind.CallOpts) (*big.Int, error) {
//	var (
//		ret0 = new(*big.Int)
//	)
//	out := ret0
//	err := _DarkNodeRegistry.contract.Call(opts, out, "numDarkNodes")
//	return *ret0, err
//}
//
//// NumDarkNodes is a free data retrieval call binding the contract method 0x060a2cfd.
////
//// Solidity: function numDarkNodes() constant returns(uint256)
//func (_DarkNodeRegistry *DarkNodeRegistrySession) NumDarkNodes() (*big.Int, error) {
//	return _DarkNodeRegistry.Contract.NumDarkNodes(&_DarkNodeRegistry.CallOpts)
//}
//
//// NumDarkNodes is a free data retrieval call binding the contract method 0x060a2cfd.
////
//// Solidity: function numDarkNodes() constant returns(uint256)
//func (_DarkNodeRegistry *DarkNodeRegistryCallerSession) NumDarkNodes() (*big.Int, error) {
//	return _DarkNodeRegistry.Contract.NumDarkNodes(&_DarkNodeRegistry.CallOpts)
//}
//
//// NumDarkNodesNextEpoch is a free data retrieval call binding the contract method 0x0620eb92.
////
//// Solidity: function numDarkNodesNextEpoch() constant returns(uint256)
//func (_DarkNodeRegistry *DarkNodeRegistryCaller) NumDarkNodesNextEpoch(opts *bind.CallOpts) (*big.Int, error) {
//	var (
//		ret0 = new(*big.Int)
//	)
//	out := ret0
//	err := _DarkNodeRegistry.contract.Call(opts, out, "numDarkNodesNextEpoch")
//	return *ret0, err
//}
//
//// NumDarkNodesNextEpoch is a free data retrieval call binding the contract method 0x0620eb92.
////
//// Solidity: function numDarkNodesNextEpoch() constant returns(uint256)
//func (_DarkNodeRegistry *DarkNodeRegistrySession) NumDarkNodesNextEpoch() (*big.Int, error) {
//	return _DarkNodeRegistry.Contract.NumDarkNodesNextEpoch(&_DarkNodeRegistry.CallOpts)
//}
//
//// NumDarkNodesNextEpoch is a free data retrieval call binding the contract method 0x0620eb92.
////
//// Solidity: function numDarkNodesNextEpoch() constant returns(uint256)
//func (_DarkNodeRegistry *DarkNodeRegistryCallerSession) NumDarkNodesNextEpoch() (*big.Int, error) {
//	return _DarkNodeRegistry.Contract.NumDarkNodesNextEpoch(&_DarkNodeRegistry.CallOpts)
//}
//
//// Deregister is a paid mutator transaction binding the contract method 0xe08b4c8a.
////
//// Solidity: function deregister(_darkNodeID bytes20) returns()
//func (_DarkNodeRegistry *DarkNodeRegistryTransactor) Deregister(opts *bind.TransactOpts, _darkNodeID [20]byte) (*types.Transaction, error) {
//	return _DarkNodeRegistry.contract.Transact(opts, "deregister", _darkNodeID)
//}
//
//// Deregister is a paid mutator transaction binding the contract method 0xe08b4c8a.
////
//// Solidity: function deregister(_darkNodeID bytes20) returns()
//func (_DarkNodeRegistry *DarkNodeRegistrySession) Deregister(_darkNodeID [20]byte) (*types.Transaction, error) {
//	return _DarkNodeRegistry.Contract.Deregister(&_DarkNodeRegistry.TransactOpts, _darkNodeID)
//}
//
//// Deregister is a paid mutator transaction binding the contract method 0xe08b4c8a.
////
//// Solidity: function deregister(_darkNodeID bytes20) returns()
//func (_DarkNodeRegistry *DarkNodeRegistryTransactorSession) Deregister(_darkNodeID [20]byte) (*types.Transaction, error) {
//	return _DarkNodeRegistry.Contract.Deregister(&_DarkNodeRegistry.TransactOpts, _darkNodeID)
//}
//
//// Epoch is a paid mutator transaction binding the contract method 0x900cf0cf.
////
//// Solidity: function epoch() returns()
//func (_DarkNodeRegistry *DarkNodeRegistryTransactor) Epoch(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _DarkNodeRegistry.contract.Transact(opts, "epoch")
//}
//
//// Epoch is a paid mutator transaction binding the contract method 0x900cf0cf.
////
//// Solidity: function epoch() returns()
//func (_DarkNodeRegistry *DarkNodeRegistrySession) Epoch() (*types.Transaction, error) {
//	return _DarkNodeRegistry.Contract.Epoch(&_DarkNodeRegistry.TransactOpts)
//}
//
//// Epoch is a paid mutator transaction binding the contract method 0x900cf0cf.
////
//// Solidity: function epoch() returns()
//func (_DarkNodeRegistry *DarkNodeRegistryTransactorSession) Epoch() (*types.Transaction, error) {
//	return _DarkNodeRegistry.Contract.Epoch(&_DarkNodeRegistry.TransactOpts)
//}
//
//// Refund is a paid mutator transaction binding the contract method 0x5a8f9b81.
////
//// Solidity: function refund(_darkNodeID bytes20) returns()
//func (_DarkNodeRegistry *DarkNodeRegistryTransactor) Refund(opts *bind.TransactOpts, _darkNodeID [20]byte) (*types.Transaction, error) {
//	return _DarkNodeRegistry.contract.Transact(opts, "refund", _darkNodeID)
//}
//
//// Refund is a paid mutator transaction binding the contract method 0x5a8f9b81.
////
//// Solidity: function refund(_darkNodeID bytes20) returns()
//func (_DarkNodeRegistry *DarkNodeRegistrySession) Refund(_darkNodeID [20]byte) (*types.Transaction, error) {
//	return _DarkNodeRegistry.Contract.Refund(&_DarkNodeRegistry.TransactOpts, _darkNodeID)
//}
//
//// Refund is a paid mutator transaction binding the contract method 0x5a8f9b81.
////
//// Solidity: function refund(_darkNodeID bytes20) returns()
//func (_DarkNodeRegistry *DarkNodeRegistryTransactorSession) Refund(_darkNodeID [20]byte) (*types.Transaction, error) {
//	return _DarkNodeRegistry.Contract.Refund(&_DarkNodeRegistry.TransactOpts, _darkNodeID)
//}
//
//// Register is a paid mutator transaction binding the contract method 0x375a8be3.
////
//// Solidity: function register(_darkNodeID bytes20, _publicKey bytes, _bond uint256) returns()
//func (_DarkNodeRegistry *DarkNodeRegistryTransactor) Register(opts *bind.TransactOpts, _darkNodeID [20]byte, _publicKey []byte, _bond *big.Int) (*types.Transaction, error) {
//	return _DarkNodeRegistry.contract.Transact(opts, "register", _darkNodeID, _publicKey, _bond)
//}
//
//// Register is a paid mutator transaction binding the contract method 0x375a8be3.
////
//// Solidity: function register(_darkNodeID bytes20, _publicKey bytes, _bond uint256) returns()
//func (_DarkNodeRegistry *DarkNodeRegistrySession) Register(_darkNodeID [20]byte, _publicKey []byte, _bond *big.Int) (*types.Transaction, error) {
//	return _DarkNodeRegistry.Contract.Register(&_DarkNodeRegistry.TransactOpts, _darkNodeID, _publicKey, _bond)
//}
//
//// Register is a paid mutator transaction binding the contract method 0x375a8be3.
////
//// Solidity: function register(_darkNodeID bytes20, _publicKey bytes, _bond uint256) returns()
//func (_DarkNodeRegistry *DarkNodeRegistryTransactorSession) Register(_darkNodeID [20]byte, _publicKey []byte, _bond *big.Int) (*types.Transaction, error) {
//	return _DarkNodeRegistry.Contract.Register(&_DarkNodeRegistry.TransactOpts, _darkNodeID, _publicKey, _bond)
//}
//
//// DarkNodeRegistryDarkNodeDeregisteredIterator is returned from FilterDarkNodeDeregistered and is used to iterate over the raw logs and unpacked data for DarkNodeDeregistered events raised by the DarkNodeRegistry contract.
//type DarkNodeRegistryDarkNodeDeregisteredIterator struct {
//	Event *DarkNodeRegistryDarkNodeDeregistered // Event containing the contract specifics and raw log
//
//	contract *bind.BoundContract // Generic contract to use for unpacking event data
//	event    string              // Event name to use for unpacking event data
//
//	logs chan types.Log        // Log channel receiving the found contract events
//	sub  ethereum.Subscription // Subscription for errors, completion and termination
//	done bool                  // Whether the subscription completed delivering logs
//	fail error                 // Occurred error to stop iteration
//}
//
//// Next advances the iterator to the subsequent event, returning whether there
//// are any more events found. In case of a retrieval or parsing error, false is
//// returned and Error() can be queried for the exact failure.
//func (it *DarkNodeRegistryDarkNodeDeregisteredIterator) Next() bool {
//	// If the iterator failed, stop iterating
//	if it.fail != nil {
//		return false
//	}
//	// If the iterator completed, deliver directly whatever's available
//	if it.done {
//		select {
//		case log := <-it.logs:
//			it.Event = new(DarkNodeRegistryDarkNodeDeregistered)
//			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//				it.fail = err
//				return false
//			}
//			it.Event.Raw = log
//			return true
//
//		default:
//			return false
//		}
//	}
//	// Iterator still in progress, wait for either a data or an error event
//	select {
//	case log := <-it.logs:
//		it.Event = new(DarkNodeRegistryDarkNodeDeregistered)
//		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//			it.fail = err
//			return false
//		}
//		it.Event.Raw = log
//		return true
//
//	case err := <-it.sub.Err():
//		it.done = true
//		it.fail = err
//		return it.Next()
//	}
//}
//
//// Error returns any retrieval or parsing error occurred during filtering.
//func (it *DarkNodeRegistryDarkNodeDeregisteredIterator) Error() error {
//	return it.fail
//}
//
//// Close terminates the iteration process, releasing any pending underlying
//// resources.
//func (it *DarkNodeRegistryDarkNodeDeregisteredIterator) Close() error {
//	it.sub.Unsubscribe()
//	return nil
//}
//
//// DarkNodeRegistryDarkNodeDeregistered represents a DarkNodeDeregistered event raised by the DarkNodeRegistry contract.
//type DarkNodeRegistryDarkNodeDeregistered struct {
//	DarkNodeID [20]byte
//	Raw        types.Log // Blockchain specific contextual infos
//}
//
//// FilterDarkNodeDeregistered is a free log retrieval operation binding the contract event 0xe3365373f442312a7d66aa1003b60fd4a0f9ed1c85fd8bb384240f8eb4fe916c.
////
//// Solidity: event DarkNodeDeregistered(_darkNodeID bytes20)
//func (_DarkNodeRegistry *DarkNodeRegistryFilterer) FilterDarkNodeDeregistered(opts *bind.FilterOpts) (*DarkNodeRegistryDarkNodeDeregisteredIterator, error) {
//
//	logs, sub, err := _DarkNodeRegistry.contract.FilterLogs(opts, "DarkNodeDeregistered")
//	if err != nil {
//		return nil, err
//	}
//	return &DarkNodeRegistryDarkNodeDeregisteredIterator{contract: _DarkNodeRegistry.contract, event: "DarkNodeDeregistered", logs: logs, sub: sub}, nil
//}
//
//// WatchDarkNodeDeregistered is a free log subscription operation binding the contract event 0xe3365373f442312a7d66aa1003b60fd4a0f9ed1c85fd8bb384240f8eb4fe916c.
////
//// Solidity: event DarkNodeDeregistered(_darkNodeID bytes20)
//func (_DarkNodeRegistry *DarkNodeRegistryFilterer) WatchDarkNodeDeregistered(opts *bind.WatchOpts, sink chan<- *DarkNodeRegistryDarkNodeDeregistered) (event.Subscription, error) {
//
//	logs, sub, err := _DarkNodeRegistry.contract.WatchLogs(opts, "DarkNodeDeregistered")
//	if err != nil {
//		return nil, err
//	}
//	return event.NewSubscription(func(quit <-chan struct{}) error {
//		defer sub.Unsubscribe()
//		for {
//			select {
//			case log := <-logs:
//				// New log arrived, parse the event and forward to the user
//				event := new(DarkNodeRegistryDarkNodeDeregistered)
//				if err := _DarkNodeRegistry.contract.UnpackLog(event, "DarkNodeDeregistered", log); err != nil {
//					return err
//				}
//				event.Raw = log
//
//				select {
//				case sink <- event:
//				case err := <-sub.Err():
//					return err
//				case <-quit:
//					return nil
//				}
//			case err := <-sub.Err():
//				return err
//			case <-quit:
//				return nil
//			}
//		}
//	}), nil
//}
//
//// DarkNodeRegistryDarkNodeRegisteredIterator is returned from FilterDarkNodeRegistered and is used to iterate over the raw logs and unpacked data for DarkNodeRegistered events raised by the DarkNodeRegistry contract.
//type DarkNodeRegistryDarkNodeRegisteredIterator struct {
//	Event *DarkNodeRegistryDarkNodeRegistered // Event containing the contract specifics and raw log
//
//	contract *bind.BoundContract // Generic contract to use for unpacking event data
//	event    string              // Event name to use for unpacking event data
//
//	logs chan types.Log        // Log channel receiving the found contract events
//	sub  ethereum.Subscription // Subscription for errors, completion and termination
//	done bool                  // Whether the subscription completed delivering logs
//	fail error                 // Occurred error to stop iteration
//}
//
//// Next advances the iterator to the subsequent event, returning whether there
//// are any more events found. In case of a retrieval or parsing error, false is
//// returned and Error() can be queried for the exact failure.
//func (it *DarkNodeRegistryDarkNodeRegisteredIterator) Next() bool {
//	// If the iterator failed, stop iterating
//	if it.fail != nil {
//		return false
//	}
//	// If the iterator completed, deliver directly whatever's available
//	if it.done {
//		select {
//		case log := <-it.logs:
//			it.Event = new(DarkNodeRegistryDarkNodeRegistered)
//			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//				it.fail = err
//				return false
//			}
//			it.Event.Raw = log
//			return true
//
//		default:
//			return false
//		}
//	}
//	// Iterator still in progress, wait for either a data or an error event
//	select {
//	case log := <-it.logs:
//		it.Event = new(DarkNodeRegistryDarkNodeRegistered)
//		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//			it.fail = err
//			return false
//		}
//		it.Event.Raw = log
//		return true
//
//	case err := <-it.sub.Err():
//		it.done = true
//		it.fail = err
//		return it.Next()
//	}
//}
//
//// Error returns any retrieval or parsing error occurred during filtering.
//func (it *DarkNodeRegistryDarkNodeRegisteredIterator) Error() error {
//	return it.fail
//}
//
//// Close terminates the iteration process, releasing any pending underlying
//// resources.
//func (it *DarkNodeRegistryDarkNodeRegisteredIterator) Close() error {
//	it.sub.Unsubscribe()
//	return nil
//}
//
//// DarkNodeRegistryDarkNodeRegistered represents a DarkNodeRegistered event raised by the DarkNodeRegistry contract.
//type DarkNodeRegistryDarkNodeRegistered struct {
//	DarkNodeID [20]byte
//	Bond       *big.Int
//	Raw        types.Log // Blockchain specific contextual infos
//}
//
//// FilterDarkNodeRegistered is a free log retrieval operation binding the contract event 0xcde427a7822769e806a704726602a3a8a44458da1352d9b5a46bcfb95142ad38.
////
//// Solidity: event DarkNodeRegistered(_darkNodeID bytes20, _bond uint256)
//func (_DarkNodeRegistry *DarkNodeRegistryFilterer) FilterDarkNodeRegistered(opts *bind.FilterOpts) (*DarkNodeRegistryDarkNodeRegisteredIterator, error) {
//
//	logs, sub, err := _DarkNodeRegistry.contract.FilterLogs(opts, "DarkNodeRegistered")
//	if err != nil {
//		return nil, err
//	}
//	return &DarkNodeRegistryDarkNodeRegisteredIterator{contract: _DarkNodeRegistry.contract, event: "DarkNodeRegistered", logs: logs, sub: sub}, nil
//}
//
//// WatchDarkNodeRegistered is a free log subscription operation binding the contract event 0xcde427a7822769e806a704726602a3a8a44458da1352d9b5a46bcfb95142ad38.
////
//// Solidity: event DarkNodeRegistered(_darkNodeID bytes20, _bond uint256)
//func (_DarkNodeRegistry *DarkNodeRegistryFilterer) WatchDarkNodeRegistered(opts *bind.WatchOpts, sink chan<- *DarkNodeRegistryDarkNodeRegistered) (event.Subscription, error) {
//
//	logs, sub, err := _DarkNodeRegistry.contract.WatchLogs(opts, "DarkNodeRegistered")
//	if err != nil {
//		return nil, err
//	}
//	return event.NewSubscription(func(quit <-chan struct{}) error {
//		defer sub.Unsubscribe()
//		for {
//			select {
//			case log := <-logs:
//				// New log arrived, parse the event and forward to the user
//				event := new(DarkNodeRegistryDarkNodeRegistered)
//				if err := _DarkNodeRegistry.contract.UnpackLog(event, "DarkNodeRegistered", log); err != nil {
//					return err
//				}
//				event.Raw = log
//
//				select {
//				case sink <- event:
//				case err := <-sub.Err():
//					return err
//				case <-quit:
//					return nil
//				}
//			case err := <-sub.Err():
//				return err
//			case <-quit:
//				return nil
//			}
//		}
//	}), nil
//}
//
//// DarkNodeRegistryNewEpochIterator is returned from FilterNewEpoch and is used to iterate over the raw logs and unpacked data for NewEpoch events raised by the DarkNodeRegistry contract.
//type DarkNodeRegistryNewEpochIterator struct {
//	Event *DarkNodeRegistryNewEpoch // Event containing the contract specifics and raw log
//
//	contract *bind.BoundContract // Generic contract to use for unpacking event data
//	event    string              // Event name to use for unpacking event data
//
//	logs chan types.Log        // Log channel receiving the found contract events
//	sub  ethereum.Subscription // Subscription for errors, completion and termination
//	done bool                  // Whether the subscription completed delivering logs
//	fail error                 // Occurred error to stop iteration
//}
//
//// Next advances the iterator to the subsequent event, returning whether there
//// are any more events found. In case of a retrieval or parsing error, false is
//// returned and Error() can be queried for the exact failure.
//func (it *DarkNodeRegistryNewEpochIterator) Next() bool {
//	// If the iterator failed, stop iterating
//	if it.fail != nil {
//		return false
//	}
//	// If the iterator completed, deliver directly whatever's available
//	if it.done {
//		select {
//		case log := <-it.logs:
//			it.Event = new(DarkNodeRegistryNewEpoch)
//			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//				it.fail = err
//				return false
//			}
//			it.Event.Raw = log
//			return true
//
//		default:
//			return false
//		}
//	}
//	// Iterator still in progress, wait for either a data or an error event
//	select {
//	case log := <-it.logs:
//		it.Event = new(DarkNodeRegistryNewEpoch)
//		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//			it.fail = err
//			return false
//		}
//		it.Event.Raw = log
//		return true
//
//	case err := <-it.sub.Err():
//		it.done = true
//		it.fail = err
//		return it.Next()
//	}
//}
//
//// Error returns any retrieval or parsing error occurred during filtering.
//func (it *DarkNodeRegistryNewEpochIterator) Error() error {
//	return it.fail
//}
//
//// Close terminates the iteration process, releasing any pending underlying
//// resources.
//func (it *DarkNodeRegistryNewEpochIterator) Close() error {
//	it.sub.Unsubscribe()
//	return nil
//}
//
//// DarkNodeRegistryNewEpoch represents a NewEpoch event raised by the DarkNodeRegistry contract.
//type DarkNodeRegistryNewEpoch struct {
//	Raw types.Log // Blockchain specific contextual infos
//}
//
//// FilterNewEpoch is a free log retrieval operation binding the contract event 0xe358419ca0dd7928a310d787a606dfae5d869f5071249efa6107105e7afc40bc.
////
//// Solidity: event NewEpoch()
//func (_DarkNodeRegistry *DarkNodeRegistryFilterer) FilterNewEpoch(opts *bind.FilterOpts) (*DarkNodeRegistryNewEpochIterator, error) {
//
//	logs, sub, err := _DarkNodeRegistry.contract.FilterLogs(opts, "NewEpoch")
//	if err != nil {
//		return nil, err
//	}
//	return &DarkNodeRegistryNewEpochIterator{contract: _DarkNodeRegistry.contract, event: "NewEpoch", logs: logs, sub: sub}, nil
//}
//
//// WatchNewEpoch is a free log subscription operation binding the contract event 0xe358419ca0dd7928a310d787a606dfae5d869f5071249efa6107105e7afc40bc.
////
//// Solidity: event NewEpoch()
//func (_DarkNodeRegistry *DarkNodeRegistryFilterer) WatchNewEpoch(opts *bind.WatchOpts, sink chan<- *DarkNodeRegistryNewEpoch) (event.Subscription, error) {
//
//	logs, sub, err := _DarkNodeRegistry.contract.WatchLogs(opts, "NewEpoch")
//	if err != nil {
//		return nil, err
//	}
//	return event.NewSubscription(func(quit <-chan struct{}) error {
//		defer sub.Unsubscribe()
//		for {
//			select {
//			case log := <-logs:
//				// New log arrived, parse the event and forward to the user
//				event := new(DarkNodeRegistryNewEpoch)
//				if err := _DarkNodeRegistry.contract.UnpackLog(event, "NewEpoch", log); err != nil {
//					return err
//				}
//				event.Raw = log
//
//				select {
//				case sink <- event:
//				case err := <-sub.Err():
//					return err
//				case <-quit:
//					return nil
//				}
//			case err := <-sub.Err():
//				return err
//			case <-quit:
//				return nil
//			}
//		}
//	}), nil
//}
//
//// DarkNodeRegistryOwnerRefundedIterator is returned from FilterOwnerRefunded and is used to iterate over the raw logs and unpacked data for OwnerRefunded events raised by the DarkNodeRegistry contract.
//type DarkNodeRegistryOwnerRefundedIterator struct {
//	Event *DarkNodeRegistryOwnerRefunded // Event containing the contract specifics and raw log
//
//	contract *bind.BoundContract // Generic contract to use for unpacking event data
//	event    string              // Event name to use for unpacking event data
//
//	logs chan types.Log        // Log channel receiving the found contract events
//	sub  ethereum.Subscription // Subscription for errors, completion and termination
//	done bool                  // Whether the subscription completed delivering logs
//	fail error                 // Occurred error to stop iteration
//}
//
//// Next advances the iterator to the subsequent event, returning whether there
//// are any more events found. In case of a retrieval or parsing error, false is
//// returned and Error() can be queried for the exact failure.
//func (it *DarkNodeRegistryOwnerRefundedIterator) Next() bool {
//	// If the iterator failed, stop iterating
//	if it.fail != nil {
//		return false
//	}
//	// If the iterator completed, deliver directly whatever's available
//	if it.done {
//		select {
//		case log := <-it.logs:
//			it.Event = new(DarkNodeRegistryOwnerRefunded)
//			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//				it.fail = err
//				return false
//			}
//			it.Event.Raw = log
//			return true
//
//		default:
//			return false
//		}
//	}
//	// Iterator still in progress, wait for either a data or an error event
//	select {
//	case log := <-it.logs:
//		it.Event = new(DarkNodeRegistryOwnerRefunded)
//		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//			it.fail = err
//			return false
//		}
//		it.Event.Raw = log
//		return true
//
//	case err := <-it.sub.Err():
//		it.done = true
//		it.fail = err
//		return it.Next()
//	}
//}
//
//// Error returns any retrieval or parsing error occurred during filtering.
//func (it *DarkNodeRegistryOwnerRefundedIterator) Error() error {
//	return it.fail
//}
//
//// Close terminates the iteration process, releasing any pending underlying
//// resources.
//func (it *DarkNodeRegistryOwnerRefundedIterator) Close() error {
//	it.sub.Unsubscribe()
//	return nil
//}
//
//// DarkNodeRegistryOwnerRefunded represents a OwnerRefunded event raised by the DarkNodeRegistry contract.
//type DarkNodeRegistryOwnerRefunded struct {
//	Owner  common.Address
//	Amount *big.Int
//	Raw    types.Log // Blockchain specific contextual infos
//}
//
//// FilterOwnerRefunded is a free log retrieval operation binding the contract event 0x8dce8f4eb4097fbdf948a109703513f52a7fabcc7e328ba7f29ad8165033e953.
////
//// Solidity: event OwnerRefunded(_owner address, _amount uint256)
//func (_DarkNodeRegistry *DarkNodeRegistryFilterer) FilterOwnerRefunded(opts *bind.FilterOpts) (*DarkNodeRegistryOwnerRefundedIterator, error) {
//
//	logs, sub, err := _DarkNodeRegistry.contract.FilterLogs(opts, "OwnerRefunded")
//	if err != nil {
//		return nil, err
//	}
//	return &DarkNodeRegistryOwnerRefundedIterator{contract: _DarkNodeRegistry.contract, event: "OwnerRefunded", logs: logs, sub: sub}, nil
//}
//
//// WatchOwnerRefunded is a free log subscription operation binding the contract event 0x8dce8f4eb4097fbdf948a109703513f52a7fabcc7e328ba7f29ad8165033e953.
////
//// Solidity: event OwnerRefunded(_owner address, _amount uint256)
//func (_DarkNodeRegistry *DarkNodeRegistryFilterer) WatchOwnerRefunded(opts *bind.WatchOpts, sink chan<- *DarkNodeRegistryOwnerRefunded) (event.Subscription, error) {
//
//	logs, sub, err := _DarkNodeRegistry.contract.WatchLogs(opts, "OwnerRefunded")
//	if err != nil {
//		return nil, err
//	}
//	return event.NewSubscription(func(quit <-chan struct{}) error {
//		defer sub.Unsubscribe()
//		for {
//			select {
//			case log := <-logs:
//				// New log arrived, parse the event and forward to the user
//				event := new(DarkNodeRegistryOwnerRefunded)
//				if err := _DarkNodeRegistry.contract.UnpackLog(event, "OwnerRefunded", log); err != nil {
//					return err
//				}
//				event.Raw = log
//
//				select {
//				case sink <- event:
//				case err := <-sub.Err():
//					return err
//				case <-quit:
//					return nil
//				}
//			case err := <-sub.Err():
//				return err
//			case <-quit:
//				return nil
//			}
//		}
//	}), nil
//}
//
//// ERC20ABI is the input ABI used to generate the binding from.
//const ERC20ABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"who\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]"
//
//// ERC20Bin is the compiled bytecode used for deploying new contracts.
//const ERC20Bin = `0x`
//
//// DeployERC20 deploys a new Ethereum contract, binding an instance of ERC20 to it.
//func DeployERC20(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ERC20, error) {
//	parsed, err := abi.JSON(strings.NewReader(ERC20ABI))
//	if err != nil {
//		return common.Address{}, nil, nil, err
//	}
//	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ERC20Bin), backend)
//	if err != nil {
//		return common.Address{}, nil, nil, err
//	}
//	return address, tx, &ERC20{ERC20Caller: ERC20Caller{contract: contract}, ERC20Transactor: ERC20Transactor{contract: contract}, ERC20Filterer: ERC20Filterer{contract: contract}}, nil
//}
//
//// ERC20 is an auto generated Go binding around an Ethereum contract.
//type ERC20 struct {
//	ERC20Caller     // Read-only binding to the contract
//	ERC20Transactor // Write-only binding to the contract
//	ERC20Filterer   // Log filterer for contract events
//}
//
//// ERC20Caller is an auto generated read-only Go binding around an Ethereum contract.
//type ERC20Caller struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// ERC20Transactor is an auto generated write-only Go binding around an Ethereum contract.
//type ERC20Transactor struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// ERC20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
//type ERC20Filterer struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// ERC20Session is an auto generated Go binding around an Ethereum contract,
//// with pre-set call and transact options.
//type ERC20Session struct {
//	Contract     *ERC20            // Generic contract binding to set the session for
//	CallOpts     bind.CallOpts     // Call options to use throughout this session
//	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
//}
//
//// ERC20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
//// with pre-set call options.
//type ERC20CallerSession struct {
//	Contract *ERC20Caller  // Generic contract caller binding to set the session for
//	CallOpts bind.CallOpts // Call options to use throughout this session
//}
//
//// ERC20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
//// with pre-set transact options.
//type ERC20TransactorSession struct {
//	Contract     *ERC20Transactor  // Generic contract transactor binding to set the session for
//	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
//}
//
//// ERC20Raw is an auto generated low-level Go binding around an Ethereum contract.
//type ERC20Raw struct {
//	Contract *ERC20 // Generic contract binding to access the raw methods on
//}
//
//// ERC20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
//type ERC20CallerRaw struct {
//	Contract *ERC20Caller // Generic read-only contract binding to access the raw methods on
//}
//
//// ERC20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
//type ERC20TransactorRaw struct {
//	Contract *ERC20Transactor // Generic write-only contract binding to access the raw methods on
//}
//
//// NewERC20 creates a new instance of ERC20, bound to a specific deployed contract.
//func NewERC20(address common.Address, backend bind.ContractBackend) (*ERC20, error) {
//	contract, err := bindERC20(address, backend, backend, backend)
//	if err != nil {
//		return nil, err
//	}
//	return &ERC20{ERC20Caller: ERC20Caller{contract: contract}, ERC20Transactor: ERC20Transactor{contract: contract}, ERC20Filterer: ERC20Filterer{contract: contract}}, nil
//}
//
//// NewERC20Caller creates a new read-only instance of ERC20, bound to a specific deployed contract.
//func NewERC20Caller(address common.Address, caller bind.ContractCaller) (*ERC20Caller, error) {
//	contract, err := bindERC20(address, caller, nil, nil)
//	if err != nil {
//		return nil, err
//	}
//	return &ERC20Caller{contract: contract}, nil
//}
//
//// NewERC20Transactor creates a new write-only instance of ERC20, bound to a specific deployed contract.
//func NewERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*ERC20Transactor, error) {
//	contract, err := bindERC20(address, nil, transactor, nil)
//	if err != nil {
//		return nil, err
//	}
//	return &ERC20Transactor{contract: contract}, nil
//}
//
//// NewERC20Filterer creates a new log filterer instance of ERC20, bound to a specific deployed contract.
//func NewERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*ERC20Filterer, error) {
//	contract, err := bindERC20(address, nil, nil, filterer)
//	if err != nil {
//		return nil, err
//	}
//	return &ERC20Filterer{contract: contract}, nil
//}
//
//// bindERC20 binds a generic wrapper to an already deployed contract.
//func bindERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
//	parsed, err := abi.JSON(strings.NewReader(ERC20ABI))
//	if err != nil {
//		return nil, err
//	}
//	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
//}
//
//// Call invokes the (constant) contract method with params as input values and
//// sets the output to result. The result type might be a single field for simple
//// returns, a slice of interfaces for anonymous returns and a struct for named
//// returns.
//func (_ERC20 *ERC20Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
//	return _ERC20.Contract.ERC20Caller.contract.Call(opts, result, method, params...)
//}
//
//// Transfer initiates a plain transaction to move funds to the contract, calling
//// its default method if one is available.
//func (_ERC20 *ERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _ERC20.Contract.ERC20Transactor.contract.Transfer(opts)
//}
//
//// Transact invokes the (paid) contract method with params as input values.
//func (_ERC20 *ERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
//	return _ERC20.Contract.ERC20Transactor.contract.Transact(opts, method, params...)
//}
//
//// Call invokes the (constant) contract method with params as input values and
//// sets the output to result. The result type might be a single field for simple
//// returns, a slice of interfaces for anonymous returns and a struct for named
//// returns.
//func (_ERC20 *ERC20CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
//	return _ERC20.Contract.contract.Call(opts, result, method, params...)
//}
//
//// Transfer initiates a plain transaction to move funds to the contract, calling
//// its default method if one is available.
//func (_ERC20 *ERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _ERC20.Contract.contract.Transfer(opts)
//}
//
//// Transact invokes the (paid) contract method with params as input values.
//func (_ERC20 *ERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
//	return _ERC20.Contract.contract.Transact(opts, method, params...)
//}
//
//// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
////
//// Solidity: function allowance(owner address, spender address) constant returns(uint256)
//func (_ERC20 *ERC20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
//	var (
//		ret0 = new(*big.Int)
//	)
//	out := ret0
//	err := _ERC20.contract.Call(opts, out, "allowance", owner, spender)
//	return *ret0, err
//}
//
//// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
////
//// Solidity: function allowance(owner address, spender address) constant returns(uint256)
//func (_ERC20 *ERC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
//	return _ERC20.Contract.Allowance(&_ERC20.CallOpts, owner, spender)
//}
//
//// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
////
//// Solidity: function allowance(owner address, spender address) constant returns(uint256)
//func (_ERC20 *ERC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
//	return _ERC20.Contract.Allowance(&_ERC20.CallOpts, owner, spender)
//}
//
//// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
////
//// Solidity: function balanceOf(who address) constant returns(uint256)
//func (_ERC20 *ERC20Caller) BalanceOf(opts *bind.CallOpts, who common.Address) (*big.Int, error) {
//	var (
//		ret0 = new(*big.Int)
//	)
//	out := ret0
//	err := _ERC20.contract.Call(opts, out, "balanceOf", who)
//	return *ret0, err
//}
//
//// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
////
//// Solidity: function balanceOf(who address) constant returns(uint256)
//func (_ERC20 *ERC20Session) BalanceOf(who common.Address) (*big.Int, error) {
//	return _ERC20.Contract.BalanceOf(&_ERC20.CallOpts, who)
//}
//
//// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
////
//// Solidity: function balanceOf(who address) constant returns(uint256)
//func (_ERC20 *ERC20CallerSession) BalanceOf(who common.Address) (*big.Int, error) {
//	return _ERC20.Contract.BalanceOf(&_ERC20.CallOpts, who)
//}
//
//// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
////
//// Solidity: function totalSupply() constant returns(uint256)
//func (_ERC20 *ERC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
//	var (
//		ret0 = new(*big.Int)
//	)
//	out := ret0
//	err := _ERC20.contract.Call(opts, out, "totalSupply")
//	return *ret0, err
//}
//
//// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
////
//// Solidity: function totalSupply() constant returns(uint256)
//func (_ERC20 *ERC20Session) TotalSupply() (*big.Int, error) {
//	return _ERC20.Contract.TotalSupply(&_ERC20.CallOpts)
//}
//
//// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
////
//// Solidity: function totalSupply() constant returns(uint256)
//func (_ERC20 *ERC20CallerSession) TotalSupply() (*big.Int, error) {
//	return _ERC20.Contract.TotalSupply(&_ERC20.CallOpts)
//}
//
//// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
////
//// Solidity: function approve(spender address, value uint256) returns(bool)
//func (_ERC20 *ERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
//	return _ERC20.contract.Transact(opts, "approve", spender, value)
//}
//
//// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
////
//// Solidity: function approve(spender address, value uint256) returns(bool)
//func (_ERC20 *ERC20Session) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
//	return _ERC20.Contract.Approve(&_ERC20.TransactOpts, spender, value)
//}
//
//// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
////
//// Solidity: function approve(spender address, value uint256) returns(bool)
//func (_ERC20 *ERC20TransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
//	return _ERC20.Contract.Approve(&_ERC20.TransactOpts, spender, value)
//}
//
//// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
////
//// Solidity: function transfer(to address, value uint256) returns(bool)
//func (_ERC20 *ERC20Transactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
//	return _ERC20.contract.Transact(opts, "transfer", to, value)
//}
//
//// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
////
//// Solidity: function transfer(to address, value uint256) returns(bool)
//func (_ERC20 *ERC20Session) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
//	return _ERC20.Contract.Transfer(&_ERC20.TransactOpts, to, value)
//}
//
//// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
////
//// Solidity: function transfer(to address, value uint256) returns(bool)
//func (_ERC20 *ERC20TransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
//	return _ERC20.Contract.Transfer(&_ERC20.TransactOpts, to, value)
//}
//
//// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
////
//// Solidity: function transferFrom(from address, to address, value uint256) returns(bool)
//func (_ERC20 *ERC20Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
//	return _ERC20.contract.Transact(opts, "transferFrom", from, to, value)
//}
//
//// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
////
//// Solidity: function transferFrom(from address, to address, value uint256) returns(bool)
//func (_ERC20 *ERC20Session) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
//	return _ERC20.Contract.TransferFrom(&_ERC20.TransactOpts, from, to, value)
//}
//
//// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
////
//// Solidity: function transferFrom(from address, to address, value uint256) returns(bool)
//func (_ERC20 *ERC20TransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
//	return _ERC20.Contract.TransferFrom(&_ERC20.TransactOpts, from, to, value)
//}
//
//// ERC20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ERC20 contract.
//type ERC20ApprovalIterator struct {
//	Event *ERC20Approval // Event containing the contract specifics and raw log
//
//	contract *bind.BoundContract // Generic contract to use for unpacking event data
//	event    string              // Event name to use for unpacking event data
//
//	logs chan types.Log        // Log channel receiving the found contract events
//	sub  ethereum.Subscription // Subscription for errors, completion and termination
//	done bool                  // Whether the subscription completed delivering logs
//	fail error                 // Occurred error to stop iteration
//}
//
//// Next advances the iterator to the subsequent event, returning whether there
//// are any more events found. In case of a retrieval or parsing error, false is
//// returned and Error() can be queried for the exact failure.
//func (it *ERC20ApprovalIterator) Next() bool {
//	// If the iterator failed, stop iterating
//	if it.fail != nil {
//		return false
//	}
//	// If the iterator completed, deliver directly whatever's available
//	if it.done {
//		select {
//		case log := <-it.logs:
//			it.Event = new(ERC20Approval)
//			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//				it.fail = err
//				return false
//			}
//			it.Event.Raw = log
//			return true
//
//		default:
//			return false
//		}
//	}
//	// Iterator still in progress, wait for either a data or an error event
//	select {
//	case log := <-it.logs:
//		it.Event = new(ERC20Approval)
//		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//			it.fail = err
//			return false
//		}
//		it.Event.Raw = log
//		return true
//
//	case err := <-it.sub.Err():
//		it.done = true
//		it.fail = err
//		return it.Next()
//	}
//}
//
//// Error returns any retrieval or parsing error occurred during filtering.
//func (it *ERC20ApprovalIterator) Error() error {
//	return it.fail
//}
//
//// Close terminates the iteration process, releasing any pending underlying
//// resources.
//func (it *ERC20ApprovalIterator) Close() error {
//	it.sub.Unsubscribe()
//	return nil
//}
//
//// ERC20Approval represents a Approval event raised by the ERC20 contract.
//type ERC20Approval struct {
//	Owner   common.Address
//	Spender common.Address
//	Value   *big.Int
//	Raw     types.Log // Blockchain specific contextual infos
//}
//
//// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
////
//// Solidity: event Approval(owner indexed address, spender indexed address, value uint256)
//func (_ERC20 *ERC20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ERC20ApprovalIterator, error) {
//
//	var ownerRule []interface{}
//	for _, ownerItem := range owner {
//		ownerRule = append(ownerRule, ownerItem)
//	}
//	var spenderRule []interface{}
//	for _, spenderItem := range spender {
//		spenderRule = append(spenderRule, spenderItem)
//	}
//
//	logs, sub, err := _ERC20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
//	if err != nil {
//		return nil, err
//	}
//	return &ERC20ApprovalIterator{contract: _ERC20.contract, event: "Approval", logs: logs, sub: sub}, nil
//}
//
//// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
////
//// Solidity: event Approval(owner indexed address, spender indexed address, value uint256)
//func (_ERC20 *ERC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {
//
//	var ownerRule []interface{}
//	for _, ownerItem := range owner {
//		ownerRule = append(ownerRule, ownerItem)
//	}
//	var spenderRule []interface{}
//	for _, spenderItem := range spender {
//		spenderRule = append(spenderRule, spenderItem)
//	}
//
//	logs, sub, err := _ERC20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
//	if err != nil {
//		return nil, err
//	}
//	return event.NewSubscription(func(quit <-chan struct{}) error {
//		defer sub.Unsubscribe()
//		for {
//			select {
//			case log := <-logs:
//				// New log arrived, parse the event and forward to the user
//				event := new(ERC20Approval)
//				if err := _ERC20.contract.UnpackLog(event, "Approval", log); err != nil {
//					return err
//				}
//				event.Raw = log
//
//				select {
//				case sink <- event:
//				case err := <-sub.Err():
//					return err
//				case <-quit:
//					return nil
//				}
//			case err := <-sub.Err():
//				return err
//			case <-quit:
//				return nil
//			}
//		}
//	}), nil
//}
//
//// ERC20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ERC20 contract.
//type ERC20TransferIterator struct {
//	Event *ERC20Transfer // Event containing the contract specifics and raw log
//
//	contract *bind.BoundContract // Generic contract to use for unpacking event data
//	event    string              // Event name to use for unpacking event data
//
//	logs chan types.Log        // Log channel receiving the found contract events
//	sub  ethereum.Subscription // Subscription for errors, completion and termination
//	done bool                  // Whether the subscription completed delivering logs
//	fail error                 // Occurred error to stop iteration
//}
//
//// Next advances the iterator to the subsequent event, returning whether there
//// are any more events found. In case of a retrieval or parsing error, false is
//// returned and Error() can be queried for the exact failure.
//func (it *ERC20TransferIterator) Next() bool {
//	// If the iterator failed, stop iterating
//	if it.fail != nil {
//		return false
//	}
//	// If the iterator completed, deliver directly whatever's available
//	if it.done {
//		select {
//		case log := <-it.logs:
//			it.Event = new(ERC20Transfer)
//			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//				it.fail = err
//				return false
//			}
//			it.Event.Raw = log
//			return true
//
//		default:
//			return false
//		}
//	}
//	// Iterator still in progress, wait for either a data or an error event
//	select {
//	case log := <-it.logs:
//		it.Event = new(ERC20Transfer)
//		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//			it.fail = err
//			return false
//		}
//		it.Event.Raw = log
//		return true
//
//	case err := <-it.sub.Err():
//		it.done = true
//		it.fail = err
//		return it.Next()
//	}
//}
//
//// Error returns any retrieval or parsing error occurred during filtering.
//func (it *ERC20TransferIterator) Error() error {
//	return it.fail
//}
//
//// Close terminates the iteration process, releasing any pending underlying
//// resources.
//func (it *ERC20TransferIterator) Close() error {
//	it.sub.Unsubscribe()
//	return nil
//}
//
//// ERC20Transfer represents a Transfer event raised by the ERC20 contract.
//type ERC20Transfer struct {
//	From  common.Address
//	To    common.Address
//	Value *big.Int
//	Raw   types.Log // Blockchain specific contextual infos
//}
//
//// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
////
//// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
//func (_ERC20 *ERC20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ERC20TransferIterator, error) {
//
//	var fromRule []interface{}
//	for _, fromItem := range from {
//		fromRule = append(fromRule, fromItem)
//	}
//	var toRule []interface{}
//	for _, toItem := range to {
//		toRule = append(toRule, toItem)
//	}
//
//	logs, sub, err := _ERC20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
//	if err != nil {
//		return nil, err
//	}
//	return &ERC20TransferIterator{contract: _ERC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
//}
//
//// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
////
//// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
//func (_ERC20 *ERC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {
//
//	var fromRule []interface{}
//	for _, fromItem := range from {
//		fromRule = append(fromRule, fromItem)
//	}
//	var toRule []interface{}
//	for _, toItem := range to {
//		toRule = append(toRule, toItem)
//	}
//
//	logs, sub, err := _ERC20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
//	if err != nil {
//		return nil, err
//	}
//	return event.NewSubscription(func(quit <-chan struct{}) error {
//		defer sub.Unsubscribe()
//		for {
//			select {
//			case log := <-logs:
//				// New log arrived, parse the event and forward to the user
//				event := new(ERC20Transfer)
//				if err := _ERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
//					return err
//				}
//				event.Raw = log
//
//				select {
//				case sink <- event:
//				case err := <-sub.Err():
//					return err
//				case <-quit:
//					return nil
//				}
//			case err := <-sub.Err():
//				return err
//			case <-quit:
//				return nil
//			}
//		}
//	}), nil
//}
//
//// ERC20BasicABI is the input ABI used to generate the binding from.
//const ERC20BasicABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"who\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]"
//
//// ERC20BasicBin is the compiled bytecode used for deploying new contracts.
//const ERC20BasicBin = `0x`
//
//// DeployERC20Basic deploys a new Ethereum contract, binding an instance of ERC20Basic to it.
//func DeployERC20Basic(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ERC20Basic, error) {
//	parsed, err := abi.JSON(strings.NewReader(ERC20BasicABI))
//	if err != nil {
//		return common.Address{}, nil, nil, err
//	}
//	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ERC20BasicBin), backend)
//	if err != nil {
//		return common.Address{}, nil, nil, err
//	}
//	return address, tx, &ERC20Basic{ERC20BasicCaller: ERC20BasicCaller{contract: contract}, ERC20BasicTransactor: ERC20BasicTransactor{contract: contract}, ERC20BasicFilterer: ERC20BasicFilterer{contract: contract}}, nil
//}
//
//// ERC20Basic is an auto generated Go binding around an Ethereum contract.
//type ERC20Basic struct {
//	ERC20BasicCaller     // Read-only binding to the contract
//	ERC20BasicTransactor // Write-only binding to the contract
//	ERC20BasicFilterer   // Log filterer for contract events
//}
//
//// ERC20BasicCaller is an auto generated read-only Go binding around an Ethereum contract.
//type ERC20BasicCaller struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// ERC20BasicTransactor is an auto generated write-only Go binding around an Ethereum contract.
//type ERC20BasicTransactor struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// ERC20BasicFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
//type ERC20BasicFilterer struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// ERC20BasicSession is an auto generated Go binding around an Ethereum contract,
//// with pre-set call and transact options.
//type ERC20BasicSession struct {
//	Contract     *ERC20Basic       // Generic contract binding to set the session for
//	CallOpts     bind.CallOpts     // Call options to use throughout this session
//	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
//}
//
//// ERC20BasicCallerSession is an auto generated read-only Go binding around an Ethereum contract,
//// with pre-set call options.
//type ERC20BasicCallerSession struct {
//	Contract *ERC20BasicCaller // Generic contract caller binding to set the session for
//	CallOpts bind.CallOpts     // Call options to use throughout this session
//}
//
//// ERC20BasicTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
//// with pre-set transact options.
//type ERC20BasicTransactorSession struct {
//	Contract     *ERC20BasicTransactor // Generic contract transactor binding to set the session for
//	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
//}
//
//// ERC20BasicRaw is an auto generated low-level Go binding around an Ethereum contract.
//type ERC20BasicRaw struct {
//	Contract *ERC20Basic // Generic contract binding to access the raw methods on
//}
//
//// ERC20BasicCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
//type ERC20BasicCallerRaw struct {
//	Contract *ERC20BasicCaller // Generic read-only contract binding to access the raw methods on
//}
//
//// ERC20BasicTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
//type ERC20BasicTransactorRaw struct {
//	Contract *ERC20BasicTransactor // Generic write-only contract binding to access the raw methods on
//}
//
//// NewERC20Basic creates a new instance of ERC20Basic, bound to a specific deployed contract.
//func NewERC20Basic(address common.Address, backend bind.ContractBackend) (*ERC20Basic, error) {
//	contract, err := bindERC20Basic(address, backend, backend, backend)
//	if err != nil {
//		return nil, err
//	}
//	return &ERC20Basic{ERC20BasicCaller: ERC20BasicCaller{contract: contract}, ERC20BasicTransactor: ERC20BasicTransactor{contract: contract}, ERC20BasicFilterer: ERC20BasicFilterer{contract: contract}}, nil
//}
//
//// NewERC20BasicCaller creates a new read-only instance of ERC20Basic, bound to a specific deployed contract.
//func NewERC20BasicCaller(address common.Address, caller bind.ContractCaller) (*ERC20BasicCaller, error) {
//	contract, err := bindERC20Basic(address, caller, nil, nil)
//	if err != nil {
//		return nil, err
//	}
//	return &ERC20BasicCaller{contract: contract}, nil
//}
//
//// NewERC20BasicTransactor creates a new write-only instance of ERC20Basic, bound to a specific deployed contract.
//func NewERC20BasicTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC20BasicTransactor, error) {
//	contract, err := bindERC20Basic(address, nil, transactor, nil)
//	if err != nil {
//		return nil, err
//	}
//	return &ERC20BasicTransactor{contract: contract}, nil
//}
//
//// NewERC20BasicFilterer creates a new log filterer instance of ERC20Basic, bound to a specific deployed contract.
//func NewERC20BasicFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC20BasicFilterer, error) {
//	contract, err := bindERC20Basic(address, nil, nil, filterer)
//	if err != nil {
//		return nil, err
//	}
//	return &ERC20BasicFilterer{contract: contract}, nil
//}
//
//// bindERC20Basic binds a generic wrapper to an already deployed contract.
//func bindERC20Basic(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
//	parsed, err := abi.JSON(strings.NewReader(ERC20BasicABI))
//	if err != nil {
//		return nil, err
//	}
//	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
//}
//
//// Call invokes the (constant) contract method with params as input values and
//// sets the output to result. The result type might be a single field for simple
//// returns, a slice of interfaces for anonymous returns and a struct for named
//// returns.
//func (_ERC20Basic *ERC20BasicRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
//	return _ERC20Basic.Contract.ERC20BasicCaller.contract.Call(opts, result, method, params...)
//}
//
//// Transfer initiates a plain transaction to move funds to the contract, calling
//// its default method if one is available.
//func (_ERC20Basic *ERC20BasicRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _ERC20Basic.Contract.ERC20BasicTransactor.contract.Transfer(opts)
//}
//
//// Transact invokes the (paid) contract method with params as input values.
//func (_ERC20Basic *ERC20BasicRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
//	return _ERC20Basic.Contract.ERC20BasicTransactor.contract.Transact(opts, method, params...)
//}
//
//// Call invokes the (constant) contract method with params as input values and
//// sets the output to result. The result type might be a single field for simple
//// returns, a slice of interfaces for anonymous returns and a struct for named
//// returns.
//func (_ERC20Basic *ERC20BasicCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
//	return _ERC20Basic.Contract.contract.Call(opts, result, method, params...)
//}
//
//// Transfer initiates a plain transaction to move funds to the contract, calling
//// its default method if one is available.
//func (_ERC20Basic *ERC20BasicTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _ERC20Basic.Contract.contract.Transfer(opts)
//}
//
//// Transact invokes the (paid) contract method with params as input values.
//func (_ERC20Basic *ERC20BasicTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
//	return _ERC20Basic.Contract.contract.Transact(opts, method, params...)
//}
//
//// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
////
//// Solidity: function balanceOf(who address) constant returns(uint256)
//func (_ERC20Basic *ERC20BasicCaller) BalanceOf(opts *bind.CallOpts, who common.Address) (*big.Int, error) {
//	var (
//		ret0 = new(*big.Int)
//	)
//	out := ret0
//	err := _ERC20Basic.contract.Call(opts, out, "balanceOf", who)
//	return *ret0, err
//}
//
//// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
////
//// Solidity: function balanceOf(who address) constant returns(uint256)
//func (_ERC20Basic *ERC20BasicSession) BalanceOf(who common.Address) (*big.Int, error) {
//	return _ERC20Basic.Contract.BalanceOf(&_ERC20Basic.CallOpts, who)
//}
//
//// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
////
//// Solidity: function balanceOf(who address) constant returns(uint256)
//func (_ERC20Basic *ERC20BasicCallerSession) BalanceOf(who common.Address) (*big.Int, error) {
//	return _ERC20Basic.Contract.BalanceOf(&_ERC20Basic.CallOpts, who)
//}
//
//// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
////
//// Solidity: function totalSupply() constant returns(uint256)
//func (_ERC20Basic *ERC20BasicCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
//	var (
//		ret0 = new(*big.Int)
//	)
//	out := ret0
//	err := _ERC20Basic.contract.Call(opts, out, "totalSupply")
//	return *ret0, err
//}
//
//// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
////
//// Solidity: function totalSupply() constant returns(uint256)
//func (_ERC20Basic *ERC20BasicSession) TotalSupply() (*big.Int, error) {
//	return _ERC20Basic.Contract.TotalSupply(&_ERC20Basic.CallOpts)
//}
//
//// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
////
//// Solidity: function totalSupply() constant returns(uint256)
//func (_ERC20Basic *ERC20BasicCallerSession) TotalSupply() (*big.Int, error) {
//	return _ERC20Basic.Contract.TotalSupply(&_ERC20Basic.CallOpts)
//}
//
//// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
////
//// Solidity: function transfer(to address, value uint256) returns(bool)
//func (_ERC20Basic *ERC20BasicTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
//	return _ERC20Basic.contract.Transact(opts, "transfer", to, value)
//}
//
//// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
////
//// Solidity: function transfer(to address, value uint256) returns(bool)
//func (_ERC20Basic *ERC20BasicSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
//	return _ERC20Basic.Contract.Transfer(&_ERC20Basic.TransactOpts, to, value)
//}
//
//// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
////
//// Solidity: function transfer(to address, value uint256) returns(bool)
//func (_ERC20Basic *ERC20BasicTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
//	return _ERC20Basic.Contract.Transfer(&_ERC20Basic.TransactOpts, to, value)
//}
//
//// ERC20BasicTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ERC20Basic contract.
//type ERC20BasicTransferIterator struct {
//	Event *ERC20BasicTransfer // Event containing the contract specifics and raw log
//
//	contract *bind.BoundContract // Generic contract to use for unpacking event data
//	event    string              // Event name to use for unpacking event data
//
//	logs chan types.Log        // Log channel receiving the found contract events
//	sub  ethereum.Subscription // Subscription for errors, completion and termination
//	done bool                  // Whether the subscription completed delivering logs
//	fail error                 // Occurred error to stop iteration
//}
//
//// Next advances the iterator to the subsequent event, returning whether there
//// are any more events found. In case of a retrieval or parsing error, false is
//// returned and Error() can be queried for the exact failure.
//func (it *ERC20BasicTransferIterator) Next() bool {
//	// If the iterator failed, stop iterating
//	if it.fail != nil {
//		return false
//	}
//	// If the iterator completed, deliver directly whatever's available
//	if it.done {
//		select {
//		case log := <-it.logs:
//			it.Event = new(ERC20BasicTransfer)
//			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//				it.fail = err
//				return false
//			}
//			it.Event.Raw = log
//			return true
//
//		default:
//			return false
//		}
//	}
//	// Iterator still in progress, wait for either a data or an error event
//	select {
//	case log := <-it.logs:
//		it.Event = new(ERC20BasicTransfer)
//		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//			it.fail = err
//			return false
//		}
//		it.Event.Raw = log
//		return true
//
//	case err := <-it.sub.Err():
//		it.done = true
//		it.fail = err
//		return it.Next()
//	}
//}
//
//// Error returns any retrieval or parsing error occurred during filtering.
//func (it *ERC20BasicTransferIterator) Error() error {
//	return it.fail
//}
//
//// Close terminates the iteration process, releasing any pending underlying
//// resources.
//func (it *ERC20BasicTransferIterator) Close() error {
//	it.sub.Unsubscribe()
//	return nil
//}
//
//// ERC20BasicTransfer represents a Transfer event raised by the ERC20Basic contract.
//type ERC20BasicTransfer struct {
//	From  common.Address
//	To    common.Address
//	Value *big.Int
//	Raw   types.Log // Blockchain specific contextual infos
//}
//
//// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
////
//// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
//func (_ERC20Basic *ERC20BasicFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ERC20BasicTransferIterator, error) {
//
//	var fromRule []interface{}
//	for _, fromItem := range from {
//		fromRule = append(fromRule, fromItem)
//	}
//	var toRule []interface{}
//	for _, toItem := range to {
//		toRule = append(toRule, toItem)
//	}
//
//	logs, sub, err := _ERC20Basic.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
//	if err != nil {
//		return nil, err
//	}
//	return &ERC20BasicTransferIterator{contract: _ERC20Basic.contract, event: "Transfer", logs: logs, sub: sub}, nil
//}
//
//// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
////
//// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
//func (_ERC20Basic *ERC20BasicFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ERC20BasicTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {
//
//	var fromRule []interface{}
//	for _, fromItem := range from {
//		fromRule = append(fromRule, fromItem)
//	}
//	var toRule []interface{}
//	for _, toItem := range to {
//		toRule = append(toRule, toItem)
//	}
//
//	logs, sub, err := _ERC20Basic.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
//	if err != nil {
//		return nil, err
//	}
//	return event.NewSubscription(func(quit <-chan struct{}) error {
//		defer sub.Unsubscribe()
//		for {
//			select {
//			case log := <-logs:
//				// New log arrived, parse the event and forward to the user
//				event := new(ERC20BasicTransfer)
//				if err := _ERC20Basic.contract.UnpackLog(event, "Transfer", log); err != nil {
//					return err
//				}
//				event.Raw = log
//
//				select {
//				case sink <- event:
//				case err := <-sub.Err():
//					return err
//				case <-quit:
//					return nil
//				}
//			case err := <-sub.Err():
//				return err
//			case <-quit:
//				return nil
//			}
//		}
//	}), nil
//}
//
//// LinkedListABI is the input ABI used to generate the binding from.
//const LinkedListABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"NULL\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes20\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"
//
//// LinkedListBin is the compiled bytecode used for deploying new contracts.
//const LinkedListBin = `0x60b361002f600b82828239805160001a6073146000811461001f57610021565bfe5b5030600052607381538281f300730000000000000000000000000000000000000000301460606040526004361060555763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663f26be3fc8114605a575b600080fd5b60606082565b6040516bffffffffffffffffffffffff19909116815260200160405180910390f35b6000815600a165627a7a72305820c95127cf3dfd959939a4479716794968359d8c70b5fef2984f000ebc4db36e310029`
//
//// DeployLinkedList deploys a new Ethereum contract, binding an instance of LinkedList to it.
//func DeployLinkedList(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *LinkedList, error) {
//	parsed, err := abi.JSON(strings.NewReader(LinkedListABI))
//	if err != nil {
//		return common.Address{}, nil, nil, err
//	}
//	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(LinkedListBin), backend)
//	if err != nil {
//		return common.Address{}, nil, nil, err
//	}
//	return address, tx, &LinkedList{LinkedListCaller: LinkedListCaller{contract: contract}, LinkedListTransactor: LinkedListTransactor{contract: contract}, LinkedListFilterer: LinkedListFilterer{contract: contract}}, nil
//}
//
//// LinkedList is an auto generated Go binding around an Ethereum contract.
//type LinkedList struct {
//	LinkedListCaller     // Read-only binding to the contract
//	LinkedListTransactor // Write-only binding to the contract
//	LinkedListFilterer   // Log filterer for contract events
//}
//
//// LinkedListCaller is an auto generated read-only Go binding around an Ethereum contract.
//type LinkedListCaller struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// LinkedListTransactor is an auto generated write-only Go binding around an Ethereum contract.
//type LinkedListTransactor struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// LinkedListFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
//type LinkedListFilterer struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// LinkedListSession is an auto generated Go binding around an Ethereum contract,
//// with pre-set call and transact options.
//type LinkedListSession struct {
//	Contract     *LinkedList       // Generic contract binding to set the session for
//	CallOpts     bind.CallOpts     // Call options to use throughout this session
//	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
//}
//
//// LinkedListCallerSession is an auto generated read-only Go binding around an Ethereum contract,
//// with pre-set call options.
//type LinkedListCallerSession struct {
//	Contract *LinkedListCaller // Generic contract caller binding to set the session for
//	CallOpts bind.CallOpts     // Call options to use throughout this session
//}
//
//// LinkedListTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
//// with pre-set transact options.
//type LinkedListTransactorSession struct {
//	Contract     *LinkedListTransactor // Generic contract transactor binding to set the session for
//	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
//}
//
//// LinkedListRaw is an auto generated low-level Go binding around an Ethereum contract.
//type LinkedListRaw struct {
//	Contract *LinkedList // Generic contract binding to access the raw methods on
//}
//
//// LinkedListCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
//type LinkedListCallerRaw struct {
//	Contract *LinkedListCaller // Generic read-only contract binding to access the raw methods on
//}
//
//// LinkedListTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
//type LinkedListTransactorRaw struct {
//	Contract *LinkedListTransactor // Generic write-only contract binding to access the raw methods on
//}
//
//// NewLinkedList creates a new instance of LinkedList, bound to a specific deployed contract.
//func NewLinkedList(address common.Address, backend bind.ContractBackend) (*LinkedList, error) {
//	contract, err := bindLinkedList(address, backend, backend, backend)
//	if err != nil {
//		return nil, err
//	}
//	return &LinkedList{LinkedListCaller: LinkedListCaller{contract: contract}, LinkedListTransactor: LinkedListTransactor{contract: contract}, LinkedListFilterer: LinkedListFilterer{contract: contract}}, nil
//}
//
//// NewLinkedListCaller creates a new read-only instance of LinkedList, bound to a specific deployed contract.
//func NewLinkedListCaller(address common.Address, caller bind.ContractCaller) (*LinkedListCaller, error) {
//	contract, err := bindLinkedList(address, caller, nil, nil)
//	if err != nil {
//		return nil, err
//	}
//	return &LinkedListCaller{contract: contract}, nil
//}
//
//// NewLinkedListTransactor creates a new write-only instance of LinkedList, bound to a specific deployed contract.
//func NewLinkedListTransactor(address common.Address, transactor bind.ContractTransactor) (*LinkedListTransactor, error) {
//	contract, err := bindLinkedList(address, nil, transactor, nil)
//	if err != nil {
//		return nil, err
//	}
//	return &LinkedListTransactor{contract: contract}, nil
//}
//
//// NewLinkedListFilterer creates a new log filterer instance of LinkedList, bound to a specific deployed contract.
//func NewLinkedListFilterer(address common.Address, filterer bind.ContractFilterer) (*LinkedListFilterer, error) {
//	contract, err := bindLinkedList(address, nil, nil, filterer)
//	if err != nil {
//		return nil, err
//	}
//	return &LinkedListFilterer{contract: contract}, nil
//}
//
//// bindLinkedList binds a generic wrapper to an already deployed contract.
//func bindLinkedList(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
//	parsed, err := abi.JSON(strings.NewReader(LinkedListABI))
//	if err != nil {
//		return nil, err
//	}
//	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
//}
//
//// Call invokes the (constant) contract method with params as input values and
//// sets the output to result. The result type might be a single field for simple
//// returns, a slice of interfaces for anonymous returns and a struct for named
//// returns.
//func (_LinkedList *LinkedListRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
//	return _LinkedList.Contract.LinkedListCaller.contract.Call(opts, result, method, params...)
//}
//
//// Transfer initiates a plain transaction to move funds to the contract, calling
//// its default method if one is available.
//func (_LinkedList *LinkedListRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _LinkedList.Contract.LinkedListTransactor.contract.Transfer(opts)
//}
//
//// Transact invokes the (paid) contract method with params as input values.
//func (_LinkedList *LinkedListRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
//	return _LinkedList.Contract.LinkedListTransactor.contract.Transact(opts, method, params...)
//}
//
//// Call invokes the (constant) contract method with params as input values and
//// sets the output to result. The result type might be a single field for simple
//// returns, a slice of interfaces for anonymous returns and a struct for named
//// returns.
//func (_LinkedList *LinkedListCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
//	return _LinkedList.Contract.contract.Call(opts, result, method, params...)
//}
//
//// Transfer initiates a plain transaction to move funds to the contract, calling
//// its default method if one is available.
//func (_LinkedList *LinkedListTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _LinkedList.Contract.contract.Transfer(opts)
//}
//
//// Transact invokes the (paid) contract method with params as input values.
//func (_LinkedList *LinkedListTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
//	return _LinkedList.Contract.contract.Transact(opts, method, params...)
//}
//
//// NULL is a free data retrieval call binding the contract method 0xf26be3fc.
////
//// Solidity: function NULL() constant returns(bytes20)
//func (_LinkedList *LinkedListCaller) NULL(opts *bind.CallOpts) ([20]byte, error) {
//	var (
//		ret0 = new([20]byte)
//	)
//	out := ret0
//	err := _LinkedList.contract.Call(opts, out, "NULL")
//	return *ret0, err
//}
//
//// NULL is a free data retrieval call binding the contract method 0xf26be3fc.
////
//// Solidity: function NULL() constant returns(bytes20)
//func (_LinkedList *LinkedListSession) NULL() ([20]byte, error) {
//	return _LinkedList.Contract.NULL(&_LinkedList.CallOpts)
//}
//
//// NULL is a free data retrieval call binding the contract method 0xf26be3fc.
////
//// Solidity: function NULL() constant returns(bytes20)
//func (_LinkedList *LinkedListCallerSession) NULL() ([20]byte, error) {
//	return _LinkedList.Contract.NULL(&_LinkedList.CallOpts)
//}
//
//// OwnableABI is the input ABI used to generate the binding from.
//const OwnableABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"
//
//// OwnableBin is the compiled bytecode used for deploying new contracts.
//const OwnableBin = `0x6060604052341561000f57600080fd5b60008054600160a060020a033316600160a060020a03199091161790556101768061003b6000396000f30060606040526004361061004b5763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416638da5cb5b8114610050578063f2fde38b1461007f575b600080fd5b341561005b57600080fd5b6100636100a0565b604051600160a060020a03909116815260200160405180910390f35b341561008a57600080fd5b61009e600160a060020a03600435166100af565b005b600054600160a060020a031681565b60005433600160a060020a039081169116146100ca57600080fd5b600160a060020a03811615156100df57600080fd5b600054600160a060020a0380831691167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a36000805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03929092169190911790555600a165627a7a72305820427be689d6ac24d9af13fc6d863519b0bb1362dcac6c786fe3bffd0777f16de10029`
//
//// DeployOwnable deploys a new Ethereum contract, binding an instance of Ownable to it.
//func DeployOwnable(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Ownable, error) {
//	parsed, err := abi.JSON(strings.NewReader(OwnableABI))
//	if err != nil {
//		return common.Address{}, nil, nil, err
//	}
//	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(OwnableBin), backend)
//	if err != nil {
//		return common.Address{}, nil, nil, err
//	}
//	return address, tx, &Ownable{OwnableCaller: OwnableCaller{contract: contract}, OwnableTransactor: OwnableTransactor{contract: contract}, OwnableFilterer: OwnableFilterer{contract: contract}}, nil
//}
//
//// Ownable is an auto generated Go binding around an Ethereum contract.
//type Ownable struct {
//	OwnableCaller     // Read-only binding to the contract
//	OwnableTransactor // Write-only binding to the contract
//	OwnableFilterer   // Log filterer for contract events
//}
//
//// OwnableCaller is an auto generated read-only Go binding around an Ethereum contract.
//type OwnableCaller struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// OwnableTransactor is an auto generated write-only Go binding around an Ethereum contract.
//type OwnableTransactor struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// OwnableFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
//type OwnableFilterer struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// OwnableSession is an auto generated Go binding around an Ethereum contract,
//// with pre-set call and transact options.
//type OwnableSession struct {
//	Contract     *Ownable          // Generic contract binding to set the session for
//	CallOpts     bind.CallOpts     // Call options to use throughout this session
//	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
//}
//
//// OwnableCallerSession is an auto generated read-only Go binding around an Ethereum contract,
//// with pre-set call options.
//type OwnableCallerSession struct {
//	Contract *OwnableCaller // Generic contract caller binding to set the session for
//	CallOpts bind.CallOpts  // Call options to use throughout this session
//}
//
//// OwnableTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
//// with pre-set transact options.
//type OwnableTransactorSession struct {
//	Contract     *OwnableTransactor // Generic contract transactor binding to set the session for
//	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
//}
//
//// OwnableRaw is an auto generated low-level Go binding around an Ethereum contract.
//type OwnableRaw struct {
//	Contract *Ownable // Generic contract binding to access the raw methods on
//}
//
//// OwnableCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
//type OwnableCallerRaw struct {
//	Contract *OwnableCaller // Generic read-only contract binding to access the raw methods on
//}
//
//// OwnableTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
//type OwnableTransactorRaw struct {
//	Contract *OwnableTransactor // Generic write-only contract binding to access the raw methods on
//}
//
//// NewOwnable creates a new instance of Ownable, bound to a specific deployed contract.
//func NewOwnable(address common.Address, backend bind.ContractBackend) (*Ownable, error) {
//	contract, err := bindOwnable(address, backend, backend, backend)
//	if err != nil {
//		return nil, err
//	}
//	return &Ownable{OwnableCaller: OwnableCaller{contract: contract}, OwnableTransactor: OwnableTransactor{contract: contract}, OwnableFilterer: OwnableFilterer{contract: contract}}, nil
//}
//
//// NewOwnableCaller creates a new read-only instance of Ownable, bound to a specific deployed contract.
//func NewOwnableCaller(address common.Address, caller bind.ContractCaller) (*OwnableCaller, error) {
//	contract, err := bindOwnable(address, caller, nil, nil)
//	if err != nil {
//		return nil, err
//	}
//	return &OwnableCaller{contract: contract}, nil
//}
//
//// NewOwnableTransactor creates a new write-only instance of Ownable, bound to a specific deployed contract.
//func NewOwnableTransactor(address common.Address, transactor bind.ContractTransactor) (*OwnableTransactor, error) {
//	contract, err := bindOwnable(address, nil, transactor, nil)
//	if err != nil {
//		return nil, err
//	}
//	return &OwnableTransactor{contract: contract}, nil
//}
//
//// NewOwnableFilterer creates a new log filterer instance of Ownable, bound to a specific deployed contract.
//func NewOwnableFilterer(address common.Address, filterer bind.ContractFilterer) (*OwnableFilterer, error) {
//	contract, err := bindOwnable(address, nil, nil, filterer)
//	if err != nil {
//		return nil, err
//	}
//	return &OwnableFilterer{contract: contract}, nil
//}
//
//// bindOwnable binds a generic wrapper to an already deployed contract.
//func bindOwnable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
//	parsed, err := abi.JSON(strings.NewReader(OwnableABI))
//	if err != nil {
//		return nil, err
//	}
//	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
//}
//
//// Call invokes the (constant) contract method with params as input values and
//// sets the output to result. The result type might be a single field for simple
//// returns, a slice of interfaces for anonymous returns and a struct for named
//// returns.
//func (_Ownable *OwnableRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
//	return _Ownable.Contract.OwnableCaller.contract.Call(opts, result, method, params...)
//}
//
//// Transfer initiates a plain transaction to move funds to the contract, calling
//// its default method if one is available.
//func (_Ownable *OwnableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _Ownable.Contract.OwnableTransactor.contract.Transfer(opts)
//}
//
//// Transact invokes the (paid) contract method with params as input values.
//func (_Ownable *OwnableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
//	return _Ownable.Contract.OwnableTransactor.contract.Transact(opts, method, params...)
//}
//
//// Call invokes the (constant) contract method with params as input values and
//// sets the output to result. The result type might be a single field for simple
//// returns, a slice of interfaces for anonymous returns and a struct for named
//// returns.
//func (_Ownable *OwnableCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
//	return _Ownable.Contract.contract.Call(opts, result, method, params...)
//}
//
//// Transfer initiates a plain transaction to move funds to the contract, calling
//// its default method if one is available.
//func (_Ownable *OwnableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _Ownable.Contract.contract.Transfer(opts)
//}
//
//// Transact invokes the (paid) contract method with params as input values.
//func (_Ownable *OwnableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
//	return _Ownable.Contract.contract.Transact(opts, method, params...)
//}
//
//// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
////
//// Solidity: function owner() constant returns(address)
//func (_Ownable *OwnableCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
//	var (
//		ret0 = new(common.Address)
//	)
//	out := ret0
//	err := _Ownable.contract.Call(opts, out, "owner")
//	return *ret0, err
//}
//
//// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
////
//// Solidity: function owner() constant returns(address)
//func (_Ownable *OwnableSession) Owner() (common.Address, error) {
//	return _Ownable.Contract.Owner(&_Ownable.CallOpts)
//}
//
//// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
////
//// Solidity: function owner() constant returns(address)
//func (_Ownable *OwnableCallerSession) Owner() (common.Address, error) {
//	return _Ownable.Contract.Owner(&_Ownable.CallOpts)
//}
//
//// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
////
//// Solidity: function transferOwnership(newOwner address) returns()
//func (_Ownable *OwnableTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
//	return _Ownable.contract.Transact(opts, "transferOwnership", newOwner)
//}
//
//// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
////
//// Solidity: function transferOwnership(newOwner address) returns()
//func (_Ownable *OwnableSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
//	return _Ownable.Contract.TransferOwnership(&_Ownable.TransactOpts, newOwner)
//}
//
//// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
////
//// Solidity: function transferOwnership(newOwner address) returns()
//func (_Ownable *OwnableTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
//	return _Ownable.Contract.TransferOwnership(&_Ownable.TransactOpts, newOwner)
//}
//
//// OwnableOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Ownable contract.
//type OwnableOwnershipTransferredIterator struct {
//	Event *OwnableOwnershipTransferred // Event containing the contract specifics and raw log
//
//	contract *bind.BoundContract // Generic contract to use for unpacking event data
//	event    string              // Event name to use for unpacking event data
//
//	logs chan types.Log        // Log channel receiving the found contract events
//	sub  ethereum.Subscription // Subscription for errors, completion and termination
//	done bool                  // Whether the subscription completed delivering logs
//	fail error                 // Occurred error to stop iteration
//}
//
//// Next advances the iterator to the subsequent event, returning whether there
//// are any more events found. In case of a retrieval or parsing error, false is
//// returned and Error() can be queried for the exact failure.
//func (it *OwnableOwnershipTransferredIterator) Next() bool {
//	// If the iterator failed, stop iterating
//	if it.fail != nil {
//		return false
//	}
//	// If the iterator completed, deliver directly whatever's available
//	if it.done {
//		select {
//		case log := <-it.logs:
//			it.Event = new(OwnableOwnershipTransferred)
//			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//				it.fail = err
//				return false
//			}
//			it.Event.Raw = log
//			return true
//
//		default:
//			return false
//		}
//	}
//	// Iterator still in progress, wait for either a data or an error event
//	select {
//	case log := <-it.logs:
//		it.Event = new(OwnableOwnershipTransferred)
//		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//			it.fail = err
//			return false
//		}
//		it.Event.Raw = log
//		return true
//
//	case err := <-it.sub.Err():
//		it.done = true
//		it.fail = err
//		return it.Next()
//	}
//}
//
//// Error returns any retrieval or parsing error occurred during filtering.
//func (it *OwnableOwnershipTransferredIterator) Error() error {
//	return it.fail
//}
//
//// Close terminates the iteration process, releasing any pending underlying
//// resources.
//func (it *OwnableOwnershipTransferredIterator) Close() error {
//	it.sub.Unsubscribe()
//	return nil
//}
//
//// OwnableOwnershipTransferred represents a OwnershipTransferred event raised by the Ownable contract.
//type OwnableOwnershipTransferred struct {
//	PreviousOwner common.Address
//	NewOwner      common.Address
//	Raw           types.Log // Blockchain specific contextual infos
//}
//
//// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
////
//// Solidity: event OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
//func (_Ownable *OwnableFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*OwnableOwnershipTransferredIterator, error) {
//
//	var previousOwnerRule []interface{}
//	for _, previousOwnerItem := range previousOwner {
//		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
//	}
//	var newOwnerRule []interface{}
//	for _, newOwnerItem := range newOwner {
//		newOwnerRule = append(newOwnerRule, newOwnerItem)
//	}
//
//	logs, sub, err := _Ownable.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
//	if err != nil {
//		return nil, err
//	}
//	return &OwnableOwnershipTransferredIterator{contract: _Ownable.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
//}
//
//// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
////
//// Solidity: event OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
//func (_Ownable *OwnableFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OwnableOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {
//
//	var previousOwnerRule []interface{}
//	for _, previousOwnerItem := range previousOwner {
//		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
//	}
//	var newOwnerRule []interface{}
//	for _, newOwnerItem := range newOwner {
//		newOwnerRule = append(newOwnerRule, newOwnerItem)
//	}
//
//	logs, sub, err := _Ownable.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
//	if err != nil {
//		return nil, err
//	}
//	return event.NewSubscription(func(quit <-chan struct{}) error {
//		defer sub.Unsubscribe()
//		for {
//			select {
//			case log := <-logs:
//				// New log arrived, parse the event and forward to the user
//				event := new(OwnableOwnershipTransferred)
//				if err := _Ownable.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
//					return err
//				}
//				event.Raw = log
//
//				select {
//				case sink <- event:
//				case err := <-sub.Err():
//					return err
//				case <-quit:
//					return nil
//				}
//			case err := <-sub.Err():
//				return err
//			case <-quit:
//				return nil
//			}
//		}
//	}), nil
//}
//
//// PausableABI is the input ABI used to generate the binding from.
//const PausableABI = "[{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"
//
//// PausableBin is the compiled bytecode used for deploying new contracts.
//const PausableBin = `0x606060405260008054600160a060020a033316600160a860020a031990911617905561033b806100306000396000f30060606040526004361061006c5763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416633f4ba83a81146100715780635c975abb146100865780638456cb59146100ad5780638da5cb5b146100c0578063f2fde38b146100ef575b600080fd5b341561007c57600080fd5b61008461010e565b005b341561009157600080fd5b61009961019e565b604051901515815260200160405180910390f35b34156100b857600080fd5b6100846101bf565b34156100cb57600080fd5b6100d3610265565b604051600160a060020a03909116815260200160405180910390f35b34156100fa57600080fd5b610084600160a060020a0360043516610274565b60005433600160a060020a0390811691161461012957600080fd5b60005474010000000000000000000000000000000000000000900460ff16151561015257600080fd5b6000805474ff0000000000000000000000000000000000000000191690557f7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b3360405160405180910390a1565b60005474010000000000000000000000000000000000000000900460ff1681565b60005433600160a060020a039081169116146101da57600080fd5b60005474010000000000000000000000000000000000000000900460ff161561020257600080fd5b6000805474ff00000000000000000000000000000000000000001916740100000000000000000000000000000000000000001790557f6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff62560405160405180910390a1565b600054600160a060020a031681565b60005433600160a060020a0390811691161461028f57600080fd5b600160a060020a03811615156102a457600080fd5b600054600160a060020a0380831691167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a36000805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03929092169190911790555600a165627a7a7230582022ab7b14390dcdcdf2219ee80d05a5fca41c657e4b3906b910b9fa3f5bd349580029`
//
//// DeployPausable deploys a new Ethereum contract, binding an instance of Pausable to it.
//func DeployPausable(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Pausable, error) {
//	parsed, err := abi.JSON(strings.NewReader(PausableABI))
//	if err != nil {
//		return common.Address{}, nil, nil, err
//	}
//	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(PausableBin), backend)
//	if err != nil {
//		return common.Address{}, nil, nil, err
//	}
//	return address, tx, &Pausable{PausableCaller: PausableCaller{contract: contract}, PausableTransactor: PausableTransactor{contract: contract}, PausableFilterer: PausableFilterer{contract: contract}}, nil
//}
//
//// Pausable is an auto generated Go binding around an Ethereum contract.
//type Pausable struct {
//	PausableCaller     // Read-only binding to the contract
//	PausableTransactor // Write-only binding to the contract
//	PausableFilterer   // Log filterer for contract events
//}
//
//// PausableCaller is an auto generated read-only Go binding around an Ethereum contract.
//type PausableCaller struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// PausableTransactor is an auto generated write-only Go binding around an Ethereum contract.
//type PausableTransactor struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// PausableFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
//type PausableFilterer struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// PausableSession is an auto generated Go binding around an Ethereum contract,
//// with pre-set call and transact options.
//type PausableSession struct {
//	Contract     *Pausable         // Generic contract binding to set the session for
//	CallOpts     bind.CallOpts     // Call options to use throughout this session
//	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
//}
//
//// PausableCallerSession is an auto generated read-only Go binding around an Ethereum contract,
//// with pre-set call options.
//type PausableCallerSession struct {
//	Contract *PausableCaller // Generic contract caller binding to set the session for
//	CallOpts bind.CallOpts   // Call options to use throughout this session
//}
//
//// PausableTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
//// with pre-set transact options.
//type PausableTransactorSession struct {
//	Contract     *PausableTransactor // Generic contract transactor binding to set the session for
//	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
//}
//
//// PausableRaw is an auto generated low-level Go binding around an Ethereum contract.
//type PausableRaw struct {
//	Contract *Pausable // Generic contract binding to access the raw methods on
//}
//
//// PausableCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
//type PausableCallerRaw struct {
//	Contract *PausableCaller // Generic read-only contract binding to access the raw methods on
//}
//
//// PausableTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
//type PausableTransactorRaw struct {
//	Contract *PausableTransactor // Generic write-only contract binding to access the raw methods on
//}
//
//// NewPausable creates a new instance of Pausable, bound to a specific deployed contract.
//func NewPausable(address common.Address, backend bind.ContractBackend) (*Pausable, error) {
//	contract, err := bindPausable(address, backend, backend, backend)
//	if err != nil {
//		return nil, err
//	}
//	return &Pausable{PausableCaller: PausableCaller{contract: contract}, PausableTransactor: PausableTransactor{contract: contract}, PausableFilterer: PausableFilterer{contract: contract}}, nil
//}
//
//// NewPausableCaller creates a new read-only instance of Pausable, bound to a specific deployed contract.
//func NewPausableCaller(address common.Address, caller bind.ContractCaller) (*PausableCaller, error) {
//	contract, err := bindPausable(address, caller, nil, nil)
//	if err != nil {
//		return nil, err
//	}
//	return &PausableCaller{contract: contract}, nil
//}
//
//// NewPausableTransactor creates a new write-only instance of Pausable, bound to a specific deployed contract.
//func NewPausableTransactor(address common.Address, transactor bind.ContractTransactor) (*PausableTransactor, error) {
//	contract, err := bindPausable(address, nil, transactor, nil)
//	if err != nil {
//		return nil, err
//	}
//	return &PausableTransactor{contract: contract}, nil
//}
//
//// NewPausableFilterer creates a new log filterer instance of Pausable, bound to a specific deployed contract.
//func NewPausableFilterer(address common.Address, filterer bind.ContractFilterer) (*PausableFilterer, error) {
//	contract, err := bindPausable(address, nil, nil, filterer)
//	if err != nil {
//		return nil, err
//	}
//	return &PausableFilterer{contract: contract}, nil
//}
//
//// bindPausable binds a generic wrapper to an already deployed contract.
//func bindPausable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
//	parsed, err := abi.JSON(strings.NewReader(PausableABI))
//	if err != nil {
//		return nil, err
//	}
//	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
//}
//
//// Call invokes the (constant) contract method with params as input values and
//// sets the output to result. The result type might be a single field for simple
//// returns, a slice of interfaces for anonymous returns and a struct for named
//// returns.
//func (_Pausable *PausableRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
//	return _Pausable.Contract.PausableCaller.contract.Call(opts, result, method, params...)
//}
//
//// Transfer initiates a plain transaction to move funds to the contract, calling
//// its default method if one is available.
//func (_Pausable *PausableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _Pausable.Contract.PausableTransactor.contract.Transfer(opts)
//}
//
//// Transact invokes the (paid) contract method with params as input values.
//func (_Pausable *PausableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
//	return _Pausable.Contract.PausableTransactor.contract.Transact(opts, method, params...)
//}
//
//// Call invokes the (constant) contract method with params as input values and
//// sets the output to result. The result type might be a single field for simple
//// returns, a slice of interfaces for anonymous returns and a struct for named
//// returns.
//func (_Pausable *PausableCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
//	return _Pausable.Contract.contract.Call(opts, result, method, params...)
//}
//
//// Transfer initiates a plain transaction to move funds to the contract, calling
//// its default method if one is available.
//func (_Pausable *PausableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _Pausable.Contract.contract.Transfer(opts)
//}
//
//// Transact invokes the (paid) contract method with params as input values.
//func (_Pausable *PausableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
//	return _Pausable.Contract.contract.Transact(opts, method, params...)
//}
//
//// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
////
//// Solidity: function owner() constant returns(address)
//func (_Pausable *PausableCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
//	var (
//		ret0 = new(common.Address)
//	)
//	out := ret0
//	err := _Pausable.contract.Call(opts, out, "owner")
//	return *ret0, err
//}
//
//// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
////
//// Solidity: function owner() constant returns(address)
//func (_Pausable *PausableSession) Owner() (common.Address, error) {
//	return _Pausable.Contract.Owner(&_Pausable.CallOpts)
//}
//
//// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
////
//// Solidity: function owner() constant returns(address)
//func (_Pausable *PausableCallerSession) Owner() (common.Address, error) {
//	return _Pausable.Contract.Owner(&_Pausable.CallOpts)
//}
//
//// Paused is a free data retrieval call binding the contract method 0x5c975abb.
////
//// Solidity: function paused() constant returns(bool)
//func (_Pausable *PausableCaller) Paused(opts *bind.CallOpts) (bool, error) {
//	var (
//		ret0 = new(bool)
//	)
//	out := ret0
//	err := _Pausable.contract.Call(opts, out, "paused")
//	return *ret0, err
//}
//
//// Paused is a free data retrieval call binding the contract method 0x5c975abb.
////
//// Solidity: function paused() constant returns(bool)
//func (_Pausable *PausableSession) Paused() (bool, error) {
//	return _Pausable.Contract.Paused(&_Pausable.CallOpts)
//}
//
//// Paused is a free data retrieval call binding the contract method 0x5c975abb.
////
//// Solidity: function paused() constant returns(bool)
//func (_Pausable *PausableCallerSession) Paused() (bool, error) {
//	return _Pausable.Contract.Paused(&_Pausable.CallOpts)
//}
//
//// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
////
//// Solidity: function pause() returns()
//func (_Pausable *PausableTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _Pausable.contract.Transact(opts, "pause")
//}
//
//// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
////
//// Solidity: function pause() returns()
//func (_Pausable *PausableSession) Pause() (*types.Transaction, error) {
//	return _Pausable.Contract.Pause(&_Pausable.TransactOpts)
//}
//
//// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
////
//// Solidity: function pause() returns()
//func (_Pausable *PausableTransactorSession) Pause() (*types.Transaction, error) {
//	return _Pausable.Contract.Pause(&_Pausable.TransactOpts)
//}
//
//// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
////
//// Solidity: function transferOwnership(newOwner address) returns()
//func (_Pausable *PausableTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
//	return _Pausable.contract.Transact(opts, "transferOwnership", newOwner)
//}
//
//// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
////
//// Solidity: function transferOwnership(newOwner address) returns()
//func (_Pausable *PausableSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
//	return _Pausable.Contract.TransferOwnership(&_Pausable.TransactOpts, newOwner)
//}
//
//// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
////
//// Solidity: function transferOwnership(newOwner address) returns()
//func (_Pausable *PausableTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
//	return _Pausable.Contract.TransferOwnership(&_Pausable.TransactOpts, newOwner)
//}
//
//// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
////
//// Solidity: function unpause() returns()
//func (_Pausable *PausableTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _Pausable.contract.Transact(opts, "unpause")
//}
//
//// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
////
//// Solidity: function unpause() returns()
//func (_Pausable *PausableSession) Unpause() (*types.Transaction, error) {
//	return _Pausable.Contract.Unpause(&_Pausable.TransactOpts)
//}
//
//// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
////
//// Solidity: function unpause() returns()
//func (_Pausable *PausableTransactorSession) Unpause() (*types.Transaction, error) {
//	return _Pausable.Contract.Unpause(&_Pausable.TransactOpts)
//}
//
//// PausableOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Pausable contract.
//type PausableOwnershipTransferredIterator struct {
//	Event *PausableOwnershipTransferred // Event containing the contract specifics and raw log
//
//	contract *bind.BoundContract // Generic contract to use for unpacking event data
//	event    string              // Event name to use for unpacking event data
//
//	logs chan types.Log        // Log channel receiving the found contract events
//	sub  ethereum.Subscription // Subscription for errors, completion and termination
//	done bool                  // Whether the subscription completed delivering logs
//	fail error                 // Occurred error to stop iteration
//}
//
//// Next advances the iterator to the subsequent event, returning whether there
//// are any more events found. In case of a retrieval or parsing error, false is
//// returned and Error() can be queried for the exact failure.
//func (it *PausableOwnershipTransferredIterator) Next() bool {
//	// If the iterator failed, stop iterating
//	if it.fail != nil {
//		return false
//	}
//	// If the iterator completed, deliver directly whatever's available
//	if it.done {
//		select {
//		case log := <-it.logs:
//			it.Event = new(PausableOwnershipTransferred)
//			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//				it.fail = err
//				return false
//			}
//			it.Event.Raw = log
//			return true
//
//		default:
//			return false
//		}
//	}
//	// Iterator still in progress, wait for either a data or an error event
//	select {
//	case log := <-it.logs:
//		it.Event = new(PausableOwnershipTransferred)
//		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//			it.fail = err
//			return false
//		}
//		it.Event.Raw = log
//		return true
//
//	case err := <-it.sub.Err():
//		it.done = true
//		it.fail = err
//		return it.Next()
//	}
//}
//
//// Error returns any retrieval or parsing error occurred during filtering.
//func (it *PausableOwnershipTransferredIterator) Error() error {
//	return it.fail
//}
//
//// Close terminates the iteration process, releasing any pending underlying
//// resources.
//func (it *PausableOwnershipTransferredIterator) Close() error {
//	it.sub.Unsubscribe()
//	return nil
//}
//
//// PausableOwnershipTransferred represents a OwnershipTransferred event raised by the Pausable contract.
//type PausableOwnershipTransferred struct {
//	PreviousOwner common.Address
//	NewOwner      common.Address
//	Raw           types.Log // Blockchain specific contextual infos
//}
//
//// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
////
//// Solidity: event OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
//func (_Pausable *PausableFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*PausableOwnershipTransferredIterator, error) {
//
//	var previousOwnerRule []interface{}
//	for _, previousOwnerItem := range previousOwner {
//		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
//	}
//	var newOwnerRule []interface{}
//	for _, newOwnerItem := range newOwner {
//		newOwnerRule = append(newOwnerRule, newOwnerItem)
//	}
//
//	logs, sub, err := _Pausable.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
//	if err != nil {
//		return nil, err
//	}
//	return &PausableOwnershipTransferredIterator{contract: _Pausable.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
//}
//
//// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
////
//// Solidity: event OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
//func (_Pausable *PausableFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PausableOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {
//
//	var previousOwnerRule []interface{}
//	for _, previousOwnerItem := range previousOwner {
//		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
//	}
//	var newOwnerRule []interface{}
//	for _, newOwnerItem := range newOwner {
//		newOwnerRule = append(newOwnerRule, newOwnerItem)
//	}
//
//	logs, sub, err := _Pausable.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
//	if err != nil {
//		return nil, err
//	}
//	return event.NewSubscription(func(quit <-chan struct{}) error {
//		defer sub.Unsubscribe()
//		for {
//			select {
//			case log := <-logs:
//				// New log arrived, parse the event and forward to the user
//				event := new(PausableOwnershipTransferred)
//				if err := _Pausable.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
//					return err
//				}
//				event.Raw = log
//
//				select {
//				case sink <- event:
//				case err := <-sub.Err():
//					return err
//				case <-quit:
//					return nil
//				}
//			case err := <-sub.Err():
//				return err
//			case <-quit:
//				return nil
//			}
//		}
//	}), nil
//}
//
//// PausablePauseIterator is returned from FilterPause and is used to iterate over the raw logs and unpacked data for Pause events raised by the Pausable contract.
//type PausablePauseIterator struct {
//	Event *PausablePause // Event containing the contract specifics and raw log
//
//	contract *bind.BoundContract // Generic contract to use for unpacking event data
//	event    string              // Event name to use for unpacking event data
//
//	logs chan types.Log        // Log channel receiving the found contract events
//	sub  ethereum.Subscription // Subscription for errors, completion and termination
//	done bool                  // Whether the subscription completed delivering logs
//	fail error                 // Occurred error to stop iteration
//}
//
//// Next advances the iterator to the subsequent event, returning whether there
//// are any more events found. In case of a retrieval or parsing error, false is
//// returned and Error() can be queried for the exact failure.
//func (it *PausablePauseIterator) Next() bool {
//	// If the iterator failed, stop iterating
//	if it.fail != nil {
//		return false
//	}
//	// If the iterator completed, deliver directly whatever's available
//	if it.done {
//		select {
//		case log := <-it.logs:
//			it.Event = new(PausablePause)
//			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//				it.fail = err
//				return false
//			}
//			it.Event.Raw = log
//			return true
//
//		default:
//			return false
//		}
//	}
//	// Iterator still in progress, wait for either a data or an error event
//	select {
//	case log := <-it.logs:
//		it.Event = new(PausablePause)
//		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//			it.fail = err
//			return false
//		}
//		it.Event.Raw = log
//		return true
//
//	case err := <-it.sub.Err():
//		it.done = true
//		it.fail = err
//		return it.Next()
//	}
//}
//
//// Error returns any retrieval or parsing error occurred during filtering.
//func (it *PausablePauseIterator) Error() error {
//	return it.fail
//}
//
//// Close terminates the iteration process, releasing any pending underlying
//// resources.
//func (it *PausablePauseIterator) Close() error {
//	it.sub.Unsubscribe()
//	return nil
//}
//
//// PausablePause represents a Pause event raised by the Pausable contract.
//type PausablePause struct {
//	Raw types.Log // Blockchain specific contextual infos
//}
//
//// FilterPause is a free log retrieval operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
////
//// Solidity: event Pause()
//func (_Pausable *PausableFilterer) FilterPause(opts *bind.FilterOpts) (*PausablePauseIterator, error) {
//
//	logs, sub, err := _Pausable.contract.FilterLogs(opts, "Pause")
//	if err != nil {
//		return nil, err
//	}
//	return &PausablePauseIterator{contract: _Pausable.contract, event: "Pause", logs: logs, sub: sub}, nil
//}
//
//// WatchPause is a free log subscription operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
////
//// Solidity: event Pause()
//func (_Pausable *PausableFilterer) WatchPause(opts *bind.WatchOpts, sink chan<- *PausablePause) (event.Subscription, error) {
//
//	logs, sub, err := _Pausable.contract.WatchLogs(opts, "Pause")
//	if err != nil {
//		return nil, err
//	}
//	return event.NewSubscription(func(quit <-chan struct{}) error {
//		defer sub.Unsubscribe()
//		for {
//			select {
//			case log := <-logs:
//				// New log arrived, parse the event and forward to the user
//				event := new(PausablePause)
//				if err := _Pausable.contract.UnpackLog(event, "Pause", log); err != nil {
//					return err
//				}
//				event.Raw = log
//
//				select {
//				case sink <- event:
//				case err := <-sub.Err():
//					return err
//				case <-quit:
//					return nil
//				}
//			case err := <-sub.Err():
//				return err
//			case <-quit:
//				return nil
//			}
//		}
//	}), nil
//}
//
//// PausableUnpauseIterator is returned from FilterUnpause and is used to iterate over the raw logs and unpacked data for Unpause events raised by the Pausable contract.
//type PausableUnpauseIterator struct {
//	Event *PausableUnpause // Event containing the contract specifics and raw log
//
//	contract *bind.BoundContract // Generic contract to use for unpacking event data
//	event    string              // Event name to use for unpacking event data
//
//	logs chan types.Log        // Log channel receiving the found contract events
//	sub  ethereum.Subscription // Subscription for errors, completion and termination
//	done bool                  // Whether the subscription completed delivering logs
//	fail error                 // Occurred error to stop iteration
//}
//
//// Next advances the iterator to the subsequent event, returning whether there
//// are any more events found. In case of a retrieval or parsing error, false is
//// returned and Error() can be queried for the exact failure.
//func (it *PausableUnpauseIterator) Next() bool {
//	// If the iterator failed, stop iterating
//	if it.fail != nil {
//		return false
//	}
//	// If the iterator completed, deliver directly whatever's available
//	if it.done {
//		select {
//		case log := <-it.logs:
//			it.Event = new(PausableUnpause)
//			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//				it.fail = err
//				return false
//			}
//			it.Event.Raw = log
//			return true
//
//		default:
//			return false
//		}
//	}
//	// Iterator still in progress, wait for either a data or an error event
//	select {
//	case log := <-it.logs:
//		it.Event = new(PausableUnpause)
//		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//			it.fail = err
//			return false
//		}
//		it.Event.Raw = log
//		return true
//
//	case err := <-it.sub.Err():
//		it.done = true
//		it.fail = err
//		return it.Next()
//	}
//}
//
//// Error returns any retrieval or parsing error occurred during filtering.
//func (it *PausableUnpauseIterator) Error() error {
//	return it.fail
//}
//
//// Close terminates the iteration process, releasing any pending underlying
//// resources.
//func (it *PausableUnpauseIterator) Close() error {
//	it.sub.Unsubscribe()
//	return nil
//}
//
//// PausableUnpause represents a Unpause event raised by the Pausable contract.
//type PausableUnpause struct {
//	Raw types.Log // Blockchain specific contextual infos
//}
//
//// FilterUnpause is a free log retrieval operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
////
//// Solidity: event Unpause()
//func (_Pausable *PausableFilterer) FilterUnpause(opts *bind.FilterOpts) (*PausableUnpauseIterator, error) {
//
//	logs, sub, err := _Pausable.contract.FilterLogs(opts, "Unpause")
//	if err != nil {
//		return nil, err
//	}
//	return &PausableUnpauseIterator{contract: _Pausable.contract, event: "Unpause", logs: logs, sub: sub}, nil
//}
//
//// WatchUnpause is a free log subscription operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
////
//// Solidity: event Unpause()
//func (_Pausable *PausableFilterer) WatchUnpause(opts *bind.WatchOpts, sink chan<- *PausableUnpause) (event.Subscription, error) {
//
//	logs, sub, err := _Pausable.contract.WatchLogs(opts, "Unpause")
//	if err != nil {
//		return nil, err
//	}
//	return event.NewSubscription(func(quit <-chan struct{}) error {
//		defer sub.Unsubscribe()
//		for {
//			select {
//			case log := <-logs:
//				// New log arrived, parse the event and forward to the user
//				event := new(PausableUnpause)
//				if err := _Pausable.contract.UnpackLog(event, "Unpause", log); err != nil {
//					return err
//				}
//				event.Raw = log
//
//				select {
//				case sink <- event:
//				case err := <-sub.Err():
//					return err
//				case <-quit:
//					return nil
//				}
//			case err := <-sub.Err():
//				return err
//			case <-quit:
//				return nil
//			}
//		}
//	}), nil
//}
//
//// PausableTokenABI is the input ABI used to generate the binding from.
//const PausableTokenABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseApproval\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseApproval\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]"
//
//// PausableTokenBin is the compiled bytecode used for deploying new contracts.
//const PausableTokenBin = `0x606060405260038054600160a860020a03191633600160a060020a0316179055610a5c8061002e6000396000f3006060604052600436106100cf5763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663095ea7b381146100d457806318160ddd1461010a57806323b872dd1461012f57806327e235e3146101575780633f4ba83a146101765780635c975abb1461018b578063661884631461019e57806370a08231146101c05780638456cb59146101df5780638da5cb5b146101f2578063a9059cbb14610221578063d73dd62314610243578063dd62ed3e14610265578063f2fde38b1461028a575b600080fd5b34156100df57600080fd5b6100f6600160a060020a03600435166024356102a9565b604051901515815260200160405180910390f35b341561011557600080fd5b61011d6102d4565b60405190815260200160405180910390f35b341561013a57600080fd5b6100f6600160a060020a03600435811690602435166044356102da565b341561016257600080fd5b61011d600160a060020a0360043516610307565b341561018157600080fd5b610189610319565b005b341561019657600080fd5b6100f6610398565b34156101a957600080fd5b6100f6600160a060020a03600435166024356103a8565b34156101cb57600080fd5b61011d600160a060020a03600435166103cc565b34156101ea57600080fd5b6101896103e7565b34156101fd57600080fd5b61020561046b565b604051600160a060020a03909116815260200160405180910390f35b341561022c57600080fd5b6100f6600160a060020a036004351660243561047a565b341561024e57600080fd5b6100f6600160a060020a036004351660243561049e565b341561027057600080fd5b61011d600160a060020a03600435811690602435166104c2565b341561029557600080fd5b610189600160a060020a03600435166104ed565b60035460009060a060020a900460ff16156102c357600080fd5b6102cd8383610588565b9392505050565b60005481565b60035460009060a060020a900460ff16156102f457600080fd5b6102ff8484846105f4565b949350505050565b60016020526000908152604090205481565b60035433600160a060020a0390811691161461033457600080fd5b60035460a060020a900460ff16151561034c57600080fd5b6003805474ff0000000000000000000000000000000000000000191690557f7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b3360405160405180910390a1565b60035460a060020a900460ff1681565b60035460009060a060020a900460ff16156103c257600080fd5b6102cd8383610776565b600160a060020a031660009081526001602052604090205490565b60035433600160a060020a0390811691161461040257600080fd5b60035460a060020a900460ff161561041957600080fd5b6003805474ff0000000000000000000000000000000000000000191660a060020a1790557f6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff62560405160405180910390a1565b600354600160a060020a031681565b60035460009060a060020a900460ff161561049457600080fd5b6102cd8383610870565b60035460009060a060020a900460ff16156104b857600080fd5b6102cd838361096b565b600160a060020a03918216600090815260026020908152604080832093909416825291909152205490565b60035433600160a060020a0390811691161461050857600080fd5b600160a060020a038116151561051d57600080fd5b600354600160a060020a0380831691167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a36003805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b600160a060020a03338116600081815260026020908152604080832094871680845294909152808220859055909291907f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259085905190815260200160405180910390a350600192915050565b6000600160a060020a038316151561060b57600080fd5b600160a060020a03841660009081526001602052604090205482111561063057600080fd5b600160a060020a038085166000908152600260209081526040808320339094168352929052205482111561066357600080fd5b600160a060020a03841660009081526001602052604090205461068c908363ffffffff610a0f16565b600160a060020a0380861660009081526001602052604080822093909355908516815220546106c1908363ffffffff610a2116565b600160a060020a03808516600090815260016020908152604080832094909455878316825260028152838220339093168252919091522054610709908363ffffffff610a0f16565b600160a060020a03808616600081815260026020908152604080832033861684529091529081902093909355908516917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9085905190815260200160405180910390a35060019392505050565b600160a060020a033381166000908152600260209081526040808320938616835292905290812054808311156107d357600160a060020a03338116600090815260026020908152604080832093881683529290529081205561080a565b6107e3818463ffffffff610a0f16565b600160a060020a033381166000908152600260209081526040808320938916835292905220555b600160a060020a0333811660008181526002602090815260408083209489168084529490915290819020547f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925915190815260200160405180910390a35060019392505050565b6000600160a060020a038316151561088757600080fd5b600160a060020a0333166000908152600160205260409020548211156108ac57600080fd5b600160a060020a0333166000908152600160205260409020546108d5908363ffffffff610a0f16565b600160a060020a03338116600090815260016020526040808220939093559085168152205461090a908363ffffffff610a2116565b600160a060020a0380851660008181526001602052604090819020939093559133909116907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9085905190815260200160405180910390a350600192915050565b600160a060020a0333811660009081526002602090815260408083209386168352929052908120546109a3908363ffffffff610a2116565b600160a060020a0333811660008181526002602090815260408083209489168084529490915290819020849055919290917f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92591905190815260200160405180910390a350600192915050565b600082821115610a1b57fe5b50900390565b6000828201838110156102cd57fe00a165627a7a72305820d13637e350e2fb1942b6d20da16972f163fbfe9b1d2fc5c9ff7f0d8a1a62fd9d0029`
//
//// DeployPausableToken deploys a new Ethereum contract, binding an instance of PausableToken to it.
//func DeployPausableToken(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *PausableToken, error) {
//	parsed, err := abi.JSON(strings.NewReader(PausableTokenABI))
//	if err != nil {
//		return common.Address{}, nil, nil, err
//	}
//	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(PausableTokenBin), backend)
//	if err != nil {
//		return common.Address{}, nil, nil, err
//	}
//	return address, tx, &PausableToken{PausableTokenCaller: PausableTokenCaller{contract: contract}, PausableTokenTransactor: PausableTokenTransactor{contract: contract}, PausableTokenFilterer: PausableTokenFilterer{contract: contract}}, nil
//}
//
//// PausableToken is an auto generated Go binding around an Ethereum contract.
//type PausableToken struct {
//	PausableTokenCaller     // Read-only binding to the contract
//	PausableTokenTransactor // Write-only binding to the contract
//	PausableTokenFilterer   // Log filterer for contract events
//}
//
//// PausableTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
//type PausableTokenCaller struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// PausableTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
//type PausableTokenTransactor struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// PausableTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
//type PausableTokenFilterer struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// PausableTokenSession is an auto generated Go binding around an Ethereum contract,
//// with pre-set call and transact options.
//type PausableTokenSession struct {
//	Contract     *PausableToken    // Generic contract binding to set the session for
//	CallOpts     bind.CallOpts     // Call options to use throughout this session
//	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
//}
//
//// PausableTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
//// with pre-set call options.
//type PausableTokenCallerSession struct {
//	Contract *PausableTokenCaller // Generic contract caller binding to set the session for
//	CallOpts bind.CallOpts        // Call options to use throughout this session
//}
//
//// PausableTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
//// with pre-set transact options.
//type PausableTokenTransactorSession struct {
//	Contract     *PausableTokenTransactor // Generic contract transactor binding to set the session for
//	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
//}
//
//// PausableTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
//type PausableTokenRaw struct {
//	Contract *PausableToken // Generic contract binding to access the raw methods on
//}
//
//// PausableTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
//type PausableTokenCallerRaw struct {
//	Contract *PausableTokenCaller // Generic read-only contract binding to access the raw methods on
//}
//
//// PausableTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
//type PausableTokenTransactorRaw struct {
//	Contract *PausableTokenTransactor // Generic write-only contract binding to access the raw methods on
//}
//
//// NewPausableToken creates a new instance of PausableToken, bound to a specific deployed contract.
//func NewPausableToken(address common.Address, backend bind.ContractBackend) (*PausableToken, error) {
//	contract, err := bindPausableToken(address, backend, backend, backend)
//	if err != nil {
//		return nil, err
//	}
//	return &PausableToken{PausableTokenCaller: PausableTokenCaller{contract: contract}, PausableTokenTransactor: PausableTokenTransactor{contract: contract}, PausableTokenFilterer: PausableTokenFilterer{contract: contract}}, nil
//}
//
//// NewPausableTokenCaller creates a new read-only instance of PausableToken, bound to a specific deployed contract.
//func NewPausableTokenCaller(address common.Address, caller bind.ContractCaller) (*PausableTokenCaller, error) {
//	contract, err := bindPausableToken(address, caller, nil, nil)
//	if err != nil {
//		return nil, err
//	}
//	return &PausableTokenCaller{contract: contract}, nil
//}
//
//// NewPausableTokenTransactor creates a new write-only instance of PausableToken, bound to a specific deployed contract.
//func NewPausableTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*PausableTokenTransactor, error) {
//	contract, err := bindPausableToken(address, nil, transactor, nil)
//	if err != nil {
//		return nil, err
//	}
//	return &PausableTokenTransactor{contract: contract}, nil
//}
//
//// NewPausableTokenFilterer creates a new log filterer instance of PausableToken, bound to a specific deployed contract.
//func NewPausableTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*PausableTokenFilterer, error) {
//	contract, err := bindPausableToken(address, nil, nil, filterer)
//	if err != nil {
//		return nil, err
//	}
//	return &PausableTokenFilterer{contract: contract}, nil
//}
//
//// bindPausableToken binds a generic wrapper to an already deployed contract.
//func bindPausableToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
//	parsed, err := abi.JSON(strings.NewReader(PausableTokenABI))
//	if err != nil {
//		return nil, err
//	}
//	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
//}
//
//// Call invokes the (constant) contract method with params as input values and
//// sets the output to result. The result type might be a single field for simple
//// returns, a slice of interfaces for anonymous returns and a struct for named
//// returns.
//func (_PausableToken *PausableTokenRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
//	return _PausableToken.Contract.PausableTokenCaller.contract.Call(opts, result, method, params...)
//}
//
//// Transfer initiates a plain transaction to move funds to the contract, calling
//// its default method if one is available.
//func (_PausableToken *PausableTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _PausableToken.Contract.PausableTokenTransactor.contract.Transfer(opts)
//}
//
//// Transact invokes the (paid) contract method with params as input values.
//func (_PausableToken *PausableTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
//	return _PausableToken.Contract.PausableTokenTransactor.contract.Transact(opts, method, params...)
//}
//
//// Call invokes the (constant) contract method with params as input values and
//// sets the output to result. The result type might be a single field for simple
//// returns, a slice of interfaces for anonymous returns and a struct for named
//// returns.
//func (_PausableToken *PausableTokenCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
//	return _PausableToken.Contract.contract.Call(opts, result, method, params...)
//}
//
//// Transfer initiates a plain transaction to move funds to the contract, calling
//// its default method if one is available.
//func (_PausableToken *PausableTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _PausableToken.Contract.contract.Transfer(opts)
//}
//
//// Transact invokes the (paid) contract method with params as input values.
//func (_PausableToken *PausableTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
//	return _PausableToken.Contract.contract.Transact(opts, method, params...)
//}
//
//// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
////
//// Solidity: function allowance(_owner address, _spender address) constant returns(uint256)
//func (_PausableToken *PausableTokenCaller) Allowance(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (*big.Int, error) {
//	var (
//		ret0 = new(*big.Int)
//	)
//	out := ret0
//	err := _PausableToken.contract.Call(opts, out, "allowance", _owner, _spender)
//	return *ret0, err
//}
//
//// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
////
//// Solidity: function allowance(_owner address, _spender address) constant returns(uint256)
//func (_PausableToken *PausableTokenSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
//	return _PausableToken.Contract.Allowance(&_PausableToken.CallOpts, _owner, _spender)
//}
//
//// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
////
//// Solidity: function allowance(_owner address, _spender address) constant returns(uint256)
//func (_PausableToken *PausableTokenCallerSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
//	return _PausableToken.Contract.Allowance(&_PausableToken.CallOpts, _owner, _spender)
//}
//
//// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
////
//// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
//func (_PausableToken *PausableTokenCaller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
//	var (
//		ret0 = new(*big.Int)
//	)
//	out := ret0
//	err := _PausableToken.contract.Call(opts, out, "balanceOf", _owner)
//	return *ret0, err
//}
//
//// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
////
//// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
//func (_PausableToken *PausableTokenSession) BalanceOf(_owner common.Address) (*big.Int, error) {
//	return _PausableToken.Contract.BalanceOf(&_PausableToken.CallOpts, _owner)
//}
//
//// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
////
//// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
//func (_PausableToken *PausableTokenCallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
//	return _PausableToken.Contract.BalanceOf(&_PausableToken.CallOpts, _owner)
//}
//
//// Balances is a free data retrieval call binding the contract method 0x27e235e3.
////
//// Solidity: function balances( address) constant returns(uint256)
//func (_PausableToken *PausableTokenCaller) Balances(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
//	var (
//		ret0 = new(*big.Int)
//	)
//	out := ret0
//	err := _PausableToken.contract.Call(opts, out, "balances", arg0)
//	return *ret0, err
//}
//
//// Balances is a free data retrieval call binding the contract method 0x27e235e3.
////
//// Solidity: function balances( address) constant returns(uint256)
//func (_PausableToken *PausableTokenSession) Balances(arg0 common.Address) (*big.Int, error) {
//	return _PausableToken.Contract.Balances(&_PausableToken.CallOpts, arg0)
//}
//
//// Balances is a free data retrieval call binding the contract method 0x27e235e3.
////
//// Solidity: function balances( address) constant returns(uint256)
//func (_PausableToken *PausableTokenCallerSession) Balances(arg0 common.Address) (*big.Int, error) {
//	return _PausableToken.Contract.Balances(&_PausableToken.CallOpts, arg0)
//}
//
//// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
////
//// Solidity: function owner() constant returns(address)
//func (_PausableToken *PausableTokenCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
//	var (
//		ret0 = new(common.Address)
//	)
//	out := ret0
//	err := _PausableToken.contract.Call(opts, out, "owner")
//	return *ret0, err
//}
//
//// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
////
//// Solidity: function owner() constant returns(address)
//func (_PausableToken *PausableTokenSession) Owner() (common.Address, error) {
//	return _PausableToken.Contract.Owner(&_PausableToken.CallOpts)
//}
//
//// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
////
//// Solidity: function owner() constant returns(address)
//func (_PausableToken *PausableTokenCallerSession) Owner() (common.Address, error) {
//	return _PausableToken.Contract.Owner(&_PausableToken.CallOpts)
//}
//
//// Paused is a free data retrieval call binding the contract method 0x5c975abb.
////
//// Solidity: function paused() constant returns(bool)
//func (_PausableToken *PausableTokenCaller) Paused(opts *bind.CallOpts) (bool, error) {
//	var (
//		ret0 = new(bool)
//	)
//	out := ret0
//	err := _PausableToken.contract.Call(opts, out, "paused")
//	return *ret0, err
//}
//
//// Paused is a free data retrieval call binding the contract method 0x5c975abb.
////
//// Solidity: function paused() constant returns(bool)
//func (_PausableToken *PausableTokenSession) Paused() (bool, error) {
//	return _PausableToken.Contract.Paused(&_PausableToken.CallOpts)
//}
//
//// Paused is a free data retrieval call binding the contract method 0x5c975abb.
////
//// Solidity: function paused() constant returns(bool)
//func (_PausableToken *PausableTokenCallerSession) Paused() (bool, error) {
//	return _PausableToken.Contract.Paused(&_PausableToken.CallOpts)
//}
//
//// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
////
//// Solidity: function totalSupply() constant returns(uint256)
//func (_PausableToken *PausableTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
//	var (
//		ret0 = new(*big.Int)
//	)
//	out := ret0
//	err := _PausableToken.contract.Call(opts, out, "totalSupply")
//	return *ret0, err
//}
//
//// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
////
//// Solidity: function totalSupply() constant returns(uint256)
//func (_PausableToken *PausableTokenSession) TotalSupply() (*big.Int, error) {
//	return _PausableToken.Contract.TotalSupply(&_PausableToken.CallOpts)
//}
//
//// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
////
//// Solidity: function totalSupply() constant returns(uint256)
//func (_PausableToken *PausableTokenCallerSession) TotalSupply() (*big.Int, error) {
//	return _PausableToken.Contract.TotalSupply(&_PausableToken.CallOpts)
//}
//
//// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
////
//// Solidity: function approve(_spender address, _value uint256) returns(bool)
//func (_PausableToken *PausableTokenTransactor) Approve(opts *bind.TransactOpts, _spender common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _PausableToken.contract.Transact(opts, "approve", _spender, _value)
//}
//
//// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
////
//// Solidity: function approve(_spender address, _value uint256) returns(bool)
//func (_PausableToken *PausableTokenSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _PausableToken.Contract.Approve(&_PausableToken.TransactOpts, _spender, _value)
//}
//
//// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
////
//// Solidity: function approve(_spender address, _value uint256) returns(bool)
//func (_PausableToken *PausableTokenTransactorSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _PausableToken.Contract.Approve(&_PausableToken.TransactOpts, _spender, _value)
//}
//
//// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
////
//// Solidity: function decreaseApproval(_spender address, _subtractedValue uint256) returns(success bool)
//func (_PausableToken *PausableTokenTransactor) DecreaseApproval(opts *bind.TransactOpts, _spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
//	return _PausableToken.contract.Transact(opts, "decreaseApproval", _spender, _subtractedValue)
//}
//
//// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
////
//// Solidity: function decreaseApproval(_spender address, _subtractedValue uint256) returns(success bool)
//func (_PausableToken *PausableTokenSession) DecreaseApproval(_spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
//	return _PausableToken.Contract.DecreaseApproval(&_PausableToken.TransactOpts, _spender, _subtractedValue)
//}
//
//// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
////
//// Solidity: function decreaseApproval(_spender address, _subtractedValue uint256) returns(success bool)
//func (_PausableToken *PausableTokenTransactorSession) DecreaseApproval(_spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
//	return _PausableToken.Contract.DecreaseApproval(&_PausableToken.TransactOpts, _spender, _subtractedValue)
//}
//
//// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
////
//// Solidity: function increaseApproval(_spender address, _addedValue uint256) returns(success bool)
//func (_PausableToken *PausableTokenTransactor) IncreaseApproval(opts *bind.TransactOpts, _spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
//	return _PausableToken.contract.Transact(opts, "increaseApproval", _spender, _addedValue)
//}
//
//// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
////
//// Solidity: function increaseApproval(_spender address, _addedValue uint256) returns(success bool)
//func (_PausableToken *PausableTokenSession) IncreaseApproval(_spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
//	return _PausableToken.Contract.IncreaseApproval(&_PausableToken.TransactOpts, _spender, _addedValue)
//}
//
//// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
////
//// Solidity: function increaseApproval(_spender address, _addedValue uint256) returns(success bool)
//func (_PausableToken *PausableTokenTransactorSession) IncreaseApproval(_spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
//	return _PausableToken.Contract.IncreaseApproval(&_PausableToken.TransactOpts, _spender, _addedValue)
//}
//
//// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
////
//// Solidity: function pause() returns()
//func (_PausableToken *PausableTokenTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _PausableToken.contract.Transact(opts, "pause")
//}
//
//// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
////
//// Solidity: function pause() returns()
//func (_PausableToken *PausableTokenSession) Pause() (*types.Transaction, error) {
//	return _PausableToken.Contract.Pause(&_PausableToken.TransactOpts)
//}
//
//// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
////
//// Solidity: function pause() returns()
//func (_PausableToken *PausableTokenTransactorSession) Pause() (*types.Transaction, error) {
//	return _PausableToken.Contract.Pause(&_PausableToken.TransactOpts)
//}
//
//// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
////
//// Solidity: function transfer(_to address, _value uint256) returns(bool)
//func (_PausableToken *PausableTokenTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _PausableToken.contract.Transact(opts, "transfer", _to, _value)
//}
//
//// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
////
//// Solidity: function transfer(_to address, _value uint256) returns(bool)
//func (_PausableToken *PausableTokenSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _PausableToken.Contract.Transfer(&_PausableToken.TransactOpts, _to, _value)
//}
//
//// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
////
//// Solidity: function transfer(_to address, _value uint256) returns(bool)
//func (_PausableToken *PausableTokenTransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _PausableToken.Contract.Transfer(&_PausableToken.TransactOpts, _to, _value)
//}
//
//// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
////
//// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
//func (_PausableToken *PausableTokenTransactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _PausableToken.contract.Transact(opts, "transferFrom", _from, _to, _value)
//}
//
//// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
////
//// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
//func (_PausableToken *PausableTokenSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _PausableToken.Contract.TransferFrom(&_PausableToken.TransactOpts, _from, _to, _value)
//}
//
//// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
////
//// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
//func (_PausableToken *PausableTokenTransactorSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _PausableToken.Contract.TransferFrom(&_PausableToken.TransactOpts, _from, _to, _value)
//}
//
//// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
////
//// Solidity: function transferOwnership(newOwner address) returns()
//func (_PausableToken *PausableTokenTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
//	return _PausableToken.contract.Transact(opts, "transferOwnership", newOwner)
//}
//
//// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
////
//// Solidity: function transferOwnership(newOwner address) returns()
//func (_PausableToken *PausableTokenSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
//	return _PausableToken.Contract.TransferOwnership(&_PausableToken.TransactOpts, newOwner)
//}
//
//// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
////
//// Solidity: function transferOwnership(newOwner address) returns()
//func (_PausableToken *PausableTokenTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
//	return _PausableToken.Contract.TransferOwnership(&_PausableToken.TransactOpts, newOwner)
//}
//
//// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
////
//// Solidity: function unpause() returns()
//func (_PausableToken *PausableTokenTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _PausableToken.contract.Transact(opts, "unpause")
//}
//
//// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
////
//// Solidity: function unpause() returns()
//func (_PausableToken *PausableTokenSession) Unpause() (*types.Transaction, error) {
//	return _PausableToken.Contract.Unpause(&_PausableToken.TransactOpts)
//}
//
//// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
////
//// Solidity: function unpause() returns()
//func (_PausableToken *PausableTokenTransactorSession) Unpause() (*types.Transaction, error) {
//	return _PausableToken.Contract.Unpause(&_PausableToken.TransactOpts)
//}
//
//// PausableTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the PausableToken contract.
//type PausableTokenApprovalIterator struct {
//	Event *PausableTokenApproval // Event containing the contract specifics and raw log
//
//	contract *bind.BoundContract // Generic contract to use for unpacking event data
//	event    string              // Event name to use for unpacking event data
//
//	logs chan types.Log        // Log channel receiving the found contract events
//	sub  ethereum.Subscription // Subscription for errors, completion and termination
//	done bool                  // Whether the subscription completed delivering logs
//	fail error                 // Occurred error to stop iteration
//}
//
//// Next advances the iterator to the subsequent event, returning whether there
//// are any more events found. In case of a retrieval or parsing error, false is
//// returned and Error() can be queried for the exact failure.
//func (it *PausableTokenApprovalIterator) Next() bool {
//	// If the iterator failed, stop iterating
//	if it.fail != nil {
//		return false
//	}
//	// If the iterator completed, deliver directly whatever's available
//	if it.done {
//		select {
//		case log := <-it.logs:
//			it.Event = new(PausableTokenApproval)
//			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//				it.fail = err
//				return false
//			}
//			it.Event.Raw = log
//			return true
//
//		default:
//			return false
//		}
//	}
//	// Iterator still in progress, wait for either a data or an error event
//	select {
//	case log := <-it.logs:
//		it.Event = new(PausableTokenApproval)
//		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//			it.fail = err
//			return false
//		}
//		it.Event.Raw = log
//		return true
//
//	case err := <-it.sub.Err():
//		it.done = true
//		it.fail = err
//		return it.Next()
//	}
//}
//
//// Error returns any retrieval or parsing error occurred during filtering.
//func (it *PausableTokenApprovalIterator) Error() error {
//	return it.fail
//}
//
//// Close terminates the iteration process, releasing any pending underlying
//// resources.
//func (it *PausableTokenApprovalIterator) Close() error {
//	it.sub.Unsubscribe()
//	return nil
//}
//
//// PausableTokenApproval represents a Approval event raised by the PausableToken contract.
//type PausableTokenApproval struct {
//	Owner   common.Address
//	Spender common.Address
//	Value   *big.Int
//	Raw     types.Log // Blockchain specific contextual infos
//}
//
//// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
////
//// Solidity: event Approval(owner indexed address, spender indexed address, value uint256)
//func (_PausableToken *PausableTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*PausableTokenApprovalIterator, error) {
//
//	var ownerRule []interface{}
//	for _, ownerItem := range owner {
//		ownerRule = append(ownerRule, ownerItem)
//	}
//	var spenderRule []interface{}
//	for _, spenderItem := range spender {
//		spenderRule = append(spenderRule, spenderItem)
//	}
//
//	logs, sub, err := _PausableToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
//	if err != nil {
//		return nil, err
//	}
//	return &PausableTokenApprovalIterator{contract: _PausableToken.contract, event: "Approval", logs: logs, sub: sub}, nil
//}
//
//// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
////
//// Solidity: event Approval(owner indexed address, spender indexed address, value uint256)
//func (_PausableToken *PausableTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *PausableTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {
//
//	var ownerRule []interface{}
//	for _, ownerItem := range owner {
//		ownerRule = append(ownerRule, ownerItem)
//	}
//	var spenderRule []interface{}
//	for _, spenderItem := range spender {
//		spenderRule = append(spenderRule, spenderItem)
//	}
//
//	logs, sub, err := _PausableToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
//	if err != nil {
//		return nil, err
//	}
//	return event.NewSubscription(func(quit <-chan struct{}) error {
//		defer sub.Unsubscribe()
//		for {
//			select {
//			case log := <-logs:
//				// New log arrived, parse the event and forward to the user
//				event := new(PausableTokenApproval)
//				if err := _PausableToken.contract.UnpackLog(event, "Approval", log); err != nil {
//					return err
//				}
//				event.Raw = log
//
//				select {
//				case sink <- event:
//				case err := <-sub.Err():
//					return err
//				case <-quit:
//					return nil
//				}
//			case err := <-sub.Err():
//				return err
//			case <-quit:
//				return nil
//			}
//		}
//	}), nil
//}
//
//// PausableTokenOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the PausableToken contract.
//type PausableTokenOwnershipTransferredIterator struct {
//	Event *PausableTokenOwnershipTransferred // Event containing the contract specifics and raw log
//
//	contract *bind.BoundContract // Generic contract to use for unpacking event data
//	event    string              // Event name to use for unpacking event data
//
//	logs chan types.Log        // Log channel receiving the found contract events
//	sub  ethereum.Subscription // Subscription for errors, completion and termination
//	done bool                  // Whether the subscription completed delivering logs
//	fail error                 // Occurred error to stop iteration
//}
//
//// Next advances the iterator to the subsequent event, returning whether there
//// are any more events found. In case of a retrieval or parsing error, false is
//// returned and Error() can be queried for the exact failure.
//func (it *PausableTokenOwnershipTransferredIterator) Next() bool {
//	// If the iterator failed, stop iterating
//	if it.fail != nil {
//		return false
//	}
//	// If the iterator completed, deliver directly whatever's available
//	if it.done {
//		select {
//		case log := <-it.logs:
//			it.Event = new(PausableTokenOwnershipTransferred)
//			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//				it.fail = err
//				return false
//			}
//			it.Event.Raw = log
//			return true
//
//		default:
//			return false
//		}
//	}
//	// Iterator still in progress, wait for either a data or an error event
//	select {
//	case log := <-it.logs:
//		it.Event = new(PausableTokenOwnershipTransferred)
//		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//			it.fail = err
//			return false
//		}
//		it.Event.Raw = log
//		return true
//
//	case err := <-it.sub.Err():
//		it.done = true
//		it.fail = err
//		return it.Next()
//	}
//}
//
//// Error returns any retrieval or parsing error occurred during filtering.
//func (it *PausableTokenOwnershipTransferredIterator) Error() error {
//	return it.fail
//}
//
//// Close terminates the iteration process, releasing any pending underlying
//// resources.
//func (it *PausableTokenOwnershipTransferredIterator) Close() error {
//	it.sub.Unsubscribe()
//	return nil
//}
//
//// PausableTokenOwnershipTransferred represents a OwnershipTransferred event raised by the PausableToken contract.
//type PausableTokenOwnershipTransferred struct {
//	PreviousOwner common.Address
//	NewOwner      common.Address
//	Raw           types.Log // Blockchain specific contextual infos
//}
//
//// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
////
//// Solidity: event OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
//func (_PausableToken *PausableTokenFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*PausableTokenOwnershipTransferredIterator, error) {
//
//	var previousOwnerRule []interface{}
//	for _, previousOwnerItem := range previousOwner {
//		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
//	}
//	var newOwnerRule []interface{}
//	for _, newOwnerItem := range newOwner {
//		newOwnerRule = append(newOwnerRule, newOwnerItem)
//	}
//
//	logs, sub, err := _PausableToken.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
//	if err != nil {
//		return nil, err
//	}
//	return &PausableTokenOwnershipTransferredIterator{contract: _PausableToken.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
//}
//
//// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
////
//// Solidity: event OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
//func (_PausableToken *PausableTokenFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PausableTokenOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {
//
//	var previousOwnerRule []interface{}
//	for _, previousOwnerItem := range previousOwner {
//		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
//	}
//	var newOwnerRule []interface{}
//	for _, newOwnerItem := range newOwner {
//		newOwnerRule = append(newOwnerRule, newOwnerItem)
//	}
//
//	logs, sub, err := _PausableToken.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
//	if err != nil {
//		return nil, err
//	}
//	return event.NewSubscription(func(quit <-chan struct{}) error {
//		defer sub.Unsubscribe()
//		for {
//			select {
//			case log := <-logs:
//				// New log arrived, parse the event and forward to the user
//				event := new(PausableTokenOwnershipTransferred)
//				if err := _PausableToken.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
//					return err
//				}
//				event.Raw = log
//
//				select {
//				case sink <- event:
//				case err := <-sub.Err():
//					return err
//				case <-quit:
//					return nil
//				}
//			case err := <-sub.Err():
//				return err
//			case <-quit:
//				return nil
//			}
//		}
//	}), nil
//}
//
//// PausableTokenPauseIterator is returned from FilterPause and is used to iterate over the raw logs and unpacked data for Pause events raised by the PausableToken contract.
//type PausableTokenPauseIterator struct {
//	Event *PausableTokenPause // Event containing the contract specifics and raw log
//
//	contract *bind.BoundContract // Generic contract to use for unpacking event data
//	event    string              // Event name to use for unpacking event data
//
//	logs chan types.Log        // Log channel receiving the found contract events
//	sub  ethereum.Subscription // Subscription for errors, completion and termination
//	done bool                  // Whether the subscription completed delivering logs
//	fail error                 // Occurred error to stop iteration
//}
//
//// Next advances the iterator to the subsequent event, returning whether there
//// are any more events found. In case of a retrieval or parsing error, false is
//// returned and Error() can be queried for the exact failure.
//func (it *PausableTokenPauseIterator) Next() bool {
//	// If the iterator failed, stop iterating
//	if it.fail != nil {
//		return false
//	}
//	// If the iterator completed, deliver directly whatever's available
//	if it.done {
//		select {
//		case log := <-it.logs:
//			it.Event = new(PausableTokenPause)
//			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//				it.fail = err
//				return false
//			}
//			it.Event.Raw = log
//			return true
//
//		default:
//			return false
//		}
//	}
//	// Iterator still in progress, wait for either a data or an error event
//	select {
//	case log := <-it.logs:
//		it.Event = new(PausableTokenPause)
//		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//			it.fail = err
//			return false
//		}
//		it.Event.Raw = log
//		return true
//
//	case err := <-it.sub.Err():
//		it.done = true
//		it.fail = err
//		return it.Next()
//	}
//}
//
//// Error returns any retrieval or parsing error occurred during filtering.
//func (it *PausableTokenPauseIterator) Error() error {
//	return it.fail
//}
//
//// Close terminates the iteration process, releasing any pending underlying
//// resources.
//func (it *PausableTokenPauseIterator) Close() error {
//	it.sub.Unsubscribe()
//	return nil
//}
//
//// PausableTokenPause represents a Pause event raised by the PausableToken contract.
//type PausableTokenPause struct {
//	Raw types.Log // Blockchain specific contextual infos
//}
//
//// FilterPause is a free log retrieval operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
////
//// Solidity: event Pause()
//func (_PausableToken *PausableTokenFilterer) FilterPause(opts *bind.FilterOpts) (*PausableTokenPauseIterator, error) {
//
//	logs, sub, err := _PausableToken.contract.FilterLogs(opts, "Pause")
//	if err != nil {
//		return nil, err
//	}
//	return &PausableTokenPauseIterator{contract: _PausableToken.contract, event: "Pause", logs: logs, sub: sub}, nil
//}
//
//// WatchPause is a free log subscription operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
////
//// Solidity: event Pause()
//func (_PausableToken *PausableTokenFilterer) WatchPause(opts *bind.WatchOpts, sink chan<- *PausableTokenPause) (event.Subscription, error) {
//
//	logs, sub, err := _PausableToken.contract.WatchLogs(opts, "Pause")
//	if err != nil {
//		return nil, err
//	}
//	return event.NewSubscription(func(quit <-chan struct{}) error {
//		defer sub.Unsubscribe()
//		for {
//			select {
//			case log := <-logs:
//				// New log arrived, parse the event and forward to the user
//				event := new(PausableTokenPause)
//				if err := _PausableToken.contract.UnpackLog(event, "Pause", log); err != nil {
//					return err
//				}
//				event.Raw = log
//
//				select {
//				case sink <- event:
//				case err := <-sub.Err():
//					return err
//				case <-quit:
//					return nil
//				}
//			case err := <-sub.Err():
//				return err
//			case <-quit:
//				return nil
//			}
//		}
//	}), nil
//}
//
//// PausableTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the PausableToken contract.
//type PausableTokenTransferIterator struct {
//	Event *PausableTokenTransfer // Event containing the contract specifics and raw log
//
//	contract *bind.BoundContract // Generic contract to use for unpacking event data
//	event    string              // Event name to use for unpacking event data
//
//	logs chan types.Log        // Log channel receiving the found contract events
//	sub  ethereum.Subscription // Subscription for errors, completion and termination
//	done bool                  // Whether the subscription completed delivering logs
//	fail error                 // Occurred error to stop iteration
//}
//
//// Next advances the iterator to the subsequent event, returning whether there
//// are any more events found. In case of a retrieval or parsing error, false is
//// returned and Error() can be queried for the exact failure.
//func (it *PausableTokenTransferIterator) Next() bool {
//	// If the iterator failed, stop iterating
//	if it.fail != nil {
//		return false
//	}
//	// If the iterator completed, deliver directly whatever's available
//	if it.done {
//		select {
//		case log := <-it.logs:
//			it.Event = new(PausableTokenTransfer)
//			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//				it.fail = err
//				return false
//			}
//			it.Event.Raw = log
//			return true
//
//		default:
//			return false
//		}
//	}
//	// Iterator still in progress, wait for either a data or an error event
//	select {
//	case log := <-it.logs:
//		it.Event = new(PausableTokenTransfer)
//		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//			it.fail = err
//			return false
//		}
//		it.Event.Raw = log
//		return true
//
//	case err := <-it.sub.Err():
//		it.done = true
//		it.fail = err
//		return it.Next()
//	}
//}
//
//// Error returns any retrieval or parsing error occurred during filtering.
//func (it *PausableTokenTransferIterator) Error() error {
//	return it.fail
//}
//
//// Close terminates the iteration process, releasing any pending underlying
//// resources.
//func (it *PausableTokenTransferIterator) Close() error {
//	it.sub.Unsubscribe()
//	return nil
//}
//
//// PausableTokenTransfer represents a Transfer event raised by the PausableToken contract.
//type PausableTokenTransfer struct {
//	From  common.Address
//	To    common.Address
//	Value *big.Int
//	Raw   types.Log // Blockchain specific contextual infos
//}
//
//// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
////
//// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
//func (_PausableToken *PausableTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PausableTokenTransferIterator, error) {
//
//	var fromRule []interface{}
//	for _, fromItem := range from {
//		fromRule = append(fromRule, fromItem)
//	}
//	var toRule []interface{}
//	for _, toItem := range to {
//		toRule = append(toRule, toItem)
//	}
//
//	logs, sub, err := _PausableToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
//	if err != nil {
//		return nil, err
//	}
//	return &PausableTokenTransferIterator{contract: _PausableToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
//}
//
//// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
////
//// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
//func (_PausableToken *PausableTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *PausableTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {
//
//	var fromRule []interface{}
//	for _, fromItem := range from {
//		fromRule = append(fromRule, fromItem)
//	}
//	var toRule []interface{}
//	for _, toItem := range to {
//		toRule = append(toRule, toItem)
//	}
//
//	logs, sub, err := _PausableToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
//	if err != nil {
//		return nil, err
//	}
//	return event.NewSubscription(func(quit <-chan struct{}) error {
//		defer sub.Unsubscribe()
//		for {
//			select {
//			case log := <-logs:
//				// New log arrived, parse the event and forward to the user
//				event := new(PausableTokenTransfer)
//				if err := _PausableToken.contract.UnpackLog(event, "Transfer", log); err != nil {
//					return err
//				}
//				event.Raw = log
//
//				select {
//				case sink <- event:
//				case err := <-sub.Err():
//					return err
//				case <-quit:
//					return nil
//				}
//			case err := <-sub.Err():
//				return err
//			case <-quit:
//				return nil
//			}
//		}
//	}), nil
//}
//
//// PausableTokenUnpauseIterator is returned from FilterUnpause and is used to iterate over the raw logs and unpacked data for Unpause events raised by the PausableToken contract.
//type PausableTokenUnpauseIterator struct {
//	Event *PausableTokenUnpause // Event containing the contract specifics and raw log
//
//	contract *bind.BoundContract // Generic contract to use for unpacking event data
//	event    string              // Event name to use for unpacking event data
//
//	logs chan types.Log        // Log channel receiving the found contract events
//	sub  ethereum.Subscription // Subscription for errors, completion and termination
//	done bool                  // Whether the subscription completed delivering logs
//	fail error                 // Occurred error to stop iteration
//}
//
//// Next advances the iterator to the subsequent event, returning whether there
//// are any more events found. In case of a retrieval or parsing error, false is
//// returned and Error() can be queried for the exact failure.
//func (it *PausableTokenUnpauseIterator) Next() bool {
//	// If the iterator failed, stop iterating
//	if it.fail != nil {
//		return false
//	}
//	// If the iterator completed, deliver directly whatever's available
//	if it.done {
//		select {
//		case log := <-it.logs:
//			it.Event = new(PausableTokenUnpause)
//			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//				it.fail = err
//				return false
//			}
//			it.Event.Raw = log
//			return true
//
//		default:
//			return false
//		}
//	}
//	// Iterator still in progress, wait for either a data or an error event
//	select {
//	case log := <-it.logs:
//		it.Event = new(PausableTokenUnpause)
//		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//			it.fail = err
//			return false
//		}
//		it.Event.Raw = log
//		return true
//
//	case err := <-it.sub.Err():
//		it.done = true
//		it.fail = err
//		return it.Next()
//	}
//}
//
//// Error returns any retrieval or parsing error occurred during filtering.
//func (it *PausableTokenUnpauseIterator) Error() error {
//	return it.fail
//}
//
//// Close terminates the iteration process, releasing any pending underlying
//// resources.
//func (it *PausableTokenUnpauseIterator) Close() error {
//	it.sub.Unsubscribe()
//	return nil
//}
//
//// PausableTokenUnpause represents a Unpause event raised by the PausableToken contract.
//type PausableTokenUnpause struct {
//	Raw types.Log // Blockchain specific contextual infos
//}
//
//// FilterUnpause is a free log retrieval operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
////
//// Solidity: event Unpause()
//func (_PausableToken *PausableTokenFilterer) FilterUnpause(opts *bind.FilterOpts) (*PausableTokenUnpauseIterator, error) {
//
//	logs, sub, err := _PausableToken.contract.FilterLogs(opts, "Unpause")
//	if err != nil {
//		return nil, err
//	}
//	return &PausableTokenUnpauseIterator{contract: _PausableToken.contract, event: "Unpause", logs: logs, sub: sub}, nil
//}
//
//// WatchUnpause is a free log subscription operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
////
//// Solidity: event Unpause()
//func (_PausableToken *PausableTokenFilterer) WatchUnpause(opts *bind.WatchOpts, sink chan<- *PausableTokenUnpause) (event.Subscription, error) {
//
//	logs, sub, err := _PausableToken.contract.WatchLogs(opts, "Unpause")
//	if err != nil {
//		return nil, err
//	}
//	return event.NewSubscription(func(quit <-chan struct{}) error {
//		defer sub.Unsubscribe()
//		for {
//			select {
//			case log := <-logs:
//				// New log arrived, parse the event and forward to the user
//				event := new(PausableTokenUnpause)
//				if err := _PausableToken.contract.UnpackLog(event, "Unpause", log); err != nil {
//					return err
//				}
//				event.Raw = log
//
//				select {
//				case sink <- event:
//				case err := <-sub.Err():
//					return err
//				case <-quit:
//					return nil
//				}
//			case err := <-sub.Err():
//				return err
//			case <-quit:
//				return nil
//			}
//		}
//	}), nil
//}
//
//// RepublicTokenABI is the input ABI used to generate the binding from.
//const RepublicTokenABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"INITIAL_SUPPLY\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseApproval\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"beneficiary\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferTokens\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseApproval\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"burner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Burn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]"
//
//// RepublicTokenBin is the compiled bytecode used for deploying new contracts.
//const RepublicTokenBin = `0x60606040526003805460a060020a60ff0219169055341561001f57600080fd5b60038054600160a060020a033316600160a060020a031990911681179091556b033b2e3c9fd0803ce800000060008181559182526001602052604090912055610dec8061006d6000396000f3006060604052600436106101115763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166306fdde038114610116578063095ea7b3146101a057806318160ddd146101d657806323b872dd146101fb57806327e235e3146102235780632ff2e9dc14610242578063313ce567146102555780633f4ba83a1461027e57806342966c68146102935780635c975abb146102a957806366188463146102bc57806370a08231146102de5780638456cb59146102fd5780638da5cb5b1461031057806395d89b411461033f578063a9059cbb14610352578063bec3fa1714610374578063d73dd62314610396578063dd62ed3e146103b8578063f2fde38b146103dd575b600080fd5b341561012157600080fd5b6101296103fc565b60405160208082528190810183818151815260200191508051906020019080838360005b8381101561016557808201518382015260200161014d565b50505050905090810190601f1680156101925780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34156101ab57600080fd5b6101c2600160a060020a0360043516602435610433565b604051901515815260200160405180910390f35b34156101e157600080fd5b6101e961045e565b60405190815260200160405180910390f35b341561020657600080fd5b6101c2600160a060020a0360043581169060243516604435610464565b341561022e57600080fd5b6101e9600160a060020a0360043516610491565b341561024d57600080fd5b6101e96104a3565b341561026057600080fd5b6102686104b3565b60405160ff909116815260200160405180910390f35b341561028957600080fd5b6102916104b8565b005b341561029e57600080fd5b610291600435610537565b34156102b457600080fd5b6101c2610600565b34156102c757600080fd5b6101c2600160a060020a0360043516602435610610565b34156102e957600080fd5b6101e9600160a060020a0360043516610634565b341561030857600080fd5b61029161064f565b341561031b57600080fd5b6103236106d3565b604051600160a060020a03909116815260200160405180910390f35b341561034a57600080fd5b6101296106e2565b341561035d57600080fd5b6101c2600160a060020a0360043516602435610719565b341561037f57600080fd5b6101c2600160a060020a036004351660243561073d565b34156103a157600080fd5b6101c2600160a060020a036004351660243561082e565b34156103c357600080fd5b6101e9600160a060020a0360043581169060243516610852565b34156103e857600080fd5b610291600160a060020a036004351661087d565b60408051908101604052600e81527f52657075626c696320546f6b656e000000000000000000000000000000000000602082015281565b60035460009060a060020a900460ff161561044d57600080fd5b6104578383610918565b9392505050565b60005481565b60035460009060a060020a900460ff161561047e57600080fd5b610489848484610984565b949350505050565b60016020526000908152604090205481565b6b033b2e3c9fd0803ce800000081565b601281565b60035433600160a060020a039081169116146104d357600080fd5b60035460a060020a900460ff1615156104eb57600080fd5b6003805474ff0000000000000000000000000000000000000000191690557f7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b3360405160405180910390a1565b600080821161054557600080fd5b600160a060020a03331660009081526001602052604090205482111561056a57600080fd5b5033600160a060020a03811660009081526001602052604090205461058f9083610b06565b600160a060020a038216600090815260016020526040812091909155546105bc908363ffffffff610b0616565b600055600160a060020a0381167fcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca58360405190815260200160405180910390a25050565b60035460a060020a900460ff1681565b60035460009060a060020a900460ff161561062a57600080fd5b6104578383610b18565b600160a060020a031660009081526001602052604090205490565b60035433600160a060020a0390811691161461066a57600080fd5b60035460a060020a900460ff161561068157600080fd5b6003805474ff0000000000000000000000000000000000000000191660a060020a1790557f6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff62560405160405180910390a1565b600354600160a060020a031681565b60408051908101604052600381527f52454e0000000000000000000000000000000000000000000000000000000000602082015281565b60035460009060a060020a900460ff161561073357600080fd5b6104578383610c12565b60035460009033600160a060020a0390811691161461075b57600080fd5b6000821161076857600080fd5b600354600160a060020a0316600090815260016020526040902054610793908363ffffffff610b0616565b600354600160a060020a0390811660009081526001602052604080822093909355908516815220546107cb908363ffffffff610d0d16565b600160a060020a03808516600081815260016020526040908190209390935560035490929116907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9085905190815260200160405180910390a350600192915050565b60035460009060a060020a900460ff161561084857600080fd5b6104578383610d1c565b600160a060020a03918216600090815260026020908152604080832093909416825291909152205490565b60035433600160a060020a0390811691161461089857600080fd5b600160a060020a03811615156108ad57600080fd5b600354600160a060020a0380831691167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a36003805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b600160a060020a03338116600081815260026020908152604080832094871680845294909152808220859055909291907f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259085905190815260200160405180910390a350600192915050565b6000600160a060020a038316151561099b57600080fd5b600160a060020a0384166000908152600160205260409020548211156109c057600080fd5b600160a060020a03808516600090815260026020908152604080832033909416835292905220548211156109f357600080fd5b600160a060020a038416600090815260016020526040902054610a1c908363ffffffff610b0616565b600160a060020a038086166000908152600160205260408082209390935590851681522054610a51908363ffffffff610d0d16565b600160a060020a03808516600090815260016020908152604080832094909455878316825260028152838220339093168252919091522054610a99908363ffffffff610b0616565b600160a060020a03808616600081815260026020908152604080832033861684529091529081902093909355908516917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9085905190815260200160405180910390a35060019392505050565b600082821115610b1257fe5b50900390565b600160a060020a03338116600090815260026020908152604080832093861683529290529081205480831115610b7557600160a060020a033381166000908152600260209081526040808320938816835292905290812055610bac565b610b85818463ffffffff610b0616565b600160a060020a033381166000908152600260209081526040808320938916835292905220555b600160a060020a0333811660008181526002602090815260408083209489168084529490915290819020547f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925915190815260200160405180910390a35060019392505050565b6000600160a060020a0383161515610c2957600080fd5b600160a060020a033316600090815260016020526040902054821115610c4e57600080fd5b600160a060020a033316600090815260016020526040902054610c77908363ffffffff610b0616565b600160a060020a033381166000908152600160205260408082209390935590851681522054610cac908363ffffffff610d0d16565b600160a060020a0380851660008181526001602052604090819020939093559133909116907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9085905190815260200160405180910390a350600192915050565b60008282018381101561045757fe5b600160a060020a033381166000908152600260209081526040808320938616835292905290812054610d54908363ffffffff610d0d16565b600160a060020a0333811660008181526002602090815260408083209489168084529490915290819020849055919290917f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92591905190815260200160405180910390a3506001929150505600a165627a7a7230582031765784f180b67cc62a25af5a1233aee1067dbefbb3bb324e66ed72aa7130b70029`
//
//// DeployRepublicToken deploys a new Ethereum contract, binding an instance of RepublicToken to it.
//func DeployRepublicToken(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *RepublicToken, error) {
//	parsed, err := abi.JSON(strings.NewReader(RepublicTokenABI))
//	if err != nil {
//		return common.Address{}, nil, nil, err
//	}
//	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(RepublicTokenBin), backend)
//	if err != nil {
//		return common.Address{}, nil, nil, err
//	}
//	return address, tx, &RepublicToken{RepublicTokenCaller: RepublicTokenCaller{contract: contract}, RepublicTokenTransactor: RepublicTokenTransactor{contract: contract}, RepublicTokenFilterer: RepublicTokenFilterer{contract: contract}}, nil
//}
//
//// RepublicToken is an auto generated Go binding around an Ethereum contract.
//type RepublicToken struct {
//	RepublicTokenCaller     // Read-only binding to the contract
//	RepublicTokenTransactor // Write-only binding to the contract
//	RepublicTokenFilterer   // Log filterer for contract events
//}
//
//// RepublicTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
//type RepublicTokenCaller struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// RepublicTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
//type RepublicTokenTransactor struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// RepublicTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
//type RepublicTokenFilterer struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// RepublicTokenSession is an auto generated Go binding around an Ethereum contract,
//// with pre-set call and transact options.
//type RepublicTokenSession struct {
//	Contract     *RepublicToken    // Generic contract binding to set the session for
//	CallOpts     bind.CallOpts     // Call options to use throughout this session
//	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
//}
//
//// RepublicTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
//// with pre-set call options.
//type RepublicTokenCallerSession struct {
//	Contract *RepublicTokenCaller // Generic contract caller binding to set the session for
//	CallOpts bind.CallOpts        // Call options to use throughout this session
//}
//
//// RepublicTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
//// with pre-set transact options.
//type RepublicTokenTransactorSession struct {
//	Contract     *RepublicTokenTransactor // Generic contract transactor binding to set the session for
//	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
//}
//
//// RepublicTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
//type RepublicTokenRaw struct {
//	Contract *RepublicToken // Generic contract binding to access the raw methods on
//}
//
//// RepublicTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
//type RepublicTokenCallerRaw struct {
//	Contract *RepublicTokenCaller // Generic read-only contract binding to access the raw methods on
//}
//
//// RepublicTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
//type RepublicTokenTransactorRaw struct {
//	Contract *RepublicTokenTransactor // Generic write-only contract binding to access the raw methods on
//}
//
//// NewRepublicToken creates a new instance of RepublicToken, bound to a specific deployed contract.
//func NewRepublicToken(address common.Address, backend bind.ContractBackend) (*RepublicToken, error) {
//	contract, err := bindRepublicToken(address, backend, backend, backend)
//	if err != nil {
//		return nil, err
//	}
//	return &RepublicToken{RepublicTokenCaller: RepublicTokenCaller{contract: contract}, RepublicTokenTransactor: RepublicTokenTransactor{contract: contract}, RepublicTokenFilterer: RepublicTokenFilterer{contract: contract}}, nil
//}
//
//// NewRepublicTokenCaller creates a new read-only instance of RepublicToken, bound to a specific deployed contract.
//func NewRepublicTokenCaller(address common.Address, caller bind.ContractCaller) (*RepublicTokenCaller, error) {
//	contract, err := bindRepublicToken(address, caller, nil, nil)
//	if err != nil {
//		return nil, err
//	}
//	return &RepublicTokenCaller{contract: contract}, nil
//}
//
//// NewRepublicTokenTransactor creates a new write-only instance of RepublicToken, bound to a specific deployed contract.
//func NewRepublicTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*RepublicTokenTransactor, error) {
//	contract, err := bindRepublicToken(address, nil, transactor, nil)
//	if err != nil {
//		return nil, err
//	}
//	return &RepublicTokenTransactor{contract: contract}, nil
//}
//
//// NewRepublicTokenFilterer creates a new log filterer instance of RepublicToken, bound to a specific deployed contract.
//func NewRepublicTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*RepublicTokenFilterer, error) {
//	contract, err := bindRepublicToken(address, nil, nil, filterer)
//	if err != nil {
//		return nil, err
//	}
//	return &RepublicTokenFilterer{contract: contract}, nil
//}
//
//// bindRepublicToken binds a generic wrapper to an already deployed contract.
//func bindRepublicToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
//	parsed, err := abi.JSON(strings.NewReader(RepublicTokenABI))
//	if err != nil {
//		return nil, err
//	}
//	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
//}
//
//// Call invokes the (constant) contract method with params as input values and
//// sets the output to result. The result type might be a single field for simple
//// returns, a slice of interfaces for anonymous returns and a struct for named
//// returns.
//func (_RepublicToken *RepublicTokenRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
//	return _RepublicToken.Contract.RepublicTokenCaller.contract.Call(opts, result, method, params...)
//}
//
//// Transfer initiates a plain transaction to move funds to the contract, calling
//// its default method if one is available.
//func (_RepublicToken *RepublicTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _RepublicToken.Contract.RepublicTokenTransactor.contract.Transfer(opts)
//}
//
//// Transact invokes the (paid) contract method with params as input values.
//func (_RepublicToken *RepublicTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
//	return _RepublicToken.Contract.RepublicTokenTransactor.contract.Transact(opts, method, params...)
//}
//
//// Call invokes the (constant) contract method with params as input values and
//// sets the output to result. The result type might be a single field for simple
//// returns, a slice of interfaces for anonymous returns and a struct for named
//// returns.
//func (_RepublicToken *RepublicTokenCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
//	return _RepublicToken.Contract.contract.Call(opts, result, method, params...)
//}
//
//// Transfer initiates a plain transaction to move funds to the contract, calling
//// its default method if one is available.
//func (_RepublicToken *RepublicTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _RepublicToken.Contract.contract.Transfer(opts)
//}
//
//// Transact invokes the (paid) contract method with params as input values.
//func (_RepublicToken *RepublicTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
//	return _RepublicToken.Contract.contract.Transact(opts, method, params...)
//}
//
//// INITIALSUPPLY is a free data retrieval call binding the contract method 0x2ff2e9dc.
////
//// Solidity: function INITIAL_SUPPLY() constant returns(uint256)
//func (_RepublicToken *RepublicTokenCaller) INITIALSUPPLY(opts *bind.CallOpts) (*big.Int, error) {
//	var (
//		ret0 = new(*big.Int)
//	)
//	out := ret0
//	err := _RepublicToken.contract.Call(opts, out, "INITIAL_SUPPLY")
//	return *ret0, err
//}
//
//// INITIALSUPPLY is a free data retrieval call binding the contract method 0x2ff2e9dc.
////
//// Solidity: function INITIAL_SUPPLY() constant returns(uint256)
//func (_RepublicToken *RepublicTokenSession) INITIALSUPPLY() (*big.Int, error) {
//	return _RepublicToken.Contract.INITIALSUPPLY(&_RepublicToken.CallOpts)
//}
//
//// INITIALSUPPLY is a free data retrieval call binding the contract method 0x2ff2e9dc.
////
//// Solidity: function INITIAL_SUPPLY() constant returns(uint256)
//func (_RepublicToken *RepublicTokenCallerSession) INITIALSUPPLY() (*big.Int, error) {
//	return _RepublicToken.Contract.INITIALSUPPLY(&_RepublicToken.CallOpts)
//}
//
//// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
////
//// Solidity: function allowance(_owner address, _spender address) constant returns(uint256)
//func (_RepublicToken *RepublicTokenCaller) Allowance(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (*big.Int, error) {
//	var (
//		ret0 = new(*big.Int)
//	)
//	out := ret0
//	err := _RepublicToken.contract.Call(opts, out, "allowance", _owner, _spender)
//	return *ret0, err
//}
//
//// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
////
//// Solidity: function allowance(_owner address, _spender address) constant returns(uint256)
//func (_RepublicToken *RepublicTokenSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
//	return _RepublicToken.Contract.Allowance(&_RepublicToken.CallOpts, _owner, _spender)
//}
//
//// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
////
//// Solidity: function allowance(_owner address, _spender address) constant returns(uint256)
//func (_RepublicToken *RepublicTokenCallerSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
//	return _RepublicToken.Contract.Allowance(&_RepublicToken.CallOpts, _owner, _spender)
//}
//
//// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
////
//// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
//func (_RepublicToken *RepublicTokenCaller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
//	var (
//		ret0 = new(*big.Int)
//	)
//	out := ret0
//	err := _RepublicToken.contract.Call(opts, out, "balanceOf", _owner)
//	return *ret0, err
//}
//
//// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
////
//// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
//func (_RepublicToken *RepublicTokenSession) BalanceOf(_owner common.Address) (*big.Int, error) {
//	return _RepublicToken.Contract.BalanceOf(&_RepublicToken.CallOpts, _owner)
//}
//
//// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
////
//// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
//func (_RepublicToken *RepublicTokenCallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
//	return _RepublicToken.Contract.BalanceOf(&_RepublicToken.CallOpts, _owner)
//}
//
//// Balances is a free data retrieval call binding the contract method 0x27e235e3.
////
//// Solidity: function balances( address) constant returns(uint256)
//func (_RepublicToken *RepublicTokenCaller) Balances(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
//	var (
//		ret0 = new(*big.Int)
//	)
//	out := ret0
//	err := _RepublicToken.contract.Call(opts, out, "balances", arg0)
//	return *ret0, err
//}
//
//// Balances is a free data retrieval call binding the contract method 0x27e235e3.
////
//// Solidity: function balances( address) constant returns(uint256)
//func (_RepublicToken *RepublicTokenSession) Balances(arg0 common.Address) (*big.Int, error) {
//	return _RepublicToken.Contract.Balances(&_RepublicToken.CallOpts, arg0)
//}
//
//// Balances is a free data retrieval call binding the contract method 0x27e235e3.
////
//// Solidity: function balances( address) constant returns(uint256)
//func (_RepublicToken *RepublicTokenCallerSession) Balances(arg0 common.Address) (*big.Int, error) {
//	return _RepublicToken.Contract.Balances(&_RepublicToken.CallOpts, arg0)
//}
//
//// Decimals is a free data retrieval call binding the contract method 0x313ce567.
////
//// Solidity: function decimals() constant returns(uint8)
//func (_RepublicToken *RepublicTokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
//	var (
//		ret0 = new(uint8)
//	)
//	out := ret0
//	err := _RepublicToken.contract.Call(opts, out, "decimals")
//	return *ret0, err
//}
//
//// Decimals is a free data retrieval call binding the contract method 0x313ce567.
////
//// Solidity: function decimals() constant returns(uint8)
//func (_RepublicToken *RepublicTokenSession) Decimals() (uint8, error) {
//	return _RepublicToken.Contract.Decimals(&_RepublicToken.CallOpts)
//}
//
//// Decimals is a free data retrieval call binding the contract method 0x313ce567.
////
//// Solidity: function decimals() constant returns(uint8)
//func (_RepublicToken *RepublicTokenCallerSession) Decimals() (uint8, error) {
//	return _RepublicToken.Contract.Decimals(&_RepublicToken.CallOpts)
//}
//
//// Name is a free data retrieval call binding the contract method 0x06fdde03.
////
//// Solidity: function name() constant returns(string)
//func (_RepublicToken *RepublicTokenCaller) Name(opts *bind.CallOpts) (string, error) {
//	var (
//		ret0 = new(string)
//	)
//	out := ret0
//	err := _RepublicToken.contract.Call(opts, out, "name")
//	return *ret0, err
//}
//
//// Name is a free data retrieval call binding the contract method 0x06fdde03.
////
//// Solidity: function name() constant returns(string)
//func (_RepublicToken *RepublicTokenSession) Name() (string, error) {
//	return _RepublicToken.Contract.Name(&_RepublicToken.CallOpts)
//}
//
//// Name is a free data retrieval call binding the contract method 0x06fdde03.
////
//// Solidity: function name() constant returns(string)
//func (_RepublicToken *RepublicTokenCallerSession) Name() (string, error) {
//	return _RepublicToken.Contract.Name(&_RepublicToken.CallOpts)
//}
//
//// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
////
//// Solidity: function owner() constant returns(address)
//func (_RepublicToken *RepublicTokenCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
//	var (
//		ret0 = new(common.Address)
//	)
//	out := ret0
//	err := _RepublicToken.contract.Call(opts, out, "owner")
//	return *ret0, err
//}
//
//// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
////
//// Solidity: function owner() constant returns(address)
//func (_RepublicToken *RepublicTokenSession) Owner() (common.Address, error) {
//	return _RepublicToken.Contract.Owner(&_RepublicToken.CallOpts)
//}
//
//// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
////
//// Solidity: function owner() constant returns(address)
//func (_RepublicToken *RepublicTokenCallerSession) Owner() (common.Address, error) {
//	return _RepublicToken.Contract.Owner(&_RepublicToken.CallOpts)
//}
//
//// Paused is a free data retrieval call binding the contract method 0x5c975abb.
////
//// Solidity: function paused() constant returns(bool)
//func (_RepublicToken *RepublicTokenCaller) Paused(opts *bind.CallOpts) (bool, error) {
//	var (
//		ret0 = new(bool)
//	)
//	out := ret0
//	err := _RepublicToken.contract.Call(opts, out, "paused")
//	return *ret0, err
//}
//
//// Paused is a free data retrieval call binding the contract method 0x5c975abb.
////
//// Solidity: function paused() constant returns(bool)
//func (_RepublicToken *RepublicTokenSession) Paused() (bool, error) {
//	return _RepublicToken.Contract.Paused(&_RepublicToken.CallOpts)
//}
//
//// Paused is a free data retrieval call binding the contract method 0x5c975abb.
////
//// Solidity: function paused() constant returns(bool)
//func (_RepublicToken *RepublicTokenCallerSession) Paused() (bool, error) {
//	return _RepublicToken.Contract.Paused(&_RepublicToken.CallOpts)
//}
//
//// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
////
//// Solidity: function symbol() constant returns(string)
//func (_RepublicToken *RepublicTokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
//	var (
//		ret0 = new(string)
//	)
//	out := ret0
//	err := _RepublicToken.contract.Call(opts, out, "symbol")
//	return *ret0, err
//}
//
//// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
////
//// Solidity: function symbol() constant returns(string)
//func (_RepublicToken *RepublicTokenSession) Symbol() (string, error) {
//	return _RepublicToken.Contract.Symbol(&_RepublicToken.CallOpts)
//}
//
//// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
////
//// Solidity: function symbol() constant returns(string)
//func (_RepublicToken *RepublicTokenCallerSession) Symbol() (string, error) {
//	return _RepublicToken.Contract.Symbol(&_RepublicToken.CallOpts)
//}
//
//// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
////
//// Solidity: function totalSupply() constant returns(uint256)
//func (_RepublicToken *RepublicTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
//	var (
//		ret0 = new(*big.Int)
//	)
//	out := ret0
//	err := _RepublicToken.contract.Call(opts, out, "totalSupply")
//	return *ret0, err
//}
//
//// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
////
//// Solidity: function totalSupply() constant returns(uint256)
//func (_RepublicToken *RepublicTokenSession) TotalSupply() (*big.Int, error) {
//	return _RepublicToken.Contract.TotalSupply(&_RepublicToken.CallOpts)
//}
//
//// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
////
//// Solidity: function totalSupply() constant returns(uint256)
//func (_RepublicToken *RepublicTokenCallerSession) TotalSupply() (*big.Int, error) {
//	return _RepublicToken.Contract.TotalSupply(&_RepublicToken.CallOpts)
//}
//
//// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
////
//// Solidity: function approve(_spender address, _value uint256) returns(bool)
//func (_RepublicToken *RepublicTokenTransactor) Approve(opts *bind.TransactOpts, _spender common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _RepublicToken.contract.Transact(opts, "approve", _spender, _value)
//}
//
//// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
////
//// Solidity: function approve(_spender address, _value uint256) returns(bool)
//func (_RepublicToken *RepublicTokenSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _RepublicToken.Contract.Approve(&_RepublicToken.TransactOpts, _spender, _value)
//}
//
//// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
////
//// Solidity: function approve(_spender address, _value uint256) returns(bool)
//func (_RepublicToken *RepublicTokenTransactorSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _RepublicToken.Contract.Approve(&_RepublicToken.TransactOpts, _spender, _value)
//}
//
//// Burn is a paid mutator transaction binding the contract method 0x42966c68.
////
//// Solidity: function burn(_value uint256) returns()
//func (_RepublicToken *RepublicTokenTransactor) Burn(opts *bind.TransactOpts, _value *big.Int) (*types.Transaction, error) {
//	return _RepublicToken.contract.Transact(opts, "burn", _value)
//}
//
//// Burn is a paid mutator transaction binding the contract method 0x42966c68.
////
//// Solidity: function burn(_value uint256) returns()
//func (_RepublicToken *RepublicTokenSession) Burn(_value *big.Int) (*types.Transaction, error) {
//	return _RepublicToken.Contract.Burn(&_RepublicToken.TransactOpts, _value)
//}
//
//// Burn is a paid mutator transaction binding the contract method 0x42966c68.
////
//// Solidity: function burn(_value uint256) returns()
//func (_RepublicToken *RepublicTokenTransactorSession) Burn(_value *big.Int) (*types.Transaction, error) {
//	return _RepublicToken.Contract.Burn(&_RepublicToken.TransactOpts, _value)
//}
//
//// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
////
//// Solidity: function decreaseApproval(_spender address, _subtractedValue uint256) returns(success bool)
//func (_RepublicToken *RepublicTokenTransactor) DecreaseApproval(opts *bind.TransactOpts, _spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
//	return _RepublicToken.contract.Transact(opts, "decreaseApproval", _spender, _subtractedValue)
//}
//
//// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
////
//// Solidity: function decreaseApproval(_spender address, _subtractedValue uint256) returns(success bool)
//func (_RepublicToken *RepublicTokenSession) DecreaseApproval(_spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
//	return _RepublicToken.Contract.DecreaseApproval(&_RepublicToken.TransactOpts, _spender, _subtractedValue)
//}
//
//// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
////
//// Solidity: function decreaseApproval(_spender address, _subtractedValue uint256) returns(success bool)
//func (_RepublicToken *RepublicTokenTransactorSession) DecreaseApproval(_spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
//	return _RepublicToken.Contract.DecreaseApproval(&_RepublicToken.TransactOpts, _spender, _subtractedValue)
//}
//
//// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
////
//// Solidity: function increaseApproval(_spender address, _addedValue uint256) returns(success bool)
//func (_RepublicToken *RepublicTokenTransactor) IncreaseApproval(opts *bind.TransactOpts, _spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
//	return _RepublicToken.contract.Transact(opts, "increaseApproval", _spender, _addedValue)
//}
//
//// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
////
//// Solidity: function increaseApproval(_spender address, _addedValue uint256) returns(success bool)
//func (_RepublicToken *RepublicTokenSession) IncreaseApproval(_spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
//	return _RepublicToken.Contract.IncreaseApproval(&_RepublicToken.TransactOpts, _spender, _addedValue)
//}
//
//// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
////
//// Solidity: function increaseApproval(_spender address, _addedValue uint256) returns(success bool)
//func (_RepublicToken *RepublicTokenTransactorSession) IncreaseApproval(_spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
//	return _RepublicToken.Contract.IncreaseApproval(&_RepublicToken.TransactOpts, _spender, _addedValue)
//}
//
//// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
////
//// Solidity: function pause() returns()
//func (_RepublicToken *RepublicTokenTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _RepublicToken.contract.Transact(opts, "pause")
//}
//
//// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
////
//// Solidity: function pause() returns()
//func (_RepublicToken *RepublicTokenSession) Pause() (*types.Transaction, error) {
//	return _RepublicToken.Contract.Pause(&_RepublicToken.TransactOpts)
//}
//
//// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
////
//// Solidity: function pause() returns()
//func (_RepublicToken *RepublicTokenTransactorSession) Pause() (*types.Transaction, error) {
//	return _RepublicToken.Contract.Pause(&_RepublicToken.TransactOpts)
//}
//
//// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
////
//// Solidity: function transfer(_to address, _value uint256) returns(bool)
//func (_RepublicToken *RepublicTokenTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _RepublicToken.contract.Transact(opts, "transfer", _to, _value)
//}
//
//// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
////
//// Solidity: function transfer(_to address, _value uint256) returns(bool)
//func (_RepublicToken *RepublicTokenSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _RepublicToken.Contract.Transfer(&_RepublicToken.TransactOpts, _to, _value)
//}
//
//// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
////
//// Solidity: function transfer(_to address, _value uint256) returns(bool)
//func (_RepublicToken *RepublicTokenTransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _RepublicToken.Contract.Transfer(&_RepublicToken.TransactOpts, _to, _value)
//}
//
//// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
////
//// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
//func (_RepublicToken *RepublicTokenTransactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _RepublicToken.contract.Transact(opts, "transferFrom", _from, _to, _value)
//}
//
//// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
////
//// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
//func (_RepublicToken *RepublicTokenSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _RepublicToken.Contract.TransferFrom(&_RepublicToken.TransactOpts, _from, _to, _value)
//}
//
//// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
////
//// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
//func (_RepublicToken *RepublicTokenTransactorSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _RepublicToken.Contract.TransferFrom(&_RepublicToken.TransactOpts, _from, _to, _value)
//}
//
//// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
////
//// Solidity: function transferOwnership(newOwner address) returns()
//func (_RepublicToken *RepublicTokenTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
//	return _RepublicToken.contract.Transact(opts, "transferOwnership", newOwner)
//}
//
//// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
////
//// Solidity: function transferOwnership(newOwner address) returns()
//func (_RepublicToken *RepublicTokenSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
//	return _RepublicToken.Contract.TransferOwnership(&_RepublicToken.TransactOpts, newOwner)
//}
//
//// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
////
//// Solidity: function transferOwnership(newOwner address) returns()
//func (_RepublicToken *RepublicTokenTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
//	return _RepublicToken.Contract.TransferOwnership(&_RepublicToken.TransactOpts, newOwner)
//}
//
//// TransferTokens is a paid mutator transaction binding the contract method 0xbec3fa17.
////
//// Solidity: function transferTokens(beneficiary address, amount uint256) returns(bool)
//func (_RepublicToken *RepublicTokenTransactor) TransferTokens(opts *bind.TransactOpts, beneficiary common.Address, amount *big.Int) (*types.Transaction, error) {
//	return _RepublicToken.contract.Transact(opts, "transferTokens", beneficiary, amount)
//}
//
//// TransferTokens is a paid mutator transaction binding the contract method 0xbec3fa17.
////
//// Solidity: function transferTokens(beneficiary address, amount uint256) returns(bool)
//func (_RepublicToken *RepublicTokenSession) TransferTokens(beneficiary common.Address, amount *big.Int) (*types.Transaction, error) {
//	return _RepublicToken.Contract.TransferTokens(&_RepublicToken.TransactOpts, beneficiary, amount)
//}
//
//// TransferTokens is a paid mutator transaction binding the contract method 0xbec3fa17.
////
//// Solidity: function transferTokens(beneficiary address, amount uint256) returns(bool)
//func (_RepublicToken *RepublicTokenTransactorSession) TransferTokens(beneficiary common.Address, amount *big.Int) (*types.Transaction, error) {
//	return _RepublicToken.Contract.TransferTokens(&_RepublicToken.TransactOpts, beneficiary, amount)
//}
//
//// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
////
//// Solidity: function unpause() returns()
//func (_RepublicToken *RepublicTokenTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _RepublicToken.contract.Transact(opts, "unpause")
//}
//
//// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
////
//// Solidity: function unpause() returns()
//func (_RepublicToken *RepublicTokenSession) Unpause() (*types.Transaction, error) {
//	return _RepublicToken.Contract.Unpause(&_RepublicToken.TransactOpts)
//}
//
//// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
////
//// Solidity: function unpause() returns()
//func (_RepublicToken *RepublicTokenTransactorSession) Unpause() (*types.Transaction, error) {
//	return _RepublicToken.Contract.Unpause(&_RepublicToken.TransactOpts)
//}
//
//// RepublicTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the RepublicToken contract.
//type RepublicTokenApprovalIterator struct {
//	Event *RepublicTokenApproval // Event containing the contract specifics and raw log
//
//	contract *bind.BoundContract // Generic contract to use for unpacking event data
//	event    string              // Event name to use for unpacking event data
//
//	logs chan types.Log        // Log channel receiving the found contract events
//	sub  ethereum.Subscription // Subscription for errors, completion and termination
//	done bool                  // Whether the subscription completed delivering logs
//	fail error                 // Occurred error to stop iteration
//}
//
//// Next advances the iterator to the subsequent event, returning whether there
//// are any more events found. In case of a retrieval or parsing error, false is
//// returned and Error() can be queried for the exact failure.
//func (it *RepublicTokenApprovalIterator) Next() bool {
//	// If the iterator failed, stop iterating
//	if it.fail != nil {
//		return false
//	}
//	// If the iterator completed, deliver directly whatever's available
//	if it.done {
//		select {
//		case log := <-it.logs:
//			it.Event = new(RepublicTokenApproval)
//			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//				it.fail = err
//				return false
//			}
//			it.Event.Raw = log
//			return true
//
//		default:
//			return false
//		}
//	}
//	// Iterator still in progress, wait for either a data or an error event
//	select {
//	case log := <-it.logs:
//		it.Event = new(RepublicTokenApproval)
//		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//			it.fail = err
//			return false
//		}
//		it.Event.Raw = log
//		return true
//
//	case err := <-it.sub.Err():
//		it.done = true
//		it.fail = err
//		return it.Next()
//	}
//}
//
//// Error returns any retrieval or parsing error occurred during filtering.
//func (it *RepublicTokenApprovalIterator) Error() error {
//	return it.fail
//}
//
//// Close terminates the iteration process, releasing any pending underlying
//// resources.
//func (it *RepublicTokenApprovalIterator) Close() error {
//	it.sub.Unsubscribe()
//	return nil
//}
//
//// RepublicTokenApproval represents a Approval event raised by the RepublicToken contract.
//type RepublicTokenApproval struct {
//	Owner   common.Address
//	Spender common.Address
//	Value   *big.Int
//	Raw     types.Log // Blockchain specific contextual infos
//}
//
//// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
////
//// Solidity: event Approval(owner indexed address, spender indexed address, value uint256)
//func (_RepublicToken *RepublicTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*RepublicTokenApprovalIterator, error) {
//
//	var ownerRule []interface{}
//	for _, ownerItem := range owner {
//		ownerRule = append(ownerRule, ownerItem)
//	}
//	var spenderRule []interface{}
//	for _, spenderItem := range spender {
//		spenderRule = append(spenderRule, spenderItem)
//	}
//
//	logs, sub, err := _RepublicToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
//	if err != nil {
//		return nil, err
//	}
//	return &RepublicTokenApprovalIterator{contract: _RepublicToken.contract, event: "Approval", logs: logs, sub: sub}, nil
//}
//
//// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
////
//// Solidity: event Approval(owner indexed address, spender indexed address, value uint256)
//func (_RepublicToken *RepublicTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *RepublicTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {
//
//	var ownerRule []interface{}
//	for _, ownerItem := range owner {
//		ownerRule = append(ownerRule, ownerItem)
//	}
//	var spenderRule []interface{}
//	for _, spenderItem := range spender {
//		spenderRule = append(spenderRule, spenderItem)
//	}
//
//	logs, sub, err := _RepublicToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
//	if err != nil {
//		return nil, err
//	}
//	return event.NewSubscription(func(quit <-chan struct{}) error {
//		defer sub.Unsubscribe()
//		for {
//			select {
//			case log := <-logs:
//				// New log arrived, parse the event and forward to the user
//				event := new(RepublicTokenApproval)
//				if err := _RepublicToken.contract.UnpackLog(event, "Approval", log); err != nil {
//					return err
//				}
//				event.Raw = log
//
//				select {
//				case sink <- event:
//				case err := <-sub.Err():
//					return err
//				case <-quit:
//					return nil
//				}
//			case err := <-sub.Err():
//				return err
//			case <-quit:
//				return nil
//			}
//		}
//	}), nil
//}
//
//// RepublicTokenBurnIterator is returned from FilterBurn and is used to iterate over the raw logs and unpacked data for Burn events raised by the RepublicToken contract.
//type RepublicTokenBurnIterator struct {
//	Event *RepublicTokenBurn // Event containing the contract specifics and raw log
//
//	contract *bind.BoundContract // Generic contract to use for unpacking event data
//	event    string              // Event name to use for unpacking event data
//
//	logs chan types.Log        // Log channel receiving the found contract events
//	sub  ethereum.Subscription // Subscription for errors, completion and termination
//	done bool                  // Whether the subscription completed delivering logs
//	fail error                 // Occurred error to stop iteration
//}
//
//// Next advances the iterator to the subsequent event, returning whether there
//// are any more events found. In case of a retrieval or parsing error, false is
//// returned and Error() can be queried for the exact failure.
//func (it *RepublicTokenBurnIterator) Next() bool {
//	// If the iterator failed, stop iterating
//	if it.fail != nil {
//		return false
//	}
//	// If the iterator completed, deliver directly whatever's available
//	if it.done {
//		select {
//		case log := <-it.logs:
//			it.Event = new(RepublicTokenBurn)
//			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//				it.fail = err
//				return false
//			}
//			it.Event.Raw = log
//			return true
//
//		default:
//			return false
//		}
//	}
//	// Iterator still in progress, wait for either a data or an error event
//	select {
//	case log := <-it.logs:
//		it.Event = new(RepublicTokenBurn)
//		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//			it.fail = err
//			return false
//		}
//		it.Event.Raw = log
//		return true
//
//	case err := <-it.sub.Err():
//		it.done = true
//		it.fail = err
//		return it.Next()
//	}
//}
//
//// Error returns any retrieval or parsing error occurred during filtering.
//func (it *RepublicTokenBurnIterator) Error() error {
//	return it.fail
//}
//
//// Close terminates the iteration process, releasing any pending underlying
//// resources.
//func (it *RepublicTokenBurnIterator) Close() error {
//	it.sub.Unsubscribe()
//	return nil
//}
//
//// RepublicTokenBurn represents a Burn event raised by the RepublicToken contract.
//type RepublicTokenBurn struct {
//	Burner common.Address
//	Value  *big.Int
//	Raw    types.Log // Blockchain specific contextual infos
//}
//
//// FilterBurn is a free log retrieval operation binding the contract event 0xcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5.
////
//// Solidity: event Burn(burner indexed address, value uint256)
//func (_RepublicToken *RepublicTokenFilterer) FilterBurn(opts *bind.FilterOpts, burner []common.Address) (*RepublicTokenBurnIterator, error) {
//
//	var burnerRule []interface{}
//	for _, burnerItem := range burner {
//		burnerRule = append(burnerRule, burnerItem)
//	}
//
//	logs, sub, err := _RepublicToken.contract.FilterLogs(opts, "Burn", burnerRule)
//	if err != nil {
//		return nil, err
//	}
//	return &RepublicTokenBurnIterator{contract: _RepublicToken.contract, event: "Burn", logs: logs, sub: sub}, nil
//}
//
//// WatchBurn is a free log subscription operation binding the contract event 0xcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5.
////
//// Solidity: event Burn(burner indexed address, value uint256)
//func (_RepublicToken *RepublicTokenFilterer) WatchBurn(opts *bind.WatchOpts, sink chan<- *RepublicTokenBurn, burner []common.Address) (event.Subscription, error) {
//
//	var burnerRule []interface{}
//	for _, burnerItem := range burner {
//		burnerRule = append(burnerRule, burnerItem)
//	}
//
//	logs, sub, err := _RepublicToken.contract.WatchLogs(opts, "Burn", burnerRule)
//	if err != nil {
//		return nil, err
//	}
//	return event.NewSubscription(func(quit <-chan struct{}) error {
//		defer sub.Unsubscribe()
//		for {
//			select {
//			case log := <-logs:
//				// New log arrived, parse the event and forward to the user
//				event := new(RepublicTokenBurn)
//				if err := _RepublicToken.contract.UnpackLog(event, "Burn", log); err != nil {
//					return err
//				}
//				event.Raw = log
//
//				select {
//				case sink <- event:
//				case err := <-sub.Err():
//					return err
//				case <-quit:
//					return nil
//				}
//			case err := <-sub.Err():
//				return err
//			case <-quit:
//				return nil
//			}
//		}
//	}), nil
//}
//
//// RepublicTokenOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the RepublicToken contract.
//type RepublicTokenOwnershipTransferredIterator struct {
//	Event *RepublicTokenOwnershipTransferred // Event containing the contract specifics and raw log
//
//	contract *bind.BoundContract // Generic contract to use for unpacking event data
//	event    string              // Event name to use for unpacking event data
//
//	logs chan types.Log        // Log channel receiving the found contract events
//	sub  ethereum.Subscription // Subscription for errors, completion and termination
//	done bool                  // Whether the subscription completed delivering logs
//	fail error                 // Occurred error to stop iteration
//}
//
//// Next advances the iterator to the subsequent event, returning whether there
//// are any more events found. In case of a retrieval or parsing error, false is
//// returned and Error() can be queried for the exact failure.
//func (it *RepublicTokenOwnershipTransferredIterator) Next() bool {
//	// If the iterator failed, stop iterating
//	if it.fail != nil {
//		return false
//	}
//	// If the iterator completed, deliver directly whatever's available
//	if it.done {
//		select {
//		case log := <-it.logs:
//			it.Event = new(RepublicTokenOwnershipTransferred)
//			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//				it.fail = err
//				return false
//			}
//			it.Event.Raw = log
//			return true
//
//		default:
//			return false
//		}
//	}
//	// Iterator still in progress, wait for either a data or an error event
//	select {
//	case log := <-it.logs:
//		it.Event = new(RepublicTokenOwnershipTransferred)
//		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//			it.fail = err
//			return false
//		}
//		it.Event.Raw = log
//		return true
//
//	case err := <-it.sub.Err():
//		it.done = true
//		it.fail = err
//		return it.Next()
//	}
//}
//
//// Error returns any retrieval or parsing error occurred during filtering.
//func (it *RepublicTokenOwnershipTransferredIterator) Error() error {
//	return it.fail
//}
//
//// Close terminates the iteration process, releasing any pending underlying
//// resources.
//func (it *RepublicTokenOwnershipTransferredIterator) Close() error {
//	it.sub.Unsubscribe()
//	return nil
//}
//
//// RepublicTokenOwnershipTransferred represents a OwnershipTransferred event raised by the RepublicToken contract.
//type RepublicTokenOwnershipTransferred struct {
//	PreviousOwner common.Address
//	NewOwner      common.Address
//	Raw           types.Log // Blockchain specific contextual infos
//}
//
//// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
////
//// Solidity: event OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
//func (_RepublicToken *RepublicTokenFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*RepublicTokenOwnershipTransferredIterator, error) {
//
//	var previousOwnerRule []interface{}
//	for _, previousOwnerItem := range previousOwner {
//		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
//	}
//	var newOwnerRule []interface{}
//	for _, newOwnerItem := range newOwner {
//		newOwnerRule = append(newOwnerRule, newOwnerItem)
//	}
//
//	logs, sub, err := _RepublicToken.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
//	if err != nil {
//		return nil, err
//	}
//	return &RepublicTokenOwnershipTransferredIterator{contract: _RepublicToken.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
//}
//
//// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
////
//// Solidity: event OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
//func (_RepublicToken *RepublicTokenFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RepublicTokenOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {
//
//	var previousOwnerRule []interface{}
//	for _, previousOwnerItem := range previousOwner {
//		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
//	}
//	var newOwnerRule []interface{}
//	for _, newOwnerItem := range newOwner {
//		newOwnerRule = append(newOwnerRule, newOwnerItem)
//	}
//
//	logs, sub, err := _RepublicToken.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
//	if err != nil {
//		return nil, err
//	}
//	return event.NewSubscription(func(quit <-chan struct{}) error {
//		defer sub.Unsubscribe()
//		for {
//			select {
//			case log := <-logs:
//				// New log arrived, parse the event and forward to the user
//				event := new(RepublicTokenOwnershipTransferred)
//				if err := _RepublicToken.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
//					return err
//				}
//				event.Raw = log
//
//				select {
//				case sink <- event:
//				case err := <-sub.Err():
//					return err
//				case <-quit:
//					return nil
//				}
//			case err := <-sub.Err():
//				return err
//			case <-quit:
//				return nil
//			}
//		}
//	}), nil
//}
//
//// RepublicTokenPauseIterator is returned from FilterPause and is used to iterate over the raw logs and unpacked data for Pause events raised by the RepublicToken contract.
//type RepublicTokenPauseIterator struct {
//	Event *RepublicTokenPause // Event containing the contract specifics and raw log
//
//	contract *bind.BoundContract // Generic contract to use for unpacking event data
//	event    string              // Event name to use for unpacking event data
//
//	logs chan types.Log        // Log channel receiving the found contract events
//	sub  ethereum.Subscription // Subscription for errors, completion and termination
//	done bool                  // Whether the subscription completed delivering logs
//	fail error                 // Occurred error to stop iteration
//}
//
//// Next advances the iterator to the subsequent event, returning whether there
//// are any more events found. In case of a retrieval or parsing error, false is
//// returned and Error() can be queried for the exact failure.
//func (it *RepublicTokenPauseIterator) Next() bool {
//	// If the iterator failed, stop iterating
//	if it.fail != nil {
//		return false
//	}
//	// If the iterator completed, deliver directly whatever's available
//	if it.done {
//		select {
//		case log := <-it.logs:
//			it.Event = new(RepublicTokenPause)
//			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//				it.fail = err
//				return false
//			}
//			it.Event.Raw = log
//			return true
//
//		default:
//			return false
//		}
//	}
//	// Iterator still in progress, wait for either a data or an error event
//	select {
//	case log := <-it.logs:
//		it.Event = new(RepublicTokenPause)
//		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//			it.fail = err
//			return false
//		}
//		it.Event.Raw = log
//		return true
//
//	case err := <-it.sub.Err():
//		it.done = true
//		it.fail = err
//		return it.Next()
//	}
//}
//
//// Error returns any retrieval or parsing error occurred during filtering.
//func (it *RepublicTokenPauseIterator) Error() error {
//	return it.fail
//}
//
//// Close terminates the iteration process, releasing any pending underlying
//// resources.
//func (it *RepublicTokenPauseIterator) Close() error {
//	it.sub.Unsubscribe()
//	return nil
//}
//
//// RepublicTokenPause represents a Pause event raised by the RepublicToken contract.
//type RepublicTokenPause struct {
//	Raw types.Log // Blockchain specific contextual infos
//}
//
//// FilterPause is a free log retrieval operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
////
//// Solidity: event Pause()
//func (_RepublicToken *RepublicTokenFilterer) FilterPause(opts *bind.FilterOpts) (*RepublicTokenPauseIterator, error) {
//
//	logs, sub, err := _RepublicToken.contract.FilterLogs(opts, "Pause")
//	if err != nil {
//		return nil, err
//	}
//	return &RepublicTokenPauseIterator{contract: _RepublicToken.contract, event: "Pause", logs: logs, sub: sub}, nil
//}
//
//// WatchPause is a free log subscription operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
////
//// Solidity: event Pause()
//func (_RepublicToken *RepublicTokenFilterer) WatchPause(opts *bind.WatchOpts, sink chan<- *RepublicTokenPause) (event.Subscription, error) {
//
//	logs, sub, err := _RepublicToken.contract.WatchLogs(opts, "Pause")
//	if err != nil {
//		return nil, err
//	}
//	return event.NewSubscription(func(quit <-chan struct{}) error {
//		defer sub.Unsubscribe()
//		for {
//			select {
//			case log := <-logs:
//				// New log arrived, parse the event and forward to the user
//				event := new(RepublicTokenPause)
//				if err := _RepublicToken.contract.UnpackLog(event, "Pause", log); err != nil {
//					return err
//				}
//				event.Raw = log
//
//				select {
//				case sink <- event:
//				case err := <-sub.Err():
//					return err
//				case <-quit:
//					return nil
//				}
//			case err := <-sub.Err():
//				return err
//			case <-quit:
//				return nil
//			}
//		}
//	}), nil
//}
//
//// RepublicTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the RepublicToken contract.
//type RepublicTokenTransferIterator struct {
//	Event *RepublicTokenTransfer // Event containing the contract specifics and raw log
//
//	contract *bind.BoundContract // Generic contract to use for unpacking event data
//	event    string              // Event name to use for unpacking event data
//
//	logs chan types.Log        // Log channel receiving the found contract events
//	sub  ethereum.Subscription // Subscription for errors, completion and termination
//	done bool                  // Whether the subscription completed delivering logs
//	fail error                 // Occurred error to stop iteration
//}
//
//// Next advances the iterator to the subsequent event, returning whether there
//// are any more events found. In case of a retrieval or parsing error, false is
//// returned and Error() can be queried for the exact failure.
//func (it *RepublicTokenTransferIterator) Next() bool {
//	// If the iterator failed, stop iterating
//	if it.fail != nil {
//		return false
//	}
//	// If the iterator completed, deliver directly whatever's available
//	if it.done {
//		select {
//		case log := <-it.logs:
//			it.Event = new(RepublicTokenTransfer)
//			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//				it.fail = err
//				return false
//			}
//			it.Event.Raw = log
//			return true
//
//		default:
//			return false
//		}
//	}
//	// Iterator still in progress, wait for either a data or an error event
//	select {
//	case log := <-it.logs:
//		it.Event = new(RepublicTokenTransfer)
//		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//			it.fail = err
//			return false
//		}
//		it.Event.Raw = log
//		return true
//
//	case err := <-it.sub.Err():
//		it.done = true
//		it.fail = err
//		return it.Next()
//	}
//}
//
//// Error returns any retrieval or parsing error occurred during filtering.
//func (it *RepublicTokenTransferIterator) Error() error {
//	return it.fail
//}
//
//// Close terminates the iteration process, releasing any pending underlying
//// resources.
//func (it *RepublicTokenTransferIterator) Close() error {
//	it.sub.Unsubscribe()
//	return nil
//}
//
//// RepublicTokenTransfer represents a Transfer event raised by the RepublicToken contract.
//type RepublicTokenTransfer struct {
//	From  common.Address
//	To    common.Address
//	Value *big.Int
//	Raw   types.Log // Blockchain specific contextual infos
//}
//
//// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
////
//// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
//func (_RepublicToken *RepublicTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RepublicTokenTransferIterator, error) {
//
//	var fromRule []interface{}
//	for _, fromItem := range from {
//		fromRule = append(fromRule, fromItem)
//	}
//	var toRule []interface{}
//	for _, toItem := range to {
//		toRule = append(toRule, toItem)
//	}
//
//	logs, sub, err := _RepublicToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
//	if err != nil {
//		return nil, err
//	}
//	return &RepublicTokenTransferIterator{contract: _RepublicToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
//}
//
//// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
////
//// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
//func (_RepublicToken *RepublicTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *RepublicTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {
//
//	var fromRule []interface{}
//	for _, fromItem := range from {
//		fromRule = append(fromRule, fromItem)
//	}
//	var toRule []interface{}
//	for _, toItem := range to {
//		toRule = append(toRule, toItem)
//	}
//
//	logs, sub, err := _RepublicToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
//	if err != nil {
//		return nil, err
//	}
//	return event.NewSubscription(func(quit <-chan struct{}) error {
//		defer sub.Unsubscribe()
//		for {
//			select {
//			case log := <-logs:
//				// New log arrived, parse the event and forward to the user
//				event := new(RepublicTokenTransfer)
//				if err := _RepublicToken.contract.UnpackLog(event, "Transfer", log); err != nil {
//					return err
//				}
//				event.Raw = log
//
//				select {
//				case sink <- event:
//				case err := <-sub.Err():
//					return err
//				case <-quit:
//					return nil
//				}
//			case err := <-sub.Err():
//				return err
//			case <-quit:
//				return nil
//			}
//		}
//	}), nil
//}
//
//// RepublicTokenUnpauseIterator is returned from FilterUnpause and is used to iterate over the raw logs and unpacked data for Unpause events raised by the RepublicToken contract.
//type RepublicTokenUnpauseIterator struct {
//	Event *RepublicTokenUnpause // Event containing the contract specifics and raw log
//
//	contract *bind.BoundContract // Generic contract to use for unpacking event data
//	event    string              // Event name to use for unpacking event data
//
//	logs chan types.Log        // Log channel receiving the found contract events
//	sub  ethereum.Subscription // Subscription for errors, completion and termination
//	done bool                  // Whether the subscription completed delivering logs
//	fail error                 // Occurred error to stop iteration
//}
//
//// Next advances the iterator to the subsequent event, returning whether there
//// are any more events found. In case of a retrieval or parsing error, false is
//// returned and Error() can be queried for the exact failure.
//func (it *RepublicTokenUnpauseIterator) Next() bool {
//	// If the iterator failed, stop iterating
//	if it.fail != nil {
//		return false
//	}
//	// If the iterator completed, deliver directly whatever's available
//	if it.done {
//		select {
//		case log := <-it.logs:
//			it.Event = new(RepublicTokenUnpause)
//			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//				it.fail = err
//				return false
//			}
//			it.Event.Raw = log
//			return true
//
//		default:
//			return false
//		}
//	}
//	// Iterator still in progress, wait for either a data or an error event
//	select {
//	case log := <-it.logs:
//		it.Event = new(RepublicTokenUnpause)
//		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//			it.fail = err
//			return false
//		}
//		it.Event.Raw = log
//		return true
//
//	case err := <-it.sub.Err():
//		it.done = true
//		it.fail = err
//		return it.Next()
//	}
//}
//
//// Error returns any retrieval or parsing error occurred during filtering.
//func (it *RepublicTokenUnpauseIterator) Error() error {
//	return it.fail
//}
//
//// Close terminates the iteration process, releasing any pending underlying
//// resources.
//func (it *RepublicTokenUnpauseIterator) Close() error {
//	it.sub.Unsubscribe()
//	return nil
//}
//
//// RepublicTokenUnpause represents a Unpause event raised by the RepublicToken contract.
//type RepublicTokenUnpause struct {
//	Raw types.Log // Blockchain specific contextual infos
//}
//
//// FilterUnpause is a free log retrieval operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
////
//// Solidity: event Unpause()
//func (_RepublicToken *RepublicTokenFilterer) FilterUnpause(opts *bind.FilterOpts) (*RepublicTokenUnpauseIterator, error) {
//
//	logs, sub, err := _RepublicToken.contract.FilterLogs(opts, "Unpause")
//	if err != nil {
//		return nil, err
//	}
//	return &RepublicTokenUnpauseIterator{contract: _RepublicToken.contract, event: "Unpause", logs: logs, sub: sub}, nil
//}
//
//// WatchUnpause is a free log subscription operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
////
//// Solidity: event Unpause()
//func (_RepublicToken *RepublicTokenFilterer) WatchUnpause(opts *bind.WatchOpts, sink chan<- *RepublicTokenUnpause) (event.Subscription, error) {
//
//	logs, sub, err := _RepublicToken.contract.WatchLogs(opts, "Unpause")
//	if err != nil {
//		return nil, err
//	}
//	return event.NewSubscription(func(quit <-chan struct{}) error {
//		defer sub.Unsubscribe()
//		for {
//			select {
//			case log := <-logs:
//				// New log arrived, parse the event and forward to the user
//				event := new(RepublicTokenUnpause)
//				if err := _RepublicToken.contract.UnpackLog(event, "Unpause", log); err != nil {
//					return err
//				}
//				event.Raw = log
//
//				select {
//				case sink <- event:
//				case err := <-sub.Err():
//					return err
//				case <-quit:
//					return nil
//				}
//			case err := <-sub.Err():
//				return err
//			case <-quit:
//				return nil
//			}
//		}
//	}), nil
//}
//
//// SafeMathABI is the input ABI used to generate the binding from.
//const SafeMathABI = "[]"
//
//// SafeMathBin is the compiled bytecode used for deploying new contracts.
//const SafeMathBin = `0x604c602c600b82828239805160001a60731460008114601c57601e565bfe5b5030600052607381538281f30073000000000000000000000000000000000000000030146060604052600080fd00a165627a7a72305820ecf558d24646bee08238c4e374bf2bf01ebba0eeae1654560944f826d95ee7d60029`
//
//// DeploySafeMath deploys a new Ethereum contract, binding an instance of SafeMath to it.
//func DeploySafeMath(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SafeMath, error) {
//	parsed, err := abi.JSON(strings.NewReader(SafeMathABI))
//	if err != nil {
//		return common.Address{}, nil, nil, err
//	}
//	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SafeMathBin), backend)
//	if err != nil {
//		return common.Address{}, nil, nil, err
//	}
//	return address, tx, &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
//}
//
//// SafeMath is an auto generated Go binding around an Ethereum contract.
//type SafeMath struct {
//	SafeMathCaller     // Read-only binding to the contract
//	SafeMathTransactor // Write-only binding to the contract
//	SafeMathFilterer   // Log filterer for contract events
//}
//
//// SafeMathCaller is an auto generated read-only Go binding around an Ethereum contract.
//type SafeMathCaller struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// SafeMathTransactor is an auto generated write-only Go binding around an Ethereum contract.
//type SafeMathTransactor struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// SafeMathFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
//type SafeMathFilterer struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// SafeMathSession is an auto generated Go binding around an Ethereum contract,
//// with pre-set call and transact options.
//type SafeMathSession struct {
//	Contract     *SafeMath         // Generic contract binding to set the session for
//	CallOpts     bind.CallOpts     // Call options to use throughout this session
//	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
//}
//
//// SafeMathCallerSession is an auto generated read-only Go binding around an Ethereum contract,
//// with pre-set call options.
//type SafeMathCallerSession struct {
//	Contract *SafeMathCaller // Generic contract caller binding to set the session for
//	CallOpts bind.CallOpts   // Call options to use throughout this session
//}
//
//// SafeMathTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
//// with pre-set transact options.
//type SafeMathTransactorSession struct {
//	Contract     *SafeMathTransactor // Generic contract transactor binding to set the session for
//	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
//}
//
//// SafeMathRaw is an auto generated low-level Go binding around an Ethereum contract.
//type SafeMathRaw struct {
//	Contract *SafeMath // Generic contract binding to access the raw methods on
//}
//
//// SafeMathCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
//type SafeMathCallerRaw struct {
//	Contract *SafeMathCaller // Generic read-only contract binding to access the raw methods on
//}
//
//// SafeMathTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
//type SafeMathTransactorRaw struct {
//	Contract *SafeMathTransactor // Generic write-only contract binding to access the raw methods on
//}
//
//// NewSafeMath creates a new instance of SafeMath, bound to a specific deployed contract.
//func NewSafeMath(address common.Address, backend bind.ContractBackend) (*SafeMath, error) {
//	contract, err := bindSafeMath(address, backend, backend, backend)
//	if err != nil {
//		return nil, err
//	}
//	return &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
//}
//
//// NewSafeMathCaller creates a new read-only instance of SafeMath, bound to a specific deployed contract.
//func NewSafeMathCaller(address common.Address, caller bind.ContractCaller) (*SafeMathCaller, error) {
//	contract, err := bindSafeMath(address, caller, nil, nil)
//	if err != nil {
//		return nil, err
//	}
//	return &SafeMathCaller{contract: contract}, nil
//}
//
//// NewSafeMathTransactor creates a new write-only instance of SafeMath, bound to a specific deployed contract.
//func NewSafeMathTransactor(address common.Address, transactor bind.ContractTransactor) (*SafeMathTransactor, error) {
//	contract, err := bindSafeMath(address, nil, transactor, nil)
//	if err != nil {
//		return nil, err
//	}
//	return &SafeMathTransactor{contract: contract}, nil
//}
//
//// NewSafeMathFilterer creates a new log filterer instance of SafeMath, bound to a specific deployed contract.
//func NewSafeMathFilterer(address common.Address, filterer bind.ContractFilterer) (*SafeMathFilterer, error) {
//	contract, err := bindSafeMath(address, nil, nil, filterer)
//	if err != nil {
//		return nil, err
//	}
//	return &SafeMathFilterer{contract: contract}, nil
//}
//
//// bindSafeMath binds a generic wrapper to an already deployed contract.
//func bindSafeMath(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
//	parsed, err := abi.JSON(strings.NewReader(SafeMathABI))
//	if err != nil {
//		return nil, err
//	}
//	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
//}
//
//// Call invokes the (constant) contract method with params as input values and
//// sets the output to result. The result type might be a single field for simple
//// returns, a slice of interfaces for anonymous returns and a struct for named
//// returns.
//func (_SafeMath *SafeMathRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
//	return _SafeMath.Contract.SafeMathCaller.contract.Call(opts, result, method, params...)
//}
//
//// Transfer initiates a plain transaction to move funds to the contract, calling
//// its default method if one is available.
//func (_SafeMath *SafeMathRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _SafeMath.Contract.SafeMathTransactor.contract.Transfer(opts)
//}
//
//// Transact invokes the (paid) contract method with params as input values.
//func (_SafeMath *SafeMathRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
//	return _SafeMath.Contract.SafeMathTransactor.contract.Transact(opts, method, params...)
//}
//
//// Call invokes the (constant) contract method with params as input values and
//// sets the output to result. The result type might be a single field for simple
//// returns, a slice of interfaces for anonymous returns and a struct for named
//// returns.
//func (_SafeMath *SafeMathCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
//	return _SafeMath.Contract.contract.Call(opts, result, method, params...)
//}
//
//// Transfer initiates a plain transaction to move funds to the contract, calling
//// its default method if one is available.
//func (_SafeMath *SafeMathTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _SafeMath.Contract.contract.Transfer(opts)
//}
//
//// Transact invokes the (paid) contract method with params as input values.
//func (_SafeMath *SafeMathTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
//	return _SafeMath.Contract.contract.Transact(opts, method, params...)
//}
//
//// StandardTokenABI is the input ABI used to generate the binding from.
//const StandardTokenABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseApproval\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseApproval\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]"
//
//// StandardTokenBin is the compiled bytecode used for deploying new contracts.
//const StandardTokenBin = `0x6060604052341561000f57600080fd5b6107228061001e6000396000f3006060604052600436106100985763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663095ea7b3811461009d57806318160ddd146100d357806323b872dd146100f857806327e235e314610120578063661884631461013f57806370a0823114610161578063a9059cbb14610180578063d73dd623146101a2578063dd62ed3e146101c4575b600080fd5b34156100a857600080fd5b6100bf600160a060020a03600435166024356101e9565b604051901515815260200160405180910390f35b34156100de57600080fd5b6100e6610255565b60405190815260200160405180910390f35b341561010357600080fd5b6100bf600160a060020a036004358116906024351660443561025b565b341561012b57600080fd5b6100e6600160a060020a03600435166103dd565b341561014a57600080fd5b6100bf600160a060020a03600435166024356103ef565b341561016c57600080fd5b6100e6600160a060020a03600435166104e9565b341561018b57600080fd5b6100bf600160a060020a0360043516602435610504565b34156101ad57600080fd5b6100bf600160a060020a03600435166024356105ff565b34156101cf57600080fd5b6100e6600160a060020a03600435811690602435166106a3565b600160a060020a03338116600081815260026020908152604080832094871680845294909152808220859055909291907f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259085905190815260200160405180910390a350600192915050565b60005481565b6000600160a060020a038316151561027257600080fd5b600160a060020a03841660009081526001602052604090205482111561029757600080fd5b600160a060020a03808516600090815260026020908152604080832033909416835292905220548211156102ca57600080fd5b600160a060020a0384166000908152600160205260409020546102f3908363ffffffff6106ce16565b600160a060020a038086166000908152600160205260408082209390935590851681522054610328908363ffffffff6106e016565b600160a060020a03808516600090815260016020908152604080832094909455878316825260028152838220339093168252919091522054610370908363ffffffff6106ce16565b600160a060020a03808616600081815260026020908152604080832033861684529091529081902093909355908516917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9085905190815260200160405180910390a35060019392505050565b60016020526000908152604090205481565b600160a060020a0333811660009081526002602090815260408083209386168352929052908120548083111561044c57600160a060020a033381166000908152600260209081526040808320938816835292905290812055610483565b61045c818463ffffffff6106ce16565b600160a060020a033381166000908152600260209081526040808320938916835292905220555b600160a060020a0333811660008181526002602090815260408083209489168084529490915290819020547f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925915190815260200160405180910390a35060019392505050565b600160a060020a031660009081526001602052604090205490565b6000600160a060020a038316151561051b57600080fd5b600160a060020a03331660009081526001602052604090205482111561054057600080fd5b600160a060020a033316600090815260016020526040902054610569908363ffffffff6106ce16565b600160a060020a03338116600090815260016020526040808220939093559085168152205461059e908363ffffffff6106e016565b600160a060020a0380851660008181526001602052604090819020939093559133909116907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9085905190815260200160405180910390a350600192915050565b600160a060020a033381166000908152600260209081526040808320938616835292905290812054610637908363ffffffff6106e016565b600160a060020a0333811660008181526002602090815260408083209489168084529490915290819020849055919290917f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92591905190815260200160405180910390a350600192915050565b600160a060020a03918216600090815260026020908152604080832093909416825291909152205490565b6000828211156106da57fe5b50900390565b6000828201838110156106ef57fe5b93925050505600a165627a7a723058206573655abd9c1785b6bafb0c01e8b67b23c0f5c7d35931a98f2fae5d3ed98cb70029`
//
//// DeployStandardToken deploys a new Ethereum contract, binding an instance of StandardToken to it.
//func DeployStandardToken(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *StandardToken, error) {
//	parsed, err := abi.JSON(strings.NewReader(StandardTokenABI))
//	if err != nil {
//		return common.Address{}, nil, nil, err
//	}
//	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(StandardTokenBin), backend)
//	if err != nil {
//		return common.Address{}, nil, nil, err
//	}
//	return address, tx, &StandardToken{StandardTokenCaller: StandardTokenCaller{contract: contract}, StandardTokenTransactor: StandardTokenTransactor{contract: contract}, StandardTokenFilterer: StandardTokenFilterer{contract: contract}}, nil
//}
//
//// StandardToken is an auto generated Go binding around an Ethereum contract.
//type StandardToken struct {
//	StandardTokenCaller     // Read-only binding to the contract
//	StandardTokenTransactor // Write-only binding to the contract
//	StandardTokenFilterer   // Log filterer for contract events
//}
//
//// StandardTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
//type StandardTokenCaller struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// StandardTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
//type StandardTokenTransactor struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// StandardTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
//type StandardTokenFilterer struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// StandardTokenSession is an auto generated Go binding around an Ethereum contract,
//// with pre-set call and transact options.
//type StandardTokenSession struct {
//	Contract     *StandardToken    // Generic contract binding to set the session for
//	CallOpts     bind.CallOpts     // Call options to use throughout this session
//	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
//}
//
//// StandardTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
//// with pre-set call options.
//type StandardTokenCallerSession struct {
//	Contract *StandardTokenCaller // Generic contract caller binding to set the session for
//	CallOpts bind.CallOpts        // Call options to use throughout this session
//}
//
//// StandardTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
//// with pre-set transact options.
//type StandardTokenTransactorSession struct {
//	Contract     *StandardTokenTransactor // Generic contract transactor binding to set the session for
//	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
//}
//
//// StandardTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
//type StandardTokenRaw struct {
//	Contract *StandardToken // Generic contract binding to access the raw methods on
//}
//
//// StandardTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
//type StandardTokenCallerRaw struct {
//	Contract *StandardTokenCaller // Generic read-only contract binding to access the raw methods on
//}
//
//// StandardTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
//type StandardTokenTransactorRaw struct {
//	Contract *StandardTokenTransactor // Generic write-only contract binding to access the raw methods on
//}
//
//// NewStandardToken creates a new instance of StandardToken, bound to a specific deployed contract.
//func NewStandardToken(address common.Address, backend bind.ContractBackend) (*StandardToken, error) {
//	contract, err := bindStandardToken(address, backend, backend, backend)
//	if err != nil {
//		return nil, err
//	}
//	return &StandardToken{StandardTokenCaller: StandardTokenCaller{contract: contract}, StandardTokenTransactor: StandardTokenTransactor{contract: contract}, StandardTokenFilterer: StandardTokenFilterer{contract: contract}}, nil
//}
//
//// NewStandardTokenCaller creates a new read-only instance of StandardToken, bound to a specific deployed contract.
//func NewStandardTokenCaller(address common.Address, caller bind.ContractCaller) (*StandardTokenCaller, error) {
//	contract, err := bindStandardToken(address, caller, nil, nil)
//	if err != nil {
//		return nil, err
//	}
//	return &StandardTokenCaller{contract: contract}, nil
//}
//
//// NewStandardTokenTransactor creates a new write-only instance of StandardToken, bound to a specific deployed contract.
//func NewStandardTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*StandardTokenTransactor, error) {
//	contract, err := bindStandardToken(address, nil, transactor, nil)
//	if err != nil {
//		return nil, err
//	}
//	return &StandardTokenTransactor{contract: contract}, nil
//}
//
//// NewStandardTokenFilterer creates a new log filterer instance of StandardToken, bound to a specific deployed contract.
//func NewStandardTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*StandardTokenFilterer, error) {
//	contract, err := bindStandardToken(address, nil, nil, filterer)
//	if err != nil {
//		return nil, err
//	}
//	return &StandardTokenFilterer{contract: contract}, nil
//}
//
//// bindStandardToken binds a generic wrapper to an already deployed contract.
//func bindStandardToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
//	parsed, err := abi.JSON(strings.NewReader(StandardTokenABI))
//	if err != nil {
//		return nil, err
//	}
//	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
//}
//
//// Call invokes the (constant) contract method with params as input values and
//// sets the output to result. The result type might be a single field for simple
//// returns, a slice of interfaces for anonymous returns and a struct for named
//// returns.
//func (_StandardToken *StandardTokenRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
//	return _StandardToken.Contract.StandardTokenCaller.contract.Call(opts, result, method, params...)
//}
//
//// Transfer initiates a plain transaction to move funds to the contract, calling
//// its default method if one is available.
//func (_StandardToken *StandardTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _StandardToken.Contract.StandardTokenTransactor.contract.Transfer(opts)
//}
//
//// Transact invokes the (paid) contract method with params as input values.
//func (_StandardToken *StandardTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
//	return _StandardToken.Contract.StandardTokenTransactor.contract.Transact(opts, method, params...)
//}
//
//// Call invokes the (constant) contract method with params as input values and
//// sets the output to result. The result type might be a single field for simple
//// returns, a slice of interfaces for anonymous returns and a struct for named
//// returns.
//func (_StandardToken *StandardTokenCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
//	return _StandardToken.Contract.contract.Call(opts, result, method, params...)
//}
//
//// Transfer initiates a plain transaction to move funds to the contract, calling
//// its default method if one is available.
//func (_StandardToken *StandardTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _StandardToken.Contract.contract.Transfer(opts)
//}
//
//// Transact invokes the (paid) contract method with params as input values.
//func (_StandardToken *StandardTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
//	return _StandardToken.Contract.contract.Transact(opts, method, params...)
//}
//
//// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
////
//// Solidity: function allowance(_owner address, _spender address) constant returns(uint256)
//func (_StandardToken *StandardTokenCaller) Allowance(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (*big.Int, error) {
//	var (
//		ret0 = new(*big.Int)
//	)
//	out := ret0
//	err := _StandardToken.contract.Call(opts, out, "allowance", _owner, _spender)
//	return *ret0, err
//}
//
//// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
////
//// Solidity: function allowance(_owner address, _spender address) constant returns(uint256)
//func (_StandardToken *StandardTokenSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
//	return _StandardToken.Contract.Allowance(&_StandardToken.CallOpts, _owner, _spender)
//}
//
//// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
////
//// Solidity: function allowance(_owner address, _spender address) constant returns(uint256)
//func (_StandardToken *StandardTokenCallerSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
//	return _StandardToken.Contract.Allowance(&_StandardToken.CallOpts, _owner, _spender)
//}
//
//// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
////
//// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
//func (_StandardToken *StandardTokenCaller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
//	var (
//		ret0 = new(*big.Int)
//	)
//	out := ret0
//	err := _StandardToken.contract.Call(opts, out, "balanceOf", _owner)
//	return *ret0, err
//}
//
//// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
////
//// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
//func (_StandardToken *StandardTokenSession) BalanceOf(_owner common.Address) (*big.Int, error) {
//	return _StandardToken.Contract.BalanceOf(&_StandardToken.CallOpts, _owner)
//}
//
//// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
////
//// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
//func (_StandardToken *StandardTokenCallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
//	return _StandardToken.Contract.BalanceOf(&_StandardToken.CallOpts, _owner)
//}
//
//// Balances is a free data retrieval call binding the contract method 0x27e235e3.
////
//// Solidity: function balances( address) constant returns(uint256)
//func (_StandardToken *StandardTokenCaller) Balances(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
//	var (
//		ret0 = new(*big.Int)
//	)
//	out := ret0
//	err := _StandardToken.contract.Call(opts, out, "balances", arg0)
//	return *ret0, err
//}
//
//// Balances is a free data retrieval call binding the contract method 0x27e235e3.
////
//// Solidity: function balances( address) constant returns(uint256)
//func (_StandardToken *StandardTokenSession) Balances(arg0 common.Address) (*big.Int, error) {
//	return _StandardToken.Contract.Balances(&_StandardToken.CallOpts, arg0)
//}
//
//// Balances is a free data retrieval call binding the contract method 0x27e235e3.
////
//// Solidity: function balances( address) constant returns(uint256)
//func (_StandardToken *StandardTokenCallerSession) Balances(arg0 common.Address) (*big.Int, error) {
//	return _StandardToken.Contract.Balances(&_StandardToken.CallOpts, arg0)
//}
//
//// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
////
//// Solidity: function totalSupply() constant returns(uint256)
//func (_StandardToken *StandardTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
//	var (
//		ret0 = new(*big.Int)
//	)
//	out := ret0
//	err := _StandardToken.contract.Call(opts, out, "totalSupply")
//	return *ret0, err
//}
//
//// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
////
//// Solidity: function totalSupply() constant returns(uint256)
//func (_StandardToken *StandardTokenSession) TotalSupply() (*big.Int, error) {
//	return _StandardToken.Contract.TotalSupply(&_StandardToken.CallOpts)
//}
//
//// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
////
//// Solidity: function totalSupply() constant returns(uint256)
//func (_StandardToken *StandardTokenCallerSession) TotalSupply() (*big.Int, error) {
//	return _StandardToken.Contract.TotalSupply(&_StandardToken.CallOpts)
//}
//
//// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
////
//// Solidity: function approve(_spender address, _value uint256) returns(bool)
//func (_StandardToken *StandardTokenTransactor) Approve(opts *bind.TransactOpts, _spender common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _StandardToken.contract.Transact(opts, "approve", _spender, _value)
//}
//
//// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
////
//// Solidity: function approve(_spender address, _value uint256) returns(bool)
//func (_StandardToken *StandardTokenSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _StandardToken.Contract.Approve(&_StandardToken.TransactOpts, _spender, _value)
//}
//
//// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
////
//// Solidity: function approve(_spender address, _value uint256) returns(bool)
//func (_StandardToken *StandardTokenTransactorSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _StandardToken.Contract.Approve(&_StandardToken.TransactOpts, _spender, _value)
//}
//
//// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
////
//// Solidity: function decreaseApproval(_spender address, _subtractedValue uint256) returns(bool)
//func (_StandardToken *StandardTokenTransactor) DecreaseApproval(opts *bind.TransactOpts, _spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
//	return _StandardToken.contract.Transact(opts, "decreaseApproval", _spender, _subtractedValue)
//}
//
//// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
////
//// Solidity: function decreaseApproval(_spender address, _subtractedValue uint256) returns(bool)
//func (_StandardToken *StandardTokenSession) DecreaseApproval(_spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
//	return _StandardToken.Contract.DecreaseApproval(&_StandardToken.TransactOpts, _spender, _subtractedValue)
//}
//
//// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
////
//// Solidity: function decreaseApproval(_spender address, _subtractedValue uint256) returns(bool)
//func (_StandardToken *StandardTokenTransactorSession) DecreaseApproval(_spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
//	return _StandardToken.Contract.DecreaseApproval(&_StandardToken.TransactOpts, _spender, _subtractedValue)
//}
//
//// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
////
//// Solidity: function increaseApproval(_spender address, _addedValue uint256) returns(bool)
//func (_StandardToken *StandardTokenTransactor) IncreaseApproval(opts *bind.TransactOpts, _spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
//	return _StandardToken.contract.Transact(opts, "increaseApproval", _spender, _addedValue)
//}
//
//// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
////
//// Solidity: function increaseApproval(_spender address, _addedValue uint256) returns(bool)
//func (_StandardToken *StandardTokenSession) IncreaseApproval(_spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
//	return _StandardToken.Contract.IncreaseApproval(&_StandardToken.TransactOpts, _spender, _addedValue)
//}
//
//// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
////
//// Solidity: function increaseApproval(_spender address, _addedValue uint256) returns(bool)
//func (_StandardToken *StandardTokenTransactorSession) IncreaseApproval(_spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
//	return _StandardToken.Contract.IncreaseApproval(&_StandardToken.TransactOpts, _spender, _addedValue)
//}
//
//// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
////
//// Solidity: function transfer(_to address, _value uint256) returns(bool)
//func (_StandardToken *StandardTokenTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _StandardToken.contract.Transact(opts, "transfer", _to, _value)
//}
//
//// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
////
//// Solidity: function transfer(_to address, _value uint256) returns(bool)
//func (_StandardToken *StandardTokenSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _StandardToken.Contract.Transfer(&_StandardToken.TransactOpts, _to, _value)
//}
//
//// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
////
//// Solidity: function transfer(_to address, _value uint256) returns(bool)
//func (_StandardToken *StandardTokenTransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _StandardToken.Contract.Transfer(&_StandardToken.TransactOpts, _to, _value)
//}
//
//// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
////
//// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
//func (_StandardToken *StandardTokenTransactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _StandardToken.contract.Transact(opts, "transferFrom", _from, _to, _value)
//}
//
//// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
////
//// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
//func (_StandardToken *StandardTokenSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _StandardToken.Contract.TransferFrom(&_StandardToken.TransactOpts, _from, _to, _value)
//}
//
//// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
////
//// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
//func (_StandardToken *StandardTokenTransactorSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
//	return _StandardToken.Contract.TransferFrom(&_StandardToken.TransactOpts, _from, _to, _value)
//}
//
//// StandardTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the StandardToken contract.
//type StandardTokenApprovalIterator struct {
//	Event *StandardTokenApproval // Event containing the contract specifics and raw log
//
//	contract *bind.BoundContract // Generic contract to use for unpacking event data
//	event    string              // Event name to use for unpacking event data
//
//	logs chan types.Log        // Log channel receiving the found contract events
//	sub  ethereum.Subscription // Subscription for errors, completion and termination
//	done bool                  // Whether the subscription completed delivering logs
//	fail error                 // Occurred error to stop iteration
//}
//
//// Next advances the iterator to the subsequent event, returning whether there
//// are any more events found. In case of a retrieval or parsing error, false is
//// returned and Error() can be queried for the exact failure.
//func (it *StandardTokenApprovalIterator) Next() bool {
//	// If the iterator failed, stop iterating
//	if it.fail != nil {
//		return false
//	}
//	// If the iterator completed, deliver directly whatever's available
//	if it.done {
//		select {
//		case log := <-it.logs:
//			it.Event = new(StandardTokenApproval)
//			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//				it.fail = err
//				return false
//			}
//			it.Event.Raw = log
//			return true
//
//		default:
//			return false
//		}
//	}
//	// Iterator still in progress, wait for either a data or an error event
//	select {
//	case log := <-it.logs:
//		it.Event = new(StandardTokenApproval)
//		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//			it.fail = err
//			return false
//		}
//		it.Event.Raw = log
//		return true
//
//	case err := <-it.sub.Err():
//		it.done = true
//		it.fail = err
//		return it.Next()
//	}
//}
//
//// Error returns any retrieval or parsing error occurred during filtering.
//func (it *StandardTokenApprovalIterator) Error() error {
//	return it.fail
//}
//
//// Close terminates the iteration process, releasing any pending underlying
//// resources.
//func (it *StandardTokenApprovalIterator) Close() error {
//	it.sub.Unsubscribe()
//	return nil
//}
//
//// StandardTokenApproval represents a Approval event raised by the StandardToken contract.
//type StandardTokenApproval struct {
//	Owner   common.Address
//	Spender common.Address
//	Value   *big.Int
//	Raw     types.Log // Blockchain specific contextual infos
//}
//
//// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
////
//// Solidity: event Approval(owner indexed address, spender indexed address, value uint256)
//func (_StandardToken *StandardTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*StandardTokenApprovalIterator, error) {
//
//	var ownerRule []interface{}
//	for _, ownerItem := range owner {
//		ownerRule = append(ownerRule, ownerItem)
//	}
//	var spenderRule []interface{}
//	for _, spenderItem := range spender {
//		spenderRule = append(spenderRule, spenderItem)
//	}
//
//	logs, sub, err := _StandardToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
//	if err != nil {
//		return nil, err
//	}
//	return &StandardTokenApprovalIterator{contract: _StandardToken.contract, event: "Approval", logs: logs, sub: sub}, nil
//}
//
//// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
////
//// Solidity: event Approval(owner indexed address, spender indexed address, value uint256)
//func (_StandardToken *StandardTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *StandardTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {
//
//	var ownerRule []interface{}
//	for _, ownerItem := range owner {
//		ownerRule = append(ownerRule, ownerItem)
//	}
//	var spenderRule []interface{}
//	for _, spenderItem := range spender {
//		spenderRule = append(spenderRule, spenderItem)
//	}
//
//	logs, sub, err := _StandardToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
//	if err != nil {
//		return nil, err
//	}
//	return event.NewSubscription(func(quit <-chan struct{}) error {
//		defer sub.Unsubscribe()
//		for {
//			select {
//			case log := <-logs:
//				// New log arrived, parse the event and forward to the user
//				event := new(StandardTokenApproval)
//				if err := _StandardToken.contract.UnpackLog(event, "Approval", log); err != nil {
//					return err
//				}
//				event.Raw = log
//
//				select {
//				case sink <- event:
//				case err := <-sub.Err():
//					return err
//				case <-quit:
//					return nil
//				}
//			case err := <-sub.Err():
//				return err
//			case <-quit:
//				return nil
//			}
//		}
//	}), nil
//}
//
//// StandardTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the StandardToken contract.
//type StandardTokenTransferIterator struct {
//	Event *StandardTokenTransfer // Event containing the contract specifics and raw log
//
//	contract *bind.BoundContract // Generic contract to use for unpacking event data
//	event    string              // Event name to use for unpacking event data
//
//	logs chan types.Log        // Log channel receiving the found contract events
//	sub  ethereum.Subscription // Subscription for errors, completion and termination
//	done bool                  // Whether the subscription completed delivering logs
//	fail error                 // Occurred error to stop iteration
//}
//
//// Next advances the iterator to the subsequent event, returning whether there
//// are any more events found. In case of a retrieval or parsing error, false is
//// returned and Error() can be queried for the exact failure.
//func (it *StandardTokenTransferIterator) Next() bool {
//	// If the iterator failed, stop iterating
//	if it.fail != nil {
//		return false
//	}
//	// If the iterator completed, deliver directly whatever's available
//	if it.done {
//		select {
//		case log := <-it.logs:
//			it.Event = new(StandardTokenTransfer)
//			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//				it.fail = err
//				return false
//			}
//			it.Event.Raw = log
//			return true
//
//		default:
//			return false
//		}
//	}
//	// Iterator still in progress, wait for either a data or an error event
//	select {
//	case log := <-it.logs:
//		it.Event = new(StandardTokenTransfer)
//		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//			it.fail = err
//			return false
//		}
//		it.Event.Raw = log
//		return true
//
//	case err := <-it.sub.Err():
//		it.done = true
//		it.fail = err
//		return it.Next()
//	}
//}
//
//// Error returns any retrieval or parsing error occurred during filtering.
//func (it *StandardTokenTransferIterator) Error() error {
//	return it.fail
//}
//
//// Close terminates the iteration process, releasing any pending underlying
//// resources.
//func (it *StandardTokenTransferIterator) Close() error {
//	it.sub.Unsubscribe()
//	return nil
//}
//
//// StandardTokenTransfer represents a Transfer event raised by the StandardToken contract.
//type StandardTokenTransfer struct {
//	From  common.Address
//	To    common.Address
//	Value *big.Int
//	Raw   types.Log // Blockchain specific contextual infos
//}
//
//// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
////
//// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
//func (_StandardToken *StandardTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*StandardTokenTransferIterator, error) {
//
//	var fromRule []interface{}
//	for _, fromItem := range from {
//		fromRule = append(fromRule, fromItem)
//	}
//	var toRule []interface{}
//	for _, toItem := range to {
//		toRule = append(toRule, toItem)
//	}
//
//	logs, sub, err := _StandardToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
//	if err != nil {
//		return nil, err
//	}
//	return &StandardTokenTransferIterator{contract: _StandardToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
//}
//
//// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
////
//// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
//func (_StandardToken *StandardTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *StandardTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {
//
//	var fromRule []interface{}
//	for _, fromItem := range from {
//		fromRule = append(fromRule, fromItem)
//	}
//	var toRule []interface{}
//	for _, toItem := range to {
//		toRule = append(toRule, toItem)
//	}
//
//	logs, sub, err := _StandardToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
//	if err != nil {
//		return nil, err
//	}
//	return event.NewSubscription(func(quit <-chan struct{}) error {
//		defer sub.Unsubscribe()
//		for {
//			select {
//			case log := <-logs:
//				// New log arrived, parse the event and forward to the user
//				event := new(StandardTokenTransfer)
//				if err := _StandardToken.contract.UnpackLog(event, "Transfer", log); err != nil {
//					return err
//				}
//				event.Raw = log
//
//				select {
//				case sink <- event:
//				case err := <-sub.Err():
//					return err
//				case <-quit:
//					return nil
//				}
//			case err := <-sub.Err():
//				return err
//			case <-quit:
//				return nil
//			}
//		}
//	}), nil
//}
