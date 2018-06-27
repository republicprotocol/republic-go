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
	"github.com/republicprotocol/republic-go/contract"
	"github.com/republicprotocol/republic-go/contract/bindings"
)

const reset = "\x1b[0m"
const green = "\x1b[32;1m"

var genesisPrivateKey, genesisTransactor = genesis()

// GenesisPrivateKey used by Ganache.
func GenesisPrivateKey() ecdsa.PrivateKey {
	return *genesisPrivateKey
}

// GenesisTransactor used by Ganache.
func GenesisTransactor() bind.TransactOpts {
	return *genesisTransactor
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

	globalGanacheCounter++
	if globalGanacheCounter > 1 {
		globalGanacheCounter--
		return false
	}

	globalGanacheCmd = exec.Command("ganache-cli", fmt.Sprintf("--account=0x%x,1000000000000000000000", crypto.FromECDSA(genesisPrivateKey)))
	// TODO: Do something better with the output than printing it to Stdout and
	// Stderr, it gets so noisy. Disabled for now. Configurable output file
	// could be nice.
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	err := globalGanacheCmd.Start()
	if err != nil {
		panic(err)
	}
	go StopOnInterrupt()

	time.Sleep(10 * time.Second)
	return true
}

// Stop will kill the local Ganache instance
func Stop() {
	globalGanacheMu.Lock()
	defer globalGanacheMu.Unlock()

	// Reference count ganache
	if globalGanacheCounter == 0 {
		return
	}
	globalGanacheCounter--
	if globalGanacheCounter > 0 {
		return
	}

	// Graceful shutdown of ganache
	globalGanacheCmd.Process.Signal(syscall.SIGTERM)
	globalGanacheCmd.Wait()
}

// StopOnInterrupt will stop Ganache upon receiving receiving a interrupt signal
func StopOnInterrupt() {
	signals := make(chan os.Signal, 1)
	defer close(signals)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	globalGanacheMu.Lock()
	defer globalGanacheMu.Unlock()

	// Reference count ganache
	if globalGanacheCounter == 0 {
		return
	}
	globalGanacheCounter = 0

	// Graceful shutdown of ganache
	globalGanacheCmd.Process.Signal(syscall.SIGTERM)
	globalGanacheCmd.Wait()
}

// Connect to a local Ganache instance.
func Connect(ganacheRPC string) (contract.Conn, error) {
	raw, err := rpc.DialContext(context.Background(), ganacheRPC)
	if err != nil {
		return contract.Conn{}, err
	}
	fmt.Printf("Ganache is listening on %shttp://localhost:8545%s...\n", green, reset)
	ethclient := ethclient.NewClient(raw)

	return contract.Conn{
		RawClient: raw,
		Client:    ethclient,
		Config: contract.Config{
			Network: contract.NetworkLocal,
			URI:     "http://localhost:8545",
		},
	}, nil
}

// StartAndConnect to a local Ganache instance and deploy all smart contracts.
func StartAndConnect() (contract.Conn, error) {
	firstConnection := Start()

	conn, err := Connect("http://localhost:8545")
	if err != nil {
		return conn, err
	}

	if firstConnection {
		// Deploy contracts and take snapshot
		if err := DeployContracts(&conn); err != nil {
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

// DeployContracts to Ganache deploys REN and DNR contracts using the genesis private key
func DeployContracts(conn *contract.Conn) error {
	return deployContracts(conn, genesisTransactor)
}

// DistributeEth transfers ETH to each of the addresses
func DistributeEth(conn contract.Conn, addresses ...common.Address) error {

	for _, address := range addresses {
		err := conn.TransferEth(context.Background(), genesisTransactor, address, big.NewInt(9000000000000000000))
		if err != nil {
			return err
		}
	}

	return nil
}

// DistributeREN transfers REN to each of the addresses
func DistributeREN(conn contract.Conn, addresses ...common.Address) error {

	republicTokenContract, err := bindings.NewRepublicToken(common.HexToAddress(conn.Config.RepublicTokenAddress), bind.ContractBackend(conn.Client))
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

// NewAccount will return a new account along with it's associated address.
func NewAccount(conn contract.Conn, eth *big.Int) (*bind.TransactOpts, common.Address, error) {
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

func deployContracts(conn *contract.Conn, transactor *bind.TransactOpts) error {

	_, republicTokenAddress, err := deployRepublicToken(context.Background(), *conn, transactor)
	if err != nil {
		panic(err)
	}
	conn.Config.RepublicTokenAddress = republicTokenAddress.String()
	_, darknodeRegistryAddress, err := deployDarknodeRegistry(context.Background(), *conn, transactor, republicTokenAddress)
	if err != nil {
		panic(err)
	}
	conn.Config.DarknodeRegistryAddress = darknodeRegistryAddress.String()
	_, orderbookAddress, err := deployOrderbook(context.Background(), *conn, transactor, republicTokenAddress, darknodeRegistryAddress)
	if err != nil {
		panic(err)
	}
	conn.Config.OrderbookAddress = orderbookAddress.String()
	_, rewardVaultAddress, err := deployRewardVault(context.Background(), *conn, transactor, darknodeRegistryAddress)
	if err != nil {
		panic(err)
	}
	conn.Config.RewardVaultAddress = rewardVaultAddress.String()
	_, renExBalancesAddress, err := deployRenExBalances(context.Background(), *conn, transactor, rewardVaultAddress)
	if err != nil {
		panic(err)
	}
	conn.Config.RenExBalancesAddress = renExBalancesAddress.String()
	_, renExSettlementAddress, err := deployRenExSettlement(context.Background(), *conn, transactor, orderbookAddress, republicTokenAddress, renExBalancesAddress)
	if err != nil {
		panic(err)
	}
	conn.Config.RenExSettlementAddress = renExSettlementAddress.String()

	return nil
}

func deployRepublicToken(ctx context.Context, conn contract.Conn, auth *bind.TransactOpts) (*bindings.RepublicToken, common.Address, error) {
	address, tx, ren, err := bindings.DeployRepublicToken(auth, conn.Client)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("cannot deploy RepublicToken: %v", err)
	}
	conn.PatchedWaitDeployed(ctx, tx)
	return ren, address, nil
}

func deployDarknodeRegistry(ctx context.Context, conn contract.Conn, auth *bind.TransactOpts, republicTokenAddress common.Address) (*bindings.DarknodeRegistry, common.Address, error) {
	// 0 aiREN
	minimumBond := big.NewInt(0)
	// 1 second
	minimumEpochInterval := big.NewInt(1)
	// 24 Darknode in a Pod
	minimumPodSize := big.NewInt(24)

	address, tx, dnr, err := bindings.DeployDarknodeRegistry(auth, conn.Client, republicTokenAddress, minimumBond, minimumPodSize, minimumEpochInterval)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("cannot deploy DarknodeRegistry: %v", err)
	}
	conn.PatchedWaitDeployed(ctx, tx)
	return dnr, address, nil
}

func deployOrderbook(ctx context.Context, conn contract.Conn, auth *bind.TransactOpts, republicTokenAddress, registryAddress common.Address) (*bindings.Orderbook, common.Address, error) {
	// 0 REN
	minimumFee := big.NewInt(0)

	address, tx, ren, err := bindings.DeployOrderbook(auth, conn.Client, minimumFee, republicTokenAddress, registryAddress)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("cannot deploy Orderbook: %v", err)
	}
	conn.PatchedWaitDeployed(ctx, tx)
	return ren, address, nil
}

func deployRewardVault(ctx context.Context, conn contract.Conn, auth *bind.TransactOpts, darknodeRegistryAddress common.Address) (*bindings.RewardVault, common.Address, error) {
	address, tx, vault, err := bindings.DeployRewardVault(auth, conn.Client, darknodeRegistryAddress)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("cannot deploy RewardVault: %v", err)
	}
	conn.PatchedWaitDeployed(ctx, tx)
	return vault, address, nil
}

func deployRenExBalances(ctx context.Context, conn contract.Conn, auth *bind.TransactOpts, rewardVaultAddress common.Address) (*bindings.RenExBalances, common.Address, error) {
	address, tx, balances, err := bindings.DeployRenExBalances(auth, conn.Client, rewardVaultAddress)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("cannot deploy RenExBalances: %v", err)
	}
	conn.PatchedWaitDeployed(ctx, tx)
	return balances, address, nil
}

func deployRenExSettlement(ctx context.Context, conn contract.Conn, auth *bind.TransactOpts, orderbookAddress, tokenAddress, balancesAddress common.Address) (*bindings.RenExSettlement, common.Address, error) {
	GWEI := 1000000000
	address, tx, accounts, err := bindings.DeployRenExSettlement(auth, conn.Client, orderbookAddress, tokenAddress, balancesAddress, big.NewInt(int64(100*GWEI)))
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("cannot deploy RenExSettlements: %v", err)
	}
	conn.PatchedWaitDeployed(ctx, tx)
	return accounts, address, nil
}

// Snapshot current Ganache state
func snapshot(conn contract.Conn) (string, error) {
	var result string
	err := conn.RawClient.CallContext(context.Background(), &result, "evm_snapshot")
	return result, err
}

// RevertToSnapshot resets Ganache state to most recent snapshot
func revertToSnapshot(conn contract.Conn, snapshotID string) error {
	var result bool
	err := conn.RawClient.CallContext(context.Background(), &result, "evm_revert", snapshotID)
	return err
}
