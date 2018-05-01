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

// ArcABI is the input ABI used to generate the binding from.
const ArcABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_secretLock\",\"type\":\"bytes32\"},{\"name\":\"_tokenAddress\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"},{\"name\":\"_feeRate\",\"type\":\"uint256\"},{\"name\":\"_validity\",\"type\":\"uint256\"},{\"name\":\"_receiver\",\"type\":\"address\"},{\"name\":\"_order\",\"type\":\"bytes\"}],\"name\":\"initiate\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_orderID\",\"type\":\"bytes32\"},{\"name\":\"_tokenAddress\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"refund\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_orderID\",\"type\":\"bytes32\"}],\"name\":\"auditSecret\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_orderID\",\"type\":\"bytes32\"},{\"name\":\"_secret\",\"type\":\"bytes32\"}],\"name\":\"redeem\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_orderID\",\"type\":\"bytes32\"}],\"name\":\"audit\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_rewardGatewayAddress\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// ArcBin is the compiled bytecode used for deploying new contracts.
const ArcBin = `0x608060405234801561001057600080fd5b50604051602080610d61833981016040525160008054600160a060020a03909216600160a060020a0319909216919091179055610d0f806100526000396000f30060806040526004361061006c5763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416631c7be96881146100715780634020ed14146100e6578063976d00f41461010d578063b31597ad14610137578063c140635b14610152575b600080fd5b604080516020601f60c4356004818101359283018490048402850184019095528184526100e494803594600160a060020a036024803582169660443596606435966084359660a4359095169536959460e49491939201919081908401838280828437509497506101a19650505050505050565b005b3480156100f257600080fd5b506100e4600435600160a060020a03602435166044356103fa565b34801561011957600080fd5b5061012560043561049b565b60408051918252519081900360200190f35b34801561014357600080fd5b506100e46004356024356104df565b34801561015e57600080fd5b5061016a6004356105f2565b60408051958652600160a060020a039485166020870152929093168483015260608401526080830191909152519081900360a00190f35b6000816040518082805190602001908083835b602083106101d35780518252601f1990920191602091820191016101b4565b5181516020939093036101000a6000190180199091169216919091179052604051920182900390912093506000925061020a915050565b6000828152600160205260409020600c015460ff16600381111561022a57fe5b1461023457600080fd5b61023e8787610667565b151561024957600080fd5b600081815260016020908152604090912060038101805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a038b16179055835161029892600690920191850190610c48565b506000818152600160205260409020600a018890556102b786866107f0565b600082815260016020818152604080842060058101869055948b036004808701919091558584018054600160a060020a0333811673ffffffffffffffffffffffffffffffffffffffff19928316179092556002880180548c84169216919091179055600b87018b9055600c909601805460ff1916909417909355835481517f29228ba70000000000000000000000000000000000000000000000000000000081528d87169481019490945290519416936329228ba7936024808501948390030190829087803b15801561038957600080fd5b505af115801561039d573d6000803e3d6000fd5b505050506040513d60208110156103b357600080fd5b5051600091825260016020526040909120600801805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0390921691909117905550505050505050565b60016000848152600160205260409020600c015460ff16600381111561041c57fe5b1461042657600080fd5b6000838152600160205260409020600b015442101561044457600080fd5b6000838152600160208190526040909120015433600160a060020a0390811691161461046f57600080fd5b6000838152600160205260409020600c01805460ff19166003179055610496828233610821565b505050565b600060026000838152600160205260409020600c015460ff1660038111156104bf57fe5b146104c957600080fd5b5060009081526001602052604090206009015490565b60016000838152600160205260409020600c015460ff16600381111561050157fe5b1461050b57600080fd5b60008281526001602052604090206002015433600160a060020a0390811691161461053557600080fd5b60408051828152905160029160208082019290918190038201816000865af1158015610565573d6000803e3d6000fd5b5050506040513d602081101561057a57600080fd5b50516000838152600160205260409020600a01541461059857600080fd5b600082815260016020526040902060098101829055600c8101805460ff191660029081179091556003820154600483015491909201546105e592600160a060020a03908116929116610821565b6105ee82610924565b5050565b60008080808060016000878152600160205260409020600c015460ff16600381111561061a57fe5b1461062457600080fd5b50505060009283525050600160205260409020600a810154600382015460028301546004840154600b909401549294600160a060020a0392831694929091169290565b600080600160a060020a03841673eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee14156106995782341491506107e9565b50604080517fdd62ed3e000000000000000000000000000000000000000000000000000000008152600160a060020a03338116600483015230811660248301529151859285929084169163dd62ed3e916044808201926020929091908290030181600087803b15801561070b57600080fd5b505af115801561071f573d6000803e3d6000fd5b505050506040513d602081101561073557600080fd5b5051101561074657600091506107e9565b604080517f23b872dd000000000000000000000000000000000000000000000000000000008152600160a060020a0333811660048301523081166024830152604482018690529151918316916323b872dd916064808201926020929091908290030181600087803b1580156107ba57600080fd5b505af11580156107ce573d6000803e3d6000fd5b505050506040513d60208110156107e457600080fd5b505191505b5092915050565b6000826103e881046103e802141561080f57506103e88282020461081b565b5060016103e883830204015b92915050565b6000600160a060020a03841673eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee141561088457604051600160a060020a0383169084156108fc029085906000818181858888f1935050505015801561087e573d6000803e3d6000fd5b5061091e565b50604080517fa9059cbb000000000000000000000000000000000000000000000000000000008152600160a060020a038381166004830152602482018590529151859283169163a9059cbb9160448083019260209291908290030181600087803b1580156108f157600080fd5b505af1158015610905573d6000803e3d6000fd5b505050506040513d602081101561091b57600080fd5b50505b50505050565b600081815260016020526040812060030154600160a060020a031673eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee14610b2c575060008181526001602090815260408083206003810154600882015460059092015483517f095ea7b3000000000000000000000000000000000000000000000000000000008152600160a060020a03938416600482015260248101919091529251911693849363095ea7b3936044808201949293918390030190829087803b1580156109e457600080fd5b505af11580156109f8573d6000803e3d6000fd5b505050506040513d6020811015610a0e57600080fd5b50506000828152600160208190526040918290206008810154600582015484517f934055400000000000000000000000000000000000000000000000000000000081526024810182905260048101958652600690930180546002958116156101000260001901169490940460448401819052600160a060020a03909216946393405540949391929091829160649091019085908015610aee5780601f10610ac357610100808354040283529160200191610aee565b820191906000526020600020905b815481529060010190602001808311610ad157829003601f168201915b50509350505050600060405180830381600087803b158015610b0f57600080fd5b505af1158015610b23573d6000803e3d6000fd5b505050506105ee565b6000828152600160208190526040918290206008810154600582015484517f934055400000000000000000000000000000000000000000000000000000000081526024810182905260048101958652600690930180546002958116156101000260001901169490940460448401819052600160a060020a0390921694639340554094919391928492829160649091019085908015610c0b5780601f10610be057610100808354040283529160200191610c0b565b820191906000526020600020905b815481529060010190602001808311610bee57829003601f168201915b505093505050506000604051808303818588803b158015610c2b57600080fd5b505af1158015610c3f573d6000803e3d6000fd5b50505050505050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10610c8957805160ff1916838001178555610cb6565b82800160010185558215610cb6579182015b82811115610cb6578251825591602001919060010190610c9b565b50610cc2929150610cc6565b5090565b610ce091905b80821115610cc25760008155600101610ccc565b905600a165627a7a72305820e32143b767cba5e14ac8ff59b97e260e966800ed1789dd21c950e7d4e829e94a0029`

// DeployArc deploys a new Ethereum contract, binding an instance of Arc to it.
func DeployArc(auth *bind.TransactOpts, backend bind.ContractBackend, _rewardGatewayAddress common.Address) (common.Address, *types.Transaction, *Arc, error) {
	parsed, err := abi.JSON(strings.NewReader(ArcABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ArcBin), backend, _rewardGatewayAddress)
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

// Audit is a free data retrieval call binding the contract method 0xc140635b.
//
// Solidity: function audit(_orderID bytes32) constant returns(bytes32, address, address, uint256, uint256)
func (_Arc *ArcCaller) Audit(opts *bind.CallOpts, _orderID [32]byte) ([32]byte, common.Address, common.Address, *big.Int, *big.Int, error) {
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
	err := _Arc.contract.Call(opts, out, "audit", _orderID)
	return *ret0, *ret1, *ret2, *ret3, *ret4, err
}

// Audit is a free data retrieval call binding the contract method 0xc140635b.
//
// Solidity: function audit(_orderID bytes32) constant returns(bytes32, address, address, uint256, uint256)
func (_Arc *ArcSession) Audit(_orderID [32]byte) ([32]byte, common.Address, common.Address, *big.Int, *big.Int, error) {
	return _Arc.Contract.Audit(&_Arc.CallOpts, _orderID)
}

// Audit is a free data retrieval call binding the contract method 0xc140635b.
//
// Solidity: function audit(_orderID bytes32) constant returns(bytes32, address, address, uint256, uint256)
func (_Arc *ArcCallerSession) Audit(_orderID [32]byte) ([32]byte, common.Address, common.Address, *big.Int, *big.Int, error) {
	return _Arc.Contract.Audit(&_Arc.CallOpts, _orderID)
}

// AuditSecret is a free data retrieval call binding the contract method 0x976d00f4.
//
// Solidity: function auditSecret(_orderID bytes32) constant returns(bytes32)
func (_Arc *ArcCaller) AuditSecret(opts *bind.CallOpts, _orderID [32]byte) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Arc.contract.Call(opts, out, "auditSecret", _orderID)
	return *ret0, err
}

// AuditSecret is a free data retrieval call binding the contract method 0x976d00f4.
//
// Solidity: function auditSecret(_orderID bytes32) constant returns(bytes32)
func (_Arc *ArcSession) AuditSecret(_orderID [32]byte) ([32]byte, error) {
	return _Arc.Contract.AuditSecret(&_Arc.CallOpts, _orderID)
}

// AuditSecret is a free data retrieval call binding the contract method 0x976d00f4.
//
// Solidity: function auditSecret(_orderID bytes32) constant returns(bytes32)
func (_Arc *ArcCallerSession) AuditSecret(_orderID [32]byte) ([32]byte, error) {
	return _Arc.Contract.AuditSecret(&_Arc.CallOpts, _orderID)
}

// Initiate is a paid mutator transaction binding the contract method 0x1c7be968.
//
// Solidity: function initiate(_secretLock bytes32, _tokenAddress address, _value uint256, _feeRate uint256, _validity uint256, _receiver address, _order bytes) returns()
func (_Arc *ArcTransactor) Initiate(opts *bind.TransactOpts, _secretLock [32]byte, _tokenAddress common.Address, _value *big.Int, _feeRate *big.Int, _validity *big.Int, _receiver common.Address, _order []byte) (*types.Transaction, error) {
	return _Arc.contract.Transact(opts, "initiate", _secretLock, _tokenAddress, _value, _feeRate, _validity, _receiver, _order)
}

// Initiate is a paid mutator transaction binding the contract method 0x1c7be968.
//
// Solidity: function initiate(_secretLock bytes32, _tokenAddress address, _value uint256, _feeRate uint256, _validity uint256, _receiver address, _order bytes) returns()
func (_Arc *ArcSession) Initiate(_secretLock [32]byte, _tokenAddress common.Address, _value *big.Int, _feeRate *big.Int, _validity *big.Int, _receiver common.Address, _order []byte) (*types.Transaction, error) {
	return _Arc.Contract.Initiate(&_Arc.TransactOpts, _secretLock, _tokenAddress, _value, _feeRate, _validity, _receiver, _order)
}

// Initiate is a paid mutator transaction binding the contract method 0x1c7be968.
//
// Solidity: function initiate(_secretLock bytes32, _tokenAddress address, _value uint256, _feeRate uint256, _validity uint256, _receiver address, _order bytes) returns()
func (_Arc *ArcTransactorSession) Initiate(_secretLock [32]byte, _tokenAddress common.Address, _value *big.Int, _feeRate *big.Int, _validity *big.Int, _receiver common.Address, _order []byte) (*types.Transaction, error) {
	return _Arc.Contract.Initiate(&_Arc.TransactOpts, _secretLock, _tokenAddress, _value, _feeRate, _validity, _receiver, _order)
}

// Redeem is a paid mutator transaction binding the contract method 0xb31597ad.
//
// Solidity: function redeem(_orderID bytes32, _secret bytes32) returns()
func (_Arc *ArcTransactor) Redeem(opts *bind.TransactOpts, _orderID [32]byte, _secret [32]byte) (*types.Transaction, error) {
	return _Arc.contract.Transact(opts, "redeem", _orderID, _secret)
}

// Redeem is a paid mutator transaction binding the contract method 0xb31597ad.
//
// Solidity: function redeem(_orderID bytes32, _secret bytes32) returns()
func (_Arc *ArcSession) Redeem(_orderID [32]byte, _secret [32]byte) (*types.Transaction, error) {
	return _Arc.Contract.Redeem(&_Arc.TransactOpts, _orderID, _secret)
}

// Redeem is a paid mutator transaction binding the contract method 0xb31597ad.
//
// Solidity: function redeem(_orderID bytes32, _secret bytes32) returns()
func (_Arc *ArcTransactorSession) Redeem(_orderID [32]byte, _secret [32]byte) (*types.Transaction, error) {
	return _Arc.Contract.Redeem(&_Arc.TransactOpts, _orderID, _secret)
}

// Refund is a paid mutator transaction binding the contract method 0x4020ed14.
//
// Solidity: function refund(_orderID bytes32, _tokenAddress address, _value uint256) returns()
func (_Arc *ArcTransactor) Refund(opts *bind.TransactOpts, _orderID [32]byte, _tokenAddress common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Arc.contract.Transact(opts, "refund", _orderID, _tokenAddress, _value)
}

// Refund is a paid mutator transaction binding the contract method 0x4020ed14.
//
// Solidity: function refund(_orderID bytes32, _tokenAddress address, _value uint256) returns()
func (_Arc *ArcSession) Refund(_orderID [32]byte, _tokenAddress common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Arc.Contract.Refund(&_Arc.TransactOpts, _orderID, _tokenAddress, _value)
}

// Refund is a paid mutator transaction binding the contract method 0x4020ed14.
//
// Solidity: function refund(_orderID bytes32, _tokenAddress address, _value uint256) returns()
func (_Arc *ArcTransactorSession) Refund(_orderID [32]byte, _tokenAddress common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Arc.Contract.Refund(&_Arc.TransactOpts, _orderID, _tokenAddress, _value)
}

// BasicTokenABI is the input ABI used to generate the binding from.
const BasicTokenABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]"

// BasicTokenBin is the compiled bytecode used for deploying new contracts.
const BasicTokenBin = `0x608060405234801561001057600080fd5b50610249806100206000396000f3006080604052600436106100565763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166318160ddd811461005b57806370a0823114610082578063a9059cbb146100a3575b600080fd5b34801561006757600080fd5b506100706100db565b60408051918252519081900360200190f35b34801561008e57600080fd5b50610070600160a060020a03600435166100e1565b3480156100af57600080fd5b506100c7600160a060020a03600435166024356100fc565b604080519115158252519081900360200190f35b60015490565b600160a060020a031660009081526020819052604090205490565b6000600160a060020a038316151561011357600080fd5b600160a060020a03331660009081526020819052604090205482111561013857600080fd5b600160a060020a033316600090815260208190526040902054610161908363ffffffff6101f516565b600160a060020a033381166000908152602081905260408082209390935590851681522054610196908363ffffffff61020716565b600160a060020a03808516600081815260208181526040918290209490945580518681529051919333909316927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef92918290030190a350600192915050565b60008282111561020157fe5b50900390565b60008282018381101561021657fe5b93925050505600a165627a7a72305820a3edee6b552fb7e383a99736cacd1afdd71336157b3079e103ca972761bf3cab0029`

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
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
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
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_BasicToken *BasicTokenSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _BasicToken.Contract.BalanceOf(&_BasicToken.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
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
// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
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
// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
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

// BurnableTokenABI is the input ABI used to generate the binding from.
const BurnableTokenABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"burner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Burn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]"

// BurnableTokenBin is the compiled bytecode used for deploying new contracts.
const BurnableTokenBin = `0x608060405234801561001057600080fd5b5061036b806100206000396000f3006080604052600436106100615763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166318160ddd811461006657806342966c681461008d57806370a08231146100a7578063a9059cbb146100c8575b600080fd5b34801561007257600080fd5b5061007b610100565b60408051918252519081900360200190f35b34801561009957600080fd5b506100a5600435610106565b005b3480156100b357600080fd5b5061007b600160a060020a0360043516610203565b3480156100d457600080fd5b506100ec600160a060020a036004351660243561021e565b604080519115158252519081900360200190f35b60015490565b600160a060020a03331660009081526020819052604081205482111561012b57600080fd5b5033600160a060020a0381166000908152602081905260409020546101509083610317565b600160a060020a03821660009081526020819052604090205560015461017c908363ffffffff61031716565b600155604080518381529051600160a060020a038316917fcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5919081900360200190a2604080518381529051600091600160a060020a038416917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9181900360200190a35050565b600160a060020a031660009081526020819052604090205490565b6000600160a060020a038316151561023557600080fd5b600160a060020a03331660009081526020819052604090205482111561025a57600080fd5b600160a060020a033316600090815260208190526040902054610283908363ffffffff61031716565b600160a060020a0333811660009081526020819052604080822093909355908516815220546102b8908363ffffffff61032916565b600160a060020a03808516600081815260208181526040918290209490945580518681529051919333909316927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef92918290030190a350600192915050565b60008282111561032357fe5b50900390565b60008282018381101561033857fe5b93925050505600a165627a7a72305820122ce7be7c34440176ba15d217d39ea33cc002a18c5863aaf4cfe76bb9a6bc7c0029`

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
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
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
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_BurnableToken *BurnableTokenSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _BurnableToken.Contract.BalanceOf(&_BurnableToken.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
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
// Solidity: event Burn(burner indexed address, value uint256)
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
// Solidity: event Burn(burner indexed address, value uint256)
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
// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
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
// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
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
const DarknodeRegistryABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"numDarkNodes\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numDarkNodesNextEpoch\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_darkNodeID\",\"type\":\"bytes20\"}],\"name\":\"isDeregistered\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_darkNodeID\",\"type\":\"bytes20\"}],\"name\":\"getPublicKey\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_darkNodeID\",\"type\":\"bytes20\"},{\"name\":\"_publicKey\",\"type\":\"bytes\"},{\"name\":\"_bond\",\"type\":\"uint256\"}],\"name\":\"register\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_darkNodeID\",\"type\":\"bytes20\"}],\"name\":\"isRegistered\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minimumEpochInterval\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_darkNodeID\",\"type\":\"bytes20\"}],\"name\":\"refund\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_darkNodeID\",\"type\":\"bytes20\"}],\"name\":\"getBond\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentEpoch\",\"outputs\":[{\"name\":\"epochhash\",\"type\":\"uint256\"},{\"name\":\"timestamp\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getDarkNodes\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes20[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"epoch\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minimumBond\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minimumDarkPoolSize\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_darkNodeID\",\"type\":\"bytes20\"}],\"name\":\"isUnregistered\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_darkNodeID\",\"type\":\"bytes20\"}],\"name\":\"deregister\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_darkNodeID\",\"type\":\"bytes20\"}],\"name\":\"getOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_token\",\"type\":\"address\"},{\"name\":\"_minimumBond\",\"type\":\"uint256\"},{\"name\":\"_minimumDarkPoolSize\",\"type\":\"uint256\"},{\"name\":\"_minimumEpochInterval\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_darkNodeID\",\"type\":\"bytes20\"},{\"indexed\":false,\"name\":\"_bond\",\"type\":\"uint256\"}],\"name\":\"DarkNodeRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_darkNodeID\",\"type\":\"bytes20\"}],\"name\":\"DarkNodeDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"OwnerRefunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"NewEpoch\",\"type\":\"event\"}]"

// DarknodeRegistryBin is the compiled bytecode used for deploying new contracts.
const DarknodeRegistryBin = `0x608060405234801561001057600080fd5b506040516080806110df83398101604081815282516020808501518386015160609096015160008054600160a060020a03909516600160a060020a03199095169490941784556005919091556006959095556007949094558183019091524360001901408083524293909201839052600891909155600991909155600381905560045561103d806100a26000396000f3006080604052600436106100f05763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663060a2cfd81146100f55780630620eb921461011c578063171f6ea81461013157806332ccd52f14610167578063375a8be3146101fe5780634f5550fc1461026a57806355cacda51461028c5780635a8f9b81146102a157806368f209eb146102c357806376671808146102e5578063879ae08414610313578063900cf0cf14610378578063aa7517e11461038d578063b31575d5146103a2578063d3841c25146103b7578063e08b4c8a146103d9578063e487eb58146103fb575b600080fd5b34801561010157600080fd5b5061010a610439565b60408051918252519081900360200190f35b34801561012857600080fd5b5061010a61043f565b34801561013d57600080fd5b506101536001606060020a031960043516610445565b604080519115158252519081900360200190f35b34801561017357600080fd5b506101896001606060020a031960043516610494565b6040805160208082528351818301528351919283929083019185019080838360005b838110156101c35781810151838201526020016101ab565b50505050905090810190601f1680156101f05780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561020a57600080fd5b5060408051602060046024803582810135601f81018590048502860185019096528585526102689583356001606060020a03191695369560449491939091019190819084018382808284375094975050933594506105429350505050565b005b34801561027657600080fd5b506101536001606060020a0319600435166107c8565b34801561029857600080fd5b5061010a610828565b3480156102ad57600080fd5b506102686001606060020a03196004351661082e565b3480156102cf57600080fd5b5061010a6001606060020a031960043516610a40565b3480156102f157600080fd5b506102fa610a60565b6040805192835260208301919091528051918290030190f35b34801561031f57600080fd5b50610328610a69565b60408051602080825283518183015283519192839290830191858101910280838360005b8381101561036457818101518382015260200161034c565b505050509050019250505060405180910390f35b34801561038457600080fd5b50610268610b21565b34801561039957600080fd5b5061010a610b95565b3480156103ae57600080fd5b5061010a610b9b565b3480156103c357600080fd5b506101536001606060020a031960043516610ba1565b3480156103e557600080fd5b506102686001606060020a031960043516610bc1565b34801561040757600080fd5b5061041d6001606060020a031960043516610c7b565b60408051600160a060020a039092168252519081900360200190f35b60035481565b60045481565b6001606060020a031981166000908152600160205260408120600301541580159061048e57506009546001606060020a0319831660009081526001602052604090206003015411155b92915050565b6001606060020a0319811660009081526001602081815260409283902060040180548451600294821615610100026000190190911693909304601f810183900483028401830190945283835260609390918301828280156105365780601f1061050b57610100808354040283529160200191610536565b820191906000526020600020905b81548152906001019060200180831161051957829003601f168201915b50505050509050919050565b8261054c81610ba1565b151561055757600080fd5b60055482101561056657600080fd5b60008054604080517fdd62ed3e000000000000000000000000000000000000000000000000000000008152600160a060020a03338116600483015230811660248301529151919092169263dd62ed3e92604480820193602093909283900390910190829087803b1580156105d957600080fd5b505af11580156105ed573d6000803e3d6000fd5b505050506040513d602081101561060357600080fd5b505182111561061157600080fd5b60008054604080517f23b872dd000000000000000000000000000000000000000000000000000000008152600160a060020a033381166004830152308116602483015260448201879052915191909216926323b872dd92606480820193602093909283900390910190829087803b15801561068b57600080fd5b505af115801561069f573d6000803e3d6000fd5b505050506040513d60208110156106b557600080fd5b505115156106c257600080fd5b6040805160a08101825233600160a060020a039081168252602080830186815260075460095401848601908152600060608601818152608087018b81526001606060020a03198d1683526001808752989092208751815473ffffffffffffffffffffffffffffffffffffffff1916971696909617865592519685019690965551600284015551600383015592518051929391926107659260048501920190610f76565b50905050610774600285610ca0565b600480546001019055604080516001606060020a0319861681526020810184905281517fcde427a7822769e806a704726602a3a8a44458da1352d9b5a46bcfb95142ad38929181900390910190a150505050565b6001606060020a031981166000908152600160205260408120600201541580159061081157506009546001606060020a0319831660009081526001602052604090206002015411155b801561048e575061082182610445565b1592915050565b60075481565b6001606060020a03198116600090815260016020526040812054829033600160a060020a0390811691161461086257600080fd5b8261086c81610445565b151561087757600080fd5b6001606060020a031984166000908152600160208190526040822001549350831161089e57fe5b6108a9600285610ccb565b6040805160a081018252600080825260208083018281528385018381526060850184815286518085018852858152608087019081526001606060020a03198c1686526001808652979095208651815473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03909116178155925196830196909655516002820155935160038501559051805192939261094a9260048501920190610f76565b505060008054604080517fa9059cbb000000000000000000000000000000000000000000000000000000008152600160a060020a03338116600483015260248201899052915191909216935063a9059cbb92604480840193602093929083900390910190829087803b1580156109bf57600080fd5b505af11580156109d3573d6000803e3d6000fd5b505050506040513d60208110156109e957600080fd5b505115156109f657600080fd5b60408051600160a060020a03331681526020810185905281517f8dce8f4eb4097fbdf948a109703513f52a7fabcc7e328ba7f29ad8165033e953929181900390910190a150505050565b6001606060020a0319166000908152600160208190526040909120015490565b60085460095482565b606080600080600354604051908082528060200260200182016040528015610a9b578160200160208202803883390190505b50925060009150610aac6002610dc4565b90505b600354821015610b1957610ac2816107c8565b1515610ada57610ad3600282610de9565b9050610aaf565b808383815181101515610ae957fe5b6001606060020a0319909216602092830290910190910152610b0c600282610de9565b6001909201919050610aaf565b509092915050565b600754600954600091014211610b3657600080fd5b5060408051808201825260001943014080825260075460098054909101602090930183905260088290559190915560045460035590517fe358419ca0dd7928a310d787a606dfae5d869f5071249efa6107105e7afc40bc90600090a150565b60055481565b60065481565b6001606060020a0319166000908152600160205260409020600201541590565b80610bcb816107c8565b1515610bd657600080fd5b6001606060020a03198216600090815260016020526040902054829033600160a060020a03908116911614610c0a57600080fd5b6007546009546001606060020a031985166000818152600160209081526040918290209390940160039093019290925560048054600019019055815190815290517fe3365373f442312a7d66aa1003b60fd4a0f9ed1c85fd8bb384240f8eb4fe916c929181900390910190a1505050565b6001606060020a031916600090815260016020526040902054600160a060020a031690565b610caa8282610e30565b15610cb457600080fd5b610cc782610cc184610e50565b83610e78565b5050565b600080610cd88484610e30565b1515610ce357600080fd5b6001606060020a031983161515610cf957610dbe565b50506001606060020a03198082166000818152602085905260408082208054600180830180546c01000000000000000000000000610100948590048102808b168952878920909401805492820282810473ffffffffffffffffffffffffffffffffffffffff1994851617909155998a168852958720805496840490940274ffffffffffffffffffffffffffffffffffffffff00199096169590951790925594909352805474ffffffffffffffffffffffffffffffffffffffffff191690558154169055905b50505050565b600080805260209190915260409020600101546c010000000000000000000000000290565b6000610df58383610e30565b1515610e0057600080fd5b506001606060020a031916600090815260209190915260409020600101546c010000000000000000000000000290565b6001606060020a0319166000908152602091909152604090205460ff1690565b60008080526020829052604090205461010090046c0100000000000000000000000002919050565b6000610e848483610e30565b15610e8e57600080fd5b610e988484610e30565b80610eab57506001606060020a03198316155b1515610eb657600080fd5b506001606060020a0319808316600090815260209490945260408085206001908101805485851680895284892080546c01000000000000000000000000998a900461010090810274ffffffffffffffffffffffffffffffffffffffff00199283161783558287018054958c028c810473ffffffffffffffffffffffffffffffffffffffff199788161790915586549b909a049a9094168a179094559690951688529287208054969093029516949094179055909252815460ff1916179055565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10610fb757805160ff1916838001178555610fe4565b82800160010185558215610fe4579182015b82811115610fe4578251825591602001919060010190610fc9565b50610ff0929150610ff4565b5090565b61100e91905b80821115610ff05760008155600101610ffa565b905600a165627a7a72305820b04a9e37d5d431f9b779dfb5b47f22640b8155859dec1b6100a3b4731a3323dd0029`

// DeployDarknodeRegistry deploys a new Ethereum contract, binding an instance of DarknodeRegistry to it.
func DeployDarknodeRegistry(auth *bind.TransactOpts, backend bind.ContractBackend, _token common.Address, _minimumBond *big.Int, _minimumDarkPoolSize *big.Int, _minimumEpochInterval *big.Int) (common.Address, *types.Transaction, *DarknodeRegistry, error) {
	parsed, err := abi.JSON(strings.NewReader(DarknodeRegistryABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(DarknodeRegistryBin), backend, _token, _minimumBond, _minimumDarkPoolSize, _minimumEpochInterval)
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

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() constant returns(epochhash uint256, timestamp uint256)
func (_DarknodeRegistry *DarknodeRegistryCaller) CurrentEpoch(opts *bind.CallOpts) (struct {
	Epochhash *big.Int
	Timestamp *big.Int
}, error) {
	ret := new(struct {
		Epochhash *big.Int
		Timestamp *big.Int
	})
	out := ret
	err := _DarknodeRegistry.contract.Call(opts, out, "currentEpoch")
	return *ret, err
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() constant returns(epochhash uint256, timestamp uint256)
func (_DarknodeRegistry *DarknodeRegistrySession) CurrentEpoch() (struct {
	Epochhash *big.Int
	Timestamp *big.Int
}, error) {
	return _DarknodeRegistry.Contract.CurrentEpoch(&_DarknodeRegistry.CallOpts)
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() constant returns(epochhash uint256, timestamp uint256)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) CurrentEpoch() (struct {
	Epochhash *big.Int
	Timestamp *big.Int
}, error) {
	return _DarknodeRegistry.Contract.CurrentEpoch(&_DarknodeRegistry.CallOpts)
}

// GetBond is a free data retrieval call binding the contract method 0x68f209eb.
//
// Solidity: function getBond(_darkNodeID bytes20) constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistryCaller) GetBond(opts *bind.CallOpts, _darkNodeID [20]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "getBond", _darkNodeID)
	return *ret0, err
}

// GetBond is a free data retrieval call binding the contract method 0x68f209eb.
//
// Solidity: function getBond(_darkNodeID bytes20) constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistrySession) GetBond(_darkNodeID [20]byte) (*big.Int, error) {
	return _DarknodeRegistry.Contract.GetBond(&_DarknodeRegistry.CallOpts, _darkNodeID)
}

// GetBond is a free data retrieval call binding the contract method 0x68f209eb.
//
// Solidity: function getBond(_darkNodeID bytes20) constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) GetBond(_darkNodeID [20]byte) (*big.Int, error) {
	return _DarknodeRegistry.Contract.GetBond(&_DarknodeRegistry.CallOpts, _darkNodeID)
}

// GetDarkNodes is a free data retrieval call binding the contract method 0x879ae084.
//
// Solidity: function getDarkNodes() constant returns(bytes20[])
func (_DarknodeRegistry *DarknodeRegistryCaller) GetDarkNodes(opts *bind.CallOpts) ([][20]byte, error) {
	var (
		ret0 = new([][20]byte)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "getDarkNodes")
	return *ret0, err
}

// GetDarkNodes is a free data retrieval call binding the contract method 0x879ae084.
//
// Solidity: function getDarkNodes() constant returns(bytes20[])
func (_DarknodeRegistry *DarknodeRegistrySession) GetDarkNodes() ([][20]byte, error) {
	return _DarknodeRegistry.Contract.GetDarkNodes(&_DarknodeRegistry.CallOpts)
}

// GetDarkNodes is a free data retrieval call binding the contract method 0x879ae084.
//
// Solidity: function getDarkNodes() constant returns(bytes20[])
func (_DarknodeRegistry *DarknodeRegistryCallerSession) GetDarkNodes() ([][20]byte, error) {
	return _DarknodeRegistry.Contract.GetDarkNodes(&_DarknodeRegistry.CallOpts)
}

// GetOwner is a free data retrieval call binding the contract method 0xe487eb58.
//
// Solidity: function getOwner(_darkNodeID bytes20) constant returns(address)
func (_DarknodeRegistry *DarknodeRegistryCaller) GetOwner(opts *bind.CallOpts, _darkNodeID [20]byte) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "getOwner", _darkNodeID)
	return *ret0, err
}

// GetOwner is a free data retrieval call binding the contract method 0xe487eb58.
//
// Solidity: function getOwner(_darkNodeID bytes20) constant returns(address)
func (_DarknodeRegistry *DarknodeRegistrySession) GetOwner(_darkNodeID [20]byte) (common.Address, error) {
	return _DarknodeRegistry.Contract.GetOwner(&_DarknodeRegistry.CallOpts, _darkNodeID)
}

// GetOwner is a free data retrieval call binding the contract method 0xe487eb58.
//
// Solidity: function getOwner(_darkNodeID bytes20) constant returns(address)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) GetOwner(_darkNodeID [20]byte) (common.Address, error) {
	return _DarknodeRegistry.Contract.GetOwner(&_DarknodeRegistry.CallOpts, _darkNodeID)
}

// GetPublicKey is a free data retrieval call binding the contract method 0x32ccd52f.
//
// Solidity: function getPublicKey(_darkNodeID bytes20) constant returns(bytes)
func (_DarknodeRegistry *DarknodeRegistryCaller) GetPublicKey(opts *bind.CallOpts, _darkNodeID [20]byte) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "getPublicKey", _darkNodeID)
	return *ret0, err
}

// GetPublicKey is a free data retrieval call binding the contract method 0x32ccd52f.
//
// Solidity: function getPublicKey(_darkNodeID bytes20) constant returns(bytes)
func (_DarknodeRegistry *DarknodeRegistrySession) GetPublicKey(_darkNodeID [20]byte) ([]byte, error) {
	return _DarknodeRegistry.Contract.GetPublicKey(&_DarknodeRegistry.CallOpts, _darkNodeID)
}

// GetPublicKey is a free data retrieval call binding the contract method 0x32ccd52f.
//
// Solidity: function getPublicKey(_darkNodeID bytes20) constant returns(bytes)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) GetPublicKey(_darkNodeID [20]byte) ([]byte, error) {
	return _DarknodeRegistry.Contract.GetPublicKey(&_DarknodeRegistry.CallOpts, _darkNodeID)
}

// IsDeregistered is a free data retrieval call binding the contract method 0x171f6ea8.
//
// Solidity: function isDeregistered(_darkNodeID bytes20) constant returns(bool)
func (_DarknodeRegistry *DarknodeRegistryCaller) IsDeregistered(opts *bind.CallOpts, _darkNodeID [20]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "isDeregistered", _darkNodeID)
	return *ret0, err
}

// IsDeregistered is a free data retrieval call binding the contract method 0x171f6ea8.
//
// Solidity: function isDeregistered(_darkNodeID bytes20) constant returns(bool)
func (_DarknodeRegistry *DarknodeRegistrySession) IsDeregistered(_darkNodeID [20]byte) (bool, error) {
	return _DarknodeRegistry.Contract.IsDeregistered(&_DarknodeRegistry.CallOpts, _darkNodeID)
}

// IsDeregistered is a free data retrieval call binding the contract method 0x171f6ea8.
//
// Solidity: function isDeregistered(_darkNodeID bytes20) constant returns(bool)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) IsDeregistered(_darkNodeID [20]byte) (bool, error) {
	return _DarknodeRegistry.Contract.IsDeregistered(&_DarknodeRegistry.CallOpts, _darkNodeID)
}

// IsRegistered is a free data retrieval call binding the contract method 0x4f5550fc.
//
// Solidity: function isRegistered(_darkNodeID bytes20) constant returns(bool)
func (_DarknodeRegistry *DarknodeRegistryCaller) IsRegistered(opts *bind.CallOpts, _darkNodeID [20]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "isRegistered", _darkNodeID)
	return *ret0, err
}

// IsRegistered is a free data retrieval call binding the contract method 0x4f5550fc.
//
// Solidity: function isRegistered(_darkNodeID bytes20) constant returns(bool)
func (_DarknodeRegistry *DarknodeRegistrySession) IsRegistered(_darkNodeID [20]byte) (bool, error) {
	return _DarknodeRegistry.Contract.IsRegistered(&_DarknodeRegistry.CallOpts, _darkNodeID)
}

// IsRegistered is a free data retrieval call binding the contract method 0x4f5550fc.
//
// Solidity: function isRegistered(_darkNodeID bytes20) constant returns(bool)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) IsRegistered(_darkNodeID [20]byte) (bool, error) {
	return _DarknodeRegistry.Contract.IsRegistered(&_DarknodeRegistry.CallOpts, _darkNodeID)
}

// IsUnregistered is a free data retrieval call binding the contract method 0xd3841c25.
//
// Solidity: function isUnregistered(_darkNodeID bytes20) constant returns(bool)
func (_DarknodeRegistry *DarknodeRegistryCaller) IsUnregistered(opts *bind.CallOpts, _darkNodeID [20]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "isUnregistered", _darkNodeID)
	return *ret0, err
}

// IsUnregistered is a free data retrieval call binding the contract method 0xd3841c25.
//
// Solidity: function isUnregistered(_darkNodeID bytes20) constant returns(bool)
func (_DarknodeRegistry *DarknodeRegistrySession) IsUnregistered(_darkNodeID [20]byte) (bool, error) {
	return _DarknodeRegistry.Contract.IsUnregistered(&_DarknodeRegistry.CallOpts, _darkNodeID)
}

// IsUnregistered is a free data retrieval call binding the contract method 0xd3841c25.
//
// Solidity: function isUnregistered(_darkNodeID bytes20) constant returns(bool)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) IsUnregistered(_darkNodeID [20]byte) (bool, error) {
	return _DarknodeRegistry.Contract.IsUnregistered(&_DarknodeRegistry.CallOpts, _darkNodeID)
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

// MinimumDarkPoolSize is a free data retrieval call binding the contract method 0xb31575d5.
//
// Solidity: function minimumDarkPoolSize() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistryCaller) MinimumDarkPoolSize(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "minimumDarkPoolSize")
	return *ret0, err
}

// MinimumDarkPoolSize is a free data retrieval call binding the contract method 0xb31575d5.
//
// Solidity: function minimumDarkPoolSize() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistrySession) MinimumDarkPoolSize() (*big.Int, error) {
	return _DarknodeRegistry.Contract.MinimumDarkPoolSize(&_DarknodeRegistry.CallOpts)
}

// MinimumDarkPoolSize is a free data retrieval call binding the contract method 0xb31575d5.
//
// Solidity: function minimumDarkPoolSize() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) MinimumDarkPoolSize() (*big.Int, error) {
	return _DarknodeRegistry.Contract.MinimumDarkPoolSize(&_DarknodeRegistry.CallOpts)
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

// NumDarkNodes is a free data retrieval call binding the contract method 0x060a2cfd.
//
// Solidity: function numDarkNodes() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistryCaller) NumDarkNodes(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "numDarkNodes")
	return *ret0, err
}

// NumDarkNodes is a free data retrieval call binding the contract method 0x060a2cfd.
//
// Solidity: function numDarkNodes() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistrySession) NumDarkNodes() (*big.Int, error) {
	return _DarknodeRegistry.Contract.NumDarkNodes(&_DarknodeRegistry.CallOpts)
}

// NumDarkNodes is a free data retrieval call binding the contract method 0x060a2cfd.
//
// Solidity: function numDarkNodes() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) NumDarkNodes() (*big.Int, error) {
	return _DarknodeRegistry.Contract.NumDarkNodes(&_DarknodeRegistry.CallOpts)
}

// NumDarkNodesNextEpoch is a free data retrieval call binding the contract method 0x0620eb92.
//
// Solidity: function numDarkNodesNextEpoch() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistryCaller) NumDarkNodesNextEpoch(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DarknodeRegistry.contract.Call(opts, out, "numDarkNodesNextEpoch")
	return *ret0, err
}

// NumDarkNodesNextEpoch is a free data retrieval call binding the contract method 0x0620eb92.
//
// Solidity: function numDarkNodesNextEpoch() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistrySession) NumDarkNodesNextEpoch() (*big.Int, error) {
	return _DarknodeRegistry.Contract.NumDarkNodesNextEpoch(&_DarknodeRegistry.CallOpts)
}

// NumDarkNodesNextEpoch is a free data retrieval call binding the contract method 0x0620eb92.
//
// Solidity: function numDarkNodesNextEpoch() constant returns(uint256)
func (_DarknodeRegistry *DarknodeRegistryCallerSession) NumDarkNodesNextEpoch() (*big.Int, error) {
	return _DarknodeRegistry.Contract.NumDarkNodesNextEpoch(&_DarknodeRegistry.CallOpts)
}

// Deregister is a paid mutator transaction binding the contract method 0xe08b4c8a.
//
// Solidity: function deregister(_darkNodeID bytes20) returns()
func (_DarknodeRegistry *DarknodeRegistryTransactor) Deregister(opts *bind.TransactOpts, _darkNodeID [20]byte) (*types.Transaction, error) {
	return _DarknodeRegistry.contract.Transact(opts, "deregister", _darkNodeID)
}

// Deregister is a paid mutator transaction binding the contract method 0xe08b4c8a.
//
// Solidity: function deregister(_darkNodeID bytes20) returns()
func (_DarknodeRegistry *DarknodeRegistrySession) Deregister(_darkNodeID [20]byte) (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.Deregister(&_DarknodeRegistry.TransactOpts, _darkNodeID)
}

// Deregister is a paid mutator transaction binding the contract method 0xe08b4c8a.
//
// Solidity: function deregister(_darkNodeID bytes20) returns()
func (_DarknodeRegistry *DarknodeRegistryTransactorSession) Deregister(_darkNodeID [20]byte) (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.Deregister(&_DarknodeRegistry.TransactOpts, _darkNodeID)
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

// Refund is a paid mutator transaction binding the contract method 0x5a8f9b81.
//
// Solidity: function refund(_darkNodeID bytes20) returns()
func (_DarknodeRegistry *DarknodeRegistryTransactor) Refund(opts *bind.TransactOpts, _darkNodeID [20]byte) (*types.Transaction, error) {
	return _DarknodeRegistry.contract.Transact(opts, "refund", _darkNodeID)
}

// Refund is a paid mutator transaction binding the contract method 0x5a8f9b81.
//
// Solidity: function refund(_darkNodeID bytes20) returns()
func (_DarknodeRegistry *DarknodeRegistrySession) Refund(_darkNodeID [20]byte) (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.Refund(&_DarknodeRegistry.TransactOpts, _darkNodeID)
}

// Refund is a paid mutator transaction binding the contract method 0x5a8f9b81.
//
// Solidity: function refund(_darkNodeID bytes20) returns()
func (_DarknodeRegistry *DarknodeRegistryTransactorSession) Refund(_darkNodeID [20]byte) (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.Refund(&_DarknodeRegistry.TransactOpts, _darkNodeID)
}

// Register is a paid mutator transaction binding the contract method 0x375a8be3.
//
// Solidity: function register(_darkNodeID bytes20, _publicKey bytes, _bond uint256) returns()
func (_DarknodeRegistry *DarknodeRegistryTransactor) Register(opts *bind.TransactOpts, _darkNodeID [20]byte, _publicKey []byte, _bond *big.Int) (*types.Transaction, error) {
	return _DarknodeRegistry.contract.Transact(opts, "register", _darkNodeID, _publicKey, _bond)
}

// Register is a paid mutator transaction binding the contract method 0x375a8be3.
//
// Solidity: function register(_darkNodeID bytes20, _publicKey bytes, _bond uint256) returns()
func (_DarknodeRegistry *DarknodeRegistrySession) Register(_darkNodeID [20]byte, _publicKey []byte, _bond *big.Int) (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.Register(&_DarknodeRegistry.TransactOpts, _darkNodeID, _publicKey, _bond)
}

// Register is a paid mutator transaction binding the contract method 0x375a8be3.
//
// Solidity: function register(_darkNodeID bytes20, _publicKey bytes, _bond uint256) returns()
func (_DarknodeRegistry *DarknodeRegistryTransactorSession) Register(_darkNodeID [20]byte, _publicKey []byte, _bond *big.Int) (*types.Transaction, error) {
	return _DarknodeRegistry.Contract.Register(&_DarknodeRegistry.TransactOpts, _darkNodeID, _publicKey, _bond)
}

// DarknodeRegistryDarkNodeDeregisteredIterator is returned from FilterDarkNodeDeregistered and is used to iterate over the raw logs and unpacked data for DarkNodeDeregistered events raised by the DarknodeRegistry contract.
type DarknodeRegistryDarkNodeDeregisteredIterator struct {
	Event *DarknodeRegistryDarkNodeDeregistered // Event containing the contract specifics and raw log

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
func (it *DarknodeRegistryDarkNodeDeregisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DarknodeRegistryDarkNodeDeregistered)
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
		it.Event = new(DarknodeRegistryDarkNodeDeregistered)
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
func (it *DarknodeRegistryDarkNodeDeregisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DarknodeRegistryDarkNodeDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DarknodeRegistryDarkNodeDeregistered represents a DarkNodeDeregistered event raised by the DarknodeRegistry contract.
type DarknodeRegistryDarkNodeDeregistered struct {
	DarkNodeID [20]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterDarkNodeDeregistered is a free log retrieval operation binding the contract event 0xe3365373f442312a7d66aa1003b60fd4a0f9ed1c85fd8bb384240f8eb4fe916c.
//
// Solidity: event DarkNodeDeregistered(_darkNodeID bytes20)
func (_DarknodeRegistry *DarknodeRegistryFilterer) FilterDarkNodeDeregistered(opts *bind.FilterOpts) (*DarknodeRegistryDarkNodeDeregisteredIterator, error) {

	logs, sub, err := _DarknodeRegistry.contract.FilterLogs(opts, "DarkNodeDeregistered")
	if err != nil {
		return nil, err
	}
	return &DarknodeRegistryDarkNodeDeregisteredIterator{contract: _DarknodeRegistry.contract, event: "DarkNodeDeregistered", logs: logs, sub: sub}, nil
}

// WatchDarkNodeDeregistered is a free log subscription operation binding the contract event 0xe3365373f442312a7d66aa1003b60fd4a0f9ed1c85fd8bb384240f8eb4fe916c.
//
// Solidity: event DarkNodeDeregistered(_darkNodeID bytes20)
func (_DarknodeRegistry *DarknodeRegistryFilterer) WatchDarkNodeDeregistered(opts *bind.WatchOpts, sink chan<- *DarknodeRegistryDarkNodeDeregistered) (event.Subscription, error) {

	logs, sub, err := _DarknodeRegistry.contract.WatchLogs(opts, "DarkNodeDeregistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DarknodeRegistryDarkNodeDeregistered)
				if err := _DarknodeRegistry.contract.UnpackLog(event, "DarkNodeDeregistered", log); err != nil {
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

// DarknodeRegistryDarkNodeRegisteredIterator is returned from FilterDarkNodeRegistered and is used to iterate over the raw logs and unpacked data for DarkNodeRegistered events raised by the DarknodeRegistry contract.
type DarknodeRegistryDarkNodeRegisteredIterator struct {
	Event *DarknodeRegistryDarkNodeRegistered // Event containing the contract specifics and raw log

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
func (it *DarknodeRegistryDarkNodeRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DarknodeRegistryDarkNodeRegistered)
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
		it.Event = new(DarknodeRegistryDarkNodeRegistered)
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
func (it *DarknodeRegistryDarkNodeRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DarknodeRegistryDarkNodeRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DarknodeRegistryDarkNodeRegistered represents a DarkNodeRegistered event raised by the DarknodeRegistry contract.
type DarknodeRegistryDarkNodeRegistered struct {
	DarkNodeID [20]byte
	Bond       *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterDarkNodeRegistered is a free log retrieval operation binding the contract event 0xcde427a7822769e806a704726602a3a8a44458da1352d9b5a46bcfb95142ad38.
//
// Solidity: event DarkNodeRegistered(_darkNodeID bytes20, _bond uint256)
func (_DarknodeRegistry *DarknodeRegistryFilterer) FilterDarkNodeRegistered(opts *bind.FilterOpts) (*DarknodeRegistryDarkNodeRegisteredIterator, error) {

	logs, sub, err := _DarknodeRegistry.contract.FilterLogs(opts, "DarkNodeRegistered")
	if err != nil {
		return nil, err
	}
	return &DarknodeRegistryDarkNodeRegisteredIterator{contract: _DarknodeRegistry.contract, event: "DarkNodeRegistered", logs: logs, sub: sub}, nil
}

// WatchDarkNodeRegistered is a free log subscription operation binding the contract event 0xcde427a7822769e806a704726602a3a8a44458da1352d9b5a46bcfb95142ad38.
//
// Solidity: event DarkNodeRegistered(_darkNodeID bytes20, _bond uint256)
func (_DarknodeRegistry *DarknodeRegistryFilterer) WatchDarkNodeRegistered(opts *bind.WatchOpts, sink chan<- *DarknodeRegistryDarkNodeRegistered) (event.Subscription, error) {

	logs, sub, err := _DarknodeRegistry.contract.WatchLogs(opts, "DarkNodeRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DarknodeRegistryDarkNodeRegistered)
				if err := _DarknodeRegistry.contract.UnpackLog(event, "DarkNodeRegistered", log); err != nil {
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

// DarknodeRegistryNewEpochIterator is returned from FilterNewEpoch and is used to iterate over the raw logs and unpacked data for NewEpoch events raised by the DarknodeRegistry contract.
type DarknodeRegistryNewEpochIterator struct {
	Event *DarknodeRegistryNewEpoch // Event containing the contract specifics and raw log

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
func (it *DarknodeRegistryNewEpochIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DarknodeRegistryNewEpoch)
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
		it.Event = new(DarknodeRegistryNewEpoch)
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
func (it *DarknodeRegistryNewEpochIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DarknodeRegistryNewEpochIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DarknodeRegistryNewEpoch represents a NewEpoch event raised by the DarknodeRegistry contract.
type DarknodeRegistryNewEpoch struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterNewEpoch is a free log retrieval operation binding the contract event 0xe358419ca0dd7928a310d787a606dfae5d869f5071249efa6107105e7afc40bc.
//
// Solidity: event NewEpoch()
func (_DarknodeRegistry *DarknodeRegistryFilterer) FilterNewEpoch(opts *bind.FilterOpts) (*DarknodeRegistryNewEpochIterator, error) {

	logs, sub, err := _DarknodeRegistry.contract.FilterLogs(opts, "NewEpoch")
	if err != nil {
		return nil, err
	}
	return &DarknodeRegistryNewEpochIterator{contract: _DarknodeRegistry.contract, event: "NewEpoch", logs: logs, sub: sub}, nil
}

// WatchNewEpoch is a free log subscription operation binding the contract event 0xe358419ca0dd7928a310d787a606dfae5d869f5071249efa6107105e7afc40bc.
//
// Solidity: event NewEpoch()
func (_DarknodeRegistry *DarknodeRegistryFilterer) WatchNewEpoch(opts *bind.WatchOpts, sink chan<- *DarknodeRegistryNewEpoch) (event.Subscription, error) {

	logs, sub, err := _DarknodeRegistry.contract.WatchLogs(opts, "NewEpoch")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DarknodeRegistryNewEpoch)
				if err := _DarknodeRegistry.contract.UnpackLog(event, "NewEpoch", log); err != nil {
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

// DarknodeRegistryOwnerRefundedIterator is returned from FilterOwnerRefunded and is used to iterate over the raw logs and unpacked data for OwnerRefunded events raised by the DarknodeRegistry contract.
type DarknodeRegistryOwnerRefundedIterator struct {
	Event *DarknodeRegistryOwnerRefunded // Event containing the contract specifics and raw log

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
func (it *DarknodeRegistryOwnerRefundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DarknodeRegistryOwnerRefunded)
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
		it.Event = new(DarknodeRegistryOwnerRefunded)
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
func (it *DarknodeRegistryOwnerRefundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DarknodeRegistryOwnerRefundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DarknodeRegistryOwnerRefunded represents a OwnerRefunded event raised by the DarknodeRegistry contract.
type DarknodeRegistryOwnerRefunded struct {
	Owner  common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterOwnerRefunded is a free log retrieval operation binding the contract event 0x8dce8f4eb4097fbdf948a109703513f52a7fabcc7e328ba7f29ad8165033e953.
//
// Solidity: event OwnerRefunded(_owner address, _amount uint256)
func (_DarknodeRegistry *DarknodeRegistryFilterer) FilterOwnerRefunded(opts *bind.FilterOpts) (*DarknodeRegistryOwnerRefundedIterator, error) {

	logs, sub, err := _DarknodeRegistry.contract.FilterLogs(opts, "OwnerRefunded")
	if err != nil {
		return nil, err
	}
	return &DarknodeRegistryOwnerRefundedIterator{contract: _DarknodeRegistry.contract, event: "OwnerRefunded", logs: logs, sub: sub}, nil
}

// WatchOwnerRefunded is a free log subscription operation binding the contract event 0x8dce8f4eb4097fbdf948a109703513f52a7fabcc7e328ba7f29ad8165033e953.
//
// Solidity: event OwnerRefunded(_owner address, _amount uint256)
func (_DarknodeRegistry *DarknodeRegistryFilterer) WatchOwnerRefunded(opts *bind.WatchOpts, sink chan<- *DarknodeRegistryOwnerRefunded) (event.Subscription, error) {

	logs, sub, err := _DarknodeRegistry.contract.WatchLogs(opts, "OwnerRefunded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DarknodeRegistryOwnerRefunded)
				if err := _DarknodeRegistry.contract.UnpackLog(event, "OwnerRefunded", log); err != nil {
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
// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
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
// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
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

// LibRewardVaultABI is the input ABI used to generate the binding from.
const LibRewardVaultABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"ETHEREUM\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// LibRewardVaultBin is the compiled bytecode used for deploying new contracts.
const LibRewardVaultBin = `0x60cd61002f600b82828239805160001a6073146000811461001f57610021565bfe5b5030600052607381538281f300730000000000000000000000000000000000000000301460806040526004361060555763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663f7cdf47c8114605a575b600080fd5b60606089565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b73eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee815600a165627a7a723058208e82ab344819676a16e7603bbc8bc581e01950e21ca004f3ae4b6ed04cd155ed0029`

// DeployLibRewardVault deploys a new Ethereum contract, binding an instance of LibRewardVault to it.
func DeployLibRewardVault(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *LibRewardVault, error) {
	parsed, err := abi.JSON(strings.NewReader(LibRewardVaultABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(LibRewardVaultBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LibRewardVault{LibRewardVaultCaller: LibRewardVaultCaller{contract: contract}, LibRewardVaultTransactor: LibRewardVaultTransactor{contract: contract}, LibRewardVaultFilterer: LibRewardVaultFilterer{contract: contract}}, nil
}

// LibRewardVault is an auto generated Go binding around an Ethereum contract.
type LibRewardVault struct {
	LibRewardVaultCaller     // Read-only binding to the contract
	LibRewardVaultTransactor // Write-only binding to the contract
	LibRewardVaultFilterer   // Log filterer for contract events
}

// LibRewardVaultCaller is an auto generated read-only Go binding around an Ethereum contract.
type LibRewardVaultCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LibRewardVaultTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LibRewardVaultTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LibRewardVaultFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LibRewardVaultFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LibRewardVaultSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LibRewardVaultSession struct {
	Contract     *LibRewardVault   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LibRewardVaultCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LibRewardVaultCallerSession struct {
	Contract *LibRewardVaultCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// LibRewardVaultTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LibRewardVaultTransactorSession struct {
	Contract     *LibRewardVaultTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// LibRewardVaultRaw is an auto generated low-level Go binding around an Ethereum contract.
type LibRewardVaultRaw struct {
	Contract *LibRewardVault // Generic contract binding to access the raw methods on
}

// LibRewardVaultCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LibRewardVaultCallerRaw struct {
	Contract *LibRewardVaultCaller // Generic read-only contract binding to access the raw methods on
}

// LibRewardVaultTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LibRewardVaultTransactorRaw struct {
	Contract *LibRewardVaultTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLibRewardVault creates a new instance of LibRewardVault, bound to a specific deployed contract.
func NewLibRewardVault(address common.Address, backend bind.ContractBackend) (*LibRewardVault, error) {
	contract, err := bindLibRewardVault(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LibRewardVault{LibRewardVaultCaller: LibRewardVaultCaller{contract: contract}, LibRewardVaultTransactor: LibRewardVaultTransactor{contract: contract}, LibRewardVaultFilterer: LibRewardVaultFilterer{contract: contract}}, nil
}

// NewLibRewardVaultCaller creates a new read-only instance of LibRewardVault, bound to a specific deployed contract.
func NewLibRewardVaultCaller(address common.Address, caller bind.ContractCaller) (*LibRewardVaultCaller, error) {
	contract, err := bindLibRewardVault(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LibRewardVaultCaller{contract: contract}, nil
}

// NewLibRewardVaultTransactor creates a new write-only instance of LibRewardVault, bound to a specific deployed contract.
func NewLibRewardVaultTransactor(address common.Address, transactor bind.ContractTransactor) (*LibRewardVaultTransactor, error) {
	contract, err := bindLibRewardVault(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LibRewardVaultTransactor{contract: contract}, nil
}

// NewLibRewardVaultFilterer creates a new log filterer instance of LibRewardVault, bound to a specific deployed contract.
func NewLibRewardVaultFilterer(address common.Address, filterer bind.ContractFilterer) (*LibRewardVaultFilterer, error) {
	contract, err := bindLibRewardVault(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LibRewardVaultFilterer{contract: contract}, nil
}

// bindLibRewardVault binds a generic wrapper to an already deployed contract.
func bindLibRewardVault(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LibRewardVaultABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LibRewardVault *LibRewardVaultRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LibRewardVault.Contract.LibRewardVaultCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LibRewardVault *LibRewardVaultRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LibRewardVault.Contract.LibRewardVaultTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LibRewardVault *LibRewardVaultRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LibRewardVault.Contract.LibRewardVaultTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LibRewardVault *LibRewardVaultCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LibRewardVault.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LibRewardVault *LibRewardVaultTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LibRewardVault.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LibRewardVault *LibRewardVaultTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LibRewardVault.Contract.contract.Transact(opts, method, params...)
}

// ETHEREUM is a free data retrieval call binding the contract method 0xf7cdf47c.
//
// Solidity: function ETHEREUM() constant returns(address)
func (_LibRewardVault *LibRewardVaultCaller) ETHEREUM(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _LibRewardVault.contract.Call(opts, out, "ETHEREUM")
	return *ret0, err
}

// ETHEREUM is a free data retrieval call binding the contract method 0xf7cdf47c.
//
// Solidity: function ETHEREUM() constant returns(address)
func (_LibRewardVault *LibRewardVaultSession) ETHEREUM() (common.Address, error) {
	return _LibRewardVault.Contract.ETHEREUM(&_LibRewardVault.CallOpts)
}

// ETHEREUM is a free data retrieval call binding the contract method 0xf7cdf47c.
//
// Solidity: function ETHEREUM() constant returns(address)
func (_LibRewardVault *LibRewardVaultCallerSession) ETHEREUM() (common.Address, error) {
	return _LibRewardVault.Contract.ETHEREUM(&_LibRewardVault.CallOpts)
}

// LinkedListABI is the input ABI used to generate the binding from.
const LinkedListABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"NULL\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes20\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// LinkedListBin is the compiled bytecode used for deploying new contracts.
const LinkedListBin = `0x60b361002f600b82828239805160001a6073146000811461001f57610021565bfe5b5030600052607381538281f300730000000000000000000000000000000000000000301460806040526004361060555763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663f26be3fc8114605a575b600080fd5b60606082565b604080516bffffffffffffffffffffffff199092168252519081900360200190f35b6000815600a165627a7a723058202c6ef008fac80160248888c2f9a11ea516528c11ff4b3124215f7aa182ccb6fb0029`

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

// OwnableABI is the input ABI used to generate the binding from.
const OwnableABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"

// OwnableBin is the compiled bytecode used for deploying new contracts.
const OwnableBin = `0x608060405234801561001057600080fd5b5060008054600160a060020a033316600160a060020a03199091161790556101778061003d6000396000f30060806040526004361061004b5763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416638da5cb5b8114610050578063f2fde38b14610081575b600080fd5b34801561005c57600080fd5b506100656100a4565b60408051600160a060020a039092168252519081900360200190f35b34801561008d57600080fd5b506100a2600160a060020a03600435166100b3565b005b600054600160a060020a031681565b60005433600160a060020a039081169116146100ce57600080fd5b600160a060020a03811615156100e357600080fd5b60008054604051600160a060020a03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a36000805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03929092169190911790555600a165627a7a723058207e2865964edcd0615b2c2f0d2789abe9b47363f2b5a5a971fea2a4746e949a410029`

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

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_Ownable *OwnableTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Ownable.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_Ownable *OwnableSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Ownable.Contract.TransferOwnership(&_Ownable.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_Ownable *OwnableTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Ownable.Contract.TransferOwnership(&_Ownable.TransactOpts, newOwner)
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
// Solidity: event OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
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
// Solidity: event OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
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
const PausableABI = "[{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"

// PausableBin is the compiled bytecode used for deploying new contracts.
const PausableBin = `0x608060405260008054600160a060020a033316600160a860020a0319909116179055610338806100306000396000f30060806040526004361061006c5763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416633f4ba83a81146100715780635c975abb146100885780638456cb59146100b15780638da5cb5b146100c6578063f2fde38b146100f7575b600080fd5b34801561007d57600080fd5b50610086610118565b005b34801561009457600080fd5b5061009d6101a3565b604080519115158252519081900360200190f35b3480156100bd57600080fd5b506100866101c4565b3480156100d257600080fd5b506100db610265565b60408051600160a060020a039092168252519081900360200190f35b34801561010357600080fd5b50610086600160a060020a0360043516610274565b60005433600160a060020a0390811691161461013357600080fd5b60005474010000000000000000000000000000000000000000900460ff16151561015c57600080fd5b6000805474ff0000000000000000000000000000000000000000191681556040517f7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b339190a1565b60005474010000000000000000000000000000000000000000900460ff1681565b60005433600160a060020a039081169116146101df57600080fd5b60005474010000000000000000000000000000000000000000900460ff161561020757600080fd5b6000805474ff00000000000000000000000000000000000000001916740100000000000000000000000000000000000000001781556040517f6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff6259190a1565b600054600160a060020a031681565b60005433600160a060020a0390811691161461028f57600080fd5b600160a060020a03811615156102a457600080fd5b60008054604051600160a060020a03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a36000805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03929092169190911790555600a165627a7a723058200bf7ebc7c04d6f50a0b5a72c43b5708c1ae7dd1b0a3c960396901d4cb30d54be0029`

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

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_Pausable *PausableTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Pausable.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_Pausable *PausableSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Pausable.Contract.TransferOwnership(&_Pausable.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_Pausable *PausableTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Pausable.Contract.TransferOwnership(&_Pausable.TransactOpts, newOwner)
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
// Solidity: event OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
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
// Solidity: event OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
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
// Solidity: event Pause()
func (_Pausable *PausableFilterer) FilterPause(opts *bind.FilterOpts) (*PausablePauseIterator, error) {

	logs, sub, err := _Pausable.contract.FilterLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return &PausablePauseIterator{contract: _Pausable.contract, event: "Pause", logs: logs, sub: sub}, nil
}

// WatchPause is a free log subscription operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
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
// Solidity: event Unpause()
func (_Pausable *PausableFilterer) FilterUnpause(opts *bind.FilterOpts) (*PausableUnpauseIterator, error) {

	logs, sub, err := _Pausable.contract.FilterLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return &PausableUnpauseIterator{contract: _Pausable.contract, event: "Unpause", logs: logs, sub: sub}, nil
}

// WatchUnpause is a free log subscription operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
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
const PausableTokenABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseApproval\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseApproval\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]"

// PausableTokenBin is the compiled bytecode used for deploying new contracts.
const PausableTokenBin = `0x608060405260038054600160a860020a03191633600160a060020a0316179055610a298061002e6000396000f3006080604052600436106100c45763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663095ea7b381146100c957806318160ddd1461010157806323b872dd146101285780633f4ba83a146101525780635c975abb14610169578063661884631461017e57806370a08231146101a25780638456cb59146101c35780638da5cb5b146101d8578063a9059cbb14610209578063d73dd6231461022d578063dd62ed3e14610251578063f2fde38b14610278575b600080fd5b3480156100d557600080fd5b506100ed600160a060020a0360043516602435610299565b604080519115158252519081900360200190f35b34801561010d57600080fd5b506101166102c4565b60408051918252519081900360200190f35b34801561013457600080fd5b506100ed600160a060020a03600435811690602435166044356102ca565b34801561015e57600080fd5b506101676102f7565b005b34801561017557600080fd5b506100ed610373565b34801561018a57600080fd5b506100ed600160a060020a0360043516602435610383565b3480156101ae57600080fd5b50610116600160a060020a03600435166103a7565b3480156101cf57600080fd5b506101676103c2565b3480156101e457600080fd5b506101ed610443565b60408051600160a060020a039092168252519081900360200190f35b34801561021557600080fd5b506100ed600160a060020a0360043516602435610452565b34801561023957600080fd5b506100ed600160a060020a0360043516602435610476565b34801561025d57600080fd5b50610116600160a060020a036004358116906024351661049a565b34801561028457600080fd5b50610167600160a060020a03600435166104c5565b60035460009060a060020a900460ff16156102b357600080fd5b6102bd838361055e565b9392505050565b60015490565b60035460009060a060020a900460ff16156102e457600080fd5b6102ef8484846105c8565b949350505050565b60035433600160a060020a0390811691161461031257600080fd5b60035460a060020a900460ff16151561032a57600080fd5b6003805474ff0000000000000000000000000000000000000000191690556040517f7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b3390600090a1565b60035460a060020a900460ff1681565b60035460009060a060020a900460ff161561039d57600080fd5b6102bd8383610748565b600160a060020a031660009081526020819052604090205490565b60035433600160a060020a039081169116146103dd57600080fd5b60035460a060020a900460ff16156103f457600080fd5b6003805474ff0000000000000000000000000000000000000000191660a060020a1790556040517f6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff62590600090a1565b600354600160a060020a031681565b60035460009060a060020a900460ff161561046c57600080fd5b6102bd8383610841565b60035460009060a060020a900460ff161561049057600080fd5b6102bd838361093a565b600160a060020a03918216600090815260026020908152604080832093909416825291909152205490565b60035433600160a060020a039081169116146104e057600080fd5b600160a060020a03811615156104f557600080fd5b600354604051600160a060020a038084169216907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a36003805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b600160a060020a03338116600081815260026020908152604080832094871680845294825280832086905580518681529051929493927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925929181900390910190a350600192915050565b6000600160a060020a03831615156105df57600080fd5b600160a060020a03841660009081526020819052604090205482111561060457600080fd5b600160a060020a038085166000908152600260209081526040808320339094168352929052205482111561063757600080fd5b600160a060020a038416600090815260208190526040902054610660908363ffffffff6109dc16565b600160a060020a038086166000908152602081905260408082209390935590851681522054610695908363ffffffff6109ee16565b600160a060020a03808516600090815260208181526040808320949094558783168252600281528382203390931682529190915220546106db908363ffffffff6109dc16565b600160a060020a038086166000818152600260209081526040808320338616845282529182902094909455805186815290519287169391927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef929181900390910190a35060019392505050565b600160a060020a033381166000908152600260209081526040808320938616835292905290812054808311156107a557600160a060020a0333811660009081526002602090815260408083209388168352929052908120556107dc565b6107b5818463ffffffff6109dc16565b600160a060020a033381166000908152600260209081526040808320938916835292905220555b600160a060020a0333811660008181526002602090815260408083209489168084529482529182902054825190815291517f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259281900390910190a35060019392505050565b6000600160a060020a038316151561085857600080fd5b600160a060020a03331660009081526020819052604090205482111561087d57600080fd5b600160a060020a0333166000908152602081905260409020546108a6908363ffffffff6109dc16565b600160a060020a0333811660009081526020819052604080822093909355908516815220546108db908363ffffffff6109ee16565b600160a060020a03808516600081815260208181526040918290209490945580518681529051919333909316927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef92918290030190a350600192915050565b600160a060020a033381166000908152600260209081526040808320938616835292905290812054610972908363ffffffff6109ee16565b600160a060020a0333811660008181526002602090815260408083209489168084529482529182902085905581519485529051929391927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259281900390910190a350600192915050565b6000828211156109e857fe5b50900390565b6000828201838110156102bd57fe00a165627a7a723058206ed6f29b6553180489719ab4b7fb9774fb87647da56742e09ba3fe97794828d40029`

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
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
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
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_PausableToken *PausableTokenSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _PausableToken.Contract.BalanceOf(&_PausableToken.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
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
// Solidity: function transferOwnership(newOwner address) returns()
func (_PausableToken *PausableTokenTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _PausableToken.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_PausableToken *PausableTokenSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _PausableToken.Contract.TransferOwnership(&_PausableToken.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_PausableToken *PausableTokenTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _PausableToken.Contract.TransferOwnership(&_PausableToken.TransactOpts, newOwner)
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
// Solidity: event Approval(owner indexed address, spender indexed address, value uint256)
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
// Solidity: event Approval(owner indexed address, spender indexed address, value uint256)
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
// Solidity: event OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
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
// Solidity: event OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
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
// Solidity: event Pause()
func (_PausableToken *PausableTokenFilterer) FilterPause(opts *bind.FilterOpts) (*PausableTokenPauseIterator, error) {

	logs, sub, err := _PausableToken.contract.FilterLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return &PausableTokenPauseIterator{contract: _PausableToken.contract, event: "Pause", logs: logs, sub: sub}, nil
}

// WatchPause is a free log subscription operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
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
// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
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
// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
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
// Solidity: event Unpause()
func (_PausableToken *PausableTokenFilterer) FilterUnpause(opts *bind.FilterOpts) (*PausableTokenUnpauseIterator, error) {

	logs, sub, err := _PausableToken.contract.FilterLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return &PausableTokenUnpauseIterator{contract: _PausableToken.contract, event: "Unpause", logs: logs, sub: sub}, nil
}

// WatchUnpause is a free log subscription operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
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
const RepublicTokenABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"INITIAL_SUPPLY\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseApproval\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"beneficiary\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferTokens\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseApproval\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"burner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Burn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]"

// RepublicTokenBin is the compiled bytecode used for deploying new contracts.
const RepublicTokenBin = `0x60806040526003805460a060020a60ff021916905534801561002057600080fd5b5060038054600160a060020a033316600160a060020a031990911681179091556b033b2e3c9fd0803ce8000000600481905560009182526020829052604090912055610e57806100716000396000f3006080604052600436106101065763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166306fdde03811461010b578063095ea7b31461019557806318160ddd146101cd57806323b872dd146101f45780632ff2e9dc1461021e578063313ce567146102335780633f4ba83a1461025e57806342966c68146102755780635c975abb1461028d57806366188463146102a257806370a08231146102c65780638456cb59146102e75780638da5cb5b146102fc57806395d89b411461032d578063a9059cbb14610342578063bec3fa1714610366578063d73dd6231461038a578063dd62ed3e146103ae578063f2fde38b146103d5575b600080fd5b34801561011757600080fd5b506101206103f6565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561015a578181015183820152602001610142565b50505050905090810190601f1680156101875780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156101a157600080fd5b506101b9600160a060020a036004351660243561042d565b604080519115158252519081900360200190f35b3480156101d957600080fd5b506101e2610458565b60408051918252519081900360200190f35b34801561020057600080fd5b506101b9600160a060020a036004358116906024351660443561045e565b34801561022a57600080fd5b506101e26104ac565b34801561023f57600080fd5b506102486104bc565b6040805160ff9092168252519081900360200190f35b34801561026a57600080fd5b506102736104c1565b005b34801561028157600080fd5b5061027360043561053d565b34801561029957600080fd5b506101b9610628565b3480156102ae57600080fd5b506101b9600160a060020a0360043516602435610638565b3480156102d257600080fd5b506101e2600160a060020a036004351661065c565b3480156102f357600080fd5b50610273610677565b34801561030857600080fd5b506103116106f8565b60408051600160a060020a039092168252519081900360200190f35b34801561033957600080fd5b50610120610707565b34801561034e57600080fd5b506101b9600160a060020a036004351660243561073e565b34801561037257600080fd5b506101b9600160a060020a0360043516602435610783565b34801561039657600080fd5b506101b9600160a060020a036004351660243561085f565b3480156103ba57600080fd5b506101e2600160a060020a0360043581169060243516610883565b3480156103e157600080fd5b50610273600160a060020a03600435166108ae565b60408051808201909152600e81527f52657075626c696320546f6b656e000000000000000000000000000000000000602082015281565b60035460009060a060020a900460ff161561044757600080fd5b6104518383610947565b9392505050565b60045481565b60035460009060a060020a900460ff161561047857600080fd5b30600160a060020a031683600160a060020a03161415151561049957600080fd5b6104a48484846109b1565b949350505050565b6b033b2e3c9fd0803ce800000081565b601281565b60035433600160a060020a039081169116146104dc57600080fd5b60035460a060020a900460ff1615156104f457600080fd5b6003805474ff0000000000000000000000000000000000000000191690556040517f7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b3390600090a1565b600160a060020a03331660009081526020819052604081205482111561056257600080fd5b5033600160a060020a03811660009081526020819052604090205461058790836109d6565b600160a060020a0382166000908152602081905260409020556001546105b3908363ffffffff6109d616565b600155604080518381529051600160a060020a038316917fcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5919081900360200190a2604080518381529051600091600160a060020a03841691600080516020610e0c8339815191529181900360200190a35050565b60035460a060020a900460ff1681565b60035460009060a060020a900460ff161561065257600080fd5b61045183836109e8565b600160a060020a031660009081526020819052604090205490565b60035433600160a060020a0390811691161461069257600080fd5b60035460a060020a900460ff16156106a957600080fd5b6003805474ff0000000000000000000000000000000000000000191660a060020a1790556040517f6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff62590600090a1565b600354600160a060020a031681565b60408051808201909152600381527f52454e0000000000000000000000000000000000000000000000000000000000602082015281565b60035460009060a060020a900460ff161561075857600080fd5b30600160a060020a031683600160a060020a03161415151561077957600080fd5b6104518383610ae1565b60035460009033600160a060020a039081169116146107a157600080fd5b600082116107ae57600080fd5b600354600160a060020a03166000908152602081905260409020546107d9908363ffffffff6109d616565b600354600160a060020a039081166000908152602081905260408082209390935590851681522054610811908363ffffffff610b0516565b600160a060020a038085166000818152602081815260409182902094909455600354815187815291519294931692600080516020610e0c83398151915292918290030190a350600192915050565b60035460009060a060020a900460ff161561087957600080fd5b6104518383610b14565b600160a060020a03918216600090815260026020908152604080832093909416825291909152205490565b60035433600160a060020a039081169116146108c957600080fd5b600160a060020a03811615156108de57600080fd5b600354604051600160a060020a038084169216907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a36003805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b600160a060020a03338116600081815260026020908152604080832094871680845294825280832086905580518681529051929493927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925929181900390910190a350600192915050565b60035460009060a060020a900460ff16156109cb57600080fd5b6104a4848484610bb6565b6000828211156109e257fe5b50900390565b600160a060020a03338116600090815260026020908152604080832093861683529290529081205480831115610a4557600160a060020a033381166000908152600260209081526040808320938816835292905290812055610a7c565b610a55818463ffffffff6109d616565b600160a060020a033381166000908152600260209081526040808320938916835292905220555b600160a060020a0333811660008181526002602090815260408083209489168084529482529182902054825190815291517f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259281900390910190a35060019392505050565b60035460009060a060020a900460ff1615610afb57600080fd5b6104518383610d24565b60008282018381101561045157fe5b600160a060020a033381166000908152600260209081526040808320938616835292905290812054610b4c908363ffffffff610b0516565b600160a060020a0333811660008181526002602090815260408083209489168084529482529182902085905581519485529051929391927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259281900390910190a350600192915050565b6000600160a060020a0383161515610bcd57600080fd5b600160a060020a038416600090815260208190526040902054821115610bf257600080fd5b600160a060020a0380851660009081526002602090815260408083203390941683529290522054821115610c2557600080fd5b600160a060020a038416600090815260208190526040902054610c4e908363ffffffff6109d616565b600160a060020a038086166000908152602081905260408082209390935590851681522054610c83908363ffffffff610b0516565b600160a060020a0380851660009081526020818152604080832094909455878316825260028152838220339093168252919091522054610cc9908363ffffffff6109d616565b600160a060020a03808616600081815260026020908152604080832033861684528252918290209490945580518681529051928716939192600080516020610e0c833981519152929181900390910190a35060019392505050565b6000600160a060020a0383161515610d3b57600080fd5b600160a060020a033316600090815260208190526040902054821115610d6057600080fd5b600160a060020a033316600090815260208190526040902054610d89908363ffffffff6109d616565b600160a060020a033381166000908152602081905260408082209390935590851681522054610dbe908363ffffffff610b0516565b600160a060020a0380851660008181526020818152604091829020949094558051868152905191933390931692600080516020610e0c83398151915292918290030190a3506001929150505600ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa165627a7a72305820aa9a3f877aa9d2a637bdcdc31a16441f0fd8e3da5576ddd590510e024da0a1890029`

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
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
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
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_RepublicToken *RepublicTokenSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _RepublicToken.Contract.BalanceOf(&_RepublicToken.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
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
// Solidity: function transferOwnership(newOwner address) returns()
func (_RepublicToken *RepublicTokenTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _RepublicToken.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_RepublicToken *RepublicTokenSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _RepublicToken.Contract.TransferOwnership(&_RepublicToken.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_RepublicToken *RepublicTokenTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _RepublicToken.Contract.TransferOwnership(&_RepublicToken.TransactOpts, newOwner)
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
// Solidity: event Approval(owner indexed address, spender indexed address, value uint256)
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
// Solidity: event Approval(owner indexed address, spender indexed address, value uint256)
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
// Solidity: event Burn(burner indexed address, value uint256)
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
// Solidity: event Burn(burner indexed address, value uint256)
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
// Solidity: event OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
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
// Solidity: event OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
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
// Solidity: event Pause()
func (_RepublicToken *RepublicTokenFilterer) FilterPause(opts *bind.FilterOpts) (*RepublicTokenPauseIterator, error) {

	logs, sub, err := _RepublicToken.contract.FilterLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return &RepublicTokenPauseIterator{contract: _RepublicToken.contract, event: "Pause", logs: logs, sub: sub}, nil
}

// WatchPause is a free log subscription operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
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
// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
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
// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
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
// Solidity: event Unpause()
func (_RepublicToken *RepublicTokenFilterer) FilterUnpause(opts *bind.FilterOpts) (*RepublicTokenUnpauseIterator, error) {

	logs, sub, err := _RepublicToken.contract.FilterLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return &RepublicTokenUnpauseIterator{contract: _RepublicToken.contract, event: "Unpause", logs: logs, sub: sub}, nil
}

// WatchUnpause is a free log subscription operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
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

// RewardGatewayABI is the input ABI used to generate the binding from.
const RewardGatewayABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"rewardVaults\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_tokenAddress\",\"type\":\"address\"},{\"name\":\"_rewardVaultAddress\",\"type\":\"address\"}],\"name\":\"updateRewardVault\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"

// RewardGatewayBin is the compiled bytecode used for deploying new contracts.
const RewardGatewayBin = `0x608060405260008054600160a060020a033316600160a060020a0319909116179055610246806100306000396000f3006080604052600436106100615763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166329228ba781146100665780638da5cb5b146100a357806397416cb4146100b8578063f2fde38b146100e1575b600080fd5b34801561007257600080fd5b50610087600160a060020a0360043516610102565b60408051600160a060020a039092168252519081900360200190f35b3480156100af57600080fd5b5061008761011d565b3480156100c457600080fd5b506100df600160a060020a036004358116906024351661012c565b005b3480156100ed57600080fd5b506100df600160a060020a0360043516610182565b600160205260009081526040902054600160a060020a031681565b600054600160a060020a031681565b60005433600160a060020a0390811691161461014757600080fd5b600160a060020a039182166000908152600160205260409020805473ffffffffffffffffffffffffffffffffffffffff191691909216179055565b60005433600160a060020a0390811691161461019d57600080fd5b600160a060020a03811615156101b257600080fd5b60008054604051600160a060020a03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a36000805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03929092169190911790555600a165627a7a72305820435780fa1fe3c55f12199c635bfab8b0775c51db111430585c1a0a584c3e5d190029`

// DeployRewardGateway deploys a new Ethereum contract, binding an instance of RewardGateway to it.
func DeployRewardGateway(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *RewardGateway, error) {
	parsed, err := abi.JSON(strings.NewReader(RewardGatewayABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(RewardGatewayBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &RewardGateway{RewardGatewayCaller: RewardGatewayCaller{contract: contract}, RewardGatewayTransactor: RewardGatewayTransactor{contract: contract}, RewardGatewayFilterer: RewardGatewayFilterer{contract: contract}}, nil
}

// RewardGateway is an auto generated Go binding around an Ethereum contract.
type RewardGateway struct {
	RewardGatewayCaller     // Read-only binding to the contract
	RewardGatewayTransactor // Write-only binding to the contract
	RewardGatewayFilterer   // Log filterer for contract events
}

// RewardGatewayCaller is an auto generated read-only Go binding around an Ethereum contract.
type RewardGatewayCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RewardGatewayTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RewardGatewayTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RewardGatewayFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RewardGatewayFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RewardGatewaySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RewardGatewaySession struct {
	Contract     *RewardGateway    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RewardGatewayCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RewardGatewayCallerSession struct {
	Contract *RewardGatewayCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// RewardGatewayTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RewardGatewayTransactorSession struct {
	Contract     *RewardGatewayTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// RewardGatewayRaw is an auto generated low-level Go binding around an Ethereum contract.
type RewardGatewayRaw struct {
	Contract *RewardGateway // Generic contract binding to access the raw methods on
}

// RewardGatewayCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RewardGatewayCallerRaw struct {
	Contract *RewardGatewayCaller // Generic read-only contract binding to access the raw methods on
}

// RewardGatewayTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RewardGatewayTransactorRaw struct {
	Contract *RewardGatewayTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRewardGateway creates a new instance of RewardGateway, bound to a specific deployed contract.
func NewRewardGateway(address common.Address, backend bind.ContractBackend) (*RewardGateway, error) {
	contract, err := bindRewardGateway(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RewardGateway{RewardGatewayCaller: RewardGatewayCaller{contract: contract}, RewardGatewayTransactor: RewardGatewayTransactor{contract: contract}, RewardGatewayFilterer: RewardGatewayFilterer{contract: contract}}, nil
}

// NewRewardGatewayCaller creates a new read-only instance of RewardGateway, bound to a specific deployed contract.
func NewRewardGatewayCaller(address common.Address, caller bind.ContractCaller) (*RewardGatewayCaller, error) {
	contract, err := bindRewardGateway(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RewardGatewayCaller{contract: contract}, nil
}

// NewRewardGatewayTransactor creates a new write-only instance of RewardGateway, bound to a specific deployed contract.
func NewRewardGatewayTransactor(address common.Address, transactor bind.ContractTransactor) (*RewardGatewayTransactor, error) {
	contract, err := bindRewardGateway(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RewardGatewayTransactor{contract: contract}, nil
}

// NewRewardGatewayFilterer creates a new log filterer instance of RewardGateway, bound to a specific deployed contract.
func NewRewardGatewayFilterer(address common.Address, filterer bind.ContractFilterer) (*RewardGatewayFilterer, error) {
	contract, err := bindRewardGateway(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RewardGatewayFilterer{contract: contract}, nil
}

// bindRewardGateway binds a generic wrapper to an already deployed contract.
func bindRewardGateway(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RewardGatewayABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RewardGateway *RewardGatewayRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _RewardGateway.Contract.RewardGatewayCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RewardGateway *RewardGatewayRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RewardGateway.Contract.RewardGatewayTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RewardGateway *RewardGatewayRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RewardGateway.Contract.RewardGatewayTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RewardGateway *RewardGatewayCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _RewardGateway.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RewardGateway *RewardGatewayTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RewardGateway.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RewardGateway *RewardGatewayTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RewardGateway.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_RewardGateway *RewardGatewayCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _RewardGateway.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_RewardGateway *RewardGatewaySession) Owner() (common.Address, error) {
	return _RewardGateway.Contract.Owner(&_RewardGateway.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_RewardGateway *RewardGatewayCallerSession) Owner() (common.Address, error) {
	return _RewardGateway.Contract.Owner(&_RewardGateway.CallOpts)
}

// RewardVaults is a free data retrieval call binding the contract method 0x29228ba7.
//
// Solidity: function rewardVaults( address) constant returns(address)
func (_RewardGateway *RewardGatewayCaller) RewardVaults(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _RewardGateway.contract.Call(opts, out, "rewardVaults", arg0)
	return *ret0, err
}

// RewardVaults is a free data retrieval call binding the contract method 0x29228ba7.
//
// Solidity: function rewardVaults( address) constant returns(address)
func (_RewardGateway *RewardGatewaySession) RewardVaults(arg0 common.Address) (common.Address, error) {
	return _RewardGateway.Contract.RewardVaults(&_RewardGateway.CallOpts, arg0)
}

// RewardVaults is a free data retrieval call binding the contract method 0x29228ba7.
//
// Solidity: function rewardVaults( address) constant returns(address)
func (_RewardGateway *RewardGatewayCallerSession) RewardVaults(arg0 common.Address) (common.Address, error) {
	return _RewardGateway.Contract.RewardVaults(&_RewardGateway.CallOpts, arg0)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_RewardGateway *RewardGatewayTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _RewardGateway.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_RewardGateway *RewardGatewaySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _RewardGateway.Contract.TransferOwnership(&_RewardGateway.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_RewardGateway *RewardGatewayTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _RewardGateway.Contract.TransferOwnership(&_RewardGateway.TransactOpts, newOwner)
}

// UpdateRewardVault is a paid mutator transaction binding the contract method 0x97416cb4.
//
// Solidity: function updateRewardVault(_tokenAddress address, _rewardVaultAddress address) returns()
func (_RewardGateway *RewardGatewayTransactor) UpdateRewardVault(opts *bind.TransactOpts, _tokenAddress common.Address, _rewardVaultAddress common.Address) (*types.Transaction, error) {
	return _RewardGateway.contract.Transact(opts, "updateRewardVault", _tokenAddress, _rewardVaultAddress)
}

// UpdateRewardVault is a paid mutator transaction binding the contract method 0x97416cb4.
//
// Solidity: function updateRewardVault(_tokenAddress address, _rewardVaultAddress address) returns()
func (_RewardGateway *RewardGatewaySession) UpdateRewardVault(_tokenAddress common.Address, _rewardVaultAddress common.Address) (*types.Transaction, error) {
	return _RewardGateway.Contract.UpdateRewardVault(&_RewardGateway.TransactOpts, _tokenAddress, _rewardVaultAddress)
}

// UpdateRewardVault is a paid mutator transaction binding the contract method 0x97416cb4.
//
// Solidity: function updateRewardVault(_tokenAddress address, _rewardVaultAddress address) returns()
func (_RewardGateway *RewardGatewayTransactorSession) UpdateRewardVault(_tokenAddress common.Address, _rewardVaultAddress common.Address) (*types.Transaction, error) {
	return _RewardGateway.Contract.UpdateRewardVault(&_RewardGateway.TransactOpts, _tokenAddress, _rewardVaultAddress)
}

// RewardGatewayOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the RewardGateway contract.
type RewardGatewayOwnershipTransferredIterator struct {
	Event *RewardGatewayOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *RewardGatewayOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardGatewayOwnershipTransferred)
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
		it.Event = new(RewardGatewayOwnershipTransferred)
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
func (it *RewardGatewayOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardGatewayOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardGatewayOwnershipTransferred represents a OwnershipTransferred event raised by the RewardGateway contract.
type RewardGatewayOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_RewardGateway *RewardGatewayFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*RewardGatewayOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _RewardGateway.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &RewardGatewayOwnershipTransferredIterator{contract: _RewardGateway.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_RewardGateway *RewardGatewayFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RewardGatewayOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _RewardGateway.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardGatewayOwnershipTransferred)
				if err := _RewardGateway.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// RewardVaultABI is the input ABI used to generate the binding from.
const RewardVaultABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_nonce\",\"type\":\"uint256\"}],\"name\":\"finalize\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_nonce\",\"type\":\"uint256\"}],\"name\":\"challengeIds\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_challenge\",\"type\":\"bytes\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_challengeID\",\"type\":\"bytes32\"},{\"name\":\"_proof\",\"type\":\"bytes\"},{\"name\":\"_rewardRoundNonce\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_nonce\",\"type\":\"uint256\"}],\"name\":\"isFinalizable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"ETHEREUM\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_nonce\",\"type\":\"uint256\"}],\"name\":\"rewardees\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_rewardeeCount\",\"type\":\"uint256\"},{\"name\":\"_challengeCount\",\"type\":\"uint256\"},{\"name\":\"_threshold\",\"type\":\"uint256\"},{\"name\":\"_dnrAddress\",\"type\":\"address\"},{\"name\":\"_tokenAddress\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// RewardVaultBin is the compiled bytecode used for deploying new contracts.
const RewardVaultBin = `0x608060405234801561001057600080fd5b5060405160a080610eee83398101604081815282516020808501518386015160608088015160809098015160e0880187526000885293870185905294860182905292959094929390929081018587028581151561006957fe5b0481526020808201869052600160a060020a03808616604080850191909152858216606094850181905285516000559285015160015584015160025591830151600355608083015160045560a083015160058054918416600160a060020a031992831617905560c0909301516006805491909316931692909217905573eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee1461011b57600a8054600160a060020a031916600160a060020a0383161790555b5050505050610dbf8061012f6000396000f3006080604052600436106100695763ffffffff60e060020a60003504166305261aea811461006e57806323f6bfb21461008857806393405540146100f0578063cfc858ff1461013e578063e533cdf51461019e578063f7cdf47c146101ca578063f87770f5146101fb575b600080fd5b34801561007a57600080fd5b50610086600435610213565b005b34801561009457600080fd5b506100a0600435610221565b60408051602080825283518183015283519192839290830191858101910280838360005b838110156100dc5781810151838201526020016100c4565b505050509050019250505060405180910390f35b6040805160206004803580820135601f810184900484028501840190955284845261008694369492936024939284019190819084018382808284375094975050933594506102a59350505050565b34801561014a57600080fd5b5060408051602060046024803582810135601f810185900485028601850190965285855261008695833595369560449491939091019190819084018382808284375094975050933594506104449350505050565b3480156101aa57600080fd5b506101b6600435610570565b604080519115158252519081900360200190f35b3480156101d657600080fd5b506101df6105a9565b60408051600160a060020a039092168252519081900360200190f35b34801561020757600080fd5b506100a06004356105c1565b61021e60008261064d565b50565b60008181526009602052604090205460609060ff16151561024157600080fd5b6000828152600960209081526040918290206004018054835181840281018401909452808452909183018282801561029957602002820191906000526020600020905b81548152600190910190602001808311610284575b50505050509050919050565b600654600160a060020a031673eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee14156102dd573481146102d857600080fd5b610434565b600a54604080517fdd62ed3e000000000000000000000000000000000000000000000000000000008152600160a060020a033381166004830152308116602483015291518493929092169163dd62ed3e916044808201926020929091908290030181600087803b15801561035057600080fd5b505af1158015610364573d6000803e3d6000fd5b505050506040513d602081101561037a57600080fd5b50511461038657600080fd5b600a54604080517f23b872dd000000000000000000000000000000000000000000000000000000008152600160a060020a033381166004830152308116602483015260448201859052915191909216916323b872dd9160648083019260209291908290030181600087803b1580156103fd57600080fd5b505af1158015610411573d6000803e3d6000fd5b505050506040513d602081101561042757600080fd5b5051151561043457600080fd5b61044060008383610b41565b5050565b610452600084848433610c22565b151561045d57600080fd5b600654600160a060020a031673eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee14156104c157600354604051600160a060020a0333169180156108fc02916000818181858888f193505050501580156104bb573d6000803e3d6000fd5b5061056b565b600a54600354604080517fa9059cbb000000000000000000000000000000000000000000000000000000008152600160a060020a03338116600483015260248201939093529051919092169163a9059cbb9160448083019260209291908290030181600087803b15801561053457600080fd5b505af1158015610548573d6000803e3d6000fd5b505050506040513d602081101561055e57600080fd5b5051151561056b57600080fd5b505050565b60045460008281526009602052604081206002015490911480156105a3575060008281526009602052604090205460ff16155b92915050565b73eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee81565b60008181526009602052604090205460609060ff1615156105e157600080fd5b6000828152600960209081526040918290206005018054835181840281018401909452808452909183018282801561029957602002820191906000526020600020905b8154600160a060020a031681526001909101906020018083116106245750505050509050919050565b600582015460048301546000838152600985016020526040812060020154600160a060020a039093169290918291606091839182911480156106a05750600087815260098901602052604090205460ff16155b15156106ab57600080fd5b6000878152600989016020526040902060038101546001909101548115156106cf57fe5b069450600093505b8760020154841015610812575b6000878152600789016020908152604080832088845290915290205460ff1615610713578460010194506106e4565b600087815260098901602052604090206003810180546004909201918790811061073957fe5b90600052602060002001604051808280546001816001161561010002031660029004801561079e5780601f1061077c57610100808354040283529182019161079e565b820191906000526020600020905b81548152906001019060200180831161078a575b5050604080519182900390912084546001818101875560009687526020808820909201929092558c865260078e0181528286208b87528152828620805460ff1916831790558c865260098e01905293206003810154930154880191505081151561080457fe5b0694506001909301926106d7565b85600160a060020a031663879ae0846040518163ffffffff1660e060020a028152600401600060405180830381600087803b15801561085057600080fd5b505af1158015610864573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052602081101561088d57600080fd5b8101908080516401000000008111156108a557600080fd5b820160208101848111156108b857600080fd5b81518560208202830111640100000000821117156108d557600080fd5b5050805160008c815260098e01602052604090206001015491975093509150508115156108fe57fe5b069450600093505b8760010154841015610b1d575b6040805188815260208082018890528251918290038301909120600090815260088b01909152205460ff161561094e57846001019450610913565b85600160a060020a031663e487eb58848781518110151561096b57fe5b60209081029091018101516040805160e060020a63ffffffff86160281526bffffffffffffffffffffffff19909216600483015251602480830193928290030181600087803b1580156109bd57600080fd5b505af11580156109d1573d6000803e3d6000fd5b505050506040513d60208110156109e757600080fd5b5051600088815260098a01602090815260408220600501805460018101825590835290822001805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03841617905590925090505b8760020154811015610abb5760008781526009890160208181526040808420600160a060020a03871685526006810183529084208b8552929091526004018054600193919085908110610a8857fe5b60009182526020808320919091015483528201929092526040019020805460ff1916911515919091179055600101610a39565b6040805188815260208082018890528251918290038301909120600090815260088b018252828120805460ff1916600190811790915586518b835260098d01909352929020909101548601811515610b0f57fe5b069450600190930192610906565b5050506000938452505050600990910160205260409020805460ff19166001179055565b82546000908152600984016020908152604082206003018054600181018083559184529282902085519193610b7c9391019190860190610cf8565b5050600483015483546000908152600985016020526040902060020154820110610c025760048301805484546000818152600987016020526040808220600290810154600194850184528284209088019590950394810194909455935487548252848220909301929092558554825291902060001943014090820155835401835561056b565b825460009081526009840160205260409020600201805482019055505050565b600082815260098601602052604081205460ff161515610c4157600080fd5b60008381526009870160209081526040808320600160a060020a0386168452600601825280832088845290915290205460ff161515610c8257506000610ce7565b610c8c8585610cf0565b1515610c9a57506000610ce7565b5060038501546000838152600987016020908152604080832060028101805495909503909455600160a060020a038516835260069093018152828220878352905220805460ff1916905560015b95945050505050565b600192915050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10610d3957805160ff1916838001178555610d66565b82800160010185558215610d66579182015b82811115610d66578251825591602001919060010190610d4b565b50610d72929150610d76565b5090565b610d9091905b80821115610d725760008155600101610d7c565b905600a165627a7a7230582086e77c13f629e0532e16e4cf17298f597df6aefb4b811c1b618d99e8e9f57e090029`

// DeployRewardVault deploys a new Ethereum contract, binding an instance of RewardVault to it.
func DeployRewardVault(auth *bind.TransactOpts, backend bind.ContractBackend, _rewardeeCount *big.Int, _challengeCount *big.Int, _threshold *big.Int, _dnrAddress common.Address, _tokenAddress common.Address) (common.Address, *types.Transaction, *RewardVault, error) {
	parsed, err := abi.JSON(strings.NewReader(RewardVaultABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(RewardVaultBin), backend, _rewardeeCount, _challengeCount, _threshold, _dnrAddress, _tokenAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &RewardVault{RewardVaultCaller: RewardVaultCaller{contract: contract}, RewardVaultTransactor: RewardVaultTransactor{contract: contract}, RewardVaultFilterer: RewardVaultFilterer{contract: contract}}, nil
}

// RewardVault is an auto generated Go binding around an Ethereum contract.
type RewardVault struct {
	RewardVaultCaller     // Read-only binding to the contract
	RewardVaultTransactor // Write-only binding to the contract
	RewardVaultFilterer   // Log filterer for contract events
}

// RewardVaultCaller is an auto generated read-only Go binding around an Ethereum contract.
type RewardVaultCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RewardVaultTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RewardVaultTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RewardVaultFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RewardVaultFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RewardVaultSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RewardVaultSession struct {
	Contract     *RewardVault      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RewardVaultCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RewardVaultCallerSession struct {
	Contract *RewardVaultCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// RewardVaultTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RewardVaultTransactorSession struct {
	Contract     *RewardVaultTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// RewardVaultRaw is an auto generated low-level Go binding around an Ethereum contract.
type RewardVaultRaw struct {
	Contract *RewardVault // Generic contract binding to access the raw methods on
}

// RewardVaultCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RewardVaultCallerRaw struct {
	Contract *RewardVaultCaller // Generic read-only contract binding to access the raw methods on
}

// RewardVaultTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RewardVaultTransactorRaw struct {
	Contract *RewardVaultTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRewardVault creates a new instance of RewardVault, bound to a specific deployed contract.
func NewRewardVault(address common.Address, backend bind.ContractBackend) (*RewardVault, error) {
	contract, err := bindRewardVault(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RewardVault{RewardVaultCaller: RewardVaultCaller{contract: contract}, RewardVaultTransactor: RewardVaultTransactor{contract: contract}, RewardVaultFilterer: RewardVaultFilterer{contract: contract}}, nil
}

// NewRewardVaultCaller creates a new read-only instance of RewardVault, bound to a specific deployed contract.
func NewRewardVaultCaller(address common.Address, caller bind.ContractCaller) (*RewardVaultCaller, error) {
	contract, err := bindRewardVault(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RewardVaultCaller{contract: contract}, nil
}

// NewRewardVaultTransactor creates a new write-only instance of RewardVault, bound to a specific deployed contract.
func NewRewardVaultTransactor(address common.Address, transactor bind.ContractTransactor) (*RewardVaultTransactor, error) {
	contract, err := bindRewardVault(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RewardVaultTransactor{contract: contract}, nil
}

// NewRewardVaultFilterer creates a new log filterer instance of RewardVault, bound to a specific deployed contract.
func NewRewardVaultFilterer(address common.Address, filterer bind.ContractFilterer) (*RewardVaultFilterer, error) {
	contract, err := bindRewardVault(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RewardVaultFilterer{contract: contract}, nil
}

// bindRewardVault binds a generic wrapper to an already deployed contract.
func bindRewardVault(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RewardVaultABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RewardVault *RewardVaultRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _RewardVault.Contract.RewardVaultCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RewardVault *RewardVaultRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RewardVault.Contract.RewardVaultTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RewardVault *RewardVaultRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RewardVault.Contract.RewardVaultTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RewardVault *RewardVaultCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _RewardVault.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RewardVault *RewardVaultTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RewardVault.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RewardVault *RewardVaultTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RewardVault.Contract.contract.Transact(opts, method, params...)
}

// ETHEREUM is a free data retrieval call binding the contract method 0xf7cdf47c.
//
// Solidity: function ETHEREUM() constant returns(address)
func (_RewardVault *RewardVaultCaller) ETHEREUM(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _RewardVault.contract.Call(opts, out, "ETHEREUM")
	return *ret0, err
}

// ETHEREUM is a free data retrieval call binding the contract method 0xf7cdf47c.
//
// Solidity: function ETHEREUM() constant returns(address)
func (_RewardVault *RewardVaultSession) ETHEREUM() (common.Address, error) {
	return _RewardVault.Contract.ETHEREUM(&_RewardVault.CallOpts)
}

// ETHEREUM is a free data retrieval call binding the contract method 0xf7cdf47c.
//
// Solidity: function ETHEREUM() constant returns(address)
func (_RewardVault *RewardVaultCallerSession) ETHEREUM() (common.Address, error) {
	return _RewardVault.Contract.ETHEREUM(&_RewardVault.CallOpts)
}

// ChallengeIds is a free data retrieval call binding the contract method 0x23f6bfb2.
//
// Solidity: function challengeIds(_nonce uint256) constant returns(bytes32[])
func (_RewardVault *RewardVaultCaller) ChallengeIds(opts *bind.CallOpts, _nonce *big.Int) ([][32]byte, error) {
	var (
		ret0 = new([][32]byte)
	)
	out := ret0
	err := _RewardVault.contract.Call(opts, out, "challengeIds", _nonce)
	return *ret0, err
}

// ChallengeIds is a free data retrieval call binding the contract method 0x23f6bfb2.
//
// Solidity: function challengeIds(_nonce uint256) constant returns(bytes32[])
func (_RewardVault *RewardVaultSession) ChallengeIds(_nonce *big.Int) ([][32]byte, error) {
	return _RewardVault.Contract.ChallengeIds(&_RewardVault.CallOpts, _nonce)
}

// ChallengeIds is a free data retrieval call binding the contract method 0x23f6bfb2.
//
// Solidity: function challengeIds(_nonce uint256) constant returns(bytes32[])
func (_RewardVault *RewardVaultCallerSession) ChallengeIds(_nonce *big.Int) ([][32]byte, error) {
	return _RewardVault.Contract.ChallengeIds(&_RewardVault.CallOpts, _nonce)
}

// IsFinalizable is a free data retrieval call binding the contract method 0xe533cdf5.
//
// Solidity: function isFinalizable(_nonce uint256) constant returns(bool)
func (_RewardVault *RewardVaultCaller) IsFinalizable(opts *bind.CallOpts, _nonce *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _RewardVault.contract.Call(opts, out, "isFinalizable", _nonce)
	return *ret0, err
}

// IsFinalizable is a free data retrieval call binding the contract method 0xe533cdf5.
//
// Solidity: function isFinalizable(_nonce uint256) constant returns(bool)
func (_RewardVault *RewardVaultSession) IsFinalizable(_nonce *big.Int) (bool, error) {
	return _RewardVault.Contract.IsFinalizable(&_RewardVault.CallOpts, _nonce)
}

// IsFinalizable is a free data retrieval call binding the contract method 0xe533cdf5.
//
// Solidity: function isFinalizable(_nonce uint256) constant returns(bool)
func (_RewardVault *RewardVaultCallerSession) IsFinalizable(_nonce *big.Int) (bool, error) {
	return _RewardVault.Contract.IsFinalizable(&_RewardVault.CallOpts, _nonce)
}

// Rewardees is a free data retrieval call binding the contract method 0xf87770f5.
//
// Solidity: function rewardees(_nonce uint256) constant returns(address[])
func (_RewardVault *RewardVaultCaller) Rewardees(opts *bind.CallOpts, _nonce *big.Int) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _RewardVault.contract.Call(opts, out, "rewardees", _nonce)
	return *ret0, err
}

// Rewardees is a free data retrieval call binding the contract method 0xf87770f5.
//
// Solidity: function rewardees(_nonce uint256) constant returns(address[])
func (_RewardVault *RewardVaultSession) Rewardees(_nonce *big.Int) ([]common.Address, error) {
	return _RewardVault.Contract.Rewardees(&_RewardVault.CallOpts, _nonce)
}

// Rewardees is a free data retrieval call binding the contract method 0xf87770f5.
//
// Solidity: function rewardees(_nonce uint256) constant returns(address[])
func (_RewardVault *RewardVaultCallerSession) Rewardees(_nonce *big.Int) ([]common.Address, error) {
	return _RewardVault.Contract.Rewardees(&_RewardVault.CallOpts, _nonce)
}

// Deposit is a paid mutator transaction binding the contract method 0x93405540.
//
// Solidity: function deposit(_challenge bytes, _value uint256) returns()
func (_RewardVault *RewardVaultTransactor) Deposit(opts *bind.TransactOpts, _challenge []byte, _value *big.Int) (*types.Transaction, error) {
	return _RewardVault.contract.Transact(opts, "deposit", _challenge, _value)
}

// Deposit is a paid mutator transaction binding the contract method 0x93405540.
//
// Solidity: function deposit(_challenge bytes, _value uint256) returns()
func (_RewardVault *RewardVaultSession) Deposit(_challenge []byte, _value *big.Int) (*types.Transaction, error) {
	return _RewardVault.Contract.Deposit(&_RewardVault.TransactOpts, _challenge, _value)
}

// Deposit is a paid mutator transaction binding the contract method 0x93405540.
//
// Solidity: function deposit(_challenge bytes, _value uint256) returns()
func (_RewardVault *RewardVaultTransactorSession) Deposit(_challenge []byte, _value *big.Int) (*types.Transaction, error) {
	return _RewardVault.Contract.Deposit(&_RewardVault.TransactOpts, _challenge, _value)
}

// Finalize is a paid mutator transaction binding the contract method 0x05261aea.
//
// Solidity: function finalize(_nonce uint256) returns()
func (_RewardVault *RewardVaultTransactor) Finalize(opts *bind.TransactOpts, _nonce *big.Int) (*types.Transaction, error) {
	return _RewardVault.contract.Transact(opts, "finalize", _nonce)
}

// Finalize is a paid mutator transaction binding the contract method 0x05261aea.
//
// Solidity: function finalize(_nonce uint256) returns()
func (_RewardVault *RewardVaultSession) Finalize(_nonce *big.Int) (*types.Transaction, error) {
	return _RewardVault.Contract.Finalize(&_RewardVault.TransactOpts, _nonce)
}

// Finalize is a paid mutator transaction binding the contract method 0x05261aea.
//
// Solidity: function finalize(_nonce uint256) returns()
func (_RewardVault *RewardVaultTransactorSession) Finalize(_nonce *big.Int) (*types.Transaction, error) {
	return _RewardVault.Contract.Finalize(&_RewardVault.TransactOpts, _nonce)
}

// Withdraw is a paid mutator transaction binding the contract method 0xcfc858ff.
//
// Solidity: function withdraw(_challengeID bytes32, _proof bytes, _rewardRoundNonce uint256) returns()
func (_RewardVault *RewardVaultTransactor) Withdraw(opts *bind.TransactOpts, _challengeID [32]byte, _proof []byte, _rewardRoundNonce *big.Int) (*types.Transaction, error) {
	return _RewardVault.contract.Transact(opts, "withdraw", _challengeID, _proof, _rewardRoundNonce)
}

// Withdraw is a paid mutator transaction binding the contract method 0xcfc858ff.
//
// Solidity: function withdraw(_challengeID bytes32, _proof bytes, _rewardRoundNonce uint256) returns()
func (_RewardVault *RewardVaultSession) Withdraw(_challengeID [32]byte, _proof []byte, _rewardRoundNonce *big.Int) (*types.Transaction, error) {
	return _RewardVault.Contract.Withdraw(&_RewardVault.TransactOpts, _challengeID, _proof, _rewardRoundNonce)
}

// Withdraw is a paid mutator transaction binding the contract method 0xcfc858ff.
//
// Solidity: function withdraw(_challengeID bytes32, _proof bytes, _rewardRoundNonce uint256) returns()
func (_RewardVault *RewardVaultTransactorSession) Withdraw(_challengeID [32]byte, _proof []byte, _rewardRoundNonce *big.Int) (*types.Transaction, error) {
	return _RewardVault.Contract.Withdraw(&_RewardVault.TransactOpts, _challengeID, _proof, _rewardRoundNonce)
}

// SafeMathABI is the input ABI used to generate the binding from.
const SafeMathABI = "[]"

// SafeMathBin is the compiled bytecode used for deploying new contracts.
const SafeMathBin = `0x604c602c600b82828239805160001a60731460008114601c57601e565bfe5b5030600052607381538281f30073000000000000000000000000000000000000000030146080604052600080fd00a165627a7a72305820d763408514ec2a0a99a793456fd496efe8a35962fcc44fed3c35edc35c33ad500029`

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

// StandardTokenABI is the input ABI used to generate the binding from.
const StandardTokenABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseApproval\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseApproval\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]"

// StandardTokenBin is the compiled bytecode used for deploying new contracts.
const StandardTokenBin = `0x608060405234801561001057600080fd5b506106ed806100206000396000f30060806040526004361061008d5763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663095ea7b3811461009257806318160ddd146100ca57806323b872dd146100f1578063661884631461011b57806370a082311461013f578063a9059cbb14610160578063d73dd62314610184578063dd62ed3e146101a8575b600080fd5b34801561009e57600080fd5b506100b6600160a060020a03600435166024356101cf565b604080519115158252519081900360200190f35b3480156100d657600080fd5b506100df610239565b60408051918252519081900360200190f35b3480156100fd57600080fd5b506100b6600160a060020a036004358116906024351660443561023f565b34801561012757600080fd5b506100b6600160a060020a03600435166024356103bf565b34801561014b57600080fd5b506100df600160a060020a03600435166104b8565b34801561016c57600080fd5b506100b6600160a060020a03600435166024356104d3565b34801561019057600080fd5b506100b6600160a060020a03600435166024356105cc565b3480156101b457600080fd5b506100df600160a060020a036004358116906024351661066e565b600160a060020a03338116600081815260026020908152604080832094871680845294825280832086905580518681529051929493927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925929181900390910190a350600192915050565b60015490565b6000600160a060020a038316151561025657600080fd5b600160a060020a03841660009081526020819052604090205482111561027b57600080fd5b600160a060020a03808516600090815260026020908152604080832033909416835292905220548211156102ae57600080fd5b600160a060020a0384166000908152602081905260409020546102d7908363ffffffff61069916565b600160a060020a03808616600090815260208190526040808220939093559085168152205461030c908363ffffffff6106ab16565b600160a060020a0380851660009081526020818152604080832094909455878316825260028152838220339093168252919091522054610352908363ffffffff61069916565b600160a060020a038086166000818152600260209081526040808320338616845282529182902094909455805186815290519287169391927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef929181900390910190a35060019392505050565b600160a060020a0333811660009081526002602090815260408083209386168352929052908120548083111561041c57600160a060020a033381166000908152600260209081526040808320938816835292905290812055610453565b61042c818463ffffffff61069916565b600160a060020a033381166000908152600260209081526040808320938916835292905220555b600160a060020a0333811660008181526002602090815260408083209489168084529482529182902054825190815291517f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259281900390910190a35060019392505050565b600160a060020a031660009081526020819052604090205490565b6000600160a060020a03831615156104ea57600080fd5b600160a060020a03331660009081526020819052604090205482111561050f57600080fd5b600160a060020a033316600090815260208190526040902054610538908363ffffffff61069916565b600160a060020a03338116600090815260208190526040808220939093559085168152205461056d908363ffffffff6106ab16565b600160a060020a03808516600081815260208181526040918290209490945580518681529051919333909316927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef92918290030190a350600192915050565b600160a060020a033381166000908152600260209081526040808320938616835292905290812054610604908363ffffffff6106ab16565b600160a060020a0333811660008181526002602090815260408083209489168084529482529182902085905581519485529051929391927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259281900390910190a350600192915050565b600160a060020a03918216600090815260026020908152604080832093909416825291909152205490565b6000828211156106a557fe5b50900390565b6000828201838110156106ba57fe5b93925050505600a165627a7a7230582021ed9693aad1354907887934ed91cfd6b0ca4671ea74217ef23cabcea6d6b3e70029`

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
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
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
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_StandardToken *StandardTokenSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _StandardToken.Contract.BalanceOf(&_StandardToken.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
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
// Solidity: event Approval(owner indexed address, spender indexed address, value uint256)
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
// Solidity: event Approval(owner indexed address, spender indexed address, value uint256)
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
// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
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
// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
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
