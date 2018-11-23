package server

import (
	"errors"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"golang.org/x/net/context"
)

// PingService ...
type PingService struct {
	running             bool
	grpcServerProvider  *provider.GRPCServer
	grpcGatewayProvider *provider.GRPCGateway
}

// NewPingService ...
func NewPingService(grpcServerProvider *provider.GRPCServer, grpcGatewayProvider *provider.GRPCGateway) *PingService {
	return &PingService{
		running:             false,
		grpcServerProvider:  grpcServerProvider,
		grpcGatewayProvider: grpcGatewayProvider,
	}
}

// Init ...
func (s *PingService) Init() error {
	RegisterPingServiceServer(s.grpcServerProvider.Server, s)
	return nil
}

func (s *PingService) Run() error {
	if err := s.grpcGatewayProvider.RegisterServices(RegisterPingServiceHandler); err != nil {
		ctxlogrus.Extract(context.Background()).WithError(err).Errorf("Could not register gateway service handlers")
		return err
	}
	return nil
}

func (s *PingService) IsRunning() bool {
	return s.running
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
