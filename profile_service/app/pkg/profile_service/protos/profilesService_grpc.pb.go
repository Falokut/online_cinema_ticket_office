// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.3
// source: app/protos/profilesService.proto

package protos

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

// ProfileServiceV1Client is the client API for ProfileServiceV1 service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProfileServiceV1Client interface {
	GetUserProfile(ctx context.Context, in *GetUserProfileRequest, opts ...grpc.CallOption) (*GetUserProfileResponce, error)
	UpdateProfilePicture(ctx context.Context, opts ...grpc.CallOption) (ProfileServiceV1_UpdateProfilePictureClient, error)
}

type profileServiceV1Client struct {
	cc grpc.ClientConnInterface
}

func NewProfileServiceV1Client(cc grpc.ClientConnInterface) ProfileServiceV1Client {
	return &profileServiceV1Client{cc}
}

func (c *profileServiceV1Client) GetUserProfile(ctx context.Context, in *GetUserProfileRequest, opts ...grpc.CallOption) (*GetUserProfileResponce, error) {
	out := new(GetUserProfileResponce)
	err := c.cc.Invoke(ctx, "/profile_service.profileServiceV1/GetUserProfile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileServiceV1Client) UpdateProfilePicture(ctx context.Context, opts ...grpc.CallOption) (ProfileServiceV1_UpdateProfilePictureClient, error) {
	stream, err := c.cc.NewStream(ctx, &ProfileServiceV1_ServiceDesc.Streams[0], "/profile_service.profileServiceV1/UpdateProfilePicture", opts...)
	if err != nil {
		return nil, err
	}
	x := &profileServiceV1UpdateProfilePictureClient{stream}
	return x, nil
}

type ProfileServiceV1_UpdateProfilePictureClient interface {
	Send(*UpdateProfilePictureRequest) error
	CloseAndRecv() (*UpdateProfilePictureResponce, error)
	grpc.ClientStream
}

type profileServiceV1UpdateProfilePictureClient struct {
	grpc.ClientStream
}

func (x *profileServiceV1UpdateProfilePictureClient) Send(m *UpdateProfilePictureRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *profileServiceV1UpdateProfilePictureClient) CloseAndRecv() (*UpdateProfilePictureResponce, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(UpdateProfilePictureResponce)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ProfileServiceV1Server is the server API for ProfileServiceV1 service.
// All implementations must embed UnimplementedProfileServiceV1Server
// for forward compatibility
type ProfileServiceV1Server interface {
	GetUserProfile(context.Context, *GetUserProfileRequest) (*GetUserProfileResponce, error)
	UpdateProfilePicture(ProfileServiceV1_UpdateProfilePictureServer) error
	mustEmbedUnimplementedProfileServiceV1Server()
}

// UnimplementedProfileServiceV1Server must be embedded to have forward compatible implementations.
type UnimplementedProfileServiceV1Server struct {
}

func (UnimplementedProfileServiceV1Server) GetUserProfile(context.Context, *GetUserProfileRequest) (*GetUserProfileResponce, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserProfile not implemented")
}
func (UnimplementedProfileServiceV1Server) UpdateProfilePicture(ProfileServiceV1_UpdateProfilePictureServer) error {
	return status.Errorf(codes.Unimplemented, "method UpdateProfilePicture not implemented")
}
func (UnimplementedProfileServiceV1Server) mustEmbedUnimplementedProfileServiceV1Server() {}

// UnsafeProfileServiceV1Server may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProfileServiceV1Server will
// result in compilation errors.
type UnsafeProfileServiceV1Server interface {
	mustEmbedUnimplementedProfileServiceV1Server()
}

func RegisterProfileServiceV1Server(s grpc.ServiceRegistrar, srv ProfileServiceV1Server) {
	s.RegisterService(&ProfileServiceV1_ServiceDesc, srv)
}

func _ProfileServiceV1_GetUserProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServiceV1Server).GetUserProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/profile_service.profileServiceV1/GetUserProfile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServiceV1Server).GetUserProfile(ctx, req.(*GetUserProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileServiceV1_UpdateProfilePicture_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ProfileServiceV1Server).UpdateProfilePicture(&profileServiceV1UpdateProfilePictureServer{stream})
}

type ProfileServiceV1_UpdateProfilePictureServer interface {
	SendAndClose(*UpdateProfilePictureResponce) error
	Recv() (*UpdateProfilePictureRequest, error)
	grpc.ServerStream
}

type profileServiceV1UpdateProfilePictureServer struct {
	grpc.ServerStream
}

func (x *profileServiceV1UpdateProfilePictureServer) SendAndClose(m *UpdateProfilePictureResponce) error {
	return x.ServerStream.SendMsg(m)
}

func (x *profileServiceV1UpdateProfilePictureServer) Recv() (*UpdateProfilePictureRequest, error) {
	m := new(UpdateProfilePictureRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ProfileServiceV1_ServiceDesc is the grpc.ServiceDesc for ProfileServiceV1 service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProfileServiceV1_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "profile_service.profileServiceV1",
	HandlerType: (*ProfileServiceV1Server)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUserProfile",
			Handler:    _ProfileServiceV1_GetUserProfile_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UpdateProfilePicture",
			Handler:       _ProfileServiceV1_UpdateProfilePicture_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "app/protos/profilesService.proto",
}
