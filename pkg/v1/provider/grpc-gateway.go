package provider

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"time"
)

// GRPCGatewayConfig ...
type GRPCGatewayConfig struct {
	Enabled bool
	Port    int
}

func NewGRPCGatewayConfigFromEnv() *GRPCGatewayConfig {
	v := viper.New()
	v.SetEnvPrefix("GRPCGATEWAY")
	v.AutomaticEnv()

	v.SetDefault("ENABLED", false)
	enabled := v.GetBool("ENABLED")

	v.SetDefault("PORT", 8080)
	port := v.GetInt("PORT")

	logrus.WithFields(logrus.Fields{
		"enabled": enabled,
		"port":    port,
	}).Debug("GRPCGateway Config Initialized")

	return &GRPCGatewayConfig{
		Enabled: enabled,
		Port:    port,
	}
}

// GRPCGateway ...
type GRPCGateway struct {
	Config *GRPCGatewayConfig

	GRPCServer *GRPCServer
	ClientConn *grpc.ClientConn
	ServeMux   *runtime.ServeMux
}

func NewGRPCGateway(config *GRPCGatewayConfig, server *GRPCServer) *GRPCGateway {
	return &GRPCGateway{
		Config:     config,
		GRPCServer: server,
	}
}

func (p *GRPCGateway) Init() error {
	return nil
}

func (p *GRPCGateway) Run() error {
	if !p.Config.Enabled {
		logrus.Info("GRPCGateway Provider Not Enabled")
		return nil
	}
	if p.GRPCServer == nil {
		return fmt.Errorf("cannot initialize GRPCGateway without GRPCServer to connect to")
	}
	if p.GRPCServer.Listener == nil {
		return fmt.Errorf("cannot run GRPCGateway without starting GRPCServer first")
	}

	serverAddr := p.GRPCServer.Listener.Addr()
	addr := fmt.Sprintf(":%d", p.Config.Port)

	logger := logrus.WithFields(logrus.Fields{
		"serverAddr": serverAddr.String(),
		"addr":       addr,
	})

	opts := []grpc_logrus.Option{
		grpc_logrus.WithDurationField(func(duration time.Duration) (key string, value interface{}) {
			return "grpc.time_ns", duration.Nanoseconds()
		}),
	}

	conn, err := grpc.DialContext(
		context.Background(),
		serverAddr.String(),
		grpc.WithUnaryInterceptor(
			grpc_middleware.ChainUnaryClient(
				grpc_opentracing.UnaryClientInterceptor(),
				grpc_prometheus.UnaryClientInterceptor,
				grpc_logrus.UnaryClientInterceptor(logger, opts...),
			),
		),
		grpc.WithStreamInterceptor(
			grpc_middleware.ChainStreamClient(
				grpc_opentracing.StreamClientInterceptor(),
				grpc_prometheus.StreamClientInterceptor,
				grpc_logrus.StreamClientInterceptor(logger, opts...),
			),
		),
	)
	if err != nil {
		logger.WithError(err).Errorf("GRPCGateway Provider Launch Failed")
		return err
	}

	p.ServeMux = runtime.NewServeMux()
	p.ClientConn = conn

	return nil
}

func (p *GRPCGateway) Close() error {
	if !p.Config.Enabled {
		return nil
	}

	if err := p.ClientConn.Close(); err != nil {
		logrus.WithError(err).Errorf("Error while closing GRPCGateway to GRPCClient connection")
		return err
	}

	return nil
}
