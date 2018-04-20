package arc_bitcoin_test

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	rpc "github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/bitcoin/client"
	"github.com/republicprotocol/republic-go/bitcoin/contracts/arc-bitcoin"
	"github.com/republicprotocol/republic-go/bitcoin/regtest"
)

const CHAIN = "regtest"

func randomBytes32() []byte {
	randString := [32]byte{}
	_, err := rand.Read(randString[:])
	if err != nil {
		panic(err)
	}
	return randString[:]
}

var _ = Describe("Bitcoin", func() {
	var rpcClient *rpc.Client

	var aliceAddr, bobAddr btcutil.Address

	BeforeSuite(func() {
		var err error
		rpcClient, err = client.ConnectToRPC(&chaincfg.RegressionNetParams, "testuser", "testpassword")
		Expect(err).ShouldNot(HaveOccurred())

		regtest.Mine(rpcClient)

		aliceAddr, err = regtest.NewAccount(rpcClient, "alice", 1000000000)
		Expect(err).ShouldNot(HaveOccurred())

		bobAddr, err = regtest.NewAccount(rpcClient, "bob", 1000000000)
		Expect(err).ShouldNot(HaveOccurred())

		fmt.Println("Alice")
		fmt.Println(aliceAddr.String())
		fmt.Println("Bob")
		fmt.Println(bobAddr.String())
	})

	AfterSuite(func() {
		rpcClient.Shutdown()
		rpcClient.WaitForShutdown()
	})

	It("can initiate a bitcoin atomic swap", func() {
		secret := randomBytes32()
		hashLock := sha256.Sum256(secret)
		BTCAtom := arc_bitcoin.NewBitcoinArc("testuser", "testpassword", CHAIN)
		err := BTCAtom.Initiate(hashLock, []byte(aliceAddr.String()), []byte(bobAddr.String()), big.NewInt(3000000), time.Now().Unix()+10000)
		Ω(err).Should(BeNil())
	})

	It("can redeem a bitcoin atomic swap with correct secret", func() {
		secret := randomBytes32()
		hashLock := sha256.Sum256(secret)
		BTCAtom := arc_bitcoin.NewBitcoinArc("testuser", "testpassword", CHAIN)
		err := BTCAtom.Initiate(hashLock, []byte(aliceAddr.String()), []byte(bobAddr.String()), big.NewInt(3000000), time.Now().Unix()+10000)
		Ω(err).Should(BeNil())
		err = BTCAtom.Redeem(secret)
		Ω(err).Should(BeNil())
	})

	It("can not redeem a bitcoin atomic swap with a wrong secret", func() {
		secret := randomBytes32()
		wrongSecret := randomBytes32()
		hashLock := sha256.Sum256(secret)
		BTCAtom := arc_bitcoin.NewBitcoinArc("testuser", "testpassword", CHAIN)
		err := BTCAtom.Initiate(hashLock, []byte(aliceAddr.String()), []byte(bobAddr.String()), big.NewInt(3000000), time.Now().Unix()+10000)
		Ω(err).Should(BeNil())
		err = BTCAtom.Redeem(wrongSecret)
		Ω(err).Should(Not(BeNil()))
	})

	It("can read a bitcoin atomic swap", func() {
		secret := randomBytes32()
		hashLock := sha256.Sum256(secret)
		BTCAtom := arc_bitcoin.NewBitcoinArc("testuser", "testpassword", CHAIN)
		to := []byte(aliceAddr.String())
		from := []byte(bobAddr.String())
		value := big.NewInt(1000000)
		expiry := time.Now().Unix() + 10000
		err := BTCAtom.Initiate(hashLock, to, from, value, expiry)
		Ω(err).Should(BeNil())
		readHashLock, readTo, _, readValue, readExpiry, readErr := BTCAtom.Audit()
		Ω(readErr).Should(BeNil())
		Ω(readHashLock).Should(Equal(hashLock[:]))
		Ω(readTo).Should(Equal(to))
		Ω(readValue).Should(Equal(value))
		Ω(readExpiry).Should(Equal(expiry))
	})

	It("can read the correct secret from a bitcoin atomic swap", func() {
		secret := randomBytes32()
		hashLock := sha256.Sum256(secret)
		BTCAtom := arc_bitcoin.NewBitcoinArc("testuser", "testpassword", CHAIN)
		err := BTCAtom.Initiate(hashLock, []byte(aliceAddr.String()), []byte(bobAddr.String()), big.NewInt(3000000), time.Now().Unix()+10000)
		Ω(err).Should(BeNil())
		err = BTCAtom.Redeem(secret)
		Ω(err).Should(BeNil())
		readSecret, err := BTCAtom.AuditSecret()
		Ω(err).Should(BeNil())
		Ω(readSecret).Should(Equal(secret[:]))
	})

	It("can not refund a bitcoin atomic swap before expiry", func() {
		secret := randomBytes32()
		hashLock := sha256.Sum256(secret)
		BTCAtom := arc_bitcoin.NewBitcoinArc("testuser", "testpassword", CHAIN)
		err := BTCAtom.Initiate(hashLock, []byte(aliceAddr.String()), []byte(bobAddr.String()), big.NewInt(3000000), time.Now().Unix()+10000)
		Ω(err).Should(BeNil())
		err = BTCAtom.Refund()
		Ω(err).Should(Not(BeNil()))
	})

	It("can refund a bitcoin atomic swap", func() {
		secret := randomBytes32()
		hashLock := sha256.Sum256(secret)
		BTCAtom := arc_bitcoin.NewBitcoinArc("testuser", "testpassword", CHAIN)
		err := BTCAtom.Initiate(hashLock, []byte(aliceAddr.String()), []byte(bobAddr.String()), big.NewInt(3000000), time.Now().Unix()+600)
		Ω(err).Should(BeNil())
		time.Sleep(30 * time.Minute)
		err = BTCAtom.Refund()
		Ω(err).Should(BeNil())
	})
})
