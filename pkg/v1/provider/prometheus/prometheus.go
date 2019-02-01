package prometheus

import (
	"fmt"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// Prometheus Provider.
// Provides metrics to be used by a Prometheus collector.
type Prometheus struct {
	provider.AbstractRunProvider

	Config *Config
}

// Creates a Prometheus Provider.
func New(config *Config) *Prometheus {
	return &Prometheus{
		Config: config,
	}
}

// Creates an HTTP service on the configured port and endpoint, where metrics are published.
func (p *Prometheus) Run() error {
	if !p.Config.Enabled {
		logrus.Info("Prometheus Provider not enabled")
		return nil
	}

	addr := fmt.Sprintf(":%d", p.Config.Port)

	logEntry := logrus.WithFields(logrus.Fields{
		"addr":     addr,
		"endpoint": p.Config.Endpoint,
	})

	mux := http.NewServeMux()
	mux.Handle(p.Config.Endpoint, promhttp.Handler())
	p.SetRunning(true)

	logEntry.Info("Prometheus Provider launched")
	if err := http.ListenAndServe(addr, mux); err != http.ErrServerClosed {
		logEntry.WithError(err).Error("Prometheus Provider launch failed")
		return err
	}
	return nil
}
