package pprof

import (
	"fmt"
	"github.azc.ext.hp.com/hp-business-platform/lib-core-go/pkg/v1/test"
	"github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/provider"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"net/http"
	"testing"
	"time"
)

func TestPprof(t *testing.T) {
	RegisterFailHandlerWithT(t, Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "Pprof provider test", test.LoadCustomReporters("../../test_provider_pprof.xml"))
}

var _ = Describe("Pprof provider", func() {
	It("Starts the Pprof HTTP service", func() {
		logrus.SetLevel(logrus.DebugLevel)
		var p *PProf
		By("Creating and initializing the provider", func() {
			p = New(&Config{
				Port:     defaultPort,
				Endpoint: defaultEndpoint,
				Enabled:  true,
			})
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
		By("Getting the profiling data", func() {
			res, err := http.Get(fmt.Sprintf("http://localhost:%d%s", defaultPort, defaultEndpoint))
			Expect(err).NotTo(HaveOccurred())
			Expect(res).NotTo(BeNil())
			Expect(res.StatusCode).To(Equal(200))
		})
	})
})
