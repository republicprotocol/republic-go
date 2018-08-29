package grpc

import (
	"context"
	"fmt"
	"net"
	"sync"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

// Server re-exports the grpc.Server type.
type Server struct {
	*grpc.Server
}

// NewServer re-exports the grpc.NewServer function.
func NewServer() *Server {
	return &Server{grpc.NewServer()}
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

type rateLimiter struct {
	mu       *sync.Mutex
	limiters map[string]*rate.Limiter
}

func (limiter rateLimiter) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		client, ok := peer.FromContext(ctx)
		if !ok {
			return nil, fmt.Errorf("fail to get peer from ctx")
		}
		if client.Addr == net.Addr(nil) {
			return nil, fmt.Errorf("fail to get peer address")
		}

		clientAddr, ok := client.Addr.(*net.TCPAddr)
		if !ok {
			return nil, fmt.Errorf("fail to read peer TCP address")
		}
		clientIP := clientAddr.IP.String()
		limiter.mu.Lock()
		defer limiter.mu.Unlock()
		if l, ok := limiter.limiters[clientIP]; ok {
			if l.Allow() {
				return handler(ctx, req)
			}
		} else {
			l := rate.NewLimiter(5, 25)
			limiter.limiters[clientIP] = l
			if l.Allow() {
				return handler(ctx, req)
			}
		}

		return nil, nil
	}
}

func (limiter rateLimiter) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		client, ok := peer.FromContext(stream.Context())
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
		limiter.mu.Lock()
		defer limiter.mu.Unlock()
		if l, ok := limiter.limiters[clientIP]; ok {
			if l.Allow() {
				return handler(srv, stream)
			}
		} else {
			l := rate.NewLimiter(5, 25)
			limiter.limiters[clientIP] = l
			if l.Allow() {
				return handler(srv, stream)
			}
		}

		return nil
	}
}
