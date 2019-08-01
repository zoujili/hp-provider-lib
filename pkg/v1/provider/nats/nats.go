package nats

import (
	"context"
	"fmt"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider/probes"
	"github.com/opentracing/opentracing-go"
	"time"

	"github.com/nats-io/go-nats"
	"github.com/sirupsen/logrus"
)

// NATS Provider.
// Enables event emitting to NATS.
type Nats struct {
	provider.AbstractProvider

	Config         *Config
	probesProvider *probes.Probes

	Conn *nats.EncodedConn
}

// Creates a NATS Provider.
// Uses the ProbesProvider to add a liveness probe.
func New(config *Config, probesProvider *probes.Probes) *Nats {
	return &Nats{
		Config:         config,
		probesProvider: probesProvider,
	}
}

// Creates an encoded connection with the NATS service.
func (p *Nats) Init() error {
	if !p.Config.Enabled {
		logrus.Info("Nats Provider Not Enabled")
		return nil
	}

	cd := &customDialer{
		ctx:             context.Background(),
		connectTimeout:  p.Config.Timeout,
		connectTimeWait: 1 * time.Second,
	}

	opts := []nats.Option{
		nats.SetCustomDialer(cd),
		nats.ReconnectWait(2 * time.Second),
	}

	logEntry := logrus.WithField("address", p.Config.URI)
	logEntry.Debug("Connecting to NATS service...")

	// Connect to NATS.
	client, err := nats.Connect(p.Config.URI, opts...)
	if err != nil {
		logEntry.WithError(err).Error("NATS client creating failed")
		return err
	}
	if !client.IsConnected() {
		err := fmt.Errorf("NATS client is not connected")
		logEntry.WithError(err).Error("NATS connection failed")
		return err
	}

	// Create encoded connection.
	if p.Conn, err = nats.NewEncodedConn(client, p.Config.Encoder); err != nil {
		logEntry.WithError(err).Error("NATS encoded connection failed")
		return err
	}

	// Add live probes if possible.
	if p.probesProvider != nil {
		p.probesProvider.AddLivenessProbes(p.livenessProbe)
	}

	return nil
}

// Closes the connection with the NATS service.
func (p *Nats) Close() error {
	if !p.Config.Enabled {
		return nil
	}

	p.Conn.Close()
	return nil
}

// Emits an Event to the NATS service. Can be called if the NATS Provider is disabled.
// Has support for tracing the emission.
// Uses the configured encoder to marshal the event.
func (p *Nats) EmitEvent(ctx context.Context, event Event, subject string) {
	if p == nil || !p.Config.Enabled {
		return
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "emit_event")
	span.SetTag("subject", subject)
	defer span.Finish()

	if err := p.Conn.Publish(subject, event); err != nil {
		logrus.WithError(err).Error("Error while emitting event")
	}
}

func (p *Nats) livenessProbe() error {
	if !p.Conn.Conn.IsConnected() {
		err := fmt.Errorf("NATS client is not connected")
		logrus.WithError(err).Error("NATS liveness probe failed")
		return err
	}

	logrus.Debug("NATS liveness probe succeeded")
	return p.AbstractProvider.Close()
}
