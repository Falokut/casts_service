// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.24.3
// source: casts_service_v1_messages.proto

package protos

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

type GetCastRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CastID int32 `protobuf:"varint,1,opt,name=CastID,json=cast_id,proto3" json:"CastID,omitempty"`
	// use , as separator. All professions will be selected for the empty professionsIDs
	ProfessionsIDs string `protobuf:"bytes,2,opt,name=professionsIDs,json=professions_ids,proto3" json:"professionsIDs,omitempty"`
}

func (x *GetCastRequest) Reset() {
	*x = GetCastRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_casts_service_v1_messages_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCastRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCastRequest) ProtoMessage() {}

func (x *GetCastRequest) ProtoReflect() protoreflect.Message {
	mi := &file_casts_service_v1_messages_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCastRequest.ProtoReflect.Descriptor instead.
func (*GetCastRequest) Descriptor() ([]byte, []int) {
	return file_casts_service_v1_messages_proto_rawDescGZIP(), []int{0}
}

func (x *GetCastRequest) GetCastID() int32 {
	if x != nil {
		return x.CastID
	}
	return 0
}

func (x *GetCastRequest) GetProfessionsIDs() string {
	if x != nil {
		return x.ProfessionsIDs
	}
	return ""
}

type Actor struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID         int32       `protobuf:"varint,1,opt,name=ID,json=id,proto3" json:"ID,omitempty"`
	Profession *Profession `protobuf:"bytes,2,opt,name=profession,proto3" json:"profession,omitempty"`
}

func (x *Actor) Reset() {
	*x = Actor{}
	if protoimpl.UnsafeEnabled {
		mi := &file_casts_service_v1_messages_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Actor) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Actor) ProtoMessage() {}

func (x *Actor) ProtoReflect() protoreflect.Message {
	mi := &file_casts_service_v1_messages_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Actor.ProtoReflect.Descriptor instead.
func (*Actor) Descriptor() ([]byte, []int) {
	return file_casts_service_v1_messages_proto_rawDescGZIP(), []int{1}
}

func (x *Actor) GetID() int32 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *Actor) GetProfession() *Profession {
	if x != nil {
		return x.Profession
	}
	return nil
}

type Cast struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Actors []*Actor `protobuf:"bytes,1,rep,name=actors,proto3" json:"actors,omitempty"`
}

func (x *Cast) Reset() {
	*x = Cast{}
	if protoimpl.UnsafeEnabled {
		mi := &file_casts_service_v1_messages_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Cast) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Cast) ProtoMessage() {}

func (x *Cast) ProtoReflect() protoreflect.Message {
	mi := &file_casts_service_v1_messages_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Cast.ProtoReflect.Descriptor instead.
func (*Cast) Descriptor() ([]byte, []int) {
	return file_casts_service_v1_messages_proto_rawDescGZIP(), []int{2}
}

func (x *Cast) GetActors() []*Actor {
	if x != nil {
		return x.Actors
	}
	return nil
}

type Profession struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID   int32  `protobuf:"varint,1,opt,name=ID,json=id,proto3" json:"ID,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *Profession) Reset() {
	*x = Profession{}
	if protoimpl.UnsafeEnabled {
		mi := &file_casts_service_v1_messages_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Profession) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Profession) ProtoMessage() {}

func (x *Profession) ProtoReflect() protoreflect.Message {
	mi := &file_casts_service_v1_messages_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Profession.ProtoReflect.Descriptor instead.
func (*Profession) Descriptor() ([]byte, []int) {
	return file_casts_service_v1_messages_proto_rawDescGZIP(), []int{3}
}

func (x *Profession) GetID() int32 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *Profession) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type Professions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Professions []*Profession `protobuf:"bytes,1,rep,name=professions,proto3" json:"professions,omitempty"`
}

func (x *Professions) Reset() {
	*x = Professions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_casts_service_v1_messages_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Professions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Professions) ProtoMessage() {}

func (x *Professions) ProtoReflect() protoreflect.Message {
	mi := &file_casts_service_v1_messages_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Professions.ProtoReflect.Descriptor instead.
func (*Professions) Descriptor() ([]byte, []int) {
	return file_casts_service_v1_messages_proto_rawDescGZIP(), []int{4}
}

func (x *Professions) GetProfessions() []*Profession {
	if x != nil {
		return x.Professions
	}
	return nil
}

type UserErrorMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *UserErrorMessage) Reset() {
	*x = UserErrorMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_casts_service_v1_messages_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserErrorMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserErrorMessage) ProtoMessage() {}

func (x *UserErrorMessage) ProtoReflect() protoreflect.Message {
	mi := &file_casts_service_v1_messages_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserErrorMessage.ProtoReflect.Descriptor instead.
func (*UserErrorMessage) Descriptor() ([]byte, []int) {
	return file_casts_service_v1_messages_proto_rawDescGZIP(), []int{5}
}

func (x *UserErrorMessage) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_casts_service_v1_messages_proto protoreflect.FileDescriptor

var file_casts_service_v1_messages_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x63, 0x61, 0x73, 0x74, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f,
	0x76, 0x31, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x0d, 0x63, 0x61, 0x73, 0x74, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x22, 0x52, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x43, 0x61, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x17, 0x0a, 0x06, 0x43, 0x61, 0x73, 0x74, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x07, 0x63, 0x61, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x12, 0x27, 0x0a, 0x0e, 0x70,
	0x72, 0x6f, 0x66, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x49, 0x44, 0x73, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0f, 0x70, 0x72, 0x6f, 0x66, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73,
	0x5f, 0x69, 0x64, 0x73, 0x22, 0x52, 0x0a, 0x05, 0x41, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x0e, 0x0a,
	0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x39, 0x0a,
	0x0a, 0x70, 0x72, 0x6f, 0x66, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x19, 0x2e, 0x63, 0x61, 0x73, 0x74, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x50, 0x72, 0x6f, 0x66, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x70, 0x72,
	0x6f, 0x66, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x34, 0x0a, 0x04, 0x43, 0x61, 0x73, 0x74,
	0x12, 0x2c, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x14, 0x2e, 0x63, 0x61, 0x73, 0x74, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x41, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x06, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x22, 0x30,
	0x0a, 0x0a, 0x50, 0x72, 0x6f, 0x66, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02,
	0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x22, 0x4a, 0x0a, 0x0b, 0x50, 0x72, 0x6f, 0x66, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12,
	0x3b, 0x0a, 0x0b, 0x70, 0x72, 0x6f, 0x66, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x63, 0x61, 0x73, 0x74, 0x73, 0x5f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x50, 0x72, 0x6f, 0x66, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52,
	0x0b, 0x70, 0x72, 0x6f, 0x66, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x2c, 0x0a, 0x10,
	0x55, 0x73, 0x65, 0x72, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42, 0x19, 0x5a, 0x17, 0x63, 0x61,
	0x73, 0x74, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_casts_service_v1_messages_proto_rawDescOnce sync.Once
	file_casts_service_v1_messages_proto_rawDescData = file_casts_service_v1_messages_proto_rawDesc
)

func file_casts_service_v1_messages_proto_rawDescGZIP() []byte {
	file_casts_service_v1_messages_proto_rawDescOnce.Do(func() {
		file_casts_service_v1_messages_proto_rawDescData = protoimpl.X.CompressGZIP(file_casts_service_v1_messages_proto_rawDescData)
	})
	return file_casts_service_v1_messages_proto_rawDescData
}

var file_casts_service_v1_messages_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_casts_service_v1_messages_proto_goTypes = []interface{}{
	(*GetCastRequest)(nil),   // 0: casts_service.GetCastRequest
	(*Actor)(nil),            // 1: casts_service.Actor
	(*Cast)(nil),             // 2: casts_service.Cast
	(*Profession)(nil),       // 3: casts_service.Profession
	(*Professions)(nil),      // 4: casts_service.Professions
	(*UserErrorMessage)(nil), // 5: casts_service.UserErrorMessage
}
var file_casts_service_v1_messages_proto_depIdxs = []int32{
	3, // 0: casts_service.Actor.profession:type_name -> casts_service.Profession
	1, // 1: casts_service.Cast.actors:type_name -> casts_service.Actor
	3, // 2: casts_service.Professions.professions:type_name -> casts_service.Profession
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_casts_service_v1_messages_proto_init() }
func file_casts_service_v1_messages_proto_init() {
	if File_casts_service_v1_messages_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_casts_service_v1_messages_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCastRequest); i {
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
		file_casts_service_v1_messages_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Actor); i {
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
		file_casts_service_v1_messages_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Cast); i {
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
		file_casts_service_v1_messages_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Profession); i {
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
		file_casts_service_v1_messages_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Professions); i {
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
		file_casts_service_v1_messages_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserErrorMessage); i {
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
			RawDescriptor: file_casts_service_v1_messages_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_casts_service_v1_messages_proto_goTypes,
		DependencyIndexes: file_casts_service_v1_messages_proto_depIdxs,
		MessageInfos:      file_casts_service_v1_messages_proto_msgTypes,
	}.Build()
	File_casts_service_v1_messages_proto = out.File
	file_casts_service_v1_messages_proto_rawDesc = nil
	file_casts_service_v1_messages_proto_goTypes = nil
	file_casts_service_v1_messages_proto_depIdxs = nil
}
