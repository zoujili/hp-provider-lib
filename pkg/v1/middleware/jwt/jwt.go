package jwt

import (
	"context"
	"fmt"
	"github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"net/http"
)

var contextKey = "jwt"

// Java Web Token (JWT) Middleware for GraphQL provider.
type JWT struct {
	middleware.Middleware

	Config *Config
}

// Creates a JWT Middleware.
func New(config *Config) *JWT {
	contextKey = config.ContextKey
	return &JWT{
		Config: config,
	}
}

// Parses a JWT from the headers and sets the token as value in the context.
// If enabled in configuration, also validates the token with its signature.
func (m *JWT) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := m.parseToken(r)
		if err != nil {
			switch err {
			case request.ErrNoTokenInRequest:
				// Token is missing. Only fail if required.
				if m.Config.Required {
					http.Error(w, "Missing authorization", http.StatusUnauthorized)
					return
				}
			default:
				// Token is invalid. Only fail if it needs to be valid.
				if m.Config.Valid {
					http.Error(w, fmt.Sprintf("Invalid authorization: %v", err), http.StatusUnauthorized)
					return
				}
			}
		}
		next.ServeHTTP(w, m.withTokenContext(r, token))
	})
}

func (m *JWT) parseToken(r *http.Request) (*jwt.Token, error) {
	return request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
		// HP JWTs are unsigned.
		return jwt.UnsafeAllowNoneSignatureType, nil
	})
}

func (m *JWT) withTokenContext(r *http.Request, token *jwt.Token) *http.Request {
	if token == nil {
		return r
	}
	ctx := context.WithValue(r.Context(), contextKey, token)
	return r.WithContext(ctx)
}

// Retrieves the JWT from the Context.
func GetToken(ctx context.Context) *jwt.Token {
	token, ok := ctx.Value(contextKey).(*jwt.Token)
	if !ok {
		return nil
	}
	return token
}

// Retrieve a specific JWT Claim from the Context.
func GetClaim(ctx context.Context, name string) string {
	token := GetToken(ctx)
	if token != nil {
		if claims, ok := GetToken(ctx).Claims.(jwt.MapClaims); ok {
			return fmt.Sprintf("%s", claims[name])
		}
	}
	return ""
}
