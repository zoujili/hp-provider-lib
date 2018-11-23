// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api.proto

package server // import "github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/examples/ping/server"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type PingRequest struct {
	In                   string   `protobuf:"bytes,1,opt,name=in,proto3" json:"in,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PingRequest) Reset()         { *m = PingRequest{} }
func (m *PingRequest) String() string { return proto.CompactTextString(m) }
func (*PingRequest) ProtoMessage()    {}
func (*PingRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_api_ec82f886005d7522, []int{0}
}
func (m *PingRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PingRequest.Unmarshal(m, b)
}
func (m *PingRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PingRequest.Marshal(b, m, deterministic)
}
func (dst *PingRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PingRequest.Merge(dst, src)
}
func (m *PingRequest) XXX_Size() int {
	return xxx_messageInfo_PingRequest.Size(m)
}
func (m *PingRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PingRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PingRequest proto.InternalMessageInfo

func (m *PingRequest) GetIn() string {
	if m != nil {
		return m.In
	}
	return ""
}

type PingResponse struct {
	Out                  string   `protobuf:"bytes,1,opt,name=out,proto3" json:"out,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PingResponse) Reset()         { *m = PingResponse{} }
func (m *PingResponse) String() string { return proto.CompactTextString(m) }
func (*PingResponse) ProtoMessage()    {}
func (*PingResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_api_ec82f886005d7522, []int{1}
}
func (m *PingResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PingResponse.Unmarshal(m, b)
}
func (m *PingResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PingResponse.Marshal(b, m, deterministic)
}
func (dst *PingResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PingResponse.Merge(dst, src)
}
func (m *PingResponse) XXX_Size() int {
	return xxx_messageInfo_PingResponse.Size(m)
}
func (m *PingResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PingResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PingResponse proto.InternalMessageInfo

func (m *PingResponse) GetOut() string {
	if m != nil {
		return m.Out
	}
	return ""
}

func init() {
	proto.RegisterType((*PingRequest)(nil), "api.PingRequest")
	proto.RegisterType((*PingResponse)(nil), "api.PingResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// PingServiceClient is the client API for PingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PingServiceClient interface {
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
}

type pingServiceClient struct {
	cc *grpc.ClientConn
}

func NewPingServiceClient(cc *grpc.ClientConn) PingServiceClient {
	return &pingServiceClient{cc}
}

func (c *pingServiceClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/api.PingService/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PingServiceServer is the server API for PingService service.
type PingServiceServer interface {
	Ping(context.Context, *PingRequest) (*PingResponse, error)
}

func RegisterPingServiceServer(s *grpc.Server, srv PingServiceServer) {
	s.RegisterService(&_PingService_serviceDesc, srv)
}

func _PingService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PingServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.PingService/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PingServiceServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _PingService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.PingService",
	HandlerType: (*PingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _PingService_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}

func init() { proto.RegisterFile("api.proto", fileDescriptor_api_ec82f886005d7522) }

var fileDescriptor_api_ec82f886005d7522 = []byte{
	// 187 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x8f, 0xc1, 0xca, 0xc2, 0x30,
	0x10, 0x84, 0xff, 0xb6, 0x3f, 0x42, 0xa3, 0x48, 0xcd, 0x49, 0x04, 0xa1, 0xf4, 0xe4, 0x25, 0x0d,
	0x28, 0xe2, 0xc5, 0x93, 0x4f, 0x20, 0xf5, 0xe6, 0x2d, 0xd5, 0x6d, 0x5d, 0xa8, 0xc9, 0x9a, 0xa4,
	0xc5, 0xc7, 0x97, 0x5a, 0x0f, 0xbd, 0xed, 0xce, 0x0c, 0xdf, 0x30, 0x2c, 0x56, 0x84, 0x39, 0x59,
	0xe3, 0x0d, 0x8f, 0x14, 0x61, 0xb6, 0x66, 0xd3, 0x33, 0xea, 0xba, 0x80, 0x57, 0x0b, 0xce, 0xf3,
	0x39, 0x0b, 0x51, 0x2f, 0x83, 0x34, 0xd8, 0xc4, 0x45, 0x88, 0x3a, 0x4b, 0xd9, 0x6c, 0xb0, 0x1d,
	0x19, 0xed, 0x80, 0x27, 0x2c, 0x32, 0xad, 0xff, 0x05, 0xfa, 0x73, 0x7b, 0x1c, 0x00, 0x17, 0xb0,
	0x1d, 0xde, 0x80, 0x0b, 0xf6, 0xdf, 0xbf, 0x3c, 0xc9, 0xfb, 0xa2, 0x11, 0x7a, 0xb5, 0x18, 0x29,
	0x03, 0x2d, 0xfb, 0x3b, 0x1d, 0xae, 0xfb, 0x0a, 0xbd, 0xf3, 0xca, 0xa3, 0xd1, 0xe2, 0x41, 0xb2,
	0xc1, 0x52, 0x54, 0x4e, 0x90, 0x35, 0x1d, 0xde, 0xc1, 0x8a, 0xda, 0x48, 0x78, 0xab, 0x27, 0x35,
	0xe0, 0x24, 0xa1, 0xae, 0xa5, 0x03, 0xdb, 0x81, 0x2d, 0x27, 0xdf, 0x0d, 0xbb, 0x4f, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x8d, 0xff, 0xe7, 0x6d, 0xd0, 0x00, 0x00, 0x00,
}
