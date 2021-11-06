// 存放所有未知的结构体, 均为手动分析复原

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.14.0
// source: unknown.proto

package channel

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

type ChannelOidb0Xf5BRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GuildId   *uint64            `protobuf:"varint,1,opt,name=guildId" json:"guildId,omitempty"`
	Bots      []*GuildMemberInfo `protobuf:"bytes,4,rep,name=bots" json:"bots,omitempty"`
	Members   []*GuildMemberInfo `protobuf:"bytes,5,rep,name=members" json:"members,omitempty"`
	AdminInfo *GuildAdminInfo    `protobuf:"bytes,25,opt,name=adminInfo" json:"adminInfo,omitempty"`
}

func (x *ChannelOidb0Xf5BRsp) Reset() {
	*x = ChannelOidb0Xf5BRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_unknown_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChannelOidb0Xf5BRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChannelOidb0Xf5BRsp) ProtoMessage() {}

func (x *ChannelOidb0Xf5BRsp) ProtoReflect() protoreflect.Message {
	mi := &file_unknown_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChannelOidb0Xf5BRsp.ProtoReflect.Descriptor instead.
func (*ChannelOidb0Xf5BRsp) Descriptor() ([]byte, []int) {
	return file_unknown_proto_rawDescGZIP(), []int{0}
}

func (x *ChannelOidb0Xf5BRsp) GetGuildId() uint64 {
	if x != nil && x.GuildId != nil {
		return *x.GuildId
	}
	return 0
}

func (x *ChannelOidb0Xf5BRsp) GetBots() []*GuildMemberInfo {
	if x != nil {
		return x.Bots
	}
	return nil
}

func (x *ChannelOidb0Xf5BRsp) GetMembers() []*GuildMemberInfo {
	if x != nil {
		return x.Members
	}
	return nil
}

func (x *ChannelOidb0Xf5BRsp) GetAdminInfo() *GuildAdminInfo {
	if x != nil {
		return x.AdminInfo
	}
	return nil
}

type ChannelOidb0Xf88Rsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Profile *GuildUserProfile `protobuf:"bytes,1,opt,name=profile" json:"profile,omitempty"`
}

func (x *ChannelOidb0Xf88Rsp) Reset() {
	*x = ChannelOidb0Xf88Rsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_unknown_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChannelOidb0Xf88Rsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChannelOidb0Xf88Rsp) ProtoMessage() {}

func (x *ChannelOidb0Xf88Rsp) ProtoReflect() protoreflect.Message {
	mi := &file_unknown_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChannelOidb0Xf88Rsp.ProtoReflect.Descriptor instead.
func (*ChannelOidb0Xf88Rsp) Descriptor() ([]byte, []int) {
	return file_unknown_proto_rawDescGZIP(), []int{1}
}

func (x *ChannelOidb0Xf88Rsp) GetProfile() *GuildUserProfile {
	if x != nil {
		return x.Profile
	}
	return nil
}

type ChannelOidb0Xfc9Rsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Profile *GuildUserProfile `protobuf:"bytes,1,opt,name=profile" json:"profile,omitempty"`
}

func (x *ChannelOidb0Xfc9Rsp) Reset() {
	*x = ChannelOidb0Xfc9Rsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_unknown_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChannelOidb0Xfc9Rsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChannelOidb0Xfc9Rsp) ProtoMessage() {}

func (x *ChannelOidb0Xfc9Rsp) ProtoReflect() protoreflect.Message {
	mi := &file_unknown_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChannelOidb0Xfc9Rsp.ProtoReflect.Descriptor instead.
func (*ChannelOidb0Xfc9Rsp) Descriptor() ([]byte, []int) {
	return file_unknown_proto_rawDescGZIP(), []int{2}
}

func (x *ChannelOidb0Xfc9Rsp) GetProfile() *GuildUserProfile {
	if x != nil {
		return x.Profile
	}
	return nil
}

type ChannelOidb0Xf57Rsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Rsp *GuildMetaRsp `protobuf:"bytes,1,opt,name=rsp" json:"rsp,omitempty"`
}

func (x *ChannelOidb0Xf57Rsp) Reset() {
	*x = ChannelOidb0Xf57Rsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_unknown_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChannelOidb0Xf57Rsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChannelOidb0Xf57Rsp) ProtoMessage() {}

func (x *ChannelOidb0Xf57Rsp) ProtoReflect() protoreflect.Message {
	mi := &file_unknown_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChannelOidb0Xf57Rsp.ProtoReflect.Descriptor instead.
func (*ChannelOidb0Xf57Rsp) Descriptor() ([]byte, []int) {
	return file_unknown_proto_rawDescGZIP(), []int{3}
}

func (x *ChannelOidb0Xf57Rsp) GetRsp() *GuildMetaRsp {
	if x != nil {
		return x.Rsp
	}
	return nil
}

type GuildMetaRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GuildId *uint64    `protobuf:"varint,3,opt,name=guildId" json:"guildId,omitempty"`
	Meta    *GuildMeta `protobuf:"bytes,4,opt,name=meta" json:"meta,omitempty"`
}

func (x *GuildMetaRsp) Reset() {
	*x = GuildMetaRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_unknown_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GuildMetaRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GuildMetaRsp) ProtoMessage() {}

func (x *GuildMetaRsp) ProtoReflect() protoreflect.Message {
	mi := &file_unknown_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GuildMetaRsp.ProtoReflect.Descriptor instead.
func (*GuildMetaRsp) Descriptor() ([]byte, []int) {
	return file_unknown_proto_rawDescGZIP(), []int{4}
}

func (x *GuildMetaRsp) GetGuildId() uint64 {
	if x != nil && x.GuildId != nil {
		return *x.GuildId
	}
	return 0
}

func (x *GuildMetaRsp) GetMeta() *GuildMeta {
	if x != nil {
		return x.Meta
	}
	return nil
}

type GuildAdminInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Admins []*GuildMemberInfo `protobuf:"bytes,2,rep,name=admins" json:"admins,omitempty"`
}

func (x *GuildAdminInfo) Reset() {
	*x = GuildAdminInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_unknown_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GuildAdminInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GuildAdminInfo) ProtoMessage() {}

func (x *GuildAdminInfo) ProtoReflect() protoreflect.Message {
	mi := &file_unknown_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GuildAdminInfo.ProtoReflect.Descriptor instead.
func (*GuildAdminInfo) Descriptor() ([]byte, []int) {
	return file_unknown_proto_rawDescGZIP(), []int{5}
}

func (x *GuildAdminInfo) GetAdmins() []*GuildMemberInfo {
	if x != nil {
		return x.Admins
	}
	return nil
}

type GuildMemberInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Title         *string `protobuf:"bytes,2,opt,name=title" json:"title,omitempty"`
	Nickname      *string `protobuf:"bytes,3,opt,name=nickname" json:"nickname,omitempty"`
	LastSpeakTime *int64  `protobuf:"varint,4,opt,name=lastSpeakTime" json:"lastSpeakTime,omitempty"` // uncertainty
	Role          *int32  `protobuf:"varint,5,opt,name=role" json:"role,omitempty"`                   // uncertainty
	TinyId        *uint64 `protobuf:"varint,8,opt,name=tinyId" json:"tinyId,omitempty"`
}

func (x *GuildMemberInfo) Reset() {
	*x = GuildMemberInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_unknown_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GuildMemberInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GuildMemberInfo) ProtoMessage() {}

func (x *GuildMemberInfo) ProtoReflect() protoreflect.Message {
	mi := &file_unknown_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GuildMemberInfo.ProtoReflect.Descriptor instead.
func (*GuildMemberInfo) Descriptor() ([]byte, []int) {
	return file_unknown_proto_rawDescGZIP(), []int{6}
}

func (x *GuildMemberInfo) GetTitle() string {
	if x != nil && x.Title != nil {
		return *x.Title
	}
	return ""
}

func (x *GuildMemberInfo) GetNickname() string {
	if x != nil && x.Nickname != nil {
		return *x.Nickname
	}
	return ""
}

func (x *GuildMemberInfo) GetLastSpeakTime() int64 {
	if x != nil && x.LastSpeakTime != nil {
		return *x.LastSpeakTime
	}
	return 0
}

func (x *GuildMemberInfo) GetRole() int32 {
	if x != nil && x.Role != nil {
		return *x.Role
	}
	return 0
}

func (x *GuildMemberInfo) GetTinyId() uint64 {
	if x != nil && x.TinyId != nil {
		return *x.TinyId
	}
	return 0
}

// 频道系统用户资料
type GuildUserProfile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TinyId    *uint64 `protobuf:"varint,2,opt,name=tinyId" json:"tinyId,omitempty"`
	Nickname  *string `protobuf:"bytes,3,opt,name=nickname" json:"nickname,omitempty"`
	AvatarUrl *string `protobuf:"bytes,6,opt,name=avatarUrl" json:"avatarUrl,omitempty"`
	// 15: avatar url info
	JoinTime *int64 `protobuf:"varint,16,opt,name=joinTime" json:"joinTime,omitempty"` // uncertainty
}

func (x *GuildUserProfile) Reset() {
	*x = GuildUserProfile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_unknown_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GuildUserProfile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GuildUserProfile) ProtoMessage() {}

func (x *GuildUserProfile) ProtoReflect() protoreflect.Message {
	mi := &file_unknown_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GuildUserProfile.ProtoReflect.Descriptor instead.
func (*GuildUserProfile) Descriptor() ([]byte, []int) {
	return file_unknown_proto_rawDescGZIP(), []int{7}
}

func (x *GuildUserProfile) GetTinyId() uint64 {
	if x != nil && x.TinyId != nil {
		return *x.TinyId
	}
	return 0
}

func (x *GuildUserProfile) GetNickname() string {
	if x != nil && x.Nickname != nil {
		return *x.Nickname
	}
	return ""
}

func (x *GuildUserProfile) GetAvatarUrl() string {
	if x != nil && x.AvatarUrl != nil {
		return *x.AvatarUrl
	}
	return ""
}

func (x *GuildUserProfile) GetJoinTime() int64 {
	if x != nil && x.JoinTime != nil {
		return *x.JoinTime
	}
	return 0
}

type GuildMeta struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GuildCode      *uint64 `protobuf:"varint,2,opt,name=guildCode" json:"guildCode,omitempty"`
	CreateTime     *int64  `protobuf:"varint,4,opt,name=createTime" json:"createTime,omitempty"`
	MaxMemberCount *int64  `protobuf:"varint,5,opt,name=maxMemberCount" json:"maxMemberCount,omitempty"`
	MemberCount    *int64  `protobuf:"varint,6,opt,name=memberCount" json:"memberCount,omitempty"`
	Name           *string `protobuf:"bytes,8,opt,name=name" json:"name,omitempty"`
	RobotMaxNum    *int32  `protobuf:"varint,11,opt,name=robotMaxNum" json:"robotMaxNum,omitempty"`
	AdminMaxNum    *int32  `protobuf:"varint,12,opt,name=adminMaxNum" json:"adminMaxNum,omitempty"`
	Profile        *string `protobuf:"bytes,13,opt,name=profile" json:"profile,omitempty"`
	AvatarSeq      *int64  `protobuf:"varint,14,opt,name=avatarSeq" json:"avatarSeq,omitempty"`
	OwnerId        *uint64 `protobuf:"varint,18,opt,name=ownerId" json:"ownerId,omitempty"`
	CoverSeq       *int64  `protobuf:"varint,19,opt,name=coverSeq" json:"coverSeq,omitempty"`
	ClientId       *int32  `protobuf:"varint,20,opt,name=clientId" json:"clientId,omitempty"`
}

func (x *GuildMeta) Reset() {
	*x = GuildMeta{}
	if protoimpl.UnsafeEnabled {
		mi := &file_unknown_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GuildMeta) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GuildMeta) ProtoMessage() {}

func (x *GuildMeta) ProtoReflect() protoreflect.Message {
	mi := &file_unknown_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GuildMeta.ProtoReflect.Descriptor instead.
func (*GuildMeta) Descriptor() ([]byte, []int) {
	return file_unknown_proto_rawDescGZIP(), []int{8}
}

func (x *GuildMeta) GetGuildCode() uint64 {
	if x != nil && x.GuildCode != nil {
		return *x.GuildCode
	}
	return 0
}

func (x *GuildMeta) GetCreateTime() int64 {
	if x != nil && x.CreateTime != nil {
		return *x.CreateTime
	}
	return 0
}

func (x *GuildMeta) GetMaxMemberCount() int64 {
	if x != nil && x.MaxMemberCount != nil {
		return *x.MaxMemberCount
	}
	return 0
}

func (x *GuildMeta) GetMemberCount() int64 {
	if x != nil && x.MemberCount != nil {
		return *x.MemberCount
	}
	return 0
}

func (x *GuildMeta) GetName() string {
	if x != nil && x.Name != nil {
		return *x.Name
	}
	return ""
}

func (x *GuildMeta) GetRobotMaxNum() int32 {
	if x != nil && x.RobotMaxNum != nil {
		return *x.RobotMaxNum
	}
	return 0
}

func (x *GuildMeta) GetAdminMaxNum() int32 {
	if x != nil && x.AdminMaxNum != nil {
		return *x.AdminMaxNum
	}
	return 0
}

func (x *GuildMeta) GetProfile() string {
	if x != nil && x.Profile != nil {
		return *x.Profile
	}
	return ""
}

func (x *GuildMeta) GetAvatarSeq() int64 {
	if x != nil && x.AvatarSeq != nil {
		return *x.AvatarSeq
	}
	return 0
}

func (x *GuildMeta) GetOwnerId() uint64 {
	if x != nil && x.OwnerId != nil {
		return *x.OwnerId
	}
	return 0
}

func (x *GuildMeta) GetCoverSeq() int64 {
	if x != nil && x.CoverSeq != nil {
		return *x.CoverSeq
	}
	return 0
}

func (x *GuildMeta) GetClientId() int32 {
	if x != nil && x.ClientId != nil {
		return *x.ClientId
	}
	return 0
}

var File_unknown_proto protoreflect.FileDescriptor

var file_unknown_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x75, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x22, 0xc8, 0x01, 0x0a, 0x13, 0x43, 0x68, 0x61,
	0x6e, 0x6e, 0x65, 0x6c, 0x4f, 0x69, 0x64, 0x62, 0x30, 0x78, 0x66, 0x35, 0x62, 0x52, 0x73, 0x70,
	0x12, 0x18, 0x0a, 0x07, 0x67, 0x75, 0x69, 0x6c, 0x64, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x07, 0x67, 0x75, 0x69, 0x6c, 0x64, 0x49, 0x64, 0x12, 0x2c, 0x0a, 0x04, 0x62, 0x6f,
	0x74, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x63, 0x68, 0x61, 0x6e, 0x6e,
	0x65, 0x6c, 0x2e, 0x47, 0x75, 0x69, 0x6c, 0x64, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x49, 0x6e,
	0x66, 0x6f, 0x52, 0x04, 0x62, 0x6f, 0x74, 0x73, 0x12, 0x32, 0x0a, 0x07, 0x6d, 0x65, 0x6d, 0x62,
	0x65, 0x72, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x63, 0x68, 0x61, 0x6e,
	0x6e, 0x65, 0x6c, 0x2e, 0x47, 0x75, 0x69, 0x6c, 0x64, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x49,
	0x6e, 0x66, 0x6f, 0x52, 0x07, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x12, 0x35, 0x0a, 0x09,
	0x61, 0x64, 0x6d, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x19, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x17, 0x2e, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x2e, 0x47, 0x75, 0x69, 0x6c, 0x64, 0x41,
	0x64, 0x6d, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x09, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x49,
	0x6e, 0x66, 0x6f, 0x22, 0x4a, 0x0a, 0x13, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x4f, 0x69,
	0x64, 0x62, 0x30, 0x78, 0x66, 0x38, 0x38, 0x52, 0x73, 0x70, 0x12, 0x33, 0x0a, 0x07, 0x70, 0x72,
	0x6f, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x63, 0x68,
	0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x2e, 0x47, 0x75, 0x69, 0x6c, 0x64, 0x55, 0x73, 0x65, 0x72, 0x50,
	0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x22,
	0x4a, 0x0a, 0x13, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x4f, 0x69, 0x64, 0x62, 0x30, 0x78,
	0x66, 0x63, 0x39, 0x52, 0x73, 0x70, 0x12, 0x33, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65,
	0x6c, 0x2e, 0x47, 0x75, 0x69, 0x6c, 0x64, 0x55, 0x73, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x66, 0x69,
	0x6c, 0x65, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x22, 0x3e, 0x0a, 0x13, 0x43,
	0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x4f, 0x69, 0x64, 0x62, 0x30, 0x78, 0x66, 0x35, 0x37, 0x52,
	0x73, 0x70, 0x12, 0x27, 0x0a, 0x03, 0x72, 0x73, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x15, 0x2e, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x2e, 0x47, 0x75, 0x69, 0x6c, 0x64, 0x4d,
	0x65, 0x74, 0x61, 0x52, 0x73, 0x70, 0x52, 0x03, 0x72, 0x73, 0x70, 0x22, 0x50, 0x0a, 0x0c, 0x47,
	0x75, 0x69, 0x6c, 0x64, 0x4d, 0x65, 0x74, 0x61, 0x52, 0x73, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x67,
	0x75, 0x69, 0x6c, 0x64, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x67, 0x75,
	0x69, 0x6c, 0x64, 0x49, 0x64, 0x12, 0x26, 0x0a, 0x04, 0x6d, 0x65, 0x74, 0x61, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x2e, 0x47, 0x75,
	0x69, 0x6c, 0x64, 0x4d, 0x65, 0x74, 0x61, 0x52, 0x04, 0x6d, 0x65, 0x74, 0x61, 0x22, 0x42, 0x0a,
	0x0e, 0x47, 0x75, 0x69, 0x6c, 0x64, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x12,
	0x30, 0x0a, 0x06, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x18, 0x2e, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x2e, 0x47, 0x75, 0x69, 0x6c, 0x64, 0x4d,
	0x65, 0x6d, 0x62, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x06, 0x61, 0x64, 0x6d, 0x69, 0x6e,
	0x73, 0x22, 0x95, 0x01, 0x0a, 0x0f, 0x47, 0x75, 0x69, 0x6c, 0x64, 0x4d, 0x65, 0x6d, 0x62, 0x65,
	0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6e,
	0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6e,
	0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x24, 0x0a, 0x0d, 0x6c, 0x61, 0x73, 0x74, 0x53,
	0x70, 0x65, 0x61, 0x6b, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d,
	0x6c, 0x61, 0x73, 0x74, 0x53, 0x70, 0x65, 0x61, 0x6b, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x72, 0x6f, 0x6c,
	0x65, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x69, 0x6e, 0x79, 0x49, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x06, 0x74, 0x69, 0x6e, 0x79, 0x49, 0x64, 0x22, 0x80, 0x01, 0x0a, 0x10, 0x47, 0x75,
	0x69, 0x6c, 0x64, 0x55, 0x73, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x74, 0x69, 0x6e, 0x79, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06,
	0x74, 0x69, 0x6e, 0x79, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x6e, 0x69, 0x63, 0x6b, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6e, 0x69, 0x63, 0x6b, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x55, 0x72, 0x6c, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x55, 0x72, 0x6c,
	0x12, 0x1a, 0x0a, 0x08, 0x6a, 0x6f, 0x69, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x10, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x08, 0x6a, 0x6f, 0x69, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x22, 0xf5, 0x02, 0x0a,
	0x09, 0x47, 0x75, 0x69, 0x6c, 0x64, 0x4d, 0x65, 0x74, 0x61, 0x12, 0x1c, 0x0a, 0x09, 0x67, 0x75,
	0x69, 0x6c, 0x64, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x67,
	0x75, 0x69, 0x6c, 0x64, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x26, 0x0a, 0x0e, 0x6d, 0x61, 0x78, 0x4d,
	0x65, 0x6d, 0x62, 0x65, 0x72, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x0e, 0x6d, 0x61, 0x78, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x43, 0x6f, 0x75, 0x6e, 0x74,
	0x12, 0x20, 0x0a, 0x0b, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x72, 0x6f, 0x62, 0x6f, 0x74, 0x4d,
	0x61, 0x78, 0x4e, 0x75, 0x6d, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x72, 0x6f, 0x62,
	0x6f, 0x74, 0x4d, 0x61, 0x78, 0x4e, 0x75, 0x6d, 0x12, 0x20, 0x0a, 0x0b, 0x61, 0x64, 0x6d, 0x69,
	0x6e, 0x4d, 0x61, 0x78, 0x4e, 0x75, 0x6d, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x61,
	0x64, 0x6d, 0x69, 0x6e, 0x4d, 0x61, 0x78, 0x4e, 0x75, 0x6d, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x72,
	0x6f, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x72, 0x6f,
	0x66, 0x69, 0x6c, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x53, 0x65,
	0x71, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x53,
	0x65, 0x71, 0x12, 0x18, 0x0a, 0x07, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x49, 0x64, 0x18, 0x12, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x07, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08,
	0x63, 0x6f, 0x76, 0x65, 0x72, 0x53, 0x65, 0x71, 0x18, 0x13, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08,
	0x63, 0x6f, 0x76, 0x65, 0x72, 0x53, 0x65, 0x71, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x49, 0x64, 0x18, 0x14, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x49, 0x64, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x3b, 0x63, 0x68, 0x61, 0x6e, 0x6e,
	0x65, 0x6c,
}

var (
	file_unknown_proto_rawDescOnce sync.Once
	file_unknown_proto_rawDescData = file_unknown_proto_rawDesc
)

func file_unknown_proto_rawDescGZIP() []byte {
	file_unknown_proto_rawDescOnce.Do(func() {
		file_unknown_proto_rawDescData = protoimpl.X.CompressGZIP(file_unknown_proto_rawDescData)
	})
	return file_unknown_proto_rawDescData
}

var file_unknown_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_unknown_proto_goTypes = []interface{}{
	(*ChannelOidb0Xf5BRsp)(nil), // 0: channel.ChannelOidb0xf5bRsp
	(*ChannelOidb0Xf88Rsp)(nil), // 1: channel.ChannelOidb0xf88Rsp
	(*ChannelOidb0Xfc9Rsp)(nil), // 2: channel.ChannelOidb0xfc9Rsp
	(*ChannelOidb0Xf57Rsp)(nil), // 3: channel.ChannelOidb0xf57Rsp
	(*GuildMetaRsp)(nil),        // 4: channel.GuildMetaRsp
	(*GuildAdminInfo)(nil),      // 5: channel.GuildAdminInfo
	(*GuildMemberInfo)(nil),     // 6: channel.GuildMemberInfo
	(*GuildUserProfile)(nil),    // 7: channel.GuildUserProfile
	(*GuildMeta)(nil),           // 8: channel.GuildMeta
}
var file_unknown_proto_depIdxs = []int32{
	6, // 0: channel.ChannelOidb0xf5bRsp.bots:type_name -> channel.GuildMemberInfo
	6, // 1: channel.ChannelOidb0xf5bRsp.members:type_name -> channel.GuildMemberInfo
	5, // 2: channel.ChannelOidb0xf5bRsp.adminInfo:type_name -> channel.GuildAdminInfo
	7, // 3: channel.ChannelOidb0xf88Rsp.profile:type_name -> channel.GuildUserProfile
	7, // 4: channel.ChannelOidb0xfc9Rsp.profile:type_name -> channel.GuildUserProfile
	4, // 5: channel.ChannelOidb0xf57Rsp.rsp:type_name -> channel.GuildMetaRsp
	8, // 6: channel.GuildMetaRsp.meta:type_name -> channel.GuildMeta
	6, // 7: channel.GuildAdminInfo.admins:type_name -> channel.GuildMemberInfo
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_unknown_proto_init() }
func file_unknown_proto_init() {
	if File_unknown_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_unknown_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChannelOidb0Xf5BRsp); i {
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
		file_unknown_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChannelOidb0Xf88Rsp); i {
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
		file_unknown_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChannelOidb0Xfc9Rsp); i {
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
		file_unknown_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChannelOidb0Xf57Rsp); i {
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
		file_unknown_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GuildMetaRsp); i {
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
		file_unknown_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GuildAdminInfo); i {
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
		file_unknown_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GuildMemberInfo); i {
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
		file_unknown_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GuildUserProfile); i {
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
		file_unknown_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GuildMeta); i {
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
			RawDescriptor: file_unknown_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_unknown_proto_goTypes,
		DependencyIndexes: file_unknown_proto_depIdxs,
		MessageInfos:      file_unknown_proto_msgTypes,
	}.Build()
	File_unknown_proto = out.File
	file_unknown_proto_rawDesc = nil
	file_unknown_proto_goTypes = nil
	file_unknown_proto_depIdxs = nil
}
