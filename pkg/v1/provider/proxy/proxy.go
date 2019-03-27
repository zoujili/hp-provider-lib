package proxy

import (
	"bytes"
	"fmt"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// Proxy Provider.
// Creates a reverse proxy for communicating with an internal service.
type Proxy struct {
	provider.AbstractRunProvider

	Config       *Config
	ReverseProxy *httputil.ReverseProxy
}

// Creates a Proxy Provider.
func New(config *Config) *Proxy {
	return &Proxy{
		Config: config,
	}
}

func (p *Proxy) Init() error {
	targetURL, err := url.Parse(p.Config.TargetURL)
	if err != nil {
		logrus.WithError(err).Errorf("%s Proxy Provider initialization failed", p.Config.Prefix)
		return err
	}
	p.ReverseProxy = httputil.NewSingleHostReverseProxy(targetURL)
	return nil
}

func (p *Proxy) Run() error {
	if !p.Config.Enabled {
		logrus.Infof("%s Proxy Provider not enabled", strings.Title(p.Config.Prefix))
		return nil
	}

	addr := fmt.Sprintf(":%d", p.Config.Port)

	logEntry := logrus.WithFields(logrus.Fields{
		"addr":     addr,
		"endpoint": p.Config.Endpoint,
	})

	mux := http.NewServeMux()
	mux.HandleFunc(p.Config.Endpoint, p.Handle)
	p.SetRunning(true)

	logEntry.Infof("%s Proxy Provider launched", strings.Title(p.Config.Prefix))
	if err := http.ListenAndServe(addr, mux); err != http.ErrServerClosed {
		logEntry.WithError(err).Errorf("%s Proxy Provider launch failed", strings.Title(p.Config.Prefix))
		return err
	}
	return nil

}

func (p *Proxy) Handle(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	logEntry := logrus.WithField("request", fmt.Sprintf("%s %s", req.Method, path))
	if p.Config.Debug {
		// Log the request body.
		if body, err := p.readBody(req); err == nil {
			logEntry = logEntry.WithField("requestBody", string(body))
		} else {
			request, _ := httputil.DumpRequest(req, true)
			logrus.WithField("request", request).WithError(err).Warn("Could not parse request body")
		}
	}

	// Wrap response in one that logs the body and status code.
	resLogger := newLoggingResponseWriter(res)

	p.ReverseProxy.ServeHTTP(res, req)

	logEntry = logEntry.WithField("response", resLogger.statusCode)

	if p.Config.Debug && resLogger.body != nil {
		// Log the response body (if set).
		logEntry = logEntry.WithField("responseBody", string(*resLogger.body))
	}
	logEntry.Debugf("Proxy to %s finished", p.Config.Prefix)
}

func (p *Proxy) handleErr(res http.ResponseWriter, err error) {
	res.WriteHeader(http.StatusBadRequest)
	if _, writeErr := res.Write([]byte(err.Error())); writeErr != nil {
		logrus.WithError(writeErr).Errorf("Could not write error body: %s", err)
	}
}

func (p *Proxy) readBody(req *http.Request) (string, error) {
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	// Since the reader for request body has now been closed, we have to replace it with a new reader.
	req.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	return string(buf), nil
}
