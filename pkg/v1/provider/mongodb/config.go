package mongodb

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

const (
	defaultURI               = "mongodb://127.0.0.1:27017"
	defaultDatabase          = "test"
	defaultTimeout           = 20
	defaultMaxPoolSize       = 16
	defaultMaxConnIdleTime   = 30
	defaultHeartbeatInterval = 15
)

// Configuration for the MongoDB Provider.
type Config struct {
	URI               string        // URI where to find the MongoDB server (including protocol and port).
	Database          string        // Database name to use.
	Timeout           time.Duration // Maximum duration to wait until the initial connection with the database is established.
	MaxPoolSize       uint16        // Maximum number of connections.
	MaxConnIdleTime   time.Duration // Maximum idle time before a connection is removed from the pool.
	HeartbeatInterval time.Duration // Interval between connection checks.
}

// Initializes the configuration from environment variables.
func NewConfigFromEnv() *Config {
	v := viper.New()
	v.SetEnvPrefix("MONGODB")
	v.AutomaticEnv()

	v.SetDefault("URI", defaultURI)
	uri := v.GetString("URI")

	v.SetDefault("DATABASE", defaultDatabase)
	database := v.GetString("DATABASE")

	v.SetDefault("TIMEOUT", defaultTimeout)
	timeout := v.GetDuration("TIMEOUT") * time.Second

	v.SetDefault("MAX_POOL_SIZE", defaultMaxPoolSize)
	maxPoolSize := uint16(v.GetInt64("MAX_POOL_SIZE"))

	v.SetDefault("MAX_CONN_IDLE_TIME", defaultMaxConnIdleTime)
	maxConnIdleTime := v.GetDuration("MAX_CONN_IDLE_TIME") * time.Second

	v.SetDefault("HEARTBEAT_INTERVAL", defaultHeartbeatInterval)
	heartbeatInterval := v.GetDuration("HEARTBEAT_INTERVAL") * time.Second

	logrus.WithFields(logrus.Fields{
		"uri":                uri,
		"database":           database,
		"timeout":            timeout,
		"max_pool_size":      maxPoolSize,
		"max_conn_idle_time": maxConnIdleTime,
		"heartbeat_interval": heartbeatInterval,
	}).Debug("MongoDB Config initialized")

	return &Config{
		URI:               uri,
		Database:          database,
		Timeout:           timeout,
		MaxPoolSize:       maxPoolSize,
		MaxConnIdleTime:   maxConnIdleTime,
		HeartbeatInterval: heartbeatInterval,
	}
}
