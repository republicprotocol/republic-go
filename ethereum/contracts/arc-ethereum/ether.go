package arc_ethereum

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/republicprotocol/republic-go/ethereum/bindings"
	"github.com/republicprotocol/republic-go/ethereum/client"
)

type EthereumArcData struct {
	ContractAddress common.Address
	Value           *big.Int
}

type EthereumArc struct {
	context context.Context
	client  client.Connection
	auth    *bind.TransactOpts
	// binding      *interop.Arc
	EthereumArcData EthereumArcData

	binding *bindings.Arc
	swapID  [32]byte
}

// NewEthereumArc returns a new EthereumArc instance
func NewEthereumArc(context context.Context, client client.Connection, auth *bind.TransactOpts, contractAddress common.Address, swapID [32]byte) (*EthereumArc, error) {
	contract, err := bindings.NewArc(contractAddress, bind.ContractBackend(client.Client))
	if err != nil {
		return nil, err
	}

	return &EthereumArc{
		context:         context,
		client:          client,
		auth:            auth,
		EthereumArcData: EthereumArcData{contractAddress, big.NewInt(0)},
		binding:         contract,
		swapID:          swapID,
	}, nil
}

// Initiate a new Arc swap by deploying an Arc contract
func (arc *EthereumArc) Initiate(hash [32]byte, to []byte, value *big.Int, expiry int64) error {
	contractAddress, tx, binding, err := bindings.DeployArc(arc.auth, arc.client.Client, hash, common.HexToAddress("0x1"), value, big.NewInt(expiry), common.BytesToAddress(to))
	if err != nil {
		return err
	}

	if err := arc.client.TransferEth(arc.context, arc.auth, contractAddress, value); err != nil {
		return err
	}

	arc.EthereumArcData = EthereumArcData{
		contractAddress, value,
	}
	arc.binding = binding

	_, err = arc.client.PatchedWaitDeployed(arc.context, tx)
	return err
}

func (arc *EthereumArc) Redeem(secret []byte) error {
	tx, err := arc.binding.Redeem(arc.auth, secret)
	if err == nil {
		_, err = arc.client.PatchedWaitMined(arc.context, tx)
	}
	return err
}

func (arc *EthereumArc) Refund() error {
	tx, err := arc.binding.Refund(arc.auth, common.HexToAddress("0x1"), arc.EthereumArcData.Value)
	if err == nil {
		_, err = arc.client.PatchedWaitMined(arc.context, tx)
	}
	return err
}

func (arc *EthereumArc) Audit() (hash [32]byte, to, from []byte, value *big.Int, expiry *big.Int, err error) {
	hash, _, toAddr, value, expiry, err := arc.binding.Audit(&bind.CallOpts{})
	// arc.client.Client.CodeAt(arc.context, arc.ethereumData.contractAddress, nil)
	// TODO: Audit values, bytecode
	return hash, arc.EthereumArcData.ContractAddress.Bytes(), toAddr.Bytes(), value, expiry, err
}

func (arc *EthereumArc) AuditSecret() ([]byte, error) {
	secret, err := arc.binding.AuditSecret(&bind.CallOpts{})
	// if err != nil {
	// 	return [32]byte{}, err
	// }
	// var secret32 [32]byte
	// copy(secret, secret32[:])
	// return secret32, nil
	return secret, err
}
