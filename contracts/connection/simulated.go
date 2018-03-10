package connection

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/republic-go/contracts/bindings"
)

// DeployREN deploys an ERC20 contract
func DeployREN(context context.Context, connection Client, auth *bind.TransactOpts) (*bindings.ERC20, common.Address) {
	// Deploy a token contract on the simulated blockchain
	address, tx, ren, err := bindings.DeployERC20(auth, connection)
	if err != nil {
		log.Fatalf("Failed to deploy: %v", err)
	}
	PatchedWaitDeployed(context, connection, tx)
	return ren, address
}

// DeployDNR deploys a Dark Node Registrar
func DeployDNR(context context.Context, connection Client, auth *bind.TransactOpts, renAddress common.Address) (*bindings.DarkNodeRegistrar, common.Address) {
	// Deploy a token contract on the simulated blockchain
	minimumBond := big.NewInt(100)
	minimumEpochInterval := big.NewInt(60)
	address, tx, dnr, err := bindings.DeployDarkNodeRegistrar(auth, connection, renAddress, minimumBond, minimumEpochInterval)
	if err != nil {
		log.Fatalf("Failed to deploy: %v", err)
	}
	PatchedWaitDeployed(context, connection, tx)
	return dnr, address
}

// Simulated will create a simulated client
func Simulated(auths ...*bind.TransactOpts) (ClientDetails, error) {

	deployerKey, err := crypto.GenerateKey()
	if err != nil {
		return ClientDetails{}, err
	}
	deployerAuth := bind.NewKeyedTransactor(deployerKey)

	// Set up Ether balances
	alloc := core.GenesisAlloc{
		deployerAuth.From: {Balance: big.NewInt(9000000000000000000)},
	}
	for _, auth := range auths {
		alloc[auth.From] = core.GenesisAccount{Balance: big.NewInt(9000000000000000000)}
	}
	sim := backends.NewSimulatedBackend(alloc)

	// Deploy contracts
	ren, renAddress := DeployREN(context.Background(), sim, deployerAuth)
	_, dnrAddress := DeployDNR(context.Background(), sim, deployerAuth, renAddress)

	// Transfer Ren to each participant
	for _, auth := range auths {
		_, err := ren.Transfer(deployerAuth, auth.From, big.NewInt(1000000000000000000))
		if err != nil {
			return ClientDetails{}, err
		}
	}
	sim.Commit()

	for _, auth := range auths {
		alloc[auth.From] = core.GenesisAccount{Balance: big.NewInt(9000000000000000000)}
	}

	return ClientDetails{
		Client:     sim,
		RenAddress: renAddress,
		DNRAddress: dnrAddress,
	}, nil
}
