package adapter

import "github.com/republicprotocol/republic-go/status"

// Status defines a structure for JSON marshalling
type Status struct {
	Address string `json:"address"`
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
	return Status{
		Address: "Hello World!",
	}, nil
}
