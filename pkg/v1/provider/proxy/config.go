package proxy

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	defaultPort      = 4040
	defaultEndpoint  = "/"
	defaultTargetURL = "http://localhost:8080"
)

// Configuration for the Proxy Provider.
type Config struct {
	Enabled   bool   // Whether or not the the HTTP service should be running.
	Debug     bool   // Whether or not to log the request and response bodies (URL and status will always be logged).
	Port      int    // Port on which to start the HTTP service.
	Prefix    string // Prefix to use in logs and to get the correct service
	Endpoint  string // Endpoint on which to expose the proxy.
	TargetURL string // URL to where the proxy requests should go.
}

// Initializes the configuration from environment variables.
func NewConfigFromEnv(prefix string) *Config {
	v := viper.New()
	v.SetEnvPrefix(prefix)
	v.AutomaticEnv()

	v.SetDefault("ENABLED", true)
	enabled := v.GetBool("ENABLED")

	v.SetDefault("DEBUG", false)
	debug := v.GetBool("DEBUG")

	v.SetDefault("PORT", defaultPort)
	port := v.GetInt("PORT")

	v.SetDefault("ENDPOINT", defaultEndpoint)
	endpoint := v.GetString("ENDPOINT")

	v.SetDefault("TARGET_URL", defaultTargetURL)
	targetURL := v.GetString("TARGET_URL")

	logrus.WithFields(logrus.Fields{
		"enabled":    enabled,
		"debug":      debug,
		"port":       port,
		"endpoint":   endpoint,
		"target_url": targetURL,
	}).Debugf("%s Proxy Config initialized", prefix)

	return &Config{
		Enabled:   enabled,
		Debug:     debug,
		Port:      port,
		Prefix:    prefix,
		Endpoint:  endpoint,
		TargetURL: targetURL,
	}
}
