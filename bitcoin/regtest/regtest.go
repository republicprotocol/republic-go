package regtest

import (
	"os"
	"os/exec"
	"time"

	rpc "github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
)

// Start a local Ganache instance.
func Start() *exec.Cmd {
	cmd := exec.Command("bitcoind", "--regtest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
	return cmd
}

func Mine(rpcClient *rpc.Client) error {

	_, err := rpcClient.Generate(1000)
	if err != nil {
		return err
	}

	tick := time.NewTicker(5 * time.Second)
	select {
	case <-tick.C:
		_, err := rpcClient.Generate(1)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewAccount(rpcClient *rpc.Client, name string, value btcutil.Amount) (btcutil.Address, error) {
	addr, err := rpcClient.GetAccountAddress(name)
	if err != nil {
		return nil, err
	}
	rpcClient.SendToAddress(addr, value)

	_, err = rpcClient.Generate(1)
	if err != nil {
		return nil, err
	}

	return addr, nil
}

// func NewAccount(conn client.Connection, eth *big.Int) (*bind.TransactOpts, common.Address, error) {
// 	ethereumPair, err := crypto.GenerateKey()
// 	if err != nil {
// 		return nil, common.Address{}, err
// 	}
// 	addr := crypto.PubkeyToAddress(ethereumPair.PublicKey)
// 	account := bind.NewKeyedTransactor(ethereumPair)
// 	if eth.Cmp(big.NewInt(0)) > 0 {
// 		if err := DistributeEth(conn, addr); err != nil {
// 			return nil, common.Address{}, err
// 		}
// 	}

// 	return account, addr, nil
// }
