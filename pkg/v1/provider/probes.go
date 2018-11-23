package provider

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// ProbesConfig ...
type ProbesConfig struct {
	Enabled           bool
	Port              int
	LivenessEndpoint  string
	ReadinessEndpoint string
}

// NewProbesConfigFromEnv ...
func NewProbesConfigFromEnv() *ProbesConfig {
	v := viper.New()
	v.SetEnvPrefix("PROBES")
	v.AutomaticEnv()

	v.SetDefault("ENABLED", true)
	enabled := v.GetBool("ENABLED")

	v.SetDefault("PORT", 8000)
	port := v.GetInt("PORT")

	v.SetDefault("LIVENESS_ENDPOINT", "/healthz")
	livenessEndpoint := v.GetString("LIVENESS_ENDPOINT")

	v.SetDefault("READINESS_ENDPOINT", "/ready")
	readinessEndpoint := v.GetString("READINESS_ENDPOINT")

	logrus.WithFields(logrus.Fields{
		"enabled":            enabled,
		"port":               port,
		"liveness_endpoint":  livenessEndpoint,
		"readiness_endpoint": readinessEndpoint,
	}).Debug("Probes Config Initialized")

	return &ProbesConfig{
		Enabled:           enabled,
		Port:              port,
		LivenessEndpoint:  livenessEndpoint,
		ReadinessEndpoint: readinessEndpoint,
	}
}

// ProbeFunc ...
type ProbeFunc func() error

// Probes ...
type Probes struct {
	Config  *ProbesConfig
	running bool

	livenessProbes  []ProbeFunc
	readinessProbes []ProbeFunc
}

// NewProbes ...
func NewProbes(config *ProbesConfig) *Probes {
	return &Probes{
		Config:  config,
		running: false,
	}
}

// Init ...
func (p *Probes) Init() error {
	return nil
}

// Run ...
func (p *Probes) Run() error {
	if !p.Config.Enabled {
		logrus.Info("Probes Provider Not Enabled")
		return nil
	}

	addr := fmt.Sprintf(":%d", p.Config.Port)

	logger := logrus.WithFields(logrus.Fields{
		"addr":               addr,
		"liveness_endpoint":  p.Config.LivenessEndpoint,
		"readiness_endpoint": p.Config.ReadinessEndpoint,
	})

	mux := http.NewServeMux()
	mux.HandleFunc(p.Config.LivenessEndpoint, p.livenessHandler)
	mux.HandleFunc(p.Config.ReadinessEndpoint, p.readinessHandler)
	p.running = true

	logger.Info("Probes Provider Launched")
	if err := http.ListenAndServe(addr, mux); err != nil {
		logger.WithError(err).Error("Probes Provider Launch Failed")
		return err
	}

	return nil
}

func (p *Probes) IsRunning() bool {
	return p.running
}

// Close ...
func (p *Probes) Close() error {
	return nil
}

func (p *Probes) livenessHandler(w http.ResponseWriter, r *http.Request) {
	for _, probe := range p.livenessProbes {
		if err := probe(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			if _, err := w.Write([]byte(err.Error())); err != nil {
				logrus.WithError(err).Warnf("Error while writing liveness data")
			}
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (p *Probes) readinessHandler(w http.ResponseWriter, r *http.Request) {
	for _, probe := range p.readinessProbes {
		if err := probe(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			if _, err := w.Write([]byte(err.Error())); err != nil {
				logrus.WithError(err).Warnf("Error while writing readiness data")
			}
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

// AddLivenessProbes ...
func (p *Probes) AddLivenessProbes(fn ProbeFunc) {
	p.livenessProbes = append(p.livenessProbes, fn)
}

// AddReadinessProbes ...
func (p *Probes) AddReadinessProbes(fn ProbeFunc) {
	p.readinessProbes = append(p.readinessProbes, fn)
}
