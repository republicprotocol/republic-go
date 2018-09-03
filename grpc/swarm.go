package grpc

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/swarm"
	"github.com/republicprotocol/republic-go/testutils"
	"golang.org/x/net/context"
	"google.golang.org/grpc/peer"
)

// RedNodeBehaviour indicates the malicious behaviours the
// red node will exhibit during swarming.
type RedNodeBehaviour int

// Values for a RedNodeBehaviour
const (
	InvalidRequests RedNodeBehaviour = iota
	InvalidNonce
	InvalidSignature
	DropMultiAddresses
	DropSignatures
)

// RedNodeTypes contains an array of all possible red-node behaviours.
var RedNodeTypes = []RedNodeBehaviour{
	InvalidRequests,
	InvalidNonce,
	InvalidSignature,
	DropMultiAddresses,
	DropSignatures,
}

// String returns a human-readable representation of RedNodeTypes.
func (behaviours RedNodeBehaviour) String() string {
	switch behaviours {
	case InvalidRequests:
		return "invalid requests"
	case InvalidNonce:
		return "invalid nonce"
	case InvalidSignature:
		return "invalid multi-address signature"
	case DropMultiAddresses:
		return "drop multi-addresses"
	case DropSignatures:
		return "drop multi-address signatures"
	default:
		return "unexpected behaviour"
	}
}

// ErrRateLimitExceeded is returned when the same client sends more than one
// request to the server within a specified rate limit.
var ErrRateLimitExceeded = errors.New("cannot process request, rate limit exceeded")

type swarmClient struct {
	addr  identity.Address
	store swarm.MultiAddressStorer
}

// NewSwarmClient returns an implementation of the swarm.Client interface that
// uses gRPC and a recycled connection pool.
func NewSwarmClient(store swarm.MultiAddressStorer, addr identity.Address) swarm.Client {
	return &swarmClient{
		addr:  addr,
		store: store,
	}
}

// Ping implements the swarm.Client interface.
func (client *swarmClient) Ping(ctx context.Context, to identity.MultiAddress, multiAddr identity.MultiAddress) error {
	conn, err := Dial(ctx, to)
	if err != nil {
		logger.Network(logger.LevelError, fmt.Sprintf("cannot dial %v: %v", to, err))
		return fmt.Errorf("cannot dial %v: %v", to, err)
	}
	defer conn.Close()

	multiAddress := getTamperedMultiAddress(multiAddr)

	request := &PingRequest{
		MultiAddress: &multiAddress,
	}

	return Backoff(ctx, func() error {
		_, err = NewSwarmServiceClient(conn).Ping(ctx, request)
		return err
	})
}

func (client *swarmClient) Pong(ctx context.Context, to identity.MultiAddress) error {
	conn, err := Dial(ctx, to)
	if err != nil {
		logger.Network(logger.LevelError, fmt.Sprintf("cannot dial %v: %v", to, err))
		return fmt.Errorf("cannot dial %v: %v", to, err)
	}
	defer conn.Close()

	multiAddr, err := client.store.MultiAddress(client.addr)
	if err != nil {
		logger.Network(logger.LevelError, fmt.Sprintf("cannot get self details: %v", err))
		return fmt.Errorf("cannot get self details: %v", err)
	}

	multiAddress := getTamperedMultiAddress(multiAddr)

	request := &PongRequest{
		MultiAddress: &multiAddress,
	}

	return Backoff(ctx, func() error {
		_, err = NewSwarmServiceClient(conn).Pong(ctx, request)
		return err
	})
}

// Query implements the swarm.Client interface.
func (client *swarmClient) Query(ctx context.Context, to identity.MultiAddress, query identity.Address) (identity.MultiAddresses, error) {
	conn, err := Dial(ctx, to)
	if err != nil {
		logger.Network(logger.LevelError, fmt.Sprintf("cannot dial %v: %v", to, err))
		return identity.MultiAddresses{}, fmt.Errorf("cannot dial %v: %v", to, err)
	}
	defer conn.Close()

	request := &QueryRequest{
		Address: query.String(),
	}

	var response *QueryResponse
	if err := Backoff(ctx, func() error {
		response, err = NewSwarmServiceClient(conn).Query(ctx, request)
		return err
	}); err != nil {
		return identity.MultiAddresses{}, err
	}

	multiAddrs := identity.MultiAddresses{}
	for _, multiAddrMsg := range response.MultiAddresses {
		multiAddr, err := identity.NewMultiAddressFromString(multiAddrMsg.MultiAddress)
		if err != nil {
			logger.Network(logger.LevelWarn, fmt.Sprintf("cannot parse %v: %v", multiAddrMsg.MultiAddress, err))
			continue
		}
		multiAddr.Nonce = multiAddrMsg.MultiAddressNonce
		multiAddr.Signature = multiAddrMsg.Signature
		multiAddrs = append(multiAddrs, multiAddr)
	}
	return multiAddrs, nil
}

// MultiAddress implements the swarm.Client interface.
func (client *swarmClient) MultiAddress() identity.MultiAddress {
	multiAddr, err := client.store.MultiAddress(client.addr)
	if err != nil {
		logger.Network(logger.LevelError, fmt.Sprintf("cannot retrieve own multiaddress: %v", err))
		return identity.MultiAddress{}
	}
	return multiAddr
}

// SwarmService is a Service that implements the gRPC SwarmService defined in
// protobuf. It delegates responsibility for handling the Ping and Query RPCs
// to a swarm.Server.
type SwarmService struct {
	server swarm.Server

	rate         time.Duration
	rateLimitsMu *sync.Mutex
	rateLimits   map[string]time.Time
}

// NewSwarmService returns a SwarmService that uses the swarm.Server as a
// delegate.
func NewSwarmService(server swarm.Server, rate time.Duration) SwarmService {
	return SwarmService{
		server: server,

		rate:         rate,
		rateLimitsMu: new(sync.Mutex),
		rateLimits:   make(map[string]time.Time),
	}
}

// Register implements the Service interface.
func (service *SwarmService) Register(server *Server) {
	RegisterSwarmServiceServer(server.Server, service)
}

// Ping is an RPC used to notify a SwarmService about the existence of a
// client. In the PingRequest, the client sends a signed identity.MultiAddress
// and the SwarmService delegates the responsibility of handling this signed
// identity.MultiAddress to its swarm.Server. If its swarm.Server accepts the
// signed identity.MultiAddress of the client it will return its own signed
// identity.MultiAddress in a PingResponse.
func (service *SwarmService) Ping(ctx context.Context, request *PingRequest) (*PingResponse, error) {
	if err := service.isRateLimited(ctx); err != nil {
		return nil, err
	}

	from, err := identity.NewMultiAddressFromString(request.GetMultiAddress().GetMultiAddress())
	if err != nil {
		logger.Network(logger.LevelError, fmt.Sprintf("cannot unmarshal multiaddress: %v", err))
		return nil, fmt.Errorf("cannot unmarshal multiaddress: %v", err)
	}
	from.Signature = request.GetMultiAddress().GetSignature()
	from.Nonce = request.GetMultiAddress().GetMultiAddressNonce()

	err = service.server.Ping(ctx, from)
	if err != nil {
		logger.Network(logger.LevelInfo, fmt.Sprintf("cannot update store with: %v", err))
		return &PingResponse{}, fmt.Errorf("cannot update store: %v", err)
	}
	return &PingResponse{}, nil
}

// Pong is an RPC used to notify a SwarmService about the existence of a
// client. In the PingRequest, the client sends a signed identity.MultiAddress
// and the SwarmService delegates the responsibility of handling this signed
// identity.MultiAddress to its swarm.Server. If its swarm.Server accepts the
// signed identity.MultiAddress of the client it will return its own signed
// identity.MultiAddress in a PongResponse.
func (service *SwarmService) Pong(ctx context.Context, request *PongRequest) (*PongResponse, error) {
	if err := service.isRateLimited(ctx); err != nil {
		return nil, err
	}

	from, err := identity.NewMultiAddressFromString(request.GetMultiAddress().GetMultiAddress())
	if err != nil {
		logger.Network(logger.LevelError, fmt.Sprintf("cannot unmarshal multiaddress: %v", err))
		return nil, fmt.Errorf("cannot unmarshal multiaddress: %v", err)
	}

	from.Signature = request.GetMultiAddress().GetSignature()
	from.Nonce = request.GetMultiAddress().GetMultiAddressNonce()

	err = service.server.Pong(ctx, from)
	if err != nil {
		logger.Network(logger.LevelInfo, fmt.Sprintf("cannot update storer with %v: %v", request.GetMultiAddress(), err))
		return &PongResponse{}, fmt.Errorf("cannot update storer: %v", err)
	}
	return &PongResponse{}, nil
}

// Query is an RPC used to find identity.MultiAddresses. In the QueryRequest,
// the client sends an identity.Address and the SwarmService will stream
// identity.MultiAddresses to the client. The SwarmService delegates
// responsibility to its swarm.Server to return identity.MultiAddresses that
// are close to the queried identity.Address.
func (service *SwarmService) Query(ctx context.Context, request *QueryRequest) (*QueryResponse, error) {
	if err := service.isRateLimited(ctx); err != nil {
		return nil, err
	}

	query := identity.Address(request.GetAddress())
	multiAddrs, err := service.server.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	multiAddrMsgs := make([]*MultiAddress, len(multiAddrs))
	for i, multiAddr := range multiAddrs {
		multiAddress := getTamperedMultiAddress(multiAddr)
		multiAddrMsgs[i] = &multiAddress
	}

	return &QueryResponse{
		MultiAddresses: multiAddrMsgs,
	}, nil
}

func (service *SwarmService) isRateLimited(ctx context.Context) error {
	client, ok := peer.FromContext(ctx)
	if !ok {
		return fmt.Errorf("fail to get peer from ctx")
	}
	if client.Addr == net.Addr(nil) {
		return fmt.Errorf("fail to get peer address")
	}

	clientAddr, ok := client.Addr.(*net.TCPAddr)
	if !ok {
		return fmt.Errorf("fail to read peer TCP address")
	}
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

func getTamperedMultiAddress(multiAddr identity.MultiAddress) MultiAddress {
	redNodeType := RedNodeTypes[rand.Intn(len(RedNodeTypes))]

	rand.Seed(time.Now().UnixNano())
	multiAddress := MultiAddress{
		Signature:         multiAddr.Signature,
		MultiAddress:      multiAddr.String(),
		MultiAddressNonce: multiAddr.Nonce,
	}

	switch redNodeType {
	case InvalidRequests:
		multiAddress.Signature = tamperSignature(multiAddr)
		multiAddress.MultiAddressNonce = tamperNonce(multiAddr)
		multiAddress.MultiAddress = tamperMultiAddress(multiAddr)
	case InvalidNonce:
		multiAddress.MultiAddressNonce = tamperNonce(multiAddr)
	case InvalidSignature:
		multiAddress.Signature = tamperSignature(multiAddr)
	case DropMultiAddresses:
		multiAddress.MultiAddress = ""
	case DropSignatures:
		multiAddress.Signature = []byte{}
	default:
	}

	log.Printf("Red-node swarmer will exhibit behaviour: %v\n", redNodeType)
	log.Printf("Red-node tampered multi-address %v to look like %v", multiAddr, multiAddress)
	return multiAddress
}

func tamperSignature(multiAddr identity.MultiAddress) []byte {
	r := rand.Intn(100)
	if r < 50 {
		randBytes := testutils.Random64Bytes()
		return randBytes[:]
	}
	multiAddr.Signature[rand.Intn(64)] = byte(rand.Intn(100))
	return multiAddr.Signature
}

func tamperMultiAddress(multiAddr identity.MultiAddress) string {
	r := rand.Intn(100)
	if r < 75 {
		multiAddr, _ := testutils.RandomMultiAddress()
		return multiAddr.String()
	}
	return multiAddr.String()
}

func tamperNonce(multiAddr identity.MultiAddress) uint64 {
	r := rand.Intn(100)
	if r < 33 {
		return multiAddr.Nonce + uint64(r)
	}
	if r < 66 {
		return multiAddr.Nonce - uint64(r)
	}
	if r < 99 {
		return 0
	}
	return multiAddr.Nonce
}
