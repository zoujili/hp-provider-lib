package main

import (
	"fitstation-hp/lib-fs-provider-go/examples/ping/server"
	"fitstation-hp/lib-fs-provider-go/pkg/v1/provider"
	"fitstation-hp/lib-fs-provider-go/pkg/v1/stack"
)

func main() {
	stack := stack.New()
	defer stack.MustClose()

	logrusConfig := provider.NewLogrusConfigFromEnv()
	logrusProvider := provider.NewLogrus(logrusConfig)
	stack.MustInit(logrusProvider)

	prometheusConfig := provider.NewPrometheusConfigFromEnv()
	prometheusProvider := provider.NewPrometheus(prometheusConfig)
	stack.MustInit(prometheusProvider)

	jaegerConfig := provider.NewJaegerConfigFromEnv()
	jaegerProvider := provider.NewJaeger(jaegerConfig)
	stack.MustInit(jaegerProvider)

	pprofConfig := provider.NewPProfConfigFromEnv()
	pprofProvider := provider.NewPProf(pprofConfig)
	stack.MustInit(pprofProvider)

	probesConfig := provider.NewProbesConfigEnv()
	probesProvider := provider.NewProbes(probesConfig)
	stack.MustInit(probesProvider)

	grpcServerConfig := provider.NewGRPCServerConfigEnv()
	grpcServerProvider := provider.NewGRPCServer(grpcServerConfig)
	stack.MustInit(grpcServerProvider)

	pingService := server.NewPingService(grpcServerProvider)
	stack.MustInit(pingService)

	stack.MustRun()
}
