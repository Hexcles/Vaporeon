// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package protos

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// JobWorkerClient is the client API for JobWorker service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type JobWorkerClient interface {
	Launch(ctx context.Context, in *Job, opts ...grpc.CallOption) (*JobId, error)
	// For simplicity, SIGKILL is always used. It blocks until the job exits.
	Kill(ctx context.Context, in *JobId, opts ...grpc.CallOption) (*Job, error)
	Query(ctx context.Context, in *JobId, opts ...grpc.CallOption) (*Job, error)
	// Stream closes when the job stops and all output has been sent.
	StreamOutput(ctx context.Context, in *JobId, opts ...grpc.CallOption) (JobWorker_StreamOutputClient, error)
	// Admin-only: kill all jobs and shutdown the server.
	Shutdown(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type jobWorkerClient struct {
	cc grpc.ClientConnInterface
}

func NewJobWorkerClient(cc grpc.ClientConnInterface) JobWorkerClient {
	return &jobWorkerClient{cc}
}

func (c *jobWorkerClient) Launch(ctx context.Context, in *Job, opts ...grpc.CallOption) (*JobId, error) {
	out := new(JobId)
	err := c.cc.Invoke(ctx, "/jobworker.JobWorker/Launch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jobWorkerClient) Kill(ctx context.Context, in *JobId, opts ...grpc.CallOption) (*Job, error) {
	out := new(Job)
	err := c.cc.Invoke(ctx, "/jobworker.JobWorker/Kill", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jobWorkerClient) Query(ctx context.Context, in *JobId, opts ...grpc.CallOption) (*Job, error) {
	out := new(Job)
	err := c.cc.Invoke(ctx, "/jobworker.JobWorker/Query", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jobWorkerClient) StreamOutput(ctx context.Context, in *JobId, opts ...grpc.CallOption) (JobWorker_StreamOutputClient, error) {
	stream, err := c.cc.NewStream(ctx, &JobWorker_ServiceDesc.Streams[0], "/jobworker.JobWorker/StreamOutput", opts...)
	if err != nil {
		return nil, err
	}
	x := &jobWorkerStreamOutputClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type JobWorker_StreamOutputClient interface {
	Recv() (*Output, error)
	grpc.ClientStream
}

type jobWorkerStreamOutputClient struct {
	grpc.ClientStream
}

func (x *jobWorkerStreamOutputClient) Recv() (*Output, error) {
	m := new(Output)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *jobWorkerClient) Shutdown(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/jobworker.JobWorker/Shutdown", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// JobWorkerServer is the server API for JobWorker service.
// All implementations must embed UnimplementedJobWorkerServer
// for forward compatibility
type JobWorkerServer interface {
	Launch(context.Context, *Job) (*JobId, error)
	// For simplicity, SIGKILL is always used. It blocks until the job exits.
	Kill(context.Context, *JobId) (*Job, error)
	Query(context.Context, *JobId) (*Job, error)
	// Stream closes when the job stops and all output has been sent.
	StreamOutput(*JobId, JobWorker_StreamOutputServer) error
	// Admin-only: kill all jobs and shutdown the server.
	Shutdown(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	mustEmbedUnimplementedJobWorkerServer()
}

// UnimplementedJobWorkerServer must be embedded to have forward compatible implementations.
type UnimplementedJobWorkerServer struct {
}

func (UnimplementedJobWorkerServer) Launch(context.Context, *Job) (*JobId, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Launch not implemented")
}
func (UnimplementedJobWorkerServer) Kill(context.Context, *JobId) (*Job, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Kill not implemented")
}
func (UnimplementedJobWorkerServer) Query(context.Context, *JobId) (*Job, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Query not implemented")
}
func (UnimplementedJobWorkerServer) StreamOutput(*JobId, JobWorker_StreamOutputServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamOutput not implemented")
}
func (UnimplementedJobWorkerServer) Shutdown(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Shutdown not implemented")
}
func (UnimplementedJobWorkerServer) mustEmbedUnimplementedJobWorkerServer() {}

// UnsafeJobWorkerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to JobWorkerServer will
// result in compilation errors.
type UnsafeJobWorkerServer interface {
	mustEmbedUnimplementedJobWorkerServer()
}

func RegisterJobWorkerServer(s grpc.ServiceRegistrar, srv JobWorkerServer) {
	s.RegisterService(&JobWorker_ServiceDesc, srv)
}

func _JobWorker_Launch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Job)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobWorkerServer).Launch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/jobworker.JobWorker/Launch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobWorkerServer).Launch(ctx, req.(*Job))
	}
	return interceptor(ctx, in, info, handler)
}

func _JobWorker_Kill_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JobId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobWorkerServer).Kill(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/jobworker.JobWorker/Kill",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobWorkerServer).Kill(ctx, req.(*JobId))
	}
	return interceptor(ctx, in, info, handler)
}

func _JobWorker_Query_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JobId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobWorkerServer).Query(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/jobworker.JobWorker/Query",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobWorkerServer).Query(ctx, req.(*JobId))
	}
	return interceptor(ctx, in, info, handler)
}

func _JobWorker_StreamOutput_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(JobId)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(JobWorkerServer).StreamOutput(m, &jobWorkerStreamOutputServer{stream})
}

type JobWorker_StreamOutputServer interface {
	Send(*Output) error
	grpc.ServerStream
}

type jobWorkerStreamOutputServer struct {
	grpc.ServerStream
}

func (x *jobWorkerStreamOutputServer) Send(m *Output) error {
	return x.ServerStream.SendMsg(m)
}

func _JobWorker_Shutdown_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobWorkerServer).Shutdown(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/jobworker.JobWorker/Shutdown",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobWorkerServer).Shutdown(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// JobWorker_ServiceDesc is the grpc.ServiceDesc for JobWorker service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var JobWorker_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "jobworker.JobWorker",
	HandlerType: (*JobWorkerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Launch",
			Handler:    _JobWorker_Launch_Handler,
		},
		{
			MethodName: "Kill",
			Handler:    _JobWorker_Kill_Handler,
		},
		{
			MethodName: "Query",
			Handler:    _JobWorker_Query_Handler,
		},
		{
			MethodName: "Shutdown",
			Handler:    _JobWorker_Shutdown_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamOutput",
			Handler:       _JobWorker_StreamOutput_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "protos/jobworker.proto",
}
