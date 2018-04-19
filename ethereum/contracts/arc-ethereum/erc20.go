package arc_ethereum

// import (
// 	"context"
// 	"math/big"

// 	"github.com/ethereum/go-ethereum/accounts/abi/bind"
// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/republicprotocol/republic-go/ethereum/bindings"
// 	"github.com/republicprotocol/republic-go/ethereum/client"
// )

// type ERC20ArcData struct {
// 	contractAddress common.Address
// 	erc20Address    common.Address
// 	value           *big.Int
// }

// type ERC20Arc struct {
// 	context   context.Context
// 	client    client.Connection
// 	auth      *bind.TransactOpts
// 	erc20Data ERC20ArcData

// 	binding      *bindings.Arc
// 	erc20Binding *bindings.ERC20
// 	swapID       [32]byte
// }

// func NewERC20Arc(context context.Context, client client.Connection, auth *bind.TransactOpts, address common.Address, erc20Address common.Address, swapID [32]byte) (*ERC20Arc, error) {
// 	binding, err := bindings.NewArc(address, bind.ContractBackend(client.Client))
// 	if err != nil {
// 		return nil, err
// 	}

// 	erc20Binding, err := bindings.NewERC20(erc20Address, bind.ContractBackend(client.Client))
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &ERC20Arc{
// 		context,
// 		client,
// 		auth,
// 		erc20Data: ERC20ArcData{
// 			contractAddress, erc20Address,
// 		},
// 		binding,
// 		erc20Binding,
// 		swapID,
// 	}, nil
// }

// // Initiate a new Arc swap by deploying an Arc contract
// func (arc *ERC20Arc) Initiate(hash [32]byte, to []byte, value *big.Int, expiry int64) error {
// 	arc.auth.Value = value
// 	contractAddress, tx, binding, err := bindings.DeployArc(arc.auth, bind.ContractBackend(arc.client.Client), hash, common.HexToAddress("0x1"), value, big.NewInt(expiry), common.BytesToAddress(to))
// 	arc.auth.Value = big.NewInt(0)
// 	arc.ethereumData = EthereumArcData{
// 		contractAddress, value,
// 	}
// 	arc.binding = binding

// 	if err == nil {
// 		_, err = arc.client.PatchedWaitDeployed(arc.context, tx)
// 		return err
// 	}

// 	return err
// }

// func (arc *ERC20Arc) Redeem(secret [32]byte) error {
// 	tx, err := arc.binding.Redeem(arc.auth, secret[:])
// 	if err == nil {
// 		_, err = arc.client.PatchedWaitMined(arc.context, tx)
// 	}
// 	return err
// }

// func (arc *ERC20Arc) Refund() error {
// 	tx, err := arc.binding.Refund(arc.auth, common.HexToAddress("0x1"), arc.ethereumData.value)
// 	if err == nil {
// 		_, err = arc.client.PatchedWaitMined(arc.context, tx)
// 	}
// 	return err
// }

// func (arc *ERC20Arc) Audit() (hash [32]byte, to, from []byte, value *big.Int, expiry *big.Int, err error) {
// 	hash, token, toAddr, value, expiry, err := arc.binding.Audit(&bind.CallOpts{})
// 	// arc.client.Client.CodeAt(arc.context, arc.ethereumData.contractAddress, nil)
// 	// TODO: Audit values, bytecode
// 	return hash, arc.ethereumData.contractAddress.Bytes(), toAddr.Bytes(), value, expiry, err
// }

// func (arc *ERC20Arc) AuditSecret() (secret []byte, err error) {
// 	return arc.binding.AuditSecret(&bind.CallOpts{})
// }

// /*

//  */

// // Initiate starts or reciprocates an atomic swap
// func (arc *Erc20Arc) Initiate(hash [32]byte, to, from []byte, value *big.Int, expiry int64) error {
// 	toAddress := common.BytesToAddress(to)

// 	// Approve ERC20 to atomic-swap
// 	tx, err := arc.erc20.Approve(arc.auth, arc.bindingAddress, value)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = arc.client.PatchedWaitMined(arc.context, tx)
// 	if err != nil {
// 		return err
// 	}

// 	// Call atomic-swap contract
// 	tx, err = arc.binding.Initiate(arc.auth, arc.swapID, value, arc.erc20Address, toAddress, hash, big.NewInt(expiry))

// 	if err != nil {
// 		return err
// 	}
// 	_, err = arc.client.PatchedWaitMined(arc.context, tx)
// 	return err
// }

// // Audit returns details about an atomic swap
// func (arc *Erc20Arc) Audit() (hash [32]byte, to, from []byte, value *big.Int, expiry int64, err error) {
// 	ret, err := arc.binding.Check(&bind.CallOpts{}, arc.swapID)
// 	return ret.SecretLock,
// 		ret.WithdrawTrader.Bytes(),
// 		ret.Erc20ContractAddress.Bytes(),
// 		ret.Erc20Value,
// 		ret.Timelock.Int64(),
// 		err
// }

// // AuditSecret returns the secret of an atomic swap if it's available
// func (arc *Erc20Arc) AuditSecret() (secret [32]byte, err error) {
// 	ret, err := arc.binding.CheckSecretKey(&bind.CallOpts{}, arc.swapID)
// 	if err != nil {
// 		return [32]byte{},

// 			err
// 	}
// 	copy(ret, secret[:])
// 	return secret, nil
// }

// // Redeem ...
// func (arc *Erc20Arc) Redeem(secret [32]byte) error {
// 	tx, err := arc.binding.Close(arc.auth, arc.swapID, secret[:])
// 	if err != nil {
// 		return err
// 	}
// 	_, err = arc.client.PatchedWaitMined(arc.context, tx)
// 	return err
// }

// // Refund will return the funds of an atomic swap, if the expiry period has passed
// func (arc *Erc20Arc) Refund() error {
// 	tx, err := arc.binding.Expire(arc.auth, arc.swapID)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = arc.client.PatchedWaitMined(arc.context, tx)
// 	return err
// }

// // // GetData returns the data required for another party to participate in an atomic swap
// // func (contract *Erc20Arc) GetData() []byte {
// // 	return contract.swapID[:]
// // }
