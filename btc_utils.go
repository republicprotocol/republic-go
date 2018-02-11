package main

import (
	"fmt"
	"errors"
	"net"
	"bytes"
	"time"

	"crypto/sha256"
	"crypto/rand"

	"encoding/json"
	"encoding/hex"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcwallet/wallet/txrules"

	"golang.org/x/crypto/ripemd160"

	rpc "github.com/btcsuite/btcd/rpcclient"
)

const txVersion = 2

const secretSize = 32

type builtContract struct {
	contract       []byte
	contractP2SH   btcutil.Address
	contractTxHash *chainhash.Hash
	contractTx     *wire.MsgTx
	contractFee    btcutil.Amount
	refundTx       *wire.MsgTx
	refundFee      btcutil.Amount
}

type contractArgs struct {
	them       *btcutil.AddressPubKeyHash
	amount     btcutil.Amount
	locktime   int64
	secretHash []byte
}

func normalizeAddress(addr string, defaultPort string) (hostport string, err error) {
	host, port, origErr := net.SplitHostPort(addr)
	if origErr == nil {
		return net.JoinHostPort(host, port), nil
	}
	addr = net.JoinHostPort(addr, defaultPort)
	_, _, err = net.SplitHostPort(addr)
	if err != nil {
		return "", origErr
	}
	return addr, nil
}

func buildContract(c *rpc.Client, args *contractArgs, chain string) (*builtContract, error) {
	var chainParams *chaincfg.Params ;
	if (chain == "testnet"){
		chainParams = &chaincfg.TestNet3Params
	} else {
		chainParams = &chaincfg.MainNetParams
	}

	refundAddr, err := getRawChangeAddress(c, chainParams)
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
	contractP2SH, err := btcutil.NewAddressScriptHash(contract, chainParams)
	if err != nil {
		return nil, err
	}
	contractP2SHPkScript, err := txscript.PayToAddrScript(contractP2SH)
	if err != nil {
		return nil, err
	}

	feePerKb, minFeePerKb, err := getFeePerKb(c)
	if err != nil {
		return nil, err
	}

	unsignedContract := wire.NewMsgTx(txVersion)
	unsignedContract.AddTxOut(wire.NewTxOut(int64(args.amount), contractP2SHPkScript))
	unsignedContract, contractFee, err := fundRawTransaction(c, unsignedContract, feePerKb)
	if err != nil {
		return nil, fmt.Errorf("fundrawtransaction: %v", err)
	}
	contractTx, complete, err := c.SignRawTransaction(unsignedContract)
	if err != nil {
		return nil, fmt.Errorf("signrawtransaction: %v", err)
	}
	if !complete {
		return nil, errors.New("signrawtransaction: failed to completely sign contract transaction")
	}

	contractTxHash := contractTx.TxHash()

	refundTx, refundFee, err := buildRefund(c, contract, contractTx, feePerKb, minFeePerKb)
	if err != nil {
		return nil, err
	}

	return &builtContract{
		contract,
		contractP2SH,
		&contractTxHash,
		contractTx,
		contractFee,
		refundTx,
		refundFee,
	}, nil
}


func walletPort(params *chaincfg.Params) string {
	switch params {
	case &chaincfg.MainNetParams:
		return "8332"
	case &chaincfg.TestNet3Params:
		return "18332"
	default:
		return ""
	}
}

func sha256Hash(x []byte) []byte {
	h := sha256.Sum256(x)
	return h[:]
}

func (cmd *extractSecretCmd) runCommand(c *rpc.Client) error {
	return cmd.runOfflineCommand()
}

// atomicSwapContract returns an output script that may be redeemed by one of
// two signature scripts:
//
//   <their sig> <their pubkey> <initiator secret> 1
//
//   <my sig> <my pubkey> 0
//
// The first signature script is the normal redemption path done by the other
// party and requires the initiator's secret.  The second signature script is
// the refund path performed by us, but the refund can only be performed after
// locktime.
func atomicSwapContract(pkhMe, pkhThem *[ripemd160.Size]byte, locktime int64, secretHash []byte) ([]byte, error) {
	b := txscript.NewScriptBuilder()

	b.AddOp(txscript.OP_IF) // Normal redeem path
	{
		// Require initiator's secret to be known to redeem the output.
		b.AddOp(txscript.OP_SHA256)
		b.AddData(secretHash)
		b.AddOp(txscript.OP_EQUALVERIFY)

		// Verify their signature is being used to redeem the output.  This
		// would normally end with OP_EQUALVERIFY OP_CHECKSIG but this has been
		// moved outside of the branch to save a couple bytes.
		b.AddOp(txscript.OP_DUP)
		b.AddOp(txscript.OP_HASH160)
		b.AddData(pkhThem[:])
	}
	b.AddOp(txscript.OP_ELSE) // Refund path
	{
		// Verify locktime and drop it off the stack (which is not done by
		// CLTV).
		b.AddInt64(locktime)
		b.AddOp(txscript.OP_CHECKLOCKTIMEVERIFY)
		b.AddOp(txscript.OP_DROP)

		// Verify our signature is being used to redeem the output.  This would
		// normally end with OP_EQUALVERIFY OP_CHECKSIG but this has been moved
		// outside of the branch to save a couple bytes.
		b.AddOp(txscript.OP_DUP)
		b.AddOp(txscript.OP_HASH160)
		b.AddData(pkhMe[:])
	}
	b.AddOp(txscript.OP_ENDIF)

	// Complete the signature check.
	b.AddOp(txscript.OP_EQUALVERIFY)
	b.AddOp(txscript.OP_CHECKSIG)

	return b.Script()
}


func getRawChangeAddress(c *rpc.Client, chainParams *chaincfg.Params) (btcutil.Address, error) {
	rawResp, err := c.RawRequest("getrawchangeaddress", nil)
	if err != nil {
		return nil, err
	}
	var addrStr string
	err = json.Unmarshal(rawResp, &addrStr)
	if err != nil {
		return nil, err
	}
	addr, err := btcutil.DecodeAddress(addrStr, chainParams)
	if err != nil {
		return nil, err
	}
	if !addr.IsForNet(chainParams) {
		return nil, fmt.Errorf("address %v is not intended for use on %v",
			addrStr, chainParams.Name)
	}
	return addr, nil
}

func getFeePerKb(c *rpc.Client) (useFee, relayFee btcutil.Amount, err error) {
	info, err := c.GetInfo()
	if err != nil {
		return 0, 0, fmt.Errorf("getinfo: %v", err)
	}
	relayFee, err = btcutil.NewAmount(info.RelayFee)
	if err != nil {
		return 0, 0, err
	}
	maxFee := info.PaytxFee
	if info.PaytxFee != 0 {
		if info.RelayFee > maxFee {
			maxFee = info.RelayFee
		}
		useFee, err = btcutil.NewAmount(maxFee)
		return useFee, relayFee, err
	}

	params := []json.RawMessage{[]byte("6")}
	estimateRawResp, err := c.RawRequest("estimatefee", params)
	if err != nil {
		return 0, 0, err
	}
	var estimateResp float64 = -1
	err = json.Unmarshal(estimateRawResp, &estimateResp)
	if err == nil && estimateResp != -1 {
		useFee, err = btcutil.NewAmount(estimateResp)
		if relayFee > useFee {
			useFee = relayFee
		}
		return useFee, relayFee, err
	}

	fmt.Println("warning: falling back to mempool relay fee policy")
	useFee, err = btcutil.NewAmount(info.RelayFee)
	return useFee, relayFee, err
}

func fundRawTransaction(c *rpc.Client, tx *wire.MsgTx, feePerKb btcutil.Amount) (fundedTx *wire.MsgTx, fee btcutil.Amount, err error) {
	var buf bytes.Buffer
	buf.Grow(tx.SerializeSize())
	tx.Serialize(&buf)
	param0, err := json.Marshal(hex.EncodeToString(buf.Bytes()))
	if err != nil {
		return nil, 0, err
	}
	param1, err := json.Marshal(struct {
		FeeRate float64 `json:"feeRate"`
	}{
		FeeRate: feePerKb.ToBTC(),
	})
	if err != nil {
		return nil, 0, err
	}
	params := []json.RawMessage{param0, param1}
	rawResp, err := c.RawRequest("fundrawtransaction", params)
	if err != nil {
		return nil, 0, err
	}
	var resp struct {
		Hex       string  `json:"hex"`
		Fee       float64 `json:"fee"`
		ChangePos float64 `json:"changepos"`
	}
	err = json.Unmarshal(rawResp, &resp)
	if err != nil {
		return nil, 0, err
	}
	fundedTxBytes, err := hex.DecodeString(resp.Hex)
	if err != nil {
		return nil, 0, err
	}
	fundedTx = &wire.MsgTx{}
	err = fundedTx.Deserialize(bytes.NewReader(fundedTxBytes))
	if err != nil {
		return nil, 0, err
	}
	feeAmount, err := btcutil.NewAmount(resp.Fee)
	if err != nil {
		return nil, 0, err
	}
	return fundedTx, feeAmount, nil
}

func calcFeePerKb(absoluteFee btcutil.Amount, serializeSize int) float64 {
	return float64(absoluteFee) / float64(serializeSize) / 1e5
}

func promptPublishTx(c *rpc.Client, tx *wire.MsgTx, name string) error {
		txHash, err := c.SendRawTransaction(tx, false)
		if err != nil {
			return fmt.Errorf("sendrawtransaction: %v", err)
		}
		fmt.Printf("Published %s transaction (%v)\n", name, txHash)
		return nil
}