package chat

import (
	"gonet/server/message"
)

type (
	ChannelMgr struct {
		channelMap       map[int64]*Channel
		playerChannelMap map[int8]map[int64]int64
	}

	IChannelMgr interface {
		Init()
		RegisterChannel(cMessageType int8, ChannelName string, nId int64)
		UnregisterChannel(channelid int64)

		AddPlayer(layerId int64, channelId int64, playerName string, gateClusterId uint32)
		RemovePlayer(playerId int64, channelId int64)
		RemoveAllChannel()
		//getChannel(channelid int64)  *Channel
		//getChannelByType(playerid int64, cMessageType int8)  *Channel
	}
)

var (
	g_wordChannelId = int64(-3000)
)

func (c *ChannelMgr) Init() {
	c.channelMap = make(map[int64]*Channel)
	c.playerChannelMap = make(map[int8]map[int64]int64)
	for i := message.CHAT_MSG_TYPE_WORLD; i < message.CHAT_MSG_TYPE_COUNT; i++ {
		c.playerChannelMap[int8(i)] = make(map[int64]int64)
	}

	c.RegisterChannel(int8(message.CHAT_MSG_TYPE_WORLD), "game", g_wordChannelId)
}

func (c *ChannelMgr) RegisterChannel(cMessageType int8, ChannelName string, nId int64) {
	// 大规模消息不能创建频道
	if cMessageType < int8(message.CHAT_MSG_TYPE_WORLD) {
		return
	}

	c.UnregisterChannel(nId)

	channel := &Channel{}
	channel.Init()
	channel.channelID = nId
	channel.messageType = cMessageType
	channel.channelName = ChannelName
	c.channelMap[nId] = channel
	return
}

func (c *ChannelMgr) UnregisterChannel(channelid int64) {
	delete(c.channelMap, channelid)
}

func (c *ChannelMgr) RemoveAllChannel() {
	for i, _ := range c.playerChannelMap {
		delete(c.playerChannelMap, i)
	}

	for i, _ := range c.channelMap {
		delete(c.channelMap, i)
	}
}

func (c *ChannelMgr) AddPlayer(playerId int64, channelId int64, playerName string, gateClusterId uint32) {
	channel := c.getChannel(channelId)
	if channel == nil {
		return
	}

	channel.AddPlayer(playerId, playerName, gateClusterId)
	c.playerChannelMap[channel.GetMessageType()][playerId] = channel.GetId()
}

func (c *ChannelMgr) RemovePlayer(playerid int64, channelid int64) {
	channel := c.getChannel(channelid)
	if channel == nil {
		return
	}

	channel.RemovePlayer(playerid)
	delete(c.playerChannelMap[channel.GetMessageType()], playerid)
}

func (c *ChannelMgr) GetChannelIdByType(playerid int64, cMessageType int8) int64 {
	channel := c.getChannelByType(playerid, cMessageType)
	if channel == nil {
		return 0
	}
	return channel.GetId()
}

func (c *ChannelMgr) getChannel(channelid int64) *Channel {
	channel, exist := c.channelMap[channelid]
	if exist {
		return channel
	}
	return nil
}

func (c *ChannelMgr) getChannelByType(playerid int64, cMessageType int8) *Channel {
	// 对于大规模消息来说，没有意义
	if cMessageType < int8(message.CHAT_MSG_TYPE_WORLD) {
		return nil
	}

	channelid, exist := c.playerChannelMap[cMessageType][playerid]
	if !exist {
		return nil
	}

	return c.getChannel(channelid)
}

func (c *ChannelMgr) SendMessageToChannel(msg *ChatMessage, channelid int64) {
	channel := c.getChannel(channelid)
	if channel == nil {
		return
	}

	channel.SendMessage(msg)
}
