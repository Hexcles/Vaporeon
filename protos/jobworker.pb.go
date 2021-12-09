// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: protos/jobworker.proto

package protos

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type JobId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid string `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
}

func (x *JobId) Reset() {
	*x = JobId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_jobworker_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JobId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobId) ProtoMessage() {}

func (x *JobId) ProtoReflect() protoreflect.Message {
	mi := &file_protos_jobworker_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobId.ProtoReflect.Descriptor instead.
func (*JobId) Descriptor() ([]byte, []int) {
	return file_protos_jobworker_proto_rawDescGZIP(), []int{0}
}

func (x *JobId) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

type Job struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    *JobId `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`       // output-only
	Owner string `protobuf:"bytes,2,opt,name=owner,proto3" json:"owner,omitempty"` // output-only
	// Required: at least args[0] should be provided. The server will find the
	// binary to launch using the usual PATH lookup.
	Args     []string               `protobuf:"bytes,3,rep,name=args,proto3" json:"args,omitempty"`
	Started  *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=started,proto3" json:"started,omitempty"`                    // output-only
	Stopped  *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=stopped,proto3" json:"stopped,omitempty"`                    // output-only; running jobs do not have this field.
	ExitCode int32                  `protobuf:"varint,6,opt,name=exit_code,json=exitCode,proto3" json:"exit_code,omitempty"` // output-only; running jobs do not have this field.
}

func (x *Job) Reset() {
	*x = Job{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_jobworker_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Job) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Job) ProtoMessage() {}

func (x *Job) ProtoReflect() protoreflect.Message {
	mi := &file_protos_jobworker_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Job.ProtoReflect.Descriptor instead.
func (*Job) Descriptor() ([]byte, []int) {
	return file_protos_jobworker_proto_rawDescGZIP(), []int{1}
}

func (x *Job) GetId() *JobId {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *Job) GetOwner() string {
	if x != nil {
		return x.Owner
	}
	return ""
}

func (x *Job) GetArgs() []string {
	if x != nil {
		return x.Args
	}
	return nil
}

func (x *Job) GetStarted() *timestamppb.Timestamp {
	if x != nil {
		return x.Started
	}
	return nil
}

func (x *Job) GetStopped() *timestamppb.Timestamp {
	if x != nil {
		return x.Stopped
	}
	return nil
}

func (x *Job) GetExitCode() int32 {
	if x != nil {
		return x.ExitCode
	}
	return 0
}

type Output struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Return raw bytes directly. Production code should inform the client of the
	// locale to decode the bytes correctly.
	// Each message may contain one or both fields, but never none.
	Stdout []byte `protobuf:"bytes,1,opt,name=stdout,proto3" json:"stdout,omitempty"`
	Stderr []byte `protobuf:"bytes,2,opt,name=stderr,proto3" json:"stderr,omitempty"`
}

func (x *Output) Reset() {
	*x = Output{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_jobworker_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Output) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Output) ProtoMessage() {}

func (x *Output) ProtoReflect() protoreflect.Message {
	mi := &file_protos_jobworker_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Output.ProtoReflect.Descriptor instead.
func (*Output) Descriptor() ([]byte, []int) {
	return file_protos_jobworker_proto_rawDescGZIP(), []int{2}
}

func (x *Output) GetStdout() []byte {
	if x != nil {
		return x.Stdout
	}
	return nil
}

func (x *Output) GetStderr() []byte {
	if x != nil {
		return x.Stderr
	}
	return nil
}

var File_protos_jobworker_proto protoreflect.FileDescriptor

var file_protos_jobworker_proto_rawDesc = []byte{
	0x0a, 0x16, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x6a, 0x6f, 0x62, 0x77, 0x6f, 0x72, 0x6b,
	0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x6a, 0x6f, 0x62, 0x77, 0x6f, 0x72,
	0x6b, 0x65, 0x72, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x1b, 0x0a, 0x05, 0x4a, 0x6f, 0x62, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x22, 0xda,
	0x01, 0x0a, 0x03, 0x4a, 0x6f, 0x62, 0x12, 0x20, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6a, 0x6f, 0x62, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2e, 0x4a,
	0x6f, 0x62, 0x49, 0x64, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x6f, 0x77, 0x6e, 0x65,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x12, 0x12,
	0x0a, 0x04, 0x61, 0x72, 0x67, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x61, 0x72,
	0x67, 0x73, 0x12, 0x34, 0x0a, 0x07, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x07, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x12, 0x34, 0x0a, 0x07, 0x73, 0x74, 0x6f, 0x70,
	0x70, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x73, 0x74, 0x6f, 0x70, 0x70, 0x65, 0x64, 0x12, 0x1b,
	0x0a, 0x09, 0x65, 0x78, 0x69, 0x74, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x08, 0x65, 0x78, 0x69, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x22, 0x38, 0x0a, 0x06, 0x4f,
	0x75, 0x74, 0x70, 0x75, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x64, 0x6f, 0x75, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x73, 0x74, 0x64, 0x6f, 0x75, 0x74, 0x12, 0x16, 0x0a,
	0x06, 0x73, 0x74, 0x64, 0x65, 0x72, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x73,
	0x74, 0x64, 0x65, 0x72, 0x72, 0x32, 0x89, 0x02, 0x0a, 0x09, 0x4a, 0x6f, 0x62, 0x57, 0x6f, 0x72,
	0x6b, 0x65, 0x72, 0x12, 0x2c, 0x0a, 0x06, 0x4c, 0x61, 0x75, 0x6e, 0x63, 0x68, 0x12, 0x0e, 0x2e,
	0x6a, 0x6f, 0x62, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2e, 0x4a, 0x6f, 0x62, 0x1a, 0x10, 0x2e,
	0x6a, 0x6f, 0x62, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2e, 0x4a, 0x6f, 0x62, 0x49, 0x64, 0x22,
	0x00, 0x12, 0x2a, 0x0a, 0x04, 0x4b, 0x69, 0x6c, 0x6c, 0x12, 0x10, 0x2e, 0x6a, 0x6f, 0x62, 0x77,
	0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2e, 0x4a, 0x6f, 0x62, 0x49, 0x64, 0x1a, 0x0e, 0x2e, 0x6a, 0x6f,
	0x62, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2e, 0x4a, 0x6f, 0x62, 0x22, 0x00, 0x12, 0x2b, 0x0a,
	0x05, 0x51, 0x75, 0x65, 0x72, 0x79, 0x12, 0x10, 0x2e, 0x6a, 0x6f, 0x62, 0x77, 0x6f, 0x72, 0x6b,
	0x65, 0x72, 0x2e, 0x4a, 0x6f, 0x62, 0x49, 0x64, 0x1a, 0x0e, 0x2e, 0x6a, 0x6f, 0x62, 0x77, 0x6f,
	0x72, 0x6b, 0x65, 0x72, 0x2e, 0x4a, 0x6f, 0x62, 0x22, 0x00, 0x12, 0x37, 0x0a, 0x0c, 0x53, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x12, 0x10, 0x2e, 0x6a, 0x6f, 0x62,
	0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2e, 0x4a, 0x6f, 0x62, 0x49, 0x64, 0x1a, 0x11, 0x2e, 0x6a,
	0x6f, 0x62, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2e, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x22,
	0x00, 0x30, 0x01, 0x12, 0x3c, 0x0a, 0x08, 0x53, 0x68, 0x75, 0x74, 0x64, 0x6f, 0x77, 0x6e, 0x12,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22,
	0x00, 0x42, 0x24, 0x5a, 0x22, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x48, 0x65, 0x78, 0x63, 0x6c, 0x65, 0x73, 0x2f, 0x56, 0x61, 0x70, 0x6f, 0x72, 0x65, 0x6f, 0x6e,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protos_jobworker_proto_rawDescOnce sync.Once
	file_protos_jobworker_proto_rawDescData = file_protos_jobworker_proto_rawDesc
)

func file_protos_jobworker_proto_rawDescGZIP() []byte {
	file_protos_jobworker_proto_rawDescOnce.Do(func() {
		file_protos_jobworker_proto_rawDescData = protoimpl.X.CompressGZIP(file_protos_jobworker_proto_rawDescData)
	})
	return file_protos_jobworker_proto_rawDescData
}

var file_protos_jobworker_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_protos_jobworker_proto_goTypes = []interface{}{
	(*JobId)(nil),                 // 0: jobworker.JobId
	(*Job)(nil),                   // 1: jobworker.Job
	(*Output)(nil),                // 2: jobworker.Output
	(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),         // 4: google.protobuf.Empty
}
var file_protos_jobworker_proto_depIdxs = []int32{
	0, // 0: jobworker.Job.id:type_name -> jobworker.JobId
	3, // 1: jobworker.Job.started:type_name -> google.protobuf.Timestamp
	3, // 2: jobworker.Job.stopped:type_name -> google.protobuf.Timestamp
	1, // 3: jobworker.JobWorker.Launch:input_type -> jobworker.Job
	0, // 4: jobworker.JobWorker.Kill:input_type -> jobworker.JobId
	0, // 5: jobworker.JobWorker.Query:input_type -> jobworker.JobId
	0, // 6: jobworker.JobWorker.StreamOutput:input_type -> jobworker.JobId
	4, // 7: jobworker.JobWorker.Shutdown:input_type -> google.protobuf.Empty
	0, // 8: jobworker.JobWorker.Launch:output_type -> jobworker.JobId
	1, // 9: jobworker.JobWorker.Kill:output_type -> jobworker.Job
	1, // 10: jobworker.JobWorker.Query:output_type -> jobworker.Job
	2, // 11: jobworker.JobWorker.StreamOutput:output_type -> jobworker.Output
	4, // 12: jobworker.JobWorker.Shutdown:output_type -> google.protobuf.Empty
	8, // [8:13] is the sub-list for method output_type
	3, // [3:8] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_protos_jobworker_proto_init() }
func file_protos_jobworker_proto_init() {
	if File_protos_jobworker_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protos_jobworker_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JobId); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_jobworker_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Job); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_jobworker_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Output); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_protos_jobworker_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protos_jobworker_proto_goTypes,
		DependencyIndexes: file_protos_jobworker_proto_depIdxs,
		MessageInfos:      file_protos_jobworker_proto_msgTypes,
	}.Build()
	File_protos_jobworker_proto = out.File
	file_protos_jobworker_proto_rawDesc = nil
	file_protos_jobworker_proto_goTypes = nil
	file_protos_jobworker_proto_depIdxs = nil
}