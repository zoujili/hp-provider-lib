package logrus

import (
	"context"
	"fmt"
	"github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/provider"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"strings"
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

func WithError(err error) *logrus.Entry {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	out := logrus.NewEntry(logrus.StandardLogger())
	var printStackTrace = func(pError stackTracer) {
		st := pError.StackTrace()
		depth := 3

		valued := fmt.Sprintf("%+v", st[0:depth])
		valued = strings.Replace(valued, "\t", "", -1)
		stack := strings.Split(valued, "\n")
		out = out.WithField("stack", stack[2:])
	}

	if err, ok := err.(stackTracer); ok {
		printStackTrace(err)
	}

	if err, ok := errors.Cause(err).(stackTracer); ok {
		printStackTrace(err)
	}

	return out.WithError(err)
}

// Retrieves a GRPC context logger.
// Calls "ctxlogrus.Extract(ctx)", but returns a proper logger (instead of no-op) if no context logger is found.
func GetContextEntry(ctx context.Context) *logrus.Entry {
	entry := ctxlogrus.Extract(ctx)
	if entry.Logger.Out == ioutil.Discard {
		return logrus.NewEntry(logrus.StandardLogger())
	}
	return entry
}

// Adds logging tags to both logEntry and span.
// Will return the logEntry as result.
func LogTags(ctx context.Context, span opentracing.Span, tags map[string]interface{}) *logrus.Entry {
	fields := make(logrus.Fields, len(tags))
	for key, val := range tags {
		if span != nil {
			span.SetTag(key, val)
		}
		fields[key] = val
	}
	return GetContextEntry(ctx).WithFields(fields)
}
