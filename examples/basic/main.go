package main

import (
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/stack"
)

func main() {
	st := stack.New()
	defer st.MustClose()

	logrusConfig := provider.NewLogrusConfigFromEnv()
	logrusProvider := provider.NewLogrus(logrusConfig)
	st.MustInit(logrusProvider)

	appConfig := provider.NewAppConfigFromEnv()
	appProvider := provider.NewApp(appConfig)
	st.MustInit(appProvider)

	prometheusConfig := provider.NewPrometheusConfigFromEnv()
	prometheusProvider := provider.NewPrometheus(prometheusConfig)
	st.MustInit(prometheusProvider)

	jaegerConfig := provider.NewJaegerConfigFromEnv()
	jaegerProvider := provider.NewJaeger(jaegerConfig, appProvider)
	st.MustInit(jaegerProvider)

	pprofConfig := provider.NewPProfConfigFromEnv()
	pprofProvider := provider.NewPProf(pprofConfig)
	st.MustInit(pprofProvider)

	probesConfig := provider.NewProbesConfigFromEnv()
	probesProvider := provider.NewProbes(probesConfig)
	st.MustInit(probesProvider)

	mongodbConfig := provider.NewMongoDBConfigFromEnv()
	mongodbProvider := provider.NewMongoDB(mongodbConfig, probesProvider, appProvider)
	st.MustInit(mongodbProvider)

	natsConfig := provider.NewNatsConfigFromEnv()
	natsProvider := provider.NewNats(natsConfig, probesProvider)
	st.MustInit(natsProvider)

	grpcServerConfig := provider.NewGRPCServerConfigFromEnv()
	grpcServerProvider := provider.NewGRPCServer(grpcServerConfig)
	st.MustInit(grpcServerProvider)

	grpcGatewayConfig := provider.NewGRPCGatewayConfigFromEnv()
	grpcGatewayProvider := provider.NewGRPCGateway(grpcGatewayConfig, grpcServerProvider)
	st.MustInit(grpcGatewayProvider)

	// Do other stuff here

	st.MustRun()
}
