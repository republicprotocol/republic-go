package ethereum_test

import (
	"context"
	"log"
	"math/big"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/republicprotocol/go-atom/ethereum"
	"github.com/republicprotocol/go-atom/ethereum/contracts"
)

var ether = big.NewInt(1000000000000000000)

const key1 = `{"version":3,"id":"7844982f-abe7-4690-8c15-34f75f847c66","address":"db205ea9d35d8c01652263d58351af75cfbcbf07","Crypto":{"ciphertext":"378dce3c1279b36b071e1c7e2540ac1271581bff0bbe36b94f919cb73c491d3a","cipherparams":{"iv":"2eb92da55cc2aa62b7ffddba891f5d35"},"cipher":"aes-128-ctr","kdf":"scrypt","kdfparams":{"dklen":32,"salt":"80d3341678f83a14024ba9c3edab072e6bd2eea6aa0fbc9e0a33bae27ffa3d6d","n":8192,"r":8,"p":1},"mac":"3d07502ea6cd6b96a508138d8b8cd2e46c3966240ff276ce288059ba4235cb0d"}}`
const key2 = `{"version":3,"id":"1bc823af-210a-4143-8eb4-306c19485622","address":"d95bd5b44a1290c91a31af1114e49b961e56b03b","Crypto":{"ciphertext":"0eb788eee71b9796390d6b3495c25d87746a7c2ddd98a641b90e7271231f6df0","cipherparams":{"iv":"35decc930518c37116e8fcf1f9933948"},"cipher":"aes-128-ctr","kdf":"scrypt","kdfparams":{"dklen":32,"salt":"a3b2f03dc27ee3c89b6d21b1e0d1973bb130524e7570bd0bf4da531313df9730","n":8192,"r":8,"p":1},"mac":"ef73d708b39832e6be65db73d088ab60a2c7711189a330324b410bc64c6bfe7a"}}`

const sim = true

func etherAtomAddress(client ethereum.Client, auth *bind.TransactOpts) common.Address {
	var address common.Address
	if sim {
		var err error
		var tx *types.Transaction
		address, tx, _, err = contracts.DeployAtomicSwapEther(auth, client)
		if err != nil {
			log.Fatalf("Failed to deploy Ether-Atom: %v", err)
		}
		ethereum.PatchedWaitDeployed(context.Background(), client, tx)
	} else {
		address = common.HexToAddress("0xbd59e72598737a08a68fe192df04c773adbbfa53")
	}

	return address
}

func erc20AtomAddress(client ethereum.Client, auth *bind.TransactOpts) common.Address {
	var address common.Address
	if sim {
		var err error
		var tx *types.Transaction
		address, tx, _, err = contracts.DeployAtomicSwapERC20(auth, client)
		if err != nil {
			log.Fatalf("Failed to deploy ERC20-Atom: %v", err)
		}
		ethereum.PatchedWaitDeployed(context.Background(), client, tx)
	} else {
		address = common.HexToAddress("...")
	}

	return address
}

func erc20Address(client ethereum.Client, auth1, auth2 *bind.TransactOpts) common.Address {
	var address common.Address
	if sim {
		var err error
		var tx *types.Transaction
		var erc20 *contracts.TestERC20
		address, tx, erc20, err = contracts.DeployTestERC20(auth1, client)
		if err != nil {
			log.Fatalf("Failed to deploy ERC20: %v", err)
		}
		ethereum.PatchedWaitDeployed(context.Background(), client, tx)
		tx, _ = erc20.Transfer(auth1, auth2.From, ether)
	} else {
		address = common.HexToAddress("...")
	}

	return address
}

func loadClient() (ethereum.Client, *bind.TransactOpts, *bind.TransactOpts) {

	auth1, err := bind.NewTransactor(strings.NewReader(key1), "password1")
	Ω(err).Should(BeNil())

	auth2, err := bind.NewTransactor(strings.NewReader(key2), "password2")
	Ω(err).Should(BeNil())

	var client ethereum.Client
	if sim {
		client = ethereum.Simulated(auth1, auth2)
	} else {
		// Connect to Infura (or use local node at 13.54.129.55:8180)
		client = ethereum.Ropsten("https://ropsten.infura.io/")
	}

	return client, auth1, auth2
}

var _ = Describe("Ethereum", func() {

	It("test utils", func() {
		_, err := ethereum.BytesTo32Bytes([]byte("this is not 32 bytes long"))
		Ω(err).Should(Not(BeNil()))
	})

})
