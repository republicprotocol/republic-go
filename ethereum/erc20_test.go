package ethereum_test

import (
	"context"
	"crypto/sha256"
	"math/big"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/republicprotocol/go-atom/ethereum"
)

var _ = Describe("ERC20", func() {

	It("works", func() {
		// Load client and accounts
		client, auth1, auth2 := loadClient()
		contractAddress := erc20AtomAddress(client, auth1)
		erc20Address := erc20Address(client, auth1, auth2)

		// ALICE has locked up tokens on another ledger with the following secret
		secret := []byte("this is the secret")
		secretHash := sha256.Sum256(secret)

		// BOB reciprocates the atomic swap with the secretHash
		user2Connection, _ := ethereum.NewERC20AtomContract(context.Background(), client, auth2, contractAddress, erc20Address, nil)
		value := big.NewInt(0).Mul(ether, big.NewInt(1))
		err := user2Connection.Initiate(secretHash[:], auth1.From.Bytes(), auth2.From.Bytes(), value, time.Now().Add(48*time.Hour).Unix())
		Ω(err).Should(BeNil())

		// ALICE checks that Bob has set up the swap correctly
		user1Connection, _ := ethereum.NewERC20AtomContract(context.Background(), client, auth1, contractAddress, erc20Address, user2Connection.GetData())
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

})
