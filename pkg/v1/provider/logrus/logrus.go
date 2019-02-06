package logrus

import (
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider"
	"github.com/sirupsen/logrus"
	"io"
)

// Logrus Provider.
type Logrus struct {
	provider.AbstractProvider

	Config *Config
}

// Creates a Logrus Provider.
func New(config *Config) *Logrus {
	return &Logrus{
		Config: config,
	}
}

// Initializes the logrus standard logger.
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

// Creates a new Logrus based logger.
func NewLogger(level logrus.Level, formatter logrus.Formatter, output io.Writer) *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(level)
	logger.SetFormatter(formatter)
	logger.SetOutput(output)
	return logger
}
