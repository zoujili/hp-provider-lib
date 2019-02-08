package provider

import (
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-core-go/pkg/v1/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestUtils(t *testing.T) {
	RegisterFailHandlerWithT(t, Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "Provider utils test", test.LoadCustomReporters("../test_provider_utils.xml"))
}

var _ = Describe("Utils", func() {
	Context("Getting a provider's name", func() {
		It("Returns the name of a provider", func() {
			Expect(Name(&TestProvider1{})).To(Equal("provider.TestProvider1"))
			Expect(Name(&TestProvider2{})).To(Equal("provider.TestProvider2"))
		})
	})
	Context("Waiting for another run provider", func() {
		When("The provider is already running", func() {
			p1 := &TestProvider1{}
			_ = p1.Run()

			It("Doesn't wait at all", func() {
				p2 := &TestProvider2{other: p1}
				err := p2.Run()
				Expect(err).ToNot(HaveOccurred())
				Expect(p2.IsRunning()).To(BeTrue())
			})
		})
		When("The provider is starting up correctly", func() {
			p1 := &TestProvider1{}
			It("Waits for the provider to start", func() {
				p2 := &TestProvider2{other: p1}
				var err error
				go func() {
					err = p2.Run()
				}()
				time.Sleep(10 * time.Millisecond) // Make sure it waits a bit.
				_ = p1.Run()
				time.Sleep(1 * time.Millisecond) // Give the provider time to find out the other provider is running.
				Expect(err).ToNot(HaveOccurred())
				Expect(p2.IsRunning()).To(BeTrue())
			})
		})
		When("The provider never starts", func() {
			p1 := &TestProvider1{}
			It("Throws an error on timeout", func() {
				p2 := &TestProvider2{other: p1}
				err := p2.Run() // This will take a second.
				Expect(err).To(HaveOccurred())
				Expect(p2.IsRunning()).To(BeFalse())
			})
		})
	})
})

type TestProvider1 struct {
	AbstractRunProvider
}

func (p *TestProvider1) Run() error {
	p.SetRunning(true)
	return nil
}

type TestProvider2 struct {
	AbstractRunProvider

	other RunProvider
}

func (p *TestProvider2) Run() error {
	if err := WaitForRunningProvider(p.other, 1); err != nil {
		return err
	}
	p.SetRunning(true)
	return nil
}
