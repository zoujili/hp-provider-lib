package gateway

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"net/http"
	"strings"
)

// A wrapper around the GRPC Gateway ServeMux.
// It removes the basePath from requests and then forwards them to the ServeMux.
type MuxWrapper struct {
	mux      *runtime.ServeMux
	basePath string
}

func NewMuxWrapper(basePath string, mux *runtime.ServeMux) *MuxWrapper {
	return &MuxWrapper{
		mux:      mux,
		basePath: basePath,
	}
}

func (m *MuxWrapper) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	req.URL.Path = "/" + strings.TrimPrefix(req.URL.Path, m.basePath)
	m.mux.ServeHTTP(res, req)
}
