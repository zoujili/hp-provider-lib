package provider

import (
	"context"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// MongoDBConfig ...
type MongoDBConfig struct {
	URI             string
	Database        string
	Timeout         time.Duration
	MaxConnsPerHost uint16
}

// NewMongoDBConfigEnv ...
func NewMongoDBConfigEnv() *MongoDBConfig {
	viper.SetDefault("MONGODB_URI", "mongodb://127.0.0.1:27017")
	viper.BindEnv("MONGODB_URI")
	uri := viper.GetString("MONGODB_URI")

	viper.SetDefault("MONGODB_DATABASE", "test")
	viper.BindEnv("MONGODB_DATABASE")
	database := viper.GetString("MONGODB_DATABASE")

	viper.SetDefault("MONGODB_TIMEOUT", 20)
	viper.BindEnv("MONGODB_TIMEOUT")
	timeout := viper.GetDuration("MONGODB_TIMEOUT") * time.Second

	viper.SetDefault("MONGODB_MAX_CONNS_PER_HOST", 16)
	viper.BindEnv("MONGODB_MAX_CONNS_PER_HOST")
	maxConnsPerHost := uint16(viper.GetInt("MONGODB_MAX_CONNS_PER_HOST"))

	logrus.WithFields(logrus.Fields{
		"uri":                uri,
		"database":           database,
		"timeout":            timeout,
		"max_conns_per_host": maxConnsPerHost,
	}).Info("MongoDB Config Initialized")

	return &MongoDBConfig{
		URI:             uri,
		Database:        database,
		Timeout:         timeout,
		MaxConnsPerHost: maxConnsPerHost,
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
	opts := []clientopt.Option{
		clientopt.MaxConnsPerHost(p.Config.MaxConnsPerHost),
		clientopt.MaxIdleConnsPerHost(p.Config.MaxConnsPerHost),
	}

	if p.appProvider != nil {
		opts = append(opts, clientopt.AppName(p.appProvider.Name()))
	}

	client, err := mongo.NewClientWithOptions(p.Config.URI, opts...)
	if err != nil {
		logrus.WithError(err).Error("MongoDB Client Creation Failed")
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.Config.Timeout)
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

	db := client.Database(p.Config.Database)

	p.Client = client
	p.Database = db

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
