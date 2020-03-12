package prometheus

import (
	"fmt"
	"github.azc.ext.hp.com/hp-business-platform/lib-core-go/pkg/v1/test"
	"github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/provider"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestPrometheus(t *testing.T) {
	RegisterFailHandlerWithT(t, Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "Prometheus provider test", test.LoadCustomReporters("../../test_provider_prometheus.xml"))
}

var _ = Describe("Prometheus provider", func() {
	It("Starts up the prometheus HTTP service", func() {
		logrus.SetLevel(logrus.DebugLevel)
		var p *Prometheus
		By("Creating and initializing the provider", func() {
			p = New(&Config{
				Enabled:  true,
				Port:     defaultPort,
				Endpoint: defaultEndpoint,
			})
			err := p.Init()
			Expect(err).ToNot(HaveOccurred())
		})
		By("Running the provider", func() {
			go func() {
				err := p.Run()
				Expect(err).ToNot(HaveOccurred())
			}()
			err := provider.WaitForRunningProvider(p, 2*time.Second)
			Expect(err).ToNot(HaveOccurred())
			Expect(p.IsRunning()).To(BeTrue())
		})
		By("Adding a gauge", func() {
			testGauge := promauto.NewGauge(prometheus.GaugeOpts{Name: "testing_gauge"})
			testGauge.Add(201)
		})
		By("Testing a request", func() {
			res, err := http.Get(fmt.Sprintf("http://localhost:%d%s", defaultPort, defaultEndpoint))
			Expect(err).ToNot(HaveOccurred())
			Expect(res).ToNot(BeNil())
			Expect(res.StatusCode).To(Equal(200))

			bytes, err := ioutil.ReadAll(res.Body)
			Expect(err).ToNot(HaveOccurred())
			body := string(bytes)
			Expect(body).To(And(
				ContainSubstring("go_goroutines "), // The exact number of routines is unreliable.
				ContainSubstring("testing_gauge 201"),
			))
		})
	})
})
