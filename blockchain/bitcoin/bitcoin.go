package bitcoin

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
)

type Conn struct {
	Client      *rpcclient.Client
	ChainParams *chaincfg.Params
}

func Connect(chain string, rpcUser, rpcPass string) (Conn, error) {
	var chainParams *chaincfg.Params
	switch chain {
	case "regtest":
		chainParams = &chaincfg.RegressionNetParams
	case "testnet":
		chainParams = &chaincfg.TestNet3Params
	default:
		chainParams = &chaincfg.MainNetParams
	}

	connect, err := normalizeAddress("localhost", walletPort(chainParams))
	if err != nil {
		return Conn{}, fmt.Errorf("wallet server address: %v", err)
	}

	connConfig := &rpcclient.ConnConfig{
		Host:         connect,
		User:         rpcUser,
		Pass:         rpcPass,
		DisableTLS:   true,
		HTTPPostMode: true,
	}

	rpcClient, err := rpcclient.New(connConfig, nil)
	if err != nil {
		return Conn{}, fmt.Errorf("rpc connect: %v", err)
	}

	// Should call the following after this function:
	/*
		defer func() {
			rpcClient.Shutdown()
			pcClient.WaitForShutdown()
		}()
	*/

	return Conn{
		Client:      rpcClient,
		ChainParams: chainParams,
	}, nil
}

func (conn Conn) FundRawTransaction(tx *wire.MsgTx) (fundedTx *wire.MsgTx, err error) {
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
	rawResp, err := conn.Client.RawRequest("fundrawtransaction", params)
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

func (conn Conn) PromptPublishTx(tx *wire.MsgTx, name string) (*chainhash.Hash, error) {
	txHash, err := conn.Client.SendRawTransaction(tx, false)
	if err != nil {
		return nil, fmt.Errorf("sendrawtransaction: %v", err)
	}
	return txHash, nil
}

func (conn Conn) WaitForConfirmations(txHash *chainhash.Hash, requiredConfirmations int64) error {
	confirmations := int64(0)
	for confirmations < requiredConfirmations {
		txDetails, err := conn.Client.GetTransaction(txHash)
		if err != nil {
			return err
		}
		confirmations = txDetails.Confirmations

		// TODO: Base delay on chain config
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (conn Conn) Shutdown() {
	conn.Client.Shutdown()
	conn.Client.WaitForShutdown()
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
