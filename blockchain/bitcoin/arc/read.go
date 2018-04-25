package arc

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/republicprotocol/republic-go/blockchain/bitcoin"
)

type readResult struct {
	contractAddress  []byte
	amount           int64
	recipientAddress []byte
	refundAddress    []byte
	secretHash       [32]byte
	lockTime         int64
}

func read(conn bitcoin.Conn, contract, contractTxBytes []byte) (readResult, error) {

	var contractTx wire.MsgTx
	err := contractTx.Deserialize(bytes.NewReader(contractTxBytes))
	if err != nil {
		return readResult{}, fmt.Errorf("failed to decode contract transaction: %v", err)
	}

	contractHash160 := btcutil.Hash160(contract)
	contractOut := -1

	for i, out := range contractTx.TxOut {
		sc, addrs, _, err := txscript.ExtractPkScriptAddrs(out.PkScript, conn.ChainParams)
		if err != nil || sc != txscript.ScriptHashTy {
			continue
		}
		if bytes.Equal(addrs[0].(*btcutil.AddressScriptHash).Hash160()[:], contractHash160) {
			contractOut = i
			break
		}
	}
	if contractOut == -1 {
		return readResult{}, errors.New("transaction does not contain the contract output")
	}

	pushes, err := txscript.ExtractAtomicSwapDataPushes(0, contract)
	if err != nil {
		return readResult{}, err
	}
	if pushes == nil {
		return readResult{}, errors.New("contract is not an atomic swap script recognized by this tool")
	}

	contractAddr, err := btcutil.NewAddressScriptHash(contract, conn.ChainParams)
	if err != nil {
		return readResult{}, err
	}
	recipientAddr, err := btcutil.NewAddressPubKeyHash(pushes.RecipientHash160[:],
		conn.ChainParams)
	if err != nil {
		return readResult{}, err
	}
	refundAddr, err := btcutil.NewAddressPubKeyHash(pushes.RefundHash160[:],
		conn.ChainParams)
	if err != nil {
		return readResult{}, err
	}

	return readResult{
		contractAddress:  contractAddr.ScriptAddress(),
		amount:           int64(btcutil.Amount(contractTx.TxOut[contractOut].Value)),
		recipientAddress: []byte(recipientAddr.EncodeAddress()),
		refundAddress:    []byte(refundAddr.EncodeAddress()),
		secretHash:       pushes.SecretHash,
		lockTime:         pushes.LockTime,
	}, nil
}
