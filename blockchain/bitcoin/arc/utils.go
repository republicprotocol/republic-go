package arc

import (
	"bytes"
	"errors"
	"fmt"

	"crypto/sha256"

	"encoding/json"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/republicprotocol/republic-go/blockchain/bitcoin"
	"golang.org/x/crypto/ripemd160"
)

const txVersion = 2

const secretSize = 32

const verify = true

type builtContract struct {
	contract       []byte
	contractP2SH   btcutil.Address
	contractTxHash *chainhash.Hash
	contractTx     *wire.MsgTx
	refundTx       *wire.MsgTx
}

type contractArgs struct {
	them       *btcutil.AddressPubKeyHash
	amount     int64
	locktime   int64
	secretHash []byte
}

func sumOutputSerializeSizes(outputs []*wire.TxOut) (serializeSize int) {
	for _, txOut := range outputs {
		serializeSize += txOut.SerializeSize()
	}
	return serializeSize
}

func inputSize(sigScriptSize int) int {
	return 32 + 4 + wire.VarIntSerializeSize(uint64(sigScriptSize)) + sigScriptSize + 4
}

func estimateRedeemSerializeSize(contract []byte, txOuts []*wire.TxOut) int {
	contractPush, err := txscript.NewScriptBuilder().AddData(contract).Script()
	if err != nil {
		panic(err)
	}
	contractPushSize := len(contractPush)

	return 12 + wire.VarIntSerializeSize(1) +
		wire.VarIntSerializeSize(uint64(len(txOuts))) +
		inputSize(redeemAtomicSwapSigScriptSize+contractPushSize) +
		sumOutputSerializeSizes(txOuts)
}

func buildContract(conn bitcoin.Conn, args *contractArgs) (*builtContract, error) {

	refundAddr, err := getRawChangeAddress(conn)
	if err != nil {
		return nil, fmt.Errorf("getrawchangeaddress: %v", err)
	}
	refundAddrH, ok := refundAddr.(interface {
		Hash160() *[ripemd160.Size]byte
	})
	if !ok {
		return nil, errors.New("unable to create hash160 from change address")
	}

	contract, err := atomicSwapContract(refundAddrH.Hash160(), args.them.Hash160(),
		args.locktime, args.secretHash)
	if err != nil {
		return nil, err
	}
	contractP2SH, err := btcutil.NewAddressScriptHash(contract, conn.ChainParams)
	if err != nil {
		return nil, err
	}
	contractP2SHPkScript, err := txscript.PayToAddrScript(contractP2SH)
	if err != nil {
		return nil, err
	}

	unsignedContract := wire.NewMsgTx(txVersion)
	unsignedContract.AddTxOut(wire.NewTxOut(int64(args.amount), contractP2SHPkScript))
	unsignedContract, err = conn.FundRawTransaction(unsignedContract)
	if err != nil {
		return nil, fmt.Errorf("fundrawtransaction: %v", err)
	}
	contractTx, complete, err := conn.Client.SignRawTransaction(unsignedContract)
	if err != nil {
		return nil, fmt.Errorf("signrawtransaction: %v", err)
	}
	if !complete {
		return nil, errors.New("signrawtransaction: failed to completely sign contract transaction")
	}

	contractTxHash := contractTx.TxHash()

	refundTx, err := buildRefund(conn, contract, contractTx)
	if err != nil {
		return nil, err
	}

	return &builtContract{
		contract,
		contractP2SH,
		&contractTxHash,
		contractTx,
		refundTx,
	}, nil
}

func sha256Hash(x []byte) []byte {
	h := sha256.Sum256(x)
	return h[:]
}

func getRawChangeAddress(conn bitcoin.Conn) (btcutil.Address, error) {
	rawResp, err := conn.Client.RawRequest("getrawchangeaddress", nil)
	if err != nil {
		return nil, err
	}
	var addrStr string
	err = json.Unmarshal(rawResp, &addrStr)
	if err != nil {
		return nil, err
	}
	addr, err := btcutil.DecodeAddress(addrStr, conn.ChainParams)
	if err != nil {
		return nil, err
	}
	if !addr.IsForNet(conn.ChainParams) {
		return nil, fmt.Errorf("address %v is not intended for use on %v",
			addrStr, conn.ChainParams.Name)
	}
	return addr, nil
}

func buildRefund(conn bitcoin.Conn, contract []byte, contractTx *wire.MsgTx) (
	refundTx *wire.MsgTx, err error) {

	contractP2SH, err := btcutil.NewAddressScriptHash(contract, conn.ChainParams)
	if err != nil {
		return nil, err
	}
	contractP2SHPkScript, err := txscript.PayToAddrScript(contractP2SH)
	if err != nil {
		return nil, err
	}

	contractTxHash := contractTx.TxHash()
	contractOutPoint := wire.OutPoint{Hash: contractTxHash, Index: ^uint32(0)}
	for i, o := range contractTx.TxOut {
		if bytes.Equal(o.PkScript, contractP2SHPkScript) {
			contractOutPoint.Index = uint32(i)
			break
		}
	}
	if contractOutPoint.Index == ^uint32(0) {
		return nil, errors.New("contract tx does not contain a P2SH contract payment")
	}

	refundAddress, err := getRawChangeAddress(conn)
	if err != nil {
		return nil, fmt.Errorf("getrawchangeaddress: %v", err)
	}
	refundOutScript, err := txscript.PayToAddrScript(refundAddress)
	if err != nil {
		return nil, err
	}

	pushes, err := txscript.ExtractAtomicSwapDataPushes(0, contract)
	if err != nil {
		return nil, err
	}
	if pushes == nil {
		return nil, fmt.Errorf("failed to extract atomic swap data")
	}

	refundAddr, err := btcutil.NewAddressPubKeyHash(pushes.RefundHash160[:], conn.ChainParams)
	if err != nil {
		return nil, err
	}

	refundTx = wire.NewMsgTx(txVersion)
	refundTx.LockTime = uint32(pushes.LockTime)
	refundTx.AddTxOut(wire.NewTxOut(0, refundOutScript))

	txIn := wire.NewTxIn(&contractOutPoint, nil, nil)
	txIn.Sequence = 0
	refundTx.AddTxIn(txIn)

	refundSig, refundPubKey, err := createSig(conn, refundTx, 0, contract, refundAddr)
	if err != nil {
		return nil, err
	}
	refundSigScript, err := refundP2SHContract(contract, refundSig, refundPubKey)
	if err != nil {
		return nil, err
	}
	refundTx.TxIn[0].SignatureScript = refundSigScript

	if verify {
		e, err := txscript.NewEngine(contractTx.TxOut[contractOutPoint.Index].PkScript,
			refundTx, 0, txscript.StandardVerifyFlags, txscript.NewSigCache(10),
			txscript.NewTxSigHashes(refundTx), contractTx.TxOut[contractOutPoint.Index].Value)
		if err != nil {
			return nil, err
		}
		err = e.Execute()
		if err != nil {
			return nil, err
		}
	}

	return refundTx, nil
}

func estimateRefundSerializeSize(contract []byte, txOuts []*wire.TxOut) int {
	contractPush, err := txscript.NewScriptBuilder().AddData(contract).Script()
	if err != nil {
		// Should never be hit since this script does exceed the limits.
		panic(err)
	}
	contractPushSize := len(contractPush)

	// 12 additional bytes are for version, locktime and expiry.
	return 12 + wire.VarIntSerializeSize(1) +
		wire.VarIntSerializeSize(uint64(len(txOuts))) +
		inputSize(refundAtomicSwapSigScriptSize+contractPushSize) +
		sumOutputSerializeSizes(txOuts)
}

func createSig(conn bitcoin.Conn, tx *wire.MsgTx, idx int,
	pkScript []byte, addr btcutil.Address) (sig, pubkey []byte, err error) {

	wif, err := conn.Client.DumpPrivKey(addr)
	if err != nil {
		return nil, nil, err
	}
	sig, err = txscript.RawTxInSignature(tx, idx, pkScript, txscript.SigHashAll, wif.PrivKey)
	if err != nil {
		return nil, nil, err
	}
	return sig, wif.PrivKey.PubKey().SerializeCompressed(), nil
}
