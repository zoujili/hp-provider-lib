package provider

import (
	"fmt"
	"io"
	"strings"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics/prometheus"
)

// JaegerConfig ...
type JaegerConfig struct {
	Enabled bool
	Host    string
	Port    int
}

// NewJaegerConfigFromEnv ...
func NewJaegerConfigFromEnv() *JaegerConfig {
	viper.SetDefault("JAEGER_ENABLED", true)
	viper.BindEnv("JAEGER_ENABLED")
	enabled := viper.GetBool("JAEGER_ENABLED")

	viper.SetDefault("JAEGER_AGENT_HOST", "127.0.0.1")
	viper.BindEnv("JAEGER_AGENT_HOST")
	host := viper.GetString("JAEGER_AGENT_HOST")

	viper.SetDefault("JAEGER_AGENT_PORT", 6831)
	viper.BindEnv("JAEGER_AGENT_PORT")
	port := viper.GetInt("JAEGER_AGENT_PORT")

	logrus.WithFields(logrus.Fields{
		"enabled": enabled,
		"port":    port,
		"host":    host,
	}).Info("Jaeger Config Initialized")

	return &JaegerConfig{
		Enabled: enabled,
		Host:    host,
		Port:    port,
	}
}

// Jaeger ...
type Jaeger struct {
	Config      *JaegerConfig
	appProvider *App

	closer io.Closer
}

// NewJaeger ...
func NewJaeger(config *JaegerConfig, appProvider *App) *Jaeger {
	return &Jaeger{
		Config:      config,
		appProvider: appProvider,
	}
}

// Init ...
func (p *Jaeger) Init() error {
	metrics := prometheus.New()
	logger := &logrusLogger{}

	tracer, closer, err := config.Configuration{
		ServiceName: p.appProvider.Name(),
		Disabled:    !p.Config.Enabled,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: fmt.Sprintf("%s:%d", p.Config.Host, p.Config.Port),
		},
	}.NewTracer(
		config.Metrics(metrics),
		config.Logger(logger),
	)

	if err != nil {
		logrus.WithError(err).Error("Jaeger Provider Initialization Failed")
		return err
	}

	opentracing.SetGlobalTracer(tracer)
	p.closer = closer

	return nil
}

// Close ...
func (p *Jaeger) Close() error {
	err := p.closer.Close()
	if err != nil {
		logrus.WithError(err).Info("Jaeger Provider Close Failed")
		return err
	}

	return nil
}

type logrusLogger struct {
}

func (l *logrusLogger) Error(msg string) {
	logrus.Error(msg)
}

func (l *logrusLogger) Infof(msg string, args ...interface{}) {
	logrus.Infof(strings.Trim(msg, "\n"), args...)
}
