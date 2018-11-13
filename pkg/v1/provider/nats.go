package provider

import (
	"context"
	"errors"
	"net"
	"time"

	nats "github.com/nats-io/go-nats"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// NatsConfig ...
type NatsConfig struct {
	URI     string
	Timeout time.Duration
}

// NewNatsConfigEnv ...
func NewNatsConfigEnv() *NatsConfig {
	viper.SetDefault("NATS_URI", "nats://127.0.0.1:4222")
	viper.BindEnv("NATS_URI")
	uri := viper.GetString("NATS_URI")

	viper.SetDefault("NATS_TIMEOUT", 20)
	viper.BindEnv("NATS_TIMEOUT")
	timeout := viper.GetDuration("NATS_TIMEOUT") * time.Second

	logrus.WithFields(logrus.Fields{
		"uri":     uri,
		"timeout": timeout,
	}).Info("Nats Config Initialized")

	return &NatsConfig{
		URI:     uri,
		Timeout: timeout,
	}
}

// Nats ...
type Nats struct {
	Config *NatsConfig

	Client *nats.Conn
}

// NewNats ...
func NewNats(config *NatsConfig) *Nats {
	return &Nats{
		Config: config,
	}
}

// Init ...
func (p *Nats) Init() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
		err = errors.New("Nats client not connected")
		logrus.WithError(err).Error("Nats Provider Initialization Failed")
		return err
	}

	p.Client = client

	logrus.Info("Nats Provider Initialized")
	return nil
}

// Close ...
func (p *Nats) Close() error {
	p.Client.Close()

	logrus.Info("Nats Provider Closed")
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
		logrus.Debug("Attempting to connect to", address)
		if ctx.Err() != nil {
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
