package graphql

import (
	"fmt"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/middleware"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider"
	"github.com/friendsofgo/graphiql"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/sirupsen/logrus"
	"net/http"
)

// GraphQL Provider.
// Uses a GraphQL schema file and root resolver to provide a GraphQL API.
// This schema file must be set before the Run phase.
type GraphQL struct {
	provider.AbstractRunProvider

	Config          *Config
	schema          *graphql.Schema
	middlewareChain []middleware.Middleware
}

// Creates a GraphQL Provider.
func New(config *Config, middlewareChain ...middleware.Middleware) *GraphQL {
	return &GraphQL{
		Config:          config,
		middlewareChain: middlewareChain,
	}
}

// Creates an HTTP service on the configured port (endpoint is always "/"), where queries can be sent.
// If enabled, also adds a GraphiQL handler to provide a GUI.
func (p *GraphQL) Run() error {
	if p.schema == nil {
		return fmt.Errorf("must set GraphQL schema")
	}

	addr := fmt.Sprintf(":%d", p.Config.Port)
	logEntry := logrus.WithField("addr", addr)

	mux := http.NewServeMux()
	mux.Handle("/", p.getHandler())
	if p.Config.GraphiQLEnabled {
		graphiqlHandler, err := graphiql.NewGraphiqlHandler("/")
		if err != nil {
			logEntry.WithError(err).Error("GraphiQL handler could not be started")
			return err
		}
		mux.Handle(p.Config.GraphiQLEndpoint, http.StripPrefix(p.Config.GraphiQLEndpoint, graphiqlHandler))
	}
	p.SetRunning(true)

	logEntry.Info("GraphQL Provider launched")
	if err := http.ListenAndServe(addr, mux); err != nil {
		logEntry.WithError(err).Error("GraphQL Provider launch failed")
		return err
	}
	return nil
}

// Allows setting the GraphQL schema file.
// Parses the schema file.
func (p *GraphQL) SetSchema(data string, resolver interface{}) error {
	schema, err := graphql.ParseSchema(data, resolver, graphql.Logger(new(graphqlLogger)))
	if err != nil {
		return err
	}
	p.schema = schema
	return nil
}

// Separate method that wraps the GraphQL HTTP handler with the configured middlewares.
func (p *GraphQL) getHandler() http.Handler {
	var handler http.Handler
	handler = &relay.Handler{Schema: p.schema}
	for _, mw := range p.middlewareChain {
		handler = mw.Handler(handler)
	}
	return handler
}
