package provider

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// LogrusConfig ...
type LogrusConfig struct {
	Level     logrus.Level
	Formatter logrus.Formatter
	Output    io.Writer
}

// NewLogrusConfigFromEnv ...
func NewLogrusConfigFromEnv() *LogrusConfig {
	viper.SetDefault("LOGRUS_LEVEL", "info")
	viper.BindEnv("LOGRUS_LEVEL")
	var level logrus.Level
	switch viper.GetString("LOGRUS_LEVEL") {
	case "panic":
		level = logrus.PanicLevel
	case "fatal":
		level = logrus.FatalLevel
	case "error":
		level = logrus.ErrorLevel
	case "warn":
		level = logrus.WarnLevel
	case "info":
		level = logrus.InfoLevel
	case "debug":
		level = logrus.DebugLevel
	}

	viper.SetDefault("LOGRUS_FORMATTER", "json")
	viper.BindEnv("LOGRUS_FORMATTER")
	var formatter logrus.Formatter
	switch viper.GetString("LOGRUS_FORMATTER") {
	case "json":
		formatter = &logrus.JSONFormatter{}
	case "text":
		formatter = &logrus.TextFormatter{}
	}

	viper.SetDefault("LOGRUS_OUTPUT", "stderr")
	viper.BindEnv("LOGRUS_OUTPUT")
	var output io.Writer
	switch viper.GetString("LOGRUS_OUTPUT") {
	case "stderr":
		output = os.Stderr
	case "stdout":
		output = os.Stdout
	}

	// One-Off logger
	logger := logrus.New()
	logger.SetLevel(level)
	logger.Formatter = formatter
	logger.SetOutput(output)
	logger.WithFields(logrus.Fields{
		"level":     viper.GetString("LOGRUS_LEVEL"),
		"formatter": viper.GetString("LOGRUS_FORMATTER"),
		"output":    viper.GetString("LOGRUS_OUTPUT"),
	}).Info("Logrus Config Initialized")

	return &LogrusConfig{
		Level:     level,
		Formatter: formatter,
		Output:    output,
	}
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
	logrus.SetLevel(p.Config.Level)

	if p.Config.Formatter != nil {
		logrus.SetFormatter(p.Config.Formatter)
	}

	if p.Config.Output != nil {
		logrus.SetOutput(p.Config.Output)
	}

	logrus.Info("Logrus Provider Initialized")
	return nil
}

// Close ...
func (p *Logrus) Close() error {
	logrus.Info("Logrus Provider Closed")
	return nil
}
