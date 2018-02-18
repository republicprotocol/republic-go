package dnr

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/republicprotocol/go-dark-node-registrar/contracts"
)

type EthereumClient struct {
	context context.Context
	client  *ethclient.Client
	auth1   *bind.TransactOpts
	auth2   *bind.CallOpts
	binding *contracts.DarkNodeRegistrar
	chainID int8
}

// type DarkNode struct {
// 	ethereumClient     EthereumClient
// 	owner              string
// 	bond               big.Int
// 	publicKey          []byte
// 	commitment         [32]byte
// 	registeredPosition big.Int
// 	registeredAt       big.Int
// 	deregisteredAt     big.Int
// }

func NewDarkNodeRegistrar(context context.Context, client *ethclient.Client, auth1 *bind.TransactOpts, auth2 *bind.CallOpts, address common.Address, data []byte, chainId int8) *EthereumClient {
	contract, err := contracts.NewDarkNodeRegistrar(address, bind.ContractBackend(client))
	if err != nil {
		log.Fatalf("%v", err)
	}
	return &EthereumClient{
		context: context,
		client:  client,
		auth1:   auth1,
		auth2:   auth2,
		binding: contract,
		chainID: chainId,
	}
}

func (ethereumClient *EthereumClient) Register(_darkNodeID [20]byte, _publicKey []byte) (*types.Transaction, error) {
	return ethereumClient.binding.Register(ethereumClient.auth1, _darkNodeID, _publicKey)
}

func (ethereumClient *EthereumClient) Deregister(_darkNodeID [20]byte) (*types.Transaction, error) {
	return ethereumClient.binding.Deregister(ethereumClient.auth1, _darkNodeID)
}

func (ethereumClient *EthereumClient) GetBond(_darkNodeID [20]byte) (*big.Int, error) {
	return ethereumClient.binding.GetBond(ethereumClient.auth2, _darkNodeID)
}

func (ethereumClient *EthereumClient) IsDarkNodeRegistered(_darkNodeID [20]byte) (bool, error) {
	return ethereumClient.binding.IsDarkNodeRegistered(ethereumClient.auth2, _darkNodeID)
}

func (ethereumClient *EthereumClient) CurrentEpoch() (struct {
	Blockhash [32]byte
	Timestamp *big.Int
}, error) {
	return ethereumClient.binding.CurrentEpoch(ethereumClient.auth2)
}

func (ethereumClient *EthereumClient) Epoch() (*types.Transaction, error) {
	return ethereumClient.binding.Epoch(ethereumClient.auth1)
}

func (ethereumClient *EthereumClient) GetCommitment(_darkNodeID [20]byte) ([32]byte, error) {
	return ethereumClient.binding.GetCommitment(ethereumClient.auth2, _darkNodeID)
}

func (ethereumClient *EthereumClient) GetOwner(_darkNodeID [20]byte) (common.Address, error) {
	return ethereumClient.binding.GetOwner(ethereumClient.auth2, _darkNodeID)
}

func (ethereumClient *EthereumClient) GetPublicKey(_darkNodeID [20]byte) ([]byte, error) {
	return ethereumClient.binding.GetPublicKey(ethereumClient.auth2, _darkNodeID)
}

func (ethereumClient *EthereumClient) GetXingOverlay() ([][20]byte, error) {
	return ethereumClient.binding.GetXingOverlay(ethereumClient.auth2)
}
func (ethereumClient *EthereumClient) MinimumBond() (*big.Int, error) {
	return ethereumClient.binding.MinimumBond(ethereumClient.auth2)
}

func (ethereumClient *EthereumClient) MinimumEpochInterval() (*big.Int, error) {
	return ethereumClient.binding.MinimumEpochInterval(ethereumClient.auth2)
}

func (ethereumClient *EthereumClient) PendingRefunds(arg0 common.Address) (*big.Int, error) {
	return ethereumClient.binding.PendingRefunds(ethereumClient.auth2, arg0)
}

func (ethereumClient *EthereumClient) Refund() (*types.Transaction, error) {
	return ethereumClient.binding.Refund(ethereumClient.auth1)
}
