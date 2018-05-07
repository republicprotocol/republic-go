package relay_test

/* import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/relay"

	"github.com/republicprotocol/republic-go/identity"
)

var _ = Describe("Relay", func() {
	var testnet TestnetEnv
	var api API

	BeforeEach(func() {
		var err error
		testnet, err = NewTestnet(10, darknodeTestnetEnv.DarknodeRegistry)
		Ω(err).ShouldNot(HaveOccurred())
		ip4, err := testnet.Relays[0].MultiAddress.ValueForProtocol(identity.IP4Code)
		Ω(err).ShouldNot(HaveOccurred())
		api = NewAPI("http://"+ip4+":4000", Insecure())
	})

	Context("when opening valid orders", func() {

		It("should not return an error", func() {
			ord := Order{
				Type:      1,
				Parity:    1,
				Expiry:    time.Now().AddDate(0, 0, 2),
				FstCode:   1,
				SndCode:   2,
				Price:     100,
				MaxVolume: 100,
				MinVolume: 100,
			}
			orderID, err := api.OpenOrder(ord)
			Ω(err).ShouldNot(HaveOccurred())
			order, err := api.GetOrder(orderID.String())
			fmt.Println(order)
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	Context("when canceling orders", func() {

		It("should not return an error", func() {

		})
	})
})
*/
