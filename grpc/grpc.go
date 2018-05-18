package grpc

import (
	"google.golang.org/grpc"
)

// Server re-exports the grpc.Server type.
type Server = grpc.Server

// NewServer re-exports the grpc.NewServer function.
func NewServer() *Server {
	return grpc.NewServer()
}
