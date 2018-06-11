package grpc

import (
	"net"

	"google.golang.org/grpc"
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
