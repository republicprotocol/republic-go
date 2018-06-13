package adapter

import "github.com/republicprotocol/republic-go/status"

// Status defines a structure for JSON marshalling
type Status struct {
	Address         string `json:"address"`
	MultiAddress    string `json:"multiAddress"`
	EthereumAddress string `json:"ethereumAddress"`
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
	addr, err := adapter.Address()
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
	return Status{
		Address:         addr,
		MultiAddress:    multiAddr.String(),
		EthereumAddress: ethAddr,
	}, nil
}
