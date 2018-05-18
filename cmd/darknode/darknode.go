package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/republicprotocol/republic-go/dispatch"

	"github.com/republicprotocol/republic-go/leveldb"
	"github.com/republicprotocol/republic-go/orderbook"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/dnr"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/ledger"
	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/darknode"
	"github.com/republicprotocol/republic-go/dht"
	"github.com/republicprotocol/republic-go/grpc"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/swarm"
)

func main() {
	defer log.Print("shutdown!")

	// Parse command-line arguments
	configParam := flag.String("config", path.Join(os.Getenv("HOME"), ".darknode/config.json"), "JSON configuration file")
	dataParam := flag.String("data", path.Join(os.Getenv("HOME"), ".darknode/data"), "Data directory")
	flag.Parse()

	// Load configuration file
	conf, err := darknode.NewConfigFromJSONFile(*configParam)
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	// Get IP-address
	ipAddr, err := getIPAddress()
	if err != nil {
		log.Fatalf("cannot get ip-address: %v", err)
	}

	// Get multi-address
	multiAddr, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/%s/tcp/18514/republic/%s", ipAddr, conf.Address))
	if err != nil {
		log.Fatalf("cannot get multiaddress: %v", err)
	}
	log.Printf("address %v", multiAddr)

	// Get ethereum bindings
	auth, _, _, _, renLedger, err := getEthereumBindings(conf.Keystore, conf.Ethereum)
	if err != nil {
		log.Fatalf("cannot get ethereum bindings: %v", err)
	}
	log.Printf("ethereum %v", auth.From.Hex())

	// Build services
	store, err := leveldb.NewStore(*dataParam)
	if err != nil {
		log.Fatalf("cannot create leveldb: %v", err)
	}
	server := grpc.NewServer()
	dht := dht.NewDHT(conf.Address, 32)
	connPool := grpc.NewConnPool(128)

	newStatus(&dht, server)
	newOrderbooker(&store, renLedger, server)
	swarmer := newSwarmer(multiAddr, &dht, &connPool, server)

	go func() {
		time.Sleep(time.Second)
		dispatch.CoBegin(func() {
			// Bootstrap into the network
			if err := swarmer.Bootstrap(context.Background(), conf.BootstrapMultiAddresses); err != nil {
				log.Printf("error during bootstrap: %v", err)
			}
			log.Printf("peers %v", len(dht.MultiAddresses()))
			for _, multiAddr := range dht.MultiAddresses() {
				log.Printf("  %v", multiAddr)
			}
		})
	}()

	// Start gRPC
	log.Printf("listening on %v:%v...", conf.Host, conf.Port)
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", conf.Host, conf.Port))
	if err != nil {
		log.Fatalf("cannot listen on %v:%v: %v", conf.Host, conf.Port, err)
	}
	if err := server.Serve(lis); err != nil {
		log.Fatalf("cannot serve on %v:%v: %v", conf.Host, conf.Port, err)
	}
}

func newStatus(dht *dht.DHT, server *grpc.Server) {
	service := grpc.NewStatusService(dht)
	service.Register(server)
}

func newSwarmer(multiAddr identity.MultiAddress, dht *dht.DHT, connPool *grpc.ConnPool, server *grpc.Server) swarm.Swarmer {
	client := grpc.NewSwarmClient(multiAddr, connPool)
	service := grpc.NewSwarmService(swarm.NewServer(client, dht))
	service.Register(server)
	return swarm.NewSwarmer(client, dht)
}

func newOrderbooker(store *leveldb.Store, renLedger cal.RenLedger, server *grpc.Server) orderbook.Orderbooker {
	orderbooker := orderbook.NewOrderbook(store, orderbook.NewSyncer(renLedger, 32))
	service := grpc.NewOrderbookService(orderbooker)
	service.Register(server)
	return orderbooker
}

func getIPAddress() (string, error) {

	out, err := exec.Command("curl", "https://ipinfo.io/ip").Output()
	out = []byte(strings.Trim(string(out), "\n "))
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func getEthereumBindings(keystore crypto.Keystore, conf ethereum.Config) (*bind.TransactOpts, cal.Darkpool, cal.DarkpoolAccounts, cal.DarkpoolFees, cal.RenLedger, error) {
	conn, err := ethereum.Connect(conf)
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
