package netgate

import (
	"gonet/actor"
	"gonet/base"
	"gonet/message"
	"gonet/rpc"
	"strconv"
	"gonet/server/common"
	"strings"
)

var(
	A_C_RegisterResponse = strings.ToLower("A_C_RegisterResponse")
	A_C_LoginResponse 	 = strings.ToLower("A_C_LoginResponse")
)

type (
	AccountProcess struct {
		actor.Actor
		m_LostTimer *common.SimpleTimer

		m_Id uint32
	}

	IAccountProcess interface {
		actor.IActor

		RegisterServer(int, string, int)
		SetSocketId(uint32)
	}
)

func (this * AccountProcess) SetSocketId(socketId uint32){
	this.m_Id = socketId
}

func (this * AccountProcess)RegisterServer(ServerType int, Ip string, Port int)  {
	SERVER.GetAccountCluster().GetCluster(this.m_Id).SendMsg("COMMON_RegisterRequest",ServerType, Ip, Port)
}

func (this *AccountProcess) Init(num int) {
	this.Actor.Init(num)
	this.m_LostTimer = common.NewSimpleTimer(3)
	this.m_LostTimer.Start()
	this.RegisterTimer(1 * 1000 * 1000 * 1000, this.Update)
	this.RegisterCall("COMMON_RegisterRequest", func() {
		port,_:=strconv.Atoi(UserNetPort)
		this.RegisterServer(int(message.SERVICE_GATESERVER), UserNetIP, port)
	})

	this.RegisterCall("COMMON_RegisterResponse", func() {
		this.m_LostTimer.Stop()
	})

	this.RegisterCall("STOP_ACTOR", func() {
		this.Stop()
	})

	this.RegisterCall("DISCONNECT", func(socketId int) {
		this.m_LostTimer.Start()
	})

	this.RegisterCall("A_G_Account_Login", func(accountId int64, socketId int) {
		SERVER.GetPlayerMgr().SendMsg("ADD_ACCOUNT", accountId, socketId)
	})

	this.Actor.Start()
}

func (this* AccountProcess) Update(){
	if this.m_LostTimer.CheckTimer(){
		SERVER.GetAccountCluster().GetCluster(this.m_Id).Start()
	}
}

func DispatchAccountPacketToClient(id int, buff []byte) bool{
	defer func(){
		if err := recover(); err != nil{
			base.TraceCode(err)
		}
	}()

	rpcPacket := rpc.UnmarshalHead(buff)
	if rpcPacket.FuncName == A_C_RegisterResponse || rpcPacket.FuncName == A_C_LoginResponse{
		SERVER.GetServer().SendById(int(rpcPacket.RpcHead.Id), rpcPacket.RpcBody)
	}else{
		socketId := SERVER.GetPlayerMgr().GetSocket(rpcPacket.RpcHead.Id)
		SERVER.GetServer().SendById(socketId, rpcPacket.RpcBody)
	}
	return true
}
