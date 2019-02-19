module github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go

require (
	github.azc.ext.hp.com/fitstation-hp/lib-fs-core-go v0.0.0
	github.com/Microsoft/go-winio v0.4.11 // indirect
	github.com/beorn7/perks v0.0.0-20180321164747-3a771d992973 // indirect
	github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1
	github.com/docker/go-connections v0.4.0
	github.com/friendsofgo/graphiql v0.2.0
	github.com/gogo/googleapis v1.1.0
	github.com/gogo/protobuf v1.2.0
	github.com/golang/protobuf v1.2.0
	github.com/graph-gophers/graphql-go v0.0.0-20181116072428-fd99376b56e9
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway v1.5.1
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/mongodb/mongo-go-driver v0.1.0
	github.com/nats-io/gnatsd v1.3.0 // indirect
	github.com/nats-io/go-nats v1.6.0
	github.com/nats-io/nuid v1.0.0 // indirect
	github.com/onsi/ginkgo v1.7.0
	github.com/onsi/gomega v1.4.3
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opentracing/opentracing-go v1.0.2
	github.com/pkg/errors v0.8.0 // indirect
	github.com/prometheus/client_golang v0.9.1
	github.com/prometheus/client_model v0.0.0-20180712105110-5c3871d89910 // indirect
	github.com/prometheus/common v0.0.0-20181116084131-1f2c4f3cd6db // indirect
	github.com/prometheus/procfs v0.0.0-20181005140218-185b4288413d // indirect
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/viper v1.2.1
	github.com/uber-go/atomic v1.3.2 // indirect
	github.com/uber/jaeger-client-go v2.15.0+incompatible
	github.com/uber/jaeger-lib v1.5.0
	go.uber.org/atomic v1.3.2 // indirect
	golang.org/x/net v0.0.0-20181201002055-351d144fa1fc
	google.golang.org/grpc v1.16.0
)

replace github.azc.ext.hp.com/fitstation-hp/lib-fs-core-go => ./../lib-fs-core-go
