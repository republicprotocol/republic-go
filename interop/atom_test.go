package atom_test

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/atom/bitcoin"
	"github.com/republicprotocol/republic-go/atom/ethereum"
	// . "github.com/republicprotocol/republic-go/atom"
)

func randomBytes32() []byte {
	randString := [32]byte{}
	_, err := rand.Read(randString[:])
	if err != nil {
		panic(err)
	}
	return randString[:]
}

var ether = big.NewInt(1000000000000000000)

const key1 = `{"version":3,"id":"7844982f-abe7-4690-8c15-34f75f847c66","address":"db205ea9d35d8c01652263d58351af75cfbcbf07","Crypto":{"ciphertext":"378dce3c1279b36b071e1c7e2540ac1271581bff0bbe36b94f919cb73c491d3a","cipherparams":{"iv":"2eb92da55cc2aa62b7ffddba891f5d35"},"cipher":"aes-128-ctr","kdf":"scrypt","kdfparams":{"dklen":32,"salt":"80d3341678f83a14024ba9c3edab072e6bd2eea6aa0fbc9e0a33bae27ffa3d6d","n":8192,"r":8,"p":1},"mac":"3d07502ea6cd6b96a508138d8b8cd2e46c3966240ff276ce288059ba4235cb0d"}}`
const key2 = `{"version":3,"id":"1bc823af-210a-4143-8eb4-306c19485622","address":"d95bd5b44a1290c91a31af1114e49b961e56b03b","Crypto":{"ciphertext":"0eb788eee71b9796390d6b3495c25d87746a7c2ddd98a641b90e7271231f6df0","cipherparams":{"iv":"35decc930518c37116e8fcf1f9933948"},"cipher":"aes-128-ctr","kdf":"scrypt","kdfparams":{"dklen":32,"salt":"a3b2f03dc27ee3c89b6d21b1e0d1973bb130524e7570bd0bf4da531313df9730","n":8192,"r":8,"p":1},"mac":"ef73d708b39832e6be65db73d088ab60a2c7711189a330324b410bc64c6bfe7a"}}`

func loadAccounts() (*bind.TransactOpts, *bind.TransactOpts) {
	auth1, err := bind.NewTransactor(strings.NewReader(key1), "password1")
	Ω(err).Should(BeNil())

	auth2, err := bind.NewTransactor(strings.NewReader(key2), "password2")
	Ω(err).Should(BeNil())

	return auth1, auth2
}

var _ = Describe("Atom", func() {
	It("should work between Bitcoin and Ethereum", func() {
		// Alice
		secret := randomBytes32()
		hashLock := sha256.Sum256(secret)
		BTCAtom := NewBTCAtomContract("testuser", "testpassword", "testnet")
		err := BTCAtom.Initiate(hashLock[:], []byte("mgTCJazbqe8JUCNQTbcVLJDv5yseRfAMVe"), []byte("mv8p79yFBUfrbWCSMPc4fNTThZS1zdPpR6"), big.NewInt(10000000), time.Now().Unix()+10000)
		Ω(err).Should(BeNil())

		// Bob
		auth1, auth2 := loadAccounts()
		client := ethereum.Ropsten("https://ropsten.infura.io/")
		contractAddress := common.HexToAddress("0xbd59e72598737a08a68fe192df04c773adbbfa53")
		user2Connection, err := ethereum.NewETHAtomContract(context.Background(), client, auth2, contractAddress, nil)
		Ω(err).Should(BeNil())
		err = user2Connection.Initiate(hashLock[:], auth1.From.Bytes(), auth2.From.Bytes(), ether, time.Now().Add(48*time.Hour).Unix())
		Ω(err).Should(BeNil())

		// Alice
		user1Connection, err := ethereum.NewETHAtomContract(context.Background(), client, auth1, contractAddress, user2Connection.GetData())
		Ω(err).Should(BeNil())
		err = user1Connection.Redeem(secret)
		Ω(err).Should(BeNil())

		// Bob
		// Account2 retrieves secret
		retSecret, err := user2Connection.ReadSecret()
		Ω(err).Should(BeNil())
		Ω(retSecret).Should(Equal(secret))
		// Redeem Bitcoins
		BTCAtomBOB := NewBTCAtomContract("testuser", "testpassword", "testnet")
		err = BTCAtomBOB.Redeem(retSecret)
		Ω(err).Should(BeNil())
	})
})
