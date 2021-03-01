// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package iam

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// ResonateIAMClient is the client API for ResonateIAM service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ResonateIAMClient interface {
	// Authenticate user by username or email and password
	Auth(ctx context.Context, in *AuthReq, opts ...grpc.CallOption) (*AuthResp, error)
	// Refresh refreshes JWT token
	Refresh(ctx context.Context, in *RefreshReq, opts ...grpc.CallOption) (*RefreshResp, error)
}

type resonateIAMClient struct {
	cc grpc.ClientConnInterface
}

func NewResonateIAMClient(cc grpc.ClientConnInterface) ResonateIAMClient {
	return &resonateIAMClient{cc}
}

func (c *resonateIAMClient) Auth(ctx context.Context, in *AuthReq, opts ...grpc.CallOption) (*AuthResp, error) {
	out := new(AuthResp)
	err := c.cc.Invoke(ctx, "/iam.ResonateIAM/Auth", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *resonateIAMClient) Refresh(ctx context.Context, in *RefreshReq, opts ...grpc.CallOption) (*RefreshResp, error) {
	out := new(RefreshResp)
	err := c.cc.Invoke(ctx, "/iam.ResonateIAM/Refresh", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ResonateIAMServer is the server API for ResonateIAM service.
// All implementations should embed UnimplementedResonateIAMServer
// for forward compatibility
type ResonateIAMServer interface {
	// Authenticate user by username or email and password
	Auth(context.Context, *AuthReq) (*AuthResp, error)
	// Refresh refreshes JWT token
	Refresh(context.Context, *RefreshReq) (*RefreshResp, error)
}

// UnimplementedResonateIAMServer should be embedded to have forward compatible implementations.
type UnimplementedResonateIAMServer struct {
}

func (UnimplementedResonateIAMServer) Auth(context.Context, *AuthReq) (*AuthResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Auth not implemented")
}
func (UnimplementedResonateIAMServer) Refresh(context.Context, *RefreshReq) (*RefreshResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Refresh not implemented")
}

// UnsafeResonateIAMServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ResonateIAMServer will
// result in compilation errors.
type UnsafeResonateIAMServer interface {
	mustEmbedUnimplementedResonateIAMServer()
}

func RegisterResonateIAMServer(s grpc.ServiceRegistrar, srv ResonateIAMServer) {
	s.RegisterService(&_ResonateIAM_serviceDesc, srv)
}

func _ResonateIAM_Auth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResonateIAMServer).Auth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/iam.ResonateIAM/Auth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResonateIAMServer).Auth(ctx, req.(*AuthReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ResonateIAM_Refresh_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResonateIAMServer).Refresh(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/iam.ResonateIAM/Refresh",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResonateIAMServer).Refresh(ctx, req.(*RefreshReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _ResonateIAM_serviceDesc = grpc.ServiceDesc{
	ServiceName: "iam.ResonateIAM",
	HandlerType: (*ResonateIAMServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Auth",
			Handler:    _ResonateIAM_Auth_Handler,
		},
		{
			MethodName: "Refresh",
			Handler:    _ResonateIAM_Refresh_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "iam/iam.proto",
}
