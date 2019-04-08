package proxy

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
)

type loggingTransport struct {
	Config *Config
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Log the reqBytes before sending.
	reqBytes, err := httputil.DumpRequestOut(req, t.Config.Debug)
	if err != nil {
		return nil, err
	}
	logEntry := logrus.WithField("request", string(reqBytes))
	logEntry.Debugf("Performing proxy request to %s", t.Config.Prefix)

	// Perform the actual reqBytes.
	res, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		logEntry.WithError(err).Warnf("Received error from %s", t.Config.Prefix)
		return res, err
	}

	// Log the resBytes.
	resBytes, err := httputil.DumpResponse(res, t.Config.Debug)
	if err != nil {
		return nil, err
	}
	logEntry.WithField("response", string(resBytes)).Infof("Finished request to %s", t.Config.Prefix)

	return res, err
}
