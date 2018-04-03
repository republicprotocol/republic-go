// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// ArcABI is the input ABI used to generate the binding from.
const ArcABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"audit\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_tokenAddress\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"refund\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"auditSecret\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_secretLock\",\"type\":\"bytes32\"},{\"name\":\"_tokenAddress\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"},{\"name\":\"_validity\",\"type\":\"uint256\"},{\"name\":\"_receiver\",\"type\":\"address\"}],\"name\":\"test2\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_secret\",\"type\":\"bytes\"}],\"name\":\"redeem\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"test\",\"outputs\":[{\"name\":\"_secretLock\",\"type\":\"bytes32\"},{\"name\":\"_tokenAddress\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"},{\"name\":\"_validity\",\"type\":\"uint256\"},{\"name\":\"_receiver\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_secretLock\",\"type\":\"bytes32\"},{\"name\":\"_tokenAddress\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"},{\"name\":\"_validity\",\"type\":\"uint256\"},{\"name\":\"_receiver\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"}]"

// ArcBin is the compiled bytecode used for deploying new contracts.
const ArcBin = `0x6060604052341561000f57600080fd5b60405160a08061099b8339810160405280805191906020018051919060200180519190602001805191906020018051915061005f9050600086868686338764010000000061074461006982021704565b505050505061012f565b600160a060020a03851615806100885750600160a060020a0385166001145b156100a757600387018054600160a060020a03191660011790556100c5565b600387018054600160a060020a031916600160a060020a0387161790555b8654600160a060020a03338116600160a060020a03199283161789556006890197909755600180890180549489169483169490941790935560028801805492909716911617909455600485019190915542016007840155506008909101805460ff19169091179055565b61085d8061013e6000396000f30060606040526004361061005e5763ffffffff60e060020a6000350416631ddc0ef08114610060578063410085df146100a95780635afe6e75146100cb57806379e9e8cc146101555780639945e3d3146101c1578063f8a8fd6d14610212575b005b341561006b57600080fd5b610073610225565b604051600160a060020a039485168152602081019390935292166040808301919091526060820192909252608001905180910390f35b34156100b457600080fd5b61005e600160a060020a0360043516602435610243565b34156100d657600080fd5b6100de610253565b60405160208082528190810183818151815260200191508051906020019080838360005b8381101561011a578082015183820152602001610102565b50505050905090810190601f1680156101475780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561016057600080fd5b610185600435600160a060020a0360243581169060443590606435906084351661026b565b604051948552600160a060020a0393841660208601526040808601939093526060850191909152909116608083015260a0909101905180910390f35b34156101cc57600080fd5b61005e60046024813581810190830135806020601f8201819004810201604051908101604052818152929190602084018383808284375094965061027695505050505050565b341561021d57600080fd5b610185610284565b60008060008061023560006102a7565b935093509350935090919293565b61024f600083836102fa565b5050565b61025b61069a565b610265600061035f565b90505b90565b939492939192909190565b61028160008261040b565b50565b600654600354600454600754600254600160a060020a0393841693169091929394565b60008080806001600886015460ff1660038111156102c157fe5b146102cb57600080fd5b5050506003820154600483015460028401546007850154600160a060020a039384169450919216909193509193565b6001600884015460ff16600381111561030f57fe5b1461031957600080fd5b600783015442101561032a57600080fd5b60088301805460ff191660031790556001830154835461035a9184918491600160a060020a03908116911661050e565b505050565b61036761069a565b816005018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156103ff5780601f106103d4576101008083540402835291602001916103ff565b820191906000526020600020905b8154815290600101906020018083116103e257829003601f168201915b50505050509050919050565b6001600883015460ff16600381111561042057fe5b1461042a57600080fd5b6002816000604051602001526040518082805190602001908083835b602083106104655780518252601f199092019160209182019101610446565b6001836020036101000a03801982511681845116808217855250505050505090500191505060206040518083038160008661646e5a03f115156104a757600080fd5b5050604051805160068401541490506104bf57600080fd5b600582018180516104d49291602001906106ac565b5060088201805460ff191660029081179091556003830154600484015491840154845461024f93600160a060020a03938416939092811691165b600061051b8585846105fa565b151561052657600080fd5b600160a060020a0385166001141561056e57600160a060020a03831684156108fc0285604051600060405180830381858888f19350505050151561056957600080fd5b6105f3565b5083600160a060020a03811663a9059cbb848660006040516020015260405160e060020a63ffffffff8516028152600160a060020a0390921660048301526024820152604401602060405180830381600087803b15156105cd57600080fd5b6102c65a03f115156105de57600080fd5b5050506040518051905015156105f357600080fd5b5050505050565b600080600160a060020a038516600114156106185760019150610692565b508383600160a060020a0382166370a082318560006040516020015260405160e060020a63ffffffff8416028152600160a060020a039091166004820152602401602060405180830381600087803b151561067257600080fd5b6102c65a03f1151561068357600080fd5b50505060405180519050101591505b509392505050565b60206040519081016040526000815290565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106106ed57805160ff191683800117855561071a565b8280016001018555821561071a579182015b8281111561071a5782518255916020019190600101906106ff565b5061072692915061072a565b5090565b61026891905b808211156107265760008155600101610730565b600160a060020a03851615806107635750600160a060020a0385166001145b1561078f5760038701805473ffffffffffffffffffffffffffffffffffffffff191660011790556107ba565b60038701805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0387161790555b8654600160a060020a0333811673ffffffffffffffffffffffffffffffffffffffff199283161789556006890197909755600180890180549489169483169490941790935560028801805492909716911617909455600485019190915542016007840155506008909101805460ff191690911790555600a165627a7a723058209300dfee41f44361813d54aded19697a9e08cfc5370190eeb4dd1fc4ddaf9e040029`

// DeployArc deploys a new Ethereum contract, binding an instance of Arc to it.
func DeployArc(auth *bind.TransactOpts, backend bind.ContractBackend, _secretLock [32]byte, _tokenAddress common.Address, _value *big.Int, _validity *big.Int, _receiver common.Address) (common.Address, *types.Transaction, *Arc, error) {
	parsed, err := abi.JSON(strings.NewReader(ArcABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ArcBin), backend, _secretLock, _tokenAddress, _value, _validity, _receiver)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Arc{ArcCaller: ArcCaller{contract: contract}, ArcTransactor: ArcTransactor{contract: contract}, ArcFilterer: ArcFilterer{contract: contract}}, nil
}

// Arc is an auto generated Go binding around an Ethereum contract.
type Arc struct {
	ArcCaller     // Read-only binding to the contract
	ArcTransactor // Write-only binding to the contract
	ArcFilterer   // Log filterer for contract events
}

// ArcCaller is an auto generated read-only Go binding around an Ethereum contract.
type ArcCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ArcTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ArcTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ArcFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ArcFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ArcSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ArcSession struct {
	Contract     *Arc              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ArcCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ArcCallerSession struct {
	Contract *ArcCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ArcTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ArcTransactorSession struct {
	Contract     *ArcTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ArcRaw is an auto generated low-level Go binding around an Ethereum contract.
type ArcRaw struct {
	Contract *Arc // Generic contract binding to access the raw methods on
}

// ArcCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ArcCallerRaw struct {
	Contract *ArcCaller // Generic read-only contract binding to access the raw methods on
}

// ArcTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ArcTransactorRaw struct {
	Contract *ArcTransactor // Generic write-only contract binding to access the raw methods on
}

// NewArc creates a new instance of Arc, bound to a specific deployed contract.
func NewArc(address common.Address, backend bind.ContractBackend) (*Arc, error) {
	contract, err := bindArc(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Arc{ArcCaller: ArcCaller{contract: contract}, ArcTransactor: ArcTransactor{contract: contract}, ArcFilterer: ArcFilterer{contract: contract}}, nil
}

// NewArcCaller creates a new read-only instance of Arc, bound to a specific deployed contract.
func NewArcCaller(address common.Address, caller bind.ContractCaller) (*ArcCaller, error) {
	contract, err := bindArc(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ArcCaller{contract: contract}, nil
}

// NewArcTransactor creates a new write-only instance of Arc, bound to a specific deployed contract.
func NewArcTransactor(address common.Address, transactor bind.ContractTransactor) (*ArcTransactor, error) {
	contract, err := bindArc(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ArcTransactor{contract: contract}, nil
}

// NewArcFilterer creates a new log filterer instance of Arc, bound to a specific deployed contract.
func NewArcFilterer(address common.Address, filterer bind.ContractFilterer) (*ArcFilterer, error) {
	contract, err := bindArc(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ArcFilterer{contract: contract}, nil
}

// bindArc binds a generic wrapper to an already deployed contract.
func bindArc(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ArcABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Arc *ArcRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Arc.Contract.ArcCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Arc *ArcRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Arc.Contract.ArcTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Arc *ArcRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Arc.Contract.ArcTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Arc *ArcCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Arc.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Arc *ArcTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Arc.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Arc *ArcTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Arc.Contract.contract.Transact(opts, method, params...)
}

// Audit is a free data retrieval call binding the contract method 0x1ddc0ef0.
//
// Solidity: function audit() constant returns(address, uint256, address, uint256)
func (_Arc *ArcCaller) Audit(opts *bind.CallOpts) (common.Address, *big.Int, common.Address, *big.Int, error) {
	var (
		ret0 = new(common.Address)
		ret1 = new(*big.Int)
		ret2 = new(common.Address)
		ret3 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
		ret3,
	}
	err := _Arc.contract.Call(opts, out, "audit")
	return *ret0, *ret1, *ret2, *ret3, err
}

// Audit is a free data retrieval call binding the contract method 0x1ddc0ef0.
//
// Solidity: function audit() constant returns(address, uint256, address, uint256)
func (_Arc *ArcSession) Audit() (common.Address, *big.Int, common.Address, *big.Int, error) {
	return _Arc.Contract.Audit(&_Arc.CallOpts)
}

// Audit is a free data retrieval call binding the contract method 0x1ddc0ef0.
//
// Solidity: function audit() constant returns(address, uint256, address, uint256)
func (_Arc *ArcCallerSession) Audit() (common.Address, *big.Int, common.Address, *big.Int, error) {
	return _Arc.Contract.Audit(&_Arc.CallOpts)
}

// AuditSecret is a free data retrieval call binding the contract method 0x5afe6e75.
//
// Solidity: function auditSecret() constant returns(bytes)
func (_Arc *ArcCaller) AuditSecret(opts *bind.CallOpts) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _Arc.contract.Call(opts, out, "auditSecret")
	return *ret0, err
}

// AuditSecret is a free data retrieval call binding the contract method 0x5afe6e75.
//
// Solidity: function auditSecret() constant returns(bytes)
func (_Arc *ArcSession) AuditSecret() ([]byte, error) {
	return _Arc.Contract.AuditSecret(&_Arc.CallOpts)
}

// AuditSecret is a free data retrieval call binding the contract method 0x5afe6e75.
//
// Solidity: function auditSecret() constant returns(bytes)
func (_Arc *ArcCallerSession) AuditSecret() ([]byte, error) {
	return _Arc.Contract.AuditSecret(&_Arc.CallOpts)
}

// Test is a free data retrieval call binding the contract method 0xf8a8fd6d.
//
// Solidity: function test() constant returns(_secretLock bytes32, _tokenAddress address, _value uint256, _validity uint256, _receiver address)
func (_Arc *ArcCaller) Test(opts *bind.CallOpts) (struct {
	SecretLock   [32]byte
	TokenAddress common.Address
	Value        *big.Int
	Validity     *big.Int
	Receiver     common.Address
}, error) {
	ret := new(struct {
		SecretLock   [32]byte
		TokenAddress common.Address
		Value        *big.Int
		Validity     *big.Int
		Receiver     common.Address
	})
	out := ret
	err := _Arc.contract.Call(opts, out, "test")
	return *ret, err
}

// Test is a free data retrieval call binding the contract method 0xf8a8fd6d.
//
// Solidity: function test() constant returns(_secretLock bytes32, _tokenAddress address, _value uint256, _validity uint256, _receiver address)
func (_Arc *ArcSession) Test() (struct {
	SecretLock   [32]byte
	TokenAddress common.Address
	Value        *big.Int
	Validity     *big.Int
	Receiver     common.Address
}, error) {
	return _Arc.Contract.Test(&_Arc.CallOpts)
}

// Test is a free data retrieval call binding the contract method 0xf8a8fd6d.
//
// Solidity: function test() constant returns(_secretLock bytes32, _tokenAddress address, _value uint256, _validity uint256, _receiver address)
func (_Arc *ArcCallerSession) Test() (struct {
	SecretLock   [32]byte
	TokenAddress common.Address
	Value        *big.Int
	Validity     *big.Int
	Receiver     common.Address
}, error) {
	return _Arc.Contract.Test(&_Arc.CallOpts)
}

// Test2 is a free data retrieval call binding the contract method 0x79e9e8cc.
//
// Solidity: function test2(_secretLock bytes32, _tokenAddress address, _value uint256, _validity uint256, _receiver address) constant returns(bytes32, address, uint256, uint256, address)
func (_Arc *ArcCaller) Test2(opts *bind.CallOpts, _secretLock [32]byte, _tokenAddress common.Address, _value *big.Int, _validity *big.Int, _receiver common.Address) ([32]byte, common.Address, *big.Int, *big.Int, common.Address, error) {
	var (
		ret0 = new([32]byte)
		ret1 = new(common.Address)
		ret2 = new(*big.Int)
		ret3 = new(*big.Int)
		ret4 = new(common.Address)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
		ret3,
		ret4,
	}
	err := _Arc.contract.Call(opts, out, "test2", _secretLock, _tokenAddress, _value, _validity, _receiver)
	return *ret0, *ret1, *ret2, *ret3, *ret4, err
}

// Test2 is a free data retrieval call binding the contract method 0x79e9e8cc.
//
// Solidity: function test2(_secretLock bytes32, _tokenAddress address, _value uint256, _validity uint256, _receiver address) constant returns(bytes32, address, uint256, uint256, address)
func (_Arc *ArcSession) Test2(_secretLock [32]byte, _tokenAddress common.Address, _value *big.Int, _validity *big.Int, _receiver common.Address) ([32]byte, common.Address, *big.Int, *big.Int, common.Address, error) {
	return _Arc.Contract.Test2(&_Arc.CallOpts, _secretLock, _tokenAddress, _value, _validity, _receiver)
}

// Test2 is a free data retrieval call binding the contract method 0x79e9e8cc.
//
// Solidity: function test2(_secretLock bytes32, _tokenAddress address, _value uint256, _validity uint256, _receiver address) constant returns(bytes32, address, uint256, uint256, address)
func (_Arc *ArcCallerSession) Test2(_secretLock [32]byte, _tokenAddress common.Address, _value *big.Int, _validity *big.Int, _receiver common.Address) ([32]byte, common.Address, *big.Int, *big.Int, common.Address, error) {
	return _Arc.Contract.Test2(&_Arc.CallOpts, _secretLock, _tokenAddress, _value, _validity, _receiver)
}

// Redeem is a paid mutator transaction binding the contract method 0x9945e3d3.
//
// Solidity: function redeem(_secret bytes) returns()
func (_Arc *ArcTransactor) Redeem(opts *bind.TransactOpts, _secret []byte) (*types.Transaction, error) {
	return _Arc.contract.Transact(opts, "redeem", _secret)
}

// Redeem is a paid mutator transaction binding the contract method 0x9945e3d3.
//
// Solidity: function redeem(_secret bytes) returns()
func (_Arc *ArcSession) Redeem(_secret []byte) (*types.Transaction, error) {
	return _Arc.Contract.Redeem(&_Arc.TransactOpts, _secret)
}

// Redeem is a paid mutator transaction binding the contract method 0x9945e3d3.
//
// Solidity: function redeem(_secret bytes) returns()
func (_Arc *ArcTransactorSession) Redeem(_secret []byte) (*types.Transaction, error) {
	return _Arc.Contract.Redeem(&_Arc.TransactOpts, _secret)
}

// Refund is a paid mutator transaction binding the contract method 0x410085df.
//
// Solidity: function refund(_tokenAddress address, _value uint256) returns()
func (_Arc *ArcTransactor) Refund(opts *bind.TransactOpts, _tokenAddress common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Arc.contract.Transact(opts, "refund", _tokenAddress, _value)
}

// Refund is a paid mutator transaction binding the contract method 0x410085df.
//
// Solidity: function refund(_tokenAddress address, _value uint256) returns()
func (_Arc *ArcSession) Refund(_tokenAddress common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Arc.Contract.Refund(&_Arc.TransactOpts, _tokenAddress, _value)
}

// Refund is a paid mutator transaction binding the contract method 0x410085df.
//
// Solidity: function refund(_tokenAddress address, _value uint256) returns()
func (_Arc *ArcTransactorSession) Refund(_tokenAddress common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Arc.Contract.Refund(&_Arc.TransactOpts, _tokenAddress, _value)
}

// Arc0ABI is the input ABI used to generate the binding from.
const Arc0ABI = "[]"

// Arc0Bin is the compiled bytecode used for deploying new contracts.
const Arc0Bin = `0x604c602c600b82828239805160001a60731460008114601c57601e565bfe5b5030600052607381538281f30073000000000000000000000000000000000000000030146060604052600080fd00a165627a7a723058205965b576db2a9cafe6e129b68746e5058ec8b540acc0d2a3a10e4909869293a90029`

// DeployArc0 deploys a new Ethereum contract, binding an instance of Arc0 to it.
func DeployArc0(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Arc0, error) {
	parsed, err := abi.JSON(strings.NewReader(Arc0ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(Arc0Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Arc0{Arc0Caller: Arc0Caller{contract: contract}, Arc0Transactor: Arc0Transactor{contract: contract}, Arc0Filterer: Arc0Filterer{contract: contract}}, nil
}

// Arc0 is an auto generated Go binding around an Ethereum contract.
type Arc0 struct {
	Arc0Caller     // Read-only binding to the contract
	Arc0Transactor // Write-only binding to the contract
	Arc0Filterer   // Log filterer for contract events
}

// Arc0Caller is an auto generated read-only Go binding around an Ethereum contract.
type Arc0Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Arc0Transactor is an auto generated write-only Go binding around an Ethereum contract.
type Arc0Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Arc0Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Arc0Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Arc0Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Arc0Session struct {
	Contract     *Arc0             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Arc0CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Arc0CallerSession struct {
	Contract *Arc0Caller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// Arc0TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Arc0TransactorSession struct {
	Contract     *Arc0Transactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Arc0Raw is an auto generated low-level Go binding around an Ethereum contract.
type Arc0Raw struct {
	Contract *Arc0 // Generic contract binding to access the raw methods on
}

// Arc0CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Arc0CallerRaw struct {
	Contract *Arc0Caller // Generic read-only contract binding to access the raw methods on
}

// Arc0TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Arc0TransactorRaw struct {
	Contract *Arc0Transactor // Generic write-only contract binding to access the raw methods on
}

// NewArc0 creates a new instance of Arc0, bound to a specific deployed contract.
func NewArc0(address common.Address, backend bind.ContractBackend) (*Arc0, error) {
	contract, err := bindArc0(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Arc0{Arc0Caller: Arc0Caller{contract: contract}, Arc0Transactor: Arc0Transactor{contract: contract}, Arc0Filterer: Arc0Filterer{contract: contract}}, nil
}

// NewArc0Caller creates a new read-only instance of Arc0, bound to a specific deployed contract.
func NewArc0Caller(address common.Address, caller bind.ContractCaller) (*Arc0Caller, error) {
	contract, err := bindArc0(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Arc0Caller{contract: contract}, nil
}

// NewArc0Transactor creates a new write-only instance of Arc0, bound to a specific deployed contract.
func NewArc0Transactor(address common.Address, transactor bind.ContractTransactor) (*Arc0Transactor, error) {
	contract, err := bindArc0(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Arc0Transactor{contract: contract}, nil
}

// NewArc0Filterer creates a new log filterer instance of Arc0, bound to a specific deployed contract.
func NewArc0Filterer(address common.Address, filterer bind.ContractFilterer) (*Arc0Filterer, error) {
	contract, err := bindArc0(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Arc0Filterer{contract: contract}, nil
}

// bindArc0 binds a generic wrapper to an already deployed contract.
func bindArc0(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Arc0ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Arc0 *Arc0Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Arc0.Contract.Arc0Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Arc0 *Arc0Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Arc0.Contract.Arc0Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Arc0 *Arc0Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Arc0.Contract.Arc0Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Arc0 *Arc0CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Arc0.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Arc0 *Arc0TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Arc0.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Arc0 *Arc0TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Arc0.Contract.contract.Transact(opts, method, params...)
}

// TokenABI is the input ABI used to generate the binding from.
const TokenABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// TokenBin is the compiled bytecode used for deploying new contracts.
const TokenBin = `0x`

// DeployToken deploys a new Ethereum contract, binding an instance of Token to it.
func DeployToken(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Token, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(TokenBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Token{TokenCaller: TokenCaller{contract: contract}, TokenTransactor: TokenTransactor{contract: contract}, TokenFilterer: TokenFilterer{contract: contract}}, nil
}

// Token is an auto generated Go binding around an Ethereum contract.
type Token struct {
	TokenCaller     // Read-only binding to the contract
	TokenTransactor // Write-only binding to the contract
	TokenFilterer   // Log filterer for contract events
}

// TokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenSession struct {
	Contract     *Token            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenCallerSession struct {
	Contract *TokenCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// TokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenTransactorSession struct {
	Contract     *TokenTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type TokenRaw struct {
	Contract *Token // Generic contract binding to access the raw methods on
}

// TokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenCallerRaw struct {
	Contract *TokenCaller // Generic read-only contract binding to access the raw methods on
}

// TokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenTransactorRaw struct {
	Contract *TokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewToken creates a new instance of Token, bound to a specific deployed contract.
func NewToken(address common.Address, backend bind.ContractBackend) (*Token, error) {
	contract, err := bindToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Token{TokenCaller: TokenCaller{contract: contract}, TokenTransactor: TokenTransactor{contract: contract}, TokenFilterer: TokenFilterer{contract: contract}}, nil
}

// NewTokenCaller creates a new read-only instance of Token, bound to a specific deployed contract.
func NewTokenCaller(address common.Address, caller bind.ContractCaller) (*TokenCaller, error) {
	contract, err := bindToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenCaller{contract: contract}, nil
}

// NewTokenTransactor creates a new write-only instance of Token, bound to a specific deployed contract.
func NewTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenTransactor, error) {
	contract, err := bindToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenTransactor{contract: contract}, nil
}

// NewTokenFilterer creates a new log filterer instance of Token, bound to a specific deployed contract.
func NewTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*TokenFilterer, error) {
	contract, err := bindToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenFilterer{contract: contract}, nil
}

// bindToken binds a generic wrapper to an already deployed contract.
func bindToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Token *TokenRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Token.Contract.TokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Token *TokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Token.Contract.TokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Token *TokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Token.Contract.TokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Token *TokenCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Token.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Token *TokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Token.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Token *TokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Token.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance( address,  address) constant returns(uint256)
func (_Token *TokenCaller) Allowance(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Token.contract.Call(opts, out, "allowance", arg0, arg1)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance( address,  address) constant returns(uint256)
func (_Token *TokenSession) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _Token.Contract.Allowance(&_Token.CallOpts, arg0, arg1)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance( address,  address) constant returns(uint256)
func (_Token *TokenCallerSession) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _Token.Contract.Allowance(&_Token.CallOpts, arg0, arg1)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf( address) constant returns(uint256)
func (_Token *TokenCaller) BalanceOf(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Token.contract.Call(opts, out, "balanceOf", arg0)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf( address) constant returns(uint256)
func (_Token *TokenSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _Token.Contract.BalanceOf(&_Token.CallOpts, arg0)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf( address) constant returns(uint256)
func (_Token *TokenCallerSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _Token.Contract.BalanceOf(&_Token.CallOpts, arg0)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer( address,  uint256) returns(bool)
func (_Token *TokenTransactor) Transfer(opts *bind.TransactOpts, arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "transfer", arg0, arg1)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer( address,  uint256) returns(bool)
func (_Token *TokenSession) Transfer(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _Token.Contract.Transfer(&_Token.TransactOpts, arg0, arg1)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer( address,  uint256) returns(bool)
func (_Token *TokenTransactorSession) Transfer(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _Token.Contract.Transfer(&_Token.TransactOpts, arg0, arg1)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom( address,  address,  uint256) returns(bool)
func (_Token *TokenTransactor) TransferFrom(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, arg2 *big.Int) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "transferFrom", arg0, arg1, arg2)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom( address,  address,  uint256) returns(bool)
func (_Token *TokenSession) TransferFrom(arg0 common.Address, arg1 common.Address, arg2 *big.Int) (*types.Transaction, error) {
	return _Token.Contract.TransferFrom(&_Token.TransactOpts, arg0, arg1, arg2)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom( address,  address,  uint256) returns(bool)
func (_Token *TokenTransactorSession) TransferFrom(arg0 common.Address, arg1 common.Address, arg2 *big.Int) (*types.Transaction, error) {
	return _Token.Contract.TransferFrom(&_Token.TransactOpts, arg0, arg1, arg2)
}
