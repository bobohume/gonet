package netgate

/*import (
	"gonet/actor"
	"gonet/base"
	"strings"
)

type (
	UserServerProcess struct {
		actor.Actor
	}
)

func (this *UserServerProcess) Init(num int) {
	this.Actor.Init(num)
	this.RegisterCall("DISCONNECT", func(socketId int) {
		SERVER.GetPlayerMgr().SendMsg("DEL_ACCOUNT", socketId)
	})

	this.Actor.Start()
}

func (this *UserServerProcess) PacketFunc(id int, buff []byte) bool{
	var io actor.CallIO
	io.Buff = buff
	io.SocketId = id

	bitstream := base.NewBitStream(io.Buff, len(io.Buff))
	funcName := bitstream.ReadString()
	funcName = strings.ToLower(funcName)
	if this.FindCall(funcName) != nil{
		this.Send(io)
		return true
	}

	return false
}*/