package go_eth_test

import (
	// "context"

	"crypto/rand"
	"encoding/hex"
	"log"
	"math/big"
	"strings"

	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/republicprotocol/go-eth"
)

var ether = big.NewInt(1000000000000000000)

const key1 = `{"version":3,"id":"7844982f-abe7-4690-8c15-34f75f847c66","address":"db205ea9d35d8c01652263d58351af75cfbcbf07","Crypto":{"ciphertext":"378dce3c1279b36b071e1c7e2540ac1271581bff0bbe36b94f919cb73c491d3a","cipherparams":{"iv":"2eb92da55cc2aa62b7ffddba891f5d35"},"cipher":"aes-128-ctr","kdf":"scrypt","kdfparams":{"dklen":32,"salt":"80d3341678f83a14024ba9c3edab072e6bd2eea6aa0fbc9e0a33bae27ffa3d6d","n":8192,"r":8,"p":1},"mac":"3d07502ea6cd6b96a508138d8b8cd2e46c3966240ff276ce288059ba4235cb0d"}}`
const key2 = `{"version":3,"id":"1bc823af-210a-4143-8eb4-306c19485622","address":"d95bd5b44a1290c91a31af1114e49b961e56b03b","Crypto":{"ciphertext":"0eb788eee71b9796390d6b3495c25d87746a7c2ddd98a641b90e7271231f6df0","cipherparams":{"iv":"35decc930518c37116e8fcf1f9933948"},"cipher":"aes-128-ctr","kdf":"scrypt","kdfparams":{"dklen":32,"salt":"a3b2f03dc27ee3c89b6d21b1e0d1973bb130524e7570bd0bf4da531313df9730","n":8192,"r":8,"p":1},"mac":"ef73d708b39832e6be65db73d088ab60a2c7711189a330324b410bc64c6bfe7a"}}`

func loadAccounts() (*bind.TransactOpts, *bind.TransactOpts) {
	auth1, err := bind.NewTransactor(strings.NewReader(key1), "password1")
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}

	auth2, err := bind.NewTransactor(strings.NewReader(key2), "password2")
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}

	return auth1, auth2
}

func randomBytes32() [32]byte {
	matchId := [32]byte{}
	_, err := rand.Read(matchId[:])
	if err != nil {
		panic(err)
	}
	return matchId
}

func hexToBytes32(str string) [32]byte {
	matchID, err := hex.DecodeString(str)
	if err != nil {
		panic(err)
	}

	var matchID32 [32]byte
	for i := range matchID {
		matchID32[i] = matchID[i]
	}

	return matchID32
}

var _ = Describe("AtomicSwapEther", func() {

	// It("works with bitcoin", func() {

	// 	auth1, _ := loadAccounts()
	// 	client := go_eth.Ropsten("http://13.54.129.55:8180")
	// 	address := common.HexToAddress("0x32Dad9E9Fe2A3eA2C2c643675A7d2A56814F554f")
	// 	connection1 := go_eth.NewEtherConnection(client, auth1, address)

	// 	// Account1 creates a secret lock and starts the atomic swap on Bitcoin
	// 	secretHash, err := hex.DecodeString("2c1b293ab8e578a96a6c92e85e4a6d6a19dc9aa1240df6a815fb4a8cfcd86228")
	// 	if err != nil {
	// 			panic(err)
	// 	}

	// 	var secretHash32 [32]byte
	// 	for i := range secretHash {
	// 		secretHash32[i] = secretHash[i]
	// 	}

	// value := big.NewInt(0).Mul(ether, big.NewInt(1))

	// 	// Account2 takes the hash from bitcoin and uses it to lock up Ether
	// 	tx, err := connection1.Open(matchId, common.HexToAddress("0xaAC4B896eC41e2672D2e1E5fbDe24119f4937E59"), 0, secretHash32, value)
	// 	if err != nil {
	// 		log.Fatalf("Failed to open Atomic Swap: %v", err)
	// 	}
	// 	bind.WaitMined(context.Background(), client, tx)
	// })

	It("can call retrieveSecretKey", func() {

		auth1, _ := loadAccounts()
		client := go_eth.Ropsten("http://13.54.129.55:8180")
		address := common.HexToAddress("0x32Dad9E9Fe2A3eA2C2c643675A7d2A56814F554f")
		connection1 := go_eth.NewEtherConnection(client, auth1, address)

		// Account1 creates a secret lock and starts the atomic swap on Bitcoin
		matchID32 := hexToBytes32("c2e61d599192b9aa1269ec98d9464170f9af3622bdb37e5f262edcbc9386b8e8")

		// Account2 retrieves secret
		retSecret, err := connection1.RetrieveSecretKey(matchID32)
		if err != nil {
			log.Fatalf("Failed to retrieve secret: %v", err)
		}
		println("!!!")
		println(hex.EncodeToString(retSecret))
		// Ω(retSecret).Should(Not(Equal("secret")))
	})

	// It("should work", func() {

	// 	auth1, auth2 := loadAccounts()

	// 	// client := go_eth.Simulated(auth1, auth2)
	// 	// address := go_eth.DeployEther(client, auth1)
	// 	client := go_eth.Ropsten("http://13.54.129.55:8180")
	// 	// Contract address
	// 	address := common.HexToAddress("0x32Dad9E9Fe2A3eA2C2c643675A7d2A56814F554f")

	// 	// Set up two connections
	// 	connection1 := go_eth.NewEtherConnection(client, auth1, address)
	// 	connection2 := go_eth.NewEtherConnection(client, auth2, address)

	// 	// Account1 creates a secret lock and starts the atomic swap on Bitcoin
	// 	secret := []byte("this is the secret")
	// 	secretHash := sha256.Sum256(secret)
	// 	matchId := [32]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
	// 	value := big.NewInt(0).Mul(ether, big.NewInt(1))

	// 	// Account2 takes the hash from bitcoin and uses it to lock up Ether
	// 	tx, err := connection2.Open(matchId, auth1.From, 0, secretHash, value)
	// 	if err != nil {
	// 		log.Fatalf("Failed to open Atomic Swap: %v", err)
	// 	}
	// 	bind.WaitMined(context.Background(), client, tx)
	// 	// client.Commit()

	// 	// Account1 checks that hash is what it should be
	// 	check, err := connection1.Check(matchId)
	// 	Ω(check.SecretLock).Should(Equal(secretHash))

	// 	// Account1 reveals secret to withdraw Ether
	// 	tx, err = connection1.Close(matchId, secret)
	// 	if err != nil {
	// 		log.Fatalf("Failed to close Atomic Swap: %v", err)
	// 	}
	// 	bind.WaitMined(context.Background(), client, tx)
	// 	// client.Commit()

	// 	// Account2 retrieves secret
	// 	retSecret, err := connection2.RetrieveSecretKey(matchId)
	// 	if err != nil {
	// 		log.Fatalf("Failed to retrieve secret: %v", err)
	// 	}
	// 	Ω(retSecret).Should(Equal(secret))

	// })

})
