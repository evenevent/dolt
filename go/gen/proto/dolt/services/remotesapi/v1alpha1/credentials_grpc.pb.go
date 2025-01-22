// Copyright 2019 Dolthub, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.28.3
// source: dolt/services/remotesapi/v1alpha1/credentials.proto

package remotesapi

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

// CredentialsServiceClient is the client API for CredentialsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CredentialsServiceClient interface {
	WhoAmI(ctx context.Context, in *WhoAmIRequest, opts ...grpc.CallOption) (*WhoAmIResponse, error)
}

type credentialsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCredentialsServiceClient(cc grpc.ClientConnInterface) CredentialsServiceClient {
	return &credentialsServiceClient{cc}
}

func (c *credentialsServiceClient) WhoAmI(ctx context.Context, in *WhoAmIRequest, opts ...grpc.CallOption) (*WhoAmIResponse, error) {
	out := new(WhoAmIResponse)
	err := c.cc.Invoke(ctx, "/dolt.services.remotesapi.v1alpha1.CredentialsService/WhoAmI", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CredentialsServiceServer is the server API for CredentialsService service.
// All implementations must embed UnimplementedCredentialsServiceServer
// for forward compatibility
type CredentialsServiceServer interface {
	WhoAmI(context.Context, *WhoAmIRequest) (*WhoAmIResponse, error)
	mustEmbedUnimplementedCredentialsServiceServer()
}

// UnimplementedCredentialsServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCredentialsServiceServer struct {
}

func (UnimplementedCredentialsServiceServer) WhoAmI(context.Context, *WhoAmIRequest) (*WhoAmIResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WhoAmI not implemented")
}
func (UnimplementedCredentialsServiceServer) mustEmbedUnimplementedCredentialsServiceServer() {}

// UnsafeCredentialsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CredentialsServiceServer will
// result in compilation errors.
type UnsafeCredentialsServiceServer interface {
	mustEmbedUnimplementedCredentialsServiceServer()
}

func RegisterCredentialsServiceServer(s grpc.ServiceRegistrar, srv CredentialsServiceServer) {
	s.RegisterService(&CredentialsService_ServiceDesc, srv)
}

func _CredentialsService_WhoAmI_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WhoAmIRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CredentialsServiceServer).WhoAmI(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dolt.services.remotesapi.v1alpha1.CredentialsService/WhoAmI",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CredentialsServiceServer).WhoAmI(ctx, req.(*WhoAmIRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CredentialsService_ServiceDesc is the grpc.ServiceDesc for CredentialsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CredentialsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "dolt.services.remotesapi.v1alpha1.CredentialsService",
	HandlerType: (*CredentialsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "WhoAmI",
			Handler:    _CredentialsService_WhoAmI_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dolt/services/remotesapi/v1alpha1/credentials.proto",
}
