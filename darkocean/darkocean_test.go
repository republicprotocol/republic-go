package darkocean_test

// . "github.com/onsi/ginkgo"
// . "github.com/onsi/gomega"
// . "github.com/republicprotocol/republic-go/darknode"

// var _ = Describe("Ocean", func() {

// 	Context("when watching for changes to the Ocean", func() {

// 		It("should signal changes once per epoch", func(done Done) {
// 			defer close(done)

// 			numberOfEpochs := 10

// 			conn, err := ganache.Connect("http://localhost:8545")
// 			Expect(err).ShouldNot(HaveOccurred())
// 			darknodeRegistry, err := contracts.NewDarkNodeRegistry(context.Background(), conn, ganache.GenesisTransactor(), &bind.CallOpts{})
// 			Expect(err).ShouldNot(HaveOccurred())
// 			darknodeRegistry.SetGasLimit(1000000)
// 			minimumEpochInterval, err := darknodeRegistry.MinimumEpochInterval()
// 			Expect(err).ShouldNot(HaveOccurred())
// 			minimumEpochIntervalInSeconds, err := minimumEpochInterval.ToUint()
// 			Expect(err).ShouldNot(HaveOccurred())
// 			minimumEpochIntervalDuration := time.Duration(minimumEpochIntervalInSeconds) * time.Second

// 			quit := make(chan struct{})

// 			// Start turning epochs in the background
// 			go func() {
// 				defer GinkgoRecover()

// 				t := time.NewTicker(minimumEpochIntervalDuration)
// 				defer t.Stop()

// 				for {
// 					select {
// 					case <-quit:
// 						return
// 					case <-t.C:
// 						_, err := darknodeRegistry.Epoch()
// 						Expect(err).ShouldNot(HaveOccurred())
// 					}
// 				}
// 			}()

// 			nodes, err := darknodeRegistry.GetAllNodes()
// 			Expect(err).Should(BeNil())
// 			epoch, err := darknodeRegistry.CurrentEpoch()
// 			Expect(err).Should(BeNil())
// 			// Start watching for updates to the Ocean
// 			ocean := NewDarkOcean(epoch.Blockhash, nodes)
// 			changes, errs := ocean.Watch(quit)
// 			for i := 0; i < numberOfEpochs; i++ {
// 				Eventually(changes, 2*minimumEpochIntervalDuration).Should(Receive())
// 			}

// 			close(quit)
// 			Expect(<-errs).ShouldNot(HaveOccurred())

// 		}, 600)
// 	})

// })
