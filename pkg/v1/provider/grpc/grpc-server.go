package grpc

import (
	"context"
	"fmt"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider"
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
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// GRPC Server Provider.
// Provides a Server that listens for GRPC traffic and forwards them to the configured handlers.
// Also adds a bunch of useful interceptors for tracing, metrics and so on.
// Relies heavily on https://github.com/grpc-ecosystem packages.
type Server struct {
	provider.AbstractRunProvider

	Config   *Config
	Listener net.Listener
	Server   *grpc.Server
}

// Creates a GRPC Server Provider.
func New(config *Config) *Server {
	return &Server{
		Config: config,
	}
}

// Creates the GRPC Server (doesn't start it yet) and adds useful interceptors.
func (p *Server) Init() error {
	logger := logrus.NewEntry(logrus.StandardLogger())

	opts := []grpc_logrus.Option{
		grpc_logrus.WithDurationField(func(duration time.Duration) (key string, value interface{}) {
			return "grpc.time_ns", duration.Nanoseconds()
		}),
	}

	// Unary and streaming have the same interceptors.
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

	// Payload is only logged by the Server if it was configured to do so.
	if p.Config.LogPayload {
		unaryInterceptors = append(unaryInterceptors, grpc_logrus.PayloadUnaryServerInterceptor(logger, p.logDeciderFunc))
		streamInterceptors = append(streamInterceptors, grpc_logrus.PayloadStreamServerInterceptor(logger, p.logDeciderFunc))
	}

	p.Server = grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryInterceptors...)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streamInterceptors...)),
	)

	return nil
}

// Creates a GRPC Listener on the configured port which is used to start the GRPC Server.
// Uses the GRPC Server reflection functionality find the available handlers.
func (p *Server) Run() error {
	addr := fmt.Sprintf(":%d", p.Config.Port)
	logEntry := logrus.WithField("addr", addr)

	reflection.Register(p.Server)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logEntry.WithError(err).Error("GRPC Listener could not be created")
		return err
	}
	p.Listener = listener
	p.SetRunning(true)

	logEntry.Info("GRPC Provider launched")
	if err := p.Server.Serve(listener); err != nil {
		logEntry.WithError(err).Error("GRPC Provider launch failed")
		return err
	}

	return nil
}

// Shuts down the GRPC Server.
func (p *Server) Close() error {
	p.Server.GracefulStop()
	return nil
}

func (p *Server) authFunc(ctx context.Context) (context.Context, error) {
	// TODO: Add support for authentication.
	return ctx, nil
}

func (p *Server) logDeciderFunc(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
	// TODO: Should we really log everything?
	return true
}
