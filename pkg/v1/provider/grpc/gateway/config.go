package gateway

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	defaultPort = 8080
)

// Configuration for the GRPC Gateway Provider.
type Config struct {
	Enabled    bool // Whether or not to enable the gateway.
	Port       int  // Port on which to start the HTTP service.
	LogPayload bool // Whether or not to enable logging of the payload. Should be disabled on production.
}

// Initializes the configuration from environment variables.
func NewConfigFromEnv() *Config {
	v := viper.New()
	v.SetEnvPrefix("GRPC_GATEWAY")
	v.AutomaticEnv()

	v.SetDefault("ENABLED", false)
	enabled := v.GetBool("ENABLED")

	v.SetDefault("PORT", defaultPort)
	port := v.GetInt("PORT")

	v.SetDefault("LOG_PAYLOAD", false)
	logPayload := v.GetBool("LOG_PAYLOAD")

	logrus.WithFields(logrus.Fields{
		"enabled":    enabled,
		"port":       port,
		"logPayload": logPayload,
	}).Debug("Gateway Config Initialized")

	return &Config{
		Enabled:    enabled,
		Port:       port,
		LogPayload: logPayload,
	}
}
