package arc

import (
	"context"
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/republicprotocol/republic-go/blockchain/arc"
	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/bindings"
)

var ETHEREUM = common.HexToAddress("0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE")

type EthereumArcData struct {
	ContractAddress common.Address `json:"contract_address"`
	Value           *big.Int       `json:"value"`
}

type Arc struct {
	order        []byte
	fee          *big.Int
	tokenAddress common.Address
	context      context.Context
	conn         ethereum.Conn
	auth         *bind.TransactOpts
	arcData      EthereumArcData
	binding      *bindings.Arc
	tokenBinding *bindings.ERC20
}

// NewArc returns a new EthereumArc instance
func NewArc(context context.Context, conn ethereum.Conn, auth *bind.TransactOpts, order []byte, tokenAddress common.Address, fee *big.Int) (arc.Arc, error) {
	arc, err := bindings.NewArc(conn.ArcAddress, conn.Client)
	if err != nil {
		return &Arc{}, err
	}

	var tokenBinding *bindings.ERC20

	if tokenAddress != ETHEREUM {
		tokenBinding, err = bindings.NewERC20(tokenAddress, conn.Client)
		if err != nil {
			return &Arc{}, err
		}
	}

	return &Arc{
		order:        order,
		context:      context,
		conn:         conn,
		auth:         auth,
		arcData:      EthereumArcData{},
		binding:      arc,
		tokenBinding: tokenBinding,
		tokenAddress: tokenAddress,
		fee:          fee,
	}, nil
}

// Initiate a new Arc swap on the Arc contract
func (arc *Arc) Initiate(hash [32]byte, from []byte, to []byte, value *big.Int, expiry int64) error {
	var tx *types.Transaction
	var err error
	if arc.tokenAddress == ETHEREUM {
		arc.auth.Value = value
		tx, err = arc.binding.Initiate(arc.auth, hash, arc.tokenAddress, value, arc.fee, big.NewInt(expiry), common.BytesToAddress(to), arc.order)
		arc.auth.Value = big.NewInt(0)
	} else {
		tx, err = arc.tokenBinding.Approve(arc.auth, arc.conn.ArcAddress, value)
	}

	if err != nil {
		return err
	}

	_, err = arc.conn.PatchedWaitDeployed(arc.context, tx)
	return err
}

func (arc *Arc) Redeem(orderID, secret [32]byte) error {
	tx, err := arc.binding.Redeem(arc.auth, orderID, secret)
	if err == nil {
		_, err = arc.conn.PatchedWaitMined(arc.context, tx)
	}
	return err
}

func (arc *Arc) Refund(orderID [32]byte) error {
	tx, err := arc.binding.Refund(arc.auth, orderID, arc.tokenAddress, arc.arcData.Value)
	if err == nil {
		_, err = arc.conn.PatchedWaitMined(arc.context, tx)
	}
	return err
}

func (arc *Arc) Audit(orderID [32]byte) (hash [32]byte, to, from []byte, value *big.Int, expiry int64, err error) {
	hash, _, toAddr, value, _expiry, err := arc.binding.Audit(&bind.CallOpts{}, orderID)
	if err != nil {
		return [32]byte{}, nil, nil, nil, 0, err
	}
	// arc.client.Client.CodeAt(arc.context, arc.ethereumData.contractAddress, nil)
	// TODO: Audit values, bytecode
	expiry = _expiry.Int64()
	return hash, arc.arcData.ContractAddress.Bytes(), toAddr.Bytes(), value, expiry, err
}

func (arc *Arc) AuditSecret(orderID [32]byte) ([32]byte, error) {
	secret, err := arc.binding.AuditSecret(&bind.CallOpts{}, orderID)
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
