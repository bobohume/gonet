package netgate

import (
	"actor"
	"base"
	"fmt"
	"message"
	"strconv"
	"server/common"
)

type (
	WorldProcess struct {
		actor.Actor
		m_LostTimer *common.SimpleTimer
	}

	IWorldlProcess interface {
		actor.IActor

		RegisterServer(int, int, string, int)
	}
)

func (this * WorldProcess)RegisterServer(ServerType int, Ip string, Port int)  {
	SERVER.GetWorldSocket().SendMsg("COMMON_RegisterRequest",ServerType, Ip, Port)
}

func (this *WorldProcess) Init(num int) {
	this.Actor.Init(num)
	this.m_LostTimer = common.NewSimpleTimer(10)
	this.m_LostTimer.Start()
	this.RegisterTimer(1 * 1000 * 1000 * 1000, this.Update)
	this.RegisterCall("COMMON_RegisterRequest", func() {
		port,_:=strconv.Atoi(UserNetPort)
		this.RegisterServer(int(message.SERVICE_GATESERVER), UserNetIP, port)
	})

	this.RegisterCall("COMMON_RegisterResponse", func() {
			//收到worldserver对自己注册的反馈
			this.m_LostTimer.Stop()
			SERVER.GetLog().Println("收到world对自己注册的反馈")
			SERVER.GetPlayerMgr().SendMsg("Account_Relink")
	})

	this.RegisterCall("G_ClientLost", func(accountId int) {
		SERVER.GetAccountSocket().SendMsg("G_ClientLost", accountId)
	})

	this.RegisterCall("DISCONNECT", func(socketId int) {
		this.m_LostTimer.Start()
	})

	this.RegisterCall("W_A_CreatePlayer", func(accountId int, playername string, sex int32) {
		SERVER.GetAccountSocket().SendMsg("W_A_CreatePlayer", accountId, playername, sex)
	})

	this.RegisterCall("W_A_DeletePlayer", func(accountId int, playerId int) {
		SERVER.GetAccountSocket().SendMsg("W_A_DeletePlayer", accountId, playerId)
	})

	this.Actor.Start()
}

func (this* WorldProcess) Update(){
	if this.m_LostTimer.CheckTimer(){
		SERVER.GetWorldSocket().Start()
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
	accountId := bitstream.ReadInt(base.Bit32)
	socketId := SERVER.GetPlayerMgr().GetAccountSocket(accountId)
	SERVER.GetServer().SendByID(socketId, bitstream.GetBytePtr())
	return false
}


