package grpc

import (
	"context"
	"fmt"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
)

type OracleClient struct {
	signature []byte
	tokens    uint64
	price     uint64
	nonce     uint64
}

// NewOracleClient returns an implementation of the OracleClient interface.
func NewOracleClient(signature []byte, tokens, price, nonce uint64) OracleClient {
	return OracleClient{
		signature: signature,
		tokens:    tokens,
		price:     price,
		nonce:     nonce,
	}
}

// UpdateMidpoint is used to connect to a peer and propogate price information
// OracleClient to the rest of the network.
func (client *OracleClient) UpdateMidpoint(ctx context.Context, to identity.MultiAddress) error {
	conn, err := Dial(ctx, to)
	if err != nil {
		logger.Network(logger.LevelError, fmt.Sprintf("cannot dial %v: %v", to, err))
		return fmt.Errorf("cannot dial %v: %v", to, err)
	}
	defer conn.Close()

	request := &UpdateMidpointRequest{
		Signature: client.signature,
		Tokens:    client.tokens,
		Price:     client.price,
		Nonce:     client.nonce,
	}
	if err := Backoff(ctx, func() error {
		_, err = NewOracleServiceClient(conn).UpdateMidpoint(ctx, request)
		return err
	}); err != nil {
		return err
	}

	return nil
}
