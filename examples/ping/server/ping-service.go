package server

import (
	"errors"
	"fitstation-hp/lib-fs-provider-go/pkg/v1/provider"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	context "golang.org/x/net/context"
)

// PingService ...
type PingService struct {
	grpcServerProvider *provider.GRPCServer
}

// NewPingService ...
func NewPingService(grpcServerProvider *provider.GRPCServer) *PingService {
	return &PingService{
		grpcServerProvider: grpcServerProvider,
	}
}

// Init ...
func (s *PingService) Init() error {
	RegisterPingServiceServer(s.grpcServerProvider.Server, s)
	return nil
}

// Close ...
func (s *PingService) Close() error {
	return nil
}

// Ping ...
func (s *PingService) Ping(ctx context.Context, request *PingRequest) (*PingResponse, error) {
	logger := ctxlogrus.Extract(ctx)
	logger.Info("hello from ping")

	if request.In == "panic" {
		panic("please panic")
	}

	if request.In == "error" {
		return nil, errors.New("please error me")
	}

	return &PingResponse{Out: request.In}, nil
}
