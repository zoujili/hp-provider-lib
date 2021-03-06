package mongodb

import (
	"context"
	"github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/provider"
	"github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/provider/app"
	"github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/provider/probes"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// MongoDB Provider.
// Provides a stable, reusable connection to a MongoDB database.
type MongoDB struct {
	provider.AbstractProvider

	Config         *Config
	probesProvider *probes.Probes
	appProvider    *app.App

	Client   *mongo.Client
	Database *mongo.Database
}

// Creates a MongoDB Provider.
// Uses the ProbesProvider to add a liveness probe.
// Uses the AppProvider to send the application name to the MongoDB client.
func New(config *Config, probesProvider *probes.Probes, appProvider *app.App) *MongoDB {
	return &MongoDB{
		Config:         config,
		probesProvider: probesProvider,
		appProvider:    appProvider,
	}
}

// Creates a MongoDB Client, connects to the database server and selects the configured database to be used.
func (p *MongoDB) Init() error {
	opts := options.Client()
	opts.ApplyURI(p.Config.URI)
	opts.SetConnectTimeout(p.Config.Timeout)
	opts.SetHeartbeatInterval(p.Config.HeartbeatInterval)
	opts.SetMaxPoolSize(p.Config.MaxPoolSize)
	opts.SetMaxConnIdleTime(p.Config.MaxConnIdleTime)

	if p.appProvider != nil {
		opts.SetAppName(p.appProvider.Name())
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.Config.Timeout)
	defer cancel()

	logEntry := logrus.WithField("address", p.Config.URI).WithField("time_out", p.Config.Timeout.String())

	logEntry.Debug("Connecting to MongoDB server...")

	// Create Client and connect to MongoDB.
	client, err := mongo.NewClient(opts)
	if err != nil {
		logEntry.WithError(err).Error("MongoDB client creation failed")
		return err
	}

	err = client.Connect(ctx)
	if err != nil {
		logEntry.WithError(err).Error("MongoDB connection failed")
		return err
	}

	// Check connection by pinging.
	err = client.Ping(ctx, nil)
	if err != nil {
		logEntry.WithError(err).Error("MongoDB ping failed")
		return err
	}

	p.Client = client
	p.Database = client.Database(p.Config.Database)

	// Add live probes if possible.
	if p.probesProvider != nil {
		p.probesProvider.AddLivenessProbes(p.livenessProbe)
	}

	return nil
}

// Close to connection with the MongoDB server.
func (p *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), p.Config.Timeout)
	defer cancel()

	err := p.Client.Disconnect(ctx)
	if err != nil {
		logrus.WithError(err).Info("MongoDB disconnecting failed")
		return err
	}

	return p.AbstractProvider.Close()
}

func (p *MongoDB) livenessProbe() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := p.Client.Ping(ctx, nil)
	if err != nil {
		logrus.WithError(err).Error("MongoDB liveness probe failed")
		return err
	}

	logrus.Debug("MongoDB liveness probe succeeded")
	return nil
}
