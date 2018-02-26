package node_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-dark-node"
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/go-identity"
)

const Configs_Path = "./test_configs/"

var _ = Describe("Dark nodes", func() {

	start := func(nodes []*node.DarkNode){
		do.ForAll(nodes, func(i int) {
			nodes[i].Start()
		})
	}

	stop := func(nodes []*node.DarkNode){
		do.ForAll(nodes, func(i int) {
			nodes[i].Stop()
		})
	}

	Context("startup", func() {

		It("should read from the config file and create a new dark node.", func() {
			config, err := node.LoadConfig("test_configs/test_config.json")
			Ω(err).ShouldNot(HaveOccurred())
			node, err := node.NewDarkNode(config)
			Ω(err).ShouldNot(HaveOccurred())
			go func() {
				defer GinkgoRecover()
				node.Start()
			}()
		})

		It("should register itself on startup", func() {

		})
	})

	Context("dark pool specification", func() {

		It("should get dark pool assignments from the registrar", func() {

		})

		It("Ping all dark nodes in the relevant dark pool", func() {

		})
	})

	Context("Order computation", func() {
		const Boostrap_Nodes_Port = 3000
		const Test_Nodes_Port = 4000


		BeforeEach(func() {
			err := generateConfigs(10 , Boostrap_Nodes_Port)
			Ω(err).ShouldNot(HaveOccurred())

			err  = generateConfigs(10 , Test_Nodes_Port)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should receive order fragment and discover (mis)matches ", func() {
		})
	})
})

// Generate config files in the test_configs folder
func generateConfigs(numberOfNodes , port int) error {
	var configs []node.Config
	for i := 0; i < numberOfNodes; i++ {
		address, keyPair, err := identity.NewAddress()
		if err != nil {
			return err
		}
		multi, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/republic/%s", port+i, address.String()))
		if err != nil {
			return err
		}

		config := node.Config{
			Host:                    "127.0.0.1",
			Port:                    fmt.Sprintf("%d", port+i),
			RepublicKeyPair:         keyPair,
			MultiAddress:            multi,
			BootstrapMultiAddresses: identity.MultiAddresses{},
		}
		configs = append(configs, config)
	}

	for i := 0; i < numberOfNodes; i++ {
		data, err := json.Marshal(configs[i])
		if err != nil {
			return err
		}
		d1 := []byte(data)
		err = ioutil.WriteFile(fmt.Sprintf("./test_configs/config-%d.json", i), d1, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

// Generate nodes from the config files
func generateNodes(numberOfNodes int )([]*node.DarkNode, error ) {
	nodes := make ([]*node.DarkNode, numberOfNodes)

	for i:=0; i < numberOfNodes; i++{
		config, err  := node.LoadConfig(fmt.Sprintf("./test_configs/config-%d.json", i))
		if err != nil {
			return nil, err
		}
		node,err  := node.NewDarkNode(config)
		if err != nil {
			return nil, err
		}
		nodes[i] = node
	}
	return nodes, nil
}

