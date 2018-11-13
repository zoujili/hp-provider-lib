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

// NewProbesConfigEnv ...
func NewProbesConfigEnv() *ProbesConfig {
	viper.SetDefault("PROBES_ENABLED", true)
	viper.BindEnv("PROBES_ENABLED")
	enabled := viper.GetBool("PROBES_ENABLED")

	viper.SetDefault("PROBES_PORT", 8000)
	viper.BindEnv("PROBES_PORT")
	port := viper.GetInt("PROBES_PORT")

	viper.SetDefault("PROBES_LIVENESS_ENDPOINT", "/healthz")
	viper.BindEnv("PROBES_LIVENESS_ENDPOINT")
	livenessEndpoint := viper.GetString("PROBES_LIVENESS_ENDPOINT")

	viper.SetDefault("PROBES_READINESS_ENDPOINT", "/ready")
	viper.BindEnv("PROBES_READINESS_ENDPOINT")
	readinessEndpoint := viper.GetString("PROBES_READINESS_ENDPOINT")

	logrus.WithFields(logrus.Fields{
		"enabled":            enabled,
		"port":               port,
		"liveness_endpoint":  livenessEndpoint,
		"readiness_endpoint": readinessEndpoint,
	}).Info("Probes Config Initialized")

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
	Config *ProbesConfig

	livenessProbes  []ProbeFunc
	readinessProbes []ProbeFunc
}

// NewProbes ...
func NewProbes(config *ProbesConfig) *Probes {
	return &Probes{
		Config: config,
	}
}

// Init ...
func (p *Probes) Init() error {
	logrus.Info("Probes Provider Initialized")
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

	logger.Info("Probes Provider Launched")
	if err := http.ListenAndServe(addr, mux); err != nil {
		logger.WithError(err).Error("Probes Provider Launch Failed")
		return err
	}

	return nil
}

// Close ...
func (p *Probes) Close() error {
	logrus.Info("Probes Provider Closed")
	return nil
}

func (p *Probes) livenessHandler(w http.ResponseWriter, r *http.Request) {
	for _, probe := range p.livenessProbes {
		if err := probe(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (p *Probes) readinessHandler(w http.ResponseWriter, r *http.Request) {
	for _, probe := range p.readinessProbes {
		if err := probe(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
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
