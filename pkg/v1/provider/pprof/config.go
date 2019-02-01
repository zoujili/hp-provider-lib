package pprof

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	defaultPort     = 9999
	defaultEndpoint = "/debug/pprof"
)

// Configuration for the PProf Provider.
type Config struct {
	Enabled  bool   // Whether or not the the HTTP service should be running.
	Port     int    // Port on which to start the HTTP service.
	Endpoint string // Endpoint on which to expose the profiler.
}

// NewConfigFromEnv ...
func NewConfigFromEnv() *Config {
	v := viper.New()
	v.SetEnvPrefix("PPROF")
	v.AutomaticEnv()

	v.SetDefault("ENABLED", true)
	enabled := v.GetBool("ENABLED")

	v.SetDefault("PORT", defaultPort)
	port := v.GetInt("PORT")

	v.SetDefault("ENDPOINT", defaultEndpoint)
	endpoint := v.GetString("ENDPOINT")

	logrus.WithFields(logrus.Fields{
		"enabled":  enabled,
		"port":     port,
		"endpoint": endpoint,
	}).Debug("PProf Config initialized")

	return &Config{
		Enabled:  enabled,
		Port:     port,
		Endpoint: endpoint,
	}
}
