package bitcoind

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/btcsuite/btcutil"
	"github.com/republicprotocol/republic-go/blockchain/bitcoin"
)

// Start a local Ganache instance.
func Start() *exec.Cmd {
	cmd := exec.Command("bitcoind", "--regtest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
	return cmd
}

func Mine(conn bitcoin.Conn) error {

	fmt.Println("initial 100")
	_, err := conn.Client.Generate(100)
	if err != nil {
		return err
	}
	fmt.Println("100 done")

	tick := time.NewTicker(2 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			_, err := conn.Client.Generate(1)
			if err != nil {
				return err
			}
		}
	}
}

func NewAccount(conn bitcoin.Conn, name string, value btcutil.Amount) (btcutil.Address, error) {
	addr, err := conn.Client.GetAccountAddress(name)
	if err != nil {
		return nil, err
	}

	if value > 0 {
		_, err = conn.Client.SendToAddress(addr, value)
		if err != nil {
			return nil, err
		}

		_, err = conn.Client.Generate(1)
		if err != nil {
			return nil, err
		}
	}

	return addr, nil
}
