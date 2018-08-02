package oracle_test

import (
	"context"
	"math/rand"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/oracle"
	"github.com/republicprotocol/republic-go/testutils"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/swarm"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var _ = Describe("Oracle", func() {

	var (
		numberOfClients          = 50
		numberOfBootstrapClients = 5
		α                        = 3
	)

	Context("when sending and receiving new oracle prices", func() {

		AfterEach(func() {
			os.RemoveAll("./tmp")
		})

		It("should be able to verify and accept the right prices", func() {
			ctx, cancelCtx := context.WithCancel(context.Background())
			defer cancelCtx()
		})
	})

})

type Tester struct {
	Key crypto.EcdsaKey
	hub testutils.
	α   int

	Swarmer       swarm.Swarmer
	SwarmerServer swarm.Server
	MultiStorer   swarm.MultiAddressStorer

	Oracler       Oracler
	OraclerServer Server
	MidPointPrice MidpointPriceStorer
}

func newTester(index, α int) {

	key, err := crypto.RandomEcdsaKey()
	Expect(err).ShouldNot(HaveOccurred())


	client, store, err := testutils.NewMockSwarmClient(serverHub, &key, clientType)
}
