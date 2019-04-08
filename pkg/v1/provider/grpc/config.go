package grpc

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	defaultPort = 3000
)

// Configuration for the GRPC Server Provider.
type Config struct {
	Port         int  // Port on which to start the GRPC service.
	LogPayload   bool // Whether or not to enable logging of the payload. Should be disabled on production.
	EnableHealth bool // Whether or not to register the health endpoint.
}

// Initializes the configuration from environment variables.
func NewConfigFromEnv() *Config {
	v := viper.New()
	v.SetEnvPrefix("GRPC")
	v.AutomaticEnv()

	v.SetDefault("PORT", defaultPort)
	port := v.GetInt("PORT")

	v.SetDefault("LOG_PAYLOAD", false)
	logPayload := v.GetBool("LOG_PAYLOAD")

	v.SetDefault("HEALTH_ENABLED", true)
	enableHealth := v.GetBool("HEALTH_ENABLED")

	logrus.WithFields(logrus.Fields{
		"port":         port,
		"logPayload":   logPayload,
		"enableHealth": enableHealth,
	}).Debug("Server Config Initialized")

	return &Config{
		Port:         port,
		LogPayload:   logPayload,
		EnableHealth: enableHealth,
	}
}
