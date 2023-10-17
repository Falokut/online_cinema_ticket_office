// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.3
// source: protos/userService.proto

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

// UserServiceV1Client is the client API for UserServiceV1 service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserServiceV1Client interface {
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponce, error)
	GetUserByEmailAndPassword(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserResponce, error)
	VerifyUserAccount(ctx context.Context, in *VerifyRequest, opts ...grpc.CallOption) (*VerifyResponce, error)
	RequestAccountVerification(ctx context.Context, in *VerificationRequest, opts ...grpc.CallOption) (*VerificationResponce, error)
	RequestChangePasswordToken(ctx context.Context, in *ChangePasswordTokenRequest, opts ...grpc.CallOption) (*ChangePasswordToken, error)
	GetUserProfile(ctx context.Context, in *GetUserProfileRequest, opts ...grpc.CallOption) (*GetUserProfileResponce, error)
	ChangePassword(ctx context.Context, in *ChangePasswordRequest, opts ...grpc.CallOption) (*ChangePasswordResponce, error)
	IsUserWithEmailExist(ctx context.Context, in *IsExistRequest, opts ...grpc.CallOption) (*IsExistResponce, error)
	UpdateProfilePicture(ctx context.Context, in *UpdateProfilePictureRequest, opts ...grpc.CallOption) (*UpdateProfilePictureResponce, error)
}

type userServiceV1Client struct {
	cc grpc.ClientConnInterface
}

func NewUserServiceV1Client(cc grpc.ClientConnInterface) UserServiceV1Client {
	return &userServiceV1Client{cc}
}

func (c *userServiceV1Client) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponce, error) {
	out := new(CreateUserResponce)
	err := c.cc.Invoke(ctx, "/user_service.userServiceV1/CreateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceV1Client) GetUserByEmailAndPassword(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserResponce, error) {
	out := new(GetUserResponce)
	err := c.cc.Invoke(ctx, "/user_service.userServiceV1/GetUserByEmailAndPassword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceV1Client) VerifyUserAccount(ctx context.Context, in *VerifyRequest, opts ...grpc.CallOption) (*VerifyResponce, error) {
	out := new(VerifyResponce)
	err := c.cc.Invoke(ctx, "/user_service.userServiceV1/VerifyUserAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceV1Client) RequestAccountVerification(ctx context.Context, in *VerificationRequest, opts ...grpc.CallOption) (*VerificationResponce, error) {
	out := new(VerificationResponce)
	err := c.cc.Invoke(ctx, "/user_service.userServiceV1/RequestAccountVerification", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceV1Client) RequestChangePasswordToken(ctx context.Context, in *ChangePasswordTokenRequest, opts ...grpc.CallOption) (*ChangePasswordToken, error) {
	out := new(ChangePasswordToken)
	err := c.cc.Invoke(ctx, "/user_service.userServiceV1/RequestChangePasswordToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceV1Client) GetUserProfile(ctx context.Context, in *GetUserProfileRequest, opts ...grpc.CallOption) (*GetUserProfileResponce, error) {
	out := new(GetUserProfileResponce)
	err := c.cc.Invoke(ctx, "/user_service.userServiceV1/GetUserProfile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceV1Client) ChangePassword(ctx context.Context, in *ChangePasswordRequest, opts ...grpc.CallOption) (*ChangePasswordResponce, error) {
	out := new(ChangePasswordResponce)
	err := c.cc.Invoke(ctx, "/user_service.userServiceV1/ChangePassword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceV1Client) IsUserWithEmailExist(ctx context.Context, in *IsExistRequest, opts ...grpc.CallOption) (*IsExistResponce, error) {
	out := new(IsExistResponce)
	err := c.cc.Invoke(ctx, "/user_service.userServiceV1/IsUserWithEmailExist", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceV1Client) UpdateProfilePicture(ctx context.Context, in *UpdateProfilePictureRequest, opts ...grpc.CallOption) (*UpdateProfilePictureResponce, error) {
	out := new(UpdateProfilePictureResponce)
	err := c.cc.Invoke(ctx, "/user_service.userServiceV1/UpdateProfilePicture", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServiceV1Server is the server API for UserServiceV1 service.
// All implementations must embed UnimplementedUserServiceV1Server
// for forward compatibility
type UserServiceV1Server interface {
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponce, error)
	GetUserByEmailAndPassword(context.Context, *GetUserRequest) (*GetUserResponce, error)
	VerifyUserAccount(context.Context, *VerifyRequest) (*VerifyResponce, error)
	RequestAccountVerification(context.Context, *VerificationRequest) (*VerificationResponce, error)
	RequestChangePasswordToken(context.Context, *ChangePasswordTokenRequest) (*ChangePasswordToken, error)
	GetUserProfile(context.Context, *GetUserProfileRequest) (*GetUserProfileResponce, error)
	ChangePassword(context.Context, *ChangePasswordRequest) (*ChangePasswordResponce, error)
	IsUserWithEmailExist(context.Context, *IsExistRequest) (*IsExistResponce, error)
	UpdateProfilePicture(context.Context, *UpdateProfilePictureRequest) (*UpdateProfilePictureResponce, error)
	mustEmbedUnimplementedUserServiceV1Server()
}

// UnimplementedUserServiceV1Server must be embedded to have forward compatible implementations.
type UnimplementedUserServiceV1Server struct {
}

func (UnimplementedUserServiceV1Server) CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponce, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedUserServiceV1Server) GetUserByEmailAndPassword(context.Context, *GetUserRequest) (*GetUserResponce, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserByEmailAndPassword not implemented")
}
func (UnimplementedUserServiceV1Server) VerifyUserAccount(context.Context, *VerifyRequest) (*VerifyResponce, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyUserAccount not implemented")
}
func (UnimplementedUserServiceV1Server) RequestAccountVerification(context.Context, *VerificationRequest) (*VerificationResponce, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestAccountVerification not implemented")
}
func (UnimplementedUserServiceV1Server) RequestChangePasswordToken(context.Context, *ChangePasswordTokenRequest) (*ChangePasswordToken, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestChangePasswordToken not implemented")
}
func (UnimplementedUserServiceV1Server) GetUserProfile(context.Context, *GetUserProfileRequest) (*GetUserProfileResponce, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserProfile not implemented")
}
func (UnimplementedUserServiceV1Server) ChangePassword(context.Context, *ChangePasswordRequest) (*ChangePasswordResponce, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangePassword not implemented")
}
func (UnimplementedUserServiceV1Server) IsUserWithEmailExist(context.Context, *IsExistRequest) (*IsExistResponce, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsUserWithEmailExist not implemented")
}
func (UnimplementedUserServiceV1Server) UpdateProfilePicture(context.Context, *UpdateProfilePictureRequest) (*UpdateProfilePictureResponce, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateProfilePicture not implemented")
}
func (UnimplementedUserServiceV1Server) mustEmbedUnimplementedUserServiceV1Server() {}

// UnsafeUserServiceV1Server may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServiceV1Server will
// result in compilation errors.
type UnsafeUserServiceV1Server interface {
	mustEmbedUnimplementedUserServiceV1Server()
}

func RegisterUserServiceV1Server(s grpc.ServiceRegistrar, srv UserServiceV1Server) {
	s.RegisterService(&UserServiceV1_ServiceDesc, srv)
}

func _UserServiceV1_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceV1Server).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user_service.userServiceV1/CreateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceV1Server).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserServiceV1_GetUserByEmailAndPassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceV1Server).GetUserByEmailAndPassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user_service.userServiceV1/GetUserByEmailAndPassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceV1Server).GetUserByEmailAndPassword(ctx, req.(*GetUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserServiceV1_VerifyUserAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceV1Server).VerifyUserAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user_service.userServiceV1/VerifyUserAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceV1Server).VerifyUserAccount(ctx, req.(*VerifyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserServiceV1_RequestAccountVerification_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerificationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceV1Server).RequestAccountVerification(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user_service.userServiceV1/RequestAccountVerification",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceV1Server).RequestAccountVerification(ctx, req.(*VerificationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserServiceV1_RequestChangePasswordToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangePasswordTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceV1Server).RequestChangePasswordToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user_service.userServiceV1/RequestChangePasswordToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceV1Server).RequestChangePasswordToken(ctx, req.(*ChangePasswordTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserServiceV1_GetUserProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceV1Server).GetUserProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user_service.userServiceV1/GetUserProfile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceV1Server).GetUserProfile(ctx, req.(*GetUserProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserServiceV1_ChangePassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangePasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceV1Server).ChangePassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user_service.userServiceV1/ChangePassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceV1Server).ChangePassword(ctx, req.(*ChangePasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserServiceV1_IsUserWithEmailExist_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsExistRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceV1Server).IsUserWithEmailExist(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user_service.userServiceV1/IsUserWithEmailExist",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceV1Server).IsUserWithEmailExist(ctx, req.(*IsExistRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserServiceV1_UpdateProfilePicture_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateProfilePictureRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceV1Server).UpdateProfilePicture(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user_service.userServiceV1/UpdateProfilePicture",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceV1Server).UpdateProfilePicture(ctx, req.(*UpdateProfilePictureRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserServiceV1_ServiceDesc is the grpc.ServiceDesc for UserServiceV1 service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserServiceV1_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user_service.userServiceV1",
	HandlerType: (*UserServiceV1Server)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUser",
			Handler:    _UserServiceV1_CreateUser_Handler,
		},
		{
			MethodName: "GetUserByEmailAndPassword",
			Handler:    _UserServiceV1_GetUserByEmailAndPassword_Handler,
		},
		{
			MethodName: "VerifyUserAccount",
			Handler:    _UserServiceV1_VerifyUserAccount_Handler,
		},
		{
			MethodName: "RequestAccountVerification",
			Handler:    _UserServiceV1_RequestAccountVerification_Handler,
		},
		{
			MethodName: "RequestChangePasswordToken",
			Handler:    _UserServiceV1_RequestChangePasswordToken_Handler,
		},
		{
			MethodName: "GetUserProfile",
			Handler:    _UserServiceV1_GetUserProfile_Handler,
		},
		{
			MethodName: "ChangePassword",
			Handler:    _UserServiceV1_ChangePassword_Handler,
		},
		{
			MethodName: "IsUserWithEmailExist",
			Handler:    _UserServiceV1_IsUserWithEmailExist_Handler,
		},
		{
			MethodName: "UpdateProfilePicture",
			Handler:    _UserServiceV1_UpdateProfilePicture_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/userService.proto",
}
