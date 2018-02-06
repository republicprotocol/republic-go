package btc

import (
	"fmt"
	"errors"
	"net"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	rpc "github.com/btcsuite/btcd/rpcclient"
)

type initiateCmd struct {
	cp2Addr *btcutil.AddressPubKeyHash
	amount  btcutil.Amount
}

type redeemCmd struct {
	contract   []byte
	contractTx *wire.MsgTx
	secret     []byte
}

type refundCmd struct {
	contract   []byte
	contractTx *wire.MsgTx
}

type extractSecretCmd struct {
	redemptionTx *wire.MsgTx
	secretHash   []byte
}

type auditContractCmd struct {
	contract   []byte
	contractTx *wire.MsgTx
}

func Open(participantAddress string, value float64, chain string) (err error, showUsage bool) {
	var chainParams *chaincfg.Params ;
	if (chain == "testnet"){
		chainParams = &chaincfg.TestNet3Params
	} else {
		chainParams = &chaincfg.MainNetParams
	}

	cp2Addr, err := btcutil.DecodeAddress(participantAddress, chainParams)
	if err != nil {
		return fmt.Errorf("failed to decode participant address: %v", err), true
	}
	if !cp2Addr.IsForNet(chainParams) {
		return fmt.Errorf("participant address is not "+
			"intended for use on %v", chainParams.Name), true
	}
	cp2AddrP2PKH, ok := cp2Addr.(*btcutil.AddressPubKeyHash)
	if !ok {
		return errors.New("participant address is not P2PKH"), true
	}

	amount, err := btcutil.NewAmount(value)
	if err != nil {
		return err, true
	}
	cmd := &initiateCmd{cp2Addr: cp2AddrP2PKH, amount: amount}


	connect, err := normalizeAddress(*connectFlag, walletPort(chainParams))
	if err != nil {
		return fmt.Errorf("wallet server address: %v", err), true
	}

	connConfig := &rpc.ConnConfig{
		Host:         connect,
		User:         *rpcuserFlag,
		Pass:         *rpcpassFlag,
		DisableTLS:   true,
		HTTPPostMode: true,
	}

	client, err := rpc.New(connConfig, nil)
	if err != nil {
		return fmt.Errorf("rpc connect: %v", err), false
	}
	defer func() {
		client.Shutdown()
		client.WaitForShutdown()
	}()

	err = cmd.runCommand(client)
	return err, false
}

func Close(contract string, contractTx string, secret []byte, chain string) {
	var chainParams *chaincfg.Params ;
	if (chain == "testnet"){
		chainParams = &chaincfg.TestNet3Params
	} else {
		chainParams = &chaincfg.MainNetParams
	}

	contract, err := hex.DecodeString(contract)
		if err != nil {
			return fmt.Errorf("failed to decode contract: %v", err), true
		}

	contractTxBytes, err := hex.DecodeString(contractTx)
		if err != nil {
			return fmt.Errorf("failed to decode contract transaction: %v", err), true
		}

	var contractTx wire.MsgTx
		err = contractTx.Deserialize(bytes.NewReader(contractTxBytes))
		if err != nil {
			return fmt.Errorf("failed to decode contract transaction: %v", err), true
		}

	secret, err := hex.DecodeString(secret)
		if err != nil {
			return fmt.Errorf("failed to decode secret: %v", err), true
		}

	cmd = &redeemCmd{contract: contract, contractTx: &contractTx, secret: secret}

	connect, err := normalizeAddress(*connectFlag, walletPort(chainParams))
	if err != nil {
		return fmt.Errorf("wallet server address: %v", err), true
	}

	connConfig := &rpc.ConnConfig{
		Host:         connect,
		User:         *rpcuserFlag,
		Pass:         *rpcpassFlag,
		DisableTLS:   true,
		HTTPPostMode: true,
	}

	client, err := rpc.New(connConfig, nil)
	if err != nil {
		return fmt.Errorf("rpc connect: %v", err), false
	}
	defer func() {
		client.Shutdown()
		client.WaitForShutdown()
	}()

	err = cmd.runCommand(client)
	return err, false
}

func Expire(contract string, contractTx string, chain string) {
	var chainParams *chaincfg.Params ;
	if (chain == "testnet"){
		chainParams = &chaincfg.TestNet3Params
	} else {
		chainParams = &chaincfg.MainNetParams
	}
	
	contract, err := hex.DecodeString(contract)
		if err != nil {
			return fmt.Errorf("failed to decode contract: %v", err), true
		}

		contractTxBytes, err := hex.DecodeString(contractTx)
		if err != nil {
			return fmt.Errorf("failed to decode contract transaction: %v", err), true
		}
		var contractTx wire.MsgTx
		err = contractTx.Deserialize(bytes.NewReader(contractTxBytes))
		if err != nil {
			return fmt.Errorf("failed to decode contract transaction: %v", err), true
		}

		cmd = &refundCmd{contract: contract, contractTx: &contractTx}

		connect, err := normalizeAddress(*connectFlag, walletPort(chainParams))
	if err != nil {
		return fmt.Errorf("wallet server address: %v", err), true
	}

	connConfig := &rpc.ConnConfig{
		Host:         connect,
		User:         *rpcuserFlag,
		Pass:         *rpcpassFlag,
		DisableTLS:   true,
		HTTPPostMode: true,
	}

	client, err := rpc.New(connConfig, nil)
	if err != nil {
		return fmt.Errorf("rpc connect: %v", err), false
	}
	defer func() {
		client.Shutdown()
		client.WaitForShutdown()
	}()

	err = cmd.runCommand(client)
	return err, false
}

func Validate(contract string, contractTx string, chain string) {
	var chainParams *chaincfg.Params ;
	if (chain == "testnet"){
		chainParams = &chaincfg.TestNet3Params
	} else {
		chainParams = &chaincfg.MainNetParams
	}

	contract, err := hex.DecodeString(contract)
		if err != nil {
			return fmt.Errorf("failed to decode contract: %v", err), true
		}

		contractTxBytes, err := hex.DecodeString(contractTx)
		if err != nil {
			return fmt.Errorf("failed to decode contract transaction: %v", err), true
		}
		var contractTx wire.MsgTx
		err = contractTx.Deserialize(bytes.NewReader(contractTxBytes))
		if err != nil {
			return fmt.Errorf("failed to decode contract transaction: %v", err), true
		}

		cmd = &auditContractCmd{contract: contract, contractTx: &contractTx}
		connect, err := normalizeAddress(*connectFlag, walletPort(chainParams))
		if err != nil {
			return fmt.Errorf("wallet server address: %v", err), true
		}
	
		connConfig := &rpc.ConnConfig{
			Host:         connect,
			User:         *rpcuserFlag,
			Pass:         *rpcpassFlag,
			DisableTLS:   true,
			HTTPPostMode: true,
		}
	
		client, err := rpc.New(connConfig, nil)
		if err != nil {
			return fmt.Errorf("rpc connect: %v", err), false
		}
		defer func() {
			client.Shutdown()
			client.WaitForShutdown()
		}()
	
		err = cmd.runCommand(client)
		return err, false
}

func RetrieveSecretKey(redemptionTx string, secretHash []byte, chain string) {
	var chainParams *chaincfg.Params ;
	if (chain == "testnet"){
		chainParams = &chaincfg.TestNet3Params
	} else {
		chainParams = &chaincfg.MainNetParams
	}

	redemptionTxBytes, err := hex.DecodeString(redemptionTx)
		if err != nil {
			return fmt.Errorf("failed to decode redemption transaction: %v", err), true
		}
		var redemptionTx wire.MsgTx
		err = redemptionTx.Deserialize(bytes.NewReader(redemptionTxBytes))
		if err != nil {
			return fmt.Errorf("failed to decode redemption transaction: %v", err), true
		}

		if len(secretHash) != ripemd160.Size {
			return errors.New("secret hash has wrong size"), true
		}

		cmd = &extractSecretCmd{redemptionTx: &redemptionTx, secretHash: secretHash}

		connect, err := normalizeAddress(*connectFlag, walletPort(chainParams))
		if err != nil {
			return fmt.Errorf("wallet server address: %v", err), true
		}
	
		connConfig := &rpc.ConnConfig{
			Host:         connect,
			User:         *rpcuserFlag,
			Pass:         *rpcpassFlag,
			DisableTLS:   true,
			HTTPPostMode: true,
		}
	
		client, err := rpc.New(connConfig, nil)
		if err != nil {
			return fmt.Errorf("rpc connect: %v", err), false
		}
		defer func() {
			client.Shutdown()
			client.WaitForShutdown()
		}()
	
		err = cmd.runCommand(client)
		return err, false
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
	default:
		return ""
	}
}
