package app

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path"
)

const (
	defaultBasePath = "/"
)

// Configuration for the App Provider.
type Config struct {
	Name     string // Application name.
	BasePath string // Base path.
}

// Initializes the configuration from environment variables.
func NewConfigFromEnv() *Config {
	v := viper.New()
	v.SetEnvPrefix("APP")
	v.AutomaticEnv()

	v.SetDefault("NAME", os.Args[0])
	name := v.GetString("NAME")

	v.SetDefault("BASE_PATH", defaultBasePath)
	basePath := path.Clean("/" + v.GetString("BASE_PATH"))

	logrus.WithFields(logrus.Fields{
		"name":      name,
		"base_path": basePath,
	}).Debug("App Config initialized")

	return &Config{
		Name:     name,
		BasePath: basePath,
	}
}
