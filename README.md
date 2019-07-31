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
$ docker run --name mongo-lfpg-example-basic -d --rm -p 27017:27017 mongo
$ docker run --name nats-lfpg-example-basic -d --rm -p 4222:4222 -p 6222:6222 -p 8222:8222 nats
$ LOGRUS_FORMATTER=text APP_NAME=basic go run examples/basic/main.go
$ docker stop mongo-lfpg-example-basic
$ docker stop nats-lfpg-example-basic
```

## Run advanced examples

```shell
// Start Jaeger backend
$ docker run --name jaeger-lfpg-example-ping \
  -d \
  --rm \
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
$ LOGRUS_FORMATTER=text APP_NAME=ping-client go run examples/ping/grpc_client/main.go hello

// Run call that panics on server side
$ LOGRUS_FORMATTER=text APP_NAME=ping-client go run examples/ping/grpc_client/main.go panic

// Run call that errors on server side
$ LOGRUS_FORMATTER=text APP_NAME=ping-client go run examples/ping/grpc_client/main.go error

// open web browser at http://127.0.0.1:16686
// to view traces

$ docker stop jaeger-lfpg-example-ping
```

## Run GraphQL example

```shell
$ docker run --name nats-lfpg-example-graphql -d --rm -p 4222:4222 -p 6222:6222 -p 8222:8222 nats
$ LOGRUS_FORMATTER=text APP_NAME=graphql-server GRAPHQL_GRAPHIQL_ENABLED=true JWT_REQUIRED=false go run examples/graphql/main.go
$ docker stop nats-lfpg-example-graphql
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

The easiest way to implement these interfaces is by extending the Abstract structs. \
Using these structs, you (for example) don't have to add the Init() method if your Provider doesn't need initialization.

The AbstractRunProvider does not provide a Run() method, since any RunProvider should always need logic in that method.

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
logrusConfig := logrus.NewConfigFromEnv()
logrusProvider := logrus.New(logrusConfig)
st.MustInit(logrusProvider)
```

NewConfigFromEnv() config:

| ENV key | ENV value | Default value | Description |
| --- | --- | --- | --- |
| LOGRUS_LEVEL | string [panic, fatal, error, warn, info, debug]| info | Minimum level to log |
| LOGRUS_FORMATTER | string [json, text, text_clr]| json | Type of log output <br>Use 'text_crl' instead of 'text' to force colors in Intellij |
| LOGRUS_OUTPUT | string [stderr, stdout]| stderr | Log output |

---

### AppProvider

Stores App info like name, version, ...

```go
appConfig := app.NewConfigEnv()
appProvider := app.NewApp(appConfig)
st.MustInit(appProvider)
```

NewConfigFromEnv() config:

| ENV key | ENV value | Default value | Description |
| --- | --- | --- | --- |
| APP_NAME | string | os.Args[0] = name of the binary | Application name |
| BASE_PATH | string | / | Application base path<br>Will be prefixed to all provider paths |

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
prometheusConfig := prometheus.NewConfigFromEnv()
prometheusProvider := prometheus.New(prometheusConfig)
st.MustInit(prometheusProvider)
```

NewConfigFromEnv() config:

| ENV key | ENV value | Default value | Description |
| --- | --- | --- | --- |
| PROMETHEUS_ENABLED | bool | true | |
| PROMETHEUS_PORT | int | 9090 | HTTP server port |
| PROMETHEUS_ENDPOINT | string | /metrics | Path to expose metrics on |

---

### JaegerProvider

Will setup global OpenTracing with Jaeger backend.

```go
jaegerConfig := jaeger.NewConfigFromEnv()
jaegerProvider := jaeger.New(jaegerConfig, appProvider)
st.MustInit(jaegerProvider)
```

NewConfigFromEnv() config:

| ENV key | ENV value | Default value | Description |
| --- | --- | --- | --- |
| JAEGER_ENABLED | bool | true | |
| JAEGER_AGENT_PORT | port | 6831 | Port on which the agent is running |
| JAEGER_AGENT_HOST | string | 127.0.0.1 | Hostname on which the agent is running |

---

### PProfProvider

Will setup a HTTP server with PPROF endpoint to allow profiling.

```go
pprofConfig := pprof.NewConfigFromEnv()
pprofProvider := pprof.New(pprofConfig)
st.MustInit(pprofProvider)
```

NewConfigFromEnv() config:

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
probesConfig := probes.NewConfigFromEnv()
probesProvider := probes.New(probesConfig, appProvider)
st.MustInit(probesProvider)
```

NewConfigFromEnv() config:

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
mongodbConfig := mongodb.NewConfigFromEnv()
mongodbProvider := mongodb.New(mongodbConfig, probesProvider, appProvider)
st.MustInit(mongodbProvider)
```

NewConfigFromEnv() config:

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
natsConfig := nats.NewConfigFromEnv()
natsProvider := nats.New(natsConfig, probesProvider)
st.MustInit(natsProvider)
```

NewConfigFromEnv() config:

| ENV key | ENV value | Default value | Description |
| --- | --- | --- | --- |
| NATS_ENABLED | bool | true | |
| NATS_URI | string | nats://127.0.0.1:4222 | Nats server URI |
| NATS_TIMEOUT | int | 20 | Max time to connect to Nats |

---

### GRPCServerProvider

Will setup a GRPC server.

```go
grpcServerConfig := grpc.NewConfigFromEnv()
grpcServerProvider := grpc.New(grpcServerConfig)
st.MustInit(grpcServerProvider)
```

NewConfigFromEnv() config:

| ENV key | ENV value | Default value | Description |
| --- | --- | --- | --- |
| GRPC_PORT | int | 3000 | GRPC server port  |
| GRPC_LOG_PAYLOAD | bool | false | Enable to log incoming and outgoing messages |

---

### GRPCGatewayProvider

Will setup a HTTP server to act as REST gateway to the GRPC Server.

```go
grpcGatewayConfig := gateway.NewConfigFromEnv()
grpcGatewayProvider := gateway.New(grpcGatewayConfig, grpcServerProvider)
st.MustInit(grpcGatewayProvider)
```

NewConfigFromEnv() config:

| ENV key | ENV value | Default value | Description |
| --- | --- | --- | --- |
| GRPC_GATEWAY_ENABLED | bool | true | |
| GRPC_GATEWAY_PORT | int | 8080 | HTTP server port |
| GRPC_GATEWAY_LOG_PAYLOAD | bool | false | Enable to log incoming and outgoing messages |

---

### GraphQLProvider

Will setup a HTTP server on which to expose a GraphQL endpoint.

The middlewareChain is optional. For more info, see the [Middlewares](#Middlewares) section.

```go
graphqlConfig := graphql.NewConfigFromEnv()
graphqlProvider := graphql.New(graphqlConfig, ...middlewareChain)
st.MustInit(graphqlProvider)
```

NewConfigFromEnv() config:

| ENV key | ENV value | Default value | Description |
| --- | --- | --- | --- |
| GRAPHQL_PORT | int | 3030 | HTTP server port |
| GRAPHQL_GRAPHIQL_ENABLED | bool | false | If set, will enable a [GraphiQL](https://github.com/graphql/graphiql) in-browser client on path '/graphiql' |

---

### ProxyProvider

Creates a reverse proxy to an external service.

```go
exampleProxyConfig := proxy.NewConfigFromEnv("EXAMPLE_SERVICE")
exampleProxy := proxy.New(exampleProxyConfig)
st.MustInit(exampleProxy)
```

NewConfigFromEnv() config:

| ENV key             | ENV value | Default value         | Description                                |
| ------------------- | --------- | --------------------- | ------------------------------------------ |
| {PREFIX}_ENABLED    | bool      | true                  | Can be used to disable the proxy           |
| {PREFIX}_DEBUG      | bool      | false                 | Can be used to enable payload logging |
| {PREFIX}_PORT       | int       | 4040                  | Port on which the proxy is listening       |
| {PREFIX}_ENDPOINT   | string    | /                     | Endpoint on which the proxy is listening   |
| {PREFIX}_TARGET_URL | string    | http://localhost:8080 | Absolute URL to the service                |

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
jwtConfig := middleware.NewConfigFromEnv()
jwtMiddleware := middleware.New(jwtConfig)
```

NewConfigFromEnv() config:

| ENV key | ENV value | Default value | Description |
| --- | --- | --- | --- |
| JWT_REQUIRED | bool | true | If true, missing JWT will lead to 401 Unauthorized error |
| JWT_VALID | bool | true | If true, invalid JWT will lead to 401 Unauthorized error |

# Examples

## Example GRPC-based service

```go
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

    // Logging
    logrusConfig := logrus.NewConfigFromEnv()
    logrusProvider := logrus.New(logrusConfig)
    st.MustInit(logrusProvider)

    // Root app
    appConfig := app.NewConfigFromEnv()
    appProvider := app.New(appConfig)
    st.MustInit(appProvider)

    // Prometheus (metrics)
    prometheusConfig := prometheus.NewConfigFromEnv()
    prometheusProvider := prometheus.New(prometheusConfig)
    st.MustInit(prometheusProvider)

    // Jaeger (tracing)
    jaegerConfig := jaeger.NewConfigFromEnv()
    jaegerProvider := jaeger.New(jaegerConfig, appProvider)
    st.MustInit(jaegerProvider)

    // PProf (profiling)
    pprofConfig := pprof.NewConfigFromEnv()
    pprofProvider := pprof.New(pprofConfig)
    st.MustInit(pprofProvider)

    // Probes (liveness/readiness for Kubernetes)
    probesConfig := probes.NewConfigFromEnv()
    probesProvider := probes.New(probesConfig, appProvider)
    st.MustInit(probesProvider)

    // MongoDB
    mongodbConfig := mongodb.NewConfigFromEnv()
    mongodbProvider := mongodb.New(mongodbConfig, probesProvider, appProvider)
    st.MustInit(mongodbProvider)

    // Nats (events)
    natsConfig := nats.NewConfigFromEnv()
    natsProvider := nats.New(natsConfig, probesProvider)
    st.MustInit(natsProvider)

    // gRPC Server
    grpcServerConfig := grpc.NewConfigFromEnv()
    grpcServerProvider := grpc.New(grpcServerConfig)
    st.MustInit(grpcServerProvider)

    // gRPC Gateway
    grpcGatewayConfig := gateway.NewConfigFromEnv()
    grpcGatewayProvider := gateway.New(grpcGatewayConfig, grpcServerProvider)
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
    "github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/middleware/jwt"
    "github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider/app"
    "github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider/graphql"
    "github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider/jaeger"
    "github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider/logrus"
    "github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider/mongodb"
    "github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider/pprof"
    "github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider/probes"
    "github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider/prometheus"
    "github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/stack"
)

func main() {
    st := stack.New()
    defer st.MustClose()

    // Logging
    logrusConfig := logrus.NewConfigFromEnv()
    logrusProvider := logrus.New(logrusConfig)
    st.MustInit(logrusProvider)

    // Root app
    appConfig := app.NewConfigFromEnv()
    appProvider := app.New(appConfig)
    st.MustInit(appProvider)

    // Prometheus (metrics)
    prometheusConfig := prometheus.NewConfigFromEnv()
    prometheusProvider := prometheus.New(prometheusConfig)
    st.MustInit(prometheusProvider)

    // Jaeger (tracing)
    jaegerConfig := jaeger.NewConfigFromEnv()
    jaegerProvider := jaeger.New(jaegerConfig, appProvider)
    st.MustInit(jaegerProvider)

    // PProf (profiling)
    pprofConfig := pprof.NewConfigFromEnv()
    pprofProvider := pprof.New(pprofConfig)
    st.MustInit(pprofProvider)

    // Probes (liveness/readiness for Kubernetes)
    probesConfig := probes.NewConfigFromEnv()
    probesProvider := probes.New(probesConfig)
    st.MustInit(probesProvider)

    // Middleware
    jwtConfig := jwt.NewConfigFromEnv()
    jwtMiddleware := jwt.New(jwtConfig)

    // GraphQL
    graphqlConfig := graphql.NewConfigFromEnv()
    graphqlProvider := graphql.New(graphqlConfig, jwtMiddleware)
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
