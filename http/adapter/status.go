package adapter

import "github.com/republicprotocol/republic-go/status"

type Status struct {
	Address string `json:"address"`
}

type StatusAdapter struct {
	status.Provider
}

func NewStatusAdapter(provider status.Provider) StatusAdapter {
	return StatusAdapter{
		Provider: provider,
	}
}

// Status returns a Status object with populated fields
func (adapter *StatusAdapter) Status() (Status, error) {
	return Status{
		Address: "Hello World!",
	}, nil
}
