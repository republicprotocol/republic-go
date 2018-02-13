package bitcoin

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

type closeResult struct {
	redeemTx     []byte
	redeemTxHash [32]byte
}

func Close(contract, contractTxBytes, secret []byte, rpcUser string, rpcPass string, chain string) (Error error, result closeResult) {
	var chainParams *chaincfg.Params
	if chain == "testnet" {
		chainParams = &chaincfg.TestNet3Params
	} else {
		chainParams = &chaincfg.MainNetParams
	}

	var contractTx wire.MsgTx
	err := contractTx.Deserialize(bytes.NewReader(contractTxBytes))
	if err != nil {
		return fmt.Errorf("failed to decode contract transaction: %v", err), closeResult{}
	}

	connect, err := normalizeAddress("localhost", walletPort(chainParams))
	if err != nil {
		return fmt.Errorf("wallet server address: %v", err), closeResult{}
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
		return fmt.Errorf("rpc connect: %v", err), closeResult{}
	}
	defer func() {
		client.Shutdown()
		client.WaitForShutdown()
	}()

	pushes, err := txscript.ExtractAtomicSwapDataPushes(contract)
	if err != nil {
		return err, closeResult{}
	}
	if pushes == nil {
		return errors.New("contract is not an atomic swap script recognized by this tool"), closeResult{}
	}
	recipientAddr, err := btcutil.NewAddressPubKeyHash(pushes.RecipientHash160[:],
		chainParams)
	if err != nil {
		return err, closeResult{}
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
		return errors.New("transaction does not contain a contract output"), closeResult{}
	}

	addr, err := getRawChangeAddress(client, chainParams)
	if err != nil {
		return fmt.Errorf("getrawchangeaddres: %v", err), closeResult{}
	}
	outScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return err, closeResult{}
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
		return err, closeResult{}
	}
	redeemSigScript, err := redeemP2SHContract(contract, redeemSig, redeemPubKey, secret)
	if err != nil {
		return err, closeResult{}
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
			panic(err)
		}
		err = e.Execute()
		if err != nil {
			panic(err)
		}
	}

	return promptPublishTx(client, redeemTx, "redeem"), closeResult{
		redeemTx:     buf.Bytes(),
		redeemTxHash: redeemTxHash,
	}
}
