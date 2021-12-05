// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: pkg/proto/v1/permission/permission.proto

package permissionproto

import (
	common "github.com/kucow/golang-grpc-base-project/pkg/proto/v1/common"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	status "google.golang.org/genproto/googleapis/rpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type FindAllPermissionsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Page int32 `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
}

func (x *FindAllPermissionsRequest) Reset() {
	*x = FindAllPermissionsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_v1_permission_permission_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindAllPermissionsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindAllPermissionsRequest) ProtoMessage() {}

func (x *FindAllPermissionsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_v1_permission_permission_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindAllPermissionsRequest.ProtoReflect.Descriptor instead.
func (*FindAllPermissionsRequest) Descriptor() ([]byte, []int) {
	return file_pkg_proto_v1_permission_permission_proto_rawDescGZIP(), []int{0}
}

func (x *FindAllPermissionsRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

type ListPermissionsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Permissions
	Data     []*Permission `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
	NextPage int32         `protobuf:"varint,2,opt,name=next_page,json=nextPage,proto3" json:"next_page,omitempty"`
}

func (x *ListPermissionsResponse) Reset() {
	*x = ListPermissionsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_v1_permission_permission_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListPermissionsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListPermissionsResponse) ProtoMessage() {}

func (x *ListPermissionsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_v1_permission_permission_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListPermissionsResponse.ProtoReflect.Descriptor instead.
func (*ListPermissionsResponse) Descriptor() ([]byte, []int) {
	return file_pkg_proto_v1_permission_permission_proto_rawDescGZIP(), []int{1}
}

func (x *ListPermissionsResponse) GetData() []*Permission {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *ListPermissionsResponse) GetNextPage() int32 {
	if x != nil {
		return x.NextPage
	}
	return 0
}

type Permission struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Output only. Id of the permission.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Output only. Name of the permission.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// Output only. Email of the permission.
	Slug string `protobuf:"bytes,3,opt,name=slug,proto3" json:"slug,omitempty"`
	// Output only. Status of the permission.
	Status int32 `protobuf:"varint,4,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *Permission) Reset() {
	*x = Permission{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_v1_permission_permission_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Permission) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Permission) ProtoMessage() {}

func (x *Permission) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_v1_permission_permission_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Permission.ProtoReflect.Descriptor instead.
func (*Permission) Descriptor() ([]byte, []int) {
	return file_pkg_proto_v1_permission_permission_proto_rawDescGZIP(), []int{2}
}

func (x *Permission) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Permission) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Permission) GetSlug() string {
	if x != nil {
		return x.Slug
	}
	return ""
}

func (x *Permission) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

// The request create new Permission
type CreatePermissionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Slug   string `protobuf:"bytes,2,opt,name=slug,proto3" json:"slug,omitempty"`
	Status int32  `protobuf:"varint,3,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *CreatePermissionRequest) Reset() {
	*x = CreatePermissionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_v1_permission_permission_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreatePermissionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreatePermissionRequest) ProtoMessage() {}

func (x *CreatePermissionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_v1_permission_permission_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreatePermissionRequest.ProtoReflect.Descriptor instead.
func (*CreatePermissionRequest) Descriptor() ([]byte, []int) {
	return file_pkg_proto_v1_permission_permission_proto_rawDescGZIP(), []int{3}
}

func (x *CreatePermissionRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreatePermissionRequest) GetSlug() string {
	if x != nil {
		return x.Slug
	}
	return ""
}

func (x *CreatePermissionRequest) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

// The request update Permission
type UpdatePermissionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name   *string `protobuf:"bytes,2,opt,name=name,proto3,oneof" json:"name,omitempty"`
	Slug   *string `protobuf:"bytes,3,opt,name=slug,proto3,oneof" json:"slug,omitempty"`
	Status *int32  `protobuf:"varint,4,opt,name=status,proto3,oneof" json:"status,omitempty"`
}

func (x *UpdatePermissionRequest) Reset() {
	*x = UpdatePermissionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_v1_permission_permission_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdatePermissionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdatePermissionRequest) ProtoMessage() {}

func (x *UpdatePermissionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_v1_permission_permission_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdatePermissionRequest.ProtoReflect.Descriptor instead.
func (*UpdatePermissionRequest) Descriptor() ([]byte, []int) {
	return file_pkg_proto_v1_permission_permission_proto_rawDescGZIP(), []int{4}
}

func (x *UpdatePermissionRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpdatePermissionRequest) GetName() string {
	if x != nil && x.Name != nil {
		return *x.Name
	}
	return ""
}

func (x *UpdatePermissionRequest) GetSlug() string {
	if x != nil && x.Slug != nil {
		return *x.Slug
	}
	return ""
}

func (x *UpdatePermissionRequest) GetStatus() int32 {
	if x != nil && x.Status != nil {
		return *x.Status
	}
	return 0
}

var File_pkg_proto_v1_permission_permission_proto protoreflect.FileDescriptor

var file_pkg_proto_v1_permission_permission_proto_rawDesc = []byte{
	0x0a, 0x28, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x2f, 0x70,
	0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2f, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x70, 0x65, 0x72, 0x6d,
	0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x20, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76,
	0x31, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2f, 0x0a, 0x19, 0x46, 0x69, 0x6e, 0x64, 0x41, 0x6c, 0x6c,
	0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x22, 0x67, 0x0a, 0x17, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x65,
	0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x2f, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x1b, 0x2e, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x12, 0x1b, 0x0a, 0x09, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x70, 0x61, 0x67, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x6e, 0x65, 0x78, 0x74, 0x50, 0x61, 0x67, 0x65, 0x22,
	0x5c, 0x0a, 0x0a, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x6c, 0x75, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x73, 0x6c, 0x75, 0x67, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x59, 0x0a,
	0x17, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x73, 0x6c, 0x75, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x73, 0x6c, 0x75, 0x67,
	0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x95, 0x01, 0x0a, 0x17, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x17, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x48, 0x00, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x17, 0x0a,
	0x04, 0x73, 0x6c, 0x75, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x04, 0x73,
	0x6c, 0x75, 0x67, 0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x48, 0x02, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x88, 0x01, 0x01, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x42, 0x07, 0x0a, 0x05,
	0x5f, 0x73, 0x6c, 0x75, 0x67, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x32, 0xcb, 0x04, 0x0a, 0x11, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x83, 0x01, 0x0a, 0x12, 0x46, 0x69, 0x6e, 0x64, 0x41,
	0x6c, 0x6c, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x2a, 0x2e,
	0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x46, 0x69, 0x6e, 0x64, 0x41, 0x6c, 0x6c, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x28, 0x2e, 0x70, 0x65, 0x72, 0x6d,
	0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x69, 0x73, 0x74,
	0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x17, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x11, 0x12, 0x0f, 0x2f, 0x76, 0x31,
	0x2f, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x6b, 0x0a, 0x12,
	0x46, 0x69, 0x6e, 0x64, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x42, 0x79,
	0x49, 0x44, 0x12, 0x18, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x55, 0x55, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x70,
	0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50,
	0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x1e, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x18, 0x12, 0x16, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x3d, 0x2a, 0x7d, 0x12, 0x6c, 0x0a, 0x10, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x28, 0x2e,
	0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x72, 0x70, 0x63, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x1a, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x14, 0x22, 0x0f, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x73, 0x3a, 0x01, 0x2a, 0x12, 0x73, 0x0a, 0x10, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x28, 0x2e, 0x70, 0x65,
	0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x72,
	0x70, 0x63, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x21, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x1b, 0x1a, 0x16, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x3d, 0x2a, 0x7d, 0x3a, 0x01, 0x2a, 0x12, 0x60, 0x0a, 0x10,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x12, 0x18, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x55,
	0x55, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x1e,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x18, 0x2a, 0x16, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x65, 0x72, 0x6d,
	0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x3d, 0x2a, 0x7d, 0x42, 0x53,
	0x5a, 0x51, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x75, 0x63,
	0x6f, 0x77, 0x2f, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2d, 0x67, 0x72, 0x70, 0x63, 0x2d, 0x62,
	0x61, 0x73, 0x65, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x70, 0x6b, 0x67, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x3b, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_proto_v1_permission_permission_proto_rawDescOnce sync.Once
	file_pkg_proto_v1_permission_permission_proto_rawDescData = file_pkg_proto_v1_permission_permission_proto_rawDesc
)

func file_pkg_proto_v1_permission_permission_proto_rawDescGZIP() []byte {
	file_pkg_proto_v1_permission_permission_proto_rawDescOnce.Do(func() {
		file_pkg_proto_v1_permission_permission_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_proto_v1_permission_permission_proto_rawDescData)
	})
	return file_pkg_proto_v1_permission_permission_proto_rawDescData
}

var file_pkg_proto_v1_permission_permission_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_pkg_proto_v1_permission_permission_proto_goTypes = []interface{}{
	(*FindAllPermissionsRequest)(nil), // 0: permissionproto.FindAllPermissionsRequest
	(*ListPermissionsResponse)(nil),   // 1: permissionproto.ListPermissionsResponse
	(*Permission)(nil),                // 2: permissionproto.Permission
	(*CreatePermissionRequest)(nil),   // 3: permissionproto.CreatePermissionRequest
	(*UpdatePermissionRequest)(nil),   // 4: permissionproto.UpdatePermissionRequest
	(*common.UUIDRequest)(nil),        // 5: commonproto.UUIDRequest
	(*status.Status)(nil),             // 6: google.rpc.Status
}
var file_pkg_proto_v1_permission_permission_proto_depIdxs = []int32{
	2, // 0: permissionproto.ListPermissionsResponse.data:type_name -> permissionproto.Permission
	0, // 1: permissionproto.PermissionService.FindAllPermissions:input_type -> permissionproto.FindAllPermissionsRequest
	5, // 2: permissionproto.PermissionService.FindPermissionByID:input_type -> commonproto.UUIDRequest
	3, // 3: permissionproto.PermissionService.CreatePermission:input_type -> permissionproto.CreatePermissionRequest
	4, // 4: permissionproto.PermissionService.UpdatePermission:input_type -> permissionproto.UpdatePermissionRequest
	5, // 5: permissionproto.PermissionService.DeletePermission:input_type -> commonproto.UUIDRequest
	1, // 6: permissionproto.PermissionService.FindAllPermissions:output_type -> permissionproto.ListPermissionsResponse
	2, // 7: permissionproto.PermissionService.FindPermissionByID:output_type -> permissionproto.Permission
	6, // 8: permissionproto.PermissionService.CreatePermission:output_type -> google.rpc.Status
	6, // 9: permissionproto.PermissionService.UpdatePermission:output_type -> google.rpc.Status
	6, // 10: permissionproto.PermissionService.DeletePermission:output_type -> google.rpc.Status
	6, // [6:11] is the sub-list for method output_type
	1, // [1:6] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_pkg_proto_v1_permission_permission_proto_init() }
func file_pkg_proto_v1_permission_permission_proto_init() {
	if File_pkg_proto_v1_permission_permission_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_proto_v1_permission_permission_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FindAllPermissionsRequest); i {
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
		file_pkg_proto_v1_permission_permission_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListPermissionsResponse); i {
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
		file_pkg_proto_v1_permission_permission_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Permission); i {
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
		file_pkg_proto_v1_permission_permission_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreatePermissionRequest); i {
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
		file_pkg_proto_v1_permission_permission_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdatePermissionRequest); i {
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
	file_pkg_proto_v1_permission_permission_proto_msgTypes[4].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pkg_proto_v1_permission_permission_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_proto_v1_permission_permission_proto_goTypes,
		DependencyIndexes: file_pkg_proto_v1_permission_permission_proto_depIdxs,
		MessageInfos:      file_pkg_proto_v1_permission_permission_proto_msgTypes,
	}.Build()
	File_pkg_proto_v1_permission_permission_proto = out.File
	file_pkg_proto_v1_permission_permission_proto_rawDesc = nil
	file_pkg_proto_v1_permission_permission_proto_goTypes = nil
	file_pkg_proto_v1_permission_permission_proto_depIdxs = nil
}