package graphql

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	defaultPort             = 3030
	defaultGraphiQLEndpoint = "/graphiql/"
)

// Configuration for the GraphQL Provider.
type Config struct {
	Port             int    // Port on which to start the HTTP service.
	GraphiQLEnabled  bool   // Whether or not to enable the GraphiQL endpoint (GUI for GraphQL messages).
	GraphiQLEndpoint string // Endpoint on which to expose the GraphiQL endpoint.
}

// Initializes the configuration from environment variables.
func NewConfigFromEnv() *Config {
	v := viper.New()
	v.SetEnvPrefix("GRAPHQL")
	v.AutomaticEnv()

	v.SetDefault("PORT", defaultPort)
	port := v.GetInt("PORT")

	v.SetDefault("GRAPHIQL_ENABLED", false)
	graphiQlEnabled := v.GetBool("GRAPHIQL_ENABLED")

	v.SetDefault("GRAPHIQL_ENDPOINT", defaultGraphiQLEndpoint)
	graphiQLEndpoint := v.GetString("GRAPHIQL_ENDPOINT")

	logrus.WithFields(logrus.Fields{
		"port":              port,
		"graphiql_enabled":  graphiQlEnabled,
		"graphiql_endpoint": graphiQLEndpoint,
	}).Debug("GraphQL Config Initialized")

	return &Config{
		Port:             port,
		GraphiQLEnabled:  graphiQlEnabled,
		GraphiQLEndpoint: graphiQLEndpoint,
	}
}
