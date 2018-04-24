package arc_bitcoin

import (
	"encoding/json"
	"math/big"

	"github.com/republicprotocol/republic-go/blockchain/arc"
	"github.com/republicprotocol/republic-go/blockchain/bitcoin"
)

type BitcoinData struct {
	ContractHash   string   `json:"contract_hash"`
	Contract       []byte   `json:"contract"`
	ContractTxHash []byte   `json:"contract_txhash"`
	ContractTx     []byte   `json:"contract_tx"`
	RefundTxHash   [32]byte `json:"refund_txhash"`
	RefundTx       []byte   `json:"refund_tx"`
	RedeemTxHash   [32]byte `json:"redeem_txhash"`
	RedeemTx       []byte   `json:"redeem_tx"`
	SecretHash     [32]byte `json:"secret_hash"`
}

type BitcoinArc struct {
	connection client.Connection
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
	arc.ledgerData.SecretHash = hash
	return nil
}

func (arc *BitcoinArc) Audit() (hash [32]byte, from, to []byte, value *big.Int, expiry int64, err error) {
	result, err := read(arc.connection, arc.ledgerData.Contract, arc.ledgerData.ContractTx)
	if err != nil {
		return [32]byte{}, []byte{}, []byte{}, big.NewInt(0), 0, err
	}
	return result.secretHash, result.refundAddress, result.recipientAddress, big.NewInt(result.amount), result.lockTime, nil
}

func (arc *BitcoinArc) Redeem(secret [32]byte) error {
	result, err := redeem(arc.connection, arc.ledgerData.Contract, arc.ledgerData.ContractTx, secret)
	if err != nil {
		return err
	}
	arc.ledgerData.RedeemTx = result.redeemTx
	arc.ledgerData.RedeemTxHash = result.redeemTxHash
	return nil
}

func (arc *BitcoinArc) AuditSecret() (secret [32]byte, err error) {
	result, err := readSecret(arc.connection, arc.ledgerData.RedeemTx, arc.ledgerData.SecretHash[:])
	if err != nil {
		return [32]byte{}, err
	}
	return result, nil
}

func (arc *BitcoinArc) Refund() error {
	return refund(arc.connection, arc.ledgerData.Contract, arc.ledgerData.ContractTx)
}

func (arc *BitcoinArc) Serialize() ([]byte, error) {
	b, err := json.Marshal(arc.ledgerData)
	return b, err
}

func (arc *BitcoinArc) Deserialize(b []byte) error {
	return json.Unmarshal(b, &arc.ledgerData)
}
