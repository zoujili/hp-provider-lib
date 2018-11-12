package stack

import (
	"fitstation-hp/lib-fs-provider-go/pkg/v1/provider"

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

// Init ...
func (s *Stack) Init(provider provider.Provider) error {
	if err := provider.Init(); err != nil {
		return err
	}

	s.providers = append(s.providers, provider)

	return nil
}

// MustInit ...
func (s *Stack) MustInit(provider provider.Provider) {
	if err := s.Init(provider); err != nil {
		panic(err)
	}
}

// Run ...
func (s *Stack) Run() error {
	for _, p := range s.providers {
		runProvider, ok := p.(provider.RunProvider)
		if ok {
			go func() {
				err := runProvider.Run()
				if err != nil {
					logrus.WithError(err).Error("Failed to run")
				}
			}()
		}
	}

	return nil
}

// MustRun ...
func (s *Stack) MustRun() {
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
}

// Close ...
func (s *Stack) Close() error {
	for i := len(s.providers) - 1; i >= 0; i-- {
		if err := s.providers[i].Close(); err != nil {
			return err
		}
	}

	return nil
}

// MustClose ...
func (s *Stack) MustClose() {
	if err := s.Close(); err != nil {
		panic(err)
	}
}
