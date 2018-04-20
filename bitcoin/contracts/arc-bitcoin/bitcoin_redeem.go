package arc_bitcoin

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	rpc "github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

type redeemResult struct {
	redeemTx     []byte
	redeemTxHash [32]byte
}

func redeem(contract, contractTxBytes, secret []byte, rpcUser string, rpcPass string, chain string) (Error error, result redeemResult) {
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
		return fmt.Errorf("failed to decode contract transaction: %v", err), redeemResult{}
	}

	connect, err := normalizeAddress("localhost", walletPort(chainParams))
	if err != nil {
		return fmt.Errorf("wallet server address: %v", err), redeemResult{}
	}

	connConfig := &rpc.ConnConfig{
		Host:         connect,
		User:         rpcUser,
		Pass:         rpcPass,
		DisableTLS:   true,
		HTTPPostMode: true,
	}

	client, err := rpc.New(connConfig, nil)
	if err != nil {
		return fmt.Errorf("rpc connect: %v", err), redeemResult{}
	}
	defer func() {
		client.Shutdown()
		client.WaitForShutdown()
	}()

	pushes, err := txscript.ExtractAtomicSwapDataPushes(0, contract)
	if err != nil {
		return err, redeemResult{}
	}
	if pushes == nil {
		return errors.New("contract is not an atomic swap script recognized by this tool"), redeemResult{}
	}
	recipientAddr, err := btcutil.NewAddressPubKeyHash(pushes.RecipientHash160[:],
		chainParams)
	if err != nil {
		return err, redeemResult{}
	}
	contractHash := btcutil.Hash160(contract)
	contractOut := -1
	for i, out := range contractTx.TxOut {
		sc, addrs, _, _ := txscript.ExtractPkScriptAddrs(out.PkScript, chainParams)
		if sc == txscript.ScriptHashTy &&
			bytes.Equal(addrs[0].(*btcutil.AddressScriptHash).Hash160()[:], contractHash) {
			contractOut = i
			break
		}
	}
	if contractOut == -1 {
		return errors.New("transaction does not contain a contract output"), redeemResult{}
	}

	addr, err := getRawChangeAddress(client, chainParams)
	if err != nil {
		return fmt.Errorf("getrawchangeaddres: %v", err), redeemResult{}
	}
	outScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return err, redeemResult{}
	}

	contractTxHash := contractTx.TxHash()
	contractOutPoint := wire.OutPoint{
		Hash:  contractTxHash,
		Index: uint32(contractOut),
	}

	redeemTx := wire.NewMsgTx(txVersion)
	redeemTx.LockTime = uint32(pushes.LockTime)
	redeemTx.AddTxIn(wire.NewTxIn(&contractOutPoint, nil, nil))
	redeemTx.AddTxOut(wire.NewTxOut(0, outScript)) // amount set below
	redeemSig, redeemPubKey, err := createSig(redeemTx, 0, contract, recipientAddr, client)
	if err != nil {
		return err, redeemResult{}
	}
	redeemSigScript, err := redeemP2SHContract(contract, redeemSig, redeemPubKey, secret)
	if err != nil {
		return err, redeemResult{}
	}
	redeemTx.TxIn[0].SignatureScript = redeemSigScript

	redeemTxHash := redeemTx.TxHash()

	var buf bytes.Buffer
	buf.Grow(redeemTx.SerializeSize())
	redeemTx.Serialize(&buf)

	if verify {
		e, err := txscript.NewEngine(contractTx.TxOut[contractOutPoint.Index].PkScript,
			redeemTx, 0, txscript.StandardVerifyFlags, txscript.NewSigCache(10),
			txscript.NewTxSigHashes(redeemTx), contractTx.TxOut[contractOut].Value)
		if err != nil {
			return err, redeemResult{}
		}
		err = e.Execute()
		if err != nil {
			return err, redeemResult{}
		}
	}

	return promptPublishTx(client, redeemTx, "redeem"), redeemResult{
		redeemTx:     buf.Bytes(),
		redeemTxHash: redeemTxHash,
	}
}
