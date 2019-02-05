package stack

import (
	"errors"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"sync"
	"testing"
	"time"
)

func TestStack(t *testing.T) {
	RegisterFailHandlerWithT(t, Fail)
	RunSpecs(t, "Stack test")
}

var _ = Describe("Stack", func() {

	BeforeEach(func() {
		// Make sure that the stack can be used multiple times (once per test).
		closeOnce = sync.Once{}
		runOnce = sync.Once{}
	})

	Describe("Provider flow", func() {
		Context("Default providers", func() {
			st := New()
			p1, p2 := &MockedProvider1{}, &MockedProvider2{}

			It("Should initialize the provider by adding it to the stack", func() {
				Expect(st.providers).To(HaveLen(0))
				st.MustInit(p1)
				st.MustInit(p2)
				Expect(st.providers).To(ConsistOf(p1, p2))
				Expect(p1.initialized).To(BeTrue())
				Expect(p2.initialized).To(BeTrue())
			})
			It("Should close the provider", func() {
				Expect(st.providers).ToNot(BeZero())
				st.MustClose()
				Expect(st.providers).ToNot(BeZero())
				Expect(p1.closed).To(BeTrue())
				Expect(p2.closed).To(BeTrue())
			})
		})
	})

	Describe("Run Provider flow", func() {
		Context("Default providers", func() {
			st := New()
			p1, p2 := &MockedRunProvider1{}, &MockedRunProvider2{}

			It("Should initialize the provider by adding it to the stack", func() {
				Expect(st.providers).To(HaveLen(0))
				st.MustInit(p1)
				st.MustInit(p2)
				Expect(st.providers).To(ConsistOf(p1, p2))
				Expect(p1.initialized).To(BeTrue())
				Expect(p2.initialized).To(BeTrue())
				Expect(p1.IsRunning()).To(BeFalse())
				Expect(p2.IsRunning()).To(BeFalse())
			})
			It("Should run the provider within a certain time limit", func() {
				go st.MustRun()
				time.Sleep(1 * time.Millisecond) // It should be running within a millisecond.
				Expect(p1.IsRunning()).To(BeTrue())
				Expect(p2.IsRunning()).To(BeTrue())
			})
			It("Should close the provider", func() {
				st.MustClose()
				Expect(p1.closed).To(BeTrue())
				Expect(p2.closed).To(BeTrue())
				Expect(p1.IsRunning()).To(BeFalse())
				Expect(p2.IsRunning()).To(BeFalse())
			})
		})
	})

	Describe("Error flow", func() {
		It("Should panic if provider initialization fails", func() {
			st := New()
			p1, p2 := &MockedProvider1{}, &MockedProviderInitErr{}

			st.MustInit(p1)
			Expect(func() { st.MustInit(p2) }).To(Panic())
			Expect(st.providers).To(ConsistOf(p1))
			Expect(p1.initialized).To(BeTrue())
			Expect(p2.initialized).To(BeFalse())
		})
		/* TODO(Peter): For some reason, the panic thrown during st.MustRun isn't being caught, which causes the whole test suite to stop due to this test.
		It("Should stop the application if provider running fails", func() {
			st := New()
			p1, p2, p3 := &MockedProvider1{}, &MockedRunProvider1{}, &MockedRunProviderRunErr{}

			st.MustInit(p1)
			st.MustInit(p2)
			st.MustInit(p3)
			Expect(st.providers).To(ConsistOf(p1, p2, p3))
			Expect(st.MustRun).To(Panic())
			Expect(p2.IsRunning()).To(BeFalse(), "Expected p2 to be closed if running p3 fails")
			Expect(p2.closed).To(BeTrue(), "Expected p2 to be closed if running p3 fails")
			Expect(p3.IsRunning()).To(BeFalse(), "Expected p3 to not be started if running it fails")
			Expect(p3.closed).To(BeTrue(), "Expected p3 to be closed if running it fails")
		})*/
		It("Should panic if provider closing fails", func() {
			st := New()
			p1, p2, p3 := &MockedProvider1{}, &MockedRunProvider1{}, &MockedProviderCloseErr{}

			st.MustInit(p1)
			st.MustInit(p2)
			st.MustInit(p3)
			Expect(st.providers).To(ConsistOf(p1, p2, p3))
			Expect(st.MustClose).To(Panic())
			Expect(p1.closed).To(BeFalse(), "Expected p3 to be closed first, so p1 should not be closed yet")
			Expect(p2.closed).To(BeFalse(), "Expected p3 to be closed first, so p2 should not be closed yet")
			Expect(p3.closed).To(BeFalse())
		})
	})

	Describe("Edge cases", func() {
		Context("Double running or closing", func() {
			st := New()
			p1, p2, p3 := &MockedRunProvider1{}, &MockedRunProvider2{}, &MockedProvider1{}

			It("Should not be possible to run providers twice", func() {
				By("Adding the first provider and running the stack", func() {
					st.MustInit(p1)
					go st.MustRun()
					time.Sleep(1 * time.Millisecond) // It should be running within a millisecond.
					Expect(p1.IsRunning()).To(BeTrue())
				})
				By("Adding a second provider and trying to run the stack again", func() {
					st.MustInit(p2)
					go st.handleInterrupt()
					time.Sleep(1 * time.Millisecond)
					go st.MustRun()
					time.Sleep(1 * time.Millisecond)
					Expect(p2.IsRunning()).To(BeFalse(), "The second provider should not be started if the stack is already running")
				})
			})
			It("Should not be possible to close providers twice", func() {
				By("Closing the stack and thus the 2 providers", func() {
					st.MustClose()
					Expect(p1.closed).To(BeTrue())
					Expect(p2.closed).To(BeTrue())
				})
				By("Adding a third provider and trying to close the stack again", func() {
					st.MustInit(p3)
					st.MustClose()
					Expect(p3.closed).To(BeFalse(), "The third provider should not be closed if the stack is already closed")
				})
			})
		})
	})
})

// Mocked provider with some extra booleans to check its status.
type AbstractMockedProvider struct {
	provider.AbstractProvider
	initialized bool
	closed      bool
}

func (p *AbstractMockedProvider) Init() error {
	p.initialized = true
	return p.AbstractProvider.Init()
}

func (p *AbstractMockedProvider) Close() error {
	p.closed = true
	return p.AbstractProvider.Close()
}

type MockedProvider1 struct {
	AbstractMockedProvider
}

type MockedProvider2 struct {
	AbstractMockedProvider
}

// Mocked provider that throws an error while initializing.
type MockedProviderInitErr struct {
	AbstractMockedProvider
}

func (p *MockedProviderInitErr) Init() error {
	return errors.New("init failed")
}

// Mocked provider that throws an error while closing.
type MockedProviderCloseErr struct {
	AbstractMockedProvider
}

func (p *MockedProviderCloseErr) Close() error {
	return errors.New("close failed")
}

// Mocked run provider with some extra booleans to check its status.
type AbstractMockedRunProvider struct {
	provider.AbstractRunProvider
	initialized bool
	closed      bool
}

func (p *AbstractMockedRunProvider) Init() error {
	p.initialized = true
	return p.AbstractRunProvider.Init()
}

func (p *AbstractMockedRunProvider) Run() error {
	logrus.Info("launching provider")
	p.SetRunning(true)
	return nil
}

func (p *AbstractMockedRunProvider) Close() error {
	p.closed = true
	return p.AbstractRunProvider.Close()
}

type MockedRunProvider1 struct {
	AbstractMockedRunProvider
}

type MockedRunProvider2 struct {
	AbstractMockedRunProvider
}

// Mocked run provider that throws an error during startup.
/*
type MockedRunProviderRunErr struct {
	AbstractMockedRunProvider
}

func (p *MockedRunProviderRunErr) Run() error {
	return errors.New("run failed")
}
*/
