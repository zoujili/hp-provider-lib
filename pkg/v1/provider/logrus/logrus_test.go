package logrus

import (
	"context"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-core-go/pkg/v1/test"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"testing"
)

func TestLogrus(t *testing.T) {
	RegisterFailHandlerWithT(t, Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "Logrus provider test", test.LoadCustomReporters("../../test_provider_logrus.xml"))
}

var _ = Describe("Logrus provider", func() {
	Context("Standard logger", func() {
		It("Initializes the default logger", func() {
			_ = os.Setenv("LOGRUS_LEVEL", "warn")
			_ = os.Setenv("LOGRUS_FORMATTER", "text_clr")
			_ = os.Setenv("LOGRUS_OUTPUT", "stdout")

			p := New(NewConfigFromEnv())

			By("Checking the default logger before initialization", func() {
				logger := logrus.StandardLogger()
				Expect(logger.GetLevel()).To(Equal(logrus.InfoLevel))
				Expect(logger.Formatter).To(Equal(&logrus.TextFormatter{}))
				Expect(logger.Out).To(Equal(os.Stderr))
			})
			By("Checking the default logger after initialization", func() {
				err := p.Init()
				Expect(err).ToNot(HaveOccurred())

				logger := logrus.StandardLogger()
				Expect(logger.GetLevel()).To(Equal(logrus.WarnLevel))
				Expect(logger.Formatter).To(Equal(&logrus.TextFormatter{ForceColors: true}))
				Expect(logger.Out).To(Equal(os.Stdout))
			})

			_ = os.Unsetenv("LOGRUS_LEVEL")
			_ = os.Unsetenv("LOGRUS_FORMATTER")
			_ = os.Unsetenv("LOGRUS_OUTPUT")
		})
	})
	Context("Separate logger", func() {
		It("Creates a new logger as configured", func() {
			txtFormatter := new(logrus.TextFormatter)
			jsonFormatter := new(logrus.JSONFormatter)

			logger1 := NewLogger(logrus.DebugLevel, txtFormatter, os.Stdout)
			Expect(logger1).ToNot(BeNil())
			Expect(logger1.GetLevel()).To(Equal(logrus.DebugLevel))
			Expect(logger1.Formatter).To(Equal(txtFormatter))

			logger2 := NewLogger(logrus.ErrorLevel, jsonFormatter, os.Stdout)
			Expect(logger2).ToNot(BeNil())
			Expect(logger2.GetLevel()).To(Equal(logrus.ErrorLevel))
			Expect(logger2.Formatter).To(Equal(jsonFormatter))
		})
	})
	Context("Environment parsing", func() {
		It("Parses the environment variables", func() {
			_ = os.Setenv("LOGRUS_LEVEL", "warn")
			_ = os.Setenv("LOGRUS_FORMATTER", "text")
			_ = os.Setenv("LOGRUS_OUTPUT", "stdout")

			lvl, formatter, writer := ParseEnv()
			Expect(lvl).To(Equal(logrus.WarnLevel))
			Expect(formatter).To(Equal(&logrus.TextFormatter{}))
			Expect(writer).To(Equal(os.Stdout))

			_ = os.Unsetenv("LOGRUS_LEVEL")
			_ = os.Unsetenv("LOGRUS_FORMATTER")
			_ = os.Unsetenv("LOGRUS_OUTPUT")
		})
		It("Uses the defaults if no environment variables are set", func() {
			lvl, formatter, writer := ParseEnv()
			Expect(lvl).To(Equal(logrus.InfoLevel))
			Expect(formatter).To(Equal(&logrus.JSONFormatter{}))
			Expect(writer).To(Equal(os.Stderr))
		})
		It("Uses the defaults if erroneous environment variables are set", func() {
			_ = os.Setenv("LOGRUS_LEVEL", "hyper")
			_ = os.Setenv("LOGRUS_FORMATTER", "random")
			_ = os.Setenv("LOGRUS_OUTPUT", "twitter feed")

			lvl, formatter, writer := ParseEnv()
			Expect(lvl).To(Equal(logrus.InfoLevel))
			Expect(formatter).To(Equal(&logrus.JSONFormatter{}))
			Expect(writer).To(Equal(os.Stderr))

			_ = os.Unsetenv("LOGRUS_LEVEL")
			_ = os.Unsetenv("LOGRUS_FORMATTER")
			_ = os.Unsetenv("LOGRUS_OUTPUT")
		})
	})
	Context("Context logger", func() {
		It("Returns a proper log entry if no context logger is set", func() {
			ctx := context.Background()
			entry := GetContextEntry(ctx)
			Expect(entry).ToNot(Equal(ctxlogrus.Extract(ctx)))
		})
		It("Returns a context log entry if a context logger is set", func() {
			ctx := context.Background()
			// Make sure we have a logger instance that isn't equal to the default logger, by using some random writer.
			customLogger := logrus.New()
			customLogger.SetOutput(&strings.Builder{})
			customEntry := logrus.NewEntry(customLogger)
			ctx = ctxlogrus.ToContext(ctx, customEntry)

			entry := GetContextEntry(ctx)
			Expect(entry).To(Equal(customEntry))
			Expect(entry).To(Equal(ctxlogrus.Extract(ctx)))
		})
	})
	Context("Logging tags", func() {
		It("Adds tags to both span and logEntry", func() {
			ctx := context.Background()
			span := mocktracer.New().StartSpan("test").(*mocktracer.MockSpan)
			entry := LogTags(ctx, span, map[string]interface{}{"key1": "value", "key2": 5})
			Expect(entry.Data).To(Equal(logrus.Fields{"key1": "value", "key2": 5}))
			Expect(span.Tags()).To(Equal(map[string]interface{}{"key1": "value", "key2": 5}))
		})
		It("Adds tags to the logEntry if the span is nil", func() {
			ctx := context.Background()
			entry := LogTags(ctx, nil, map[string]interface{}{"key1": "value", "key2": 5})
			Expect(entry.Data).To(Equal(logrus.Fields{"key1": "value", "key2": 5}))
		})
	})
})
