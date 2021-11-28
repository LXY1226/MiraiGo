// Code generated by protoc-gen-golite. DO NOT EDIT.
// source: pb/channel/MsgResponsesSvr.proto

package channel

type BatchGetMsgRspCountReq struct {
	GuildMsgList []*GuildMsg `protobuf:"bytes,1,rep"`
}

func (x *BatchGetMsgRspCountReq) GetGuildMsgList() []*GuildMsg {
	if x != nil {
		return x.GuildMsgList
	}
	return nil
}

type BatchGetMsgRspCountRsp struct {
	GuildMsgInfoList []*GuildMsgInfo `protobuf:"bytes,1,rep"`
}

func (x *BatchGetMsgRspCountRsp) GetGuildMsgInfoList() []*GuildMsgInfo {
	if x != nil {
		return x.GuildMsgInfoList
	}
	return nil
}

type SvrChannelMsg struct {
	ChannelId *uint64  `protobuf:"varint,1,opt"`
	Id        []*MsgId `protobuf:"bytes,2,rep"`
}

func (x *SvrChannelMsg) GetChannelId() uint64 {
	if x != nil && x.ChannelId != nil {
		return *x.ChannelId
	}
	return 0
}

func (x *SvrChannelMsg) GetId() []*MsgId {
	if x != nil {
		return x.Id
	}
	return nil
}

type ChannelMsgInfo struct {
	ChannelId *uint64        `protobuf:"varint,1,opt"`
	RespData  []*MsgRespData `protobuf:"bytes,2,rep"`
}

func (x *ChannelMsgInfo) GetChannelId() uint64 {
	if x != nil && x.ChannelId != nil {
		return *x.ChannelId
	}
	return 0
}

func (x *ChannelMsgInfo) GetRespData() []*MsgRespData {
	if x != nil {
		return x.RespData
	}
	return nil
}

type EmojiReaction struct {
	EmojiId        *string `protobuf:"bytes,1,opt"`
	EmojiType      *uint64 `protobuf:"varint,2,opt"`
	Cnt            *uint64 `protobuf:"varint,3,opt"`
	IsClicked      *bool   `protobuf:"varint,4,opt"`
	IsDefaultEmoji *bool   `protobuf:"varint,10001,opt"`
}

func (x *EmojiReaction) GetEmojiId() string {
	if x != nil && x.EmojiId != nil {
		return *x.EmojiId
	}
	return ""
}

func (x *EmojiReaction) GetEmojiType() uint64 {
	if x != nil && x.EmojiType != nil {
		return *x.EmojiType
	}
	return 0
}

func (x *EmojiReaction) GetCnt() uint64 {
	if x != nil && x.Cnt != nil {
		return *x.Cnt
	}
	return 0
}

func (x *EmojiReaction) GetIsClicked() bool {
	if x != nil && x.IsClicked != nil {
		return *x.IsClicked
	}
	return false
}

func (x *EmojiReaction) GetIsDefaultEmoji() bool {
	if x != nil && x.IsDefaultEmoji != nil {
		return *x.IsDefaultEmoji
	}
	return false
}

type GuildMsg struct {
	GuildId        *uint64          `protobuf:"varint,1,opt"`
	ChannelMsgList []*SvrChannelMsg `protobuf:"bytes,2,rep"`
}

func (x *GuildMsg) GetGuildId() uint64 {
	if x != nil && x.GuildId != nil {
		return *x.GuildId
	}
	return 0
}

func (x *GuildMsg) GetChannelMsgList() []*SvrChannelMsg {
	if x != nil {
		return x.ChannelMsgList
	}
	return nil
}

type GuildMsgInfo struct {
	GuildId            *uint64           `protobuf:"varint,1,opt"`
	ChannelMsgInfoList []*ChannelMsgInfo `protobuf:"bytes,2,rep"`
}

func (x *GuildMsgInfo) GetGuildId() uint64 {
	if x != nil && x.GuildId != nil {
		return *x.GuildId
	}
	return 0
}

func (x *GuildMsgInfo) GetChannelMsgInfoList() []*ChannelMsgInfo {
	if x != nil {
		return x.ChannelMsgInfoList
	}
	return nil
}

type MsgCnt struct {
	Id            *MsgId           `protobuf:"bytes,1,opt"`
	EmojiReaction []*EmojiReaction `protobuf:"bytes,2,rep"`
}

func (x *MsgCnt) GetId() *MsgId {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *MsgCnt) GetEmojiReaction() []*EmojiReaction {
	if x != nil {
		return x.EmojiReaction
	}
	return nil
}

type MsgId struct {
	Version *uint64 `protobuf:"varint,1,opt"`
	Seq     *uint64 `protobuf:"varint,2,opt"`
}

func (x *MsgId) GetVersion() uint64 {
	if x != nil && x.Version != nil {
		return *x.Version
	}
	return 0
}

func (x *MsgId) GetSeq() uint64 {
	if x != nil && x.Seq != nil {
		return *x.Seq
	}
	return 0
}

type MsgRespData struct {
	Id  *MsgId `protobuf:"bytes,1,opt"`
	Cnt []byte `protobuf:"bytes,2,opt"`
}

func (x *MsgRespData) GetId() *MsgId {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *MsgRespData) GetCnt() []byte {
	if x != nil {
		return x.Cnt
	}
	return nil
}
