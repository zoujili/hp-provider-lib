package probes

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Configuration for the Probes Provider.
type Config struct {
	Enabled           bool   // Whether or not the the HTTP service should be running.
	Port              int    // Port on which to start the HTTP service.
	LivenessEndpoint  string // Endpoint on which to expose the liveness status.
	ReadinessEndpoint string // Endpoint on which to expose the readiness status.
}

// Initializes the configuration from environment variables.
func NewConfigFromEnv() *Config {
	v := viper.New()
	v.SetEnvPrefix("PROBES")
	v.AutomaticEnv()

	v.SetDefault("ENABLED", true)
	enabled := v.GetBool("ENABLED")

	v.SetDefault("PORT", 8000)
	port := v.GetInt("PORT")

	v.SetDefault("LIVENESS_ENDPOINT", "/healthz")
	livenessEndpoint := v.GetString("LIVENESS_ENDPOINT")

	v.SetDefault("READINESS_ENDPOINT", "/ready")
	readinessEndpoint := v.GetString("READINESS_ENDPOINT")

	logrus.WithFields(logrus.Fields{
		"enabled":            enabled,
		"port":               port,
		"liveness_endpoint":  livenessEndpoint,
		"readiness_endpoint": readinessEndpoint,
	}).Debug("Probes Config initialized")

	return &Config{
		Enabled:           enabled,
		Port:              port,
		LivenessEndpoint:  livenessEndpoint,
		ReadinessEndpoint: readinessEndpoint,
	}
}
