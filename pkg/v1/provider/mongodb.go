package provider

import (
	"context"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// MongoDBConfig ...
type MongoDBConfig struct {
	URI               string
	Database          string
	Timeout           time.Duration
	MaxPoolSize       uint16
	MaxConnIdleTime   time.Duration
	HeartbeatInterval time.Duration
}

// NewMongoDBConfigFromEnv ...
func NewMongoDBConfigFromEnv() *MongoDBConfig {
	v := viper.New()
	v.SetEnvPrefix("MONGODB")
	v.AutomaticEnv()

	v.SetDefault("URI", "mongodb://127.0.0.1:27017")
	uri := v.GetString("URI")

	v.SetDefault("DATABASE", "test")
	database := v.GetString("DATABASE")

	v.SetDefault("TIMEOUT", 20)
	timeout := v.GetDuration("TIMEOUT") * time.Second

	v.SetDefault("MAX_POOL_SIZE", 16)
	maxPoolSize := uint16(v.GetInt64("MAX_POOL_SIZE"))

	v.SetDefault("MAX_CONN_IDLE_TIME", 30)
	maxConnIdleTime := v.GetDuration("MAX_CONN_IDLE_TIME") * time.Second

	v.SetDefault("HEARTBEAT_INTERVAL", 15)
	heartbeatInterval := v.GetDuration("HEARTBEAT_INTERVAL") * time.Second

	logrus.WithFields(logrus.Fields{
		"uri":                uri,
		"database":           database,
		"timeout":            timeout,
		"max_pool_size":      maxPoolSize,
		"max_conn_idle_time": maxConnIdleTime,
		"heartbeat_interval": heartbeatInterval,
	}).Debug("MongoDB Config Initialized")

	return &MongoDBConfig{
		URI:               uri,
		Database:          database,
		Timeout:           timeout,
		MaxPoolSize:       maxPoolSize,
		MaxConnIdleTime:   maxConnIdleTime,
		HeartbeatInterval: heartbeatInterval,
	}
}

// MongoDB ...
type MongoDB struct {
	Config         *MongoDBConfig
	probesProvider *Probes
	appProvider    *App

	Client   *mongo.Client
	Database *mongo.Database
}

// NewMongoDB ...
func NewMongoDB(config *MongoDBConfig, probesProvider *Probes, appProvider *App) *MongoDB {
	return &MongoDB{
		Config:         config,
		probesProvider: probesProvider,
		appProvider:    appProvider,
	}
}

// Init ...
func (p *MongoDB) Init() error {
	opts := options.Client()
	opts.SetConnectTimeout(p.Config.Timeout)
	opts.SetMaxPoolSize(p.Config.MaxPoolSize)
	opts.SetMaxConnIdleTime(p.Config.MaxConnIdleTime)
	opts.SetHeartbeatInterval(p.Config.HeartbeatInterval)

	if p.appProvider != nil {
		opts.SetAppName(p.appProvider.Name())
	}

	client, err := mongo.NewClientWithOptions(p.Config.URI, opts)
	if err != nil {
		logrus.WithError(err).Error("MongoDB Client Creation Failed")
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.Config.Timeout*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		logrus.WithError(err).Error("MongoDB Connect Failed")
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		logrus.WithError(err).Error("MongoDB Ping Failed")
		return err
	}

	p.Client = client
	p.Database = client.Database(p.Config.Database)

	if p.probesProvider != nil {
		p.probesProvider.AddLivenessProbes(p.livenessProbe)
	}

	return nil
}

// Close ...
func (p *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), p.Config.Timeout)
	defer cancel()

	err := p.Client.Disconnect(ctx)
	if err != nil {
		logrus.WithError(err).Info("MongoDB Provider Close Failed")
		return err
	}

	return nil
}

func (p *MongoDB) livenessProbe() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := p.Client.Ping(ctx, nil)
	if err != nil {
		logrus.WithError(err).Error("MongoDB LivenessProbe Failed")
		return err
	}

	logrus.Debug("MongoDB LivenessProbe Succeeded")
	return nil
}
