package authz

import (
	"context"
	"net/http"

	"github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/middleware/authz/client"
	"github.com/antihax/optional"
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

// Interceptor authz Interceptor implement the access control
type Interceptor struct {
	authzClient *client.AuthorizationApiService
	u           IUserGetter
	org         IOrganizationGetter
	skipper     Skipper
}

const defaultAuthzServiceAddr = "https://hpbp.hpbp.io/hpbp-authz/v1"

// NewInterceptor .
func NewInterceptor(confFuncs ...ConfigFunc) *Interceptor {
	v := viper.New()
	v.SetEnvPrefix("AUTHZ")
	v.AutomaticEnv()
	v.SetDefault("SERVICE_ADDR", defaultAuthzServiceAddr)
	c := &Config{
		AuthzServiceAddr: v.GetString("SERVICE_ADDR"),
	}
	for _, confFun := range confFuncs {
		confFun(c)
	}
	apiClient := client.NewAPIClient(&client.Configuration{
		BasePath: c.AuthzServiceAddr,
	})
	return &Interceptor{
		authzClient: apiClient.AuthorizationApi,
		u:           c.UserGetter,
		org:         c.OrganizationGetter,
	}
}

func (a *Interceptor) auth(ctx context.Context, userID, externalID, method, path string) (bool, error) {
	if a.skipper != nil && a.skipper(path) {
		return true, nil
	}
	result, resp, err := a.authzClient.AuthorizeRequest(ctx, client.AuthorizationRequest{
		Principal: &client.AuthorizationRequestPrincipal{
			Id:    userID,
			Type_: "user",
		},
		Resource: &client.AuthorizationRequestResource{},
		Request: &client.AuthorizationRequestRequest{
			Method:                 method,
			Path:                   path,
			ExternalOrganizationId: externalID,
		},
	}, &client.AuthorizationApiAuthorizeRequestOpts{
		Rbac: optional.NewString("true"),
		Abac: optional.NewString("false"),
	})
	if resp != nil && resp.StatusCode == http.StatusForbidden {
		return false, nil
	}
	return result.Allow, err
}

// UnaryServerInterceptor authz Interceptor for unary request in grpc Service
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

// StreamServerInterceptor authz Interceptor for stream request in grpc Service, it will not support organization this time
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
