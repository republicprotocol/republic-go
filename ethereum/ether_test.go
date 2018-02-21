package ethereum_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/republicprotocol/go-atom/ethereum"
)

var _ = Describe("Ether", func() {

	It("can swap with an arbitrary ledger", func() {
		// Load client and accounts
		client, auth1, auth2 := loadClient()
		contractAddress := etherAtomAddress(client, auth1)

		// ALICE has locked up tokens on another ledger with the following secret
		secret := []byte("this is the secret")
		secretHash := sha256.Sum256(secret)

		// BOB reciprocates the atomic swap with the secretHash
		user2Connection, _ := ethereum.NewETHAtomContract(context.Background(), client, auth2, contractAddress, nil)
		value := big.NewInt(0).Mul(ether, big.NewInt(1))
		err := user2Connection.Initiate(secretHash[:], auth1.From.Bytes(), auth2.From.Bytes(), value, time.Now().Add(48*time.Hour).Unix())
		Ω(err).Should(BeNil())

		// ALICE checks that Bob has set up the swap correctly
		user1Connection, _ := ethereum.NewETHAtomContract(context.Background(), client, auth1, contractAddress, user2Connection.GetData())
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
		// Load client and accounts
		client, auth1, auth2 := loadClient()
		contractAddress := etherAtomAddress(client, auth1)

		secret := []byte("this is the secret")
		secretHash := sha256.Sum256(secret)

		// BOB reciprocates the atomic swap with the secretHash
		user2Connection, _ := ethereum.NewETHAtomContract(context.Background(), client, auth2, contractAddress, nil)
		err := user2Connection.Initiate(secretHash[:], auth1.From.Bytes(), auth2.From.Bytes(), ether, time.Now().Add(48*time.Hour).Unix())
		Ω(err).Should(BeNil())

		//ALICE tries to reem with wrong password, and then right password
		user1Connection, _ := ethereum.NewETHAtomContract(context.Background(), client, auth1, contractAddress, user2Connection.GetData())
		wrongSecret := []byte("this is NOT the secret")
		err = user1Connection.Redeem(wrongSecret)
		// Error should NOT be nil
		Ω(err).Should(Not(BeNil()))

		err = user1Connection.Redeem(secret)
		Ω(err).Should(BeNil())
	})

	It("can't expire before timeout", func() {
		// Load client and accounts
		client, auth1, auth2 := loadClient()
		contractAddress := etherAtomAddress(client, auth1)

		secret := []byte("this is the secret")
		secretHash := sha256.Sum256(secret)

		// BOB reciprocates the atomic swap with the secretHash
		user2Connection, _ := ethereum.NewETHAtomContract(context.Background(), client, auth2, contractAddress, nil)
		err := user2Connection.Initiate(secretHash[:], auth1.From.Bytes(), auth2.From.Bytes(), ether, time.Now().Add(48*time.Hour).Unix())
		Ω(err).Should(BeNil())

		// Should NOT be able to refund
		err = user2Connection.Refund()
		Ω(err).Should(Not(BeNil()))

		// ALICE redeems by revealing secret
		user1Connection, _ := ethereum.NewETHAtomContract(context.Background(), client, auth1, contractAddress, user2Connection.GetData())
		err = user1Connection.Redeem(secret)
		Ω(err).Should(BeNil())
	})

	// It("can expire after timeout", func() {
	// 	// Load client and accounts
	// 	client, auth1, auth2 := loadClient()
	//  contractAddress := etherAtomAddress(client, auth1)

	// 	secret := []byte("this is the secret")
	// 	secretHash := sha256.Sum256(secret)

	// 	timeout := time.Now()

	// 	// BOB reciprocates the atomic swap with the secretHash
	// 	user2Connection, _ := ethereum.NewETHAtomContract(context.Background(), client, auth2, contractAddress, nil)
	// 	err := user2Connection.Initiate(secretHash[:], auth1.From.Bytes(), auth2.From.Bytes(), ether, timeout.Unix())
	// 	Ω(err).Should(BeNil())

	// 	time.Sleep(time.Second * 1)
	// 	// Should be able to refund
	// 	// auth1.GasLimit = 250000
	// 	user1Connection, _ := ethereum.NewETHAtomContract(context.Background(), client, auth1, contractAddress, user2Connection.GetData())

	// 	err = user1Connection.Refund()
	// 	Ω(err).Should(BeNil())
	// })

	It("can't expire an expired order", func() {
		hash, err := hex.DecodeString("bdcf34f235c01c96c6110b30ea4cec8b78e90304271db5cef8b1b1eb34f0f71b")
		Ω(err).Should(BeNil())

		// Load client and accounts
		client, auth, _ := loadClient()
		contractAddress := etherAtomAddress(client, auth)

		connection, _ := ethereum.NewETHAtomContract(context.Background(), client, auth, contractAddress, hash)

		// Should NOT be able to refund
		err = connection.Refund()
		Ω(err).Should(Not(BeNil()))
		// err = connection.Redeem([]byte("this is the secret"))
	})

})
