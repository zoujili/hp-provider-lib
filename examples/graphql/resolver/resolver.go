package resolver

import (
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider"
)

type RootResolver struct {
	graphqlProvider *provider.GraphQL
}

func NewRootResolver(graphqlProvider *provider.GraphQL) *RootResolver {
	return &RootResolver{
		graphqlProvider: graphqlProvider,
	}
}

func (r *RootResolver) Init() error {
	schema := `
		schema {
			query: Query
		}

		type Query {
			ping(): String!
		}
	`
	return r.graphqlProvider.SetSchema(schema, r)
}

func (r *RootResolver) Close() error {
	return nil
}

func (r *RootResolver) Ping() string {
	return "pong"
}
