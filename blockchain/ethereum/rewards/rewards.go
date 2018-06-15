package rewards

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/republicprotocol/republic-go/blockchain/ethereum"
)

// DarknodeRegistry is the dark node interface
type RewardVaultContract struct {
	network      ethereum.Network
	context      context.Context
	conn         ethereum.Conn
	transactOpts *bind.TransactOpts
	callOpts     *bind.CallOpts
	binding      *RewardVault
	address      common.Address
}

func NewRewardVaultContract(context context.Context, conn ethereum.Conn, transactOpts *bind.TransactOpts, callOpts *bind.CallOpts) (RewardVaultContract, error) {
	contract, err := NewRewardVault(common.HexToAddress(conn.Config.DarknodeRegistryAddress), bind.ContractBackend(conn.Client))
	if err != nil {
		return RewardVaultContract{}, err
	}
	return RewardVaultContract{
		network:      conn.Config.Network,
		context:      context,
		conn:         conn,
		transactOpts: transactOpts,
		callOpts:     callOpts,
		binding:      contract,
		//address:      common.HexToAddress(conn.Config.RewardsContractAddress),
	}, nil
}

func (reward *RewardVaultContract) Finalize(round *big.Int) error {
	tx, err := reward.binding.Finalize(reward.transactOpts, round)
	if err != nil {
		return err
	}

	_, err = reward.conn.PatchedWaitMined(reward.context, tx)
	return err
}

func (reward *RewardVaultContract) IsFinalizable(round *big.Int) (bool, error) {
	return reward.binding.IsFinalizable(reward.callOpts, round)
}

func (reward *RewardVaultContract) Rewardees(round *big.Int) ([]string, error) {
	addresses, err := reward.binding.Rewardees(reward.callOpts, round)
	if err != nil {
		return nil, err
	}
	addrs := make([]string, len(addresses))
	for i := range addrs {
		addrs[i] = addresses[i].Hex()
	}

	return addrs, nil
}

func (reward *RewardVaultContract) Challenges(round *big.Int) ([][32]byte, error) {
	return reward.binding.ChallengeIds(reward.callOpts, round)
}

func (reward *RewardVaultContract) Solve(round *big.Int, proofs [][16]byte) error {
	tx, err := reward.binding.SubmitProof(reward.transactOpts, round, proofs)
	if err != nil {
		return err
	}

	_, err = reward.conn.PatchedWaitMined(reward.context, tx)
	return err
}
