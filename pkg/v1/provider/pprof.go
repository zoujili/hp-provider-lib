package provider

import (
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// PProfConfig ...
type PProfConfig struct {
	Enabled  bool
	Port     int
	Endpoint string
}

// NewPProfConfigFromEnv ...
func NewPProfConfigFromEnv() *PProfConfig {
	viper.SetDefault("PPROF_ENABLED", true)
	viper.BindEnv("PPROF_ENABLED")
	enabled := viper.GetBool("PPROF_ENABLED")

	viper.SetDefault("PPROF_PORT", 9999)
	viper.BindEnv("PPROF_PORT")
	port := viper.GetInt("PPROF_PORT")

	viper.SetDefault("PPROF_ENDPOINT", "/debug/pprof")
	viper.BindEnv("PPROF_ENDPOINT")
	endpoint := viper.GetString("PPROF_ENDPOINT")

	logrus.WithFields(logrus.Fields{
		"enabled":  enabled,
		"port":     port,
		"endpoint": endpoint,
	}).Info("PProf Config Initialized")

	return &PProfConfig{
		Enabled:  enabled,
		Port:     port,
		Endpoint: endpoint,
	}
}

// PProf ...
type PProf struct {
	Config *PProfConfig
}

// NewPProf ...
func NewPProf(config *PProfConfig) *PProf {
	return &PProf{
		Config: config,
	}
}

// Init ...
func (p *PProf) Init() error {
	logrus.Info("PProf Provider Initialized")
	return nil
}

// Run ...
func (p *PProf) Run() error {
	if !p.Config.Enabled {
		logrus.Info("PProf Provider Not Enabled")
		return nil
	}

	addr := fmt.Sprintf(":%d", p.Config.Port)

	logger := logrus.WithFields(logrus.Fields{
		"addr":     addr,
		"endpoint": p.Config.Endpoint,
	})

	mux := http.NewServeMux()
	mux.HandleFunc(p.Config.Endpoint+"/", pprof.Index)
	mux.HandleFunc(p.Config.Endpoint+"/cmdline", pprof.Cmdline)
	mux.HandleFunc(p.Config.Endpoint+"/profile", pprof.Profile)
	mux.HandleFunc(p.Config.Endpoint+"/symbol", pprof.Symbol)
	mux.HandleFunc(p.Config.Endpoint+"/trace", pprof.Trace)

	logger.Info("PProf Provider Launched")
	if err := http.ListenAndServe(addr, mux); err != nil {
		logger.WithError(err).Error("PProf Provider Launch Failed")
	}

	return nil
}

// Close ...
func (p *PProf) Close() error {
	logrus.Info("PProf Provider Closed")
	return nil
}
