// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.3
// source: api.proto

package protos

import (
	context "context"
	httpbody "google.golang.org/genproto/googleapis/api/httpbody"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ImageStorageServiceV1Client is the client API for ImageStorageServiceV1 service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ImageStorageServiceV1Client interface {
	UploadImage(ctx context.Context, in *UploadImageRequest, opts ...grpc.CallOption) (*UploadImageResponce, error)
	StreamingUploadImage(ctx context.Context, opts ...grpc.CallOption) (ImageStorageServiceV1_StreamingUploadImageClient, error)
	GetImage(ctx context.Context, in *ImageRequest, opts ...grpc.CallOption) (*httpbody.HttpBody, error)
	IsImageExist(ctx context.Context, in *ImageRequest, opts ...grpc.CallOption) (*ImageExistResponce, error)
	DeleteImage(ctx context.Context, in *ImageRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ReplaceImage(ctx context.Context, in *ReplaceImageRequest, opts ...grpc.CallOption) (*ReplaceImageResponce, error)
}

type imageStorageServiceV1Client struct {
	cc grpc.ClientConnInterface
}

func NewImageStorageServiceV1Client(cc grpc.ClientConnInterface) ImageStorageServiceV1Client {
	return &imageStorageServiceV1Client{cc}
}

func (c *imageStorageServiceV1Client) UploadImage(ctx context.Context, in *UploadImageRequest, opts ...grpc.CallOption) (*UploadImageResponce, error) {
	out := new(UploadImageResponce)
	err := c.cc.Invoke(ctx, "/image_storage_service.imageStorageServiceV1/UploadImage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *imageStorageServiceV1Client) StreamingUploadImage(ctx context.Context, opts ...grpc.CallOption) (ImageStorageServiceV1_StreamingUploadImageClient, error) {
	stream, err := c.cc.NewStream(ctx, &ImageStorageServiceV1_ServiceDesc.Streams[0], "/image_storage_service.imageStorageServiceV1/StreamingUploadImage", opts...)
	if err != nil {
		return nil, err
	}
	x := &imageStorageServiceV1StreamingUploadImageClient{stream}
	return x, nil
}

type ImageStorageServiceV1_StreamingUploadImageClient interface {
	Send(*StreamingUploadImageRequest) error
	CloseAndRecv() (*UploadImageResponce, error)
	grpc.ClientStream
}

type imageStorageServiceV1StreamingUploadImageClient struct {
	grpc.ClientStream
}

func (x *imageStorageServiceV1StreamingUploadImageClient) Send(m *StreamingUploadImageRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *imageStorageServiceV1StreamingUploadImageClient) CloseAndRecv() (*UploadImageResponce, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(UploadImageResponce)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *imageStorageServiceV1Client) GetImage(ctx context.Context, in *ImageRequest, opts ...grpc.CallOption) (*httpbody.HttpBody, error) {
	out := new(httpbody.HttpBody)
	err := c.cc.Invoke(ctx, "/image_storage_service.imageStorageServiceV1/GetImage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *imageStorageServiceV1Client) IsImageExist(ctx context.Context, in *ImageRequest, opts ...grpc.CallOption) (*ImageExistResponce, error) {
	out := new(ImageExistResponce)
	err := c.cc.Invoke(ctx, "/image_storage_service.imageStorageServiceV1/IsImageExist", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *imageStorageServiceV1Client) DeleteImage(ctx context.Context, in *ImageRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/image_storage_service.imageStorageServiceV1/DeleteImage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *imageStorageServiceV1Client) ReplaceImage(ctx context.Context, in *ReplaceImageRequest, opts ...grpc.CallOption) (*ReplaceImageResponce, error) {
	out := new(ReplaceImageResponce)
	err := c.cc.Invoke(ctx, "/image_storage_service.imageStorageServiceV1/ReplaceImage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ImageStorageServiceV1Server is the server API for ImageStorageServiceV1 service.
// All implementations must embed UnimplementedImageStorageServiceV1Server
// for forward compatibility
type ImageStorageServiceV1Server interface {
	UploadImage(context.Context, *UploadImageRequest) (*UploadImageResponce, error)
	StreamingUploadImage(ImageStorageServiceV1_StreamingUploadImageServer) error
	GetImage(context.Context, *ImageRequest) (*httpbody.HttpBody, error)
	IsImageExist(context.Context, *ImageRequest) (*ImageExistResponce, error)
	DeleteImage(context.Context, *ImageRequest) (*emptypb.Empty, error)
	ReplaceImage(context.Context, *ReplaceImageRequest) (*ReplaceImageResponce, error)
	mustEmbedUnimplementedImageStorageServiceV1Server()
}

// UnimplementedImageStorageServiceV1Server must be embedded to have forward compatible implementations.
type UnimplementedImageStorageServiceV1Server struct {
}

func (UnimplementedImageStorageServiceV1Server) UploadImage(context.Context, *UploadImageRequest) (*UploadImageResponce, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadImage not implemented")
}
func (UnimplementedImageStorageServiceV1Server) StreamingUploadImage(ImageStorageServiceV1_StreamingUploadImageServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamingUploadImage not implemented")
}
func (UnimplementedImageStorageServiceV1Server) GetImage(context.Context, *ImageRequest) (*httpbody.HttpBody, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetImage not implemented")
}
func (UnimplementedImageStorageServiceV1Server) IsImageExist(context.Context, *ImageRequest) (*ImageExistResponce, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsImageExist not implemented")
}
func (UnimplementedImageStorageServiceV1Server) DeleteImage(context.Context, *ImageRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteImage not implemented")
}
func (UnimplementedImageStorageServiceV1Server) ReplaceImage(context.Context, *ReplaceImageRequest) (*ReplaceImageResponce, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReplaceImage not implemented")
}
func (UnimplementedImageStorageServiceV1Server) mustEmbedUnimplementedImageStorageServiceV1Server() {}

// UnsafeImageStorageServiceV1Server may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ImageStorageServiceV1Server will
// result in compilation errors.
type UnsafeImageStorageServiceV1Server interface {
	mustEmbedUnimplementedImageStorageServiceV1Server()
}

func RegisterImageStorageServiceV1Server(s grpc.ServiceRegistrar, srv ImageStorageServiceV1Server) {
	s.RegisterService(&ImageStorageServiceV1_ServiceDesc, srv)
}

func _ImageStorageServiceV1_UploadImage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadImageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImageStorageServiceV1Server).UploadImage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/image_storage_service.imageStorageServiceV1/UploadImage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImageStorageServiceV1Server).UploadImage(ctx, req.(*UploadImageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ImageStorageServiceV1_StreamingUploadImage_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ImageStorageServiceV1Server).StreamingUploadImage(&imageStorageServiceV1StreamingUploadImageServer{stream})
}

type ImageStorageServiceV1_StreamingUploadImageServer interface {
	SendAndClose(*UploadImageResponce) error
	Recv() (*StreamingUploadImageRequest, error)
	grpc.ServerStream
}

type imageStorageServiceV1StreamingUploadImageServer struct {
	grpc.ServerStream
}

func (x *imageStorageServiceV1StreamingUploadImageServer) SendAndClose(m *UploadImageResponce) error {
	return x.ServerStream.SendMsg(m)
}

func (x *imageStorageServiceV1StreamingUploadImageServer) Recv() (*StreamingUploadImageRequest, error) {
	m := new(StreamingUploadImageRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _ImageStorageServiceV1_GetImage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ImageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImageStorageServiceV1Server).GetImage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/image_storage_service.imageStorageServiceV1/GetImage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImageStorageServiceV1Server).GetImage(ctx, req.(*ImageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ImageStorageServiceV1_IsImageExist_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ImageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImageStorageServiceV1Server).IsImageExist(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/image_storage_service.imageStorageServiceV1/IsImageExist",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImageStorageServiceV1Server).IsImageExist(ctx, req.(*ImageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ImageStorageServiceV1_DeleteImage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ImageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImageStorageServiceV1Server).DeleteImage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/image_storage_service.imageStorageServiceV1/DeleteImage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImageStorageServiceV1Server).DeleteImage(ctx, req.(*ImageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ImageStorageServiceV1_ReplaceImage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReplaceImageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImageStorageServiceV1Server).ReplaceImage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/image_storage_service.imageStorageServiceV1/ReplaceImage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImageStorageServiceV1Server).ReplaceImage(ctx, req.(*ReplaceImageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ImageStorageServiceV1_ServiceDesc is the grpc.ServiceDesc for ImageStorageServiceV1 service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ImageStorageServiceV1_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "image_storage_service.imageStorageServiceV1",
	HandlerType: (*ImageStorageServiceV1Server)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UploadImage",
			Handler:    _ImageStorageServiceV1_UploadImage_Handler,
		},
		{
			MethodName: "GetImage",
			Handler:    _ImageStorageServiceV1_GetImage_Handler,
		},
		{
			MethodName: "IsImageExist",
			Handler:    _ImageStorageServiceV1_IsImageExist_Handler,
		},
		{
			MethodName: "DeleteImage",
			Handler:    _ImageStorageServiceV1_DeleteImage_Handler,
		},
		{
			MethodName: "ReplaceImage",
			Handler:    _ImageStorageServiceV1_ReplaceImage_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamingUploadImage",
			Handler:       _ImageStorageServiceV1_StreamingUploadImage_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "api.proto",
}