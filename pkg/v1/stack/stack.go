package stack

import (
	"fitstation-hp/lib-fs-provider-go/pkg/v1/provider"
	"os"
	"os/signal"
	"sync"

	"github.com/sirupsen/logrus"
)

// Stack ...
type Stack struct {
	providers []provider.Provider
}

// New ...
func New() *Stack {
	return &Stack{}
}

// MustInit ...
func (s *Stack) MustInit(provider provider.Provider) {
	if err := provider.Init(); err != nil {
		panic(err)
	}

	s.providers = append(s.providers, provider)
}

var runOnce sync.Once

// MustRun ...
func (s *Stack) MustRun() {
	runOnce.Do(func() {
		for _, p := range s.providers {
			runProvider, ok := p.(provider.RunProvider)
			if ok {
				go func() {
					err := runProvider.Run()
					if err != nil {
						logrus.WithError(err).Panic("Failed to run")
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
			if err := s.providers[i].Close(); err != nil {
				panic(err)
			}
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
