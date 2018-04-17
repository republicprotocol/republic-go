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

// AtomicSwapERC20ABI is the input ABI used to generate the binding from.
const AtomicSwapERC20ABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"name\":\"_secretKey\",\"type\":\"bytes\"}],\"name\":\"close\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"check\",\"outputs\":[{\"name\":\"timelock\",\"type\":\"uint256\"},{\"name\":\"erc20Value\",\"type\":\"uint256\"},{\"name\":\"erc20ContractAddress\",\"type\":\"address\"},{\"name\":\"withdrawTrader\",\"type\":\"address\"},{\"name\":\"secretLock\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"name\":\"_erc20Value\",\"type\":\"uint256\"},{\"name\":\"_erc20ContractAddress\",\"type\":\"address\"},{\"name\":\"_withdrawTrader\",\"type\":\"address\"},{\"name\":\"_secretLock\",\"type\":\"bytes32\"},{\"name\":\"_timelock\",\"type\":\"uint256\"}],\"name\":\"open\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"expire\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"checkSecretKey\",\"outputs\":[{\"name\":\"secretKey\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_withdrawTrader\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_secretLock\",\"type\":\"bytes32\"}],\"name\":\"Open\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"Expire\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_secretKey\",\"type\":\"bytes\"}],\"name\":\"Close\",\"type\":\"event\"}]"

// AtomicSwapERC20Bin is the compiled bytecode used for deploying new contracts.
const AtomicSwapERC20Bin = `0x6060604052341561000f57600080fd5b610e808061001e6000396000f3006060604052600436106100535763ffffffff60e060020a6000350416631a26720c8114610058578063399e0792146100b05780635d44011e14610104578063c644179814610135578063f200e4041461014b575b600080fd5b341561006357600080fd5b6100ae600480359060446024803590810190830135806020601f820181900481020160405190810160405281815292919060208401838380828437509496506101d895505050505050565b005b34156100bb57600080fd5b6100c6600435610538565b6040519485526020850193909352600160a060020a0391821660408086019190915291166060840152608083019190915260a0909101905180910390f35b341561010f57600080fd5b6100ae600435602435600160a060020a036044358116906064351660843560a435610683565b341561014057600080fd5b6100ae6004356109c7565b341561015657600080fd5b610161600435610c0e565b60405160208082528190810183818151815260200191508051906020019080838360005b8381101561019d578082015183820152602001610185565b50505050905090810190601f1680156101ca5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6101e0610d64565b600083600160008281526001602052604090205460ff16600381111561020257fe5b1461020c57600080fd5b84846002816000604051602001526040518082805190602001908083835b602083106102495780518252601f19909201916020918201910161022a565b6001836020036101000a03801982511681845116808217855250505050505090500191505060206040518083038160008661646e5a03f1151561028b57600080fd5b505060405180516000848152602081905260409020600501541490506102b057600080fd5b600087815260208190526040908190209060e09051908101604090815282548252600180840154602080850191909152600280860154600160a060020a03908116858701526003870154811660608701526004870154166080860152600586015460a086015260068601805495969560c088019591948116156101000260001901169190910491601f8301819004810201905190810160405280929190818152602001828054600181600116156101000203166002900480156103b45780601f10610389576101008083540402835291602001916103b4565b820191906000526020600020905b81548152906001019060200180831161039757829003601f168201915b5050509190925250505060008881526020819052604090209095506006018680516103e3929160200190610da7565b506000878152600160205260409020805460ff191660021790556060850151935083600160a060020a031663a9059cbb8660800151876020015160006040516020015260405160e060020a63ffffffff8516028152600160a060020a0390921660048301526024820152604401602060405180830381600087803b151561046957600080fd5b6102c65a03f1151561047a57600080fd5b50505060405180519050151561048f57600080fd5b7f692fd10a275135b9a2a2f5819db3d9965a5129ea2ad3640a0156dbce2fc81bdd878760405182815260406020820181815290820183818151815260200191508051906020019080838360005b838110156104f45780820151838201526020016104dc565b50505050905090810190601f1680156105215780820380516001836020036101000a031916815260200191505b50935050505060405180910390a150505050505050565b6000806000806000610548610d64565b600087815260208190526040908190209060e09051908101604090815282548252600180840154602080850191909152600280860154600160a060020a03908116858701526003870154811660608701526004870154166080860152600586015460a086015260068601805495969560c088019591948116156101000260001901169190910491601f83018190048102019051908101604052809291908181526020018280546001816001161561010002031660029004801561064c5780601f106106215761010080835404028352916020019161064c565b820191906000526020600020905b81548152906001019060200180831161062f57829003601f168201915b505050505081525050905080600001518160200151826060015183608001518460a00151939b929a50909850965090945092505050565b600061068d610d64565b876000808281526001602052604090205460ff1660038111156106ac57fe5b146106b657600080fd5b60008981526001602052604081205460ff1660038111156106d357fe5b146106dd57600080fd5b86925082600160a060020a031663dd62ed3e333060006040516020015260405160e060020a63ffffffff8516028152600160a060020a03928316600482015291166024820152604401602060405180830381600087803b151561073f57600080fd5b6102c65a03f1151561075057600080fd5b5050506040518051891115905061076657600080fd5b82600160a060020a03166323b872dd33308b60006040516020015260405160e060020a63ffffffff8616028152600160a060020a0393841660048201529190921660248201526044810191909152606401602060405180830381600087803b15156107d057600080fd5b6102c65a03f115156107e157600080fd5b5050506040518051905015156107f657600080fd5b60e06040519081016040528085815260200189815260200133600160a060020a0316815260200188600160a060020a0316815260200187600160a060020a031681526020018660001916815260200160006040518059106108545750595b818152601f19601f83011681016020016040529050905260008a815260208190526040902090925082908151815560208201518160010155604082015160028201805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055606082015160038201805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055608082015160048201805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a039290921691909117905560a0820151600582015560c08201518160060190805161094c929160200190610da7565b5050506000898152600160208190526040909120805460ff1916828002179055507f6ed79a08bf5c8a7d4a330df315e4ac386627ecafbe5d2bfd6654237d967b24f3898787604051928352600160a060020a0390911660208301526040808301919091526060909101905180910390a1505050505050505050565b6109cf610d64565b600082600160008281526001602052604090205460ff1660038111156109f157fe5b146109fb57600080fd5b600084815260208190526040902054849042901115610a1957600080fd5b600085815260208190526040908190209060e09051908101604090815282548252600180840154602080850191909152600280860154600160a060020a03908116858701526003870154811660608701526004870154166080860152600586015460a086015260068601805495969560c088019591948116156101000260001901169190910491601f830181900481020190519081016040528092919081815260200182805460018160011615610100020316600290048015610b1d5780601f10610af257610100808354040283529160200191610b1d565b820191906000526020600020905b815481529060010190602001808311610b0057829003601f168201915b505050919092525050506000868152600160205260409020805460ff1916600317905593506060840151925082600160a060020a031663a9059cbb8560400151866020015160006040516020015260405160e060020a63ffffffff8516028152600160a060020a0390921660048301526024820152604401602060405180830381600087803b1515610bae57600080fd5b6102c65a03f11515610bbf57600080fd5b505050604051805190501515610bd457600080fd5b7fbddd9b693ea862fad6ecf78fd51c065be26fda94d1f3cad3a7d691453a38a7358560405190815260200160405180910390a15050505050565b610c16610e25565b610c1e610d64565b82600260008281526001602052604090205460ff166003811115610c3e57fe5b14610c4857600080fd5b600084815260208190526040908190209060e09051908101604090815282548252600180840154602080850191909152600280860154600160a060020a03908116858701526003870154811660608701526004870154166080860152600586015460a086015260068601805495969560c088019591948116156101000260001901169190910491601f830181900481020190519081016040528092919081815260200182805460018160011615610100020316600290048015610d4c5780601f10610d2157610100808354040283529160200191610d4c565b820191906000526020600020905b815481529060010190602001808311610d2f57829003601f168201915b50505050508152505091508160c00151949350505050565b60e06040519081016040908152600080835260208301819052908201819052606082018190526080820181905260a082015260c08101610da2610e25565b905290565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10610de857805160ff1916838001178555610e15565b82800160010185558215610e15579182015b82811115610e15578251825591602001919060010190610dfa565b50610e21929150610e37565b5090565b60206040519081016040526000815290565b610e5191905b80821115610e215760008155600101610e3d565b905600a165627a7a723058207764410e28c28bfc24ac500e1b987db7436c5d68441661261a187be104caa5320029`

// DeployAtomicSwapERC20 deploys a new Ethereum contract, binding an instance of AtomicSwapERC20 to it.
func DeployAtomicSwapERC20(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *AtomicSwapERC20, error) {
	parsed, err := abi.JSON(strings.NewReader(AtomicSwapERC20ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AtomicSwapERC20Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AtomicSwapERC20{AtomicSwapERC20Caller: AtomicSwapERC20Caller{contract: contract}, AtomicSwapERC20Transactor: AtomicSwapERC20Transactor{contract: contract}, AtomicSwapERC20Filterer: AtomicSwapERC20Filterer{contract: contract}}, nil
}

// AtomicSwapERC20 is an auto generated Go binding around an Ethereum contract.
type AtomicSwapERC20 struct {
	AtomicSwapERC20Caller     // Read-only binding to the contract
	AtomicSwapERC20Transactor // Write-only binding to the contract
	AtomicSwapERC20Filterer   // Log filterer for contract events
}

// AtomicSwapERC20Caller is an auto generated read-only Go binding around an Ethereum contract.
type AtomicSwapERC20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomicSwapERC20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type AtomicSwapERC20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomicSwapERC20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AtomicSwapERC20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomicSwapERC20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AtomicSwapERC20Session struct {
	Contract     *AtomicSwapERC20  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AtomicSwapERC20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AtomicSwapERC20CallerSession struct {
	Contract *AtomicSwapERC20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// AtomicSwapERC20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AtomicSwapERC20TransactorSession struct {
	Contract     *AtomicSwapERC20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// AtomicSwapERC20Raw is an auto generated low-level Go binding around an Ethereum contract.
type AtomicSwapERC20Raw struct {
	Contract *AtomicSwapERC20 // Generic contract binding to access the raw methods on
}

// AtomicSwapERC20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AtomicSwapERC20CallerRaw struct {
	Contract *AtomicSwapERC20Caller // Generic read-only contract binding to access the raw methods on
}

// AtomicSwapERC20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AtomicSwapERC20TransactorRaw struct {
	Contract *AtomicSwapERC20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewAtomicSwapERC20 creates a new instance of AtomicSwapERC20, bound to a specific deployed contract.
func NewAtomicSwapERC20(address common.Address, backend bind.ContractBackend) (*AtomicSwapERC20, error) {
	contract, err := bindAtomicSwapERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AtomicSwapERC20{AtomicSwapERC20Caller: AtomicSwapERC20Caller{contract: contract}, AtomicSwapERC20Transactor: AtomicSwapERC20Transactor{contract: contract}, AtomicSwapERC20Filterer: AtomicSwapERC20Filterer{contract: contract}}, nil
}

// NewAtomicSwapERC20Caller creates a new read-only instance of AtomicSwapERC20, bound to a specific deployed contract.
func NewAtomicSwapERC20Caller(address common.Address, caller bind.ContractCaller) (*AtomicSwapERC20Caller, error) {
	contract, err := bindAtomicSwapERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AtomicSwapERC20Caller{contract: contract}, nil
}

// NewAtomicSwapERC20Transactor creates a new write-only instance of AtomicSwapERC20, bound to a specific deployed contract.
func NewAtomicSwapERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*AtomicSwapERC20Transactor, error) {
	contract, err := bindAtomicSwapERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AtomicSwapERC20Transactor{contract: contract}, nil
}

// NewAtomicSwapERC20Filterer creates a new log filterer instance of AtomicSwapERC20, bound to a specific deployed contract.
func NewAtomicSwapERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*AtomicSwapERC20Filterer, error) {
	contract, err := bindAtomicSwapERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AtomicSwapERC20Filterer{contract: contract}, nil
}

// bindAtomicSwapERC20 binds a generic wrapper to an already deployed contract.
func bindAtomicSwapERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AtomicSwapERC20ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AtomicSwapERC20 *AtomicSwapERC20Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AtomicSwapERC20.Contract.AtomicSwapERC20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AtomicSwapERC20 *AtomicSwapERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AtomicSwapERC20.Contract.AtomicSwapERC20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AtomicSwapERC20 *AtomicSwapERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AtomicSwapERC20.Contract.AtomicSwapERC20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AtomicSwapERC20 *AtomicSwapERC20CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AtomicSwapERC20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AtomicSwapERC20 *AtomicSwapERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AtomicSwapERC20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AtomicSwapERC20 *AtomicSwapERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AtomicSwapERC20.Contract.contract.Transact(opts, method, params...)
}

// Check is a free data retrieval call binding the contract method 0x399e0792.
//
// Solidity: function check(_swapID bytes32) constant returns(timelock uint256, erc20Value uint256, erc20ContractAddress address, withdrawTrader address, secretLock bytes32)
func (_AtomicSwapERC20 *AtomicSwapERC20Caller) Check(opts *bind.CallOpts, _swapID [32]byte) (struct {
	Timelock             *big.Int
	Erc20Value           *big.Int
	Erc20ContractAddress common.Address
	WithdrawTrader       common.Address
	SecretLock           [32]byte
}, error) {
	ret := new(struct {
		Timelock             *big.Int
		Erc20Value           *big.Int
		Erc20ContractAddress common.Address
		WithdrawTrader       common.Address
		SecretLock           [32]byte
	})
	out := ret
	err := _AtomicSwapERC20.contract.Call(opts, out, "check", _swapID)
	return *ret, err
}

// Check is a free data retrieval call binding the contract method 0x399e0792.
//
// Solidity: function check(_swapID bytes32) constant returns(timelock uint256, erc20Value uint256, erc20ContractAddress address, withdrawTrader address, secretLock bytes32)
func (_AtomicSwapERC20 *AtomicSwapERC20Session) Check(_swapID [32]byte) (struct {
	Timelock             *big.Int
	Erc20Value           *big.Int
	Erc20ContractAddress common.Address
	WithdrawTrader       common.Address
	SecretLock           [32]byte
}, error) {
	return _AtomicSwapERC20.Contract.Check(&_AtomicSwapERC20.CallOpts, _swapID)
}

// Check is a free data retrieval call binding the contract method 0x399e0792.
//
// Solidity: function check(_swapID bytes32) constant returns(timelock uint256, erc20Value uint256, erc20ContractAddress address, withdrawTrader address, secretLock bytes32)
func (_AtomicSwapERC20 *AtomicSwapERC20CallerSession) Check(_swapID [32]byte) (struct {
	Timelock             *big.Int
	Erc20Value           *big.Int
	Erc20ContractAddress common.Address
	WithdrawTrader       common.Address
	SecretLock           [32]byte
}, error) {
	return _AtomicSwapERC20.Contract.Check(&_AtomicSwapERC20.CallOpts, _swapID)
}

// CheckSecretKey is a free data retrieval call binding the contract method 0xf200e404.
//
// Solidity: function checkSecretKey(_swapID bytes32) constant returns(secretKey bytes)
func (_AtomicSwapERC20 *AtomicSwapERC20Caller) CheckSecretKey(opts *bind.CallOpts, _swapID [32]byte) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _AtomicSwapERC20.contract.Call(opts, out, "checkSecretKey", _swapID)
	return *ret0, err
}

// CheckSecretKey is a free data retrieval call binding the contract method 0xf200e404.
//
// Solidity: function checkSecretKey(_swapID bytes32) constant returns(secretKey bytes)
func (_AtomicSwapERC20 *AtomicSwapERC20Session) CheckSecretKey(_swapID [32]byte) ([]byte, error) {
	return _AtomicSwapERC20.Contract.CheckSecretKey(&_AtomicSwapERC20.CallOpts, _swapID)
}

// CheckSecretKey is a free data retrieval call binding the contract method 0xf200e404.
//
// Solidity: function checkSecretKey(_swapID bytes32) constant returns(secretKey bytes)
func (_AtomicSwapERC20 *AtomicSwapERC20CallerSession) CheckSecretKey(_swapID [32]byte) ([]byte, error) {
	return _AtomicSwapERC20.Contract.CheckSecretKey(&_AtomicSwapERC20.CallOpts, _swapID)
}

// Close is a paid mutator transaction binding the contract method 0x1a26720c.
//
// Solidity: function close(_swapID bytes32, _secretKey bytes) returns()
func (_AtomicSwapERC20 *AtomicSwapERC20Transactor) Close(opts *bind.TransactOpts, _swapID [32]byte, _secretKey []byte) (*types.Transaction, error) {
	return _AtomicSwapERC20.contract.Transact(opts, "close", _swapID, _secretKey)
}

// Close is a paid mutator transaction binding the contract method 0x1a26720c.
//
// Solidity: function close(_swapID bytes32, _secretKey bytes) returns()
func (_AtomicSwapERC20 *AtomicSwapERC20Session) Close(_swapID [32]byte, _secretKey []byte) (*types.Transaction, error) {
	return _AtomicSwapERC20.Contract.Close(&_AtomicSwapERC20.TransactOpts, _swapID, _secretKey)
}

// Close is a paid mutator transaction binding the contract method 0x1a26720c.
//
// Solidity: function close(_swapID bytes32, _secretKey bytes) returns()
func (_AtomicSwapERC20 *AtomicSwapERC20TransactorSession) Close(_swapID [32]byte, _secretKey []byte) (*types.Transaction, error) {
	return _AtomicSwapERC20.Contract.Close(&_AtomicSwapERC20.TransactOpts, _swapID, _secretKey)
}

// Expire is a paid mutator transaction binding the contract method 0xc6441798.
//
// Solidity: function expire(_swapID bytes32) returns()
func (_AtomicSwapERC20 *AtomicSwapERC20Transactor) Expire(opts *bind.TransactOpts, _swapID [32]byte) (*types.Transaction, error) {
	return _AtomicSwapERC20.contract.Transact(opts, "expire", _swapID)
}

// Expire is a paid mutator transaction binding the contract method 0xc6441798.
//
// Solidity: function expire(_swapID bytes32) returns()
func (_AtomicSwapERC20 *AtomicSwapERC20Session) Expire(_swapID [32]byte) (*types.Transaction, error) {
	return _AtomicSwapERC20.Contract.Expire(&_AtomicSwapERC20.TransactOpts, _swapID)
}

// Expire is a paid mutator transaction binding the contract method 0xc6441798.
//
// Solidity: function expire(_swapID bytes32) returns()
func (_AtomicSwapERC20 *AtomicSwapERC20TransactorSession) Expire(_swapID [32]byte) (*types.Transaction, error) {
	return _AtomicSwapERC20.Contract.Expire(&_AtomicSwapERC20.TransactOpts, _swapID)
}

// Open is a paid mutator transaction binding the contract method 0x5d44011e.
//
// Solidity: function open(_swapID bytes32, _erc20Value uint256, _erc20ContractAddress address, _withdrawTrader address, _secretLock bytes32, _timelock uint256) returns()
func (_AtomicSwapERC20 *AtomicSwapERC20Transactor) Open(opts *bind.TransactOpts, _swapID [32]byte, _erc20Value *big.Int, _erc20ContractAddress common.Address, _withdrawTrader common.Address, _secretLock [32]byte, _timelock *big.Int) (*types.Transaction, error) {
	return _AtomicSwapERC20.contract.Transact(opts, "open", _swapID, _erc20Value, _erc20ContractAddress, _withdrawTrader, _secretLock, _timelock)
}

// Open is a paid mutator transaction binding the contract method 0x5d44011e.
//
// Solidity: function open(_swapID bytes32, _erc20Value uint256, _erc20ContractAddress address, _withdrawTrader address, _secretLock bytes32, _timelock uint256) returns()
func (_AtomicSwapERC20 *AtomicSwapERC20Session) Open(_swapID [32]byte, _erc20Value *big.Int, _erc20ContractAddress common.Address, _withdrawTrader common.Address, _secretLock [32]byte, _timelock *big.Int) (*types.Transaction, error) {
	return _AtomicSwapERC20.Contract.Open(&_AtomicSwapERC20.TransactOpts, _swapID, _erc20Value, _erc20ContractAddress, _withdrawTrader, _secretLock, _timelock)
}

// Open is a paid mutator transaction binding the contract method 0x5d44011e.
//
// Solidity: function open(_swapID bytes32, _erc20Value uint256, _erc20ContractAddress address, _withdrawTrader address, _secretLock bytes32, _timelock uint256) returns()
func (_AtomicSwapERC20 *AtomicSwapERC20TransactorSession) Open(_swapID [32]byte, _erc20Value *big.Int, _erc20ContractAddress common.Address, _withdrawTrader common.Address, _secretLock [32]byte, _timelock *big.Int) (*types.Transaction, error) {
	return _AtomicSwapERC20.Contract.Open(&_AtomicSwapERC20.TransactOpts, _swapID, _erc20Value, _erc20ContractAddress, _withdrawTrader, _secretLock, _timelock)
}

// AtomicSwapERC20CloseIterator is returned from FilterClose and is used to iterate over the raw logs and unpacked data for Close events raised by the AtomicSwapERC20 contract.
type AtomicSwapERC20CloseIterator struct {
	Event *AtomicSwapERC20Close // Event containing the contract specifics and raw log

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
func (it *AtomicSwapERC20CloseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtomicSwapERC20Close)
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
		it.Event = new(AtomicSwapERC20Close)
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
func (it *AtomicSwapERC20CloseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtomicSwapERC20CloseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtomicSwapERC20Close represents a Close event raised by the AtomicSwapERC20 contract.
type AtomicSwapERC20Close struct {
	SwapID    [32]byte
	SecretKey []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterClose is a free log retrieval operation binding the contract event 0x692fd10a275135b9a2a2f5819db3d9965a5129ea2ad3640a0156dbce2fc81bdd.
//
// Solidity: event Close(_swapID bytes32, _secretKey bytes)
func (_AtomicSwapERC20 *AtomicSwapERC20Filterer) FilterClose(opts *bind.FilterOpts) (*AtomicSwapERC20CloseIterator, error) {

	logs, sub, err := _AtomicSwapERC20.contract.FilterLogs(opts, "Close")
	if err != nil {
		return nil, err
	}
	return &AtomicSwapERC20CloseIterator{contract: _AtomicSwapERC20.contract, event: "Close", logs: logs, sub: sub}, nil
}

// WatchClose is a free log subscription operation binding the contract event 0x692fd10a275135b9a2a2f5819db3d9965a5129ea2ad3640a0156dbce2fc81bdd.
//
// Solidity: event Close(_swapID bytes32, _secretKey bytes)
func (_AtomicSwapERC20 *AtomicSwapERC20Filterer) WatchClose(opts *bind.WatchOpts, sink chan<- *AtomicSwapERC20Close) (event.Subscription, error) {

	logs, sub, err := _AtomicSwapERC20.contract.WatchLogs(opts, "Close")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtomicSwapERC20Close)
				if err := _AtomicSwapERC20.contract.UnpackLog(event, "Close", log); err != nil {
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

// AtomicSwapERC20ExpireIterator is returned from FilterExpire and is used to iterate over the raw logs and unpacked data for Expire events raised by the AtomicSwapERC20 contract.
type AtomicSwapERC20ExpireIterator struct {
	Event *AtomicSwapERC20Expire // Event containing the contract specifics and raw log

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
func (it *AtomicSwapERC20ExpireIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtomicSwapERC20Expire)
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
		it.Event = new(AtomicSwapERC20Expire)
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
func (it *AtomicSwapERC20ExpireIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtomicSwapERC20ExpireIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtomicSwapERC20Expire represents a Expire event raised by the AtomicSwapERC20 contract.
type AtomicSwapERC20Expire struct {
	SwapID [32]byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterExpire is a free log retrieval operation binding the contract event 0xbddd9b693ea862fad6ecf78fd51c065be26fda94d1f3cad3a7d691453a38a735.
//
// Solidity: event Expire(_swapID bytes32)
func (_AtomicSwapERC20 *AtomicSwapERC20Filterer) FilterExpire(opts *bind.FilterOpts) (*AtomicSwapERC20ExpireIterator, error) {

	logs, sub, err := _AtomicSwapERC20.contract.FilterLogs(opts, "Expire")
	if err != nil {
		return nil, err
	}
	return &AtomicSwapERC20ExpireIterator{contract: _AtomicSwapERC20.contract, event: "Expire", logs: logs, sub: sub}, nil
}

// WatchExpire is a free log subscription operation binding the contract event 0xbddd9b693ea862fad6ecf78fd51c065be26fda94d1f3cad3a7d691453a38a735.
//
// Solidity: event Expire(_swapID bytes32)
func (_AtomicSwapERC20 *AtomicSwapERC20Filterer) WatchExpire(opts *bind.WatchOpts, sink chan<- *AtomicSwapERC20Expire) (event.Subscription, error) {

	logs, sub, err := _AtomicSwapERC20.contract.WatchLogs(opts, "Expire")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtomicSwapERC20Expire)
				if err := _AtomicSwapERC20.contract.UnpackLog(event, "Expire", log); err != nil {
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

// AtomicSwapERC20OpenIterator is returned from FilterOpen and is used to iterate over the raw logs and unpacked data for Open events raised by the AtomicSwapERC20 contract.
type AtomicSwapERC20OpenIterator struct {
	Event *AtomicSwapERC20Open // Event containing the contract specifics and raw log

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
func (it *AtomicSwapERC20OpenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtomicSwapERC20Open)
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
		it.Event = new(AtomicSwapERC20Open)
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
func (it *AtomicSwapERC20OpenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtomicSwapERC20OpenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtomicSwapERC20Open represents a Open event raised by the AtomicSwapERC20 contract.
type AtomicSwapERC20Open struct {
	SwapID         [32]byte
	WithdrawTrader common.Address
	SecretLock     [32]byte
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterOpen is a free log retrieval operation binding the contract event 0x6ed79a08bf5c8a7d4a330df315e4ac386627ecafbe5d2bfd6654237d967b24f3.
//
// Solidity: event Open(_swapID bytes32, _withdrawTrader address, _secretLock bytes32)
func (_AtomicSwapERC20 *AtomicSwapERC20Filterer) FilterOpen(opts *bind.FilterOpts) (*AtomicSwapERC20OpenIterator, error) {

	logs, sub, err := _AtomicSwapERC20.contract.FilterLogs(opts, "Open")
	if err != nil {
		return nil, err
	}
	return &AtomicSwapERC20OpenIterator{contract: _AtomicSwapERC20.contract, event: "Open", logs: logs, sub: sub}, nil
}

// WatchOpen is a free log subscription operation binding the contract event 0x6ed79a08bf5c8a7d4a330df315e4ac386627ecafbe5d2bfd6654237d967b24f3.
//
// Solidity: event Open(_swapID bytes32, _withdrawTrader address, _secretLock bytes32)
func (_AtomicSwapERC20 *AtomicSwapERC20Filterer) WatchOpen(opts *bind.WatchOpts, sink chan<- *AtomicSwapERC20Open) (event.Subscription, error) {

	logs, sub, err := _AtomicSwapERC20.contract.WatchLogs(opts, "Open")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtomicSwapERC20Open)
				if err := _AtomicSwapERC20.contract.UnpackLog(event, "Open", log); err != nil {
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

// ERC20ABI is the input ABI used to generate the binding from.
const ERC20ABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"ok\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"ok\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"who\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"ok\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]"

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
// Solidity: function approve(spender address, value uint256) returns(ok bool)
func (_ERC20 *ERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns(ok bool)
func (_ERC20 *ERC20Session) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Approve(&_ERC20.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns(ok bool)
func (_ERC20 *ERC20TransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Approve(&_ERC20.TransactOpts, spender, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns(ok bool)
func (_ERC20 *ERC20Transactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns(ok bool)
func (_ERC20 *ERC20Session) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Transfer(&_ERC20.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns(ok bool)
func (_ERC20 *ERC20TransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Transfer(&_ERC20.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns(ok bool)
func (_ERC20 *ERC20Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns(ok bool)
func (_ERC20 *ERC20Session) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.TransferFrom(&_ERC20.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns(ok bool)
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
// Solidity: event Approval(owner indexed address, spender indexed address, value uint256)
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
// Solidity: event Approval(owner indexed address, spender indexed address, value uint256)
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
// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
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
// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
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
