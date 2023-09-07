// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023 Datadog, Inc.

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.4
// source: fixtures_test.proto

package grpc

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

const (
	Fixture_Ping_FullMethodName       = "/grpc.Fixture/Ping"
	Fixture_StreamPing_FullMethodName = "/grpc.Fixture/StreamPing"
)

// FixtureClient is the client API for Fixture service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FixtureClient interface {
	Ping(ctx context.Context, in *FixtureRequest, opts ...grpc.CallOption) (*FixtureReply, error)
	StreamPing(ctx context.Context, opts ...grpc.CallOption) (Fixture_StreamPingClient, error)
}

type fixtureClient struct {
	cc grpc.ClientConnInterface
}

func NewFixtureClient(cc grpc.ClientConnInterface) FixtureClient {
	return &fixtureClient{cc}
}

func (c *fixtureClient) Ping(ctx context.Context, in *FixtureRequest, opts ...grpc.CallOption) (*FixtureReply, error) {
	out := new(FixtureReply)
	err := c.cc.Invoke(ctx, Fixture_Ping_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fixtureClient) StreamPing(ctx context.Context, opts ...grpc.CallOption) (Fixture_StreamPingClient, error) {
	stream, err := c.cc.NewStream(ctx, &Fixture_ServiceDesc.Streams[0], Fixture_StreamPing_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &fixtureStreamPingClient{stream}
	return x, nil
}

type Fixture_StreamPingClient interface {
	Send(*FixtureRequest) error
	Recv() (*FixtureReply, error)
	grpc.ClientStream
}

type fixtureStreamPingClient struct {
	grpc.ClientStream
}

func (x *fixtureStreamPingClient) Send(m *FixtureRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *fixtureStreamPingClient) Recv() (*FixtureReply, error) {
	m := new(FixtureReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// FixtureServer is the server API for Fixture service.
// All implementations must embed UnimplementedFixtureServer
// for forward compatibility
type FixtureServer interface {
	Ping(context.Context, *FixtureRequest) (*FixtureReply, error)
	StreamPing(Fixture_StreamPingServer) error
	mustEmbedUnimplementedFixtureServer()
}

// UnimplementedFixtureServer must be embedded to have forward compatible implementations.
type UnimplementedFixtureServer struct {
}

func (UnimplementedFixtureServer) Ping(context.Context, *FixtureRequest) (*FixtureReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedFixtureServer) StreamPing(Fixture_StreamPingServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamPing not implemented")
}
func (UnimplementedFixtureServer) mustEmbedUnimplementedFixtureServer() {}

// UnsafeFixtureServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FixtureServer will
// result in compilation errors.
type UnsafeFixtureServer interface {
	mustEmbedUnimplementedFixtureServer()
}

func RegisterFixtureServer(s grpc.ServiceRegistrar, srv FixtureServer) {
	s.RegisterService(&Fixture_ServiceDesc, srv)
}

func _Fixture_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FixtureRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FixtureServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Fixture_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FixtureServer).Ping(ctx, req.(*FixtureRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Fixture_StreamPing_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FixtureServer).StreamPing(&fixtureStreamPingServer{stream})
}

type Fixture_StreamPingServer interface {
	Send(*FixtureReply) error
	Recv() (*FixtureRequest, error)
	grpc.ServerStream
}

type fixtureStreamPingServer struct {
	grpc.ServerStream
}

func (x *fixtureStreamPingServer) Send(m *FixtureReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *fixtureStreamPingServer) Recv() (*FixtureRequest, error) {
	m := new(FixtureRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Fixture_ServiceDesc is the grpc.ServiceDesc for Fixture service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Fixture_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.Fixture",
	HandlerType: (*FixtureServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Fixture_Ping_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamPing",
			Handler:       _Fixture_StreamPing_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "fixtures_test.proto",
}
