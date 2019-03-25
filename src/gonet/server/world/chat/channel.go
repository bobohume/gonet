package chat

type(
	Channel struct{
		m_cMessageType int8	//消息类型
		m_nChannelID int64	//ID
		m_strChannelName string //名称
		m_playerMap map[int64] *player
	}

	IChannel interface {
		Init()
		GetId() int64
		GetMessageType() int8
		HasPlayer(int64) bool
		AddPlayer(int64, int64, string, int)
		RemovePlayer(int64)
		SendMessage(*ChatMessage)
	}
)


func (this *Channel) Init(){
	this.m_playerMap = make(map[int64] *player)
}

func (this *Channel) GetId() int64{
	return this.m_nChannelID
}

func (this *Channel) GetMessageType() int8{
	return this.m_cMessageType
}

func (this *Channel) AddPlayer(accountId, playerId int64, playername string, socketId int){
	this.m_playerMap[playerId] = &player{accountId, playerId, playername, socketId}
}

func (this *Channel) RemovePlayer(playerid int64) {
	delete(this.m_playerMap, playerid)
}

func (this *Channel) HasPlayer(playerid int64) bool{
	_, exist := this.m_playerMap[playerid]
	if exist{
		return true
	}
	return false
}

func (this *Channel) SendMessage(msg *ChatMessage){
	for _, v := range this.m_playerMap{
		SendMessage(msg, v)
	}
}