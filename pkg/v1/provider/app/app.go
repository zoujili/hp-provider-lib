package app

import (
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-core-go/pkg/v1/version"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider"
	"github.com/sirupsen/logrus"
)

// Application Provider.
// Maintains basic info for the application (like name and version).
type App struct {
	provider.AbstractProvider

	Config *Config
}

// Creates an App Provider.
func New(config *Config) *App {
	return &App{
		Config: config,
	}
}

// App Provider doesn't need initialization, since version is set during compilation and name via environment variables.
func (p *App) Init() error {
	logrus.WithFields(logrus.Fields{
		"name":    p.Name(),
		"version": p.Version().String(),
	}).Info("App Provider initialized")
	return nil
}

// Returns the Application name.
func (p *App) Name() string {
	return p.Config.Name
}

// Returns the Application version.
func (p *App) Version() version.Version {
	return version.CurrentVersion()
}
