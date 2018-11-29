package resolver

import (
	"context"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/middleware"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider"
	"github.com/sirupsen/logrus"
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

func (r *RootResolver) Ping(ctx context.Context) string {
	logrus.WithFields(logrus.Fields{
		"claim_audience":   middleware.GetJWTClaim(ctx, "aud"),
		"claim_expires_at": middleware.GetJWTClaim(ctx, "exp"),
		"claim_issued_at":  middleware.GetJWTClaim(ctx, "iat"),
		"claim_issuer":     middleware.GetJWTClaim(ctx, "iss"),
		"jwt_valid":        middleware.GetJWTToken(ctx).Valid,
	}).Info("Received message")
	return "pong"
}
