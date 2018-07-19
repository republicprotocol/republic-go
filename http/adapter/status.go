package adapter

import (
	"encoding/hex"

	"github.com/republicprotocol/republic-go/status"
)

// Status defines a structure for JSON marshalling
type Status struct {
	Network                 string            `json:"network"`
	MultiAddress            string            `json:"multiAddress"`
	EthereumNetwork         string            `json:"ethereumNetwork"`
	EthereumAddress         string            `json:"ethereumAddress"`
	DarknodeRegistryAddress string            `json:"darknodeRegistryAddress"`
	RewardVaultAddress      string            `json:"rewardVaultAddress"`
	PublicKey               string            `json:"publicKey"`
	InfuraURL               string            `json:"infura"`
	Tokens                  map[string]string `json:"tokens"`
	Peers                   int               `json:"peers"`
}

// StatusAdapter defines a struct which has status reading capability
type StatusAdapter struct {
	status.Reader
}

// NewStatusAdapter returns an adapter which contains reading statuses
func NewStatusAdapter(reader status.Reader) StatusAdapter {
	return StatusAdapter{
		Reader: reader,
	}
}

// Status returns a Status object with populated fields
func (adapter *StatusAdapter) Status() (Status, error) {
	network, err := adapter.Network()
	if err != nil {
		return Status{}, err
	}
	multiAddr, err := adapter.MultiAddress()
	if err != nil {
		return Status{}, err
	}
	ethAddr, err := adapter.EthereumAddress()
	if err != nil {
		return Status{}, err
	}
	ethNetwork, err := adapter.EthereumNetwork()
	if err != nil {
		return Status{}, err
	}
	darknodeRegistryAddr, err := adapter.DarknodeRegistryAddress()
	if err != nil {
		return Status{}, err
	}
	rewardVaultAddr, err := adapter.RewardVaultAddress()
	if err != nil {
		return Status{}, err
	}
	peers, err := adapter.Peers()
	if err != nil {
		return Status{}, err
	}
	infuraURL, err := adapter.InfuraURL()
	if err != nil {
		return Status{}, err
	}
	tokens, err := adapter.Tokens()
	if err != nil {
		return Status{}, err
	}
	pk, err := adapter.PublicKey()
	if err != nil {
		return Status{}, err
	}
	hexPk := "0x" + hex.EncodeToString(pk)
	return Status{
		Network:                 network,
		MultiAddress:            multiAddr.String(),
		EthereumNetwork:         ethNetwork,
		EthereumAddress:         ethAddr,
		DarknodeRegistryAddress: darknodeRegistryAddr,
		RewardVaultAddress:      rewardVaultAddr,
		PublicKey:               hexPk,
		InfuraURL:               infuraURL,
		Tokens:                  tokens,
		Peers:                   peers,
	}, nil
}
