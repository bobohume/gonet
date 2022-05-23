package chat

import (
	"gonet/actor"
	"gonet/base"
	"gonet/common/cluster"
	"gonet/rpc"
	"gonet/server/gm"
	"gonet/server/message"
	"time"

	"golang.org/x/net/context"
)

const (
	CHAT_PENDING_TIME_NORAML    = 1
	CHAT_PENDING_TIME_PRIVATE   = 1
	CHAT_PENDING_TIME_WORLDPLUS = 1
)

type (
	ChatMessage struct {
		Sender      int64
		Recver      int64
		MessageType int8
		Message     string
		SenderName  string
		RecverName  string
	}

	stPlayerChatRecord struct {
		lastTime    int64
		pendingTime int64
	}

	player struct {
		playerId      int64
		playerName    string
		gateClusterId uint32
	}

	ChatMgr struct {
		actor.Actor
		cluster.Stub
		channelManager ChannelMgr
		playerChatMap  map[int64]*stPlayerChatRecord
	}

	IChatMgr interface {
		actor.IActor

		SendMessageTo(msg *ChatMessage, playerId int64)
		SendMessageToAll(msg *ChatMessage)
		GetChannelManager() *ChannelMgr
		//setPlayerChatLastTime(int64, int8, int64)
		//getPlayerChatLastTime(int64, int8) int64
		//getPlayerChatPendingTime(int64, int8) int64
	}
)

var (
	MGR ChatMgr
)

func (c *ChatMgr) Init() {
	c.Actor.Init()
	c.playerChatMap = make(map[int64]*stPlayerChatRecord)
	c.channelManager.Init()
	actor.MGR.RegisterActor(c)
	c.Stub.InitStub(rpc.STUB_ChatMgr)
	c.Actor.Start()
}

func (c *ChatMgr) OnStubRegister(ctx context.Context) {
	//这里可以是加载db数据
	base.LOG.Println("Stub Chat register sucess")
}

func (c *ChatMgr) OnStubUnRegister(ctx context.Context) {
	//lease一致性这里要清理缓存数据了
	base.LOG.Println("Stub Chat unregister sucess")
}

func (c *ChatMgr) GetChannelManager() *ChannelMgr {
	return &c.channelManager
}

func (c *ChatMgr) SendMessageTo(msg *ChatMessage, playerId int64) {
	player := c.channelManager.getChannel(g_wordChannelId).GetPlayer(playerId)
	if player != nil {
		SendMessage(msg, player)
	}
}

func SendMessage(msg *ChatMessage, player *player) {
	gm.SendToClient(rpc.RpcHead{ClusterId: player.gateClusterId}, &message.ChatMessageResponse{
		PacketHead:  message.BuildPacketHead(player.playerId, rpc.SERVICE_GATE),
		Sender:      msg.Sender,
		SenderName:  msg.SenderName,
		Recver:      msg.Recver,
		RecverName:  msg.RecverName,
		MessageType: int32(msg.MessageType),
		Message:     msg.Message,
	})
}

func (c *ChatMgr) SendMessageToAll(msg *ChatMessage) {
	cluster.MGR.SendMsg(rpc.RpcHead{DestServerType: rpc.SERVICE_GATE, SendType: rpc.SEND_BOARD_CAST},
		"Chat_SendMessageAll", msg)
}

func (c *ChatMgr) setPlayerChatLastTime(playerid int64, cMessageType int8, nTime int64) {
	v := int64(playerid)
	v = (v << 8) | int64(cMessageType)

	c.playerChatMap[v] = &stPlayerChatRecord{}
	c.playerChatMap[v].lastTime = nTime
}

func (c *ChatMgr) getPlayerChatLastTime(playerid int64, cMessageType int8) int64 {
	v := int64(playerid)
	v = (v << 8) | int64(cMessageType)

	data, exist := c.playerChatMap[v]
	if exist {
		return data.lastTime
	}
	return 0
}

func (c *ChatMgr) getPlayerChatPendingTime(playerid int64, cMessageType int8) int64 {
	v := int64(playerid)
	v = (v << 8) | int64(cMessageType)

	data, exist := c.playerChatMap[v]
	if !exist {
		return 0
	}

	if data.pendingTime == 0 {
		switch cMessageType {
		case int8(message.CHAT_MSG_TYPE_PRIVATE):
			data.pendingTime = CHAT_PENDING_TIME_PRIVATE
		case int8(message.CHAT_MSG_TYPE_WORLD):
			data.pendingTime = CHAT_PENDING_TIME_WORLDPLUS
		default:
			data.pendingTime = CHAT_PENDING_TIME_NORAML
		}
	}

	return data.pendingTime
}

//聊天信息
func (c *ChatMgr) ChatMessageRequest(ctx context.Context, packet *message.ChatMessageRequest) {
	playerId := packet.GetSender()

	msg := &ChatMessage{}
	msg.Sender = playerId
	msg.SenderName = gm.SIMPLEMGR.GetPlayerName(msg.Sender)
	msg.Message = packet.GetMessage()
	msg.Recver = packet.GetRecver()
	msg.MessageType = int8(packet.GetMessageType())
	msg.RecverName = gm.SIMPLEMGR.GetPlayerName(msg.Recver)
	//替换屏蔽字库
	//data.ReplaceBanWord(msg.Message, "*")

	// 检查发送时间间隔
	pendingTime := c.getPlayerChatPendingTime(playerId, msg.MessageType)
	lastTime := c.getPlayerChatLastTime(playerId, msg.MessageType)
	nCurTime := time.Now().Unix()

	if nCurTime-lastTime < pendingTime {
		return
	}

	c.setPlayerChatLastTime(playerId, msg.MessageType, nCurTime)
	//writelog

	channelId := c.GetChannelManager().GetChannelIdByType(playerId, msg.MessageType)

	if msg.MessageType == int8(message.CHAT_MSG_TYPE_PRIVATE) && msg.Recver != msg.Sender { // 不能给自己发点对点消息
		c.SendMessageTo(msg, msg.Recver)
	} else if msg.MessageType == int8(message.CHAT_MSG_TYPE_WORLD) {
		//c.SendMessageToAll(msg)
		c.channelManager.SendMessageToChannel(msg, channelId)
	} else {
		if channelId == 0 {
			return
		}

		c.channelManager.SendMessageToChannel(msg, channelId)
	}
}

//注册频道
func (c *ChatMgr) RegisterChannel(ctx context.Context, messageType int8, channelId int64) {
	c.GetChannelManager().RegisterChannel(messageType, "", channelId)

	if 0 == channelId {
		return
	}

	if messageType == int8(message.CHAT_MSG_TYPE_ORG) {
	}
}

//销毁频道
func (c *ChatMgr) UnRegisterChannel(ctx context.Context, channelId int64) {
	c.GetChannelManager().UnregisterChannel(channelId)
}

//添加玩家到频道
func (c *ChatMgr) AddPlayerToChannel(ctx context.Context, playerId int64, channelId int64, playerName string, gateClusterId uint32) {
	c.GetChannelManager().AddPlayer(playerId, channelId, playerName, gateClusterId)
}

//删除玩家到频道
func (c *ChatMgr) RemovePlayerToChannel(ctx context.Context, playerId int64, channelId int64) {
	c.GetChannelManager().RemovePlayer(playerId, channelId)
}
