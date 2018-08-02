package oracle_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/leveldb"
	. "github.com/republicprotocol/republic-go/oracle"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/swarm"
	"github.com/republicprotocol/republic-go/testutils"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var _ = Describe("Oracle", func() {

	var (
		numberOfTester        = 10
		numberOfBootrapTester = 5
		α                     = 3
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

			price := randMidpointPrice()
			price.Signature, err = RenOraclerKey.Sign(price.Hash())
			Expect(err).ShouldNot(HaveOccurred())

			for i := 0; i < numberOfBootrapTester; i++ {
				err = RenOracler.UpdateMidpoint(ctx, testers[i].Multi, price)
				Expect(err).ShouldNot(HaveOccurred())
			}

			for i := range testers {
				storedPrice, err := testers[i].MidPointPriceStore.MidpointPrice()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(storedPrice.Equals(price)).Should(BeTrue())
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
