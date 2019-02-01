package app

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

// Configuration for the App Provider.
type Config struct {
	Name string // Application name.
}

// Initializes the configuration from environment variables.
func NewConfigFromEnv() *Config {
	v := viper.New()
	v.SetEnvPrefix("APP")
	v.AutomaticEnv()

	v.SetDefault("NAME", os.Args[0])
	name := v.GetString("NAME")

	logrus.WithFields(logrus.Fields{
		"name": name,
	}).Debug("App Config initialized")

	return &Config{
		Name: name,
	}
}
