package arc_test

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"math/big"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/republicprotocol/republic-go/blockchain/bitcoin"
	btcarc "github.com/republicprotocol/republic-go/blockchain/bitcoin/arc"
	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	etharc "github.com/republicprotocol/republic-go/blockchain/ethereum/arc"
	"github.com/republicprotocol/republic-go/blockchain/test"
	"github.com/republicprotocol/republic-go/blockchain/test/bitcoind"
	"github.com/republicprotocol/republic-go/blockchain/test/ganache"
)

const CHAIN = "regtest"
const RPC_USERNAME = "testuser"
const RPC_PASSWORD = "testpassword"

var _ = Describe("ARC", func() {

	var aliceEthAcc, bobEthAcc *bind.TransactOpts
	var aliceEthAddr, bobEthAddr common.Address
	var ethConnection ethereum.Conn
	var swapID [32]byte
	var btcValue, ethValue *big.Int
	var validity int64

	var aliceBtcConnection, bobBtcConnection bitcoin.Conn
	var aliceBtcAddr, bobBtcAddr string // btcutil.Address

	// Don't run on CI
	test.Context("BTC-ETH", func() {

		BeforeSuite(func() {

			{ // BITCOIN
				var err error
				minerBtcConnection, err := bitcoin.Connect("regtest", RPC_USERNAME, RPC_PASSWORD)
				Expect(err).ShouldNot(HaveOccurred())

				_aliceBtcAddr, err := bitcoind.NewAccount(minerBtcConnection, "alice", 100000000)
				Expect(err).ShouldNot(HaveOccurred())
				aliceBtcAddr = _aliceBtcAddr.EncodeAddress()
				_bobBtcAddr, err := bitcoind.NewAccount(minerBtcConnection, "bob", 0)
				Expect(err).ShouldNot(HaveOccurred())
				bobBtcAddr = _bobBtcAddr.EncodeAddress()

				aliceBtcConnection, err = bitcoin.Connect("regtest", RPC_USERNAME, RPC_PASSWORD)
				Expect(err).ShouldNot(HaveOccurred())
				bobBtcConnection, err = bitcoin.Connect("regtest", RPC_USERNAME, RPC_PASSWORD)
				Expect(err).ShouldNot(HaveOccurred())

				go func() {
					defer minerBtcConnection.Shutdown()
					err = bitcoind.Mine(minerBtcConnection)
					Expect(err).ShouldNot(HaveOccurred())
				}()

				btcValue = big.NewInt(0.03 * 100000000) // 0.03 BTC
			}

			{ // ETHEREUM

				var err error
				ethConnection, err = ganache.Connect("http://localhost:8545")
				Expect(err).ShouldNot(HaveOccurred())

				aliceEthAcc, aliceEthAddr, err = ganache.NewAccount(ethConnection, big.NewInt(10))
				Expect(err).ShouldNot(HaveOccurred())
				aliceEthAcc.GasLimit = 3000000
				bobEthAcc, bobEthAddr, err = ganache.NewAccount(ethConnection, big.NewInt(10))
				Expect(err).ShouldNot(HaveOccurred())
				bobEthAcc.GasLimit = 3000000

				ethValue = big.NewInt(0.42 * 1000000000000000000) // 0.42 eth
				validity = int64((time.Hour * 24).Seconds())

				swapID[0] = 0x13
			}
		})

		AfterSuite(func() {
			aliceBtcConnection.Shutdown()
			bobBtcConnection.Shutdown()
		})

		It("can perform ETH-ETH arc swap", func() {

			var btcArcData, ethArcData []byte

			var aliceSecret [32]byte
			var secretHash [32]byte

			{ // Alice can initiate swap on ethereum
				aliceEthArc, err := etharc.NewArc(context.Background(), ethConnection, aliceEthAcc, swapID)
				Expect(err).ShouldNot(HaveOccurred())

				// Genreate random secret
				_, err = rand.Read(aliceSecret[:])
				if err != nil {
					panic(err)
				}
				secretHash = sha256.Sum256(aliceSecret[:])

				err = aliceEthArc.Initiate(secretHash, aliceEthAddr.Bytes(), bobEthAddr.Bytes(), ethValue, time.Now().Unix()+10000)
				Expect(err).ShouldNot(HaveOccurred())
				ethArcData, err = aliceEthArc.Serialize()
				Expect(err).ShouldNot(HaveOccurred())
			}

			{ // Bob audits Alice's ethereum contract, uploads his bitcoin script
				bobEthArc, err := etharc.NewArc(context.Background(), ethConnection, bobEthAcc, swapID)
				Expect(err).ShouldNot(HaveOccurred())
				err = bobEthArc.Deserialize(ethArcData)
				Expect(err).ShouldNot(HaveOccurred())

				_secretHash, _, _to, _value, _, err := bobEthArc.Audit()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(_to).Should(Equal(bobEthAddr.Bytes()))
				Expect(_value).Should(Equal(ethValue))
				Expect(_secretHash).Should(Equal(secretHash))
				// Expect(_expiry.Int64()).Should(Equal(expiry))

				bobBtcArc := btcarc.NewArc(bobBtcConnection)
				err = bobBtcArc.Initiate(_secretHash, []byte(bobBtcAddr), []byte(aliceBtcAddr), btcValue, time.Now().Unix()+validity)
				Expect(err).ShouldNot(HaveOccurred())
				btcArcData, err = bobBtcArc.Serialize()
				Expect(err).ShouldNot(HaveOccurred())
			}

			{ // Alice audits Bob's bitcoin script, redeems it with her password
				aliceBtcArc := btcarc.NewArc(aliceBtcConnection)
				err := aliceBtcArc.Deserialize(btcArcData)
				Expect(err).ShouldNot(HaveOccurred())

				_secretHash, _, _to, _value, _, err := aliceBtcArc.Audit()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(_secretHash).Should(Equal(secretHash))
				Expect(_to).Should(Equal([]byte(aliceBtcAddr)))
				Expect(_value).Should(Equal(btcValue))

				err = aliceBtcArc.Redeem(aliceSecret)
				Expect(err).ShouldNot(HaveOccurred())

				btcArcData, err = aliceBtcArc.Serialize()
				Expect(err).ShouldNot(HaveOccurred())
			}

			{ // Bob audits Alice's password on bitcoin, redeems the ethereum swap with it

				bobBtcArc := btcarc.NewArc(bobBtcConnection)
				err := bobBtcArc.Deserialize(btcArcData)
				Expect(err).ShouldNot(HaveOccurred())

				_secret, err := bobBtcArc.AuditSecret()
				Expect(err).ShouldNot(HaveOccurred())

				bobEthArc, err := etharc.NewArc(context.Background(), ethConnection, bobEthAcc, swapID)
				Expect(err).ShouldNot(HaveOccurred())
				err = bobEthArc.Deserialize(ethArcData)
				Expect(err).ShouldNot(HaveOccurred())

				err = bobEthArc.Redeem(_secret)
				Expect(err).ShouldNot(HaveOccurred())
			}
		})
	})
})
