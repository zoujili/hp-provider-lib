package main

import (
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider/app"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider/grpc"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider/grpc/gateway"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider/jaeger"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider/logrus"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider/mongodb"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider/nats"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider/pprof"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider/probes"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider/prometheus"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/stack"
)

func main() {
	st := stack.New()
	defer st.MustClose()

	logrusConfig := logrus.NewConfigFromEnv()
	logrusProvider := logrus.New(logrusConfig)
	st.MustInit(logrusProvider)

	appConfig := app.NewConfigFromEnv()
	appProvider := app.New(appConfig)
	st.MustInit(appProvider)

	prometheusConfig := prometheus.NewConfigFromEnv()
	prometheusProvider := prometheus.New(prometheusConfig)
	st.MustInit(prometheusProvider)

	jaegerConfig := jaeger.NewConfigFromEnv()
	jaegerProvider := jaeger.New(jaegerConfig, appProvider)
	st.MustInit(jaegerProvider)

	pprofConfig := pprof.NewConfigFromEnv()
	pprofProvider := pprof.New(pprofConfig)
	st.MustInit(pprofProvider)

	probesConfig := probes.NewConfigFromEnv()
	probesProvider := probes.New(probesConfig)
	st.MustInit(probesProvider)

	mongodbConfig := mongodb.NewConfigFromEnv()
	mongodbProvider := mongodb.New(mongodbConfig, probesProvider, appProvider)
	st.MustInit(mongodbProvider)

	natsConfig := nats.NewConfigFromEnv()
	natsProvider := nats.New(natsConfig, probesProvider)
	st.MustInit(natsProvider)

	grpcServerConfig := grpc.NewConfigFromEnv()
	grpcServerProvider := grpc.New(grpcServerConfig)
	st.MustInit(grpcServerProvider)

	grpcGatewayConfig := gateway.NewConfigFromEnv()
	grpcGatewayProvider := gateway.New(grpcGatewayConfig, grpcServerProvider)
	st.MustInit(grpcGatewayProvider)

	// Do other stuff here

	st.MustRun()
}
