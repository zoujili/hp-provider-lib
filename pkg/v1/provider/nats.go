package provider

import (
	"context"
	"errors"
	"github.com/opentracing/opentracing-go"
	"net"
	"time"

	nats "github.com/nats-io/go-nats"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Nats event
type NatsEvent interface {
	Marshal() (data []byte, err error)
}

// NatsConfig ...
type NatsConfig struct {
	Enabled bool
	URI     string
	Timeout time.Duration
}

// NewNatsConfigFromEnv ...
func NewNatsConfigFromEnv() *NatsConfig {
	v := viper.New()
	v.SetEnvPrefix("NATS")
	v.AutomaticEnv()

	v.SetDefault("ENABLED", true)
	enabled := v.GetBool("ENABLED")

	v.SetDefault("URI", "nats://127.0.0.1:4222")
	uri := v.GetString("URI")

	v.SetDefault("TIMEOUT", 20)
	timeout := v.GetDuration("TIMEOUT") * time.Second

	logrus.WithFields(logrus.Fields{
		"uri":     uri,
		"timeout": timeout,
	}).Debug("Nats Config Initialized")

	return &NatsConfig{
		Enabled: enabled,
		URI:     uri,
		Timeout: timeout,
	}
}

// Nats ...
type Nats struct {
	Config         *NatsConfig
	probesProvider *Probes

	Client *nats.Conn
}

// NewNats ...
func NewNats(config *NatsConfig, probesProvider *Probes) *Nats {
	return &Nats{
		Config:         config,
		probesProvider: probesProvider,
	}
}

// Init ...
func (p *Nats) Init() error {
	if !p.Config.Enabled {
		logrus.Info("Nats Provider Not Enabled")
		return nil
	}

	ctx := context.Background()

	cd := &customDialer{
		ctx:             ctx,
		connectTimeout:  p.Config.Timeout,
		connectTimeWait: 1 * time.Second,
	}

	opts := []nats.Option{
		nats.SetCustomDialer(cd),
		nats.ReconnectWait(2 * time.Second),
	}

	client, err := nats.Connect(p.Config.URI, opts...)
	if err != nil {
		logrus.WithError(err).Error("Nats Provider Initialization Failed")
		return err
	}

	if !client.IsConnected() {
		err = errors.New("nats client not connected")
		logrus.WithError(err).Error("Nats Provider Initialization Failed")
		return err
	}

	p.Client = client

	if p.probesProvider != nil {
		p.probesProvider.AddLivenessProbes(p.livenessProbe)
	}

	return nil
}

// Close ...
func (p *Nats) Close() error {
	if !p.Config.Enabled {
		return nil
	}

	p.Client.Close()

	return nil
}

type customDialer struct {
	ctx             context.Context
	nc              *nats.Conn
	connectTimeout  time.Duration
	connectTimeWait time.Duration
}

func (cd *customDialer) Dial(network, address string) (net.Conn, error) {
	ctx, cancel := context.WithTimeout(cd.ctx, cd.connectTimeout)
	defer cancel()

	for {
		logrus.Debug("Attempting to connect to: ", address)
		if ctx.Err() != nil {
			logrus.WithError(ctx.Err()).Error("Nats ctx error")
			return nil, ctx.Err()
		}

		select {
		case <-cd.ctx.Done():
			return nil, cd.ctx.Err()
		default:
			d := &net.Dialer{}
			if conn, err := d.DialContext(ctx, network, address); err == nil {
				logrus.Debug("Connected to NATS successfully")
				return conn, nil
			}

			time.Sleep(cd.connectTimeWait)
		}
	}
}

// Send an event to Nats. Can be called even if Nats is disabled.
func (p *Nats) EmitEvent(ctx context.Context, event NatsEvent, subject string) {
	if p == nil || !p.Config.Enabled {
		return
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "emit_event")
	span.SetTag("subject", subject)
	defer span.Finish()

	data, err := event.Marshal()
	if err != nil {
		logrus.WithError(err).Error("Unable to marshal event")
		return
	}
	if err := p.Client.Publish(subject, data); err != nil {
		logrus.WithError(err).Error("Error while emitting event")
	}
}

func (p *Nats) livenessProbe() error {
	if !p.Client.IsConnected() {
		err := errors.New("nats client not connected")
		logrus.WithError(err).Error("Nats LivenessProbe Failed")
		return err
	}

	logrus.Debug("Nats LivenessProbe Succeeded")

	return nil
}
