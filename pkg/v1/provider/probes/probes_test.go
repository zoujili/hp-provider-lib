package probes

import (
	"fmt"
	"github.azc.ext.hp.com/hp-business-platform/lib-core-go/pkg/v1/test"
	"github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/provider"
	"github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/provider/app"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"net/http"
	"testing"
	"time"
)

func TestProbes(t *testing.T) {
	RegisterFailHandlerWithT(t, Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "Probes provider test", test.LoadCustomReporters("../../test_provider_probes.xml"))
}

var _ = Describe("Probes provider", func() {

	port := 8080
	livenessEndpoint := "/health"
	readinessEndpoint := "/ready"

	Context("The probes HTTP service is started", func() {
		var p *Probes

		It("Starts the probes HTTP service", func() {
			logrus.SetLevel(logrus.DebugLevel)
			By("Creating and initializing the provider", func() {
				p = New(&Config{
					Enabled:           true,
					Port:              port,
					LivenessEndpoint:  livenessEndpoint,
					ReadinessEndpoint: readinessEndpoint,
				}, app.New(app.NewConfigFromEnv()))
				err := p.Init()
				Expect(err).NotTo(HaveOccurred())
			})
			By("Running the provider", func() {
				go func() {
					err := p.Run()
					Expect(err).NotTo(HaveOccurred())
				}()
				err := provider.WaitForRunningProvider(p, 2*time.Second)
				Expect(err).NotTo(HaveOccurred())
				Expect(p.IsRunning()).To(BeTrue())
			})
		})

		Context("All probes succeed", func() {
			It("Returns an OK response", func() {
				By("Calling the liveness endpoint", func() {
					res, err := http.Get(fmt.Sprintf("http://localhost:%d%s", port, livenessEndpoint))
					Expect(err).NotTo(HaveOccurred())
					Expect(res).NotTo(BeNil())
					Expect(res.StatusCode).To(Equal(200))
				})
				By("Calling the readiness endpoint", func() {
					res, err := http.Get(fmt.Sprintf("http://localhost:%d%s", port, readinessEndpoint))
					Expect(err).NotTo(HaveOccurred())
					Expect(res).NotTo(BeNil())
					Expect(res.StatusCode).To(Equal(200))
				})
			})
		})

		Context("One of the probes fails", func() {
			It("Returns an error response", func() {
				By("Adding the broken probes", func() {
					p.AddLivenessProbes(brokenProbe)
					p.AddReadinessProbes(brokenProbe)
					Expect(p.livenessProbes).To(HaveLen(1))
					Expect(p.readinessProbes).To(HaveLen(1))
				})
				By("Calling the liveness endpoint", func() {
					res, err := http.Get(fmt.Sprintf("http://localhost:%d%s", port, livenessEndpoint))
					Expect(err).NotTo(HaveOccurred())
					Expect(res).NotTo(BeNil())
					Expect(res.StatusCode).To(Equal(503))
				})
				By("Calling the readiness endpoint", func() {
					res, err := http.Get(fmt.Sprintf("http://localhost:%d%s", port, readinessEndpoint))
					Expect(err).NotTo(HaveOccurred())
					Expect(res).NotTo(BeNil())
					Expect(res.StatusCode).To(Equal(503))
				})
			})
		})
	})
})

func brokenProbe() error {
	return fmt.Errorf("broken")
}
