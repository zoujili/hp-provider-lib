package nats

import (
	"context"
	"github.com/sirupsen/logrus"
	"net"
	"time"
)

// Custom dialer that allows a bit more customization in connecting to the NATS service.
// It is almost a direct copy of the example provided by NATS on https://nats.io/documentation/additional_documentation/custom_dialer.
type customDialer struct {
	ctx             context.Context
	connectTimeout  time.Duration
	connectTimeWait time.Duration
}

func (cd *customDialer) Dial(network, address string) (net.Conn, error) {
	ctx, cancel := context.WithTimeout(cd.ctx, cd.connectTimeout)
	defer cancel()

	logEntry := logrus.WithField("address", address)
	logEntry.Debug("Connecting to NATS...")

	// While loop.
	for {
		if ctx.Err() != nil {
			logEntry.WithError(ctx.Err()).Error("NATS context error")
			return nil, ctx.Err()
		}

		select {
		case <-cd.ctx.Done():
			logEntry.WithError(ctx.Err()).Error("NATS context error")
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
