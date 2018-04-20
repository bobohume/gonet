package chat

import "sync/atomic"

type(
	Channel struct{
		m_cMessageType int8	//消息类型
		m_nChannelID int	//ID
		m_strChannelName string //名称
		m_playerMap map[int] *player
	}

	IChannel interface {
		GetId() int
		GetMessageType() int8
		HasPlayer(int) bool
		AddPlayer(int, string)
		RemovePlayer(int)
		GetPlayerList() []int
	}
)

type (
	ChannelMgr struct {
		m_hmChannelMap	map[int] *Channel
		m_hmPlayerChannelMap	map[int64] int
		m_nChannelIDSeed	int32
	}

	IChannelMgr interface {
		GetPlayerList(int)  []int
		GetChannelId(int)
		GetChannelIdByType(int, int8)
		RegisterChannel(int8, string) int
		UnregisterChannel(int)

		AddPlayer(int, int, string)
		RemovePlayer(int, int)
		RemoveAllChannel()
		BuildChannelID() int

		getChannel(int) *Channel
		getChannelByType(int, int8) *Channel
	}
)

func (this *Channel) GetId() int{
	return this.m_nChannelID
}

func (this *Channel) GetMessageType() int8{
	return this.m_cMessageType
}

func (this *Channel) AddPlayer(playerid int, playername string){
	this.m_playerMap[playerid] = &player{playerid, playername}
}

func (this *Channel) RemovePlayer(playerid int) {
	delete(this.m_playerMap, playerid)
}

func (this *Channel) HasPlayer(playerid int) bool{
	_, exist := this.m_playerMap[playerid]
	if exist{
		return true
	}
	return false
}

func (this *Channel) GetPlayerList() []int{
	playerList := make([]int, 0)
	for i,_ := range this.m_playerMap{
		playerList = append(playerList, i)
	}
	return playerList
}

func (this *ChannelMgr) getChannel(channelid int)  *Channel{
	pChannel, exist := this.m_hmChannelMap[channelid]
	if exist{
		return pChannel
	}
	return nil
}

func (this *ChannelMgr) getChannelByType(playerid int, cMessageType int8)  *Channel{
	// 对于大规模消息来说，没有意义
	if cMessageType < CHAT_MSG_WORLD{
		return nil
	}

	v := int64(playerid)
	v = (v << 8) | int64(cMessageType)

	channelid, exist := this.m_hmPlayerChannelMap[v]
	if !exist{
		return nil
	}

	return this.getChannel(channelid)
}

func (this *ChannelMgr) GetPlayerList(channelid int)  []int{
	playerList := make([]int, 0)
	pChannel := this.getChannel(channelid)
	if pChannel != nil{
		playerList = pChannel.GetPlayerList()
	}
	return playerList
}

func (this *ChannelMgr) GetChannelId(channelid int) int{
	pChannel := this.getChannel(channelid)
	if pChannel == nil{
		return 0
	}
	return pChannel.GetId()
}

func (this *ChannelMgr) GetChannelIdByType(playerid int, cMessageType int8) int{
	pChannel := this.getChannelByType(playerid, cMessageType)
	if pChannel == nil{
		return 0
	}
	return pChannel.GetId()
}

func (this *ChannelMgr) BuildChannelID() int{
	return  int(atomic.AddInt32(&this.m_nChannelIDSeed, 1))
}

func (this *ChannelMgr) RegisterChannel(cMessageType int8, ChannelName string) int{
	// 大规模消息不能创建频道
	if cMessageType < CHAT_MSG_WORLD {
		return 0
	}

	nId := this.BuildChannelID()
	this.UnregisterChannel(nId)

	pChannel := &Channel{}
	pChannel.m_nChannelID = nId
	pChannel.m_cMessageType = cMessageType
	pChannel.m_strChannelName = ChannelName
	this.m_hmChannelMap[nId] = pChannel
	return nId
}

func (this *ChannelMgr) UnregisterChannel(channelid int) {
	delete(this.m_hmChannelMap, channelid)
}

func (this *ChannelMgr) RemoveAllChannel(channelid int) {
	for i,_ := range this.m_hmPlayerChannelMap{
		delete(this.m_hmPlayerChannelMap, i)
	}

	for i,_ := range this.m_hmChannelMap{
		delete(this.m_hmChannelMap, i)
	}
}

func (this *ChannelMgr) AddPlayer(playerid int, channelid int, playername string) {
	pChannel := this.getChannel(channelid)
	if pChannel == nil{
		return
	}

	pChannel.AddPlayer(playerid, playername)
	v := int64(playerid)
	v = (v << 8) | int64(pChannel.GetMessageType())
	this.m_hmPlayerChannelMap[v] = pChannel.GetId()
}

func (this *ChannelMgr) RemovePlayer(playerid int, channelid int){
	pChannel := this.getChannel(channelid)
	if pChannel == nil{
		return
	}

	pChannel.RemovePlayer(playerid)
	v := int64(playerid)
	v = (v << 8) | int64(pChannel.GetMessageType())
	delete(this.m_hmPlayerChannelMap, v)
}


