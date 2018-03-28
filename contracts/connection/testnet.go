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
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/republicprotocol/republic-go/contracts/bindings"
)

// The key that contains all ther ren and testnet eth
// Fixed so that the tests can easily request eth and ren
var ganachePrivateKey, ganacheTransactor = genesisKey()

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
func StartTestnet(debug bool) {
	var wg sync.WaitGroup

	wg.Add(1)
	cmd := startGanache(ganachePrivateKey, debug)
	defer cmd.Process.Kill()
	go waitGanache(cmd, &wg)

	time.Sleep(5 * time.Second)

	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		cmd.Process.Kill()
		panic(err)
	}

	_, err = DeployContracts(client, ganacheTransactor)
	if err != nil {
		cmd.Process.Kill()
		panic(err)
	}

	wg.Wait()
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

func waitGanache(cmd *exec.Cmd, wg *sync.WaitGroup) {
	defer wg.Done()
	cmd.Wait()
}

func distributeEth(conn ClientDetails, addresses ...ecdsa.PublicKey) error {

	if conn.Chain != ChainGanache {
		return errors.New("must be using ganache to distribute eth")
	}

	transactor := &bind.TransactOpts{
		From:     ganacheTransactor.From,
		Nonce:    ganacheTransactor.Nonce,
		Signer:   ganacheTransactor.Signer,
		Value:    big.NewInt(1000000000000000000),
		GasPrice: ganacheTransactor.GasPrice,
		GasLimit: 30000,
		Context:  ganacheTransactor.Context,
	}

	for _, address := range addresses {
		to := crypto.PubkeyToAddress(address)

		bound := bind.NewBoundContract(to, abi.ABI{}, nil, conn.Client, nil)

		tx, err := bound.Transfer(transactor)
		if err != nil {
			return err
		}
		conn.PatchedWaitMined(context.Background(), tx)
	}

	return nil
}

func distributeRen(conn ClientDetails, addresses ...ecdsa.PublicKey) error {

	if conn.Chain != ChainGanache {
		return errors.New("must be using ganache to distribute ren")
	}

	renContract, err := bindings.NewERC20(conn.RenAddress, bind.ContractBackend(conn.Client))
	if err != nil {
		return err
	}

	// Transfer Ren to each participant
	for _, address := range addresses {
		_, err := renContract.Transfer(ganacheTransactor, crypto.PubkeyToAddress(address), big.NewInt(1000000000000000000))
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

// DeployDarkNodeRegistrar

// import (
// 	"context"
// 	"flag"
// 	"fmt"
// 	"log"
// 	"math/big"

// 	"github.com/ethereum/go-ethereum/accounts/abi/bind"
// 	"github.com/ethereum/go-ethereum/common"
// 	base58 "github.com/jbenet/go-base58"
// 	"github.com/republicprotocol/republic-go/contracts/bindings"
// 	"github.com/republicprotocol/republic-go/contracts/connection"
// 	node "github.com/republicprotocol/republic-go/dark-node"
// )

// var config *node.Config

// func oldMain() {
// 	err := parseCommandLineFlags()
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	auth := bind.NewKeyedTransactor(config.EthereumKey.PrivateKey)

// 	client, err := FromURI("https://ropsten.infura.io/", "ropsten")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// REPLACE REN ADDRESS HERE
// 	renContract := common.HexToAddress("")
// 	address, tx, _, err := bindings.DeployDarkNodeRegistrar(auth, client.Client, renContract, big.NewInt(1000), big.NewInt(60))
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	_, err = PatchedWaitDeployed(context.Background(), client.Client, tx)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	fmt.Printf("[%v] Contract deployed at %s%v%s\n", base58.Encode(config.KeyPair.ID()), green, address.Hex(), reset)
// }

// DeployREN deploys an ERC20 contract
func DeployREN(context context.Context, conn ClientDetails, auth *bind.TransactOpts) (*bindings.TestERC20, common.Address, error) {
	// Deploy a token contract on the simulated blockchain
	address, tx, ren, err := bindings.DeployTestERC20(auth, conn.Client)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("Failed to deploy REN: %v", err)
	}
	conn.PatchedWaitDeployed(context, tx)
	return ren, address, nil
}

// DeployDNR deploys a Dark Node Registrar
func DeployDNR(context context.Context, conn ClientDetails, auth *bind.TransactOpts, renAddress common.Address) (*bindings.DarkNodeRegistrar, common.Address, error) {
	// Deploy a token contract on the simulated blockchain
	minimumBond := big.NewInt(100)
	minimumEpochInterval := big.NewInt(60)
	address, tx, dnr, err := bindings.DeployDarkNodeRegistrar(auth, conn.Client, renAddress, minimumBond, minimumEpochInterval)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("Failed to deploy DNR: %v", err)
	}
	conn.PatchedWaitDeployed(context, tx)
	return dnr, address, nil
}

// DeployContracts deploys the REN and DNR contracts
func DeployContracts(client *ethclient.Client, transactor *bind.TransactOpts) (ClientDetails, error) {

	conn := ClientDetails{
		Client: client,
		Chain:  ChainGanache,
	}

	// Deploy contracts
	_, _renAddress, err := DeployREN(context.Background(), conn, transactor)
	if err != nil {
		return ClientDetails{}, err
	}
	// sim.Commit()

	// t, err := ren.TotalSupply(&bind.CallOpts{})
	// if err != nil {
	// 	return ClientDetails{}, err
	// }

	_, _dnrAddress, err := DeployDNR(context.Background(), conn, transactor, _renAddress)
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
