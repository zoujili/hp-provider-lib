package tenant

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.azc.ext.hp.com/hp-business-platform/hpbp-utils/errors"
	grpcProvider "github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/provider/grpc"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Grpc-Metadata-Tenant in HTTP header
const HeaderTenant = "Tenant"

type tenantInterceptorKey struct{}

func FromTenantInterceptorContext(ctx context.Context) (tenant string, ok bool) {
	tenant, ok = ctx.Value(tenantInterceptorKey{}).(string)
	return
}

func CustomTenantInterceptorOpts() grpcProvider.CustomOpts {
	return grpcProvider.CustomOpts{
		UnaryInterceptor:  []grpc.UnaryServerInterceptor{UnaryServerInterceptor()},
		StreamInterceptor: []grpc.StreamServerInterceptor{StreamServerInterceptor()},
	}
}

// InjectTenant tenant global check
func InjectTenant(h http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(strings.ToLower(r.RequestURI), "/health") {
			h.ServeHTTP(w, r)
			return
		}
		tenantID := r.Header.Get("X-HPBP-Tenant-ID")
		if tenantID == "" {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			appError := errors.NotFoundError(nil).WithMessage("Could not get the X-HPBP-Tenant-ID in request headers")
			w.WriteHeader(appError.Code)
			json.NewEncoder(w).Encode(appError)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), tenantInterceptorKey{}, tenantID))

		h.ServeHTTP(w, r)
	}

	// converts a function to an implementation of an interface
	return http.HandlerFunc(f)
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if tenants := md.Get(HeaderTenant); len(tenants) > 0 {
				ctx = context.WithValue(ctx, tenantInterceptorKey{}, tenants[0])
				return handler(ctx, req)
			}
		}
		return nil, status.Error(codes.Unauthenticated, "Grpc-Metadata-Tenant Missing")
	}
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if tenants := md.Get(HeaderTenant); len(tenants) > 0 {
				newCtx := context.WithValue(ctx, tenantInterceptorKey{}, tenants[0])
				wrappedStream := grpc_middleware.WrapServerStream(ss)
				wrappedStream.WrappedContext = newCtx
				return handler(srv, wrappedStream)
			}
		}
		return status.Error(codes.Unauthenticated, "Grpc-Metadata-Tenant Missing")
	}
}
