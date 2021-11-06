package client

import (
	"fmt"
	"runtime/debug"
	"sync/atomic"

	"github.com/Mrs4s/MiraiGo/message"
)

type EventHandler struct {
	PrivateMessageHandler           func(*QQClient, *message.PrivateMessage)
	TempMessageHandler              func(*QQClient, *TempMessageEvent)
	GroupMessageHandler             func(*QQClient, *message.GroupMessage)
	SelfPrivateMessageHandler       func(*QQClient, *message.PrivateMessage)
	SelfGroupMessageHandler         func(*QQClient, *message.GroupMessage)
	GroupMuteEventHandler           func(*QQClient, *GroupMuteEvent)
	GroupRecalledHandler            func(*QQClient, *GroupMessageRecalledEvent)
	FriendRecalledHandler           func(*QQClient, *FriendMessageRecalledEvent)
	JoinGroupHandler                func(*QQClient, *GroupInfo)
	LeaveGroupHandler               func(*QQClient, *GroupLeaveEvent)
	MemberJoinedHandler             func(*QQClient, *MemberJoinGroupEvent)
	MemberLeavedHandler             func(*QQClient, *MemberLeaveGroupEvent)
	MemberCardUpdatedHandler        func(*QQClient, *MemberCardUpdatedEvent)
	GroupNameUpdatedHandler         func(*QQClient, *GroupNameUpdatedEvent)
	PermissionChangedHandler        func(*QQClient, *MemberPermissionChangedEvent)
	GroupInvitedHandler             func(*QQClient, *GroupInvitedRequest)
	JoinRequestHandler              func(*QQClient, *UserJoinGroupRequest)
	FriendRequestHandler            func(*QQClient, *NewFriendRequest)
	NewFriendHandler                func(*QQClient, *NewFriendEvent)
	DisconnectHandler               func(*QQClient, *ClientDisconnectedEvent)
	LogHandler                      func(*QQClient, *LogEvent)
	ServerUpdatedHandler            func(*QQClient, *ServerUpdatedEvent) bool
	GroupNotifyHandler              func(*QQClient, INotifyEvent)
	FriendNotifyHandler             func(*QQClient, INotifyEvent)
	MemberTitleUpdatedHandler       func(*QQClient, *MemberSpecialTitleUpdatedEvent)
	OfflineFileHandler              func(*QQClient, *OfflineFileEvent)
	OtherClientStatusChangedHandler func(*QQClient, *OtherClientStatusChangedEvent)
	GroupDigestHandler              func(*QQClient, *GroupDigestEvent)
}

var nopHandlers = EventHandler{
	PrivateMessageHandler:           func(*QQClient, *message.PrivateMessage) {},
	TempMessageHandler:              func(*QQClient, *TempMessageEvent) {},
	GroupMessageHandler:             func(*QQClient, *message.GroupMessage) {},
	SelfPrivateMessageHandler:       func(*QQClient, *message.PrivateMessage) {},
	SelfGroupMessageHandler:         func(*QQClient, *message.GroupMessage) {},
	GroupMuteEventHandler:           func(*QQClient, *GroupMuteEvent) {},
	GroupRecalledHandler:            func(*QQClient, *GroupMessageRecalledEvent) {},
	FriendRecalledHandler:           func(*QQClient, *FriendMessageRecalledEvent) {},
	JoinGroupHandler:                func(*QQClient, *GroupInfo) {},
	LeaveGroupHandler:               func(*QQClient, *GroupLeaveEvent) {},
	MemberJoinedHandler:             func(*QQClient, *MemberJoinGroupEvent) {},
	MemberLeavedHandler:             func(*QQClient, *MemberLeaveGroupEvent) {},
	MemberCardUpdatedHandler:        func(*QQClient, *MemberCardUpdatedEvent) {},
	GroupNameUpdatedHandler:         func(*QQClient, *GroupNameUpdatedEvent) {},
	PermissionChangedHandler:        func(*QQClient, *MemberPermissionChangedEvent) {},
	GroupInvitedHandler:             func(*QQClient, *GroupInvitedRequest) {},
	JoinRequestHandler:              func(*QQClient, *UserJoinGroupRequest) {},
	FriendRequestHandler:            func(*QQClient, *NewFriendRequest) {},
	NewFriendHandler:                func(*QQClient, *NewFriendEvent) {},
	DisconnectHandler:               func(*QQClient, *ClientDisconnectedEvent) {},
	LogHandler:                      func(*QQClient, *LogEvent) {},
	ServerUpdatedHandler:            func(*QQClient, *ServerUpdatedEvent) bool { return true },
	GroupNotifyHandler:              func(*QQClient, INotifyEvent) {},
	FriendNotifyHandler:             func(*QQClient, INotifyEvent) {},
	MemberTitleUpdatedHandler:       func(*QQClient, *MemberSpecialTitleUpdatedEvent) {},
	OfflineFileHandler:              func(*QQClient, *OfflineFileEvent) {},
	OtherClientStatusChangedHandler: func(*QQClient, *OtherClientStatusChangedEvent) {},
	GroupDigestHandler:              func(*QQClient, *GroupDigestEvent) {},
}

func (c *QQClient) dispatchLogEvent(e *LogEvent) {
	if e == nil {
		return
	}
	for _, f := range c.eventHandler.logHandler {
		cover(func() {
			f(c, e)
		})
	}
}

func (c *QQClient) onGroupMessageReceipt(id string, f ...func(*QQClient, *groupMessageReceiptEvent)) {
	if len(f) == 0 {
		c.groupMessageReceiptHandler.Delete(id)
		return
	}
	atomic.AddUint64(&c.stat.MessageSent, 1)
	c.groupMessageReceiptHandler.LoadOrStore(id, f[0])
}

func (c *QQClient) dispatchGroupMessageReceiptEvent(e *groupMessageReceiptEvent) {
	c.groupMessageReceiptHandler.Range(func(_, f interface{}) bool {
		go f.(func(*QQClient, *groupMessageReceiptEvent))(c, e)
		return true
	})
}

func cover(f func()) {
	defer func() {
		if pan := recover(); pan != nil {
			fmt.Printf("event error: %v\n%s", pan, debug.Stack())
		}
	}()
	f()
}
