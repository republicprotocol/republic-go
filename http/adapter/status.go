package adapter

type StatusAdapter struct {
}

func NewStatusAdapter() StatusAdapter {
	return StatusAdapter{}
}

func (adapter *StatusAdapter) GenerateResponse() string {
	return "Hello World!"
}
