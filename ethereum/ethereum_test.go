package ethereum_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-atom/ethereum"

	// . "github.com/republicprotocol/go-atom/ethereum"

	"context"

	"crypto/sha256"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var ether = big.NewInt(1000000000000000000)

const key1 = `{"version":3,"id":"7844982f-abe7-4690-8c15-34f75f847c66","address":"db205ea9d35d8c01652263d58351af75cfbcbf07","Crypto":{"ciphertext":"378dce3c1279b36b071e1c7e2540ac1271581bff0bbe36b94f919cb73c491d3a","cipherparams":{"iv":"2eb92da55cc2aa62b7ffddba891f5d35"},"cipher":"aes-128-ctr","kdf":"scrypt","kdfparams":{"dklen":32,"salt":"80d3341678f83a14024ba9c3edab072e6bd2eea6aa0fbc9e0a33bae27ffa3d6d","n":8192,"r":8,"p":1},"mac":"3d07502ea6cd6b96a508138d8b8cd2e46c3966240ff276ce288059ba4235cb0d"}}`
const key2 = `{"version":3,"id":"1bc823af-210a-4143-8eb4-306c19485622","address":"d95bd5b44a1290c91a31af1114e49b961e56b03b","Crypto":{"ciphertext":"0eb788eee71b9796390d6b3495c25d87746a7c2ddd98a641b90e7271231f6df0","cipherparams":{"iv":"35decc930518c37116e8fcf1f9933948"},"cipher":"aes-128-ctr","kdf":"scrypt","kdfparams":{"dklen":32,"salt":"a3b2f03dc27ee3c89b6d21b1e0d1973bb130524e7570bd0bf4da531313df9730","n":8192,"r":8,"p":1},"mac":"ef73d708b39832e6be65db73d088ab60a2c7711189a330324b410bc64c6bfe7a"}}`

func loadAccounts(client *ethclient.Client) (*bind.TransactOpts, *bind.TransactOpts) {
	auth1, err := bind.NewTransactor(strings.NewReader(key1), "password1")
	Ω(err).Should(BeNil())

	auth2, err := bind.NewTransactor(strings.NewReader(key2), "password2")
	Ω(err).Should(BeNil())

	// Ensure that the bigger acount is sending to the smaller account
	bal1, _ := client.BalanceAt(context.Background(), auth1.From, nil)
	bal2, _ := client.BalanceAt(context.Background(), auth2.From, nil)
	if bal1.Cmp(bal2) > 0 {
		auth1, auth2 = auth2, auth1
	}

	return auth1, auth2
}

var _ = Describe("Ethereum", func() {

	It("can swap with an arbitrary ledger", func() {
		// Connect to Infura (or use local node at 13.54.129.55:8180)
		client := ethereum.Ropsten("https://ropsten.infura.io/")
		// Load accounts
		auth1, auth2 := loadAccounts(client)
		// Contract address
		contractAddress := common.HexToAddress("0x32Dad9E9Fe2A3eA2C2c643675A7d2A56814F554f")

		// ALICE has locked up tokens on another ledger with the following secret
		secret := []byte("this is the secret")
		secretHash := sha256.Sum256(secret)

		// BOB reciprocates the atomic swap with the secretHash
		user2Connection := ethereum.NewETHAtomContract(context.Background(), client, auth2, contractAddress, nil)
		value := big.NewInt(0).Mul(ether, big.NewInt(1))
		err := user2Connection.Initiate(secretHash[:], auth1.From.Bytes(), auth2.From.Bytes(), value, time.Now().Add(48*time.Hour).Unix())
		Ω(err).Should(BeNil())

		// ALICE checks that Bob has set up the swap correctly
		user1Connection := ethereum.NewETHAtomContract(context.Background(), client, auth1, contractAddress, user2Connection.GetData())
		// Checks that the hash is right
		retrievedHash, to, _, readValue, expiry, err := user1Connection.Read()
		Ω(err).Should(BeNil())
		Ω(retrievedHash).Should(Equal(secretHash[:]))
		Ω(to).Should(Equal(auth1.From.Bytes()))
		// Ω(from).Should(Equal(auth2.From.Bytes()))
		Ω(value).Should(Equal(readValue))
		Ω(expiry).Should(BeNumerically(">=", time.Now().Add(time.Hour*23).Unix()))
		// ALICE redeems the ether by revealing the secret
		err = user1Connection.Redeem(secret)
		Ω(err).Should(BeNil())

		// BOB retrieves the secret to use on the other ledger
		retSecret, err := user2Connection.ReadSecret()
		Ω(err).Should(BeNil())
		Ω(retSecret).Should(Equal(secret))
	})

	It("can't redeem with a wrong password", func() {
		client := ethereum.Ropsten("https://ropsten.infura.io/")
		auth1, auth2 := loadAccounts(client)
		contractAddress := common.HexToAddress("0x32Dad9E9Fe2A3eA2C2c643675A7d2A56814F554f")

		secret := []byte("this is the secret")
		secretHash := sha256.Sum256(secret)

		// BOB reciprocates the atomic swap with the secretHash
		user2Connection := ethereum.NewETHAtomContract(context.Background(), client, auth2, contractAddress, nil)
		value := big.NewInt(0).Mul(ether, big.NewInt(1))
		err := user2Connection.Initiate(secretHash[:], auth1.From.Bytes(), auth2.From.Bytes(), value, time.Now().Add(48*time.Hour).Unix())
		Ω(err).Should(BeNil())

		//ALICE tries to reem with wrong password, and then right password
		user1Connection := ethereum.NewETHAtomContract(context.Background(), client, auth1, contractAddress, user2Connection.GetData())
		wrongSecret := []byte("this is NOT the secret")
		err = user1Connection.Redeem(wrongSecret)
		// Error should NOT be nil
		Ω(err).Should(Not(BeNil()))

		err = user1Connection.Redeem(secret)
		Ω(err).Should(BeNil())
	})

	It("can't expire before timeout", func() {
		client := ethereum.Ropsten("https://ropsten.infura.io/")
		auth1, auth2 := loadAccounts(client)
		contractAddress := common.HexToAddress("0x32Dad9E9Fe2A3eA2C2c643675A7d2A56814F554f")

		secret := []byte("this is the secret")
		secretHash := sha256.Sum256(secret)

		// BOB reciprocates the atomic swap with the secretHash
		user2Connection := ethereum.NewETHAtomContract(context.Background(), client, auth2, contractAddress, nil)
		value := big.NewInt(0).Mul(ether, big.NewInt(1))
		err := user2Connection.Initiate(secretHash[:], auth1.From.Bytes(), auth2.From.Bytes(), value, time.Now().Add(48*time.Hour).Unix())
		Ω(err).Should(BeNil())

		// Should NOT be able to refund
		err = user2Connection.Refund()
		Ω(err).Should(Not(BeNil()))

		// ALICE redeems by revealing secret
		user1Connection := ethereum.NewETHAtomContract(context.Background(), client, auth1, contractAddress, user2Connection.GetData())
		err = user1Connection.Redeem(secret)
		Ω(err).Should(BeNil())
	})

})
