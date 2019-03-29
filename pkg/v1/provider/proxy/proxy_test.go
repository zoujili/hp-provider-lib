package proxy

import (
	"fmt"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-core-go/pkg/v1/test"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
	"testing"
	"time"
)

func TestProxy(t *testing.T) {
	RegisterFailHandlerWithT(t, Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "Proxy provider test", test.LoadCustomReporters("../../test_provider_proxy.xml"))
}

var _ = Describe("Proxy provider", func() {
	It("Starts up the proxy HTTP service", func() {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
		var p *Proxy
		By("Creating and initializing the provider", func() {
			p = New(&Config{
				Enabled:   true,
				Debug:     true,
				Port:      defaultPort,
				Endpoint:  defaultEndpoint,
				TargetURL: defaultTargetURL + "/testing",
				Prefix:    "TEST_SERVICE",
			})
			err := p.Init()
			Expect(err).ToNot(HaveOccurred())
		})
		By("Running the provider", func() {
			go func() {
				startService()
				err := p.Run()
				Expect(err).ToNot(HaveOccurred())
			}()
			err := provider.WaitForRunningProvider(p, 2*time.Second)
			Expect(err).ToNot(HaveOccurred())
			Expect(p.IsRunning()).To(BeTrue())
		})
		By("Testing a GET request to default path", func() {
			res, err := http.Get(fmt.Sprintf("http://localhost:%d%s", defaultPort, defaultEndpoint))
			Expect(err).ToNot(HaveOccurred())
			Expect(res).ToNot(BeNil())
			Expect(res.StatusCode).To(Equal(200))

			bytes, err := ioutil.ReadAll(res.Body)
			Expect(err).ToNot(HaveOccurred())
			body := string(bytes)
			Expect(body).To(Equal("GET:"))
		})
		By("Testing a GET request to sub path", func() {
			res, err := http.Get(fmt.Sprintf("http://localhost:%d%s", defaultPort, defaultEndpoint+"sub/somewhere"))
			Expect(err).ToNot(HaveOccurred())
			Expect(res).ToNot(BeNil())
			Expect(res.StatusCode).To(Equal(200))

			bytes, err := ioutil.ReadAll(res.Body)
			Expect(err).ToNot(HaveOccurred())
			body := string(bytes)
			Expect(body).To(Equal("GET:"))
		})
		By("Testing a POST request", func() {
			reqBody := "bla"
			res, err := http.Post(fmt.Sprintf("http://localhost:%d%s", defaultPort, defaultEndpoint), "text/plain", strings.NewReader(reqBody))
			Expect(err).ToNot(HaveOccurred())
			Expect(res).ToNot(BeNil())
			Expect(res.StatusCode).To(Equal(200))

			bytes, err := ioutil.ReadAll(res.Body)
			Expect(err).ToNot(HaveOccurred())
			body := string(bytes)
			Expect(body).To(Equal("POST:" + reqBody))
		})
	})
})

func startService() {
	// Create a service to connect to, which does nothing else but respond with the request body, prefixed with the HTTP method
	mux := http.NewServeMux()
	mux.HandleFunc("/testing/", func(res http.ResponseWriter, req *http.Request) {
		request, err := httputil.DumpRequest(req, true)
		if err != nil {
			logrus.WithError(err).Error("Error while dumping request")
		}
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			logrus.WithError(err).Error("Error while reading request body")
		}
		body = []byte(req.Method + ":" + string(body))
		_, err = res.Write(body)
		if err != nil {
			logrus.WithError(err).Error("Error while writing response body")
		}
		logrus.WithFields(logrus.Fields{
			"request":        string(request),
			"responseStatus": http.StatusOK,
			"responseBody":   string(body),
		}).Debug("Request handled")
	})

	go func() {
		if err := http.ListenAndServe(":8080", mux); err != http.ErrServerClosed {
			panic(err.Error())
		}
	}()
}
