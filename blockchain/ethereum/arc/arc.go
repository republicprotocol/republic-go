package arc

import (
	"context"
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/republicprotocol/republic-go/blockchain/arc"
	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/bindings"
)

type EthereumArcData struct {
	ContractAddress common.Address `json:"contract_address"`
	Value           *big.Int       `json:"value"`
}

type Arc struct {
	context context.Context
	conn    ethereum.Conn
	auth    *bind.TransactOpts
	arcData EthereumArcData

	binding *bindings.Arc
	swapID  [32]byte
}

// NewArc returns a new EthereumArc instance
func NewArc(context context.Context, conn ethereum.Conn, auth *bind.TransactOpts, swapID [32]byte) (arc.Arc, error) {
	return &Arc{
		context: context,
		conn:    conn,
		auth:    auth,
		arcData: EthereumArcData{},
		binding: nil,
		swapID:  swapID,
	}, nil
}

// Initiate a new Arc swap by deploying an Arc contract
func (arc *Arc) Initiate(hash [32]byte, from []byte, to []byte, value *big.Int, expiry int64) error {
	contractAddress, tx, binding, err := bindings.DeployArc(arc.auth, arc.conn.Client, hash, common.HexToAddress("0x1"), value, big.NewInt(expiry), common.BytesToAddress(to))
	if err != nil {
		return err
	}

	if err := arc.conn.TransferEth(arc.context, arc.auth, contractAddress, value); err != nil {
		return err
	}

	arc.arcData = EthereumArcData{
		contractAddress, value,
	}
	arc.binding = binding

	_, err = arc.conn.PatchedWaitDeployed(arc.context, tx)
	return err
}

func (arc *Arc) Redeem(secret [32]byte) error {
	tx, err := arc.binding.Redeem(arc.auth, secret)
	if err == nil {
		_, err = arc.conn.PatchedWaitMined(arc.context, tx)
	}
	return err
}

func (arc *Arc) Refund() error {
	tx, err := arc.binding.Refund(arc.auth, common.HexToAddress("0x1"), arc.arcData.Value)
	if err == nil {
		_, err = arc.conn.PatchedWaitMined(arc.context, tx)
	}
	return err
}

func (arc *Arc) Audit() (hash [32]byte, to, from []byte, value *big.Int, expiry int64, err error) {
	hash, _, toAddr, value, _expiry, err := arc.binding.Audit(&bind.CallOpts{})
	if err != nil {
		return [32]byte{}, nil, nil, nil, 0, err
	}
	// arc.client.Client.CodeAt(arc.context, arc.ethereumData.contractAddress, nil)
	// TODO: Audit values, bytecode
	expiry = _expiry.Int64()
	return hash, arc.arcData.ContractAddress.Bytes(), toAddr.Bytes(), value, expiry, err
}

func (arc *Arc) AuditSecret() ([32]byte, error) {
	secret, err := arc.binding.AuditSecret(&bind.CallOpts{})
	// if err != nil {
	// 	return [32]byte{}, err
	// }
	// var secret32 [32]byte
	// copy(secret, secret32[:])
	// return secret32, nil
	return secret, err
}

func (arc *Arc) Serialize() ([]byte, error) {
	b, err := json.Marshal(arc.arcData)
	return b, err
}

func (arc *Arc) Deserialize(b []byte) error {
	if err := json.Unmarshal(b, &arc.arcData); err != nil {
		return err
	}

	contract, err := bindings.NewArc(arc.arcData.ContractAddress, bind.ContractBackend(arc.conn.Client))
	if err != nil {
		return err
	}
	arc.binding = contract
	return nil
}
