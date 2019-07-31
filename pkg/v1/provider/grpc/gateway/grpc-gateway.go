package gateway

import (
	"context"
	"fmt"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider/app"
	server "github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider/grpc"
	"github.com/gogo/gateway"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net/http"
	"time"
)

// GRPC Gateway Provider.
// Provides a gateway that allows clients to perform REST calls to the GRPC server.
// Needs to know the individual GRPC providers to know where to send its messages.
type Gateway struct {
	provider.AbstractRunProvider

	Config      *Config
	server      *server.Server
	appProvider *app.App

	client *grpc.ClientConn
	mux    *runtime.ServeMux
}

// Creates a GRPC Gateway Provider.
// Relies on the server to know where to forward the REST messages.
func New(config *Config, server *server.Server, appProvider *app.App) *Gateway {
	return &Gateway{
		Config:      config,
		server:      server,
		appProvider: appProvider,
	}
}

func (p *Gateway) Run() error {
	if !p.Config.Enabled {
		logrus.Info("GRPC Gateway Provider not enabled")
		return nil
	}

	if err := provider.WaitForRunningProvider(p.server, 2); err != nil {
		return err
	}

	basePath := p.appProvider.Config.BasePath + "/"
	serverAddr := p.server.Listener.Addr().String()
	addr := fmt.Sprintf(":%d", p.Config.Port)

	logEntry := logrus.WithFields(logrus.Fields{
		"basePath":   basePath,
		"serverAddr": serverAddr,
		"addr":       addr,
	})

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

	conn, err := grpc.DialContext(
		context.Background(),
		serverAddr,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(unaryInterceptors...)),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(streamInterceptors...)),
	)
	if err != nil {
		logEntry.WithError(err).Errorf("GRPC Gateway could not connect to GRPC server")
		return err
	}

	jsonpb := &gateway.JSONPb{
		EmitDefaults: true,
		Indent:       "  ",
		OrigName:     true,
	}

	p.mux = runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, jsonpb),
		runtime.WithProtoErrorHandler(runtime.DefaultHTTPProtoErrorHandler),
	)
	p.client = conn
	p.SetRunning(true)

	logEntry.Info("GRPC Gateway Provider launched")
	if err := http.ListenAndServe(addr, NewMuxWrapper(basePath, p.mux)); err != nil {
		logEntry.WithError(err).Error("GRPC Gateway Provider launch failed")
		return err
	}

	return nil
}

// Used to register the GRPC providers.
// The Gateway isn't able to use the same reflection based functionality as the GRPC Provider, therefor this is needed.
func (p *Gateway) RegisterServices(functions ...func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error) error {
	if !p.Config.Enabled {
		return nil
	}
	if err := provider.WaitForRunningProvider(p.server, 2); err != nil {
		return err
	}

	for _, function := range functions {
		if err := function(context.Background(), p.mux, p.client); err != nil {
			return err
		}
	}
	return nil
}

// Closes the connection to the GRPC Provider.
func (p *Gateway) Close() error {
	if !p.Config.Enabled {
		return nil
	}

	if p.client != nil {
		if err := p.client.Close(); err != nil {
			logrus.WithError(err).Errorf("Error while closing GRPC Gateway connection")
			return err
		}
	}

	return nil
}

func (p *Gateway) logDeciderFunc(ctx context.Context, fullMethodName string) bool {
	// TODO: Should we really log everything?
	return true
}
