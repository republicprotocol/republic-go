package go_eth_test

import (
	"context"
	"log"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"crypto/sha256"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/republicprotocol/go-eth"
)

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

var _ = Describe("AtomicSwapEther", func() {

	It("should work", func() {

		auth1, auth2 := loadAccounts()

		// client := go_eth.Simulated(auth1, auth2)
		// address := go_eth.DeployEther(client, auth1)
		client := go_eth.Ropsten()
		address := common.HexToAddress("0x1510deefb6d61714ead883905c0649b5aa66bd92")

		// Set up two connections
		connection1 := go_eth.NewEtherConnection(client, auth1, address)
		connection2 := go_eth.NewEtherConnection(client, auth2, address)

		// Account1 creates a secret lock and starts the atomic swap on Bitcoin
		secret := []byte("this is the secret")
		secretHash := sha256.Sum256(secret)

		// Account2 takes the hash from bitcoin and uses it to lock up Ether
		matchId, tx, err := connection2.Open(auth2.From, 0, secretHash)
		if err != nil {
			log.Fatalf("Failed to open Atomic Swap: %v", err)
		}
		bind.WaitMined(context.Background(), client, tx)
		// client.Commit()

		// Account1 checks that hash is what it should be
		check, err := connection1.Check(matchId)
		Ω(check.SecretLock).Should(Equal(secretHash))

		// Account1 reveals secret to withdraw Ether
		tx, err = connection1.Close(matchId, secret)
		if err != nil {
			log.Fatalf("Failed to close Atomic Swap: %v", err)
		}
		bind.WaitMined(context.Background(), client, tx)
		// client.Commit()

		// Account2 retrieves secret
		retSecret, err := connection2.RetrieveSecretKey(matchId)
		if err != nil {
			log.Fatalf("Failed to retrieve secret: %v", err)
		}
		Ω(retSecret).Should(Equal(secret))

		// state, err := connection2.GetState(matchId)
		// if err != nil {
		// 	log.Fatalf("Failed to get Swap state: %v", err)
		// }

		// Ω(state).Should(Equal(uint8(2)))
	})

	// It("can expire an order", func() {

	// 	auth1, auth2 := loadAccounts()

	// 	// client := go_eth.Simulated(auth1, auth2)
	// 	// address := go_eth.DeployEther(client, auth1)
	// 	client := go_eth.Ropsten()
	// 	address := common.HexToAddress("0x1510deefb6d61714ead883905c0649b5aa66bd92")

	// 	// Set up two connections
	// 	connection1 := go_eth.NewEtherConnection(client, auth1, address)
	// 	connection2 := go_eth.NewEtherConnection(client, auth2, address)

	// 	// Account1 creates a secret lock and starts the atomic swap on Bitcoin
	// 	secret := []byte("this is the secret")
	// 	secretHash := sha256.Sum256(secret)

	// 	// Account2 takes the hash from bitcoin and uses it to lock up Ether
	// 	matchId, _, err := connection2.Open(auth2.From, 0, secretHash)
	// 	if err != nil {
	// 		log.Fatalf("Failed to open Atomic Swap: %v", err)
	// 	}
	// 	// client.Commit()

	// 	// client.AdjustTime(48 * time.Hour)
	// 	// client.Commit()

	// 	_, err = connection2.Expire(matchId)
	// 	if err != nil {
	// 		log.Fatalf("Failed to expire Atomic Swap: %v", err)
	// 	}
	// 	// client.Commit()

	// 	// Account1 checks that hash is what it should be
	// 	check, err := connection1.Check(matchId)
	// 	Ω(check.SecretLock).Should(Equal(secretHash))

	// 	// Account1 reveals secret to withdraw Ether
	// 	connection1.Close(matchId, secret)
	// 	// client.Commit()

	// 	// Account2 retrieves secret
	// 	retSecret, err := connection2.RetrieveSecretKey(matchId)
	// 	if err != nil {
	// 		log.Fatalf("Failed to retrieve secret: %v", err)
	// 	}
	// 	// Ω(retSecret).Should(Equal(secret))
	// 	Ω(retSecret).Should(Equal([]byte("")))

	// 	// state, err := connection2.GetState(matchId)
	// 	// if err != nil {
	// 	// 	log.Fatalf("Failed to get Swap state: %v", err)
	// 	// }
	// 	// Ω(state).Should(Equal(uint8(3)))
	// })

})
