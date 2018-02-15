// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

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

// AtomicSwapEtherABI is the input ABI used to generate the binding from.
const AtomicSwapEtherABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"name\":\"_withdrawTrader\",\"type\":\"address\"},{\"name\":\"_secretLock\",\"type\":\"bytes32\"},{\"name\":\"_timelock\",\"type\":\"uint256\"}],\"name\":\"open\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"name\":\"_secretKey\",\"type\":\"bytes\"}],\"name\":\"close\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"check\",\"outputs\":[{\"name\":\"timelock\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"withdrawTrader\",\"type\":\"address\"},{\"name\":\"secretLock\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"expire\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"checkSecretKey\",\"outputs\":[{\"name\":\"secretKey\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_withdrawTrader\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_secretLock\",\"type\":\"bytes32\"}],\"name\":\"Open\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"Expire\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_secretKey\",\"type\":\"bytes\"}],\"name\":\"Close\",\"type\":\"event\"}]"

// AtomicSwapEtherBin is the compiled bytecode used for deploying new contracts.
const AtomicSwapEtherBin = `0x6060604052341561000f57600080fd5b610bff8061001e6000396000f30060606040526004361061006c5763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416630eed854881146100715780631a26720c14610090578063399e0792146100e6578063c644179814610131578063f200e40414610147575b600080fd5b61008e600435600160a060020a03602435166044356064356101d4565b005b341561009b57600080fd5b61008e600480359060446024803590810190830135806020601f8201819004810201604051908101604052818152929190602084018383808284375094965061039095505050505050565b34156100f157600080fd5b6100fc60043561068a565b6040519384526020840192909252600160a060020a031660408084019190915260608301919091526080909101905180910390f35b341561013c57600080fd5b61008e6004356107bf565b341561015257600080fd5b61015d6004356109a0565b60405160208082528190810183818151815260200191508051906020019080838360005b83811015610199578082015183820152602001610181565b50505050905090810190601f1680156101c65780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6101dc610aea565b846000808281526001602052604090205460ff1660038111156101fb57fe5b1461020557600080fd5b60c06040519081016040528084815260200134815260200133600160a060020a0316815260200186600160a060020a031681526020018560001916815260200160006040518059106102545750595b818152601f19601f830116810160200160405290509052600087815260208190526040902090925082908151815560208201518160010155604082015160028201805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055606082015160038201805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03929092169190911790556080820151600482015560a082015181600501908051610318929160200190610b26565b5050506000868152600160208190526040909120805460ff1916828002179055507f6ed79a08bf5c8a7d4a330df315e4ac386627ecafbe5d2bfd6654237d967b24f3868686604051928352600160a060020a0390911660208301526040808301919091526060909101905180910390a1505050505050565b610398610aea565b82600160008281526001602052604090205460ff1660038111156103b857fe5b146103c257600080fd5b83836002816000604051602001526040518082805190602001908083835b602083106103ff5780518252601f1990920191602091820191016103e0565b6001836020036101000a03801982511681845116808217855250505050505090500191505060206040518083038160008661646e5a03f1151561044157600080fd5b5050604051805160008481526020819052604090206004015414905061046657600080fd5b600086815260208190526040908190209060c09051908101604090815282548252600180840154602080850191909152600280860154600160a060020a039081168587015260038701541660608601526004860154608086015260058601805495969560a088019591948116156101000260001901169190910491601f83018190048102019051908101604052809291908181526020018280546001816001161561010002031660029004801561055e5780601f106105335761010080835404028352916020019161055e565b820191906000526020600020905b81548152906001019060200180831161054157829003601f168201915b50505091909252505050600087815260208190526040902090945060050185805161058d929160200190610b26565b506000868152600160205260409020805460ff191660021790556060840151600160a060020a03166108fc85602001519081150290604051600060405180830381858888f1935050505015156105e257600080fd5b7f692fd10a275135b9a2a2f5819db3d9965a5129ea2ad3640a0156dbce2fc81bdd868660405182815260406020820181815290820183818151815260200191508051906020019080838360005b8381101561064757808201518382015260200161062f565b50505050905090810190601f1680156106745780820380516001836020036101000a031916815260200191505b50935050505060405180910390a1505050505050565b600080600080610698610aea565b600086815260208190526040908190209060c09051908101604090815282548252600180840154602080850191909152600280860154600160a060020a039081168587015260038701541660608601526004860154608086015260058601805495969560a088019591948116156101000260001901169190910491601f8301819004810201905190810160405280929190818152602001828054600181600116156101000203166002900480156107905780601f1061076557610100808354040283529160200191610790565b820191906000526020600020905b81548152906001019060200180831161077357829003601f168201915b505050505081525050905080600001518160200151826060015183608001519450945094509450509193509193565b6107c7610aea565b81600160008281526001602052604090205460ff1660038111156107e757fe5b146107f157600080fd5b600083815260208190526040902054839042101561080e57600080fd5b600084815260208190526040908190209060c09051908101604090815282548252600180840154602080850191909152600280860154600160a060020a039081168587015260038701541660608601526004860154608086015260058601805495969560a088019591948116156101000260001901169190910491601f8301819004810201905190810160405280929190818152602001828054600181600116156101000203166002900480156109065780601f106108db57610100808354040283529160200191610906565b820191906000526020600020905b8154815290600101906020018083116108e957829003601f168201915b5050509190925250505060008581526001602052604090819020805460ff19166003179055909350830151600160a060020a03166108fc84602001519081150290604051600060405180830381858888f19350505050151561096757600080fd5b7fbddd9b693ea862fad6ecf78fd51c065be26fda94d1f3cad3a7d691453a38a7358460405190815260200160405180910390a150505050565b6109a8610ba4565b6109b0610aea565b82600260008281526001602052604090205460ff1660038111156109d057fe5b146109da57600080fd5b600084815260208190526040908190209060c09051908101604090815282548252600180840154602080850191909152600280860154600160a060020a039081168587015260038701541660608601526004860154608086015260058601805495969560a088019591948116156101000260001901169190910491601f830181900481020190519081016040528092919081815260200182805460018160011615610100020316600290048015610ad25780601f10610aa757610100808354040283529160200191610ad2565b820191906000526020600020905b815481529060010190602001808311610ab557829003601f168201915b50505050508152505091508160a00151949350505050565b60c0604051908101604090815260008083526020830181905290820181905260608201819052608082015260a08101610b21610ba4565b905290565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10610b6757805160ff1916838001178555610b94565b82800160010185558215610b94579182015b82811115610b94578251825591602001919060010190610b79565b50610ba0929150610bb6565b5090565b60206040519081016040526000815290565b610bd091905b80821115610ba05760008155600101610bbc565b905600a165627a7a723058203db1239f88567fe3a932105fe77c1e6ac66e514162bc868749712128e1d5f6790029`

// DeployAtomicSwapEther deploys a new Ethereum contract, binding an instance of AtomicSwapEther to it.
func DeployAtomicSwapEther(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *AtomicSwapEther, error) {
	parsed, err := abi.JSON(strings.NewReader(AtomicSwapEtherABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AtomicSwapEtherBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AtomicSwapEther{AtomicSwapEtherCaller: AtomicSwapEtherCaller{contract: contract}, AtomicSwapEtherTransactor: AtomicSwapEtherTransactor{contract: contract}, AtomicSwapEtherFilterer: AtomicSwapEtherFilterer{contract: contract}}, nil
}

// AtomicSwapEther is an auto generated Go binding around an Ethereum contract.
type AtomicSwapEther struct {
	AtomicSwapEtherCaller     // Read-only binding to the contract
	AtomicSwapEtherTransactor // Write-only binding to the contract
	AtomicSwapEtherFilterer   // Log filterer for contract events
}

// AtomicSwapEtherCaller is an auto generated read-only Go binding around an Ethereum contract.
type AtomicSwapEtherCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomicSwapEtherTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AtomicSwapEtherTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomicSwapEtherFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AtomicSwapEtherFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomicSwapEtherSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AtomicSwapEtherSession struct {
	Contract     *AtomicSwapEther  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AtomicSwapEtherCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AtomicSwapEtherCallerSession struct {
	Contract *AtomicSwapEtherCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// AtomicSwapEtherTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AtomicSwapEtherTransactorSession struct {
	Contract     *AtomicSwapEtherTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// AtomicSwapEtherRaw is an auto generated low-level Go binding around an Ethereum contract.
type AtomicSwapEtherRaw struct {
	Contract *AtomicSwapEther // Generic contract binding to access the raw methods on
}

// AtomicSwapEtherCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AtomicSwapEtherCallerRaw struct {
	Contract *AtomicSwapEtherCaller // Generic read-only contract binding to access the raw methods on
}

// AtomicSwapEtherTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AtomicSwapEtherTransactorRaw struct {
	Contract *AtomicSwapEtherTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAtomicSwapEther creates a new instance of AtomicSwapEther, bound to a specific deployed contract.
func NewAtomicSwapEther(address common.Address, backend bind.ContractBackend) (*AtomicSwapEther, error) {
	contract, err := bindAtomicSwapEther(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AtomicSwapEther{AtomicSwapEtherCaller: AtomicSwapEtherCaller{contract: contract}, AtomicSwapEtherTransactor: AtomicSwapEtherTransactor{contract: contract}, AtomicSwapEtherFilterer: AtomicSwapEtherFilterer{contract: contract}}, nil
}

// NewAtomicSwapEtherCaller creates a new read-only instance of AtomicSwapEther, bound to a specific deployed contract.
func NewAtomicSwapEtherCaller(address common.Address, caller bind.ContractCaller) (*AtomicSwapEtherCaller, error) {
	contract, err := bindAtomicSwapEther(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AtomicSwapEtherCaller{contract: contract}, nil
}

// NewAtomicSwapEtherTransactor creates a new write-only instance of AtomicSwapEther, bound to a specific deployed contract.
func NewAtomicSwapEtherTransactor(address common.Address, transactor bind.ContractTransactor) (*AtomicSwapEtherTransactor, error) {
	contract, err := bindAtomicSwapEther(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AtomicSwapEtherTransactor{contract: contract}, nil
}

// NewAtomicSwapEtherFilterer creates a new log filterer instance of AtomicSwapEther, bound to a specific deployed contract.
func NewAtomicSwapEtherFilterer(address common.Address, filterer bind.ContractFilterer) (*AtomicSwapEtherFilterer, error) {
	contract, err := bindAtomicSwapEther(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AtomicSwapEtherFilterer{contract: contract}, nil
}

// bindAtomicSwapEther binds a generic wrapper to an already deployed contract.
func bindAtomicSwapEther(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AtomicSwapEtherABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AtomicSwapEther *AtomicSwapEtherRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AtomicSwapEther.Contract.AtomicSwapEtherCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AtomicSwapEther *AtomicSwapEtherRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AtomicSwapEther.Contract.AtomicSwapEtherTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AtomicSwapEther *AtomicSwapEtherRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AtomicSwapEther.Contract.AtomicSwapEtherTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AtomicSwapEther *AtomicSwapEtherCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AtomicSwapEther.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AtomicSwapEther *AtomicSwapEtherTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AtomicSwapEther.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AtomicSwapEther *AtomicSwapEtherTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AtomicSwapEther.Contract.contract.Transact(opts, method, params...)
}

// Check is a free data retrieval call binding the contract method 0x399e0792.
//
// Solidity: function check(_swapID bytes32) constant returns(timelock uint256, value uint256, withdrawTrader address, secretLock bytes32)
func (_AtomicSwapEther *AtomicSwapEtherCaller) Check(opts *bind.CallOpts, _swapID [32]byte) (struct {
	Timelock       *big.Int
	Value          *big.Int
	WithdrawTrader common.Address
	SecretLock     [32]byte
}, error) {
	ret := new(struct {
		Timelock       *big.Int
		Value          *big.Int
		WithdrawTrader common.Address
		SecretLock     [32]byte
	})
	out := ret
	err := _AtomicSwapEther.contract.Call(opts, out, "check", _swapID)
	return *ret, err
}

// Check is a free data retrieval call binding the contract method 0x399e0792.
//
// Solidity: function check(_swapID bytes32) constant returns(timelock uint256, value uint256, withdrawTrader address, secretLock bytes32)
func (_AtomicSwapEther *AtomicSwapEtherSession) Check(_swapID [32]byte) (struct {
	Timelock       *big.Int
	Value          *big.Int
	WithdrawTrader common.Address
	SecretLock     [32]byte
}, error) {
	return _AtomicSwapEther.Contract.Check(&_AtomicSwapEther.CallOpts, _swapID)
}

// Check is a free data retrieval call binding the contract method 0x399e0792.
//
// Solidity: function check(_swapID bytes32) constant returns(timelock uint256, value uint256, withdrawTrader address, secretLock bytes32)
func (_AtomicSwapEther *AtomicSwapEtherCallerSession) Check(_swapID [32]byte) (struct {
	Timelock       *big.Int
	Value          *big.Int
	WithdrawTrader common.Address
	SecretLock     [32]byte
}, error) {
	return _AtomicSwapEther.Contract.Check(&_AtomicSwapEther.CallOpts, _swapID)
}

// CheckSecretKey is a free data retrieval call binding the contract method 0xf200e404.
//
// Solidity: function checkSecretKey(_swapID bytes32) constant returns(secretKey bytes)
func (_AtomicSwapEther *AtomicSwapEtherCaller) CheckSecretKey(opts *bind.CallOpts, _swapID [32]byte) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _AtomicSwapEther.contract.Call(opts, out, "checkSecretKey", _swapID)
	return *ret0, err
}

// CheckSecretKey is a free data retrieval call binding the contract method 0xf200e404.
//
// Solidity: function checkSecretKey(_swapID bytes32) constant returns(secretKey bytes)
func (_AtomicSwapEther *AtomicSwapEtherSession) CheckSecretKey(_swapID [32]byte) ([]byte, error) {
	return _AtomicSwapEther.Contract.CheckSecretKey(&_AtomicSwapEther.CallOpts, _swapID)
}

// CheckSecretKey is a free data retrieval call binding the contract method 0xf200e404.
//
// Solidity: function checkSecretKey(_swapID bytes32) constant returns(secretKey bytes)
func (_AtomicSwapEther *AtomicSwapEtherCallerSession) CheckSecretKey(_swapID [32]byte) ([]byte, error) {
	return _AtomicSwapEther.Contract.CheckSecretKey(&_AtomicSwapEther.CallOpts, _swapID)
}

// Close is a paid mutator transaction binding the contract method 0x1a26720c.
//
// Solidity: function close(_swapID bytes32, _secretKey bytes) returns()
func (_AtomicSwapEther *AtomicSwapEtherTransactor) Close(opts *bind.TransactOpts, _swapID [32]byte, _secretKey []byte) (*types.Transaction, error) {
	return _AtomicSwapEther.contract.Transact(opts, "close", _swapID, _secretKey)
}

// Close is a paid mutator transaction binding the contract method 0x1a26720c.
//
// Solidity: function close(_swapID bytes32, _secretKey bytes) returns()
func (_AtomicSwapEther *AtomicSwapEtherSession) Close(_swapID [32]byte, _secretKey []byte) (*types.Transaction, error) {
	return _AtomicSwapEther.Contract.Close(&_AtomicSwapEther.TransactOpts, _swapID, _secretKey)
}

// Close is a paid mutator transaction binding the contract method 0x1a26720c.
//
// Solidity: function close(_swapID bytes32, _secretKey bytes) returns()
func (_AtomicSwapEther *AtomicSwapEtherTransactorSession) Close(_swapID [32]byte, _secretKey []byte) (*types.Transaction, error) {
	return _AtomicSwapEther.Contract.Close(&_AtomicSwapEther.TransactOpts, _swapID, _secretKey)
}

// Expire is a paid mutator transaction binding the contract method 0xc6441798.
//
// Solidity: function expire(_swapID bytes32) returns()
func (_AtomicSwapEther *AtomicSwapEtherTransactor) Expire(opts *bind.TransactOpts, _swapID [32]byte) (*types.Transaction, error) {
	return _AtomicSwapEther.contract.Transact(opts, "expire", _swapID)
}

// Expire is a paid mutator transaction binding the contract method 0xc6441798.
//
// Solidity: function expire(_swapID bytes32) returns()
func (_AtomicSwapEther *AtomicSwapEtherSession) Expire(_swapID [32]byte) (*types.Transaction, error) {
	return _AtomicSwapEther.Contract.Expire(&_AtomicSwapEther.TransactOpts, _swapID)
}

// Expire is a paid mutator transaction binding the contract method 0xc6441798.
//
// Solidity: function expire(_swapID bytes32) returns()
func (_AtomicSwapEther *AtomicSwapEtherTransactorSession) Expire(_swapID [32]byte) (*types.Transaction, error) {
	return _AtomicSwapEther.Contract.Expire(&_AtomicSwapEther.TransactOpts, _swapID)
}

// Open is a paid mutator transaction binding the contract method 0x0eed8548.
//
// Solidity: function open(_swapID bytes32, _withdrawTrader address, _secretLock bytes32, _timelock uint256) returns()
func (_AtomicSwapEther *AtomicSwapEtherTransactor) Open(opts *bind.TransactOpts, _swapID [32]byte, _withdrawTrader common.Address, _secretLock [32]byte, _timelock *big.Int) (*types.Transaction, error) {
	return _AtomicSwapEther.contract.Transact(opts, "open", _swapID, _withdrawTrader, _secretLock, _timelock)
}

// Open is a paid mutator transaction binding the contract method 0x0eed8548.
//
// Solidity: function open(_swapID bytes32, _withdrawTrader address, _secretLock bytes32, _timelock uint256) returns()
func (_AtomicSwapEther *AtomicSwapEtherSession) Open(_swapID [32]byte, _withdrawTrader common.Address, _secretLock [32]byte, _timelock *big.Int) (*types.Transaction, error) {
	return _AtomicSwapEther.Contract.Open(&_AtomicSwapEther.TransactOpts, _swapID, _withdrawTrader, _secretLock, _timelock)
}

// Open is a paid mutator transaction binding the contract method 0x0eed8548.
//
// Solidity: function open(_swapID bytes32, _withdrawTrader address, _secretLock bytes32, _timelock uint256) returns()
func (_AtomicSwapEther *AtomicSwapEtherTransactorSession) Open(_swapID [32]byte, _withdrawTrader common.Address, _secretLock [32]byte, _timelock *big.Int) (*types.Transaction, error) {
	return _AtomicSwapEther.Contract.Open(&_AtomicSwapEther.TransactOpts, _swapID, _withdrawTrader, _secretLock, _timelock)
}

// AtomicSwapEtherCloseIterator is returned from FilterClose and is used to iterate over the raw logs and unpacked data for Close events raised by the AtomicSwapEther contract.
type AtomicSwapEtherCloseIterator struct {
	Event *AtomicSwapEtherClose // Event containing the contract specifics and raw log

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
func (it *AtomicSwapEtherCloseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtomicSwapEtherClose)
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
		it.Event = new(AtomicSwapEtherClose)
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
func (it *AtomicSwapEtherCloseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtomicSwapEtherCloseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtomicSwapEtherClose represents a Close event raised by the AtomicSwapEther contract.
type AtomicSwapEtherClose struct {
	SwapID    [32]byte
	SecretKey []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterClose is a free log retrieval operation binding the contract event 0x692fd10a275135b9a2a2f5819db3d9965a5129ea2ad3640a0156dbce2fc81bdd.
//
// Solidity: event Close(_swapID bytes32, _secretKey bytes)
func (_AtomicSwapEther *AtomicSwapEtherFilterer) FilterClose(opts *bind.FilterOpts) (*AtomicSwapEtherCloseIterator, error) {

	logs, sub, err := _AtomicSwapEther.contract.FilterLogs(opts, "Close")
	if err != nil {
		return nil, err
	}
	return &AtomicSwapEtherCloseIterator{contract: _AtomicSwapEther.contract, event: "Close", logs: logs, sub: sub}, nil
}

// WatchClose is a free log subscription operation binding the contract event 0x692fd10a275135b9a2a2f5819db3d9965a5129ea2ad3640a0156dbce2fc81bdd.
//
// Solidity: event Close(_swapID bytes32, _secretKey bytes)
func (_AtomicSwapEther *AtomicSwapEtherFilterer) WatchClose(opts *bind.WatchOpts, sink chan<- *AtomicSwapEtherClose) (event.Subscription, error) {

	logs, sub, err := _AtomicSwapEther.contract.WatchLogs(opts, "Close")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtomicSwapEtherClose)
				if err := _AtomicSwapEther.contract.UnpackLog(event, "Close", log); err != nil {
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

// AtomicSwapEtherExpireIterator is returned from FilterExpire and is used to iterate over the raw logs and unpacked data for Expire events raised by the AtomicSwapEther contract.
type AtomicSwapEtherExpireIterator struct {
	Event *AtomicSwapEtherExpire // Event containing the contract specifics and raw log

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
func (it *AtomicSwapEtherExpireIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtomicSwapEtherExpire)
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
		it.Event = new(AtomicSwapEtherExpire)
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
func (it *AtomicSwapEtherExpireIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtomicSwapEtherExpireIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtomicSwapEtherExpire represents a Expire event raised by the AtomicSwapEther contract.
type AtomicSwapEtherExpire struct {
	SwapID [32]byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterExpire is a free log retrieval operation binding the contract event 0xbddd9b693ea862fad6ecf78fd51c065be26fda94d1f3cad3a7d691453a38a735.
//
// Solidity: event Expire(_swapID bytes32)
func (_AtomicSwapEther *AtomicSwapEtherFilterer) FilterExpire(opts *bind.FilterOpts) (*AtomicSwapEtherExpireIterator, error) {

	logs, sub, err := _AtomicSwapEther.contract.FilterLogs(opts, "Expire")
	if err != nil {
		return nil, err
	}
	return &AtomicSwapEtherExpireIterator{contract: _AtomicSwapEther.contract, event: "Expire", logs: logs, sub: sub}, nil
}

// WatchExpire is a free log subscription operation binding the contract event 0xbddd9b693ea862fad6ecf78fd51c065be26fda94d1f3cad3a7d691453a38a735.
//
// Solidity: event Expire(_swapID bytes32)
func (_AtomicSwapEther *AtomicSwapEtherFilterer) WatchExpire(opts *bind.WatchOpts, sink chan<- *AtomicSwapEtherExpire) (event.Subscription, error) {

	logs, sub, err := _AtomicSwapEther.contract.WatchLogs(opts, "Expire")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtomicSwapEtherExpire)
				if err := _AtomicSwapEther.contract.UnpackLog(event, "Expire", log); err != nil {
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

// AtomicSwapEtherOpenIterator is returned from FilterOpen and is used to iterate over the raw logs and unpacked data for Open events raised by the AtomicSwapEther contract.
type AtomicSwapEtherOpenIterator struct {
	Event *AtomicSwapEtherOpen // Event containing the contract specifics and raw log

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
func (it *AtomicSwapEtherOpenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtomicSwapEtherOpen)
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
		it.Event = new(AtomicSwapEtherOpen)
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
func (it *AtomicSwapEtherOpenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtomicSwapEtherOpenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtomicSwapEtherOpen represents a Open event raised by the AtomicSwapEther contract.
type AtomicSwapEtherOpen struct {
	SwapID         [32]byte
	WithdrawTrader common.Address
	SecretLock     [32]byte
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterOpen is a free log retrieval operation binding the contract event 0x6ed79a08bf5c8a7d4a330df315e4ac386627ecafbe5d2bfd6654237d967b24f3.
//
// Solidity: event Open(_swapID bytes32, _withdrawTrader address, _secretLock bytes32)
func (_AtomicSwapEther *AtomicSwapEtherFilterer) FilterOpen(opts *bind.FilterOpts) (*AtomicSwapEtherOpenIterator, error) {

	logs, sub, err := _AtomicSwapEther.contract.FilterLogs(opts, "Open")
	if err != nil {
		return nil, err
	}
	return &AtomicSwapEtherOpenIterator{contract: _AtomicSwapEther.contract, event: "Open", logs: logs, sub: sub}, nil
}

// WatchOpen is a free log subscription operation binding the contract event 0x6ed79a08bf5c8a7d4a330df315e4ac386627ecafbe5d2bfd6654237d967b24f3.
//
// Solidity: event Open(_swapID bytes32, _withdrawTrader address, _secretLock bytes32)
func (_AtomicSwapEther *AtomicSwapEtherFilterer) WatchOpen(opts *bind.WatchOpts, sink chan<- *AtomicSwapEtherOpen) (event.Subscription, error) {

	logs, sub, err := _AtomicSwapEther.contract.WatchLogs(opts, "Open")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtomicSwapEtherOpen)
				if err := _AtomicSwapEther.contract.UnpackLog(event, "Open", log); err != nil {
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
