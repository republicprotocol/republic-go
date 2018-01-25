package identity

import "fmt"

// Errors returned by the package.
var (
	ErrFailToDecode       = fmt.Errorf("fail to decode the string")
	ErrWrongAddressLength = fmt.Errorf("wrong address length")
)
