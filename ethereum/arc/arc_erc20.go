package arc

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/republicprotocol/republic-go/contracts/bindings"
	"github.com/republicprotocol/republic-go/contracts/connection"
)

type Erc20Data struct {
	contractAddress [20]byte
	contract        []byte
}

type Erc20Arc struct {
	context                  context.Context
	client                   *connection.Client
	auth                     *bind.TransactOpts
	binding                  *bindings.DarkNodeRegistrar
	tokenBinding             *bindings.ERC20
	darkNodeRegistrarAddress common.Address
}

func NewErc20Arc(context context.Context, clientDetails *connection.ClientDetails, auth *bind.TransactOpts) (Arc, error) {
	contract, err := bindings.NewDarkNodeRegistrar(clientDetails.DNRAddress, bind.ContractBackend(clientDetails.Client))
	if err != nil {
		return nil, err
	}
	renContract, err := bindings.NewERC20(clientDetails.RenAddress, bind.ContractBackend(clientDetails.Client))
	if err != nil {
		return nil, err
	}
	return &Erc20DarkNodeRegistrar{
		context:                  context,
		client:                   &clientDetails.Client,
		auth:                     auth,
		binding:                  contract,
		tokenBinding:             renContract,
		darkNodeRegistrarAddress: clientDetails.DNRAddress,
	}, nil
}

func (arc *Erc20Arc) Initiate(hash, to, from []byte, value *big.Int, expiry int64) error {

}

func (arc *Erc20Arc) Redeem(secret []byte) error {

}

func (arc *Erc20Arc) Refund() error {

}

func (arc *Erc20Arc) Audit() (hash, to, from []byte, value *big.Int, expiry int64, err error) {

}

func (arc *Erc20Arc) AuditSecret() (secret []byte, err error) {

}
