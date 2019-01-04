# lib-fs-provider-go

FitStation provider library for golang, containing commonly used components to use within services.

Use the Stack to control providers (initialization, running and closing) by adding them. \
Examples can be found at the [bottom of this document](#Examples).

## Installation

If you want to use this in a service, just "go get it".
```shell
go get -u github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go
```

To work on this library, fork it and clone the repository outside of your $GOPATH. \
It has to be outside, since Go 1.11 by default doesn't support modules inside the $GOPATH.

## Run basic examples

```shell
$ LOGRUS_FORMATTER=text APP_NAME=basic go run examples/basic/main.go
```

## Run advanced examples

```shell
// Start Jaeger backend
$ docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.7

// Run Server
$ LOGRUS_FORMATTER=text APP_NAME=ping-server go run examples/ping/server/cmd/main.go

// Run normal call
$ LOGRUS_FORMATTER=text APP_NAME=ping-client go run examples/ping/client/main.go hello

// Run call that panics on server side
$ LOGRUS_FORMATTER=text APP_NAME=ping-client go run examples/ping/client/main.go panic

// Run call that errors on server side
$ LOGRUS_FORMATTER=text APP_NAME=ping-client go run examples/ping/client/main.go error

// open web browser at http://127.0.0.1:16686
// to view traces
```

## Providers

Providers are used to add common functionality to the services. \
They are initialized by creating a config and using that to create the service.

Anything can be a provider, as long as it has the Init() and Close() methods. \
Some providers are runnable, they also need to Run() and IsRunning() methods.

```go
type Provider interface {
    Init() error
    Close() error
}

type RunProvider interface {
    Provider

    Run() error
    IsRunning() bool
}
```

#### Initialization and launching

Providers are always first initialized (as soon as st.MustInit() is called). \
Only once st.MustRun() is called (after all providers are initialized), the runnable providers wil be launched. \

The Run() methods are called one at a time (in the same order as initialization), but in separate go routines. \
This means slower providers might finish later, even if they were launched earlier.

If a provider is dependant on another provider, use the provider.WaitForRunningProvider() method to wait for it to start. \
This method is normally called inside the Run() method of the provider.

---

### LogrusProvider

Will configure logrus global logger to use the defined configuration.

```go
logrusConfig := provider.NewLogrusConfigFromEnv()
logrusProvider := provider.NewLogrus(logrusConfig)
st.MustInit(logrusProvider)
```

NewLogrusConfigFromEnv() config:

| ENV key | ENV value | Default value | Description |
| --- | --- | --- | --- |
| LOGRUS_LEVEL | string [panic, fatal, error, warn, info, debug]| info | Minimum level to log |
| LOGRUS_FORMATTER | string [json, text, text_clr]| json | Type of log output <br>Use 'text_crl' instead of 'text' to force colors in Intellij |
| LOGRUS_OUTPUT | string [stderr, stdout]| stderr | Log output |

---

### AppProvider

Stores App info like name, version, ...

```go
appConfig := provider.NewAppConfigEnv()
appProvider := provider.NewApp(appConfig)
st.MustInit(appProvider)
```

NewAppConfigEnv() config:

| ENV key | ENV value | Default value | Description |
| --- | --- | --- | --- |
| APP_NAME | string | os.Args[0] = name of the binary | Application name |

App provider exposes methods

```go
appProvider.Name()
appProvider.Version() // Injected by compiler
st.MustInit(appProvider)
```

---

### PrometheusProvider

Will setup a HTTP server with Prometheus endpoint to expose metrics.

```go
prometheusConfig := provider.NewPrometheusConfigFromEnv()
prometheusProvider := provider.NewPrometheus(prometheusConfig)
st.MustInit(prometheusProvider)
```

NewPrometheusConfigFromEnv() config:

| ENV key | ENV value | Default value | Description |
| --- | --- | --- | --- |
| PROMETHEUS_ENABLED | bool | true | |
| PROMETHEUS_PORT | int | 9090 | HTTP server port |
| PROMETHEUS_ENDPOINT | string | /metrics | Path to expose metrics on |

---

### JaegerProvider

Will setup global OpenTracing with Jaeger backend.

```go
jaegerConfig := provider.NewJaegerConfigFromEnv()
jaegerProvider := provider.NewJaeger(jaegerConfig, appProvider)
st.MustInit(jaegerProvider)
```

NewJaegerConfigFromEnv() config:

| ENV key | ENV value | Default value | Description |
| --- | --- | --- | --- |
| JAEGER_ENABLED | bool | true | |
| JAEGER_AGENT_PORT | port | 6831 | Port on which the agent is running |
| JAEGER_AGENT_HOST | string | 127.0.0.1 | Hostname on which the agent is running |

---

### PProfProvider

Will setup a HTTP server with PPROF endpoint to allow profiling.

```go
pprofConfig := provider.NewPProfConfigFromEnv()
pprofProvider := provider.NewPProf(pprofConfig)
st.MustInit(pprofProvider)
```

NewPProfConfigFromEnv() config:

| ENV key | ENV value | Default value | Description |
| --- | --- | --- | --- |
| PPROF_ENABLED | bool | true | |
| PPROF_PORT | int | 9999 | HTTP server port |
| PPROF_ENDPOINT | string | /debug/pprof | Path to expose profiling data on |

---

### ProbesProvider

Will setup a HTTP server with Liveness and Readiness endpoints. \
These are mainly used by Kubernetes to check the state of the application.

```go
probesConfig := provider.NewProbesConfigEnv()
probesProvider := provider.NewProbes(probesConfig)
st.MustInit(probesProvider)
```

NewProbesConfigEnv() config:

| ENV key | ENV value | Default value | Description |
| --- | --- | --- | --- |
| PROBES_ENABLED | bool | true | |
| PROBES_PORT | int | 8000 | HTTP server port |
| PROBES_LIVENESS_ENDPOINT | string | /healthz | Path to expose health on |
| PROBES_READINESS_ENDPOINT | string | /ready | Path to expose readiness on |

You can easily add probes to this provider.

```go
probesProvider.AddLivenessProbes(func() error {
    // return err if probe should fail, or nil on success
})

probesProvider.AddReadinessProbes(func() error {
    // return err if probe should fail, or nil on success
})
```

---

### MongoDBProvider

Will setup a reusable connection to MongoDB server.

Usage of the probeProvider is optional. If set, the probe will return successfully when its able to ping the DB.

```go
mongodbConfig := provider.NewMongoDBConfigEnv()
mongodbProvider := provider.NewMongoDB(mongodbConfig, probesProvider, appProvider)
st.MustInit(mongodbProvider)
```

NewMongoDBConfigEnv() config:

| ENV key | ENV value | Default value | Description |
| --- | --- | --- | --- |
| MONGODB_URI | string | mongodb://127.0.0.1:27017 | MongoDB server URI |
| MONGODB_DATABASE | string | test | Database name |
| MONGODB_TIMEOUT | int | 20 | Max time for initial connection |
| MONGODB_MAX_POOL_SIZE | int | 16 | Max connection pool size |
| MONGODB_MAX_CONN_IDLE_TIME | int | 30 | Max time for idle connections to be stopped |
| MONGODB_HEARTBEAT_INTERVAL | int | 15 | Interval between connection checks |

---

### NatsProvider

Will setup a reusable connection to Nats (events).

Usage of the probeProvider is optional, if set: the probe will return successfully when its connected to Nats.

```go
natsConfig := provider.NewNatsConfigEnv()
natsProvider := provider.NewNats(natsConfig, probesProvider)
st.MustInit(natsProvider)
```

NewNatsConfigEnv() config:

| ENV key | ENV value | Default value | Description |
| --- | --- | --- | --- |
| NATS_ENABLED | bool | true | |
| NATS_URI | string | nats://127.0.0.1:4222 | Nats server URI |
| NATS_TIMEOUT | int | 20 | Max time to connect to Nats |

---

### GRPCServerProvider

Will setup a GRPC server.

```go
grpcServerConfig := provider.NewGRPCServerConfigEnv()
grpcServerProvider := provider.NewGRPCServer(grpcServerConfig)
st.MustInit(grpcServerProvider)
```

NewGRPCServerConfigEnv() config:

| ENV key | ENV value | Default value | Description |
| --- | --- | --- | --- |
| GRPCSERVER_PORT | int | 3000 | GRPC server port  |
| GRPCSERVER_LOG_PAYLOAD | bool | false | Enable to log incoming and outgoing messages |

---

### GRPCGatewayProvider

Will setup a HTTP server to act as REST gateway to the GRPC Server.

```go
grpcGatewayConfig := provider.NewGRPCGatewayConfigFromEnv()
grpcGatewayProvider := provider.NewGRPCGateway(grpcGatewayConfig, grpcServerProvider)
st.MustInit(grpcGatewayProvider)
```

NewGRPCGatewayConfigFromEnv() config:

| ENV key | ENV value | Default value | Description |
| --- | --- | --- | --- |
| GRPCGATEWAY_ENABLED | bool | true | |
| GRPCGATEWAY_PORT | int | 8080 | HTTP server port |
| GRPCGATEWAY_LOG_PAYLOAD | bool | false | Enable to log incoming and outgoing messages |

---

### GraphQLProvider

Will setup a HTTP server on which to expose a GraphQL endpoint.

The middlewareChain is optional. For more info, see the [Middlewares](#Middlewares) section.

```go
graphqlConfig := provider.NewGraphQLConfigFromEnv()
graphqlProvider := provider.NewGraphQL(graphqlConfig, ...middlewareChain)
st.MustInit(graphqlProvider)
```

NewGraphQLConfigFromEnv() config:

| ENV key | ENV value | Default value | Description |
| --- | --- | --- | --- |
| GRAPHQL_PORT | int | 3030 | HTTP server port |
| GRAPHQL_GRAPHIQL_ENABLED | bool | false | If set, will enable a [GraphiQL](https://github.com/graphql/graphiql) in-browser client on path '/graphiql' |

---

## Middlewares
Middlewares are used to add extra functionality around an existing HTTP handler. \
Only some providers support middlewares.

```go
type Middleware interface {
    Handler(next http.Handler) http.Handler
}
```

Middlewares are chained in the order they are given to the provider. \
Each middleware calls the next handler once it's finished.

Middlewares are not providers and thus do not need the stack to know them.

---

### JWT Middleware

Will decode a JWT in the authorization header and pass this to the context (key: "jwt").

```go
jwtConfig := middleware.NewJWTConfigFromEnv()
jwtMiddleware := middleware.NewJWT(jwtConfig)
```

NewJWTConfigFromEnv() config:

| ENV key | ENV value | Default value | Description |
| --- | --- | --- | --- |
| JWT_REQUIRED | bool | true | If true, missing JWT will lead to 401 Unauthorized error |
| JWT_VALID | bool | true | If true, invalid JWT will lead to 401 Unauthorized error |

# Examples

## Example GRPC-based service

```go
package main

import (
    "github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider"
    "github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/stack"
)

func main() {
    st := stack.New()
    defer st.MustClose()

    // Logging
    logrusConfig := provider.NewLogrusConfigFromEnv()
    logrusProvider := provider.NewLogrus(logrusConfig)
    st.MustInit(logrusProvider)

    // Root app
    appConfig := provider.NewAppConfigFromEnv()
    appProvider := provider.NewApp(appConfig)
    st.MustInit(appProvider)

    // Prometheus (metrics)
    prometheusConfig := provider.NewPrometheusConfigFromEnv()
    prometheusProvider := provider.NewPrometheus(prometheusConfig)
    st.MustInit(prometheusProvider)

    // Jaeger (tracing)
    jaegerConfig := provider.NewJaegerConfigFromEnv()
    jaegerProvider := provider.NewJaeger(jaegerConfig, appProvider)
    st.MustInit(jaegerProvider)

    // PProf (profiling)
    pprofConfig := provider.NewPProfConfigFromEnv()
    pprofProvider := provider.NewPProf(pprofConfig)
    st.MustInit(pprofProvider)

    // Probes (liveness/readiness for Kubernetes)
    probesConfig := provider.NewProbesConfigFromEnv()
    probesProvider := provider.NewProbes(probesConfig)
    st.MustInit(probesProvider)

    // MongoDB
    mongodbConfig := provider.NewMongoDBConfigFromEnv()
    mongodbProvider := provider.NewMongoDB(mongodbConfig, probesProvider, appProvider)
    st.MustInit(mongodbProvider)

    // Nats (events)
    natsConfig := provider.NewNatsConfigFromEnv()
    natsProvider := provider.NewNats(natsConfig, probesProvider)
    st.MustInit(natsProvider)

    // gRPC server
    grpcServerConfig := provider.NewGRPCServerConfigFromEnv()
    grpcServerProvider := provider.NewGRPCServer(grpcServerConfig)
    st.MustInit(grpcServerProvider)

    // gRPC gateway
    grpcGatewayConfig := provider.NewGRPCGatewayConfigFromEnv()
    grpcGatewayProvider := provider.NewGRPCGateway(grpcGatewayConfig, grpcServerProvider)
    st.MustInit(grpcGatewayProvider)

    // Resources

    exampleResource := resource.NewExampleResourceMongo(mongodbProvider)
    st.MustInit(exampleResource)

    // Controllers

    exampleController := controller.NewExampleController(exampleResource, natsProvider)
    st.MustInit(exampleController)

    // Handler
    handler := example.NewHandler(grpcServerProvider, grpcGatewayProvider, exampleController)
    st.MustInit(handler)

    st.MustRun()
}
```

## Example GraphQL-based service
```go
package main

import (
    "github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/middleware"
    "github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider"
    "github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/stack"
)

func main() {
    st := stack.New()
    defer st.MustClose()

    // Logging
    logrusConfig := provider.NewLogrusConfigFromEnv()
    logrusProvider := provider.NewLogrus(logrusConfig)
    st.MustInit(logrusProvider)

    // Root app
    appConfig := provider.NewAppConfigFromEnv()
    appProvider := provider.NewApp(appConfig)
    st.MustInit(appProvider)

    // Prometheus (metrics)
    prometheusConfig := provider.NewPrometheusConfigFromEnv()
    prometheusProvider := provider.NewPrometheus(prometheusConfig)
    st.MustInit(prometheusProvider)

    // Jaeger (tracing)
    jaegerConfig := provider.NewJaegerConfigFromEnv()
    jaegerProvider := provider.NewJaeger(jaegerConfig, appProvider)
    st.MustInit(jaegerProvider)

    // PProf (profiling)
    pprofConfig := provider.NewPProfConfigFromEnv()
    pprofProvider := provider.NewPProf(pprofConfig)
    st.MustInit(pprofProvider)

    // Probes (liveness/readiness for Kubernetes)
    probesConfig := provider.NewProbesConfigFromEnv()
    probesProvider := provider.NewProbes(probesConfig)
    st.MustInit(probesProvider)

    // Nats (events)
    natsConfig := provider.NewNatsConfigFromEnv()
    natsProvider := provider.NewNats(natsConfig, probesProvider)
    st.MustInit(natsProvider)

    // Middleware
    jwtConfig := middleware.NewJTWConfigFromEnv()
    jwtMiddleware := middleware.NewJWT(jwtConfig)

    // GraphQL
    graphqlConfig := provider.NewGraphQLConfigFromEnv()
    graphqlProvider := provider.NewGraphQL(graphqlConfig, jwtMiddleware)
    st.MustInit(graphqlProvider)

    // Resources

    exampleResourceConfig := resource.NewResourceServiceConfig("EXAMPLE_SERVICE")
    exampleResource := resource.NewExampleResourceService(exampleResourceConfig)
    st.MustInit(exampleResource)

    // Root resolver
    rootResolver := resolver.NewRootResolver(exampleResource)
    st.MustInit(rootResolver)

    // Handler
    handler := example.NewHandler(graphqlProvider, rootResolver)
    st.MustInit(handler)

    st.MustRun()
}
```
