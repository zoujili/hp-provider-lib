package proxy

import (
	"context"
	"fmt"
	"github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/provider"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

// Proxy Provider.
// Creates a reverse proxy for communicating with an internal service.
type Proxy struct {
	provider.AbstractRunProvider

	Config       *Config
	ReverseProxy *httputil.ReverseProxy

	srv *http.Server
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

	p.srv = &http.Server{Addr: addr, Handler: mux}
	p.SetRunning(true)

	logEntry.Infof("%s Proxy launched", strings.Title(p.Config.Prefix))
	if err := p.srv.ListenAndServe(); err != http.ErrServerClosed {
		logEntry.WithError(err).Errorf("%s Proxy launch failed", strings.Title(p.Config.Prefix))
		return err
	}

	return nil
}

func (p *Proxy) Close() error {
	if !p.Config.Enabled || p.srv == nil {
		return p.AbstractRunProvider.Close()
	}

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Millisecond)
	if err := p.srv.Shutdown(ctx); err != nil {
		logrus.WithError(err).Errorf("Error while closing %s Proxy server", strings.Title(p.Config.Prefix))
	}

	return p.AbstractRunProvider.Close()
}
