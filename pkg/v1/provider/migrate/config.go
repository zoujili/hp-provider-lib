package migrate

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	defaultPath = "scripts/json"
)

type Config struct {
	Directory string
}

func NewConfigFromEnv() *Config {
	v := viper.New()
	v.SetEnvPrefix("MIGRATIONS")
	v.AutomaticEnv()

	v.SetDefault("DIRECTORY", defaultPath)
	directory := v.GetString("DIRECTORY")

	logrus.WithFields(logrus.Fields{
		"directory": directory,
	}).Debug("Migrate Config Initialized")

	return &Config{
		Directory: directory,
	}
}
