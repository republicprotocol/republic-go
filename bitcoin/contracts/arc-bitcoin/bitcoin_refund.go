package arc_bitcoin

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/republicprotocol/republic-go/bitcoin/client"
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

	rpcClient, err := client.ConnectToRPC(chainParams, rpcUser, rpcPass)
	if err != nil {
		return err
	}
	defer func() {
		rpcClient.Shutdown()
		rpcClient.WaitForShutdown()
	}()

	pushes, err := txscript.ExtractAtomicSwapDataPushes(0, contract)
	if err != nil {
		return err
	}
	if pushes == nil {
		return errors.New("contract is not an atomic swap script recognized by this tool")
	}

	refundTx, err := buildRefund(rpcClient, contract, &contractTx, chainParams)
	if err != nil {
		return err
	}

	txHash, err := client.PromptPublishTx(rpcClient, refundTx, "refund")
	if err != nil {
		return err
	}

	client.WaitForConfirmations(rpcClient, txHash, 1)

	return nil
}
