package client

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	rpc "github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
)

func ConnectToRPC(chainParams *chaincfg.Params, rpcUser, rpcPass string) (*rpc.Client, error) {
	connect, err := normalizeAddress("localhost", walletPort(chainParams))
	if err != nil {
		return nil, fmt.Errorf("wallet server address: %v", err)
	}

	connConfig := &rpc.ConnConfig{
		Host:         connect,
		User:         rpcUser,
		Pass:         rpcPass,
		DisableTLS:   true,
		HTTPPostMode: true,
	}

	rpcClient, err := rpc.New(connConfig, nil)
	if err != nil {
		return nil, fmt.Errorf("rpc connect: %v", err)
	}

	// Should call the following after this function:
	/*
		defer func() {
			rpcClient.Shutdown()
			pcClient.WaitForShutdown()
		}()
	*/

	return rpcClient, nil
}

func FundRawTransaction(c *rpc.Client, tx *wire.MsgTx) (fundedTx *wire.MsgTx, err error) {
	var buf bytes.Buffer
	buf.Grow(tx.SerializeSize())
	tx.Serialize(&buf)
	param0, err := json.Marshal(hex.EncodeToString(buf.Bytes()))
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	params := []json.RawMessage{param0}
	rawResp, err := c.RawRequest("fundrawtransaction", params)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Hex       string  `json:"hex"`
		Fee       float64 `json:"fee"`
		ChangePos float64 `json:"changepos"`
	}
	err = json.Unmarshal(rawResp, &resp)
	if err != nil {
		return nil, err
	}
	fundedTxBytes, err := hex.DecodeString(resp.Hex)
	if err != nil {
		return nil, err
	}
	fundedTx = &wire.MsgTx{}
	err = fundedTx.Deserialize(bytes.NewReader(fundedTxBytes))
	if err != nil {
		return nil, err
	}
	return fundedTx, nil
}

func PromptPublishTx(c *rpc.Client, tx *wire.MsgTx, name string) (*chainhash.Hash, error) {
	txHash, err := c.SendRawTransaction(tx, false)
	if err != nil {
		return nil, fmt.Errorf("sendrawtransaction: %v", err)
	}
	return txHash, nil
}

func WaitForConfirmations(c *rpc.Client, txHash *chainhash.Hash, requiredConfirmations int64) error {
	confirmations := int64(0)
	for confirmations < requiredConfirmations {
		txDetails, err := c.GetTransaction(txHash)
		if err != nil {
			return err
		}
		confirmations = txDetails.Confirmations

		// TODO: Base delay on chain config
		time.Sleep(1 * time.Second)
	}
	return nil
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

func walletPort(params *chaincfg.Params) string {
	switch params {
	case &chaincfg.MainNetParams:
		return "8332"
	case &chaincfg.TestNet3Params:
		return "18332"
	case &chaincfg.RegressionNetParams:
		return "18443"
	default:
		return ""
	}
}
