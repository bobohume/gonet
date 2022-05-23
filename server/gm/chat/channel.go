package chat

type (
	Channel struct {
		messageType int8   //消息类型
		channelID   int64  //ID
		channelName string //名称
		playerMap   map[int64]*player
	}

	IChannel interface {
		Init()
		GetId() int64
		GetMessageType() int8
		HasPlayer(playerId int64) bool
		GetPlayer(playerId int64) *player
		AddPlayer(playerId int64, playername string, gateClusterId uint32)
		RemovePlayer(playerId int64)
		SendMessage(msg *ChatMessage)
	}
)

func (c *Channel) Init() {
	c.playerMap = make(map[int64]*player)
}

func (c *Channel) GetId() int64 {
	return c.channelID
}

func (c *Channel) GetMessageType() int8 {
	return c.messageType
}

func (c *Channel) AddPlayer(playerId int64, playername string, gateClusterId uint32) {
	c.playerMap[playerId] = &player{playerId, playername, gateClusterId}
}

func (c *Channel) RemovePlayer(playerId int64) {
	delete(c.playerMap, playerId)
}

func (c *Channel) HasPlayer(playerId int64) bool {
	_, exist := c.playerMap[playerId]
	if exist {
		return true
	}
	return false
}

func (c *Channel) GetPlayer(playerId int64) *player {
	player, exist := c.playerMap[playerId]
	if exist {
		return player
	}
	return nil
}

func (c *Channel) SendMessage(msg *ChatMessage) {
	for _, v := range c.playerMap {
		SendMessage(msg, v)
	}
}
