package provider

import (
	"io"

	"github.com/sirupsen/logrus"
)

// LogrusConfig ...
type LogrusConfig struct {
	Level     logrus.Level
	Formatter logrus.Formatter
	Output    io.Writer
}

// NewLogrusConfigEnv ...
func NewLogrusConfigEnv() *LogrusConfig {
	return &LogrusConfig{}
}

// Logrus ...
type Logrus struct {
	Config *LogrusConfig
}

// NewLogrus ...
func NewLogrus(config *LogrusConfig) *Logrus {
	return &Logrus{
		Config: config,
	}
}

// Init ...
func (p *Logrus) Init() error {
	logrus.Info("Logrus Provider Initialized")
	return nil
}

// Close ...
func (p *Logrus) Close() error {
	logrus.Info("Logrus Provider Closed")
	return nil
}
