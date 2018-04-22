package arc_bitcoin

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/republicprotocol/republic-go/bitcoin/client"
)

func initiate(participantAddress string, value int64, chain string, rpcUser string, rpcPass string, hash []byte, lockTime int64) (err error, result BitcoinData) {
	var chainParams *chaincfg.Params
	if chain == "regtest" {
		chainParams = &chaincfg.RegressionNetParams
	} else if chain == "testnet" {
		chainParams = &chaincfg.TestNet3Params
	} else {
		chainParams = &chaincfg.MainNetParams
	}

	cp2Addr, err := btcutil.DecodeAddress(participantAddress, chainParams)
	if err != nil {
		return fmt.Errorf("failed to decode participant address: %v", err), BitcoinData{}
	}
	if !cp2Addr.IsForNet(chainParams) {
		return fmt.Errorf("participant address is not "+
			"intended for use on %v", chainParams.Name), BitcoinData{}
	}
	cp2AddrP2PKH, ok := cp2Addr.(*btcutil.AddressPubKeyHash)
	if !ok {
		return errors.New("participant address is not P2PKH"), BitcoinData{}
	}

	rpcClient, err := client.ConnectToRPC(chainParams, rpcUser, rpcPass)
	if err != nil {
		return err, BitcoinData{}
	}
	defer func() {
		rpcClient.Shutdown()
		rpcClient.WaitForShutdown()
	}()

	b, err := buildContract(rpcClient, &contractArgs{
		them:       cp2AddrP2PKH,
		amount:     value,
		locktime:   lockTime,
		secretHash: hash,
	}, chain)
	if err != nil {
		return err, BitcoinData{}
	}

	var contractBuf bytes.Buffer
	contractBuf.Grow(b.contractTx.SerializeSize())
	b.contractTx.Serialize(&contractBuf)

	var refundBuf bytes.Buffer
	refundBuf.Grow(b.refundTx.SerializeSize())
	b.refundTx.Serialize(&refundBuf)

	txHash, err := client.PromptPublishTx(rpcClient, b.contractTx, "contract")

	if err != nil {
		return err, BitcoinData{}
	}

	client.WaitForConfirmations(rpcClient, txHash, 1)

	refundTx := *b.refundTx
	return nil, BitcoinData{
		contract:       b.contract,
		contractHash:   b.contractP2SH.EncodeAddress(),
		contractTx:     contractBuf.Bytes(),
		contractTxHash: b.contractTxHash.CloneBytes(),
		refundTx:       refundBuf.Bytes(),
		refundTxHash:   refundTx.TxHash(),
	}
}
