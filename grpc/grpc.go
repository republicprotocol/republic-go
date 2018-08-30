package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

// Server re-exports the grpc.Server type.
type Server struct {
	*grpc.Server
}

func NewServer() *Server {
	return &Server{grpc.NewServer()}
}

// NewServer re-exports the grpc.NewServer function.
func NewServerwithLimiter(unaryLimiter, streamLimiter *RateLimiter) *Server {
	unaryInterceptor := grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		clientIP, err := addressFromContext(ctx)
		if err != nil {
			return nil, err
		}
		if unaryLimiter.Allow(clientIP) {
			return handler(ctx, req)
		}
		log.Println(clientIP, "hit the unary rate limit")

		return nil, errors.New("429: Too Many Requests")
	})

	streamInterceptor := grpc.StreamInterceptor(func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		clientIP, err := addressFromContext(stream.Context())
		if err != nil {
			return err
		}
		if streamLimiter.Allow(clientIP) {
			return handler(srv, stream)
		}
		log.Println(clientIP, "hit the stream rate limit")

		return errors.New("429: Too Many Requests")
	})

	return &Server{grpc.NewServer(unaryInterceptor, streamInterceptor)}
}

// Start the Server listening on a TCP connection at the given binding address.
func (server *Server) Start(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return server.Serve(lis)
}

// A Service can register to a Server. Registration must happen before the
// Server is started. The Service will be available when the Server is started.
type Service interface {
	Register(server *Server)
}

func addressFromContext(ctx context.Context) (string, error) {
	client, ok := peer.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("fail to get peer from ctx")
	}
	if client.Addr == net.Addr(nil) {
		return "", fmt.Errorf("fail to get peer address")
	}
	clientAddr, ok := client.Addr.(*net.TCPAddr)
	if !ok {
		return "", fmt.Errorf("fail to read peer TCP address")
	}

	return clientAddr.IP.String(), nil
}
