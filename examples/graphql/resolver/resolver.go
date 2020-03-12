package resolver

import (
	"context"
	"github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/middleware/jwt"
	"github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/provider"
	"github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/provider/graphql"
	"github.com/sirupsen/logrus"
)

type RootResolver struct {
	provider.AbstractProvider

	graphqlProvider *graphql.GraphQL
}

func NewRootResolver(graphqlProvider *graphql.GraphQL) *RootResolver {
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

func (r *RootResolver) Ping(ctx context.Context) string {
	logrus.WithFields(logrus.Fields{
		"claim_audience":   jwt.GetClaim(ctx, "aud"),
		"claim_expires_at": jwt.GetClaim(ctx, "exp"),
		"claim_issued_at":  jwt.GetClaim(ctx, "iat"),
		"claim_issuer":     jwt.GetClaim(ctx, "iss"),
		"jwt_valid":        jwt.GetToken(ctx).Valid,
	}).Info("Received message")
	return "pong"
}
