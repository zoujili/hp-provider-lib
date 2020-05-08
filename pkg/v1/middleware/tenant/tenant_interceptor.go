package tenant

import (
	"context"
	grpcProvider "github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/provider/grpc"
	"github.com/grpc-ecosystem/go-grpc-middleware"
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
