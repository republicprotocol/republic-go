package main_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/dnr"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/ledger"
	"github.com/republicprotocol/republic-go/blockchain/test/ganache"
	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/darknode"
	"github.com/republicprotocol/republic-go/dht"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/grpc"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/leveldb"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/ome"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/smpc"
	"github.com/republicprotocol/republic-go/stackint"
	"github.com/republicprotocol/republic-go/stream"
	"github.com/republicprotocol/republic-go/swarm"
)

var _ = Describe("Darknode integration", func() {

	mu := new(sync.Mutex)

	n := 24
	bootstraps := 5

	var genesis ethereum.Conn
	var configs []darknode.Config

	var darknodes []darknode.Darknode
	var servers []*grpc.Server
	var swarmers []swarm.Swarmer
	var smpcers []smpc.Smpcer
	var omes []ome.Ome
	var stores []leveldb.Store

	BeforeEach(func() {
		var err error
		mu.Lock()

		genesis, err = ganache.StartAndConnect()
		Expect(err).ShouldNot(HaveOccurred())

		configs, err = newDarknodeConfigs(n, bootstraps)
		Expect(err).ShouldNot(HaveOccurred())
		darknodes, servers, swarmers, smpcers, omes, stores, err = newDarknodes(genesis, configs)
		Expect(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		defer mu.Unlock()

		for _, server := range servers {
			server.Stop()
		}

		for _, smpcer := range smpcers {
			err := smpcer.Shutdown()
			Expect(err).ShouldNot(HaveOccurred())
		}

		for _, store := range stores {
			err := store.Close()
			Expect(err).ShouldNot(HaveOccurred())
		}

		os.RemoveAll("./db.out/")
	})

	Context("when bootstrapping into a network", func() {

		It("should be able to query the super majority of nodes", func() {
			Expect(nil).To(BeNil())
		})

		It("should be responsive to RPCs", func() {
			Expect(nil).To(BeNil())
		})

	})

	Context("when the epoch changes", func() {

		It("should reach consensus on the configuration of the pods", func() {
			Expect(nil).To(BeNil())
		})

		It("should continue computations in previous epoch", func() {
			Expect(nil).To(BeNil())
		})

		It("should run computations in new epoch", func() {
			Expect(nil).To(BeNil())
		})

	})

	Context("when orders are opened", func() {

		FIt("should confirm matching orders", func(done Done) {
			defer close(done)

			By("booting darknodes")
			go func() {
				dispatch.CoForAll(configs, func(i int) {
					dispatch.CoBegin(func() {
						// Start the SMPC
						smpcers[i].Start()
					}, func() {
						defer GinkgoRecover()

						// Bootstrap into the network
						err := swarmers[i].Bootstrap(context.Background(), configs[i].BootstrapMultiAddresses)
						Expect(err).ShouldNot(HaveOccurred())
					}, func() {
						// Synchronizing the OME
						for {
							if err := omes[i].Sync(); err != nil {
								log.Printf("cannot sync ome: %v", err)
								continue
							}
							time.Sleep(time.Second * 2)
						}
					}, func() {
						// Synchronizing the Darknode
						for {
							if err := darknodes[i].Sync(); err != nil {
								log.Printf("cannot sync darknode: %v", err)
								continue
							}
							time.Sleep(time.Second * 2)
						}
					})
				})
			}()

			// Start gRPC server
			go func() {
				defer GinkgoRecover()

				dispatch.CoForAll(configs, func(i int) {
					log.Printf("listening on %v:%v...", configs[i].Host, configs[i].Port)
					lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", configs[i].Host, configs[i].Port))
					if err != nil {
						log.Fatalf("cannot listen on %v:%v: %v", configs[i].Host, configs[i].Port, err)
					}
					if err := servers[i].Serve(lis); err != nil {
						log.Fatalf("cannot serve on %v:%v: %v", configs[i].Host, configs[i].Port, err)
					}
				})
			}()

			err := openOrderPairs(genesis, configs)
			Expect(err).ShouldNot(HaveOccurred())

			<-time.NewTimer(time.Minute).C
		}, 60 /* 60 second timeout */)

		It("should not confirm mismatching orders", func() {
			Expect(nil).To(BeNil())
		})

	})

	Context("when orders are canceled", func() {

		It("should not confirm canceled orders", func() {
			Expect(nil).To(BeNil())
		})

	})

	Context("when orders are confirmed", func() {

		It("should not reconfirm orders", func() {
			Expect(nil).To(BeNil())
		})

	})

	Context("when faults occur under the threshold", func() {

		It("should not block the computation", func() {
			Expect(nil).To(BeNil())
		})

		It("should not corrupt the state of the faulty nodes", func() {
			Expect(nil).To(BeNil())
		})

	})

	Context("when faults occur above the threshold", func() {

		It("should block the computation", func() {
			Expect(nil).To(BeNil())
		})

		It("should unblock the computation after recovery", func() {
			Expect(nil).To(BeNil())
		})

		It("should not corrupt the state of the faulty nodes", func() {
			Expect(nil).To(BeNil())
		})
	})

})

func newDarknodeConfigs(n int, bootstraps int) ([]darknode.Config, error) {
	configs := []darknode.Config{}

	for i := 0; i < n; i++ {
		keystore, err := crypto.RandomKeystore()
		if err != nil {
			return configs, err
		}

		addr := identity.Address(keystore.Address())
		configs = append(configs, darknode.Config{
			Keystore:                keystore,
			Host:                    "0.0.0.0",
			Port:                    fmt.Sprintf("%d", 18514+i),
			Address:                 addr,
			BootstrapMultiAddresses: identity.MultiAddresses{},
			Logs: logger.Options{
				Plugins: []logger.PluginOptions{
					{
						File: &logger.FilePluginOptions{
							Path: fmt.Sprintf("%v.out", addr),
						},
					},
				},
			},
			Ethereum: ethereum.Config{
				Network:                 ethereum.NetworkGanache,
				URI:                     "http://localhost:8545",
				RepublicTokenAddress:    ethereum.RepublicTokenAddressOnGanache.String(),
				DarknodeRegistryAddress: ethereum.DarknodeRegistryAddressOnGanache.String(),
				RenLedgerAddress:        ethereum.RenLedgerAddressOnGanache.String(),
			},
		})
	}

	for i := 0; i < n; i++ {
		for j := 0; j < bootstraps; j++ {
			if i == j {
				continue
			}
			bootstrapMultiAddr, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/%v/tcp/%v/republic/%v", configs[j].Host, configs[j].Port, configs[j].Address))
			if err != nil {
				return configs, err
			}
			configs[i].BootstrapMultiAddresses = append(configs[i].BootstrapMultiAddresses, bootstrapMultiAddr)
		}
	}

	return configs, nil

}

func newDarknodes(genesis ethereum.Conn, configs []darknode.Config) ([]darknode.Darknode, []*grpc.Server, []swarm.Swarmer, []smpc.Smpcer, []ome.Ome, []leveldb.Store, error) {

	darknodes := []darknode.Darknode{}
	servers := []*grpc.Server{}
	swarmers := []swarm.Swarmer{}
	smpcers := []smpc.Smpcer{}
	omes := []ome.Ome{}
	stores := []leveldb.Store{}

	var wg sync.WaitGroup

	for i, config := range configs {
		addr := config.Address
		multiAddr, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/%v/tcp/%v/republic/%v", config.Host, config.Port, addr))
		if err != nil {
			return nil, nil, nil, nil, nil, nil, err
		}

		// New Ethereum bindings
		auth, darkPool, darkPoolAccounts, darkPoolFees, renLedger, err := newEthereumBindings(config.Keystore, config.Ethereum)
		if err != nil {
			return nil, nil, nil, nil, nil, nil, err
		}

		// New crypter for signing and verification
		crypter := darknode.NewCrypter(config.Keystore, darkPool, len(configs)/2, time.Second*2)

		// New database for persistent storage
		store, err := leveldb.NewStore(fmt.Sprintf("./db.out/%v", addr))
		if err != nil {
			return nil, nil, nil, nil, nil, nil, err
		}
		stores = append(stores, store)

		// New DHT
		dht := dht.NewDHT(addr, len(configs))

		// New gRPC components
		servers = append(servers, grpc.NewServer())
		connPool := grpc.NewConnPool(len(configs) / 2)

		statusService := grpc.NewStatusService(&dht)
		statusService.Register(servers[i])

		swarmClient := grpc.NewSwarmClient(multiAddr, &connPool)
		swarmService := grpc.NewSwarmService(swarm.NewServer(swarmClient, &dht))
		swarmService.Register(servers[i])
		swarmers = append(swarmers, swarm.NewSwarmer(swarmClient, &dht))

		orderbook := orderbook.NewOrderbook(config.Keystore.RsaKey, orderbook.NewSyncer(renLedger, 32), &stores[i])
		orderbookService := grpc.NewOrderbookService(orderbook)
		orderbookService.Register(servers[i])

		streamClient := grpc.NewStreamClient(&crypter, addr, &connPool)
		streamService := grpc.NewStreamService(&crypter, addr)
		streamService.Register(servers[i])
		streamer := stream.NewStreamRecycler(stream.NewStreamer(addr, streamClient, &streamService))

		// New secure multi-party computer
		smpcers = append(smpcers, smpc.NewSmpcer(swarmers[i], streamer, 1))

		// New OME
		omes = append(omes, ome.NewOme(ome.NewRanker(1, 0), ome.NewComputer(orderbook, smpcers[i]), orderbook, smpcers[i]))

		// New Darknode
		darknodes = append(darknodes, darknode.NewDarknode(config.Address, darkPool, darkPoolAccounts, darkPoolFees))
		darknodes[i].SetEpochListener(omes[i])

		// Distribute ETH to the Darknode
		ganache.DistributeEth(genesis, auth.From)

		wg.Add(1)
		go func(config darknode.Config) {
			defer wg.Done()

			// Register the Darknode
			pubKeyBytes, err := crypto.BytesFromRsaPublicKey(&config.Keystore.RsaKey.PublicKey)
			if err != nil {
				panic(fmt.Sprintf("cannot get bytes from rsa public key: %v", err))
			}
			zero := stackint.Zero()
			tx, err := darkPool.(*dnr.DarknodeRegistry).Register(addr.ID(), pubKeyBytes, &zero)
			if err != nil {
				panic(fmt.Sprintf("cannot register darknode: %v", err))
			}
			if _, err := genesis.PatchedWaitMined(context.Background(), tx); err != nil {
				panic(fmt.Sprintf("cannot wait for register: %v", err))
			}

			// Call the epoch and ignore errors
			time.Sleep(time.Second)
			_, _ = darkPool.(*dnr.DarknodeRegistry).TriggerEpoch()
		}(config)
	}
	wg.Wait()

	return darknodes, servers, swarmers, smpcers, omes, stores, nil

}

func newEthereumBindings(keystore crypto.Keystore, config ethereum.Config) (*bind.TransactOpts, cal.Darkpool, cal.DarkpoolAccounts, cal.DarkpoolFees, cal.RenLedger, error) {
	conn, err := ethereum.Connect(config)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("cannot connect to ethereum: %v", err)
	}

	auth := bind.NewKeyedTransactor(keystore.EcdsaKey.PrivateKey)
	auth.GasPrice = big.NewInt(1000000000)
	darkpool, err := dnr.NewDarknodeRegistry(context.Background(), conn, auth, &bind.CallOpts{})
	if err != nil {
		fmt.Println(fmt.Errorf("cannot bind to darkpool: %v", err))
		return auth, nil, nil, nil, nil, err
	}

	renLedger, err := ledger.NewRenLedgerContract(context.Background(), conn, auth, &bind.CallOpts{})
	if err != nil {
		fmt.Println(fmt.Errorf("cannot bind to ren ledger: %v", err))
		return auth, nil, nil, nil, nil, err
	}

	return auth, &darkpool, nil, nil, &renLedger, nil

}

func openOrderPairs(genesis ethereum.Conn, configs []darknode.Config) error {

	keystore, err := crypto.RandomKeystore()
	if err != nil {
		return err
	}

	multiAddr, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/0.0.0.0/tcp/3000/republic/%v", keystore.Address()))
	if err != nil {
		return err
	}

	dht := dht.NewDHT(multiAddr.Address(), len(configs))
	connPool := grpc.NewConnPool(len(configs) / 2)
	orderbookClient := grpc.NewOrderbookClient(&connPool)
	swarmClient := grpc.NewSwarmClient(multiAddr, &connPool)
	swarmer := swarm.NewSwarmer(swarmClient, &dht)
	swarmer.Bootstrap(context.Background(), configs[0].BootstrapMultiAddresses)

	auth, darkPool, _, _, renLedger, err := newEthereumBindings(keystore, configs[0].Ethereum)
	if err != nil {
		return fmt.Errorf("cannot get ethereum binding: %v", err)
	}

	ganache.DistributeEth(genesis, auth.From)

	numberOfOrderPairs := 1
	for n := 0; n < numberOfOrderPairs; n++ {
		one := order.CoExp{
			Co:  200,
			Exp: 26,
		}

		buyOrder := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.TokensETHREN, one, one, one, rand.Int63())
		sellOrder := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.TokensETHREN, one, one, one, rand.Int63())

		orders := []order.Order{buyOrder, sellOrder}

		for _, ord := range orders {
			signatureData := crypto.Keccak256([]byte("Republic Protocol: open: "), ord.ID[:])
			signatureData = crypto.Keccak256([]byte("\x19Ethereum Signed Message:\n32"), signatureData)
			signature, err := keystore.Sign(signatureData)
			if err != nil {
				return fmt.Errorf("cannot sign: %v", err)
			}

			signature65 := [65]byte{}
			copy(signature65[:], signature)
			if err := renLedger.OpenBuyOrder(signature65, ord.ID); err != nil {
				return fmt.Errorf("cannot open order on ren ledger: %v", err)
			}

			pods, err := darkPool.Pods()
			if err != nil {
				return fmt.Errorf("cannot get pods: %v", err)
			}

			for _, pod := range pods {
				n := int64(len(pod.Darknodes))
				k := int64(2 * (len(pod.Darknodes) + 1) / 3)
				hash := base64.StdEncoding.EncodeToString(pod.Hash[:])
				ordFragments, err := ord.Split(n, k)
				if err != nil {
					return fmt.Errorf("cannot split order to %v: %v", hash, err)
				}

				for i, ordFragment := range ordFragments {
					pubKey, err := darkPool.PublicKey(pod.Darknodes[i])
					if err != nil {
						return fmt.Errorf("cannot get public key: %v", err)
					}

					encryptedFragment, err := ordFragment.Encrypt(pubKey)
					if err != nil {
						return fmt.Errorf("cannot encrypt: %v", err)
					}

					darknodeMultiAddr, err := swarmer.Query(context.Background(), pod.Darknodes[i], -1)
					if err != nil {
						return fmt.Errorf("cannot query: %v", err)
					}

					if err := orderbookClient.OpenOrder(context.Background(), darknodeMultiAddr, encryptedFragment); err != nil {
						return fmt.Errorf("cannot open order fragment: %v", err)
					}
				}
			}
		}
	}
	return nil
}
