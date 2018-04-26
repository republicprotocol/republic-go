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
	"sync"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/bindings"
	"github.com/republicprotocol/republic-go/blockchain/test"
)

var genesisPrivateKey, genesisTransactor = genesis()

// GenesisPrivateKey used by Ganache.
func GenesisPrivateKey() ecdsa.PrivateKey {
	return *genesisPrivateKey
}

// GenesisTransactor used by Ganache.
func GenesisTransactor() bind.TransactOpts {
	return *genesisTransactor
}

// WatchForInterrupt will stop Ganache upon receiving receiving a interrupt signal
func WatchForInterrupt() {
	cmd := globalGanacheCmd
	signals := make(chan os.Signal, 1)
	defer close(signals)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	<-signals
	fmt.Println("ganache is shutting down...")
	if err := cmd.Process.Kill(); err != nil {
		fmt.Printf("ganache cannot shutdown: %v", err)
		return
	}
	if err := cmd.Wait(); err != nil {
		fmt.Printf("ganache cannot shutdown: %v", err)
		return
	}
	fmt.Printf("ganache shutdown")
}

var globalGanacheMu = &sync.Mutex{}
var globalGanacheCounter uint64
var globalGanacheCmd *exec.Cmd
var globalGanacheSnapshot string
var globalGanacheStopSchedule time.Time

// Start a local Ganache instance.
func Start() bool {
	globalGanacheMu.Lock()
	defer globalGanacheMu.Unlock()

	if globalGanacheCounter > 0 {
		globalGanacheCounter++
		return false
	}

	cmd := exec.Command("ganache-cli", fmt.Sprintf("--account=0x%x,1000000000000000000000", crypto.FromECDSA(genesisPrivateKey)))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()

	go WatchForInterrupt()

	// Wait for ganache to boot
	var delay time.Duration
	if test.GetCIEnv() {
		delay = 10 * time.Second
	} else {
		delay = 4 * time.Second
	}
	time.Sleep(delay)

	globalGanacheCounter++
	globalGanacheCmd = cmd

	return true
}

func Wait() error {
	return globalGanacheCmd.Wait()
}

// Stop will kill the local Ganache instance
func Stop() {
	globalGanacheMu.Lock()
	defer globalGanacheMu.Unlock()
	if globalGanacheCounter > 1 {
		globalGanacheCounter--
		return
	}

	globalGanacheCounter--

	// Stop ganache
	if err := globalGanacheCmd.Process.Kill(); err != nil {
		panic(err)
	}
	globalGanacheCmd.Wait()
}

// Connect to a local Ganache instance.
func Connect(ganacheRPC string) (ethereum.Conn, error) {
	raw, err := rpc.DialContext(context.Background(), ganacheRPC)
	if err != nil {
		return ethereum.Conn{}, err
	}
	ethclient := ethclient.NewClient(raw)

	return ethereum.Conn{
		RawClient:               raw,
		Client:                  ethclient,
		Network:                 ethereum.NetworkGanache,
		DarknodeRegistryAddress: ethereum.DarknodeRegistryAddressOnGanache,
		RepublicTokenAddress:    ethereum.RepublicTokenAddressOnGanache,
	}, nil
}

// StartAndConnect to a local Ganache instance and deploy all smart contracts.
func StartAndConnect() (ethereum.Conn, error) {
	firstConnection := Start()

	conn, err := Connect("http://localhost:8545")
	if err != nil {
		return conn, err
	}

	if firstConnection {
		// Deploy contracts and take snapshot

		if err := DeployContracts(conn); err != nil {
			return conn, err
		}
		snapshot, err := snapshot(conn)
		if err != nil {
			return conn, err
		}
		globalGanacheSnapshot = snapshot
	} else {
		// Roll back to snapshot
		if err := revertToSnapshot(conn, globalGanacheSnapshot); err != nil {
			return conn, err
		}

		// Take snapshot again -- avoid problem with Ganache freezing
		snapshot, err := snapshot(conn)
		if err != nil {
			return conn, err
		}
		globalGanacheSnapshot = snapshot
	}
	return conn, nil
}

// Snapshot current Ganache state
func snapshot(conn ethereum.Conn) (string, error) {
	var result string
	err := conn.RawClient.CallContext(context.Background(), &result, "evm_snapshot")
	return result, err
}

// RevertToSnapshot resets Ganache state to most recent snapshot
func revertToSnapshot(conn ethereum.Conn, snaptshotID string) error {
	var result bool
	err := conn.RawClient.CallContext(context.Background(), &result, "evm_revert", snaptshotID)
	return err
}

// DeployContracts to Ganache deploys REN and DNR contracts using the genesis private key
func DeployContracts(conn ethereum.Conn) error {
	return deployContracts(conn, genesisTransactor)
}

// DistributeEth transfers ETH to each of the addresses
func DistributeEth(conn ethereum.Conn, addresses ...common.Address) error {

	for _, address := range addresses {
		err := conn.TransferEth(context.Background(), genesisTransactor, address, big.NewInt(1000000000000000000))
		if err != nil {
			return err
		}
	}

	return nil
}

// DistributeREN transfers REN to each of the addresses
func DistributeREN(conn ethereum.Conn, addresses ...common.Address) error {

	republicTokenContract, err := bindings.NewRepublicToken(conn.RepublicTokenAddress, bind.ContractBackend(conn.Client))
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

func NewAccount(conn ethereum.Conn, eth *big.Int) (*bind.TransactOpts, common.Address, error) {
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

func deployContracts(conn ethereum.Conn, transactor *bind.TransactOpts) error {

	_, republicTokenAddress, err := deployRepublicToken(context.Background(), conn, transactor)
	if err != nil {
		panic(err)
	}

	_, darkNodeRegistryAddress, err := deployDarkNodeRegistry(context.Background(), conn, transactor, republicTokenAddress)
	if err != nil {
		panic(err)
	}

	if republicTokenAddress != ethereum.RepublicTokenAddressOnGanache {
		return fmt.Errorf("RepublicToken address has changed: expected: %s, got: %s", ethereum.RepublicTokenAddressOnGanache.Hex(), republicTokenAddress.Hex())
	}
	if darkNodeRegistryAddress != ethereum.DarknodeRegistryAddressOnGanache {
		return fmt.Errorf("DarknodeRegistry address has changed: expected: %s, got: %s", ethereum.DarknodeRegistryAddressOnGanache.Hex(), darkNodeRegistryAddress.Hex())
	}

	return nil
}

func deployRepublicToken(ctx context.Context, conn ethereum.Conn, auth *bind.TransactOpts) (*bindings.RepublicToken, common.Address, error) {
	address, tx, ren, err := bindings.DeployRepublicToken(auth, conn.Client)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("cannot deploy RepublicToken: %v", err)
	}
	conn.PatchedWaitDeployed(ctx, tx)
	return ren, address, nil
}

func deployDarkNodeRegistry(ctx context.Context, conn ethereum.Conn, auth *bind.TransactOpts, republicTokenAddress common.Address) (*bindings.DarkNodeRegistry, common.Address, error) {
	// 1 aiREN
	minimumBond := big.NewInt(1)
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
