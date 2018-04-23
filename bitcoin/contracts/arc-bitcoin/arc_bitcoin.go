package arc_bitcoin

import (
	"encoding/json"
	"math/big"

	"github.com/republicprotocol/republic-go/bitcoin/client"
	arc "github.com/republicprotocol/republic-go/interop/arc"
)

type BitcoinData struct {
	contractHash   string
	contract       []byte
	contractTxHash []byte
	contractTx     []byte
	refundTxHash   [32]byte
	refundTx       []byte
	redeemTxHash   [32]byte
	redeemTx       []byte
}

type BitcoinArc struct {
	connection client.Connection
	secretHash [32]byte
	ledgerData BitcoinData
}

// NewBitcoinArc returns an arc object
func NewBitcoinArc(connection client.Connection) arc.Arc {
	return &BitcoinArc{
		connection: connection,
	}
}

func (arc *BitcoinArc) Initiate(hash [32]byte, from, to []byte, value *big.Int, expiry int64) (err error) {
	result, err := initiate(arc.connection, string(to), value.Int64(), hash[:], expiry)
	if err != nil {
		return err
	}
	arc.ledgerData = result
	arc.secretHash = hash
	return nil
}

func (arc *BitcoinArc) Audit() (hash [32]byte, from, to []byte, value *big.Int, expiry int64, err error) {
	result, err := read(arc.connection, arc.ledgerData.contract, arc.ledgerData.contractTx)
	if err != nil {
		return [32]byte{}, []byte{}, []byte{}, big.NewInt(0), 0, err
	}
	return result.secretHash, result.refundAddress, result.recipientAddress, big.NewInt(result.amount), result.lockTime, nil
}

func (arc *BitcoinArc) Redeem(secret [32]byte) error {
	result, err := redeem(arc.connection, arc.ledgerData.contract, arc.ledgerData.contractTx, secret)
	if err != nil {
		return err
	}
	arc.ledgerData.redeemTx = result.redeemTx
	arc.ledgerData.redeemTxHash = result.redeemTxHash
	return nil
}

func (arc *BitcoinArc) AuditSecret() (secret [32]byte, err error) {
	result, err := readSecret(arc.connection, arc.ledgerData.redeemTx, arc.secretHash[:])
	if err != nil {
		return [32]byte{}, err
	}
	return result, nil
}

func (arc *BitcoinArc) Refund() error {
	return refund(arc.connection, arc.ledgerData.contract, arc.ledgerData.contractTx)
}

func (arc *BitcoinArc) Serialize() ([]byte, error) {
	b, err := json.Marshal(arc.ledgerData)
	return b, err
}

func (arc *BitcoinArc) Deserialize(b []byte) error {
	return json.Unmarshal(b, &arc.ledgerData)
}
