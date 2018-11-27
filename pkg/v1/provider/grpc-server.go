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
	Port       int
	LogPayload bool
}

// NewGRPCServerConfigFromEnv ...
func NewGRPCServerConfigFromEnv() *GRPCServerConfig {
	v := viper.New()
	v.SetEnvPrefix("GRPCSERVER")
	v.AutomaticEnv()

	v.SetDefault("PORT", 3000)
	port := v.GetInt("PORT")

	v.SetDefault("LOG_PAYLOAD", false)
	logPayload := v.GetBool("LOG_PAYLOAD")

	logrus.WithFields(logrus.Fields{
		"port":       port,
		"logPayload": logPayload,
	}).Debug("GRPCServer Config Initialized")

	return &GRPCServerConfig{
		Port:       port,
		LogPayload: logPayload,
	}
}

// GRPCServer ...
type GRPCServer struct {
	Config  *GRPCServerConfig
	running bool

	Listener net.Listener
	Server   *grpc.Server
}

// NewGRPCServer ...
func NewGRPCServer(config *GRPCServerConfig) *GRPCServer {
	return &GRPCServer{
		Config:  config,
		running: false,
	}
}

// Init ...
func (p *GRPCServer) Init() error {
	logger := logrus.NewEntry(logrus.StandardLogger())

	opts := []grpc_logrus.Option{
		grpc_logrus.WithDurationField(func(duration time.Duration) (key string, value interface{}) {
			return "grpc.time_ns", duration.Nanoseconds()
		}),
	}

	unaryInterceptors := []grpc.UnaryServerInterceptor{
		grpc_ctxtags.UnaryServerInterceptor(),
		grpc_opentracing.UnaryServerInterceptor(),
		grpc_prometheus.UnaryServerInterceptor,
		grpc_logrus.UnaryServerInterceptor(logger, opts...),
		grpc_auth.UnaryServerInterceptor(p.authFunc),
		grpc_recovery.UnaryServerInterceptor(),
	}
	streamInterceptors := []grpc.StreamServerInterceptor{
		grpc_ctxtags.StreamServerInterceptor(),
		grpc_opentracing.StreamServerInterceptor(),
		grpc_prometheus.StreamServerInterceptor,
		grpc_logrus.StreamServerInterceptor(logger, opts...),
		grpc_auth.StreamServerInterceptor(p.authFunc),
		grpc_recovery.StreamServerInterceptor(),
	}

	if p.Config.LogPayload {
		decider := func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
			// TODO: Move the decider somewhere else (maybe to service and use a config?)
			return true
		}
		unaryInterceptors = append(unaryInterceptors, grpc_logrus.PayloadUnaryServerInterceptor(logger, decider))
		streamInterceptors = append(streamInterceptors, grpc_logrus.PayloadStreamServerInterceptor(logger, decider))
	}

	p.Server = grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryInterceptors...)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streamInterceptors...)),
	)

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
	p.Listener = listener
	p.running = true

	logger.Info("GRPCServer Provider Launched")
	if err := p.Server.Serve(listener); err != nil {
		logger.WithError(err).Error("GRPCServer Provider Launch Failed")
		return err
	}

	return nil
}

func (p *GRPCServer) IsRunning() bool {
	return p.running
}

// Close ...
func (p *GRPCServer) Close() error {
	p.Server.GracefulStop()

	return nil
}

//type tokenInfo string

func (p *GRPCServer) authFunc(ctx context.Context) (context.Context, error) {
	// TODO

	return ctx, nil
}
