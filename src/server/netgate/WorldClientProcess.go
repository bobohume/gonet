package netgate

import (
	"actor"
	"base"
)

type (
	WorldClientProcess struct {
		actor.Actor
	}
)

func (this *WorldClientProcess) Init(num int) {
	this.Actor.Init(num)
	this.Actor.Start()
}

func (this *WorldClientProcess) PacketFunc(id int, buff []byte) bool{
	bitstream := base.NewBitStream(buff, len(buff))
	bitstream.ReadString()//统一格式包头名字
	accountId := bitstream.ReadInt(base.Bit32)
	socketId := SERVER.GetPlayerMgr().GetAccountSocket(accountId)
	SERVER.GetServer().SendByID(socketId, bitstream.GetBytePtr())
	return false
}