package prometheus

import (
	"context"
	"fmt"
	"github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/provider"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// Prometheus Provider.
// Provides metrics to be used by a Prometheus collector.
type Prometheus struct {
	provider.AbstractRunProvider

	Config *Config

	srv *http.Server
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
	p.srv = &http.Server{Addr: addr, Handler: mux}
	p.SetRunning(true)

	logEntry.Info("Prometheus Provider launched")
	if err := p.srv.ListenAndServe(); err != http.ErrServerClosed {
		logEntry.WithError(err).Error("Prometheus Provider launch failed")
		return err
	}

	return nil
}

func (p *Prometheus) Close() error {
	if !p.Config.Enabled || p.srv == nil {
		return nil
	}

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Millisecond)
	if err := p.srv.Shutdown(ctx); err != nil {
		logrus.WithError(err).Error("Error while closing Prometheus server")
	}

	return p.AbstractRunProvider.Close()
}
