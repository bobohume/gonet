package chat

import (
	"actor"
	"message"
	player2 "server/world/player"
	"time"
	"server/world"
)

const(
	CHAT_MSG_WORLD	 = iota
	CHAT_MSG_TYPE_ORG	 =iota

	CHAT_PENDING_TIME_NORAML = 3000
	CHAT_PENDING_TIME_PRIVATE = 1000
	CHAT_PENDING_TIME_WORLDPLUS = 30000
)

type(
	ChatMessage struct{
		Sender	int
		Recver	int
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
		playerId int
		playerName string
	}

	ChatMgr struct {
		actor.Actor
		m_channelManager ChannelMgr
		m_playerChatMap map[int64] *stPlayerChatRecord
	}

	IChatMgr interface {
		actor.IActor

		SendMessageTo(*ChatMessage, int)
		SendMessageToChannel(*ChatMessage, int)
		SendMessage(*ChatMessage, []int)
		SendMessageToAll(*ChatMessage)
		GetChannelManager() *ChannelMgr

		setPlayerChatLastTime(int, int8, int64)
		getPlayerChatLastTime(int, int8) int64
		getPlayerChatPendingTime(int, int8) int64
	}
)

var(
	CHATMGR ChatMgr
)

func (this *ChatMgr) Init(num int) {
	this.Actor.Init(num)

	this.m_playerChatMap = make(map[int64] *stPlayerChatRecord)

	//聊天信息
	this.RegisterCall("C_W_ChatMessage", func(caller *actor.Caller, packet *message.C_W_ChatMessage){
		playerId := int(*packet.Sender)
		accountId := int(*packet.GetPacketHead().Id)
		if accountId == 0{
			return
		}

		msg := &ChatMessage{}
		msg.Sender = playerId
		msg.Message = *packet.Message
		msg.Recver = int(*packet.Recver)
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
			this.SendMessageToAll(msg)
		}else{
			if channelId == 0 {
				return
			}

			this.SendMessageToChannel(msg, channelId)
		}
	})

	//注册频道
	this.RegisterCall("RegisterChannel", func(caller *actor.Caller, messageType int8) {
		channelId := this.GetChannelManager().RegisterChannel(messageType, "")

		if 0 == channelId{
			return
		}

		if messageType == int8(message.CHAT_MSG_TYPE_ORG){
		}
	})

	//销毁频道
	this.RegisterCall("UnRegisterChannel", func(caller *actor.Caller, channelId int) {
		this.GetChannelManager().UnregisterChannel(channelId)
	})

	//添加玩家到频道
	this.RegisterCall("AddPlayerToChannel", func(caller *actor.Caller, playerId, channelId int, playerName string) {
		this.GetChannelManager().AddPlayer(playerId, channelId, playerName)
	})

	//删除玩家到频道
	this.RegisterCall("RemovePlayerToChannel", func(caller *actor.Caller, playerId, channelId int) {
		this.GetChannelManager().RemovePlayer(playerId, channelId)
	})
	
	this.Actor.Start()
}

func (this *ChatMgr) GetChannelManager() *ChannelMgr{
	return &this.m_channelManager
}

func (this *ChatMgr) SendMessageTo(chat *ChatMessage, playerId int){

}

func (this *ChatMgr) SendMessageToChannel(chat *ChatMessage, channelid int){
	if channelid == 0{
		return
	}

	playerList := this.GetChannelManager().GetPlayerList(channelid)
	this.SendMessage(chat, playerList)
}

func (this *ChatMgr) SendMessage(chat *ChatMessage, playerList []int){

}

func (this *ChatMgr) SendMessageToAll(chat *ChatMessage){
	world.SERVER.GetServerMgr().SendMsg(0, "Chat_SendMessageAll", )
}


func (this *ChatMgr) setPlayerChatLastTime(playerid int, cMessageType int8, nTime int64){
	v := int64(playerid)
	v = (v << 8) | int64(cMessageType)

	this.m_playerChatMap[v].nLastTime = nTime
}

func (this *ChatMgr) getPlayerChatLastTime(playerid int, cMessageType int8) int64{
	v := int64(playerid)
	v = (v << 8) | int64(cMessageType)

	pData, exist := this.m_playerChatMap[v]
	if exist{
		return  pData.nLastTime
	}
	return 0
}

func (this *ChatMgr) getPlayerChatPendingTime(playerid int, cMessageType int8) int64{
	v := int64(playerid)
	v = (v << 8) | int64(cMessageType)

	pData := this.m_playerChatMap[v]
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