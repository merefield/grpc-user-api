// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package user

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// TagServiceClient is the client API for TagService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TagServiceClient interface {
	SearchGenres(ctx context.Context, in *Query, opts ...grpc.CallOption) (*SearchResults, error)
}

type tagServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTagServiceClient(cc grpc.ClientConnInterface) TagServiceClient {
	return &tagServiceClient{cc}
}

func (c *tagServiceClient) SearchGenres(ctx context.Context, in *Query, opts ...grpc.CallOption) (*SearchResults, error) {
	out := new(SearchResults)
	err := c.cc.Invoke(ctx, "/user.TagService/SearchGenres", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TagServiceServer is the server API for TagService service.
// All implementations should embed UnimplementedTagServiceServer
// for forward compatibility
type TagServiceServer interface {
	SearchGenres(context.Context, *Query) (*SearchResults, error)
}

// UnimplementedTagServiceServer should be embedded to have forward compatible implementations.
type UnimplementedTagServiceServer struct {
}

func (UnimplementedTagServiceServer) SearchGenres(context.Context, *Query) (*SearchResults, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchGenres not implemented")
}

// UnsafeTagServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TagServiceServer will
// result in compilation errors.
type UnsafeTagServiceServer interface {
	mustEmbedUnimplementedTagServiceServer()
}

func RegisterTagServiceServer(s grpc.ServiceRegistrar, srv TagServiceServer) {
	s.RegisterService(&_TagService_serviceDesc, srv)
}

func _TagService_SearchGenres_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Query)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TagServiceServer).SearchGenres(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.TagService/SearchGenres",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TagServiceServer).SearchGenres(ctx, req.(*Query))
	}
	return interceptor(ctx, in, info, handler)
}

var _TagService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "user.TagService",
	HandlerType: (*TagServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SearchGenres",
			Handler:    _TagService_SearchGenres_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user/tag.proto",
}
