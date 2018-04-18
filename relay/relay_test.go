package relay_test

import (
//	"errors"
	// "fmt"
//	"sync"
	"time"

//	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

//	"github.com/ethereum/go-ethereum/accounts/abi/bind"
//	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/ethereum/contracts"
	// "github.com/republicprotocol/republic-go/contracts/dnr"
	"github.com/republicprotocol/republic-go/darknode"
	"github.com/republicprotocol/republic-go/identity"
	// "github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/order"
	. "github.com/republicprotocol/republic-go/relay"
	"github.com/republicprotocol/republic-go/stackint"
)

// var dnrOuterLock = new(sync.Mutex)
// var dnrInnerLock = new(sync.Mutex)

var epochDNR contracts.DarkNodeRegistry
// var nodes []*node.DarkNode

var Prime, _ = stackint.FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111")

// var _ = BeforeSuite(func() {
// 	var err error
// 	epochDNR, err = dnr.TestnetDNR(nil)
// 	if err != nil {
// 		panic(err)
// 	}
// 	By("generate nodes")
// 	nodes, err = generateNodes(5)
// 	Ω(err).ShouldNot(HaveOccurred())
// 	err = registerNodes(nodes)
// 	Ω(err).ShouldNot(HaveOccurred())

// 	By("start node service")
// 	startNodeServices(nodes)

// 	By("start node background workers")
// 	startNodeBackgroundWorkers(nodes)

// 	By("bootstrap nodes")
// 	bootstrapNodes(nodes)
// })

// var _ = Describe("Relay", func() {

// 	var err error
// 	epochDNR, err = dnr.TestnetDNR(nil)
// 	if err != nil {
// 		panic(err)
// 	}

// 	Context("when sending full orders", func() {

// 		It("should not return an error", func() {
// 			pools, trader := getPoolsAndTrader()

// 			sendOrder := getFullOrder()

// 			err = SendOrderToDarkOcean(sendOrder, &trader, pools, getBootstrapNodes())
// 			Ω(err).ShouldNot(HaveOccurred())
// 		})
// 	})

// 	Context("when sending fragmented orders that do not have sufficient fragments", func() {

// 		It("should return an error", func() {
// 			pools, trader := getPoolsAndTrader()

// 			sendOrder := getFragmentedOrder()

// 			err = SendOrderFragmentsToDarkOcean(sendOrder, &trader, pools, getBootstrapNodes())
// 			Ω(err).Should(HaveOccurred())
// 			Expect(err.Error()).To(ContainSubstring("number of fragments do not match pool size"))
// 		})
// 	})

// 	Context("when sending fragmented orders that have sufficient fragments for atleast one dark pool", func() {

// 		It("should not return an error", func() {
// 			pools, trader := getPoolsAndTrader()

// 			sendOrder, err := generateFragmentedOrderForDarkPool(pools[0])
// 			Ω(err).ShouldNot(HaveOccurred())

// 			fmt.Println("Before relay ...", sendOrder.DarkPools)
// 			err = SendOrderFragmentsToDarkOcean(sendOrder, &trader, pools, getBootstrapNodes())
// 			Ω(err).ShouldNot(HaveOccurred())
// 		})
// 	})

// 	Context("when canceling orders", func() {

// 		It("should not return an error", func() {
// 			pools, trader := getPoolsAndTrader()
			
// 			orderID := []byte("vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQyV6ryi1wDSM=")

// 			err = CancelOrder(orderID, &trader, pools, getBootstrapNodes())
// 			Ω(err).ShouldNot(HaveOccurred())
// 		})
// 	})
// })

// var _ = AfterSuite(func() {
// 	err := deregisterNodes(nodes)
// 	Ω(err).ShouldNot(HaveOccurred())
// 	stopNodes(nodes)
// })

// func generateNodes(numberOfNodes int) ([]*node.DarkNode, error) {
// 	// Generate nodes from the config files
// 	nodes := make([]*node.DarkNode, numberOfNodes)
// 	for i := 0; i < numberOfNodes; i++ {
// 		var err error
// 		var config *node.Config
// 		if i < 5 {
// 			config, err = node.LoadConfig(fmt.Sprintf("../test/config/bootstrap-node-%d.json", i+1))
// 		} else {
// 			config, err = node.LoadConfig(fmt.Sprintf("../test/config/node-%d.json", i-6))
// 		}
// 		if err != nil {
// 			return nil, err
// 		}
// 		auth := bind.NewKeyedTransactor(config.EthereumKey.PrivateKey)
// 		dnr, err := dnr.TestnetDNR(auth)
// 		if err != nil {
// 			return nil, err
// 		}
// 		node, err := node.NewDarkNode(*config, dnr)
// 		if err != nil {
// 			return nil, err
// 		}
// 		nodes[i] = node
// 	}
// 	return nodes, nil
// }

// func registerNodes(nodes []*node.DarkNode) error {
// 	dnrOuterLock.Lock()
// 	dnrInnerLock.Lock()
// 	defer dnrInnerLock.Unlock()
// 	for _, node := range nodes {
// 		isRegistered, err := node.DarkNodeRegistry.IsRegistered(nodes[0].NetworkOptions.MultiAddress.ID())
// 		if isRegistered {
// 			return errors.New("already registered")
// 		}

// 		bond := stackint.FromUint(10)
// 		err = node.DarkNodeRegistry.ApproveRen(&bond)
// 		if err != nil {
// 			return err
// 		}

// 		node.DarkNodeRegistry.SetGasLimit(300000)
// 		_, err = node.DarkNodeRegistry.Register(node.ID, []byte{}, &bond)
// 		node.DarkNodeRegistry.SetGasLimit(0)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	_, err := epochDNR.WaitForEpoch()
// 	return err
// }

// func deregisterNodes(nodes []*node.DarkNode) error {
// 	defer dnrOuterLock.Unlock()
// 	dnrInnerLock.Lock()
// 	defer dnrInnerLock.Unlock()
// 	for _, node := range nodes {
// 		node.DarkNodeRegistry.SetGasLimit(300000)
// 		_, err := node.DarkNodeRegistry.Deregister(node.ID)
// 		node.DarkNodeRegistry.SetGasLimit(0)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// 	epochDNR.SetGasLimit(300000)
// 	_, err := epochDNR.WaitForEpoch()
// 	epochDNR.SetGasLimit(0)
// 	if err != nil {
// 		panic(err)
// 	}
// 	for _, node := range nodes {
// 		node.DarkNodeRegistry.SetGasLimit(300000)
// 		_, err := node.DarkNodeRegistry.Refund(node.ID)
// 		node.DarkNodeRegistry.SetGasLimit(0)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// 	_, err = epochDNR.WaitForEpoch()
// 	return err
// }

// func startNodeServices(nodes []*node.DarkNode) {
// 	for i := range nodes {
// 		go func(i int) {
// 			defer GinkgoRecover()
// 			nodes[i].StartServices()
// 		}(i)
// 	}
// 	time.Sleep(time.Millisecond * time.Duration(10*len(nodes)))
// }

// func startNodeBackgroundWorkers(nodes []*node.DarkNode) {
// 	for i := range nodes {
// 		nodes[i].StartBackgroundWorkers()
// 	}
// 	time.Sleep(time.Millisecond * time.Duration(10*len(nodes)))
// }

// func bootstrapNodes(nodes []*node.DarkNode) {
// 	do.CoForAll(nodes, func(i int) {
// 		nodes[i].Bootstrap()
// 	})
// 	do.CoForAll(nodes, func(i int) {
// 		nodes[i].Bootstrap()
// 	})
// }

// func stopNodes(nodes []*node.DarkNode) {
// 	for i := range nodes {
// 		nodes[i].Stop()
// 	}
// }

// getPools return dark pools from a mock dnr
func getPools(dnr contracts.DarkNodeRegistry) darknode.Pools {
	// log, err := logger.NewLogger(logger.Options{})
	// if err != nil {
	// 	panic(fmt.Sprintf("cannot get logger: %v", err))
	// }

	ocean := darknode.NewOcean(dnr)
	// if err != nil {
	// 	panic(fmt.Sprintf("cannot get dark ocean: %v", err))
	// }
	return ocean.GetPools();
}

func getFullOrder() order.Order {
	fullOrder := order.Order{}

	defaultStackVal, _ := stackint.FromString("179761232312312")

	fullOrder.ID = []byte("vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQyV6ryi1wDSM=")
	fullOrder.Type = 2
	fullOrder.Parity = 1
	fullOrder.Expiry = time.Time{}
	fullOrder.FstCode = order.CurrencyCodeETH
	fullOrder.SndCode = order.CurrencyCodeBTC
	fullOrder.Price = defaultStackVal
	fullOrder.MaxVolume = defaultStackVal
	fullOrder.MinVolume = defaultStackVal
	fullOrder.Nonce = defaultStackVal
	return fullOrder
}

func getFragmentedOrder() OrderFragments {
	defaultStackVal, _ := stackint.FromString("179761232312312")

	fragmentedOrder := OrderFragments{}
	fragmentSet := map[string][]*order.Fragment{}
	fragments := []*order.Fragment{}

	var err error
	fragments, err = order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, defaultStackVal, defaultStackVal, defaultStackVal, defaultStackVal).Split(2, 1, &Prime)
	Ω(err).ShouldNot(HaveOccurred())

	fragmentSet["vrZhWU3VV9LRIM="] = fragments

	fragmentedOrder.ID = []byte("vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQyV6ryi1wDSM=")
	fragmentedOrder.Type = 2
	fragmentedOrder.Parity = 1
	fragmentedOrder.Expiry = time.Time{}
	fragmentedOrder.DarkPools = fragmentSet

	return fragmentedOrder
}

func generateFragmentedOrderForDarkPool(pool *darknode.Pool) (OrderFragments, error) {
	sendOrder := getFullOrder()
	fragments, err := sendOrder.Split(int64(pool.Size()), int64(pool.Size()*2/3), &Prime)
	if err != nil {
		return OrderFragments{}, err
	}
	fragmentSet := map[string][]*order.Fragment{}
	fragmentOrder := getFragmentedOrder()
	fragmentSet[GeneratePoolID(pool)] = fragments
	fragmentOrder.DarkPools = fragmentSet
	return fragmentOrder, nil
}

func getPoolsAndTrader() (darknode.Pools, identity.MultiAddress) {
	// trader, err := identity.NewMultiAddressFromString("/ip4/127.0.0.1/tcp/80/republic/8MGfbzAMS59Gb4cSjpm34soGNYsM2f")
	trader, err := identity.NewMultiAddressFromString("/ip4/0.0.0.0/tcp/3003/republic/8MJNCQhMrUCHuAk977igrdJk3tSzkT")
	Ω(err).ShouldNot(HaveOccurred())

	return getPools(epochDNR), trader
}

func getBootstrapNodes() []string {
	return []string{
		"/ip4/0.0.0.0/tcp/3003/republic/8MJNCQhMrUCHuAk977igrdJk3tSzkT",
		"/ip4/0.0.0.0/tcp/3000/republic/8MJxpBsezEGKPZBbhFE26HwDFxMtFu",
		"/ip4/0.0.0.0/tcp/3001/republic/8MGB2cj2HbQFepRVs43Ghct5yCRS9C",
		"/ip4/0.0.0.0/tcp/3002/republic/8MGVBvrQJji8ecEf3zmb8SXFCx1PaR",
		"/ip4/0.0.0.0/tcp/3004/republic/8MK6bq5m7UfE1mzRNunJTFH6zTbyss",
	}
}
