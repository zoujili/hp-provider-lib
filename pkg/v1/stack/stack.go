package stack

import (
	p "github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider"
	"os"
	"os/signal"
	"sync"

	"github.com/sirupsen/logrus"
)

// Stack ...
type Stack struct {
	providers []p.Provider
}

// New ...
func New() *Stack {
	return &Stack{}
}

// MustInit ...
func (s *Stack) MustInit(provider p.Provider) {
	logger := p.NewLogger(p.ParseEnv())

	name := p.Name(provider)
	logger.Debugf("%s Initializing...", name)

	if err := provider.Init(); err != nil {
		logger.WithError(err).Errorf("Error during %s initialization", name)
		panic(err)
	}

	s.providers = append(s.providers, provider)

	logger.Infof("%s Initialized", name)
}

var runOnce sync.Once

// MustRun ...
func (s *Stack) MustRun() {
	runOnce.Do(func() {
		for _, pr := range s.providers {
			runProvider, ok := pr.(p.RunProvider)
			if ok {
				go func() {
					name := p.Name(runProvider)
					logrus.Debugf("%s Launching...", name)

					err := runProvider.Run()
					if err != nil {
						logrus.WithError(err).Panicf("%s Failed to run", name)
					}
				}()
			}
		}
		s.handleInterrupt()
	})
}

var closeOnce sync.Once

// MustClose ...
func (s *Stack) MustClose() {
	closeOnce.Do(func() {
		for i := len(s.providers) - 1; i >= 0; i-- {
			name := p.Name(s.providers[i])
			logrus.Info(name + " Closing...")

			if err := s.providers[i].Close(); err != nil {
				panic(err)
			}

			logrus.Info(name + " Closed")
		}
	})
}

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
