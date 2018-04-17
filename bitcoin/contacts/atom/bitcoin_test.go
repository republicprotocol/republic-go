package arc_bitcoin_test

import (
	"crypto/rand"
	"crypto/sha256"
	"log"
	"math/big"
	"time"

	"github.com/btcsuite/btcd/rpcclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/bitcoin/contracts/atom"
)

func randomBytes32() []byte {
	randString := [32]byte{}
	_, err := rand.Read(randString[:])
	if err != nil {
		panic(err)
	}
	return randString[:]
}

func setupTestNode() {
	// create new client instance
	client, err := rpcclient.New(&rpcclient.ConnConfig{
		HTTPPostMode: true,
		DisableTLS:   true,
		Host:         "127.0.0.1:8332",
		User:         "testuser",
		Pass:         "testpassword",
	}, nil)

	if err != nil {
		log.Fatalf("error creating new btc client: %v", err)
	}

	err = client.Connect(10)
	if err != nil {
		log.Fatalf("error connecting to new btc client: %v", err)
	}

	// list accounts
	_, err = client.ListAccounts()
	if err != nil {
		log.Fatalf("error listing accounts: %v", err)
	}

}

var _ = Describe("Bitcoin", func() {

	setupTestNode()

	It("can initiate a bitcoin atomic swap", func() {
		secret := randomBytes32()
		hashLock := sha256.Sum256(secret)
		BTCAtom := NewBTCAtomContract("testuser", "testpassword", "mainnet")
		err := BTCAtom.Initiate(hashLock[:], []byte("mgTCJazbqe8JUCNQTbcVLJDv5yseRfAMVe"), []byte("mv8p79yFBUfrbWCSMPc4fNTThZS1zdPpR6"), big.NewInt(1000000), time.Now().Unix()+10000)
		Ω(err).Should(BeNil())
	})

	It("can redeem a bitcoin atomic swap", func() {
		secret := randomBytes32()
		hashLock := sha256.Sum256(secret)
		BTCAtom := NewBTCAtomContract("testuser", "testpassword", "mainnet")
		err := BTCAtom.Initiate(hashLock[:], []byte("mgTCJazbqe8JUCNQTbcVLJDv5yseRfAMVe"), []byte("mv8p79yFBUfrbWCSMPc4fNTThZS1zdPpR6"), big.NewInt(1000000), time.Now().Unix()+10000)
		Ω(err).Should(BeNil())
		err = BTCAtom.Redeem(secret)
		Ω(err).Should(BeNil())
	})

	It("can read a bitcoin atomic swap", func() {
		secret := randomBytes32()
		hashLock := sha256.Sum256(secret)
		BTCAtom := NewBTCAtomContract("testuser", "testpassword", "mainnet")
		to := []byte("mgTCJazbqe8JUCNQTbcVLJDv5yseRfAMVe")
		from := []byte("mv8p79yFBUfrbWCSMPc4fNTThZS1zdPpR6")
		value := big.NewInt(1000000)
		expiry := time.Now().Unix() + 10000
		err := BTCAtom.Initiate(hashLock[:], to, from, value, expiry)
		Ω(err).Should(BeNil())
		readHashLock, readTo, readFrom, readValue, readExpiry, readErr := BTCAtom.Read()
		Ω(readErr).Should(BeNil())
		Ω(readHashLock).Should(Equal(hashLock[:]))
		Ω(readTo).Should(Equal(to))
		Ω(readFrom).Should(Equal(from))
		Ω(readValue).Should(Equal(value))
		Ω(readExpiry).Should(Equal(expiry))
	})

	It("can read secret a bitcoin atomic swap", func() {
		secret := randomBytes32()
		hashLock := sha256.Sum256(secret)
		BTCAtom := NewBTCAtomContract("testuser", "testpassword", "mainnet")
		err := BTCAtom.Initiate(hashLock[:], []byte("mgTCJazbqe8JUCNQTbcVLJDv5yseRfAMVe"), []byte("mv8p79yFBUfrbWCSMPc4fNTThZS1zdPpR6"), big.NewInt(1000000), time.Now().Unix()+10000)
		Ω(err).Should(BeNil())
		err = BTCAtom.Redeem(secret)
		Ω(err).Should(BeNil())
		readSecret, err := BTCAtom.ReadSecret()
		Ω(err).Should(BeNil())
		Ω(readSecret).Should(Equal(secret[:]))
	})

	It("can refund a bitcoin atomic swap", func() {
		secret := randomBytes32()
		hashLock := sha256.Sum256(secret)
		BTCAtom := NewBTCAtomContract("testuser", "testpassword", "mainnet")
		err := BTCAtom.Initiate(hashLock[:], []byte("mgTCJazbqe8JUCNQTbcVLJDv5yseRfAMVe"), []byte("mv8p79yFBUfrbWCSMPc4fNTThZS1zdPpR6"), big.NewInt(1000000), time.Now().Unix()+1000)
		Ω(err).Should(BeNil())
		err = BTCAtom.Refund()
		Ω(err).Should(BeNil())
	})
})
