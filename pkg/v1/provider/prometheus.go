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
	Enabled  bool
	Port     int
	Endpoint string
}

// NewPrometheusConfigFromEnv ...
func NewPrometheusConfigFromEnv() *PrometheusConfig {
	viper.SetDefault("PROMETHEUS_ENABLED", true)
	viper.BindEnv("PROMETHEUS_ENABLED")
	enabled := viper.GetBool("PROMETHEUS_ENABLED")

	viper.SetDefault("PROMETHEUS_PORT", 9090)
	viper.BindEnv("PROMETHEUS_PORT")
	port := viper.GetInt("PROMETHEUS_PORT")

	viper.SetDefault("PROMETHEUS_ENDPOINT", "/metrics")
	viper.BindEnv("PROMETHEUS_ENDPOINT")
	endpoint := viper.GetString("PROMETHEUS_ENDPOINT")

	logrus.WithFields(logrus.Fields{
		"enabled":  enabled,
		"port":     port,
		"endpoint": endpoint,
	}).Info("Prometheus Config Initialized")

	return &PrometheusConfig{
		Enabled:  enabled,
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
	if !p.Config.Enabled {
		logrus.Info("Prometheus Provider Not Enabled")
		return nil
	}

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
		return err
	}

	return nil
}

// Close ...
func (p *Prometheus) Close() error {
	logrus.Info("Prometheus Provider Closed")
	return nil
}
