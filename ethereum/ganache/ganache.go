package ganache

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/exec"

	"github.com/ethereum/go-ethereum/accounts/abi"
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

// DeployContracts to Ganache deploys REN and DNR contracts using the genesis private key
func DeployContracts(conn client.Connection) error {
	return deployContracts(conn, genesisTransactor)
}

// DistributeEth transfers ETH to each of the addresses
func DistributeEth(conn client.Connection, addresses ...common.Address) error {

	transactor := &bind.TransactOpts{
		From:     genesisTransactor.From,
		Nonce:    genesisTransactor.Nonce,
		Signer:   genesisTransactor.Signer,
		Value:    big.NewInt(1000000000000000000),
		GasPrice: genesisTransactor.GasPrice,
		GasLimit: 30000,
		Context:  genesisTransactor.Context,
	}

	for _, address := range addresses {
		bound := bind.NewBoundContract(address, abi.ABI{}, nil, conn.Client, nil)
		tx, err := bound.Transfer(transactor)
		if err != nil {
			return err
		}
		conn.PatchedWaitMined(context.Background(), tx)
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

	 // Deploy contracts
	_, republicTokenAddress, err := deployRepublicToken(context.Background(), conn, transactor)
	if err != nil {
		return err
	}
	_, darkNodeRegistryAddress, err := deployDarkNodeRegistry(context.Background(), conn, transactor, republicTokenAddress)
	if err != nil {
		return err
	}
	_, hyperdriveRegistryAdrress, err := deployHyperdriveRegistry(context.Background(), conn, transactor)
	if err != nil {
		return err
	}

	// Check contract address
	if republicTokenAddress != client.RepublicTokenAddressOnGanache {
		return fmt.Errorf("RepublicToken address has changed: expected: %s, got: %s", client.RepublicTokenAddressOnGanache.Hex(), republicTokenAddress.Hex())
	}
	if darkNodeRegistryAddress != client.DarkNodeRegistryAddressOnGanache {
		return fmt.Errorf("DarkNodeRegistry address has changed: expected: %s, got: %s", client.DarkNodeRegistryAddressOnGanache.Hex(), darkNodeRegistryAddress.Hex())
	}
	if hyperdriveRegistryAdrress != client.HyperDriveRegistryAdreessOnGanache {
		return fmt.Errorf("HyperdriveRegistry address has changed: expected: %s, got: %s", client.HyperDriveRegistryAdreessOnGanache.Hex(), hyperdriveRegistryAdrress.Hex())
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

func deployHyperdriveRegistry(ctx context.Context, conn client.Connection, auth *bind.TransactOpts) (*bindings.HyperdriveEpoch, common.Address, error) {
	// 1 second
	minimumEpochInterval := big.NewInt(1)

	address, tx, ren, err := bindings.DeployHyperdriveEpoch(auth, conn.Client, minimumEpochInterval)
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
