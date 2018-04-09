package relay

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os/exec"
	"strings"
	"sync"

	"github.com/republicprotocol/republic-go/order"
)

// Relay consists of configuration values (?)
type Relay struct {
	keyFile    string
	passPhrase string
}

// SendOrderToDarkOcean will fragment and send orders to the dark ocean
func SendOrderToDarkOcean(order order.Order) {
	fmt.Println(order.Type)

	// Get trader address and dark pools
	// TODO: Change to keyFile and passPhrase
	traderMultiAddress, pools, err := getTraderAddressAndDarkPools("../cmd/trader/test/keystore.json", "divya")
	if err != nil {
		log.Fatalf("cannot get trader multi address and dark pools: %v", err)
	}

	sendOrder(order, pools, *traderMultiAddress)
}

// SendOrderFragmentsToDarkOcean will send order fragments to the dark ocean
func SendOrderFragmentsToDarkOcean(order order.Order) {

}

// CancelOrder will cancel orders that aren't confirmed or settled in the dark ocean
func CancelOrder(order order.ID) {
	fmt.Println(order.String())
}

func getTraderAddressAndDarkPools(openKeyFile, openPassphrase string) (*identity.MultiAddress, dark.Pools, error) {
	// Get key and traderMultiAddress
	key, traderMultiAddress, err := getKeyAndAddress(openKeyFile, openPassphrase)
	if err != nil {
		return nil, nil, err
	}

	// Get dark pools
	pools, err := getDarkPools(key)
	if err != nil {
		return nil, nil, err
	}

	return traderMultiAddress, pools, nil

}

// Returns key and multi address
func getKeyAndAddress(filename, passphrase string) (*keystore.Key, *identity.MultiAddress, error) {

	// Get our IP address
	ipInfoOut, err := exec.Command("curl", "https://ipinfo.io/ip").Output()
	if err != nil {
		return nil, nil, err
	}
	ipAddress := strings.Trim(string(ipInfoOut), "\n ")

	// Read data from keystore file and generate the key
	encryptedKey, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot read keystore file: %v", err)
	}

	key, err := keystore.DecryptKey(encryptedKey, passphrase)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot decrypt key with provided passphrase: %v", err)
	}

	id, err := identity.NewKeyPairFromPrivateKey(key.PrivateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot generate id from key %v", err)
	}

	traderMultiAddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/%s/tcp/80/republic/%s", ipAddress, id.Address().String()))
	if err != nil {
		return nil, nil, fmt.Errorf("cannot obtain trader multi address %v", err)
	}
	log.Println("Trader Address: ", id.Address().String())

	return key, &traderMultiAddress, nil
}

func getDarkPools(key *keystore.Key) (dark.Pools, error) {
	// Logger
	logs, err := getLogger()
	if err != nil {
		return nil, fmt.Errorf("cannot get logger: %v", err)
	}

	// Get Dark Node Registry
	clientDetails, err := connection.FromURI("https://ropsten.infura.io/", connection.ChainRopsten)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch dark node registry: %v", err)
	}

	auth := bind.NewKeyedTransactor(key.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch dark node registry: %v", err)
	}

	// Gas Price
	auth.GasPrice = big.NewInt(6000000000)

	registrar, err := dnr.NewDarkNodeRegistry(context.Background(), &clientDetails, auth, &bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("cannot fetch dark node registry: %v", err)
	}

	// print(registrar.client)
	// Get the Dark Ocean
	ocean, err := dark.NewOcean(logs, 5, registrar)
	if err != nil {
		return nil, fmt.Errorf("cannot read dark ocean: %v", err)
	}

	// return the dark pools
	return ocean.GetPools(), nil
}

func sendOrder(openOrder order.Order, pools dark.Pools, traderMultiAddress identity.MultiAddress) {
	// Buy or sell ?
	if openOrder.Parity == order.ParityBuy {
		log.Println("sending buy order : ", base58.Encode(openOrder.ID))
	} else {
		log.Println("sending sell order : ", base58.Encode(openOrder.ID))
	}

	var wg sync.WaitGroup
	wg.Add(len(pools))
	// For every dark pool
	for i := range pools {
		go func(darkPool *dark.Pool) {
			defer wg.Done()
			// Split order into (number of nodes in each pool) * 2/3 fragments
			shares, err := openOrder.Split(int64(darkPool.Size()), int64(darkPool.Size()*2/3), &prime)
			if err != nil {
				log.Println("cannot split orders: ", err)
				return
			}
			sendSharesToDarkPool(darkPool, traderMultiAddress, shares)
		}(pools[i])
	}
	wg.Wait()
}

func getLogger() (*logger.Logger, error) {
	return logger.NewLogger(logger.Options{
		Plugins: []logger.PluginOptions{
			logger.PluginOptions{
				File: &logger.FilePluginOptions{
					Path: "stdout",
				},
			},
		}})
}


// Send the shares across all nodes within the Dark Pool
func sendSharesToDarkPool(pool *dark.Pool, multi identity.MultiAddress, shares []*order.Fragment) {
	j := 0
	pool.ForAll(func(n *dark.Node) {

		// Get multiaddress
		multiaddress, err := getMultiAddress(n.ID.Address(), multi)
		if err != nil {
			log.Fatalf("cannot read multi-address: %v", err)
		}

		// Create a client
		client, err := rpc.NewClient(multiaddress, multi)
		if err != nil {
			log.Fatalf("cannot connect to client: %v", err)
		}

		// Send fragment to node
		err = client.OpenOrder(&rpc.OrderSignature{}, rpc.SerializeOrderFragment(shares[j]))
		if err != nil {
			log.Println(err)
			log.Printf("%sCoudln't send order fragment to %v%s\n", red, base58.Encode(n.ID), reset)
			return
		}
		j++
	})
}

// Function to obtain multiaddress of a node by sending requests to bootstrap nodes
func getMultiAddress(address identity.Address, traderMultiAddress identity.MultiAddress) (identity.MultiAddress, error) {
	BootstrapMultiAddresses := []string{
		"/ip4/52.77.88.84/tcp/18514/republic/8MGzXN7M1ucxvtumVjQ7Ybb7xQ8TUw",
		"/ip4/52.79.194.108/tcp/18514/republic/8MGBUdoFFd8VsfAG5bQSAptyjKuutE",
		"/ip4/52.59.176.141/tcp/18514/republic/8MHmrykz65HimBPYaVgm8bTSpRUoXA",
		"/ip4/52.21.44.236/tcp/18514/republic/8MKFT9CDQQru1hYqnaojXqCQU2Mmuk",
		"/ip4/52.41.118.171/tcp/18514/republic/8MGb8k337pp2GSh6yG8iv2GK6FbNHN",
	}

	serializedTarget := rpc.SerializeAddress(address)
	for _, peer := range BootstrapMultiAddresses {

		bootStrapMultiAddress, err := identity.NewMultiAddressFromString(peer)
		if err != nil {
			return traderMultiAddress, err
		}

		client, err := rpc.NewClient(bootStrapMultiAddress, traderMultiAddress)
		if err != nil {
			return traderMultiAddress, err
		}

		candidates, err := client.QueryPeersDeep(serializedTarget)
		if err != nil {
			return traderMultiAddress, err
		}

		for candidate := range candidates {
			deserializedCandidate, err := rpc.DeserializeMultiAddress(candidate)
			if err != nil {
				return traderMultiAddress, err
			}
			if address == deserializedCandidate.Address() {
				fmt.Println("Found the target : ", address)
				return deserializedCandidate, nil
			}
		}
	}
	return traderMultiAddress, nil
}
