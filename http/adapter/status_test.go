package adapter_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/http/adapter"

	"github.com/republicprotocol/republic-go/http"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/status"
	"github.com/republicprotocol/republic-go/testutils"
)

var _ = Describe("Status adapter", func() {
	var prov status.Provider

	// populateProvider will populate the fields in the provider with valid
	// values.
	populateProvider := func(prov status.Provider) {
		prov.WriteDarknodeRegistryAddress("0x000000000000000000")
		prov.WriteEthereumAddress("0xeeeeeeeeeeeeeeeeee")
		prov.WriteEthereumNetwork("falcon")
		prov.WriteInfuraURL("https://kovan.infura.io")
		multiAddr, err := testutils.RandomMultiAddress()
		Expect(err).ShouldNot(HaveOccurred())
		prov.WriteMultiAddress(multiAddr)
		prov.WriteNetwork("falcon")
		prov.WritePublicKey([]byte{byte(103)})
		prov.WriteRewardVaultAddress("0x123456789012345678")
		prov.WriteTokens(map[string]string{"REN": "083", "DGX": "012", "ABC": "223"})
	}

	// assertStatus will assert that all the fields in the status match the
	// fields in the Reader object.
	assertStatus := func(status Status, reader status.Reader) {
		providerDarknodeRegistryAddress, err := reader.DarknodeRegistryAddress()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(status.DarknodeRegistryAddress).To(Equal(providerDarknodeRegistryAddress))

		providerEthereumAddress, err := reader.EthereumAddress()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(status.EthereumAddress).To(Equal(providerEthereumAddress))

		providerInfuraURL, err := reader.InfuraURL()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(status.InfuraURL).To(Equal(providerInfuraURL))

		providerMultiAddress, err := reader.MultiAddress()
		Expect(err).ShouldNot(HaveOccurred())
		providerMultiAddressStr := ""
		if !providerMultiAddress.IsNil() {
			providerMultiAddressStr = providerMultiAddress.String()
		}
		Expect(status.MultiAddress).To(Equal(providerMultiAddressStr))

		providerNetwork, err := reader.Network()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(status.Network).To(Equal(providerNetwork))

		providerPublicKey, err := reader.PublicKey()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(status.PublicKey).To(Equal(fmt.Sprintf("0x%x", providerPublicKey)))

		providerRewardVaultAddress, err := reader.RewardVaultAddress()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(status.RewardVaultAddress).To(Equal(providerRewardVaultAddress))

		providerTokens, err := reader.Tokens()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(status.Tokens).To(Equal(providerTokens))
	}

	// sendRequestAndAssertSuccess will send a GET http request to retrieve the
	// status of the server and assert that the response is as expected.
	sendRequestAndAssertSuccess := func(statusAdapter StatusAdapter) {
		body := bytes.NewBuffer([]byte{})
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "http://localhost/status", body)
		server := http.NewStatusServer(statusAdapter)

		server.ServeHTTP(w, r)

		Expect(w.Code).To(Equal(200))

		var response Status
		err := json.Unmarshal(w.Body.Bytes(), &response)
		Expect(err).ShouldNot(HaveOccurred())
		assertStatus(response, statusAdapter.Reader)
	}

	BeforeEach(func() {
		swarmer := testutils.NewMockSwarmer()
		prov = status.NewProvider(&swarmer)
		populateProvider(prov)
	})

	Context("when checking status of the provider", func() {
		It("should return the status as it is in the provider", func() {
			statusAdapter := NewStatusAdapter(prov)

			sendRequestAndAssertSuccess(statusAdapter)
		})

		Context("when multi-address is nil", func() {
			It("should not return an error code", func() {
				statusAdapter := NewStatusAdapter(prov)

				// Update provider's multi-address to be nil.
				prov.WriteMultiAddress(identity.MultiAddress{})

				// The response of the server's status must not result in a
				// panic or nil-pointer error.
				sendRequestAndAssertSuccess(statusAdapter)
			})
		})
	})
})
