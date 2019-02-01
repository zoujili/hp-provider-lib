package pprof

import (
	"fmt"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider"
	"net/http"
	"net/http/pprof"

	"github.com/sirupsen/logrus"
)

// PProf Provider.
// Provides profiling data to be used by a Google PProf visualization/analysis tool.
type PProf struct {
	provider.AbstractRunProvider

	Config *Config
}

// Creates a PProf Provider.
func New(config *Config) *PProf {
	return &PProf{
		Config: config,
	}
}

// Creates an HTTP service on the configured port and endpoint, where profiling data is published.
func (p *PProf) Run() error {
	if !p.Config.Enabled {
		logrus.Info("PProf Provider not enabled")
		return nil
	}

	addr := fmt.Sprintf(":%d", p.Config.Port)

	logEntry := logrus.WithFields(logrus.Fields{
		"addr":     addr,
		"endpoint": p.Config.Endpoint,
	})

	mux := http.NewServeMux()
	mux.HandleFunc(p.Config.Endpoint+"/", pprof.Index)
	mux.HandleFunc(p.Config.Endpoint+"/cmdline", pprof.Cmdline)
	mux.HandleFunc(p.Config.Endpoint+"/profile", pprof.Profile)
	mux.HandleFunc(p.Config.Endpoint+"/symbol", pprof.Symbol)
	mux.HandleFunc(p.Config.Endpoint+"/trace", pprof.Trace)
	p.SetRunning(true)

	logEntry.Info("PProf Provider launched")
	if err := http.ListenAndServe(addr, mux); err != nil {
		logEntry.WithError(err).Error("PProf Provider launch failed")
	}

	return nil
}
