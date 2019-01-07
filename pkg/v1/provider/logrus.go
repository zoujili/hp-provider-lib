package provider

import (
	"fmt"
	"io"
	"os"
	"path"
	"reflect"
	"runtime"

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
	level, formatter, output := ParseEnv()

	// One-Off logger
	logger := NewLogger(level, formatter, output)
	logger.WithFields(logrus.Fields{
		"level":     level,
		"formatter": reflect.TypeOf(formatter).Elem().String(),
		"output":    reflect.TypeOf(output).Elem().String(),
	}).Debug("Logrus Config Initialized")

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

	return nil
}

// Close ...
func (p *Logrus) Close() error {
	return nil
}

// ParseEnv ...
func ParseEnv() (logrus.Level, logrus.Formatter, io.Writer) {
	v := viper.New()
	v.SetEnvPrefix("LOGRUS")
	v.AutomaticEnv()

	v.SetDefault("LEVEL", "info")
	level := logrus.InfoLevel
	switch v.GetString("LEVEL") {
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

	v.SetDefault("FORMATTER", "json")
	var formatter logrus.Formatter
	switch v.GetString("FORMATTER") {
	case "json":
		formatter = &logrus.JSONFormatter{}
	case "text":
		formatter = &logrus.TextFormatter{}
	case "text_clr":
		formatter = &logrus.TextFormatter{ForceColors: true}
	}

	v.SetDefault("OUTPUT", "stderr")
	var output io.Writer
	switch v.GetString("OUTPUT") {
	case "stderr":
		output = os.Stderr
	case "stdout":
		output = os.Stdout
	}

	return level, formatter, output
}

// NewLogger ...
func NewLogger(level logrus.Level, formatter logrus.Formatter, output io.Writer) *logrus.Logger {
	logrus.AddHook(&SourceHook{})

	logger := logrus.New()
	logger.SetLevel(level)
	logger.Formatter = formatter
	logger.SetOutput(output)
	return logger
}

// Hook to add the source of a message to the log (only for WARN and higher).
type SourceHook struct {
	logrus.Hook
}

func (h SourceHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
	}
}

func (h SourceHook) Fire(entry *logrus.Entry) error {
	if pc, file, line, ok := runtime.Caller(10); ok {
		name := runtime.FuncForPC(pc).Name()
		entry.Data["source"] = fmt.Sprintf("%s:%v:%s", path.Base(file), line, path.Base(name))
	}
	return nil
}
