// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package placement

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

// PlacementClient is the client API for Placement service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PlacementClient interface {
	ReportAppStatus(ctx context.Context, opts ...grpc.CallOption) (Placement_ReportAppStatusClient, error)
}

type placementClient struct {
	cc grpc.ClientConnInterface
}

func NewPlacementClient(cc grpc.ClientConnInterface) PlacementClient {
	return &placementClient{cc}
}

func (c *placementClient) ReportAppStatus(ctx context.Context, opts ...grpc.CallOption) (Placement_ReportAppStatusClient, error) {
	stream, err := c.cc.NewStream(ctx, &Placement_ServiceDesc.Streams[0], "/v1.placement.Placement/ReportAppStatus", opts...)
	if err != nil {
		return nil, err
	}
	x := &placementReportAppStatusClient{stream}
	return x, nil
}

type Placement_ReportAppStatusClient interface {
	Send(*Host) error
	Recv() (*PlacementOrder, error)
	grpc.ClientStream
}

type placementReportAppStatusClient struct {
	grpc.ClientStream
}

func (x *placementReportAppStatusClient) Send(m *Host) error {
	return x.ClientStream.SendMsg(m)
}

func (x *placementReportAppStatusClient) Recv() (*PlacementOrder, error) {
	m := new(PlacementOrder)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PlacementServer is the server API for Placement service.
// All implementations should embed UnimplementedPlacementServer
// for forward compatibility
type PlacementServer interface {
	ReportAppStatus(Placement_ReportAppStatusServer) error
}

// UnimplementedPlacementServer should be embedded to have forward compatible implementations.
type UnimplementedPlacementServer struct {
}

func (UnimplementedPlacementServer) ReportAppStatus(Placement_ReportAppStatusServer) error {
	return status.Errorf(codes.Unimplemented, "method ReportAppStatus not implemented")
}

// UnsafePlacementServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PlacementServer will
// result in compilation errors.
type UnsafePlacementServer interface {
	mustEmbedUnimplementedPlacementServer()
}

func RegisterPlacementServer(s grpc.ServiceRegistrar, srv PlacementServer) {
	s.RegisterService(&Placement_ServiceDesc, srv)
}

func _Placement_ReportAppStatus_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(PlacementServer).ReportAppStatus(&placementReportAppStatusServer{stream})
}

type Placement_ReportAppStatusServer interface {
	Send(*PlacementOrder) error
	Recv() (*Host, error)
	grpc.ServerStream
}

type placementReportAppStatusServer struct {
	grpc.ServerStream
}

func (x *placementReportAppStatusServer) Send(m *PlacementOrder) error {
	return x.ServerStream.SendMsg(m)
}

func (x *placementReportAppStatusServer) Recv() (*Host, error) {
	m := new(Host)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Placement_ServiceDesc is the grpc.ServiceDesc for Placement service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Placement_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.placement.Placement",
	HandlerType: (*PlacementServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ReportAppStatus",
			Handler:       _Placement_ReportAppStatus_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "pkg/core/v1/placement/placement.proto",
}
