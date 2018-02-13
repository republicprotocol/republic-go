package bitcoin

import (
	"fmt"
	"errors"
	"bytes"

	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	rpc "github.com/btcsuite/btcd/rpcclient"
)

func Open(participantAddress string, value int64, chain string, rpcuser string, rpcpass string, hash []byte, lockTime int64) (err error, result openResult) {
	var chainParams *chaincfg.Params ;
	if (chain == "testnet"){
		chainParams = &chaincfg.TestNet3Params
	} else {
		chainParams = &chaincfg.MainNetParams
	}

	cp2Addr, err := btcutil.DecodeAddress(participantAddress, chainParams)
	if err != nil {
		return fmt.Errorf("failed to decode participant address: %v", err), openResult{}
	}
	if !cp2Addr.IsForNet(chainParams) {
		return fmt.Errorf("participant address is not "+
			"intended for use on %v", chainParams.Name), openResult{}
	}
	cp2AddrP2PKH, ok := cp2Addr.(*btcutil.AddressPubKeyHash)
	if !ok {
		return errors.New("participant address is not P2PKH"), openResult{}
	}

	connect, err := normalizeAddress("localhost", walletPort(chainParams))
	if err != nil {
		return fmt.Errorf("wallet server address: %v", err), openResult{}
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
		return fmt.Errorf("rpc connect: %v", err), openResult{}
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
		return err, openResult{}
	}

	var contractBuf bytes.Buffer
	contractBuf.Grow(b.contractTx.SerializeSize())
	b.contractTx.Serialize(&contractBuf)

	var refundBuf bytes.Buffer
	refundBuf.Grow(b.refundTx.SerializeSize())
	b.refundTx.Serialize(&refundBuf)

	if err := promptPublishTx(c, b.contractTx, "contract"); err != nil {
		return err, openResult{}
	}

	return nil, openResult{
		contract: b.contract,
		contractHash: b.contractP2SH.EncodeAddress(),
		contractTx: contractBuf.Bytes(),
		contractTxHash: b.contractTxHash,
		refundTx: refundBuf.Bytes(),
		refundTxHash: &b.refundTx.TxHash(),
	}
}