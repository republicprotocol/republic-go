package bitcoin_test

import (
	"crypto/rand"
	"crypto/sha256"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/republicprotocol/go-atom/bitcoin"
)

func randomBytes32() []byte {
	randString := [32]byte{}
	_, err := rand.Read(randString[:])
	if err != nil {
		panic(err)
	}
	return randString[:]
}

var _ = Describe("Bitcoin", func() {

	secret := randomBytes32()
	hashLock := sha256.Sum256(secret)

	It("can initiate a bitcoin atomic swap", func() {

		BTCAtom := NewBTCAtomContract("testuser", "testpassword", "testnet")
		err := BTCAtom.Initiate(hashLock[:], []byte("mgTCJazbqe8JUCNQTbcVLJDv5yseRfAMVe"), []byte("mv8p79yFBUfrbWCSMPc4fNTThZS1zdPpR6"), 10000000, 1518570598)
		Ω(err).Should(BeNil())

	})

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

	// It("can call retrieveSecretKey", func() {

	// 	auth1, _ := loadAccounts()
	// 	client := go_eth.Ropsten("http://13.54.129.55:8180")
	// 	address := common.HexToAddress("0x32Dad9E9Fe2A3eA2C2c643675A7d2A56814F554f")
	// 	connection1 := go_eth.NewEtherConnection(client, auth1, address)

	// 	// Account1 creates a secret lock and starts the atomic swap on Bitcoin
	// 	matchID32 := hexToBytes32("c2e61d599192b9aa1269ec98d9464170f9af3622bdb37e5f262edcbc9386b8e8")

	// 	// Account2 retrieves secret
	// 	retSecret, err := connection1.RetrieveSecretKey(matchID32)
	// 	if err != nil {
	// 		log.Fatalf("Failed to retrieve secret: %v", err)
	// 	}
	// 	println("!!!")
	// 	println(hex.EncodeToString(retSecret))
	// 	// Ω(retSecret).Should(Not(Equal("secret")))
	// })

	// It("should work", func() {

	// 	auth1, auth2 := loadAccounts()
	// 	client := ethereum.Ropsten("http://13.54.129.55:8180")
	// 	// Contract address
	// 	address := common.HexToAddress("0x32Dad9E9Fe2A3eA2C2c643675A7d2A56814F554f")

	// 	// Set up two connections
	// 	connection1 := ethereum.NewETHAtomContract(client, auth1, address, nil)

	// 	// Account1 creates a secret lock and starts the atomic swap on Bitcoin
	// 	secret := []byte("this is the secret")
	// 	secretHash := sha256.Sum256(secret)
	// 	value1 := big.NewInt(0).Mul(ether, big.NewInt(1))
	// 	err := connection1.Initiate(secretHash[:], auth1.From.Bytes(), auth2.From.Bytes(), value1, time.Now().Add(48*time.Hour).Unix())
	// 	if err != nil {
	// 		log.Fatalf("Failed to open Atomic Swap: %v", err)
	// 	}

	// 	connection2 := ethereum.NewETHAtomContract(client, auth2, address, connection1.GetData())
	// 	// Account2 checks that hash is what it should be
	// 	_, _, _, _, _, err = connection2.Read()
	// 	if err != nil {
	// 		log.Fatalf("Failed: %v", err)
	// 	}
	// 	// Ω(hash).Should(Equal(secretHash))

	// 	// Account2 takes the hash from bitcoin and uses it to lock up Ether
	// 	value2 := big.NewInt(0).Mul(ether, big.NewInt(1))
	// 	err = connection2.Initiate(secretHash[:], auth2.From.Bytes(), auth1.From.Bytes(), value2, time.Now().Add(24*time.Hour).Unix())
	// 	if err != nil {
	// 		log.Fatalf("Failed to open Atomic Swap: %v", err)
	// 	}

	// 	// Account1 reveals secret to withdraw Ether
	// 	err = connection1.Redeem(secret)
	// 	if err != nil {
	// 		log.Fatalf("Failed to close Atomic Swap: %v", err)
	// 	}

	// 	// Account2 retrieves secret
	// 	retSecret, err := connection2.ReadSecret()
	// 	if err != nil {
	// 		log.Fatalf("Failed to retrieve secret: %v", err)
	// 	}
	// 	Ω(retSecret).Should(Equal(secret))

	// })

})
