package grpc

import (
	"context"
	"fmt"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
)

type OracleClient struct {
	signature []byte
}

// NewOracleClient returns an implementation of the OracleClient interface.
func NewOracleClient(signature []byte) OracleClient {
	return OracleClient{
		signature: signature,
	}
}

// UpdateMidpoint is used to send updated midpoint information to a given
// multiaddress.
func (client *OracleClient) UpdateMidpoint(ctx context.Context, to identity.MultiAddress, tokens, price, nonce uint64) error {
	conn, err := Dial(ctx, to)
	if err != nil {
		logger.Network(logger.LevelError, fmt.Sprintf("cannot dial %v: %v", to, err))
		return fmt.Errorf("cannot dial %v: %v", to, err)
	}
	defer conn.Close()

	request := &UpdateMidpointRequest{
		Signature: client.signature,
		Tokens:    tokens,
		Price:     price,
		Nonce:     nonce,
	}
	if err := Backoff(ctx, func() error {
		_, err = NewOracleServiceClient(conn).UpdateMidpoint(ctx, request)
		return err
	}); err != nil {
		return err
	}

	return nil
}
