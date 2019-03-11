package connection

import (
	"context"
	"fmt"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider/probes"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"time"
)

// GRPC Connection Provider.
// Provides a stable connection to a GRPC server.
type Connection struct {
	provider.AbstractRunProvider

	Config         *Config
	Conn           *grpc.ClientConn
	Health         grpc_health_v1.HealthClient
	probesProvider *probes.Probes
}

// Creates a GRPC Connection Provider.
func New(config *Config, probesProvider *probes.Probes) *Connection {
	return &Connection{
		Config:         config,
		probesProvider: probesProvider,
	}
}

// Establishes the gRPC connection.
func (p *Connection) Run() error {
	addr := fmt.Sprintf("%s:%d", p.Config.Host, p.Config.Port)
	logEntry := logrus.WithFields(logrus.Fields{
		"service": p.Config.Prefix,
		"addr":    addr,
	})
	logEntry.Info("Establishing GRPC connection")

	ctx, cancel := context.WithTimeout(context.Background(), p.Config.Timeout)
	defer cancel()

	logOpts := []grpc_logrus.Option{
		grpc_logrus.WithDurationField(func(duration time.Duration) (key string, value interface{}) {
			return "grpc.time_ns", duration.Nanoseconds()
		}),
	}

	// Unary and streaming have the same interceptors.
	unaryInterceptors := []grpc.UnaryClientInterceptor{
		grpc_opentracing.UnaryClientInterceptor(),
		grpc_prometheus.UnaryClientInterceptor,
		grpc_logrus.UnaryClientInterceptor(logEntry, logOpts...),
	}
	streamInterceptors := []grpc.StreamClientInterceptor{
		grpc_opentracing.StreamClientInterceptor(),
		grpc_prometheus.StreamClientInterceptor,
		grpc_logrus.StreamClientInterceptor(logEntry, logOpts...),
	}

	// Payload is only logged by the server if it was configured to do so.
	if p.Config.LogPayload {
		unaryInterceptors = append(unaryInterceptors, grpc_logrus.PayloadUnaryClientInterceptor(logEntry, p.logDeciderFunc))
		streamInterceptors = append(streamInterceptors, grpc_logrus.PayloadStreamClientInterceptor(logEntry, p.logDeciderFunc))
	}

	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithBlock(),
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{}),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(unaryInterceptors...)),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(streamInterceptors...)),
	)
	if err != nil {
		logEntry.WithError(err).Error("GRPC connection could not be created")
		return err
	}

	p.Conn = conn
	logEntry.Info("GRPC connection opened")
	p.initHealthClient()
	return nil
}

func (p *Connection) Close() error {
	if p.Conn != nil {
		if err := p.Conn.Close(); err != nil {
			logrus.WithError(err).Error("Could not close GRPC connection")
			return err
		}
	}

	return nil
}

func (p *Connection) CheckHealth(ctx context.Context) error {
	if !p.Config.EnableHealth {
		return nil
	}
	req := &grpc_health_v1.HealthCheckRequest{}
	res, err := p.Health.Check(context.Background(), req)
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"service": p.Config.Prefix,
		"status":  res.Status,
	}).Debug("GRPC Connection health check performed")

	if res.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		return fmt.Errorf("unhealthy response from GRPC server: %s", res.Status.String())
	}
	return nil
}

func (p *Connection) logDeciderFunc(ctx context.Context, fullMethodName string) bool {
	// TODO: Should we really log everything?
	return true
}

func (p *Connection) initHealthClient() {
	if !p.Config.EnableHealth {
		logrus.WithField("service", p.Config.Prefix).Debug("GRPC Connection health disabled.")
	}
	p.Health = grpc_health_v1.NewHealthClient(p.Conn)
}
