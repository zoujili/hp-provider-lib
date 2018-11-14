package main

import (
	"fitstation-hp/lib-fs-provider-go/pkg/v1/provider"
	"fitstation-hp/lib-fs-provider-go/pkg/v1/stack"
)

func main() {
	stack := stack.New()
	defer stack.MustClose()

	logrusConfig := provider.NewLogrusConfigFromEnv()
	logrusProvider := provider.NewLogrus(logrusConfig)
	stack.MustInit(logrusProvider)

	appConfig := provider.NewAppConfigFromEnv()
	appProvider := provider.NewApp(appConfig)
	stack.MustInit(appProvider)

	prometheusConfig := provider.NewPrometheusConfigFromEnv()
	prometheusProvider := provider.NewPrometheus(prometheusConfig)
	stack.MustInit(prometheusProvider)

	jaegerConfig := provider.NewJaegerConfigFromEnv()
	jaegerProvider := provider.NewJaeger(jaegerConfig, appProvider)
	stack.MustInit(jaegerProvider)

	pprofConfig := provider.NewPProfConfigFromEnv()
	pprofProvider := provider.NewPProf(pprofConfig)
	stack.MustInit(pprofProvider)

	probesConfig := provider.NewProbesConfigFromEnv()
	probesProvider := provider.NewProbes(probesConfig)
	stack.MustInit(probesProvider)

	mongodbConfig := provider.NewMongoDBConfigFromEnv()
	mongodbProvider := provider.NewMongoDB(mongodbConfig, probesProvider, appProvider)
	stack.MustInit(mongodbProvider)

	natsConfig := provider.NewNatsConfigFromEnv()
	natsProvider := provider.NewNats(natsConfig, probesProvider)
	stack.MustInit(natsProvider)

	grpcServerConfig := provider.NewGRPCServerConfigFromEnv()
	grpcServerProvider := provider.NewGRPCServer(grpcServerConfig)
	stack.MustInit(grpcServerProvider)

	// Do other stuff here

	stack.MustRun()
}
