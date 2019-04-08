// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: api.proto

/*
Package server is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package gen

import (
	"io"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/utilities"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

var _ codes.Code
var _ io.Reader
var _ status.Status
var _ = runtime.String
var _ = utilities.NewDoubleArray

var (
	filter_PingService_Ping_0 = &utilities.DoubleArray{Encoding: map[string]int{}, Base: []int(nil), Check: []int(nil)}
)

func request_PingService_Ping_0(ctx context.Context, marshaler runtime.Marshaler, client PingServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq PingRequest
	var metadata runtime.ServerMetadata

	if err := runtime.PopulateQueryParameters(&protoReq, req.URL.Query(), filter_PingService_Ping_0); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.Ping(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

// RegisterPingServiceHandlerFromEndpoint is same as RegisterPingServiceHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterPingServiceHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return RegisterPingServiceHandler(ctx, mux, conn)
}

// RegisterPingServiceHandler registers the http handlers for service PingService to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterPingServiceHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterPingServiceHandlerClient(ctx, mux, NewPingServiceClient(conn))
}

// RegisterPingServiceHandlerClient registers the http handlers for service PingService
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "PingServiceClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "PingServiceClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "PingServiceClient" to call the correct interceptors.
func RegisterPingServiceHandlerClient(ctx context.Context, mux *runtime.ServeMux, client PingServiceClient) error {

	mux.Handle("GET", pattern_PingService_Ping_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_PingService_Ping_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_PingService_Ping_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

var (
	pattern_PingService_Ping_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0}, []string{"ping"}, ""))
)

var (
	forward_PingService_Ping_0 = runtime.ForwardResponseMessage
)
