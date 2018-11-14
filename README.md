# lib-fs-provider-go

FitStation provider library for golang, that provide generic an consitent setup of common components.

## Install

Clone repository into $GOPATH/src/fitstation-hp/lib-fs-provider-go

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

### LogrusProvider

Will configure logrus global logger to use the defined configuration

```go
logrusConfig := provider.NewLogrusConfigFromEnv()
logrusProvider := provider.NewLogrus(logrusConfig)
stack.MustInit(logrusProvider)
```

NewLogrusConfigFromEnv() config:

| ENV key | ENV value | Default value |
| --- | --- | --- |
| LOGRUS_LEVEL | string [panic, fatal, error, warn, info, debug]| info |
| LOGRUS_FORMATTER | string [text, json]| json |
| LOGRUS_OUTPUT | string [stderr, stdout]| stderr |

---

### AppProvider

Stores App info like name, version, ...

```go
appConfig := provider.NewAppConfigEnv()
appProvider := provider.NewApp(appConfig)
stack.MustInit(appProvider)
```

NewAppConfigEnv() config:

| ENV key | ENV value | Default value |
| --- | --- | --- |
| APP_NAME | string | os.Args[0] = name of the binary |

App provider exposes methods

```go
appProvider.Name()
appProvider.Version() // Injected by compiler
```

---

### PrometheusProvider

Will setup a http server with prometheus endpoint to expose metrics

```go
prometheusConfig := provider.NewPrometheusConfigFromEnv()
prometheusProvider := provider.NewPrometheus(prometheusConfig)
stack.MustInit(prometheusProvider)
```

NewPrometheusConfigFromEnv() config:

| ENV key | ENV value | Default value |
| --- | --- | --- |
| PROMETHEUS_ENABLED | bool | true |
| PROMETHEUS_PORT | int | 9090 |
| PROMETHEUS_ENDPOINT | string | /metrics |

---

### JaegerProvider

Will setup global opentracing with jaeger backend

```go
jaegerConfig := provider.NewJaegerConfigFromEnv()
jaegerProvider := provider.NewJaeger(jaegerConfig, appProvider)
stack.MustInit(jaegerProvider)
```

NewJaegerConfigFromEnv() config:

| ENV key | ENV value | Default value |
| --- | --- | --- |
| JAEGER_ENABLED | bool | true |
| JAEGER_AGENT_PORT | port | 6831 |
| JAEGER_AGENT_HOST | string | 127.0.0.1 |

---

### PProfProvider

Will setup pprof endpoint to allow profiling

```go
pprofConfig := provider.NewPProfConfigFromEnv()
pprofProvider := provider.NewPProf(pprofConfig)
stack.MustInit(pprofProvider)
```

NewPProfConfigFromEnv() config:

| ENV key | ENV value | Default value |
| --- | --- | --- |
| PPROF_ENABLED | bool | true |
| PPROF_PORT | int | 9999 |
| PPROF_ENDPOINT | string | /debug/pprof |

---

### ProbesProvider

Will setup Liveness and Readiness endpoints

```go
probesConfig := provider.NewProbesConfigEnv()
probesProvider := provider.NewProbes(probesConfig)
stack.MustInit(probesProvider)
```

NewProbesConfigEnv() config:

| ENV key | ENV value | Default value |
| --- | --- | --- |
| PROBES_ENABLED | bool | true |
| PROBES_PORT | int | 8000 |
| PROBES_LIVENESS_ENDPOINT | string | /healthz |
| PROBES_READINESS_ENDPOINT | string | /ready |

You can easaly add probes to this provider

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

Will setup a reusable connection to MongoDB

Usage of the probeProvider is optional, if set: the probe will return succesfully when its able to ping the DB

```go
mongodbConfig := provider.NewMongoDBConfigEnv()
mongodbProvider := provider.NewMongoDB(mongodbConfig, probesProvider, appProvider)
stack.MustInit(mongodbProvider)
```

NewMongoDBConfigEnv() config:

| ENV key | ENV value | Default value | Desc |
| --- | --- | --- | --- |
| MONGODB_URI | string | mongodb://127.0.0.1:27017 | uri |
| MONGODB_DATABASE | string | test | Database to use |
| MONGODB_TIMEOUT | int | 20 | Max time to connect to DB |
| MONGODB_MAX_CONNS_PER_HOST | int | 16 | Max connection pool size |

---

### NatsProvider

Will setup a reusable connection to Nats

Usage of the probeProvider is optional, if set: the probe will return succesfully when its connected to nats

```go
natsConfig := provider.NewNatsConfigEnv()
natsProvider := provider.NewNats(natsConfig, probesProvider)
stack.MustInit(natsProvider)
```

NewNatsConfigEnv() config:

| ENV key | ENV value | Default value | Desc |
| --- | --- | --- | --- |
| NATS_URI | string | nats://127.0.0.1:4222 | uri |
| NATS_TIMEOUT | int | 20 | Max time to connect to Nats |

---

### GRPCServerProvider

Will setup a reusable GRPC server configured with middleware

```go
grpcServerConfig := provider.NewGRPCServerConfigEnv()
grpcServerProvider := provider.NewGRPCServer(grpcServerConfig)
stack.MustInit(grpcServerProvider)
```

NewGRPCServerConfigEnv() config:

| ENV key | ENV value | Default value |
| --- | --- | --- |
| GRPCSERVER_PORT | int | 3000 |

---

## Example usage

```go
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

    appConfig := provider.NewAppConfigEnv()
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

    probesConfig := provider.NewProbesConfigEnv()
    probesProvider := provider.NewProbes(probesConfig)
    stack.MustInit(probesProvider)

    mongodbConfig := provider.NewMongoDBConfigEnv()
    mongodbProvider := provider.NewMongoDB(mongodbConfig, probesProvider, appProvider)
    stack.MustInit(mongodbProvider)

    natsConfig := provider.NewNatsConfigEnv()
    natsProvider := provider.NewNats(natsConfig, probesProvider)
    stack.MustInit(natsProvider)

    grpcServerConfig := provider.NewGRPCServerConfigEnv()
    grpcServerProvider := provider.NewGRPCServer(grpcServerConfig)
    stack.MustInit(grpcServerProvider)

    // Do other stuff here

    stack.MustRun()

}
```