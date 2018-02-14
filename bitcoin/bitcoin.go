package bitcoin

import (
	"math/big"
)

type BTCAtomContract struct {
	rpcUser    string
	rpcPass    string
	chain      string
	secretHash []byte
	ledgerData BitcoinData
}

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

func NewBTCAtomContract(rpcUser string, rpcPass string, chain string) *BTCAtomContract {
	return &BTCAtomContract{
		rpcUser: rpcUser,
		rpcPass: rpcPass,
		chain:   chain,
	}
}

func (contract *BTCAtomContract) Initiate(hash, to, from []byte, value *big.Int, expiry int64) (err error) {
	err, result := initiate(string(to), value.Int64(), contract.chain, contract.rpcUser, contract.rpcPass, hash, expiry)
	if err != nil {
		return err
	}
	contract.ledgerData = result
	contract.secretHash = hash
	return nil
}

func (contract *BTCAtomContract) Read() (hash, to, from []byte, value *big.Int, expiry int64, err error) {
	err, result := read(contract.ledgerData.contract, contract.ledgerData.contractTx, contract.chain, contract.rpcUser, contract.rpcPass)
	if err != nil {
		return []byte{}, []byte{}, []byte{}, big.NewInt(0), 0, err
	}
	return result.secretHash, result.recipientAddress, result.refundAddress, big.NewInt(result.amount), result.lockTime, nil
}

func (contract *BTCAtomContract) Redeem(secret []byte) error {
	err, result := redeem(contract.ledgerData.contract, contract.ledgerData.contractTx, secret, contract.rpcUser, contract.rpcPass, contract.chain)
	if err != nil {
		return err
	}
	contract.ledgerData.redeemTx = result.redeemTx
	contract.ledgerData.redeemTxHash = result.redeemTxHash
	return nil
}

func (contract *BTCAtomContract) ReadSecret() (secret []byte, err error) {
	err, result := readSecret(contract.ledgerData.redeemTx, contract.secretHash, contract.rpcUser, contract.rpcPass)
	if err != nil {
		return []byte{}, err
	}
	return result, nil
}

func (contract *BTCAtomContract) Refund() error {
	return refund(contract.ledgerData.contract, contract.ledgerData.contractTx, contract.chain, contract.rpcUser, contract.rpcPass)
}
