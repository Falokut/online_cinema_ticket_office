// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.24.3
// source: api.proto

package protos

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_api_proto protoreflect.FileDescriptor

var file_api_proto_rawDesc = []byte{
	0x0a, 0x09, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x61, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x1a, 0x0e, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32,
	0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x98, 0x09, 0x0a, 0x10, 0x61,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x56, 0x31, 0x12,
	0x6c, 0x0a, 0x0d, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x12, 0x25, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x56, 0x6f, 0x69, 0x64, 0x22, 0x1d,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x17, 0x22, 0x12, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x2d, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x3a, 0x01, 0x2a, 0x12, 0x86, 0x01,
	0x0a, 0x1f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x12, 0x29, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x61,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x56,
	0x6f, 0x69, 0x64, 0x22, 0x21, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1b, 0x12, 0x19, 0x2f, 0x76, 0x31,
	0x2f, 0x76, 0x65, 0x72, 0x66, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x65, 0x0a, 0x0d, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79,
	0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x25, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79,
	0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15,
	0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x56, 0x6f, 0x69, 0x64, 0x22, 0x16, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x10, 0x22, 0x0b, 0x2f,
	0x76, 0x31, 0x2f, 0x2f, 0x76, 0x65, 0x72, 0x69, 0x66, 0x79, 0x3a, 0x01, 0x2a, 0x12, 0x61, 0x0a,
	0x06, 0x53, 0x69, 0x67, 0x6e, 0x49, 0x6e, 0x12, 0x1e, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x49, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x63, 0x65, 0x22, 0x16, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x10,
	0x22, 0x0b, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x69, 0x67, 0x6e, 0x2d, 0x69, 0x6e, 0x3a, 0x01, 0x2a,
	0x12, 0x6b, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44,
	0x12, 0x24, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x49, 0x44, 0x22, 0x19, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x13, 0x12, 0x0e, 0x2f, 0x76, 0x31, 0x2f,
	0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2d, 0x69, 0x64, 0x3a, 0x01, 0x2a, 0x12, 0x56, 0x0a,
	0x06, 0x4c, 0x6f, 0x67, 0x6f, 0x75, 0x74, 0x12, 0x1e, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4c, 0x6f, 0x67, 0x6f, 0x75, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x56, 0x6f, 0x69, 0x64, 0x22, 0x15,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0f, 0x22, 0x0a, 0x2f, 0x76, 0x31, 0x2f, 0x6c, 0x6f, 0x67, 0x6f,
	0x75, 0x74, 0x3a, 0x01, 0x2a, 0x12, 0x9c, 0x01, 0x0a, 0x1a, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x2b, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x50, 0x61, 0x73,
	0x73, 0x77, 0x6f, 0x72, 0x64, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x2c, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f,
	0x72, 0x64, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x63, 0x65, 0x22,
	0x23, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1d, 0x12, 0x1b, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x68, 0x61,
	0x6e, 0x67, 0x65, 0x2d, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x2f, 0x72, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x6f, 0x0a, 0x0e, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x50, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x26, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x50,
	0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15,
	0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x56, 0x6f, 0x69, 0x64, 0x22, 0x1e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x18, 0x22, 0x13, 0x2f,
	0x76, 0x31, 0x2f, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x2d, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f,
	0x72, 0x64, 0x3a, 0x01, 0x2a, 0x12, 0x74, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x53,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x26, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c,
	0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x24, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x41, 0x6c, 0x6c, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x63, 0x65, 0x22, 0x14, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0e, 0x12, 0x0c, 0x2f,
	0x76, 0x31, 0x2f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x78, 0x0a, 0x11, 0x54,
	0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73,
	0x12, 0x29, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x54, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x53, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x61, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x56, 0x6f,
	0x69, 0x64, 0x22, 0x21, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1b, 0x22, 0x16, 0x2f, 0x76, 0x31, 0x2f,
	0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x74, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61,
	0x74, 0x65, 0x3a, 0x01, 0x2a, 0x42, 0x9f, 0x01, 0x5a, 0x19, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x73, 0x92, 0x41, 0x80, 0x01, 0x12, 0x57, 0x0a, 0x0f, 0x41, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x20, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x22, 0x3f, 0x0a, 0x07, 0x46, 0x61,
	0x6c, 0x6f, 0x6b, 0x75, 0x74, 0x12, 0x1a, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x46, 0x61, 0x6c, 0x6f, 0x6b, 0x75,
	0x74, 0x1a, 0x18, 0x74, 0x69, 0x6d, 0x75, 0x72, 0x2e, 0x73, 0x69, 0x6e, 0x65, 0x6c, 0x6e, 0x69,
	0x6b, 0x40, 0x79, 0x61, 0x6e, 0x64, 0x65, 0x78, 0x2e, 0x72, 0x75, 0x32, 0x03, 0x31, 0x2e, 0x30,
	0x2a, 0x01, 0x01, 0x32, 0x10, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2f, 0x6a, 0x73, 0x6f, 0x6e, 0x3a, 0x10, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x2f, 0x6a, 0x73, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_api_proto_goTypes = []interface{}{
	(*CreateAccountRequest)(nil),        // 0: account_service.CreateAccountRequest
	(*VerificationTokenRequest)(nil),    // 1: account_service.VerificationTokenRequest
	(*VerifyAccountRequest)(nil),        // 2: account_service.VerifyAccountRequest
	(*SignInRequest)(nil),               // 3: account_service.SignInRequest
	(*GetAccountIDRequest)(nil),         // 4: account_service.GetAccountIDRequest
	(*LogoutRequest)(nil),               // 5: account_service.LogoutRequest
	(*ChangePasswordTokenRequest)(nil),  // 6: account_service.ChangePasswordTokenRequest
	(*ChangePasswordRequest)(nil),       // 7: account_service.ChangePasswordRequest
	(*GetAllSessionsRequest)(nil),       // 8: account_service.GetAllSessionsRequest
	(*TerminateSessionsRequest)(nil),    // 9: account_service.TerminateSessionsRequest
	(*Void)(nil),                        // 10: account_service.Void
	(*AccessResponce)(nil),              // 11: account_service.AccessResponce
	(*AccountID)(nil),                   // 12: account_service.AccountID
	(*ChangePasswordTokenResponce)(nil), // 13: account_service.ChangePasswordTokenResponce
	(*AllSessionsResponce)(nil),         // 14: account_service.AllSessionsResponce
}
var file_api_proto_depIdxs = []int32{
	0,  // 0: account_service.accountServiceV1.CreateAccount:input_type -> account_service.CreateAccountRequest
	1,  // 1: account_service.accountServiceV1.RequestAccountVerificationToken:input_type -> account_service.VerificationTokenRequest
	2,  // 2: account_service.accountServiceV1.VerifyAccount:input_type -> account_service.VerifyAccountRequest
	3,  // 3: account_service.accountServiceV1.SignIn:input_type -> account_service.SignInRequest
	4,  // 4: account_service.accountServiceV1.GetAccountID:input_type -> account_service.GetAccountIDRequest
	5,  // 5: account_service.accountServiceV1.Logout:input_type -> account_service.LogoutRequest
	6,  // 6: account_service.accountServiceV1.RequestChangePasswordToken:input_type -> account_service.ChangePasswordTokenRequest
	7,  // 7: account_service.accountServiceV1.ChangePassword:input_type -> account_service.ChangePasswordRequest
	8,  // 8: account_service.accountServiceV1.GetAllSessions:input_type -> account_service.GetAllSessionsRequest
	9,  // 9: account_service.accountServiceV1.TerminateSessions:input_type -> account_service.TerminateSessionsRequest
	10, // 10: account_service.accountServiceV1.CreateAccount:output_type -> account_service.Void
	10, // 11: account_service.accountServiceV1.RequestAccountVerificationToken:output_type -> account_service.Void
	10, // 12: account_service.accountServiceV1.VerifyAccount:output_type -> account_service.Void
	11, // 13: account_service.accountServiceV1.SignIn:output_type -> account_service.AccessResponce
	12, // 14: account_service.accountServiceV1.GetAccountID:output_type -> account_service.AccountID
	10, // 15: account_service.accountServiceV1.Logout:output_type -> account_service.Void
	13, // 16: account_service.accountServiceV1.RequestChangePasswordToken:output_type -> account_service.ChangePasswordTokenResponce
	10, // 17: account_service.accountServiceV1.ChangePassword:output_type -> account_service.Void
	14, // 18: account_service.accountServiceV1.GetAllSessions:output_type -> account_service.AllSessionsResponce
	10, // 19: account_service.accountServiceV1.TerminateSessions:output_type -> account_service.Void
	10, // [10:20] is the sub-list for method output_type
	0,  // [0:10] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_api_proto_init() }
func file_api_proto_init() {
	if File_api_proto != nil {
		return
	}
	file_messages_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_proto_goTypes,
		DependencyIndexes: file_api_proto_depIdxs,
	}.Build()
	File_api_proto = out.File
	file_api_proto_rawDesc = nil
	file_api_proto_goTypes = nil
	file_api_proto_depIdxs = nil
}
