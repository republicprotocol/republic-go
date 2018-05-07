package relay

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/dnr"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/darkocean"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/rpc/client"
	"github.com/republicprotocol/republic-go/rpc/dht"
	"github.com/republicprotocol/republic-go/rpc/relayer"
	"github.com/republicprotocol/republic-go/rpc/smpcer"
	"github.com/republicprotocol/republic-go/rpc/swarmer"
	"google.golang.org/grpc"
)

// TestnetEnv stores all relays of the testnet in an exported Relays field.
type TestnetEnv struct {
	Relays Relays
}

// NewTestnet will create a testnet that sets up new Relays.
func NewTestnet(numberOfRelays int, darknodeRegistry dnr.DarknodeRegistry, port int, bootstrapAddresses identity.MultiAddresses) (TestnetEnv, error) {
	relays, err := NewRelays(numberOfRelays, darknodeRegistry, port, bootstrapAddresses)
	if err != nil {
		return TestnetEnv{}, fmt.Errorf("cannot create new relays: %v", err)
	}
	return TestnetEnv{
		Relays: relays,
	}, nil
}

// Run all Relays. This will start the gRPC services and bootstrap the
// Relays. Calls to TestnetEnv.Run are blocking, until a call to
// TestnetEnv.Teardown is made. Errors returned by the Relays while running
// are ignored.
// FIXME: Store the errors in a buffer that can be inspected after the test.
func (env *TestnetEnv) Run(port int) {
	dispatch.CoForAll(env.Relays, func(i int) {
		// Create gRPC server and TCP listener
		server := grpc.NewServer()
		listener, err := net.Listen("tcp", fmt.Sprintf("%v:%v", "127.0.0.1", fmt.Sprintf("%d", port+len(env.Relays)+i)))
		if err != nil {
			log.Fatal(err)
		}

		relay := env.Relays[i]
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		if err := relay.ListenAndServe("127.0.0.1", fmt.Sprintf("%d", port+i)); err != nil {
			log.Fatalf("error serving http: %v", err)
		}

		relay.Register(server)
		if err := server.Serve(listener); err != nil {
			log.Fatalf("error serving grpc: %v", err)
		}

		relay.Sync(context.Background(), 100)

		for err := range relay.swarmerClient.Bootstrap(ctx, relay.Config.BootstrapMultiAddresses, -1) {
			if strings.Contains(err.Error(), darkocean.ErrInvalidRegistration.Error()) {
				// try again
				log.Printf("error while bootstrapping the relay: %v trying again ..", err)
				time.Sleep(1 * time.Second)
				for errs := range relay.swarmerClient.Bootstrap(ctx, relay.Config.BootstrapMultiAddresses, -1) {
					log.Printf("error while bootstrapping the relay: %v", errs)
				}
			}
			log.Printf("error while bootstrapping the relay: %v", err)
		}

		// done := make(chan struct{})
		// entries := relay.orderbook.Listen(done)
		// defer close(done)

	})
}

// NewRelays configured for a local test environment.
func NewRelays(numberOfRelays int, darknodeRegistry dnr.DarknodeRegistry, port int, bootstrapAddresses identity.MultiAddresses) (Relays, error) {
	var err error

	relays := make(Relays, numberOfRelays)
	configs := make([]Config, numberOfRelays)
	for i := 0; i < numberOfRelays; i++ {
		configs[i], err = NewLocalConfig("127.0.0.1", fmt.Sprintf("%d", port+i), bootstrapAddresses)
		if err != nil {
			return nil, err
		}
	}

	for i := 0; i < numberOfRelays; i++ {
		book := orderbook.NewOrderbook()
		crypter := darkocean.NewCrypter(configs[i].Keystore, darknodeRegistry, 128, time.Minute)
		dht := dht.NewDHT(configs[i].MultiAddress.Address(), 100)
		connPool := client.NewConnPool(100)
		relayerClient := relayer.NewClient(&crypter, &dht, &connPool)
		smpcerClient := smpcer.NewClient(&crypter, configs[i].MultiAddress, &connPool)
		swarmerClient := swarmer.NewClient(&crypter, configs[i].MultiAddress, &dht, &connPool)
		relays[i] = NewRelay(configs[i], darknodeRegistry, &book, &relayerClient, &smpcerClient, &swarmerClient)
	}

	return relays, nil
}

// NewLocalConfig will return newly generated multiaddress and config that are
// constructed using the host and port that are passed as arguments to the
// method.
func NewLocalConfig(host, port string, bootstrapAddresses identity.MultiAddresses) (Config, error) {
	keystore, err := crypto.RandomKeystore()
	if err != nil {
		return Config{}, err
	}
	auth := bind.NewKeyedTransactor(keystore.EcdsaKey.PrivateKey)
	multiAddr, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/%s/tcp/%s/republic/%s", host, port, keystore.Address()))
	if err != nil {
		return Config{}, err
	}
	return Config{
		Token:                   "token",
		EthereumAddress:         auth.From.String(),
		Keystore:                keystore,
		MultiAddress:            multiAddr,
		BootstrapMultiAddresses: bootstrapAddresses,
		Ethereum: ethereum.Config{
			Network:                 ethereum.NetworkGanache,
			URI:                     "http://localhost:8545",
			RepublicTokenAddress:    ethereum.RepublicTokenAddressOnGanache.String(),
			DarknodeRegistryAddress: ethereum.DarknodeRegistryAddressOnGanache.String(),
			HyperdriveAddress:       ethereum.HyperdriveAddressOnRopsten.String(),
			ArcAddress:              ethereum.ArcAddressOnRopsten.String(),
		},
	}, nil
}
