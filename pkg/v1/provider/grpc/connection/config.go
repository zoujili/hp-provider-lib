package connection

import "github.com/spf13/viper"

const (
	defaultHost = "127.0.0.1"
	defaultPort = 3000
)

// Configuration for the GRPC Connection Provider.
type Config struct {
	Host       string // Host on which to connect to the GRPC service.
	Port       int    // Port on which to connect to the GRPC service.
	LogPayload bool   // Whether or not to enable logging of the payload. Should be disabled on production.
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

	return &Config{
		Host:       host,
		Port:       port,
		LogPayload: logPayload,
	}
}
