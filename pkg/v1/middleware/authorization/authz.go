package authorization

import (
	"context"
	"net/http"

	"github.azc.ext.hp.com/hp-business-platform/lib-hpbp-rest-go/gen/authz_service/client"
	"github.azc.ext.hp.com/hp-business-platform/lib-hpbp-rest-go/gen/authz_service/client/authorization"
	"github.azc.ext.hp.com/hp-business-platform/lib-hpbp-rest-go/gen/authz_service/models"
	"github.com/go-openapi/strfmt"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// IUserGetter for get the user_id from the context.Context
type IUserGetter interface {
	UserID(ctx context.Context) string
}

// IOrganizationGetter for get the organization information from the request
type IOrganizationGetter interface {
	ExternalOrganizationID(ctx context.Context, req interface{}) string
}

// Interceptor authorization Interceptor implement the access control
type Interceptor struct {
	authzClient authorization.ClientService
	u           IUserGetter
	org         IOrganizationGetter
	skipper     Skipper
}

const defaultAuthzServiceHost = "hpbp.hpbp.io"

// NewInterceptor .
func NewInterceptor(confFuncs ...ConfigFunc) *Interceptor {
	v := viper.New()
	v.SetEnvPrefix("AUTHZ")
	v.AutomaticEnv()
	v.SetDefault("SERVICE_HOST", defaultAuthzServiceHost)
	c := &Config{
		AuthzServiceHost: v.GetString("SERVICE_HOST"),
	}
	for _, confFun := range confFuncs {
		confFun(c)
	}
	var authzClient authorization.ClientService
	if c.AuthzClient != nil {
		authzClient = c.AuthzClient
	} else {
		authzClient = client.NewHTTPClientWithConfig(
			strfmt.Default,
			client.DefaultTransportConfig().WithHost(c.AuthzServiceHost),
		).Authorization
	}

	return &Interceptor{
		authzClient: authzClient,
		u:           c.UserGetter,
		org:         c.OrganizationGetter,
		skipper:     c.Skipper,
	}
}

func (a *Interceptor) auth(ctx context.Context, userID, externalID, method, path string) (bool, error) {
	if a.skipper != nil && a.skipper(path) {
		return true, nil
	}
	var user = "user"
	var rbac = "true"
	var abac = "false"
	params := authorization.NewAuthorizeRequestParams().WithBody(&models.AuthorizationRequest{
		Principal: &models.AuthorizationRequestPrincipal{
			ID:   userID,
			Type: &user,
		},
		Resource: &models.AuthorizationRequestResource{},
		Request: &models.AuthorizationRequestRequest{
			Method:                 method,
			Path:                   path,
			ExternalOrganizationID: externalID,
		},
	}).WithAbac(&abac).WithRbac(&rbac).WithContext(ctx)
	resp, err := a.authzClient.AuthorizeRequest(params)

	if err != nil {
		if _, ok := err.(*authorization.AuthorizeRequestForbidden); ok {
			return false, nil
		}
		return false, err
	}
	return resp.Payload.Allow, err
}

// UnaryServerInterceptor authorization Interceptor for unary request in grpc Service
func (a *Interceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		path := info.FullMethod
		method := http.MethodPost
		userID := a.u.UserID(ctx)
		organizationID := a.org.ExternalOrganizationID(ctx, req)
		authResult, err := a.auth(ctx, userID, organizationID, method, path)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		if !authResult {
			return nil, status.Error(codes.PermissionDenied, "Invalid User Permission")
		}
		return handler(ctx, req)
	}
}

// StreamServerInterceptor authorization Interceptor for stream request in grpc Service, it will not support organization this time
func (a *Interceptor) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		path := info.FullMethod
		method := http.MethodPost
		userID := a.u.UserID(context.TODO())
		organizationID := a.org.ExternalOrganizationID(context.TODO(), nil)
		authResult, err := a.auth(context.TODO(), userID, organizationID, method, path)
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
		if !authResult {
			return status.Error(codes.PermissionDenied, "Invalid User Permission")
		}
		return handler(srv, ss)
	}
}
