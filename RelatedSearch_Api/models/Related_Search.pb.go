// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1-devel
// 	protoc        v3.20.1
// source: Related_Search.proto

package models

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

type RelatedSearchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Query string `protobuf:"bytes,1,opt,name=query,proto3" json:"query,omitempty"`
}

func (x *RelatedSearchRequest) Reset() {
	*x = RelatedSearchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Related_Search_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RelatedSearchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RelatedSearchRequest) ProtoMessage() {}

func (x *RelatedSearchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_Related_Search_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RelatedSearchRequest.ProtoReflect.Descriptor instead.
func (*RelatedSearchRequest) Descriptor() ([]byte, []int) {
	return file_Related_Search_proto_rawDescGZIP(), []int{0}
}

func (x *RelatedSearchRequest) GetQuery() string {
	if x != nil {
		return x.Query
	}
	return ""
}

type RelatedSearchResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Related []string `protobuf:"bytes,1,rep,name=related,proto3" json:"related,omitempty"`
}

func (x *RelatedSearchResponse) Reset() {
	*x = RelatedSearchResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Related_Search_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RelatedSearchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RelatedSearchResponse) ProtoMessage() {}

func (x *RelatedSearchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_Related_Search_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RelatedSearchResponse.ProtoReflect.Descriptor instead.
func (*RelatedSearchResponse) Descriptor() ([]byte, []int) {
	return file_Related_Search_proto_rawDescGZIP(), []int{1}
}

func (x *RelatedSearchResponse) GetRelated() []string {
	if x != nil {
		return x.Related
	}
	return nil
}

var File_Related_Search_proto protoreflect.FileDescriptor

var file_Related_Search_proto_rawDesc = []byte{
	0x0a, 0x14, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x22, 0x2c,
	0x0a, 0x14, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x65, 0x64, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x22, 0x31, 0x0a, 0x15,
	0x52, 0x65, 0x6c, 0x61, 0x74, 0x65, 0x64, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x65, 0x64,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x65, 0x64, 0x32,
	0x60, 0x0a, 0x0d, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x65, 0x64, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68,
	0x12, 0x4f, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x65, 0x64, 0x53, 0x65,
	0x61, 0x72, 0x63, 0x68, 0x12, 0x1c, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x52, 0x65,
	0x6c, 0x61, 0x74, 0x65, 0x64, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x52, 0x65, 0x6c, 0x61,
	0x74, 0x65, 0x64, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x42, 0x0e, 0x5a, 0x0c, 0x2e, 0x2e, 0x2f, 0x2e, 0x2e, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_Related_Search_proto_rawDescOnce sync.Once
	file_Related_Search_proto_rawDescData = file_Related_Search_proto_rawDesc
)

func file_Related_Search_proto_rawDescGZIP() []byte {
	file_Related_Search_proto_rawDescOnce.Do(func() {
		file_Related_Search_proto_rawDescData = protoimpl.X.CompressGZIP(file_Related_Search_proto_rawDescData)
	})
	return file_Related_Search_proto_rawDescData
}

var file_Related_Search_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_Related_Search_proto_goTypes = []interface{}{
	(*RelatedSearchRequest)(nil),  // 0: models.RelatedSearchRequest
	(*RelatedSearchResponse)(nil), // 1: models.RelatedSearchResponse
}
var file_Related_Search_proto_depIdxs = []int32{
	0, // 0: models.RelatedSearch.GetRelatedSearch:input_type -> models.RelatedSearchRequest
	1, // 1: models.RelatedSearch.GetRelatedSearch:output_type -> models.RelatedSearchResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_Related_Search_proto_init() }
func file_Related_Search_proto_init() {
	if File_Related_Search_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_Related_Search_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RelatedSearchRequest); i {
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
		file_Related_Search_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RelatedSearchResponse); i {
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
			RawDescriptor: file_Related_Search_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_Related_Search_proto_goTypes,
		DependencyIndexes: file_Related_Search_proto_depIdxs,
		MessageInfos:      file_Related_Search_proto_msgTypes,
	}.Build()
	File_Related_Search_proto = out.File
	file_Related_Search_proto_rawDesc = nil
	file_Related_Search_proto_goTypes = nil
	file_Related_Search_proto_depIdxs = nil
}