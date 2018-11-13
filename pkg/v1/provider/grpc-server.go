package provider

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// GRPCServerConfig ...
type GRPCServerConfig struct {
	Port int
}

// NewGRPCServerConfigEnv ...
func NewGRPCServerConfigEnv() *GRPCServerConfig {
	viper.SetDefault("GRPCSERVER_PORT", 3000)
	viper.BindEnv("GRPCSERVER_PORT")
	port := viper.GetInt("GRPCSERVER_PORT")

	logrus.WithFields(logrus.Fields{
		"port": port,
	}).Info("GRPCServer Config Initialized")

	return &GRPCServerConfig{
		Port: port,
	}
}

// GRPCServer ...
type GRPCServer struct {
	Config *GRPCServerConfig

	Server *grpc.Server
}

// NewGRPCServer ...
func NewGRPCServer(config *GRPCServerConfig) *GRPCServer {
	return &GRPCServer{
		Config: config,
	}
}

// Init ...
func (p *GRPCServer) Init() error {
	logrusEntry := logrus.NewEntry(logrus.StandardLogger())

	opts := []grpc_logrus.Option{
		grpc_logrus.WithDurationField(func(duration time.Duration) (key string, value interface{}) {
			return "grpc.time_ns", duration.Nanoseconds()
		}),
	}

	p.Server = grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_ctxtags.UnaryServerInterceptor(),
				grpc_opentracing.UnaryServerInterceptor(),
				grpc_prometheus.UnaryServerInterceptor,
				grpc_logrus.UnaryServerInterceptor(logrusEntry, opts...),
				grpc_auth.UnaryServerInterceptor(p.authFunc),
				grpc_recovery.UnaryServerInterceptor(),
			),
		),
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				grpc_ctxtags.StreamServerInterceptor(),
				grpc_opentracing.StreamServerInterceptor(),
				grpc_prometheus.StreamServerInterceptor,
				grpc_logrus.StreamServerInterceptor(logrusEntry, opts...),
				grpc_auth.StreamServerInterceptor(p.authFunc),
				grpc_recovery.StreamServerInterceptor(),
			),
		),
	)

	logrus.Info("GRPCServer Provider Initialized")
	return nil
}

// Run ...
func (p *GRPCServer) Run() error {
	addr := fmt.Sprintf(":%d", p.Config.Port)

	logger := logrus.WithFields(logrus.Fields{
		"addr": addr,
	})

	reflection.Register(p.Server)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.WithError(err).Error("GRPCServer Provider Launch Failed")
		return err
	}

	logger.Info("GRPCServer Provider Launched")
	if err := p.Server.Serve(listener); err != nil {
		logger.WithError(err).Error("GRPCServer Provider Launch Failed")
		return err
	}

	return nil
}

// Close ...
func (p *GRPCServer) Close() error {
	p.Server.GracefulStop()

	logrus.Info("GRPCServer Provider Closed")
	return nil
}

type tokenInfo string

func (p *GRPCServer) authFunc(ctx context.Context) (context.Context, error) {
	// TODO

	return ctx, nil
}
