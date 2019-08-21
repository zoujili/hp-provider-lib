package nats

import (
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

const (
	defaultURI     = "nats://127.0.0.1:4222"
	defaultTimeout = 20
	defaultEncoder = nats.JSON_ENCODER
)

// Configuration for the NATS Provider.
type Config struct {
	Enabled bool          // Whether or not a connection should be made and events should be emitted.
	URI     string        // URI where to find the NATS service.
	Timeout time.Duration // Maximum duration to wait until connection with the NATS service is established.
	Encoder string        // NATS Encoder (see go-nats/enc.go).
}

// Initializes the configuration from environment variables.
func NewConfigFromEnv() *Config {
	v := viper.New()
	v.SetEnvPrefix("NATS")
	v.AutomaticEnv()

	v.SetDefault("ENABLED", true)
	enabled := v.GetBool("ENABLED")

	v.SetDefault("URI", defaultURI)
	uri := v.GetString("URI")

	v.SetDefault("TIMEOUT", defaultTimeout)
	timeout := v.GetDuration("TIMEOUT") * time.Second

	v.SetDefault("ENCODER", defaultEncoder)
	encoder := v.GetString("ENCODER")

	logrus.WithFields(logrus.Fields{
		"enabled": enabled,
		"uri":     uri,
		"timeout": timeout,
		"encoder": encoder,
	}).Debug("NATS Config initialized")

	return &Config{
		Enabled: enabled,
		URI:     uri,
		Timeout: timeout,
		Encoder: encoder,
	}
}
