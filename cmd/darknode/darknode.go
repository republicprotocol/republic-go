package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/republicprotocol/go-dark-node"
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/go-identity"
)

//const PATH = "/home/ubuntu/"
const PATH = ""

var config *node.Config

func main() {

	// Parse command line arguments and fill the node.Config.
	if err := parseCommandLineFlags(); err != nil {
		log.Println(err)
		flag.Usage()
		return
	}

	// Create profiling logs for cpu and memory usage.
	if config.Dev {
		f, err := os.Create(fmt.Sprintf("%scpu.log", PATH))
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}

		go func() {
			time.Sleep(time.Hour)

			pprof.StopCPUProfile()

			f, err := os.Create(fmt.Sprintf("%smem.log", PATH))
			if err != nil {
				log.Fatal("could not create memory profile: ", err)
			}
			runtime.GC() // get up-to-date statistics
			if err := pprof.WriteHeapProfile(f); err != nil {
				log.Fatal("could not write memory profile: ", err)
			}
			f.Close()
		}()
	}

	// Create a new node.node.
	node, err := node.NewDarkNode(config)
	if err != nil {
		node.Configuration.Logger.Error(err)
	}

	// Start the dark node.
	options := do.CoBegin(func() do.Option {
		return do.Err(node.StartListening())
	}, func() do.Option {
		time.Sleep(time.Second)
		return do.Err(node.Start())
	})

	for _, option := range options {
		if option.Err != nil {
			node.Configuration.Logger.Error(err)
		}
	}
}

// Parse the config file path and read config from it.
func parseCommandLineFlags() error {
	confFilename := flag.String("config", fmt.Sprintf("%sdefault-config.json", PATH), "Path to the JSON configuration file")

	flag.Parse()

	conf, err := node.LoadConfig(*confFilename)
	if err != nil {
		conf, err = LoadDefaultConfig()
		if err != nil {
			return err
		}
		config = conf
		return nil
	}
	config = conf
	return nil
}

// LoadDefaultConfig loads a default config if no config is provided
func LoadDefaultConfig() (*node.Config, error) {
	address, keyPair, err := identity.NewAddress()
	if err != nil {
		return &node.Config{}, err
	}
	out, err := exec.Command("curl", "https://ipinfo.io/ip").Output()
	out = []byte(strings.Trim(string(out), "\n "))
	if err != nil {
		return &node.Config{}, err
	}
	multiAddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/%s/tcp/18514/republic/%s", out, address))
	if err != nil {
		return &node.Config{}, err
	}

	// 5 bootstraps nodes set up by the team.
	bootstrapNodes := []string{
		"/ip4/52.21.44.236/tcp/18514/republic/8MGg76n7RfC6tuw23PYf85VFyM8Zto",
		"/ip4/52.41.118.171/tcp/18514/republic/8MJ38m8Nzknh3gVj7QiMjuejmHBMSf",
		"/ip4/52.59.176.141/tcp/18514/republic/8MKDGUTgKtkymyKTH28xeMxiCnJ9xy",
		"/ip4/52.77.88.84/tcp/18514/republic/8MHarRJdvWd7SsTJE8vRVfj2jb5cWS",
		"/ip4/52.79.194.108/tcp/18514/republic/8MKZ8JwCU9m9affPWHZ9rxp2azXNnE",
	}
	filePlugin := node.NewFilePlugin("darknode.log")
	websocketPlugin := node.NewWebSocketPlugin(fmt.Sprintf("%s", out), "8080", "", "")

	config := &node.Config{
		Host:                    "0.0.0.0",
		Port:                    "18514",
		RepublicKeyPair:         keyPair,
		MultiAddress:            multiAddress,
		BootstrapMultiAddresses: make([]identity.MultiAddress, len(bootstrapNodes)),
		Logger:                  node.NewLogger(filePlugin, websocketPlugin),
		Dev:                     false,
	}
	for i, bootstrapNode := range bootstrapNodes {
		multi, err := identity.NewMultiAddressFromString(bootstrapNode)
		if err != nil {
			return &node.Config{}, err
		}
		config.BootstrapMultiAddresses[i] = multi
	}
	err = saveConfigFile(config)
	return config, err
}

func saveConfigFile(config *node.Config) error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	d1 := []byte(data)
	err = ioutil.WriteFile(fmt.Sprintf("%sdefault-config.json", PATH), d1, 0644)
	return err
}
