package node_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-dark-node"
	"github.com/republicprotocol/go-identity"
)

var _ = Describe("Dark nodes", func() {
	const Boostrap_Nodes_Port = 3000
	const Test_Nodes_Port = 4000

	var bootstrapNodes, nodes []*node.DarkNode
	var err error

	startListening := func(nodes []*node.DarkNode) {
		for _, node := range nodes{
			go func() {
				defer GinkgoRecover()

				node.Start()
			}()
		}
	}

	stopListening := func(nodes []*node.DarkNode) {
		for _, node := range nodes{
			go func() {
				defer GinkgoRecover()

				node.Stop()
			}()
		}
	}

	for _, numberOfBootstrapNodes := range []int{4} {
		for _, numberOfNodes := range []int{8} {
			Context("nodes start up", func() {
				BeforeEach(func() {
					err = generateConfigs(numberOfBootstrapNodes, Boostrap_Nodes_Port , []*node.DarkNode{})
					Ω(err).ShouldNot(HaveOccurred())
					bootstrapNodes, err = generateNodes(numberOfBootstrapNodes)
					Ω(err).ShouldNot(HaveOccurred())

					err = generateConfigs(numberOfNodes, Test_Nodes_Port, bootstrapNodes)
					Ω(err).ShouldNot(HaveOccurred())
					nodes, err = generateNodes(numberOfNodes)
					Ω(err).ShouldNot(HaveOccurred())

					startListening(bootstrapNodes)
				})

				AfterEach(func() {
					stopListening(nodes)
					stopListening(bootstrapNodes)
				})

				It("should be able to run startup successfully", func() {
					startListening(nodes)
				})

				It("should register itself during startup", func() {

				})

				It("should receive order fragment and discover (mis)matches ", func() {

				})
			})
		}
	}
})

// Generate config files in the test_configs folder
func generateConfigs(numberOfNodes, port int, bootstrap []*node.DarkNode) error {
	var configs []node.Config
	bootstrapNodes := make ([]identity.MultiAddress, len(bootstrap))
	for i, j  := range bootstrap{
		bootstrapNodes[i] = j.Configuration.MultiAddress
	}
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
			RSAKeyPair:              keyPair,
			MultiAddress:            multi,
			BootstrapMultiAddresses: bootstrapNodes,
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
func generateNodes(numberOfNodes int) ([]*node.DarkNode, error) {
	nodes := make([]*node.DarkNode, numberOfNodes)

	for i := 0; i < numberOfNodes; i++ {
		config, err := node.LoadConfig(fmt.Sprintf("./test_configs/config-%d.json", i))
		if err != nil {
			return nil, err
		}
		node, err := node.NewDarkNode(config)
		if err != nil {
			return nil, err
		}
		nodes[i] = node
	}
	return nodes, nil
}
