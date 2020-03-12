package stack

import (
	p "github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/provider"
	lp "github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/provider/logrus"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
)

var runOnce sync.Once
var closeOnce sync.Once

// Stack manages all providers.
type Stack struct {
	logger    *logrus.Logger
	providers []p.Provider
}

// Creates a new Stack.
func New() *Stack {
	return &Stack{
		// Stack uses its own Logger, since it already logs before the Logrus Provider has been initialized.
		logger:    lp.NewLogger(lp.ParseEnv()),
		providers: make([]p.Provider, 0),
	}
}

// Initializes a given Provider. Panics on failure.
func (s *Stack) MustInit(provider p.Provider) {
	name := p.Name(provider)
	s.logger.Debugf("%s initializing...", name)

	if err := provider.Init(); err != nil {
		s.logger.WithError(err).Panicf("Error during %s initialization", name)
	}

	s.providers = append(s.providers, provider)
	s.logger.Infof("%s initialized", name)
}

// Loops through all Providers and runs all RunProvider instances.
// If any run fails, will automatically close all providers and shut down the application.
func (s *Stack) MustRun() {
	// RunOnce makes sure the Stack isn't started twice.
	runOnce.Do(func() {
		for _, provider := range s.providers {
			if runProvider, ok := provider.(p.RunProvider); ok {
				go s.launch(runProvider)
			}
		}
		s.handleInterrupt()
	})
}

// Loops through all Providers (backwards) and closes all of them. Panics on failure.
func (s *Stack) MustClose() {
	// CloseOnce makes sure the Stack isn't stopped twice.
	closeOnce.Do(func() {
		for i := len(s.providers) - 1; i >= 0; i-- {
			name := p.Name(s.providers[i])
			s.logger.Debugf(" %s closing...", name)

			if err := s.providers[i].Close(); err != nil {
				s.logger.WithError(err).Panicf("%s failed to close", name)
			}

			s.logger.Infof("%s closed", name)
		}
	})
}

// Launches a RunProvider.
// The run method of Provider is a blocking call, thus this method should be called in a separate routine.
func (s *Stack) launch(provider p.RunProvider) {
	name := p.Name(provider)
	s.logger.Debugf("%s launching...", name)

	if err := provider.Run(); err != nil {
		s.logger.WithError(err).Panicf("%s failed to run", name)
	}
}

// Handles any panic inside the MustRun() method by closing all providers.
func (s *Stack) handleInterrupt() {
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan struct{})
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		s.MustClose()
		close(cleanupDone)
	}()
	<-cleanupDone
}
