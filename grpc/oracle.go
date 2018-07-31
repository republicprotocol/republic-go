package grpc

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/oracle"
	"github.com/republicprotocol/republic-go/swarm"
	"google.golang.org/grpc/peer"
)

type oracleClient struct {
	addr  identity.Address
	store swarm.MultiAddressStorer
}

// NewOracleClient returns an implementation of the OracleClient interface.
func NewOracleClient(addr identity.Address, store swarm.MultiAddressStorer) oracle.Client {
	return &oracleClient{
		addr:  addr,
		store: store,
	}
}

// UpdateMidpoint is used to send updated midpoint information to a given
// multiaddress.
func (client *oracleClient) UpdateMidpoint(ctx context.Context, to identity.MultiAddress, midpointPrice oracle.MidpointPrice) error {
	conn, err := Dial(ctx, to)
	if err != nil {
		logger.Network(logger.LevelError, fmt.Sprintf("cannot dial %v: %v", to, err))
		return fmt.Errorf("cannot dial %v: %v", to, err)
	}
	defer conn.Close()

	request := &UpdateMidpointRequest{
		Signature: midpointPrice.Signature,
		Tokens:    midpointPrice.Tokens,
		Price:     midpointPrice.Price,
		Nonce:     midpointPrice.Nonce,
	}
	if err := Backoff(ctx, func() error {
		_, err = NewOracleServiceClient(conn).UpdateMidpoint(ctx, request)
		return err
	}); err != nil {
		return err
	}

	return nil
}

func (client *oracleClient) MultiAddress() identity.MultiAddress {
	multiAddr, err := client.store.MultiAddress(client.addr)
	if err != nil {
		logger.Network(logger.LevelError, fmt.Sprintf("cannot retrieve own multiaddress: %v", err))
		return identity.MultiAddress{}
	}
	return multiAddr
}

type OracleService struct {
	server oracle.Server

	rate         time.Duration
	rateLimitsMu *sync.Mutex
	rateLimits   map[string]time.Time
}

func NewOracleService(server oracle.Server, rate time.Duration) OracleService {
	return OracleService{
		server:       server,
		rate:         rate,
		rateLimitsMu: new(sync.Mutex),
		rateLimits:   make(map[string]time.Time),
	}
}

func (service *OracleService) Register(server *Server) {
	RegisterOracleServiceServer(server.Server, service)
}

func (service *OracleService) UpdateMidpoint(ctx context.Context, request *UpdateMidpointRequest) (*UpdateMidpointResponse, error) {
	if err := service.isRateLimited(ctx); err != nil {
		return nil, err
	}

	if request.Signature == nil || len(request.Signature) == 0 || request.Tokens == 0 || request.Price == 0 || request.Nonce == 0 {
		return nil, fmt.Errorf("invalid midpoint data request")
	}

	midpointPrice := oracle.MidpointPrice{
		Signature: request.Signature,
		Tokens:    request.Tokens,
		Price:     request.Price,
		Nonce:     request.Nonce,
	}

	return &UpdateMidpointResponse{}, service.server.UpdateMidpoint(ctx, midpointPrice)
}

func (service *OracleService) isRateLimited(ctx context.Context) error {
	client, ok := peer.FromContext(ctx)
	if !ok {
		return fmt.Errorf("failed to get peer from ctx")
	}
	if client.Addr == net.Addr(nil) {
		return fmt.Errorf("failed to get peer address")
	}

	clientAddr := client.Addr.(*net.TCPAddr)
	clientIP := clientAddr.IP.String()

	service.rateLimitsMu.Lock()
	defer service.rateLimitsMu.Unlock()

	if lastPing, ok := service.rateLimits[clientIP]; ok {
		if service.rate > time.Since(lastPing) {
			return ErrRateLimitExceeded
		}
	}

	service.rateLimits[clientIP] = time.Now()
	return nil
}
