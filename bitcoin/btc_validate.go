package bitcoin

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

type validateResult struct {
	contractAddress  []byte
	amount           int64
	recipientAddress []byte
	refundAddress    []byte
	secretHash       []byte
	lockTime         int64
}

func Validate(contract, contractTxBytes []byte, chain string, rpcuser string, rpcpass string) (Error error, result validateResult) {
	var chainParams *chaincfg.Params
	if chain == "testnet" {
		chainParams = &chaincfg.TestNet3Params
	} else {
		chainParams = &chaincfg.MainNetParams
	}

	var contractTx wire.MsgTx
	err := contractTx.Deserialize(bytes.NewReader(contractTxBytes))
	if err != nil {
		return fmt.Errorf("failed to decode contract transaction: %v", err), validateResult{}
	}

	contractHash160 := btcutil.Hash160(contract)
	contractOut := -1

	for i, out := range contractTx.TxOut {
		sc, addrs, _, err := txscript.ExtractPkScriptAddrs(out.PkScript, chainParams)
		if err != nil || sc != txscript.ScriptHashTy {
			continue
		}
		if bytes.Equal(addrs[0].(*btcutil.AddressScriptHash).Hash160()[:], contractHash160) {
			contractOut = i
			break
		}
	}
	if contractOut == -1 {
		return errors.New("transaction does not contain the contract output"), validateResult{}
	}

	pushes, err := txscript.ExtractAtomicSwapDataPushes(contract)
	if err != nil {
		return err, validateResult{}
	}
	if pushes == nil {
		return errors.New("contract is not an atomic swap script recognized by this tool"), validateResult{}
	}

	contractAddr, err := btcutil.NewAddressScriptHash(contract, chainParams)
	if err != nil {
		return err, validateResult{}
	}
	recipientAddr, err := btcutil.NewAddressPubKeyHash(pushes.RecipientHash160[:],
		chainParams)
	if err != nil {
		return err, validateResult{}
	}
	refundAddr, err := btcutil.NewAddressPubKeyHash(pushes.RefundHash160[:],
		chainParams)
	if err != nil {
		return err, validateResult{}
	}

	return nil, validateResult{
		contractAddress:  contractAddr.ScriptAddress(),
		amount:           int64(btcutil.Amount(contractTx.TxOut[contractOut].Value)),
		recipientAddress: recipientAddr.ScriptAddress(),
		refundAddress:    refundAddr.ScriptAddress(),
		secretHash:       pushes.SecretHash[:],
		lockTime:         pushes.LockTime,
	}
}
