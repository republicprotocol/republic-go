package arc_bitcoin

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	rpc "github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func refund(contract, contractTxBytes []byte, chain, rpcUser, rpcPass string) (Error error) {
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
		return fmt.Errorf("failed to decode contract transaction: %v", err)
	}

	connect, err := normalizeAddress("localhost", walletPort(chainParams))
	if err != nil {
		return fmt.Errorf("wallet server address: %v", err)
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
		return fmt.Errorf("rpc connect: %v", err)
	}
	defer func() {
		client.Shutdown()
		client.WaitForShutdown()
	}()

	pushes, err := txscript.ExtractAtomicSwapDataPushes(0, contract)
	if err != nil {
		return err
	}
	if pushes == nil {
		return errors.New("contract is not an atomic swap script recognized by this tool")
	}

	refundTx, err := buildRefund(client, contract, &contractTx, chainParams)
	if err != nil {
		return err
	}

	return promptPublishTx(client, refundTx, "refund")
}
