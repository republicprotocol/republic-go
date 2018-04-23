package arc_bitcoin

import (
	"math/big"

	interop "github.com/republicprotocol/republic-go/interop"
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
	rpcUser    string
	rpcPass    string
	chain      string
	secretHash [32]byte
	ledgerData BitcoinData
}

// NewBitcoinArc returns an arc object
func NewBitcoinArc(rpcUser, rpcPass, chain string) interop.Arc {
	return &BitcoinArc{
		rpcUser: rpcUser,
		rpcPass: rpcPass,
		chain:   chain,
	}
}

func (arc *BitcoinArc) Initiate(hash [32]byte, from, to []byte, value *big.Int, expiry int64) (err error) {
	err, result := initiate(string(to), value.Int64(), arc.chain, arc.rpcUser, arc.rpcPass, hash[:], expiry)
	if err != nil {
		return err
	}
	arc.ledgerData = result
	arc.secretHash = hash
	return nil
}

func (arc *BitcoinArc) Audit() (hash [32]byte, from, to []byte, value *big.Int, expiry int64, err error) {
	err, result := read(arc.ledgerData.contract, arc.ledgerData.contractTx, arc.chain, arc.rpcUser, arc.rpcPass)
	if err != nil {
		return [32]byte{}, []byte{}, []byte{}, big.NewInt(0), 0, err
	}
	return result.secretHash, result.refundAddress, result.recipientAddress, big.NewInt(result.amount), result.lockTime, nil
}

func (arc *BitcoinArc) Redeem(secret []byte) error {
	err, result := redeem(arc.ledgerData.contract, arc.ledgerData.contractTx, secret, arc.rpcUser, arc.rpcPass, arc.chain)
	if err != nil {
		return err
	}
	arc.ledgerData.redeemTx = result.redeemTx
	arc.ledgerData.redeemTxHash = result.redeemTxHash
	return nil
}

func (arc *BitcoinArc) AuditSecret() (secret []byte, err error) {
	err, result := readSecret(arc.ledgerData.redeemTx, arc.secretHash[:], arc.rpcUser, arc.rpcPass)
	if err != nil {
		return []byte{}, err
	}
	return result, nil
}

func (arc *BitcoinArc) Refund() error {
	return refund(arc.ledgerData.contract, arc.ledgerData.contractTx, arc.chain, arc.rpcUser, arc.rpcPass)
}
