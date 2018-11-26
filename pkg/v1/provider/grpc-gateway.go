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
	"net/http"
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
	Config  *GRPCGatewayConfig
	running bool

	GRPCServer *GRPCServer
	ClientConn *grpc.ClientConn
	ServeMux   *runtime.ServeMux
}

func NewGRPCGateway(config *GRPCGatewayConfig, server *GRPCServer) *GRPCGateway {
	return &GRPCGateway{
		Config:     config,
		running:    false,
		GRPCServer: server,
	}
}

func (p *GRPCGateway) Init() error {
	if p.GRPCServer == nil {
		return fmt.Errorf("cannot use GRPCGateway without GRPCServer")
	}
	return nil
}

func (p *GRPCGateway) Run() error {
	if !p.Config.Enabled {
		logrus.Info("GRPCGateway Provider Not Enabled")
		return nil
	}

	if err := WaitForRunningProvider(p.GRPCServer, 2); err != nil {
		return err
	}

	serverAddr := p.GRPCServer.Listener.Addr().String()
	addr := fmt.Sprintf(":%d", p.Config.Port)

	logger := logrus.WithFields(logrus.Fields{
		"serverAddr": serverAddr,
		"addr":       addr,
	})

	opts := []grpc_logrus.Option{
		grpc_logrus.WithDurationField(func(duration time.Duration) (key string, value interface{}) {
			return "grpc.time_ns", duration.Nanoseconds()
		}),
	}

	conn, err := grpc.DialContext(
		context.Background(),
		serverAddr,
		grpc.WithInsecure(),
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
	p.running = true

	logger.Info("GRPCGateway Provider Launched")
	if err := http.ListenAndServe(addr, p.ServeMux); err != nil {
		logger.WithError(err).Error("GRPCGateway Provider Launch Failed")
		return err
	}

	return nil
}

func (p *GRPCGateway) IsRunning() bool {
	return p.running
}

func (p *GRPCGateway) RegisterServices(functions ...func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error) error {
	if !p.Config.Enabled {
		return nil
	}
	if err := WaitForRunningProvider(p, 2); err != nil {
		return err
	}

	for _, function := range functions {
		if err := function(context.Background(), p.ServeMux, p.ClientConn); err != nil {
			return err
		}
	}
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
