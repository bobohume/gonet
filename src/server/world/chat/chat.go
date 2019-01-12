package chat

import (
	"actor"
	"github.com/golang/protobuf/proto"
	"message"
	"server/world"
	player2 "server/world/player"
	"time"
)

const(
	CHAT_PENDING_TIME_NORAML = 1
	CHAT_PENDING_TIME_PRIVATE = 1
	CHAT_PENDING_TIME_WORLDPLUS = 1
)

type(
	ChatMessage struct{
		Sender	int64
		Recver	int64
		MessageType int8
		Message string
		SenderName string
		RecverName string
	}

	stPlayerChatRecord struct {
		nLastTime int64
		nPendingTime int64
	}

	player struct{
		accountId int64
		playerId int64
		playerName string
		sockeId int
	}

	ChatMgr struct {
		actor.Actor
		m_channelManager ChannelMgr
		m_playerChatMap map[int64] *stPlayerChatRecord
	}

	IChatMgr interface {
		actor.IActor

		SendMessageTo(*ChatMessage, int64)
		SendMessageToChannel(*ChatMessage, int64)
		SendMessageToAll(*ChatMessage)
		GetChannelManager() *ChannelMgr

		setPlayerChatLastTime(int64, int8, int64)
		getPlayerChatLastTime(int64, int8) int64
		getPlayerChatPendingTime(int64, int8) int64
	}
)

var(
	CHATMGR ChatMgr
)

func (this *ChatMgr) Init(num int) {
	this.Actor.Init(num)

	this.m_playerChatMap = make(map[int64] *stPlayerChatRecord)
	actor.MGR().AddActor(this)
	this.m_channelManager.Init()
	//聊天信息
	this.RegisterCall("C_W_ChatMessage", func(packet *message.C_W_ChatMessage){
		playerId := *packet.Sender
		accountId := *packet.GetPacketHead().Id
		if accountId == 0{
			return
		}

		msg := &ChatMessage{}
		msg.Sender = playerId
		msg.Message = *packet.Message
		msg.Recver = *packet.Recver
		msg.MessageType = int8(*packet.MessageType)
		msg.RecverName = player2.PLAYERSIMPLEMGR.GetPlayerName(msg.Recver)
		//替换屏蔽字库
		//data.ReplaceBanWord(msg.Message, "*")

		// 检查发送时间间隔
		nPendingTime := this.getPlayerChatPendingTime(playerId, msg.MessageType)
		nLastTime := this.getPlayerChatLastTime(playerId, msg.MessageType)
		nCurTime := time.Now().Unix()

		if nCurTime - nLastTime < nPendingTime{
			return
		}

		this.setPlayerChatLastTime(playerId, msg.MessageType, nCurTime)
		//writelog

		channelId := this.GetChannelManager().GetChannelIdByType(playerId, msg.MessageType)

		if msg.MessageType == int8(message.CHAT_MSG_TYPE_PRIVATE) &&
			msg.Recver != msg.Sender{// 不能给自己发点对点消息
			this.SendMessageTo(msg, msg.Recver)
		}else if msg.MessageType == int8(message.CHAT_MSG_TYPE_WORLD){
			//this.SendMessageToAll(msg)
			this.m_channelManager.SendMessageToChannel(msg, channelId)
		}else{
			if channelId == 0 {
				return
			}

			this.m_channelManager.SendMessageToChannel(msg, channelId)
		}
	})

	//注册频道
	this.RegisterCall("RegisterChannel", func(messageType int8) {
		channelId := this.GetChannelManager().RegisterChannel(messageType, "")

		if 0 == channelId{
			return
		}

		if messageType == int8(message.CHAT_MSG_TYPE_ORG){
		}
	})

	//销毁频道
	this.RegisterCall("UnRegisterChannel", func(channelId int64) {
		this.GetChannelManager().UnregisterChannel(channelId)
	})

	//添加玩家到频道
	this.RegisterCall("AddPlayerToChannel", func(accoudId, playerId int64, channelId int64, playerName string, socketId int) {
		if channelId == -3000{
			channelId = g_wordChannelId
		}
		this.GetChannelManager().AddPlayer(accoudId, playerId, channelId, playerName, socketId)
	})

	//删除玩家到频道
	this.RegisterCall("RemovePlayerToChannel", func(playerId int64, channelId int64) {
		this.GetChannelManager().RemovePlayer(playerId, channelId)
	})
	
	this.Actor.Start()
}

func (this *ChatMgr) GetChannelManager() *ChannelMgr{
	return &this.m_channelManager
}

func (this *ChatMgr) SendMessageTo(chat *ChatMessage, playerId int64){

}

func SendMessage(chat *ChatMessage, player *player){
	world.SendToClient(player.sockeId, &message.W_C_ChatMessage{
		PacketHead:message.BuildPacketHead(player.accountId, int(message.SERVICE_CLIENT)),
		Sender:proto.Int64(chat.Sender),
		SenderName:proto.String(chat.SenderName),
		Recver:proto.Int64(chat.Recver),
		RecverName:proto.String(chat.RecverName),
		MessageType:proto.Int32(int32(chat.MessageType)),
		Message:proto.String(chat.Message),
	})
}

func (this *ChatMgr) SendMessageToAll(chat *ChatMessage){
	world.SERVER.GetServerMgr().SendMsg("Chat_SendMessageAll", )
}


func (this *ChatMgr) setPlayerChatLastTime(playerid int64, cMessageType int8, nTime int64){
	v := int64(playerid)
	v = (v << 8) | int64(cMessageType)

	this.m_playerChatMap[v] = &stPlayerChatRecord{}
	this.m_playerChatMap[v].nLastTime = nTime
}

func (this *ChatMgr) getPlayerChatLastTime(playerid int64, cMessageType int8) int64{
	v := int64(playerid)
	v = (v << 8) | int64(cMessageType)

	pData, exist := this.m_playerChatMap[v]
	if exist{
		return  pData.nLastTime
	}
	return 0
}

func (this *ChatMgr) getPlayerChatPendingTime(playerid int64, cMessageType int8) int64{
	v := int64(playerid)
	v = (v << 8) | int64(cMessageType)

	pData, exist := this.m_playerChatMap[v]
	if !exist{
		return 0
	}

	if pData.nPendingTime == 0{
		switch cMessageType {
		case int8(message.CHAT_MSG_TYPE_PRIVATE):
			pData.nPendingTime = CHAT_PENDING_TIME_PRIVATE
		case  int8(message.CHAT_MSG_TYPE_WORLD):
			pData.nPendingTime = CHAT_PENDING_TIME_WORLDPLUS
		default:
			pData.nPendingTime = CHAT_PENDING_TIME_NORAML
		}
	}

	return pData.nPendingTime
}