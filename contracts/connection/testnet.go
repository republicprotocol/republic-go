package connection

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/republicprotocol/republic-go/contracts/bindings"
)

// The key that contains all ther ren and testnet eth
// Fixed so that the tests can easily request eth and ren
var genesisPrivateKey, genesisTransactor = genesisKey()

// GenesisAuth is used to create a connection to contracts not specifci to a particular node
var GenesisAuth = genesisTransactor

var details, _ = ConnectToTestnet()

// Because the genesis key is fixed, these should be as well
var renAddress = common.HexToAddress("0x8DE2a0D285cd6fDB47ABAe34024a6EED79ef0E92")
var dnrAddress = common.HexToAddress("0xbf195e17802736ff4e19275b961bb1C2D45f2C8d")

// ConnectToTestnet connects to the rpc client at the local port 8545
func ConnectToTestnet() (ClientDetails, error) {

	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		return ClientDetails{}, err
	}

	return ClientDetails{
		Client:     client,
		DNRAddress: dnrAddress,
		RenAddress: renAddress,
		Chain:      ChainGanache,
	}, nil
}

// StartTestnet starts ganache on port 8545
func StartTestnet(debug bool, wg *sync.WaitGroup) *exec.Cmd {

	wg.Add(1)
	cmd := startGanache(genesisPrivateKey, debug)

	return cmd
}

// DeployContractsToGanache deploys REN and DNR contracts using the genesis private key
func DeployContractsToGanache(uri string) error {
	client, err := ethclient.Dial(uri)
	if err != nil {
		return err
	}

	_, err = deployContracts(client, genesisTransactor)
	return err
}

func startGanache(privateKey *ecdsa.PrivateKey, debug bool) *exec.Cmd {

	arg := fmt.Sprintf("--account=0x%x,1000000000000000000000", crypto.FromECDSA(privateKey))
	cmd := exec.Command("ganache-cli", arg)

	if debug {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr

	cmd.Start()

	return cmd
}

// DistributeEth transfers testnet ETH to each of the addresses
func DistributeEth(conn ClientDetails, addresses ...common.Address) error {

	if conn.Chain != ChainGanache {
		return errors.New("must be using ganache to distribute eth")
	}

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

// DistributeRen transfers testnet REN to each of the addresses
func DistributeRen(conn ClientDetails, addresses ...common.Address) error {

	if conn.Chain != ChainGanache {
		return errors.New("must be using ganache to distribute ren")
	}

	renContract, err := bindings.NewRepublicToken(conn.RenAddress, bind.ContractBackend(conn.Client))
	if err != nil {
		return err
	}

	// Transfer Ren to each participant
	for _, address := range addresses {
		_, err := renContract.Transfer(genesisTransactor, address, big.NewInt(0).Mul(big.NewInt(100), big.NewInt(1000000000000000000)))

		if err != nil {
			return err
		}
	}

	return nil
}

func genesisKey() (*ecdsa.PrivateKey, *bind.TransactOpts) {
	// deployerKey, err := crypto.GenerateKey()
	deployerKey, err := crypto.HexToECDSA("2aba04ee8a322b8648af2a784144181a0c793f1a2e80519418f3d20bbfb22249")
	if err != nil {
		panic("couldn't read genesis key")
	}
	// fmt.Printf("%x", crypto.FromECDSA(deployerKey))
	deployerAuth := bind.NewKeyedTransactor(deployerKey)
	return deployerKey, deployerAuth
}

// deployREN deploys an ERC20 contract
func deployREN(context context.Context, conn ClientDetails, auth *bind.TransactOpts) (*bindings.RepublicToken, common.Address, error) {
	// Deploy a token contract on the simulated blockchain
	address, tx, ren, err := bindings.DeployRepublicToken(auth, conn.Client)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("Failed to deploy REN: %v", err)
	}
	conn.PatchedWaitDeployed(context, tx)
	return ren, address, nil
}

// deployDNR deploys a Dark Node Registrar
func deployDNR(context context.Context, conn ClientDetails, auth *bind.TransactOpts, renAddress common.Address) (*bindings.DarkNodeRegistry, common.Address, error) {
	// Deploy a token contract on the simulated blockchain
	// 10 aiREN
	minimumBond := big.NewInt(10)
	// One minute
	minimumEpochInterval := big.NewInt(1)
	minimumPoolSize := big.NewInt(74)
	address, tx, dnr, err := bindings.DeployDarkNodeRegistry(auth, conn.Client, renAddress, minimumBond, minimumPoolSize, minimumEpochInterval)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("Failed to deploy DNR: %v", err)
	}
	conn.PatchedWaitDeployed(context, tx)
	return dnr, address, nil
}

// deployContracts deploys the REN and DNR contracts
func deployContracts(client *ethclient.Client, transactor *bind.TransactOpts) (ClientDetails, error) {

	conn := ClientDetails{
		Client: client,
		Chain:  ChainGanache,
	}

	// Deploy contracts
	_, _renAddress, err := deployREN(context.Background(), conn, transactor)
	if err != nil {
		return ClientDetails{}, err
	}
	// sim.Commit()

	// t, err := ren.TotalSupply(&bind.CallOpts{})
	// if err != nil {
	// 	return ClientDetails{}, err
	// }

	_, _dnrAddress, err := deployDNR(context.Background(), conn, transactor, _renAddress)
	if err != nil {
		return ClientDetails{}, err
	}

	if _renAddress != renAddress || _dnrAddress != dnrAddress {
		return ClientDetails{}, errors.New("ganache contract addresses have changed")
	}

	conn.DNRAddress = _dnrAddress
	conn.RenAddress = _renAddress
	return conn, nil
}
