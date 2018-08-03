package oracle_test

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/oracle"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/leveldb"
	"github.com/republicprotocol/republic-go/swarm"
	"github.com/republicprotocol/republic-go/testutils"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var _ = Describe("Oracle", func() {

	var (
		numberOfTester         = 100
		numberOfBootsrapTester = 5
		α                      = 5
	)

	Context("when sending and receiving new oracle prices", func() {
		var hub map[identity.Address]Server
		var testers []Tester
		var RenOracler Client
		var RenOraclerKey crypto.EcdsaKey
		var RenAddress identity.Address
		var err error

		BeforeEach(func() {
			hub = map[identity.Address]Server{}
			testers = make([]Tester, numberOfTester)

			// Create the RenOracler
			RenOraclerKey, err = crypto.RandomEcdsaKey()
			Expect(err).ShouldNot(HaveOccurred())
			RenAddress = identity.Address(RenOraclerKey.Address())
			Expect(err).ShouldNot(HaveOccurred())
			RenOracler, _, err = testutils.NewMockOracleClient(RenAddress, hub)
			Expect(err).ShouldNot(HaveOccurred())

			// Create tester
			for i := range testers {
				testers[i] = newTester(α, hub, RenAddress)
			}

			// Connect each other
			for i := range testers {
				for j := range testers {
					Expect(testers[i].MultiStorer.PutMultiAddress(testers[j].Multi)).ShouldNot(HaveOccurred())
				}
			}

		})

		AfterEach(func() {
			os.RemoveAll("./tmp")
		})

		It("should be able to verify and accept the right prices", func() {
			ctx, cancelCtx := context.WithCancel(context.Background())
			defer cancelCtx()

			for times := 0; times < 5; times++ {
				price := testutils.RandMidpointPrice()
				price.Signature, err = RenOraclerKey.Sign(price.Hash())
				Expect(err).ShouldNot(HaveOccurred())

				for i := 0; i < numberOfBootsrapTester; i++ {
					err = RenOracler.UpdateMidpoint(ctx, testers[i].Multi, price)
					Expect(err).ShouldNot(HaveOccurred())
				}

				receivedTester := 0
				for i := range testers {
					storedPrice, err := testers[i].MidPointPriceStore.MidpointPrice()
					Expect(err).ShouldNot(HaveOccurred())
					if storedPrice.Equals(price) {
						receivedTester++
					}
				}
				log.Printf("%d out of %d tester received the new midpoint price", receivedTester, numberOfTester)
				Expect(receivedTester).Should(BeNumerically(">=", numberOfTester*9/10))

				// Update the nonce
				time.Sleep(1 * time.Second)
			}
		})
	})

})

type Tester struct {
	Multi       identity.MultiAddress
	Key         crypto.EcdsaKey
	Hub         map[identity.Address]Server
	α           int
	MultiStorer swarm.MultiAddressStorer

	OracleClient       Client
	Oracler            Oracler
	OracleServer       Server
	MidPointPriceStore MidpointPriceStorer
}

func newTester(α int, hub map[identity.Address]Server, oracleAddress identity.Address) Tester {

	key, err := crypto.RandomEcdsaKey()
	addr := identity.Address(key.Address())
	Expect(err).ShouldNot(HaveOccurred())
	multi, err := addr.MultiAddress()
	Expect(err).ShouldNot(HaveOccurred())
	multi.Signature, err = key.Sign(multi.Hash())
	Expect(err).ShouldNot(HaveOccurred())

	// Create leveldb store and store own multiAddress.
	db, err := leveldb.NewStore(fmt.Sprintf("./tmp/swarmer-%v.out", key.Address()), 72*time.Hour)
	Expect(err).ShouldNot(HaveOccurred())

	if err != nil {
		return Tester{}
	}
	multiStore := db.SwarmMultiAddressStore()

	oracleClient, midPointPriceStore, err := testutils.NewMockOracleClient(addr, hub)
	Expect(err).ShouldNot(HaveOccurred())
	oracler := NewOracler(oracleClient, &key, multiStore, α)
	server := NewServer(oracler, oracleAddress, multiStore, midPointPriceStore, α)
	hub[addr] = server

	return Tester{
		Multi:       multi,
		Key:         key,
		Hub:         hub,
		α:           α,
		MultiStorer: multiStore,

		OracleClient:       oracleClient,
		Oracler:            oracler,
		OracleServer:       server,
		MidPointPriceStore: midPointPriceStore,
	}
}
