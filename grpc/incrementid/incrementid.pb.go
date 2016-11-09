// Code generated by protoc-gen-go.
// source: incrementid.proto
// DO NOT EDIT!

/*
Package incrementid is a generated protocol buffer package.

It is generated from these files:
	incrementid.proto

It has these top-level messages:
	IncrIdNameRequest
	IncrIdReply
	IncrBoolReply
	IncrIdNameValueRequest
	NoneReply
*/
package incrementid

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

type IncrIdNameRequest struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *IncrIdNameRequest) Reset()                    { *m = IncrIdNameRequest{} }
func (m *IncrIdNameRequest) String() string            { return proto.CompactTextString(m) }
func (*IncrIdNameRequest) ProtoMessage()               {}
func (*IncrIdNameRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type IncrIdReply struct {
	Id uint64 `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
}

func (m *IncrIdReply) Reset()                    { *m = IncrIdReply{} }
func (m *IncrIdReply) String() string            { return proto.CompactTextString(m) }
func (*IncrIdReply) ProtoMessage()               {}
func (*IncrIdReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type IncrBoolReply struct {
	Ret bool `protobuf:"varint,1,opt,name=ret" json:"ret,omitempty"`
}

func (m *IncrBoolReply) Reset()                    { *m = IncrBoolReply{} }
func (m *IncrBoolReply) String() string            { return proto.CompactTextString(m) }
func (*IncrBoolReply) ProtoMessage()               {}
func (*IncrBoolReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type IncrIdNameValueRequest struct {
	Name  string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Value uint64 `protobuf:"varint,2,opt,name=value" json:"value,omitempty"`
}

func (m *IncrIdNameValueRequest) Reset()                    { *m = IncrIdNameValueRequest{} }
func (m *IncrIdNameValueRequest) String() string            { return proto.CompactTextString(m) }
func (*IncrIdNameValueRequest) ProtoMessage()               {}
func (*IncrIdNameValueRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type NoneReply struct {
}

func (m *NoneReply) Reset()                    { *m = NoneReply{} }
func (m *NoneReply) String() string            { return proto.CompactTextString(m) }
func (*NoneReply) ProtoMessage()               {}
func (*NoneReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func init() {
	proto.RegisterType((*IncrIdNameRequest)(nil), "incrementid.IncrIdNameRequest")
	proto.RegisterType((*IncrIdReply)(nil), "incrementid.IncrIdReply")
	proto.RegisterType((*IncrBoolReply)(nil), "incrementid.IncrBoolReply")
	proto.RegisterType((*IncrIdNameValueRequest)(nil), "incrementid.IncrIdNameValueRequest")
	proto.RegisterType((*NoneReply)(nil), "incrementid.NoneReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for IncrementId service

type IncrementIdClient interface {
	GetIncrId(ctx context.Context, in *IncrIdNameRequest, opts ...grpc.CallOption) (*IncrIdReply, error)
	GetIncrIdByAmount(ctx context.Context, in *IncrIdNameValueRequest, opts ...grpc.CallOption) (*IncrIdReply, error)
	CheckIncrKeyExist(ctx context.Context, in *IncrIdNameRequest, opts ...grpc.CallOption) (*IncrBoolReply, error)
	CreateIncrKey(ctx context.Context, in *IncrIdNameValueRequest, opts ...grpc.CallOption) (*NoneReply, error)
}

type incrementIdClient struct {
	cc *grpc.ClientConn
}

func NewIncrementIdClient(cc *grpc.ClientConn) IncrementIdClient {
	return &incrementIdClient{cc}
}

func (c *incrementIdClient) GetIncrId(ctx context.Context, in *IncrIdNameRequest, opts ...grpc.CallOption) (*IncrIdReply, error) {
	out := new(IncrIdReply)
	err := grpc.Invoke(ctx, "/incrementid.IncrementId/GetIncrId", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *incrementIdClient) GetIncrIdByAmount(ctx context.Context, in *IncrIdNameValueRequest, opts ...grpc.CallOption) (*IncrIdReply, error) {
	out := new(IncrIdReply)
	err := grpc.Invoke(ctx, "/incrementid.IncrementId/GetIncrIdByAmount", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *incrementIdClient) CheckIncrKeyExist(ctx context.Context, in *IncrIdNameRequest, opts ...grpc.CallOption) (*IncrBoolReply, error) {
	out := new(IncrBoolReply)
	err := grpc.Invoke(ctx, "/incrementid.IncrementId/CheckIncrKeyExist", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *incrementIdClient) CreateIncrKey(ctx context.Context, in *IncrIdNameValueRequest, opts ...grpc.CallOption) (*NoneReply, error) {
	out := new(NoneReply)
	err := grpc.Invoke(ctx, "/incrementid.IncrementId/CreateIncrKey", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for IncrementId service

type IncrementIdServer interface {
	GetIncrId(context.Context, *IncrIdNameRequest) (*IncrIdReply, error)
	GetIncrIdByAmount(context.Context, *IncrIdNameValueRequest) (*IncrIdReply, error)
	CheckIncrKeyExist(context.Context, *IncrIdNameRequest) (*IncrBoolReply, error)
	CreateIncrKey(context.Context, *IncrIdNameValueRequest) (*NoneReply, error)
}

func RegisterIncrementIdServer(s *grpc.Server, srv IncrementIdServer) {
	s.RegisterService(&_IncrementId_serviceDesc, srv)
}

func _IncrementId_GetIncrId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IncrIdNameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IncrementIdServer).GetIncrId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/incrementid.IncrementId/GetIncrId",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IncrementIdServer).GetIncrId(ctx, req.(*IncrIdNameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IncrementId_GetIncrIdByAmount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IncrIdNameValueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IncrementIdServer).GetIncrIdByAmount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/incrementid.IncrementId/GetIncrIdByAmount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IncrementIdServer).GetIncrIdByAmount(ctx, req.(*IncrIdNameValueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IncrementId_CheckIncrKeyExist_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IncrIdNameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IncrementIdServer).CheckIncrKeyExist(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/incrementid.IncrementId/CheckIncrKeyExist",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IncrementIdServer).CheckIncrKeyExist(ctx, req.(*IncrIdNameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IncrementId_CreateIncrKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IncrIdNameValueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IncrementIdServer).CreateIncrKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/incrementid.IncrementId/CreateIncrKey",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IncrementIdServer).CreateIncrKey(ctx, req.(*IncrIdNameValueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _IncrementId_serviceDesc = grpc.ServiceDesc{
	ServiceName: "incrementid.IncrementId",
	HandlerType: (*IncrementIdServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetIncrId",
			Handler:    _IncrementId_GetIncrId_Handler,
		},
		{
			MethodName: "GetIncrIdByAmount",
			Handler:    _IncrementId_GetIncrIdByAmount_Handler,
		},
		{
			MethodName: "CheckIncrKeyExist",
			Handler:    _IncrementId_CheckIncrKeyExist_Handler,
		},
		{
			MethodName: "CreateIncrKey",
			Handler:    _IncrementId_CreateIncrKey_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "incrementid.proto",
}

func init() { proto.RegisterFile("incrementid.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 278 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x8c, 0x92, 0xd1, 0x4a, 0xf3, 0x40,
	0x10, 0x85, 0xd3, 0xfc, 0xfd, 0xc5, 0x4c, 0xa8, 0x98, 0x41, 0x4a, 0x28, 0x28, 0xba, 0x5e, 0xe8,
	0x55, 0x2f, 0xf4, 0x09, 0x4c, 0x91, 0x12, 0x84, 0x80, 0x41, 0xbc, 0x8f, 0xd9, 0x01, 0x17, 0x93,
	0xdd, 0xba, 0xdd, 0x88, 0x79, 0x4a, 0x5f, 0x49, 0xb2, 0x2b, 0x31, 0x52, 0xad, 0xde, 0xcd, 0xce,
	0x7c, 0x9c, 0x33, 0x67, 0x58, 0x88, 0x84, 0x2c, 0x35, 0xd5, 0x24, 0x8d, 0xe0, 0xf3, 0x95, 0x56,
	0x46, 0x61, 0x38, 0x68, 0xb1, 0x33, 0x88, 0x52, 0x59, 0xea, 0x94, 0x67, 0x45, 0x4d, 0x39, 0x3d,
	0x37, 0xb4, 0x36, 0x88, 0x30, 0x96, 0x45, 0x4d, 0xf1, 0xe8, 0x78, 0x74, 0x1e, 0xe4, 0xb6, 0x66,
	0x87, 0x10, 0x3a, 0x30, 0xa7, 0x55, 0xd5, 0xe2, 0x1e, 0xf8, 0x82, 0x5b, 0x60, 0x9c, 0xfb, 0x82,
	0xb3, 0x13, 0x98, 0x74, 0xe3, 0x44, 0xa9, 0xca, 0x01, 0xfb, 0xf0, 0x4f, 0x93, 0xb1, 0xc4, 0x6e,
	0xde, 0x95, 0x2c, 0x81, 0xe9, 0xa7, 0xd5, 0x7d, 0x51, 0x35, 0xdb, 0xfc, 0xf0, 0x00, 0xfe, 0xbf,
	0x74, 0x4c, 0xec, 0x5b, 0x0f, 0xf7, 0x60, 0x21, 0x04, 0x99, 0x92, 0x64, 0x2d, 0x2e, 0xde, 0x7c,
	0xb7, 0x93, 0xcd, 0x92, 0x72, 0x5c, 0x42, 0xb0, 0x24, 0xe3, 0x3c, 0xf0, 0x68, 0x3e, 0x4c, 0xbe,
	0x91, 0x71, 0x16, 0x7f, 0x33, 0xb7, 0xb2, 0xcc, 0xc3, 0x3b, 0x88, 0x7a, 0xa1, 0xa4, 0xbd, 0xaa,
	0x55, 0x23, 0x0d, 0x9e, 0xfe, 0x20, 0x38, 0x4c, 0xb2, 0x55, 0xf5, 0x16, 0xa2, 0xc5, 0x23, 0x95,
	0x4f, 0x5d, 0xf7, 0x86, 0xda, 0xeb, 0x57, 0xb1, 0x36, 0xbf, 0xae, 0x39, 0xdb, 0x98, 0xf7, 0x27,
	0x66, 0x1e, 0x66, 0x30, 0x59, 0x68, 0x2a, 0x0c, 0x7d, 0x68, 0xfe, 0x6d, 0xc9, 0xe9, 0x17, 0xa8,
	0xbf, 0x27, 0xf3, 0x1e, 0x76, 0xec, 0x0f, 0xb9, 0x7c, 0x0f, 0x00, 0x00, 0xff, 0xff, 0x41, 0xf6,
	0xc2, 0xdb, 0x36, 0x02, 0x00, 0x00,
}
