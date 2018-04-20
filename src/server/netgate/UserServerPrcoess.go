package netgate

import (
	"actor"
	"base"
	"strings"
	"message"
)

type (
	UserServerProcess struct {
		actor.Actor
	}
)

func (this *UserServerProcess) Init(num int) {
	this.Actor.Init(num)
	this.RegisterCall("DISCONNECT", func(caller *actor.Caller, socketId int) {
		SERVER.GetPlayerMgr().SendMsg(caller.SocketId, "DEL_ACCOUNT", caller.SocketId)
	})

	this.Actor.Start()
}

func (this *UserServerProcess) PacketFunc(id int, buff []byte) bool{
	packetId,_ := message.Decode(buff)
	packet := message.GetPakcet(packetId)
	if packet != nil{
		return false
	}

	var io actor.CallIO
	io.Buff = buff
	io.SocketId = id

	bitstream := base.NewBitStream(io.Buff, len(io.Buff))
	funcName := bitstream.ReadString()
	funcName = strings.ToLower(funcName)
	pFunc := this.FindCall(funcName)
	if pFunc != nil{
		this.Send(io)
		return true
	}

	return false
}