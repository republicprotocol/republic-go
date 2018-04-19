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
const ArcABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"audit\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_tokenAddress\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"refund\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"auditSecret\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_secret\",\"type\":\"bytes\"}],\"name\":\"redeem\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_secretLock\",\"type\":\"bytes32\"},{\"name\":\"_tokenAddress\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"},{\"name\":\"_validity\",\"type\":\"uint256\"},{\"name\":\"_receiver\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"}]"

// ArcBin is the compiled bytecode used for deploying new contracts.
const ArcBin = `0x608060405234801561001057600080fd5b5060405160a080610936833981016040908152815160208301519183015160608401516080909401519193909161005a60008686868633876401000000006106e461006482021704565b505050505061012a565b600160a060020a03851615806100835750600160a060020a0385166001145b156100a257600387018054600160a060020a03191660011790556100c0565b600387018054600160a060020a031916600160a060020a0387161790555b8654600160a060020a03338116600160a060020a03199283161789556006890197909755600180890180549489169483169490941790935560028801805492909716911617909455600485019190915542016007840155506008909101805460ff19169091179055565b6107fd806101396000396000f3006080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416631ddc0ef08114610063578063410085df146100af5780635afe6e75146100d35780639945e3d31461015d575b005b34801561006f57600080fd5b506100786101b6565b60408051958652600160a060020a039485166020870152929093168483015260608401526080830191909152519081900360a00190f35b3480156100bb57600080fd5b50610061600160a060020a03600435166024356101d9565b3480156100df57600080fd5b506100e86101e9565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561012257818101518382015260200161010a565b50505050905090810190601f16801561014f5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561016957600080fd5b506040805160206004803580820135601f81018490048402850184019095528484526100619436949293602493928401919081908401838280828437509497506101fb9650505050505050565b60008060008060006101c86000610209565b945094509450945094509091929394565b6101e560008383610264565b5050565b60606101f560006102c9565b90505b90565b610206600082610363565b50565b6000808080806001600887015460ff16600381111561022457fe5b1461022e57600080fd5b50505060068301546003840154600285015460048601546007870154939550600160a060020a0392831694509116939590929450565b6001600884015460ff16600381111561027957fe5b1461028357600080fd5b600783015442101561029457600080fd5b60088301805460ff19166003179055600183015483546102c49184918491600160a060020a03908116911661046a565b505050565b600581018054604080516020601f60026000196101006001881615020190951694909404938401819004810282018101909252828152606093909290918301828280156103575780601f1061032c57610100808354040283529160200191610357565b820191906000526020600020905b81548152906001019060200180831161033a57829003601f168201915b50505050509050919050565b6001600883015460ff16600381111561037857fe5b1461038257600080fd5b6002816040518082805190602001908083835b602083106103b45780518252601f199092019160209182019101610395565b51815160209384036101000a600019018019909216911617905260405191909301945091925050808303816000865af11580156103f5573d6000803e3d6000fd5b5050506040513d602081101561040a57600080fd5b505160068301541461041b57600080fd5b8051610430906005840190602084019061064c565b5060088201805460ff19166002908117909155600383015460048401549184015484546101e593600160a060020a03938416939092811691165b6000610477858584610571565b151561048257600080fd5b600160a060020a038516600114156104d057604051600160a060020a0384169085156108fc029086906000818181858888f193505050501580156104ca573d6000803e3d6000fd5b5061056a565b50604080517fa9059cbb000000000000000000000000000000000000000000000000000000008152600160a060020a038481166004830152602482018690529151869283169163a9059cbb9160448083019260209291908290030181600087803b15801561053d57600080fd5b505af1158015610551573d6000803e3d6000fd5b505050506040513d602081101561056757600080fd5b50505b5050505050565b600080600160a060020a0385166001141561059b578383600160a060020a03163110159150610644565b8490508381600160a060020a03166370a08231856040518263ffffffff167c01000000000000000000000000000000000000000000000000000000000281526004018082600160a060020a0316600160a060020a03168152602001915050602060405180830381600087803b15801561061357600080fd5b505af1158015610627573d6000803e3d6000fd5b505050506040513d602081101561063d57600080fd5b5051101591505b509392505050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061068d57805160ff19168380011785556106ba565b828001600101855582156106ba579182015b828111156106ba57825182559160200191906001019061069f565b506106c69291506106ca565b5090565b6101f891905b808211156106c657600081556001016106d0565b600160a060020a03851615806107035750600160a060020a0385166001145b1561072f5760038701805473ffffffffffffffffffffffffffffffffffffffff1916600117905561075a565b60038701805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0387161790555b8654600160a060020a0333811673ffffffffffffffffffffffffffffffffffffffff199283161789556006890197909755600180890180549489169483169490941790935560028801805492909716911617909455600485019190915542016007840155506008909101805460ff191690911790555600a165627a7a72305820a90e5579f68a00c6879aad790cc12754ee7109990ad114acf23112e2447ebca40029`

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
// Solidity: function audit() constant returns(bytes32, address, address, uint256, uint256)
func (_Arc *ArcCaller) Audit(opts *bind.CallOpts) ([32]byte, common.Address, common.Address, *big.Int, *big.Int, error) {
	var (
		ret0 = new([32]byte)
		ret1 = new(common.Address)
		ret2 = new(common.Address)
		ret3 = new(*big.Int)
		ret4 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
		ret3,
		ret4,
	}
	err := _Arc.contract.Call(opts, out, "audit")
	return *ret0, *ret1, *ret2, *ret3, *ret4, err
}

// Audit is a free data retrieval call binding the contract method 0x1ddc0ef0.
//
// Solidity: function audit() constant returns(bytes32, address, address, uint256, uint256)
func (_Arc *ArcSession) Audit() ([32]byte, common.Address, common.Address, *big.Int, *big.Int, error) {
	return _Arc.Contract.Audit(&_Arc.CallOpts)
}

// Audit is a free data retrieval call binding the contract method 0x1ddc0ef0.
//
// Solidity: function audit() constant returns(bytes32, address, address, uint256, uint256)
func (_Arc *ArcCallerSession) Audit() ([32]byte, common.Address, common.Address, *big.Int, *big.Int, error) {
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

// LibArcABI is the input ABI used to generate the binding from.
const LibArcABI = "[]"

// LibArcBin is the compiled bytecode used for deploying new contracts.
const LibArcBin = `0x604c602c600b82828239805160001a60731460008114601c57601e565bfe5b5030600052607381538281f30073000000000000000000000000000000000000000030146080604052600080fd00a165627a7a72305820c9ef7f8792694783ae3c471a949964713dbe61967b8d1218c473ca83b8fd996b0029`

// DeployLibArc deploys a new Ethereum contract, binding an instance of LibArc to it.
func DeployLibArc(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *LibArc, error) {
	parsed, err := abi.JSON(strings.NewReader(LibArcABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(LibArcBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LibArc{LibArcCaller: LibArcCaller{contract: contract}, LibArcTransactor: LibArcTransactor{contract: contract}, LibArcFilterer: LibArcFilterer{contract: contract}}, nil
}

// LibArc is an auto generated Go binding around an Ethereum contract.
type LibArc struct {
	LibArcCaller     // Read-only binding to the contract
	LibArcTransactor // Write-only binding to the contract
	LibArcFilterer   // Log filterer for contract events
}

// LibArcCaller is an auto generated read-only Go binding around an Ethereum contract.
type LibArcCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LibArcTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LibArcTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LibArcFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LibArcFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LibArcSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LibArcSession struct {
	Contract     *LibArc           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LibArcCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LibArcCallerSession struct {
	Contract *LibArcCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// LibArcTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LibArcTransactorSession struct {
	Contract     *LibArcTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LibArcRaw is an auto generated low-level Go binding around an Ethereum contract.
type LibArcRaw struct {
	Contract *LibArc // Generic contract binding to access the raw methods on
}

// LibArcCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LibArcCallerRaw struct {
	Contract *LibArcCaller // Generic read-only contract binding to access the raw methods on
}

// LibArcTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LibArcTransactorRaw struct {
	Contract *LibArcTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLibArc creates a new instance of LibArc, bound to a specific deployed contract.
func NewLibArc(address common.Address, backend bind.ContractBackend) (*LibArc, error) {
	contract, err := bindLibArc(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LibArc{LibArcCaller: LibArcCaller{contract: contract}, LibArcTransactor: LibArcTransactor{contract: contract}, LibArcFilterer: LibArcFilterer{contract: contract}}, nil
}

// NewLibArcCaller creates a new read-only instance of LibArc, bound to a specific deployed contract.
func NewLibArcCaller(address common.Address, caller bind.ContractCaller) (*LibArcCaller, error) {
	contract, err := bindLibArc(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LibArcCaller{contract: contract}, nil
}

// NewLibArcTransactor creates a new write-only instance of LibArc, bound to a specific deployed contract.
func NewLibArcTransactor(address common.Address, transactor bind.ContractTransactor) (*LibArcTransactor, error) {
	contract, err := bindLibArc(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LibArcTransactor{contract: contract}, nil
}

// NewLibArcFilterer creates a new log filterer instance of LibArc, bound to a specific deployed contract.
func NewLibArcFilterer(address common.Address, filterer bind.ContractFilterer) (*LibArcFilterer, error) {
	contract, err := bindLibArc(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LibArcFilterer{contract: contract}, nil
}

// bindLibArc binds a generic wrapper to an already deployed contract.
func bindLibArc(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LibArcABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LibArc *LibArcRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LibArc.Contract.LibArcCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LibArc *LibArcRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LibArc.Contract.LibArcTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LibArc *LibArcRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LibArc.Contract.LibArcTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LibArc *LibArcCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LibArc.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LibArc *LibArcTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LibArc.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LibArc *LibArcTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LibArc.Contract.contract.Transact(opts, method, params...)
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
