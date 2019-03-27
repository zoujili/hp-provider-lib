package proxy

import "net/http"

// Middleware response writer that allows the proxy to log responses.
type loggingResponseWriter struct {
	http.ResponseWriter

	statusCode int
	body       *[]byte
}

func newLoggingResponseWriter(res http.ResponseWriter) loggingResponseWriter {
	return loggingResponseWriter{res, 0, nil}
}

func (w loggingResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.WriteHeader(statusCode)
}

func (w loggingResponseWriter) Write(buf []byte) (int, error) {
	w.body = &buf
	return w.Write(buf)
}
