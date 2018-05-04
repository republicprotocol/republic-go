package rpc_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/rpc"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/rpc/status"
	"github.com/republicprotocol/republic-go/shamir"
	"github.com/republicprotocol/republic-go/stackint"
	"google.golang.org/grpc"
)

var _ = Describe("rpc", func() {

	Context("NewRPC method", func() {

		It("should return an RPC object", func() {
			rpc, err := createNewRPC()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rpc).ShouldNot(BeNil())
		})
	})

	Context("OpenOrder method", func() {

		It("should not return an error", func() {
			rpc, err := createNewRPC()
			Expect(err).ShouldNot(HaveOccurred())

			fragment, err := createNewOrderFragment()
			Expect(err).ShouldNot(HaveOccurred())

			Expect(rpc.OpenOrder([]byte{}, fragment)).Should(BeNil())
		})

		It("should call onOpenOrder delegate", func() {
			rpc, err := createNewRPC()
			Expect(err).ShouldNot(HaveOccurred())

			fragment, err := createNewOrderFragment()
			Expect(err).ShouldNot(HaveOccurred())

			flag := false

			delegate := func([]byte, order.Fragment) error {
				flag = true
				return nil
			}

			rpc.OnOpenOrder(delegate)

			Expect(flag).To(BeFalse())
			Expect(rpc.OpenOrder([]byte{}, fragment)).Should(BeNil())
			Expect(flag).To(BeTrue())
		})
	})

	Context("CancelOrder method", func() {

		It("should not return an error", func() {
			rpc, err := createNewRPC()
			Expect(err).ShouldNot(HaveOccurred())

			testOrder := createNewOrder()
			Expect(rpc.CancelOrder([]byte{}, testOrder.ID)).Should(BeNil())
		})

		It("should call onCancelOrder delegate", func() {
			rpc, err := createNewRPC()
			Expect(err).ShouldNot(HaveOccurred())

			flag := false
			delegate := func([]byte, order.ID) error {
				flag = true
				return nil
			}
			rpc.OnCancelOrder(delegate)

			testOrder := createNewOrder()

			Expect(flag).To(BeFalse())
			Expect(rpc.CancelOrder([]byte{}, testOrder.ID)).Should(BeNil())
			Expect(flag).To(BeTrue())
		})
	})

	Context("OnOpenOrder method", func() {

		It("should set the rpc's onOpenOrder with delegate", func() {
			rpc, err := createNewRPC()
			Expect(err).ShouldNot(HaveOccurred())

			delegate := func([]byte, order.Fragment) error { return nil }

			oldRpc := rpc
			rpc.OnOpenOrder(delegate)
			Expect(rpc).ShouldNot(Equal(oldRpc))
		})
	})

	Context("OnCancelOrder method", func() {

		It("should set the rpc's OnCancelOrder with delegate", func() {
			rpc, err := createNewRPC()
			Expect(err).ShouldNot(HaveOccurred())

			delegate := func([]byte, order.ID) error { return nil }

			oldRpc := rpc
			rpc.OnCancelOrder(delegate)
			Expect(rpc).ShouldNot(Equal(oldRpc))
		})
	})

	Context("RelayerClient method", func() {

		It("should return rpc's RelayerClient", func() {
			rpc, err := createNewRPC()
			Expect(err).ShouldNot(HaveOccurred())

			Expect(rpc.RelayerClient()).ShouldNot(BeNil())

		})
	})

	Context("Relayer method", func() {

		It("should return rpc's Relayer", func() {
			rpc, err := createNewRPC()
			Expect(err).ShouldNot(HaveOccurred())

			Expect(rpc.Relayer()).ShouldNot(BeNil())

		})
	})

	Context("SmpcerClient method", func() {

		It("should return rpc's SmpcerClient", func() {
			rpc, err := createNewRPC()
			Expect(err).ShouldNot(HaveOccurred())

			Expect(rpc.SmpcerClient()).ShouldNot(BeNil())

		})
	})

	Context("Smpcer method", func() {

		It("should return rpc's Smpcer", func() {
			rpc, err := createNewRPC()
			Expect(err).ShouldNot(HaveOccurred())

			Expect(rpc.Smpcer()).ShouldNot(BeNil())

		})
	})

	Context("SwarmerClient method", func() {

		It("should return rpc's SwarmerClient", func() {
			rpc, err := createNewRPC()
			Expect(err).ShouldNot(HaveOccurred())

			Expect(rpc.SwarmerClient()).ShouldNot(BeNil())

		})
	})

	Context("Swarmer method", func() {

		It("should return rpc's Swarmer", func() {
			rpc, err := createNewRPC()
			Expect(err).ShouldNot(HaveOccurred())

			Expect(rpc.Swarmer()).ShouldNot(BeNil())

		})
	})

	Context("Status service", func() {
		It("should be able to return status of the node", func() {
			rpc, err := createNewRPC()
			Expect(err).ShouldNot(HaveOccurred())

			// Create rpc server
			server := grpc.NewServer()
			listener, err := net.Listen("tcp", "0.0.0.0:8080")
			Expect(err).ShouldNot(HaveOccurred())

			rpc.Relayer().Register(server)
			rpc.Smpcer().Register(server)
			rpc.Swarmer().Register(server)
			rpc.RegisterStatus(server)

			go func() {
				log.Println("listening at 0.0.0.0:8080...")
				if err = server.Serve(listener); err != nil {
					return
				}
			}()

			// Create rpc client
			conn, err := grpc.Dial("0.0.0.0:8080", grpc.WithInsecure())
			Expect(err).ShouldNot(HaveOccurred())
			statusClient := status.NewStatusClient(conn)
			status, err := statusClient.Status(context.Background(), &status.StatusRequest{})
			Expect(err).ShouldNot(HaveOccurred())

			Expect(status.Peers).Should(Equal(int64(0)))
			Expect(status.Bootstrapped).Should(BeFalse())
			Expect(len(status.Address)).Should(Equal(30))

			time.Sleep(1 * time.Second)
			server.Stop()
		})
	})
})

func createNewRPC() (RPC, error) {
	crypter := crypto.NewWeakCrypter()
	keystore, err := crypto.RandomKeystore()
	if err != nil {
		return RPC{}, err
	}
	multiaddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/80/republic/%s", keystore.Address()))
	if err != nil {
		return RPC{}, err
	}
	orderbook := orderbook.NewOrderbook(5)
	rpc := NewRPC(&crypter, multiaddress, &orderbook)
	return *rpc, nil
}

// Create a test order
func createNewOrder() order.Order {
	price := 1000000000000

	testOrder := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour),
		order.CurrencyCodeETH, order.CurrencyCodeBTC, stackint.FromUint(uint(price)), stackint.FromUint(uint(price)),
		stackint.FromUint(uint(price)), stackint.FromUint(1))
	return *testOrder
}

// Create a test order fragment
func createNewOrderFragment() (order.Fragment, error) {

	testOrder := createNewOrder()
	prime, err := stackint.FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111")
	if err != nil {
		return order.Fragment{}, err
	}

	share := shamir.Share{
		Key:   int64(1),
		Value: prime,
	}

	fragment := order.NewFragment(
		testOrder.ID,
		testOrder.Type,
		testOrder.Parity,
		share,
		share,
		share,
		share,
		share,
	)

	return *fragment, nil
}
