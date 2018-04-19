package rpc

import (
	"context"

	"google.golang.org/grpc"
)

type HyperdriveService struct {
}

func NewHyperdriveService() HyperdriveService {
	return HyperdriveService{}
}

func (service *HyperdriveService) Register(server *grpc.Server) {
	RegisterHyperdriveServer(server, service)
}

func (service *HyperdriveService) SendTx(ctx context.Context, tx *Tx) (*Nothing, error) {
	return &Nothing{}, nil
}

func (service *HyperdriveService) SyncBlock(nothing *Nothing, stream Hyperdrive_SyncBlockServer) error {
	return nil
}

func (service *HyperdriveService) Drive(stream Hyperdrive_DriveServer) error {
	return nil
}
