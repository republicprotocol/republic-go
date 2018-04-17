package arc

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/republicprotocol/republic-go/ethereum/bindings"
	"github.com/republicprotocol/republic-go/ethereum/client"
)

type EthereumData struct {
	contractAddress [20]byte
	value *big.Int
}

type EthereumArc struct {
	context context.Context
	client  *connection.Client
	auth    *bind.TransactOpts
	binding *bindings.Arc
	ethereumData EthereumData 
}

func NewEthereumArc(context context.Context, Client Client, auth *bind.TransactOpts) Arc {
	return &EthereumDarkNodeRegistrar{
		context: context,
		client:  Client,
		auth:    auth,
	}
}

func (arc *EthereumArc) Initiate(hash [32]byte, to []byte, value *big.Int, expiry int64) error {
	arc.ethereumData.contractAddress, tx, arc.binding, err := bindings.DeployArc(auth, bind.ContractBackend(arc.client), hash, common.HexToAddress("0x1"), value, expiry, to)
	arc.ethereumData.value = value;

	if err == nil {
		_, err = connection.PatchedWaitDeployed(arc.context, *arc.client, tx)
		return err;
	}

	return err;
}

func (arc *EthereumArc) Redeem(secret []byte) error {
	tx, err := arc.binding.Redeem(arc.auth, secret);
	if err == nil {
		_, err = connection.PatchedWaitMined(arc.context, *arc.client, tx);
	}
	return err;
}

func (arc *EthereumArc) Refund() error {
	tx, err := arc.binding.Refund(arc.auth, common.HexToAddress("0x1"), arc.ethereumData.value);
}

func (arc *EthereumArc) Audit() (hash, to, from []byte, value *big.Int, expiry *big.Int, err error) {
	add1, num1, addr2, num2, err := arc.binding.Audit(*bind.CallOpts{})
	// return 
}

func (arc *EthereumArc) AuditSecret() (secret []byte, err error) {
	return arc.binding.AuditSecret(bind.CallOpts{});
}
