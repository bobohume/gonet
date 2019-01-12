package netgate

import (
	"actor"
	"base"
	"fmt"
	"message"
	"server/common"
	"strconv"
)

type (
	WorldProcess struct {
		actor.Actor
		m_LostTimer *common.SimpleTimer

		m_SocketId int
	}

	IWorldlProcess interface {
		actor.IActor

		RegisterServer(int, string, int)
		SetSocketId(int)
	}
)

func (this * WorldProcess) RegisterServer(ServerType int, Ip string, Port int)  {
	SERVER.GetDispatchMgr().SendMsg(this.m_SocketId, "COMMON_RegisterRequest",ServerType, Ip, Port)
}

func (this * WorldProcess) SetSocketId(socketId int){
	this.m_SocketId = socketId
}

func (this *WorldProcess) Init(num int) {
	this.Actor.Init(num)
	this.m_LostTimer = common.NewSimpleTimer(10)
	this.m_LostTimer.Start()
	this.m_SocketId = 0
	this.RegisterTimer(1 * 1000 * 1000 * 1000, this.Update)
	this.RegisterCall("COMMON_RegisterRequest", func() {
		port,_:=strconv.Atoi(UserNetPort)
		this.RegisterServer(int(message.SERVICE_GATESERVER), UserNetIP, port)
	})

	this.RegisterCall("COMMON_RegisterResponse", func() {
			//收到worldserver对自己注册的反馈
			this.m_LostTimer.Stop()
			SERVER.GetLog().Println("收到world对自己注册的反馈")
	})

	this.RegisterCall("STOP_ACTOR", func() {
		this.Stop()
	})

	this.RegisterCall("DISCONNECT", func(socketId int) {
		this.m_LostTimer.Start()
	})

	this.Actor.Start()
}

func (this* WorldProcess) Update(){
	if this.m_LostTimer.CheckTimer(){
		SERVER.GetDispatchMgr().GetSocket(this.m_SocketId).Start()
	}
}

func DispatchPacketToClient(id int, buff []byte) bool{
	defer func(){
		if err := recover(); err != nil{
			fmt.Println("WorldClientProcess PacketFunc", err)
		}
	}()

	bitstream := base.NewBitStream(buff, len(buff))
	bitstream.ReadString()//统一格式包头名字
	accountId := bitstream.ReadInt64(base.Bit64)
	socketId := SERVER.GetPlayerMgr().GetSocket(accountId)
	SERVER.GetServer().SendByID(socketId, bitstream.GetBytePtr())
	return false
}


