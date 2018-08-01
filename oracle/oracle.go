package oracle

import (
	"context"
	"fmt"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/swarm"
)

type Client interface {
	UpdateMidpoint(ctx context.Context, to identity.MultiAddress, midpointPrice MidpointPrice) error
	MultiAddress() identity.MultiAddress
}

type Oracler interface {
	UpdateMidpoint(ctx context.Context, midpointPrice MidpointPrice) error
	MultiAddress() identity.MultiAddress
}

type oracler struct {
	client          Client
	key             *crypto.EcdsaKey
	multiAddrStorer swarm.MultiAddressStorer
	α               int
}

// NewOracler will return an object that implements the Oracler interface.
func NewOracler(client Client, key *crypto.EcdsaKey, multiAddrStorer swarm.MultiAddressStorer, α int) Oracler {
	return &oracler{
		client:          client,
		key:             key,
		multiAddrStorer: multiAddrStorer,
		α:               α,
	}
}

func (oracler *oracler) UpdateMidpoint(ctx context.Context, midpointPrice MidpointPrice) error {
	randomMultiAddrs, err := swarm.RandomMultiAddrs(oracler.multiAddrStorer, oracler.MultiAddress().Address(), oracler.α)
	if err != nil {
		return err
	}

	var errs = make([]error, len(randomMultiAddrs))
	dispatch.CoForAll(randomMultiAddrs, func(i int) {
		errs[i] = oracler.client.UpdateMidpoint(ctx, randomMultiAddrs[i], midpointPrice)
		if errs[i] != nil {
			logger.Error(fmt.Sprintf("cannot send midpoint price to %v: %v", randomMultiAddrs[i].Address(), err))
		}
	})

	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}

func (oracler *oracler) MultiAddress() identity.MultiAddress {
	return oracler.client.MultiAddress()
}

type Server interface {
	UpdateMidpoint(ctx context.Context, midpointPrice MidpointPrice) error
}

type server struct {
	oracler             Oracler
	oracleAddr          identity.Address
	multiAddrStorer     swarm.MultiAddressStorer
	midpointPriceStorer MidpointPriceStorer
	α                   int
}

func NewServer(oracler Oracler, oracleAddr identity.Address, multiAddrStorer swarm.MultiAddressStorer, midpointPriceStorer MidpointPriceStorer, α int) Server {
	return &server{
		oracler:             oracler,
		oracleAddr:          oracleAddr,
		multiAddrStorer:     multiAddrStorer,
		midpointPriceStorer: midpointPriceStorer,
		α:                   α,
	}
}

func (server *server) UpdateMidpoint(ctx context.Context, midpointPrice MidpointPrice) error {
	verifier := crypto.NewEcdsaVerifier(server.oracleAddr.String())
	if err := verifier.Verify(midpointPrice.Hash(), midpointPrice.Signature); err != nil {
		return fmt.Errorf("failed to verify midpoint price signature: %v", err)
	}
	oldPrice, err := server.midpointPriceStorer.MidpointPrice(midpointPrice.Tokens)
	if err != nil {
		return err
	}

	if oldPrice.Equals(MidpointPrice{}) || midpointPrice.Nonce > oldPrice.Nonce {
		err := server.midpointPriceStorer.PutMidpointPrice(midpointPrice)
		if err != nil {
			return fmt.Errorf("cannot store midpoint price: %v", err)
		}
		return server.oracler.UpdateMidpoint(ctx, midpointPrice)
	}

	return nil
}

func (server *server) MultiAddress() identity.MultiAddress {
	return identity.MultiAddress{}
}
