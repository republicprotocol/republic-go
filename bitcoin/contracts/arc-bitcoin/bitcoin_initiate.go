package arc_bitcoin

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	rpc "github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
)

func initiate(participantAddress string, value int64, chain string, rpcuser string, rpcpass string, hash []byte, lockTime int64) (err error, result BitcoinData) {
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

	connect, err := normalizeAddress("localhost", walletPort(chainParams))
	if err != nil {
		return fmt.Errorf("wallet server address: %v", err), BitcoinData{}
	}

	connConfig := &rpc.ConnConfig{
		Host:         connect,
		User:         rpcuser,
		Pass:         rpcpass,
		DisableTLS:   true,
		HTTPPostMode: true,
	}

	client, err := rpc.New(connConfig, nil)
	if err != nil {
		return fmt.Errorf("rpc connect: %v", err), BitcoinData{}
	}
	defer func() {
		client.Shutdown()
		client.WaitForShutdown()
	}()

	b, err := buildContract(client, &contractArgs{
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

	if err := promptPublishTx(client, b.contractTx, "contract"); err != nil {
		return err, BitcoinData{}
	}

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
