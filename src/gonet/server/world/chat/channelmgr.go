package chat

import (
	"gonet/base"
	"gonet/message"
)

type (
	ChannelMgr struct {
		m_hmChannelMap	map[int64] *Channel
		m_hmPlayerChannelMap map[int8] map[int64] int64
	}

	IChannelMgr interface {
		Init()
		GetChannelId(int64) int64
		GetChannelIdByType(int64, int8) int64
		RegisterChannel(int8, string) int64
		UnregisterChannel(int64)

		AddPlayer(int64, int64, int64, string, int)
		RemovePlayer(int64, int64)
		RemoveAllChannel()
		BuildChannelID() int64

		getChannel(int64) *Channel
		getChannelByType(int64, int8) *Channel
	}
)

var (
	g_wordChannelId int64
)


func (this *ChannelMgr) getChannel(channelid int64)  *Channel{
	pChannel, exist := this.m_hmChannelMap[channelid]
	if exist{
		return pChannel
	}
	return nil
}

func (this *ChannelMgr) Init() {
	this.m_hmChannelMap	= make(map[int64] *Channel)
	this.m_hmPlayerChannelMap = make(map [int8] map[int64] int64)
	for i := message.CHAT_MSG_TYPE_WORLD; i < message.CHAT_MSG_TYPE_COUNT; i++{
		this.m_hmPlayerChannelMap[int8(i)] = make(map[int64] int64)
	}
	g_wordChannelId = this.RegisterChannel(int8(message.CHAT_MSG_TYPE_WORLD), "world")
}

func (this *ChannelMgr) getChannelByType(playerid int64, cMessageType int8)  *Channel{
	// 对于大规模消息来说，没有意义
	if cMessageType < int8(message.CHAT_MSG_TYPE_WORLD){
		return nil
	}

	channelid, exist := this.m_hmPlayerChannelMap[cMessageType][playerid]
	if !exist{
		return nil
	}

	return this.getChannel(channelid)
}

func (this *ChannelMgr) GetChannelIdByType(playerid int64, cMessageType int8) int64{
	pChannel := this.getChannelByType(playerid, cMessageType)
	if pChannel == nil{
		return 0
	}
	return pChannel.GetId()
}

func (this *ChannelMgr) BuildChannelID() int64{
	return base.UUID.UUID()
}

func (this *ChannelMgr) RegisterChannel(cMessageType int8, ChannelName string) int64{
	// 大规模消息不能创建频道
	if cMessageType < int8(message.CHAT_MSG_TYPE_WORLD) {
		return 0
	}

	nId := this.BuildChannelID()
	this.UnregisterChannel(nId)

	pChannel := &Channel{}
	pChannel.Init()
	pChannel.m_nChannelID = nId
	pChannel.m_cMessageType = cMessageType
	pChannel.m_strChannelName = ChannelName
	this.m_hmChannelMap[nId] = pChannel
	return nId
}

func (this *ChannelMgr) UnregisterChannel(channelid int64) {
	delete(this.m_hmChannelMap, channelid)
}

func (this *ChannelMgr) RemoveAllChannel() {
	for i,_ := range this.m_hmPlayerChannelMap{
		delete(this.m_hmPlayerChannelMap, i)
	}

	for i,_ := range this.m_hmChannelMap{
		delete(this.m_hmChannelMap, i)
	}
}

func (this *ChannelMgr) AddPlayer(accountid, playerid int64, channelid int64, playername string, socketId int) {
	pChannel := this.getChannel(channelid)
	if pChannel == nil{
		return
	}

	pChannel.AddPlayer(accountid, playerid, playername, socketId)
	this.m_hmPlayerChannelMap[pChannel.GetMessageType()][playerid] = pChannel.GetId()
}

func (this *ChannelMgr) RemovePlayer(playerid int64, channelid int64){
	pChannel := this.getChannel(channelid)
	if pChannel == nil{
		return
	}

	pChannel.RemovePlayer(playerid)
	delete(this.m_hmPlayerChannelMap[pChannel.GetMessageType()], playerid)
}


func (this *ChannelMgr) SendMessageToChannel(msg *ChatMessage, channelid int64){
	pChannel := this.getChannel(channelid)
	if pChannel == nil{
		return
	}

	pChannel.SendMessage(msg)
}



