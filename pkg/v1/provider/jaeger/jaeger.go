package jaeger

import (
	"fmt"
	"github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/provider"
	"github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/provider/app"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	"io"
)

// Jaeger Provider.
// Enables OpenTracing support in the application, which is sent to a Jaeger agent.
type Jaeger struct {
	provider.AbstractProvider

	Config      *Config
	appProvider *app.App

	closer io.Closer
}

// Creates a Jaeger Provider.
// Uses the AppProvider to send the service name to Jaeger.
func New(config *Config, appProvider *app.App) *Jaeger {
	return &Jaeger{
		Config:      config,
		appProvider: appProvider,
	}
}

// Creates the global tracer that reports tracing data to Jaeger.
func (p *Jaeger) Init() error {
	metrics := prometheus.New()

	// Initialize the tracing configuration.
	conf := config.Configuration{
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
	}
	// Use the configuration to create a new tracer.
	tracer, closer, err := conf.NewTracer(
		config.Metrics(metrics),
		config.Logger(&LogrusLogger{}),
		config.ZipkinSharedRPCSpan(true),
	)
	if err != nil {
		logrus.WithError(err).Error("Jaeger Tracer Provider launch failed")
		return err
	}

	opentracing.SetGlobalTracer(tracer)
	p.closer = closer

	return nil
}

// Closes the connection to Jaeger, using the closer created during startup.
func (p *Jaeger) Close() error {
	err := p.closer.Close()
	if err != nil {
		logrus.WithError(err).Info("Error while closing Jaeger tracer")
		return err
	}

	return p.AbstractProvider.Close()
}
