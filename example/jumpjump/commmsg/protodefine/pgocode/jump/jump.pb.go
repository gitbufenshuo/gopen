// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: jump.proto

package jump

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

type JumpMSGOne struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Kind     string `protobuf:"bytes,1,opt,name=Kind,proto3" json:"Kind,omitempty"`          //     string // move jump choose login doatt underatt
	Uid      string `protobuf:"bytes,2,opt,name=uid,proto3" json:"uid,omitempty"`            //  UID      string
	Which    int64  `protobuf:"varint,3,opt,name=which,proto3" json:"which,omitempty"`       // Which    int64 // 哪一个
	MoveValX int64  `protobuf:"varint,4,opt,name=moveValX,proto3" json:"moveValX,omitempty"` // MoveValX int64
	MoveValZ int64  `protobuf:"varint,5,opt,name=moveValZ,proto3" json:"moveValZ,omitempty"` // MoveValZ int64
	M        bool   `protobuf:"varint,6,opt,name=m,proto3" json:"m,omitempty"`               //M        bool
	PosX     int64  `protobuf:"varint,7,opt,name=PosX,proto3" json:"PosX,omitempty"`
	PosY     int64  `protobuf:"varint,8,opt,name=PosY,proto3" json:"PosY,omitempty"`
	PosZ     int64  `protobuf:"varint,9,opt,name=PosZ,proto3" json:"PosZ,omitempty"`
}

func (x *JumpMSGOne) Reset() {
	*x = JumpMSGOne{}
	if protoimpl.UnsafeEnabled {
		mi := &file_jump_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JumpMSGOne) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JumpMSGOne) ProtoMessage() {}

func (x *JumpMSGOne) ProtoReflect() protoreflect.Message {
	mi := &file_jump_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JumpMSGOne.ProtoReflect.Descriptor instead.
func (*JumpMSGOne) Descriptor() ([]byte, []int) {
	return file_jump_proto_rawDescGZIP(), []int{0}
}

func (x *JumpMSGOne) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *JumpMSGOne) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

func (x *JumpMSGOne) GetWhich() int64 {
	if x != nil {
		return x.Which
	}
	return 0
}

func (x *JumpMSGOne) GetMoveValX() int64 {
	if x != nil {
		return x.MoveValX
	}
	return 0
}

func (x *JumpMSGOne) GetMoveValZ() int64 {
	if x != nil {
		return x.MoveValZ
	}
	return 0
}

func (x *JumpMSGOne) GetM() bool {
	if x != nil {
		return x.M
	}
	return false
}

func (x *JumpMSGOne) GetPosX() int64 {
	if x != nil {
		return x.PosX
	}
	return 0
}

func (x *JumpMSGOne) GetPosY() int64 {
	if x != nil {
		return x.PosY
	}
	return 0
}

func (x *JumpMSGOne) GetPosZ() int64 {
	if x != nil {
		return x.PosZ
	}
	return 0
}

type JumpMSGTurn struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Turn int64         `protobuf:"varint,1,opt,name=turn,proto3" json:"turn,omitempty"` // Turn int64
	List []*JumpMSGOne `protobuf:"bytes,2,rep,name=list,proto3" json:"list,omitempty"`  // List []JumpMSGOne
}

func (x *JumpMSGTurn) Reset() {
	*x = JumpMSGTurn{}
	if protoimpl.UnsafeEnabled {
		mi := &file_jump_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JumpMSGTurn) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JumpMSGTurn) ProtoMessage() {}

func (x *JumpMSGTurn) ProtoReflect() protoreflect.Message {
	mi := &file_jump_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JumpMSGTurn.ProtoReflect.Descriptor instead.
func (*JumpMSGTurn) Descriptor() ([]byte, []int) {
	return file_jump_proto_rawDescGZIP(), []int{1}
}

func (x *JumpMSGTurn) GetTurn() int64 {
	if x != nil {
		return x.Turn
	}
	return 0
}

func (x *JumpMSGTurn) GetList() []*JumpMSGOne {
	if x != nil {
		return x.List
	}
	return nil
}

var File_jump_proto protoreflect.FileDescriptor

var file_jump_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x6a, 0x75, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xca, 0x01, 0x0a,
	0x0a, 0x4a, 0x75, 0x6d, 0x70, 0x4d, 0x53, 0x47, 0x4f, 0x6e, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x4b,
	0x69, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4b, 0x69, 0x6e, 0x64, 0x12,
	0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x69,
	0x64, 0x12, 0x14, 0x0a, 0x05, 0x77, 0x68, 0x69, 0x63, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x05, 0x77, 0x68, 0x69, 0x63, 0x68, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x6f, 0x76, 0x65, 0x56,
	0x61, 0x6c, 0x58, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x6d, 0x6f, 0x76, 0x65, 0x56,
	0x61, 0x6c, 0x58, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x6f, 0x76, 0x65, 0x56, 0x61, 0x6c, 0x5a, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x6d, 0x6f, 0x76, 0x65, 0x56, 0x61, 0x6c, 0x5a, 0x12,
	0x0c, 0x0a, 0x01, 0x6d, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x01, 0x6d, 0x12, 0x12, 0x0a,
	0x04, 0x50, 0x6f, 0x73, 0x58, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x50, 0x6f, 0x73,
	0x58, 0x12, 0x12, 0x0a, 0x04, 0x50, 0x6f, 0x73, 0x59, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x04, 0x50, 0x6f, 0x73, 0x59, 0x12, 0x12, 0x0a, 0x04, 0x50, 0x6f, 0x73, 0x5a, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x04, 0x50, 0x6f, 0x73, 0x5a, 0x22, 0x42, 0x0a, 0x0b, 0x4a, 0x75, 0x6d,
	0x70, 0x4d, 0x53, 0x47, 0x54, 0x75, 0x72, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x75, 0x72, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74, 0x75, 0x72, 0x6e, 0x12, 0x1f, 0x0a, 0x04,
	0x6c, 0x69, 0x73, 0x74, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x4a, 0x75, 0x6d,
	0x70, 0x4d, 0x53, 0x47, 0x4f, 0x6e, 0x65, 0x52, 0x04, 0x6c, 0x69, 0x73, 0x74, 0x42, 0x0e, 0x5a,
	0x0c, 0x70, 0x67, 0x6f, 0x63, 0x6f, 0x64, 0x65, 0x2f, 0x6a, 0x75, 0x6d, 0x70, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_jump_proto_rawDescOnce sync.Once
	file_jump_proto_rawDescData = file_jump_proto_rawDesc
)

func file_jump_proto_rawDescGZIP() []byte {
	file_jump_proto_rawDescOnce.Do(func() {
		file_jump_proto_rawDescData = protoimpl.X.CompressGZIP(file_jump_proto_rawDescData)
	})
	return file_jump_proto_rawDescData
}

var file_jump_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_jump_proto_goTypes = []interface{}{
	(*JumpMSGOne)(nil),  // 0: JumpMSGOne
	(*JumpMSGTurn)(nil), // 1: JumpMSGTurn
}
var file_jump_proto_depIdxs = []int32{
	0, // 0: JumpMSGTurn.list:type_name -> JumpMSGOne
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_jump_proto_init() }
func file_jump_proto_init() {
	if File_jump_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_jump_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JumpMSGOne); i {
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
		file_jump_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JumpMSGTurn); i {
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
			RawDescriptor: file_jump_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_jump_proto_goTypes,
		DependencyIndexes: file_jump_proto_depIdxs,
		MessageInfos:      file_jump_proto_msgTypes,
	}.Build()
	File_jump_proto = out.File
	file_jump_proto_rawDesc = nil
	file_jump_proto_goTypes = nil
	file_jump_proto_depIdxs = nil
}
