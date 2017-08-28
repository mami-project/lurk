// Code generated by protoc-gen-go.
// source: publisher.proto
// DO NOT EDIT!

/*
Package publisher is a generated protocol buffer package.

It is generated from these files:
	publisher.proto

It has these top-level messages:
	Request
	Empty
*/
package publisher

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

type Request struct {
	Der              []byte  `protobuf:"bytes,1,opt,name=der" json:"der,omitempty"`
	LogURL           *string `protobuf:"bytes,2,opt,name=LogURL" json:"LogURL,omitempty"`
	LogPublicKey     *string `protobuf:"bytes,3,opt,name=LogPublicKey" json:"LogPublicKey,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Request) Reset()                    { *m = Request{} }
func (m *Request) String() string            { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()               {}
func (*Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Request) GetDer() []byte {
	if m != nil {
		return m.Der
	}
	return nil
}

func (m *Request) GetLogURL() string {
	if m != nil && m.LogURL != nil {
		return *m.LogURL
	}
	return ""
}

func (m *Request) GetLogPublicKey() string {
	if m != nil && m.LogPublicKey != nil {
		return *m.LogPublicKey
	}
	return ""
}

type Empty struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func init() {
	proto.RegisterType((*Request)(nil), "Request")
	proto.RegisterType((*Empty)(nil), "Empty")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Publisher service

type PublisherClient interface {
	SubmitToCT(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Empty, error)
	SubmitToSingleCT(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Empty, error)
}

type publisherClient struct {
	cc *grpc.ClientConn
}

func NewPublisherClient(cc *grpc.ClientConn) PublisherClient {
	return &publisherClient{cc}
}

func (c *publisherClient) SubmitToCT(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/Publisher/SubmitToCT", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *publisherClient) SubmitToSingleCT(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/Publisher/SubmitToSingleCT", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Publisher service

type PublisherServer interface {
	SubmitToCT(context.Context, *Request) (*Empty, error)
	SubmitToSingleCT(context.Context, *Request) (*Empty, error)
}

func RegisterPublisherServer(s *grpc.Server, srv PublisherServer) {
	s.RegisterService(&_Publisher_serviceDesc, srv)
}

func _Publisher_SubmitToCT_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PublisherServer).SubmitToCT(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Publisher/SubmitToCT",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PublisherServer).SubmitToCT(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Publisher_SubmitToSingleCT_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PublisherServer).SubmitToSingleCT(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Publisher/SubmitToSingleCT",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PublisherServer).SubmitToSingleCT(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

var _Publisher_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Publisher",
	HandlerType: (*PublisherServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SubmitToCT",
			Handler:    _Publisher_SubmitToCT_Handler,
		},
		{
			MethodName: "SubmitToSingleCT",
			Handler:    _Publisher_SubmitToSingleCT_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "publisher.proto",
}

func init() { proto.RegisterFile("publisher.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 155 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2f, 0x28, 0x4d, 0xca,
	0xc9, 0x2c, 0xce, 0x48, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57, 0xb2, 0xe1, 0x62, 0x0f,
	0x4a, 0x2d, 0x2c, 0x4d, 0x2d, 0x2e, 0x11, 0xe2, 0xe6, 0x62, 0x4e, 0x49, 0x2d, 0x92, 0x60, 0x54,
	0x60, 0xd4, 0xe0, 0x11, 0xe2, 0xe3, 0x62, 0xf3, 0xc9, 0x4f, 0x0f, 0x0d, 0xf2, 0x91, 0x60, 0x52,
	0x60, 0xd4, 0xe0, 0x14, 0x12, 0xe1, 0xe2, 0xf1, 0xc9, 0x4f, 0x0f, 0x00, 0xe9, 0x4e, 0xf6, 0x4e,
	0xad, 0x94, 0x60, 0x06, 0x89, 0x2a, 0xb1, 0x73, 0xb1, 0xba, 0xe6, 0x16, 0x94, 0x54, 0x1a, 0x85,
	0x72, 0x71, 0x06, 0xc0, 0x4c, 0x16, 0x52, 0xe0, 0xe2, 0x0a, 0x2e, 0x4d, 0xca, 0xcd, 0x2c, 0x09,
	0xc9, 0x77, 0x0e, 0x11, 0xe2, 0xd0, 0x83, 0x5a, 0x20, 0xc5, 0xa6, 0x07, 0x56, 0xac, 0xc4, 0x20,
	0xa4, 0xc6, 0x25, 0x00, 0x53, 0x11, 0x9c, 0x99, 0x97, 0x9e, 0x93, 0x8a, 0x5d, 0x1d, 0x20, 0x00,
	0x00, 0xff, 0xff, 0x19, 0x25, 0xca, 0x2f, 0xaf, 0x00, 0x00, 0x00,
}