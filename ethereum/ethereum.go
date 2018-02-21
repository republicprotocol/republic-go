package ethereum

import (
	"errors"
	// "github.com/ethereum/go-ethereum/common"
	// "github.com/ethereum/go-ethereum/core/types"
)

// BytesTo32Bytes ...
func BytesTo32Bytes(bytes []byte) ([32]byte, error) {
	var bytes32 [32]byte
	if len(bytes) != 32 {
		return bytes32, errors.New("Expected 32 bytes")
	}
	for i := 0; i < 32; i++ {
		bytes32[i] = bytes[i]
	}

	return bytes32, nil
}
