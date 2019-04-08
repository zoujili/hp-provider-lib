package jaeger

import (
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	"strings"
)

// Jaeger logger that uses Logrus.
type LogrusLogger struct {
	jaeger.Logger
}

// Performs debug level logging.
// It logs debug instead of info, since every trace is logged using this method.
func (l *LogrusLogger) Infof(msg string, args ...interface{}) {
	logrus.Debugf(strings.Trim(msg, "\n"), args...)
}

// Performs error level logging.
func (l *LogrusLogger) Error(msg string) {
	logrus.Error(msg)
}
