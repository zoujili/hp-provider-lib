package probes

import (
	"fmt"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider"
	"net/http"
	"net/http/httputil"

	"github.com/sirupsen/logrus"
)

// Probe function. Called by the probe handlers to determine the status of the application.
type ProbeFunc func() error

// Probes Provider.
// Provides "probes" to be used by Kubernetes.
// This obviously means you need to configure Kubernetes to work with these probes.
// More info about this can be found here: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes
//
// Liveness probes are used to describe if the application is running properly.
// If these display errors, Kubernetes will try to restart the container.
//
// Readiness probes are used to describe if the application is accepting messages.
// If these display errors, Kubernetes/Istio will remove the Pod from the load-balancers.
type Probes struct {
	provider.AbstractRunProvider

	Config *Config

	livenessProbes  []ProbeFunc
	readinessProbes []ProbeFunc
}

// Creates a Probes Provider.
func New(config *Config) *Probes {
	return &Probes{
		Config: config,
	}
}

// Creates an HTTP service on the configured port and endpoints, where the statuses are published.
func (p *Probes) Run() error {
	if !p.Config.Enabled {
		logrus.Info("Probes Provider not enabled")
		return nil
	}

	addr := fmt.Sprintf(":%d", p.Config.Port)

	logEntry := logrus.WithFields(logrus.Fields{
		"addr":               addr,
		"liveness_endpoint":  p.Config.LivenessEndpoint,
		"readiness_endpoint": p.Config.ReadinessEndpoint,
	})

	mux := http.NewServeMux()
	mux.HandleFunc(p.Config.LivenessEndpoint, p.livenessHandler)
	mux.HandleFunc(p.Config.ReadinessEndpoint, p.readinessHandler)
	p.SetRunning(true)

	logEntry.Info("Probes Provider Launched")
	if err := http.ListenAndServe(addr, mux); err != nil {
		logEntry.WithError(err).Error("Probes Provider launch failed")
		return err
	}

	return nil
}

// This handler will check each liveness probe for errors.
// Only if no errors have occurred, it will respond with an 200 OK. Otherwise there will be a 503.
func (p *Probes) livenessHandler(res http.ResponseWriter, req *http.Request) {
	reqDump, _ := httputil.DumpRequest(req, false)
	logrus.WithField("req", string(reqDump)).Debug("Handling liveness request")
	for _, probe := range p.livenessProbes {
		if err := probe(); err != nil {
			res.WriteHeader(http.StatusServiceUnavailable)
			if _, err := res.Write([]byte(err.Error())); err != nil {
				logrus.WithError(err).Warnf("Error while writing liveness data")
			}
			return
		}
	}
	res.WriteHeader(http.StatusOK)
}

// This handler will check each readiness probe for errors.
// Only if no errors have occurred, it will respond with an 200 OK. Otherwise there will be a 503.
func (p *Probes) readinessHandler(res http.ResponseWriter, req *http.Request) {
	reqDump, _ := httputil.DumpRequest(req, false)
	logrus.WithField("req", string(reqDump)).Debug("Handling readiness request")
	for _, probe := range p.readinessProbes {
		if err := probe(); err != nil {
			res.WriteHeader(http.StatusServiceUnavailable)
			if _, err := res.Write([]byte(err.Error())); err != nil {
				logrus.WithError(err).Warnf("Error while writing readiness data")
			}
			return
		}
	}
	res.WriteHeader(http.StatusOK)
}

// Allows adding extra liveness probes to the handler.
func (p *Probes) AddLivenessProbes(fn ProbeFunc) {
	p.livenessProbes = append(p.livenessProbes, fn)
}

// Allows adding extra readiness probes to the handler.
func (p *Probes) AddReadinessProbes(fn ProbeFunc) {
	p.readinessProbes = append(p.readinessProbes, fn)
}
