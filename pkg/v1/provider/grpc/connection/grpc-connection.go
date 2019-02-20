package connection

import (
	"context"
	"fmt"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"
)

// GRPC Connection Provider.
// Provides a dialed in connection that can be used to create GRPC clients from proto files.
type Connection struct {
	provider.AbstractRunProvider

	Config *Config
	Conn   *grpc.ClientConn
}

// Created a GRPC Connection Provider
func New(config *Config) *Connection {
	return &Connection{
		Config: config,
	}
}

// Creates the GRPC connection
func (p *Connection) Run() error {
	addr := fmt.Sprintf("%s:%d", p.Config.Host, p.Config.Port)
	logEntry := logrus.WithField("addr", addr)

	opts := []grpc_logrus.Option{
		grpc_logrus.WithDurationField(func(duration time.Duration) (key string, value interface{}) {
			return "grpc.time_ns", duration.Nanoseconds()
		}),
	}

	// Unary and streaming have the same interceptors.
	unaryInterceptors := []grpc.UnaryClientInterceptor{
		grpc_opentracing.UnaryClientInterceptor(),
		grpc_prometheus.UnaryClientInterceptor,
		grpc_logrus.UnaryClientInterceptor(logEntry, opts...),
	}
	streamInterceptors := []grpc.StreamClientInterceptor{
		grpc_opentracing.StreamClientInterceptor(),
		grpc_prometheus.StreamClientInterceptor,
		grpc_logrus.StreamClientInterceptor(logEntry, opts...),
	}

	// Payload is only logged by the server if it was configured to do so.
	if p.Config.LogPayload {
		unaryInterceptors = append(unaryInterceptors, grpc_logrus.PayloadUnaryClientInterceptor(logEntry, p.logDeciderFunc))
		streamInterceptors = append(streamInterceptors, grpc_logrus.PayloadStreamClientInterceptor(logEntry, p.logDeciderFunc))
	}

	conn, err := grpc.DialContext(context.Background(),
		addr,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(unaryInterceptors...)),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(streamInterceptors...)),
	)
	if err != nil {
		logEntry.WithError(err).Error("GRPC connection could not be created")
		return err
	}

	p.Conn = conn
	p.SetRunning(true)
	logEntry.Info("GRPC connection opened")
	return nil
}

func (p *Connection) Close() error {
	if err := p.Conn.Close(); err != nil {
		logrus.WithError(err).Error("Could not close GRPC connection")
		return err
	}

	return nil
}

func (p *Connection) logDeciderFunc(ctx context.Context, fullMethodName string) bool {
	// TODO: Should we really log everything?
	return true
}
