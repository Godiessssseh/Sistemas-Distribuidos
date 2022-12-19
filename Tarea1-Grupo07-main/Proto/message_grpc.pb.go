// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: Proto/message.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MessageServiceClient is the client API for MessageService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MessageServiceClient interface {
	// Aca se dice que la funcion se llama Itercambio y recibe como parametro un dato de tipo mensaje y
	// retorna otro dato de tipo mensaje. NO NECESARIAMENTE EL RETURN DEBE SER DEL MISMO TIPO DE DATO QUE
	// EL INPUT
	Intercambio(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Message, error)
	Revision(ctx context.Context, in *Revisando, opts ...grpc.CallOption) (*Revisando, error)
}

type messageServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMessageServiceClient(cc grpc.ClientConnInterface) MessageServiceClient {
	return &messageServiceClient{cc}
}

func (c *messageServiceClient) Intercambio(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Message, error) {
	out := new(Message)
	err := c.cc.Invoke(ctx, "/grpc.MessageService/Intercambio", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messageServiceClient) Revision(ctx context.Context, in *Revisando, opts ...grpc.CallOption) (*Revisando, error) {
	out := new(Revisando)
	err := c.cc.Invoke(ctx, "/grpc.MessageService/Revision", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MessageServiceServer is the server API for MessageService service.
// All implementations must embed UnimplementedMessageServiceServer
// for forward compatibility
type MessageServiceServer interface {
	// Aca se dice que la funcion se llama Itercambio y recibe como parametro un dato de tipo mensaje y
	// retorna otro dato de tipo mensaje. NO NECESARIAMENTE EL RETURN DEBE SER DEL MISMO TIPO DE DATO QUE
	// EL INPUT
	Intercambio(context.Context, *Message) (*Message, error)
	Revision(context.Context, *Revisando) (*Revisando, error)
	mustEmbedUnimplementedMessageServiceServer()
}

// UnimplementedMessageServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMessageServiceServer struct {
}

func (UnimplementedMessageServiceServer) Intercambio(context.Context, *Message) (*Message, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Intercambio not implemented")
}
func (UnimplementedMessageServiceServer) Revision(context.Context, *Revisando) (*Revisando, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Revision not implemented")
}
func (UnimplementedMessageServiceServer) mustEmbedUnimplementedMessageServiceServer() {}

// UnsafeMessageServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MessageServiceServer will
// result in compilation errors.
type UnsafeMessageServiceServer interface {
	mustEmbedUnimplementedMessageServiceServer()
}

func RegisterMessageServiceServer(s grpc.ServiceRegistrar, srv MessageServiceServer) {
	s.RegisterService(&MessageService_ServiceDesc, srv)
}

func _MessageService_Intercambio_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageServiceServer).Intercambio(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.MessageService/Intercambio",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageServiceServer).Intercambio(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

func _MessageService_Revision_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Revisando)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageServiceServer).Revision(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.MessageService/Revision",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageServiceServer).Revision(ctx, req.(*Revisando))
	}
	return interceptor(ctx, in, info, handler)
}

// MessageService_ServiceDesc is the grpc.ServiceDesc for MessageService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MessageService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.MessageService",
	HandlerType: (*MessageServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Intercambio",
			Handler:    _MessageService_Intercambio_Handler,
		},
		{
			MethodName: "Revision",
			Handler:    _MessageService_Revision_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "Proto/message.proto",
}
