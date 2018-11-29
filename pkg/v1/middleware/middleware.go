package middleware

import "net/http"

// Middleware for HTTP communication.
// Can be used to wrap functionality around a HTTP handler.
type Middleware interface {
	Handler(next http.Handler) http.Handler
}
