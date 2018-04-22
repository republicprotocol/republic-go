package arc_bitcoin

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

type readResult struct {
	contractAddress  []byte
	amount           int64
	recipientAddress []byte
	refundAddress    []byte
	secretHash       [32]byte
	lockTime         int64
}

func read(contract, contractTxBytes []byte, chain string, rpcuser string, rpcpass string) (Error error, result readResult) {
	var chainParams *chaincfg.Params
	if chain == "regtest" {
		chainParams = &chaincfg.RegressionNetParams
	} else if chain == "testnet" {
		chainParams = &chaincfg.TestNet3Params
	} else {
		chainParams = &chaincfg.MainNetParams
	}

	var contractTx wire.MsgTx
	err := contractTx.Deserialize(bytes.NewReader(contractTxBytes))
	if err != nil {
		return fmt.Errorf("failed to decode contract transaction: %v", err), readResult{}
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
		return errors.New("transaction does not contain the contract output"), readResult{}
	}

	pushes, err := txscript.ExtractAtomicSwapDataPushes(0, contract)
	if err != nil {
		return err, readResult{}
	}
	if pushes == nil {
		return errors.New("contract is not an atomic swap script recognized by this tool"), readResult{}
	}

	contractAddr, err := btcutil.NewAddressScriptHash(contract, chainParams)
	if err != nil {
		return err, readResult{}
	}
	recipientAddr, err := btcutil.NewAddressPubKeyHash(pushes.RecipientHash160[:],
		chainParams)
	if err != nil {
		return err, readResult{}
	}
	refundAddr, err := btcutil.NewAddressPubKeyHash(pushes.RefundHash160[:],
		chainParams)
	if err != nil {
		return err, readResult{}
	}

	return nil, readResult{
		contractAddress:  contractAddr.ScriptAddress(),
		amount:           int64(btcutil.Amount(contractTx.TxOut[contractOut].Value)),
		recipientAddress: []byte(recipientAddr.EncodeAddress()),
		refundAddress:    []byte(refundAddr.EncodeAddress()),
		secretHash:       pushes.SecretHash,
		lockTime:         pushes.LockTime,
	}
}
