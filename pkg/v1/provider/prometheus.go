package provider

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// PrometheusConfig ...
type PrometheusConfig struct {
	Port     int
	Endpoint string
}

// NewPrometheusConfigFromEnv ...
func NewPrometheusConfigFromEnv() *PrometheusConfig {
	viper.SetDefault("PROMETHEUS_PORT", 9090)
	viper.BindEnv("PROMETHEUS_PORT")
	port := viper.GetInt("PROMETHEUS_PORT")

	viper.SetDefault("PROMETHEUS_ENDPOINT", "/metrics")
	viper.BindEnv("PROMETHEUS_ENDPOINT")
	endpoint := viper.GetString("PROMETHEUS_ENDPOINT")

	return &PrometheusConfig{
		Port:     port,
		Endpoint: endpoint,
	}
}

// Prometheus ...
type Prometheus struct {
	Config *PrometheusConfig
}

// NewPrometheus ...
func NewPrometheus(config *PrometheusConfig) *Prometheus {
	return &Prometheus{
		Config: config,
	}
}

// Init ...
func (p *Prometheus) Init() error {
	logrus.Info("Prometheus Provider Initialized")
	return nil
}

// Run ...
func (p *Prometheus) Run() error {
	addr := fmt.Sprintf(":%d", p.Config.Port)

	logger := logrus.WithFields(logrus.Fields{
		"addr":     addr,
		"endpoint": p.Config.Endpoint,
	})

	mux := http.NewServeMux()
	mux.Handle(p.Config.Endpoint, promhttp.Handler())

	logger.Info("Prometheus Provider Launched")
	if err := http.ListenAndServe(addr, mux); err != nil {
		logger.WithError(err).Error("Prometheus Provider Launch Failed")
	}

	return nil
}

// Close ...
func (p *Prometheus) Close() error {
	logrus.Info("Prometheus Provider Closed")
	return nil
}
