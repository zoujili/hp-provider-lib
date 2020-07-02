package tenant

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func FromTenantInterceptor(ctx context.Context) (string, bool) {
	tenantInterceptor := ctx.Value(tenantInterceptorKey{})
	if tenantInterceptor == nil {
		return "", false
	}
	tenantID, ok := tenantInterceptor.(string)
	return tenantID, ok
}

func TenantUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		tenantID, ok := FromTenantInterceptor(ctx)
		if !ok {
			return status.Error(codes.Unauthenticated, "Grpc-Metadata-Tenant Missing")
		}
		ctx = metadata.AppendToOutgoingContext(ctx, XHpbpTenantID, tenantID)
		ctx = context.WithValue(ctx, tenantInterceptorKey{}, tenantID)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func TenantStreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		tenantID, ok := FromTenantInterceptor(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "Grpc-Metadata-Tenant Missing")
		}
		ctx = metadata.AppendToOutgoingContext(ctx, XHpbpTenantID, tenantID)
		ctx = context.WithValue(ctx, tenantInterceptorKey{}, tenantID)
		return streamer(ctx, desc, cc, method, opts...)
	}
}
