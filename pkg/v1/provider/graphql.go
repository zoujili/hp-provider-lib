package provider

import (
	"fmt"
	"github.com/friendsofgo/graphiql"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

type GraphQLConfig struct {
	Port            int
	GraphiQLEnabled bool
}

func NewGraphQLConfigFromEnv() *GraphQLConfig {
	v := viper.New()
	v.SetEnvPrefix("GRAPHQL")
	v.AutomaticEnv()

	v.SetDefault("PORT", 3030)
	port := v.GetInt("PORT")

	v.SetDefault("GRAPHIQL_ENABLED", false)
	graphiQlEnabled := v.GetBool("GRAPHIQL_ENABLED")

	logrus.WithFields(logrus.Fields{
		"port":            port,
		"graphiqlEnabled": graphiQlEnabled,
	}).Debug("GRPCServer Config Initialized")

	return &GraphQLConfig{
		Port:            port,
		GraphiQLEnabled: graphiQlEnabled,
	}
}

type GraphQL struct {
	Config  *GraphQLConfig
	schema  *graphql.Schema
	running bool
}

func NewGraphQL(config *GraphQLConfig) *GraphQL {
	return &GraphQL{
		Config:  config,
		running: false,
	}
}

func (p *GraphQL) Init() error {
	return nil
}

func (p *GraphQL) Run() error {
	if p.schema == nil {
		return fmt.Errorf("must set GraphQL schema")
	}

	addr := fmt.Sprintf(":%d", p.Config.Port)
	logger := logrus.WithFields(logrus.Fields{
		"addr": addr,
	})

	mux := http.NewServeMux()
	mux.Handle("/", &relay.Handler{Schema: p.schema})
	if p.Config.GraphiQLEnabled {
		graphiqlHandler, err := graphiql.NewGraphiqlHandler("/")
		if err != nil {
			logger.WithError(err).Error("GraphQL Provider Launch Failed")
			return err
		}
		mux.Handle("/graphiql/", http.StripPrefix("/graphiql/", graphiqlHandler))
	}

	logger.Info("GraphQL Provider Launched")
	if err := http.ListenAndServe(addr, mux); err != nil {
		logger.WithError(err).Error("GraphQL Provider Launch Failed")
		return err
	}
	return nil
}

func (p *GraphQL) IsRunning() bool {
	return p.running
}

func (p *GraphQL) Close() error {
	return nil
}

func (p *GraphQL) SetSchema(data string, resolver interface{}) error {
	schema, err := graphql.ParseSchema(data, resolver)
	if err != nil {
		return err
	}
	p.schema = schema
	return nil
}
