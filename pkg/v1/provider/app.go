package provider

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"fitstation-hp/lib-fs-core-go/pkg/v1/version"
)

// AppConfig ...
type AppConfig struct {
	Name string
}

// NewAppConfigEnv ...
func NewAppConfigEnv() *AppConfig {
	viper.SetDefault("APP_NAME", os.Args[0])
	viper.BindEnv("APP_NAME")
	name := viper.GetString("APP_NAME")

	logrus.WithFields(logrus.Fields{
		"name": name,
	}).Info("App Config Initialized")

	return &AppConfig{
		Name: name,
	}
}

// App ...
type App struct {
	Config *AppConfig
}

// NewApp ...
func NewApp(config *AppConfig) *App {
	return &App{
		Config: config,
	}
}

// Init ...
func (p *App) Init() error {
	logrus.WithFields(logrus.Fields{
		"name":    p.Name(),
		"version": p.Version().String(),
	}).Info("App Provider Initialized")
	return nil
}

// Close ...
func (p *App) Close() error {
	logrus.Info("App Provider Closed")
	return nil
}

// Name ...
func (p *App) Name() string {
	return p.Config.Name
}

// Version ...
func (p *App) Version() version.Version {
	return version.CurrentVersion()
}
