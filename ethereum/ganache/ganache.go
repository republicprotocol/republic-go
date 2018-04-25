package ganache

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/republicprotocol/republic-go/ethereum/bindings"
	"github.com/republicprotocol/republic-go/ethereum/client"
)

var genesisPrivateKey, genesisTransactor = genesis()

// GenesisPrivateKey used by Ganache.
func GenesisPrivateKey() *ecdsa.PrivateKey {
	return genesisPrivateKey
}

// GenesisTransactor used by Ganache.
func GenesisTransactor() *bind.TransactOpts {
	return genesisTransactor
}

// Start a local Ganache instance.
func Start() *exec.Cmd {
	cmd := exec.Command("ganache-cli", fmt.Sprintf("--account=0x%x,1000000000000000000000", crypto.FromECDSA(genesisPrivateKey)))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	go func() {
		<-signals
		fmt.Printf("ganache is shutting down...\n")
		if err := cmd.Process.Kill(); err != nil {
			fmt.Printf("ganache cannot shutdown: %v", err)
			return
		}
		if err := cmd.Wait(); err != nil {
			fmt.Printf("ganache cannot shutdown: %v", err)
			return
		}
		fmt.Printf("ganache shutdown")
	}()

	// Wait for ganache to boot
	time.Sleep(2 * time.Second)

	return cmd
}

// Connect to a local Ganache instance.
func Connect(ganacheRPC string) (client.Connection, error) {
	ethclient, err := ethclient.Dial(ganacheRPC)
	if err != nil {
		return client.Connection{}, err
	}
	return client.Connection{
		Client:     ethclient,
		DNRAddress: client.DarkNodeRegistryAddressOnGanache,
		RenAddress: client.RepublicTokenAddressOnGanache,
		Network:    client.NetworkGanache,
	}, nil
}

// StartAndConnect to a local Ganache instance and deploy all smart contracts.
func StartAndConnect() (*exec.Cmd, client.Connection, error) {
	cmd := Start()
	time.Sleep(time.Second)
	conn, err := Connect("http://localhost:8545")
	if err != nil {
		return cmd, conn, err
	}
	if err := DeployContracts(conn); err != nil {
		return cmd, conn, err
	}
	return cmd, conn, nil
}

// DeployContracts to Ganache deploys REN and DNR contracts using the genesis private key
func DeployContracts(conn client.Connection) error {
	return deployContracts(conn, genesisTransactor)
}

// DistributeEth transfers ETH to each of the addresses
func DistributeEth(conn client.Connection, addresses ...common.Address) error {

	for _, address := range addresses {
		err := conn.TransferEth(context.Background(), genesisTransactor, address, big.NewInt(1000000000000000000))
		if err != nil {
			return err
		}
	}

	return nil
}

// DistributeREN transfers REN to each of the addresses
func DistributeREN(conn client.Connection, addresses ...common.Address) error {

	republicTokenContract, err := bindings.NewRepublicToken(conn.RenAddress, bind.ContractBackend(conn.Client))
	if err != nil {
		return err
	}

	for _, address := range addresses {
		if _, err := republicTokenContract.Transfer(genesisTransactor, address, big.NewInt(0).Mul(big.NewInt(100), big.NewInt(1000000000000000000))); err != nil {
			return err
		}
	}

	return nil
}

func NewAccount(conn client.Connection, eth *big.Int) (*bind.TransactOpts, common.Address, error) {
	ethereumPair, err := crypto.GenerateKey()
	if err != nil {
		return nil, common.Address{}, err
	}
	addr := crypto.PubkeyToAddress(ethereumPair.PublicKey)
	account := bind.NewKeyedTransactor(ethereumPair)
	if eth.Cmp(big.NewInt(0)) > 0 {
		if err := DistributeEth(conn, addr); err != nil {
			return nil, common.Address{}, err
		}
	}

	return account, addr, nil
}

func genesis() (*ecdsa.PrivateKey, *bind.TransactOpts) {
	deployerKey, err := crypto.HexToECDSA("2aba04ee8a322b8648af2a784144181a0c793f1a2e80519418f3d20bbfb22249")
	if err != nil {
		log.Fatalf("cannot read genesis key: %v", err)
		return nil, nil
	}
	deployerAuth := bind.NewKeyedTransactor(deployerKey)
	return deployerKey, deployerAuth
}

func deployContracts(conn client.Connection, transactor *bind.TransactOpts) error {

	_, republicTokenAddress, err := deployRepublicToken(context.Background(), conn, transactor)
	if err != nil {
		return err
	}

	_, darkNodeRegistryAddress, err := deployDarkNodeRegistry(context.Background(), conn, transactor, republicTokenAddress)
	if err != nil {
		return err
	}

	if republicTokenAddress != client.RepublicTokenAddressOnGanache {
		return fmt.Errorf("RepublicToken address has changed: expected: %s, got: %s", client.RepublicTokenAddressOnGanache.Hex(), republicTokenAddress.Hex())
	}
	if darkNodeRegistryAddress != client.DarkNodeRegistryAddressOnGanache {
		return fmt.Errorf("DarkNodeRegistry address has changed: expected: %s, got: %s", client.DarkNodeRegistryAddressOnGanache.Hex(), darkNodeRegistryAddress.Hex())
	}

	return nil
}

func deployRepublicToken(ctx context.Context, conn client.Connection, auth *bind.TransactOpts) (*bindings.RepublicToken, common.Address, error) {
	address, tx, ren, err := bindings.DeployRepublicToken(auth, conn.Client)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("cannot deploy RepublicToken: %v", err)
	}
	conn.PatchedWaitDeployed(ctx, tx)
	return ren, address, nil
}

func deployDarkNodeRegistry(ctx context.Context, conn client.Connection, auth *bind.TransactOpts, republicTokenAddress common.Address) (*bindings.DarkNodeRegistry, common.Address, error) {
	// 0 REN
	minimumBond := big.NewInt(0)
	// 1 second
	minimumEpochInterval := big.NewInt(1)
	// 24 Darknode in a pool
	minimumPoolSize := big.NewInt(24)

	address, tx, dnr, err := bindings.DeployDarkNodeRegistry(auth, conn.Client, republicTokenAddress, minimumBond, minimumPoolSize, minimumEpochInterval)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("cannot deploy DarkNodeRegistry: %v", err)
	}
	conn.PatchedWaitDeployed(ctx, tx)
	return dnr, address, nil
}
