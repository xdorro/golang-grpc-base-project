// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: api/proto/common/common.proto

package common_proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type UUIDRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The resource name of the book to be deleted, for example:
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *UUIDRequest) Reset() {
	*x = UUIDRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_common_common_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UUIDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UUIDRequest) ProtoMessage() {}

func (x *UUIDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_common_common_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UUIDRequest.ProtoReflect.Descriptor instead.
func (*UUIDRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_common_common_proto_rawDescGZIP(), []int{0}
}

func (x *UUIDRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type SlugRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The resource name of the book to be deleted, for example:
	Slug string `protobuf:"bytes,1,opt,name=slug,proto3" json:"slug,omitempty"`
}

func (x *SlugRequest) Reset() {
	*x = SlugRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_common_common_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SlugRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SlugRequest) ProtoMessage() {}

func (x *SlugRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_common_common_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SlugRequest.ProtoReflect.Descriptor instead.
func (*SlugRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_common_common_proto_rawDescGZIP(), []int{1}
}

func (x *SlugRequest) GetSlug() string {
	if x != nil {
		return x.Slug
	}
	return ""
}

var File_api_proto_common_common_proto protoreflect.FileDescriptor

var file_api_proto_common_common_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x1d, 0x0a,
	0x0b, 0x55, 0x55, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x21, 0x0a, 0x0b,
	0x53, 0x6c, 0x75, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x73,
	0x6c, 0x75, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x73, 0x6c, 0x75, 0x67, 0x42,
	0x4a, 0x5a, 0x48, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x78, 0x64,
	0x6f, 0x72, 0x72, 0x6f, 0x2f, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2d, 0x67, 0x72, 0x70, 0x63,
	0x2d, 0x62, 0x61, 0x73, 0x65, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x3b, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_api_proto_common_common_proto_rawDescOnce sync.Once
	file_api_proto_common_common_proto_rawDescData = file_api_proto_common_common_proto_rawDesc
)

func file_api_proto_common_common_proto_rawDescGZIP() []byte {
	file_api_proto_common_common_proto_rawDescOnce.Do(func() {
		file_api_proto_common_common_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_proto_common_common_proto_rawDescData)
	})
	return file_api_proto_common_common_proto_rawDescData
}

var file_api_proto_common_common_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_api_proto_common_common_proto_goTypes = []interface{}{
	(*UUIDRequest)(nil), // 0: common_proto.UUIDRequest
	(*SlugRequest)(nil), // 1: common_proto.SlugRequest
}
var file_api_proto_common_common_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_proto_common_common_proto_init() }
func file_api_proto_common_common_proto_init() {
	if File_api_proto_common_common_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_proto_common_common_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UUIDRequest); i {
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
		file_api_proto_common_common_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SlugRequest); i {
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
			RawDescriptor: file_api_proto_common_common_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_proto_common_common_proto_goTypes,
		DependencyIndexes: file_api_proto_common_common_proto_depIdxs,
		MessageInfos:      file_api_proto_common_common_proto_msgTypes,
	}.Build()
	File_api_proto_common_common_proto = out.File
	file_api_proto_common_common_proto_rawDesc = nil
	file_api_proto_common_common_proto_goTypes = nil
	file_api_proto_common_common_proto_depIdxs = nil
}
