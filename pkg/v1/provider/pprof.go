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
	v := viper.New()
	v.SetEnvPrefix("PPROF")
	v.AutomaticEnv()

	v.SetDefault("ENABLED", true)
	enabled := v.GetBool("ENABLED")

	v.SetDefault("PORT", 9999)
	port := v.GetInt("PORT")

	v.SetDefault("ENDPOINT", "/debug/pprof")
	endpoint := v.GetString("ENDPOINT")

	logrus.WithFields(logrus.Fields{
		"enabled":  enabled,
		"port":     port,
		"endpoint": endpoint,
	}).Debug("PProf Config Initialized")

	return &PProfConfig{
		Enabled:  enabled,
		Port:     port,
		Endpoint: endpoint,
	}
}

// PProf ...
type PProf struct {
	Config  *PProfConfig
	running bool
}

// NewPProf ...
func NewPProf(config *PProfConfig) *PProf {
	return &PProf{
		Config:  config,
		running: false,
	}
}

// Init ...
func (p *PProf) Init() error {
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
	p.running = true

	logger.Info("PProf Provider Launched")
	if err := http.ListenAndServe(addr, mux); err != nil {
		logger.WithError(err).Error("PProf Provider Launch Failed")
	}

	return nil
}

func (p *PProf) IsRunning() bool {
	return p.running
}

// Close ...
func (p *PProf) Close() error {
	return nil
}
