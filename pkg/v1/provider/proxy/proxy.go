package proxy

import (
	"fmt"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider"
	"github.com/sirupsen/logrus"
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
		logrus.WithField("target_url", p.Config.TargetURL).WithError(err).Errorf("%s Proxy Provider initialization failed", p.Config.Prefix)
		return err
	}
	p.ReverseProxy = httputil.NewSingleHostReverseProxy(targetURL)
	p.ReverseProxy.Transport = &loggingTransport{Config: p.Config}
	logrus.WithField("target_url", targetURL).Infof("%s Proxy initialized", strings.Title(p.Config.Prefix))
	return nil
}

func (p *Proxy) Run() error {
	if !p.Config.Enabled {
		logrus.Infof("%s Proxy not enabled", strings.Title(p.Config.Prefix))
		return nil
	}

	addr := fmt.Sprintf(":%d", p.Config.Port)

	logEntry := logrus.WithFields(logrus.Fields{
		"addr":     addr,
		"endpoint": p.Config.Endpoint,
	})

	mux := http.NewServeMux()
	mux.HandleFunc(p.Config.Endpoint, func(res http.ResponseWriter, req *http.Request) {
		req.Host = req.URL.Host
		p.ReverseProxy.ServeHTTP(res, req)
	})
	p.SetRunning(true)

	logEntry.Infof("%s Proxy launched", strings.Title(p.Config.Prefix))
	if err := http.ListenAndServe(addr, mux); err != http.ErrServerClosed {
		logEntry.WithError(err).Errorf("%s Proxy launch failed", strings.Title(p.Config.Prefix))
		return err
	}
	return nil

}
