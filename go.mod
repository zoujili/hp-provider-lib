module github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go

require (
	github.azc.ext.hp.com/fitstation-hp/lib-fs-core-go v0.0.0
	github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/friendsofgo/graphiql v0.2.1
	github.com/gogo/googleapis v1.1.0
	github.com/gogo/protobuf v1.2.1
	github.com/golang/protobuf v1.3.0
	github.com/golang/snappy v0.0.1 // indirect
	github.com/graph-gophers/graphql-go v0.0.0-20190225005345-3e8838d4614c
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway v1.8.1
	github.com/nats-io/gnatsd v1.4.1 // indirect
	github.com/nats-io/go-nats v1.7.2
	github.com/nats-io/nkeys v0.0.2 // indirect
	github.com/nats-io/nuid v1.0.0 // indirect
	github.com/onsi/ginkgo v1.7.0
	github.com/onsi/gomega v1.4.3
	github.com/opentracing/opentracing-go v1.0.2
	github.com/pkg/errors v0.8.1 // indirect
	github.com/prometheus/client_golang v0.9.2
	github.com/sirupsen/logrus v1.3.0
	github.com/spf13/viper v1.3.1
	github.com/uber-go/atomic v1.3.2 // indirect
	github.com/uber/jaeger-client-go v2.15.0+incompatible
	github.com/uber/jaeger-lib v1.5.0
	github.com/xdg/scram v0.0.0-20180814205039-7eeb5667e42c // indirect
	github.com/xdg/stringprep v1.0.0 // indirect
	go.mongodb.org/mongo-driver v1.0.0-rc1
	go.uber.org/atomic v1.3.2 // indirect
	golang.org/x/net v0.0.0-20181220203305-927f97764cc3
	google.golang.org/grpc v1.19.0
)

replace github.azc.ext.hp.com/fitstation-hp/lib-fs-core-go => ./../lib-fs-core-go
