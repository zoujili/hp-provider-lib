package jwt

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Configuration for the JWT Middleware.
type Config struct {
	ContextKey string // Context key on which the JWT is saved.
	Required   bool   // If required is true, an error will be thrown if no JWT is available.
	Valid      bool   // If valid is true, an error will be thrown if the JWT is not valid.
}

// Initializes the configuration from environment variables.
func NewConfigFromEnv() *Config {
	v := viper.New()
	v.SetEnvPrefix("JWT")
	v.AutomaticEnv()

	v.SetDefault("CONTEXT_KEY", contextKey)
	contextKey := v.GetString("CONTEXT_KEY")

	v.SetDefault("REQUIRED", true)
	required := v.GetBool("REQUIRED")

	v.SetDefault("VALID", true)
	valid := v.GetBool("VALID")

	logrus.WithFields(logrus.Fields{
		"required":   required,
		"valid":      valid,
		"contextKey": contextKey,
	}).Debug("JWT Config initialized")

	return &Config{
		ContextKey: contextKey,
		Required:   required,
		Valid:      valid,
	}
}
