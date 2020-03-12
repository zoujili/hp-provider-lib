package app

import (
	"github.azc.ext.hp.com/hp-business-platform/lib-core-go/pkg/v1/version"
	"github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/provider"
	"github.com/sirupsen/logrus"
	"path"
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

// Appends the given elements to the base path and returns a cleaned URL path.
// The resulting path will always end with a "/".
func (p *App) ParsePath(elem ...string) string {
	res := p.ParseEndpoint(elem...)
	if res != "/" {
		res += "/"
	}
	return res
}

// Appends the given elements to the base path and returns a cleaned URL path.
// The resulting path will not end with a "/", unless that's the only character it contains (root path).
func (p *App) ParseEndpoint(elem ...string) string {
	elem = append([]string{p.Config.BasePath}, elem...)
	return path.Join(elem...)
}
