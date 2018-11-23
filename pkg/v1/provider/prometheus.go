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
	v := viper.New()
	v.SetEnvPrefix("PROMETHEUS")
	v.AutomaticEnv()

	v.SetDefault("ENABLED", true)
	enabled := v.GetBool("ENABLED")

	v.SetDefault("PORT", 9090)
	port := v.GetInt("PORT")

	v.SetDefault("ENDPOINT", "/metrics")
	endpoint := v.GetString("ENDPOINT")

	logrus.WithFields(logrus.Fields{
		"enabled":  enabled,
		"port":     port,
		"endpoint": endpoint,
	}).Debug("Prometheus Config Initialized")

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
	return nil
}
