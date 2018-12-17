package main

import (
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/examples/graphql/resolver"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/middleware"
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

	natsConfig := provider.NewNatsConfigFromEnv()
	natsProvider := provider.NewNats(natsConfig, probesProvider)
	st.MustInit(natsProvider)

	jwtConfig := middleware.NewJWTConfigFromEnv()
	jwtMiddleware := middleware.NewJWT(jwtConfig)

	graphqlConfig := provider.NewGraphQLConfigFromEnv()
	graphqlProvider := provider.NewGraphQL(graphqlConfig, jwtMiddleware)
	st.MustInit(graphqlProvider)

	// Do other stuff here

	rootResolver := resolver.NewRootResolver(graphqlProvider)
	st.MustInit(rootResolver)

	st.MustRun()
}
