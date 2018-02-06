// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package atomic_swap

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Atomic_swapABI is the input ABI used to generate the binding from.
const Atomic_swapABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_refundTime\",\"type\":\"uint256\"},{\"name\":\"_hashedSecret\",\"type\":\"bytes20\"},{\"name\":\"_initiator\",\"type\":\"address\"}],\"name\":\"participate\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_hashedSecret\",\"type\":\"bytes20\"}],\"name\":\"refund\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_secret\",\"type\":\"bytes32\"},{\"name\":\"_hashedSecret\",\"type\":\"bytes20\"}],\"name\":\"redeem\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_refundTime\",\"type\":\"uint256\"},{\"name\":\"_hashedSecret\",\"type\":\"bytes20\"},{\"name\":\"_participant\",\"type\":\"address\"}],\"name\":\"initiate\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes20\"}],\"name\":\"swaps\",\"outputs\":[{\"name\":\"initTimestamp\",\"type\":\"uint256\"},{\"name\":\"refundTime\",\"type\":\"uint256\"},{\"name\":\"hashedSecret\",\"type\":\"bytes20\"},{\"name\":\"secret\",\"type\":\"bytes32\"},{\"name\":\"initiator\",\"type\":\"address\"},{\"name\":\"participant\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"emptied\",\"type\":\"bool\"},{\"name\":\"state\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_refundTime\",\"type\":\"uint256\"}],\"name\":\"Refunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_redeemTime\",\"type\":\"uint256\"}],\"name\":\"Redeemed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_initiator\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_participator\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_hashedSecret\",\"type\":\"bytes20\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Participated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_initTimestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"_refundTime\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"_hashedSecret\",\"type\":\"bytes20\"},{\"indexed\":false,\"name\":\"_participant\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_initiator\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_funds\",\"type\":\"uint256\"}],\"name\":\"Initiated\",\"type\":\"event\"}]"

// Atomic_swapBin is the compiled bytecode used for deploying new contracts.
const Atomic_swapBin = `6060604052341561000f57600080fd5b6112b08061001e6000396000f30060606040526004361061006d576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806327f10ee0146100725780635a8f9b81146100c15780637ec27ba1146100f3578063eb8ae1ed14610132578063f325cc9114610181575b600080fd5b6100bf60048080359060200190919080356bffffffffffffffffffffffff191690602001909190803573ffffffffffffffffffffffffffffffffffffffff1690602001909190505061028f565b005b34156100cc57600080fd5b6100f160048080356bffffffffffffffffffffffff191690602001909190505061060d565b005b34156100fe57600080fd5b6101306004808035600019169060200190919080356bffffffffffffffffffffffff19169060200190919050506109a9565b005b61017f60048080359060200190919080356bffffffffffffffffffffffff191690602001909190803573ffffffffffffffffffffffffffffffffffffffff16906020019091905050610e02565b005b341561018c57600080fd5b6101b160048080356bffffffffffffffffffffffff19169060200190919050506111c3565b604051808a8152602001898152602001886bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200187600019166000191681526020018673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018481526020018315151515815260200182600281111561027357fe5b60ff168152602001995050505050505050505060405180910390f35b816000600281111561029d57fe5b600080836bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200190815260200160002060070160019054906101000a900460ff1660028111156102e957fe5b1415156102f557600080fd5b83600080856bffffffffffffffffffffffff19166bffffffffffffffffffffffff191681526020019081526020016000206001018190555042600080856bffffffffffffffffffffffff19166bffffffffffffffffffffffff191681526020019081526020016000206000018190555033600080856bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200190815260200160002060050160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555081600080856bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200190815260200160002060040160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555034600080856bffffffffffffffffffffffff19166bffffffffffffffffffffffff191681526020019081526020016000206006018190555082600080856bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200190815260200160002060020160006101000a81548173ffffffffffffffffffffffffffffffffffffffff02191690836c01000000000000000000000000900402179055506002600080856bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200190815260200160002060070160016101000a81548160ff0219169083600281111561053d57fe5b02179055507fd2ae53b489af667edae5247fc948af6bab36cb0ec684fb14cc114c95905f4f1982338534604051808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001836bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200182815260200194505050505060405180910390a150505050565b80600080826bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200190815260200160002060010154600080836bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200190815260200160002060000154014211151561068457600080fd5b60001515600080836bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200190815260200160002060070160009054906101000a900460ff1615151415156106d757600080fd5b6002808111156106e357fe5b600080846bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200190815260200160002060070160019054906101000a900460ff16600281111561072f57fe5b14156107fc57600080836bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200190815260200160002060050160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc600080856bffffffffffffffffffffffff19166bffffffffffffffffffffffff19168152602001908152602001600020600601549081150290604051600060405180830381858888f1935050505015156107fb57600080fd5b5b6001600281111561080957fe5b600080846bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200190815260200160002060070160019054906101000a900460ff16600281111561085557fe5b141561092257600080836bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200190815260200160002060040160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc600080856bffffffffffffffffffffffff19166bffffffffffffffffffffffff19168152602001908152602001600020600601549081150290604051600060405180830381858888f19350505050151561092157600080fd5b5b6001600080846bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200190815260200160002060070160006101000a81548160ff0219169083151502179055507f3d2a04f53164bedf9a8a46353305d6b2d2261410406df3b41f99ce6489dc003c426040518082815260200191505060405180910390a15050565b8082816bffffffffffffffffffffffff191660038260006040516020015260405180826000191660001916815260200191505060206040518083038160008661646e5a03f115156109f957600080fd5b5050604051516c01000000000000000000000000026bffffffffffffffffffffffff1916141515610a2957600080fd5b600080836bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200190815260200160002060010154600080846bffffffffffffffffffffffff19166bffffffffffffffffffffffff19168152602001908152602001600020600001540142101515610a9f57600080fd5b60001515600080846bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200190815260200160002060070160009054906101000a900460ff161515141515610af257600080fd5b600280811115610afe57fe5b600080856bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200190815260200160002060070160019054906101000a900460ff166002811115610b4a57fe5b1415610c1757600080846bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200190815260200160002060040160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc600080866bffffffffffffffffffffffff19166bffffffffffffffffffffffff19168152602001908152602001600020600601549081150290604051600060405180830381858888f193505050501515610c1657600080fd5b5b60016002811115610c2457fe5b600080856bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200190815260200160002060070160019054906101000a900460ff166002811115610c7057fe5b1415610d3d57600080846bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200190815260200160002060050160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc600080866bffffffffffffffffffffffff19166bffffffffffffffffffffffff19168152602001908152602001600020600601549081150290604051600060405180830381858888f193505050501515610d3c57600080fd5b5b6001600080856bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200190815260200160002060070160006101000a81548160ff0219169083151502179055507f82498456531a1065f689ba348ce20bda781238c424cf36748dd40bc282831e03426040518082815260200191505060405180910390a183600080856bffffffffffffffffffffffff19166bffffffffffffffffffffffff19168152602001908152602001600020600301816000191690555050505050565b8160006002811115610e1057fe5b600080836bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200190815260200160002060070160019054906101000a900460ff166002811115610e5c57fe5b141515610e6857600080fd5b83600080856bffffffffffffffffffffffff19166bffffffffffffffffffffffff191681526020019081526020016000206001018190555042600080856bffffffffffffffffffffffff19166bffffffffffffffffffffffff191681526020019081526020016000206000018190555082600080856bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200190815260200160002060020160006101000a81548173ffffffffffffffffffffffffffffffffffffffff02191690836c010000000000000000000000009004021790555081600080856bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200190815260200160002060050160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555033600080856bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200190815260200160002060040160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506001600080856bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200190815260200160002060070160016101000a81548160ff0219169083600281111561107857fe5b021790555034600080856bffffffffffffffffffffffff19166bffffffffffffffffffffffff19168152602001908152602001600020600601819055507f8f52c15a8dda5af727677001f0ceb13df57f42198fd5cdde86a06aee80333a4b600080856bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200190815260200160002060000154858585333460405180878152602001868152602001856bffffffffffffffffffffffff19166bffffffffffffffffffffffff191681526020018473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001828152602001965050505050505060405180910390a150505050565b60006020528060005260406000206000915090508060000154908060010154908060020160009054906101000a90046c0100000000000000000000000002908060030154908060040160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060050160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060060154908060070160009054906101000a900460ff16908060070160019054906101000a900460ff169050895600a165627a7a72305820909651ea9c300a1e4b3ebef7a89d8891965c4789b55c22fc63ab3f9aee02eaed0029`

// DeployAtomic_swap deploys a new Ethereum contract, binding an instance of Atomic_swap to it.
func DeployAtomic_swap(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Atomic_swap, error) {
	parsed, err := abi.JSON(strings.NewReader(Atomic_swapABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(Atomic_swapBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Atomic_swap{Atomic_swapCaller: Atomic_swapCaller{contract: contract}, Atomic_swapTransactor: Atomic_swapTransactor{contract: contract}}, nil
}

// Atomic_swap is an auto generated Go binding around an Ethereum contract.
type Atomic_swap struct {
	Atomic_swapCaller     // Read-only binding to the contract
	Atomic_swapTransactor // Write-only binding to the contract
}

// Atomic_swapCaller is an auto generated read-only Go binding around an Ethereum contract.
type Atomic_swapCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Atomic_swapTransactor is an auto generated write-only Go binding around an Ethereum contract.
type Atomic_swapTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Atomic_swapSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Atomic_swapSession struct {
	Contract     *Atomic_swap      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Atomic_swapCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Atomic_swapCallerSession struct {
	Contract *Atomic_swapCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// Atomic_swapTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Atomic_swapTransactorSession struct {
	Contract     *Atomic_swapTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// Atomic_swapRaw is an auto generated low-level Go binding around an Ethereum contract.
type Atomic_swapRaw struct {
	Contract *Atomic_swap // Generic contract binding to access the raw methods on
}

// Atomic_swapCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Atomic_swapCallerRaw struct {
	Contract *Atomic_swapCaller // Generic read-only contract binding to access the raw methods on
}

// Atomic_swapTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Atomic_swapTransactorRaw struct {
	Contract *Atomic_swapTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAtomic_swap creates a new instance of Atomic_swap, bound to a specific deployed contract.
func NewAtomic_swap(address common.Address, backend bind.ContractBackend) (*Atomic_swap, error) {
	contract, err := bindAtomic_swap(address, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Atomic_swap{Atomic_swapCaller: Atomic_swapCaller{contract: contract}, Atomic_swapTransactor: Atomic_swapTransactor{contract: contract}}, nil
}

// NewAtomic_swapCaller creates a new read-only instance of Atomic_swap, bound to a specific deployed contract.
func NewAtomic_swapCaller(address common.Address, caller bind.ContractCaller) (*Atomic_swapCaller, error) {
	contract, err := bindAtomic_swap(address, caller, nil)
	if err != nil {
		return nil, err
	}
	return &Atomic_swapCaller{contract: contract}, nil
}

// NewAtomic_swapTransactor creates a new write-only instance of Atomic_swap, bound to a specific deployed contract.
func NewAtomic_swapTransactor(address common.Address, transactor bind.ContractTransactor) (*Atomic_swapTransactor, error) {
	contract, err := bindAtomic_swap(address, nil, transactor)
	if err != nil {
		return nil, err
	}
	return &Atomic_swapTransactor{contract: contract}, nil
}

// bindAtomic_swap binds a generic wrapper to an already deployed contract.
func bindAtomic_swap(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Atomic_swapABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Atomic_swap *Atomic_swapRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Atomic_swap.Contract.Atomic_swapCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Atomic_swap *Atomic_swapRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Atomic_swap.Contract.Atomic_swapTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Atomic_swap *Atomic_swapRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Atomic_swap.Contract.Atomic_swapTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Atomic_swap *Atomic_swapCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Atomic_swap.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Atomic_swap *Atomic_swapTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Atomic_swap.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Atomic_swap *Atomic_swapTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Atomic_swap.Contract.contract.Transact(opts, method, params...)
}

// Swaps is a free data retrieval call binding the contract method 0xf325cc91.
//
// Solidity: function swaps( bytes20) constant returns(initTimestamp uint256, refundTime uint256, hashedSecret bytes20, secret bytes32, initiator address, participant address, value uint256, emptied bool, state uint8)
func (_Atomic_swap *Atomic_swapCaller) Swaps(opts *bind.CallOpts, arg0 [20]byte) (struct {
	InitTimestamp *big.Int
	RefundTime    *big.Int
	HashedSecret  [20]byte
	Secret        [32]byte
	Initiator     common.Address
	Participant   common.Address
	Value         *big.Int
	Emptied       bool
	State         uint8
}, error) {
	ret := new(struct {
		InitTimestamp *big.Int
		RefundTime    *big.Int
		HashedSecret  [20]byte
		Secret        [32]byte
		Initiator     common.Address
		Participant   common.Address
		Value         *big.Int
		Emptied       bool
		State         uint8
	})
	out := ret
	err := _Atomic_swap.contract.Call(opts, out, "swaps", arg0)
	return *ret, err
}

// Swaps is a free data retrieval call binding the contract method 0xf325cc91.
//
// Solidity: function swaps( bytes20) constant returns(initTimestamp uint256, refundTime uint256, hashedSecret bytes20, secret bytes32, initiator address, participant address, value uint256, emptied bool, state uint8)
func (_Atomic_swap *Atomic_swapSession) Swaps(arg0 [20]byte) (struct {
	InitTimestamp *big.Int
	RefundTime    *big.Int
	HashedSecret  [20]byte
	Secret        [32]byte
	Initiator     common.Address
	Participant   common.Address
	Value         *big.Int
	Emptied       bool
	State         uint8
}, error) {
	return _Atomic_swap.Contract.Swaps(&_Atomic_swap.CallOpts, arg0)
}

// Swaps is a free data retrieval call binding the contract method 0xf325cc91.
//
// Solidity: function swaps( bytes20) constant returns(initTimestamp uint256, refundTime uint256, hashedSecret bytes20, secret bytes32, initiator address, participant address, value uint256, emptied bool, state uint8)
func (_Atomic_swap *Atomic_swapCallerSession) Swaps(arg0 [20]byte) (struct {
	InitTimestamp *big.Int
	RefundTime    *big.Int
	HashedSecret  [20]byte
	Secret        [32]byte
	Initiator     common.Address
	Participant   common.Address
	Value         *big.Int
	Emptied       bool
	State         uint8
}, error) {
	return _Atomic_swap.Contract.Swaps(&_Atomic_swap.CallOpts, arg0)
}

// Initiate is a paid mutator transaction binding the contract method 0xeb8ae1ed.
//
// Solidity: function initiate(_refundTime uint256, _hashedSecret bytes20, _participant address) returns()
func (_Atomic_swap *Atomic_swapTransactor) Initiate(opts *bind.TransactOpts, _refundTime *big.Int, _hashedSecret [20]byte, _participant common.Address) (*types.Transaction, error) {
	return _Atomic_swap.contract.Transact(opts, "initiate", _refundTime, _hashedSecret, _participant)
}

// Initiate is a paid mutator transaction binding the contract method 0xeb8ae1ed.
//
// Solidity: function initiate(_refundTime uint256, _hashedSecret bytes20, _participant address) returns()
func (_Atomic_swap *Atomic_swapSession) Initiate(_refundTime *big.Int, _hashedSecret [20]byte, _participant common.Address) (*types.Transaction, error) {
	return _Atomic_swap.Contract.Initiate(&_Atomic_swap.TransactOpts, _refundTime, _hashedSecret, _participant)
}

// Initiate is a paid mutator transaction binding the contract method 0xeb8ae1ed.
//
// Solidity: function initiate(_refundTime uint256, _hashedSecret bytes20, _participant address) returns()
func (_Atomic_swap *Atomic_swapTransactorSession) Initiate(_refundTime *big.Int, _hashedSecret [20]byte, _participant common.Address) (*types.Transaction, error) {
	return _Atomic_swap.Contract.Initiate(&_Atomic_swap.TransactOpts, _refundTime, _hashedSecret, _participant)
}

// Participate is a paid mutator transaction binding the contract method 0x27f10ee0.
//
// Solidity: function participate(_refundTime uint256, _hashedSecret bytes20, _initiator address) returns()
func (_Atomic_swap *Atomic_swapTransactor) Participate(opts *bind.TransactOpts, _refundTime *big.Int, _hashedSecret [20]byte, _initiator common.Address) (*types.Transaction, error) {
	return _Atomic_swap.contract.Transact(opts, "participate", _refundTime, _hashedSecret, _initiator)
}

// Participate is a paid mutator transaction binding the contract method 0x27f10ee0.
//
// Solidity: function participate(_refundTime uint256, _hashedSecret bytes20, _initiator address) returns()
func (_Atomic_swap *Atomic_swapSession) Participate(_refundTime *big.Int, _hashedSecret [20]byte, _initiator common.Address) (*types.Transaction, error) {
	return _Atomic_swap.Contract.Participate(&_Atomic_swap.TransactOpts, _refundTime, _hashedSecret, _initiator)
}

// Participate is a paid mutator transaction binding the contract method 0x27f10ee0.
//
// Solidity: function participate(_refundTime uint256, _hashedSecret bytes20, _initiator address) returns()
func (_Atomic_swap *Atomic_swapTransactorSession) Participate(_refundTime *big.Int, _hashedSecret [20]byte, _initiator common.Address) (*types.Transaction, error) {
	return _Atomic_swap.Contract.Participate(&_Atomic_swap.TransactOpts, _refundTime, _hashedSecret, _initiator)
}

// Redeem is a paid mutator transaction binding the contract method 0x7ec27ba1.
//
// Solidity: function redeem(_secret bytes32, _hashedSecret bytes20) returns()
func (_Atomic_swap *Atomic_swapTransactor) Redeem(opts *bind.TransactOpts, _secret [32]byte, _hashedSecret [20]byte) (*types.Transaction, error) {
	return _Atomic_swap.contract.Transact(opts, "redeem", _secret, _hashedSecret)
}

// Redeem is a paid mutator transaction binding the contract method 0x7ec27ba1.
//
// Solidity: function redeem(_secret bytes32, _hashedSecret bytes20) returns()
func (_Atomic_swap *Atomic_swapSession) Redeem(_secret [32]byte, _hashedSecret [20]byte) (*types.Transaction, error) {
	return _Atomic_swap.Contract.Redeem(&_Atomic_swap.TransactOpts, _secret, _hashedSecret)
}

// Redeem is a paid mutator transaction binding the contract method 0x7ec27ba1.
//
// Solidity: function redeem(_secret bytes32, _hashedSecret bytes20) returns()
func (_Atomic_swap *Atomic_swapTransactorSession) Redeem(_secret [32]byte, _hashedSecret [20]byte) (*types.Transaction, error) {
	return _Atomic_swap.Contract.Redeem(&_Atomic_swap.TransactOpts, _secret, _hashedSecret)
}

// Refund is a paid mutator transaction binding the contract method 0x5a8f9b81.
//
// Solidity: function refund(_hashedSecret bytes20) returns()
func (_Atomic_swap *Atomic_swapTransactor) Refund(opts *bind.TransactOpts, _hashedSecret [20]byte) (*types.Transaction, error) {
	return _Atomic_swap.contract.Transact(opts, "refund", _hashedSecret)
}

// Refund is a paid mutator transaction binding the contract method 0x5a8f9b81.
//
// Solidity: function refund(_hashedSecret bytes20) returns()
func (_Atomic_swap *Atomic_swapSession) Refund(_hashedSecret [20]byte) (*types.Transaction, error) {
	return _Atomic_swap.Contract.Refund(&_Atomic_swap.TransactOpts, _hashedSecret)
}

// Refund is a paid mutator transaction binding the contract method 0x5a8f9b81.
//
// Solidity: function refund(_hashedSecret bytes20) returns()
func (_Atomic_swap *Atomic_swapTransactorSession) Refund(_hashedSecret [20]byte) (*types.Transaction, error) {
	return _Atomic_swap.Contract.Refund(&_Atomic_swap.TransactOpts, _hashedSecret)
}
