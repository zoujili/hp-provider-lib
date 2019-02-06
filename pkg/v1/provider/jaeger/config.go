package jaeger

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	defaultHost = "127.0.0.1"
	defaultPort = 6831
)

// Config ...
type Config struct {
	Enabled bool   // Whether or not a connection should be made and events should be emitted.
	Host    string // Hostname where to find the Jaeger agent.
	Port    int    // Port where to find the Jaeger agent.
}

// Initializes the configuration from environment variables.
func NewConfigFromEnv() *Config {
	v := viper.New()
	v.SetEnvPrefix("JAEGER")
	v.AutomaticEnv()

	v.SetDefault("ENABLED", true)
	enabled := v.GetBool("ENABLED")

	v.SetDefault("AGENT_HOST", defaultHost)
	host := v.GetString("AGENT_HOST")

	v.SetDefault("AGENT_PORT", defaultPort)
	port := v.GetInt("AGENT_PORT")

	logrus.WithFields(logrus.Fields{
		"enabled": enabled,
		"port":    port,
		"host":    host,
	}).Debug("Jaeger Config Initialized")

	return &Config{
		Enabled: enabled,
		Host:    host,
		Port:    port,
	}
}
