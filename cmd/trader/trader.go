package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/jbenet/go-base58"
	"github.com/republicprotocol/republic-go/contracts/connection"
	"github.com/republicprotocol/republic-go/contracts/dnr"
	"github.com/republicprotocol/republic-go/dark"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/network/rpc"
	"github.com/republicprotocol/republic-go/order"
)

var Prime, _ = big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)

const reset = "\x1b[0m"
const red = "\x1b[31;1m"

type OrderBook struct {
	LastUpdateId int             `json:"lastUpdateId"`
	Bids         [][]interface{} `json:"bids"`
	Asks         [][]interface{} `json:"asks"`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("expected command : open, cancel or binance")
	}

	openCommand := flag.NewFlagSet("open", flag.ExitOnError)
	openOrderFile := openCommand.String("file", "", "path of order file")
	openKeyFile := openCommand.String("keystore", "", "path of key file")
	openPassphrase := openCommand.String("passphrase", "", "passphrase")

	cancelCommand := flag.NewFlagSet("cancel", flag.ExitOnError)
	cancelOrderFile := cancelCommand.String("file", "", "path of order file")
	cancelKeyFile := cancelCommand.String("keystore", "", "path of key file")
	cancelPassphrase := cancelCommand.String("passphrase", "", "passphrase")

	binanceCommand := flag.NewFlagSet("binance", flag.ExitOnError)
	binanceOrderFile := binanceCommand.String("file", "", "path of order file")
	binanceKeyFile := binanceCommand.String("keystore", "", "path of key file")
	binancePassphrase := binanceCommand.String("passphrase", "", "passphrase")
	binanceNumberOfOrders := binanceCommand.Int("count", 10, "number of orders")

	if len(os.Args) == 1 {
		fmt.Println("usage: trader <command> [<args>]")
		return
	}

	switch os.Args[1] {
	case "open":
		openCommand.Parse(os.Args[2:])
	case "cancel":
		cancelCommand.Parse(os.Args[2:])
	case "binance":
		binanceCommand.Parse(os.Args[2:])
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}

	if openCommand.Parsed() {
		openOrder(*openOrderFile, *openKeyFile, *openPassphrase)
		return
	}
	if cancelCommand.Parsed() {
		cancelOrder(*cancelOrderFile, *cancelKeyFile, *cancelPassphrase)
		return
	}
	if binanceCommand.Parsed() {
		binanceOrder(*binanceOrderFile, *binanceKeyFile, *binancePassphrase, *binanceNumberOfOrders)
		return
	}

}

// openOrder will send order to all nodes in the dark ocean
func openOrder(openOrderFile, openKeyFile, openPassphrase string) {

	// Get trader address and dark pools
	traderMultiAddress, pools, err := getTraderAddressAndDarkPools(openKeyFile, openPassphrase)
	if err != nil {
		log.Fatalf("cannot get trader multi address and dark pools: %v", err)
	}

	// Get orders from JSON file
	orders, err := order.NewOrdersFromJSONFile(openOrderFile)
	if err != nil {
		log.Fatalf("cannot read orders from file: %v", err)
	}

	sendOrders(orders, pools, *traderMultiAddress)
}

func cancelOrder(closeOrderFile, closeKeyFile, closePassphrase string) {

	// Get trader address and dark pools
	traderMultiAddress, pools, err := getTraderAddressAndDarkPools(closeKeyFile, closePassphrase)
	if err != nil {
		log.Fatalf("cannot get trader multi address and dark pools: %v", err)
	}

	// TODO : uncomment this code when cancelOrders is functioning
	//
	// Get orders from JSON file
	// orders, err := order.NewOrdersFromJSONFile(closeOrderFile)
	// if err != nil {
	// 	log.Fatalf("cannot read orders from file: %v", err)
	// }

	// TODO: cancel only orders that are retrieved from file
	cancelTraderOrder(pools, *traderMultiAddress)
}

func binanceOrder(binanceOrderFile, binanceKeyFile, binancePassphrase string, binanceNumberOfOrders int) {

	// Get orders details from Binance
	resp, err := http.Get(fmt.Sprintf("https://api.binance.com/api/v1/depth?symbol=ETHBTC&limit=%v", binanceNumberOfOrders))
	if err != nil {
		log.Fatalf("cannot read data from binance: %v", err)
	}
	defer resp.Body.Close()

	orderBook := new(OrderBook)
	if err := json.NewDecoder(resp.Body).Decode(orderBook); err != nil {
		log.Fatalf("cannot parse orders from binance: %v", err)
	}

	writeOrders := make([]*order.Order, len(orderBook.Bids)+len(orderBook.Asks))

	// Generate orders from the Binance data
	sellOrders := make([]*order.Order, len(orderBook.Asks))
	for i, j := range orderBook.Asks {

		price, amount, err := parsePriceAndAmount(j)
		if err != nil {
			log.Fatalf("cannot parse price and amount: %v", err)
		}

		order := order.NewOrder(order.TypeLimit, order.ParitySell, time.Time{},
			order.CurrencyCodeETH, order.CurrencyCodeBTC, big.NewInt(int64(price)), big.NewInt(int64(amount)),
			big.NewInt(int64(amount)), big.NewInt(1))
		sellOrders[i] = order
	}

	buyOrders := make([]*order.Order, len(orderBook.Bids))
	for i, j := range orderBook.Asks { //change asks/bids

		price, amount, err := parsePriceAndAmount(j)
		if err != nil {
			log.Fatalf("cannot parse price and amount: %v", err)
		}

		order := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Time{},
			order.CurrencyCodeETH, order.CurrencyCodeBTC, big.NewInt(int64(price)), big.NewInt(int64(amount)),
			big.NewInt(int64(amount)), big.NewInt(1))
		buyOrders[i] = order
	}

	j := 0
	for i := 0; i < len(buyOrders); i++ {
		writeOrders[j] = buyOrders[i]
		j++
		writeOrders[j] = sellOrders[i]
		j++
	}

	if err := order.WriteOrdersToJSONFile(binanceOrderFile, writeOrders); err != nil {
		log.Fatalf("cannot write binance orders to file: %v", err)
	}
}

func sendOrders(orders []order.Order, pools dark.Pools, traderMultiAddress identity.MultiAddress) {
	for _, ord := range orders {
		// Buy or sell ?
		if ord.Parity == order.ParityBuy {
			log.Println("sending buy order : ", base58.Encode(ord.ID))
		} else {
			log.Println("sending sell order : ", base58.Encode(ord.ID))
		}

		var wg sync.WaitGroup
		wg.Add(len(pools))
		// For every dark pool
		for i := range pools {
			go func(darkPool *dark.Pool) {
				defer wg.Done()
				// Split order into (number of nodes in each pool) * 2/3 fragments
				shares, err := ord.Split(int64(darkPool.Size()), int64(darkPool.Size()*2/3), Prime)
				if err != nil {
					log.Println("cannot split orders: ", err)
					return
				}
				sendSharesToDarkPool(darkPool, traderMultiAddress, shares)
			}(pools[i])
		}
		wg.Wait()
	}
}

func cancelTraderOrder(pools dark.Pools, traderAddress identity.MultiAddress) {
	// For every Dark Pool
	for i := range pools {

		// Cancel orders for all nodes in the pool
		pools[i].ForAll(func(n *dark.Node) {

			// Get multiaddress
			multiaddress, err := getMultiAddress(n.ID.Address(), traderAddress)
			if err != nil {
				log.Fatalf("cannot read multi-address: %v", err)
			}

			// Create a client
			client, err := rpc.NewClient(multiaddress, traderAddress)
			if err != nil {
				log.Fatalf("cannot connect to client: %v", err)
			}

			// Close order
			err = client.CancelOrder(&rpc.OrderSignature{})
			if err != nil {
				log.Println(err)
				log.Printf("%sCoudln't cancel order to %v%s\n", red, base58.Encode(n.ID), reset)
				return
			}
		})
	}
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
	clientDetails, err := connection.FromURI("https://ropsten.infura.io/", "ropsten")
	if err != nil {
		return nil, fmt.Errorf("cannot fetch dark node registry: %v", err)
	}

	auth := bind.NewKeyedTransactor(key.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch dark node registry: %v", err)
	}

	// Gas Price
	auth.GasPrice = big.NewInt(6000000000)

	registrar, err := dnr.NewEthereumDarkNodeRegistrar(context.Background(), &clientDetails, auth, &bind.CallOpts{})

	// Get the Dark Ocean
	ocean, err := dark.NewOcean(logs, 5, registrar)
	if err != nil {
		return nil, fmt.Errorf("cannot read dark ocean: %v", err)
	}

	// return the dark pools
	return ocean.GetPools(), nil
}

func parsePriceAndAmount(j []interface{}) (float64, float64, error) {
	price, err := strconv.ParseFloat(j[0].(string), 10)
	if err != nil {
		return 0, 0, fmt.Errorf("cannot parse price into a big int: %v", err)
	}

	amount, err := strconv.ParseFloat(j[1].(string), 10)
	if err != nil {
		return 0, 0, fmt.Errorf("cannot parse amount into a big int: %v", err)
	}
	return amount * 1000000000000, price * 1000000000000, nil
}
