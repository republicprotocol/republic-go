package arc_bitcoin

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/republicprotocol/republic-go/blockchain/bitcoin"
)

type redeemResult struct {
	redeemTx     []byte
	redeemTxHash [32]byte
}

func redeem(connection client.Connection, contract, contractTxBytes []byte, secret [32]byte) (redeemResult, error) {
	var contractTx wire.MsgTx
	err := contractTx.Deserialize(bytes.NewReader(contractTxBytes))
	if err != nil {
		return redeemResult{}, fmt.Errorf("failed to decode contract transaction: %v", err)
	}

	pushes, err := txscript.ExtractAtomicSwapDataPushes(0, contract)
	if err != nil {
		return redeemResult{}, err
	}
	if pushes == nil {
		return redeemResult{}, errors.New("contract is not an atomic swap script recognized by this tool")
	}
	recipientAddr, err := btcutil.NewAddressPubKeyHash(pushes.RecipientHash160[:],
		connection.ChainParams)
	if err != nil {
		return redeemResult{}, err
	}
	contractHash := btcutil.Hash160(contract)
	contractOut := -1
	for i, out := range contractTx.TxOut {
		sc, addrs, _, _ := txscript.ExtractPkScriptAddrs(out.PkScript, connection.ChainParams)
		if sc == txscript.ScriptHashTy &&
			bytes.Equal(addrs[0].(*btcutil.AddressScriptHash).Hash160()[:], contractHash) {
			contractOut = i
			break
		}
	}
	if contractOut == -1 {
		return redeemResult{}, errors.New("transaction does not contain a contract output")
	}

	addr, err := getRawChangeAddress(connection)
	if err != nil {
		return redeemResult{}, fmt.Errorf("getrawchangeaddres: %v", err)
	}
	outScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return redeemResult{}, err
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
	redeemSig, redeemPubKey, err := createSig(connection, redeemTx, 0, contract, recipientAddr)
	if err != nil {
		return redeemResult{}, err
	}
	redeemSigScript, err := redeemP2SHContract(contract, redeemSig, redeemPubKey, secret)
	if err != nil {
		return redeemResult{}, err
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
			return redeemResult{}, err
		}
		err = e.Execute()
		if err != nil {
			return redeemResult{}, err
		}
	}

	txHash, err := connection.PromptPublishTx(redeemTx, "redeem")
	if err != nil {
		return redeemResult{}, err
	}

	connection.WaitForConfirmations(txHash, 1)

	return redeemResult{
		redeemTx:     buf.Bytes(),
		redeemTxHash: redeemTxHash,
	}, nil
}
