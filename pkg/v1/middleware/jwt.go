package middleware

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

type JWTConfig struct {
	Required bool
	Valid    bool
}

func NewJTWConfigFromEnv() *JWTConfig {
	v := viper.New()
	v.SetEnvPrefix("JWT")
	v.AutomaticEnv()

	v.SetDefault("REQUIRED", true)
	required := v.GetBool("REQUIRED")

	v.SetDefault("VALID", true)
	valid := v.GetBool("VALID")

	logrus.WithFields(logrus.Fields{
		"required": required,
		"valid":    valid,
	}).Debug("JWTConfig Config Initialized")

	return &JWTConfig{
		Required: required,
		Valid:    valid,
	}
}

type JWTMiddleware struct {
	Middleware
	Config *JWTConfig
}

func NewJWT(config *JWTConfig) *JWTMiddleware {
	return &JWTMiddleware{
		Config: config,
	}
}

// JWT middleware ...
// Validates and unpacks a JWT and sets it as "jwt" value in the context
func (m *JWTMiddleware) Handler(next http.Handler) http.Handler {
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

func (m *JWTMiddleware) parseToken(r *http.Request) (*jwt.Token, error) {
	return request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
		// HP JWTs are unsigned.
		return jwt.UnsafeAllowNoneSignatureType, nil
	})
}

func (m *JWTMiddleware) withTokenContext(r *http.Request, token *jwt.Token) *http.Request {
	if token == nil {
		return r
	}
	ctx := context.WithValue(r.Context(), "jwt", token)
	return r.WithContext(ctx)
}

func GetJWTToken(ctx context.Context) *jwt.Token {
	token, ok := ctx.Value("jwt").(*jwt.Token)
	if !ok {
		return nil
	}
	return token
}

func GetJWTClaim(ctx context.Context, name string) string {
	claims, ok := GetJWTToken(ctx).Claims.(jwt.MapClaims)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s", claims[name])
}
