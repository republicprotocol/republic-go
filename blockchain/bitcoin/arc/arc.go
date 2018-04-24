package arc

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

type Arc struct {
	conn       bitcoin.Conn
	ledgerData BitcoinData
}

// NewArc returns an arc object
func NewArc(conn bitcoin.Conn) arc.Arc {
	return &Arc{
		conn: conn,
	}
}

func (arc *Arc) Initiate(hash [32]byte, from, to []byte, value *big.Int, expiry int64) (err error) {
	result, err := initiate(arc.conn, string(to), value.Int64(), hash[:], expiry)
	if err != nil {
		return err
	}
	arc.ledgerData = result
	arc.ledgerData.SecretHash = hash
	return nil
}

func (arc *Arc) Audit() (hash [32]byte, from, to []byte, value *big.Int, expiry int64, err error) {
	result, err := read(arc.conn, arc.ledgerData.Contract, arc.ledgerData.ContractTx)
	if err != nil {
		return [32]byte{}, []byte{}, []byte{}, big.NewInt(0), 0, err
	}
	return result.secretHash, result.refundAddress, result.recipientAddress, big.NewInt(result.amount), result.lockTime, nil
}

func (arc *Arc) Redeem(secret [32]byte) error {
	result, err := redeem(arc.conn, arc.ledgerData.Contract, arc.ledgerData.ContractTx, secret)
	if err != nil {
		return err
	}
	arc.ledgerData.RedeemTx = result.redeemTx
	arc.ledgerData.RedeemTxHash = result.redeemTxHash
	return nil
}

func (arc *Arc) AuditSecret() (secret [32]byte, err error) {
	result, err := readSecret(arc.conn, arc.ledgerData.RedeemTx, arc.ledgerData.SecretHash[:])
	if err != nil {
		return [32]byte{}, err
	}
	return result, nil
}

func (arc *Arc) Refund() error {
	return refund(arc.conn, arc.ledgerData.Contract, arc.ledgerData.ContractTx)
}

func (arc *Arc) Serialize() ([]byte, error) {
	b, err := json.Marshal(arc.ledgerData)
	return b, err
}

func (arc *Arc) Deserialize(b []byte) error {
	return json.Unmarshal(b, &arc.ledgerData)
}
