package connection

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	defaultHost = "127.0.0.1"
	defaultPort = 3000
)

// Configuration for the GRPC Connection Provider.
type Config struct {
	Prefix       string // GRPC Connection prefix, used for environment variables and in some bits of logging (like health).
	Host         string // Host on which to connect to the GRPC service.
	Port         int    // Port on which to connect to the GRPC service.
	LogPayload   bool   // Whether or not to enable logging of the payload. Should be disabled on production.
	EnableHealth bool   // Whether or not to enable checking the health of the connection.
}

// Initializes the configuration from environment variables.
func NewConfigFromEnv(prefix string) *Config {
	fsv := viper.New()
	fsv.SetEnvPrefix("FIT_STATION")
	fsv.AutomaticEnv()

	v := viper.New()
	v.SetEnvPrefix(prefix)
	v.AutomaticEnv()

	hostDefault := defaultHost
	if host := fsv.GetString("HOST"); host != "" {
		hostDefault = host
	}
	v.SetDefault("HOST", hostDefault)
	host := v.GetString("HOST")

	v.SetDefault("PORT", defaultPort)
	port := v.GetInt("PORT")

	v.SetDefault("LOG_PAYLOAD", false)
	logPayload := v.GetBool("LOG_PAYLOAD")

	v.SetDefault("HEALTH_ENABLED", true)
	enableHealth := v.GetBool("HEALTH_ENABLED")

	logrus.WithFields(logrus.Fields{
		"prefix":       prefix,
		"host":         host,
		"port":         port,
		"logPayload":   logPayload,
		"enableHealth": enableHealth,
	}).Debug("GRPC Connection Config initialized")

	return &Config{
		Prefix:       prefix,
		Host:         host,
		Port:         port,
		LogPayload:   logPayload,
		EnableHealth: enableHealth,
	}
}
