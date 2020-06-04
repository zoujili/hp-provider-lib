package authorization

import (
	"context"
	"github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/middleware/jwt"
	grpcProvider "github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/provider/grpc"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type authorizationInterceptorKey struct{}

func FromInterceptorContext(ctx context.Context) *jwt.JwtOperator {
	return ctx.Value(authorizationInterceptorKey{}).(*jwt.JwtOperator)
}

func CustomAuthorizationInterceptorOpts() (opt grpcProvider.CustomOpts) {
	return grpcProvider.CustomOpts{
		UnaryInterceptor:  []grpc.UnaryServerInterceptor{UnaryServerInterceptor()},
		StreamInterceptor: []grpc.StreamServerInterceptor{StreamServerInterceptor()},
	}
}
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		operator, err := jwt.NewJwtOperator(ctx)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, err.Error())
		}

		ctx = metadata.AppendToOutgoingContext(ctx, jwt.Authorization, operator.Token())
		ctx = context.WithValue(ctx, authorizationInterceptorKey{}, operator)
		return handler(ctx, req)
	}
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		operator, err := jwt.NewJwtOperator(ctx)
		if err != nil {
			return status.Errorf(codes.Unauthenticated, err.Error())
		}

		newCtx := metadata.AppendToOutgoingContext(ctx, jwt.Authorization, operator.Token())
		newCtx = context.WithValue(ctx, authorizationInterceptorKey{}, operator)
		wrappedStream := grpc_middleware.WrapServerStream(ss)
		wrappedStream.WrappedContext = newCtx
		return handler(srv, wrappedStream)
	}
}
