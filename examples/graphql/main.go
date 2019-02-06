package main

import (
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/examples/graphql/resolver"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/middleware/jwt"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider/app"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider/graphql"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider/jaeger"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider/logrus"
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

	natsConfig := nats.NewConfigFromEnv()
	natsProvider := nats.New(natsConfig, probesProvider)
	st.MustInit(natsProvider)

	jwtConfig := jwt.NewConfigFromEnv()
	jwtMiddleware := jwt.New(jwtConfig)

	graphqlConfig := graphql.NewConfigFromEnv()
	graphqlProvider := graphql.New(graphqlConfig, jwtMiddleware)
	st.MustInit(graphqlProvider)

	// Do other stuff here

	rootResolver := resolver.NewRootResolver(graphqlProvider)
	st.MustInit(rootResolver)

	st.MustRun()
}
